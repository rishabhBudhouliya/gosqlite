package db

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Record []interface{}

func CreateRecord(payload []byte) Record {
	var record Record
	if len(payload) == 0 {
		fmt.Print("Header can't be empty!!")
		return record
	}
	headerSize, j := ProcessVarint(payload)

	hb, body := payload[j:headerSize], payload[headerSize:]

	for len(hb) > 0 {
		column, n := ProcessVarint(hb)
		if n < 0 {
			return nil
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
	return record
}
