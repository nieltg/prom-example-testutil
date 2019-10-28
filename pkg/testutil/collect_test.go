package testutil

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocktestutil "github.com/nieltg/prom-example-testutil/test/mock_testutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

var collectCounterA = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "some_total",
	Help: "A value that represents a counter.",
})
var collectGathererA prometheus.Gatherer = prometheus.NewPedanticRegistry()

func mockGlobalCollector(collector collector) func() {
	originalCollector := globalCollector
	globalCollector = collector

	return func() {
		globalCollector = originalCollector
	}
}

func TestMustCollect(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	collector := mocktestutil.NewMockcollector(controller)
	collector.EXPECT().MustCollect(collectCounterA).Return(collectGathererA)

	defer mockGlobalCollector(collector)()
	gatherer := MustCollect(collectCounterA)
	t.Run("return", func(t *testing.T) {
		assert.Equal(t, collectGathererA, gatherer)
	})
}

func newCollectorWithRegistererGatherer(r registererGatherer) collector {
	return &collectorImpl{
		newRegistryFunc: func() registererGatherer {
			return r
		},
	}
}

func Test_collectorImpl_MustCollect(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	registererGatherer := mocktestutil.NewMockregistererGatherer(controller)
	registererGatherer.EXPECT().MustRegister(collectCounterA)

	collector := newCollectorWithRegistererGatherer(registererGatherer)
	gatherer := collector.MustCollect(collectCounterA)
	t.Run("return", func(t *testing.T) {
		assert.Equal(t, registererGatherer, gatherer)
	})
}
