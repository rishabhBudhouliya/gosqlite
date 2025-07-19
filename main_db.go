package main

import (
	"github.com/rishabhBudhouliya/gosqlite/db"
)

func main() {
	pager := db.Pager{}
	pager.Open("sample.db")
	// for first page, get header and initiatlize a btree with page 1 at root
	result, _ := pager.GetPage(1, 4096)
	db.Read(1, result)
}
