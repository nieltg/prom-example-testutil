package testutil

import prommodel "github.com/prometheus/client_model/go"

type filterer interface {
	FilterMetricsByName(
		metrics []*prommodel.MetricFamily, name string) []*prommodel.MetricFamily
}

type filtererImpl struct{}

func (filtererImpl) FilterMetricsByName(
	metrics []*prommodel.MetricFamily, name string) []*prommodel.MetricFamily {
	if metrics[0].GetName() == name {
		return metrics
	}
	return nil
}

var globalFilterer filterer = &filtererImpl{}

// FilterMetricsByName ...
func FilterMetricsByName(
	metrics []*prommodel.MetricFamily, name string) []*prommodel.MetricFamily {
	return globalFilterer.FilterMetricsByName(metrics, name)
}
