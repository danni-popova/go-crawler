package queue

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
)

func TestQueue(t *testing.T) {
	q := NewQueue(logrus.New())
	testInsert := "some-string"

	q.Insert(testInsert)
	// The inserted value should be the one to be removed next
	testRemove, err := q.Remove()
	require.Equal(t, testInsert, testRemove)
	require.NoError(t, err)

	// The queue should now be empty and return an error
	_, err = q.Remove()
	require.Error(t, err)
}
