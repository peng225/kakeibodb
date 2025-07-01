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
		id, err := strconv.ParseInt(strEvent[0], 10, 64)
		require.NoError(t, err)
		date, err := model.ParseDate(strEvent[1])
		require.NoError(t, err)
		money, err := strconv.ParseInt(strEvent[2], 10, 32)
		require.NoError(t, err)
		events[i] = model.NewEventWithID(
			id, *date, int32(money), strEvent[3], nil,
		)
		for _, tag := range strings.Split(strEvent[4], ",") {
			events[i].AddTag(model.Tag(tag))
		}
	}
	return events
}

func parseTagList(t *testing.T, rawTagList []byte) []*model.TagWithID {
	t.Helper()
	strTagLines := strings.Split(strings.TrimSpace(string(rawTagList)), "\n")
	require.LessOrEqual(t, 1, len(strTagLines))
	// Skip header
	strTagLines = strTagLines[1:]
	tags := make([]*model.TagWithID, len(strTagLines))
	for i, strTagLine := range strTagLines {
		strTag := strings.Fields(strTagLine)
		require.Len(t, strTag, 2, strTagLine)
		id, err := strconv.ParseInt(strTag[0], 10, 64)
		require.NoError(t, err)
		tags[i] = model.NewTagWithID(
			id, model.Tag(strTag[1]),
		)
	}
	return tags
}

func parsePatternList(t *testing.T, rawPatternList []byte) []*model.PatternWithID {
	t.Helper()
	strPatternLines := strings.Split(strings.TrimSpace(string(rawPatternList)), "\n")
	require.LessOrEqual(t, 1, len(strPatternLines))
	// Skip header
	strPatternLines = strPatternLines[1:]
	patterns := make([]*model.PatternWithID, len(strPatternLines))
	for i, strPatternLine := range strPatternLines {
		strPattern := strings.Fields(strPatternLine)
		require.Len(t, strPattern, 3, strPatternLine)
		id, err := strconv.ParseInt(strPattern[0], 10, 64)
		require.NoError(t, err)
		patterns[i] = model.NewPatternWithID(
			id, strPattern[1], nil,
		)
		for _, tag := range strings.Split(strPattern[2], ",") {
			patterns[i].AddTag(model.Tag(tag))
		}
	}
	return patterns
}
