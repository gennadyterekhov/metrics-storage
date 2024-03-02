package iohelpler

import (
	"compress/gzip"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"io"
)

func ReadFromReaderPlain(reader *io.Reader) []byte {
	//bytesBuffer := bytes.NewBuffer(nil)
	bytesSlice := make([]byte, 0)
	_, err := (*reader).Read(bytesSlice)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when reading")
	}

	return bytesSlice
}

func ReadFromReader(reader *io.Reader) []byte {
	readBytes, err := io.ReadAll(*reader)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when reading")
	}
	logger.ZapSugarLogger.Debugln("readBytes", string(readBytes))

	return readBytes
}

func ReadFromReaderOrDie(reader io.Reader) []byte {
	readBytes, err := io.ReadAll(reader)
	if err != nil {
		logger.ZapSugarLogger.Panicln("error when reading")
	}
	logger.ZapSugarLogger.Debugln("readBytes", string(readBytes))

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
