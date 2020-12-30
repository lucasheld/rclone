package yenc

import (
	"bytes"
	"strconv"
)

type Decoder struct {
	Writer *bytes.Buffer
	Line   int
	Size   int
	Name   string
	Crc    string
	Total  int
	Parts  []*Part
}

func (d *Decoder) Part() *Part {
	lenParts := len(d.Parts)
	if lenParts > 0 {
		return d.Parts[len(d.Parts)-1]
	} else {
		return nil
	}
}

func (d *Decoder) parseParam(param []byte) (key string, value string) {
	paramSplit := bytes.SplitN(param, []byte{'='}, 2)
	key = string(paramSplit[0])
	value = string(paramSplit[1])
	return key, value
}

func (d *Decoder) parseLineBegin(params [][]byte) error {
	for _, param := range params {
		key, value := d.parseParam(param)

		var err error
		switch key {
		case "part":
			// add new part
			d.Parts = append(d.Parts, &Part{})
			d.Part().Part, err = strconv.Atoi(value)
			break
		case "total":
			d.Total, err = strconv.Atoi(value)
			break
		case "line":
			d.Line, err = strconv.Atoi(value)
			break
		case "size":
			d.Size, err = strconv.Atoi(value)
			break
		case "name":
			d.Name = value
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) parseLinePart(params [][]byte) error {
	for _, param := range params {
		key, value := d.parseParam(param)

		var err error
		switch key {
		case "begin":
			d.Part().Begin, err = strconv.Atoi(value)
			break
		case "end":
			d.Part().End, err = strconv.Atoi(value)
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) parseLineEnd(params [][]byte) error {
	for _, param := range params {
		key, value := d.parseParam(param)

		var err error
		switch key {
		case "size":
			part := d.Part()
			if part != nil {
				part.Size, err = strconv.Atoi(value)
			} else {
				d.Size, err = strconv.Atoi(value)
			}
			break
		case "part":
			d.Part().Part, err = strconv.Atoi(value)
			break
		case "pcrc32":
			d.Part().Crc = value
			break
		case "crc32":
			d.Crc = value
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) Decode(data []byte) error {
	lines := bytes.Split(data, []byte{'\n'})

	for _, line := range lines {
		line = bytes.TrimRight(line, "\r")
		line = bytes.Trim(line, " ")

		if line == nil {
			continue
		}

		if string(line[0:2]) == "=y" {
			split := bytes.Split(line[2:], []byte{' '})
			keyword := string(split[0])
			params := split[1:]

			switch keyword {
			case "begin":
				err := d.parseLineBegin(params)
				if err != nil {
					return err
				}
				break
			case "part":
				err := d.parseLinePart(params)
				if err != nil {
					return err
				}
				break
			case "end":
				err := d.parseLineEnd(params)
				if err != nil {
					return err
				}
				break
			}
		} else {
			for i := 0; i < len(line); i++ {
				char := line[i]

				if char == '=' {
					i++
					char = line[i]
					char = char - 64
				}
				char = char - 42
				err := d.Writer.WriteByte(char)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func NewDecoder(writer *bytes.Buffer) *Decoder {
	return &Decoder{
		Writer: writer,
	}
}
