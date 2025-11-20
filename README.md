# git2midi

A Go program that transforms Git commit history into musical compositions by converting commit properties into MIDI events. Supports both MIDI and audio format output (MP3, WAV, OGG, FLAC, AAC, M4A).

## Overview

`git2midi` reads the commit history of any Git repository (local path or remote URL) and generates a valid MIDI file (Format 0 or Format 1) where musical notes are derived from commit properties such as hash, timestamp, author, and message length. Remote repositories are automatically cloned to a temporary directory and cleaned up after processing.

## Features

- **Commit-to-Music Mapping**:
  - **Pitch**: Derived from commit hash using **pentatonic minor scale** for modern, harmonious sound
  - **Velocity**: Derived from commit message length (40-127 for musical range)
  - **Duration**: Configurable note duration (default: 120 ticks for faster playback)
  - **Tempo**: Configurable BPM (default: 140 for modern feel)
  - **Rhythm**: Automatic rhythm variations (staccato and emphasis patterns)
  - **Authors**: Can be mapped to different MIDI channels (per-author mode)
  - **Commit Limiting**: Limit commits processed to keep music length reasonable (recommended: 500-2000)
  - **Sampling**: Evenly sample commits across history for better representation

- **MIDI Standards Compliance**:
  - Supports MIDI Format 0 (single track) and Format 1 (multi-track)
  - Proper variable-length encoding for delta times
  - Valid MIDI file structure (header chunk + track chunks)
  - Compatible with all major DAWs and MIDI players

- **Repository Support**:
  - Works with local Git repositories
  - Supports remote repository URLs (GitHub, GitLab, etc.)
  - Automatic temporary cloning and cleanup for URLs
  - Supports HTTP, HTTPS, Git, SSH protocols

- **Audio Format Support**:
  - Generate MIDI files (.mid)
  - Convert to MP3, WAV, OGG, FLAC, AAC, M4A (requires ffmpeg)
  - Automatic format detection from file extension
  - High-quality audio conversion

- **CLI Interface**:
  - Simple command-line flags
  - Flexible configuration options
  - Clear error messages

## Installation

### Prerequisites

- Go 1.20 or later (for building from source)
- Git (for reading repository history)
- ffmpeg (optional, for audio format conversion)

### Installation Options

**Option 1: Download Pre-built Binaries (Recommended)**

