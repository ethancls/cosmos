package controller

import (
	"github.com/ethancls/kyle-server/server/telemetry"
)

type metrics struct {
	*telemetry.UpdateChannelMetrics
}

func newMetrics(updateChannelMetrics *telemetry.UpdateChannelMetrics) (*metrics, error) {
	return &metrics{
		updateChannelMetrics,
	}, nil
}
