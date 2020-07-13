/* Exercise 5.19: Use panic and recover to write a function that contains no return
statement yet returns a non-zero value.
*/
package main

import "fmt"

var value interface{}

func noReturn(i int) {
	defer func() {
		value = recover()
	}()

	if i > 0 {
		panic(fmt.Sprintf("Artificial panic. i=%d. Set i to 0 for natural panic.", i))
	}

	i += i / i
}

func main() {
	for i := 2; i >= 0; i-- {
		noReturn(i)

		fmt.Printf("Return from recover(): %s\n", value)
	}
}
