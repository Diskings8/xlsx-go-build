package script

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

type Generate struct {
	savePath string // 生成文件的保存路径
	allType  string // 文件当中的数据类型
	xlsxFileName string //

	writeHeadData string // 生成文件的内容
	writeBodyData string // 生成文件的内容

	fileName   string // 导出的文件名
	structName string // 导出结构体

	headRowMap map[string]int // title-tag的映射信息

	headRow       []int            //title-tag的所在行信息
	headSheetData []map[int]string //组装struct-head数据前的信息

	bodyRow       []int      //value-data的所在行信息
	bodySheetData [][]string //组装struct-data数据前的信息

	keyIdCol	int // key-Id所在列
}

// 不生成共用的文件，采用goroute
func OneProOpenXlsx(savePath, allType, fileName string) {
	gt := Generate{savePath: savePath, allType: allType}
	if err := gt.doInitAndOpen(fileName); err != nil {
		fmt.Println(err)
	}
}

// 测试输出调用的函数
func (gen *Generate) print() {
	fmt.Println(gen.headRowMap)
	fmt.Println(gen.headRow)
	fmt.Println(gen.headSheetData)
	fmt.Println(gen.bodyRow)
	fmt.Println(gen.bodySheetData)
	fmt.Println(gen.writeBodyData)
}

// 初始化数据
func (gen *Generate) doInitAndOpen(fileName string) error {
	gen.xlsxFileName = fileName
	gen.headSheetData = make([]map[int]string, 0)
	gen.bodySheetData = make([][]string, 0)
	gen.headRow = make([]int, 0)
	gen.bodyRow = make([]int, 0)
	gen.headRowMap = make(map[string]int)
	gen.keyIdCol = 1
	return gen.openXlsxFile(fileName)
}

// 初始化单元，并遍历该表中的各个子表
func (gen *Generate) openXlsxFile(fileName string) error {
	wb, err := xlsx.OpenFile(fileName)
	if err != nil {
		return fmt.Errorf("ReadExcel|xlsx.OpenFile is err :%v", err)
	}
	// 遍历工作表
	for _, sheet := range wb.Sheets {
		// 执行自定义导出规则
		if err := gen.customRules(sheet); err != nil {
			return err
		}
		// 执行head数据组装
		if err := gen.SplicingHeadData(); err != nil {
			return err
		}
		//执行body数据组装
		if err := gen.SplicingBodyData(); err != nil {
			return err
		}

		//判断是否有head-cache 和 body-cache(可以为空) 数据读入成功
		if gen.writeHeadData == "" {
			return fmt.Errorf("ReadExcel|gen.head-data is nil")
		}
		//执行head数据文件落地
		if err := gen.mergeDataWrite(); err != nil {
			return err
		}

		// gen.print()
	}
	return nil
}

// 个人自定义配置表头
func (gen *Generate) customRules(sheet *xlsx.Sheet) error {
	gen.findTitle(sheet)

	if err := gen.setBaseInfo(sheet); err != nil {
		return err
	}
	// 遍历列
	if gen.checkIdUnique(sheet) {
		text := fmt.Sprintf("%s sheet has Duplicate id in %s file", sheet.Name, gen.xlsxFileName)
		return errors.New(text)
	}
	// 排除第一列
	for i := 1; i < sheet.MaxCol; i++ {
		//判断某一列的数据类型是否为空或者是否没有相关Type配置
		if !gen.filterNoExport(sheet, i) {
			continue
		}
		cellHeadData := make(map[int]string)
		// 遍历行-头信息部分
		for _, j := range gen.headRow {
			cellHeadData[j] = sheet.Cell(j, i).Value
		}
		gen.headSheetData = append(gen.headSheetData, cellHeadData)
	}
	for _, i := range gen.bodyRow {
		cellBodyData := make([]string, 0)
		for j := 0; j < sheet.MaxCol; j++ {
			if !gen.filterNoExport(sheet, j) {
				continue
			}
			cellBodyData = append(cellBodyData, sheet.Cell(i, j).Value)
		}
		gen.bodySheetData = append(gen.bodySheetData, cellBodyData)
	}
	return nil
}

// 不导出的不处理
// 如果没有准确的数据类型TYPE和需要导出的声明EXPORT，则默认为不导出
// 兼容默认TYPE=int
func (gen *Generate) filterNoExport(sheet *xlsx.Sheet, j int) bool {
	if _, ok := gen.headRowMap["TYPE"]; !ok {
		return false
	}
	if exportRow, ok := gen.headRowMap["EXPORT"]; !ok {
		return false
	} else if sheet.Cell(exportRow, j).Value != "all" && sheet.Cell(exportRow, j).Value != "s" {
		return false
	}
	return true
}
//
func (gen *Generate)findKeyCol(sheet *xlsx.Sheet){
	if keyRow, ok := gen.headRowMap["KEY"]; !ok {
		for i := 0; i<  sheet.MaxCol; i++{
			if sheet.Cell(keyRow, i).Value == "id" || sheet.Cell(keyRow, i).Value == "Id" || sheet.Cell(keyRow, i).Value == "ID" {
				gen.keyIdCol = i
				break
			}
		}
	}
}

