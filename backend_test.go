package metrics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	newrelic "github.com/newrelic/go-agent"
)

func TestBackendFactory_okAppNil(t *testing.T) {
	app = nil
	cfg := &config.Backend{
		URLPattern: "/my_endpoint",
		Host: []string{
			"localhost:8080",
		},
		Timeout: time.Second,
	}

	expectedError := errors.New("expected error")
	bf := BackendFactory("segm", func(_ *config.Backend) proxy.Proxy {
		return func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) { return nil, expectedError }
	})

	if resp, err := bf(cfg)(context.Background(), nil); resp != nil || err != expectedError {
		t.Errorf("unexpected response: resp = %v, error = %v", resp, err)
	}
}

func TestBackendFactory_okNRApp(t *testing.T) {
	cfg := &config.Backend{
		URLPattern: "/my_endpoint",
		Host: []string{
			"localhost:8080",
		},
		Timeout: time.Second,
	}
	nrApp := newApp()
	defer func() { app = nil }()
	app = nrApp

	totalCalls := 0
	expectedError := errors.New("expected error")

	bf := BackendFactory("segm", func(_ *config.Backend) proxy.Proxy {
		return func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) { return nil, expectedError }
	})

	txn := newTx()
	txn.startSegmentNow = func() newrelic.SegmentStartTime {
		totalCalls++
		return newrelic.SegmentStartTime{}
	}

	if resp, err := bf(cfg)(context.WithValue(context.Background(), nrCtxKey, txn), nil); resp != nil || err != expectedError {
		t.Errorf("unexpected response: resp = %v, error = %v", resp, err)
	}

	if resp, err := bf(cfg)(context.Background(), nil); resp != nil || err != expectedError {
		t.Errorf("unexpected response: resp = %v, error = %v", resp, err)
	}

	if totalCalls != 1 {
		t.Errorf("unexpected number of calls to the txn end. have: %d, want: 1", totalCalls)
	}
}
