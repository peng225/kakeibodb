package test

import (
	"kakeibodb/internal/model"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "embed"

	"github.com/stretchr/testify/require"
)

//go:embed setup.sql
var setupSQL []byte

//go:embed cleanup.sql
var cleanupSQL []byte

func dbSetup(t *testing.T) {
	t.Helper()
	require.Eventually(t, func() bool {
		_, _, err := runCommand(setupSQL, "mysql", "-h", "127.0.0.1",
			"--port", strconv.Itoa(dbPort), "-B", "-u", "root")
		return err != nil
	}, 2*time.Second, 100*time.Millisecond)
}

func dbCleanup(t *testing.T) {
	t.Helper()
	_, stderr, err := runCommand(cleanupSQL, "mysql", "-h", "127.0.0.1",
		"--port", strconv.Itoa(dbPort), "-B", "-u", "root")
	require.NoError(t, err, string(stderr))
}

func TestEvent(t *testing.T) {
	dbSetup(t)
	t.Cleanup(func() { dbCleanup(t) })
	_, stderr, err := runKakeiboDB("event", "load", "-d", "event")
	require.NoError(t, err, string(stderr))
	var stdout []byte
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events := parseEventList(t, stdout)
	i := slices.IndexFunc(events, func(e *model.EventWithID) bool {
		return e.GetDesc() == "クレジットカード"
	})
	require.NotEqual(t, -1, i)
	creditEventID := strconv.FormatInt(int64(events[i].GetID()), 10)

	_, stderr, err = runKakeiboDB("event", "load", "--credit",
		"--parentEventID", creditEventID, "-f", "credit/cmeisai1.csv")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events = parseEventList(t, stdout)
	i = slices.IndexFunc(events, func(e *model.EventWithID) bool {
		return e.GetDesc() == "クレジットカード"
	})
	require.Equal(t, -1, i)
	i = slices.IndexFunc(events, func(e *model.EventWithID) bool {
		return e.GetDesc() == "チョコ"
	})
	require.NotEqual(t, -1, i)
}

func getEventsWithAllTags(t *testing.T, tags ...string) []*model.EventWithID {
	stdout, stderr, err := runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events := parseEventList(t, stdout)
	seq := func(yield func(*model.EventWithID) bool) {
		for _, e := range events {
			for _, tag := range e.GetTags() {
				if !slices.ContainsFunc(tags, func(t string) bool {
					return t == tag.String()
				}) {
					return
				}
			}
			if !yield(e) {
				return
			}
		}
	}
	return slices.Collect(seq)
}

func TestTag(t *testing.T) {
	dbSetup(t)
	t.Cleanup(func() { dbCleanup(t) })
	_, stderr, err := runKakeiboDB("event", "load", "-d", "event")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "create", "-t", "foo")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "create", "-t", "bar")
	require.NoError(t, err, string(stderr))
	var stdout []byte
	stdout, stderr, err = runKakeiboDB("tag", "list")
	require.NoError(t, err, string(stderr))
	tagList := parseTagList(t, stdout)
	require.Equal(t, "foo", tagList[0].String())
	require.Equal(t, "bar", tagList[1].String())
	_, stderr, err = runKakeiboDB("event", "addTag", "--eventID", "1", "--tagNames", "foo,bar")
	require.NoError(t, err, string(stderr))
	// Idempotency check.
	_, stderr, err = runKakeiboDB("event", "addTag", "--eventID", "1", "--tagNames", "foo,bar")
	require.NoError(t, err, string(stderr))
	eventsWithAllTags := getEventsWithAllTags(t, "foo", "bar")
	require.Len(t, eventsWithAllTags, 1)
	require.Equal(t, int64(1), eventsWithAllTags[0].GetID())

	// Remove tags from a event.
	_, stderr, err = runKakeiboDB("event", "removeTag", "--eventID", "1", "-t", "foo")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "removeTag", "--eventID", "1", "-t", "bar")
	require.NoError(t, err, string(stderr))
	// Idempotency check.
	_, stderr, err = runKakeiboDB("event", "removeTag", "--eventID", "1", "-t", "bar")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events := parseEventList(t, stdout)
	i := slices.IndexFunc(events, func(e *model.EventWithID) bool {
		for _, tag := range e.GetTags() {
			if tag.String() == "foo" || tag.String() == "bar" {
				return true
			}
		}
		return false
	})
	require.Equal(t, -1, i)

	// Delete tags.
	_, stderr, err = runKakeiboDB("tag", "delete", "--tagID", "1")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "delete", "--tagID", "2")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("tag", "list")
	require.NoError(t, err, string(stderr))
	tags := parseTagList(t, stdout)
	require.Empty(t, tags)
}

func TestUserEnv(t *testing.T) {
	dbSetup(t)
	t.Cleanup(func() { dbCleanup(t) })
	os.Setenv("KAKEIBODB_USER", "test")
	t.Cleanup(func() { os.Unsetenv("KAKEIBODB_USER") })
	completeArgs := append([]string{"event", "list"}, commonOptions[:4]...)
	_, stderr, err := runCommand(nil, kakeibodb, completeArgs...)
	require.NoError(t, err, string(stderr))
}

func getPatternsWithAllTags(t *testing.T, tags ...string) []*model.PatternWithID {
	stdout, stderr, err := runKakeiboDB("pattern", "list")
	require.NoError(t, err, string(stderr))
	patterns := parsePatternList(t, stdout)
	seq := func(yield func(*model.PatternWithID) bool) {
		for _, p := range patterns {
			for _, tag := range p.GetTags() {
				if !slices.ContainsFunc(tags, func(t string) bool {
					return t == tag.String()
				}) {
					return
				}
			}
			if !yield(p) {
				return
			}
		}
	}
	return slices.Collect(seq)
}

