package main

import (
	"sync"

	"github.com/danni-popova/go-crawler/parser"
	"github.com/danni-popova/go-crawler/queue"
	"github.com/danni-popova/go-crawler/reader"
	"github.com/sirupsen/logrus"
)

type Crawler interface {
	Run()
	Setup(url string) error
}

// Everything that a single crawler needs to run
type crawler struct {
	startingUrl string            // The URL to start from
	pageReader  reader.PageReader // Reader to connect to the pages
	parser      parser.Parser     // Parser to extract the links
	queue       queue.Queue       // A queue to store and get URLs from
	log         *logrus.Logger    // Logger for debugging
}

func NewCrawler(url string, pr reader.PageReader, prs parser.Parser, q queue.Queue, log *logrus.Logger) Crawler {
	return &crawler{
		startingUrl: url,
		pageReader:  pr,
		parser:      prs,
		queue:       q,
		log:         log,
	}
}

// Setup visits the starting URL and fills the queue with all the found pages
func (w crawler) Setup(url string) error {
	// Get starting page contents
	contents, err := w.pageReader.GetPageContents("")
	if err != nil {
		return err
	}

	// Process starting page
	err = w.parser.ProcessPage(url, contents)
	if err != nil {
		return err
	}
	return nil
}

// Run creates and run the goroutines that pick up pages from the queue and process them
func (w crawler) Run() {
	// Create a wait group to know when all the work is done
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go w.work(&wg)
	}

	// Wait for all the go routines to finish
	wg.Wait()
}

func (w crawler) work(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		url, err := w.queue.Remove()
		if err != nil {
			w.log.Warn(err)
			return
		}

		contents, err := w.pageReader.GetPageContents(url)
		if err != nil {
			continue
		}

		err = w.parser.ProcessPage(url, contents)
		if err != nil {
			w.log.Error(err)
		}
	}
}
