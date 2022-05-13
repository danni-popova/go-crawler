package robots

import (
	"crawler/mocks"
	"io"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testRobotsPage = "# robotstxt.org/\n\nUser-agent: *\nDisallow: /docs/\nDisallow: /referral/\nDisallow: /-staging-referral/\nDisallow: /install/"
)

func TestRobotsList(t *testing.T) {
	startingUrl := "https://monzo.com"

	rc := io.NopCloser(strings.NewReader(testRobotsPage))
	mockPageReader := mocks.NewPageReader(t)
	mockPageReader.On("GetPageContents", mock.Anything).Return(rc, nil)

	rl := NewRobotsList(startingUrl, logrus.New(), mockPageReader)
	require.Equal(t, true, rl.Contains(startingUrl+"/docs/"))
	require.Equal(t, true, rl.Contains(startingUrl+"/referral/"))
	require.Equal(t, true, rl.Contains(startingUrl+"/-staging-referral/"))
	require.Equal(t, true, rl.Contains(startingUrl+"/install/"))
}
