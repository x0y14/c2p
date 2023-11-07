package Gen

import (
	"c2p/python/parse"
	"fmt"
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
			fmt.Println(got)
			if diff := cmp.Diff(tt.expect, got); diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}

func TestGen2(t *testing.T) {
	fizzbuzz := []*parse.Node{
		parse.NewFunctionDefineFieldNode(
			parse.NewIdentNode("main"),
			nil,
			parse.NewBlockFieldNode([]*parse.Node{
				parse.NewWhileFieldNode(
					parse.NewAssignFieldNode( // i = 1
						parse.NewIdentNode("i"),
						parse.NewIntNode(1),
					),
					parse.NewBinaryFieldNode( // i <= 100
						parse.Le,
						parse.NewIdentNode("i"),
						parse.NewIntNode(100),
					),
					parse.NewAssignFieldNode( // i = i+1
						parse.NewIdentNode("i"),
						parse.NewBinaryFieldNode(
							parse.Add,
							parse.NewIdentNode("i"),
							parse.NewIntNode(1),
						),
					),
					parse.NewBlockFieldNode([]*parse.Node{
						parse.NewIfElseFieldNode( // (i % 3 == 0 && i % 5 == 0)
							parse.NewBinaryFieldNode(
								parse.And,
								parse.NewBinaryFieldNode( // i % 3 == 0
									parse.Eq,
									parse.NewBinaryFieldNode(
										parse.Mod,
										parse.NewIdentNode("i"),
										parse.NewIntNode(3),
									),
									parse.NewIntNode(0),
								),
								parse.NewBinaryFieldNode( // i % 5 == 0
									parse.Eq,
									parse.NewBinaryFieldNode(
										parse.Mod,
										parse.NewIdentNode("i"),
										parse.NewIntNode(5),
									),
									parse.NewIntNode(0),
								),
							),
							parse.NewBlockFieldNode([]*parse.Node{
								// printf("FizzBuzz\n");
								parse.NewCallFieldNode(
									parse.NewIdentNode("printf"),
									parse.NewPolynomialFieldNode(
										parse.CallArgs,
										[]*parse.Node{
											parse.NewStringNode("FizzBuzz\n"),
										},
									),
								),
							}),
							parse.NewIfElseFieldNode( // (i % 3 == 0)
								parse.NewBinaryFieldNode(
									parse.Eq,
									parse.NewBinaryFieldNode(
										parse.Mod,
										parse.NewIdentNode("i"),
										parse.NewIntNode(3),
									),
									parse.NewIntNode(0),
								),
								parse.NewBlockFieldNode([]*parse.Node{
									// printf("Fizz\n");
									parse.NewCallFieldNode(
										parse.NewIdentNode("printf"),
										parse.NewPolynomialFieldNode(
											parse.CallArgs,
											[]*parse.Node{
												parse.NewStringNode("Fizz\n"),
											},
										),
									),
								}),
								parse.NewIfElseFieldNode( // i % 5 == 0
									parse.NewBinaryFieldNode(
										parse.Eq,
										parse.NewBinaryFieldNode(
											parse.Mod,
											parse.NewIdentNode("i"),
											parse.NewIntNode(5),
										),
										parse.NewIntNode(0),
									),
									parse.NewBlockFieldNode([]*parse.Node{
										parse.NewCallFieldNode(
											parse.NewIdentNode("printf"),
											parse.NewPolynomialFieldNode(
												parse.CallArgs,
												[]*parse.Node{
													parse.NewStringNode("Buzz\n"),
												},
											),
										),
									}),
									parse.NewBlockFieldNode([]*parse.Node{
										parse.NewCallFieldNode(
											parse.NewIdentNode("printf"),
											parse.NewPolynomialFieldNode(
												parse.CallArgs,
												[]*parse.Node{
													parse.NewStringNode("%d\n"),
													parse.NewIdentNode("i"),
												},
											),
										),
									}),
								),
							),
						),
					}),
				),
				parse.NewReturnFieldNode([]*parse.Node{
					parse.NewIntNode(0),
				}),
			}),
		),
	}
	source, err := Gen(fizzbuzz)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(source)
}
