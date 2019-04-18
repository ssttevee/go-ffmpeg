package ffmpeg

// Status represents the status of a job
type Status interface {
	isStatus()
}

// Progress represents an ffmpeg progress line
type Progress struct {
	Frame   int64
	Fps     float64
	Time    float64
	Bitrate string
	Speed   float64
}

func (*Progress) isStatus() {}

// Error represents an error that occurred during a job
type Error struct {
	Arguments []string
	error
}

func (*Error) isStatus() {}

// Done represents the completion of a job
type Done struct {
}

func (*Done) isStatus() {}
