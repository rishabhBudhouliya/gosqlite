package db

import (
	"os"

	"golang.org/x/exp/mmap"
)

const (
	PageSize    = 4096
	PagerOpen   = 0
	PagerReader = 1
)

type Pager struct {
	// a file based lock?
	// ErrCode uint8
	// TODO: file handle
	f  *os.File
	mm *mmap.ReaderAt
	// TODO: page cache - we rely on mmap for caching pages
	// Cache map[int][]byte
}

func (p *Pager) Open(filename string) error {
	// open the file
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	mm, err := mmap.Open(filename)
	if err == nil {
		p.f = f
		p.mm = mm
	} else {
		f.Close()
		return err
	}
	return nil
}

func (p *Pager) GetPage(pageNum int, pageSize int) ([]byte, error) {
	// page number will start from 1 to N (sequential)
	// from page 1, will will get the database header and sqlite_schema table stored a btree leaf node

	// make page index 1-indexed
	pageOffset := (int64(pageNum-1) * int64(pageSize))
	buff := make([]byte, pageSize)
	// determine the size of byte[] to be loaded in memory
	// determine the offset
	_, err := p.mm.ReadAt(buff, pageOffset)
	return buff, err
}
