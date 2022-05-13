package robots

import (
	"io"
	"strings"

	"github.com/danni-popova/go-crawler/reader"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type RobotsList interface {
	Contains(url string) bool
	Print()
}

type robotsList struct {
	log          *logrus.Logger
	pr           reader.PageReader
	startingUrl  string
	disallowList []string
}

const (
	robotsUrl = "/robots.txt"
)

func NewRobotsList(startingUrl string, log *logrus.Logger, pr reader.PageReader) RobotsList {
	rl := robotsList{
		startingUrl:  startingUrl,
		log:          log,
		pr:           pr,
		disallowList: []string{},
	}

	rl.fillRobotsList()
	return &rl
}

func (r *robotsList) fillRobotsList() {
	body, err := r.pr.GetPageContents(robotsUrl)

	if err != nil {
		r.log.Warnf("Coudln't get robots page, disallow list will be empty")
		return
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, body)
	if err != nil {
		r.log.Warnf("Couldn't read the robots page body, disallow list will be empty")
		return
	}

	err = r.parseFile(buf.String())
	if err != nil {
		r.log.Warnf("Couldn't parse robots page body, disallow list will be empty")
		return
	}
}

func (r *robotsList) parseFile(robotsContents string) error {
	lines := strings.Split(robotsContents, "\n")
	// TODO: I am assuming that the case here is User-agent: * to save some time going further into the robots parsing
	for _, line := range lines {
		if strings.Contains(line, "Disallow: ") {
			result := strings.SplitAfter(line, "Disallow: ")
			r.log.Tracef("Adding to disallow list: " + result[1])
			r.disallowList = append(r.disallowList, r.startingUrl+result[1])
		}
	}
	return nil
}

func (r *robotsList) Contains(url string) bool {
	return slices.Contains(r.disallowList, url)
}

func (r *robotsList) Print() {
	r.log.Info("Pages in the disallow list:")
	for _, d := range r.disallowList {
		r.log.Info(d)
	}
}
