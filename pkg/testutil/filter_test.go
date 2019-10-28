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

var filterMetricsNameA = "name-a"
var filterMetricsNameB = "name-b"
var filterMetricsA = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &filterMetricsNameA},
}
var filterMetricsB = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &filterMetricsNameA},
	&prommodel.MetricFamily{Name: &filterMetricsNameB},
}

func TestFilterMetricsByName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	filterer := mocktestutil.NewMockfilterer(controller)
	filterer.EXPECT().FilterMetricsByName(filterMetricsA, "name").Return(
		filterMetricsB)

	defer mockGlobalFilterer(filterer)()
	out := FilterMetricsByName(filterMetricsA, "name")
	t.Run("return", func(t *testing.T) {
		assert.Equal(t, filterMetricsB, out)
	})
}

func Test_filtererImpl_FilterMetricsByName(t *testing.T) {
	out := filtererImpl{}.FilterMetricsByName(filterMetricsA, filterMetricsNameA)
	assert.Equal(t, filterMetricsA, out)
}

func Test_filtererImpl_FilterMetricsByName_reject(t *testing.T) {
	out := filtererImpl{}.FilterMetricsByName(filterMetricsA, "different")
	assert.Nil(t, out)
}

func Test_filtererImpl_FilterMetricsByName_rejectMany(t *testing.T) {
	out := filtererImpl{}.FilterMetricsByName(filterMetricsB, filterMetricsNameA)
	assert.Equal(t, filterMetricsA, out)
}

func Test_filtererImpl_FilterMetricsByName_multipleNames(t *testing.T) {
	out := filtererImpl{}.FilterMetricsByName(
		filterMetricsB, filterMetricsNameA, filterMetricsNameB)
	assert.Equal(t, filterMetricsB, out)
}
