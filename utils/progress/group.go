package progress

import (
	"sync"

	"github.com/iochen/mudl/utils/size"
)

// Group is a thread-safe structure to record the progress of a group of progresses
type Group struct {
	sync.RWMutex
	progressList []*Progress
}

func (g *Group)Append(pg *Progress) {
	g.Lock()
	defer g.Unlock()
	g.progressList = append(g.progressList,pg)
}

func (g *Group)Status() (already, total size.Size, unknown bool) {
	g.RLock()
	defer g.RUnlock()
	for i:=0;i<len(g.progressList);i++ {
		already += g.progressList[i].already
		total += g.progressList[i].total
	}
	return
}
