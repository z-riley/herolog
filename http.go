package main // TODO: make this "herolog" when example code finalised

import (
	"bytes"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

// LogHTTPWriter attempts to send logs to a server via HTTP POST.
type LogHTTPWriter struct {
	serverURL   string         // The server URL to which logs will be sent
	errorLogger zerolog.Logger // Logger for internal errors
}

// NewLogHTTPWriter creates a new instance of LogHTTPWriter.
func NewLogHTTPWriter(serverURL string) *LogHTTPWriter {
	return &LogHTTPWriter{
		serverURL:   serverURL,
		errorLogger: zerolog.New(os.Stderr).With().Caller().Timestamp().Logger(),
	}
}

// Write sends the data as a log entry to the server. Logs to stderr if the HTTP request fails.
func (w *LogHTTPWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", w.serverURL, bytes.NewBuffer(p))
	if err != nil {
		w.errorLogger.Warn().Err(err).Msg("Failed to create HTTP request for logging")
		return len(p), err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.errorLogger.Warn().Err(err).Msg("Failed to send log via HTTP")
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	return len(p), nil // Proceed as if the data was handled, regardless of HTTP errors
}
