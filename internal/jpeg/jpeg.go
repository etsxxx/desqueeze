package jpeg

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/tajtiattila/metadata/exif"
)

func Desqueeze(inputPath string, outputPath string, multiply float64, quality int) (err error) {
	// read EXIF
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	exifData, err := exif.Decode(input)
	if err != nil {
		return err
	}
	input.Seek(0, 0)

	// read image data
	img, err := jpeg.Decode(input)
	if err != nil {
		return err
	}
	input.Close()

	// desqueeze
	rct := img.Bounds()
	floatDx := float64(rct.Dx()) * multiply
	outDx := uint(floatDx)
	outDy := uint(rct.Dy())
	fmt.Printf("Desqueeze width: %dpx -> %dpx\n", rct.Dx(), outDx)
	m := resize.Resize(outDx, outDy, img, resize.Lanczos3)

	// create tmpfile
	tmp, err := ioutil.TempFile("", filepath.Base(outputPath))
	if err != nil {
		return err
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	if err = jpeg.Encode(tmp, m, &jpeg.Options{Quality: quality}); err != nil {
		return err
	}
	tmp.Seek(0, 0)

	// copy image data from tmpfile and write EXIF
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()
	if err = exif.Copy(out, tmp, exifData); err != nil {
		return err
	}

	return nil
}
