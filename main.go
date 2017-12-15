package main

import (
	"fmt"
	"os"

	"github.com/sscaling/monkey/repl"
)

func main() {
	fmt.Println("Monkey")

	repl.Start(os.Stdin, os.Stdout)
}
