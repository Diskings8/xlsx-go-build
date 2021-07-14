package script

import (
	"fmt"
	"os"
	"strings"
)
// 拼装struct-head内容
func (gen *Generate) SpENDHeadData() error {
		structData := fmt.Sprintf(structBegin, firstRuneToUpper(gen.structName))
		lTYPE := gen.headRowMap[TYPE]
		lDESC := gen.headRowMap[DESC]
		lKEY := gen.headRowMap[KEY]
		lEXPORT := gen.headRowMap[EXPORT]
		for _, value := range gen.headSheetData {
			err := gen.CheckType(value[lTYPE])
			if err != nil {
				return err
			}
			extChangeStructValue := extTypeChange(value[lTYPE])
			switch value[lEXPORT] {
			case OPALL:
				structData += fmt.Sprintf(structValue, firstRuneToUpper(value[lKEY]), extChangeStructValue, value[lKEY], value[lKEY])
				if value[lDESC] != "" {
					structData += fmt.Sprintf(structRemarks, value[lDESC])
				}
				structData += fmt.Sprintf(structValueEnd)
			case OPS:
				structData += fmt.Sprintf(structValueForServer, firstRuneToUpper(value[lKEY]), extChangeStructValue, value[lKEY])
				if value[lDESC] != "" {
					structData += fmt.Sprintf(structRemarks, value[lDESC])
				}
				structData += fmt.Sprintf(structValueEnd)
			case "c":
				continue
			case "":
				continue
			default:
				return fmt.Errorf("SplicingData|value[EXPORT]:\"%v\" is not in s,c,all", value[lEXPORT])
			}
		}
		structData += structEnd
		gen.writeHeadData = structData

	return nil
}

// 拼装struct-body内容
func (gen *Generate) SpENDBodyData() error {
	STRUCTNAME := gen.structName
	IdCol := gen.keyIdCol - 1
	lKEY := gen.headRowMap[KEY]
	lTYPE := gen.headRowMap[TYPE]
	bodyData := fmt.Sprintf(structFuncHead2, ToLower(STRUCTNAME), firstRuneToUpper(STRUCTNAME), firstRuneToUpper(STRUCTNAME))
	bodyData += structBody1
	for _, onedata := range gen.bodySheetData {
		bodyData += fmt.Sprintf(structSwitchCase1, onedata[IdCol])
		bodyData += fmt.Sprintf(structSwitchCase2, firstRuneToUpper(STRUCTNAME))
		for k, j := range onedata {
			changedata := extTypeChangeWithValue(gen.headSheetData[k][lTYPE], j)
			bodyData += fmt.Sprintf(" %s:%s,", firstRuneToUpper(gen.headSheetData[k][lKEY]), changedata)
		}
		bodyData = strings.TrimRight(bodyData, ",")
		bodyData += structSwitchCase3
	}
	bodyData += structSwitchCase4
	bodyData += structSwitchCase5
	bodyData += structFuncEnd
	gen.writeBodyData = bodyData

	return nil
}

// 检测解析出来的字段类型是否符合要求
func (gen *Generate) CheckType(dataType string) error {
	switch dataType {
	case "":
		return nil
	default:
		res := strings.Index(gen.allType, dataType)
		if res == -1 {
			return fmt.Errorf("CheckType|struct:\"%v\" dataType:\"%v\" is not in provide dataType", gen.structName, dataType)
		}
		return nil
	}
}

// 拼装好的struct-head写入新的文件
func (gen *Generate) mergeDataWrite() error {

		datahead := gen.writeHeadData
		databody := gen.writeBodyData
		str := strings.Split(gen.savePath, "\\")
		if len(str) == 0 {
			return fmt.Errorf("WriteNewFile|len(str) is 0")
		}
		header := fmt.Sprintf(header1, str[len(str)-1])
		header += header2
		data := header + datahead + databody
		fw, err := os.OpenFile(gen.savePath+"\\"+gen.fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("WriteNewFile|OpenFile is err:%v", err)
		}
		defer func() { _ = fw.Close() }()
		_, err = fw.Write([]byte(data))
		if err != nil {
			return fmt.Errorf("WriteNewFile|Write is err:%v", err)
		}

	return nil
}
