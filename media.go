package ffmpeg

// InputMedia represents some input media
type InputMedia interface {
	inputURL() string
}

// OutputMedia represents some output media
type OutputMedia interface {
	outputURL() string
}

type mediaFile string

func (m mediaFile) inputURL() string {
	return string(m)
}

func (m mediaFile) outputURL() string {
	return string(m)
}
