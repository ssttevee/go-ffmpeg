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

type Configuration struct {
	ffmpeg         string
	ffmpegVersion  string
	ffprobe        string
	ffprobeVersion string
}

func DefaultConfiguration() (*Configuration, error) {
	defaultConfigOnce.Do(configurationFromEnvironment)
	if defaultConfigError != nil {
		return nil, defaultConfigError
	}

	return &defaultConfig, nil
}

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
