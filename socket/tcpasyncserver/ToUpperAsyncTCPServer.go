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
	for {
		data := make([]byte, config.SERVER_RECV_LEN)
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		strData := string(data)
		fmt.Println("Received:", strData)

		upper := strings.ToUpper(strData)
		_, err = conn.Write([]byte(upper))
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", upper)
	}
}
