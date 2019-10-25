package testutil

import (
	"os"

	prommodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var newEncoder = expfmt.NewEncoder

var printMetrics = func(metrics []*prommodel.MetricFamily) error {
	encoder := newEncoder(os.Stdout, expfmt.FmtText)
	if len(metrics) > 0 {
		encoder.Encode(metrics[0])
	}
	return nil
}

// MustPrintMetrics prints metrics or panic if error has occured.
func MustPrintMetrics(metrics []*prommodel.MetricFamily) {
	if err := printMetrics(metrics); err != nil {
		panic(err)
	}
}

// PrintMetrics prints metrics, otherwise return error.
func PrintMetrics(metrics []*prommodel.MetricFamily) error {
	return printMetrics(metrics)
}
