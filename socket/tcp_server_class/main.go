package main

import (
	"tcp_server_class/tcp_server"
	"fmt"
	"strings"
	"syscall"

	DEATH "gopkg.in/vrecan/death.v3"
)

func main() {
	server := tcp_server.CreateTCPServer("", 10050, 1024)

	err := server.Start(onAccept)
	if err != nil {
		fmt.Println(err)
		return
	}

	death := DEATH.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	_ = death.WaitForDeath()

	server.Stop()

	tcp_server.DestroyServer(server)
}

func onAccept(_ *tcp_server.ClientConn, readChan chan []byte, writeChan chan []byte) {
	for {
		data := <-readChan
		fmt.Println("recv:", string(data))
		writeChan <- []byte(strings.ToUpper(string(data)))
	}
}
