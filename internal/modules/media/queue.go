package media

import (
	"sync"
)

type queue struct {
	prevItems [][2]string
	playHead  [2]string
	nextItems [][2]string
	m         sync.Mutex
}

func newQueue() *queue {
	return &queue{}
}

func (q *queue) enqueue(url string, title string) {
	if url == "" {
		return
	}

	q.m.Lock()
	defer q.m.Unlock()

	if len(q.playHead[0]) == 0 {
		q.playHead[0] = url
		q.playHead[1] = title
	} else {
		q.nextItems = append(q.nextItems, [2]string{url, title})
	}
}

func (q *queue) Get() (string, string) {
	return q.playHead[0], q.playHead[1]
}

func (q *queue) GetUrl() string {
	return q.playHead[0]
}

func (q *queue) HasNext() bool {
	return len(q.nextItems) > 0
}

func (q *queue) Prev() (string, string) {
	q.m.Lock()
	defer q.m.Unlock()

	if len(q.prevItems) > 0 {
		// Seek prev
		if len(q.playHead[0]) > 0 {
			q.nextItems = append([][2]string{q.playHead}, q.nextItems...)
		}

		q.playHead = q.prevItems[len(q.prevItems)-1]
		q.prevItems = q.prevItems[:len(q.prevItems)-1]
	}
	return q.playHead[0], q.playHead[1]
}

func (q *queue) Next() (string, string) {
	q.m.Lock()
	defer q.m.Unlock()

	if len(q.playHead[0]) > 0 {
		q.prevItems = append(q.prevItems, q.playHead)
	}

	q.playHead = [2]string{}

	if len(q.nextItems) > 0 {
		q.playHead = q.nextItems[0]
		if len(q.nextItems) > 1 {
			q.nextItems = q.nextItems[1:]
		} else {
			q.nextItems = nil
		}
	}
	return q.playHead[0], q.playHead[1]
}

func (q *queue) Empty() {
	q.m.Lock()
	defer q.m.Unlock()

	q.prevItems = [][2]string{}
	q.playHead = [2]string{}
	q.nextItems = [][2]string{}
}
