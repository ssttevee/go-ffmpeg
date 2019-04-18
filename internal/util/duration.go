package util

import (
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses an duration string in the format of hh:mm:ss
func ParseDuration(dur string) time.Duration {
	parts := strings.Split(dur, ":")

	if len(parts) != 3 {
		return 0
	}

	var d time.Duration
	h, _ := strconv.ParseInt(parts[0], 10, 64)
	d += time.Duration(h) * time.Hour

	m, _ := strconv.ParseFloat(parts[1], 64)
	d += time.Duration(m) * time.Minute

	s, _ := strconv.ParseFloat(parts[2], 64)
	d += time.Duration(s) * time.Second

	return d
}
