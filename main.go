package main

import (
	"os"
)

func main() {
	if len(os.Args) > 1 {
		StartRestController(os.Args[1])
	} else {
		StartRestController("8080")
	}
}
