// +build !windows,!plan9,!nacl,!js

package ffmpeg

import (
	"fmt"
	"io"
	"os"
)

// AddInputReader adds an input reader
func (j *Job) AddInputReader(r io.Reader, options ...CliOption) (*Metadata, error) {
	return j.addInput(&inputReader{
		opts: options,
		r:    r,
	})
}

// AddOutputReader adds an output reader
func (j *Job) AddOutputReader(w io.Writer, options ...CliOption) {
	j.addOutput(&outputWriter{
		opts: options,
		w:    w,
	})
}

func (j *Job) buildArgs(args []string) ([]string, []*os.File, error) {
	var extra []*os.File

	for _, m := range j.inputs {
		args = append(args, m.options()...)

		var url string
		switch v := m.(type) {
		case *inputReader:
			f, err := v.reader()
			if err != nil {
				return nil, nil, err
			}

			extra = append(extra, f)
			url = fmt.Sprintf("/proc/self/fd/%d", 2+len(extra))
		case *mediaFile:
			url = v.file()
		}

		args = append(args, "-i", url)
	}

	for _, m := range j.outputs {
		args = append(args, m.options()...)

		var url string
		switch v := m.(type) {
		case *mediaFile:
			url = v.file()
		case *outputWriter:
			f, err := v.writer()
			if err != nil {
				return nil, nil, err
			}

			extra = append(extra, f)
			url = fmt.Sprintf("/proc/self/fd/%d", 2+len(extra))
		}

		args = append(args, url)
	}

	return args, extra, nil
}
