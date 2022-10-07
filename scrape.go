package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"
)

var (
	newsPageScrapper *colly.Collector
	articleScrapper  *colly.Collector
)

func scrapeNewsPageForSanctions(element *colly.HTMLElement) {

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

			absUrl := e.Request.AbsoluteURL(newsURL)

			// Ignore normal news articles
			if !foundASanction {
				return
			}

			tempItem.Link = &feeds.Link{Href: absUrl}

			// Reset flag
			foundASanction = false

			// Push the item on the feed
			sanctionFeed.Items = append(sanctionFeed.Items, tempItem)

			// Add article for scraping
			if articleScrapper != nil {
				articleScrapper.Visit(absUrl)
			}

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
		}
	})
}

func scrape() {
	defer executionTime(time.Now(), "Full scrape took")

	// Scrapper init
	articleScrapper = colly.NewCollector(
		colly.AllowedDomains("www.dataprotection.ro", "dataprotection.ro"),
		colly.MaxDepth(2),
		colly.Async(true),
	)

	articleScrapper.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	articleScrapper.OnHTML("div[id=rectangle_scroll]", scrapeSanctionArticle)

	// Scrapper init
	newsPageScrapper = colly.NewCollector(
		colly.AllowedDomains("www.dataprotection.ro", "dataprotection.ro"),
		colly.MaxDepth(2),
		colly.Async(true),
	)

	newsPageScrapper.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	var sanctionsFound []*feeds.Item

	newsPageScrapper.OnHTML("div[id=rectangle_scroll]", scrapeNewsPageForSanctions)

	// Start news page scrapper
	newsPageScrapper.Visit(WATCHDOG_URL)

	if lastSanctionsFound < len(sanctionsFound) {
		fmt.Printf("Found %d new sanctions. Updating the news feed.\n", len(sanctionsFound)-lastSanctionsFound)

		feedMutex.Lock()
		sanctionFeed.Items = sanctionsFound
		feedMutex.Unlock()

		lastSanctionsFound = len(sanctionsFound)
	}

	// Wait for the news scrapper to finish
	newsPageScrapper.Wait()

	fmt.Printf("Found %d total sanctions!\n", len(sanctionsFound))

	/*
		for _, l := range badLinks {
			articleScrapper.Visit(l)
		}*/

	// Wait for article scrapper. The news scrapper will add articles for scraping
	articleScrapper.Wait()

	fmt.Println("Finished scraping articles for fines.")
}

func scrapeSanctionArticle(element *colly.HTMLElement) {
	companyName := "unknown"
	var sanctionSum float64 = 0
	article := element.Request.URL.String()

	sanction_document, err := app.Dao().FindCollectionByNameOrId("sanctions")

	if err != nil {
		fmt.Println("Can't retrive sanction document: ", err)
		return
	}

	record, _ := app.Dao().FindFirstRecordByData(sanction_document, "url", article)

	// Don't continue if the sanction data has been confirmed by a human
	if record != nil {
		if record.GetBoolDataValue("human_verified") {
			//fmt.Println("Skiping article(verified): ", article)
			return
		}

		if record.GetStringDataValue("company_name") != "unknown" && record.GetStringDataValue("fine_amount") != "0" {
			//fmt.Println("Skiping article(found): ", article)
			return
		}

		//printSanction(record)
	}

	element.ForEach("p, li", func(index int, e *colly.HTMLElement) {
		match := extractCompanyName(e.Text)

		if len(match) > 0 {
			companyName = strings.TrimSpace(match)
		}

		sanction := extractSanctionSum(e.Text)

		if len(sanction) > 0 {
			// Make sure no words were included
			sanction = leaveOnlyNumbers([]byte(sanction))

			// Remove '.' from string
			sanction = strings.ReplaceAll(sanction, ".", "")

			// Transform ',' into '.' for conversion
			sanction = strings.ReplaceAll(sanction, ",", ".")

			// Convert string to int64
			var err error
			sanctionSum, err = strconv.ParseFloat(sanction, 64)

			if err != nil {
				fmt.Println("Can't parse fine amount: ", sanction, article)
				return
			}
		}
	})

	if record == nil {
		saveSanction(sanction_document, companyName, article, sanctionSum)
	} else {
		updateSanction(record, sanction_document, companyName, article, sanctionSum)
	}
}

func webScrapper() {
	if app == nil {
		fmt.Println("App not ready...")
		return
	}

	fmt.Println("Waiting for database to initialize..")

	for app.Dao() == nil {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}

	dataExpiry := time.NewTicker(time.Duration(WATCHDOG_CHECK_INTERVAL * int(time.Second)))

	fmt.Printf("[i] Cyber Shift GDPR Bot is live. Checking page every: %d seconds\n", WATCHDOG_CHECK_INTERVAL)

	for ; true; <-dataExpiry.C {
		scrape()
	}
}

func executionTime(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("\n\n%s took %s\n", name, elapsed)
}
