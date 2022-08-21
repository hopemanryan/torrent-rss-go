package rssScrapper

import (
	"fmt"
	"log"
	"strings"

	"github.com/anacrolix/torrent"

	"github.com/gocolly/colly"
)

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

// func addListeners(collyInstance *colly.Collector, tc *torrent.Client) {

// 	collyInstance.OnHTML(".featured-list", func(e *colly.HTMLElement) {
// 		links := e.ChildAttrs("a", "href")

// 		for _, link := range links {
// 			if strings.Contains(link, "/torrent") {
// 				collyInstance.Visit(fmt.Sprintf("%s/%s", baseURL, link))
// 			}
// 		}
// 	})

// 	collyInstance.OnHTML(".torrentdown1", func(e *colly.HTMLElement) {
// 		magent := e.Attr("href")
// 		t, _ := tc.AddMagnet(magent)
// 		t.DownloadAll()

// 	})

// 	collyInstance.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL)
// 	})
// }

func (s *Scrapper) AddListeners() {
	var count = 0
	s.browser.OnHTML(".featured-list", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")

		for _, link := range links {
			if strings.Contains(link, "/torrent") {
				s.browser.Visit(fmt.Sprintf("%s/%s", baseURL, link))
			}
		}
	})

	s.browser.OnHTML(".torrentdown1", func(e *colly.HTMLElement) {
		magent := e.Attr("href")
		if count < 2 {
			count = count + 1
			s.allLinks = append(s.allLinks, magent)
		}

	})

	s.browser.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
}
func (s *Scrapper) StartScrap() {
	s.browser.Visit(s.url)

	defer s.torrectClient.Close()

	for _, link := range s.allLinks {
		t, _ := s.torrectClient.AddMagnet(link)
		<-t.GotInfo()
		t.DownloadAll()
	}
	s.torrectClient.WaitAll()
	log.Print("Files Downloaded")

}
