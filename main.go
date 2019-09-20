package main

import (
	"image/color"
	"image/gif"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rthornton128/goncurses"
	"github.com/theMomax/asciiify/asciiif"
	"github.com/tomnomnom/xtermcolor"
)

func main() {
	if len(os.Args) == 0 {
		log.Println("Self recursion!")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigs {

		}
	}()

	stdscr, err := goncurses.Init()
	if err != nil {
		panic(err)
	}
	defer goncurses.End()

	goncurses.Cursor(0)
	stdscr.ScrollOk(false)

	gifFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer gifFile.Close()

	gif, err := gif.DecodeAll(gifFile)
	if err != nil {
		panic(err)
	}

	vid := asciiif.DecodeGIFStreamed(gif)

	goncurses.StartColor()

	for i := int16(0); i < 256; i++ {
		goncurses.InitPair(i, i, 0)
	}

	for img := range vid {
		for y, row := range img.Image {
			for x, pix := range row {
				stdscr.ColorOn(int16(xtermcolor.FromColor(color.RGBA{pix.R, pix.G, pix.B, pix.A})))
				stdscr.MoveAddChar(y, x, goncurses.Char(pix.Char))
			}
		}
		stdscr.Refresh()
		time.Sleep(time.Duration(img.Delay) * 10 * time.Millisecond)
	}

	/*
		fmt.Println(os.Args)
		if len(os.Args) > 1 {
			syscall.Exec(os.Args[0], os.Args[1:], os.Environ())
		} else {
			syscall.Exec(os.Args[0], []string{}, os.Environ())
		}*/
}
