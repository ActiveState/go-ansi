//  pcb.go
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

// Character structure
type pcbChar struct {
	positionX       int
	positionY       int
	colorBackground int
	colorForeground int
	currentChar     int
}

func pcboard(inputFileBuffer []byte, inputFileSize int64, fontName string, bits int, icecolors bool) image.Image {
	// some type declarations
	var f font
	columns := 80
	var loop int

	// font selection
	alSelectFont(&f, fontName)

	// libgd image pointers
	var imPCB draw.Image

	// process PCBoard
	var char, currentChar, nextChar int
	var colorBackground, colorForeground int = 0, 7
	var posX, posY, posXMax, posYMax int

	// PCB buffer structure array definition
	var pcbBuffer []pcbChar

	// reset loop
	loop = 0
	structIndex := 0

	for loop < int(inputFileSize) {
		currentChar = int(inputFileBuffer[loop])
		nextChar = int(inputFileBuffer[loop+1])

		if posX == 80 {
			posY++
			posX = 0
		}

		// CR + LF
		if currentChar == 13 && nextChar == 10 {
			posY++
			posX = 0
			loop++
		}

		// LF
		if currentChar == 10 {
			posY++
			posX = 0
		}

		// Tab
		if currentChar == 9 {
			posX += 8
		}

		// Sub
		if currentChar == 26 {
			break
		}

		// PCB sequence
		if currentChar == 64 && nextChar == 88 {
			colorBackground = int(inputFileBuffer[loop+2])
			if colorBackground >= 65 {
				colorBackground -= 55
			} else {
				colorBackground -= 48
			}
			if !icecolors && colorBackground > 7 {
				colorBackground -= 8
			}
			colorForeground = int(inputFileBuffer[loop+3])
			if colorForeground >= 65 {
				colorForeground -= 55
			} else {
				colorForeground -= 48
			}
			loop += 3
		} else if currentChar == 64 && nextChar == 67 &&
			inputFileBuffer[loop+2] == 'L' && inputFileBuffer[loop+3] == 'S' {
			// erase display
			posX = 0
			posY = 0

			posXMax = 0
			posYMax = 0

			loop += 4
		} else if currentChar == 64 && nextChar == 80 && inputFileBuffer[loop+2] == 'O' && inputFileBuffer[loop+3] == 'S' && inputFileBuffer[loop+4] == ':' {
			// cursor position
			if inputFileBuffer[loop+6] == '@' {
				posX = int(((inputFileBuffer[loop+5]) - 48)) - 1
				loop += 5
			} else {
				posX = int((10*((inputFileBuffer[loop+5])-48) + (inputFileBuffer[loop+6]) - 48)) - 1
				loop += 6
			}
		} else if currentChar != 10 && currentChar != 13 && currentChar != 9 {
			// record number of columns and lines used
			if posX > posXMax {
				posXMax = posX
			}

			if posY > posYMax {
				posYMax = posY
			}

			var newChar pcbChar

			// write current character in pcbChar structure
			newChar.positionX = posX
			newChar.positionY = posY
			newChar.colorBackground = colorBackground
			newChar.colorForeground = colorForeground
			newChar.currentChar = currentChar

			pcbBuffer = append(pcbBuffer, newChar)

			structIndex++
			posX++
		}
		loop++
	}
	posXMax++
	posYMax++

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

	imPCB = image.NewRGBA(image.Rect(0, 0, columns*bits, posYMax*f.sizeY))

	if imPCB == nil {
		fmt.Printf("\nError, can't allocate buffer image memory.\n\n")
		os.Exit(6)
	}

	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(imPCB, imPCB.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	// render PCB
	for loop = 0; loop < structIndex; loop++ {
		// grab our chars out of the structure
		posX = pcbBuffer[loop].positionX
		posY = pcbBuffer[loop].positionY
		colorBackground = pcbBuffer[loop].colorBackground
		colorForeground = pcbBuffer[loop].colorForeground
		char = pcbBuffer[loop].currentChar

		alDrawChar(imPCB, f.data, bits, f.sizeY, posX, posY, colors[colorBackground], colors[colorForeground], byte(char))
	}

	return imPCB
}
