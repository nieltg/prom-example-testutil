package testutil

import (
	"os"

	prommodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var newEncoder = expfmt.NewEncoder

var printMetrics = func(metrics []*prommodel.MetricFamily) error {
	encoder := newEncoder(os.Stdout, expfmt.FmtText)
	for _, metric := range metrics {
		if err := encoder.Encode(metric); err != nil {
			return err
		}
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
