package main

import (
	"fmt"
	"os"

	options "github.com/qeesung/image2ascii/convert"
	"github.com/theMomax/asciiify/convert"
)

func main() {
	g := convert.NewGIFConverter()

	convertOptions := options.DefaultOptions

	asciiif, err := g.FromFile(os.Args[1], &convertOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, i := range asciiif.Image {
		fmt.Println("=========================================================")
		fmt.Println(i)
	}
}
