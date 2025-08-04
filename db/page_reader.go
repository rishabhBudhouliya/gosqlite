package db

import (
	"bytes"
	"encoding/binary"
)

type Header struct {
	// upper bound on size of a database page
	PageSize uint
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
