package testutil

import prommodel "github.com/prometheus/client_model/go"

type filterer interface {
	FilterMetricsByName(
		metrics []*prommodel.MetricFamily, name string) []*prommodel.MetricFamily
}

type filtererImpl struct{}

func (filtererImpl) FilterMetricsByName(
	metrics []*prommodel.MetricFamily,
	name string,
) (filteredMetrics []*prommodel.MetricFamily) {
	for _, metric := range metrics {
		if metric.GetName() == name {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}
	return
}

var globalFilterer filterer = &filtererImpl{}

// FilterMetricsByName ...
func FilterMetricsByName(
	metrics []*prommodel.MetricFamily, name string) []*prommodel.MetricFamily {
	return globalFilterer.FilterMetricsByName(metrics, name)
}
