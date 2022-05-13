package reader

import (
	"io"
	"net/http"
	"net/url"

	"github.com/danni-popova/go-crawler/list"
	"github.com/sirupsen/logrus"
)

const (
	USERAGENT = "DanniCrawl v0.1"
)

// pageReader interface abstracts the network calls and mock them for testing purposes
type pageReader struct {
	startingUrl string
	visitedUrls list.SyncedList
	client      http.Client
	log         *logrus.Logger
}

type PageReader interface {
	GetPageContents(url string) (io.ReadCloser, error)
}

func NewPageReader(startingUrl string, visited list.SyncedList, client http.Client, log *logrus.Logger) PageReader {
	return &pageReader{
		startingUrl: startingUrl,
		visitedUrls: visited,
		client:      client,
		log:         log,
	}
}

// GetPageContents tries to send a GET request and returns the contents of a page
func (p *pageReader) GetPageContents(urlString string) (io.ReadCloser, error) {
	p.visitedUrls.Add(urlString)
	requestUrl := p.startingUrl + urlString

	_, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", USERAGENT)
	response, err := p.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}
