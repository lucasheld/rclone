package yenc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestSinglepartEncoder(t *testing.T) {
	input, err := ioutil.ReadFile("testdata/singlepart/testfile.txt")
	assert.NoError(t, err)

	outputBuffer := new(bytes.Buffer)
	encoder := NewEncoder(outputBuffer)
	encoder.Size = len(input)
	encoder.Name = "testfile.txt"

	err = encoder.Encode(input)
	assert.NoError(t, err)

	assert.Equal(t, encoder.Size, len(input))
	assert.Equal(t, encoder.Crc, "ded29f4f")

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)

	outputExpected, err := ioutil.ReadFile("testdata/singlepart/00000005.ntx")
	outputExpected = bytes.ReplaceAll(outputExpected, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
	assert.NoError(t, err)

	assert.Equal(t, outputExpected, output)
}

func TestMultipartEncoder(t *testing.T) {
	input, err := ioutil.ReadFile("testdata/multipart/joystick.jpg")
	assert.NoError(t, err)

	chunk := input[0:11250]

	outputBuffer := new(bytes.Buffer)
	encoder := NewEncoder(outputBuffer)
	encoder.Size = len(input)
	encoder.Name = "joystick.jpg"
	// encoder.Crc = encoder.CalcCrc(input)
	encoder.PPart = 1
	// encoder.PTotal = 2
	encoder.PBegin = 1
	encoder.PEnd = 11250
	encoder.PSize = len(chunk)

	err = encoder.Encode(chunk)
	assert.NoError(t, err)

	assert.Equal(t, encoder.Line, 128)
	assert.Equal(t, encoder.Crc, "")
	assert.Equal(t, encoder.PTotal, 0)
	assert.Equal(t, encoder.PCrc, "bfae5c0b")

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)
	outputExpected, err := ioutil.ReadFile("testdata/multipart/00000020.ntx")
	outputExpected = bytes.ReplaceAll(outputExpected, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
	assert.NoError(t, err)
	assert.Equal(t, outputExpected, output)

	chunk = input[11250:19338]

	outputBuffer = new(bytes.Buffer)
	encoder = NewEncoder(outputBuffer)
	encoder.Size = len(input)
	encoder.Name = "joystick.jpg"
	// encoder.Crc = encoder.CalcCrc(input)
	encoder.PPart = 2
	// encoder.PTotal = 2
	encoder.PBegin = 11251
	encoder.PEnd = 19338
	encoder.PSize = len(chunk)

	err = encoder.Encode(chunk)
	assert.NoError(t, err)

	assert.Equal(t, encoder.Line, 128)
	assert.Equal(t, encoder.Crc, "")
	assert.Equal(t, encoder.PTotal, 0)
	assert.Equal(t, encoder.PCrc, "aca76043")

	output, err = ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)
	outputExpected, err = ioutil.ReadFile("testdata/multipart/00000021.ntx")
	outputExpected = bytes.ReplaceAll(outputExpected, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
	assert.NoError(t, err)
	assert.Equal(t, outputExpected, output)
}
