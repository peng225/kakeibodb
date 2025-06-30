package test

import (
	"kakeibodb/internal/model"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func parseEventList(t *testing.T, rawEventList []byte) []*model.EventWithID {
	t.Helper()
	strEventLines := strings.Split(strings.TrimSpace(string(rawEventList)), "\n")
	require.LessOrEqual(t, 1, len(strEventLines))
	// Skip header
	strEventLines = strEventLines[1:]
	events := make([]*model.EventWithID, len(strEventLines))
	for i, strEventLine := range strEventLines {
		strEvent := strings.Fields(strEventLine)
		require.Len(t, strEvent, 5, strEventLine)
		id, err := strconv.ParseInt(strEvent[0], 10, 32)
		require.NoError(t, err)
		date, err := model.ParseDate(strEvent[1])
		require.NoError(t, err)
		money, err := strconv.ParseInt(strEvent[2], 10, 32)
		require.NoError(t, err)
		events[i] = model.NewEventWithID(
			int32(id), *date, int32(money), strEvent[3], nil,
		)
		for _, tag := range strings.Split(strEvent[4], ",") {
			events[i].AddTag(model.Tag(tag))
		}
	}
	return events
}
