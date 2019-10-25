package testutil

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func mockPrintMetrics(f func(metrics []*prommodel.MetricFamily) error) func() {
	originalPrintMetrics := printMetrics
	printMetrics = f

	return func() {
		printMetrics = originalPrintMetrics
	}
}

func TestMustPrintMetrics(t *testing.T) {
	var inParam []*prommodel.MetricFamily
	defer mockPrintMetrics(func(metrics []*prommodel.MetricFamily) error {
		inParam = metrics
		return nil
	})()

	metrics := []*prommodel.MetricFamily{}
	MustPrintMetrics(metrics)
	assert.Equal(t, metrics, inParam)
}

func TestMustPrintMetrics_panic(t *testing.T) {
	expectedPanicValue := fmt.Errorf("sample error")
	unmockFunc := mockPrintMetrics(func(metrics []*prommodel.MetricFamily) error {
		return expectedPanicValue
	})

	assert.PanicsWithValue(t, expectedPanicValue, func() {
		defer unmockFunc()
		MustPrintMetrics(nil)
	})
}

func TestPrintMetrics(t *testing.T) {
	var inParam []*prommodel.MetricFamily
	expectedErr := fmt.Errorf("sample error")
	defer mockPrintMetrics(func(metrics []*prommodel.MetricFamily) error {
		inParam = metrics
		return expectedErr
	})()

	metrics := []*prommodel.MetricFamily{}
	err := PrintMetrics(metrics)

	t.Run("parameter", func(t *testing.T) {
		assert.Equal(t, metrics, inParam)
	})
	t.Run("error", func(t *testing.T) {
		assert.Equal(t, expectedErr, err)
	})
}

func Example_printMetrics() {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "metric1",
		Help: "metric1 help.",
	})
	counter.Inc()

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(counter)

	metrics, _ := registry.Gather()
	_ = printMetrics(metrics)
	// Output:
	// # HELP metric1 metric1 help.
	// # TYPE metric1 counter
	// metric1 1
}

func Example_printMetrics_nil() {
	_ = printMetrics(nil)
	// Output:
}
