package nzb

import (
	"encoding/xml"
	"time"
)

type Nzb struct {
	XMLName xml.Name `xml:"http://www.newzbin.com/DTD/2003/nzb nzb"`
	Files   []*File  `xml:"file"`
}

type File struct {
	Poster   string     `xml:"poster,attr"`
	Date     int        `xml:"date,attr"`
	Subject  string     `xml:"subject,attr"`
	Groups   []string   `xml:"groups>group"`
	Segments []*Segment `xml:"segments>segment"`
}

type Segment struct {
	Bytes  int    `xml:"bytes,attr"`
	Number int    `xml:"number,attr"`
	Id     string `xml:",chardata"`
}

func (nzb *Nzb) AddFile(poster string, subject string, groups []string) (file *File) {
	date := int(time.Now().Unix())
	file = &File{
		Poster:   poster,
		Date:     date,
		Subject:  subject,
		Groups:   groups,
		Segments: []*Segment{},
	}
	nzb.Files = append(nzb.Files, file)
	return file
}

func (file *File) AddSegment(bytes int, id string) (segment *Segment) {
	number := len(file.Segments) + 1
	segment = &Segment{
		Bytes:  bytes,
		Number: number,
		Id:     id,
	}
	file.Segments = append(file.Segments, segment)
	return segment
}
