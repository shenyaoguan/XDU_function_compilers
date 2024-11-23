package main

import (
	"compilers/interpreter"
	"compilers/lexer"
	"compilers/parser"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// 从命令行参数获取文件路径
	if len(os.Args) < 2 {
		log.Fatal("Usage: mygo <filename.mygo>")
	}
	filePath := os.Args[1]

	// 读取文件内容
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// 词法分析：将文件内容转换为 token 列表
	lex := lexer.New(string(fileContent))

	// 语法分析：将 token 列表解析为语法树
	parser := parser.New(lex)
	statements := parser.ParseProgram()

	// 创建并运行解释器
	interpreter := interpreter.NewInterpreter()
	interpreter.RunProgram(statements)

	// 输出绘图结果
	fmt.Println("Program executed. Check the generated plot in 'output.png'.")
}
