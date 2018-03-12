package metrics

import (
	"context"
	"net/http"

	"github.com/devopsfaith/krakend/proxy"
	"github.com/newrelic/go-agent"
)

// HTTPClientFactory includes a http.RoundTripper for NewRelic instrumentation
func HTTPClientFactory(cf proxy.HTTPClientFactory) proxy.HTTPClientFactory {
	return func(ctx context.Context) *http.Client {
		client := cf(ctx)

		if tx, ok := ctx.Value(nrCtxKey).(newrelic.Transaction); ok {
			client.Transport = newrelic.NewRoundTripper(tx, client.Transport)
		}

		return client
	}
}
