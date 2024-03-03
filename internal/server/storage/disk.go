package storage

import (
	"encoding/json"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"os"
)

func (strg *MemStorage) Save(filename string) error {
	data, err := json.MarshalIndent(strg, "", "   ")
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when saving metrics to disk")
		return err
	}

	return os.WriteFile(filename, data, 0666)
}

func (strg *MemStorage) Load(fname string) error {
	fbytes, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	// прочитайте файл с помощью os.ReadFile
	// десериализуйте данные используя json.Unmarshal
	// ...
	tempSettings := &MemStorage{}
	json.Unmarshal(fbytes, tempSettings)
}
