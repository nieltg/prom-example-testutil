package testutil

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocktestutil "github.com/nieltg/prom-example-testutil/test/mock_testutil"
	prommodel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func mockGlobalFilterer(f filterer) func() {
	originalGlobalFilterer := globalFilterer
	globalFilterer = f

	return func() {
		globalFilterer = originalGlobalFilterer
	}
}

var metricsNameA = "name-a"
var metricsNameB = "name-b"
var metricsA = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &metricsNameA},
}
var metricsB = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &metricsNameA},
	&prommodel.MetricFamily{Name: &metricsNameB},
}

func TestFilterMetricsByName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	filterer := mocktestutil.NewMockfilterer(controller)
	filterer.EXPECT().FilterMetricsByName(metricsA, "name").Return(metricsB)
	defer mockGlobalFilterer(filterer)()

	out := FilterMetricsByName(metricsA, "name")
	t.Run("return", func(t *testing.T) {
		assert.Equal(t, metricsB, out)
	})
}

func Test_filtererImpl_FilterMetricsByName(t *testing.T) {
	out := filtererImpl{}.FilterMetricsByName(metricsA, metricsNameA)
	assert.Equal(t, metricsA, out)
}

func Test_filtererImpl_FilterMetricsByName_reject(t *testing.T) {
	out := filtererImpl{}.FilterMetricsByName(metricsA, "different")
	assert.Nil(t, out)
}

func Test_filtererImpl_FilterMetricsByName_rejectMany(t *testing.T) {
	out := filtererImpl{}.FilterMetricsByName(metricsB, metricsNameA)
	assert.Equal(t, metricsA, out)
}
