package ffmpeg

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

// ProbeError represents an error emitted by ffprobe
type ProbeError struct {
	Code    int64  `json:"code"`
	Message string `json:"string"`
}

// Error returns an error message
func (e *ProbeError) Error() string {
	return "probe error: " + e.Message
}

// Metadata represents the output of ffprobe
type Metadata struct {
	Streams []Stream    `json:"streams"`
	Format  *Format     `json:"format"`
	Error   *ProbeError `json:"error"`
}

// Stream represents stream metadata
type Stream struct {
	Index              int         `json:"index"`
	ID                 string      `json:"id"`
	CodecName          string      `json:"codec_name"`
	CodecLongName      string      `json:"codec_long_name"`
	Profile            string      `json:"profile"`
	CodecType          string      `json:"codec_type"`
	CodecTimeBase      string      `json:"codec_time_base"`
	CodecTagString     string      `json:"codec_tag_string"`
	CodecTag           string      `json:"codec_tag"`
	Width              int         `json:"width"`
	Height             int         `json:"height"`
	CodedWidth         int         `json:"coded_width"`
	CodedHeight        int         `json:"coded_height"`
	HasBFrames         int         `json:"has_b_frames"`
	SampleAspectRatio  string      `json:"sample_aspect_ratio"`
	DisplayAspectRatio string      `json:"display_aspect_ratio"`
	PixFmt             string      `json:"pix_fmt"`
	Level              int         `json:"level"`
	ChromaLocation     string      `json:"chroma_location"`
	Refs               int         `json:"refs"`
	QuarterSample      string      `json:"quarter_sample"`
	DivxPacked         string      `json:"divx_packed"`
	RFrameRrate        string      `json:"r_frame_rate"`
	AvgFrameRate       string      `json:"avg_frame_rate"`
	SampleRate         string      `json:"sample_rate"`
	TimeBase           string      `json:"time_base"`
	DurationTs         int         `json:"duration_ts"`
	Duration           string      `json:"duration"`
	Disposition        Disposition `json:"disposition"`
	BitRate            string      `json:"bit_rate"`
}

// Disposition represents stream disposition
type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
}

// Format represents video format
type Format struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	NbFrames       int    `json:"nb_frames"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
	Tags           Tags   `json:"tags"`
}

// Tags represents format tags
type Tags struct {
	MajorBrand       string `json:"major_brand"`
	MinorVersion     string `json:"minor_version"`
	CompatibleBrands string `json:"compatible_brands"`
	Encoder          string `json:"encoder"`
}

func (c *Configuration) newProbeCommand(url string) *exec.Cmd {
	return exec.Command(c.ffprobe, "-i", url, "-print_format", "json", "-show_format", "-show_streams", "-show_error")
}

// Probe reads metadata from the url using ffprobe and returns an
// input media to be added to a job as well as the aformentioned metadata.
func (c *Configuration) Probe(url string) (InputMedia, *Metadata, error) {
	cmd := c.newProbeCommand(url)

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

	return mediaFile(url), &metadata, nil
}
