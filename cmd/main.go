package main

import (
	rssScrapper "github.com/hopemanryan/torrent-rss/scrapper"
)

func main() {
	scrapper := *rssScrapper.NewScrapper()
	scrapper.AddListeners()
	scrapper.StartScrap()

}
