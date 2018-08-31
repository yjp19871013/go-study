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

		lineLen := len(line)
		for written := 0; written < lineLen; {
			var toWrite string
			if lineLen-written > config.SERVER_RECV_LEN {
				toWrite = line[written : written+config.SERVER_RECV_LEN]
			} else {
				toWrite = line[written:]
			}

			n, err := conn.Write([]byte(toWrite))
			checkError(err)

			fmt.Println("Write:", toWrite)

			msg := make([]byte, config.SERVER_RECV_LEN)
			n, err = conn.Read(msg)
			checkError(err)

			fmt.Println("Response:", string(msg))

			written += n
		}
	}

}
