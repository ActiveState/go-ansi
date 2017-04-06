// Package ansigo provides library-level access to ANSi file formats and SAUCE records
package ansigo

import (
	"bytes"
	"image"

	"github.com/nfnt/resize"
)

// Parse takes a buffer of ANSi data and returns an Image.image
func Parse(inputFileBuffer []byte, inputFileSize int64, fontName string, bits int, columns int, mode string, icecolors bool, fext string, scaleFactor float32) image.Image {

	var outputImg image.Image
	adjustedSize := inputFileSize
	buf := bytes.NewReader(inputFileBuffer)
	record := readRecord(buf)

	// if we find a SAUCE record, update bool flag
	fileHasSAUCE := (record != nil && string(record.Sauceinf.ID[:]) == SauceID)

	// adjust the file size if file contains a SAUCE record
	if fileHasSAUCE {
		adjustedSize -= 129
		if record.Sauceinf.Comments > 0 {
			adjustedSize -= int64(5 + 64*record.Sauceinf.Comments)
		}
	}

	// create the output file by invoking the appropiate function
	if fext == ".pcb" {
		// params: input, output, font, bits
		outputImg = pcboard(inputFileBuffer, adjustedSize, fontName, bits, icecolors)
	} else if fext == ".bin" {
		// params: input, output, columns, font, bits, icecolors
		outputImg = binfile(inputFileBuffer, adjustedSize, columns, fontName, bits, icecolors)
	} else if fext == ".adf" {
		// params: input, output, bits
		outputImg = artworx(inputFileBuffer, adjustedSize)
	} else if fext == ".idf" {
		// params: input, output, bits
		outputImg = icedraw(inputFileBuffer, adjustedSize)
	} else if fext == ".tnd" {
		outputImg = tundra(inputFileBuffer, adjustedSize, columns, fontName, bits)
	} else if fext == ".xb" {
		// params: input, output, bits
		outputImg = xbin(inputFileBuffer, adjustedSize)
	} else {
		// params: input, output, font, bits, icecolors, fext
		outputImg = ansi(inputFileBuffer, adjustedSize, fontName, bits, mode, icecolors, fext)
	}

	if scaleFactor != 1.0 && outputImg != nil {
		scaledHeight := float32(outputImg.Bounds().Max.Y) * scaleFactor
		scaledWidth := float32(outputImg.Bounds().Max.X) * scaleFactor
		outputImg = resize.Resize(uint(scaledWidth), uint(scaledHeight), outputImg, resize.NearestNeighbor)
	}

	return outputImg
}

// GetSauce returns a sauce record for a given file if it exists
func GetSauce(fileName string) Sauce {
	return *readFileName(fileName)
}
