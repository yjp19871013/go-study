package main

import (
	"com.fs/tcp_client"
	"fmt"
)

func main() {
	client := tcp_client.CreateTCPClient("localhost", 10050, 10)

	listenChan := make(chan int)
	client.StartConnect(listenChan)

	for {
		state := <-listenChan
		if state == tcp_client.TcpStateConnected {
			for i := 0; i < 10; i++ {
				data := []byte("abcdef")
				err := client.Write(data)
				if err != nil {
					fmt.Println(err)
					return
				}

				recvData := make([]byte, 1024)
				_, err = client.Read(recvData)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(recvData))
			}

			break
		}
	}

	client.StopConnect()

	tcp_client.DestroyTCPClient(client)
}