Download the latest release for your platform:
- [Windows](https://github.com/klejdi94/git2midi/releases/latest) - `git2midi-windows-amd64.exe`
- [Linux](https://github.com/klejdi94/git2midi/releases/latest) - `git2midi-linux-amd64`
- [macOS](https://github.com/klejdi94/git2midi/releases/latest) - `git2midi-darwin-amd64`

**Option 2: Install via Go**

```bash
go install github.com/klejdi94/git2midi@latest
```

**Option 3: Build from Source**

```bash
git clone https://github.com/klejdi94/git2midi.git
cd git2midi
go build -o git2midi
```

**For Audio Format Support:**

Install ffmpeg:
- **Windows**: Download from [ffmpeg.org](https://ffmpeg.org/download.html) or use `choco install ffmpeg`
- **Linux**: `sudo apt-get install ffmpeg` (Ubuntu/Debian) or `sudo dnf install ffmpeg` (Fedora)
- **macOS**: `brew install ffmpeg`

## Usage

### Basic Usage

Generate a MIDI file from the current directory's Git repository:

```bash
./git2midi -repo . -out commits.mid
```

Generate a MIDI file directly from a Git repository URL (no local clone needed):

```bash
./git2midi -repo https://github.com/user/repo.git -out commits.mid
```

### Command-Line Flags

- `-repo <path|url>`: Path to local Git repository or Git repository URL (default: current directory `.`)
  - Supports local paths: `.`, `/path/to/repo`, `../other-repo`
  - Supports URLs: `https://github.com/user/repo.git`, `http://...`, `git://...`, `ssh://...`, `git@github.com:user/repo.git`
  - URLs are automatically cloned to a temporary directory and cleaned up after processing
- `-out <path>`: Output file path (default: `commits.mid`)
  - Supports MIDI: `.mid`, `.midi`
  - Supports Audio (requires ffmpeg): `.mp3`, `.wav`, `.ogg`, `.flac`, `.aac`, `.m4a`
  - Format is automatically detected from file extension
- `-bpm <number>`: Tempo in beats per minute (default: `140` for modern feel)
- `-ticks <number>`: Ticks per quarter note (default: `480`)
- `-dur <number>`: Duration of each note in ticks (default: `120` for faster playback)
- `-limit <number>`: Maximum number of commits to process (default: `0` = all commits)
  - Recommended: `500-2000` for large repositories to keep music length reasonable
  - Use with `-sample` to evenly distribute commits across history
- `-sample`: Evenly sample commits instead of taking first N (useful with `-limit`)
- `-mode <mode>`: Generation mode - `single-track` or `per-author` (default: `single-track`)

### Examples

**Generate from a huge repository (Spring Framework example):**
```bash
./git2midi -repo https://github.com/spring-projects/spring-framework.git -out spring.mp3 -limit 2000 -sample -bpm 140
```

**Generate from Linux kernel (massive repository):**
```bash
./git2midi -repo https://github.com/torvalds/linux.git -out linux.mp3 -limit 1500 -sample
```

**Generate from a large repository as MP3:**
```bash
./git2midi -repo https://github.com/user/huge-repo.git -out sampled.mp3 -limit 1500 -sample
```

**Generate from a GitHub repository URL (with commit limit for large repos):**
```bash
./git2midi -repo https://github.com/torvalds/linux.git -out linux.mid -limit 1000 -sample
```

**Generate from a local repository with modern defaults:**
```bash
./git2midi -repo /path/to/repo -out modern.mid
```

**Generate a fast-paced composition:**
```bash
./git2midi -repo . -out fast.mid -bpm 180 -dur 60
```

**Generate with per-author tracks (limited commits):**
```bash
./git2midi -repo . -out authors.mid -mode per-author -limit 800
```

**Generate as MP3 audio file (requires ffmpeg):**
```bash
./git2midi -repo . -out commits.mp3 -limit 1000
```

**Generate as WAV audio file:**
```bash
./git2midi -repo https://github.com/spring-projects/spring-framework.git -out spring.wav -limit 2000 -sample
```

**Generate as OGG audio file:**
```bash
./git2midi -repo . -out commits.ogg
```

**Create separate tracks for each author:**
```bash
./git2midi -repo . -out authors.mid -mode per-author -bpm 100
```

**Generate a slow, ambient piece:**
```bash
./git2midi -repo . -out ambient.mid -bpm 60 -dur 960
```

**High-resolution timing:**
```bash
./git2midi -repo . -out precise.mid -ticks 960 -dur 480
```

### Playing MIDI Files from Terminal

**Quick Reference:**
- **Windows**: `start commits.mid`
- **Linux (Ubuntu/Debian)**: `sudo apt-get install timidity && timidity commits.mid`
- **Linux (Fedora/RHEL)**: `sudo dnf install timidity++ && timidity commits.mid`
- **macOS**: `brew install timidity && timidity commits.mid`

**Detailed Instructions:**

**Windows (PowerShell/CMD):**
```powershell
# Open with default MIDI player (Windows Media Player)
start commits.mid

# Or using PowerShell
Start-Process commits.mid
```

**Windows (with third-party tools):**

If you have **Timidity++** installed:
```bash
timidity commits.mid
```

If you have **FluidSynth** installed:
```bash
fluidsynth -a winmm commits.mid
```

**Linux Terminal:**

**Option 1: Timidity++ (Recommended - Simple and reliable)**
```bash
# Install on Ubuntu/Debian
sudo apt-get install timidity timidity-daemon

# Install on Fedora/RHEL
sudo dnf install timidity++

# Install on Arch Linux
sudo pacman -S timidity++

# Play MIDI file
timidity commits.mid

# Play with specific soundfont
timidity -x "soundfont /usr/share/sounds/sf2/FluidR3_GM.sf2" commits.mid
```

**Option 2: FluidSynth (High-quality software synthesizer)**
```bash
# Install on Ubuntu/Debian
sudo apt-get install fluidsynth fluid-soundfont-gm

# Install on Fedora/RHEL
sudo dnf install fluidsynth soundfont-fluid

# Install on Arch Linux
sudo pacman -S fluidsynth soundfont-fluid

# Play MIDI file (requires ALSA)
fluidsynth -a alsa -m alsa_seq /usr/share/sounds/sf2/FluidR3_GM.sf2 commits.mid

# Or with PulseAudio
fluidsynth -a pulseaudio -m alsa_seq /usr/share/sounds/sf2/FluidR3_GM.sf2 commits.mid
```

**Option 3: aplaymidi (ALSA MIDI player)**
```bash
# Install on Ubuntu/Debian
sudo apt-get install alsa-utils

# Install on Fedora/RHEL
sudo dnf install alsa-utils

# Install on Arch Linux
sudo pacman -S alsa-utils

# Play MIDI file (requires MIDI sequencer running)
aplaymidi commits.mid

# Or with specific port
aplaymidi -p 128:0 commits.mid
```

**Option 4: Python with pygame (No installation needed if Python is available)**
```bash
# Install pygame if needed
pip3 install pygame

# Play MIDI file
python3 -c "import pygame; pygame.mixer.init(); pygame.mixer.music.load('commits.mid'); pygame.mixer.music.play(); import time; time.sleep(60)"
```

**Option 5: VLC (If already installed)**
```bash
# Play MIDI file
vlc commits.mid

# Or in background
vlc --intf dummy commits.mid
```

**Option 6: WildMIDI (Lightweight option)**
```bash
# Install on Ubuntu/Debian
sudo apt-get install wildmidi

# Install on Arch Linux
sudo pacman -S wildmidi

# Play MIDI file
wildmidi commits.mid
```

**Option 7: Using xdg-open (Opens with default application)**
```bash
# Opens with system default MIDI player
xdg-open commits.mid
```

**macOS Terminal:**
```bash
# Using timidity (install via Homebrew)
brew install timidity

# Play MIDI file
timidity commits.mid

# Or open with default player
open commits.mid
```

**Alternative: Online Players**
- Upload the MIDI file to online MIDI players
- Use DAW software (Ableton, FL Studio, Reaper, etc.)
- Use VLC media player: `vlc commits.mid`

## How It Works

### Commit-to-Music Mapping

1. **Pitch Selection** (Modern Pentatonic Scale):
   - The commit hash is hashed using FNV-1a
   - Notes are mapped to a **pentatonic minor scale** (C, D#, F, G, A#) for a modern, pleasant sound
   - Focused range: C4 to C6 (MIDI notes 60-84) for musical coherence
   - Ensures deterministic but harmonically pleasing pitch selection

2. **Velocity Mapping**:
   - Commit message length determines note velocity
   - Shorter messages → lower velocity (quieter)
   - Longer messages → higher velocity (louder)
   - Mapped to range 40-127 for musical expressiveness

3. **Temporal Structure** (Modern Rhythm):
   - Each commit becomes a note event
   - Notes are played sequentially in commit order (oldest first)
   - **Rhythm variations**: Every 4th note is shorter (staccato), every 8th note is longer (emphasis)
   - Notes overlap slightly (75% of duration) for a modern, flowing sound
   - Default duration: 120 ticks (faster than before) for snappier playback
   - Default tempo: 140 BPM for a modern feel

4. **Commit Limiting & Sampling**:
   - Use `-limit` to cap the number of commits processed (prevents hours-long files)
   - Recommended: 500-2000 commits for reasonable music length (2-10 minutes)
   - Use `-sample` with `-limit` to evenly distribute commits across history
   - Without sampling, takes the first N commits chronologically

5. **Author Separation** (per-author mode):
   - Each unique author gets their own MIDI track
   - Authors are assigned to different MIDI channels (0-15)
   - Enables polyphonic composition with author-specific voices
   - All tracks use the same modern rhythm and scale patterns

### MIDI File Structure

The generated MIDI files follow the standard MIDI file format:

- **Header Chunk** (`MThd`):
  - Format type (0 or 1)
  - Number of tracks
  - Time division (ticks per quarter note)

- **Track Chunks** (`MTrk`):
  - Variable-length delta times
  - MIDI events (Note On/Off, Tempo, End of Track)
  - Proper meta events for tempo and track termination

### Technical Details

- **Variable-Length Encoding**: MIDI delta times use variable-length quantities where each byte contains 7 bits of data and 1 continuation bit
- **Event Timing**: All events use relative timing (delta times) from the previous event
- **Tempo Events**: Set Tempo meta events (0xFF 0x51) specify microseconds per quarter note
- **Note Events**: Standard MIDI Note On (0x9n) and Note Off (0x8n) events

## Project Structure

```
git2midi/
├── main.go              # CLI entry point
├── go.mod               # Go module definition
├── LICENSE              # MIT License
├── .gitignore          # Git ignore rules
├── README.md           # This file
├── audio/              # Audio conversion package
│   ├── converter.go    # Audio format conversion (ffmpeg)
│   └── errors.go       # Audio errors
├── config/             # Configuration package
│   └── config.go       # Configuration and validation
├── git/                # Git package
│   ├── commits.go      # Commit data structures
│   └── log.go          # Git log parsing
├── midi/               # MIDI package
│   ├── writer.go       # MIDI file writing
│   ├── track.go        # Track management
│   ├── events.go       # MIDI event construction
│   ├── varlen.go       # Variable-length encoding
│   └── varlen_test.go  # Tests for encoding
└── music/              # Music generation package
    ├── generator.go    # Music generation logic
    └── errors.go       # Music errors
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./midi
```

## Roadmap & Extensions

Potential enhancements for future versions:

- **Advanced Mapping Algorithms**:
  - Map commit frequency to rhythm patterns
  - Use file changes to determine note duration
  - Map branch structure to musical phrases

- **Musical Features**:
  - Scale quantization (map notes to specific scales)
  - Chord generation from related commits
  - Dynamic tempo changes based on commit activity

- **Output Formats**:
  - Support for MIDI Format 2 (pattern-based)
  - Export to MusicXML
  - Generate audio files directly

- **Visualization**:
  - Generate sheet music
  - Create visual representations of the composition
  - Interactive playback controls

- **Repository Analysis**:
  - Filter commits by date range
  - Focus on specific branches
  - Weight commits by significance (merge commits, etc.)

- **Musical Styles**:
  - Preset configurations for different genres
  - Custom mapping functions
  - Pattern templates

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is open source. See LICENSE file for details.

## Acknowledgments

- Built with Go standard library only (no external MIDI dependencies)
- Follows MIDI 1.0 specification
- Inspired by the concept of data sonification

