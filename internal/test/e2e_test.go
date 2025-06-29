package test

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "embed"

	"github.com/stretchr/testify/require"
)

const (
	kakeibodb = "../../kakeibodb"
	dbPort    = 3307
)

var commonOptions = []string{"--dbname", "testdb",
	"--dbport", strconv.Itoa(dbPort), "-u", "test"}

func runCommand(stdin []byte, cmd string, args ...string) ([]byte, []byte, error) {
	c := exec.Command(cmd, args...)
	c.Stdin = bytes.NewReader(stdin)

	var stdoutBuf, stderrBuf bytes.Buffer
	c.Stdout = &stdoutBuf
	c.Stderr = &stderrBuf
	err := c.Run()

	return stdoutBuf.Bytes(), stderrBuf.Bytes(), err
}

func runKakeiboDB(args ...string) ([]byte, []byte, error) {
	completeArgs := append(args, commonOptions...)
	return runCommand(nil, kakeibodb, completeArgs...)
}

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
	stdout, stderr, err = runCommand(stdout, "grep", "クレジット")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "awk", "{print $1}")
	require.NoError(t, err, string(stderr))
	creditEventID := strings.TrimSpace(string(stdout))
	_, stderr, err = runKakeiboDB("event", "load", "--credit",
		"--parentEventID", creditEventID, "-f", "credit/cmeisai1.csv")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	_, _, err = runCommand(stdout, "grep", "クレジット")
	require.Error(t, err)
	_, stderr, err = runCommand(stdout, "grep", "チョコ")
	require.NoError(t, err, string(stderr))
}

func getEventsWithTags(t *testing.T, tags ...string) []byte {
	stdout, stderr, err := runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	for _, tag := range tags {
		stdout, stderr, err = runCommand(stdout, "grep", tag)
		require.NoError(t, err, string(stderr))
	}
	return stdout
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
	_, stderr, err = runCommand(stdout, "grep", "foo")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runCommand(stdout, "grep", "bar")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "addTag", "--eventID", "1", "--tagNames", "foo,bar")
	require.NoError(t, err, string(stderr))
	// Idempotency check.
	_, stderr, err = runKakeiboDB("event", "addTag", "--eventID", "1", "--tagNames", "foo,bar")
	require.NoError(t, err, string(stderr))
	stdout = getEventsWithTags(t, "foo", "bar")
	stdout, stderr, err = runCommand(stdout, "awk", "{print $1}")
	require.NoError(t, err, string(stderr))
	eventID := strings.TrimSpace(string(stdout))
	require.Equal(t, "1", eventID)

	// Remove tags from a event.
	_, stderr, err = runKakeiboDB("event", "removeTag", "--eventID", "1", "-t", "foo")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "removeTag", "--eventID", "1", "-t", "bar")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	_, _, err = runCommand(stdout, "grep", "-e", "foo", "-e", "bar")
	require.Error(t, err)

	// Delete tags.
	_, stderr, err = runKakeiboDB("tag", "delete", "--tagID", "1")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "delete", "--tagID", "2")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("tag", "list")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runCommand(stdout, "grep", "-e", "foo", "-e", "bar")
	require.Error(t, err, string(stderr))
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

func getPatternsWithTags(t *testing.T, tags ...string) []byte {
	stdout, stderr, err := runKakeiboDB("pattern", "list")
	require.NoError(t, err, string(stderr))
	for _, tag := range tags {
		stdout, stderr, err = runCommand(stdout, "grep", tag)
		require.NoError(t, err, string(stderr))
	}
	return stdout
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
	var stdout []byte
	stdout = getPatternsWithTags(t, "fruit", "yellow")
	stdout, stderr, err = runCommand(stdout, "grep", "バナ")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "awk", "{print $1}")
	require.NoError(t, err, string(stderr))
	patternID := strings.TrimSpace(string(stdout))
	require.Equal(t, "1", patternID)

	_, stderr, err = runKakeiboDB("event", "applyPattern",
		"--from", "2022-01-04", "--to", "2022-02-03")
	require.NoError(t, err, string(stderr))
	// Idempotency check.
	_, stderr, err = runKakeiboDB("event", "applyPattern",
		"--from", "2022-01-04", "--to", "2022-02-03")
	require.NoError(t, err, string(stderr))
	stdout = getEventsWithTags(t, "fruit", "yellow")
	_, stderr, err = runCommand(stdout, "grep", "バナ")
	require.NoError(t, err, string(stderr))

	// Remove tags from a pattern.
	_, stderr, err = runKakeiboDB("pattern", "removeTag",
		"--patternID", "1", "-t", "fruit")
	require.NoError(t, err, string(stderr))
	stdout = getPatternsWithTags(t, "yellow")
	stdout, stderr, err = runCommand(stdout, "grep", "バナ")
	require.NoError(t, err, string(stderr))
	_, _, err = runCommand(stdout, "grep", "fruit")
	require.Error(t, err)

	// Delete a pattern.
	_, stderr, err = runKakeiboDB("pattern", "delete", "--patternID", "1")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("pattern", "list")
	require.NoError(t, err, string(stderr))
	_, _, err = runCommand(stdout, "grep", "バナ")
	require.Error(t, err)
}

func TestSplit(t *testing.T) {
	dbSetup(t)
	t.Cleanup(func() { dbCleanup(t) })
	_, stderr, err := runKakeiboDB("event", "load", "-d", "event")
	require.NoError(t, err, string(stderr))
	var stdout []byte
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "grep", "クレジット")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "awk", "{print $1}")
	require.NoError(t, err, string(stderr))
	creditEventID := strings.TrimSpace(string(stdout))
	_, stderr, err = runKakeiboDB("event", "load", "--credit",
		"--parentEventID", creditEventID, "-f", "credit/cmeisai1.csv")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("tag", "create", "-t", "candy")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "addTag", "--eventID", "10", "--tagNames", "candy")
	require.NoError(t, err, string(stderr))

	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "grep", "飴")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "awk", "{print $3}")
	require.NoError(t, err, string(stderr))
	beforeCandyMoney, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	require.NoError(t, err)

	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "split", "--eventID", "10",
		"--date", "2021-12-04", "--money", "-30", "--desc", "はちみつのど飴")
	require.NoError(t, err, string(stderr))
	_, stderr, err = runKakeiboDB("event", "split", "--eventID", "10",
		"--date", "2021/12/05", "--money", "-30", "--desc", "きんかんのど飴")
	require.NoError(t, err, string(stderr))

	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "grep", "飴")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "awk", "{s += $3} END {print s}")
	require.NoError(t, err, string(stderr))
	afterCandyMoneySum, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	require.NoError(t, err)
	require.Equal(t, beforeCandyMoney, afterCandyMoneySum)

	os.Setenv("KAKEIBODB_SPLIT_BASE_TAG_NAME", "candy")
	t.Cleanup(func() { os.Unsetenv("KAKEIBODB_SPLIT_BASE_TAG_NAME") })
	_, stderr, err = runKakeiboDB("event", "split",
		"--date", "2021-12-06", "--money", "-40", "--desc", "ミルク飴")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runKakeiboDB("event", "list")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "grep", "飴")
	require.NoError(t, err, string(stderr))
	stdout, stderr, err = runCommand(stdout, "awk", "{s += $3} END {print s}")
	require.NoError(t, err, string(stderr))
	afterCandyMoneySum, err = strconv.Atoi(strings.TrimSpace(string(stdout)))
	require.NoError(t, err)
	require.Equal(t, beforeCandyMoney, afterCandyMoneySum)
}
