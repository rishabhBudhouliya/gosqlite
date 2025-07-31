package db

import "fmt"

// TODO: do we need the reference of root page?
type database struct {
	pager  Pager
	header *Header
}

// when we create the new database in memory instance, it should also contain
// the master table?
// this gives the user the access to virtual memory that can get any page from the
// binary file
func newDatabase(fileName string) *database {
	pager := Pager{}
	pager.Open(fileName)
	// get root page to get
	// only get the first 100 bytes to read the header
	rootHeader, _ := pager.GetPage(1, 100)
	h, _ := ParseHeader(rootHeader)
	return &database{
		pager:  pager,
		header: h,
	}
}

// a btree is a collection of pages
type BPage struct {
	bh      *BtreeHeader
	offsets []int
	cells   []Cell
	isLeaf  bool
}

type sqliteSchema struct {
	internalType string
	name         string
	tableName    string
	rootPage     int
	sql          string
}

// what is the abstraction of a table in our codebase? a btree reference.
func (db *database) createSchemaTable() {
	// for schema table
	// page offset is predetermined and page size is always predetermined atp
	pageContent, err := db.pager.GetPage(1, int(db.header.PageSize))
	if err != nil {
		fmt.Errorf("Unable to fetch page for schema table: ", err)
	}
	// we've already read the first 100 bytes at this point
	// at max 12 bytes in case it's an interior page
	bh, err := ParseBtreeHeader(pageContent[100:112])
	if err != nil {
		fmt.Errorf("Unable to parse btree header: ", err)
	}
	var schemaPage *BPage
	if bh.Type == 0x0d {
		schemaPage = NewBtreeLeafPage(bh, pageContent[8:8+(2*bh.Cells)], pageContent)
	} else if bh.Type == 0x05 {
		schemaPage = NewBtreeInteriorPage(bh, pageContent[12:12+(2*bh.Cells)], pageContent)
	}
	PagetoTable(schemaPage)
}

func PageToTable(bPage *BPage) {
	if bPage == nil {
		fmt.Print("can't convert a nil page")
	}
	if bPage.isLeaf {
		// we can read it as it
	}
}

// func (db *database) CreateTableContent(tableName string) {
// 	//get page for table content
// 	db.pager.GetPage()
// }
