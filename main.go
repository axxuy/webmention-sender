package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/axxuy/webmention-sender/feed"
	"github.com/axxuy/webmention-sender/webmention"
)

type Feeds []string

func (f *Feeds) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *Feeds) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var feedUrls Feeds
	var interval int
	flag.Var(&feedUrls, "feed", "Url of the rss feed to monitor")
	flag.IntVar(&interval, "interval", 6, "Time in hours since the feed was last checked")
	firstRun := flag.Bool("first-run", false, "Is this the first time you have checked this feed?")
	verbose := flag.Bool("verbose", false, "List extra information")
	slog.SetLogLoggerLevel(slog.LevelDebug)
	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
	flag.Parse()
	if len(feedUrls) == 0 {
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
	for _, feedUrl := range feedUrls {
		doFeedWebmentions(feedUrl, &lastRun)
	}

}

func doFeedWebmentions(feedUrl string, lastRun *time.Time) {

	entries, err := feed.Fetch(feedUrl, lastRun)
	if err != nil {
		log.Fatal("Could not retrieve feed: " + err.Error())
	}
	slog.Info("Retrieved feed", "feed", feedUrl, "numEntries", len(entries))
	endpoints := make(map[string]*webmention.Endpoint)
	for _, entry := range entries {
		sentLinks := make(map[*url.URL]bool)
		for _, link := range entry.Links {
			//Links may occur multiple times in a document; only send a mention for one
			_, ok := sentLinks[link]
			if ok {
				continue
			} else {
				sentLinks[link] = true
			}
			host := link.Hostname()
			//Do we already have an endpoint for this domain?
			endpoint, ok := endpoints[host]
			if !ok {
				//Try to get one
				endpoint, err = webmention.GetWebmentionEndpoint(link)
				slog.Info("Looked up link for webmention endpoint", "feed", feedUrl, "url", link, "error", err)
				//If that fails record it in the table
				if err != nil {
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
			slog.Info("Sent webmention", "target", link, "source", entry.Url, "error", err)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v sending to %v\n", err.Error(), link.String())
			}
		}
	}
}
