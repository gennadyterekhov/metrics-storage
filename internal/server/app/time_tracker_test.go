package app

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestInterval(t *testing.T) {
	filename := "test_interval.json"
	config.Conf.StoreInterval = 2
	config.Conf.FileStorage = filename
	StartTrackingIntervals()

	time.Sleep(time.Second * 3)

	assert.Equal(t, 1, TimeTrackerInstance.ActionFulfilled)

	if err := os.Remove(filename); err != nil {
		t.Error(err)
	}

}
