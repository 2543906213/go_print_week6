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

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

type Creat struct {
	dicpath      string
	batchid_Path string
	labelPath    string
	qrcodePath   string
	logo         string
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewCreate(batch_id string) (*Creat, error) {
	// 获取当前可执行文件的路径
	exePath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("无法获取当前文件路径：%v", err)
	}
	// fmt.Println(exePath)
	// 获取项目根目录的路径
	dicPath := filepath.Join(exePath) + string(os.PathSeparator) + "Files" + string(os.PathSeparator)
	batchid_Path := dicPath + batch_id
	// fmt.Println(batchid_Path)
	// fmt.Println(filepath.Join(batchid_Path, "qrcode.png"))
	// fmt.Println(filepath.Join(batchid_Path, "logo.png"))

	// 创建新文件夹
	err = os.MkdirAll(batchid_Path, 0755) // 0755 是文件夹的权限模式
	if err != nil {
		return nil, fmt.Errorf("创建文件夹失败:%v", err)
	}

	return &Creat{
		dicpath:      dicPath,
		batchid_Path: batchid_Path,
		labelPath:    filepath.Join(batchid_Path, "label.png"),
		qrcodePath:   filepath.Join(batchid_Path, "qrcode.png"),
		logo:         filepath.Join(dicPath, "logo.png"),
	}, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Creat) CreatePic(RowData map[string]interface{}, template models.Template, templateDtl []models.TemplateDtl) bool {
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

	//加载了一种字体并将其大小设置为50（对于bar_code）
	if err := dc.LoadFontFace(fonts, 50); err != nil {
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
	//绘制文本
	if templateDtl != nil {
		//一行的数量
		same_line_count := 0

		for _, row := range templateDtl {
			//将templateDtl里面的值全部取出来
			field := row.Field
			field_desc := row.FieldDesc
			desc_type, err := strconv.Atoi(row.DescType)
			ftype, err := strconv.Atoi(row.Type)
			field_left, err := strconv.Atoi(row.FieldLeft)
			field_up, err := strconv.Atoi(row.FieldUp)
			field_right, err := strconv.Atoi(row.FieldRight)
			//field_down, err := strconv.Atoi(row.FieldDown)
			same_line, err := strconv.Atoi(row.SameLine)

			//文字统一往后移一点
			field_left = int(field_left) + 25
			if same_line == 1 {
				same_line_count = same_line_count + 1
			} else {
				same_line_count = 0
			}

			//对于没有logo的标签需要整体将文字向下移
			if template.ShowLogo == "0" {
				field_up = field_up + 50
			}

			//绘制文本信息
			var font_size float64
			var rowText string

			//设置文本字体
			if desc_type == 0 {
				//python里面font_text包含了fonts和fonts_size两个参数，go里面不行
				template_FontsSize, err := strconv.Atoi(template.FontsSize)
				if err != nil {
					fmt.Println("template.FontsSize数据转换错误")
					return false
				}
				font_size = float64(template_FontsSize)
				if RowData[field] == nil {
					rowText = string(field_desc) + ": "
				} else {
					rowText = string(field_desc) + ": " + fmt.Sprintf("%v", RowData[field])
				}
				// fmt.Println(rowText)
			} else {
				font_size = 60
				if RowData[field] == nil {
					rowText = ""
				} else {
					rowText = fmt.Sprintf("%v", RowData[field])
				}
			}

			switch {
			case ftype == 1:
				if same_line_count > 0 {
					if field == "fdIqc" {
						field_up = field_up + 20
					}
					field_left = int(field_left) + 220*same_line_count
					drawText(dc, field_left, field_up, rowText, "black", fonts, font_size)
				} else {
					drawText(dc, field_left, field_up, rowText, "black", fonts, font_size)
				}
			case ftype == 2:
				if field == "bar_code" {
					drawText(dc, field_right+10, field_up+10, rowText, "black", fonts, 50)
				} else {
					drawText(dc, field_left, field_up, rowText, "black", fonts, font_size)
				}
			case ftype == 4 || ftype == 5:
				rowText = string(field_desc)
				drawText(dc, field_left, field_up, rowText, "black", fonts, font_size)
			}

			//保存图片
			file, err := os.Create(c.labelPath)
			if err != nil {
				fmt.Println("保存文字主图错误")
				return false
			}
			defer file.Close()
			err = png.Encode(file, img)
			if err != nil {
				log.Fatal(err)
				return false
			}
		}
	}
	fmt.Println("文本打印完毕！")
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//绘制条码

	for _, row := range templateDtl {
		//将templateDtl里面的值全部取出来
		field := row.Field
		//field_desc := row.FieldDesc
		//desc_type, err := strconv.Atoi(row.DescType)
		ftype, err := strconv.Atoi(row.Type)
		field_left, err := strconv.Atoi(row.FieldLeft)
		field_up, err := strconv.Atoi(row.FieldUp)
		field_right, err := strconv.Atoi(row.FieldRight)
		field_down, err := strconv.Atoi(row.FieldDown)
		//same_line, err := strconv.Atoi(row.SameLine)
		if err != nil {
			fmt.Println("数据转换错误")
			return false
		}

		if ftype == 2 && RowData[field] != "" {

			//生成条形码，获取条形码数据
			barcode_str := RowData[field].(string)
			// 选择条形码类型，这里使用 code128
			EAN, err := code128.Encode(barcode_str)
			if err != nil {
				fmt.Println("编码条形码时出错: ", err)
				return false
			}

			//设置条形码的大小，即在标签上画出放条形码的区域
			var box [4]int
			if field == "bar_code" {
				box = [4]int{field_left, field_up, field_right, field_down}
			} else {
				template_FontsSize, err := strconv.Atoi(template.FontsSize)
				if err != nil {
					fmt.Println("template.FontsSize数据转换错误")
					return false
				}
				if int(template_FontsSize) > 32 && int(template_FontsSize) < 36 {
					field_up = int(field_up) + 5
					field_down = int(field_down) + 5
				}
				if int(template_FontsSize) > 36 && int(template_FontsSize) < 44 {
					field_up = int(field_up) + 10
					field_down = int(field_down) + 10
				}
				if int(template_FontsSize) >= 44 {
					field_up = int(field_up) + 20
					field_down = int(field_down) + 20
				}
				box = [4]int{int(field_left), int(field_up) + 45, int(field_right), int(field_down) + 10}
			}

			//定义了保存条形码图像的文件路径和文件名
			pic_name := c.batchid_Path + string(os.PathSeparator) + "barcode-" + string(field) + ".png"

			// 创建一个文件来保存生成的条形码图像
			file, err := os.Create(pic_name)
			if err != nil {
				log.Fatal(err)
				return false
			}
			defer file.Close()

			//设置图片像素大小
			barcodeWidth := box[2] - box[0]
			barcodeHeight := box[3] - box[1] - 15

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
			defer source_img.Close()
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
			//设置条码图片的长宽
			// imgSize := tmp_img.Bounds() //图片的长和宽
			// maxSize := imgSize.Dx()     //图片的长
			// minSize := imgSize.Dy()     //图片的宽
			// //剪切图片
			// rect := image.Rect(0, 0, maxSize, minSize)           //先定义裁剪区域
			// region := image.NewRGBA(rect)                        // 创建一个新的空白图像来存储裁剪后的图像
			// draw.Draw(region, rect, tmp_img, rect.Min, draw.Src) // 将裁剪区域的像素复制到新图像
			// //将条码图片剪切出来之后重新设置大小，uint用于无符号整数类型，避免负数
			// resizedregion := resize.Resize(uint(box[2]-box[0]), uint(box[3]-box[1]), region, resize.Lanczos3)

			// 创建一个新的图像，用来存储最终的合成结果
			resultImg := image.NewRGBA(base_img.Bounds())
			// 将无条码的标签绘制到新的图像上
			draw.Draw(resultImg, base_img.Bounds(), base_img, image.Point{}, draw.Src)

			// 定义条码粘贴位置
			pastePointX := box[0]
			pastePointY := box[1] - 35
			//对于没有logo的标签需要将图片粘贴点向下移
			if template.ShowLogo == "0" {
				pastePointY = pastePointY + 50
			}
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
		}
	}
	fmt.Println("条码打印完毕!")
	//fmt.Println("有条码的标签图像已保存为", c.labelPath)

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//生成二维码
	if template.ShowQrcode == "10" {
		template_QrcodeLeft, err := strconv.Atoi(template.QrcodeLeft)
		template_QrcodeUp, err := strconv.Atoi(template.QrcodeUp)
		template_QrcodeRight, err := strconv.Atoi(template.QrcodeRight)
		template_QrcodeDown, err := strconv.Atoi(template.QrcodeDown)
		// 检查是否需要生成 QR 码
		var qrcodeData string
		if RowData["qrcode"] != "" {
			qrcodeData = RowData["qrcode"].(string)
		} else {
			qrcodeData = ""
		}
		// 生成 QR 码
		qr, err := qr.Encode(qrcodeData, 0, qr.Auto) //设置错误纠正级别为0，即 qr.Low
		if err != nil {
			fmt.Println("生成qr码失败: ", err)
			return false
		}
		// 将 QR 码调整为指定的大小
		qr, err = barcode.Scale(qr, 220, 220)
		if err != nil {
			fmt.Println("调整qr码大小失败: ", err)
			return false
		}
		// 添加边框
		border := 1
		imgBounds := qr.Bounds()
		newWidth := imgBounds.Dx() + 2*border
		newHeight := imgBounds.Dy() + 2*border
		// 创建一个新的图像作为背景
		newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		// 在新图像上绘制边框
		draw.Draw(newImg, newImg.Bounds(), image.NewUniform(image.White), image.Point{}, draw.Src)
		// 绘制 QR 码图像
		draw.Draw(newImg, image.Rect(border, border, imgBounds.Dx()+border, imgBounds.Dy()+border), qr, image.Point{}, draw.Over)

		// 保存合成后的二维码图像到指定路径
		outFile, err := os.Create(c.qrcodePath)
		if err != nil {
			fmt.Println("生成和保存 QR 码失败: ", err)
			return false
		}
		defer outFile.Close()
		// 将最终图像保存为 PNG 文件
		err = png.Encode(outFile, newImg)
		if err != nil {
			fmt.Println("无法保存二维码图像的PNG文件:", err)
			return false
		}
		// fmt.Println("QR 码图像已保存为", c.qrcodePath)

		//将二维码打印在标签上
		//打开放标签的图片
		source_img, err := os.Open(c.labelPath)
		if err != nil {
			fmt.Println("打开标签图片出错: ", err)
			return false
		}
		// 解码标签图片
		base_img, err := png.Decode(source_img)
		if err != nil {
			fmt.Println("无法解码标签图片: ", err)
			return false
		}
		//设置大小
		var box [4]int
		box = [4]int{int(template_QrcodeLeft), int(template_QrcodeUp), int(template_QrcodeRight), int(template_QrcodeDown)}

		//打开二维码图片
		tmp, err := os.Open(c.qrcodePath)
		if err != nil {
			fmt.Println("打开二维码图片出错: ", err)
			return false
		}
		defer tmp.Close()
		//解码二维码图片，为了获取长宽
		tmp_img, err := png.Decode(tmp)
		if err != nil {
			fmt.Println("无法解码二维码图片: ", err)
			return false
		}
		// //设置二维码图片的长宽
		// imgSize := tmp_img.Bounds() //图片的长和宽
		// maxSize := imgSize.Dx()     //图片的长
		// minSize := imgSize.Dy()     //图片的宽
		// //剪切图片
		// rect := image.Rect(0, 0, maxSize, minSize)           //先定义裁剪区域
		// region := image.NewRGBA(rect)                        // 创建一个新的空白图像来存储裁剪后的图像
		// draw.Draw(region, rect, tmp_img, rect.Min, draw.Src) // 将裁剪区域的像素复制到新图像
		// //将二维码图片剪切出来之后重新设置大小，uint用于无符号整数类型，避免负数
		// resizedregion := resize.Resize(uint(box[2]-box[0]), uint(box[3]-box[1]), region, resize.Lanczos3)

		// 创建一个新的图像，用来存储最终的合成结果
		resultImg := image.NewRGBA(base_img.Bounds())
		// 将无二维码的标签绘制到新的图像上
		draw.Draw(resultImg, base_img.Bounds(), base_img, image.Point{}, draw.Src)
		// 定义粘贴位置
		pastePosition := image.Pt(box[0], box[1])
		// 获取二维码图片的大小
		pasteRect := image.Rectangle{pastePosition, pastePosition.Add(tmp_img.Bounds().Size())}
		// 将二维码图像粘贴到无二维码的标签图像的指定位置
		draw.Draw(resultImg, pasteRect, tmp_img, image.Point{}, draw.Over)
		// 保存合成后的图像到指定路径
		outFile, err = os.Create(c.labelPath)
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
	}
	fmt.Println("二维码打印完毕!")

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// //绘制二维码
	// qr_code_num := 0
	// for _, row := range templateDtl {
	// 	//将templateDtl里面的值全部取出来
	// 	// field := row.Field
	// 	// field_desc := row.FieldDesc
	// 	// desc_type, err := strconv.Atoi(row.DescType)
	// 	ftype, err := strconv.Atoi(row.Type)
	// 	field_left, err := strconv.Atoi(row.FieldLeft)
	// 	field_up, err := strconv.Atoi(row.FieldUp)
	// 	// field_right, err := strconv.Atoi(row.FieldRight)
	// 	// field_down, err := strconv.Atoi(row.FieldDown)
	// 	// same_line, err := strconv.Atoi(row.SameLine)
	// 	if err != nil {
	// 		fmt.Println("数据转换错误")
	// 		return false
	// 	}
	// 	//try
	// 	var qrcodeData string
	// 	if RowData[qrcodeData] != "" {
	// 		qrcodeData = RowData[qrcodeData].(string)
	// 	} else {
	// 		qrcodeData = ""
	// 	}
	// 	if ftype == 3 {
	// 		template_QrcodeLeft, err := strconv.Atoi(template.QrcodeLeft)
	// 		template_QrcodeUp, err := strconv.Atoi(template.QrcodeUp)
	// 		//template_QrcodeRight, err := strconv.Atoi(template.QrcodeRight)
	// 		//template_QrcodeDown, err := strconv.Atoi(template.QrcodeDown)
	// 		// 生成 QR 码
	// 		qr, err := qr.Encode(qrcodeData, 0, qr.Auto) //设置错误纠正级别为0，即 qr.Low
	// 		if err != nil {
	// 			fmt.Println("生成qr码失败: ", err)
	// 			return false
	// 		}
	// 		// 将 QR 码调整为指定的大小，版本为2的二维码宽为5换算过来就是5*25
	// 		qr, err = barcode.Scale(qr, 125, 125)
	// 		if err != nil {
	// 			fmt.Println("调整qr码大小失败: ", err)
	// 			return false
	// 		}
	// 		// 添加边框
	// 		border := 1
	// 		imgBounds := qr.Bounds()
	// 		newWidth := imgBounds.Dx() + 2*border
	// 		newHeight := imgBounds.Dy() + 2*border
	// 		// 创建一个新的图像作为背景
	// 		newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	// 		// 在新图像上绘制边框
	// 		draw.Draw(newImg, newImg.Bounds(), image.NewUniform(image.White), image.Point{}, draw.Src)
	// 		// 绘制 QR 码图像
	// 		draw.Draw(newImg, image.Rect(border, border, imgBounds.Dx()+border, imgBounds.Dy()+border), qr, image.Point{}, draw.Over)

	// 		// 保存合成后的二维码图像到指定路径
	// 		outFile, err := os.Create(c.qrcodePath)
	// 		if err != nil {
	// 			fmt.Println("生成和保存 QR 码失败: ", err)
	// 			return false
	// 		}
	// 		defer outFile.Close()
	// 		// 将最终图像保存为 PNG 文件
	// 		err = png.Encode(outFile, newImg)
	// 		if err != nil {
	// 			fmt.Println("无法保存二维码图像的PNG文件:", err)
	// 			return false
	// 		}
	// 		fmt.Println("QR 码图像已保存为", c.qrcodePath)

	// 		//将二维码打印在标签上
	// 		//打开放标签的图片
	// 		source_img, err := os.Open(c.labelPath)
	// 		if err != nil {
	// 			fmt.Println("打开标签图片出错: ", err)
	// 			return false
	// 		}
	// 		// 解码标签图片
	// 		base_img, err := png.Decode(source_img)
	// 		if err != nil {
	// 			fmt.Println("无法解码标签图片: ", err)
	// 			return false
	// 		}
	// 		qr_code_num = qr_code_num + 1
	// 		//设置大小
	// 		var box [4]int
	// 		box = [4]int{int(field_left), int(field_up), int(field_left) + 200, int(field_up) + 200}
	// 		// 打开二维码图片
	// 		tmp, err := os.Open(c.qrcodePath)
	// 		if err != nil {
	// 			fmt.Println("打开二维码图片出错: ", err)
	// 			return false
	// 		}
	// 		// 解码二维码图片，为了获取长宽
	// 		tmp_img, err := png.Decode(tmp)
	// 		if err != nil {
	// 			fmt.Println("无法解码二维码图片: ", err)
	// 			return false
	// 		}
	// 		// 设置二维码图片的长宽
	// 		imgSize := tmp_img.Bounds() //图片的长和宽
	// 		maxSize := imgSize.Dx()     //图片的长
	// 		minSize := imgSize.Dy()     //图片的宽
	// 		// 剪切图片
	// 		rect := image.Rect(0, 0, maxSize, minSize) //先定义裁剪区域

	// 		region := image.NewRGBA(rect)                        // 创建一个新的空白图像来存储裁剪后的图像
	// 		draw.Draw(region, rect, tmp_img, rect.Min, draw.Src) // 将裁剪区域的像素复制到新图像
	// 		//将二维码图片剪切出来之后重新设置大小，uint用于无符号整数类型，避免负数
	// 		resizedregion := resize.Resize(uint(box[2]-box[0]), uint(box[3]-box[1]), region, resize.Lanczos3)
	// 		// 创建一个新的图像，用来存储最终的合成结果
	// 		resultImg := image.NewRGBA(base_img.Bounds())
	// 		// 将无二维码的标签绘制到新的图像上
	// 		draw.Draw(resultImg, base_img.Bounds(), base_img, image.Point{}, draw.Src)
	// 		// 定义粘贴位置
	// 		pastePosition := image.Pt(int(template_QrcodeLeft), int(template_QrcodeUp))
	// 		// 获取二维码图片的大小
	// 		pasteRect := image.Rectangle{pastePosition, pastePosition.Add(resizedregion.Bounds().Size())}
	// 		// 将二维码图像粘贴到无二维码的标签图像的指定位置
	// 		draw.Draw(resultImg, pasteRect, resizedregion, image.Point{}, draw.Over)
	// 		// 保存合成后的图像到指定路径
	// 		outFile, err = os.Create(c.labelPath)
	// 		if err != nil {
	// 			fmt.Println("无法创建输出文件: ", err)
	// 			return false
	// 		}
	// 		defer outFile.Close()
	// 		// 将最终图像保存为 PNG 文件
	// 		err = png.Encode(outFile, resultImg)
	// 		if err != nil {
	// 			fmt.Println("无法保存有条码的PNG文件:", err)
	// 			return false
	// 		}
	// 		fmt.Println("标签图像已保存为", c.labelPath)
	// 	}
	// 	if template.ShowQrcode == "10" {
	// 		if qr_code_num == 2 {
	// 			break
	// 		}
	// 	} else {
	// 		if qr_code_num == 3 {
	// 			break
	// 		}
	// 	}
	// }
	// fmt.Println("二维码打印成功!")
	// fmt.Println("有二维码的标签图像已保存为", c.labelPath)
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//生成logo
	if template.ShowLogo == "10" {
		template_LogoLeft, err := strconv.Atoi(template.LogoLeft)
		template_LogoUp, err := strconv.Atoi(template.LogoUp)
		template_LogoRight, err := strconv.Atoi(template.LogoRight)
		template_LogoDown, err := strconv.Atoi(template.LogoDown)
		logo := template.Logo
		if logo == "" {
			logo = "logo.png"
		}
		c.logo = c.dicpath + logo

		//打开放标签的图片
		source_img, err := os.Open(c.labelPath)
		if err != nil {
			fmt.Println("打开标签图片出错: ", err)
			return false
		}
		defer source_img.Close()
		// 解码标签图片
		base_img, err := png.Decode(source_img)
		if err != nil {
			fmt.Println("无法解码标签图片: ", err)
			return false
		}
		//设置大小
		var box [4]int
		box = [4]int{int(template_LogoLeft), int(template_LogoUp), int(template_LogoRight), int(template_LogoDown)}

		// 打开logo图片
		tmp, err := os.Open(c.logo)
		if err != nil {
			fmt.Println("打开logo图片出错: ", err)
			return false
		}
		defer tmp.Close()
		// 解码logo图片，为了获取长宽
		tmp_img, err := png.Decode(tmp)
		if err != nil {
			fmt.Println("无法解码logo图片: ", err)
			return false
		}
		// 设置logo图片的长宽
		imgSize := tmp_img.Bounds() //图片的长和宽
		maxSize := imgSize.Dx()     //图片的长
		minSize := imgSize.Dy()     //图片的宽

		// 剪切图片
		rect := image.Rect(0, 0, maxSize, minSize)           //先定义裁剪区域
		region := image.NewRGBA(rect)                        // 创建一个新的空白图像来存储裁剪后的图像
		draw.Draw(region, rect, tmp_img, rect.Min, draw.Src) // 将裁剪区域的像素复制到新图像
		//将logo图片剪切出来之后重新设置大小，uint用于无符号整数类型，避免负数
		resizedregion := resize.Resize(uint(box[2]-box[0])-60, uint(box[3]-box[1])-15, region, resize.Lanczos3)
		// 创建一个新的图像，用来存储最终的合成结果
		resultImg := image.NewRGBA(base_img.Bounds())
		// 将无logo的标签绘制到新的图像上
		draw.Draw(resultImg, base_img.Bounds(), base_img, image.Point{}, draw.Src)
		// 定义粘贴位置
		pastePosition := image.Pt(int(template_LogoLeft), int(template_LogoUp))
		// 获取logo图片的大小
		pasteRect := image.Rectangle{pastePosition, pastePosition.Add(resizedregion.Bounds().Size())}
		// 将logo图像粘贴到无logo的标签图像的指定位置
		draw.Draw(resultImg, pasteRect, resizedregion, image.Point{}, draw.Over)
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
		fmt.Println("logo打印完毕!")
		// fmt.Println("有logo的标签图像已保存为", c.labelPath)
		outFile.Close()
		
	}
	return true
}
