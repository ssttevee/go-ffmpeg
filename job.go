package ffmpeg

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type jobInput struct {
	input   InputMedia
	options []CliOption
}

type jobOutput struct {
	output  OutputMedia
	options []CliOption
}

// Job represents an ffmpeg job
type Job struct {
	cfg *Configuration

	globalOptions []CliOption

	inputs  []*jobInput
	outputs []*jobOutput
}

// NewJob creates a new job
func (c *Configuration) NewJob(options ...CliOption) *Job {
	return &Job{
		cfg:           c,
		globalOptions: options,
	}
}

// AddInput adds an input
func (j *Job) AddInput(input InputMedia, options ...CliOption) {
	j.inputs = append(j.inputs, &jobInput{
		input:   input,
		options: options,
	})
}

// AddInputFile adds an input file
func (j *Job) AddInputFile(url string, options ...CliOption) (*Metadata, error) {
	input, metadata, err := j.cfg.Probe(url)
	if err != nil {
		return nil, err
	}

	j.AddInput(input, options...)

	return metadata, nil
}

func (j *Job) addOutput(output OutputMedia, options []CliOption) {
	j.outputs = append(j.outputs, &jobOutput{
		output:  output,
		options: options,
	})
}

// AddOutputFile adds an output file
func (j *Job) AddOutputFile(file string, options ...CliOption) {
	j.addOutput(mediaFile(file), options)
}

// Start starts the job
func (j *Job) Start(ctx context.Context) (*os.Process, <-chan Status, error) {
	return j.StartDebug(ctx, nil)
}

// StartDebug starts the job and writes all output to w
func (j *Job) StartDebug(ctx context.Context, w io.Writer) (*os.Process, <-chan Status, error) {
	var args []string

	for _, option := range j.globalOptions {
		args = append(args, option.args()...)
	}

	args = append(args, "-hide_banner", "-stats")

	args, extra, err := j.buildArgs(args)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	if w != nil {
		w.Write([]byte(strings.Join(cmd.Args, " ") + "\n\n"))
	}

	go parseProgress(stderr, w, sendStatus)

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	go func() {
		defer j.cleanup()

		if err := cmd.Wait(); err != nil {
			sendFinalStatus(&Error{
				Arguments: args,
				error:     err,
			})
		} else {
			sendFinalStatus(&Done{})
		}
	}()

	return cmd.Process, statusChan, nil
}

func flattenOptions(options []CliOption) []string {
	var args []string
	for _, option := range options {
		args = append(args, option.args()...)
	}

	return args
}
