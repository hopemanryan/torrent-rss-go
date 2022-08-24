package main

import (
	"time"

	"github.com/go-co-op/gocron"
	localDB "github.com/hopemanryan/torrent-rss/db"
	rssScrapper "github.com/hopemanryan/torrent-rss/scrapper"
)

func main() {
	db := localDB.NewDb()
	s := gocron.NewScheduler(time.UTC)
	scrapper := *rssScrapper.NewScrapper()
	scrapper.AddListeners()
	scrapper.StartScrap(db)
	s.Every(1).Days().At("07:00").Do(func() {
		scrapper.StartScrap(db)
	})

}
