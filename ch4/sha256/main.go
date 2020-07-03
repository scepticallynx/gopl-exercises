// Exercise 4.1: Write a function that counts the number of bits that are different in two SHA256 hashes.
// Exercise 4.2: Write a program that prints the SHA256 hash of its standard input
// by default but supports a command-line flag to print the SHA384 or SHA512 hash instead.
package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	input := os.Stdin
	stat, _ := input.Stat()

	// checking for empty pipe
	if stat.Mode()&os.ModeNamedPipe == 0 {
		log.Fatal("No input provided")
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	var chosenSHA string
	flag.StringVar(&chosenSHA, "sha", "256", "Select SHA: 256, 384, 512. Default: 256")
	flag.Parse()

	switch chosenSHA {
	case "512":
		fmt.Printf("SHA512 sum: %x\n", sha512.Sum512(bytes))
	case "384":
		fmt.Printf("SHA384 sum: %x\n", sha512.Sum384(bytes))
	default:
		fmt.Printf("SHA256 sum: %x\n", sha256.Sum256(bytes))
	}
}
