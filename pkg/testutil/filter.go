package testutil

import prommodel "github.com/prometheus/client_model/go"

type filterer interface {
	FilterMetricsByName(
		metrics []*prommodel.MetricFamily,
		names ...string,
	) []*prommodel.MetricFamily
}

type filtererImpl struct{}

func (filtererImpl) FilterMetricsByName(
	metrics []*prommodel.MetricFamily,
	names ...string,
) (filteredMetrics []*prommodel.MetricFamily) {
	for _, metric := range metrics {
		for _, name := range names {
			if metric.GetName() == name {
				filteredMetrics = append(filteredMetrics, metric)
			}
		}
	}
	return
}

var globalFilterer filterer = &filtererImpl{}

// FilterMetricsByName filters metrics by specified names.
func FilterMetricsByName(
	metrics []*prommodel.MetricFamily,
	names ...string,
) []*prommodel.MetricFamily {
	return globalFilterer.FilterMetricsByName(metrics, names...)
}
