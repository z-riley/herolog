package main // TODO: make this "herolog" when example code finalised

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func main() {

	logServerURL := "http://0.0.0.0:2021"
	httpWriter := NewLogHTTPWriter(logServerURL)
	multiWriter := io.MultiWriter(httpWriter, os.Stdout)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log = zerolog.New(multiWriter).With().Timestamp().Logger()
	log.Info().Msg("This is a log message that will be sent to the server and printed to stdout")

	// Listen for messages
	ch := make(chan string)
	done := make(chan bool)
	defer close(ch)
	go func() {
		for msg := range ch {
			println(msg)
			log.Info().Str("i", msg).Msg("testing")
		}
		done <- true // Signal completion

	}()

	// Simulate work
	for i := 0; i < 100; i++ {
		ch <- fmt.Sprintf("%d", i)
		// time.Sleep(10 * time.Millisecond)
	}

	<-done // Wait for the signal
}
