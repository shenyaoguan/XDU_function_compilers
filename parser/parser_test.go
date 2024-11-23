package parser

import (
	"compilers/lexer"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected []Statement
	}{
		{
			// Test "ORIGIN IS ..." statement
			input: "ORIGIN IS (100, 200);",
			expected: []Statement{
				&OriginStatement{
					X: &ConstantExpression{"100"},
					Y: &ConstantExpression{"200"},
				},
			},
		},
		{
			// Test "SCALE IS ..." statement
			input: "SCALE IS (1.5, 2.5);",
			expected: []Statement{
				&ScaleStatement{
					X: &ConstantExpression{"1.5"},
					Y: &ConstantExpression{"2.5"},
				},
			},
		},
		{
			// Test "ROT IS ..." statement
			input: "ROT IS 45;",
			expected: []Statement{
				&RotStatement{
					Angle: &ConstantExpression{"45"},
				},
			},
		},
		{
			// Test assignment statement
			input: "myVar = 100;",
			expected: []Statement{
				&AssignmentStatement{
					Identifier: "myVar",
					Value:      &ConstantExpression{"100"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			parser := New(l)

			statements := parser.ParseProgram()

			if !reflect.DeepEqual(statements, tt.expected) {
				t.Errorf("For input %s, expected statements %+v, but got %+v", tt.input, tt.expected, statements)
			}
		})
	}
}
