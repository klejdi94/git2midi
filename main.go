package main

import (
	"flag"
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/klejdi94/git2midi/audio"
	"github.com/klejdi94/git2midi/config"
	"github.com/klejdi94/git2midi/git"
	"github.com/klejdi94/git2midi/music"
)

const (
	AppName = "git2midi"
	Version = "1.0.0"
)

func main() {
	cfg := parseFlags()

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *config.Config {
	cfg := config.NewConfig()

	flag.StringVar(&cfg.RepoPath, "repo", config.DefaultRepoPath,
		"Path to Git repository or Git repository URL (http://, https://, git://, ssh://, or git@)")
	flag.StringVar(&cfg.OutputPath, "out", config.DefaultOutputPath,
		"Output file path (MIDI or audio format: .mid, .mp3, .wav, .ogg, .flac, .aac, .m4a)")
	flag.IntVar(&cfg.BPM, "bpm", config.DefaultBPM,
		fmt.Sprintf("Tempo in BPM (default: %d for modern feel)", config.DefaultBPM))
	flag.IntVar(&cfg.Ticks, "ticks", config.DefaultTicks,
		"Ticks per quarter note")
	flag.IntVar(&cfg.Duration, "dur", config.DefaultDuration,
		fmt.Sprintf("Duration of each note in ticks (default: %d for faster playback)", config.DefaultDuration))
	flag.IntVar(&cfg.MaxCommits, "limit", 0,
		"Maximum number of commits to process (0 = all, recommended: 500-2000 for large repos)")
	flag.BoolVar(&cfg.Sample, "sample", false,
		"Evenly sample commits instead of taking first N (useful with -limit)")

	modeStr := flag.String("mode", "single-track",
		"Mode: 'single-track' or 'per-author'")

	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", AppName, Version)
		os.Exit(0)
	}

	mode, err := config.ParseMode(*modeStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	cfg.Mode = mode

	return cfg
}

func run(cfg *config.Config) error {
	fmt.Printf("Reading commits from: %s\n", cfg.RepoPath)
	commits, err := git.ParseLog(cfg.RepoPath)
	if err != nil {
		return fmt.Errorf("failed to read git log: %w", err)
	}

	if len(commits) == 0 {
		return fmt.Errorf("no commits found in repository")
	}

	originalCount := len(commits)

	if cfg.MaxCommits > 0 && len(commits) > cfg.MaxCommits {
		if cfg.Sample {
			commits, err = git.SampleCommits(commits, cfg.MaxCommits)
			if err != nil {
				return fmt.Errorf("failed to sample commits: %w", err)
			}
			fmt.Printf("Sampled %d commits from %d total\n", len(commits), originalCount)
		} else {
			commits = git.LimitCommits(commits, cfg.MaxCommits)
			fmt.Printf("Limited to first %d commits from %d total\n", len(commits), originalCount)
		}
	} else {
		fmt.Printf("Found %d commits\n", len(commits))
	}

	genCfg := &music.Config{
		BPM:      cfg.BPM,
		Ticks:    cfg.Ticks,
		Duration: cfg.Duration,
		Mode:     music.Mode(cfg.Mode),
	}

	generator := music.NewGenerator(genCfg)

	fmt.Printf("Generating MIDI composition...\n")
	writer, err := generator.Generate(commits)
	if err != nil {
		return fmt.Errorf("failed to generate MIDI: %w", err)
	}

	// Determine output format from extension
	outputExt := strings.ToLower(filepath.Ext(cfg.OutputPath))
	isAudioFormat := outputExt != "" && outputExt != ".mid" && outputExt != ".midi"

	var midiPath string
	if isAudioFormat {
		// Generate MIDI first, then convert
		midiPath = strings.TrimSuffix(cfg.OutputPath, outputExt) + ".mid"
	} else {
		midiPath = cfg.OutputPath
	}

	fmt.Printf("Writing MIDI file to: %s\n", midiPath)
	if err := writer.WriteFile(midiPath); err != nil {
		return fmt.Errorf("failed to write MIDI file: %w", err)
	}

	if isAudioFormat {
		fmt.Printf("Converting to %s format...\n", strings.TrimPrefix(outputExt, "."))
		converter := audio.NewConverter("")
		if !converter.IsAvailable() {
			fmt.Fprintf(os.Stderr, "Warning: ffmpeg not found. Audio conversion skipped.\n")
			fmt.Fprintf(os.Stderr, "Install ffmpeg to convert MIDI to audio formats.\n")
			fmt.Printf("MIDI file saved as: %s\n", midiPath)
			return nil
		}

		format := strings.TrimPrefix(outputExt, ".")
		if err := converter.Convert(midiPath, cfg.OutputPath, format); err != nil {
			return fmt.Errorf("failed to convert to audio: %w", err)
		}

		// Remove temporary MIDI file if conversion successful
		if err := os.Remove(midiPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove temporary MIDI file: %v\n", err)
		}

		fmt.Printf("Successfully generated audio file: %s\n", cfg.OutputPath)
	} else {
		fmt.Printf("Successfully generated MIDI file with %d track(s)\n", writer.TrackCount())
	}

	return nil
}
