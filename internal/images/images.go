package images

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"sync"
)

func ProcessThumbnails(path string) {
	wg := sync.WaitGroup{}
	for _, i := range []int{1, 4} {
		wg.Add(1)
		fileName := fmt.Sprintf("%s/%d.jpeg", path, i)
		thumbnailName := fmt.Sprintf("%s/thumbnails/%d-sm.jpeg", path, i)
		go func() {
			defer wg.Done()
			createThumbnail(fileName, thumbnailName)
		}()
	}
	wg.Wait()
	fmt.Println("finished processing")
}

func createThumbnail(fileName, thumbnailName string) {
	imageGrid := makeImageGrid(fileName)
	resizedImageGrid := resize(imageGrid, 0.2)
	convertedImage := convert(resizedImageGrid)
	save(thumbnailName, convertedImage)
}

func save(filePath string, im *image.NRGBA) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jpeg.Encode(file, im.SubImage(im.Rect), nil)
}

func imageFromGrid(grid [][]color.Color) *image.NRGBA {
	rect := image.Rect(0, 0, len(grid), len(grid[0]))
	return &image.NRGBA{Rect: rect, Pix: make([]uint8, rect.Dx()*rect.Dy()*4), Stride: rect.Dx() * 4}
}

func convert(grid [][]color.Color) *image.NRGBA {
	newImage := imageFromGrid(grid)
	for x := 0; x < len(grid); x++ {
		for y := 0; y < len(grid[0]); y++ {
			q := grid[x]
			if q == nil {
				continue
			}
			p := grid[x][y]
			if p == nil {
				continue
			}
			original, ok := color.NRGBAModel.Convert(p).(color.NRGBA)
			if ok {
				newImage.Set(x, y, original)
			}
		}
	}
	return newImage
}

func makeImageGrid(fileName string) [][]color.Color {
	file, _ := os.Open(fileName)
	defer file.Close()
	im, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("bad image:", fileName)
		panic(err)
	}
	size := im.Bounds().Size()
	var grid [][]color.Color

	for i := 0; i < size.X; i++ {
		var y []color.Color
		for j := 0; j < size.Y; j++ {
			y = append(y, im.At(i, j))
		}
		grid = append(grid, y)
	}
	return grid
}

func resize(grid [][]color.Color, scale float64) [][]color.Color {
	xLen := int(float64(len(grid)) * scale)
	yLen := int(float64(len(grid[0])) * scale)
	resized := make([][]color.Color, xLen)
	for i := 0; i < len(resized); i++ {
		resized[i] = make([]color.Color, yLen)
	}
	for j := 0; j < xLen; j++ {
		for k := 0; k < yLen; k++ {
			jp := int(math.Floor(float64(j) / scale))
			kp := int(math.Floor(float64(k) / scale))
			resized[j][k] = grid[jp][kp]
		}
	}
	return resized
}
