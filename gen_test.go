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
    int i;
    for (i = 1; i <= 100; i=i+1) {
        if (i % 3 == 0 && i % 5 == 0) {
            printf("FizzBuzz\n");
        } else if (i % 3 == 0) {
            printf("Fizz\n");
        } else if (i % 5 == 0) {
            printf("Buzz\n");
        } else {
            printf("%d\n", i);
        }
    }
    return 0;
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
