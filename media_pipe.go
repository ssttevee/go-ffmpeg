// +build !windows,!plan9,!nacl,!js,!nopipe

package ffmpeg

import (
	"io"
)

type inputReader struct {
	r   io.Reader
	buf io.ReadCloser
}

func (i *inputReader) inputURL() string {
	panic("not implemented")
}

func (i *inputReader) Read(buf []byte) (int, error) {
	if i.buf != nil {
		n, err := i.buf.Read(buf)
		if err != io.EOF {
			return n, err
		}

		if err := i.buf.Close(); err != nil {
			return n, err
		}

		i.buf = nil
		return n, nil
	}

	return i.r.Read(buf)
}

func (i *inputReader) Close() error {
	if i.buf == nil {
		return nil
	}

	return i.buf.Close()
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
