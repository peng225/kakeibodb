package event

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type BankCSV struct {
	fp *os.File
	sc *bufio.Scanner
}

func NewBankCSV() *BankCSV {
	return &BankCSV{}
}

func (c *BankCSV) Open(filePath string) {
	var err error
	c.fp, err = os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *BankCSV) Close() {
	c.fp.Close()
}

func (c *BankCSV) Read() []string {
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
