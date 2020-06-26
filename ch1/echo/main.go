// Echo3 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func echo1() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

func echo2() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

func measure(funcName string, f func()) {
	startTime := time.Now()
	f()
	endTime := time.Since(startTime)

	fmt.Printf("Execution of %q took %fs\n", funcName, endTime.Seconds())
}

func main() {
	// Exercise 1.1: Modify the echo program to also print os.Args[0], the name of
	// the command that invoke it
	fmt.Println(os.Args[0])

	// Exercise 1.2: Modify the echo program to print the index and value of each
	// of its arguments, one per line
	for n, arg := range os.Args[1:] {
		fmt.Println(n, arg)
	}

	// Exercise 1.3: Experiment with versions of "echo". Measure each variant.
	measure("echo1", echo1)
	measure("echo2", echo2)
	measure("echo3", echo3)
}
