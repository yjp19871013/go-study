package msg_queue

import (
	"encoding/json"
	"fmt"
)

type MsgQueueHandler interface {
	OnStart()
	OnStop()
	OnMsgRecv(msg *Message)
}

type Message struct {
	Msg      string
	JsonData []byte
}

func formMessage(msg string, data interface{}) (*Message, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println("SendMsg error", err.Error())
		return nil, err
	}

	message := new(Message)
	message.Msg = msg
	message.JsonData = jsonStr

	return message, nil
}

type MsgQueue struct {
	handler MsgQueueHandler

	msgChan  chan *Message
	stopChan chan bool
}

func InitMsgQueue(handler MsgQueueHandler) *MsgQueue {
	if handler == nil {
		return nil
	}

	q := new(MsgQueue)
	q.handler = handler
	q.msgChan = make(chan *Message)
	q.stopChan = make(chan bool)

	return q
}

func DestroyMsgQueue(q *MsgQueue) {
	if q == nil {
		return
	}

	q.stopChan = nil
	q.msgChan = nil
	q.handler = nil
	q = nil
}

func (q *MsgQueue) Start() {
	go q.run()
}

func (q *MsgQueue) Stop() {
	q.stopChan <- true
}

func (q *MsgQueue) SendMsg(msg string, data interface{}) {
	message, err := formMessage(msg, data)
	if err != nil {
		return
	}

	q.msgChan <- message
}

func (q *MsgQueue) run() {
	q.handler.OnStart()

	for {
		select {
		case stop := <-q.stopChan:
			if stop {
				q.handler.OnStop()
			}
		case msg := <-q.msgChan:
			q.dealWithMsg(msg)
		}
	}
}

func (q *MsgQueue) dealWithMsg(msg *Message) {
	if msg == nil {
		return
	}

	q.handler.OnMsgRecv(msg)
}
