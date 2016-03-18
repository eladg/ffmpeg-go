package main

import (
  "fmt"
  "log"
  "flag"

  ffmpeg "github.com/eladg/ffmpeg-go"
)

var (
  path        string
)

func main() {
  flag.StringVar (&path,"binpath", "/usr/local/bin/ffmpeg", "ffmpeg executable path")
  flag.Parse()

  ff, err := ffmpeg.NewFFmpeg(path)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Found ffmpeg.")
  fmt.Println("> Binary:       ", ff.BinPath)
  fmt.Println("> Version:      ", ff.Version)
  fmt.Println("> Configuration:")
  for _, c := range ff.Configuration {
    fmt.Printf("  > %s\n",c)
  }
  fmt.Println("")
}
