package token

// TokenType 是标记的类型
type TokenType string

// Token 表示一个词法单元
type Token struct {
	Type    TokenType // 标记类型
	Literal string    // 标记的字面值
}

// List of token types.
const (
	// Special Tokens
	ILLEGAL TokenType = "ILLEGAL" // Invalid token
	EOF     TokenType = "EOF"     // End of file

	// Identifiers and Literals
	ID       TokenType = "ID"       // Identifier
	CONST_ID TokenType = "CONST_ID" // Constant (e.g., numbers)

	// Operators
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	MUL    TokenType = "*"
	DIV    TokenType = "/"
	POWER  TokenType = "**"
	ASSIGN TokenType = "="

	// Delimiters
	COMMA     TokenType = ","
	SEMICO    TokenType = ";"
	L_BRACKET TokenType = "("
	R_BRACKET TokenType = ")"

	// Keywords
	ORIGIN TokenType = "ORIGIN"
	IS     TokenType = "IS"
	SCALE  TokenType = "SCALE"
	ROT    TokenType = "ROT"
	FOR    TokenType = "FOR"
	FROM   TokenType = "FROM"
	TO     TokenType = "TO"
	STEP   TokenType = "STEP"
	DRAW   TokenType = "DRAW"

	// Comments
	COMMENT TokenType = "COMMENT"
)

// New 创建一个新的 Token
func New(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}
