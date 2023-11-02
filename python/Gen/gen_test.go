package Gen

import (
	"c2p/python/parse"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGen(t *testing.T) {
	tests := []struct {
		name   string
		in     []*parse.Node
		expect string
	}{
		{
			"hello world",
			[]*parse.Node{
				{
					Kind: parse.FunctionDefine, FunctionDefineField: &parse.FunctionDefineField{
						Ident:  parse.NewIdentNode("main"),
						Params: nil,
						Block: parse.NewBlockFieldNode(
							[]*parse.Node{
								parse.NewCallFieldNode(
									parse.NewIdentNode("print"),
									parse.NewPolynomialFieldNode(parse.CallArgs, []*parse.Node{
										parse.NewIdentNode("hello world"),
									})),
							}),
					},
				},
			},
			`def main():
	print("hello world")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gen(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.expect, got); diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}
