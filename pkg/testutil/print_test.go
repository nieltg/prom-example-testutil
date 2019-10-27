package testutil

import (
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	mockexpfmt "github.com/nieltg/prom-example-testutil/test/mock_expfmt"
	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
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

func Example_printMetrics_nil() {
	_ = printMetrics(nil)
	// Output:
}

func Example_printMetrics_multiple() {
	counter1 := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "metric1",
		Help: "metric1 help.",
	})
	counter1.Inc()
	counter2 := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "metric2",
		Help: "metric2 help.",
	})
	counter2.Inc()

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(counter1)
	registry.MustRegister(counter2)

	metrics, _ := registry.Gather()
	_ = printMetrics(metrics)
	// Output:
	// # HELP metric1 metric1 help.
	// # TYPE metric1 counter
	// metric1 1
	// # HELP metric2 metric2 help.
	// # TYPE metric2 counter
	// metric2 1
}

func mockNewEncoder(f func(w io.Writer, format expfmt.Format) expfmt.Encoder) func() {
	originalNewEncoder := newEncoder
	newEncoder = f

	return func() {
		newEncoder = originalNewEncoder
	}
}

func Test_printMetrics(t *testing.T) {
	metrics := []*prommodel.MetricFamily{&prommodel.MetricFamily{}}

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockEncoder := mockexpfmt.NewMockEncoder(controller)
	mockEncoder.EXPECT().Encode(metrics[0]).Return(nil).Times(1)
	mockNewEncoder(func(w io.Writer, format expfmt.Format) expfmt.Encoder {
		return mockEncoder
	})

	_ = printMetrics(metrics)
}

func Test_printMetrics_error(t *testing.T) {
	expecterErr := fmt.Errorf("sample error")

	mockEncoder := mockexpfmt.NewMockEncoder(gomock.NewController(t))
	mockEncoder.EXPECT().Encode(gomock.Any()).Return(expecterErr)
	mockNewEncoder(func(w io.Writer, format expfmt.Format) expfmt.Encoder {
		return mockEncoder
	})

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "metric1",
		Help: "metric1 help.",
	})
	counter.Inc()

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(counter)

	metrics, _ := registry.Gather()
	assert.EqualError(t, printMetrics(metrics), expecterErr.Error())
}
