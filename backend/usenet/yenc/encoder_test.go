package yenc

import (
	"bufio"
	"bytes"
	"fmt"
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

	inputBuffer := new(bytes.Buffer)
	inputBuffer.Write(input)
	inputReader := bufio.NewReader(inputBuffer)

	outputBuffer := new(bytes.Buffer)
	encoder := NewEncoder(outputBuffer, size)
	encoder.Name = "joystick.jpg"

	chunksize := 11250
	chunkBuffer := make([]byte, chunksize)

	number := 1
	begin := 1
	for {
		n, err := inputReader.Read(chunkBuffer)
		if err != nil {
			break
		}
		chunk := chunkBuffer[0:n]

		part := &Part{
			Part:  number,
			Begin: begin,
			End:   begin+n-1,
		}

		err = encoder.Encode(part, chunk)
		assert.NoError(t, err)

		assert.Equal(t, encoder.Size, size)

		output, err := ioutil.ReadAll(outputBuffer)
		assert.NoError(t, err)
		filename := fmt.Sprintf("testdata/multipart/0000002%d.ntx", number-1)
		outputExpected, err := ioutil.ReadFile(filename)
		outputExpected = bytes.ReplaceAll(outputExpected, []byte{' ', '\r', '\n'}, []byte{'\r', '\n'})
		assert.NoError(t, err)
		assert.Equal(t, outputExpected, output)

		number++
		begin = begin + n
	}
}
