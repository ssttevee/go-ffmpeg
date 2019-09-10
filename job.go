package ffmpeg

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// Job represents an ffmpeg job
type Job struct {
	cfg *Configuration

	globalOptions []CliOption

	inputs  []inputMedia
	outputs []outputMedia
}

// NewJob creates a new job
func (c *Configuration) NewJob(options ...CliOption) *Job {
	return &Job{
		cfg:           c,
		globalOptions: options,
	}
}

func (j *Job) addInput(input inputMedia) (*Metadata, error) {
	metadata, err := input.probe(j.cfg.ffprobe)
	if err != nil {
		return nil, err
	}

	j.inputs = append(j.inputs, input)

	return metadata, nil
}

// AddInputFile adds an input file
func (j *Job) AddInputFile(file string, options ...CliOption) (*Metadata, error) {
	return j.addInput(&mediaFile{
		opts: options,
		path: file,
	})
}

func (j *Job) addOutput(output outputMedia) {
	j.outputs = append(j.outputs, output)
}

// AddOutputFile adds an output file
func (j *Job) AddOutputFile(file string, options ...CliOption) {
	j.outputs = append(j.outputs, &mediaFile{
		opts: options,
		path: file,
	})
}

// Start starts the job
func (j *Job) Start(ctx context.Context) (<-chan Status, error) {
	return j.start(ctx, nil)
}

// StartDebug starts the job and writes all output to w
func (j *Job) StartDebug(ctx context.Context, w io.Writer) (<-chan Status, error) {
	return j.start(ctx, w)
}

func (j *Job) start(ctx context.Context, debug io.Writer) (<-chan Status, error) {
	var args []string

	for _, option := range j.globalOptions {
		args = append(args, option.args()...)
	}

	args = append(args, "-hide_banner", "-stats")

	args, extra, err := j.buildArgs(args)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, j.cfg.ffmpeg, args...)
	cmd.ExtraFiles = extra

	var mu sync.Mutex
	var finished bool

	statusChan := make(chan Status, 1)
	sendStatusUnsafe := func(status Status) {
		// flush channel to prevent the user from reading out of date progress
		select {
		case <-statusChan:
		default:
		}

		statusChan <- status
	}

	sendStatus := func(status Status) {
		mu.Lock()
		defer mu.Unlock()

		if finished {
			return
		}

		sendStatusUnsafe(status)
	}

	sendFinalStatus := func(status Status) {
		mu.Lock()
		defer mu.Unlock()

		defer close(statusChan)

		finished = true

		sendStatusUnsafe(status)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if debug != nil {
		debug.Write([]byte(strings.Join(cmd.Args, " ") + "\n\n"))
	}

	go parseProgress(stderr, debug, sendStatus)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			sendFinalStatus(&Error{
				Arguments: args,
				error:     err,
			})
		} else {
			sendFinalStatus(&Done{})
		}
	}()

	return statusChan, nil
}
