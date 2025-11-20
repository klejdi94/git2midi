package music

import (
	"hash/fnv"
	"sort"

	"github.com/klejdi94/git2midi/git"
	"github.com/klejdi94/git2midi/midi"
)

// Generator generates MIDI music from Git commits.
type Generator struct {
	config *Config
}

// Config holds configuration for music generation.
type Config struct {
	BPM      int
	Ticks    int
	Duration int
	Mode     Mode
}

// Mode represents the generation mode.
type Mode int

const (
	// ModeSingleTrack generates a single MIDI track.
	ModeSingleTrack Mode = iota
	// ModePerAuthor generates separate tracks per author.
	ModePerAuthor
)

// NewGenerator creates a new music generator with the given configuration.
func NewGenerator(cfg *Config) *Generator {
	return &Generator{
		config: cfg,
	}
}

// Generate creates MIDI tracks from the given commits.
func (g *Generator) Generate(commits []git.Commit) (*midi.Writer, error) {
	if len(commits) == 0 {
		return nil, ErrNoCommits
	}

	format := uint16(0)
	if g.config.Mode == ModePerAuthor {
		format = 1
	}

	writer := midi.NewWriter(format, uint16(g.config.Ticks))
	tempo := midi.BPMToMicrosecondsPerQuarter(g.config.BPM)

	switch g.config.Mode {
	case ModeSingleTrack:
		if err := g.generateSingleTrack(writer, commits, tempo); err != nil {
			return nil, err
		}
	case ModePerAuthor:
		if err := g.generatePerAuthorTracks(writer, commits, tempo); err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidMode
	}

	return writer, nil
}

// generateSingleTrack generates a single MIDI track from all commits.
func (g *Generator) generateSingleTrack(writer *midi.Writer, commits []git.Commit, tempo uint32) error {
	track := midi.NewTrack()
	track.AddTempo(0, tempo)

	currentTime := uint32(0)
	for i, commit := range commits {
		pitch := g.hashToPitch(commit.Hash)
		velocity := g.messageToVelocity(commit.Message)
		deltaTime := g.calculateRhythm(i, uint32(g.config.Duration))

		track.AddNoteOn(currentTime, 0, pitch, velocity)
		track.AddNoteOff(deltaTime, 0, pitch, 64)

	currentTime = deltaTime * 3 / 4
	}

	track.AddEndOfTrack(0)
	writer.AddTrack(track)
	return nil
}

func (g *Generator) generatePerAuthorTracks(writer *midi.Writer, commits []git.Commit, tempo uint32) error {
	authorCommits := make(map[string][]git.Commit)
	for _, commit := range commits {
		authorCommits[commit.Author] = append(authorCommits[commit.Author], commit)
	}

	authors := make([]string, 0, len(authorCommits))
	for author := range authorCommits {
		authors = append(authors, author)
	}
	sort.Strings(authors)

	for channel, author := range authors {
		if channel > 15 {
			channel = 15
		}

		track := midi.NewTrack()
		track.AddTempo(0, tempo)

		currentTime := uint32(0)
		for i, commit := range authorCommits[author] {
			pitch := g.hashToPitch(commit.Hash)
			velocity := g.messageToVelocity(commit.Message)
			deltaTime := g.calculateRhythm(i, uint32(g.config.Duration))

			track.AddNoteOn(currentTime, byte(channel), pitch, velocity)
			track.AddNoteOff(deltaTime, byte(channel), pitch, 64)

			currentTime = deltaTime * 3 / 4
		}

		track.AddEndOfTrack(0)
		writer.AddTrack(track)
	}

	return nil
}

// hashToPitch converts a Git commit hash to a MIDI note using a pentatonic minor scale.
// Maps to MIDI range C4-C6 (60-84) for a focused, musical range.
func (g *Generator) hashToPitch(hash string) byte {
	h := fnv.New32a()
	h.Write([]byte(hash))
	hashValue := h.Sum32()

	pentatonic := []byte{0, 3, 5, 7, 10}
	scaleDegree := hashValue % uint32(len(pentatonic))
	octaveOffset := (hashValue / uint32(len(pentatonic))) % 3

	note := 60 + int(pentatonic[scaleDegree]) + int(octaveOffset)*12

	if note > 84 {
		note = 84
	}
	if note < 60 {
		note = 60
	}

	return byte(note)
}

// messageToVelocity converts commit message length to MIDI velocity (40-127).
func (g *Generator) messageToVelocity(message string) byte {
	const (
		maxMessageLength = 200
		minVelocity      = 40
		maxVelocity      = 127
	)

	length := len(message)
	if length > maxMessageLength {
		length = maxMessageLength
	}

	velocity := minVelocity + (length*(maxVelocity-minVelocity))/maxMessageLength
	if velocity > maxVelocity {
		velocity = maxVelocity
	}
	if velocity < minVelocity {
		velocity = minVelocity
	}

	return byte(velocity)
}

// calculateRhythm calculates note duration with rhythm variations.
func (g *Generator) calculateRhythm(index int, baseDuration uint32) uint32 {
	if index == 0 {
		return baseDuration
	}

	if index%4 == 0 {
		return baseDuration * 3 / 4
	}
	if index%8 == 0 {
		return baseDuration * 5 / 4
	}

	return baseDuration
}

