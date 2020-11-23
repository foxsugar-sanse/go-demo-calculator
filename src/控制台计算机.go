package main

// Python demo重写之控制台计算机

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"giuhub.com/foxsugar-sanse/GO语言小Demo/控制台计算机/src/model/cmdmodel"
	"giuhub.com/foxsugar-sanse/GO语言小Demo/控制台计算机/src/model/container/stack"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Uivm struct {
	// 管理和生成界面
	uiView 				string
	uiViewPlus 			string
	uiViewTerm 			string
	uiManage 			string
	comParameter 		string
	// 创建bytes读写缓冲区
	comParameterDeRune 	bytes.Buffer
	// 计算器用于动态显示数据
	showStr 			bytes.Buffer
	// 未处理前是一个数组
	noneshowStr			[]string
	// 长度限制，默认23个字符
	lengthLimit 		int
	// 传入的数据类型
	types 				interface{}
	// 控制台用于显示和计算数据的列表
	contsList 			*list.List
	// 用于绑定视图的map
	viewMap 			map[int]string
	// 控制展示什么样的界面
	showviews 			int

	MathcFuctionObject 	*MathConts
	object 				*UivmOperator
	TermDataCreate

	// 提供一些帮助信息
	__help__ 			string

	// 多任务通信管道
	multitasking 		chan int
}

type UivmOperator struct {
	operatorRemain 		string
	operator02	   		string
}

type TermDataCreate struct {
	// term模式
	termOf 				bool
	termInfomation 		string
	// 简写-全名-函数/方法名
	cmdStringDe 		map[string]string
}

type MathConts struct {
	// 管理计算器运算
	oldNum 				int64
	oldFloat 			float64
	oldUNum 			uint64
	operatorTop 		map[int]string
	//用户输入的公式
	expressionString	bytes.Buffer
}

func (M *MathConts) MCinit() {
	// 运算符优先级
	oldOperatorTop := map[int]string{0:"()",1:"*^",2:"/",3:"+-"}
	M.operatorTop = oldOperatorTop
}

func (U *Uivm) Init() {
	// init args
	U.contsList = list.New()
	// 运算符格式化
	U.object = &UivmOperator{"%s", "½"}

	U.termOf = false
	U.lengthLimit = 23
	U.noneshowStr = make([]string, U.lengthLimit)
	U.comParameterDeRune = bytes.Buffer{}
	U.showStr = bytes.Buffer{}
	U.multitasking = make(chan int)

	U.uiView = ("|-----Compute V1.0------|\n"+
		"|%s|\n"+ //23长度
		"|[ C ] [ X ] [ "+U.object.operatorRemain+" ] [ * ]|\n"+
		"|[ 1 ] [ 2 ] [ 3 ] [ + ]|\n"+
		"|[ 4 ] [ 5 ] [ 6 ] [ - ]|\n"+
		"|[ 7 ] [ 8 ] [ 9 ] [ / ]|\n"+
		"|[ 0 ] [ . ] [    =    ]|\n"+
		"|-----------------------|\n")
	// Plus计算机视图
	U.uiViewPlus = ("|--------Compute V2.0---------|\n"+
		"|%s|\n"+ //30长度
		"|[ C ] [ X ] [ "+U.object.operatorRemain+" ] [ * ] [ √ ]|\n"+
		"|[ 1 ] [ 2 ] [ 3 ] [ + ] [ ^ ]|\n"+
		"|[ 4 ] [ 5 ] [ 6 ] [ - ] [ "+U.object.operator02+" ]|\n"+
		"|[ 7 ] [ 8 ] [ 9 ] [ / ] [ π ]|\n"+
		"|[ 0 ] [ . ] [    =    ] [ E ]|\n"+
		"|[ ( ] [ ) ]------------------|\n")
	U.uiViewTerm = "Powershell@root❱❱❱%s "
	NomalViewMap := map[int]string{0:U.uiView, 1:U.uiViewPlus, 2:U.uiViewTerm}
	U.viewMap = NomalViewMap
	U.__help__ += fmt.Sprintf("    %s | %s\n","h 打印帮助信息", "C 清空显示数据")
	U.__help__ += fmt.Sprintf("    %s | %s\n","q 退出当前程序", "X 删除显示数据的后一位数")
	U.__help__ += fmt.Sprintf("    %s\n","p 开启全面功能")
	U.__help__ += fmt.Sprintf("    %s\n","r 回到简约模式")
	U.__help__ += fmt.Sprintf("    %s\n","s 打开命令行模式")
	U.__help__ += fmt.Sprintf("\n按回车或者等待100秒退出🔙")
	//h 打印帮助信息 | C 清空显示数据
	//q 退出当前程序 | X 删除显示数据的后一位数
	//p 开启全面功能
	//r 回到简约模式
	//s 打开命令行模式

}

