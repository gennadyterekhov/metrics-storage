package iohelpler

import (
	"compress/gzip"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"io"
)

func ReadFromReaderOrDie(reader io.Reader) []byte {
	readBytes, err := io.ReadAll(reader)
	if err != nil {
		logger.ZapSugarLogger.Panicln("error when reading")
	}

	return readBytes
}

func ReadFromGzipReaderOrDie(reader io.Reader) []byte {
	gzipBodyReader, err := gzip.NewReader(reader)
	if err != nil {
		logger.ZapSugarLogger.Panicln("error when creating gzip reader")
	}
	err = gzipBodyReader.Close()
	if err != nil {
		logger.ZapSugarLogger.Panicln("error when closing")
	}
	bts := ReadFromReaderOrDie(gzipBodyReader)

	return bts
}

func ReadFromReadCloserOrDie(reader io.ReadCloser) []byte {
	readBytes := ReadFromReaderOrDie(reader)
	err := reader.Close()
	if err != nil {
		logger.ZapSugarLogger.Panicln("error when closing")
	}

	return readBytes
}

func ReadFromGzipReadCloserOrDie(reader io.ReadCloser) []byte {
	readBytes := ReadFromGzipReaderOrDie(reader)
	err := reader.Close()
	if err != nil {
		logger.ZapSugarLogger.Panicln("error when closing")
	}

	return readBytes
}
