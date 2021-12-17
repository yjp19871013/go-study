package main

import (
	"math/rand"
	"sync"
	"time"
)

const (
	statusRetrieverPeriodSec = 1
)

type StatusInfo struct {
	Status string
}

type WaitErrorCallback func(task *StatusRetrieverTask, err error)
type GetStatusCallback func(task *StatusRetrieverTask, statusInfo *StatusInfo)

type StatusRetrieverTask struct {
	BaseUrl      string
	Token        string
	ResourceID   string
	WantedStatus string
	WaitErrorCallback
	GetStatusCallback
}

type statusRetriever struct {
	done             chan interface{}
	taskChan         chan *StatusRetrieverTask
	taskSelfOverChan chan *StatusRetrieverTask
	taskDoneMap      map[*StatusRetrieverTask]chan interface{}
	taskDoneMapMutex sync.Mutex
}

func newStatusRetriever() *statusRetriever {
	retriever := new(statusRetriever)
	retriever.done = make(chan interface{})
	retriever.taskChan = make(chan *StatusRetrieverTask)
	retriever.taskSelfOverChan = make(chan *StatusRetrieverTask)
	retriever.taskDoneMap = make(map[*StatusRetrieverTask]chan interface{})

	go retriever.scheduleTask()

	return retriever
}

func destroyStatusRetriever(retriever *statusRetriever) {
	if retriever == nil {
		return
	}

	retriever.taskDoneMapMutex.Lock()
	for _, taskDone := range retriever.taskDoneMap {
		retriever.taskDoneMapMutex.Unlock()
		taskDone <- true
		close(taskDone)
		retriever.taskDoneMapMutex.Lock()
	}
	retriever.taskDoneMap = nil
	retriever.taskDoneMapMutex.Unlock()

	retriever.done <- true
	close(retriever.done)
	retriever.done = nil

	close(retriever.taskChan)
	retriever.taskChan = nil

	retriever = nil
}

func (retriever *statusRetriever) AddRetrieverTask(task *StatusRetrieverTask) {
	retriever.taskDoneMapMutex.Lock()
	retriever.taskDoneMap[task] = make(chan interface{})
	retriever.taskDoneMapMutex.Unlock()

	retriever.taskChan <- task
}

func (retriever *statusRetriever) RemoveRetrieverTask(task *StatusRetrieverTask) {
	retriever.taskDoneMapMutex.Lock()
	done, ok := retriever.taskDoneMap[task]
	if !ok {
		retriever.taskDoneMapMutex.Unlock()
		return
	}

	delete(retriever.taskDoneMap, task)
	retriever.taskDoneMapMutex.Unlock()

	done <- true
	close(done)
	done = nil
}

func (retriever *statusRetriever) scheduleTask() {
	for {
		select {
		case <-retriever.done:
			return
		case task := <-retriever.taskSelfOverChan:
			retriever.RemoveRetrieverTask(task)
		case task := <-retriever.taskChan:
			go retriever.dealTask(task)
		}
	}
}

func (retriever *statusRetriever) dealTask(task *StatusRetrieverTask) {
	retriever.taskDoneMapMutex.Lock()
	done, ok := retriever.taskDoneMap[task]
	if !ok {
		retriever.taskDoneMapMutex.Unlock()
		return
	}
	retriever.taskDoneMapMutex.Unlock()

	defer func() {
		retriever.taskDoneMapMutex.Lock()
		_, ok := retriever.taskDoneMap[task]
		if ok {
			delete(retriever.taskDoneMap, task)
		}
		retriever.taskDoneMapMutex.Unlock()
	}()

	for {
		select {
		case <-done:
			return
		default:
			statusInfo, err := GetStatus(task.BaseUrl, task.Token, task.ResourceID)
			if err != nil {
				if task.WaitErrorCallback != nil {
					task.WaitErrorCallback(task, err)
				}

				time.Sleep(statusRetrieverPeriodSec * time.Second)
				continue
			}

			if task.WantedStatus == statusInfo.Status {
				if task.GetStatusCallback != nil {
					task.GetStatusCallback(task, statusInfo)
				}

				return
			}

			time.Sleep(statusRetrieverPeriodSec * time.Second)
		}
	}
}

func GetStatus(baseUrl string, token string, resourceID string) (*StatusInfo, error) {
	rand.Seed(time.Now().Unix())
	num := rand.Intn(2)

	var status string
	if num == 0 {
		status = "success"
	} else {
		status = "failure"
	}

	return &StatusInfo{Status: status}, nil
}
