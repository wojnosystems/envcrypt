package main

import (
	"io"
	"os"
)

// prepareStreamsForEncrypt by reading the paths and creating a reader and writer
func prepareStreamsForEncrypt(inputPath, outputPath string) (inputStream io.ReadCloser, outputStream io.WriteCloser, err error) {
	// perform clean-up on error
	defer func() {
		if err != nil {
			if inputStream != nil {
				_ = inputStream.Close()
			}
			if outputStream != nil {
				_ = outputStream.Close()
			}
		}
	}()
	inputStream, err = os.OpenFile(inputPath, os.O_RDONLY, 0400)
	if err != nil {
		return
	}
	outputStream, err = os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	return
}
