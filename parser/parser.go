package parser

import (
	"compilers/lexer"
	"compilers/token"
	"fmt"
	"math"
	"strconv"
)

// Statement is an interface for all statement types
type Statement interface {
}

// Expression is an interface for all expression types
type Expression interface {
	Evaluate(t float64, variables map[string]float64) []float64
}

// VariableExpression represents a variable in an expression
type VariableExpression struct {
	Name string
}

func (v *VariableExpression) Evaluate(t float64, variables map[string]float64) []float64 {
	if val, ok := variables[v.Name]; ok {
		return []float64{val}
	}
	panic(fmt.Sprintf("Undefined variable: %v", v.Name))
}

// ConstantExpression represents a constant value in an expression
type ConstantExpression struct {
	Value string
}

func (c *ConstantExpression) Evaluate(t float64, variables map[string]float64) []float64 {
	if c.Value == "PI" {
		return []float64{math.Pi}
	} else if c.Value == "E" {
		return []float64{math.E}
	} else if c.Value == "T" {
		return []float64{t}
	}
	val, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert constant expression to float: %v", c.Value))
	}
	return []float64{val}
}

// BinaryExpression represents a binary operation in an expression
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator token.TokenType
}

func (b *BinaryExpression) Evaluate(t float64, variables map[string]float64) []float64 {
	left := b.Left.Evaluate(t, variables)[0]
	right := b.Right.Evaluate(t, variables)[0]
	var result []float64
	switch b.Operator {
	case token.PLUS:
		result = append(result, left+right)
	case token.MINUS:
		result = append(result, left-right)
	case token.MUL:
		result = append(result, left*right)
	case token.DIV:
		result = append(result, left/right)
	case token.COMMA:
		result = append(result, left, right)
	}
	return result
}

// FunctionCallExpression represents a function call in an expression
type FunctionCallExpression struct {
	Name      string
	Arguments []Expression
}

func (f *FunctionCallExpression) Evaluate(t float64, variables map[string]float64) []float64 {
	var args []float64
	for _, arg := range f.Arguments {
		args = append(args, arg.Evaluate(t, variables)[0])
	}
	return applyFunction(f.Name, args)
}

func applyFunction(name string, args []float64) []float64 {
	var result []float64
	switch name {
	case "SIN":
		for _, arg := range args {
			result = append(result, math.Sin(arg))
		}
	case "COS":
		for _, arg := range args {
			result = append(result, math.Cos(arg))
		}
	case "TAN":
		for _, arg := range args {
			result = append(result, math.Tan(arg))
		}
	case "SQRT":
		for _, arg := range args {
			result = append(result, math.Sqrt(arg))
		}
	case "EXP":
		for _, arg := range args {
			result = append(result, math.Exp(arg))
		}
	case "LN":
		for _, arg := range args {
			result = append(result, math.Log(arg))
		}
	default:
		panic(fmt.Sprintf("Unknown function: %s", name))
	}
	return result
}

// 语句类型
type OriginStatement struct {
	X Expression
	Y Expression
}

type ScaleStatement struct {
	X Expression
	Y Expression
}

type RotStatement struct {
	Angle Expression
}

type AssignmentStatement struct {
	Identifier string
	Value      Expression
}

type ForStatement struct {
	LoopVar string
	Start   Expression
	End     Expression
	Step    Expression
	Body    Statement
}

type CommentStatement struct {
}

// Parser 结构体用于解析输入
type Parser struct {
	lexer    *lexer.Lexer
	curToken token.Token
}

// New 创建一个新的语法分析器
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	return p
}

// nextToken 移动到下一个 token
func (p *Parser) nextToken() {
	p.curToken = p.lexer.NextToken()
}

// ParseProgram 解析程序
func (p *Parser) ParseProgram() []Statement {
	var statements []Statement
	for p.curToken.Type != token.EOF {
		statements = append(statements, p.parseStatement())
		if p.curToken.Type == token.SEMICO || p.curToken.Type == token.COMMENT {
			p.nextToken()
		} else {
			break
		}
	}
	return statements
}

// parseStatement 解析一个语句
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case token.ORIGIN:
		return p.parseOriginStatement()
	case token.SCALE:
		return p.parseScaleStatement()
	case token.ROT:
		return p.parseRotStatement()
	case token.ID:
		return p.parseAssignmentStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.TAN, token.SIN, token.COS, token.SQRT, token.EXP, token.LN:
		return p.parseFunctionCall()
	case token.COMMENT:
		return p.parseCommentStatement()
	default:
		p.error("Unexpected token in statement: " + p.curToken.Literal)
		return nil
	}
}

// parseOriginStatement 解析 ORIGIN IS 语句
func (p *Parser) parseOriginStatement() *OriginStatement {
	p.nextToken() // skip ORIGIN
	p.expect(token.IS)
	p.expect(token.L_BRACKET)
	x := p.parseExpression()
	p.expect(token.COMMA)
	y := p.parseExpression()
	p.expect(token.R_BRACKET)
	return &OriginStatement{X: x, Y: y}
}

// parseScaleStatement 解析 SCALE IS 语句
func (p *Parser) parseScaleStatement() *ScaleStatement {
	p.nextToken() // skip SCALE
	p.expect(token.IS)
	p.expect(token.L_BRACKET)
	x := p.parseExpression()
	p.expect(token.COMMA)
	y := p.parseExpression()
	p.expect(token.R_BRACKET)
	return &ScaleStatement{X: x, Y: y}
}

