package loader

import (
	"bufio"
	"io"
	"strings"
)

type CSV struct {
	r  io.Reader
	sc *bufio.Scanner
}

func newCSV(r io.Reader) *CSV {
	return &CSV{
		r:  r,
		sc: bufio.NewScanner(r),
	}
}

func (c *CSV) read() []string {
	if c.sc == nil {
		c.sc = bufio.NewScanner(c.r)
	}
	var result []string
	if c.sc.Scan() {
		line := c.sc.Text()
		result = strings.Split(line, ",")
	}
	return result
}
