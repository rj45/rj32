package data

import (
	"os"
)

var StdoutWriter BusHandlerFunc = stdoutWriter

func stdoutWriter(bus Bus) Bus {
	if bus.WE() {
		char := bus.Data()
		buf := []byte{byte(char)}
		if char > 127 {
			buf = append(buf, byte(char>>8))
		}
		_, err := os.Stdout.Write(buf)
		if err != nil {
			panic(err)
		}

		return bus.SetAck(true)
	}

	return bus
}