func (U *Uivm) termOptions() {
	// 控制台模式
	// 排序命令
	//U.termOf = false
	//U.showviews = 0
	// 分隔字符串（cmd args data）
	var recvReturnString string
	stringsSplit := strings.Split(U.termInfomation, " ")
	switch len(stringsSplit) {
	case 3:
		recvReturnString = cmdmodel.PushCmd(stringsSplit[0], stringsSplit[1], stringsSplit[2])
	case 2:
		recvReturnString = cmdmodel.PushCmd(stringsSplit[0], stringsSplit[1], "")
	case 1:
		recvReturnString = cmdmodel.PushCmd(stringsSplit[0], "", "")
	case 0:
		// 无输入，不做反应
	}
	//if recvReturnString != ""{
	//	fmt.Printf("%s\n", recvReturnString)
	//} else if recvReturnString == "exit" {
	//	U.termOf = false
	//	U.showviews = 0
	//}
	switch recvReturnString {
	case "exit":
		U.termOf = false
		U.showviews = 0
	default:
		fmt.Printf("%s\n", recvReturnString)
	}
}

func (U *Uivm) uiViews(showviews int) {

	// 规范长度
	if len(U.comParameterDeRune.String()) > U.lengthLimit {
		// 使用科学计数法
		// 超过只打印从后面开始的几个数
		// 满足没有运算符，才会转为科学计数法显示
		okNum := 0
		for i := 0; i < len(U.comParameterDeRune.String()); i++ {
			if unicode.IsDigit(rune(U.comParameterDeRune.String()[i])) {okNum++;}
		}
		if okNum == len(U.comParameterDeRune.String()) {
			// 使用科学计数法
			exportNum, _ := strconv.ParseFloat(U.comParameterDeRune.String(),64)
			exportNumTwo := fmt.Sprintf("%E",exportNum)
			// 判断科学计数法之后是否还是会过长
			if len(exportNumTwo) > U.lengthLimit {
				// 省略后方输入的数
				var newBytesBuffer bytes.Buffer
				for i := 0; i < U.lengthLimit; i++ {
					newBytesBuffer.WriteByte(U.comParameterDeRune.String()[i])
				}
				U.showStr = newBytesBuffer
			} else {
				U.showStr.Reset()
				U.showStr.WriteString(exportNumTwo)
			}
		} else {
			// 省略后方输入的数
			var newBytesBuffer bytes.Buffer
			for i := 0; i < U.lengthLimit; i++ {
				newBytesBuffer.WriteByte(U.comParameterDeRune.String()[i])
			}
			U.showStr = newBytesBuffer
		}
	}
	// 组织要显示的数据
	forNum := U.lengthLimit - len(U.comParameterDeRune.String())
	var beSplitString string
	if len(U.showStr.String()) != 0 {
		for i := 0; i < U.lengthLimit - len(U.showStr.String()); i++ {
			beSplitString += fmt.Sprintf("%s"," ")
		}
		beStr := U.showStr.String()
		U.showStr.Reset()
		U.showStr.WriteString(beSplitString);U.showStr.WriteString(beStr)
	}else {
		for i := 0; i < forNum; i++ {
			U.showStr.WriteString(" ")
		}
		if len(U.comParameterDeRune.String()) != 0 {
			U.showStr.WriteString(U.comParameterDeRune.String())
		}
	}
	if showviews != 0 {
		for key, value := range U.viewMap {
			if key == showviews && showviews !=2 {
				fmt.Printf(value,U.showStr.String(),"%")
			} else if key == showviews && showviews == 2 {
				fmt.Printf(value,"/")
			}
		}
	} else {
		fmt.Printf(U.uiView,U.showStr.String(),"%")
	}
	defer func() {U.showStr.Reset()}()
}

