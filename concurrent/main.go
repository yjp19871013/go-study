package main

import (
	"fmt"
	"time"
)

func Foo(i int, ch chan int) {
	fmt.Printf("%d will sleep\n", i)
	time.Sleep(5 * time.Second)
	fmt.Printf("%d wake up\n", i)
	ch <- 1
}

func main() {
	ch := make(chan int)

	for i := 0; i < 5; i++ {
		go Foo(i, ch)
	}

	count := 0
	for count < 5 {
		count += <-ch
	}
}
