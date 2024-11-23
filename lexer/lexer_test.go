package lexer

import (
	"compilers/token"
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []token.Token
	}{
		{
			// Basic arithmetic and control flow test
			input: "ORIGIN IS (100, 100); SCALE IS (100, 100/3); ROT IS PI/4;",
			expected: []token.Token{
				{Type: token.ORIGIN, Literal: "ORIGIN"},
				{Type: token.IS, Literal: "IS"},
				{Type: token.L_BRACKET, Literal: "("},
				{Type: token.CONST_ID, Literal: "100"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.CONST_ID, Literal: "100"},
				{Type: token.R_BRACKET, Literal: ")"},
				{Type: token.SEMICO, Literal: ";"},
				{Type: token.SCALE, Literal: "SCALE"},
				{Type: token.IS, Literal: "IS"},
				{Type: token.L_BRACKET, Literal: "("},
				{Type: token.CONST_ID, Literal: "100"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.CONST_ID, Literal: "100/3"},
				{Type: token.R_BRACKET, Literal: ")"},
				{Type: token.SEMICO, Literal: ";"},
				{Type: token.ROT, Literal: "ROT"},
				{Type: token.IS, Literal: "IS"},
				{Type: token.CONST_ID, Literal: "PI"},
				{Type: token.DIV, Literal: "/"},
				{Type: token.CONST_ID, Literal: "4"},
				{Type: token.SEMICO, Literal: ";"},
			},
		},
		{
			// Test comments
			input: "// This is a comment\nORIGIN IS (0, 0);",
			expected: []token.Token{
				{Type: token.COMMENT, Literal: "// This is a comment"},
				{Type: token.ORIGIN, Literal: "ORIGIN"},
				{Type: token.IS, Literal: "IS"},
				{Type: token.L_BRACKET, Literal: "("},
				{Type: token.CONST_ID, Literal: "0"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.CONST_ID, Literal: "0"},
				{Type: token.R_BRACKET, Literal: ")"},
				{Type: token.SEMICO, Literal: ";"},
			},
		},
		{
			// Test operators and keywords
			input: "+ - * / ** =",
			expected: []token.Token{
				{Type: token.PLUS, Literal: "+"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.MUL, Literal: "*"},
				{Type: token.DIV, Literal: "/"},
				{Type: token.POWER, Literal: "**"},
				{Type: token.ASSIGN, Literal: "="},
			},
		},
		{
			// Test identifier recognition
			input: "myVar PI someVar",
			expected: []token.Token{
				{Type: token.ID, Literal: "myVar"},
				{Type: token.CONST_ID, Literal: "PI"},
				{Type: token.ID, Literal: "someVar"},
			},
		},
		{
			// Test numeric values
			input: "123 45.67 0.89",
			expected: []token.Token{
				{Type: token.CONST_ID, Literal: "123"},
				{Type: token.CONST_ID, Literal: "45.67"},
				{Type: token.CONST_ID, Literal: "0.89"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := New(tt.input)
			var tokens []token.Token
			for {
				tok := lexer.NextToken()
				if tok.Type == token.EOF {
					break
				}
				tokens = append(tokens, tok)
			}

			if !reflect.DeepEqual(tokens, tt.expected) {
				t.Errorf("For input %s, expected tokens %+v, but got %+v", tt.input, tt.expected, tokens)
			}
		})
	}
}
