package parser

import (
	"compilers/lexer"
	"compilers/token"
	"fmt"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected []Statement
	}{
		// Test "ORIGIN IS ..." statement
		{
			input: "ORIGIN IS (100, 200);",
			expected: []Statement{
				&OriginStatement{
					X: &ConstantExpression{"100"},
					Y: &ConstantExpression{"200"},
				},
			},
		},
		// Test "SCALE IS ..." statement
		{
			input: "SCALE IS (1.5, 2.5);",
			expected: []Statement{
				&ScaleStatement{
					X: &ConstantExpression{"1.5"},
					Y: &ConstantExpression{"2.5"},
				},
			},
		},
		// Test "ROT IS ..." statement
		{
			input: "ROT IS 45;",
			expected: []Statement{
				&RotStatement{
					Angle: &ConstantExpression{"45"},
				},
			},
		},
		// Test assignment statement
		{
			input: "myVar = 100;",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "myVar",
					Value:      &ConstantExpression{"100"},
				},
			},
		},
		// Test "FOR" loop statement
		{
			input: "FOR T FROM 0 TO 120 STEP 1 DRAW (T, 3*T);",
			expected: []Statement{
				&ForStatement{
					LoopVar: "T",
					Start:   &ConstantExpression{"0"},
					End:     &ConstantExpression{"120"},
					Step:    &ConstantExpression{"1"},
					Body: &AssignmentStatement{
						Identifier: "DRAW",
						Value: &BinaryExpression{
							Left:     &ConstantExpression{"T"},
							Operator: token.COMMA,
							Right: &BinaryExpression{
								Left:     &ConstantExpression{"3"},
								Operator: token.MUL,
								Right:    &ConstantExpression{"T"},
							},
						},
					},
				},
			},
		},
		// Test "Sin" function call
		{
			input: "SIN(30);",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "Sin",
					Value: &FunctionCallExpression{
						Name:      "Sin",
						Arguments: []Expression{&ConstantExpression{"30"}},
					},
				},
			},
		},
		// Test "Cos" function call
		{
			input: "COS(45);",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "Cos",
					Value: &FunctionCallExpression{
						Name:      "Cos",
						Arguments: []Expression{&ConstantExpression{"45"}},
					},
				},
			},
		},
		// Test "Tan" function call
		{
			input: "TAN(60);",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "Tan",
					Value: &FunctionCallExpression{
						Name:      "Tan",
						Arguments: []Expression{&ConstantExpression{"60"}},
					},
				},
			},
		},
		// Test "Sqrt" function call
		{
			input: "SQRT(9);",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "Sqrt",
					Value: &FunctionCallExpression{
						Name:      "Sqrt",
						Arguments: []Expression{&ConstantExpression{"9"}},
					},
				},
			},
		},
		// Test "Exp" function call
		{
			input: "EXP(1);",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "Exp",
					Value: &FunctionCallExpression{
						Name:      "Exp",
						Arguments: []Expression{&ConstantExpression{"1"}},
					},
				},
			},
		},
		// Test "Ln" function call
		{
			input: "LN(2);",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "Ln",
					Value: &FunctionCallExpression{
						Name:      "Ln",
						Arguments: []Expression{&ConstantExpression{"2"}},
					},
				},
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Initialize lexer and parser
			l := lexer.New(tt.input)
			parser := New(l)

			// Parse the program and get the resulting statements
			statements := parser.ParseProgram()

			// Check if the parsed statements match the expected output
			if !reflect.DeepEqual(statements, tt.expected) {
				t.Errorf("For input %s, expected statements %+v, but got %+v", tt.input, tt.expected, statements)
				fmt.Printf("Expected: %+v\n", tt.expected)
				fmt.Printf("Got: %+v\n", statements)
			}
		})
	}
}
