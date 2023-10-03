package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

func main() {

	// cli flags
	var folderPath string
	flag.StringVar(&folderPath, "f", "", "Path to image folder")
	flag.Parse()

	// no run without folder
	if folderPath == "" {
		fmt.Println("Please enter a valid path")
		return
	}

	// get some images
	pngFiles, err := getPNGFiles(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// basic canvas setup
	const canvasWidth = 500
	const canvasHeight = 500
	dc := gg.NewContext(canvasWidth, canvasHeight)
	dc.SetRGBA(0.0, 0.0, 0.0, 1.0)
	dc.Clear()

	// loop through images
	for _, pngFile := range pngFiles {
		// some variables
		alpha := 0.1 + rand.Float64()*0.9
		//BUG: scale := 0.2

		// load image files
		img, err := gg.LoadImage(filepath.Clean(filepath.Join(folderPath, pngFile)))
		if err != nil {
			fmt.Println("Error (loading out dir):", err)
			continue
		}

		img = applyOpacity(img, alpha)

		// add images and transform
		dc.Push()
		dc.SetRGBA(1, 1, 1, alpha)
		//BUG: dc.Scale(scale, scale)

		// randomise coordinates of images within canvas bounds
		x := rand.Float64() * (canvasWidth - float64(img.Bounds().Size().X))
		y := rand.Float64() * (canvasHeight - float64(img.Bounds().Size().Y))
		// DEBUG: fmt.Println("Coordinates (x, y): ", x, y)

		// rotate each image around its center
		dc.RotateAbout(rand.Float64()*8*gg.Radians(360), x+(float64(img.Bounds().Size().X)/2), y+(float64(img.Bounds().Size().Y)/2))

		// draw
		dc.DrawImage(img, int(x), int(y))

		// restore context
		dc.Pop()
	}

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Working directory: ", cwd)

	// create output dir for render if it doesn't exist
	outputDir := filepath.Join(cwd, "/render/")

	if err := os.Mkdir(outputDir, 0755); os.IsExist(err) {
		fmt.Println("Output dir found at: ", outputDir)
	} else {
		fmt.Println("Creating output dir at: ", outputDir)
	}

	// Save image in render with generic name
	outputPath := filepath.Join(outputDir, "output.png")

	// Check if file already exists and clean if needed
	fmt.Println("Running clean at: ", outputPath)
	if _, err := os.Stat(outputPath); err == nil {
		// If file exists then delete it
		if err := os.Remove(outputPath); err != nil {
			fmt.Println("Deleting existing output.png:", err)
			return
		}
	}

	// handle any errors on saving
	if err := dc.SavePNG(outputPath); err != nil {
		fmt.Println("Error saving...")
	} else {
		fmt.Println("Saved to: ", outputPath)
	}

}

// traverse a folder and grab png images
func getPNGFiles(folderPath string) ([]string, error) {
	var pngFiles []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".png" {
			pngFiles = append(pngFiles, info.Name())
		}
		return nil
	})
	return pngFiles, err
}

// take an image and change its opacity
func applyOpacity(img image.Image, alpha float64) image.Image {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			rgba.SetRGBA(x, y, color.RGBA{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(alpha * float64(a>>8)),
			})
		}
	}

	return rgba
}
