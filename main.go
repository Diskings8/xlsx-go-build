package main

import (
	script "Script/script"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"sync"
)

var (
	savePath = flag.String("savePath", "", "Path to save the makefile")
	readPath = flag.String("readPath", "", "The path of reading Excel")
	allType  = flag.String("allType", "", "Specified field type")
	wg       sync.WaitGroup
)

func HandleReadExcel(readPath, savePath, allType string) error {
	if savePath == "" || allType == "" {
		return fmt.Errorf("ReadExcel|savePath or allType is nil")
	}
	// 获取目录下所有文件
	files, err := ioutil.ReadDir(readPath)
	if err != nil {
		return fmt.Errorf("ReadExcel|ReadDir is err:%v", err)
	}
	var xlsxFileNum = 0
	// 获取目录下所有符合规则的文件
	for _, file := range files {
		if path.Ext(file.Name()) != ".xlsx" {
			continue
		}
		filename := readPath + "\\" + file.Name()
		wg.Add(1)
		go func() {
			defer wg.Done()
			script.OneProOpenXlsx(savePath, allType, filename)
		}()
		xlsxFileNum++
	}
	wg.Wait()
	fmt.Printf("had deal %d xslx file\n", xlsxFileNum)
	return nil
}

// func testData() {
// 	obj := &data.Rap{}
// 	v, _ := obj.Get(1)
// 	fmt.Println(v.RecycleReward)
// }
func runHandle() {
	flag.Parse()
	if *savePath == "" || *readPath == "" || *allType == "" {
		fmt.Println("savePath, readPath or allType is nil")
	}
	err := HandleReadExcel(*readPath, *savePath, *allType)
	if err != nil {
		fmt.Printf("something err:%v\n", err)
	}
}

// 读取外部参数
func main() {
	runHandle()
}
