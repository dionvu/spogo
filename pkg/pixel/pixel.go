package pixel

import (
	"image"
	"image/color"
	"net/http"
)

// FetchImageFromURL downloads an image from the provided URL.
func FetchImageFromURL(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Pixelate function to pixelate an image
func Pixelate(img image.Image, blockSize int) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	pixelatedImg := image.NewRGBA(bounds)

	for y := 0; y < h; y += blockSize {
		for x := 0; x < w; x += blockSize {
			avgColor := averageBlockColor(img, x, y, blockSize)
			fillBlock(pixelatedImg, x, y, blockSize, avgColor)
		}
	}

	return pixelatedImg
}

func averageBlockColor(img image.Image, startX, startY, blockSize int) color.Color {
	var r, g, b, count int
	bounds := img.Bounds()
	endX := min(startX+blockSize, bounds.Max.X)
	endY := min(startY+blockSize, bounds.Max.Y)

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			c := img.At(x, y)
			r1, g1, b1, _ := c.RGBA()
			r += int(r1)
			g += int(g1)
			b += int(b1)
			count++
		}
	}

	return color.RGBA{
		R: uint8(r / count >> 8),
		G: uint8(g / count >> 8),
		B: uint8(b / count >> 8),
		A: 255,
	}
}

func fillBlock(img *image.RGBA, startX, startY, blockSize int, c color.Color) {
	bounds := img.Bounds()
	endX := min(startX+blockSize, bounds.Max.X)
	endY := min(startY+blockSize, bounds.Max.Y)

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			img.Set(x, y, c)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
