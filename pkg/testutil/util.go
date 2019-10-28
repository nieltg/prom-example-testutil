package testutil

import "github.com/prometheus/client_golang/prometheus"

// CollectAndPrint collects metrics, filters by names, and prints.
func CollectAndPrint(collector prometheus.Collector, names ...string) {
	metrics, _ := globalCollector.MustCollect(collector).Gather()
	filteredMetrics := globalFilterer.FilterMetricsByName(metrics, names...)
	globalPrinter.PrintMetrics(filteredMetrics)
}
