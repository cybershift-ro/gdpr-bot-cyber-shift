package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"
)

var (
	WATCHDOG_URL            = "https://www.dataprotection.ro/?page=allnews"
	WATCHDOG_DEBUG          = false
	WATCHDOG_CHECK_INTERVAL = 1800 // 30 minutes
)

var keywords = []string{
	"încălcarea RGPD",
	"Amendă pentru",
	"Sancţiune pentru",
	"Sancţiune",
	"sancţiune",
	"încălcarea",
	"Amendă",
	"amendă",
	"amenda",
}

var compiledKeyWords = []*regexp.Regexp{}

var feedMutex sync.Mutex

var sanctionFeed = &feeds.Feed{
	Title:       "Cyber Shift GDPR Bot",
	Link:        &feeds.Link{Href: "https://gdpr.cybershift.dev"},
	Description: "O listă cu toate amenzile anunțate de către Autoritatea Naţională de Supraveghere a Prelucrării Datelor cu Caracter Personal.",
	Author:      &feeds.Author{Name: "Cyber Shift", Email: "office@cybershift.dev"},
	Created:     time.Now(),
}

var lastSanctionsFound = 0

func compileAllKeywords() {
	for _, triggers := range keywords {
		compiled := regexp.MustCompile(triggers)
		compiledKeyWords = append(compiledKeyWords, compiled)
	}
}

func isSaction(title string) bool {

	for _, keyword := range compiledKeyWords {
		if keyword.MatchString(title) {
			return true
		}
	}

	return false
}

func emptyFeedItem() *feeds.Item {
	return &feeds.Item{
		Title:       "Default Title",
		Link:        &feeds.Link{Href: "http://cybershift.dev"},
		Description: "Cyber Shift GDPR Bot",
		Author:      &feeds.Author{Name: "Cyber Shift", Email: "office@cybershift.dev"},
		Created:     time.Unix(0, 0),
	}
}

func scrape() {
	// Scrapper init
	c := colly.NewCollector(
		colly.AllowedDomains("www.dataprotection.ro", "dataprotection.ro"),
	)

	totalSanctions := 0

	var sanctionsFound []*feeds.Item

	c.OnHTML("div[id=rectangle_scroll]", func(element *colly.HTMLElement) {

		// Start with an empty feed item and fill it as we go
		tempItem := emptyFeedItem()

		foundASanction := false

		element.ForEach("p", func(index int, e *colly.HTMLElement) {
			if len(e.Text) == 0 {
				return
			}

			// Skip <p>'s that only have whitespaces
			if len(strings.TrimSpace(e.Text)) == 0 {
				return
			}

			// Check if it's a link
			newsURL := e.ChildAttr("a", "href")
			if len(newsURL) != 0 {
				absUrl := "https://" + "www.dataprotection.ro/" + newsURL

				// Ignore normal news articles
				if !foundASanction {
					return
				}

				tempItem.Link = &feeds.Link{Href: absUrl}

				// Reset flag
				foundASanction = false

				// Push the item on the feed
				sanctionFeed.Items = append(sanctionFeed.Items, tempItem)

				// Reset temp item
				tempItem = emptyFeedItem()

				return
			}

			// Dates are in format 27/08/2021
			if len(e.Text) == 10 {
				stamp, err := time.Parse("02/01/2006", e.Text)

				if err != nil {
					return
				}

				tempItem.Created = stamp
				return
			}

			// Else the <p> it's a title
			if isSaction(e.Text) {
				tempItem.Title = e.Text

				foundASanction = true

				totalSanctions += 1
			}

		})
	})

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("Looking for new sanctions on: ", req.URL.String())
	})

	c.Visit(WATCHDOG_URL)

	if lastSanctionsFound < len(sanctionsFound) {
		fmt.Println("Updating the news feed")

		feedMutex.Lock()
		sanctionFeed.Items = sanctionsFound
		feedMutex.Unlock()
	}

	fmt.Printf("Found %d total sanctions!", totalSanctions)
}

func main() {
	// Make sure to compile all keywords
	compileAllKeywords()

	router := gin.Default()

	router.GET("/rss", serveRSSFeed)
	router.GET("/json", serveJSONFeed)
	router.GET("/atom", serveAtomFeed)

	// Start scapper in a separate routine
	go webScrapper()

	router.Run(":3001")
}

func webScrapper() {
	// Timer
	dataExpiry := time.NewTicker(time.Duration(WATCHDOG_CHECK_INTERVAL * int(time.Second)))

	fmt.Printf("[i] Cyber Shift GDPR Bot is live. Checking page every: %d seconds\n", WATCHDOG_CHECK_INTERVAL)

	for ; true; <-dataExpiry.C {
		scrape()
	}
}

func serveRSSFeed(c *gin.Context) {
	feedMutex.Lock()
	rss, err := sanctionFeed.ToRss()
	feedMutex.Unlock()

	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
	}

	// Set the right type
	c.Writer.Header().Set("Content-Type", "application/rss+xml")

	c.String(http.StatusOK, rss)
}

func serveJSONFeed(c *gin.Context) {
	feedMutex.Lock()
	json, err := sanctionFeed.ToJSON()
	feedMutex.Unlock()

	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
	}

	// Set the right type
	c.Writer.Header().Set("Content-Type", "application/json")

	c.String(http.StatusOK, json)
}

func serveAtomFeed(c *gin.Context) {
	feedMutex.Lock()
	atom, err := sanctionFeed.ToAtom()
	feedMutex.Unlock()

	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
	}

	// Set the right type
	c.Writer.Header().Set("Content-Type", "application/atom+xml")

	c.String(http.StatusOK, atom)
}
