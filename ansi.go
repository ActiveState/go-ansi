//  ansi.go
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
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"
)

// Character structure
type ansiChar struct {
	positionX       int
	positionY       int
	colorBackground int
	colorForeground int
	colorFg24       color.RGBA
	colorBg24       color.RGBA
	currentChar     byte
	bold            bool
	italics         bool
	underline       bool
}

// Ansi takes an inputFileBuffer with .ans data and returns an image buffer
func ansi(inputFileBuffer []byte, inputFileSize int64, fontName string, bits int, mode string, icecolors bool, fext string) image.Image {
	var f font

	columns := 80

	isDizFile := false
	ced := false
	transparent := false
	workbench := false

	// font selection
	alSelectFont(&f, fontName)

	// to deal with the bits flag, we declared handy bool types
	if mode == "ced" {
		ced = true
	} else if mode == "transparent" {
		transparent = true
	} else if mode == "workbench" {
		workbench = true
	}

	// check if current file has a .diz extension
	if fext == ".diz" {
		isDizFile = true
	}

	// image buffer
	var imANSi draw.Image

	// ANSi processing loops
	var loop int

	// character definitions
	var currentChar, nextChar, character byte
	var ansiSequenceChar byte

	// default color values
	colorBackground := 0
	colorForeground := 7

	var fgcolor color.RGBA
	var bgcolor color.RGBA

	// text attributes
	var bold, underline, italics, blink bool = false, false, false, false

	// positions
	var positionX, positionY, positionXMax, positionYMax int = 0, 0, 0, 0
	var savedPositionY, savedPositionX int = 0, 0

	// sequence parsing variables
	var seqArrayCount int
	var seqArray []string
	var seqGrab string

	// ANSi buffer structure array definition
	var structIndex int
	var ansiBuffer []ansiChar

	fg24 := color.RGBA{0, 0, 0, 0}
	bg24 := color.RGBA{0, 0, 0, 0}

	// ANSi interpreter
	for loop < int(inputFileSize)-1 {
		currentChar = inputFileBuffer[loop]
		nextChar = inputFileBuffer[loop+1]

		// TODO Make these properly scoped
		const wrapColumn80 bool = true
		if positionX == 80 && wrapColumn80 {
			positionY++
			positionX = 0
		}

		// CR + LF
		if currentChar == 13 && nextChar == 10 {
			positionY++
			positionX = 0
			loop++
		}

		// LF
		if currentChar == 10 {
			positionY++
			positionX = 0
		}

		// tab
		if currentChar == 9 {
			positionX += 8
		}

		// TODO Make this scoped properly
		const subBreak bool = true
		// sub
		if currentChar == 26 && subBreak {
			break
		}

		// ANSi sequence
		if currentChar == 27 && nextChar == 91 {
			for ansiSequenceLoop := 0; ansiSequenceLoop < 15; ansiSequenceLoop++ {
				ansiSequenceChar = inputFileBuffer[loop+2+ansiSequenceLoop]

				// cursor position
				if ansiSequenceChar == 'H' || ansiSequenceChar == 'f' {
					// create substring from the sequence's content

					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])

					// create sequence content array
					seqArray = strings.Split(seqGrab, ";")
					seqArrayCount = len(seqArray)

					if seqArrayCount > 1 {
						// convert grabbed sequence content to integers
						seqLine, _ := strconv.Atoi(seqArray[0])
						seqColumn, _ := strconv.Atoi(seqArray[1])

						// finally set the positions
						positionY = seqLine - 1
						positionX = seqColumn - 1
					} else {
						// no coordinates specified? we move to the home position
						positionY = 0
						positionX = 0
					}
					loop += ansiSequenceLoop + 2
					break
				}

				// cursor up
				if ansiSequenceChar == 'A' {
					// create substring from the sequence's content
					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])
					//println("A")

					// now get escape sequence's position value
					seqLine, _ := strconv.Atoi(seqGrab)

					if seqLine == 0 {
						seqLine = 1
					}

					positionY = positionY - seqLine

					loop += ansiSequenceLoop + 2
					break
				}

				// cursor down
				if ansiSequenceChar == 'B' {
					// create substring from the sequence's content
					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])
					//println("B")

					// now get escape sequence's position value
					seqLine, _ := strconv.Atoi(seqGrab)

					if seqLine == 0 {
						seqLine = 1
					}

					positionY = positionY + seqLine

					loop += ansiSequenceLoop + 2
					break
				}

				// cursor forward
				if ansiSequenceChar == 'C' {
					// create substring from the sequence's content
					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])
					//println("C")

					// now get escape sequence's position value
					seqColumn, _ := strconv.Atoi(seqGrab)

					if seqColumn == 0 {
						seqColumn = 1
					}

					positionX = positionX + seqColumn

					if positionX > 80 {
						positionX = 80
					}

					loop += ansiSequenceLoop + 2
					break
				}

				// cursor backward
				if ansiSequenceChar == 'D' {
					// create substring from the sequence's content
					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])
					//println("D")

					// now get escape sequence's content length
					seqColumn, _ := strconv.Atoi(seqGrab)

					if seqColumn == 0 {
						seqColumn = 1
					}

					positionX = positionX - seqColumn

					if positionX < 0 {
						positionX = 0
					}

					loop += ansiSequenceLoop + 2
					break
				}

				// save cursor position
				if ansiSequenceChar == 's' {
					savedPositionY = positionY
					savedPositionX = positionX
					//println("s")

					loop += ansiSequenceLoop + 2
					break
				}

				// restore cursor position
				if ansiSequenceChar == 'u' {
					positionY = savedPositionY
					positionX = savedPositionX
					//println("u")

					loop += ansiSequenceLoop + 2
					break
				}

				// erase display
				if ansiSequenceChar == 'J' {
					// create substring from the sequence's content
					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])
					//println("J")
					// convert grab to an integer
					eraseDisplayInt, _ := strconv.Atoi(seqGrab)

					if eraseDisplayInt == 2 {
						positionX = 0
						positionY = 0

						positionXMax = 0
						positionYMax = 0

						// reset ansi buffer
						ansiBuffer = nil
						structIndex = 0
					}
					loop += ansiSequenceLoop + 2
					break
				}

				// set graphics mode
				if ansiSequenceChar == 'm' {
					//println("m")
					// create substring from the sequence's content
					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])

					// create sequence content array
					seqArray = strings.Split(seqGrab, ";")
					seqArrayCount = len(seqArray)

					//					fmt.Printf("SEQARRAY: %v", seqArray)
					// a loophole in limbo
					for seqGraphicsLoop := 0; seqGraphicsLoop < seqArrayCount; seqGraphicsLoop++ {
						// convert split content value to integer
						seqValue, _ := strconv.Atoi(seqArray[seqGraphicsLoop])

						//println("SEQVALUE", seqValue)

						if seqValue == 0 {
							colorBackground = 0
							colorForeground = 7
							bold = false
							underline = false
							italics = false
							blink = false
						}

						if seqValue == 1 {
							if !workbench {
								colorForeground += 8
							}
							bold = true
						}

						if seqValue == 3 {
							italics = true
						}

						if seqValue == 4 {
							underline = true
						}

						if seqValue == 5 {
							if !workbench {
								colorBackground += 8
							}
							blink = true
						}

						if seqValue > 29 && seqValue < 38 {
							colorForeground = seqValue - 30

							if bold {
								colorForeground += 8
							}
						}

						if seqValue > 39 && seqValue < 48 {
							colorBackground = seqValue - 40

							if blink && icecolors {
								colorBackground += 8
							}
						}
					}

					loop += ansiSequenceLoop + 2
					break
				}

				// 24-bit ANSI support
				// Sets a temporary color for this sequence that overrides the foreground or background colors
				if ansiSequenceChar == 't' {
					// create substring from the sequence's content
					seqGrab = string(inputFileBuffer[loop+2 : loop+2+ansiSequenceLoop])

					// create sequence content array
					seqArray = strings.Split(seqGrab, ";")
					seqArrayCount = len(seqArray)

					if seqArrayCount == 4 {
						r, _ := strconv.Atoi(seqArray[1])
						g, _ := strconv.Atoi(seqArray[2])
						b, _ := strconv.Atoi(seqArray[3])

						if seqArray[0] == "0" {
							bg24 = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
						} else if seqArray[0] == "1" {
							fg24 = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
						}
					}

					loop += ansiSequenceLoop + 2
					break
				}

				// cursor (de)activation (Amiga ANSi)
				if ansiSequenceChar == 'p' {
					loop += ansiSequenceLoop + 2
					break
				}

				// skipping set mode and reset mode sequences
				if ansiSequenceChar == 'h' || ansiSequenceChar == 'l' {
					loop += ansiSequenceLoop + 2
					break
				}
			}
		} else if currentChar != 10 && currentChar != 13 && currentChar != 9 {
			// record number of columns and lines used
			if positionX > positionXMax {
				positionXMax = positionX
			}

			if positionY > positionYMax {
				positionYMax = positionY
			}

			// write current character in ansiChar structure
			if !f.isAmigaFont || (currentChar != 12 && currentChar != 13) {
				var newChar ansiChar

				newChar.colorBackground = colorBackground
				newChar.colorForeground = colorForeground
				newChar.colorFg24 = fg24
				newChar.colorBg24 = bg24
				newChar.currentChar = currentChar
				newChar.bold = bold
				newChar.italics = italics
				newChar.underline = underline
				newChar.positionX = positionX
				newChar.positionY = positionY

				ansiBuffer = append(ansiBuffer, newChar)

				fg24 = color.RGBA{0, 0, 0, 0}
				bg24 = color.RGBA{0, 0, 0, 0}

				structIndex++
				positionX++
			}
		}
		loop++
	}

	// allocate image buffer memory
	positionXMax++
	positionYMax++

	if ced {
		columns = 78
	}

	if isDizFile {
		columns = min(positionXMax, 80)
	}

	imANSi = image.NewRGBA(image.Rect(0, 0, columns*bits, (positionYMax)*f.sizeY))
	black := color.RGBA{0, 0, 0, 255}

	var colors [16]color.RGBA

	var cedBackground, cedForeground color.RGBA

	if ced {
		cedBackground = color.RGBA{170, 170, 170, 255}
		cedForeground = color.RGBA{0, 0, 0, 255}
		draw.Draw(imANSi, imANSi.Bounds(), &image.Uniform{cedBackground}, image.ZP, draw.Src)
	} else if workbench {

		if transparent {
			draw.Draw(imANSi, imANSi.Bounds(), image.Transparent, image.ZP, draw.Src)
		} else {
			draw.Draw(imANSi, imANSi.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)
		}

		colors[0] = color.RGBA{170, 170, 170, 255}
		colors[1] = color.RGBA{0, 0, 0, 255}
		colors[2] = color.RGBA{255, 255, 255, 255}
		colors[3] = color.RGBA{102, 136, 187, 255}
		colors[4] = color.RGBA{0, 0, 255, 255}
		colors[5] = color.RGBA{255, 0, 255, 255}
		colors[6] = color.RGBA{0, 255, 255, 255}
		colors[7] = color.RGBA{255, 255, 255, 255}
		colors[8] = color.RGBA{170, 170, 170, 255}
		colors[9] = color.RGBA{0, 0, 0, 255}
		colors[10] = color.RGBA{255, 255, 255, 255}
		colors[11] = color.RGBA{102, 136, 187, 255}
		colors[12] = color.RGBA{0, 0, 255, 255}
		colors[13] = color.RGBA{255, 0, 255, 255}
		colors[14] = color.RGBA{0, 255, 255, 255}
		colors[15] = color.RGBA{255, 255, 255, 255}
	} else {
		// Allocate standard ANSi color palette
		if transparent {
			draw.Draw(imANSi, imANSi.Bounds(), image.Transparent, image.ZP, draw.Src)
		} else {
			draw.Draw(imANSi, imANSi.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)
		}

		colors[0] = color.RGBA{0, 0, 0, 255}
		colors[1] = color.RGBA{170, 0, 0, 255}
		colors[2] = color.RGBA{0, 170, 0, 255}
		colors[3] = color.RGBA{170, 85, 0, 255}
		colors[4] = color.RGBA{0, 0, 170, 255}
		colors[5] = color.RGBA{170, 0, 170, 255}
		colors[6] = color.RGBA{0, 170, 170, 255}
		colors[7] = color.RGBA{170, 170, 170, 255}
		colors[8] = color.RGBA{85, 85, 85, 255}
		colors[9] = color.RGBA{255, 85, 85, 255}
		colors[10] = color.RGBA{85, 255, 85, 255}
		colors[11] = color.RGBA{255, 255, 85, 255}
		colors[12] = color.RGBA{85, 85, 255, 255}
		colors[13] = color.RGBA{255, 85, 255, 255}
		colors[14] = color.RGBA{85, 255, 255, 255}
		colors[15] = color.RGBA{255, 255, 255, 255}
	}

	// even more definitions, sigh
	ansiBufferItems := structIndex

	// render ANSi
	for loop := 0; loop < ansiBufferItems; loop++ {
		// grab ANSi char from our structure array
		colorBackground = ansiBuffer[loop].colorBackground
		colorForeground = ansiBuffer[loop].colorForeground
		character = ansiBuffer[loop].currentChar
		bold = ansiBuffer[loop].bold
		italics = ansiBuffer[loop].italics
		underline = ansiBuffer[loop].underline
		positionX = ansiBuffer[loop].positionX
		positionY = ansiBuffer[loop].positionY

		if ansiBuffer[loop].colorBg24.A > 0 {
			bgcolor = ansiBuffer[loop].colorBg24
		} else {
			bgcolor = colors[colorBackground]
		}
		if ansiBuffer[loop].colorFg24.A > 0 {
			fgcolor = ansiBuffer[loop].colorFg24
		} else {
			fgcolor = colors[colorForeground]
		}

		if ced {
			alDrawChar(imANSi, f.data, bits, f.sizeY,
				positionX, positionY, cedBackground, cedForeground, character)
		} else {
			alDrawChar(imANSi, f.data, bits, f.sizeY,
				positionX, positionY, bgcolor, fgcolor, character)
		}

	}

	return imANSi
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
