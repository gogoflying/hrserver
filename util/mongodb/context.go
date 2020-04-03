package mongodb

import (
	"container/heap"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

// session
type session struct {
	*mgo.Session
	ref   int
	index int
}

// session heap
type sessionHeap []*session

func (h sessionHeap) Len() int {
	return len(h)
}

func (h sessionHeap) Less(i, j int) bool {
	return h[i].ref < h[j].ref
}

func (h sessionHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *sessionHeap) Push(s interface{}) {
	s.(*session).index = len(*h)
	*h = append(*h, s.(*session))
}

func (h *sessionHeap) Pop() interface{} {
	l := len(*h)
	s := (*h)[l-1]
	s.index = -1
	*h = (*h)[:l-1]
	return s
}

type dialContext struct {
	sync.Mutex
	sessions sessionHeap
}

// goroutine safe
func dial(url string, sessionNum int) (*dialContext, error) {
	c, err := dialWithTimeout(url, sessionNum, 10*time.Second, 0)
	return c, err
}

// goroutine safe
func dialWithTimeout(url string, sessionNum int, dialTimeout time.Duration, timeout time.Duration) (*dialContext, error) {
	if sessionNum <= 0 {
		sessionNum = 100
		//log.Release("invalid sessionNum, reset to %v", sessionNum)
	}

	s, err := mgo.DialWithTimeout(url, dialTimeout)
	if err != nil {
		return nil, err
	}
	s.SetSyncTimeout(timeout)
	s.SetSocketTimeout(timeout)

	c := new(dialContext)

	// sessions
	c.sessions = make(sessionHeap, sessionNum)
	c.sessions[0] = &session{s, 0, 0}
	for i := 1; i < sessionNum; i++ {
		c.sessions[i] = &session{s.New(), 0, i}
	}
	heap.Init(&c.sessions)

	return c, nil
}

// goroutine safe
func (c *dialContext) Close() {
	c.Lock()
	for _, s := range c.sessions {
		s.Close()
		if s.ref != 0 {
			//			log.Error("session ref = %v", s.ref)
		}
	}
	c.Unlock()
}

// goroutine safe
func (c *dialContext) Ref() *session {
	c.Lock()
	s := c.sessions[0]
	if s.ref == 0 {
		s.Refresh()
	}
	s.ref++
	heap.Fix(&c.sessions, 0)
	c.Unlock()

	return s
}

// goroutine safe
func (c *dialContext) UnRef(s *session) {
	c.Lock()
	s.ref--
	heap.Fix(&c.sessions, s.index)
	c.Unlock()
}