func (gen *Generate)checkIdUnique(sheet *xlsx.Sheet) bool{
	colIds := make([]string,0)
	for i := 0; i<  sheet.MaxRow; i++{
		val := sheet.Cell(i, gen.keyIdCol).Value
		if  val!= "" && sheet.Cell(i,0).Value == "VALUE"{
			colIds = append(colIds, val)
		}
	}
	return ContainsDuplicate(colIds)
}

// 头信息映射建立
func (gen *Generate) findTitle(sheet *xlsx.Sheet) {
	j := 0
	for i := 0; i < sheet.MaxRow; i++ {
		celValue := sheet.Cell(i, 0).Value
		switch celValue {
		case "STRUCTNAME", "TYPE", "DESC", "KEY", "EXPORT", "FILENAME":
			// 对上述类型记录字段以便生成struct的声明，同时记录所在行 进行数据读取
			gen.headRow = append(gen.headRow, i)
			gen.headRowMap[celValue] = i
			j++
		case "VALUE":
			// 对上述类型记录所在行 进行数据读取
			gen.bodyRow = append(gen.bodyRow, i)
		}
	}
}

// 设置基础导出文件所需信息
// 基本的导出结构体信息名和导出的data数据文件名必须要有
func (gen *Generate) setBaseInfo(sheet *xlsx.Sheet) error {

	structNameRow, ok := gen.headRowMap["STRUCTNAME"]
	if !ok {
		return errors.New("not structName")
	}
	gen.structName = sheet.Cell(structNameRow, 1).Value

	fileNameRow, ok := gen.headRowMap["FILENAME"]
	if !ok {
		return errors.New("not fileName")
	}
	gen.fileName = sheet.Cell(fileNameRow, 1).Value

	return nil
}

// 拼装struct-head内容
func (gen *Generate) SplicingHeadData() error {
	structData := fmt.Sprintf(structBegin, firstRuneToUpper(gen.structName))
	TYPE := gen.headRowMap["TYPE"]
	DESC := gen.headRowMap["DESC"]
	KEY := gen.headRowMap["KEY"]
	EXPORT := gen.headRowMap["EXPORT"]
	for _, value := range gen.headSheetData {
		err := gen.CheckType(value[TYPE])
		if err != nil {
			return err
		}
		extChangeStructValue := extTypeChange(value[TYPE])
		switch value[EXPORT] {
		case "all":
			structData += fmt.Sprintf(structValue, firstRuneToUpper(value[KEY]), extChangeStructValue, value[KEY], value[KEY])
			if value[DESC] != "" {
				structData += fmt.Sprintf(structRemarks, value[DESC])
			}
			structData += fmt.Sprintf(structValueEnd)
		case "s":
			structData += fmt.Sprintf(structValueForServer, firstRuneToUpper(value[KEY]), extChangeStructValue, value[KEY])
			if value[DESC] != "" {
				structData += fmt.Sprintf(structRemarks, value[DESC])
			}
			structData += fmt.Sprintf(structValueEnd)
		case "c":
			continue
		case "":
			continue
		default:
			return fmt.Errorf("SplicingData|value[EXPORT]:\"%v\" is not in s,c,all", value[EXPORT])
		}
	}
	structData += structEnd
	gen.writeHeadData += structData
	return nil
}

// 拼装struct-body内容
func (gen *Generate) SplicingBodyData() error {
	STRUCTNAME := gen.structName
	IdCol := gen.keyIdCol - 1
	KEY := gen.headRowMap["KEY"]
	TYPE := gen.headRowMap["TYPE"]
	bodyData := fmt.Sprintf(structFuncHead2, ToLower(STRUCTNAME), firstRuneToUpper(STRUCTNAME), firstRuneToUpper(STRUCTNAME))
	bodyData += structBody1
	for _, onedata := range gen.bodySheetData {
		bodyData += fmt.Sprintf(structSwitchCase1, onedata[IdCol])
		bodyData += fmt.Sprintf(structSwitchCase2, firstRuneToUpper(STRUCTNAME))
		for k, j := range onedata {
			changedata := extTypeChangeWithValue(gen.headSheetData[k][TYPE], j)
			bodyData += fmt.Sprintf(" %s:%s,", firstRuneToUpper(gen.headSheetData[k][KEY]), changedata)
		}
		bodyData = strings.TrimRight(bodyData, ",")
		bodyData += structSwitchCase3
	}
	bodyData += structSwitchCase4
	bodyData += structSwitchCase5
	bodyData += structFuncEnd
	gen.writeBodyData += bodyData
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
