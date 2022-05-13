package list

import (
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

type SyncedList interface {
	Add(url string)
	Contains(url string) bool
	PrintContents()
}

type list struct {
	list map[string]bool
	mu   sync.Mutex
	log  *logrus.Logger
}

func NewSyncedList(log *logrus.Logger) SyncedList {
	return &list{
		list: make(map[string]bool),
		mu:   sync.Mutex{},
		log:  log,
	}
}

// Add takes a strings and saves it to the list unless it's already in there
func (l *list) Add(url string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.list[url] {
		l.list[url] = true
	}
}

// Contains checks if a string is contained in the synced list
func (l *list) Contains(url string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.list[url]
}

// PrintContents logs every string in the list
func (l *list) PrintContents() {
	l.log.Infof("Visited %d URLs", len(l.list))
	l.log.Info("Printing the contents of the crawled list:")
	for _, k := range maps.Keys(l.list) {
		l.log.Info(k)
	}
}
