package event

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CSV struct {
	fp *os.File
	sc *bufio.Scanner
}

func NewCSV() *CSV {
	return &CSV{}
}

func (c *CSV) Open(filePath string) error {
	var err error
	c.fp, err = os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	return nil
}

func (c *CSV) Close() error {
	err := c.fp.Close()
	if err != nil {
		return fmt.Errorf("failed to close: %w", err)
	}
	return nil
}

func (c *CSV) Read() []string {
	if c.sc == nil {
		c.sc = bufio.NewScanner(c.fp)
	}
	var result []string
	if c.sc.Scan() {
		line := c.sc.Text()
		result = strings.Split(line, ",")
	}
	return result
}
