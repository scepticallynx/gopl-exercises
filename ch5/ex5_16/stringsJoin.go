/*Exercise 5.16: Write a variadic version of strings.Join
 */
package main

import "fmt"

func join(separator string, s ...string) (result string) {
	for n := 0; n < len(s); n++ {
		if n == len(s)-1 {
			separator = ""
		}
		result = result + s[n] + separator
	}

	return
}

func main() {
	fmt.Println(join(";", []string{"one", "2", "test this stuff", "join"}...))
}
