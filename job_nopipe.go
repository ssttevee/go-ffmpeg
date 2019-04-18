// +build windows,plan9,nacl,js

package ffmpeg

import (
	"os"
)

func (j *Job) buildArgs(args []string) ([]string, []*os.File, error) {
	var extra []*os.File

	for _, m := range j.inputs {
		args = append(args, m.options()...)

		var url string
		switch v := m.(type) {
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
		}

		args = append(args, url)
	}

	return args, extra, nil
}
