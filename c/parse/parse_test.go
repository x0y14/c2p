package parse

import (
	"c2p/c/tokenize"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{
			"fizzbuzz nopre",
			`
int main(void) {
    int i;
    for (i = 1; i <= 100; i=i+1) {
        if (1 == 1) {
			return 1;
		} else if (i % 3 == 0) {
			return 2;
        }
    }
    return 0;
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize.Tokenize(tt.in)
			if err != nil {
				t.Errorf("failed tokenize: %v", err)
			}
			nodes, err := Parse(tokens)
			if err != nil {
				t.Errorf("failed parse: %v", err)
			}
			log.Println(nodes)
		})
	}
}