func (U *Uivm) uiManages(MathcFuctionObject *MathConts) {
	for true {
		U.MathcFuctionObject = MathcFuctionObject
		// 调用视图
		U.uiViews(U.showviews)
		// 从控制台接收要运算的数
		// 控制台模式
		if U.termOf == false {
			fmt.Print("请输入，输入h打印帮助信息：")
			fmt.Scan(&U.comParameter)
			U.processInput()

			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

		} else if U.termOf == true {
			//fmt.Scan(&U.TermDataCreate.termInfomation)
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			U.TermDataCreate.termInfomation = input.Text()
			U.termOptions()

		}
		//input := bufio.NewScanner(os.Stdin)
		//input.Scan()

	}
}

func (U *Uivm) processInput() {
	// 根据运算符的权重来筛选
		for _, i := range U.comParameter {
			if i == 'q' && len(U.comParameter) == 1  {
				os.Exit(1)
			} else if i == 's' && len(U.comParameter) == 1 {
				// 控制台
				U.showviews = 2
				U.termOf = true
			} else if i == 'p' && len(U.comParameter) == 1 {
				// Plus
				U.showviews = 1
				U.lengthLimit = 29
				oldnoneshowStr := make([]string, U.lengthLimit)
				// 更换模式造成的显示改变
				for next,result := range U.noneshowStr {
					oldnoneshowStr[next] = result
				}
				U.noneshowStr = oldnoneshowStr
			} else if i == 'r' && len(U.comParameter) == 1 {
				// 恢复原来的显示模式
				U.showviews = 0
				U.lengthLimit = 23
			} else if i == '=' {
				// 运算
				// 替换运算成功后的值
				U.comParameterDeRune.Reset()
				U.comParameterDeRune.WriteString(U.MathcFuctionObject.Conts())
			} else if i == 'C' && len(U.comParameter) == 1 {
				// 清除
				//var newBytesBuffer bytes.Buffer
				//U.comParameterDeRune = newBytesBuffer
				U.comParameterDeRune.Reset()
			} else if i == 'X' && len(U.comParameter) == 1 {
				// 删除一位数
				var newBytesBuffer bytes.Buffer
				//U.comParameterDeRune = newBytesBuffer
				for i := 0; i < len(U.comParameterDeRune.String())-1; i++ {
					newBytesBuffer.WriteByte(U.comParameterDeRune.String()[i])
				}
				// 显示和计算用的缓冲区进行统一
				U.comParameterDeRune = newBytesBuffer
				U.MathcFuctionObject.expressionString = newBytesBuffer
			} else if i == 'h' && len(U.comParameter) == 1{
				chanControl1 := make(chan int)
				runOK		 := 0
				go U.readKeyBoardEntry(chanControl1, &runOK)
				for i := 0; i < 101; i++ {
					if runOK == 1 {
						break
					}
					time.Sleep(time.Second * 1)
				}
				if runOK != 1 {
					chanControl1 <- 1
				}
				defer func() {close(chanControl1)}()
			} else {
				U.MathcFuctionObject.expressionString.WriteRune(i)
				U.contsList.PushBack(i)
				U.comParameterDeRune.WriteRune(i)
			}
		}


}

func (U *Uivm) readKeyBoardEntry(chanControl2 chan int, runok *int) bool {
	for 1==1 {
		fmt.Printf(U.__help__)
		exitstr := true
		go func() {
			j,_ := <-chanControl2
			if j == 1 {exitstr = false}
		}()
		if exitstr == false {
			return exitstr
		}
		robotgo.EventHook(hook.KeyDown, []string{"enter"}, func(e hook.Event) {
			robotgo.EventEnd()
		})

		s := robotgo.EventStart()
		<-robotgo.EventProcess(s)

		ok := robotgo.AddEvents("enter")
		if ok {
			*runok = 1
			return false
		}
	}
	return false
}

