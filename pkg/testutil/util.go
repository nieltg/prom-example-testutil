package testutil

import "github.com/prometheus/client_golang/prometheus"

// CollectAndPrint collects metrics, filters by names, and prints.
func CollectAndPrint(collector prometheus.Collector, names ...string) {
	GatherAndPrint(globalCollector.MustCollect(collector), names...)
}

// GatherAndPrint gathers metrics, filters by names, and prints.
func GatherAndPrint(gatherer prometheus.Gatherer, names ...string) {
	metrics, _ := gatherer.Gather()
	filteredMetrics := globalFilterer.FilterMetricsByName(metrics, names...)
	globalPrinter.PrintMetrics(filteredMetrics)
}
