package tasks

import (
	"io"
)

// NewLoggedWriteCloser creates a new LoggedWriteCloser
// that wraps source
func NewLoggedWriteCloser(source io.WriteCloser, details, log io.Writer) *LoggedWriteCloser {
	return &LoggedWriteCloser{
		Writer: io.MultiWriter(source, details, log),
		source: source,
	}
}

// LoggedWriteCloser is a wrapper
// io.WriteCloser that writes everything it's written
// both to the source writer and to stdout
type LoggedWriteCloser struct {
	io.Writer
	source io.WriteCloser
}

// Close the source io.WriteCloser
func (log *LoggedWriteCloser) Close() error {
	return log.source.Close()
}
