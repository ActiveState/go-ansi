//  bin.go
//  go-ansi
//
// Copyright (C) 2017 ActiveState Software Inc.
// Written by Pete Garcin (@rawktron)
//
// 	Based on ansilove/C
//  Copyright (C) 2011-2017 Stefan Vogt, Brian Cassidy, and Frederic Cambus.
//  All rights reserved.
//  ansilove/C is licensed under the BSD-2 License.
//
//  go-ansi is licensed under the BSD 3-Clause License.
//  See the file LICENSE for details.
//

package goansi

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
)

// binary processes inputFileBuffer and generates an image
func binfile(inputFileBuffer []byte, inputFileSize int64, columns int, fontName string, bits int, icecolors bool) image.Image {
	// some type declarations
	var f font

	// font selection
	alSelectFont(&f, fontName)

	// libgd image pointers
	var imBinary draw.Image

	imBinary = image.NewRGBA(image.Rect(0, 0, columns*bits, (int(inputFileSize)/2)/columns*f.sizeY))

	if imBinary == nil {
		fmt.Printf("\nError, can't allocate buffer image memory.\n\n")
		os.Exit(6)
	}

	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(imBinary, imBinary.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	// allocate color palette
	var colors [16]color.RGBA

	colors[0] = color.RGBA{0, 0, 0, 255}
	colors[1] = color.RGBA{0, 0, 170, 255}
	colors[2] = color.RGBA{0, 170, 0, 255}
	colors[3] = color.RGBA{0, 170, 170, 255}
	colors[4] = color.RGBA{170, 0, 0, 255}
	colors[5] = color.RGBA{170, 0, 170, 255}
	colors[6] = color.RGBA{170, 85, 0, 255}
	colors[7] = color.RGBA{170, 170, 170, 255}
	colors[8] = color.RGBA{85, 85, 85, 255}
	colors[9] = color.RGBA{85, 85, 255, 255}
	colors[10] = color.RGBA{85, 255, 85, 255}
	colors[11] = color.RGBA{85, 255, 255, 255}
	colors[12] = color.RGBA{255, 85, 85, 255}
	colors[13] = color.RGBA{255, 85, 255, 255}
	colors[14] = color.RGBA{255, 255, 85, 255}
	colors[15] = color.RGBA{255, 255, 255, 255}

	// process binary
	var character, attribute, colorBackground, colorForeground int
	var loop, positionX, positionY int = 0, 0, 0

	for loop < int(inputFileSize) {
		if positionX == columns {
			positionX = 0
			positionY++
		}

		character = int(inputFileBuffer[loop])
		attribute = int(inputFileBuffer[loop+1])

		colorBackground = (attribute & 240) >> 4
		colorForeground = (attribute & 15)

		if colorBackground > 8 && !icecolors {
			colorBackground -= 8
		}

		alDrawChar(imBinary, f.data, bits, f.sizeY,
			positionX, positionY, colors[colorBackground], colors[colorForeground], byte(character))

		positionX++
		loop += 2
	}

	return imBinary
}
