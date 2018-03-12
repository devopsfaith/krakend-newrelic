package metrics

import (
	"context"
	"net/http"
	"testing"

	"github.com/devopsfaith/krakend/proxy"
	"github.com/newrelic/go-agent"
)

func TestHTTPClientFactory_ok(t *testing.T) {
	txn := newTx()
	txn.startSegmentNow = func() newrelic.SegmentStartTime {
		return newrelic.SegmentStartTime{}
	}

	client1 := HTTPClientFactory(proxy.NewHTTPClient)(context.Background())
	switch client1.Transport.(type) {
	case nil:
	default:
		t.Errorf("unexpected client type %t", client1.Transport)
	}

	client2 := HTTPClientFactory(proxy.NewHTTPClient)(context.WithValue(context.Background(), nrCtxKey, txn))
	switch client2.Transport.(type) {
	case http.RoundTripper:
	default:
		t.Errorf("unexpected client type %t", client2.Transport)
	}
}
