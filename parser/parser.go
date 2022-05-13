package parser

import (
	"io"
	"strings"

	"github.com/danni-popova/go-crawler/list"
	"github.com/danni-popova/go-crawler/queue"
	"github.com/danni-popova/go-crawler/robots"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Parser interface {
	ProcessPage(url string, body io.Reader) error
}

type parser struct {
	log          *logrus.Logger
	disallowList robots.RobotsList
	visitedUrls  list.SyncedList
	queue        queue.Queue
}

func NewParser(queue queue.Queue, robots robots.RobotsList, visited list.SyncedList, log *logrus.Logger) Parser {
	return &parser{
		log:          log,
		disallowList: robots,
		visitedUrls:  visited,
		queue:        queue,
	}
}

// ProcessPage takes the body of a page and adds the found URLs to the processing queue
func (p *parser) ProcessPage(url string, body io.Reader) error {
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			p.log.Tracef("End of page: %s", url)
			return nil
		default:
			token := tokenizer.Token()
			if p.isAnchorTag(tokenType, token) {
				cl, ok := p.extractUrlFromToken(token)
				processedUrl, ok := p.processUrl(cl)
				if ok && !p.disallowList.Contains(processedUrl) && !p.visitedUrls.Contains(processedUrl) {
					p.queue.Insert(processedUrl)
				}
			}
		}
	}
}

func (p *parser) isAnchorTag(tokenType html.TokenType, token html.Token) bool {
	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

func (p *parser) extractUrlFromToken(token html.Token) (string, bool) {
	for _, attr := range token.Attr {
		// If the key of the element is href, we've found a link to a website
		if attr.Key == "href" {
			link := attr.Val
			return link, true
		}
	}
	return "", false
}

func (p *parser) processUrl(urlString string) (string, bool) {
	switch {
	case strings.HasPrefix(urlString, "http://"):
		return "", false
	case strings.HasPrefix(urlString, "https://"):
		return "", false
	case strings.Contains(urlString, "mailto:"):
		return "", false
	default:
	}

	if urlString == "#" || urlString == "/" {
		return urlString, false
	}

	// Add trailing slash to avoid duplicates
	if !strings.HasSuffix(urlString, "/") {
		urlString = urlString + "/"
	}
	return urlString, true
}
