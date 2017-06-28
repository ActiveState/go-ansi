//  pngw.go
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
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

// WritePng takes an image and filename and encodes to png
func WritePng(fileName string, img image.Image, scaleFactor float32) {

	scaledImg := img

	if scaleFactor != 1.0 {
		scaledHeight := float32(img.Bounds().Max.Y) * scaleFactor
		scaledWidth := float32(img.Bounds().Max.X) * scaleFactor
		scaledImg = resize.Resize(uint(scaledWidth), uint(scaledHeight), img, resize.NearestNeighbor)
	}

	// create output file
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, scaledImg); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
