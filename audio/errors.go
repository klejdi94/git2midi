package audio

import "errors"

var (
	// ErrFFmpegNotFound is returned when ffmpeg is not found in PATH.
	ErrFFmpegNotFound = errors.New("ffmpeg not found in PATH")
)

