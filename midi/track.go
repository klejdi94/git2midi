package midi

import "bytes"

// Track represents a MIDI track with events.
type Track struct {
	events []TrackEvent
}

// TrackEvent represents a single MIDI event with its delta time.
type TrackEvent struct {
	DeltaTime uint32
	Data      []byte
}

// NewTrack creates a new empty track.
func NewTrack() *Track {
	return &Track{
		events: make([]TrackEvent, 0),
	}
}

// AddEvent adds an event to the track with the specified delta time.
func (t *Track) AddEvent(deltaTime uint32, data []byte) {
	t.events = append(t.events, TrackEvent{
		DeltaTime: deltaTime,
		Data:      data,
	})
}

// AddNoteOn adds a Note On event with delta time.
func (t *Track) AddNoteOn(deltaTime uint32, channel, note, velocity byte) {
	t.AddEvent(deltaTime, NoteOn(channel, note, velocity))
}

// AddNoteOff adds a Note Off event with delta time.
func (t *Track) AddNoteOff(deltaTime uint32, channel, note, velocity byte) {
	t.AddEvent(deltaTime, NoteOff(channel, note, velocity))
}

// AddTempo adds a Set Tempo meta event with delta time.
func (t *Track) AddTempo(deltaTime uint32, tempo uint32) {
	t.AddEvent(deltaTime, SetTempo(tempo))
}

// AddEndOfTrack adds an End of Track event with delta time.
func (t *Track) AddEndOfTrack(deltaTime uint32) {
	t.AddEvent(deltaTime, EndOfTrack())
}

// Encode encodes the track into MIDI track chunk format.
func (t *Track) Encode() []byte {
	var buf bytes.Buffer

	for _, event := range t.events {
		buf.Write(EncodeVarLen(event.DeltaTime))
		buf.Write(event.Data)
	}

	trackData := buf.Bytes()
	trackLength := uint32(len(trackData))

	result := make([]byte, 8+len(trackData))
	copy(result[0:4], []byte("MTrk"))
	result[4] = byte(trackLength >> 24)
	result[5] = byte(trackLength >> 16)
	result[6] = byte(trackLength >> 8)
	result[7] = byte(trackLength)
	copy(result[8:], trackData)

	return result
}

// EventCount returns the number of events in the track.
func (t *Track) EventCount() int {
	return len(t.events)
}
