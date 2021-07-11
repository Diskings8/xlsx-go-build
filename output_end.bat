@echo off
chcp 65001

set OUTTYPE="s"
:: 设置前端当前目录环境
set BASEPATH=%~dp0
echo %BASEPATH%

::生成objs.go文件的路径
set savePath=%BASEPATH%data
::目标excel文件路径

set readPath=%BASEPATH%xlsx

::所有的字段类型
set allType=int,IntSlice,IntSlice2,IntSlice3,IntMap,string,StringSlice,float64,Condition,Conditions,ItemInfo,ItemInfos,ItemInfosSlice,PropInfo,PropInfos,ProbItem,ProbItems,HmsTime,HmsTimes,Defenderweights
echo 开始生成objs.go文件 注意:生成文件时需要关闭对应的excel文件
::%readPath%\generateStruct -savePath=%savePath% -readPath=%readPath% -allType=%allType%
:: generateStruct -savePath=%savePath% -readPath=%readPath% -allType=%allType%
@set GROOT=%BASEPATH%;%GROOT%
Script.exe -savePath=%savePath% -readPath=%readPath% -allType=%allType% -outType=%OUTTYPE%

echo 生成完毕，按任意键继续
PAUSE
::
::	项目中所有的字段类型(注意区分类型的大小写)
::	int		整型 		例如:1
::	IntSlice		整型的一维数组 	例如:1,2,3
::	IntSlice2		整型的二维数组  	例如:[1,2,3;4,5,6]

::	string		字符型  		例如:"无生无灭炉"
::	StringSlice		字符型的一维数组  	例如:"无生无灭炉","雕花青铜炉"
::	float64		浮点型  		例如:1.5
