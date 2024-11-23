package interpreter

import (
	"compilers/parser"
	"compilers/semantic"
	"compilers/token"
	"fmt"
	"math"
	"strconv"
)

// Interpreter 解释器结构体，负责执行解析的语法树
type Interpreter struct {
	state  *semantic.State
	parser *parser.Parser
}

// NewInterpreter 创建一个新的解释器实例
func NewInterpreter(p *parser.Parser) *Interpreter {
	return &Interpreter{
		state:  semantic.NewState(),
		parser: p,
	}
}

// Interpret 执行程序
func (i *Interpreter) Interpret() {
	statements := i.parser.ParseProgram()
	for _, stmt := range statements {
		i.executeStatement(stmt)
	}
}

// 执行语句
func (i *Interpreter) executeStatement(stmt parser.Statement) {
	switch stmt := stmt.(type) {
	case *parser.OriginStatement:
		i.executeOriginStatement(stmt)
	case *parser.ScaleStatement:
		i.executeScaleStatement(stmt)
	case *parser.RotStatement:
		i.executeRotStatement(stmt)
	case *parser.AssignmentStatement:
		i.executeAssignmentStatement(stmt)
	case *parser.ForStatement:
		i.executeForStatement(stmt)
	}
}

// 执行 ORIGIN 语句
func (i *Interpreter) executeOriginStatement(stmt *parser.OriginStatement) {
	// 计算表达式的值
	x := i.evaluateExpression(stmt.X)
	y := i.evaluateExpression(stmt.Y)
	// 更新坐标系的原点
	i.state.ApplyOrigin(x[0], y[0])
	fmt.Printf("Origin set to: (%f, %f)\n", x[0], y[0])
}

// 执行 SCALE 语句
func (i *Interpreter) executeScaleStatement(stmt *parser.ScaleStatement) {
	// 计算表达式的值
	x := i.evaluateExpression(stmt.X)
	y := i.evaluateExpression(stmt.Y)
	// 更新比例因子
	i.state.ApplyScale(x[0], y[0])
	fmt.Printf("Scale set to: (%f, %f)\n", x[0], y[0])
}

// 执行 ROT 语句
func (i *Interpreter) executeRotStatement(stmt *parser.RotStatement) {
	// 计算角度的值
	angle := i.evaluateExpression(stmt.Angle)
	angle1 := angle[0] * 180 / math.Pi
	// 更新旋转角度
	i.state.ApplyRotation(angle[0])
	fmt.Printf("Rotation set to: %f radians\n", angle1)
}

// 执行赋值语句
func (i *Interpreter) executeAssignmentStatement(stmt *parser.AssignmentStatement) {
	// 计算右侧表达式的值
	value := i.evaluateExpression(stmt.Value)
	// 输出结果
	fmt.Printf("Assignment: %s = %v\n", stmt.Identifier, value)
}

// 执行 FOR 语句
func (i *Interpreter) executeForStatement(stmt *parser.ForStatement) {
	// 计算 FOR 语句中的起始值、终止值和步长
	start := i.evaluateExpression(stmt.Start)[0]
	end := i.evaluateExpression(stmt.End)[0]
	step := i.evaluateExpression(stmt.Step)[0]
	var drawExpr parser.Expression
	// 计算绘制表达式
	if assignStmt, ok := stmt.Body.(*parser.AssignmentStatement); ok {
		drawExpr = assignStmt.Value
	} else {
		panic("Expected AssignmentStatement in FOR loop body")
	}

	// 执行循环
	i.state.ParseForStatement(start, end, step, drawExpr)
}

// 计算表达式的值
func (i *Interpreter) evaluateExpression(expr parser.Expression) []float64 {
	switch expr := expr.(type) {
	case *parser.ConstantExpression:
		// 常量表达式
		if expr.Value == "PI" {
			return []float64{math.Pi}
		} else if expr.Value == "E" {
			return []float64{math.E}
		}
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse float: %v", err))
		}
		return []float64{val}
	case *parser.BinaryExpression:
		// 二元表达式
		left := i.evaluateExpression(expr.Left)[0]
		right := i.evaluateExpression(expr.Right)[0]
		var result float64
		switch expr.Operator {
		case token.PLUS:
			result = left + right
		case token.MINUS:
			result = left - right
		case token.MUL:
			result = left * right
		case token.DIV:
			result = left / right
		}
		return []float64{result}
	case *parser.FunctionCallExpression:
		// 函数调用表达式
		var args []float64
		for _, arg := range expr.Arguments {
			args = append(args, i.evaluateExpression(arg)[0])
		}
		// 调用函数
		return i.state.ApplyFunction(expr.Name, args)
	case *parser.VariableExpression:
		if val, ok := i.state.Variables[expr.Name]; ok {
			return []float64{val}
		}
		panic(fmt.Sprintf("Undefined variable: %v", expr.Name))
	default:
		// 错误处理
		panic("Unknown expression type")
	}
}
