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

package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	goansi "github.com/ActiveState/go-ansi"
)

// Version - Package Version
const Version = "1.0.0"

// ExitSuccess - 0
const ExitSuccess = 0

// ExitFailure - 1
const ExitFailure = 1

// Exported Constant Settings Values
const (
	SubBreak  = true
	WrapCol80 = true
)

func showHelp() {
	fmt.Print("\nSUPPORTED FILE TYPES:\n" +
		"  ANS  BIN  ADF  IDF  XB  PCB  TND  ASC  NFO  DIZ\n" +
		"  Files with custom suffix default to the ANSi renderer.\n\n" +
		"PC FONTS:\n" +
		"  80x25              icelandic\n" +
		"  80x50              latin1\n" +
		"  baltic             latin2\n" +
		"  cyrillic           nordic\n" +
		"  french-canadian    portuguese\n" +
		"  greek              russian\n" +
		"  greek-869          terminus\n" +
		"  hebrew             turkish\n\n" +
		"AMIGA FONTS:\n" +
		"  amiga              topaz\n" +
		"  microknight        topaz+\n" +
		"  microknight+       topaz500\n" +
		"  mosoul             topaz500+\n" +
		"  pot-noodle\n\n" +
		"DOCUMENTATION:\n" +
		"  Detailed help is available at the go-ansi repository on GitHub.\n" +
		"  <https://github.com/ActiveState/go-ansi>\n\n")
}

func listExamples() {
	fmt.Print("\nEXAMPLES:\n")
	fmt.Print("  go-ansi file.ans (output path/name identical to input, no options)\n" +
		"  go-ansi -i file.ans (enable iCE colors)\n" +
		"  go-ansi -r file.ans (adds Retina @2x output file)\n" +
		"  go-ansi -o dir/file file.ans (custom path/name for output)\n" +
		"  go-ansi -s file.bin (just display SAUCE record, don't generate output)\n" +
		"  go-ansi -m transparent file.ans (render with transparent background)\n" +
		"  go-ansi -f amiga file.txt (custom font)\n" +
		"  go-ansi -f 80x50 -b 9 -c 320 -i file.bin (custom font, bits, columns, icecolors)\n" +
		"\n")
}

func versionInfo() {
	fmt.Print("All rights reserved.\n" +
		"\nFork me on GitHub: <https://github.com/ActiveState/go-ansi>\n" +
		"Bug reports: <https://github.com/ActiveState/go-ansi/issues>\n\n" +
		"This is free software, released under the 3-Clause BSD license.\n" +
		"<https://github.com/ActiveState/go-ansi/blob/master/LICENSE>\n\n")
}

// following the IEEE Std 1003.1 for utility conventions
func synopsis() {
	fmt.Print("\nSYNOPSIS:\n" +
		"  go-ansi [options] file\n" +
		"  go-ansi -e | -h | -v\n\n" +
		"OPTIONS:\n" +
		"  -b bits     set to 9 to render 9th column of block characters (default: 8)\n" +
		"  -c columns  adjust number of columns for BIN files (default: 160)\n" +
		"  -e          print a list of examples\n" +
		"  -f font     select font (default: 80x25)\n" +
		"  -h          show help\n" +
		"  -i          enable iCE colors\n" +
		"  -m mode     set rendering mode for ANS files:\n" +
		"                ced            black on gray, with 78 columns\n" +
		"                transparent    render with transparent background\n" +
		"                workbench      use Amiga Workbench palette\n" +
		"  -o file     specify output filename/path\n" +
		"  -r          creates additional Retina @2x output file\n" +
		"  -s          show SAUCE record without generating output\n" +
		"  -v          show version information\n" +
		"\n")
}

