package ffmpeg

import (
	"bufio"
	"bytes"
	"github.com/ssttevee/go-ffmpeg/internal/util"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var equalsPattern = regexp.MustCompile(`(\w+)=\s*([^ ]+)`)

func splitProgressLine(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a cr terminated line
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func parseProgress(stderr io.Reader, tee io.Writer, updateStatus func(Status)) {
	scanner := bufio.NewScanner(stderr)
	scanner.Split(splitProgressLine)
	for scanner.Scan() {
		line := scanner.Text()
		if tee != nil {
			tee.Write([]byte(line + "\n"))
		}

		if !strings.HasPrefix(line, "frame=") {
			continue
		}

		matches := equalsPattern.FindAllStringSubmatch(line, -1)
		if matches == nil {
			continue
		}

		var progress Progress
		for _, match := range matches {
			if len(match) > 1 {
				switch match[1] {
				case "frame":
					progress.Frame, _ = strconv.ParseInt(match[2], 10, 64)
				case "fps":
					progress.Fps, _ = strconv.ParseFloat(match[2], 64)
				case "time":
					progress.Time = util.ParseDuration(match[2]).Seconds()
				case "bitrate":
					progress.Bitrate = match[2]
				case "speed":
					progress.Speed, _ = strconv.ParseFloat(match[2], 64)
				}
			}
		}

		updateStatus(&progress)
	}
}
