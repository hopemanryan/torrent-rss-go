package scrappertorrentclient

import (
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
)

var defaultDownloadDir = "./download"

type TorrentClient struct {
	Client *torrent.Client
}

func NewTorrentScrapperClient() *TorrentClient {
	defaultConfig := torrent.NewDefaultClientConfig()
	defaultConfig.DataDir = defaultDownloadDir
	c, _ := torrent.NewClient(defaultConfig)
	return &TorrentClient{
		Client: c,
	}
}

func (t *TorrentClient) AddMagnet(magent string) {

	newMagent, _ := t.Client.AddMagnet(magent)
	<-newMagent.GotInfo()
	newMagent.DownloadAll()
	for newMagent.BytesCompleted() != newMagent.Info().TotalLength() {
		fmt.Printf("%d / %d \n", newMagent.BytesCompleted(), newMagent.Info().TotalLength())
		time.Sleep(time.Second * 5)
	}

}
