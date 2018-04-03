package metrics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/newrelic/go-agent"
)

func TestProxyFactory_okAppNil(t *testing.T) {
	app = nil

	cfg := &config.EndpointConfig{
		Endpoint: "/my_endpoint",
		Timeout:  time.Second,
		Method:   "GET",
	}

	errorExpected := errors.New("expected error")
	pf := proxy.FactoryFunc(func(_ *config.EndpointConfig) (proxy.Proxy, error) {
		return proxy.NoopProxy, errorExpected
	})

	if _, err := ProxyFactory("segm", pf)(cfg); err != errorExpected {
		t.Errorf("unexpected error: %v", err)
	}

}

func TestProxyFactory_okNRApp(t *testing.T) {
	nrApp := newApp()
	defer func() { app = nil }()
	app = &Application{nrApp, Config{InstrumentationRate: 100}}

	cfg := &config.EndpointConfig{
		Endpoint: "/my_endpoint",
		Timeout:  time.Second,
		Method:   "GET",
	}

	expectedResponse := &proxy.Response{
		Data: map[string]interface{}{
			"key": "result",
		}}

	pf := proxy.FactoryFunc(func(_ *config.EndpointConfig) (proxy.Proxy, error) {
		return func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
			return expectedResponse, nil
		}, nil
	})

	totalCalls := 0
	txn := newTx()
	txn.startSegmentNow = func() newrelic.SegmentStartTime {
		totalCalls++
		return newrelic.SegmentStartTime{}
	}

	pr, err := ProxyFactory("segm", pf)(cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	if resp, err := pr(context.WithValue(context.Background(), nrCtxKey, txn), nil); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	} else if resp != expectedResponse {
		t.Errorf("unexpected response: %v", resp)
	}

	if totalCalls != 1 {
		t.Errorf("wrong number of segments, got: %d, wanted 1", totalCalls)
	}
}

func TestNewProxyMiddleware_okAppNil(t *testing.T) {
	app = nil

	expectedResponse := &proxy.Response{
		Data: map[string]interface{}{
			"key": "result",
		}}

	p := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return expectedResponse, nil
	}

	totalCalls := 0
	txn := newTx()
	txn.startSegmentNow = func() newrelic.SegmentStartTime {
		totalCalls++
		return newrelic.SegmentStartTime{}
	}

	pr := NewProxyMiddleware("segm")(p)

	if resp, err := pr(context.WithValue(context.Background(), nrCtxKey, txn), nil); err != nil {
		t.Errorf("unexpected error, %s", err.Error())
	} else if resp != expectedResponse {
		t.Errorf("unexpected response: %v", resp)
	}

	if totalCalls != 0 {
		t.Errorf("wrong number of segments, got: %d, wanted 0", totalCalls)
	}
}

func TestNewProxyMiddleware_koPanicTooManyProxies(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, too many proxies for this proxy middleware")
		}
	}()

	nrApp := newApp()
	defer func() { app = nil }()
	app = &Application{nrApp, Config{InstrumentationRate: 100}}

	expectedResponse := &proxy.Response{
		Data: map[string]interface{}{
			"key": "result",
		}}

	p := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return expectedResponse, nil
	}

	totalCalls := 0
	txn := newTx()
	txn.startSegmentNow = func() newrelic.SegmentStartTime {
		totalCalls++
		return newrelic.SegmentStartTime{}
	}

	NewProxyMiddleware("segm")(p, p)(context.WithValue(context.Background(), nrCtxKey, txn), nil)
}

func TestNewProxyMiddleware_koPanicNoProxies(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, not enough proxies for this endpoint")
		}
	}()

	nrApp := newApp()
	defer func() { app = nil }()
	app = &Application{nrApp, Config{InstrumentationRate: 100}}

	totalCalls := 0
	txn := newTx()
	txn.startSegmentNow = func() newrelic.SegmentStartTime {
		totalCalls++
		return newrelic.SegmentStartTime{}
	}

	NewProxyMiddleware("segm")()(context.WithValue(context.Background(), nrCtxKey, txn), nil)
}
