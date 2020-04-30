package write

// 阻塞
func Nil() {
	var ch chan interface{}
	ch <- 10
}

// 阻塞
func Full() {
	ch := make(chan interface{}, 1)
	ch <- 10
	ch <- 11
}

// 可写入
func NotFull() {
	ch := make(chan interface{}, 1)
	ch <- 10
}

// panic send on closed channel
func Closed() {
	ch := make(chan interface{}, 1)
	close(ch)
	ch <- 10
}
