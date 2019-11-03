package stream_io

import (
	"encoding/binary"
	"fmt"
	"io"
)

func WriteString(writer io.Writer, data string) (err error) {
	dataAsByteSlice := []byte(data)
	var dataLen int32
	dataLen = int32(len(dataAsByteSlice))
	err = binary.Write(writer, binary.BigEndian, dataLen)
	if err != nil {
		return
	}
	_, err = writer.Write(dataAsByteSlice)
	return
}

func ReadString(reader io.Reader, maxSize int) (data string, err error) {
	// read the size
	var sliceSize int32
	err = binary.Read(reader, binary.BigEndian, &sliceSize)
	if err != nil {
		return
	}
	if int32(maxSize) < sliceSize {
		err = fmt.Errorf("buffer too small, expected: %d but got: %d", maxSize, sliceSize)
		return
	}
	buffer := make([]byte, sliceSize)
	_, err = reader.Read(buffer)
	return string(buffer), err
}
