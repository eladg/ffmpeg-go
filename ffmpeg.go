package ffmpeg

import (
  "fmt"
  "os"
  "bufio"
  "regexp"
  "os/exec"
  "strings"
)

type Stream struct {
  FilePath string
}

type FFmpegOption struct {
  Name  string
  Value string 
}

type FFmpegVideoStream struct {
}

type FFmpegAudioStream struct {
}

type FFmpeg struct {
  BinPath       string
  Version       string
  Configuration []string
  Simulate      bool
  Inputs        []Stream
  Outputs       []Stream
  Options       []FFmpegOption

  Video         *FFmpegVideoStream
  Audio         *FFmpegAudioStream
}

func NewFFmpeg(binpath string) (*FFmpeg, error) {
  f := new(FFmpeg)
  var err error;

  f.BinPath = binpath
  if _, err = os.Stat(binpath); os.IsNotExist(err) {
    return nil, fmt.Errorf("could not find ffmpeg binary (%s)",binpath)
  }

  if f.Version, f.Configuration, err = version_run(binpath) ; err != nil {
    return nil, fmt.Errorf("Could not initialize ffmpeg from binary (%s)",binpath)
  }

  f.Simulate = true

  return f, nil
}

func version_run(path string) (string, []string, error) {
  var (
    cmdOut []byte
    err    error
  )
  cmdName := "ffmpeg"
  cmdArgs := []string{"-version"}
  
  if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
    return "", []string{}, fmt.Errorf("There was an error running ffmpeg -version command (%v)",err)
  }

  // version
  output := strings.Split(string(cmdOut),"\n")
  rc, _ := regexp.Compile(`\d\.\d\.\d`)
  version := rc.FindAllString(output[0], -1)[0]

  // configuration
  rg, _ := regexp.Compile(`--[\w*-=/]*\s*`)
  configuration := rg.FindAllString(output[2], -1)

  return version, configuration, nil
}

func build() {
  // docker build current directory
  cmdName := "ffmpeg"
  cmdArgs := []string{"build", "."}

  cmd := exec.Command(cmdName, cmdArgs...)
  cmdReader, err := cmd.StdoutPipe()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
    os.Exit(1)
  }

  scanner := bufio.NewScanner(cmdReader)
  go func() {
    for scanner.Scan() {
      fmt.Printf("docker build out | %s\n", scanner.Text())
    }
  }()

  err = cmd.Start()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
    os.Exit(1)
  }

  err = cmd.Wait()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
    os.Exit(1)
  }
}