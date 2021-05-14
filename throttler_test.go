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

	tests := []struct{
		Name string
		ServerRateLimit float64
		ConnectionRateLimit float64
		NumberOfConnections int
		Message []byte
		ExpectedTimeTakeInSeconds float64
	} {
		{
			Name: "Should handle bandwidth limit for multiple connections",
			ServerRateLimit: 1000,
			ConnectionRateLimit: 10,
			NumberOfConnections: 10,
			Message: []byte("1231112311"),  // 10 byte message
			ExpectedTimeTakeInSeconds: 1,
		},
		{
			Name: "Should slow requests to handle bandwidth limit for multiple connections",
			ServerRateLimit: 1000,
			ConnectionRateLimit: 10,
			NumberOfConnections: 10,
			Message: []byte("12311123111231112311"),  // 20 byte message
			ExpectedTimeTakeInSeconds: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			//Given
			throttler := NewThrottler()
			throttler.SetServerRateLimit(tt.ServerRateLimit) // 1000 bytes/sec
			throttler.SetConnectionRateLimit(tt.ConnectionRateLimit) // 10 byte/sec
			l, _ := net.Listen("tcp", ":3323")
			defer l.Close()

			go func() {
				for {
					conn, err := l.Accept()
					if err != nil {
						break
					}
					go func() {
						byteReader := bytes.NewBuffer(tt.Message)
						throttler.Throttle(conn, byteReader)
					}()
				}
			}()
			allComplete := make(chan bool, 1)
			startTime := time.Now()

			//When
			for i := 0; i < tt.NumberOfConnections; i++ {
				go func() {
					startTime := time.Now()
					c, err := net.Dial("tcp", ":3323")
					if err != nil {
						t.Errorf("could not start client: %v", err.Error())
					}
					defer c.Close()
					expectedResponse := tt.Message
					currentIndex := 0
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
					timeTaken := time.Now().Sub(startTime).Seconds()
					if timeTaken  > tt.ExpectedTimeTakeInSeconds || timeTaken  < tt.ExpectedTimeTakeInSeconds - 1 {
						t.Errorf("incorrect timing to receive entire response: took %v expected %v", timeTaken, tt.ExpectedTimeTakeInSeconds)
						allComplete <- false
					}
					allComplete <- true
				}()
			}

			//Then
			numberComplete := 0
			for {
				if <- allComplete {
					numberComplete ++
				}

				if numberComplete == 10 || time.Now().Sub(startTime).Seconds() > tt.ExpectedTimeTakeInSeconds {
					break
				}
			}
			if numberComplete != 10 {
				t.Errorf("did not all complete")
			}
		})
	}

}
