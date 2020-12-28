package yenc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestSinglepart(t *testing.T) {
	outputBuffer := new(bytes.Buffer)
	encoder := NewEncoder(outputBuffer)
	encoder.Name = "testfile.txt"

	input, err := ioutil.ReadFile("testdata/singlepart/testfile.txt")
	assert.NoError(t, err)

	err = encoder.Encode(input)
	assert.NoError(t, err)
	err = encoder.Close()
	assert.NoError(t, err)

	assert.Equal(t, encoder.Size, 584)
	assert.Equal(t, encoder.Crc32, "ded29f4f")

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)

	outputExpected, err := ioutil.ReadFile("testdata/singlepart/00000005.ntx")
	outputExpected = bytes.ReplaceAll(outputExpected, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
	assert.NoError(t, err)

	//for i := range outputExpected {
	//	oe := outputExpected[i]
	//	o := output[i]
	//	assert.Equal(t, int(oe), int(o))
	//}

	//assert.Equal(t, outputExpected, output)
	assert.Equal(t, string(outputExpected), string(output))
}
