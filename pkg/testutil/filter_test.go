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

func TestFilterMetricsByName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	metricsInName := "metricsIn"
	metricsIn := []*prommodel.MetricFamily{
		&prommodel.MetricFamily{Name: &metricsInName},
	}
	metricsOutName := "metricsOut"
	metricsOut := []*prommodel.MetricFamily{
		&prommodel.MetricFamily{Name: &metricsOutName},
	}

	filterer := mocktestutil.NewMockfilterer(controller)
	filterer.EXPECT().FilterMetricsByName(metricsIn, "name").Return(metricsOut)
	defer mockGlobalFilterer(filterer)()

	out := FilterMetricsByName(metricsIn, "name")
	t.Run("return", func(t *testing.T) {
		assert.Equal(t, metricsOut, out)
	})
}
