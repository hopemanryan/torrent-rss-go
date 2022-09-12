package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-co-op/gocron"
	localDB "github.com/hopemanryan/torrent-rss/db"
	rssScrapper "github.com/hopemanryan/torrent-rss/scrapper"
)

func main() {
	db := localDB.NewDb()
	defer func(Storage *mongo.Client, ctx context.Context) {
		err := Storage.Disconnect(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}(db.Storage, db.CTX)
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
