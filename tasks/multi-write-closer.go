package tasks

import (
	"io"
)

// NewMultiWriteCloser creates a new MultiWriteCloser
// that wraps source
func NewMultiWriteCloser(closer io.WriteCloser, writers ...io.Writer) *MultiWriteCloser {
	writers = append(writers, closer)
	return &MultiWriteCloser{
		Writer: io.MultiWriter(writers...),
		closer: closer,
	}
}

// MultiWriteCloser is a wrapper
// io.WriteCloser that writes everything it's written
// both to the source writer and to stdout
type MultiWriteCloser struct {
	io.Writer
	closer io.WriteCloser
}

// Close the source io.WriteCloser
func (mwc *MultiWriteCloser) Close() error {
	return mwc.closer.Close()
}
