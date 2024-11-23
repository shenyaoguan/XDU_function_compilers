package semantic

import (
	"compilers/parser"
	"compilers/token"
	"fmt"
	"git.sr.ht/~sbinet/gg"
	"math"
)

// State 定义当前坐标系的状态
type State struct {
	Variables map[string]float64 // 变量表
	OriginX   float64            // 原点横坐标
	OriginY   float64            // 原点纵坐标
	ScaleX    float64            // 横坐标比例因子
	ScaleY    float64            // 纵坐标比例因子
	Rotation  float64            // 旋转角度，弧度制
}

// NewState 返回一个初始状态
func NewState() *State {
	return &State{
		Variables: make(map[string]float64),
		OriginX:   0,
		OriginY:   0,
		ScaleX:    1,
		ScaleY:    1,
		Rotation:  0,
	}
}

// ApplyScale 应用 SCALE 语句，修改比例因子
func (s *State) ApplyScale(xFactor, yFactor float64) {
	s.ScaleX = xFactor
	s.ScaleY = yFactor
}

// ApplyOrigin 应用 ORIGIN 语句，修改原点
func (s *State) ApplyOrigin(x, y float64) {
	s.OriginX = x
	s.OriginY = y
}

// ApplyRotation 应用 ROT 语句，修改旋转角度
func (s *State) ApplyRotation(angle float64) {
	s.Rotation = angle
}

// TransformPoint 根据当前的坐标变换状态转换一个点的坐标
func (s *State) TransformPoint(x, y float64) (float64, float64) {
	// 先应用缩放
	x *= s.ScaleX
	y *= s.ScaleY

	// 计算旋转后的坐标
	rotX := x*math.Cos(s.Rotation) + y*math.Sin(s.Rotation)
	rotY := y*math.Cos(s.Rotation) + x*math.Sin(s.Rotation)

	// 最后应用平移（ORIGIN）
	rotX += s.OriginX
	rotY += s.OriginY

	return rotX, rotY
}

// FunctionCall 对应的函数调用，例如 Sin, Cos 等
type FunctionCall struct {
	Name      string
	Arguments []float64
}

// ApplyFunction 应用函数调用（例如 Sin, Cos）
func (s *State) ApplyFunction(fn string, args []float64) []float64 {
	var result []float64
	switch fn {
	case "Sin":
		for _, arg := range args {
			result = append(result, math.Sin(arg))
		}
	case "Cos":
		for _, arg := range args {
			result = append(result, math.Cos(arg))
		}
	case "Tan":
		for _, arg := range args {
			result = append(result, math.Tan(arg))
		}
	case "Sqrt":
		for _, arg := range args {
			result = append(result, math.Sqrt(arg))
		}
	case "Exp":
		for _, arg := range args {
			result = append(result, math.Exp(arg))
		}
	case "Ln":
		for _, arg := range args {
			result = append(result, math.Log(arg))
		}
	}
	return result
}

// ForStatement 表示 FOR 循环语句
type ForStatement struct {
	LoopVar  string     // 循环变量名
	Start    float64    // 起点
	End      float64    // 终点
	Step     float64    // 步长
	DrawExpr Expression // 表达式，生成 (横坐标, 纵坐标)
}

// Expression 表达式接口，所有表达式实现此接口
type Expression interface {
	Evaluate(t float64) []float64 // 计算表达式的值，t 作为动态参数
}

// ConstantExpression 常量表达式
type ConstantExpression struct {
	Value float64
}

// Evaluate 计算常量表达式的值
func (c *ConstantExpression) Evaluate(t float64) []float64 {
	return []float64{c.Value}
}

// BinaryExpression 二元表达式
type BinaryExpression struct {
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

// Evaluate 计算二元表达式的值
func (b *BinaryExpression) Evaluate(t float64) []float64 {
	left := b.Left.Evaluate(t)[0]
	right := b.Right.Evaluate(t)[0]
	var result float64
	switch b.Operator {
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
}

// FunctionCallExpression 函数调用表达式（例如 Sin(T), Cos(T)）
type FunctionCallExpression struct {
	FunctionName string
	Arguments    []Expression
}

// Evaluate 计算函数调用的值，t 作为参数传入
func (f *FunctionCallExpression) Evaluate(t float64) []float64 {
	var args []float64
	for _, arg := range f.Arguments {
		args = append(args, arg.Evaluate(t)[0]) // 获取每个参数的值
	}

	// 调用函数
	return applyFunction(f.FunctionName, args)
}

// applyFunction 应用函数（例如 Sin, Cos）
func applyFunction(fn string, args []float64) []float64 {
	var result []float64
	switch fn {
	case "Sin":
		for _, arg := range args {
			result = append(result, math.Sin(arg))
		}
	case "Cos":
		for _, arg := range args {
			result = append(result, math.Cos(arg))
		}
	}
	return result
}

// ParseForStatement 解析 FOR T FROM 起点 TO 终点 STEP 步长 DRAW (横坐标, 纵坐标)
func (s *State) ParseForStatement(start, end, step float64, drawExpr parser.Expression) {
	const width = 800
	const height = 600

	// Create a new image context
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	// Loop from start to end, incrementing by step
	for t := start; t <= end; t += step {
		// Evaluate the DRAW expression
		result := drawExpr.Evaluate(t, s.Variables)
		x, y := result[0], result[1]

		// Transform the point according to the current state
		transformedX, transformedY := s.TransformPoint(x, y)

		// Draw the point
		dc.DrawPoint(transformedX, transformedY, 2)
		fmt.Println("Drawing point:", transformedX, transformedY)
		dc.Fill()
	}

	// Save the image to a file
	err := dc.SavePNG("output.png")
	if err != nil {
		fmt.Println("Failed to save image:", err)
	}
}
