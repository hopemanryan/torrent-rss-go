package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron"
	localDB "github.com/hopemanryan/torrent-rss/db"
	rssScrapper "github.com/hopemanryan/torrent-rss/scrapper"
)

func main() {
	db := localDB.NewDb()
	defer db.Storage.Disconnect(db.CTX)
	s := gocron.NewScheduler(time.UTC)
	scrapper := *rssScrapper.NewScrapper()
	scrapper.AddListeners()
	scrapper.StartScrap(db)
	s.Every(1).Days().At("07:00").Do(func() {

		scrapper.StartScrap(db)
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

}

// mongo gets disconnected for some reason
