package main

import (
	"fmt"
	"github.com/GeertJohan/go.leptonica"
	"github.com/GeertJohan/go.tesseract"
	"log"
)

func main() {
	fmt.Println(tesseract.Version())

	t, err := tesseract.NewTess("/usr/local/share/tessdata", "nld")
	if err != nil {
		log.Fatalf("Error while initializing Tess: %s\n", err)
	}
	pix, err := leptonica.NewPixFromFile("/home/geertjohan/input.jpg")
	if err != nil {
		log.Fatalf("Error while getting pix from file: %s\n", err)
	}
	t.SetImagePix(pix)
	// t.SetInputName("/home/geertjohan/input.jpg")

	fmt.Println(t.GetText())

	t.DumpVariables()
}
