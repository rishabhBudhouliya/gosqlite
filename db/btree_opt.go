package db

import (
	"encoding/binary"
	"sort"
)

type BasePage struct {
	bh      *BtreeHeader
	offsets []uint16
	content []byte
}

type Page interface {
	CellCount() int
	GetRowID(i int) int64
}

type LeafPage struct {
	BasePage
}

type InteriorPage struct {
	BasePage
}

func NewPage(bh *BtreeHeader, cellPointer []byte, pageContent []byte) Page {
	base := BasePage{
		bh:      bh,
		offsets: parseCellPointers(cellPointer),
		content: pageContent,
	}
	switch bh.Type {
	case 0x0d:
		return &LeafPage{BasePage: base}
	case 0x05:
		return &InteriorPage{BasePage: base}
	default:
		return nil
	}
}

func (p *BasePage) CellCount() int {
	return len(p.offsets)
}

func (p *LeafPage) GetRowID(i int) int64 {
	if i < 0 || i >= len(p.offsets) {
		return 0
	}
	offset := p.offsets[i]
	if offset != 0 {
		cell := p.content[offset:]
		_, startRowid := ProcessVarint(cell)
		rowId, _ := ProcessVarint(cell[startRowid:])
		return rowId
	}
	return 0
}

func (p *InteriorPage) GetRowID(i int) int64 {
	if i < 0 || i >= len(p.offsets) {
		return 0
	}
	offset := p.offsets[i]
	cell := p.content[offset:]
	cell = cell[4:]
	rowId, _ := ProcessVarint(cell)
	return rowId
}

func (p *InteriorPage) GetLeftmostPageID(i int) uint32 {
	if i < 0 || i >= len(p.offsets) {
		return 0
	}
	offset := p.offsets[i]
	cell := p.content[offset:]
	leftPageNumber := binary.BigEndian.Uint32(cell[:4])
	return leftPageNumber
}

// for the given index of offset, it will try to find the cell and parse the record
func (p *LeafPage) GetRecord(i int) Record {
	if i < 0 || i >= len(p.offsets) {
		return nil
	}
	offset := p.offsets[i]
	cell := p.content[offset:]
	payloadSize, startRowid := ProcessVarint(cell)
	_, startPayload := ProcessVarint(cell[startRowid:])
	payload := cell[startRowid+startPayload : (int64(startRowid+startPayload) + payloadSize)]
	return CreateRecord(payload)
}

func rSearch(pager *Pager, p Page, rowID int64) Record {
	switch page := p.(type) {
	case *LeafPage:
		index := sort.Search(page.CellCount(), func(i int) bool {
			return page.GetRowID(i) >= rowID
		})
		if index < page.CellCount() {
			if page.GetRowID(index) == rowID {
				return page.GetRecord(index)
			}
		}
		// didn't find the record on this page
		return nil
	case *InteriorPage:
		index := sort.Search(page.CellCount(), func(i int) bool {
			return page.GetRowID(i) >= rowID
		})
		if index < page.CellCount()-1 {
			pageContent, _ := pager.GetPage(page.GetLeftmostPageID(index), 4096)
			bh, _ := ParseBtreeHeader(pageContent[:12])
			if bh.Type == 0x0d {
				return rSearch(pager, NewPage(bh, pageContent[8:8+(2*bh.Cells)], pageContent), rowID)
			} else if bh.Type == 0x05 {
				return rSearch(pager, NewPage(bh, pageContent[12:12+(2*bh.Cells)], pageContent), rowID)
			}
		} else if index >= page.CellCount()-1 {
			if page.bh.RightMostPointer != nil {
				pageContent, _ := pager.GetPage(*page.bh.RightMostPointer, 4096)
				bh, _ := ParseBtreeHeader(pageContent[:12])
				if bh.Type == 0x0d {
					return rSearch(pager, NewPage(bh, pageContent[8:8+(2*bh.Cells)], pageContent), rowID)
				} else if bh.Type == 0x05 {
					return rSearch(pager, NewPage(bh, pageContent[12:12+(2*bh.Cells)], pageContent), rowID)
				}
			}
		}
	}
	return nil
}

func parseCellPointers(cp []byte) []uint16 {
	n := len(cp)
	result := make([]uint16, n/2)
	for i := range n / 2 {
		result[i] = binary.BigEndian.Uint16(cp[2*i : 2*i+2])
	}
	return result
}

type BtreeHeader struct {
	Type             uint
	Cells            uint16
	RightMostPointer *uint32 // 4 byte page number
}
