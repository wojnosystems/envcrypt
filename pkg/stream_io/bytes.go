package stream_io

import (
	"encoding/binary"
	"fmt"
	"io"
)

func WriteByteSlice(writer io.Writer, data []byte) (err error) {
	var dataLen int32
	dataLen = int32(len(data))
	err = binary.Write(writer, binary.BigEndian, dataLen)
	if err != nil {
		return
	}
	_, err = writer.Write(data)
	return
}

func ReadIntoByteSlice(reader io.Reader, buffer []byte) (bufferRet []byte, err error) {
	// read the size
	var sliceSize int32
	err = binary.Read(reader, binary.BigEndian, &sliceSize)
	if err != nil {
		return
	}
	if int32(cap(buffer)) < sliceSize {
		err = fmt.Errorf("buffer too small, expected: %d but got: %d", sliceSize, cap(buffer))
		return
	}
	_, err = reader.Read(buffer[0:sliceSize])
	return buffer[0:sliceSize], nil
}

func ReadByteSlice(reader io.Reader, maxBufferSize int) (bufferRet []byte, err error) {
	// read the size
	var sliceSize int32
	err = binary.Read(reader, binary.BigEndian, &sliceSize)
	if err != nil {
		return
	}
	if maxBufferSize < int(sliceSize) {
		err = fmt.Errorf("buffer too small, expected: %d but got: %d", sliceSize, maxBufferSize)
		return
	}
	bufferRet = make([]byte, sliceSize)
	_, err = reader.Read(bufferRet)
	return
}
