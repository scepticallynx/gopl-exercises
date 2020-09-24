/*Ex 8.1: modify clock2 to accept a port number. Write clockwall program that
reads time from several servers at once and outputs it to table.

clockwall NewYork=localhost:8010 Tokyo=localhost:8020
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func netcat(output io.Writer, addr string) {
	for {

	}
}

var (
	missingTZSeparator       = errors.New("missing '=' separator between timezone and server")
	missingHostPortSeparator = errors.New("missing ':' separator between host and port")
	invalidHost              = errors.New("invalid host value")
	invalidPort              = errors.New("invalid port value")
)

// parseArg parses each argument specified to clockwall
//
// Argument is expected to be a string in the following format:
//  Timezone=host:port
//
// It returns timezone, host and port, or parsing error.
//
// For current use case valid ports are in range 8000-8050
func parseArg(arg string) (tz, host, port string, err error) {
	if !strings.ContainsRune(arg, '=') {
		err = missingTZSeparator
		return
	}

	if !strings.ContainsRune(arg, ':') {
		err = missingHostPortSeparator
		return
	}

	_split := strings.Split(arg, "=")
	tz = _split[0]

	host, port, err = net.SplitHostPort(_split[1])
	if err != nil {
		return tz, "", "", err
	}

	if host != "localhost" {
		if net.ParseIP(host) == nil {
			err = fmt.Errorf("%v: %q", invalidHost, host)
			return
		}
	}

	portValue, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		return tz, host, "", invalidPort
	}

	if portValue < 8000 || portValue > 8050 {
		return tz, host, "", invalidPort
	}

	return
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("clockwall: not a single timezone is specified")
	}

	timezones := make(map[string]string)
	for n, arg := range os.Args[1:] {
		if tz, h, p, err := parseArg(arg); err != nil {
			fmt.Println("argument #:", n, "error:", err)
		} else {
			timezones[tz] = net.JoinHostPort(h, p)
		}
	}

	if len(timezones) == 0 {
		os.Exit(0)
	}

	output := func(w io.Writer, tz string, r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if _, err := fmt.Fprintf(w, "%s: %s\t", tz, scanner.Bytes()); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	for tz, addr := range timezones {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer conn.Close()

		go output(os.Stdout, tz, conn)
	}

	for {
		fmt.Fprintln(os.Stdout)
		time.Sleep(1 * time.Second)
	}
}
