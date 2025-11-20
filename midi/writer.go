package midi

import (
	"fmt"
	"os"
)

// Writer handles writing MIDI files.
type Writer struct {
	format   uint16
	division uint16
	tracks   []*Track
}

// NewWriter creates a new MIDI writer.
func NewWriter(format uint16, division uint16) *Writer {
	return &Writer{
		format:   format,
		division: division,
		tracks:   make([]*Track, 0),
	}
}

// AddTrack adds a track to the writer.
func (w *Writer) AddTrack(track *Track) {
	w.tracks = append(w.tracks, track)
}

// WriteFile writes the MIDI file to the specified path.
func (w *Writer) WriteFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	header := w.encodeHeader()
	if _, err := file.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, track := range w.tracks {
		trackData := track.Encode()
		if _, err := file.Write(trackData); err != nil {
			return fmt.Errorf("failed to write track: %w", err)
		}
	}

	return nil
}

func (w *Writer) encodeHeader() []byte {
	header := make([]byte, 14)
	copy(header[0:4], []byte("MThd"))
	header[4] = 0x00
	header[5] = 0x00
	header[6] = 0x00
	header[7] = 0x06
	header[8] = byte(w.format >> 8)
	header[9] = byte(w.format)
	numTracks := uint16(len(w.tracks))
	header[10] = byte(numTracks >> 8)
	header[11] = byte(numTracks)
	header[12] = byte(w.division >> 8)
	header[13] = byte(w.division)

	return header
}

// GetDivision returns the ticks per quarter note.
func (w *Writer) GetDivision() uint16 {
	return w.division
}

// TrackCount returns the number of tracks in the writer.
func (w *Writer) TrackCount() int {
	return len(w.tracks)
}
