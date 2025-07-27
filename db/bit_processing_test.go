package db

import (
	"fmt"
	"testing"
)

func TestProcessVarint(t *testing.T) {
	testTable := []struct {
		name          string
		input         []byte
		expectedValue int64
		expectedBytes int
	}{
		{
			name:          "Single byte, value 0",
			input:         []byte{0x00},
			expectedValue: 0,
			expectedBytes: 1,
		},
		{
			name:          "Two bytes, value 128",
			input:         []byte{0x81, 0x00},
			expectedValue: 128,
			expectedBytes: 2,
		},
		{
			name:          "Three bytes, value 128",
			input:         []byte{0x81, 0x01, 0x01},
			expectedValue: 129,
			expectedBytes: 3,
		},
		{
			name:          "Two bytes, value 240",
			input:         []byte{0x81, 0x70},
			expectedValue: 240,
			expectedBytes: 2,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			actualValue, actualBytes := ProcessVarint(tc.input)

			if actualValue == tc.expectedValue {
				fmt.Print("Value parsed is correct!")
			} else {
				t.Errorf("Error in parsing var int, expected: %d vs actual: %d ", tc.expectedValue, actualValue)
			}

			if actualBytes == tc.expectedBytes {
				fmt.Print("Number of bytes parsed correctly!")
			} else {
				fmt.Printf("Error, expected bytes: %d", tc.expectedBytes)
			}
		})
	}
}
