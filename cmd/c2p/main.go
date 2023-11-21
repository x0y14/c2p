package main

import (
	"c2p"
	"c2p/c/parse"
	"c2p/c/tokenize"
	"c2p/python/Gen"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("input file")
	}
	target := os.Args[1]
	bytes, err := os.ReadFile(target)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	// c
	tokens, err := tokenize.Tokenize(string(bytes))
	if err != nil {
		log.Fatalf("failed to tokenize c: %s", err)
	}
	cNodes, err := parse.Parse(tokens)
	if err != nil {
		log.Fatalf("failed to parse c: %s", err)
	}

	// c -> p
	pNodes, err := c2p.CToP(cNodes)
	if err != nil {
		log.Fatalf("failed to convert c->p: %s", err)
	}

	// p
	pCode, err := Gen.Gen(pNodes)
	if err != nil {
		log.Fatalf("failed to generate p: %s", err)
	}

	fmt.Println(pCode)
}
