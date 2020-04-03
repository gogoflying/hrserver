package msgq

import (
	"fmt"
	"time"

	"sync"
	"testing"
)

type TestData struct {
	Data string
}

func NewTestData(s string) *TestData {
	return &TestData{
		Data: s,
	}
}

func (t *TestData) HandlerMsg() error {
	fmt.Println(t.Data)
	time.Sleep(time.Second)
	return nil
}

func Test_List(t *testing.T) {

	d := NewTestData("123")

	msg := NewMsgQ(10)

	msg.list_Put(d)

	e := msg.list_Get()

	t.Log(e.Value.(*TestData).Data)

}

func Test_Bench(t *testing.T) {

	var wg sync.WaitGroup

	msg := NewMsgQ(300)

	go func() {

		count := 3000

		for {
			count--

			if count < 0 {
				break
			}
			s := fmt.Sprint(count)
			if err := msg.PutMsg(NewTestData(s)); err != nil {
				t.Log(err)
			}
		}

	}()

	wg.Add(1)
	go msg.Do()
	wg.Wait()

	return
}
