package event

import (
	"bufio"
	"log"
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

func (c *CSV) Open(filePath string) {
	var err error
	c.fp, err = os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *CSV) Close() {
	c.fp.Close()
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
