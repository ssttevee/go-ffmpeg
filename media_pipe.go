// +build !windows,!plan9,!nacl,!js

package ffmpeg

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
)

type inputReader struct {
	opts []CliOption
	r    io.Reader

	buf bytes.Buffer
}

func (i *inputReader) options() []string {
	var args []string
	for _, option := range i.opts {
		args = append(args, option.args()...)
	}
	return args
}

func (i *inputReader) reader() (*os.File, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer w.Close()

		if _, err := i.buf.WriteTo(w); err != nil {
			log.Println(err)
			return
		}

		if _, err := io.Copy(w, i.r); err != nil {
			log.Println(err)
			return
		}
	}()

	return r, nil
}

func (i *inputReader) probe(ffprobe string) (*Metadata, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	defer r.Close()

	cmd := exec.Command(ffprobe, "-i", "/proc/self/fd/3", "-print_format", "json", "-show_format", "-show_streams", "-show_error")
	cmd.ExtraFiles = append(cmd.ExtraFiles, r)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	wait := make(chan struct{})
	go func() {
		defer close(wait)

		buf := make([]byte, 1<<12)
		for {
			n, err := i.r.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println(err)
				return
			}

			i.buf.Write(buf[:n])

			if _, err := w.Write(buf[:n]); err != nil {
				return
			}
		}
	}()

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	var metadata Metadata
	if err := json.Unmarshal(buf.Bytes(), &metadata); err != nil {
		return nil, err
	}

	<-wait

	return &metadata, nil
}

type outputWriter struct {
	opts []CliOption
	w    io.Writer
}

func (o *outputWriter) options() []string {
	var args []string
	for _, option := range o.opts {
		args = append(args, option.args()...)
	}
	return args
}

func (o *outputWriter) writer() (*os.File, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer w.Close()

		if _, err := io.Copy(o.w, r); err != nil {
			log.Println(err)
			return
		}
	}()

	return w, nil
}
