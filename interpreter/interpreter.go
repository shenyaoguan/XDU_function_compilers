package interpreter

import (
	"compilers/parser"
	"compilers/token"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"os"
)

// Interpreter 用于执行语法树并计算语义
type Interpreter struct {
	variables map[string]float64
}

// NewInterpreter 创建一个新的解释器
func NewInterpreter() *Interpreter {
	return &Interpreter{
		variables: make(map[string]float64),
	}
}

// RunProgram 执行语法树中的所有语句
func (i *Interpreter) RunProgram(statements []parser.Statement) {
	for _, stmt := range statements {
		switch stmt := stmt.(type) {
		case *parser.OriginStatement:
			// 处理 ORIGIN 语句
			fmt.Printf("Origin: (%s, %s)\n", stmt.X, stmt.Y)
		case *parser.ScaleStatement:
			// 处理 SCALE 语句
			fmt.Printf("Scale: (%s, %s)\n", stmt.X, stmt.Y)
		case *parser.RotStatement:
			// 处理 ROT 语句
			fmt.Printf("Rotation: %s\n", stmt.Angle)
		case *parser.AssignmentStatement:
			// 处理赋值语句
			value := i.evaluateExpression(stmt.Value)
			i.variables[stmt.Identifier] = value
			fmt.Printf("Assign: %s = %f\n", stmt.Identifier, value)
		case *parser.ForStatement:
			// 处理 FOR 循环
			i.executeForStatement(stmt)
		default:
			fmt.Println("Unknown statement")
		}
	}
}

// 执行 FOR 循环语句
func (i *Interpreter) executeForStatement(forStmt *parser.ForStatement) {
	start := i.evaluateExpression(forStmt.Start)
	end := i.evaluateExpression(forStmt.End)
	step := i.evaluateExpression(forStmt.Step)

	fmt.Printf("For %s = %f to %f step %f:\n", forStmt.LoopVar, start, end, step)

	for j := start; j <= end; j += step {
		i.variables[forStmt.LoopVar] = j
		i.evaluateStatement(forStmt.Body)
	}
}

// 计算表达式的值
func (i *Interpreter) evaluateExpression(expr parser.Expression) float64 {
	switch expr := expr.(type) {
	case *parser.ConstantExpression:
		var value float64
		fmt.Sscanf(expr.Value, "%f", &value)
		return value
	case *parser.BinaryExpression:
		left := i.evaluateExpression(expr.Left)
		right := i.evaluateExpression(expr.Right)
		switch expr.Operator {
		case token.PLUS:
			return left + right
		case token.MINUS:
			return left - right
		case token.MUL:
			return left * right
		case token.DIV:
			return left / right
		case token.COMMA:
			// 处理绘图中的 x, y 坐标
			return left // 可以根据需要扩展绘图计算
		default:
			fmt.Println("Unknown operator")
			return 0
		}
	default:
		fmt.Println("Unknown expression type")
		return 0
	}
}

// 执行单个语句
func (i *Interpreter) evaluateStatement(stmt parser.Statement) {
	switch stmt := stmt.(type) {
	case *parser.AssignmentStatement:
		// 根据是否为绘图语句来区分处理
		if stmt.Identifier == "DRAW" {
			i.handleDrawStatement(stmt) // 处理绘图
		} else {
			// 普通赋值语句
			value := i.evaluateExpression(stmt.Value)
			i.variables[stmt.Identifier] = value
			fmt.Printf("Assigned: %s = %.6f\n", stmt.Identifier, value)
		}
	case *parser.ForStatement:
		// FOR 循环语句
		i.executeForStatement(stmt)
	default:
		fmt.Println("Unknown statement")
	}
}

// 处理绘图语句
func (i *Interpreter) handleDrawStatement(stmt *parser.AssignmentStatement) {
	// 假设 DRAW 是一个二元表达式
	if binaryExpr, ok := stmt.Value.(*parser.BinaryExpression); ok {
		// 解析 x 和 y 坐标
		x := i.evaluateExpression(binaryExpr.Left)
		y := i.evaluateExpression(binaryExpr.Right)

		// 打印绘图坐标
		fmt.Printf("Drawing at (%.2f, %.2f)\n", x, y)

		// 使用 gonum/plot 绘制图形
		p, err := plot.New()
		if err != nil {
			fmt.Println("Error creating plot:", err)
			return
		}

		// 设置标题和标签
		p.Title.Text = "Plotting Coordinates"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		// 创建一个数据集，绘制坐标点
		points := make(plotter.XYs, 1)
		points[0].X = x
		points[0].Y = y

		// 创建一个散点图，并将数据集添加到图表中
		scatter, err := plotter.NewScatter(points)
		if err != nil {
			fmt.Println("Error creating scatter plot:", err)
			return
		}

		p.Add(scatter)

		// 保存图形为 PNG 文件
		if err := p.Save(4*vg.Inch, 4*vg.Inch, "output.png"); err != nil {
			fmt.Println("Error saving plot:", err)
		} else {
			fmt.Println("Plot saved as 'output.png'")
		}
	}
}
