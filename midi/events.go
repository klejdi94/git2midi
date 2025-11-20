package midi

// NoteOn creates a Note On event.
func NoteOn(channel, note, velocity byte) []byte {
	if channel > 15 {
		channel = 15
	}
	if note > 127 {
		note = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return []byte{0x90 | channel, note, velocity}
}

// NoteOff creates a Note Off event.
func NoteOff(channel, note, velocity byte) []byte {
	if channel > 15 {
		channel = 15
	}
	if note > 127 {
		note = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return []byte{0x80 | channel, note, velocity}
}

// SetTempo creates a Set Tempo meta event.
func SetTempo(tempo uint32) []byte {
	return []byte{
		0xFF, 0x51, 0x03,
		byte(tempo >> 16),
		byte(tempo >> 8),
		byte(tempo),
	}
}

// EndOfTrack creates an End of Track meta event.
func EndOfTrack() []byte {
	return []byte{0xFF, 0x2F, 0x00}
}

// BPMToMicrosecondsPerQuarter converts BPM to microseconds per quarter note.
func BPMToMicrosecondsPerQuarter(bpm int) uint32 {
	if bpm <= 0 {
		bpm = 120
	}
	return 60000000 / uint32(bpm)
}
