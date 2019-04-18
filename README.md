# go-ffmpeg

## Usage

### Configuration

Find FFmpeg from your PATH

```go
cfg, _ := ffmpeg.DefaultConfiguration()
```

Or point to a specific installation of FFmpeg

```go
cfg, _ := ffmpeg.NewConfiguration("/path/to/ffmpeg", "/path/to/ffprobe")
```

### Transcoding

```go
job := cfg.NewJob()
job.AddInputFile("video.mp4")
job.AddOutputFile("out.avi")

statusChan, _ := job.Start(context.Background())
for status := range statusChan {
    switch v := status.(type) {
    case *ffmpeg.Progress:
        log.Printf("%#v", v)
    case *ffmpeg.Done:
        log.Printf("done")
        return
    case *ffmpeg.Error:
        log.Fatalf(v.Error())
    }
}
```

### FFprobe

Whenever an input is added to a job, ffprobe is used to validate it and the result is returned.

```go
metadata, _ := job.AddInputFile("video.mp4")
fmt.Printf("the video is %v seconds long", metadata.Format.Duration)
```

### Streams (no windows support)

```go
res, _ := http.Get("http://download.blender.org/peach/bigbuckbunny_movies/big_buck_bunny_480p_surround-fix.avi")
job.AddInputReader(res.Body)
```

```go
file, _ := os.Create("out.mp4")
job.AddOutputWriter(file)
```

### Options

Each of the following methods accept `ffmpeg.CliOption`s from `ffmpeg.Option()` and `ffmpeg.Flag()`

-   `Configuration.NewJob()`
-   `Job.AddInputFile()`
-   `Job.AddInputReader()`
-   `Job.AddOutputFile()`
-   `Job.AddOutputWriter()`

Need to overwrite a file?

```go
cfg.NewJob(ffmpeg.Flag("-y"))
```

Need to export using a specific codec and apply a filter?

```go
job.AddInputFile(
    "out.mp4",
    ffmpeg.Option("-codec:v", "libx264"),
    ffmpeg.Option("-filter:v", "scale=640:360"),
)
```

### Debugging

If there a problem with ffmpeg, you can read it's raw output like this:

```go
job.StartDebug(context.Background(), os.Stderr)
```

## Streams on windows

A somewhat convoluted way to get a data stream back from ffmpeg is to setup a http server and output to it.

```go
go http.ListenAndServe(":8080", func(w http.ResponseWriter, r *http.Request) {
    // read output data from r.Body
})

job.AddOutputFile("http://localhost:8080")
```