// Simple error checking to reduce duplicated code
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Printf("go-ansi %s - ANSi / ASCII art to PNG converter\n"+
		"Copyright (C) 2017 ActiveState Software Inc. Written by Pete Garcin.\n", Version)

	// SAUCE record related bool types
	justDisplaySAUCE := false
	fileHasSAUCE := false

	// retina output bool type
	createRetinaRep := false

	// iCE colors bool type
	icecolors := false

	// analyze options and do what has to be done
	fileIsBinary := false
	fileIsANSi := false
	fileIsPCBoard := false
	fileIsTundra := false

	var mode string
	var fontName string

	var input, output string
	var retinaout string

	var outputFile string
	// default to 8 if bits option is not specified
	bits := 8
	// default to 160 if columns option is not specified
	columns := 160

	// Define command line flags for parsing
	flag.IntVar(&bits, "b", 8, "-b bits")
	flag.IntVar(&columns, "c", 160, "-c columns")
	var exFl = flag.Bool("e", false, "-e show examples")
	flag.StringVar(&fontName, "f", "80x25", "-f font")
	var helpFl = flag.Bool("h", false, "-h show help")
	flag.BoolVar(&icecolors, "i", false, "-i enable iCE colors")
	flag.StringVar(&mode, "m", "", "-m mode")
	flag.StringVar(&output, "o", "", "-o file")
	flag.BoolVar(&createRetinaRep, "r", false, "-r")
	flag.BoolVar(&justDisplaySAUCE, "s", false, "-s")
	var verFl = flag.Bool("v", false, "-v")

	// Parse command line args
	flag.Parse()

	// Error checking on values
	if !(bits == 8 || bits == 9) {
		fmt.Print("\nInvalid value for bits.\n\n")
		os.Exit(ExitFailure)
	}

	if !(columns >= 1 && columns <= 8192) {
		fmt.Print("\nInvalid value for columns.\n\n")
		os.Exit(ExitFailure)
	}

	if *exFl {
		listExamples()
		os.Exit(ExitSuccess)
	}

	if *helpFl {
		showHelp()
		os.Exit(ExitSuccess)
	}

	if *verFl {
		versionInfo()
		os.Exit(ExitSuccess)
	}

	if len(flag.Args()) == 1 {
		input = flag.Arg(0)
	} else {
		synopsis()
		os.Exit(ExitSuccess)
	}

	// let's check the file for a valid SAUCE record
	record := goansi.GetSauce(input)

	// if we find a SAUCE record, update bool flag
	if string(record.Sauceinf.ID[:]) == goansi.SauceID {
		fileHasSAUCE = true
	}

	if !justDisplaySAUCE {
		// create output file name if output is not specified
		var outputName string

		if output == "" {
			outputName = input
		} else {
			outputName = output
		}

		// appending ".png" extension to output file name
		outputFile = outputName + ".png"

		if createRetinaRep {
			retinaout = outputName + "@2x.png"
		}

		// display name of input and output files
		fmt.Printf("\nInput File: %s\n", input)
		fmt.Printf("Output File: %s\n", outputFile)

		if createRetinaRep {
			fmt.Printf("Retina Output File: %s\n", retinaout)
		}

		// get file extension
		fext := strings.ToLower(filepath.Ext(input))

		// Open File
		f, err := os.Open(input)
		check(err)
		// Get file size
		fi, err := f.Stat()
		check(err)
		inputFileSize := fi.Size()
		// Read File
		inputFileBuffer := make([]byte, inputFileSize)
		_, err = f.Read(inputFileBuffer)
		check(err)
		// close input file, we don't need it anymore
		f.Close()

		var outputImg image.Image
		// create the output file by invoking the appropiate function
		if fext == ".pcb" {
			fileIsPCBoard = true
		} else if fext == ".bin" {
			fileIsBinary = true
		} else if fext == ".tnd" {
			fileIsTundra = true
		} else {
			fileIsANSi = true
		}

		// CLI does image resizing inside the pngw pkg to avoid parsing the file twice
		outputImg = goansi.Parse(inputFileBuffer, inputFileSize, fontName, bits, columns, mode, icecolors, fext, 1.0)

		if outputImg != nil {
			goansi.WritePng(outputFile, outputImg, 1.0)
			if createRetinaRep {
				goansi.WritePng(retinaout, outputImg, 2.0)
			}
		}

		// gather information and report to the command line
		if fileIsANSi || fileIsBinary ||
			fileIsPCBoard || fileIsTundra {
			fmt.Printf("Font: %s\n", fontName)
			fmt.Printf("Bits: %d\n", bits)
		}
		if icecolors && (fileIsANSi || fileIsBinary) {
			fmt.Printf("iCE Colors: enabled\n")
		}
		if fileIsBinary {
			fmt.Printf("Columns: %d\n", columns)
		}
	}
	// TODO SAUCE SUPPORT
	// either display SAUCE or tell us if there is no record
	if !fileHasSAUCE {
		fmt.Printf("\nFile %s does not have a SAUCE record.\n", input)
	} else {
		fmt.Printf("\nId: %s v%s\n", record.Sauceinf.ID, record.Sauceinf.Version)
		fmt.Printf("Title: %s\n", record.Sauceinf.Title)
		fmt.Printf("Author: %s\n", record.Sauceinf.Author)
		fmt.Printf("Group: %s\n", record.Sauceinf.Group)
		fmt.Printf("Date: %s\n", record.Sauceinf.Date)
		fmt.Printf("Datatype: %d\n", record.Sauceinf.DataType)
		fmt.Printf("Filetype: %d\n", record.Sauceinf.FileType)
		if record.Sauceinf.Flags != 0 {
			fmt.Printf("Flags: %d\n", record.Sauceinf.Flags)
		}
		if record.Sauceinf.Tinfo1 != 0 {
			fmt.Printf("Tinfo1: %d\n", record.Sauceinf.Tinfo1)
		}
		if record.Sauceinf.Tinfo2 != 0 {
			fmt.Printf("Tinfo2: %d\n", record.Sauceinf.Tinfo2)
		}
		if record.Sauceinf.Tinfo3 != 0 {
			fmt.Printf("Tinfo3: %d\n", record.Sauceinf.Tinfo3)
		}
		if record.Sauceinf.Tinfo4 != 0 {
			fmt.Printf("Tinfo4: %d\n", record.Sauceinf.Tinfo4)
		}

		fmt.Printf("Num comments: %d", record.Sauceinf.Comments)
		if record.Sauceinf.Comments > 0 && len(record.CommentLines) > 0 {
			fmt.Printf("Comments: ")
			for i := 0; i < int(record.Sauceinf.Comments); i++ {
				fmt.Printf("%s\n", record.CommentLines[i])
			}
		}
	}

	os.Exit(ExitSuccess)
}
