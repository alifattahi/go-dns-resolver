package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Exported counters
	DNSCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_cache_hits_total",
			Help: "Total number of DNS cache hits",
		},
	)

	DNSCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_cache_misses_total",
			Help: "Total number of DNS cache misses",
		},
	)
	FailedDNSResolver = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_failed_resolve_total",
			Help: "Total number of DNS Resolves Failed",
		},
	)
)

// Register function to register all metrics with Prometheus
func Register() {
	prometheus.MustRegister(DNSCacheHits)
	prometheus.MustRegister(DNSCacheMisses)
	prometheus.MustRegister(FailedDNSResolver)
}
