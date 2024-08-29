package createfile

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/fogleman/gg"
)

type CreatTableBox struct {
	rpath      string
	labelPath  string
	qrcodePath string
	logo       string
}

func NewCreatTableBox() (*CreatTableBox, error) {
	// 获取当前可执行文件的路径
	exePath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("无法获取当前文件路径：%v", err)
	}
	fmt.Println(exePath)
	// 获取项目根目录的路径
	dicPath := filepath.Join(exePath) + string(os.PathSeparator) + "Files" + string(os.PathSeparator)
	fmt.Println(dicPath)
	fmt.Println(filepath.Join(dicPath, "qrcode.png"))
	fmt.Println(filepath.Join(dicPath, "logo.png"))

	return &CreatTableBox{
		rpath:      dicPath,
		labelPath:  filepath.Join(dicPath, "label.png"),
		qrcodePath: filepath.Join(dicPath, "qrcode.png"),
		logo:       filepath.Join(dicPath, "logo.png"),
	}, nil
}

func (c *CreatTableBox) CreatTableBoxPic(page_data []map[string]interface{}, printData []string, idx int, page int, serial_number int) bool {

	// 创建主图
	white := color.RGBA{255, 255, 255, 255}                                      //颜色
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1500))                           //创建图像
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src) //确保图像的背景完全是白色（避免任何透明区域）

	// win10-字体
	// 检查操作系统类型并设置字体路径
	var fonts string
	// var fonts_code string
	if runtime.GOOS == "windows" {
		fonts = "C:/Windows/Fonts/msyh.ttc"
		// fonts_code = "C:/Windows/Fonts/msyhbd.ttc"
	} else if runtime.GOOS == "darwin" { // MacOS 的 `runtime.GOOS` 是 "darwin"
		fonts = "/System/Library/Fonts/Supplemental/Arial Unicode.ttf"
	}

	//创建了一个可以在图像上绘图的对象dc
	dc := gg.NewContextForRGBA(img)

	//加载了一种字体并将其大小设置为22（对于整个TableBox）
	if err := dc.LoadFontFace(fonts, 22); err != nil {
		fmt.Println("加载字体错误")
		return false
	}

	//线条宽度
	width_line := 2.0

	//获取概要数据
	var sumy_data map[string]interface{}
	idx_sumy := 0
	for _, row := range printData {
		idx_sumy = idx_sumy + 1
		if idx_sumy == idx {
			err := json.Unmarshal([]byte(row), &sumy_data)
			if err != nil {
				fmt.Println("JSON解析失败:", err)
				return false
			}
		}
	}

	box_no := sumy_data["box_no"]
	box_num := sumy_data["box_num"]
	customs_code := sumy_data["customs_code"]
	customs_box_id := sumy_data["customs_box_id"]
	customs_company := sumy_data["customs_company"]
	customs_company_code := sumy_data["customs_company_code"]
	var entry_code interface{}
	if sumy_data["entry_code"] != nil {
		entry_code = sumy_data["entry_code"]
	} else {
		entry_code = nil
	}

	//保存图片，路径为c.labelPath
	file, err := os.Create(c.labelPath)
	if err != nil {
		fmt.Println("保存空白主图错误")
		return false
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
		return false
	}
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	page_desc := fmt.Sprintf("第 %d 页", page)

	//标题
	barcode_desc := customs_code.(string) + box_no.(string)
	title := customs_company.(string) + customs_company_code.(string)
	drawText(dc, 460, 50, title, "black", fonts, 30)

	//箱号
	drawText(dc, 820, 90, box_num.(string), "black", fonts, 32)

	//页码
	drawText(dc, 820, 130, page_desc, "black", fonts, 28)

	//条码文字
	drawText(dc, 40, 130, barcode_desc, "black", fonts, 28)

	//表头
	idx_desc := "序号"
	pro_sno_desc := "型号"
	pro_fagent_desc := "厂牌"
	num_desc := "数量"
	production_place_desc := "产地"
	entryCode_desc := "入仓单号"

	//表格横线
	h := 150.0
	row_h := 50.0
	for range page_data {
		h = h + row_h
		DrawLine(dc, 40, h, 960, h, color.Gray{Y: 128}, width_line)
	}

	//表头
	drawText(dc, 45, 180, idx_desc, "black", fonts, 22)
	drawText(dc, 180, 180, pro_sno_desc, "black", fonts, 22)
	drawText(dc, 400, 180, pro_fagent_desc, "black", fonts, 22)
	drawText(dc, 560, 180, num_desc, "black", fonts, 22)
	drawText(dc, 660, 180, production_place_desc, "black", fonts, 22)
	drawText(dc, 800, 180, entryCode_desc, "black", fonts, 22)

	data_h := 155
	data_code_h := 155

	str_len := 14
	for key, value := range page_data {
		serial_number = serial_number + 1
		data_h = int(data_h) + int(row_h)
		data_code_h = int(data_code_h) + int(row_h)
		data_h = int(data_h)
		data_code_h = int(data_code_h)
		fmt.Println(key, value)
		row_data_h := int(data_h) + 30

		serial_number_str := fmt.Sprintf("%d", serial_number)
		drawText(dc, 50, row_data_h, serial_number_str, "black", fonts, 22)
		drawText(dc, 550, row_data_h, value["num"].(string), "black", fonts, 22)
		drawText(dc, 625, row_data_h, value["production_place"].(string), "black", fonts, 24)
		drawText(dc, 750, row_data_h, entry_code.(string), "black", fonts, 28)

		pro_sno := value["pro_sno"].(string)
		pro_sno_len := len(pro_sno)
		if pro_sno_len > str_len {
			pro_sno_one := pro_sno[0:str_len]
			pro_sno_two := pro_sno[str_len:]
			drawText(dc, 95, row_data_h-10, pro_sno_one, "black", fonts, 22)
			drawText(dc, 95, row_data_h+12, pro_sno_two, "black", fonts, 22)
		} else {
			drawText(dc, 95, row_data_h, pro_sno, "black", fonts, 22)
		}

		pro_fagent := value["pro_fagent"].(string)
		pro_fagent_len := len(pro_fagent)
		if pro_fagent_len > str_len {
			pro_fagent_one := pro_fagent[0:str_len]
			pro_fagent_two := pro_fagent[str_len:]
			drawText(dc, 320, row_data_h-10, pro_fagent_one, "black", fonts, 22)
			drawText(dc, 320, row_data_h+12, pro_fagent_two, "black", fonts, 22)
		} else {
			drawText(dc, 320, row_data_h, pro_fagent, "black", fonts, 22)
		}
	}

	//底部坐标
	table_h := float64(len(page_data))*row_h + 210

	//外部框
	DrawLine(dc, 40, 150, 960, 150, color.Gray{Y: 128}, width_line)
	DrawLine(dc, 40, 150, 40, table_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, 960, 150, 960, table_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, 40, table_h, 960, table_h, color.Gray{Y: 128}, width_line)

	//列竖线
	DrawLine(dc, 90, 150, 90, table_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, 310, 150, 310, table_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, 540, 150, 540, table_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, 620, 150, 620, table_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, 740, 150, 740, table_h, color.Gray{Y: 128}, width_line)

	// 保存图片，路径为c.labelPath
	file, err = os.Create(c.labelPath)
	if err != nil {
		fmt.Println("保存空白主图错误")
		return false
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// 绘制条码
	customs_box_id_str := customs_box_id.(string)
	// 选择条形码类型，这里使用 code128
	EAN, err := code128.Encode(customs_box_id_str)
	if err != nil {
		fmt.Println("编码条形码时出错: ", err)
		return false
	}
	pic_name := c.rpath + "barcode-" + customs_box_id_str

	//图片放入的区域，坐标点：左 上 右 下
	var box [4]int
	box = [4]int{30, 20, 460, 100}

	// 创建一个文件来保存生成的条形码图像
	file, err = os.Create(pic_name)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()

	//设置图片像素大小
	barcodeWidth := box[2] - box[0]
	barcodeHeight := box[3] - box[1]

	//设置为图片像素大小的图片
	rowean, _ := barcode.Scale(EAN, barcodeWidth, barcodeHeight)

	//条码图片左边距
	leftMargin := 30
	//裁剪条形码的静区（空白区域）
	ean := cropQuietZone(rowean, leftMargin)

	// 将图像编码为 PNG 格式并保存到文件
	err = png.Encode(file, ean)
	if err != nil {
		log.Fatal(err)
		return false
	}

	//打开放标签的图片
	source_img, err := os.Open(c.labelPath)
	if err != nil {
		fmt.Println("打开标签图片出错: ", err)
		return false
	}
	//解码标签图片
	base_img, err := png.Decode(source_img)
	if err != nil {
		fmt.Println("无法解码标签图片: ", err)
		return false
	}
	//打开条码图片
	tmp, err := os.Open(pic_name)
	if err != nil {
		fmt.Println("打开条码图片出错: ", err)
		return false
	}
	//解码条码图片，为了获取长宽
	tmp_img, err := png.Decode(tmp)
	if err != nil {
		fmt.Println("无法解码条码图片: ", err)
		return false
	}

	// 创建一个新的图像，用来存储最终的合成结果
	resultImg := image.NewRGBA(base_img.Bounds())
	// 将无条码的标签绘制到新的图像上
	draw.Draw(resultImg, base_img.Bounds(), base_img, image.Point{}, draw.Src)

	// 定义条码粘贴位置
	pastePointX := box[0]
	pastePointY := box[1]

	pastePosition := image.Pt(pastePointX, pastePointY)

	// 获取条码图片的大小
	pasteRect := image.Rectangle{pastePosition, pastePosition.Add(tmp_img.Bounds().Size())}
	// 将条码图像粘贴到无条码的标签图像的指定位置
	draw.Draw(resultImg, pasteRect, tmp_img, image.Point{}, draw.Over)
	// 保存合成后的图像到指定路径
	outFile, err := os.Create(c.labelPath)
	if err != nil {
		fmt.Println("无法创建输出文件: ", err)
		return false
	}
	defer outFile.Close()
	// 将最终图像保存为 PNG 文件
	err = png.Encode(outFile, resultImg)
	if err != nil {
		fmt.Println("无法保存有条码的PNG文件:", err)
		return false
	}

	// //等返回数据了再继续写page_data和RowData数据类型有问题 8.22
	return true
}
