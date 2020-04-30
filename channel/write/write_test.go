package write

import "testing"

func TestNil(t *testing.T) {
	Nil()
}

func TestFull(t *testing.T) {
	Full()
}

func TestNotFull(t *testing.T) {
	NotFull()
}

func TestClosed(t *testing.T) {
	Closed()
}
