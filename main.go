package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
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
	//training()
	argsWithProg := os.Args[1:]
	fmt.Println(argsWithProg[0])
	fmt.Println(fromPy(argsWithProg[0]))

}

func forMelbie(path string) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	imagePrediction(file, "", "")
}

func fromPy(path string) string {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	result, _ := imagePrediction(file, "", "")
	return result
}

func training() {
	done := timed("Running Program")
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	folder, err := os.ReadDir("./train")
	if err != nil {
		fmt.Println("Error reading training directory. Make sure it exists.")
		os.Exit(1)
	}

	total := 0
	totalFailed := 0
	tally := make(map[string][]int)
	for _, entry := range folder {
		/* DEBUG CODE */
		//if entry.Name() != "3210" {
		//	continue
		//}

		if entry.IsDir() {
			subfolder, err := os.ReadDir(fmt.Sprintf("./train/%s", entry.Name()))
			if err != nil {
				fmt.Println("Error reading sub folder in training directory.")
				os.Exit(1)
			}

			count := 0
			failed := 0
			var wg sync.WaitGroup

			for _, testImage := range subfolder {
				/* DEBUG CODE */
				//if testImage.Name() != "00005.png" && testImage.Name() != "00004.png" && testImage.Name() != "00003.png" && testImage.Name() != "00002.png" && testImage.Name() != "00001.png" {
				//	continue
				//}
				wg.Add(1)
				count += 1
				go func(testImage os.DirEntry) {
					defer wg.Done()
					file, err := os.Open(fmt.Sprintf("./train/%s/%s", entry.Name(), testImage.Name()))

					if err != nil {
						fmt.Println("Error: File could not be opened")
						os.Exit(1)
					}

					_, failure := imagePrediction(file, entry.Name(), testImage.Name())
					if failure {
						failed += 1
					}
				}(testImage)
			}

			wg.Wait()

			tally[entry.Name()] = []int{count, failed}
			total += count
			totalFailed += failed
		}

		fmt.Println()
		fmt.Print("Foldername :  Count, Failed")
		fmt.Println()
		for key, element := range tally {
			fmt.Print(key, " : ")
			fmt.Println(element[0], " ", element[1])
		}
	}

	fmt.Println("\nTotal tests ran: ", total)
	fmt.Println("Total tests failed: ", totalFailed)

	done()
}

