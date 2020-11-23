package stack


// 栈的实现，基于空接口切片
//
//type Item interface {
//}
//
//// ItemStack the stack of items
//type ItemStack struct {
//	items []Item
//}
//
//// New Create a new ItemStack
//func (s *ItemStack) New() *ItemStack {
//	s.items = []Item{}
//	return s
//}
//
//// Push adds an Item to the top of the stack
//func (s *ItemStack) Push(t Item) {
//	s.items = append(s.items, t)
//}
//
//// Pop removes an Item from the top of the stack
//func (s *ItemStack) Pop() *Item {
//	item := s.items[len(s.items)-1] // 后进先出
//	s.items = s.items[0:len(s.items)-1]
//	return &item
//
//}

import "sync"

type Item interface{}

type Stack struct {
	slices []Item
	lock sync.RWMutex
}

// 初始化栈并返回一个可以操作的对象
func (s *Stack) Init() *Stack {
	s.slices = []Item{}
	return s
}

// 往切片中添加数据
func (s *Stack) Push(data Item) {
	s.lock.Lock()
	s.slices = append(s.slices, data)
	s.lock.Unlock()
}

// 根据先进后出的原则弹出栈中的元素
func (s *Stack) Pop() Item {
	s.lock.Lock()
	slice := s.slices[len(s.slices)-1]
	s.slices = s.slices[0:len(s.slices)-1]
	s.lock.Unlock()
	return slice
}

// 判断栈是否为空，不为空返回true，否则返回false
func (s *Stack) IsEmpty() bool {
	return len(s.slices) == 0
}

// 如果栈不为空的话返回栈顶的第一个元素
func (s *Stack) Peek() Item {
	return s.slices[len(s.slices)-1]
}

// 返回元素栈的长度
func (s *Stack) Length() int {
	return len(s.slices)
}