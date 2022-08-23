package main

import (
	localDB "github.com/hopemanryan/torrent-rss/db"
	rssScrapper "github.com/hopemanryan/torrent-rss/scrapper"
)

func main() {
	db := localDB.NewDb()
	scrapper := *rssScrapper.NewScrapper()
	scrapper.AddListeners()
	scrapper.StartScrap(db)

}
