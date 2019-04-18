package ffmpeg

import (
	"strings"
)

// CliOption is an command line option for ffmpeg
type CliOption interface {
	args() []string
}

type flag string

func (f flag) args() []string {
	return []string{string(f)}
}

// Flag creates a CliOption for boolean options
func Flag(name string) CliOption {
	if !strings.HasPrefix(name, "-") {
		name = "-" + name
	}

	return flag(name)
}

type option [2]string

func (o option) args() []string {
	return []string{o[0], o[1]}
}

// Option creates a CliOption
func Option(name, value string) CliOption {
	if !strings.HasPrefix(name, "-") {
		name = "-" + name
	}

	return option{name, value}
}
