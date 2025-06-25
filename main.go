package main

import (
	"kakeibodb/internal/cmd"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
