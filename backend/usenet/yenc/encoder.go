package yenc

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

// characters in the binary input data that should be escaped
// TODO: Careful writers of encoders will encode TAB (09h) SPACES (20h) if they would appear in the first or last column of a line. Implementors who write directly to a TCP stream will care about the doubling of dots in the first column - or also encode a DOT in the first column.
const (
	null   byte = 0x00
	tab    byte = 0x09
	lf     byte = 0x0A
	cr     byte = 0x0D
	space  byte = 0x20
	dot    byte = 0x2E
	escape byte = 0x3D
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
		// read input character
		char := data[i]

		// calculate output character
		char = char + 42
		// escape output character if necessary
		switch char {
		case escape, null, lf, cr, tab, dot:
			e.Writer.WriteByte(escape)
			currentLineLength++
			char = char + 64
		}

		//if currentLineLength == 0 && char == dot {
		//	e.Writer.WriteByte(escape)
		//	currentLineLength++
		//	char = byte(math.Mod(float64(char + 64), 256))
		//}
		//if (currentLineLength == 0 || currentLineLength >= e.Line) && char == space {
		//	e.Writer.WriteByte(escape)
		//	currentLineLength++
		//	char = byte(math.Mod(float64(char + 64), 256))
		//}

		// write output character to output stream
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
