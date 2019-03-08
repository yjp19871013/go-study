package main

import (
	"fmt"
	"os"
)

func main() {
	name := os.Getenv("MY_PROGRAM_NAME")
	version := os.Getenv("MY_PROGRAM_VERSION")
	fmt.Println("Name: ", name)
	fmt.Println("Version: ", version)
}
