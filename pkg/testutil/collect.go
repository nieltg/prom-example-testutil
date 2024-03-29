package testutil

import (
	"github.com/prometheus/client_golang/prometheus"
)

type registererGatherer interface {
	prometheus.Registerer
	prometheus.Gatherer
}

type collector interface {
	MustCollect(...prometheus.Collector) prometheus.Gatherer
}

type collectorImpl struct {
	newRegistryFunc func() registererGatherer
}

var globalCollector collector = &collectorImpl{
	newRegistryFunc: func() registererGatherer {
		return prometheus.NewPedanticRegistry()
	},
}

func (impl *collectorImpl) MustCollect(
	collectors ...prometheus.Collector,
) prometheus.Gatherer {
	registry := impl.newRegistryFunc()
	registry.MustRegister(collectors...)
	return registry
}

// MustCollect collects metrics from collectors with pedantic checks.
func MustCollect(collectors ...prometheus.Collector) prometheus.Gatherer {
	return globalCollector.MustCollect(collectors...)
}
