package ffmpeg

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

type inputMedia interface {
	options() []string
	probe(ffprobe string) (*Metadata, error)
}

type outputMedia interface {
	options() []string
}

type mediaFile struct {
	opts []CliOption
	path string
}

func (m *mediaFile) options() []string {
	var args []string
	for _, option := range m.opts {
		args = append(args, option.args()...)
	}

	return args
}

func (m *mediaFile) file() string {
	return m.path
}

func (m *mediaFile) probe(ffprobe string) (*Metadata, error) {
	cmd := exec.Command(ffprobe, "-i", m.path, "-print_format", "json", "-show_format", "-show_streams", "-show_error")

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var metadata Metadata
	if err := json.Unmarshal(buf.Bytes(), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}
