package main

import (
	"fmt"
	"github.com/GeertJohan/go.tesseract"
	"log"
)

func main() {
	fmt.Println(tesseract.Version())

	t, err := tesseract.NewTess("/usr/local/share/tessdata", "nld")
	if err != nil {
		log.Fatalf("Error while initializing Tess: %s\n", err)
	}
	t.SetInputName("/home/geertjohan/input.jpg")
	fmt.Println(t.GetText())

	t.DumpVariables()
}
