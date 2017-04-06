//  sauce.go
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
	"encoding/binary"
	"io"
	"os"
)

// SauceID is the SAUCE record identifier tagged onto the file data
const SauceID = "SAUCE"

// SauceInfo contains the bulk of the SAUCE record
type SauceInfo struct {
	ID       [5]byte
	Version  [2]byte
	Title    [35]byte
	Author   [20]byte
	Group    [20]byte
	Date     [8]byte
	FileSize int32
	DataType byte
	FileType byte
	Tinfo1   uint16
	Tinfo2   uint16
	Tinfo3   uint16
	Tinfo4   uint16
	Comments byte
	Flags    byte
	Filler   [22]byte
}

// Sauce - container structure for sauceInfo and variable length comments
// This feels like a bit of a hack
type Sauce struct {
	Sauceinf     SauceInfo
	CommentLines []string
}

// Internal constants
const recordSize = 128
const commentSize = 64
const commentID = "COMNT"

// TODO: Put this check function into a utils file
// Simple error checking to reduce duplicated code
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ReadFileName reads SAUCE via a filename.
func readFileName(fileName string) *Sauce {
	file, err := os.Open(fileName)
	check(err)

	record := readFile(file)
	file.Close()

	return record
}

// ReadFile - Read SAUCE via a FILE pointer.
// TODO: This doesn't trap OOM errors
// TODO: This seems redundant now, can we eliminate this?
func readFile(file *os.File) *Sauce {
	record := readRecord(file)
	return record
}

// ReadRecord parses a SAUCE record from a data stream
func readRecord(stream io.ReadSeeker) *Sauce {
	_, err := stream.Seek(0-recordSize, 2)

	if err != nil {
		return nil
	}

	var record Sauce
	var sinfo SauceInfo
	err = binary.Read(stream, binary.LittleEndian, &sinfo)
	check(err)

	if string(sinfo.ID[:]) == SauceID {
		var comments []string
		if sinfo.Comments > 0 {
			comments = readComments(stream, int(sinfo.Comments))
		}
		record = Sauce{Sauceinf: sinfo, CommentLines: comments}
	} else {
		record = Sauce{}
	}

	return &record
}

func readComments(stream io.ReadSeeker, comments int) []string {
	var commentLines []string

	_, err := stream.Seek(0-(recordSize+5+commentSize*int64(comments)), 2)
	check(err)

	ID := make([]byte, 6)
	stream.Read(ID)
	idString := string(ID[:6])

	if idString != commentID {
		return nil
	}

	// TODO Error checking
	for i := 0; i < comments; i++ {
		buf := make([]byte, commentSize+1)

		stream.Read(buf)

		commentLines = append(commentLines, string(buf[:]))
	}

	return commentLines
}
