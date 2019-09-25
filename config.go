package ffmpeg

import (
	"sync"

	"github.com/ssttevee/go-ffmpeg/internal/util"
)

var (
	defaultConfig      Configuration
	defaultConfigError error
	defaultConfigOnce  sync.Once
)

func configurationFromEnvironment() {
	defaultConfig.ffmpeg, defaultConfig.ffmpegVersion, defaultConfigError = util.FindBinary("ffmpeg")
	if defaultConfigError != nil {
		return
	}

	defaultConfig.ffprobe, defaultConfig.ffprobeVersion, defaultConfigError = util.FindBinary("ffprobe")
	if defaultConfigError != nil {
		return
	}
}

// Configuration represents valid paths to ffmpeg and ffprobe
type Configuration struct {
	ffmpeg         string
	ffmpegVersion  string
	ffprobe        string
	ffprobeVersion string
}

// DefaultConfiguration looks for and
// returns a configuration from the environment
func DefaultConfiguration() (*Configuration, error) {
	defaultConfigOnce.Do(configurationFromEnvironment)
	if defaultConfigError != nil {
		return nil, defaultConfigError
	}

	return &defaultConfig, nil
}

// NewConfiguration validates the given paths to
// ffmpeg and ffprobe and returns a configuration
func NewConfiguration(ffmpeg, ffprobe string) (*Configuration, error) {
	ffmpegVersion, err := util.TestBinary("ffmpeg", ffmpeg)
	if err != nil {
		return nil, err
	}

	ffprobeVersion, err := util.TestBinary("ffprobe", ffprobe)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		ffmpeg:         ffmpeg,
		ffmpegVersion:  ffmpegVersion,
		ffprobe:        ffprobe,
		ffprobeVersion: ffprobeVersion,
	}, nil
}
