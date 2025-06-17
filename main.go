package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/axxuy/webmention-sender/feed"
	"github.com/axxuy/webmention-sender/webmention"
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
	doFeedWebmentions(feedUrl, &lastRun)

}

func doFeedWebmentions(feedUrl string, lastRun *time.Time) {

	entries, err := feed.Fetch(feedUrl, lastRun)
	if err != nil {
		log.Fatal("Could not retrieve feed: " + err.Error())
	}
	endpoints := make(map[string]*webmention.Endpoint)
	for _, entry := range entries {
		for _, link := range entry.Links {
			host := link.Hostname()
			//Do we already have an endpoint for this domain?
			endpoint, ok := endpoints[host]
			if !ok {
				//Try to get one
				endpoint, err := webmention.GetWebmentionEndpoint(link)
				//If that fails record it in the table
				if err != nil || endpoint == nil {
					endpoints[host] = nil
					continue
				}
				endpoints[host] = endpoint
			}
			//Have we previously failed to get an endpoint for this domain?
			if endpoint == nil {
				continue
			}

			err = endpoint.Send(link, entry.Url)
			fmt.Fprintf(os.Stderr, "Error: %v sending to %v\n", err.Error(), link.String())
		}
	}
}
