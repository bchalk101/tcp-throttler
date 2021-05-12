package throttler

import (
	"context"
	"golang.org/x/time/rate"
	"io"
	"net"
)

type Throttler struct {
	serverRateLimit     *rate.Limiter
	connectionRateLimit *rate.Limiter
	connectionUsage     map[net.Conn]int
}

func NewThrottler() Throttler {
	return Throttler{}
}

// Set server bandwidth limit in bytes per second
func (t *Throttler) SetServerRateLimit(limit float64) {
	t.serverRateLimit = rate.NewLimiter(rate.Limit(limit), 1)
}

// Set connection bandwidth limit in bytes per second
func (t *Throttler) SetConnectionRateLimit(limit float64) {
	t.connectionRateLimit = rate.NewLimiter(rate.Limit(limit), 1)
}

func (t Throttler) Throttle(writer io.Writer, reader io.Reader) error {
	for {
		err := t.serverRateLimit.Wait(context.TODO())
		if err != nil {
			return err
		}
		response := make([]byte, 1)
		_, err = reader.Read(response)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = writer.Write(response)
		if err != nil {
			return err
		}
	}
	return nil
}
