package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
)

var grayImg *image.Gray

func RGBAtoGray(pix color.Color) uint8 {
	r, g, b, _ := pix.RGBA()
	return uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
}

func numInRange(num, bottom, top uint8) bool {
	return bottom < num && num < top
}

func PhotoToAscii(grayImg *image.Gray) string {
	result := ""
	for y := grayImg.Bounds().Min.Y; y < grayImg.Bounds().Max.Y; y++ {
		for x := grayImg.Bounds().Min.X; x < grayImg.Bounds().Max.X; x++ {
			switch pixel := RGBAtoGray(grayImg.At(x, y)); true {
			case numInRange(pixel, 0, 30):
				result += " "
			case numInRange(pixel, 30, 50):
				result += ";"
			case numInRange(pixel, 50, 75):
				result += "."
			case numInRange(pixel, 75, 125):
				result += "o"
			case numInRange(pixel, 125, 175):
				result += "+"
			case numInRange(pixel, 175, 200):
				result += "*"
			case numInRange(pixel, 200, 255):
				result += "M"
			case pixel > 225:
				result += "8"
			}
		}

		result += "\n"
	}

	return result
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", PhotoToAscii(grayImg))
}

func main() {
	var file *os.File
	var err error

	if len(os.Args) > 1 {
		file, err = os.Open(os.Args[1])
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		file, err = os.Open(scanner.Text())
	}

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	grayImg = image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	// Web server
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}
