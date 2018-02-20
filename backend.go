package metrics

import (
	"context"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	newrelic "github.com/newrelic/go-agent"
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

func NewBackend(segmentName string, next proxy.Proxy) proxy.Proxy {
	if app == nil {
		return next
	}
	return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
		tx, ok := ctx.Value(nrCtxKey).(newrelic.Transaction)
		if ok {
			defer newrelic.StartSegment(tx, segmentName).End()
		}
		return next(ctx, req)
	}
}
