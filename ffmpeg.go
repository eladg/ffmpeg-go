package ffmpeg

import (
  "fmt"
  "os"
  "bufio"
  "regexp"
  "os/exec"
  "strings"
  "strconv"
)

type Point struct {
  X uint16
  Y uint16
}

type Size struct {
  Width  uint16
  Height uint16
}

func NewSize(s string) (*Size, error) {
  sizearr      := strings.Split(s, "x")

  width, err_w := strconv.Atoi(sizearr[0])
  if err_w != nil {
    return nil, err_w
  }

  height, err_h := strconv.Atoi(sizearr[1])
  if err_h != nil {
    return nil, err_h
  }

  return &Size{uint16(width), uint16(height)}
}

func (s *Size)String() string {
  return fmt.Sprintf("%dx%d",s.Width, s.Height)
}

func (s *Size)ColonString() string {
  return fmt.Sprintf("%d:%d",s.Width, s.Height)
}

const OptionOverwrite            = "y"
const OptionDuration             = "t"
const OptionStartTime            = "ss"
const OptionPixelFormat          = "pix_fmt"
const OptionFormat               = "f"
const OptionLogLevel             = "loglevel"
const OptionStreamInput          = "i"
const OptionStreamOutput         = ""
const OptionVideoFrameRate       = "r"
const OptionVideoSize            = "s"
const OptionVideoEncoding        = "c:v"
const OptionVideoGroupOfPictures = "g"
const OptionVideoQscale          = "q:v"
const OptionVideoDisable         = "vn"
const OptionVideoBitrate         = "b:v"
const OptionAudioSampleRate      = "ar"
const OptionAudioEncoding        = "c:a"
const OptionAudioChannels        = "ac"
const OptionAudioQscale          = "q:a"
const OptionAudioDisable         = "an"
const OptionAudioBitrate         = "b:a"

type Option struct {
  Name  string
  Value string
}

func NewOption(name, value string) (*Option) {
  return &Option{name, value}
}

func (o *Option) String() string {
  if o.Value == "" {
    return o.Name
  } else {
    return o.Name + " " + o.Value
  }
}

func OptionsToString(options []*Option) (string) {
  str := ""
  for _, opt := range options {
    if opt.Name == OptionStreamOutput {
      str = str + fmt.Sprintf("%s",opt)
    } else {
      str = str + fmt.Sprintf("-%s ",opt)
    }
  }
  return str
}

type Stream struct {
  Format      string
  PixelFormat string
  Size        *Size
  Source      string
  Destination string
  Duration    float64
  StartTime   float64
}

func (s *Stream) toFFmpegOptionString() (string) {
  options := make([]*Option, 0)

  if s.Format != "" {
    options = append(options, NewOption(OptionFormat, s.Format))
  }
  if s.PixelFormat != "" {
    options = append(options, NewOption(OptionPixelFormat, s.PixelFormat))
  }
  if s.Size != nil {
    options = append(options, NewOption(OptionVideoSize, fmt.Sprintf("%s",s.Size)))
  }
  if s.Duration != 0 {
    dur := strconv.FormatFloat(s.Duration, 'f', -1, 32)
    options = append(options, NewOption(OptionDuration, dur))
  }
  if s.StartTime != 0 {
    time := strconv.FormatFloat(s.StartTime, 'f', -1, 32)
    options = append(options, NewOption(OptionStartTime, time))
  }
  if s.Source != "" {
    options = append(options, NewOption(OptionStreamInput, s.Source))
  }
  if s.Destination != "" {
    options = append(options, NewOption(OptionStreamOutput, s.Destination))
  }

  return OptionsToString(options)
}

type VideoFilter struct {
  Id      string
  Name    string
  Values  []string
  Inputs  []string
  Outputs []string
  Filter  *VideoFilter
}

type AudioFilter struct {
  Id      string
  Name    string
  Values  []string
  Inputs  []string
  Outputs []string
  Filter  *VideoFilter
}

type VideoSettings struct {
  Filter *VideoFilter
  FrameRate       float64
  BitRate         uint32
  Encoding        string
  GroupOfPictures uint16
  QScale          uint8
  Disabled        bool
}

func (s *VideoSettings) toFFmpegOptionString() (string) {
  options := make([]*Option, 0)
  if s.FrameRate != 0.0 {
    fr := strconv.FormatFloat(s.FrameRate, 'f', -1, 32)
    options = append(options, NewOption(OptionVideoFrameRate, fr))
  }
  if s.BitRate != 0 {
    options = append(options, NewOption(OptionVideoBitrate, strconv.Itoa(int(s.BitRate))))
  }
  if s.Encoding != "" {
    options = append(options, NewOption(OptionVideoEncoding, s.Encoding))
  }
  if s.GroupOfPictures != 0 {
    options = append(options, NewOption(OptionVideoGroupOfPictures, strconv.Itoa(int(s.GroupOfPictures))))
  }
  if s.QScale != 0 {
    options = append(options, NewOption(OptionVideoQscale, strconv.Itoa(int(s.QScale))))
  }
  if s.Disabled != false {
    options = append(options, NewOption(OptionVideoDisable, ""))
  }

  return OptionsToString(options)
}

