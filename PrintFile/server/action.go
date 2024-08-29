package server

import (
	"encoding/json"
	"fmt"
	createfile "go_print_week6/createFile"
	"go_print_week6/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 包含 PrintFile 属性，用于存储根路径和打印文件路径。
type Action struct {
	exePath   string
	PrintFile string
}

// 打印结果
type PrintResult struct {
	status  int
	message string
}

// 1. 初始化方法设定了根路径和打印文件的路径。
func (a *Action) PrintPath(batch_id string) {
	// 获取当前文件的根路径
	exePath, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取到当前文件路径：", err)
		return
	}

	// 创建一个新路径，设置为打印文件的路径
	a.exePath = exePath
	a.PrintFile = filepath.Join(exePath, "Files", batch_id, "label.png")
	// fmt.Println("打印文件路径:", a.PrintFile)
}

func (a *Action) PrintCode(batch_id string) PrintResult {

	result := PrintResult{0, ""}

	//2. 解析 batch_id 以获取打印批次的信息。
	// 使用 strings.Split 函数将 batch_id 按照 "GZHUADOUPRINT" 分隔符进行分割。
	// batchInfo 变量保存分割后的字符串切片。
	batchInfo := strings.Split(batch_id, "GZHUADOUPRINT")
	batchInfoLen := len(batchInfo)
	// 初始化变量
	batch_id = batchInfo[0]
	printName := ""
	// gzFontSize := 0
	// 检查 batchInfo 的长度，以确定是否存在第二部分的字符串。
	// 根据条件设置 printName 和 gzFontSize 变量的值。
	if batchInfoLen > 1 {
		printName = batchInfo[1]
		// gzFontSize = 1
	}
	// 打印当前时间和id
	// fmt.Println(printName, gzFontSize)
	// fmt.Println("batch_id为: ", batch_id)
	// fmt.Printf("1: %s: %s \n", time.Now().Format("2006-01-02 15:04:05"), batch_id)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//3. 通过 HttpConnetion 类从服务器获取批次数据
	url := "http://erp-internal-api.ickey.com.cn/barcode/print-code/print?id=" + batch_id

	pdata, err := http.Get(url) //从接口获取信息
	if err != nil {
		result = PrintResult{0, "无法获取打印信息"}
		fmt.Println("无法获取打印信息: ", err)
		return result
	}
	defer pdata.Body.Close()
	//pdata, err := http.Fetch(url)
	if err != nil {
		fmt.Printf("%s:%s:接口连接异常", time.Now().Format("2006-01-02 15:04:05"), batch_id)
		result = PrintResult{0, "接口连接超时"}
		fmt.Println(result)
		return result
	}

	// 读取响应体数据
	body, err := io.ReadAll(pdata.Body)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
	}
	// 关闭响应体
	defer pdata.Body.Close()

	//4.匹配 JSON 数据结构、解析 JSON 数据
	var data models.PrintData

	// 解析 JSON 数据
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println("JSON解析失败:", err)
		return result
	}
	if len(data.Printdata) == 0 &&
		data.Template == (models.Template{}) &&
		len(data.TemplateDtl) == 0 &&
		len(data.BoxData) == 0 &&
		data.Box == "" &&
		data.PrintNum == 0 {
		fmt.Println("打印数据为空")
	}
	//fmt.Println("JSON解析成功:", data)

	//5.处理JSON数据，将JSON数据分类
	printData := data.Printdata
	template := data.Template
	templateDtl := data.TemplateDtl

	//6.处理print_num、label_style、box_data、box、template_name
	print_num := data.PrintNum
	if data.PrintNum == 0 {
		print_num = 1
	}

	label_style := template.LabelStyle
	if template.LabelStyle == "" {
		label_style = "60"
	}
	box_data := data.BoxData
	//box := data.Box
	template_name := template.TemplateName

	//7.根据标签类型分别处理
	switch template_name {
	case "EMS快递面单":
		label_style = "100"
	case "德邦快递面单":
		label_style = "100"
	case "顺丰快递面单":
		label_style = "100"
	case "地区异常标识面单":
		label_style = "100"
	case "到货异常条码模板":
		label_style = "100"
	case "装箱单打印模板":
		page_size := 25
		idx := 1
		for _, value := range box_data {
			var RowData []map[string]interface{}
			//等返回数据了再继续写page_data和RowData数据类型有问题 8.2
			err = json.Unmarshal([]byte(value), &RowData)
			if err != nil {
				fmt.Println("JSON解析失败:", err)
				return result
			}
			//当前行号
			var row_index int
			var page_data []map[string]interface{}
			for _, item := range RowData {
				row_index++
				page_data = append(page_data, item)

				//计算页码（20条一页）
				var page int
				if row_index%page_size == 0 {
					page = int((row_index) / page_size)
				} else {
					page = int(int(row_index)/page_size + 1)
				}
				var serial_number int
				if row_index == len(RowData) || row_index%page_size == 0 {
					c, err := createfile.NewCreatTableBox()
					if err != nil {
						fmt.Println(err)
					}
					serial_number = row_index - len(page_data)

					// page_data 列表数据
					// print_data 概要数据
					// idx 当前批次的一个打印任务
					// page 计算页码
					// serial_number 序号
					res := c.CreatTableBoxPic(page_data, printData, idx, page, serial_number)
					page_data = nil
					if res {
						//连接打印机
						p := createfile.Print{}
						printName, err = p.GetDefaultPrinter()
						if err != nil {
							fmt.Println("获取打印机失败 \n")
						}
						fmt.Printf("默认打印机为：%s \n", printName)
						//驱动打印机打印
						for i := 0; i < print_num; i++ {
							p.PrintFile(a.PrintFile, printName)
						}
					}
				}
			}
			batchid_Path := filepath.Join(a.exePath, "Files", batch_id)
			err = os.RemoveAll(batchid_Path)
			if err != nil {
				fmt.Printf("无法删除文件夹 %s \n", batchid_Path)
			}
		}
	default:
		switch label_style {
		//常规面模版
		case "10":
			for _, index := range printData {
				var RowData map[string]interface{}
				err = json.Unmarshal([]byte(index), &RowData)
				if err != nil {
					fmt.Println("JSON解析失败:", err)
					return result
				}
				c, err := createfile.NewCreate(batch_id)
				if err != nil {
					fmt.Println(err)
				}
				res := c.CreatePic(RowData, template, templateDtl)
				if res {
					// fmt.Printf("3: %s :%s \n", time.Now().Format("2006-01-02 15:04:05"), batch_id)
					//连接打印机
					p := createfile.Print{}
					printName, err = p.GetDefaultPrinter()
					if err != nil {
						fmt.Println("获取打印机失败 \n")
					}
					fmt.Printf("默认打印机为：%s \n", printName)
					//驱动打印机打印
					for i := 0; i < print_num; i++ {
						p.PrintFile(a.PrintFile, printName)
						// fmt.Printf("4: %s :%s \n", time.Now().Format("2006-01-02 15:04:05"), batch_id)
					}
				}
				batchid_Path := filepath.Join(a.exePath, "Files", batch_id)
				err = os.RemoveAll(batchid_Path)
				if err != nil {
					fmt.Printf("无法删除文件夹 %s %v \n", batchid_Path, err)
				}
			}
		//常规表格
		case "20":
			for _, index := range printData {
				var RowData map[string]interface{}
				err = json.Unmarshal([]byte(index), &RowData)
				if err != nil {
					fmt.Println("JSON解析失败:", err)
					return result
				}
				c, err := createfile.NewCreateTable()
				if err != nil {
					fmt.Println(err)
				}
				res := c.CreateTablePic(RowData, template, templateDtl)
				if res {
					fmt.Printf("3: %s :%s \n", time.Now().Format("2006-01-02 15:04:05"), batch_id)
					//连接打印机
					p := createfile.Print{}
					printName, err = p.GetDefaultPrinter()
					if err != nil {
						fmt.Println("获取打印机失败 \n")
					}
					fmt.Printf("默认打印机为：%s \n", printName)
					//驱动打印机打印
					for i := 0; i < print_num; i++ {
						p.PrintFile(a.PrintFile, printName)
						fmt.Printf("4: %s :%s \n", time.Now().Format("2006-01-02 15:04:05"), batch_id)
					}
				}
				batchid_Path := filepath.Join(a.exePath, "Files", batch_id)
				err = os.RemoveAll(batchid_Path)
				if err != nil {
					fmt.Printf("无法删除文件夹 %s \n", batchid_Path)
				}
			}
		//自定义表格
		case "30":
		}
	}
	//更新打印状态
	//url = "https://erp-internal-api.ickey.com.cn/barcode/print-code/update-status?print_type=10&id=" + batch_id
	//http.Fetch(url)
	result = PrintResult{1, "打印成功！"}
	fmt.Println(result)

	//返回数据
	return result
}
