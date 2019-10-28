package testutil

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocktestutil "github.com/nieltg/prom-example-testutil/test/mock_testutil"
	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"
)

var utilCounterA = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "counterA",
	Help: "counterA help.",
})
var utilMetricsNameA = "name-a"
var utilMetricsNameB = "name-b"
var utilMetricsA = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &filterMetricsNameA},
}
var utilMetricsB = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &filterMetricsNameB},
}

func TestCollectAndPrint(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	registerer := mocktestutil.NewMockregistererGatherer(controller)
	collector := mocktestutil.NewMockcollector(controller)
	filterer := mocktestutil.NewMockfilterer(controller)
	printer := mocktestutil.NewMockprinter(controller)

	collector.EXPECT().MustCollect(utilCounterA).Return(registerer)
	registerer.EXPECT().Gather().Return(utilMetricsA, nil)
	filterer.EXPECT().FilterMetricsByName(utilMetricsA, "n0").Return(utilMetricsB)
	printer.EXPECT().PrintMetrics(utilMetricsB)

	defer mockGlobalCollector(collector)()
	defer mockGlobalFilterer(filterer)()
	defer mockGlobalPrinter(printer)()
	CollectAndPrint(utilCounterA, "n0")
}
