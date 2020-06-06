// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package thunder

import (
	"bytes"
	"github.com/just-coding-0/thunder/internal/bytesconv"
)

var (
	strStar  = []byte("*")
	strColon = []byte(":")
)

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func countParams(path string) uint16 {
	var n uint16
	s := bytesconv.StringToBytes(path)
	n += uint16(bytes.Count(s, strColon))
	n += uint16(bytes.Count(s, strStar))
	return n
}
