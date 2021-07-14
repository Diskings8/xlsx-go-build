package script

import (
	"fmt"
	"os"
	"strings"
)

func (gen *Generate) SpFROHeadData() error {
	structData := ""
	lDESC := gen.headRowMap[DESC]
	lKEY := gen.headRowMap[KEY]
	lEXPORT := gen.headRowMap[EXPORT]
	for _, value := range gen.headSheetData {
		switch value[lEXPORT] {
		case OPALL, OPF:
			structData += fmt.Sprintf(frodesc, value[lKEY])
			if value[lDESC] != "" {
				structData += fmt.Sprintf(frodesc, value[lDESC])
			}
			structData += "\n"
		default:
			return fmt.Errorf("SplicingData|value[EXPORT]:\"%v\" is not in s,c,all", value[lEXPORT])
		}
		}
	gen.writeHeadData = structData
	return nil
}

// 拼装struct-body内容
func (gen *Generate) SpFROBodyData() error {
	STRUCTNAME := gen.structName
	IdCol := gen.keyIdCol - 1
	lKEY := gen.headRowMap[KEY]
	bodyData := fmt.Sprintf(frohead, ToLower(STRUCTNAME))
	for yk := 0; yk <len(gen.bodySheetData); yk++{
		onedata := gen.bodySheetData[yk]
			bodyData += fmt.Sprintf(froid, onedata[IdCol])
		for k, j := range onedata {
			bodyData += fmt.Sprintf(" %s:%s,", firstRuneToUpper(gen.headSheetData[k][lKEY]), j)
		}
		bodyData = strings.TrimRight(bodyData, "")
		if yk == len(gen.bodySheetData) -1 {
			bodyData += "}"
		}else{
			bodyData += "},"
		}

	}
	bodyData += fmt.Sprintf(froend,ToLower(STRUCTNAME))
	gen.writeBodyData = bodyData

	return nil
}

func (gen *Generate) froDataWrite() error {

	datahead := gen.writeHeadData
	databody := gen.writeBodyData
	str := strings.Split(gen.savePath, "\\")
	if len(str) == 0 {
		return fmt.Errorf("WriteNewFile|len(str) is 0")
	}
	data := datahead + databody
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


