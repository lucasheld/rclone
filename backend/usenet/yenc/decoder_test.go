package yenc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestSinglepartDecoder(t *testing.T) {
	input, err := ioutil.ReadFile("testdata/singlepart/00000005.ntx")
	assert.NoError(t, err)

	outputBuffer := new(bytes.Buffer)
	decoder := NewDecoder(outputBuffer)

	err = decoder.Decode(input)
	assert.NoError(t, err)

	outputExpected, err := ioutil.ReadFile("testdata/singlepart/testfile.txt")
	assert.NoError(t, err)

	assert.Equal(t, decoder.Line, 128)
	assert.Equal(t, decoder.Size, len(outputExpected))
	assert.Equal(t, decoder.Name, "testfile.txt")
	assert.Equal(t, decoder.Crc, "ded29f4f")
	assert.Equal(t, decoder.Total, 0)
	assert.Equal(t, len(decoder.Parts), 0)

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)

	assert.Equal(t, outputExpected, output)
}
