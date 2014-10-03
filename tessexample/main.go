package main

import (
	"fmt"
	"log"

	"github.com/GeertJohan/go.leptonica"
	"github.com/GeertJohan/go.tesseract"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	// print the version
	fmt.Println(tesseract.Version())

	// create new tess instance and point it to the tessdata location. Set language to english.
	t, err := tesseract.NewTess("/usr/local/share/tessdata", "eng")
	if err != nil {
		log.Fatalf("Error while initializing Tess: %s\n", err)
	}
	defer t.Close()

	// open a new Pix from file with leptonica
	pix, err := leptonica.NewPixFromFile("./differentFonts.png")
	if err != nil {
		log.Fatalf("Error while getting pix from file: %s\n", err)
	}

	// set the page seg mode to autodetect
	t.SetPageSegMode(tesseract.PSM_AUTO_OSD)

	err = t.APISetVariable("tessedit_char_whitelist", ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~`+"`")
	if err != nil {
		log.Fatalf("Failed to APISetVariable: %s\n", err)
	}

	// set the image to the tesseract instance
	t.SetImagePix(pix)

	// retrieve text from the tesseract instance
	fmt.Println(t.Text())

	// // retrieve text from the tesseract instance
	// fmt.Println(t.HOCRText(0))

	// retrieve text from the tesseract instance
	fmt.Println(t.BoxText(0))

	// // retrieve text from the tesseract instance
	// fmt.Println(t.UNLVText())

	// dump variables for info
	// t.DumpVariables()

	spew.Dump(t.AvailableLanguages())
}
