package imgutil

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const BLACK = 0
const WHITE = 1
const TRANSPARENT = 2

type ImgViewer struct {
	img [][]int
	h   int
	w   int
}

func (c *ImgViewer) Draw(scale int) error {
	black := color.RGBA{0, 0, 0, 1}
	white := color.RGBA{255, 255, 255, 1}
	transparent := color.RGBA{0, 0, 0, 0}
	w := c.w
	h := c.h

	rect := image.Rect(0, 0, w*scale, h*scale)
	cimg := image.NewGray(rect)
	draw.Draw(cimg, cimg.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			color := transparent
			if c.img[x][y] == BLACK {
				color = black
			} else if c.img[x][y] == WHITE {
				color = white
			}
			draw.Draw(cimg, image.Rect(x*scale, y*scale, x*scale+scale, y*scale+scale), &image.Uniform{color}, image.ZP, draw.Src)
		}
	}

	out, err := os.Create("./result.png")
	if err != nil {
		return err
	}
	err = png.Encode(out, cimg)
	if err != nil {
		return err
	}
	return nil
}

func Parse(img string, w, h int) *ImgViewer {
	layerLength := w * h
	layerCount := len(img) / layerLength
	parsedImg := make([][]int, w)
	for i := 0; i < w; i++ {
		parsedImg[i] = make([]int, h)
		for j := 0; j < h; j++ {
			parsedImg[i][j] = TRANSPARENT
		}
	}

	for pixel := 0; pixel < layerLength; pixel++ {
		colorValue := TRANSPARENT
		for layer := 0; layer < layerCount; layer++ {
			color := img[layer*layerLength+pixel]
			if color == '0' {
				colorValue = BLACK
				break
			} else if color == '1' {
				colorValue = WHITE
				break
			}
		}
		parsedImg[pixel%w][pixel/w] = colorValue
	}

	return &ImgViewer{
		img: parsedImg,
		h:   h,
		w:   w,
	}
}
