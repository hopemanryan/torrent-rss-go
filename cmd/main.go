package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	scrapper := *rssScrapper.NewScrapper()
	scrapper.AddListeners()
	s := gocron.NewScheduler(time.UTC)
	fmt.Printf("%v", s)

	s.Every(1).Days().At("07:00").Do(func() {

		scrapper.StartScrap(db)
	})

	fmt.Printf("%v", s.Jobs())
	s.StartAsync()
	http.ListenAndServe(":8090", nil)

}