func imagePrediction(file io.Reader, expected string, imageName string) (string, bool) {
	debugString := fmt.Sprintf("\n\nNEW TEST: %s/%s\n", expected, imageName)

	pixels, err := getPixels(file)
	//Pixels = [0:127][0:127]=

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

	comparisonSides := make(map[string][]float64)

	comparisonSides["tl-tr"] = comparePixelArray(TopLeft.right, TopRight.left)
	comparisonSides["tl-br"] = comparePixelArray(TopLeft.right, BottomRight.left)
	comparisonSides["tl-bl"] = comparePixelArray(TopLeft.right, BottomLeft.left)

	comparisonSides["tr-tl"] = comparePixelArray(TopRight.right, TopLeft.left)
	comparisonSides["tr-br"] = comparePixelArray(TopRight.right, BottomRight.left)
	comparisonSides["tr-bl"] = comparePixelArray(TopRight.right, BottomLeft.left)

	comparisonSides["bl-tl"] = comparePixelArray(BottomLeft.right, TopLeft.left)
	comparisonSides["bl-tr"] = comparePixelArray(BottomLeft.right, TopRight.left)
	comparisonSides["bl-br"] = comparePixelArray(BottomLeft.right, BottomRight.left)

	comparisonSides["br-tl"] = comparePixelArray(BottomRight.right, TopLeft.left)
	comparisonSides["br-tr"] = comparePixelArray(BottomRight.right, TopRight.left)
	comparisonSides["br-bl"] = comparePixelArray(BottomRight.right, BottomLeft.left)

	/* DEBUG CODE */
	debugString += "\n\nKey :  Difference, Std Dev\n"
	for key, element := range comparisonSides {
		debugString += key + " : "
		debugString += fmt.Sprintf("%f", element[0]) + " " + fmt.Sprintf("%f", element[1]) + "\n"
	}

	keys := make([]float64, len(comparisonSides))

	i := 0
	for _, k := range comparisonSides {
		keys[i] = k[1] - k[0]
		i++
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	comparison2 := make(map[float64]string)
	for key, element := range comparisonSides {
		comparison2[(element[1] - element[0])] = key
	}

	biggestKey := comparison2[keys[0]]
	usedTile1 := biggestKey[:2]
	usedTile2 := biggestKey[3:]

	secondKey := ""
	for idx := range keys[1:] {
		if !strings.Contains(comparison2[keys[idx]], usedTile1) && !strings.Contains(comparison2[keys[idx]], usedTile2) {
			secondKey = comparison2[keys[idx]]
			break
		}
	}

	if secondKey == "" || biggestKey == "" {
		//fmt.Println(fmt.Sprintf("\n\n\n %s/%s", imageName, expected))
		//fmt.Println(comparison2)
		//fmt.Println(secondKey)
		//fmt.Println(biggestKey)
		//fmt.Println(keys)
		//fmt.Println(secondBiggest)
		//fmt.Println(comparisonSides)
		return "", true
	}

	var Top Subimage
	var Bottom Subimage

	switch secondKey[:2] {
	case "tl":
		Top.top = append(Top.top, TopLeft.top...)
		Top.bottom = append(Top.bottom, TopLeft.bottom...)
	case "tr":
		Top.top = append(Top.top, TopRight.top...)
		Top.bottom = append(Top.bottom, TopRight.bottom...)
	case "bl":
		Top.top = append(Top.top, BottomLeft.top...)
		Top.bottom = append(Top.bottom, BottomLeft.bottom...)
	case "br":
		Top.top = append(Top.top, BottomRight.top...)
		Top.bottom = append(Top.bottom, BottomRight.bottom...)
	}

	switch secondKey[3:] {
	case "tl":
		Top.top = append(Top.top, TopLeft.top...)
		Top.bottom = append(Top.bottom, TopLeft.bottom...)
	case "tr":
		Top.top = append(Top.top, TopRight.top...)
		Top.bottom = append(Top.bottom, TopRight.bottom...)
	case "bl":
		Top.top = append(Top.top, BottomLeft.top...)
		Top.bottom = append(Top.bottom, BottomLeft.bottom...)
	case "br":
		Top.top = append(Top.top, BottomRight.top...)
		Top.bottom = append(Top.bottom, BottomRight.bottom...)
	}

	switch biggestKey[:2] {
	case "tl":
		Bottom.top = append(Bottom.top, TopLeft.top...)
		Bottom.bottom = append(Bottom.bottom, TopLeft.bottom...)
	case "tr":
		Bottom.top = append(Bottom.top, TopRight.top...)
		Bottom.bottom = append(Bottom.bottom, TopRight.bottom...)
	case "bl":
		Bottom.top = append(Bottom.top, BottomLeft.top...)
		Bottom.bottom = append(Bottom.bottom, BottomLeft.bottom...)
	case "br":
		Bottom.top = append(Bottom.top, BottomRight.top...)
		Bottom.bottom = append(Bottom.bottom, BottomRight.bottom...)
	}

	switch biggestKey[3:] {
	case "tl":
		Bottom.top = append(Bottom.top, TopLeft.top...)
		Bottom.bottom = append(Bottom.bottom, TopLeft.bottom...)
	case "tr":
		Bottom.top = append(Bottom.top, TopRight.top...)
		Bottom.bottom = append(Bottom.bottom, TopRight.bottom...)
	case "bl":
		Bottom.top = append(Bottom.top, BottomLeft.top...)
		Bottom.bottom = append(Bottom.bottom, BottomLeft.bottom...)
	case "br":
		Bottom.top = append(Bottom.top, BottomRight.top...)
		Bottom.bottom = append(Bottom.bottom, BottomRight.bottom...)
	}

	comparisonTops := make(map[string][]float64)

	comparisonTops["secondKey"] = comparePixelArray(Top.bottom, Bottom.top)
	comparisonTops["biggestKey"] = comparePixelArray(Bottom.bottom, Top.top)

	/* DEBUG CODE */
	debugString += "\n\nKey :  Difference, Std Dev\n"
	for key, element := range comparisonTops {
		if key == "secondKey" {
			debugString += secondKey + " : "
		} else {
			debugString += biggestKey + " : "
		}
		debugString += fmt.Sprintf("%f", element[0]) + " " + fmt.Sprintf("%f", element[1]) + "\n"
	}

	result := ""

	correctAns := ""
	if (comparisonTops["secondKey"][1] - comparisonTops["secondKey"][0]) > (comparisonTops["biggestKey"][1] - comparisonTops["biggestKey"][0]) {
		result = output(secondKey, biggestKey)
		correctAns += secondKey + " " + biggestKey
	} else {
		result = output(biggestKey, secondKey)
		correctAns += biggestKey + " " + secondKey
	}

	correctString := Melbie(comparisonSides["br-tl"]) + Melbie(comparisonSides["br-bl"]) + Melbie(comparisonSides["tl-tr"])
	correctString += Melbie(comparisonSides["tr-br"]) + Melbie(comparisonSides["bl-tl"]) + Melbie(comparisonSides["bl-tr"])
	correctString += Melbie(comparisonSides["bl-br"]) + Melbie(comparisonSides["br-tr"]) + Melbie(comparisonSides["tl-br"])
	correctString += Melbie(comparisonSides["tl-bl"]) + Melbie(comparisonSides["tr-tl"]) + Melbie(comparisonSides["tr-bl"]) + correctAns + "\n"
	f, _ := os.OpenFile("correct.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if _, err := f.Write([]byte(correctString)); err != nil {
		fmt.Println(err)
	}

	if result == expected {
		return result, false
	} else {
		//True if invalid
		if getMiddle5x5(result, pixels) {
			//Remove this result from possible results and try again
			if !getMiddle5x5(reverseStr(result), pixels) { //if valid
				return reverseStr(result), result == expected
			} else if !getMiddle5x5(result[2:]+result[:2], pixels) {
				return result[2:] + result[:2], result == expected
			} else if !getMiddle5x5(reverseStr(result[2:]+result[:2]), pixels) {
				return reverseStr(result[2:] + result[:2]), result == expected
			} else if !getMiddle5x5(string(result[0])+string(result[2])+string(result[1])+string(result[3]), pixels) {
				return string(result[0]) + string(result[2]) + string(result[1]) + string(result[3]), result == expected
			} else if !getMiddle5x5(string(result[3])+string(result[1])+string(result[2])+string(result[0]), pixels) {
				return string(result[3]) + string(result[1]) + string(result[2]) + string(result[0]), result == expected
			} else if !getMiddle5x5(string(result[0])+string(result[3])+string(result[1])+string(result[2]), pixels) {
				return string(result[0]) + string(result[3]) + string(result[1]) + string(result[2]), result == expected
			} else if !getMiddle5x5(string(result[2])+string(result[1])+string(result[3])+string(result[0]), pixels) {
				return string(result[2]) + string(result[1]) + string(result[3]) + string(result[0]), result == expected
			} else if !getMiddle5x5(string(result[2])+string(result[0])+string(result[3])+string(result[1]), pixels) {
				return string(result[2]) + string(result[0]) + string(result[3]) + string(result[1]), result == expected
			}

		}

		//debugString += fmt.Sprintf("FAILED TEST: expected:%s | result:%s", expected, result)
		//fmt.Println(debugString)
		return result, true
	}

}

// MISC
func reverseStr(str string) (result string) {
	// iterate over str and prepend to result
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func Melbie(element []float64) string {
	return fmt.Sprintf("%f", element[0]) + "," + fmt.Sprintf("%f", element[1]) + ","
}

//END MISC

func output(s1 string, s2 string) string {
	out := ""
	s := s1 + "-" + s2
	//TL-TR-RL-RB
	i := strings.Index(s, "tl") / 3
	out += strconv.FormatInt(int64(i), 10)

	i = strings.Index(s, "tr") / 3
	out += strconv.FormatInt(int64(i), 10)

	i = strings.Index(s, "bl") / 3
	out += strconv.FormatInt(int64(i), 10)

	i = strings.Index(s, "br") / 3
	out += strconv.FormatInt(int64(i), 10)

	return out
}

func reverseOutput(expected string) string {
	s := ""
	idx := strings.Index(expected, "0")
	if idx == 0 {
		s += "tl"
	} else if idx == 1 {
		s += "tr"
	} else if idx == 2 {
		s += "bl"
	} else {
		s += "br"
	}
	s += "-"

	idx = strings.Index(expected, "1")
	if idx == 0 {
		s += "tl"
	} else if idx == 1 {
		s += "tr"
	} else if idx == 2 {
		s += "bl"
	} else {
		s += "br"
	}
	s += " "

	idx = strings.Index(expected, "2")
	if idx == 0 {
		s += "tl"
	} else if idx == 1 {
		s += "tr"
	} else if idx == 2 {
		s += "bl"
	} else {
		s += "br"
	}
	s += "-"

	idx = strings.Index(expected, "3")
	if idx == 0 {
		s += "tl"
	} else if idx == 1 {
		s += "tr"
	} else if idx == 2 {
		s += "bl"
	} else {
		s += "br"
	}

	return s
}

func (img *Subimage) countPixels() uint {
	//fmt.Println(len(img.top) + len(img.bottom) + len(img.left) + len(img.right))
	return uint(len(img.top) + len(img.bottom) + len(img.left) + len(img.right))
}

func stdDev(pixels []float64) float64 {
	var sum, mean, sd float64

	for _, num := range pixels {
		sum += num
	}
	mean = sum / float64(len(pixels))

	for _, num := range pixels {
		// The use of Pow math function func Pow(x, y float64) float64
		sd += math.Pow(num-mean, 2)
	}
	// The use of Sqrt math function func Sqrt(x float64) float64
	sd = math.Sqrt(sd / 10)

	return math.Round(sd*100) / 100
}

// Returns Array [Difference, Std Dev]
func comparePixelArray(s1 []Pixel, s2 []Pixel) []float64 {
	var sum float64
	for idx := range s1 {
		sum += similarPixels(s1[idx], s2[idx])
	}

	var s1Grayed []float64
	for _, pixel := range s1 {
		gray := 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)
		s1Grayed = append(s1Grayed, gray)
	}

	var s2Grayed []float64
	for _, pixel := range s2 {
		gray := 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)
		s2Grayed = append(s2Grayed, gray)
	}

	return []float64{math.Round(sum/64*100) / 100, (stdDev(s1Grayed) + stdDev(s2Grayed)) / 2}
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

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s...\n", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

// Grab middle + of pixels and compare, we want low difference and low standard deviation on non white but we need to filter out white space.
// Call middle plus 5x5 of each corner in middle.
func getMiddle5x5(output string, pixels [][]Pixel) bool {
	var tl, tr, bl, br [5][5]Pixel
	if output[0] == '3' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tl[i][j] = pixels[i][j]
			}
		}
	} else if output[0] == '2' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tl[i][j] = pixels[i][j+59]
			}
		}
	} else if output[0] == '1' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tl[i][j] = pixels[i+59][j]
			}
		}
	} else {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tl[i][j] = pixels[i+59][j+59]
			}
		}
	}

	// tr
	if output[1] == '3' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tr[i][j] = pixels[i][j+64]
			}
		}
	} else if output[1] == '2' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tr[i][j] = pixels[i+64][j]
			}
		}
	} else if output[1] == '1' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tr[i][j] = pixels[i][j+123]
			}
		}
	} else {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				tr[i][j] = pixels[i+59][j+123]
			}
		}
	}

	//bl
	if output[2] == '3' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				bl[i][j] = pixels[i+64][j]
			}
		}
	} else if output[2] == '2' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				bl[i][j] = pixels[i+64][j+59]
			}
		}
	} else if output[2] == '1' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				bl[i][j] = pixels[i+123][j]
			}
		}
	} else {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				bl[i][j] = pixels[i+123][j+59]
			}
		}
	}

	//br
	if output[3] == '3' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				br[i][j] = pixels[i+64][j+64]
			}
		}
	} else if output[3] == '2' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				br[i][j] = pixels[i+64][j+123]
			}
		}
	} else if output[3] == '1' {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				br[i][j] = pixels[i+123][j+64]
			}
		}
	} else {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				br[i][j] = pixels[i+123][j+123]
			}
		}
	}

	return checkInner(tl, tr, bl, br)
}

