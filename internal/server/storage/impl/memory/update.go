package memorystorage

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Update(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) (err error) {
	for _, m := range gaugeMetrics {
		s.gaugeStorage.Set(m.ID, m)
	}

	for _, m := range counterMetrics {
		s.counterStorage.Set(m.ID, m)
	}

	return nil
}
