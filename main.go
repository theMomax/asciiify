package main

import (
	"bytes"
	"image/color"
	"image/gif"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rthornton128/goncurses"
	"github.com/theMomax/asciiify/asciiif"
	"github.com/tomnomnom/xtermcolor"
)

//go:generate go-bindata "Happy B-Day 716.gif"

func main() {
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

	gifData, err := Asset("Happy B-Day 716.gif")
	if err != nil {
		panic(err)
	}

	gif, err := gif.DecodeAll(bytes.NewReader(gifData))
	if err != nil {
		panic(err)
	}

	// don't make it repeating
	gif.LoopCount = -1

	vid := asciiif.DecodeGIFStreamed(gif)

	goncurses.StartColor()

	for i := int16(0); i < 256; i++ {
		goncurses.InitPair(i, i, 0)
	}

	for img := range vid {
		start := time.Now()
		for y, row := range img.Image {
			for x, pix := range row {
				stdscr.ColorOn(int16(xtermcolor.FromColor(color.RGBA{pix.R, pix.G, pix.B, pix.A})))
				stdscr.MoveAddChar(y, x, goncurses.Char(pix.Char))
			}
		}
		stdscr.Refresh()
		took := time.Since(start)
		delay := time.Duration(img.Delay) * 10 * time.Millisecond
		if took < delay {
			time.Sleep(delay - took)
		}
	}
}
