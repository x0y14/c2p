package c2p

import (
	cparse "c2p/c/parse"
	"c2p/c/tokenize"
	pparse "c2p/python/parse"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestCToP(t *testing.T) {
	tests := []struct {
		name   string
		in     []*cparse.Node
		expect []*pparse.Node
	}{
		{
			"return",
			[]*cparse.Node{
				// int main(void) {
				//     return 0
				// }
				cparse.NewFunctionDefineFieldNode(
					cparse.NewIdentNode("int"),
					cparse.NewIdentNode("main"),
					cparse.NewIdentNode("void"),
					cparse.NewBlockFieldNode([]*cparse.Node{
						cparse.NewReturnFieldNode(
							cparse.NewLiteralFieldNode(cparse.LInt, 0, 0, "")),
					}),
				),
			},
			[]*pparse.Node{
				// def main():
				//     return 0
				pparse.NewFunctionDefineFieldNode(
					pparse.NewIdentNode("main"),
					nil,
					pparse.NewBlockFieldNode([]*pparse.Node{
						pparse.NewReturnFieldNode([]*pparse.Node{
							pparse.NewLiteralFieldNode(pparse.LInt, 0, 0, ""),
						}),
					}),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CToP(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.expect, got); diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}

func TestCToPFizzBuzz(t *testing.T) {
	fizzbuzz := `
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
`
	tokens, err := tokenize.Tokenize(fizzbuzz)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := cparse.Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}
	got, err := CToP(nodes)
	if err != nil {
		t.Fatal(err)
	}

	expect := []*pparse.Node{
		pparse.NewFunctionDefineFieldNode(
			pparse.NewIdentNode("main"),
			nil,
			pparse.NewBlockFieldNode([]*pparse.Node{
				pparse.NewWhileFieldNode(
					pparse.NewAssignFieldNode( // i = 1
						pparse.NewIdentNode("i"),
						pparse.NewIntNode(1),
					),
					pparse.NewBinaryFieldNode( // i <= 100
						pparse.Le,
						pparse.NewIdentNode("i"),
						pparse.NewIntNode(100),
					),
					pparse.NewAssignFieldNode( // i = i+1
						pparse.NewIdentNode("i"),
						pparse.NewBinaryFieldNode(
							pparse.Add,
							pparse.NewIdentNode("i"),
							pparse.NewIntNode(1),
						),
					),
					pparse.NewBlockFieldNode([]*pparse.Node{
						pparse.NewIfElseFieldNode( // (i % 3 == 0 && i % 5 == 0)
							pparse.NewBinaryFieldNode(
								pparse.And,
								pparse.NewBinaryFieldNode( // i % 3 == 0
									pparse.Eq,
									pparse.NewBinaryFieldNode(
										pparse.Mod,
										pparse.NewIdentNode("i"),
										pparse.NewIntNode(3),
									),
									pparse.NewIntNode(0),
								),
								pparse.NewBinaryFieldNode( // i % 5 == 0
									pparse.Eq,
									pparse.NewBinaryFieldNode(
										pparse.Mod,
										pparse.NewIdentNode("i"),
										pparse.NewIntNode(5),
									),
									pparse.NewIntNode(0),
								),
							),
							pparse.NewBlockFieldNode([]*pparse.Node{
								// printf("FizzBuzz\n");
								pparse.NewCallFieldNode(
									pparse.NewIdentNode("printf"),
									pparse.NewPolynomialFieldNode(
										pparse.CallArgs,
										[]*pparse.Node{
											pparse.NewStringNode("FizzBuzz\n"),
										},
									),
								),
							}),
							pparse.NewIfElseFieldNode( // (i % 3 == 0)
								pparse.NewBinaryFieldNode(
									pparse.Eq,
									pparse.NewBinaryFieldNode(
										pparse.Mod,
										pparse.NewIdentNode("i"),
										pparse.NewIntNode(3),
									),
									pparse.NewIntNode(0),
								),
								pparse.NewBlockFieldNode([]*pparse.Node{
									// printf("Fizz\n");
									pparse.NewCallFieldNode(
										pparse.NewIdentNode("printf"),
										pparse.NewPolynomialFieldNode(
											pparse.CallArgs,
											[]*pparse.Node{
												pparse.NewStringNode("Fizz\n"),
											},
										),
									),
								}),
								pparse.NewIfElseFieldNode( // i % 5 == 0
									pparse.NewBinaryFieldNode(
										pparse.Eq,
										pparse.NewBinaryFieldNode(
											pparse.Mod,
											pparse.NewIdentNode("i"),
											pparse.NewIntNode(5),
										),
										pparse.NewIntNode(0),
									),
									pparse.NewBlockFieldNode([]*pparse.Node{
										pparse.NewCallFieldNode(
											pparse.NewIdentNode("printf"),
											pparse.NewPolynomialFieldNode(
												pparse.CallArgs,
												[]*pparse.Node{
													pparse.NewStringNode("Buzz\n"),
												},
											),
										),
									}),
									pparse.NewBlockFieldNode([]*pparse.Node{
										pparse.NewCallFieldNode(
											pparse.NewIdentNode("printf"),
											pparse.NewPolynomialFieldNode(
												pparse.CallArgs,
												[]*pparse.Node{
													pparse.NewStringNode("%d\n"),
													pparse.NewIdentNode("i"),
												},
											),
										),
									}),
								),
							),
						),
					}),
				),
				pparse.NewReturnFieldNode([]*pparse.Node{
					pparse.NewIntNode(0),
				}),
			}),
		),
	}

	if diff := cmp.Diff(expect, got); diff != "" {
		t.Errorf("%v", diff)
	}
}
