package yenc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestSinglepart(t *testing.T) {
	input, err := ioutil.ReadFile("testdata/singlepart/testfile.txt")
	assert.NoError(t, err)
	size := len(input)

	outputBuffer := new(bytes.Buffer)
	encoder := NewEncoder(outputBuffer, size)
	encoder.Name = "testfile.txt"


	part := &Part{}
	err = encoder.Encode(part, input)
	assert.NoError(t, err)

	assert.Equal(t, encoder.Size, size)
	assert.Equal(t, part.Crc, "ded29f4f")

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)

	outputExpected, err := ioutil.ReadFile("testdata/singlepart/00000005.ntx")
	outputExpected = bytes.ReplaceAll(outputExpected, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
	assert.NoError(t, err)

	assert.Equal(t, outputExpected, output)
}

func TestMultipart(t *testing.T) {
	input, err := ioutil.ReadFile("testdata/multipart/joystick.jpg")
	assert.NoError(t, err)
	size := len(input)

	outputBuffer := new(bytes.Buffer)
	encoder := NewEncoder(outputBuffer, size)
	encoder.Name = "joystick.jpg"


	part1 := &Part{
		Part:  1,
		Begin: 1,
		End:   11250,
	}
	err = encoder.Encode(part1, input[part1.Begin-1:part1.End])
	assert.NoError(t, err)

	assert.Equal(t, encoder.Size, size)
	assert.Equal(t, part1.Size(), 11250)
	assert.Equal(t, part1.Crc, "bfae5c0b")

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)
	outputExpected1, err := ioutil.ReadFile("testdata/multipart/00000020.ntx")
	outputExpected1 = bytes.ReplaceAll(outputExpected1, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
	assert.NoError(t, err)
	assert.Equal(t, outputExpected1, output)


	part2 := &Part{
		Part:  2,
		Begin: part1.End+1,
		End:   encoder.Size,
	}
	err = encoder.Encode(part2, input[part2.Begin-1:part2.End])
	assert.NoError(t, err)

	assert.Equal(t, encoder.Size, size)
	assert.Equal(t, part2.Size(), 8088)
	assert.Equal(t, part2.Crc, "aca76043")

	output, err = ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)
	outputExpected2, err := ioutil.ReadFile("testdata/multipart/00000021.ntx")
	outputExpected2 = bytes.ReplaceAll(outputExpected2, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
	assert.NoError(t, err)
	assert.Equal(t, outputExpected2, output)
}
