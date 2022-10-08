package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"os"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

type Subimage struct {
	top    []Pixel
	left   []Pixel
	right  []Pixel
	bottom []Pixel
}

func main() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open("./00000.png")

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := getPixels(file)
	//Pixels = [0:127][0:127]

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	//Grab 63 | 64
	var TopLeft Subimage
	var TopRight Subimage
	var BottomLeft Subimage
	var BottomRight Subimage

	for idxY, row := range pixels {
		for idxX, pixel := range row {
			if idxY == 0 {
				if idxX < 64 {
					TopLeft.top = append(TopLeft.top, pixel)
				} else {
					TopRight.top = append(TopRight.top, pixel)
				}
			}

			if idxY == 64 {
				if idxX < 64 {
					BottomLeft.top = append(BottomLeft.top, pixel)
				} else {
					BottomRight.top = append(BottomRight.top, pixel)
				}
			}

			if idxY == 63 {
				if idxX < 64 {
					TopLeft.bottom = append(TopLeft.bottom, pixel)
				} else {
					TopRight.bottom = append(TopRight.bottom, pixel)
				}
			}

			if idxY == 127 {
				if idxX < 64 {
					BottomLeft.bottom = append(BottomLeft.bottom, pixel)
				} else {
					BottomRight.bottom = append(BottomRight.bottom, pixel)
				}
			}

			if idxX == 63 {
				if idxY < 64 {
					TopLeft.right = append(TopLeft.right, pixel)
				} else {
					BottomLeft.right = append(BottomLeft.right, pixel)
				}
			}

			if idxX == 64 {
				if idxY < 64 {
					TopRight.left = append(TopRight.left, pixel)
				} else {
					BottomRight.left = append(BottomRight.left, pixel)
				}
			}

			if idxX == 0 {
				if idxY < 64 {
					TopLeft.left = append(TopLeft.left, pixel)
				} else {
					BottomLeft.left = append(BottomLeft.left, pixel)
				}
			}

			if idxX == 127 {
				if idxY < 64 {
					TopRight.right = append(TopRight.right, pixel)
				} else {
					BottomRight.right = append(BottomRight.right, pixel)
				}
			}
		}
	}

	//fmt.Println(TopLeft)
	//fmt.Println(TopRight)
	//fmt.Println(BottomLeft)
	//fmt.Println(BottomRight)
	//TopLeft.countPixels()
	//TopRight.countPixels()
	//BottomLeft.countPixels()
	//BottomRight.countPixels()

	//similarPixels(TopLeft.right[0], TopLeft.right[50])
	fmt.Println(compareSides(TopLeft.right, TopRight.left))

}

func (img *Subimage) countPixels() uint {
	fmt.Println(len(img.top) + len(img.bottom) + len(img.left) + len(img.right))
	return uint(len(img.top) + len(img.bottom) + len(img.left) + len(img.right))
}

func compareSides(s1 []Pixel, s2 []Pixel) float64 {
	var sum float64
	for idx := range s1 {
		sum += similarPixels(s1[idx], s2[idx])
	}

	return sum / 64
}

/* Percent difference between two pixels */
func similarPixels(p1 Pixel, p2 Pixel) float64 {
	algo1 := math.Sqrt(math.Pow(float64(p1.R-p2.R), 2)+math.Pow(float64(p1.G-p2.G), 2)+math.Pow(float64(p1.B-p2.B), 2)) / 441.67 * 100

	return algo1
}

func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	//Optimizations down the road, we can set this to a predefined width since we know images are 128x128 we can prevent reslicing
	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}
