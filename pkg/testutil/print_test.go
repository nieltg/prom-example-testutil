package testutil

import (
	"fmt"
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

func TestMustPrintMetrics_panic(t *testing.T) {
	expectedPanicValue := fmt.Errorf("sample error")
	printMetrics = func(metrics []*prommodel.MetricFamily) error {
		return expectedPanicValue
	}

	assert.PanicsWithValue(t, expectedPanicValue, func() {
		MustPrintMetrics(nil)
	})
}

func TestPrintMetrics(t *testing.T) {
	var inParam []*prommodel.MetricFamily
	expectedErr := fmt.Errorf("sample error")
	printMetrics = func(metrics []*prommodel.MetricFamily) error {
		inParam = metrics
		return expectedErr
	}

	metrics := []*prommodel.MetricFamily{}
	err := PrintMetrics(metrics)

	t.Run("parameter", func(t *testing.T) {
		assert.Equal(t, metrics, inParam)
	})
	t.Run("error", func(t *testing.T) {
		assert.Equal(t, expectedErr, err)
	})
}
