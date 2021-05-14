package throttler

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net"
)

type Throttler struct {
	serverRateLimit     rate.Limit
	serverRateLimiter   *rate.Limiter
	connectionRateLimit rate.Limit
	connectionLimiter   map[net.Conn]*rate.Limiter
}

func NewThrottler() Throttler {
	return Throttler{
		connectionLimiter: map[net.Conn]*rate.Limiter{},
	}
}

// Set server bandwidth limit in bytes per second
func (t *Throttler) SetServerRateLimit(limit float64) {
	t.serverRateLimit = rate.Limit(limit)
	if t.serverRateLimiter != nil {
		t.serverRateLimiter.SetLimit(t.serverRateLimit)
		return
	}
	t.serverRateLimiter = rate.NewLimiter(t.serverRateLimit, 1)
}

// Set connection bandwidth limit in bytes per second, may be re-called during run time
func (t *Throttler) SetConnectionRateLimit(limit float64) {
	t.connectionRateLimit = rate.Limit(limit)
	for _, connection := range t.connectionLimiter {
		connection.SetLimit(t.connectionRateLimit)
	}
}

// Main function for throttling connections, takes the connection and a reader as inputs. Throttle writes a byte at a time, sending as many bytes as allowed with in the second.
func (t Throttler) Throttle(conn net.Conn, reader io.Reader) error {
	if t.connectionRateLimit == 0 || t.serverRateLimit == 0 {
		return fmt.Errorf("connection limit and server limit must be set before starting to throttle")
	}
	for {
		err := t.waitRateLimit(conn)
		if err != nil {
			return err
		}
		err = t.writeBytes(conn, reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (t Throttler) waitRateLimit(conn net.Conn) error {
	err := t.serverRateLimiter.Wait(context.TODO())
	if err != nil {
		return err
	}
	limiter, ok := t.connectionLimiter[conn]
	if !ok {
		limiter = t.registerConnection(conn)
	}
	err = limiter.Wait(context.TODO())
	return err
}

func (t *Throttler) registerConnection(conn net.Conn) *rate.Limiter {
	t.connectionLimiter[conn] = rate.NewLimiter(t.connectionRateLimit, 1)
	return t.connectionLimiter[conn]
}

func (t Throttler) writeBytes(writer io.Writer, reader io.Reader) error {
	response := make([]byte, 1)
	_, err := reader.Read(response)
	if err != nil {
		return err
	}
	_, err = writer.Write(response)
	if err != nil {
		return err
	}
	return nil
}
