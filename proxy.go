package metrics

import (
	"context"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	newrelic "github.com/newrelic/go-agent"
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

func NewProxyMiddleware(segmentName string) proxy.Middleware {
	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}
		return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
			tx, ok := ctx.Value(nrCtxKey).(newrelic.Transaction)
			if ok {
				defer newrelic.StartSegment(tx, segmentName).End()
			}
			return next[0](ctx, req)
		}
	}
}
