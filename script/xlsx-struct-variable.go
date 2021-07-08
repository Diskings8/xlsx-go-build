package script

import (
	"strings"
)

// struct-head-data
var (
	structBegin          = "type %s struct {\n"                             // 结构体开始
	structValue          = "    %-15s %-8s	`col:\"%-15s\" client:\"%-8s\"`" // 结构体的内容
	structValueForServer = "    %-15s %-8s	`col:\"%-15s\"`"                 // 服务端使用的结构体内容
	structRemarks        = "	 // %10s"                                      // 结构体备注
	structValueEnd       = "\n"                                             // 结构体内容结束
	structEnd            = "}\n"                                            // 结构体结束
	header               = "package %s\n\r"                                 // 文件头
)

// struct-body-data
var (
	structFuncHead1   = "import \"errors\" \n"
	structFuncHead2   = "func (%s *%s)Get(id int) (*%s, error){\n"
	structBody1       = "    switch id {\n"
	structSwitchCase1 = "        case %s : \n"
	structSwitchCase2 = "            return &%s{"
	structSwitchCase3 = " }, nil\n"
	structSwitchCase4 = "        default: return nil, errors.New(\"no data\")\n"
	structSwitchCase5 = "    }\n"
	structFuncEnd     = "}"
	structheader      = "package %s\n\n"
)

// 字符串首字母转换成大写
func firstRuneToUpper(str string) string {
	data := []byte(str)
	for k, v := range data {
		if k == 0 {
			first := []byte(strings.ToUpper(string(v)))
			newData := data[1:]
			data = append(first, newData...)
			break
		}
	}
	return string(data[:])
}

// 字符串全部小写
func ToLower(str string) string {
	return strings.ToLower(str)
}
