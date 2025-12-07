package database

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

func Test_Metrics(t *testing.T) {
	t.Parallel()

	conn, err := sqlx.Open("mysql", "toll:toll@tcp(db:3306)/toll")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	metrics := NewStatsMetrics("test-metrics", conn)

	// Test Describe().
	descChan := make(chan<- *prometheus.Desc)

	go func() {
		metrics.Describe(descChan)
		close(descChan)
	}()

	// Test Collect().
	metricChan := make(chan<- prometheus.Metric)

	go func() {
		metrics.Collect(metricChan)
		close(metricChan)
	}()
}
