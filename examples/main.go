package main

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/zac460/herolog"
)

var log zerolog.Logger

func main() {

	multiWriter := io.MultiWriter(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		herolog.NewLogHTTPWriter("http://0.0.0.0:2021"),
	)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log = zerolog.New(multiWriter).With().Timestamp().Logger()

	// Simulate work
	for i := 0; i < 1000; i++ {
		log.Info().Int("i", i).Msg("testing")
		time.Sleep(50 * time.Microsecond)
	}

}
