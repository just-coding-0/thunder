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

func lastChar(str string) byte {
	if len(str) > 0 {
		return str[len(str)-1]
	}
	return ' '
}
