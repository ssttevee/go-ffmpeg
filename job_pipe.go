// +build !windows,!plan9,!nacl,!js,!nopipe

package ffmpeg

import (
	"fmt"
	"io"
	"os"
)

// AddInputReader adds an input reader
func (j *Job) AddInputReader(r io.Reader, options ...CliOption) (*Metadata, error) {
	input, metadata, err := j.cfg.ProbeReader(r)
	if err != nil {
		return nil, err
	}

	j.AddInput(input, options...)

	return metadata, nil
}

// AddOutputReader adds an output reader
func (j *Job) AddOutputReader(w io.Writer, options ...CliOption) {
	j.addOutput(&outputWriter{w}, options)
}

func (j *Job) buildArgs(args []string) ([]string, []*os.File, error) {
	var extra []*os.File

	for _, m := range j.inputs {
		args = append(args, flattenOptions(m.options)...)

		var url string
		if r, ok := m.input.(io.Reader); ok {
			pr, pw, err := os.Pipe()
			if err != nil {
				return nil, nil, err
			}

			go io.Copy(pw, r)

			extra = append(extra, pr)
			url = fmt.Sprintf("/proc/self/fd/%d", 2+len(extra))
		} else {
			url = m.input.inputURL()
		}

		args = append(args, "-i", url)
	}

	for _, m := range j.outputs {
		args = append(args, flattenOptions(m.options)...)

		var url string
		if w, ok := m.output.(io.Writer); ok {
			pr, pw, err := os.Pipe()
			if err != nil {
				return nil, nil, err
			}

			go io.Copy(w, pr)

			extra = append(extra, pw)
			url = fmt.Sprintf("/proc/self/fd/%d", 2+len(extra))
		} else {
			url = m.output.outputURL()
		}

		args = append(args, url)
	}

	return args, extra, nil
}
