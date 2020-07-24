package mudl_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/iochen/mudl"
	"github.com/iochen/mudl/utils/download"
	"github.com/iochen/mudl/utils/progress"
)

func TestTask_Download(t *testing.T) {
	logger := log.New(os.Stdout,"",log.Flags())
	task := mudl.Task{
		URL:"https://dl.google.com/go/go1.14.6.src.tar.gz",
		Filename:"go1.14.6.src.tar.gz",
		ThreadNum:3,
		Logger:*logger,
		Config: download.Default(),
		Progress:&progress.Group{},
	}
	fmt.Println(task.Download())
}