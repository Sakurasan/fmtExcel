package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	_ "fmtExcel/routers"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "Path to an XLSX file")
var sheetIndex = flag.Int("i", 0, "Index of sheet to convert, zero based")
var delimiter = flag.String("d", ";", "Delimiter to use between fields")

func main() {
	beego.Run()

	// xlFile, error := xlsx.OpenFile("static/upload/excel9.xlsx")
	// if error != nil {
	// 	fmt.Println("打开失败")
	// 	return
	// }
	// sheetLen := len(xlFile.Sheets)
	// switch {
	// case sheetLen == 0:
	// 	fmt.Errorf("This XLSX file contains no sheets.")
	// 	return
	// 	// case  >= sheetLen:
	// 	// 	fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
	// 	// 	return
	// }

	// sheet := xlFile.Sheets[0]
	// var ColList []string
	// var ColsList [][]string
	// var temp [][]string

	// for index, row := range sheet.Rows {
	// 	if index == 0 {
	// 		continue
	// 	}
	// 	var vals []string
	// 	if row != nil {
	// 		for _, cell := range row.Cells {
	// 			str, err := cell.FormattedValue()
	// 			if err != nil {
	// 				vals = append(vals, err.Error())
	// 			}
	// 			vals = append(vals, fmt.Sprintf("%q", str))
	// 		}
	// 		// outputf(strings.Join(vals, *delimiter) + "\n")
	// 		fmt.Println("列？", vals)
	// 		if vals[1] != "" {
	// 			ColList = append(ColList, vals[0]) //过滤依据
	// 			ColsList = append(ColsList, vals)
	// 		}
	// 	}
	// }
	// fmt.Println("\n---------------\n第一列", ColList)

	// fmt.Println("===============================\n行集合", ColsList, "\n")
	// var kFlag string
	// kFlag = ColList[0] //命名标志
	// for i, v := range ColList {

	// 	if kFlag == v {
	// 		temp = append(temp, ColsList[i])
	// 		fmt.Println("归档T:", ColsList[i])

	// 	} else {
	// 		fmt.Println("标志位->", kFlag)
	// 		temp = nil
	// 		kFlag = v
	// 		temp = append(temp, ColsList[i])
	// 		fmt.Println("**************************\n归档F:", ColsList[i])

	// 	}
	// 	ExcelWriter(kFlag, temp)
	// }

	// ExcelWriter00()
	return

}

//将行的Cells直接读取成字符串数组
func GetCellValues(r *xlsx.Row) (cells []string) {
	for _, cell := range r.Cells {
		txt := strings.Replace(strings.Replace(cell.Value, "\n", "", -1), " ", "", -1)
		cells = append(cells, txt)
	}
	return
}

//读取列的数组
// func GteList(xlFile *xlsx.TimeToExcelTime) {

// }

// 写文件
func ExcelWriter(name string, exceldata [][]string) {
	f, err := os.Create(fmt.Sprintf(`%s.xls`, name[1:len(name)-1]))
	// f, err := os.Create("test.xls")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)
	w.Write([]string{"分公司代码", "", "流水日期", "交易金额", "商户手续费", "机构分润"})
	for i, v := range exceldata {
		if v != nil {
			w.Write(exceldata[i])
		}
	}
	// w.Write([]string{"1", "张三", "23"})

	w.Flush()
}

// func ExcelWriter00() {
// 	f, err := os.Create("test.xls")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
// 	w := csv.NewWriter(f)
// 	w.Write([]string{"编号", "姓名", "年龄"})
// 	w.Write([]string{"1", "张三", "23"})
// 	w.Write([]string{"2", "李四", "24"})
// 	w.Write([]string{"3", "王五", "25"})
// 	w.Write([]string{"4", "赵六", "26"})
// 	w.Flush()
// }

func IsExist(name string) bool {
	path := fmt.Sprintf("%s.xls", name)
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
	// 或者
	//return err == nil || !os.IsNotExist(err)
	// 或者
	//return !os.IsNotExist(err)
}
