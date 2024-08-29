package createfile

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime"

	"unsafe"

	"github.com/nfnt/resize"
	"golang.org/x/sys/windows"
)

type Print struct {
}

func (p *Print) PrintFile(printfile string, printerModel string) {

	// 加白框保证图片确保在打印范围的中间
	// 打开原始图像文件
	file, err := os.Open(printfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// 解码图像
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	// 获取原始图像的尺寸
	imgSize := img.Bounds()        //图片的长和宽
	originalWidth := imgSize.Dx()  //图片的宽
	originalHeight := imgSize.Dy() //图片的高

	// 框的宽度a
	a := 15 // 可以根据需要调整

	// 计算缩小后的目标宽度和高度，保持宽高比
	scaleFactor := float64(originalWidth-2*a) / float64(originalWidth)
	targetWidth := uint(float64(originalWidth) * scaleFactor)
	targetHeight := uint(float64(originalHeight) * scaleFactor)

	// 缩小图像（保持宽高比）
	resizedImg := resize.Resize(targetWidth, targetHeight, img, resize.Lanczos3)

	// 创建带有白色背景的最终图像
	finalImg := image.NewRGBA(image.Rect(0, 0, originalWidth, originalHeight))
	white := color.RGBA{255, 255, 255, 255} // 白色
	draw.Draw(finalImg, finalImg.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// 计算缩小后的图像在最终图像中的居中位置
	offsetX := (originalWidth - int(targetWidth)) / 2
	offsetY := (originalHeight - int(targetHeight)) / 2

	// 将缩小后的图像绘制到带白色背景的最终图像上
	draw.Draw(finalImg, image.Rect(offsetX, offsetY, offsetX+int(targetWidth), offsetY+int(targetHeight)), resizedImg, image.Point{}, draw.Over)

	//将图像旋转90度
	rotatedImg := rotate90(finalImg)

	// 保存最终图像
	outFile, err := os.Create(printfile)
	if err != nil {
		log.Fatal(err)
	}

	// 将最终图像保存为 PNG 文件
	err = png.Encode(outFile, rotatedImg)
	if err != nil {
		fmt.Println("无法保存有条码的PNG文件:", err)
	}
	// fmt.Println("有logo的标签图像已保存为", c.labelPath)
	outFile.Close()

	// 使用命令行打印图片
	if runtime.GOOS == "windows" {
		if err := printPngFile(printfile, printerModel); err != nil {
			fmt.Println("打印PNG文件失败:", err)
		}
	} else {
		fmt.Errorf("该系统不支持打印: %s", runtime.GOOS)
	}
}

// 利用win系统自带的图像处理功能，执行PNG文件打印操作
func printPngFile(filePath, printerName string) error {
	cmd := exec.Command("rundll32", "C:\\Windows\\System32\\shimgvw.dll,ImageView_PrintTo", filePath, printerName)
	return cmd.Run()
}

// 获取默认打印机
const (
	spoolerDLL = "winspool.drv"
)

var (
	modwinspool           = windows.NewLazyDLL(spoolerDLL)
	procGetDefaultPrinter = modwinspool.NewProc("GetDefaultPrinterW")
)

func (p *Print) GetDefaultPrinter() (printerModel string, err error) {
	var size uint32 = 512
	buffer := make([]uint16, size)
	ret, _, err := procGetDefaultPrinter.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return "", err
	}
	return windows.UTF16ToString(buffer), nil
	// return printerModel, err
}
