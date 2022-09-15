package main

import (
	"context"
	"fmt"
	torrentScrapper "github.com/hopemanryan/torrent-rss/torrent"
	"net/http"

	redisScrapper "github.com/hopemanryan/torrent-rss/redis"
)

var ctx = context.Background()

func main() {

	client := torrentScrapper.NewTorrentScrapperClient()

	redisInstance := redisScrapper.ConnectToRedis()
	fmt.Printf("%v", redisInstance)
	println("Redis connected")
	subscriber := redisInstance.Subscribe(ctx, "new-magent")

	ch := subscriber.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		go client.AddMagnet(msg.Payload)
	}

	http.ListenAndServe(":8090", nil)

}
