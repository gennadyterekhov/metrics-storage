package iohelpler

import (
	"compress/gzip"
	"io"
	"os"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

func GetFileContents(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ReadFromReaderOrDie(reader io.Reader) []byte {
	readBytes, err := io.ReadAll(reader)
	if err != nil {
		logger.Custom.Panicln("error when reading", err.Error())
	}

	return readBytes
}

func ReadFromGzipReaderOrDie(reader io.Reader) []byte {
	gzipBodyReader, err := gzip.NewReader(reader)
	if err != nil {
		logger.Custom.Panicln("error when creating gzip reader", err.Error())
	}
	err = gzipBodyReader.Close()
	if err != nil {
		logger.Custom.Panicln("error when closing", err.Error())
	}
	bts := ReadFromReaderOrDie(gzipBodyReader)

	return bts
}

func ReadFromReadCloserOrDie(reader io.ReadCloser) []byte {
	readBytes := ReadFromReaderOrDie(reader)
	CloseOrPanic(reader)

	return readBytes
}

func ReadFromGzipReadCloserOrDie(reader io.ReadCloser) []byte {
	readBytes := ReadFromGzipReaderOrDie(reader)
	CloseOrPanic(reader)

	return readBytes
}

func CloseOrPanic(reader io.ReadCloser) {
	err := reader.Close()
	if err != nil {
		logger.Custom.Panicln("error when closing", err.Error())
	}
}
