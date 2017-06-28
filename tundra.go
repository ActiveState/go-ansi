//  tundra.go
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

func tundra(inputFileBuffer []byte, inputFileSize int64, columns int, fontName string, bits int) image.Image {
	// some type declarations
	var f font

	// font selection
	alSelectFont(&f, fontName)

	// libgd image pointers
	var imTundra draw.Image

	// extract tundra header
	tundraVersion := inputFileBuffer[0]

	// need to add check for "TUNDRA24" string in the header
	if tundraVersion != 24 {
		fmt.Printf("\nInput file is not a TUNDRA file.\n\n")
		os.Exit(4)
	}

	// read tundra file a first time to find the image size
	var character, loop, positionX, positionY int

	var colorBackground, colorForeground color.RGBA

	loop = 9

	for loop < int(inputFileSize) {
		if positionX == 80 {
			positionX = 0
			positionY++
		}

		character = int(inputFileBuffer[loop])

		if character == 1 {
			positionY = (int(inputFileBuffer[loop+1]) << 24) + (int(inputFileBuffer[loop+2]) << 16) + (int(inputFileBuffer[loop+3]) << 8) + int(inputFileBuffer[loop+4])

			positionX = (int(inputFileBuffer[loop+5]) << 24) + (int(inputFileBuffer[loop+6]) << 16) + (int(inputFileBuffer[loop+7]) << 8) + int(inputFileBuffer[loop+8])

			loop += 8
		}

		if character == 2 {
			character = int(inputFileBuffer[loop+1])
			loop += 5
		}

		if character == 4 {
			character = int(inputFileBuffer[loop+1])
			loop += 5
		}

		if character == 6 {
			character = int(inputFileBuffer[loop+1])
			loop += 9
		}

		if character != 1 && character != 2 && character != 4 && character != 6 {
			positionX++
		}

		loop++
	}
	positionY++

	imTundra = image.NewRGBA(image.Rect(0, 0, columns*bits, positionY*f.sizeY))

	if imTundra == nil {
		fmt.Printf("\nError, can't allocate buffer image memory.\n\n")
		os.Exit(6)
	}

	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(imTundra, imTundra.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	// process tundra
	positionX = 0
	positionY = 0

	loop = 9

	for loop < int(inputFileSize) {
		if positionX == columns {
			positionX = 0
			positionY++
		}

		character = int(inputFileBuffer[loop])

		if character == 1 {
			positionY = (int(inputFileBuffer[loop+1]) << 24) + (int(inputFileBuffer[loop+2]) << 16) + (int(inputFileBuffer[loop+3]) << 8) + int(inputFileBuffer[loop+4])

			positionX = (int(inputFileBuffer[loop+5]) << 24) + (int(inputFileBuffer[loop+6]) << 16) + (int(inputFileBuffer[loop+7]) << 8) + int(inputFileBuffer[loop+8])

			loop += 8
		}

		if character == 2 {
			colorForeground = color.RGBA{inputFileBuffer[loop+3], inputFileBuffer[loop+4], inputFileBuffer[loop+5], 255}

			character = int(inputFileBuffer[loop+1])

			loop += 5
		}

		if character == 4 {
			colorBackground = color.RGBA{inputFileBuffer[loop+3], inputFileBuffer[loop+4], inputFileBuffer[loop+5], 255}

			character = int(inputFileBuffer[loop+1])

			loop += 5
		}

		if character == 6 {
			colorForeground = color.RGBA{inputFileBuffer[loop+3], inputFileBuffer[loop+4], inputFileBuffer[loop+5], 255}
			colorBackground = color.RGBA{inputFileBuffer[loop+7], inputFileBuffer[loop+8], inputFileBuffer[loop+9], 255}

			character = int(inputFileBuffer[loop+1])

			loop += 9
		}

		if character != 1 && character != 2 && character != 4 && character != 6 {
			alDrawChar(imTundra, f.data, bits, f.sizeY, positionX, positionY, colorBackground, colorForeground, byte(character))

			positionX++
		}

		loop++
	}

	return imTundra
}
