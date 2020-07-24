package mudl

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/iochen/mudl/utils/download"
	"github.com/iochen/mudl/utils/progress"
)

type thread struct {
	ID       int
	Task     *Task
	RangeL   int64
	RangeR   int64
	Progress *progress.Progress
}

type Task struct {
	sync.WaitGroup
	URL       string
	Filename  string
	Config    *download.Task
	ThreadNum int
	threads   []*thread
	length    int64
	Progress  *progress.Group
	log.Logger
}

func (thread *thread) Start() {
	defer thread.Task.Done()
	thread.Task.Logger.Printf("[thread %d]starting...\n",thread.ID)
	fn := thread.Task.Filename + ".mudlid." + strconv.Itoa(thread.ID)
	file, err := os.Create(fn)
	if err != nil {
		thread.Task.Logger.Fatalf("[thread %d]cannot create file: %s\n",thread.ID,fn)
	}
	defer file.Close()
	task := download.Default()
	task.Request = thread.Task.Config.Request.Clone(context.Background())
	task.AddHeader("Range",
		"bytes="+strconv.Itoa(int(thread.RangeL))+"-"+strconv.Itoa(int(thread.RangeR)))
	fmt.Println("bytes="+strconv.Itoa(int(thread.RangeL))+"-"+strconv.Itoa(int(thread.RangeR)))
	thread.Task.Logger.Printf("[thread %d]downloading...\n",thread.ID)
	resp, err := task.Download()
	if err != nil {
		thread.Task.Logger.Printf("[thread %d]error: %s, wait to retry...\n",thread.ID,err.Error())
		time.Sleep(1*time.Second)
		thread.Task.Logger.Printf("[thread %d]retrying...\n",thread.ID)
		thread.Start()
	}
	defer resp.Body.Close()
	buf := bufio.NewWriter(file)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		thread.Task.Logger.Fatalln(err)
	}
	defer buf.Flush()
	thread.Task.Logger.Printf("[thread %d]finished!\n",thread.ID)
}

func div(l int64,t int) []int64 {
	div := make([]int64,t)
	part := l / int64(t)
	for i:=0;i<t-1;i++ {
		div[i] = part
	}
	div[t-1] = l - part * (int64(t)-1)
	return div
}

func (t *Task) Download() error {
	if t.ThreadNum < 1 {
		return errors.New("wrong thread number")
	}
	err := t.Config.URL(t.URL)
	if err != nil {
		return err
	}
	resp, err := t.Config.Download()
	if err != nil {
		return err
	}
	if resp.Header.Get("Accept-Ranges") == "" {
		t.ThreadNum = 1
	}
	sL := resp.Header.Get("Content-Length")
	if sL == "" {
		t.ThreadNum = 1
	} else {
		l, err := strconv.Atoi(sL)
		if err != nil {
			return err
		}
		t.length = int64(l)
	}
	resp.Body.Close()

	div := div(t.length,t.ThreadNum)
	var last int64
	t.threads = make([]*thread,t.ThreadNum)
	for i := 0; i < t.ThreadNum; i++ {
		t.Add(1)
		p := &progress.Progress{}
		t.Progress.Append(p)
		t.threads[i] = &thread{
			ID:i,
			Task:t,
			RangeL:last,
			RangeR:last + div[i]-1,
			Progress:p,
		}
		go t.threads[i].Start()
		last = last + div[i]
	}
	t.Wait()

	file, err := os.Create(t.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < t.ThreadNum; i++ {
		fn := t.Filename + ".mudlid." + strconv.Itoa(i)
		f, err := os.OpenFile(fn, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		_, err = io.Copy(file, f)
		if err != nil {
			return err
		}
		f.Close()
		err = os.Remove(fn)
		if err != nil {
			return err
		}
	}
	return nil
}
