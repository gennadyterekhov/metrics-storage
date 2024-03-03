package storage

import (
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"os"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	container.MetricsRepository.AddCounter("c1", 1)
	container.MetricsRepository.AddCounter("c2", 2)
	container.MetricsRepository.AddCounter("c3", 3)
	container.MetricsRepository.SetGauge("g1", 1.5)
	container.MetricsRepository.SetGauge("g2", 2.6)
	container.MetricsRepository.SetGauge("g3", 3.7)

	filename := `metrics.json`

	if err := container.MetricsRepository.Save(filename); err != nil {
		t.Error(err)
	}
	var result MemStorage
	if err := (&result).Load(filename); err != nil {
		t.Error(err)
	}
	if !container.MetricsRepository.IsEqual(&result) {
		t.Errorf(`%+v не равно %+v`, container.MetricsRepository, result)
	}

	if err := os.Remove(filename); err != nil {
		t.Error(err)
	}
}
