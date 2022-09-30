package rssScrapper

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	localDB "github.com/hopemanryan/torrent-rss/db"
	redisScrapper "github.com/hopemanryan/torrent-rss/redis"
)

var baseURL = "https://www.1377x.to"
var skipQuality = []string{"720p", "WEBRip", "WEB-x264"}

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

	scrapInstance := *colly.NewCollector(colly.AllowURLRevisit())

	scrapper := Scrapper{
		Url:     fmt.Sprintf("%s/trending/w/tv/", baseURL),
		Browser: &scrapInstance,
	}

	return &scrapper
}

func (s *Scrapper) AddListeners() {
	s.Browser.OnHTML(".featured-list", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")

		for _, link := range links {
			if strings.Contains(link, "/torrent") {
				if checkUrlContains(link) {
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
				s.AllLinks = append(s.AllLinks, Link{
					Name: originalName,
					Url:  magent,
				})
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

func checkUrlContains(link string) bool {
	var isValid = true
	for _, quality := range skipQuality {
		if strings.Contains(link, quality) {
			isValid = false
		}
	}

	return isValid
}
