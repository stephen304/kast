package media

import (
	"fmt"
	"log"
	"sync"
)

type queue struct {
	prevItems []string
	playHead  string
	nextItems []string
	m         sync.Mutex
}

func newQueue() *queue {
	return &queue{}
}

func (q *queue) enqueue(media string) {
	q.m.Lock()
	defer q.m.Unlock()

	if len(q.playHead) == 0 {
		q.playHead = media
	} else {
		q.nextItems = append(q.nextItems, media)
	}
	fmt.Println("Enqueue")
	fmt.Printf("%+v\n", q.prevItems)
	fmt.Printf("%+v\n", q.playHead)
	fmt.Printf("%+v\n", q.nextItems)
}

func (q *queue) Get() string {
	return q.playHead
}

func (q *queue) HasNext() bool {
	return len(q.nextItems) > 0
}

func (q *queue) Prev() string {
	q.m.Lock()
	defer q.m.Unlock()

	if len(q.prevItems) > 0 {
		// Seek prev
		if len(q.playHead) > 0 {
			q.nextItems = append([]string{q.playHead}, q.nextItems...)
		}

		q.playHead = q.prevItems[len(q.prevItems)-1]
		q.prevItems = q.prevItems[:len(q.prevItems)-1]

	} else {
		log.Print("Nothing to reqind to")
	}
	fmt.Println("Prev")
	fmt.Printf("%+v\n", q.prevItems)
	fmt.Printf("%+v\n", q.playHead)
	fmt.Printf("%+v\n", q.nextItems)
	return q.playHead
}

func (q *queue) Next() string {
	q.m.Lock()
	defer q.m.Unlock()

	if len(q.playHead) > 0 {
		q.prevItems = append(q.prevItems, q.playHead)
	}

	q.playHead = ""

	if len(q.nextItems) > 0 {
		q.playHead = q.nextItems[0]
		if len(q.nextItems) > 1 {
			q.nextItems = q.nextItems[1:]
		} else {
			q.nextItems = nil
		}
	}
	fmt.Println("Next")
	fmt.Printf("%+v\n", q.prevItems)
	fmt.Printf("%+v\n", q.playHead)
	fmt.Printf("%+v\n", q.nextItems)
	return q.playHead
}

func (q *queue) Empty() {
	q.m.Lock()
	defer q.m.Unlock()

	q.prevItems = []string{}
	q.playHead = ""
	q.nextItems = []string{}
}
