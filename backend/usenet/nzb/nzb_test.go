package nzb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNzb(t *testing.T) {
	nzb, err := ParseNzb("testdata/test.nzb")
	assert.NoError(t, err)

	assert.Equal(t, len(nzb.Files), 46)
	nzbFile := nzb.Files[0]
	assert.Equal(t, nzbFile.Poster, "NewsUP <NewsUP@somewhere.cbr>")
	assert.Equal(t, nzbFile.Date, 1487587920)
	assert.Equal(t, nzbFile.Subject, "[02/45] - \"ubuntu-mate-16.04.2-desktop-amd64.iso.part01.rar\" (1/66)")

	assert.Equal(t, len(nzbFile.Groups), 1)
	assert.Equal(t, nzbFile.Groups[0], "alt.binaries.test")

	assert.Equal(t, len(nzbFile.Segments), 66)
	nzbSegment := nzbFile.Segments[0]
	assert.Equal(t, nzbSegment.Bytes, 787135)
	assert.Equal(t, nzbSegment.Number, 1)
	assert.Equal(t, nzbSegment.Id, "VNdHifYKSAtPRNYRoIhyQLghUETQLDaC@UFK")
}
