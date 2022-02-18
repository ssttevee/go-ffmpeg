package util

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func which(name string) (string, error) {
	tool := "which"
	if runtime.GOOS == "windows" {
		tool = "where"
	}

	cmd := exec.Command(tool, name)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

// TestBinary tests an ffmpeg-family binary
func TestBinary(name, path string) (version string, _ error) {
	cmd := exec.Command(path, "-version")

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	pattern, err := regexp.Compile("^" + regexp.QuoteMeta(name) + " version (.+) Copyright .* the FFmpeg developers\r?\n")
	if err != nil {
		return "", err
	}

	matches := pattern.FindStringSubmatch(buf.String())
	if matches == nil {
		return "", errors.New("invalid " + name + " binary")
	}

	return matches[1], nil
}

// FindBinary finds a in the path and tests if it is an ffmpeg-family binary
func FindBinary(name string) (path string, version string, err error) {
	path, err = which(name)
	if err != nil {
		return "", "", err
	}

	version, err = TestBinary(name, path)
	if err != nil {
		return "", "", err
	}

	return
}
