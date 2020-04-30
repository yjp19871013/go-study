package main

import (
	"fmt"
	"strings"
	"time"
)

func toUpper(done <-chan interface{}, str string) <-chan string {
	strChan := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)

		defer close(strChan)

		for {
			select {
			case <-done:
				return
			case strChan <- strings.ToUpper(str):
			}
		}
	}()

	return strChan
}

func main() {
	done := make(chan interface{})
	defer close(done)

	toUpperChan := toUpper(done, "aaBBcc")
	fmt.Println(<-toUpperChan)
}
