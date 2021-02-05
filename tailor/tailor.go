// Tideland Go Library - Tailor
//
// Copyright (C) 2014-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package tailor

//--------------------
// IMPORTS
//--------------------

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

const (
	defaultbuffsz   = 4096
	defaultPollTime = time.Second
	delimiter       = '\n'
)

var (
	delimiters = []byte{delimiter}
)

// Error messages ...
const (
	ErrNoSource      = "cannot start scroller: no source"
	ErrNoTarget      = "cannot start scroller: no target"
	ErrNegativeLines = "negative number of lines not allowed: %d"
)

// Tailor scrolls and filters a ReadSeeker line by line and
// writes the data into a Writer.
type Tailor struct {
	shallStop chan bool
	source    io.ReadSeeker
	target    io.Writer
	buffsz    int
	reader    *bufio.Reader
	writer    *bufio.Writer
}

// New creates a Tailor for the given source and target.
func New(source io.ReadSeeker, target io.Writer, buffsz int) *Tailor {
	return &Tailor{
		source:    source,
		target:    target,
		buffsz:    buffsz,
		shallStop: make(chan bool),
		reader:    bufio.NewReaderSize(source, buffsz),
		writer:    bufio.NewWriter(target),
	}
}

// Stop will causes the tail operation to ends.
// The log file is readed until EOF before stopping.
func (s *Tailor) Stop() {
	s.shallStop <- true
}

// Start is the goroutine for reading, filtering and writing.
// The returned chan emit an error in case of failure, otherwise
// it's closed upon tail completion without emitting nothing.
func (s *Tailor) Start() chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

	InitialPositioning:
		/*if err := s.seekInitial(); err != nil {
			errs <- err
			return
		}*/
		shallStop := false
		select {
		case <-s.shallStop:
			shallStop = true
		default:
		}

		for {
			line, readErr := s.readLine()
			//fmt.Println(string(line))
			_, writeErr := s.writer.Write(line)
			if writeErr != nil {
				errs <- writeErr
				return
			}
			if readErr != nil {
				if readErr != io.EOF {
					errs <- readErr
					return
				}
				break
			}
		}

		if writeErr := s.writer.Flush(); writeErr != nil {
			errs <- writeErr
			return
		}

		if !shallStop {
			time.Sleep(200 * time.Millisecond)
			goto InitialPositioning
		}

	}()

	return errs
}

// seekInitial sets the initial position to start reading. This
// position depends on the number lines and the filter st.
func (s *Tailor) seekInitial() error {
	offset, err := s.source.Seek(0, os.SEEK_END)
	if err != nil {
		return err
	}
	seekPos := int64(0)
	found := 0
	buffer := make([]byte, s.buffsz)
SeekLoop:
	for offset > 0 {
		// bufferf partly filled, check if large enough.
		space := cap(buffer) - len(buffer)
		if space < s.buffsz {
			// Grow buffer.
			newBuffer := make([]byte, len(buffer), cap(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
			space = cap(buffer) - len(buffer)
		}
		if int64(space) > offset {
			// Use exactly the right amount of space if there's
			// only a small amount remaining.
			space = int(offset)
		}
		// Copy remaining data to the end of the buffer.
		copy(buffer[space:cap(buffer)], buffer)
		buffer = buffer[0 : len(buffer)+space]
		offset -= int64(space)
		_, err := s.source.Seek(offset, os.SEEK_SET)
		if err != nil {
			return err
		}
		// Read into the beginning of the buffer.
		_, err = io.ReadFull(s.source, buffer[0:space])
		if err != nil {
			return err
		}
		// Find the end of the last line in the buffer.
		// This will discard any unterminated line at the end
		// of the file.
		end := bytes.LastIndex(buffer, delimiters)
		if end == -1 {
			// No end of line found - discard incomplete
			// line and continue looking. If this happens
			// at the beginning of the file, we don't care
			// because we're going to stop anyway.
			buffer = buffer[:0]
			continue
		}
		end++
		for {
			start := bytes.LastIndex(buffer[0:end-1], delimiters)
			if start == -1 && offset >= 0 {
				break
			}
			start++
			found++
			seekPos = offset + int64(start)
			break SeekLoop
		}
		// Leave the last line in the buffer. It's not
		// clear if it is complete or not.
		buffer = buffer[0:end]
	}
	// Final positioning.
	s.source.Seek(seekPos, os.SEEK_SET)
	return nil
}

// readLine reads the next valid line from the reader, even if it is
// larger than the reader buffer.
func (s *Tailor) readLine() ([]byte, error) {
	for {
		slice, err := s.reader.ReadSlice(delimiter)
		if err == nil {
			return slice, nil
		}
		line := append([]byte(nil), slice...)
		for err == bufio.ErrBufferFull {
			slice, err = s.reader.ReadSlice(delimiter)
			line = append(line, slice...)
		}
		switch err {
		case nil:
			return line, nil
		case io.EOF:
			// Reached EOF without a delimiter,
			// so step back for next time.
			s.source.Seek(-int64(len(line)), os.SEEK_CUR)
			return nil, err
		default:
			return nil, err
		}
	}
}

// EOF
