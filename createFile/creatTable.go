package createfile

import (
	"fmt"
	"go_print_week6/models"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/fogleman/gg"
)

type CreatTable struct {
	rpath      string
	labelPath  string
	qrcodePath string
	logo       string
}

func NewCreateTable() (*CreatTable, error) {
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

	return &CreatTable{
		rpath:      dicPath,
		labelPath:  filepath.Join(dicPath, "label.png"),
		qrcodePath: filepath.Join(dicPath, "qrcode.png"),
		logo:       filepath.Join(dicPath, "logo.png"),
	}, nil
}

func (c *CreatTable) CreateTablePic(RowData map[string]interface{}, template models.Template, templateDtl []models.TemplateDtl) bool {
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 创建主图
	label_width, err := strconv.Atoi(template.LabelWidth)                        //宽
	label_height, err := strconv.Atoi(template.LabelHeight)                      //高
	white := color.RGBA{255, 255, 255, 255}                                      //颜色
	img := image.NewRGBA(image.Rect(0, 0, label_width, label_height))            //创建图像
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src) //确保图像的背景完全是白色（避免任何透明区域）

	// 检查操作系统类型并设置字体路径
	var fonts string
	if runtime.GOOS == "windows" {
		fonts = "C:/Windows/Fonts/msyh.ttc"
	} else if runtime.GOOS == "darwin" { // MacOS 的 `runtime.GOOS` 是 "darwin"
		fonts = "/System/Library/Fonts/Supplemental/Arial Unicode.ttf"
	}

	//创建了一个可以在图像上绘图的对象dc
	dc := gg.NewContextForRGBA(img)

	//加载了一种字体并将其大小设置为46（对于整个table）
	if err := dc.LoadFontFace(fonts, 46); err != nil {
		fmt.Println("加载字体错误")
		return false
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
	//线条宽度
	width_line := 5.0
	//标题
	label_mark := template.LabelMark
	//每一行的间距
	default_row := 100.0

	//文字距离上边的距离
	font_h := 20.0

	//开始元素距离画布顶端的距离
	default_h := 50.0
	row_w := 20.0

	row_h := default_h

	//第一条横线
	DrawLine(dc, row_w, row_h, float64(label_width)-row_w, row_h, color.Gray{Y: 128}, width_line)

	//标签标题
	drawText(dc, int(row_w+0.4*float64(label_width)), int(row_h+font_h), string(label_mark), "black", fonts, 50)

	//第二条横线
	row_h = row_h + default_row
	DrawLine(dc, row_w, row_h, float64(label_width)-row_w, row_h, color.Gray{Y: 128}, width_line)

	//计算后续的线条
	row_count := len(templateDtl)

	for _, row := range templateDtl {
		field := row.Field
		field_desc := row.FieldDesc

		str_data := RowData[field].(string)
		if len(str_data) > 28 {
			//添加内容
			drawText(dc, int(row_w+5), int(row_h+font_h), string(field_desc), "black", fonts, 46)

			str1 := str_data[0:28]
			str2 := str_data[28:]
			drawText(dc, int(float64(label_width)*0.24)+5, int(row_h+font_h), str1, "black", fonts, 46)
			drawText(dc, int(float64(label_width)*0.24)+5, int(row_h+font_h+default_row-5), str2, "black", fonts, 46)

			//画5条横线
			row_h = row_h + (default_row * 2)
			DrawLine(dc, row_w, row_h, float64(label_width)-row_w, row_h, color.Gray{Y: 128}, width_line)
			row_count = row_count + 1
		} else {

			//添加内容
			drawText(dc, int(row_w+5), int(row_h+font_h), string(field_desc), "black", fonts, 46)
			drawText(dc, int(float64(label_width)*0.24)+5, int(row_h+font_h), str_data, "black", fonts, 46)

			//画5条横线
			row_h = row_h + default_row
			DrawLine(dc, row_w, row_h, float64(label_width)-row_w, row_h, color.Gray{Y: 128}, width_line)
		}
	}
	//表格竖线
	DrawLine(dc, row_w, default_h, row_w, float64(row_count)*default_row+default_row+default_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, float64(label_width)-row_w*0.24, default_h+default_row, float64(label_width)-row_w*0.24, float64(row_count)*default_row+default_row+default_h, color.Gray{Y: 128}, width_line)
	DrawLine(dc, float64(label_width)-row_w, default_h, float64(label_width)-row_w, float64(row_count)*default_row+default_row+default_h, color.Gray{Y: 128}, width_line)

	//保存图片

	//返回
	return true
}
