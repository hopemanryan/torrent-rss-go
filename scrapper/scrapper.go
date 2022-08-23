package rssScrapper

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/anacrolix/torrent"
	localDB "github.com/hopemanryan/torrent-rss/db"

	"github.com/gocolly/colly"
)

var limit = 2
var baseURL = "https://www.1377x.to"

type Scrapper struct {
	url           string
	browser       *colly.Collector
	torrectClient *torrent.Client
	allLinks      []string
}

func NewScrapper() *Scrapper {
	scrapInstnace := *colly.NewCollector()

	c, _ := torrent.NewClient(nil)

	scrapper := Scrapper{
		url:           fmt.Sprintf("%s/trending/w/tv/", baseURL),
		browser:       &scrapInstnace,
		torrectClient: c,
	}

	return &scrapper
}

func (s *Scrapper) AddListeners() {
	var count = 0
	s.browser.OnHTML(".featured-list", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")

		for _, link := range links {
			if strings.Contains(link, "/torrent") && count < limit {
				s.browser.Visit(fmt.Sprintf("%s/%s", baseURL, link))
			}
		}
	})

	s.browser.OnHTML(".torrentdown1", func(e *colly.HTMLElement) {
		magent := e.Attr("href")
		if count < limit {
			s.allLinks = append(s.allLinks, magent)
		}

	})

	s.browser.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
}
func (s *Scrapper) StartScrap(db *localDB.DB) {
	s.browser.Visit(s.url)

	defer s.torrectClient.Close()

	for _, link := range s.allLinks {
		t, _ := s.torrectClient.AddMagnet(link)
		<-t.GotInfo()
		info := cleanName(t.Info().Name)

		// check why file is downloading even though it returns true
		isDownloaded := db.CheckDownloadedName(info)
		if !isDownloaded {
			db.SaveFile(info, t.Info().Name)
			t.DownloadAll()

		} else {
			t.DisallowDataDownload()
		}

	}
	s.torrectClient.WaitAll()

	log.Print("Files Downloaded")

}

func cleanName(dirtyName string) string {
	re := regexp.MustCompile(`S(\d+)E(\d+)`)
	split := re.Split(dirtyName, -1)
	return split[0]
}
