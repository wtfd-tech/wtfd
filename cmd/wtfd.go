package main

import (
	"gitlab.fachschaften.org/foss-ag/wtfd/internal"
	"log"
)

func main() {
	log.Fatal(wtfd.Server())
}
