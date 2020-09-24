/*Exercise 8.3: In netcat3 modify the main goroutine to close only the write half
of the connection (CloseWrite) so that the program will continue to print the
final echoes from the reverb1 server even after the standard input has been closed.
*/
package main

import (
	"io"
	"log"
	"net"
	"os"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	c, ok := conn.(*net.TCPConn)
	if !ok {
		log.Printf("not a TCP connection: %T", c)
		return
	}
	c.CloseWrite()
	<-done // wait for background goroutine to finish
}