type AudioSettings struct {
  Filter *AudioFilter
  SamepleRate     uint32
  BitRate         uint32
  Channels        uint8
  Encoding        string
  QScale          uint8
  Disabled        bool
}

func (s *AudioSettings) toFFmpegOptionString() (string) {
  options := make([]*Option, 0)
  if s.SamepleRate != 0 {
    options = append(options, NewOption(OptionAudioSampleRate, strconv.FormatUint(uint64(s.SamepleRate),10)))
  }
  if s.BitRate != 0 {
    options = append(options, NewOption(OptionAudioBitrate, strconv.FormatUint(uint64(s.BitRate),10)))
  }
  if s.Channels != 0 {
    options = append(options, NewOption(OptionAudioChannels, strconv.FormatUint(uint64(s.Channels),10)))
  }
  if s.Encoding != "" {
    options = append(options, NewOption(OptionAudioEncoding, s.Encoding))
  }
  if s.QScale != 0 {
    options = append(options, NewOption(OptionAudioQscale, strconv.FormatUint(uint64(s.QScale),10)))
  }
  if s.Disabled != false {
    options = append(options, NewOption(OptionAudioDisable, ""))
  }

  return OptionsToString(options)
}

type FFmpeg struct {
  BinPath       string
  Version       string
  Configuration []string
  Options       []*Option
  Inputs        []*Stream
  Outputs       []*Stream

  VideoSettings *VideoSettings
  AudioSettings *AudioSettings

  Simulate      bool
  Overwrite     bool
  LogLevel      string
}

func NewFFmpeg(binpath string) (*FFmpeg, error) {
  f := new(FFmpeg)
  var err error;

  f.BinPath = binpath
  if _, err = os.Stat(binpath); os.IsNotExist(err) {
    return nil, fmt.Errorf("could not find ffmpeg binary (%s)",binpath)
  }

  if f.Version, f.Configuration, err = f.dryVersionRun() ; err != nil {
    return nil, fmt.Errorf("Could not initialize ffmpeg from binary (%s)",binpath)
  }

  // default options
  f.SetOverwrite(true)
  f.HideBunner()

  return f, nil
}

func (f *FFmpeg) AddGlobalOption(opt *Option) {
  f.Options = append(f.Options, opt)
}

func (f *FFmpeg) AddInput(s *Stream) {
  f.Inputs = append(f.Inputs, s)
}

func (f *FFmpeg) AddOutput(s *Stream) {
  f.Outputs = append(f.Outputs, s)
}

func (f *FFmpeg) Command() (string) {
  cmd := make([]string,5)
  cmd = append(cmd, OptionsToString(f.Options))             // globals

  inputs := ""                                              // inputs
  for _, inpt := range f.Inputs { inputs = inputs + inpt.toFFmpegOptionString() }
  cmd = append(cmd, inputs)

  cmd = append(cmd, f.VideoSettings.toFFmpegOptionString()) // video settings
  cmd = append(cmd, f.AudioSettings.toFFmpegOptionString()) // audio settings

  outputs := ""                                             // outputs
  for _, outp := range f.Outputs { outputs = outputs + outp.toFFmpegOptionString() }
  cmd = append(cmd, outputs)

  return strings.Join(cmd,"")
}

func (f *FFmpeg) HideBunner() {
  f.AddGlobalOption(&Option{"hide_banner",""})
}

func (f *FFmpeg) SetOverwrite(overwrite bool) {
  if overwrite {
    f.AddGlobalOption(&Option{"y",""})
  } else {
    f.AddGlobalOption(&Option{"n",""})
  }
}

func (f *FFmpeg) SetLogLevel(level string) {
  f.AddGlobalOption(&Option{"loglevel",level})
}

func (f *FFmpeg) dryVersionRun() (string, []string, error) {
  var (
    cmdOut []byte
    err    error
  )
  cmdName := f.BinPath
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

func (f *FFmpeg) Run() {
  cmdName := f.BinPath
  cmdArgs := strings.Fields(f.Command())

  cmd := exec.Command(cmdName,cmdArgs...)
  cmdReader, err := cmd.StdoutPipe()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
    os.Exit(1)
  }

  scanner := bufio.NewScanner(cmdReader)
  go func() {
    for scanner.Scan() {
      fmt.Printf("ffmpeg returned | %s\n", scanner.Text())
    }
  }()

  err = cmd.Start()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error running FFmpeg", err)
    os.Exit(1)
  }

  err = cmd.Wait()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error waiting for FFmpeg", err)
    os.Exit(1)
  }
}
