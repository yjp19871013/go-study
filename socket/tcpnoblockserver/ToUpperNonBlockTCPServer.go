package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"go-study/socket/config"
)

func main() {
	address := config.SERVER_IP + ":" + strconv.Itoa(config.SERVER_PORT)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	readChan := make(chan string)
	writeChan := make(chan string)
	stopChan := make(chan bool)

	go readConn(conn, readChan, stopChan)
	go writeConn(conn, writeChan, stopChan)

	for {
		select {
		case readStr := <-readChan:
			upper := strings.ToUpper(readStr)
			writeChan <- upper
		case stop := <-stopChan:
			if stop {
				break
			}
		}
	}
}

func readConn(conn net.Conn, readChan chan<- string, stopChan chan<- bool) {
	for {
		data := make([]byte, config.SERVER_RECV_LEN)
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		strData := string(data)
		fmt.Println("Received:", strData)

		readChan <- strData
	}

	stopChan <- true
}

func writeConn(conn net.Conn, writeChan <-chan string, stopChan chan<- bool) {
	for {
		strData := <-writeChan
		_, err := conn.Write([]byte(strData))
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", strData)
	}

	stopChan <- true
}
