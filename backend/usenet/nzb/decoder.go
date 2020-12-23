package nzb

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func (d *Decoder) Decode(nzb *Nzb) error {
	decoder := xml.NewDecoder(d.r)
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&nzb)
	if err != nil {
		return err
	}
	return nil
}
