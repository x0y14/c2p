package c2p

import (
	"c2p/c/parse"
	"c2p/c/tokenize"
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{
			"",
			`
int main(void) {
	return 10000;
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize.Tokenize(tt.in)
			if err != nil {
				t.Errorf("failed to tokenize: %v", err)
			}

			nodes, err := parse.Parse(tokens)
			if err != nil {
				t.Errorf("failed to parse: %v", err)
			}

			code, err := Gen(nodes)
			if err != nil {
				t.Errorf("failed to gen: %v", err)
			}

			fmt.Println(code)
		})
	}
}