func TestPattern(t *testing.T) {
	dbSetup(t)
	t.Cleanup(func() { dbCleanup(t) })
	_, stderr, err := runKakeiboDB("event", "load", "-d", "event")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "create", "-t", "fruit")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "create", "-t", "yellow")
	require.NoError(t, err, string(stderr))

	_, stderr, err = runKakeiboDB("pattern", "create", "-k", "バナ")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("pattern", "addTag", "--patternID", "1",
		"--tagNames", "fruit,yellow")
	require.NoError(t, err, string(stderr))
	// Idempotency check.
	_, stderr, err = runKakeiboDB("pattern", "addTag", "--patternID", "1",
		"--tagNames", "fruit,yellow")
	require.NoError(t, err, string(stderr))
	patternsWithAllTags := getPatternsWithAllTags(t, "fruit", "yellow")
	require.NotEmpty(t, patternsWithAllTags)
	require.Equal(t, "バナ", patternsWithAllTags[0].GetKey())
	require.Equal(t, int64(1), patternsWithAllTags[0].GetID())

	_, stderr, err = runKakeiboDB("event", "applyPattern",
		"--from", "2022-01-04", "--to", "2022-02-03")
	require.NoError(t, err, string(stderr))
	// Idempotency check.
	_, stderr, err = runKakeiboDB("event", "applyPattern",
		"--from", "2022-01-04", "--to", "2022-02-03")
	require.NoError(t, err, string(stderr))
	eventsWithAllTags := getEventsWithAllTags(t, "fruit", "yellow")
	for _, ewt := range eventsWithAllTags {
		require.Contains(t, ewt.GetDesc(), "バナ")
	}

	// Remove tags from a pattern.
	_, stderr, err = runKakeiboDB("pattern", "removeTag",
		"--patternID", "1", "-t", "fruit")
	require.NoError(t, err, string(stderr))
	// Idempotency check.
	_, stderr, err = runKakeiboDB("pattern", "removeTag",
		"--patternID", "1", "-t", "fruit")
	require.NoError(t, err, string(stderr))
	patternsWithAllTags = getPatternsWithAllTags(t, "fruit")
	require.Empty(t, patternsWithAllTags)
	patternsWithAllTags = getPatternsWithAllTags(t, "yellow")
	require.NotEmpty(t, patternsWithAllTags)

	// Delete a pattern.
	_, stderr, err = runKakeiboDB("pattern", "delete", "--patternID", "1")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err := runKakeiboDB("pattern", "list")
	require.NoError(t, err, string(stderr))
	patterns := parsePatternList(t, stdout)
	require.Empty(t, patterns)
}

func TestSplit(t *testing.T) {
	dbSetup(t)
	t.Cleanup(func() { dbCleanup(t) })
	_, stderr, err := runKakeiboDB("event", "load", "-d", "event")
	require.NoError(t, err, string(stderr))
	var stdout []byte
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events := parseEventList(t, stdout)
	i := slices.IndexFunc(events, func(e *model.EventWithID) bool {
		return e.GetDesc() == "クレジットカード"
	})
	require.NotEqual(t, -1, i)
	creditEventID := strconv.FormatInt(int64(events[i].GetID()), 10)
	_, stderr, err = runKakeiboDB("event", "load", "--credit",
		"--parentEventID", creditEventID, "-f", "credit/cmeisai1.csv")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "create", "-t", "candy")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "addTag", "--eventID", "10", "--tagNames", "candy")
	require.NoError(t, err, string(stderr))

	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events = parseEventList(t, stdout)
	seq := func(yield func(*model.EventWithID) bool) {
		for _, e := range events {
			if strings.Contains(e.GetDesc(), "飴") && !yield(e) {
				return
			}
		}
	}
	candyEvents := slices.Collect(seq)
	require.Len(t, candyEvents, 1)
	beforeCandyMoney := candyEvents[0].GetMoney()

	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "split", "--eventID", "10",
		"--date", "2021-12-04", "--money", "-30", "--desc", "はちみつのど飴")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "split", "--eventID", "10",
		"--date", "2021/12/05", "--money", "-30", "--desc", "きんかんのど飴")
	require.NoError(t, err, string(stderr))

	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events = parseEventList(t, stdout)
	seq = func(yield func(*model.EventWithID) bool) {
		for _, e := range events {
			if strings.Contains(e.GetDesc(), "飴") && !yield(e) {
				return
			}
		}
	}
	candyEvents = slices.Collect(seq)
	require.Len(t, candyEvents, 3)
	afterCandyMoneySum := int32(0)
	for _, e := range candyEvents {
		afterCandyMoneySum += e.GetMoney()
	}
	require.Equal(t, beforeCandyMoney, afterCandyMoneySum)

	os.Setenv("KAKEIBODB_SPLIT_BASE_TAG_NAME", "candy")
	t.Cleanup(func() { os.Unsetenv("KAKEIBODB_SPLIT_BASE_TAG_NAME") })
	_, stderr, err = runKakeiboDB("event", "split",
		"--date", "2021-12-06", "--money", "-40", "--desc", "ミルク飴")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	events = parseEventList(t, stdout)
	seq = func(yield func(*model.EventWithID) bool) {
		for _, e := range events {
			if strings.Contains(e.GetDesc(), "飴") && !yield(e) {
				return
			}
		}
	}
	// The original candy event must be deleted because its remaining money is 0.
	candyEvents = slices.Collect(seq)
	require.Len(t, candyEvents, 3)
	afterCandyMoneySum = int32(0)
	for _, e := range candyEvents {
		afterCandyMoneySum += e.GetMoney()
	}
	require.Equal(t, beforeCandyMoney, afterCandyMoneySum)
}
