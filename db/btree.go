package db

import (
	"encoding/binary"
)

// TODO: should we have varints type in go?
type Cell struct {
	PayloadSize  int64 // varints
	Rowid        int64 // varints
	Payload      []byte
	PageOverflow int32
}

type BtreeLeafPage struct {
	// each key contains an (sorted by key) offset within the page that lets us know about the cell location
	Offsets []int
	Cells   []Cell
}

func NewBtreeLeafPage(cellPointer []byte, pageContent []byte) *BtreeLeafPage {
	// pointer array with offsets, those offsets give you cells
	// 2*i, 2*i + 2
	offsets := parseCellPointers(cellPointer)
	// fmt.Print(offsets)
	var cells []Cell
	for _, v := range offsets {
		payloadSize, startRowid := processVarint(pageContent[v:])
		rowId, startPayload := processVarint(pageContent[startRowid:])
		payload := pageContent[startPayload:payloadSize]
		c := Cell{
			PayloadSize: payloadSize,
			Rowid:       rowId,
			Payload:     payload,
		}
		cells = append(cells, c)
	}
	return &BtreeLeafPage{
		Offsets: offsets,
		Cells:   cells,
	}
}

// 3901 3983 3779
func parseCellPointers(cp []byte) []int {
	n := len(cp)
	var result []int
	for i := range n / 2 {
		start := binary.BigEndian.Uint16(cp[2*i : 2*i+2])
		result = append(result, int(start))
	}
	return result
}

// what's a varint?
// a base 128 integer.
// 128 = 128 ^ 1 + 128 ^ 0
// how is 80 represented?
// 80 = 128 ^ 0 + payloadBits = 0 + 80 = 80

func processVarint(b []byte) (int64, int) {
	var x int64
	for i := range b {
		if i < 8 {
			// keep adding the last 7 bits of the current byte unless the MSB is 1
			x = (x << 7) | int64(b[i]&0x7f)
			if x&0x80 == 0 {
				return x, i + 1
			}
		}
		// 9th bit reached, take it as is
		if i == 8 {
			x = (x << 8) | int64(b[i])
			return x, i + 1
		}
	}
	// should not reach here
	return 0, -1
}
