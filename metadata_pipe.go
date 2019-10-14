// +build !windows,!plan9,!nacl,!js,!nopipe

package ffmpeg

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"

	"github.com/ssttevee/go-disk-buffer"
)

// ProbeReader reads metadata from the input stream using ffprobe and
// returns an input media to be added to a job as well as the aformentioned metadata.
func (c *Configuration) ProbeReader(r io.Reader) (InputMedia, *Metadata, error) {
	cmd := c.newProbeCommand("-")

	head := buffer.New(1 << 27) // 128MiB
	cmd.Stdin = io.TeeReader(r, head)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); !ok && err != nil {
		return nil, nil, err
	}

	var metadata Metadata
	if err := json.Unmarshal(buf.Bytes(), &metadata); err != nil {
		return nil, nil, err
	}

	if metadata.Error != nil {
		return nil, nil, metadata.Error
	}

	return &inputReader{
		r:   r,
		buf: head,
	}, &metadata, nil
}
