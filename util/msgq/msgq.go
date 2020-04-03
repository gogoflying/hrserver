package msgq

import (
	"container/list"
	"errors"
	"sync"
	"time"
	"util/log"
)

const (
	Scan_List_Intervel = 5 //s,遍历list间隔
)

type IMsgQ interface {
	HandlerMsg() error
}

var (
	Err_CHANFULL = errors.New("chan has full")
)

type MsgQ struct {
	msqList *list.List
	lock    sync.RWMutex
	chanMsg chan IMsgQ
}

func NewMsgQ(mLen uint32) *MsgQ {
	msgQ := new(MsgQ)
	msgQ.msqList = list.New()
	msgQ.chanMsg = make(chan IMsgQ, mLen)
	return msgQ
}

func (m *MsgQ) PutMsg(msg IMsgQ) (err error) {
	select {
	case m.chanMsg <- msg:
	default:
		m.list_Put(msg) //如果chan满了，暂时放置到list里面
		err = Err_CHANFULL
	}
	return
}

func (m *MsgQ) Do() {
	defer func() {
		if r := recover(); r != nil {
			log.GetLog().LogError(r)
		}
	}()

	go m.list_Scan()

	for {
		select {
		case v := <-m.chanMsg:
			if err := v.HandlerMsg(); err != nil {
				continue
			}
		}
	}
}

func (m *MsgQ) OnClose() {
	close(m.chanMsg)
}

func (m *MsgQ) put_Chan(msgQ IMsgQ) (err error) {
	select {
	case m.chanMsg <- msgQ:
	default:
		err = Err_CHANFULL
	}
	return
}
func (m *MsgQ) list_Scan() {

	defer func() {
		if r := recover(); r != nil {
			log.GetLog().LogError(r)
		}
	}()

	for {
		m.lock.Lock()
		for e := m.msqList.Front(); e != nil; {
			if err := m.put_Chan(e.Value.(IMsgQ)); err != nil {
				break
			}
			en := e.Next()
			m.msqList.Remove(e)
			e = en
		}
		m.lock.Unlock()
		time.Sleep(time.Second * Scan_List_Intervel)
	}

}

func (m *MsgQ) list_Put(msg IMsgQ) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.msqList.PushBack(msg)
}

func (m *MsgQ) list_Get() (e *list.Element) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.msqList.Front()
}

func (m *MsgQ) list_Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.msqList.Len()
}

func (m *MsgQ) list_Remove(e *list.Element) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.msqList.Remove(e)
	return
}
