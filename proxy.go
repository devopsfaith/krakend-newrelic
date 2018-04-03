package metrics

import (
	"context"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/newrelic/go-agent"
)

const nrCtxKey = "newRelicTransaction"

// ProxyFactory creates an instrumented proxy factory
func ProxyFactory(segmentName string, next proxy.Factory) proxy.FactoryFunc {
	if app == nil {
		return next.New
	}
	return proxy.FactoryFunc(func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		next, err := next.New(cfg)
		if err != nil {
			return proxy.NoopProxy, err
		}
		return NewProxyMiddleware(segmentName)(next), nil
	})
}

// NewProxyMiddleware adds NewRelic segmentation
func NewProxyMiddleware(segmentName string) proxy.Middleware {
	if app == nil {
		return proxy.EmptyMiddleware
	}
	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}
		if len(next) == 0 {
			panic(proxy.ErrNotEnoughProxies)
		}
		return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
			tx, ok := ctx.Value(nrCtxKey).(newrelic.Transaction)
			if !ok {
				return next[0](ctx, req)
			}

			segment := newrelic.StartSegment(tx, segmentName)
			resp, err := next[0](ctx, req)
			segment.End()

			return resp, err
		}
	}
}
