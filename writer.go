package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/png"
	"io"
	"sync"

	"golang.org/x/image/draw"
)

type icondir struct {
	Reserved  uint16
	ImgType   uint16
	NumImages uint16
}

type icondirentry struct {
	Width    uint8
	Height   uint8
	Pallette uint8
	Reserved uint8
	Planes   uint16
	Bpp      uint16
	DataSize uint32
	Offset   uint32
}

func Encode(iw io.Writer, im image.Image, size int) error {
	if size != 16 && size != 32 && size != 48 && size != 64 && size != 256 {
		err := errors.New("Unsupported icon size request")
		return err
	}

	img, err := scaleImage(im, size)

	pngbuffer := new(bytes.Buffer)
	pngwriter := bufio.NewWriter(pngbuffer)

	err = png.Encode(pngwriter, img)
	if err != nil {
		return err
	}

	err = pngwriter.Flush()
	if err != nil {
		return err
	}

	header := icondir{
		Reserved:  0,
		ImgType:   1,
		NumImages: 1,
	}

	entry := icondirentry{
		Pallette: 0,
		Reserved: 0,
		Planes:   1,
		Bpp:      32,
		Offset:   22,
	}

	entry.DataSize = uint32(len(pngbuffer.Bytes()))

	bounds := img.Bounds()
	entry.Width = uint8(bounds.Dx())
	entry.Height = uint8(bounds.Dy())

	bytebuff := new(bytes.Buffer)
	err = binary.Write(bytebuff, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	err = binary.Write(bytebuff, binary.LittleEndian, entry)
	if err != nil {
		return err
	}

	_, err = iw.Write(bytebuff.Bytes())
	if err != nil {
		return err
	}

	_, err = iw.Write(pngbuffer.Bytes())
	if err != nil {
		return err
	}

	return err
}

func scaleImage(img image.Image, size int) (image.Image, error) {
	out := image.NewRGBA(image.Rect(0, 0, size, size))

	var wg sync.WaitGroup

	srcRegions := splitImage(img.Bounds(), 4)
	outRegions := splitImage(out.Bounds(), 4)

	for i, r := range srcRegions { // Launch a goroutine for each region
		ro := outRegions[i]
		wg.Add(1)
		go func(r, ro image.Rectangle) {
			defer wg.Done()
			draw.CatmullRom.Scale(out, ro, img, r, draw.Src, nil)
		}(r, ro)
	}

	wg.Wait()
	return out, nil
}

func splitImage(r image.Rectangle, n int) []image.Rectangle {
	var regions []image.Rectangle
	midX := (r.Min.X + r.Max.X) / 2
	midY := (r.Min.Y + r.Max.Y) / 2
	regions = append(regions, image.Rect(r.Min.X, r.Min.Y, midX, midY))
	regions = append(regions, image.Rect(midX, r.Min.Y, r.Max.X, midY))
	regions = append(regions, image.Rect(r.Min.X, midY, midX, r.Max.Y))
	regions = append(regions, image.Rect(midX, midY, r.Max.X, r.Max.Y))
	return regions
}
