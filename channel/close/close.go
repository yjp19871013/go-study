package close

//panic: close of nil channel
func Nil() {
	var ch chan interface{}
	close(ch)
}

//panic: close of closed channel
func Closed() {
	ch := make(chan interface{})
	close(ch)
	close(ch)
}
