package yenc

import "bytes"

// http://www.yenc.org/yenc-draft.1.3.txt
// http://www.yenc.org/develop.htm

const (
	null  byte = 0x00
	tab   byte = 0x09
	lf    byte = 0x0A
	cr    byte = 0x0D
	space byte = 0x20
	dot   byte = 0x2E
	equal byte = 0x3D

	offset1 byte = 42
	offset2 byte = 64
)

type Yenc struct {
	Writer *bytes.Buffer
	Line   int
	Size   int
	Name   string
	Crc    string

	// additional attributes for multi-part binaries
	PPart  int
	PTotal int
	PBegin int
	PEnd   int
	PSize  int
	PCrc   string
}

func (y *Yenc) IsMultipart() bool {
	return y.PPart > 0
}
