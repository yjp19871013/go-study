package read

import "fmt"

// 阻塞
func Nil() {
	var ch chan interface{}
	<-ch
}

// 阻塞
func Empty() {
	ch := make(chan interface{})
	<-ch
}

// 阻塞
func CloseNotEmpty() {
	ch := make(chan interface{})
	ch <- 10
	close(ch)
	fmt.Println(<-ch)
}

// 读取0值
func CloseEmpty() {
	ch := make(chan interface{})
	close(ch)
	fmt.Println(<-ch)
}
