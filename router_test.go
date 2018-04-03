package metrics

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

func TestMiddleware_ok(t *testing.T) {
	totalCalls := 0
	nrApp := newApp()
	defer func() { app = nil }()
	nrApp.startTransaction = func(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
		totalCalls++
		return newTx()
	}

	app = &Application{nrApp, Config{InstrumentationRate: 100}}
	handler, err := Middleware()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/my_endpoint", handler, func(c *gin.Context) {
		if txn := nrgin.Transaction(c); txn == nil {
			t.Error("nil transaction")
			c.Status(999)
			return
		}
		c.Status(http.StatusTeapot)
	})
	req, _ := http.NewRequest("GET", "/my_endpoint", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusTeapot {
		t.Error("unexpected status code")
		return
	}

	if totalCalls != 1 {
		t.Errorf("unexpected number of calls to the txn generator. have: %d, wanted: 1", totalCalls)
	}
}

func TestMiddleware_koNoApp(t *testing.T) {
	app = nil
	if _, err := Middleware(); err != errNoApp {
		t.Error("Should have given errNoApp error")
	}
}

func TestHandlerFactory_okAppNil(t *testing.T) {
	app = nil
	cfg := &config.EndpointConfig{
		Endpoint: "/my_endpoint",
		Timeout:  time.Second,
		Method:   "GET",
	}

	handlerFunc := func(_ *config.EndpointConfig, _ proxy.Proxy) gin.HandlerFunc {
		return func(c *gin.Context) {
			if txn := nrgin.Transaction(c); txn != nil {
				c.AbortWithStatus(999)
				return
			}
			c.JSON(http.StatusTeapot, gin.H{"sample": "data"})
		}
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/my_endpoint", HandlerFactory(handlerFunc)(cfg, proxy.NoopProxy))
	req, _ := http.NewRequest("GET", "/my_endpoint", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusTeapot {
		t.Error("unexpected status code")
		return
	}

	buff := &bytes.Buffer{}
	buff.ReadFrom(w.Result().Body)
	w.Result().Body.Close()

	if buff.String() != `{"sample":"data"}` {
		t.Errorf("unexpected body: %s", buff.String())
	}
}

func TestHandlerFactory_okNRApp(t *testing.T) {
	nrApp := newApp()
	defer func() { app = nil }()
	totalCalls := 0
	nrApp.startTransaction = func(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
		totalCalls++
		return newTx()
	}
	app = &Application{nrApp, Config{InstrumentationRate: 100}}

	expectedErr := errors.New("expect me")
	expectedProxy := func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
		if req != nil {
			t.Error("unexpected request")
		}
		return nil, expectedErr
	}

	handler := func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		return func(c *gin.Context) {
			if txn := nrgin.Transaction(c); txn == nil {
				c.AbortWithStatus(999)
				return
			}
			res, err := p(c, nil)
			if res != nil {
				c.AbortWithStatus(998)
				return
			}
			if err != expectedErr {
				c.AbortWithStatus(997)
				return
			}
			c.JSON(http.StatusTeapot, gin.H{"sample": "data"})
		}
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	mw, err := Middleware()
	if err != nil {
		t.Error(err)
		return
	}
	router.GET("/my_endpoint", mw, HandlerFactory(handler)(&config.EndpointConfig{
		Endpoint: "endpointName",
	}, expectedProxy))
	req, _ := http.NewRequest("GET", "/my_endpoint", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusTeapot {
		t.Errorf("unexpected status code. wanted: %d, have: %d", http.StatusTeapot, w.Result().StatusCode)
		return
	}

	buff := &bytes.Buffer{}
	buff.ReadFrom(w.Result().Body)
	w.Result().Body.Close()

	if buff.String() != `{"sample":"data"}` {
		t.Error("unexpected body")
	}

	if totalCalls != 1 {
		t.Errorf("unexpected number of calls to the txn generator. have: %d, wanted: 1", totalCalls)
	}
}

type sampleApplication struct {
	startTransaction   func(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction
	recordCustomEvent  func(eventType string, params map[string]interface{}) error
	recordCustomMetric func(name string, value float64) error
	waitForConnection  func(timeout time.Duration) error
	shutdown           func(timeout time.Duration)
}

func (s sampleApplication) StartTransaction(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	return s.startTransaction(name, w, r)
}
func (s sampleApplication) RecordCustomEvent(eventType string, params map[string]interface{}) error {
	return s.recordCustomEvent(eventType, params)
}
func (s sampleApplication) RecordCustomMetric(name string, value float64) error {
	return s.recordCustomMetric(name, value)
}
func (s sampleApplication) WaitForConnection(timeout time.Duration) error {
	return s.waitForConnection(timeout)
}
func (s sampleApplication) Shutdown(timeout time.Duration) { s.shutdown(timeout) }

func newApp() sampleApplication {
	return sampleApplication{
		startTransaction:   func(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction { return nil },
		recordCustomEvent:  func(eventType string, params map[string]interface{}) error { return nil },
		recordCustomMetric: func(name string, value float64) error { return nil },
		waitForConnection:  func(timeout time.Duration) error { return nil },
		shutdown:           func(timeout time.Duration) {},
	}
}

type transaction struct {
	http.ResponseWriter
	end             func() error
	ignore          func() error
	setName         func(name string) error
	noticeError     func(err error) error
	addAttribute    func(key string, value interface{}) error
	startSegmentNow func() newrelic.SegmentStartTime
}

func (tx transaction) End() error {
	return tx.end()
}

func (tx transaction) Ignore() error {
	return tx.ignore()
}

func (tx transaction) SetName(name string) error {
	return tx.setName(name)
}

func (tx transaction) NoticeError(err error) error {
	return tx.noticeError(err)
}

func (tx transaction) AddAttribute(key string, value interface{}) error {
	return tx.addAttribute(key, value)
}

func (tx transaction) StartSegmentNow() newrelic.SegmentStartTime {
	return tx.startSegmentNow()
}

func newTx() transaction {
	return transaction{
		ResponseWriter:  httptest.NewRecorder(),
		end:             func() error { return nil },
		ignore:          func() error { return nil },
		setName:         func(name string) error { return nil },
		noticeError:     func(err error) error { return nil },
		addAttribute:    func(key string, value interface{}) error { return nil },
		startSegmentNow: func() newrelic.SegmentStartTime { return newrelic.SegmentStartTime{} },
	}
}
