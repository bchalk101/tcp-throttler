package throttler

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_Throttler(t *testing.T) {
	t.Run("Should limit bandwidth to 1byte per second for entire server", func(t *testing.T) {
		//Given
		throttler := NewThrottler()
		throttler.SetServerRateLimit(1) // 1 byte/sec
		//When
		writer := &bytes.Buffer{}
		byteReader := bytes.NewReader([]byte("1"))
		throttler.Throttle(writer, byteReader)
		//Then
		if string(writer.Bytes()) != "1" {
			fmt.Println(string(writer.Bytes()))
			t.Errorf("did not read byte")
		}
	})

	t.Run("Should limit bandwidth for connection", func(t *testing.T) {
		//Given
		throttler := NewThrottler()
		throttler.SetServerRateLimit(1) // 1 byte/sec
	})
}
