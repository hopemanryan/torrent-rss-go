package scrappertorrentclient

import "github.com/anacrolix/torrent"

var defaultDownloadDir = "./download"

type TorrentClient struct {
	client *torrent.Client
}

func NewTorrentScrapperClient() *TorrentClient {
	defaultConfig := torrent.NewDefaultClientConfig()
	defaultConfig.DataDir = defaultDownloadDir
	c, _ := torrent.NewClient(defaultConfig)
	return &TorrentClient{
		client: c,
	}
}
