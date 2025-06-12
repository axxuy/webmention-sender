package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/axxuy/webmention-sender/feed"
)

func main() {
	var feedUrl string
	var interval int
	flag.StringVar(&feedUrl, "feed", "", "Url of the rss feed to monitor")
	flag.IntVar(&interval, "interval", 6, "Time in hours since the feed was last checked")
	firstRun := flag.Bool("first-run", false, "Is this the first time you have checked this feed?")
	flag.Parse()
	if feedUrl == "" {
		log.Fatal("No feed given")
	}
	var lastRun time.Time
	if *firstRun {
		lastRun = time.Time{}
	} else {
		delta := -time.Hour * time.Duration(interval)
		now := time.Now()
		lastRun = now.Add(delta)
	}
	entries, err := feed.Fetch(feedUrl, &lastRun)
	if err != nil {
		log.Fatal("Could not retrieve feed: " + err.Error())
	}
	fmt.Println(entries)

}
