package test

import (
	"bytes"
	"os/exec"
	"strconv"
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
