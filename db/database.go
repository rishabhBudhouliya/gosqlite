package db

import "sort"

// TODO: do we need the reference of root page?
type Database struct {
	Pager  Pager
	Header *Header
}

const (
	pageSize = 4096
)

// when we create the new database in memory instance, it should also contain
// the master table?
// this gives the user the access to virtual memory that can get any page from the
// binary file
func NewDatabase(fileName string) *Database {
	pager := Pager{}
	pager.Open(fileName)
	// get root page to get
	// only get the first 100 bytes to read the header
	rootHeader, _ := pager.GetPage(1, 100)
	h, _ := ParseHeader(rootHeader)
	return &Database{
		Pager:  pager,
		Header: h,
	}
}

func Search(rootPageID uint32, pager *Pager, rowID int64) Record {
	currentPagerID := rootPageID
	for {
		p := fetchAndParsePage(pager, currentPagerID, pageSize)
		var record Record
		found := false
		switch page := p.(type) {
		case *LeafPage:
			found = true
			index := sort.Search(page.CellCount(), func(i int) bool {
				return page.GetRowID(i) >= rowID
			})
			if index < page.CellCount() {
				if page.GetRowID(index) == rowID {
					return page.GetRecord(index)
				}
			}
			return nil
		case *InteriorPage:
			index := sort.Search(page.CellCount(), func(i int) bool {
				return page.GetRowID(i) >= rowID
			})
			if index < page.CellCount()-1 {
				currentPagerID = page.GetLeftmostPageID(index)
			} else if index >= page.CellCount()-1 {
				if page.bh.RightMostPointer != nil {
					currentPagerID = *page.bh.RightMostPointer
				} else {
					found = true
				}
			}
		}
		if found {
			return record
		}
	}
}

func fetchAndParsePage(pager *Pager, pageID uint32, pageSize int) Page {
	pageContent, _ := pager.GetPage(pageID, pageSize)
	var bh *BtreeHeader
	btHeader := pageContent
	if pageID == 1 {
		bh, _ = ParseBtreeHeader(btHeader[100:112])
		if bh.Type == 0x0d {
			btHeader = btHeader[108:]
		} else {
			btHeader = btHeader[112:]
		}
	} else {
		bh, _ = ParseBtreeHeader(pageContent[:12])
		if bh.Type == 0x0d {
			btHeader = btHeader[8:]
		} else {
			btHeader = btHeader[12:]
		}
	}
	return NewPage(bh, btHeader[:(2*bh.Cells)], pageContent)
}

type SqliteSchema struct {
	InternalType string
	Name         string
	TableName    string
	RootPage     int
	Sql          string
}

// what is the abstraction of a table in our codebase? a btree reference.
// func (db *database) createSchemaTable() {
// 	// for schema table
// 	// page offset is predetermined and page size is always predetermined atp
// 	pageContent, err := db.pager.GetPage(1, int(db.header.PageSize))
// 	if err != nil {
// 		fmt.Errorf("Unable to fetch page for schema table: ", err)
// 	}
// 	// we've already read the first 100 bytes at this point
// 	// at max 12 bytes in case it's an interior page
// 	bh, err := ParseBtreeHeader(pageContent[100:112])
// 	if err != nil {
// 		fmt.Errorf("Unable to parse btree header: ", err)
// 	}
// 	var schemaPage *BPage
// 	if bh.Type == 0x0d {
// 		schemaPage = NewBtreeLeafPage(bh, pageContent[8:8+(2*bh.Cells)], pageContent)
// 	} else if bh.Type == 0x05 {
// 		schemaPage = NewBtreeInteriorPage(bh, pageContent[12:12+(2*bh.Cells)], pageContent)
// 	}
// 	// PagetoTable(schemaPage)
// }

// func (db *database) CreateTableContent(tableName string) {
// 	//get page for table content
// 	db.pager.GetPage()
// }