// Takes in 4 boxes // returns true IF NOT VALID
func checkInner(tl [5][5]Pixel, tr [5][5]Pixel, bl [5][5]Pixel, br [5][5]Pixel) bool {
	tlRes := comparePixelMatrix(tl)
	trRes := comparePixelMatrix(tr)
	blRes := comparePixelMatrix(bl)
	brRes := comparePixelMatrix(br)

	mean := (tlRes[0] + trRes[0] + blRes[0] + brRes[0]) / 4.0
	variance := math.Sqrt((math.Pow(mean-tlRes[0], 2) + math.Pow(mean-trRes[0], 2) + math.Pow(mean-blRes[0], 2) + math.Pow(mean-brRes[0], 2)) / 4.0)
	//fmt.Println(variance)

	//20 = 7270
	//15 = 7658
	//17.5 = 7363
	if variance > mean*.20 {
		//fmt.Println("ITS INVALID!!!")
		return true
	}

	return false
}

func comparePixelMatrix(s1 [5][5]Pixel) []float64 {
	var avg float64
	for _, row := range s1 {
		for _, pixel := range row {
			avg += 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)
		}
	}
	avg /= 25.0

	var s1Grayed []float64
	for _, row := range s1 {
		for _, pixel := range row {
			gray := 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)
			s1Grayed = append(s1Grayed, gray)
		}
	}

	return []float64{math.Round(avg*100) / 100, stdDev(s1Grayed)}
}
