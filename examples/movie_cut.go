package main

import (
  "fmt"
  "flag"
  "log"

  ffmpeg "github.com/eladg/ffmpeg-go"
)

var (
  path        string
  input       string
  output      string
  start_time  float64
  end_time    float64
)

func main() {
  flag.StringVar (&path,       "binpath", "/usr/local/bin/ffmpeg", "ffmpeg executable path")
  flag.StringVar (&input,      "input",   "fixtures/input.mp4",    "input filename/path to cut")
  flag.StringVar (&output,     "output",  "fixtures/output.mp4",   "timmed version of the file")
  flag.Float64Var(&start_time, "start",   5.0,                     "Input video Start time. Notice: For complex format such as H.264, cuts are only possible on i-frames (https://en.wikipedia.org/wiki/Inter_frame)\nFFmpeg will not warn about 'inaccurate' cuts, so results may vary depends on the input file.")
  flag.Float64Var(&end_time,   "end",     99999.9,                 "Input video end time")
  flag.Parse()

  ff, err := ffmpeg.NewFFmpeg(path)
  if err != nil {
    log.Fatal(err)
  }
  
  ff.AddInput(&ffmpeg.Stream{
    Source: input,
    StartTime: start_time,
  })

  ff.VideoSettings = &ffmpeg.VideoSettings{
    Encoding: "copy",
  }

  ff.AudioSettings = &ffmpeg.AudioSettings{
    Encoding: "copy",
  }

  ff.AddOutput(&ffmpeg.Stream{
    Destination: output,
    Duration: end_time - start_time,
  })
  
  ff.SetOverwrite(true)
  ff.SetLogLevel("info")

  fmt.Println(ff.Command())
}