func (M *MathConts) Conts() string{
	// 负责管理计算机运算的一个方法
	//var addStr bytes.Buffer
	//for i := M.oldListNum.Front(); i != nil; i = i.Next() {
	//	addStr.WriteString(i.Value.(string))
	//}
	// 采用bytes转换字符为字符串
	postfix := M.infix2ToPostfix(M.expressionString.String())
	//清空数据
	defer func() {
		M.expressionString.Reset()
	}()
	return strconv.Itoa(M.calculate(postfix))
}

// 后缀表达式的运算
func (M *MathConts) calculate(postfix [100]string) int {
	var stack stack.Stack
	stack.Init()
	fixLen := len(postfix)
	for i := 0; i < fixLen && postfix[i] != ""; i++ {
		nextChar := postfix[i]
		// 数字：直接压栈
		if unicode.IsDigit(rune(nextChar[0])) {
			stack.Push(nextChar)
		} else {
			// 操作符：取出两个数字计算值，再将结果压栈
			num1, _ := strconv.Atoi(stack.Pop().(string))
			num2, _ := strconv.Atoi(stack.Pop().(string))
			switch nextChar {
			case "+":
				stack.Push(strconv.Itoa(num2 + num1))
			case "-":
				stack.Push(strconv.Itoa(num2 - num1))
			case "*":
				stack.Push(strconv.Itoa(num2 * num1))
			case "/":
				stack.Push(strconv.Itoa(num2 / num1))
			case "%":
				stack.Push(strconv.Itoa(num2 % num1))
			case "^":
				stack.Push(strconv.Itoa(num2 ^ num1))
			}
		}
	}
	result, _ := strconv.Atoi(stack.Peek().(string))
	return result
}

// 中缀表达式转后缀表达式
func (M *MathConts) infix2ToPostfix(exp string) [100]string {
	//stack := stack.ItemStack{}
	var stack stack.Stack
	stack.Init()
	//postfix := ""
	var postfix [100]string
	expLen := len(exp)
	indexArrayStrings := 0
	// 遍历整个表达式
	for i := 0; i < expLen; i++ {

		char := string(exp[i])
		switch char {
		case " ":
			continue
		case "(":
			// 左括号直接入栈
			stack.Push("(")
		case ")":
			// 右括号则弹出元素直到遇到左括号
			for !stack.IsEmpty() {
				preChar := stack.Peek().(string)
				if preChar == "(" {
					stack.Pop() // 弹出 "("
					break
				}
				postfix[indexArrayStrings] = preChar
				indexArrayStrings++
				stack.Pop()
			}

			// 数字则直接输出
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			j := i
			digit := ""
			for ; j < expLen && unicode.IsDigit(rune(exp[j])); j++ {
				digit += string(exp[j])
			}
			//postfix += digit
			postfix[indexArrayStrings] = digit
			indexArrayStrings++
			i = j - 1 // i 向前跨越一个整数，由于执行了一步多余的 j++，需要减 1

		default:
			// 操作符：遇到高优先级的运算符，不断弹出，直到遇见更低优先级运算符
			for !stack.IsEmpty() {
				top := stack.Peek().(string)
				if top == "(" || M.isLower(top, char) {
					break
				}
				//postfix += top
				postfix[indexArrayStrings] = top
				indexArrayStrings++
				stack.Pop()
			}
			// 低优先级的运算符入栈
			stack.Push(char)
		}
	}

	// 栈不空则全部输出
	for !stack.IsEmpty() {
		//postfix += stack.Pop()
		postfix[indexArrayStrings] = stack.Pop().(string)
		indexArrayStrings++
	}

	return postfix
}

//运算符优先级的比较
func (M *MathConts) isLower(top string, newTop string) bool {
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

func main() {
	test := Uivm{}
	mathc := MathConts{}
	mathc.MCinit()
	test.Init()
	test.uiManages(&mathc)
}