# Installation Guide

## Quick Install

The easiest way to install git2midi is using Go's install command:

```bash
go install github.com/klejdi94/git2midi@latest
```

This will install the `git2midi` binary to your `$GOPATH/bin` or `$HOME/go/bin` directory.

## Using as a Go Module

Add git2midi to your Go project:

```bash
go get github.com/klejdi94/git2midi@latest
```

Then import the packages you need:

```go
import (
    "github.com/klejdi94/git2midi/git"
    "github.com/klejdi94/git2midi/music"
    "github.com/klejdi94/git2midi/midi"
    "github.com/klejdi94/git2midi/config"
    "github.com/klejdi94/git2midi/audio"
)
```

## Available Packages

- `github.com/klejdi94/git2midi/git` - Git repository parsing
- `github.com/klejdi94/git2midi/music` - Music generation from commits
- `github.com/klejdi94/git2midi/midi` - MIDI file writing
- `github.com/klejdi94/git2midi/config` - Configuration management
- `github.com/klejdi94/git2midi/audio` - Audio format conversion

## Version Pinning

To use a specific version:

```bash
go get github.com/klejdi94/git2midi@v1.0.0
```

## Requirements

- Go 1.20 or later
- Git (for reading repositories)
- ffmpeg (optional, for audio conversion)

