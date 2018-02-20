package metrics

import (
	"context"
	"net/http"

	"github.com/devopsfaith/krakend/proxy"
	newrelic "github.com/newrelic/go-agent"
)

func HTTPClientFactory(cf proxy.HTTPClientFactory) proxy.HTTPClientFactory {
	return func(ctx context.Context) *http.Client {
		client := cf(ctx)

		if tx, ok := ctx.Value(nrCtxKey).(newrelic.Transaction); ok {
			client.Transport = newrelic.NewRoundTripper(tx, client.Transport)
		}

		return client
	}
}
