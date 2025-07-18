package test

import (
	"kakeibodb/internal/model"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func parseEventList(t *testing.T, rawEventList []byte) []*model.Event {
	t.Helper()
	strEventLines := strings.Split(strings.TrimSpace(string(rawEventList)), "\n")
	require.LessOrEqual(t, 1, len(strEventLines))
	// Skip header
	strEventLines = strEventLines[1:]
	events := make([]*model.Event, len(strEventLines))
	for i, strEventLine := range strEventLines {
		strEvent := strings.Fields(strEventLine)
		require.Len(t, strEvent, 5, strEventLine)
		id, err := strconv.ParseInt(strEvent[0], 10, 64)
		require.NoError(t, err)
		date, err := model.ParseDate(strEvent[1])
		require.NoError(t, err)
		money, err := strconv.ParseInt(strEvent[2], 10, 32)
		require.NoError(t, err)
		events[i] = model.NewEvent(
			id, *date, int32(money), strEvent[3], nil,
		)
		if len(strEvent) == 4 {
			continue
		}
		for _, tagName := range strings.Split(strEvent[4], ",") {
			events[i].AddTag(tagName)
		}
	}
	return events
}

func parseTagList(t *testing.T, rawTagList []byte) []*model.Tag {
	t.Helper()
	strTagLines := strings.Split(strings.TrimSpace(string(rawTagList)), "\n")
	require.LessOrEqual(t, 1, len(strTagLines))
	// Skip header
	strTagLines = strTagLines[1:]
	tags := make([]*model.Tag, len(strTagLines))
	for i, strTagLine := range strTagLines {
		strTag := strings.Fields(strTagLine)
		require.Len(t, strTag, 2, strTagLine)
		id, err := strconv.ParseInt(strTag[0], 10, 64)
		require.NoError(t, err)
		tags[i] = model.NewTag(
			id, strTag[1],
		)
	}
	return tags
}

func parsePatternList(t *testing.T, rawPatternList []byte) []*model.Pattern {
	t.Helper()
	strPatternLines := strings.Split(strings.TrimSpace(string(rawPatternList)), "\n")
	require.LessOrEqual(t, 1, len(strPatternLines))
	// Skip header
	strPatternLines = strPatternLines[1:]
	patterns := make([]*model.Pattern, len(strPatternLines))
	for i, strPatternLine := range strPatternLines {
		strPattern := strings.Fields(strPatternLine)
		require.GreaterOrEqual(t, len(strPattern), 2)
		id, err := strconv.ParseInt(strPattern[0], 10, 64)
		require.NoError(t, err)
		patterns[i] = model.NewPattern(
			id, strPattern[1], nil,
		)
		if len(strPattern) == 2 {
			continue
		}
		for _, tagName := range strings.Split(strPattern[2], ",") {
			patterns[i].AddTag(tagName)
		}
	}
	return patterns
}
