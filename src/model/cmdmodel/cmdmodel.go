package cmdmodel

import(
	"os/exec"
	"os"
)

type cmdmodelDtaObject struct {
	information 		string
	cmdConts			string
	userData 			string
	processKeyValueDict map[string]string
	processArrayArgs 	[1][10]string
	// 简写-全名-对应参数
	// {{ss,ll,dd},{select,ls,delete},{{-h,-d,-s},{-l},{-h}}}
	// 使用数组对应关系无法实现
	//processKeyValueArray   [3][3][3]string
	cmdIsHas 			bool
	// select命令相关
	__help__			string
	__version__ 		string
}

func (C *cmdmodelDtaObject)Init() {
	C.cmdIsHas 			  = false
	C.processKeyValueDict = map[string]string{"sl":"select","cls":"clear"} 		// 排序功能的简写与全名
	C.processArrayArgs	  = [1][10]string{{"select","-h","-mp","-ho","-sj","-v","-ex"}} //参数
	C.__help__			  = "排序和查询工具-Power By XiaoHui@github.com/foxsugar-sanse"
	C.__version__		  = "0.10"
}

func (C *cmdmodelDtaObject) control() string{
	for key, value := range C.processKeyValueDict {
		if C.cmdConts == key {
			// 为true则证明有此命令
			C.cmdIsHas = true
			switch value {
			case "select":
				if C.sElect() != "" {
					// 不为空则有要return的数据
					return C.sElect()
				}
			case "clear":
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
				// 该命令无参数
			}
		}
	}

	if C.cmdIsHas == false {
		return "无此命令\n"
	}
	// 返回命令执行时间与资源消耗
	return ""
}

func (C *cmdmodelDtaObject) sElect() string{
	switch C.information {
	case "-h":
		// 调用帮助选项
		return C.__help__
	case "-v":
		// 调用版本选项
		return C.__version__
	case "-mp":
		// 调用冒泡排序选项
	case "-sj":
		// 一会想想看
	case "-ho":
		// 调用猴子排序选项
	case "-ex":
		// 退出终端模式
		return "exit"
	default:
		return "参数错误\n"
	}

	return ""
}

// 暴露的公共函数
func PushCmd(cmdConts,cmdInfromation,userData string) string{
	c := cmdmodelDtaObject{cmdConts:cmdConts, information:cmdInfromation, userData: userData}
	c.Init()
	return c.control()
}
