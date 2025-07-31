package db

import (
	"encoding/binary"
)

// TODO: should we have varints type in go?
type Cell struct {
	LeftPageId   int32
	PayloadSize  int64 // varints
	Rowid        int64 // varints
	Payload      []byte
	PageOverflow int32
}

type BtreeLeafPage struct {
	// each key contains an (sorted by key) offset slice within the page that lets us know about the cell location
	Offsets []int
	Cells   []Cell
}

type BtreeInteriorPage struct {
	Offsets []int
	Cells   []Cell
}

/*
This function doesn't deserialize the payload
*/
func NewBtreeLeafPage(bh *BtreeHeader, cellPointer []byte, pageContent []byte) *BPage {
	offsets := parseCellPointers(cellPointer)
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
	return &BPage{
		bh:      bh,
		offsets: offsets,
		cells:   cells,
		isLeaf:  (bh.Type == 0x0d),
	}
}

func NewBtreeInteriorPage(bh *BtreeHeader, cellPointer []byte, pageContent []byte) *BPage {
	offsets := parseCellPointers(cellPointer)
	var cells []Cell
	for _, v := range offsets {
		cell := pageContent[v:]
		leftPageNumber, startRowId := ProcessVarint(cell)
		cell = cell[startRowId:]
		rowId, _ := ProcessVarint(cell)
		c := Cell{
			LeftPageId: int32(leftPageNumber),
			Rowid:      rowId,
		}
		cells = append(cells, c)
	}
	return &BPage{
		bh:      bh,
		offsets: offsets,
		cells:   cells,
		isLeaf:  (bh.Type == 0x0d),
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
