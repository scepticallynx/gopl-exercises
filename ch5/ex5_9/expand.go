/*Exercise 5.9: Write a function expand(s string, f func(string) string) string
that replaces each substring "$foo" within s by the text returned by f("foo").
*/
package main

import (
	"fmt"
	"strings"
)

func expand(s string, f func(string) string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if word[0] != '$' { // HasPrefix
			continue
		}

		words[i] = f(word[1:])
	}

	return strings.Join(words, " ")
}

// example of f func: in-place modification
func reverse(s string) string {
	_s := []byte(s)
	for i := 0; i < len(_s)/2; i++ {
		_s[i], _s[len(_s)-i-1] = _s[len(_s)-i-1], _s[i]
	}

	return string(_s)
}

// example of f: KEY=VAL
func getVar(s string) string {
	vars := map[string]string{"foo": "foo_val", "bar": "bar_val"}

	return vars[s]
}

func main() {
	var testString = `this $foo is substituted, and $foo too`

	fmt.Printf("In : %s\nOut: %s\n", testString, expand(testString, getVar))
	fmt.Printf("In : %s\nOut: %s\n", testString, expand(testString, reverse))
}
