package main

import (
	"fmt"
	"giuhub.com/foxsugar-sanse/GO语言小Demo/控制台计算机/src/model/container/stack"
	"unicode"
)

func main() {
	centerExpression := "11+2+(3*4)"
	fmt.Print(infix2ToPostfix(centerExpression))

}

// 中缀表达式转后缀表达式
func infix2ToPostfix(exp string) string {
	var s stack.Stack
	s.Init()
	postfix := ""
	expLen := len(exp)

	// 遍历整个表达式
	for i := 0; i < expLen; i++ {

		char := string(exp[i])

		switch char {
		case " ":
			continue
		case "(":
			// 左括号直接入栈
			s.Push("(")
		case ")":
			// 右括号则弹出元素直到遇到左括号
			for !s.IsEmpty() {
				preChar := s.Peek().(string)
				if preChar == "(" {
					s.Pop() // 弹出 "("
					break
				}
				postfix += preChar
				//postfix = fmt.Sprintf("%s%s",postfix,preChar)
				s.Pop()
			}

			// 数字则直接输出
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			j := i
			digit := ""
			/* 大于或等于10以上要做的，因为单次只能遍历一个字符，为了确保能够正确读取10及以上的数
				需要做一个判断：之后的字符是否为10进制的数，不是则和遍历的字符拼接，为运算符则退出
			 */
			for ; j < expLen && unicode.IsDigit(rune(exp[j])); j++ {
				digit += string(exp[j])
			}
			postfix += digit
			i = j - 1 // i 向前跨越一个整数，由于执行了一步多余的 j++，需要减 1

		default:
			// 操作符：遇到高优先级的运算符，不断弹出，直到遇见更低优先级运算符
			for !s.IsEmpty() {
				top := s.Peek().(string)
				if top == "(" || isLower(top, char) {
					break
				}
				postfix += top
				s.Pop()
			}
			// 低优先级的运算符入栈
			s.Push(char)
		}
	}

	// 栈不空则全部输出
	for !s.IsEmpty() {
		postfix += s.Pop().(string)
		//postfix = fmt.Sprintf("%s%s",postfix, s.Pop())
	}

	return postfix
}

// 比较运算符栈栈顶 top 和新运算符 newTop 的优先级高低
func isLower(top string, newTop string) bool {
	// 注意 a + b + c 的后缀表达式是 ab + c +，不是 abc + +
	switch top {
	case "+", "-":
		if newTop == "*" || newTop == "/" {
			return true
		}
	case "(":
		return true
	}
	return false
}