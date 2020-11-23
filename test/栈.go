package main

import (
	"fmt"
	"giuhub.com/foxsugar-sanse/GO语言小Demo/控制台计算机/src/model/container/stack"
)

func main() {
	var s stack.Stack
	s.Init()
	s.Push(1)
	s.Push(2)
	s.Push("nmsl")
	if s.IsEmpty() == true {
		fmt.Println("栈不为空！")
	}
	fmt.Println("s = ",s.Peek())
	fmt.Println(s.Length())
	indexNums := s.Length()
	for i := 0; i < indexNums; i++ {
		//fmt.Println(reflect.TypeOf(*s.Pop()))
		fmt.Println("s = ",s.Pop())
	}
	if s.Peek() == nil {
		fmt.Println(s.Peek())
	}

}
