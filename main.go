package main

import (
	"flag"
	"log"
)

func main() {
	var feedUrl string
	var interval int
	flag.StringVar(&feedUrl, "feed", "", "Url of the rss feed to monitor")
	flag.IntVar(&interval, "interval", 6, "Time in hours since the feed was last checked")
	flag.Parse()
	if feedUrl == "" {
		log.Fatal("No feed given")
	}
}
