package main

import (
	"fmt"

	redisScrapper "github.com/hopemanryan/torrent-rss/redis"
	torrentScrapper "github.com/hopemanryan/torrent-rss/torrent"
)

func main() {

	client := torrentScrapper.NewTorrentScrapperClient()
	fmt.Printf("%v", client)

	redisInstance := redisScrapper.ConnectToRedis()
	fmt.Printf("%v", redisInstance)

	subscriber := redisInstance.Subscribe("send-user-data")
	for {
		msg, err := subscriber.ReceiveMessage()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v", msg)

	}

}
