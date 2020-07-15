/*Exercise 7.5: The LimitReader from io accepts an io.Reader r and a number of
bytes n, and returns another Reader that reads from r but reports an end-of-file
condition after n bytes. Implement it.
	func LimitReader(r io.Reader, n int64) io.Reader
*/
package main

import "io"

type limitReader struct {
	R io.Reader
	N int64
}

func (r *limitReader) Read(b []byte) (n int, err error) {
	if r.N <= 0 {
		return 0, io.EOF
	}

	if int64(len(b)) > r.N {
		b = b[:r.N]
	}

	n, err = r.R.Read(b)
	r.N -= int64(n)
	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &limitReader{r, n}
}
