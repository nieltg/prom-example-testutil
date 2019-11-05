package testutil

import (
	"io"
	"os"

	prommodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type printer interface {
	PrintMetrics(metrics []*prommodel.MetricFamily) error
}

type printerImpl struct {
	newEncoderFunc func(w io.Writer, format expfmt.Format) expfmt.Encoder
}

func (impl *printerImpl) PrintMetrics(metrics []*prommodel.MetricFamily) error {
	encoder := impl.newEncoderFunc(os.Stdout, expfmt.FmtText)
	for _, metric := range metrics {
		if err := encoder.Encode(metric); err != nil {
			return err
		}
	}
	return nil
}

var globalPrinter printer = &printerImpl{
	newEncoderFunc: expfmt.NewEncoder,
}

// MustPrintMetrics prints metrics or panic if error has occurred.
func MustPrintMetrics(metrics []*prommodel.MetricFamily) {
	if err := globalPrinter.PrintMetrics(metrics); err != nil {
		panic(err)
	}
}

// PrintMetrics prints metrics, otherwise return error.
func PrintMetrics(metrics []*prommodel.MetricFamily) error {
	return globalPrinter.PrintMetrics(metrics)
}
