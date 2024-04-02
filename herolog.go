package herolog

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

// LogHTTPWriter attempts to send logs to a server via HTTP POST.
type LogHTTPWriter struct {
	Writer io.Writer

	serverURL     string         // The server URL to which logs will be sent
	warnOnHttpErr bool           // If true, generate warning message when HTTP logging fails
	errorLogger   zerolog.Logger // Logger for internal errors
}

// NewLogHTTPWriter creates a new instance of LogHTTPWriter.
func NewLogHTTPWriter(serverURL string, warnOnHttpErr bool) *LogHTTPWriter {
	return &LogHTTPWriter{
		serverURL:     serverURL,
		warnOnHttpErr: warnOnHttpErr,
		errorLogger:   zerolog.New(os.Stderr).With().Caller().Timestamp().Logger(),
	}
}

// Write sends the data as a log entry to the server. Logs to stderr if the HTTP request fails.
func (w *LogHTTPWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", w.serverURL, bytes.NewBuffer(p))
	if err != nil {
		if w.warnOnHttpErr {
			w.errorLogger.Warn().Err(err).Msg("Failed to create HTTP request for logging")
		}
		return len(p), err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if w.warnOnHttpErr {
			w.errorLogger.Warn().Err(err).Msg("Failed to send log via HTTP")
		}
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	return len(p), nil // Proceed as if the data was handled, regardless of HTTP errors
}
