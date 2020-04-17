package main

import (
	"tcp_client_class/tcp_client"
	"fmt"
)

func main() {
	client := tcp_client.CreateTCPClient("10.3.1.132", 10050, 10)

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
