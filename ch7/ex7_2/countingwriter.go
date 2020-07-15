/* Exercise 7.2: Write a function CountingWriter with the signature below that,
given an io.Writer, returns a new Writer that wraps the original, and a pointer
to an int64 variable that at any moment contains the number of bytes written to
the new Writer.
	func CountingWriter(w io.Writer) (io.Writer, *int64)
*/
package main

import (
	"io"
)

type countingWriter struct {
	w io.Writer
	n int64
}

// Write writes b to internal writer of countingWriter (w), updating number of
// bytes written.
func (w *countingWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		return n, err
	}

	w.n = int64(n)

	return n, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var cw countingWriter

	cw.w = w
	return &cw, &cw.n
}
