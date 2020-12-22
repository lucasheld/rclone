package nzb

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestDecodeNzb(t *testing.T) {
	file, err := os.Open("testdata/decode.nzb")
	assert.NoError(t, err)
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	assert.NoError(t, err)

	nzb, err := DecodeNzb(content)
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

func TestEncodeNzb(t *testing.T) {
	poster := "NewsUP <NewsUP@somewhere.cbr>"
	groups := []string{
		"alt.binaries.test",
		"alt.binaries.ath",
	}

	nzb := &Nzb{}
	file := nzb.AddFile(
		poster,
		"[02/45] - \"ubuntu-mate-16.04.2-desktop-amd64.iso.part01.rar\" (1/66)",
		groups,
	)
	file.AddSegment(787135, "VNdHifYKSAtPRNYRoIhyQLghUETQLDaC@UFK")
	file.AddSegment(793807, "zmvIGqoSHmSwAVMlkJnPppKRtwezvAIL@UFK")
	file = nzb.AddFile(
		poster,
		"[01/45] - &quot;ubuntu-mate-16.04.2-desktop-amd64.iso.par2&quot; (1/1)",
		groups,
	)
	file.AddSegment(51058, "BcvtngyhrBKDHVsgCmTcgeCnbTaxzkSL@UFK")

	content, err := EncodeNzb(nzb)
	assert.NoError(t, err)
	contentStr := string(content)

	fileExpected, err := os.Open("testdata/encode.nzb")
	assert.NoError(t, err)
	dataExpected, err := ioutil.ReadAll(fileExpected)
	assert.NoError(t, err)
	contentStrExpected := string(dataExpected)
	contentStrExpected = strings.ReplaceAll(contentStrExpected, "1487587920", strconv.Itoa(nzb.Files[0].Date))
	contentStrExpected = strings.ReplaceAll(contentStrExpected, "1487587921", strconv.Itoa(nzb.Files[1].Date))

	assert.Equal(t, contentStr, contentStrExpected)
}
