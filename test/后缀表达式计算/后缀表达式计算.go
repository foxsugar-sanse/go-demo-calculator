package main

import (
	"fmt"
	"giuhub.com/foxsugar-sanse/GO语言小Demo/控制台计算机/src/model/container/stack"
	"strconv"
	"unicode"
)

func main() {
	var args [100]string
	args[0] = "12"
	args[1] = "12"
	args[2] = "+"
	fmt.Print(calculate(args))
}

func calculate(postfix [100]string) int {
	var s stack.Stack
	s.Init()
	fixLen := len(postfix)
	for i := 0; i < fixLen && postfix[i] != ""; i++ {
		//nextChar := string(postfix[i])
		nextChar := postfix[i]
		// 数字：直接压栈
		if unicode.IsDigit(rune(nextChar[0])) {
			s.Push(nextChar)
		} else {
			// 操作符：取出两个数字计算值，再将结果压栈
			num1, _ := strconv.Atoi(s.Pop().(string))
			num2, _ := strconv.Atoi(s.Pop().(string))
			switch nextChar {
			case "+":
				s.Push(strconv.Itoa(num1 + num2))
			case "-":
				s.Push(strconv.Itoa(num1 - num2))
			case "*":
				s.Push(strconv.Itoa(num1 * num2))
			case "/":
				s.Push(strconv.Itoa(num1 / num2))
			}
		}
	}
	result, _ := strconv.Atoi(s.Peek().(string))
	return result
}
