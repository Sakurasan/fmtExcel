package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.TplName = "main.html"
}

func (this *MainController) Post() {
	this.TplName = "main.html"
	cate := this.Input().Get("cate")
	namecate := this.Input().Get("namecate")
	timefile := time.Now().Format("20060102150405")

	f, h, err := this.GetFile("excel")
	if err != nil {
		log.Println("getfile err ", err)
		this.Redirect("/", 302)
	}
	defer f.Close()
	// 设置保存目录
	mpath := "static/upload/" + timefile + "/"
	// 创建目录
	os.MkdirAll(mpath, 0755)
	fpath := mpath + h.Filename
	for i := 0; ; i++ {
		_, err = os.Stat(fpath)
		if err == nil {
			fpath = mpath + strconv.Itoa(i) + h.Filename
		} else {
			break
		}
	}

	this.SaveToFile("excel", fpath) //保存位置在 static/upload, 没有文件夹要先创建
	fmt.Println("-----------------------")
	this.Data["App"] = fmt.Sprintf("过滤依据%v列\n命名依据%v列\n,上传文件为%v", cate, namecate, h.Filename)

	catecode, _ := strconv.Atoi(this.Input().Get("cate"))
	namecatecode, _ := strconv.Atoi(this.Input().Get("namecate"))
	f.Close()
	fmtExcel(catecode, namecatecode, timefile, h.Filename)
	fmt.Println("cmd ok?")
	// exec_shell("./static/upload/tar.sh " + timefile)
	exec_shell(fmt.Sprintf("./static/upload/tar.sh ./static/upload/%v ", timefile))
	fmt.Println("cmd is OK!")
	this.Data["Url"] = timefile
}

func Ziper([]string) {
	// 创建一个缓冲区用来保存压缩文件内容
	buf := new(bytes.Buffer)

	// 创建一个压缩文档
	w := zip.NewWriter(buf)

	// 将文件加入压缩文档
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}

	// 关闭压缩文档
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 将压缩文档内容写入文件
	f, err := os.OpenFile("file.zip", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	buf.WriteTo(f)
}

//格式化 Excel
func fmtExcel(num, name int, timepath, filename string) {
	// xlFile, error := xlsx.OpenFile("static/upload/" + pathname + "/" + filename)
	xlFile, error := xlsx.OpenFile(fmt.Sprintf("static/upload/%v/%v", timepath, filename))
	// xlFile, error := xlsx.OpenFile("static/upload/9.xlsx")
	if error != nil {
		fmt.Println("打开失败")
		return
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		fmt.Errorf("This XLSX file contains no sheets.")
		return
		// case  >= sheetLen:
		// 	fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
		// 	return
	}

	sheet := xlFile.Sheets[0]
	var ColList []string //第一 列表
	var ColsList [][]string
	var temp [][]string //

	for index, row := range sheet.Rows {
		if index == 0 {
			continue
		}
		var vals []string
		if row != nil {
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					vals = append(vals, err.Error())
				}
				vals = append(vals, fmt.Sprintf("%s", str)) //
			}
			// outputf(strings.Join(vals, *delimiter) + "\n")
			fmt.Println("单行？", vals)
			if vals[1] != "" {
				ColList = append(ColList, vals[num]) //筛选列,标志位
				ColsList = append(ColsList, vals)    //筛选容器,二维数组
			}
		}
	}
	fmt.Println("\n---------------\n列,单一行", ColList)

	fmt.Println("===============================\n行集合", ColsList, "\n")
	var kFlag string
	kFlag = ColList[name] //命名标志
	for i, v := range ColList {

		if kFlag == v {
			temp = append(temp, ColsList[i])
			fmt.Println("归档T:", ColsList[i])

		} else {
			fmt.Println("标志位->", kFlag)
			temp = nil
			kFlag = v
			temp = append(temp, ColsList[i])
			fmt.Println("**************************\n归档F:", ColsList[i])
		}
		ExcelWriter(timepath, kFlag, temp)
	}

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
func ExcelWriter(path, name string, exceldata [][]string) {

	fmt.Println("写文件->", path, "+", name)
	// fmt.Println("static/update/" + path + "/" + name + ".xls")
	// fmt.Println("路径->", fmt.Sprintf(`static/update/%v/%s.xls`, path, name))
	f, err := os.Create(fmt.Sprintf(`static/upload/%v/%s.csv`, path, name)) //name[1:len(name)-1]
	// f, err := os.Create(path + name + ".xls")
	// f, err := os.Create("static/update/"+path+"/"+ name + ".xls")
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

func convCode(code string) int {
	if code == "A" {
		return 0
	} else if code == "B" {
		return 1
	}
	return 0
}

func exec_shell(s string) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}
