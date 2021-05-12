package throttler

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

func Test_Throttler(t *testing.T) {
	t.Run("Should write a byte to response", func(t *testing.T) {
		//Given
		throttler := NewThrottler()
		throttler.SetServerRateLimit(1) // 1 byte/sec
		//When
		writer := &bytes.Buffer{}
		byteReader := bytes.NewReader([]byte("1"))
		throttler.writeBytes(writer, byteReader)
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

func Test_IntThrottler(t *testing.T) {
	throttler := NewThrottler()
	throttler.SetServerRateLimit(2)
	throttler.SetConnectionRateLimit(1)
	l, _ := net.Listen("tcp", ":3323")

	go func() {
		for {
			conn, _ := l.Accept()
			defer l.Close()
			byteReader := bytes.NewBuffer([]byte("123111231112311123111231112311"))
			throttler.Throttle(conn, byteReader)
		}
	}()

	startTime := time.Now()
	c, err := net.Dial("tcp", ":3323")
	if err != nil {
		t.Errorf("could not start client: %v", err.Error())
	}
	expectedResponse := []byte("123111231112311123111231112311")
	currentIndex := 0
	defer c.Close()
	for {
		response, _ := bufio.NewReader(c).ReadByte()
		if response != expectedResponse[currentIndex] {
			t.Errorf("not receing the correct bytes: %v", string(response))
		}
		currentIndex++
		if currentIndex == len(expectedResponse) {
			break
		}
	}
	endTime := time.Now()
	if endTime.Sub(startTime).Seconds() > 30 || endTime.Sub(startTime).Seconds() < 29 {
		t.Errorf("took too long to get the expected response")
	}
}
