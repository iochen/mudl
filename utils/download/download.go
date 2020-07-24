package download

import (
	"net/http"
	"net/url"
	"time"
)

type Task struct {
	Client  *http.Client
	Request *http.Request
}

func Default() *Task {
	req, _ := http.NewRequest("GET", "", nil)
	return &Task{
		Request: req,
		Client:  &http.Client{},
	}
}

func (t *Task)URL(dst string) error {
	u, err := url.Parse(dst)
	if err != nil {
		return err
	}
	t.Request.URL = u
	return nil
}

func (t *Task)Method(m string) *Task {
	t.Request.Method = m
	return t
}

func (t *Task)UserAgent(ua string) *Task {
	t.Request.Header.Set("User-Agent",ua)
	return t
}

func (t *Task)Timeout(d time.Duration) *Task {
	t.Client.Timeout = d
	return t
}

func (t *Task)AddHeader(key,value string) *Task {
	t.Request.Header.Add(key,value)
	return t
}

func (t *Task)RemoveHeader(key string) *Task {
	t.Request.Header.Del(key)
	return t
}

func (t *Task)AddCookie(name,value string) *Task {
	t.Request.AddCookie(&http.Cookie{
		Name:name,
		Value:value,
	})
	return t
}

func (t *Task)Download() (*http.Response,error) {
	if t.Request.Method == "" {
		t.Request.Method= "GET"
	}
	return t.Client.Do(t.Request)
}