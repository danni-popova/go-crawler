package queue

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type Queue interface {
	Insert(url string)
	Remove() (string, error)
}

type queue struct {
	log   *logrus.Logger
	mu    sync.Mutex
	items []string
}

func NewQueue(log *logrus.Logger) Queue {
	return &queue{
		log:   log,
		mu:    sync.Mutex{},
		items: []string{},
	}
}

func (q *queue) Insert(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.contains(url) {
		q.items = append(q.items, url)
	}
}

func (q *queue) Remove() (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) > 0 {
		first, remaining := q.items[0], q.items[1:]
		q.items = remaining
		return first, nil
	}

	return "", errors.New("queue is empty")
}

func (q *queue) contains(url string) bool {
	return slices.Contains(q.items, url)
}
