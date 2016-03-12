# ffmpeg-go
> **WORK IN PROGRESS: THIS IS A VERY EARLY STAGE - A DESIGN IDEA FOR GO FFmpeg WRAPPER**

FFmpeg-go is ffmpeg for all of your Go needs

FFmpeg
-----------------
###### Private:
- binPath (string)
- configuration ([]string)
- execute (string)

###### Public:
- Simulate (bool)
- Version (string)
- Inputs ([]Filepath)
- Outputs ([]Filepath)
- GlobalFlags ([]string)
- VideoStream (FFmpegVideo)
- AudioStream (FFmpegAudio)

###### Methods:

FFmpegVideoStream
-----------------
- Info (FFmpegStreamInfo)
- Attributes ([]FFmpegFlag)
- Filter ([]FFmpegFilter)

###### Methods:
- Rate
- Encoding
- Gop
- Qscale
- Disable
- Copy
- Bitrate
- StartTime
- Duration
- Format

FFmpegAudioStream
-----------------
- Info (FFmpegStreamInfo)
- Attributes ([]FFmpegFlag)
- Filter ([]FFmpegFilter)

###### Methods:
- Encoding
- Qscale
- Disable
- Copy
- SampleRate
- Channels
- Quality
- BitRate
- StartTime
- Duration

FFmpegStreamInfo
-----------------

FFmpegFlag
-----------------
- Name
- Value

###### Methods:
String() string

FFmpegFilter
-----------------
- Id
- Name
- Values
- InputStreams
- OutputStreams
- RecursiveFilter



