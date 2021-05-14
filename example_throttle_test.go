package throttler

import (
	"bytes"
	"net"
)

// This example shows how to use the throttler to throttle the bandwidth for tcp connections
func ExampleThrottler_Throttle() {
	// Setup the throttler
	throttler := NewThrottler()
	throttler.SetServerRateLimit(1000) // 1000 bytes/sec or 1kb/sec
	throttler.SetConnectionRateLimit(10) // 10 byte/sec

	// Create a tcp server
	l, _ := net.Listen("tcp", ":3323")
	defer l.Close()

	// Accept tcp connections and throttle the response
	logFile := []byte("LOTS AND LOTS OF LOGS")
	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}
		go func() {
			byteReader := bytes.NewBuffer(logFile)
			throttler.Throttle(conn, byteReader)
		}()
	}

	// Output:
	// Will serve the log file at 10 bytes per second for each connection, limited over all at 1kb per second for all connections
}