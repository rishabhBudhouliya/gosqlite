package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	// upper bound on size of a database page
	PageSize uint
}

type BtreeHeader struct {
	Type             uint
	Cells            uint16
	RightMostPointer *uint32 // 4 byte page number
}

func ParseBTreePage(pageContent []byte) interface{} {
	// determine first if it's the root page
	// second, determine if it's an interior or leaf page
	// iterate till you find the leaf page
	header, _ := parseHeader(pageContent)
	if header != nil {
		// we are at root page
		bHeader := pageContent[100:]
		bh, err := parseBtreeHeader(bHeader)
		if err != nil {
			fmt.Print(err)
			return nil
		}
		if bh.Type == 0x0d {
			return NewBtreeLeafPage(bHeader[8:8+(2*bh.Cells)], pageContent)
		} else if bh.Type == 0x05 {
			return NewBtreeInteriorPage(bHeader[12:12+(2*bh.Cells)], pageContent)
		}
	} else {
		// not the root page
		bh, err := parseBtreeHeader(pageContent)
		if err != nil {
			fmt.Print(err)
			return nil
		}
		if bh.Type == 0x0d {
			return NewBtreeLeafPage(pageContent[8:8+(2*bh.Cells)], pageContent)
		} else if bh.Type == 0x05 {
			return NewBtreeInteriorPage(pageContent[12:12+(2*bh.Cells)], pageContent)
		}
	}
	return nil
}

func Read(pageId int, pageContent []byte) {
	// special provisions for page 1
	h, err := parseHeader(pageContent)
	if err != nil {
		fmt.Print(err)
	}
	// for now just print the database header
	fmt.Println(h)
	bHeader := pageContent[100:]
	// now, let's read the btree header
	bh, err := parseBtreeHeader(bHeader)
	if err != nil {
		fmt.Print(err)
	}
	page := NewBtreeLeafPage(bHeader[8:8+(2*bh.Cells)], pageContent)
	// parseCellRecord(page.Cells[0])
	// 3901 3983 3779
	// fmt.Print(page.Offsets)
	fmt.Printf("Parsing page offset: %d", page.Offsets[0])
	fmt.Printf("Parsing cell: %d", page.Cells[0])
	parseCellRecord(page.Cells[0])

	// for _, c := range page.Cells {
	// 	parseCellRecord(c)
	// }

	// decision plane on what type of page to read
	switch t := bh.Type; t {
	case 0x02:
		return // interior index b-tree page
	case 0x05:
		return // interior table b tree page
	case 0x0a:
		return // leaf index b tree page
	case 0x0d:
		return // leaf table b tree page
	default:
		return
	}
}

func ParseCellRecord(cell Cell) Record {
	return CreateRecord(cell.Payload)
}

func ParseBtreeHeader(b []byte) (*BtreeHeader, error) {
	bh := struct {
		PageType         [1]byte
		StartFreeBlock   [2]byte
		Cells            [2]byte
		CellStartArea    [2]byte
		FragmentedFree   [1]byte
		RightMostPointer [4]byte // applicable for interior btree pages
	}{}
	err := binary.Read(bytes.NewBuffer(b), binary.BigEndian, &bh)
	rPointer := binary.BigEndian.Uint32(bh.RightMostPointer[:]) // can be nil
	bTreeHeader := BtreeHeader{
		Type:             uint(bh.PageType[0]),
		Cells:            binary.BigEndian.Uint16(bh.Cells[:]),
		RightMostPointer: &rPointer,
	}
	return &bTreeHeader, err
}

func ParseHeader(pageContent []byte) (*Header, error) {
	header := struct {
		Magic                [16]byte
		PageSize             uint16
		WriteVersion         uint8
		ReadVersion          uint8
		ReservedSpace        uint8
		MaxFraction          uint8
		MinFraction          uint8
		LeafFraction         uint8
		ChangeCounter        uint32
		_                    uint32
		_                    uint32
		_                    uint32
		SchemaCookie         uint32
		SchemaFormat         uint32
		_                    uint32
		_                    uint32
		TextEncoding         uint32
		_                    uint32
		_                    uint32
		_                    uint32
		ReservedForExpansion [20]byte
		_                    uint32
		_                    uint32
	}{}
	err := binary.Read(bytes.NewBuffer(pageContent), binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}
	h := Header{}
	h.PageSize = uint(header.PageSize)
	return &h, nil
}
