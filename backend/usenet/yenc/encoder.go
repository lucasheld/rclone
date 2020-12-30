package yenc

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

type Encoder struct {
	Yenc
}

func (e *Encoder) CalcCrc(data []byte) string {
	return fmt.Sprintf("%x", crc32.ChecksumIEEE(data))
}

func (e *Encoder) writeBody(data []byte) error {
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
			err := e.Writer.WriteByte(equal)
			if err != nil {
				return err
			}
			currentLineLength++
			char = char + 64
		}
		err := e.Writer.WriteByte(char)
		if err != nil {
			return err
		}
		currentLineLength++

		if currentLineLength >= e.Line {
			_, err := e.Writer.Write([]byte("\r\n"))
			if err != nil {
				return err
			}
			currentLineLength = 0
		}
	}

	if currentLineLength > 0 {
		_, err := e.Writer.Write([]byte("\r\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) writeHeader() error {
	var header string
	if e.IsMultipart() {
		// TODO: add total and adjust test files
		header = fmt.Sprintf("=ybegin part=%d line=%d size=%d name=%s\r\n"+
			"=ypart begin=%d end=%d\r\n",
			e.PPart, e.Line, e.Size, e.Name,
			e.PBegin, e.PEnd)
	} else {
		header = fmt.Sprintf("=ybegin line=%d size=%d name=%s\r\n",
			e.Line, e.Size, e.Name)
	}

	_, err := fmt.Fprintf(e.Writer, header)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) writeTrailer() error {
	var trailer string
	if e.IsMultipart() {
		// TODO: add crc32 of the entire encoded binary and adjust test files?
		trailer = fmt.Sprintf("=yend size=%d part=%d pcrc32=%s\r\n",
			e.PSize, e.PPart, e.PCrc)
	} else {
		trailer = fmt.Sprintf("=yend size=%d crc32=%s\r\n",
			e.Size, e.Crc)
	}

	_, err := fmt.Fprintf(e.Writer, trailer)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) Encode(data []byte) error {
	crc := e.CalcCrc(data)
	if e.IsMultipart() {
		e.PCrc = crc
	} else {
		e.Crc = crc
	}

	err := e.writeHeader()
	if err != nil {
		return err
	}

	err = e.writeBody(data)
	if err != nil {
		return err
	}

	err = e.writeTrailer()
	if err != nil {
		return err
	}

	return nil
}

func NewEncoder(writer *bytes.Buffer) *Encoder {
	return &Encoder{Yenc{
		Writer: writer,
		Line:   128,
	}}
}
