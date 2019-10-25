package testutil

import (
	prommodel "github.com/prometheus/client_model/go"
)

var printMetrics = func(metrics []*prommodel.MetricFamily) error {
	return nil
}

// MustPrintMetrics prints metrics or panic if error has occured.
func MustPrintMetrics(metrics []*prommodel.MetricFamily) {
	if err := printMetrics(metrics); err != nil {
		panic(err)
	}
}
