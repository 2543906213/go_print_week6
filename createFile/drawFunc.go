package createfile

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/fogleman/gg"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 用于绘制文本内容的函数
func drawText(dc *gg.Context, x, y int, rowText string, color string, fontPath string, fontSize float64) {
	//加载字体
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Printf("加载字体失败: %v", err)
		return
	}
	// 设置颜色
	switch color {
	case "black":
		dc.SetRGB(0, 0, 0)
	case "white":
		dc.SetRGB(255, 255, 255)
	case "red":
		dc.SetRGB(255, 0, 0)
	case "green":
		dc.SetRGB(0, 255, 0)
	case "blue":
		dc.SetRGB(0, 0, 255)
	default:
		dc.SetRGB(0, 0, 0) // 默认颜色为黑色
	}
	//绘制文本
	dc.DrawString(rowText, float64(x), float64(y))
	dc.Stroke()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// DrawLine 函数接受起点和终点的坐标、颜色和线条宽度作为参数，并在给定的图像上绘制一条线
func DrawLine(dc *gg.Context, x1, y1, x2, y2 float64, col color.Color, width float64) {
	// 设置线条颜色
	dc.SetColor(col)
	// 设置线条宽度
	dc.SetLineWidth(width)
	// 绘制线条
	dc.DrawLine(x1, y1, x2, y2)
	// 填充线条
	dc.Stroke()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// cropQuietZone 检测并裁剪条形码图像中的空白区域，仅在左边保留边距，右边去掉所有空白区域
func cropQuietZone(img image.Image, leftMargin int) image.Image {
	bounds := img.Bounds()
	left, right := bounds.Max.X, bounds.Min.X

	// 遍历图像的每一列，检测实际内容的边界
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if img.At(x, y) != color.White { // 检测非白色像素
				if x < left {
					left = x
				}
				if x > right {
					right = x
				}
			}
		}
	}

	// 保留左边的边距
	if left-leftMargin < bounds.Min.X {
		left = bounds.Min.X
	} else {
		left -= leftMargin
	}

	// 确保裁剪后的边界不会超出图像的实际范围
	if left < bounds.Min.X {
		left = bounds.Min.X
	}

	// 右边界从实际内容的边界开始
	if right > bounds.Max.X {
		right = bounds.Max.X
	}

	// 裁剪出条形码内容区域
	croppedBounds := image.Rect(left, bounds.Min.Y, right, bounds.Max.Y)
	croppedImg := image.NewRGBA(croppedBounds)
	draw.Draw(croppedImg, croppedImg.Bounds(), img, image.Point{X: left, Y: bounds.Min.Y}, draw.Src)

	return croppedImg
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// resizeImage 将图像调整为指定的宽度和高度
func resizeImage(img image.Image, width, height int) image.Image {
	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(resizedImg, resizedImg.Bounds(), image.White, image.Point{}, draw.Src)

	// 将条形码内容图像调整为目标尺寸
	draw.Draw(resizedImg, resizedImg.Bounds(), img, img.Bounds().Min, draw.Src)

	return resizedImg
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// rotate90 将图像顺时针旋转90度
func rotate90(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建一个新的图像，宽度和高度调换
	newImg := image.NewRGBA(image.Rect(0, 0, height, width))

	// 顺时针旋转图像90度
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImg.Set(y, width-1-x, img.At(x, y))
		}
	}

	return newImg
}
