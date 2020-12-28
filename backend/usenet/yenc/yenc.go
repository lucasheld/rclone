package yenc

// http://www.yenc.org/yenc-draft.1.3.txt
// http://www.yenc.org/develop.htm

const (
	null  byte = 0x00
	tab   byte = 0x09
	lf    byte = 0x0A
	cr    byte = 0x0D
	space byte = 0x20
	dot   byte = 0x2E
	equal byte = 0x3D
)

type Part struct {
	Part  int
	Begin int
	End   int
	Crc   string
}

func (p Part) Size() int {
	return p.End - p.Begin + 1
}

func (p *Part) IsMultipart() bool {
	return p.Part > 0
}
