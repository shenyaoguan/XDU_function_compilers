package graphics

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// DrawPlot 绘制一个简单的图形
func DrawPlot(x, y float64) {
	// 创建一个 500x500 的图像
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	// 设置背景为白色
	for i := 0; i < 500; i++ {
		for j := 0; j < 500; j++ {
			img.Set(i, j, color.White)
		}
	}
	// 绘制一个点，坐标 (x, y)
	img.Set(int(x), int(y), color.Black)

	// 保存图像
	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}
