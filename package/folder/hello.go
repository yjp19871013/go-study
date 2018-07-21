package myPackage

import (
	"fmt"
)

func init() {
	fmt.Println("In init")
}

func Hello(word string) {
	fmt.Println("Hello", word)
}
