// +build !windows,!plan9,!nacl,!js,!nopipe

package ffmpeg

import (
	"bytes"
	"io"
)

type inputReader struct {
	r   io.Reader
	buf *bytes.Buffer
}

func (i *inputReader) inputURL() string {
	panic("not implemented")
}

func (i *inputReader) Read(buf []byte) (int, error) {
	if i.buf != nil {
		n, err := i.buf.Read(buf)
		if err == io.EOF {
			i.buf = nil
			return n, nil
		}

		return n, err
	}

	return i.r.Read(buf)
}

type outputWriter struct {
	w io.Writer
}

func (o *outputWriter) outputURL() string {
	panic("not implemented")
}

func (o *outputWriter) Write(buf []byte) (int, error) {
	return o.w.Write(buf)
}
