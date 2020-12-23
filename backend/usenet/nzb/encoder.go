package nzb

import (
	"encoding/xml"
	"io"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e Encoder) Encode(nzb *Nzb) error {
	header := []byte(xml.Header + "<!DOCTYPE nzb PUBLIC \"-//newzBin//DTD NZB 1.1//EN\" \"http://www.newzbin.com/DTD/nzb/nzb-1.1.dtd\">" + "\n")
	_, err := e.w.Write(header)
	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(e.w)
	encoder.Indent("", "\t")
	err = encoder.Encode(nzb)
	if err != nil {
		return err
	}
	return nil
}
