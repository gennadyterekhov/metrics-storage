package repositories

import (
	"context"

	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

// deprecated
type MetricsRepository interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)

	GetAll() (map[string]float64, map[string]int64)

	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
}

type RepositoryInterface interface {
	GetAll(ctx context.Context) (map[string]float64, map[string]int64)

	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)

	SetGauge(ctx context.Context, name string, value float64)
	AddCounter(ctx context.Context, name string, value int64)
	SaveToDisk(ctx context.Context, fileStorage string) error
	GetAllGauges(ctx context.Context) map[string]float64
	GetAllCounters(ctx context.Context) map[string]int64
}

type Repository struct {
	stor storage.StorageInterface
}

func (r Repository) GetAllGauges(ctx context.Context) map[string]float64 {
	return r.stor.GetAllGauges(ctx)
}

func (r Repository) GetAllCounters(ctx context.Context) map[string]int64 {
	return r.stor.GetAllCounters(ctx)
}

func (r Repository) SaveToDisk(ctx context.Context, fileStorage string) error {
	return r.stor.SaveToDisk(ctx, fileStorage)
}

func New(stor storage.StorageInterface) Repository {
	return Repository{
		stor: stor,
	}
}

func (r Repository) GetAll(ctx context.Context) (map[string]float64, map[string]int64) {
	return r.stor.GetAllGauges(ctx), r.stor.GetAllCounters(ctx)
}

func (r Repository) GetGauge(ctx context.Context, name string) (float64, error) {
	return r.stor.GetGauge(ctx, name)
}

func (r Repository) GetCounter(ctx context.Context, name string) (int64, error) {
	return r.stor.GetCounter(ctx, name)
}

func (r Repository) SetGauge(ctx context.Context, name string, value float64) {
	r.stor.SetGauge(ctx, name, value)
}

func (r Repository) AddCounter(ctx context.Context, name string, value int64) {
	r.stor.AddCounter(ctx, name, value)
}

func (r Repository) Clear() {
	r.stor.Clear()
}

func (r Repository) GetCounterOrZero(ctx context.Context, name string) int64 {
	v, err := r.stor.GetCounter(ctx, name)
	if err != nil {
		return 0
	}
	return v
}

func (r Repository) GetGaugeOrZero(ctx context.Context, name string) float64 {
	v, err := r.stor.GetGauge(ctx, name)
	if err != nil {
		return 0
	}
	return v
}
