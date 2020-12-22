package nzb

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
)

type Nzb struct {
	XMLName xml.Name `xml:"http://www.newzbin.com/DTD/2003/nzb nzb"`
	Files   []File   `xml:"file"`
}

type File struct {
	Poster   string    `xml:"poster,attr"`
	Date     int       `xml:"date,attr"`
	Subject  string    `xml:"subject,attr"`
	Groups   []string  `xml:"groups>group"`
	Segments []Segment `xml:"segments>segment"`
}

type Segment struct {
	Bytes  int    `xml:"bytes,attr"`
	Number int    `xml:"number,attr"`
	Id     string `xml:",chardata"`
}

func EncodeNzb(nzb *Nzb) (content []byte, err error) {
	buffer := new(bytes.Buffer)

	buffer.WriteString(xml.Header)
	buffer.WriteString("<!DOCTYPE nzb PUBLIC \"-//newzBin//DTD NZB 1.1//EN\" \"http://www.newzbin.com/DTD/nzb/nzb-1.1.dtd\">" + "\n")

	encoder := xml.NewEncoder(buffer)
	encoder.Indent("", "\t")
	err = encoder.Encode(&nzb)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DecodeNzb(content []byte) (nzb *Nzb, err error) {
	buffer := new(bytes.Buffer)
	buffer.Write(content)

	decoder := xml.NewDecoder(buffer)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&nzb)
	if err != nil {
		return nil, err
	}
	return nzb, nil
}
