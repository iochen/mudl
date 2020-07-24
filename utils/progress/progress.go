package progress

import (
	"errors"
	"sync"

	"github.com/iochen/mudl/utils/size"
)

// Progress is a thread-safe structure to record the progress
type Progress struct {
	sync.RWMutex
	already      size.Size
	total        size.Size
}

var (
	ErrSizeNotCorrect = errors.New("already written size more than total size")
)

func (pg *Progress) checkSize() error {
	pg.RLock()
	defer pg.RUnlock()
	if pg.already > pg.total && pg.already >= 0 {
		return ErrSizeNotCorrect
	}
	return nil
}

func (pg *Progress) SetTotal(s size.Size) error {
	if s <0 {
		s = -1
	}
	pg.Lock()
	defer pg.Unlock()
	pg.total = s
	return pg.checkSize()
}

func (pg *Progress) SetAlready(s size.Size) error {
	pg.Lock()
	defer pg.Unlock()
	pg.already = s
	return pg.checkSize()
}

func (pg *Progress) Write(b []byte) (int, error) {
	pg.Lock()
	defer pg.Unlock()
	l := len(b)
	pg.already += size.Size(l)
	return l, pg.checkSize()
}

func (pg *Progress) Status() (already, total size.Size) {
	pg.RLock()
	defer pg.RUnlock()
	already, total = pg.already, pg.total
	return
}
