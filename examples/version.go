package main

import (
  "fmt"
  "log"

  ffmpeg "github.com/eladg/ffmpeg-go"
)

func main() {
  path := "/usr/local/bin/ffmpeg"
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