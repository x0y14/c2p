package main

import (
	"c2p"
	"c2p/c/parse"
	"c2p/c/tokenize"
	"c2p/python/Gen"
	"fmt"
	"syscall/js"
)

func CToPython(_ js.Value, args []js.Value) any {
	// c
	tokens, err := tokenize.Tokenize(args[0].String())
	if err != nil {
		//log.Fatalf("failed to tokenize c: %s", err)
		return map[string]any{"status": 1, "result": fmt.Sprintf("failed to tokenize c: %s", err)}
	}
	cNodes, err := parse.Parse(tokens)
	if err != nil {
		//log.Fatalf("failed to parse c: %s", err)
		return map[string]any{"status": 1, "result": fmt.Sprintf("failed to parse c: %s", err)}
	}

	// c -> p
	pNodes, err := c2p.CToP(cNodes)
	if err != nil {
		//log.Fatalf("failed to convert c->p: %s", err)
		return map[string]any{"status": 1, "result": fmt.Sprintf("failed to convert c->p: %s", err)}
	}

	// p
	pCode, err := Gen.Gen(pNodes)
	if err != nil {
		//log.Fatalf("failed to generate p: %s", err)
		return map[string]any{"status": 1, "result": fmt.Sprintf("failed to generate p: %s", err)}
	}

	return map[string]any{"status": 0, "result": pCode}
}

func main() {
	c := make(chan struct{})
	js.Global().Set("CToP", js.FuncOf(CToPython))
	<-c
}
