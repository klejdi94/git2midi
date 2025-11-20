package midi

// Variable-length encoding for MIDI delta times.
// MIDI uses a variable-length quantity where:
// - Each byte has 7 bits of data and 1 continuation bit
// - The MSB (bit 7) is 1 if more bytes follow, 0 for the last byte
// - Values are stored in big-endian order

// EncodeVarLen encodes a 32-bit value as a MIDI variable-length quantity.
// Returns the encoded bytes.
func EncodeVarLen(value uint32) []byte {
	if value == 0 {
		return []byte{0x00}
	}

	var result []byte
	var buffer [4]byte
	pos := 0

	// Build the value in reverse order
	for value > 0 {
		buffer[pos] = byte(value & 0x7F)
		value >>= 7
		pos++
	}

	// Reverse and set continuation bits
	for i := pos - 1; i >= 0; i-- {
		if i > 0 {
			buffer[i] |= 0x80 // Set continuation bit
		}
		result = append(result, buffer[i])
	}

	return result
}

// DecodeVarLen decodes a MIDI variable-length quantity from bytes.
// Returns the decoded value and the number of bytes consumed.
func DecodeVarLen(data []byte) (uint32, int) {
	var value uint32
	var bytesRead int

	for i, b := range data {
		bytesRead++
		value = (value << 7) | uint32(b&0x7F)
		if b&0x80 == 0 {
			break
		}
		if i >= 3 {
			// Prevent overflow - max 4 bytes for 32-bit value
			break
		}
	}

	return value, bytesRead
}
