package yenc

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

type Encoder struct {
	Writer *bytes.Buffer
	Line   int
	Size   int
	Name   string
}

func NewEncoder(writer *bytes.Buffer, inputFileSize int) *Encoder {
	return &Encoder{
		Writer: writer,
		Line:   128,
		Size:   inputFileSize,
	}
}

func (e *Encoder) writeHeader(part *Part) error {
	content := "=ybegin"
	if part.IsMultipart() {
		content = fmt.Sprintf("%s part=%d", content, part.Part)
	}
	content = fmt.Sprintf("%s line=%d size=%d name=%s\r\n", content, e.Line, e.Size, e.Name)
	if part.IsMultipart() {
		content = fmt.Sprintf("%s=ypart begin=%d end=%d\r\n", content, part.Begin, part.End)
	}

	_, err := e.Writer.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) writeTrailer(part *Part) error {
	content := ""
	if part.IsMultipart() {
		content = fmt.Sprintf("=yend size=%d part=%d pcrc32=%s\r\n", part.Size(), part.Part, part.Crc)
	} else {
		content = fmt.Sprintf("=yend size=%d crc32=%s\r\n", e.Size, part.Crc)
	}

	_, err := e.Writer.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) Encode(part *Part, data []byte) error {
	part.Crc = fmt.Sprintf("%x", crc32.ChecksumIEEE(data))

	err := e.writeHeader(part)
	if err != nil {
		return err
	}

	currentLineLength := 0
	for i := range data {
		char := data[i]
		char = char + 42

		firstColumn := currentLineLength == 0
		lastColumn := currentLineLength == e.Line - 1

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

	err = e.writeTrailer(part)
	if err != nil {
		return err
	}

	return nil
}
