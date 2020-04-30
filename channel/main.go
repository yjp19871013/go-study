package main

import (
	"fmt"
	"strings"
	"time"
)

func str(done <-chan interface{}, str string) <-chan string {
	strChan := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)

		defer close(strChan)

		select {
		case <-done:
			return
		case strChan <- str:
		}
	}()

	return strChan
}

func toUpper(done <-chan interface{}, str <-chan string) <-chan string {
	strChan := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)

		defer close(strChan)

		for {
			select {
			case <-done:
				return
			case strChan <- strings.ToUpper(<-str):
			}
		}
	}()

	return strChan
}

func appendStr(done <-chan interface{}, oldStr <-chan string, appendStr string) <-chan string {
	strChan := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)

		defer close(strChan)

		for {
			select {
			case <-done:
				return
			case strChan <- <-oldStr + appendStr:
			}
		}
	}()

	return strChan
}

func main() {
	done := make(chan interface{})

	strChan := str(done, "aaa111bbb222CCC")
	fmt.Println(1)

	upperChan := toUpper(done, strChan)
	fmt.Println(2)

	appendStrChan := appendStr(done, upperChan, "dddd")
	fmt.Println(3)

	fmt.Println(<-appendStrChan)

	close(done)

	done = make(chan interface{})
	fmt.Println(<-appendStr(done, toUpper(done, str(done, "aaa111bbb222CCC")), "dddd"))
	close(done)
}
