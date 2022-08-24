package rssScrapper

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/gocolly/colly"
	localDB "github.com/hopemanryan/torrent-rss/db"
)

var limit = 20
var TorrentLimitToken = "TORRENT_LIMIT"
var VideoQualityToken = "QUIALITY"
var baseURL = "https://www.1377x.to"
var defaultDownloadDir = "./download"
var videoQueality string

type Scrapper struct {
	Url           string
	Browser       *colly.Collector
	TorrentClient *torrent.Client
	AllLinks      []string
}

func NewScrapper() *Scrapper {

	getLimitFromEnv()
	getVideoQualityFromEnv()

	println(limit)

	scrapInstance := *colly.NewCollector()

	defaultConfig := torrent.NewDefaultClientConfig()
	defaultConfig.DataDir = defaultDownloadDir
	c, _ := torrent.NewClient(defaultConfig)

	scrapper := Scrapper{
		Url:           fmt.Sprintf("%s/trending/w/tv/", baseURL),
		Browser:       &scrapInstance,
		TorrentClient: c,
	}

	return &scrapper
}

func (s *Scrapper) AddListeners() {
	var count = 0
	s.Browser.OnHTML(".featured-list", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")

		for _, link := range links {
			if strings.Contains(link, "/torrent") {
				if videoQueality != "" && !strings.Contains(link, videoQueality) {
					return
				}
				s.Browser.Visit(fmt.Sprintf("%s/%s", baseURL, link))
			}
		}
	})

	s.Browser.OnHTML(".torrentdown1", func(e *colly.HTMLElement) {
		magent := e.Attr("href")
		if count < limit {
			count = count + 1
			s.AllLinks = append(s.AllLinks, magent)
		}

	})

}
func (s *Scrapper) StartScrap(db *localDB.DB) {
	s.Browser.Visit(s.Url)
	defer s.TorrentClient.Close()

	for _, link := range s.AllLinks {
		print(link)
		t, _ := s.TorrentClient.AddMagnet(link)
		<-t.GotInfo()
		info := cleanName(t.Info().Name)

		// check why file is downloading even though it returns true
		isDownloaded := db.CheckDownloadedName(info)
		if !isDownloaded {
			db.SaveFile(info, t.Info().Name)
			t.DownloadAll()

			fmt.Printf("Total Length: %d", t.Info().TotalLength())
			for t.BytesCompleted() != t.Info().TotalLength() {
				fmt.Printf("%d / %d \n", t.BytesCompleted(), t.Info().TotalLength())
				time.Sleep(time.Second * 5)

			}

		}

	}
	s.TorrentClient.WaitAll()

	log.Print("Files Downloaded")

}

func cleanName(dirtyName string) string {
	re := regexp.MustCompile(`S(\d+)E(\d+)`)
	split := re.Split(dirtyName, -1)
	return split[0]
}

func getLimitFromEnv() {
	osLimit := os.Getenv(TorrentLimitToken)

	if osLimit != "" {
		i, err := strconv.Atoi(osLimit)
		if err != nil {
			limit = i
		}
	}
}

func getVideoQualityFromEnv() {
	videoQueality = os.Getenv(VideoQualityToken)
}
