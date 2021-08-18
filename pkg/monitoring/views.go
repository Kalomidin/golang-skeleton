package monitoring

import (
	"context"
	"sync"

	"github.com/ridebeam/go-common/monitor"
)

const (
	MetricSomeMetric = "example-metric"
)

var (
	onceServiceViews sync.Once
)

func RegisterServiceViews(ctx context.Context) {
	onceServiceViews.Do(func() {
		monitor.Measurement(ctx).DescribeMeter(monitor.MeterCounter, MetricSomeMetric, "some interesting text",
			nil,
			nil,
			nil,
		)
	})
}
