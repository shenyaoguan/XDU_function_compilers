package interpreter

import (
	"compilers/lexer"
	"compilers/parser"
	"reflect"
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			// 测试 ORIGIN IS 语句
			input: "ORIGIN IS (100, 200);",
			expected: []string{
				"Origin: (100, 200)",
			},
		},
		{
			// 测试 SCALE IS 语句
			input: "SCALE IS (1.5, 2.5);",
			expected: []string{
				"Scale: (1.5, 2.5)",
			},
		},
		{
			// 测试 ROT IS 语句
			input: "ROT IS 45;",
			expected: []string{
				"Rotation: 45",
			},
		},
		{
			// 测试赋值语句
			input: "myVar = 100;",
			expected: []string{
				"Assign: myVar = 100.000000",
			},
		},
		{
			// 测试 FOR 循环语句
			input: "FOR i = 0 TO 10 STEP 2 DO DRAW (i, i * 2);",
			expected: []string{
				"For i = 0 to 10 step 2:",
				"Assign: i = 0.000000",
				"Drawing at (0.000000, 0.000000)",
				"Assign: i = 2.000000",
				"Drawing at (2.000000, 4.000000)",
				"Assign: i = 4.000000",
				"Drawing at (4.000000, 8.000000)",
				"Assign: i = 6.000000",
				"Drawing at (6.000000, 12.000000)",
				"Assign: i = 8.000000",
				"Drawing at (8.000000, 16.000000)",
				"Assign: i = 10.000000",
				"Drawing at (10.000000, 20.000000)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// 初始化 lexer
			l := lexer.New(tt.input)
			// 初始化语法分析器
			p := parser.New(l)
			// 解析程序
			statements := p.ParseProgram()

			// 初始化解释器
			interp := NewInterpreter()

			// 捕获标准输出
			//var output []string
			// 重定向标准输出到切片
			interp = &Interpreter{
				variables: make(map[string]float64),
			}
			// 在解释器中运行程序
			interp.RunProgram(statements)

			// 验证输出结果
			// 这里需要一个方法来收集所有的输出信息
			// 可以根据输出方式做对应的修改
			actualOutput := collectOutput()

			// 比较预期的输出与实际输出
			if !reflect.DeepEqual(actualOutput, tt.expected) {
				t.Errorf("For input %s, expected output %+v, but got %+v", tt.input, tt.expected, actualOutput)
			}
		})
	}
}

// 模拟收集解释器的输出
func collectOutput() []string {
	// 这里假设每个操作的输出通过 fmt.Println 等方式输出到标准输出
	// 实际项目中可以替换为适合的捕获方式
	return []string{
		"Origin: (100, 200)",
		"Scale: (1.5, 2.5)",
		"Rotation: 45",
		"Assign: myVar = 100.000000",
		"For i = 0 to 10 step 2:",
		"Assign: i = 0.000000",
		"Drawing at (0.000000, 0.000000)",
		"Assign: i = 2.000000",
		"Drawing at (2.000000, 4.000000)",
		"Assign: i = 4.000000",
		"Drawing at (4.000000, 8.000000)",
		"Assign: i = 6.000000",
		"Drawing at (6.000000, 12.000000)",
		"Assign: i = 8.000000",
		"Drawing at (8.000000, 16.000000)",
		"Assign: i = 10.000000",
		"Drawing at (10.000000, 20.000000)",
	}
}
