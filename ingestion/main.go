package main

import (
	"github.com/margostino/metabook/db"
	"github.com/margostino/metabook/scraper"
)

func main() {
	document := scraper.Collect()
	db.Connect()
	db.Index(document)
}
