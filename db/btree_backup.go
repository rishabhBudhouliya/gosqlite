package db

import "sort"

const (
	pageSize = 4096
)

type BTree struct {
	pager      Pager
	rootPageID uint32
	limit      int
}

func NewBTree(pager Pager) *BTree {
	return &BTree{
		pager: pager,
	}
}

func (bt *BTree) Search(rowID int64) {
	pageContent, _ := bt.pager.GetPage(bt.rootPageID, pageSize)
	var bh *BtreeHeader
	if bt.rootPageID == 1 {
		bh, _ = ParseBtreeHeader(pageContent[100:112])
	} else {
		bh, _ = ParseBtreeHeader(pageContent[:12])
	}
	var page *BPage
	if bh.Type == 0x0d {
		page = NewBtreeLeafPage(bh, pageContent[8:8+(2*bh.Cells)], pageContent)
	} else if bh.Type == 0x05 {
		page = NewBtreeInteriorPage(bh, pageContent[12:12+(2*bh.Cells)], pageContent)
	}
	// until page is not leaf
	for !page.isLeaf {
		index := sort.Search(len(page.cells), func(i int) bool {
			return page.cells[i].Rowid > rowID
		})

		if index >= len(page.cells)-1 {
			if page.bh.RightMostPointer != nil {
				pageContent, _ = bt.pager.GetPage(*bh.RightMostPointer, 4096)
				bh, _ = ParseBtreeHeader(pageContent[:12])
				if bh.Type == 0x0d {
					page = NewBtreeLeafPage(bh, pageContent[8:8+(2*bh.Cells)], pageContent)
				} else if bh.Type == 0x05 {
					page = NewBtreeInteriorPage(bh, pageContent[12:12+(2*bh.Cells)], pageContent)
				}
			}
		}
	}
}
