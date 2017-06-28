//  drawchar.go
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
)

// AlDrawChar - shared method for drawing ANSI characters into an image buffer
func alDrawChar(im draw.Image, font []byte, bits int, fontSizeY int, positionX int, positionY int, colorBackground color.RGBA, colorForeground color.RGBA, character byte) {
	x := positionX * bits
	y := positionY * fontSizeY

	draw.Draw(im, image.Rect(x, y, x+bits, y+fontSizeY), &image.Uniform{colorBackground}, image.ZP, draw.Src)

	for line := 0; line < fontSizeY; line++ {
		for column := 0; column < bits; column++ {
			if (font[line+int(character)*fontSizeY] & (0x80 >> uint(column))) != 0 {
				im.Set(positionX*bits+column, positionY*fontSizeY+line, colorForeground)
				if bits == 9 && column == 7 && character > 191 && character < 224 {
					im.Set(positionX*bits+8, positionY*fontSizeY+line, colorForeground)
				}
			}
		}
	}
}
