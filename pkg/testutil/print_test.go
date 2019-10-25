package testutil

import (
	"testing"

	prommodel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestMustPrintMetrics(t *testing.T) {
	var inParam []*prommodel.MetricFamily
	printMetrics = func(metrics []*prommodel.MetricFamily) error {
		inParam = metrics
		return nil
	}

	metrics := []*prommodel.MetricFamily{}
	MustPrintMetrics(metrics)
	assert.Equal(t, metrics, inParam)
}