// parseRotStatement 解析 ROT IS 语句
func (p *Parser) parseRotStatement() *RotStatement {
	p.nextToken() // skip ROT
	p.expect(token.IS)
	angle := p.parseExpression()
	return &RotStatement{Angle: angle}
}

// parseForStatement 解析 FOR 语句
func (p *Parser) parseForStatement() *ForStatement {
	p.nextToken() // skip FOR

	loopVar := p.curToken.Literal
	p.nextToken() // skip loop variable

	p.expect(token.FROM)
	start := p.parseExpression()

	p.expect(token.TO)
	end := p.parseExpression()

	p.expect(token.STEP)
	step := p.parseExpression()

	p.expect(token.DRAW)
	body := p.parseDrawStatement()

	return &ForStatement{
		LoopVar: loopVar,
		Start:   start,
		End:     end,
		Step:    step,
		Body:    body,
	}
}

func (p *Parser) parseDrawStatement() *AssignmentStatement {
	identifier := "DRAW"
	p.expect(token.L_BRACKET)
	value1 := p.parseExpression()
	p.expect(token.COMMA)
	value2 := p.parseExpression()
	value := &BinaryExpression{
		Left:     value1,
		Operator: token.COMMA,
		Right:    value2,
	}
	p.expect(token.R_BRACKET)
	return &AssignmentStatement{Identifier: identifier, Value: value}
}

func (p *Parser) parseCommentStatement() *CommentStatement {
	return &CommentStatement{}
}

// parseAssignmentStatement 解析赋值语句
func (p *Parser) parseAssignmentStatement() *AssignmentStatement {
	identifier := p.curToken.Literal
	p.nextToken() // skip identifier

	// 确保当前 token 是赋值操作符 "="
	p.expect(token.ASSIGN)

	// 解析右侧的表达式
	value := p.parseExpression()

	return &AssignmentStatement{Identifier: identifier, Value: value}
}

// parseExpression 解析表达式
func (p *Parser) parseExpression() Expression {
	left := p.parseTerm()

	// Check for operators in the form of Expression PLUS Term | Expression MINUS Term
	for p.curToken.Type == token.PLUS || p.curToken.Type == token.MINUS {
		operator := p.curToken.Type
		p.nextToken()
		right := p.parseTerm()
		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
		p.nextToken()
	}
	return left
}

// parseTerm 解析乘法和除法
func (p *Parser) parseTerm() Expression {
	left := p.parseFactor()
	p.nextToken()

	// Check for MUL or DIV operators
	for p.curToken.Type == token.MUL || p.curToken.Type == token.DIV {
		operator := p.curToken.Type
		p.nextToken()
		right := p.parseFactor()
		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
		p.nextToken()
	}
	return left
}

// parseFactor 解析一元加法/减法和原子
func (p *Parser) parseFactor() Expression {
	switch p.curToken.Type {
	case token.PLUS, token.MINUS:
		operator := p.curToken.Type
		p.nextToken()
		factor := p.parseComponent()
		return &BinaryExpression{
			Left:     &ConstantExpression{"0"}, // 表示 "+Factor" 或 "-Factor"
			Operator: operator,
			Right:    factor,
		}
	default:
		return p.parseComponent()
	}
}

// parseComponent 解析原子表达式（包括函数调用）
func (p *Parser) parseComponent() Expression {
	// 处理数字常量或标识符（变量）
	switch p.curToken.Type {
	case token.CONST_ID:
		return &ConstantExpression{Value: p.curToken.Literal}
	case token.ID:
		// 判断当前 ID 是否是已知函数名
		if isFunction(p.curToken.Literal) {
			return p.parseFunctionCall()
		}
		// 处理普通的标识符（变量名）
		return &ConstantExpression{Value: p.curToken.Literal}
	case token.L_BRACKET:
		// 解析括号内的表达式
		p.nextToken()
		expr := p.parseExpression()
		p.expect(token.R_BRACKET)
		return expr
	case token.TAN, token.SIN, token.COS, token.SQRT, token.EXP, token.LN:
		return p.parseFunctionCall()
	default:
		p.error("Unexpected token in component: " + p.curToken.Literal)
		return nil
	}
}

// isFunction 检查给定的标识符是否为已知的函数名
func isFunction(name string) bool {
	switch name {
	case "SIN", "COS", "TAN", "SQRT", "EXP", "LN":
		return true
	default:
		return false
	}
}

// parseFunctionCall 解析函数调用表达式
func (p *Parser) parseFunctionCall() *FunctionCallExpression {
	funcName := p.curToken.Literal
	p.nextToken() // skip function name
	p.expect(token.L_BRACKET)

	var arguments []Expression
	for p.curToken.Type != token.R_BRACKET {
		arguments = append(arguments, p.parseExpression())
		if p.curToken.Type == token.COMMA {
			p.nextToken() // skip comma
		}
	}

	//p.expect(token.R_BRACKET)
	return &FunctionCallExpression{
		Name:      funcName,
		Arguments: arguments,
	}
}

// expect 检查当前 token 类型是否匹配，如果不匹配则报错
func (p *Parser) expect(t token.TokenType) {
	if p.curToken.Type == t {
		p.nextToken()
	} else {
		p.error(fmt.Sprintf("Expected %s, got %s", t, p.curToken.Type))
	}
}

// error 报告语法错误
func (p *Parser) error(msg string) {
	panic(msg)
}
