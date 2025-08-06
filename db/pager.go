package db

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
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
	f    *os.File
	data []byte
}

func (p *Pager) Open(filename string) error {
	// open the file
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	fileSize := stat.Size()

	data, err := unix.Mmap(int(f.Fd()), 0, int(fileSize), unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		f.Close()
		return err
	}
	if err := unix.Madvise(data, unix.MADV_SEQUENTIAL); err != nil {
		fmt.Printf("Madvise warning: %v", err)
	}
	p.f = f
	p.data = data
	return nil
}

func (p *Pager) GetPage(pageNum uint32, pageSize int) ([]byte, error) {
	start := int64(pageNum-1) * int64(pageSize)
	end := start + int64(pageSize)
	if start < 0 || end > int64(len(p.data)) {
		return nil, io.EOF
	}
	pageSlice := p.data[start:end]
	return pageSlice, nil
}

func (p *Pager) Close() error {
	if err := unix.Munmap(p.data); err != nil {
		p.f.Close()
		return err
	}
	return p.f.Close()
}
