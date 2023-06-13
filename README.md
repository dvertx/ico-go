# ICO Decoding and Encoding in Go

ico-go is a Go package to encode and decode Windows ICO file. Encoding creates ICO files with a PNG embedded image, as described in [https://en.wikipedia.org/wiki/ICO_(file_format)](https://en.wikipedia.org/wiki/ICO_(file_format)). This means that the resulting ICO will have transparencies 
consistent with its source image. Encoding implements Catmull-Rom scaling filter in goroutines for speedy processing.

The code is easy to read and modify to suit your needs. Currently only includes the encoding part. Will update with the decoding part later.

## Installation

Install with `go get -u github.com/dvertx/ico-go` or by manually cloning this repository into `$GOPATH/src/github.com/dvertx/`

## Usage

To encode, pass 3 parameters consisting of io.Writer, image.Image, and size. Size should be a single integer value of either 16, 32, 48, 64, or 256. In response ico-go will create an ICO file with embedded PNG image sizes of 16x16, 32x32, 48x48, 64x64, or 256x256 for you. Input image can be of almost any size and dimensions, and doesn't need to be square. The input image will be scaled and adjusted into a square shape accordingly.

## Example

```go
package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/dvertx/ico-go"
)

func main() {
	f, err := os.Open("sample.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	img, err := png.Decode(f)

	outFile, err := os.Create("test.ico")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer outFile.Close()

	err = ico.Encode(outFile, img, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
```
