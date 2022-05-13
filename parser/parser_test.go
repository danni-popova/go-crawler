package parser

import (
	"bufio"
	"crawler/mocks"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	testHtml = "<!DOCTYPE html><html lang=\"en\"> <head> <meta charset=\"utf-8\"> <title>title</title> <link rel=\"stylesheet\" href=\"style.css\"> <script src=\"script.js\"></script> </head> <body> <a href=\"/flex\">Mozno Flex</a> <a href=\"/flex/\">Monzo Flex again</a> <a href=\"/isa/\">Cash ISAs</a> <a class=\"c-footer__social-link\" href=\"https://www.facebook.com/monzobank\"> <img src=\"/static/images/facebook.svg\" alt=\"Facebook\"> </a></body></html>"
)

// TODO: create a simple html page string for the test instead of using the file
func TestParser_ProcessPage(t *testing.T) {
	logger := logrus.New()

	mockQueue := mocks.NewQueue(t)
	mockQueue.On("Insert", mock.Anything).Return()

	mockRobots := mocks.NewRobotsList(t)
	mockRobots.On("Contains", mock.Anything).Return(false)

	mockVisited := mocks.NewSyncedList(t)
	mockVisited.On("Contains", mock.Anything).Return(false)

	contents := strings.NewReader(testHtml)

	prs := NewParser(mockQueue, mockRobots, mockVisited, logger)
	err := prs.ProcessPage("https://monzo.com", bufio.NewReader(contents))
	require.NoError(t, err)

	// Only the pages from the same domain should be added to the queue
	mockQueue.AssertCalled(t, "Insert", "/flex/")
	mockQueue.AssertCalled(t, "Insert", "/isa/")
	mockQueue.AssertNotCalled(t, "Insert", "https://www.facebook.com/monzobank/")
}
