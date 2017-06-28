//  artworx.go
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

// Artworx processes inputFileBuffer and generates an image
func artworx(inputFileBuffer []byte, inputFileSize int64) image.Image {
	// some type declarations
	var f font

	// libgd image pointers
	var imADF draw.Image

	imADF = image.NewRGBA(image.Rect(0, 0, 640, (((int(inputFileSize)-192-4096-1)/2)/80)*16))

	if imADF == nil {
		fmt.Printf("\nError, can't allocate buffer image memory.\n\n")
		os.Exit(6)
	}

	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(imADF, imADF.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	// ADF color palette array
	adfColors := [16]int{0, 1, 2, 3, 4, 5, 20, 7, 56, 57, 58, 59, 60, 61, 62, 63}
	var colors [16]color.RGBA
	f.data = inputFileBuffer[193 : 193+4096]

	var loop int
	var index int
	// process ADF palette
	for loop = 0; loop < 16; loop++ {
		index = (adfColors[loop] * 3) + 1

		colors[loop] = color.RGBA{(inputFileBuffer[index]<<2 | inputFileBuffer[index]>>4),
			(inputFileBuffer[index+1]<<2 | inputFileBuffer[index+1]>>4),
			(inputFileBuffer[index+2]<<2 | inputFileBuffer[index+2]>>4), 255}
	}

	// process ADF
	var positionX, positionY int = 0, 0
	var character, attribute, colorForeground, colorBackground int
	loop = 192 + 4096 + 1

	for loop < int(inputFileSize) {
		if positionX == 80 {
			positionX = 0
			positionY++
		}

		character = int(inputFileBuffer[loop])
		attribute = int(inputFileBuffer[loop+1])

		colorBackground = (attribute & 240) >> 4
		colorForeground = attribute & 15

		alDrawChar(imADF, f.data, 8, 16, positionX, positionY, colors[colorBackground], colors[colorForeground], byte(character))

		positionX++
		loop += 2
	}

	return imADF
}
