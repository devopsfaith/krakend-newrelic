package metrics

import (
	"context"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/newrelic/go-agent"
)

// BackendFactory creates an instrumented backend factory
func BackendFactory(segmentName string, next proxy.BackendFactory) proxy.BackendFactory {
	if app == nil {
		return next
	}
	return func(cfg *config.Backend) proxy.Proxy {
		return NewBackend(segmentName, next(cfg))
	}
}

// NewBackend includes NewRelic segmentation
func NewBackend(segmentName string, next proxy.Proxy) proxy.Proxy {
	if app == nil {
		return next
	}
	return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
		tx, ok := ctx.Value(nrCtxKey).(newrelic.Transaction)
		if !ok {
			return next(ctx, req)
		}

		segment := newrelic.StartSegment(tx, segmentName)
		resp, err := next(ctx, req)
		segment.End()

		return resp, err
	}
}
