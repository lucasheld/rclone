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
	assert.Equal(t, decoder.PPart, 0)
	assert.Equal(t, decoder.PTotal, 0)
	assert.Equal(t, decoder.PBegin, 0)
	assert.Equal(t, decoder.PEnd, 0)
	assert.Equal(t, decoder.PSize, 0)
	assert.Equal(t, decoder.PCrc, "")

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)

	assert.Equal(t, outputExpected, output)
}

func TestMultipartDecoder(t *testing.T) {
	outputExpected, err := ioutil.ReadFile("testdata/multipart/joystick.jpg")
	assert.NoError(t, err)
	sizeTotal := len(outputExpected)

	input, err := ioutil.ReadFile("testdata/multipart/00000020.ntx")
	assert.NoError(t, err)

	outputBuffer := new(bytes.Buffer)
	decoder := NewDecoder(outputBuffer)
	err = decoder.Decode(input)
	assert.NoError(t, err)

	assert.Equal(t, decoder.Line, 128)
	assert.Equal(t, decoder.Size, sizeTotal)
	assert.Equal(t, decoder.Name, "joystick.jpg")
	assert.Equal(t, decoder.Crc, "")
	assert.Equal(t, decoder.PPart, 1)
	assert.Equal(t, decoder.PTotal, 0)
	assert.Equal(t, decoder.PBegin, 1)
	assert.Equal(t, decoder.PEnd, 11250)
	assert.Equal(t, decoder.PCrc, "bfae5c0b")
	assert.Equal(t, decoder.PSize, 11250)

	output, err := ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)

	assert.Equal(t, outputExpected[decoder.PBegin-1:decoder.PEnd], output)

	input, err = ioutil.ReadFile("testdata/multipart/00000021.ntx")
	assert.NoError(t, err)

	outputBuffer = new(bytes.Buffer)
	decoder = NewDecoder(outputBuffer)
	err = decoder.Decode(input)
	assert.NoError(t, err)

	assert.Equal(t, decoder.Line, 128)
	assert.Equal(t, decoder.Size, sizeTotal)
	assert.Equal(t, decoder.Name, "joystick.jpg")
	assert.Equal(t, decoder.Crc, "")
	assert.Equal(t, decoder.PPart, 2)
	assert.Equal(t, decoder.PTotal, 0)
	assert.Equal(t, decoder.PBegin, 11251)
	assert.Equal(t, decoder.PEnd, 19338)
	assert.Equal(t, decoder.PCrc, "aca76043")
	assert.Equal(t, decoder.PSize, 8088)

	output, err = ioutil.ReadAll(outputBuffer)
	assert.NoError(t, err)

	assert.Equal(t, outputExpected[decoder.PBegin-1:decoder.PEnd], output)
}
