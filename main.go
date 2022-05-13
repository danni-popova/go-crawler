package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/danni-popova/go-crawler/list"
	"github.com/danni-popova/go-crawler/parser"
	"github.com/danni-popova/go-crawler/queue"
	"github.com/danni-popova/go-crawler/reader"
	"github.com/danni-popova/go-crawler/robots"
	log "github.com/sirupsen/logrus"
)

var (
	startingPage = "https://monzo.com"
)

func main() {
	parseFlags()

	// Set up logger settings
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "19:57:38.5621",
	})

	// Set up http client settings
	client := http.Client{Timeout: time.Second * 5}

	// Make a list to store the (successfully) visited pages
	visitedUrls := list.NewSyncedList(logger)

	// Make a queue for the workers to grab URLs from
	q := queue.NewQueue(logger)

	// Setup page reader
	pr := reader.NewPageReader(startingPage, visitedUrls, client, logger)

	// Try to get the robots list for the starting page
	robotsList := robots.NewRobotsList(startingPage, logger, pr)

	// Create parser
	prs := parser.NewParser(q, robotsList, visitedUrls, logger)

	// Finally, create and set up the crawler
	w := NewCrawler(startingPage, pr, prs, q, logger)
	err := w.Setup(startingPage)
	if err != nil {
		// TODO: There's a lot more that can be done for input validation, but I that's not as interesting as the rest
		logger.Fatal("Could not lookup the page you provided. Please make sure it's in this format https://monzo.com")
	}
	w.Run()
	visitedUrls.PrintContents()
}

func parseFlags() {
	flag.StringVar(&startingPage, "starting-url", "https://monzo.com", "The starting page for the crawling e.g. https://monzo.com")
	flag.Parse()
}
