// +build windows plan9 nacl js nopipe

package ffmpeg

import (
	"os"
)

func (j *Job) buildArgs(args []string) ([]string, []*os.File, error) {
	var extra []*os.File

	for _, m := range j.inputs {
		args = append(args, flattenOptions(m.options)...)
		args = append(args, "-i", m.input.inputURL())
	}

	for _, m := range j.outputs {
		args = append(args, flattenOptions(m.options)...)
		args = append(args, m.output.outputURL())
	}

	return args, extra, nil
}

func (*Job) cleanup() error {
	return nil
}
