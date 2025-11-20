package midi

import "testing"

func TestEncodeVarLen(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected []byte
	}{
		{
			name:     "zero",
			input:    0,
			expected: []byte{0x00},
		},
		{
			name:     "single byte value",
			input:    127,
			expected: []byte{0x7F},
		},
		{
			name:     "two byte value",
			input:    128,
			expected: []byte{0x81, 0x00},
		},
		{
			name:     "two byte value max",
			input:    16383,
			expected: []byte{0xFF, 0x7F},
		},
		{
			name:     "three byte value",
			input:    16384,
			expected: []byte{0x81, 0x80, 0x00},
		},
		{
			name:     "common delta time",
			input:    480,
			expected: []byte{0x83, 0x60},
		},
		{
			name:     "large value",
			input:    2097151,
			expected: []byte{0xFF, 0xFF, 0x7F},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeVarLen(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("byte %d: got 0x%02X, want 0x%02X", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestDecodeVarLen(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expected     uint32
		expectedRead int
	}{
		{
			name:         "zero",
			input:        []byte{0x00},
			expected:     0,
			expectedRead: 1,
		},
		{
			name:         "single byte",
			input:        []byte{0x7F},
			expected:     127,
			expectedRead: 1,
		},
		{
			name:         "two byte",
			input:        []byte{0x81, 0x00},
			expected:     128,
			expectedRead: 2,
		},
		{
			name:         "two byte max",
			input:        []byte{0xFF, 0x7F},
			expected:     16383,
			expectedRead: 2,
		},
		{
			name:         "three byte",
			input:        []byte{0x81, 0x80, 0x00},
			expected:     16384,
			expectedRead: 3,
		},
		{
			name:         "with extra bytes",
			input:        []byte{0x7F, 0x42, 0x99},
			expected:     127,
			expectedRead: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, bytesRead := DecodeVarLen(tt.input)
			if result != tt.expected {
				t.Errorf("value: got %d, want %d", result, tt.expected)
			}
			if bytesRead != tt.expectedRead {
				t.Errorf("bytes read: got %d, want %d", bytesRead, tt.expectedRead)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	values := []uint32{0, 1, 127, 128, 255, 256, 16383, 16384, 100000, 2097151}

	for _, val := range values {
		t.Run("", func(t *testing.T) {
			encoded := EncodeVarLen(val)
			decoded, _ := DecodeVarLen(encoded)
			if decoded != val {
				t.Errorf("round trip failed: %d -> %v -> %d", val, encoded, decoded)
			}
		})
	}
}
