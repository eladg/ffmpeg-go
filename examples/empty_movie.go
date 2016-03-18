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
  output      string
  size        string
  duration    float64
)

func main() {
  flag.StringVar(&path,      "binpath",  "/usr/local/bin/ffmpeg", "FFmpeg executable path")
  flag.StringVar(&output,    "output",   "fixtures/output.mp4",   "Empty blank file")
  flag.StringVar(&size,      "size",     "640x360",               "Width & Height of the empty video file")
  flag.Float64Var(&duration, "duration", 30,                      "Duration of the empty clip" )
  flag.Parse()

  sizearr      := strings.Split(size, "x")
  width, err_w := strconv.Atoi(sizearr[0])
  height,err_h := strconv.Atoi(sizearr[1])
  if err_h != nil || err_w != nil {
    fmt.Printf("Error parsing frame size input: %s",size)
    os.Exit(1)
  }

  ff := &ffmpeg.FFmpeg{
    BinPath: path,
    VideoSettings: &ffmpeg.VideoSettings{
      FrameRate:       30,
      GroupOfPictures: 30,
      QScale:          1,
      Encoding:        "libx264",
    },
    AudioSettings: &ffmpeg.AudioSettings{
      Disabled:        true,
    },
    Inputs: []*ffmpeg.Stream{
      &ffmpeg.Stream{
        Source: "/dev/zero",
        Format: "rawvideo",
        PixelFormat: "rgb24",
        Size: &ffmpeg.Size{uint16(width), uint16(height)},
      },
    },
    Outputs: []*ffmpeg.Stream{
      &ffmpeg.Stream{
        Destination: output,
        Duration: duration,
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
