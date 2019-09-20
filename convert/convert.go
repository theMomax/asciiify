package convert

import (
	"image"
	"image/gif"
	"os"

	"github.com/qeesung/image2ascii/convert"
)

// DefaultOptions for the underlying image2ascii converter
var DefaultOptions = convert.DefaultOptions

// ASCIIIF is the ascii-equivalent to a GIF: the ASCII Interchange Format
type ASCIIIF struct {
	Image []string // The successive images.
	Delay []int    // The successive delay times, one per frame, in 100ths of a second.
	// LoopCount controls the number of times an animation will be
	// restarted during display.
	// A LoopCount of 0 means to loop forever.
	// A LoopCount of -1 means to show each frame only once.
	// Otherwise, the animation is looped LoopCount+1 times.
	LoopCount int
}

// GIFConverter provides methods to convert GIF to ASCIIIF
type GIFConverter struct {
	convert.ImageConverter
}

// NewGIFConverter creates a GIFConverter from the default ImageConverter
func NewGIFConverter() *GIFConverter {
	return &GIFConverter{
		*convert.NewImageConverter(),
	}
}

// FromFile is a convenience function, that reads the GIF from filepath, before
// passing it to FromGIF.
func (g *GIFConverter) FromFile(filepath string, options *convert.Options) (*ASCIIIF, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gif, err := gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}

	return g.FromGIF(gif, options), nil
}

// FromGIF converts the given gif to an ASCIIIF structure using the
// ImageConverter stored in g. The options are passed down as given.
func (g *GIFConverter) FromGIF(gif *gif.GIF, options *convert.Options) *ASCIIIF {
	if len(gif.Image) == 0 {
		return nil
	}

	a := ASCIIIF{
		Delay:     gif.Delay,
		LoopCount: gif.LoopCount,
	}

	images := make([]image.Image, len(gif.Image))

	// fallback for empty Config required by Go 1.5
	if gif.Config.Width == 0 {
		gif.Config.Width = gif.Image[0].Bounds().Dx()
	}
	if gif.Config.Height == 0 {
		gif.Config.Height = gif.Image[0].Bounds().Dy()
	}

	r := image.Rect(0, 0, gif.Config.Width, gif.Config.Height)

	for i, p := range gif.Image {
		images[i] = p.SubImage(r)
	}

	asciiImages := make([]string, len(images))

	for i, im := range images {
		asciiImages[i] = g.Image2ASCIIString(im, options)
	}

	a.Image = asciiImages

	return &a
}
