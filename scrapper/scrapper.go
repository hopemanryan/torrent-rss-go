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
var videoQuality = "1080p"

type Link struct {
	Name string
	Url  string
}
type Scrapper struct {
	Url      string
	Browser  *colly.Collector
	AllLinks []Link
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

		var url = e.Request.URL.Path
		rawFileName := getRawTvShowName(url)
		if rawFileName != "" {
			originalName := getCleanTvShowName(rawFileName)
			if originalName != "" {
				if count < limit {
					count = count + 1
					s.AllLinks = append(s.AllLinks, Link{
						Name: originalName,
						Url:  magent,
					})
				}
			}

		}

	})

}
func (s *Scrapper) StartScrap(db *localDB.DB) {
	err := s.Browser.Visit(s.Url)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, link := range s.AllLinks {
		re, _ := regexp.Compile(`S\d\dE\d\d`)
		seasonAndEpisode := re.FindString(link.Url)
		// check why file is downloading even though it returns true
		isDownloaded := db.CheckDownloadedName(link.Name, seasonAndEpisode)

		if !isDownloaded {

			clinet := redisScrapper.ConnectToRedis()
			ctx := context.Background()
			clinet.Publish(ctx, "new-magent", link.Url)
			db.SaveFile(link.Name, seasonAndEpisode)
		}
	}

}

func getRawTvShowName(url string) string {
	rawShowReg := regexp.MustCompile(`((//torrent/)\d*/)(.*)`)
	res := rawShowReg.FindStringSubmatch(url)
	if len(res) > 3 {
		return res[3]
	}
	return ""
}

func getCleanTvShowName(rawString string) string {
	originalFileReg := regexp.MustCompile(`(.*)-S(\d+)E(\d+)`)
	var orirginalFileGroups = originalFileReg.FindStringSubmatch(rawString)
	if len(orirginalFileGroups) >= 2 {
		return orirginalFileGroups[1]
	}
	return ""
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
