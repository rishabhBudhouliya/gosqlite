package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

type Header struct {
	// upper bound on size of a database page
	PageSize uint
}

type BtreeHeader struct {
	Type  uint
	Cells uint16
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

type Record struct {
	Header []byte
	Body   []byte
}

func (r *Record) ParseRecord(payload []byte) {
	if len(payload) == 0 {
		fmt.Print("Header can't be empty!!")
		return
	}
	// payload has to be divided into two parts

	// header
	// body
	// reading the first varint will give us the boundary of header and body
	headerSize, j := ProcessVarint(payload)

	var record []interface{}

	hb, body := payload[j:headerSize], payload[headerSize:]

	for len(hb) > 0 {
		column, n := ProcessVarint(hb)
		if n < 0 {
			return
		}
		hb = hb[n:]
		switch column {
		case 0:
			record = append(record, nil)
		case 1:
			record = append(record, int64(int8(body[0])))
			body = body[1:]
		case 2:
			record = append(record, int64(binary.BigEndian.Uint16(body[:2])))
			body = body[2:]
		case 3:
			record = append(record, ReadTwos24Bit(body[:3]))
			body = body[3:]
		case 4:
			record = append(record, int64(binary.BigEndian.Uint32(body[:4])))
			body = body[4:]
		case 5:
			record = append(record, ReadTwos48Bit(body[:6]))
			body = body[6:]
		case 6:
			record = append(record, binary.BigEndian.Uint64(body[:8]))
			body = body[8:]
		case 7:
			record = append(record, math.Float64frombits(binary.BigEndian.Uint64(body[:8])))
			body = body[8:]
		case 8:
			record = append(record, int64(0))
		case 9:
			record = append(record, int64(1))
		case 10, 11:
			fmt.Print("reserved for internal use")
		default:
			if column >= 12 {
				n := (column - 12) / 2
				record = append(record, body[:n])
				body = body[n:]
			} else if column >= 13 {
				n := (column - 13) / 2
				data := string(body[:n])
				record = append(record, data)
				body = body[n:]
			}
		}
	}
	for _, v := range record {
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Int64 {
			fmt.Println(v)
		} else {
			fmt.Println(string(v.([]byte)))
		}
	}
	// fmt.Printf("here's the record: %v ", record)
}

func parseCellRecord(cell Cell) {
	r := Record{}
	fmt.Printf("\nI am parsing cell with row id: %d\n", cell.Rowid)
	r.ParseRecord(cell.Payload)
}

func parseBtreeHeader(b []byte) (*BtreeHeader, error) {
	bh := struct {
		PageType         [1]byte
		StartFreeBlock   [2]byte
		Cells            [2]byte
		CellStartArea    [2]byte
		FragmentedFree   [1]byte
		RightMostPointer [4]byte // applicable for interior btree pages
	}{}
	err := binary.Read(bytes.NewBuffer(b), binary.BigEndian, &bh)
	bTreeHeader := BtreeHeader{
		Type:  uint(bh.PageType[0]),
		Cells: binary.BigEndian.Uint16(bh.Cells[:]),
	}
	return &bTreeHeader, err
}

func parseHeader(pageContent []byte) (*Header, error) {
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
