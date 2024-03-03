package storage

import (
	"os"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	MetricsRepository.AddCounter("c1", 1)
	MetricsRepository.AddCounter("c2", 2)
	MetricsRepository.AddCounter("c3", 3)
	MetricsRepository.SetGauge("g1", 1.5)
	MetricsRepository.SetGauge("g2", 2.6)
	MetricsRepository.SetGauge("g3", 3.7)

	filename := `metrics.json`

	if err := MetricsRepository.Save(filename); err != nil {
		t.Error(err)
	}
	var result MemStorage
	if err := (&result).Load(filename); err != nil {
		t.Error(err)
	}
	if !MetricsRepository.IsEqual(&result) {
		t.Errorf(`%+v не равно %+v`, MetricsRepository, result)
	}

	if err := os.Remove(filename); err != nil {
		t.Error(err)
	}
}
