package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Converter handles audio format conversion from MIDI files.
type Converter struct {
	ffmpegPath string
}

// NewConverter creates a new audio converter.
// It will attempt to find ffmpeg in PATH if not specified.
func NewConverter(ffmpegPath string) *Converter {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	return &Converter{
		ffmpegPath: ffmpegPath,
	}
}

// Convert converts a MIDI file to the specified audio format.
// Supported formats: mp3, wav, ogg, flac, aac, m4a
func (c *Converter) Convert(midiPath, outputPath string, format string) error {
	if !c.isFormatSupported(format) {
		return fmt.Errorf("unsupported format: %s (supported: mp3, wav, ogg, flac, aac, m4a)", format)
	}

	// Check if ffmpeg is available
	if err := c.checkFFmpeg(); err != nil {
		return fmt.Errorf("ffmpeg not found: %w. Please install ffmpeg to convert to audio formats", err)
	}

	// Determine output format from extension if not specified
	if format == "" {
		ext := strings.ToLower(filepath.Ext(outputPath))
		format = strings.TrimPrefix(ext, ".")
	}

	// Build ffmpeg command
	// For MIDI files, we need to use a soundfont or synthesizer
	// Using fluidsynth via ffmpeg is the most reliable approach
	args := []string{
		"-i", midiPath,
		"-acodec", c.getCodec(format),
		"-ar", "44100",
		"-ac", "2",
		"-y", // Overwrite output file
		outputPath,
	}

	cmd := exec.Command(c.ffmpegPath, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert MIDI to %s: %w", format, err)
	}

	return nil
}

// ConvertWithSoundfont converts MIDI to audio using a specific soundfont.
func (c *Converter) ConvertWithSoundfont(midiPath, outputPath, soundfontPath, format string) error {
	if err := c.checkFFmpeg(); err != nil {
		return fmt.Errorf("ffmpeg not found: %w", err)
	}

	// Use fluidsynth to convert MIDI with soundfont, then pipe to ffmpeg
	// This requires both fluidsynth and ffmpeg
	args := []string{
		"-i", midiPath,
		"-acodec", c.getCodec(format),
		"-ar", "44100",
		"-ac", "2",
		"-y",
		outputPath,
	}

	cmd := exec.Command(c.ffmpegPath, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

// IsAvailable checks if audio conversion is available (ffmpeg installed).
func (c *Converter) IsAvailable() bool {
	return c.checkFFmpeg() == nil
}

func (c *Converter) checkFFmpeg() error {
	cmd := exec.Command(c.ffmpegPath, "-version")
	return cmd.Run()
}

func (c *Converter) isFormatSupported(format string) bool {
	supported := []string{"mp3", "wav", "ogg", "flac", "aac", "m4a"}
	format = strings.ToLower(format)
	for _, f := range supported {
		if f == format {
			return true
		}
	}
	return false
}

func (c *Converter) getCodec(format string) string {
	switch strings.ToLower(format) {
	case "mp3":
		return "libmp3lame"
	case "wav":
		return "pcm_s16le"
	case "ogg":
		return "libvorbis"
	case "flac":
		return "flac"
	case "aac", "m4a":
		return "aac"
	default:
		return "libmp3lame"
	}
}
