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

func mockCollector(collector collector) func() {
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

	defer mockCollector(collector)()
	gatherer := MustCollect(collectCounterA)
	t.Run("return", func(t *testing.T) {
		assert.Equal(t, collectGathererA, gatherer)
	})
}
