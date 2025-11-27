package worker

import (
	"fmt"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/storage"
)

type ArchiveWorker struct {
	storage storage.Storage
	done    chan bool
	ticker  *time.Ticker
}

func NewArchiveWorker(st storage.Storage) *ArchiveWorker {
	return &ArchiveWorker{
		storage: st,
		done:    make(chan bool),
		ticker:  time.NewTicker(5 * time.Minute),
	}
}

func (w *ArchiveWorker) Start() {
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.archiveOldEvents()
			case <-w.done:
				w.ticker.Stop()
				return
			}
		}
	}()
}

func (w *ArchiveWorker) Stop() {
	w.done <- true
}

func (w *ArchiveWorker) archiveOldEvents() {
	cutoffDate := time.Now().AddDate(0, 0, -30)
	count := w.storage.ArchiveOldEvents(cutoffDate)
	fmt.Printf("Archived %d events older than %s\n",
		count, cutoffDate.Format("2006-01-02"))
}
