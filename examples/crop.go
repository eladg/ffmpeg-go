package main

import (
  "flag"
  "strings"
  "strconv"
  "os"
  "fmt"

  ffmpeg "github.com/eladg/ffmpeg-go"
)

var (
  path        string
  input       string
  output      string
  crop        string
)

func ParseSizeSring() (uint16,uint16) {
  
}

func main() {
  flag.StringVar(&path,     "binpath",   "/usr/local/bin/ffmpeg", "FFmpeg executable path")
  flag.StringVar(&input,    "input",     "fixtures/input.mp4",    "input filename/path to cut")
  flag.StringVar(&output,   "output",    "fixtures/output.mp4",   "Re-encoded file output")
  flag.StringVar(&cropSize, "crop-size", "640x360",               "crop size")
  flag.StringVar(&crop,     "crop-size", "640x360",               "crop size")
  flag.Parse()

  croparr      := strings.Split(s, "x")
  width, err_w := strconv.Atoi(croparr[0])
  if err_w != nil {
    return nil, err_w
  }

  height, err_h := strconv.Atoi(croparr[1])
  if err_h != nil {
    return nil, err_h
  }


  crop := &ffmpeg.VideoFilter{
    Name: "crop"
    Values: []{"640","360","0","0"}
    Inputs: []{"in"}
    Outputs: []{"out"}
  }

  ff := &ffmpeg.FFmpeg{
    BinPath: path,
    VideoSettings: &ffmpeg.VideoSettings{
      FrameRate:       30/1.001,
      QScale:          1,
      Encoding:        "libx264",
    },
    AudioSettings: &ffmpeg.AudioSettings{
      Encoding:        "copy",
    },
    Inputs: []*ffmpeg.Stream{
      &ffmpeg.Stream{
        Source: input,
        StartTime: start_time,
      },
    },

    Outputs: []*ffmpeg.Stream{
      &ffmpeg.Stream{
        Destination: output,
      },
    },
  }

  ff.SetOverwrite(true)
  ff.SetLogLevel("info")

  fmt.Println("Config:")
  fmt.Printf("  > Binpath: %s\n",path)
  fmt.Printf("  > Frame Size: %dx%d\n",width, height)
  fmt.Printf("  > Duration: %f\n",duration)
  fmt.Printf("  > Output File: %s\n",output)
  fmt.Println("")
  fmt.Printf("Executing: %s %s\n\n",path, ff.Command())
  ff.Run()
  fmt.Println("Done.")
}
