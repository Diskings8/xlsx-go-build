package script

import (
	"strings"
)

// struct-head-data
var (
	header1 = "package %s\n\n" // 文件头
	header2 = "import \"errors\" \n"

	structBegin          = "type %s struct {\n"                             // 结构体开始
	structValue          = "    %-15s %-8s	`col:\"%-15s\" client:\"%-8s\"`" // 结构体的内容
	structValueForServer = "    %-15s %-8s	`col:\"%-15s\"`"                 // 服务端使用的结构体内容
	structRemarks        = "	 // %10s"                                      // 结构体备注
	structValueEnd       = "\n"                                             // 结构体内容结束
	structEnd            = "}\n"                                            // 结构体结束

	structFuncHead2   = "\nfunc (%s *%s)Get(id int) (*%s, error){\n"
	structBody1       = "    switch id {\n"
	structSwitchCase1 = "        case %s : \n"
	structSwitchCase2 = "            return &%s{"
	structSwitchCase3 = " }, nil\n"
	structSwitchCase4 = "        default: return nil, errors.New(\"no data\")\n"
	structSwitchCase5 = "    }\n"
	structFuncEnd     = "}"
)
var (
	ENDFILENAME = "END_FILE"
	ENDSTRUCT = "END_NAME"

	FROFILENAME = "FRO_FILE"
	FROSRTUCT = "FRO_NAME"

	TYPE = "TYPE"
	DESC = "DESC"
	KEY = "KEY"
	EXPORT = "EXPORT"

	VALUE = "VALUE"

	OPS = "s"
	OPF = "c"
	OPALL = "all"
)

var(
	frohead = "// \n let %s: object = {\n"
	frodesc = "// \"%s\""
	froid = "\n        %s:{"
	froend = "\n};\n export default %s;\n"

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
func ContainsDuplicate(nums []string) bool {
	m := make(map[string]interface{}, len(nums))
	for _, v := range(nums){
		if _, ok := m[v];ok{
			return true
		}
		m[v] = true
	}
	return false
}


func valueWithOutBracket(value string) string {
	// value = strings.Trim(value, "[")
	// value = strings.Trim(value, "]")
	value = strings.Replace(value, "[", "{", -1)
	value = strings.Replace(value, "]", "}", -1)
	return strings.TrimSpace(value)
}

// 对于自定义类型的数据转换为go内置类型
func extTypeChangeWithValue(dataType, value string) string {
	value = valueWithOutBracket(value)
	switch dataType {
	case "":
		return value
	case "IntSlice":
		return "[]int" + value + ""
	case "IntSlice2":
		return "[][]int" + value + ""
	case "StringSlice":
		return "[]string" + value + ""
	default:
		return value
	}
}

// 对于自定义的类型转换为go内置类型
func extTypeChange(dataType string) string {
	switch dataType {
	case "":
		return "int"
	case "IntSlice":
		return "[]int"
	case "IntSlice2":
		return "[][]int{}"
	case "StringSlice":
		return "[]string"
	case "float":
		return "float64"
	default:
		return dataType
	}
}
