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

func (e *Encoder) calcCrc(data []byte) string {
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

type SinglepartEncoder struct {
	*Encoder
	Crc string
}

func (e *SinglepartEncoder) writeHeader() error {
	_, err := fmt.Fprintf(e.Writer,
		"=ybegin line=%d size=%d name=%s\r\n",
		e.Line, e.Size, e.Name)
	if err != nil {
		return err
	}
	return nil
}

func (e *SinglepartEncoder) writeTrailer() error {
	_, err := fmt.Fprintf(e.Writer,
		"=yend size=%d crc32=%s\r\n",
		e.Size, e.Crc)
	if err != nil {
		return err
	}
	return nil
}

func (e *SinglepartEncoder) Encode(data []byte) error {
	e.Size = len(data)
	e.Crc = e.calcCrc(data)

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

type MultipartEncoder struct {
	*Encoder
}

func (e *MultipartEncoder) writeHeader(part *Part) error {
	// TODO: add total and adjust test files
	_, err := fmt.Fprintf(e.Writer,
		"=ybegin part=%d line=%d size=%d name=%s\r\n"+
			"=ypart begin=%d end=%d\r\n",
		part.Part, e.Line, e.Size, e.Name,
		part.Begin, part.End)
	if err != nil {
		return err
	}
	return nil
}
func (e *MultipartEncoder) writeTrailer(part *Part) error {
	// TODO: add crc32 of the entire encoded binary and adjust test files?
	_, err := fmt.Fprintf(e.Writer,
		"=yend size=%d part=%d pcrc32=%s\r\n",
		part.Size, part.Part, part.Crc)
	if err != nil {
		return err
	}
	return nil
}

func (e *MultipartEncoder) Encode(part *Part, data []byte) error {
	part.Size = part.End - part.Begin + 1
	part.Crc = e.calcCrc(data)

	err := e.writeHeader(part)
	if err != nil {
		return err
	}

	err = e.writeBody(data)
	if err != nil {
		return err
	}

	err = e.writeTrailer(part)
	if err != nil {
		return err
	}

	return nil
}

func NewSinglepartEncoder(writer *bytes.Buffer) *SinglepartEncoder {
	return &SinglepartEncoder{
		Encoder: &Encoder{
			Writer: writer,
			Line:   128,
		},
	}
}

func NewMultipartEncoder(writer *bytes.Buffer, inputFileSize int) *MultipartEncoder {
	return &MultipartEncoder{
		Encoder: &Encoder{
			Writer: writer,
			Line:   128,
			Size:   inputFileSize,
		},
	}
}
