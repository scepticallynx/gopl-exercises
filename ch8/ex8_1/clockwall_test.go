package main

import "testing"

func TestParseArg(t *testing.T) {
	testCases := []struct {
		input          string
		tz, host, port string
		err            error
	}{
		{
			input: "NewYork=localhost:8010",
			tz:    "NewYork",
			host:  "localhost",
			port:  "8010",
			err:   nil,
		},
		{
			input: "Tokyo=localhost:8020",
			tz:    "Tokyo",
			host:  "localhost",
			port:  "8020",
			err:   nil,
		},
		{
			input: "Kyiv",
			tz:    "",
			host:  "",
			port:  "",
			err:   missingTZSeparator,
		},
		{
			input: "Kyiv=localhost",
			tz:    "",
			host:  "",
			port:  "",
			err:   missingHostPortSeparator,
		},
		{
			input: "Kyiv=192.168.1.1:9999999",
			tz:    "Kyiv",
			host:  "192.168.1.1",
			port:  "",
			err:   invalidPort,
		},
	}

	for n, test := range testCases {
		tz, host, port, err := parseArg(test.input)
		if err != nil {
			if test.err == nil {
				t.Errorf("Case %2d: unexpected err: %v\n", n, err)
			} else if test.err.Error() != err.Error() {
				t.Errorf("Case %2d: expected err: %v, got: %v\n", n, test.err, err)
			}
		} else {
			if test.err != nil {
				t.Errorf("Case %2d: expected err: %v\n", n, test.err)
			}
		}

		if tz != test.tz {
			t.Errorf("Case %2d: expected tz: %q, got: %q\n", n, test.tz, tz)
		}

		if host != test.host {
			t.Errorf("Case %2d: expected host: %q, got: %q\n", n, test.host, host)
		}

		if port != test.port {
			t.Errorf("Case %2d: expected port: %q, got: %q\n", n, test.port, port)
		}
	}
}
