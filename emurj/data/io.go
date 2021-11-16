package data

import (
	"os"
)

var StdoutWriter BusHandlerFunc = stdoutWriter

func stdoutWriter(bus Bus) Bus {
	if bus.WE() {
		char := bus.Data()
		buf := []byte{byte(char)}

		// emulate serial/telnet handling of LF
		if buf[0] == 10 {
			buf = []byte{033, 'D'}
		}

		_, err := os.Stdout.Write(buf)
		if err != nil {
			panic(err)
		}

		return bus.SetAck(true)
	}

	return bus
}
