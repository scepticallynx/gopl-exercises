package popcount

// pc[i] is the population count of i.
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(
		pc[byte(x>>(0*8))] +
			pc[byte(x>>(1*8))] +
			pc[byte(x>>(2*8))] +
			pc[byte(x>>(3*8))] +
			pc[byte(x>>(4*8))] +
			pc[byte(x>>(5*8))] +
			pc[byte(x>>(6*8))] +
			pc[byte(x>>(7*8))])
}

// Exercise 2.3: Rewrite PopCount to use loop instead of single expression. Compare
// the performance of the two versions.
func PopCountLoop(x uint64) int {
	var b byte
	for i := 0; i < 8; i++ {
		b += pc[byte(x>>(i*8))]
	}
	return int(b)
}

// Exercise 2.4: Write a version of PopCount that counts bits by shifting its argument
// through 64 bit position, testing the rightmost bit each time. Compare its performance
// to the table-lookup version.

// Exercise 2.5:
