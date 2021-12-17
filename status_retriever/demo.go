package main

import (
	"fmt"
	"time"
)

func main() {
	retriever := newStatusRetriever()
	defer destroyStatusRetriever(retriever)

	retriever.AddRetrieverTask(&StatusRetrieverTask{
		BaseUrl:      "http://example.com",
		Token:        "token",
		ResourceID:   "xxx",
		WantedStatus: "success",
		WaitErrorCallback: func(task *StatusRetrieverTask, err error) {
			fmt.Println("Receiver error:", err)
		},
		GetStatusCallback: func(task *StatusRetrieverTask, statusInfo *StatusInfo) {
			fmt.Println("Receiver wanted status:", statusInfo.Status)
		},
	})

	time.Sleep(5 * time.Second)
}
