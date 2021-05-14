# TCP Throttler

[![GoDoc](https://godoc.org/github.com/bchalk101/tcp-throttler?status.png)](https://godoc.org/github.com/bchalk101/tcp-throttler)

Simple package to enable tcp bandwidth throttling

It is possible to set bandwidth limit per server and bandwidth limit per connection. It is also possible to make the changes during runtime.

For a 30s transfer sample, consumed bandwidth is accurate to at least +/- 5%

Example:

```go
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
```

## Testing
This package consists of both unit tests and an integration test. The unit tests focusing on a couple of the simpler functionality, however it does not test the actual limiting. 

The integration tests the full functionality of the throttler, setting up a tcp server and monitoring as it accepts connections. These tests however are longer to run, with one running up to 30 seconds to verify that the throttler handles longer time frames.

Running Unit Tests: `go test -run Unit .`

Running Int Tests: `go test -run Int .`