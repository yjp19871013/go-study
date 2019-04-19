package msg_queue

import (
	"encoding/json"
	"fmt"
	"sync"
)

type MsgQueueHandler interface {
	OnStart(q *MsgQueue)
	OnStop(q *MsgQueue)
	OnMsgRecv(q *MsgQueue, msg *Message)
	OnDefaultRun(q *MsgQueue)
}

type Message struct {
	Msg      string
	JsonData []byte
}

func formMessage(msg string, data interface{}) (*Message, error) {
	jsonStr := make([]byte, 0)
	if data != nil {
		dataJson, err := json.Marshal(data)
		if err != nil {
			fmt.Println("formMessage error", err.Error())
			return nil, err
		}

		jsonStr = dataJson
	}

	message := new(Message)
	message.Msg = msg
	message.JsonData = jsonStr

	return message, nil
}

type MsgQueue struct {
	handler MsgQueueHandler

	isStart bool

	observers      []*MsgQueue
	observersMutex sync.Mutex

	msgChan  chan *Message
	stopChan chan bool
}

func InitMsgQueue(handler MsgQueueHandler) *MsgQueue {
	return InitMsgQueueWithMsgBufferSize(handler, 0)
}

func InitMsgQueueWithMsgBufferSize(handler MsgQueueHandler, bufferSize int) *MsgQueue {
	if handler == nil {
		return nil
	}

	q := new(MsgQueue)
	q.handler = handler
	q.observers = make([]*MsgQueue, 0)

	if bufferSize > 0 {
		q.msgChan = make(chan *Message, bufferSize)
	} else {
		q.msgChan = make(chan *Message)
	}

	q.stopChan = make(chan bool)

	return q
}

func DestroyMsgQueue(q *MsgQueue) {
	if q == nil {
		return
	}

	q.stopChan = nil
	q.msgChan = nil
	q.observers = nil
	q.handler = nil
	q = nil
}

func (q *MsgQueue) Start() {
	q.isStart = true

	go q.run()
}

func (q *MsgQueue) Stop() {
	if !q.isStart {
		return
	}

	q.stopChan <- true
	<-q.stopChan
}

func (q *MsgQueue) AddObserver(observer *MsgQueue) {
	if observer == nil {
		return
	}

	q.observersMutex.Lock()
	defer q.observersMutex.Unlock()

	q.observers = append(q.observers, observer)
}

func (q *MsgQueue) DeleteObserver(observer *MsgQueue) {
	if observer == nil {
		return
	}

	q.observersMutex.Lock()
	defer q.observersMutex.Unlock()

	for i, o := range q.observers {
		if o == observer {
			q.observers = append(q.observers[:i], q.observers[i+1:]...)
			break
		}
	}
}

func (q *MsgQueue) NotifyObservers(msg string, data interface{}) {
	q.observersMutex.Lock()
	defer q.observersMutex.Unlock()

	for _, o := range q.observers {
		o.SendMsg(msg, data)
	}
}

func (q *MsgQueue) SendMsg(msg string, data interface{}) {
	message, err := formMessage(msg, data)
	if err != nil {
		fmt.Println("SendMsg error:", err)
		return
	}
	fmt.Println("2222222", msg)
	q.msgChan <- message
}

func (q *MsgQueue) run() {
	q.handler.OnStart(q)

	for {
		select {
		case stop := <-q.stopChan:
			if stop {
				q.handler.OnStop(q)
				q.stopChan <- true
				return
			}
			q.stopChan <- true
		case msg := <-q.msgChan:
			q.handler.OnMsgRecv(q, msg)
		default:
			q.handler.OnDefaultRun(q)
		}
	}
}
