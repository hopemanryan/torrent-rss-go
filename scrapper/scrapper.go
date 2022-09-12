package rssScrapper

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	localDB "github.com/hopemanryan/torrent-rss/db"
	redisScrapper "github.com/hopemanryan/torrent-rss/redis"
)

var limit = 20
var TorrentLimitToken = "TORRENT_LIMIT"
var VideoQualityToken = "QUIALITY"
var baseURL = "https://www.1377x.to"
var defaultDownloadDir = "./download"
var videoQuality = "1080p"

type Scrapper struct {
	Url      string
	Browser  *colly.Collector
	AllLinks []string
}

func NewScrapper() *Scrapper {

	getLimitFromEnv()
	getVideoQualityFromEnv()

	println(limit)

	scrapInstance := *colly.NewCollector()

	scrapper := Scrapper{
		Url:     fmt.Sprintf("%s/trending/w/tv/", baseURL),
		Browser: &scrapInstance,
	}

	return &scrapper
}

func (s *Scrapper) AddListeners() {
	var count = 0
	s.Browser.OnHTML(".featured-list", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")

		for _, link := range links {
			if strings.Contains(link, "/torrent") {
				if strings.Contains(link, videoQuality) {
					log.Printf("Visit webpage: %s ", link)
					err := s.Browser.Visit(fmt.Sprintf("%s/%s", baseURL, link))

					if err != nil {
						log.Printf("%s , url could not open", link)
					}
				}
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
	err := s.Browser.Visit(s.Url)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, link := range s.AllLinks {
		info := cleanName(link)
		re, _ := regexp.Compile(`S\d\dE\d\d`)
		seasonAndEpisode := re.FindString(link)
		// check why file is downloading even though it returns true
		isDownloaded := db.CheckDownloadedName(info, seasonAndEpisode)

		if !isDownloaded {

			clinet := redisScrapper.ConnectToRedis()
			ctx := context.Background()
			clinet.Publish(ctx, "new-magent", link)
			db.SaveFile(info, link, seasonAndEpisode)
		}
	}

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
	videoQuality = os.Getenv(VideoQualityToken)
}
