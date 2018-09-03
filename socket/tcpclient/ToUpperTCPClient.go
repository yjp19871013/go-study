package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"go-study/socket/config"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	address := config.SERVER_IP + ":" + strconv.Itoa(config.SERVER_PORT)

	conn, err := net.Dial("tcp", address)
	checkError(err)

	defer conn.Close()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		_, err = conn.Write([]byte(line))
		checkError(err)

		fmt.Println("Write:", line)

		recvLen := len(line)
		msg := make([]byte, recvLen)
		received := 0
		for recvLen >= 0 {
			start := received
			received, err = conn.Read(msg[start:])
			checkError(err)
			recvLen -= received
		}

		fmt.Println("Response:", string(msg))
	}

}
