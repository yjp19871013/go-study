package msg_queue

import (
	"encoding/json"
	"fmt"
	"sync"
)

type MsgQueueHandler interface {
	// 在消息线程进入消息循环前执行
	OnStart(q *MsgQueue)

	// 在消息线程退出消息循环前执行
	OnStop(q *MsgQueue)

	// 消息响应函数
	OnMsgRecv(q *MsgQueue, msg *Message)

	// 消息循环中的default调用
	OnDefaultRun(q *MsgQueue)
}

type Message struct {
	// 消息名称
	Msg string

	// 消息数据
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

// 初始化无buffer的消息队列
func InitMsgQueue(handler MsgQueueHandler) *MsgQueue {
	return InitMsgQueueWithMsgBufferSize(handler, 0)
}

// 创建有buffer的消息队列，一般能够提供更好的并发
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

// 销毁消息队列
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

// 启动消息循环
func (q *MsgQueue) Start(useDefaultRun bool) {
	q.isStart = true

	if useDefaultRun {
		go q.runWithDefault()
	} else {
		go q.runWithoutDefault()
	}
}

// 停止消息循环
func (q *MsgQueue) Stop() {
	if !q.isStart {
		return
	}

	q.stopChan <- true
	<-q.stopChan

	q.isStart = false
}

// 添加观察者
func (q *MsgQueue) AddObserver(observer *MsgQueue) {
	if observer == nil {
		return
	}

	q.observersMutex.Lock()
	defer q.observersMutex.Unlock()

	q.observers = append(q.observers, observer)
}

// 删除观察者
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

// 通知观察者
func (q *MsgQueue) NotifyObservers(msg string, data interface{}) {
	q.observersMutex.Lock()
	defer q.observersMutex.Unlock()

	for _, o := range q.observers {
		o.SendMsg(msg, data)
	}
}

// 发送消息
func (q *MsgQueue) SendMsg(msg string, data interface{}) {
	message, err := formMessage(msg, data)
	if err != nil {
		fmt.Println("SendMsg error:", err)
		return
	}

	q.msgChan <- message
}

// 使用default的消息循环
func (q *MsgQueue) runWithDefault() {
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

// 不使用default的消息循环
func (q *MsgQueue) runWithoutDefault() {
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
		}
	}
}
