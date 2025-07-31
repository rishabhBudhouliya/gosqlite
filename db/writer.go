package db

import "fmt"

// the intention is the ability to insert a record if given a row id

// to be able to add a record
// step 1. retreive all available table names
// step 2. show their schema? or use their schema to create a potential record like
// (123, "honeydew", "golden brown")
// step 3. the function should take in the page number and record
// step 4. read page contents, determine if the free space is available, insert by
// step 5. create a cell
// step 5. in case of a split, a lot of things have to be considered
// step 6. new leaf and interior pages have to be created. also, may have to update the root.
// step 7. should we then go and update page 1 too?

func Insert(pageNum int, record Record) {
	//
	pager := Pager{}
	page, err := pager.GetPage(pageNum, 4096)
	if err != nil {
		fmt.Print(err)
	}
	// on receiving the page contents, it should start iterating and finding the leaf page
	btreePage := ParseBTreePage(page)
	// from the page, if it's an interior page, we need to traverse to leaf page

}
