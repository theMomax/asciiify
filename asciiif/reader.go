package asciiif

import (
	"encoding/json"
	"image"
	"image/gif"
	"io"

	"github.com/qeesung/image2ascii/ascii"
	"github.com/qeesung/image2ascii/convert"
)

var c = newGIFConverter()

// ASCIIIF is the ascii-equivalent to a GIF: the ASCII Interchange Format.
type ASCIIIF struct {
	Image [][][]ascii.CharPixel // The successive images.
	Delay []int                 // The successive delay times, one per frame, in 100ths of a second.
	// LoopCount controls the number of times an animation will be
	// restarted during display.
	// A LoopCount of 0 means to loop forever.
	// A LoopCount of -1 means to show each frame only once.
	// Otherwise, the animation is looped LoopCount+1 times.
	LoopCount int
}

// ASCIIIFrame represents a single image taken form an ASCIIIF and the delay to
// the succeeding frame.
type ASCIIIFrame struct {
	Image [][]ascii.CharPixel
	Delay int
}

// DecodeAll decodes an ASCIIIF from the given reader. The reader shall read
// from a json-encoded ASCIIIF as produced by EncodeALL.
func DecodeAll(r io.Reader) (*ASCIIIF, error) {
	decoder := json.NewDecoder(r)
	a := ASCIIIF{}
	err := decoder.Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// DecodeGIF decodes an ASCIIIF from the given GIF and options. The GIF will
// be decoded using the default converter provided by the image2ascii library.
// If no options are provided, the convert package's default options will be
// used.
func DecodeGIF(gif *gif.GIF, options ...*convert.Options) *ASCIIIF {

	a := ASCIIIF{
		Delay:     gif.Delay,
		LoopCount: gif.LoopCount,
	}

	convertOptions, r, ok := adapt(gif, options...)
	if !ok {
		return nil
	}

	images := make([]image.Image, len(gif.Image))

	for i, p := range gif.Image {
		images[i] = p.SubImage(r)
	}

	asciiImages := make([][][]ascii.CharPixel, len(images))

	for i, im := range images {
		asciiImages[i] = c.Image2CharPixelMatrix(im, convertOptions)
	}

	a.Image = asciiImages

	return &a
}

// DecodeGIFAsync is an asynchronous version of DecodeGIF. It returns the
// ASCIIIF's LoopCount immediately and closes the returned channel, after all
// frames have been converted and sent.
func DecodeGIFAsync(gif *gif.GIF, options ...*convert.Options) (int, <-chan *ASCIIIFrame) {
	f := make(chan *ASCIIIFrame, 10)

	go func() {
		defer close(f)

		convertOptions, r, ok := adapt(gif, options...)
		if !ok {
			return
		}

		for i := range gif.Image {
			d := 0
			if len(gif.Delay) > i {
				d = gif.Delay[i]
			}

			frame := ASCIIIFrame{
				Image: c.Image2CharPixelMatrix(gif.Image[i].SubImage(r), convertOptions),
				Delay: d,
			}
			f <- &frame
		}
	}()

	return gif.LoopCount, f
}

// DecodeGIFStreamed streams the complete resulting ASCIIIF asynchronously. The
// returned channel is closed, when all frames have been sent for as many times
// as specified by the gif's LoopCount. The returned channel may be closed by
// the receiving side.
func DecodeGIFStreamed(gif *gif.GIF, options ...*convert.Options) (frames <-chan *ASCIIIFrame) {
	f := make(chan *ASCIIIFrame, 10)

	go func() {
		defer close(f)
		defer recover() // handle send on closed channel

		infinite := gif.LoopCount == 0
		once := gif.LoopCount == -1
		count := gif.LoopCount

		cache := make([]*ASCIIIFrame, len(gif.Image))

		_, frames := DecodeGIFAsync(gif, options...)

		for infinite || once || count >= 0 {

			for i := range cache {
				if cache[i] == nil {
					// frames being closed at this point should not be possible
					cache[i] = <-frames
				}

				f <- cache[i]
			}

			once = false
			count--
		}
	}()
	return f
}

// gifConverter provides methods to convert GIF to ASCIIIF
type gifConverter struct {
	convert.ImageConverter
}

// newGIFConverter creates a gifConverter from the default ImageConverter
func newGIFConverter() *gifConverter {
	return &gifConverter{
		*convert.NewImageConverter(),
	}
}

func adapt(gif *gif.GIF, options ...*convert.Options) (opt *convert.Options, rect image.Rectangle, ok bool) {
	convertOptions := &convert.DefaultOptions
	if len(options) > 0 {
		convertOptions = options[0]
	}

	if len(gif.Image) == 0 {
		return nil, image.Rect(0, 0, 0, 0), false
	}

	width := gif.Config.Width
	height := gif.Config.Height

	// fallback for empty Config required by Go 1.5
	if width == 0 {
		width = gif.Image[0].Bounds().Dx()
	}
	if height == 0 {
		height = gif.Image[0].Bounds().Dy()
	}

	return convertOptions, image.Rect(0, 0, width, height), true
}
