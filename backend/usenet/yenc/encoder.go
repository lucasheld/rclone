package yenc

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

const (
	null  byte = 0x00
	tab   byte = 0x09
	lf    byte = 0x0A
	cr    byte = 0x0D
	space byte = 0x20
	dot   byte = 0x2E
	equal byte = 0x3D
)

// http://www.yenc.org/yenc-draft.1.3.txt
// http://www.yenc.org/ydecode-c.txt

type Encoder struct {
	Writer *bytes.Buffer
	Line   int
	Size   int
	Name   string
	Crc32  string
}

func NewEncoder(writer *bytes.Buffer) *Encoder {
	return &Encoder{
		Writer: writer,
		Line:   128,
	}
}

func (e *Encoder) writeHeader() error {
	_, err := fmt.Fprintf(e.Writer, "=ybegin line=%d size=%d name=%s\r\n", e.Line, e.Size, e.Name)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) writeTrailer() error {
	_, err := fmt.Fprintf(e.Writer, "=yend size=%d crc32=%s\r\n", e.Size, e.Crc32)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) Encode(data []byte) error {
	e.Size = len(data)
	e.Crc32 = fmt.Sprintf("%x", crc32.ChecksumIEEE(data))

	err := e.writeHeader()
	if err != nil {
		return err
	}

	currentLineLength := 0
	for i := range data {
		char := data[i]
		char = char + 42

		firstColumn := currentLineLength == 0
		lastColumn := currentLineLength == e.Line-1

		// TODO: remove tab and dot from from default escape characters and adjust test files
		escapeChar := false
		if char == null || char == lf || char == cr || char == equal || char == tab || char == dot {
			escapeChar = true
		} else if firstColumn && (char == tab || char == space || char == dot) {
			escapeChar = true
		} else if lastColumn && (char == tab || char == space) {
			escapeChar = true
		}

		if escapeChar {
			e.Writer.WriteByte(equal)
			currentLineLength++
			char = char + 64
		}
		e.Writer.WriteByte(char)
		currentLineLength++

		if currentLineLength >= e.Line {
			e.Writer.Write([]byte("\r\n"))
			currentLineLength = 0
		}
	}

	if currentLineLength > 0 {
		e.Writer.Write([]byte("\r\n"))
	}
	return nil
}

func (e *Encoder) Close() error {
	err := e.writeTrailer()
	if err != nil {
		return err
	}
	return nil
}
