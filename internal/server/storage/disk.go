package storage

import (
	"encoding/json"
	"os"
)

func (settings MemStorage) Save(fname string) error {
	// сериализуем структуру в JSON формат
	data, err := json.MarshalIndent(settings, "", "   ")
	if err != nil {
		return err
	}
	// сохраняем данные в файл
	return os.WriteFile(fname, data, 0666)
}

func (settings *MemStorage) Load(fname string) error {
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
