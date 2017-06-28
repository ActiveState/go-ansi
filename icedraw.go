//  icedraw.go
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
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
)

func icedraw(inputFileBuffer []byte, inputFileSize int64) image.Image {
	// extract relevant part of the IDF header, 16-bit little-endian unsigned short
	var byteBuf = []byte{inputFileBuffer[8], inputFileBuffer[9]}
	x2 := binary.LittleEndian.Uint16(byteBuf)

	// libgd image pointers
	var imIDF draw.Image

	var loop int
	var index int
	var colors [16]color.RGBA

	offset := inputFileSize - 48 - 4096

	fontData := inputFileBuffer[offset : offset+4096]

	// process IDF
	loop = 12
	var idfSequenceLength, idfSequenceLoop int

	// dynamically allocated memory buffer for IDF data
	var idfBuffer []byte

	var idfData, idfDataLength int16

	for loop < (int(inputFileSize) - 4096 - 48) {
		var byteBuf = []byte{inputFileBuffer[loop], inputFileBuffer[loop+1]}
		idfData = int16(binary.LittleEndian.Uint16(byteBuf))

		// RLE compressed data
		if idfData == 1 {
			var byteBuf = []byte{inputFileBuffer[loop+2], inputFileBuffer[loop+3]}
			idfDataLength = int16(binary.LittleEndian.Uint16(byteBuf))
			idfSequenceLength = int(idfDataLength & 255)

			for idfSequenceLoop = 0; idfSequenceLoop < idfSequenceLength; idfSequenceLoop++ {
				idfBuffer = append(idfBuffer, inputFileBuffer[loop+4])
				idfBuffer = append(idfBuffer, inputFileBuffer[loop+5])
			}
			loop += 4
		} else {
			// normal character
			idfBuffer = append(idfBuffer, inputFileBuffer[loop])
			idfBuffer = append(idfBuffer, inputFileBuffer[loop+1])
		}
		loop += 2
	}

	// create IDF instance
	imIDF = image.NewRGBA(image.Rect(0, 0, int((x2+1)*8), len(idfBuffer)/2/80*16))

	if imIDF == nil {
		fmt.Printf("\nError, can't allocate buffer image memory.\n\n")
		os.Exit(6)
	}

	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(imIDF, imIDF.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	// process IDF palette
	for loop = 0; loop < 16; loop++ {
		index = (loop * 3) + int(inputFileSize) - 48
		r := (inputFileBuffer[index]<<2 | inputFileBuffer[index]>>4)
		g := (inputFileBuffer[index+1]<<2 | inputFileBuffer[index+1]>>4)
		b := (inputFileBuffer[index+2]<<2 | inputFileBuffer[index+2]>>4)
		colors[loop] = color.RGBA{r, g, b, 255}
	}

	// render IDF
	var positionX, positionY int
	var character, attribute, colorForeground, colorBackground int

	for loop = 0; loop < len(idfBuffer); loop += 2 {
		if positionX == int(x2+1) {
			positionX = 0
			positionY++
		}

		character = int(idfBuffer[loop])
		attribute = int(idfBuffer[loop+1])

		colorBackground = (attribute & 240) >> 4
		colorForeground = attribute & 15

		alDrawChar(imIDF, fontData, 8, 16, positionX, positionY, colors[colorBackground], colors[colorForeground], byte(character))

		positionX++
	}

	// return IDF image
	return imIDF
}
