//  xbin.go
//  ansigo
//
// Copyright (C) 2017 ActiveState Software Inc.
// Written by Pete Garcin (@rawktron)
//
// 	Based on ansilove/C
//  Copyright (C) 2011-2017 Stefan Vogt, Brian Cassidy, and Frederic Cambus.
//  All rights reserved.
//
//  This source code is licensed under the BSD 3-Clause License.
//  See the file LICENSE for details.
//

package ansigo

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
)

// Xbin processes inputFileBuffer and outputs image data
func xbin(inputFileBuffer []byte, inputFileSize int64) image.Image {
	var f font

	if string(inputFileBuffer[0:4]) == "XBIN\x1a" {
		fmt.Print("\nNot an XBin.\n\n")
		os.Exit(4)
	}

	var xbinWidth, xbinHeight, xbinFontSize, xbinFlags int

	xbinWidth = (int(inputFileBuffer[6]) << 8) + int(inputFileBuffer[5])
	xbinHeight = (int(inputFileBuffer[8]) << 8) + int(inputFileBuffer[7])
	xbinFontSize = int(inputFileBuffer[9])
	xbinFlags = int(inputFileBuffer[10])

	var imXBIN draw.Image

	imXBIN = image.NewRGBA(image.Rect(0, 0, 8*int(xbinWidth), int(xbinFontSize)*int(xbinHeight)))

	if imXBIN == nil {
		fmt.Printf("\nError, can't allocate buffer image memory.\n\n")
		os.Exit(6)
	}

	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(imXBIN, imXBIN.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	var colors [16]color.RGBA
	offset := 11

	// palette
	if (xbinFlags & 1) == 1 {
		var index int

		for loop := 0; loop < 16; loop++ {
			index = (loop * 3) + offset

			colors[loop] = color.RGBA{(inputFileBuffer[index]<<2 | inputFileBuffer[index]>>4),
				(inputFileBuffer[index+1]<<2 | inputFileBuffer[index+1]>>4),
				(inputFileBuffer[index+2]<<2 | inputFileBuffer[index+2]>>4), 255}
		}

		offset += 48
	} else {
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
	}

	// font
	if (xbinFlags & 2) == 2 {
		var numchars int

		if (xbinFlags & 0x10) != 0 {
			numchars = 512
		} else {
			numchars = 256
		}

		f.data = inputFileBuffer[offset : offset+(int(xbinFontSize)*numchars)]
		f.sizeY = int(xbinFontSize)
		f.sizeX = 8
		f.isAmigaFont = false

		offset += (int(xbinFontSize) * numchars)
	} else {
		// using default 80x25 font
		alSelectFont(&f, "80x25")
	}

	var positionX, positionY int = 0, 0
	var character, attribute, colorForeground, colorBackground int

	// read compressed xbin
	if (xbinFlags & 4) == 4 {
		for offset < int(inputFileSize) && positionY != int(xbinHeight) {
			ctype := inputFileBuffer[offset] & 0xC0
			counter := (inputFileBuffer[offset] & 0x3F) + 1

			character = -1
			attribute = -1

			offset++
			for i := counter; i > 0; i-- {
				// none
				if ctype == 0 {
					character = int(inputFileBuffer[offset])
					attribute = int(inputFileBuffer[offset+1])
					offset += 2
				} else if ctype == 0x40 {
					// char
					if character == -1 {
						character = int(inputFileBuffer[offset])
						offset++
					}
					attribute = int(inputFileBuffer[offset])
					offset++
				} else if ctype == 0x80 {
					// attr
					if attribute == -1 {
						attribute = int(inputFileBuffer[offset])
						offset++
					}
					character = int(inputFileBuffer[offset])
					offset++
				} else {
					// both
					if character == -1 {
						character = int(inputFileBuffer[offset])
						offset++
					}
					if attribute == -1 {
						attribute = int(inputFileBuffer[offset])
						offset++
					}
				}

				colorBackground = (attribute & 240) >> 4
				colorForeground = attribute & 15

				alDrawChar(imXBIN, f.data, 8, 16, positionX, positionY, colors[colorBackground], colors[colorForeground], byte(character))

				positionX++

				if positionX == int(xbinWidth) {
					positionX = 0
					positionY++
				}
			}
		}
	} else {
		// read uncompressed xbin
		for offset < int(inputFileSize) && positionY != int(xbinHeight) {
			if positionX == int(xbinWidth) {
				positionX = 0
				positionY++
			}

			character = int(inputFileBuffer[offset])
			attribute = int(inputFileBuffer[offset+1])

			colorBackground = (attribute & 240) >> 4
			colorForeground = attribute & 15

			alDrawChar(imXBIN, f.data, 8, int(xbinFontSize), positionX, positionY, colors[colorBackground], colors[colorForeground], byte(character))

			positionX++
			offset += 2
		}
	}

	return imXBIN
}
