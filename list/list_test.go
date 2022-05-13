package list

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCrawledList(t *testing.T) {
	urlToAdd := "www.test.com"
	cr := NewSyncedList(logrus.New())
	require.Equal(t, false, cr.Contains(urlToAdd))

	cr.Add(urlToAdd)
	require.Equal(t, true, cr.Contains(urlToAdd))
}
