package main

import (
	"sync"
	"time"

	"github.com/gorilla/feeds"
)

var feedMutex sync.Mutex

var lastSanctionsFound = 0

var sanctionFeed = &feeds.Feed{
	Title:       "Cyber Shift GDPR Bot",
	Link:        &feeds.Link{Href: "https://gdpr.cybershift.dev"},
	Description: "O listă cu toate amenzile anunțate de către Autoritatea Naţională de Supraveghere a Prelucrării Datelor cu Caracter Personal.",
	Author:      &feeds.Author{Name: "Cyber Shift", Email: "office@cybershift.dev"},
	Created:     time.Now(),
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
