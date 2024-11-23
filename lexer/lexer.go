package lexer

import (
	"compilers/token"
	"unicode"
)

// Lexer represents the lexical analyzer.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
}

// New creates a new Lexer instance.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character from the input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 represents EOF
	} else {
		l.ch = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken lexes the next token and advances the lexer state.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '/':
		if l.peekChar() == '/' { // 处理注释 //
			tok = token.New(token.COMMENT, l.readComment())
		} else { // 单独的 /
			tok = token.New(token.DIV, string(l.ch))
		}
	case '-':
		if l.peekChar() == '-' { // 处理注释 --
			tok = token.New(token.COMMENT, l.readComment())
		} else {
			tok = token.New(token.MINUS, string(l.ch))
		}
	case '+':
		tok = token.New(token.PLUS, string(l.ch))
	case '*':
		if l.peekChar() == '*' { // 处理幂运算符 **
			l.readChar()
			tok = token.New(token.POWER, "**")
		} else {
			tok = token.New(token.MUL, string(l.ch))
		}
	case '=':
		tok = token.New(token.ASSIGN, string(l.ch))
	case ',':
		tok = token.New(token.COMMA, string(l.ch))
	case ';':
		tok = token.New(token.SEMICO, string(l.ch))
	case '(':
		tok = token.New(token.L_BRACKET, string(l.ch))
	case ')':
		tok = token.New(token.R_BRACKET, string(l.ch))
	case 0:
		tok = token.New(token.EOF, "")
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			tokenType := lookupKeyword(ident)
			tok = token.New(tokenType, ident)
			return tok
		} else if isDigit(l.ch) {
			value := l.readNumber()
			tok = token.New(token.CONST_ID, value)
			return tok
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace skips over white space characters.
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// readComment reads the entire comment (either // or --).
func (l *Lexer) readComment() string {
	start := l.position
	for l.ch != '\n' && l.ch != 0 { // 注释内容直到换行符或文件结束
		l.readChar()
	}
	return l.input[start:l.position]
}

// readIdentifier reads an identifier (variable, function name, etc.).
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readNumber reads a number (integer or float).
// readNumber reads a number (integer or float), handling cases like "100/3".
func (l *Lexer) readNumber() string {
	start := l.position
	// First, read the integer part
	for isDigit(l.ch) {
		l.readChar()
	}
	// Handle fractional part
	if l.ch == '.' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	// Handle division operator within numbers
	if l.ch == '/' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[start:l.position]
}

// peekChar returns the next character without advancing the position.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPosition])
}

// isLetter checks if a character is a letter.
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isDigit checks if a character is a digit.
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// lookupKeyword checks if an identifier is a keyword.
func lookupKeyword(ident string) token.TokenType {
	// Check for specific constants like PI, E, etc.
	constants := map[string]token.TokenType{
		"PI": token.CONST_ID,
		"E":  token.CONST_ID,
		// Add any other predefined constants here
	}

	if tok, ok := constants[ident]; ok {
		return tok
	}

	// Check for keywords
	keywords := map[string]token.TokenType{
		"ORIGIN": token.ORIGIN,
		"IS":     token.IS,
		"SCALE":  token.SCALE,
		"ROT":    token.ROT,
		"FOR":    token.FOR,
		"FROM":   token.FROM,
		"TO":     token.TO,
		"STEP":   token.STEP,
		"DRAW":   token.DRAW,
	}

	if tok, ok := keywords[ident]; ok {
		return tok
	}

	// Default to ID for generic identifiers
	return token.ID
}
