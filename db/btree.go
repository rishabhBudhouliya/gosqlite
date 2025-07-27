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
		cell := pageContent[v:]
		payloadSize, startRowid := ProcessVarint(cell)
		rowId, startPayload := ProcessVarint(cell[startRowid:])
		payload := cell[startRowid+startPayload : (int64(startRowid+startPayload) + payloadSize)]
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
