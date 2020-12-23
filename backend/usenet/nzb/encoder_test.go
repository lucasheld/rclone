package nzb

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

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

	buffer := &bytes.Buffer{}
	encoder := NewEncoder(buffer)
	err := encoder.Encode(nzb)
	assert.NoError(t, err)
	contentStr := string(buffer.Bytes())

	dataExpected, err := ioutil.ReadFile("testdata/encode.nzb")
	assert.NoError(t, err)
	contentStrExpected := string(dataExpected)
	contentStrExpected = strings.ReplaceAll(contentStrExpected, "1487587920", strconv.Itoa(nzb.Files[0].Date))
	contentStrExpected = strings.ReplaceAll(contentStrExpected, "1487587921", strconv.Itoa(nzb.Files[1].Date))

	assert.Equal(t, contentStr, contentStrExpected)
}
