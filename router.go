package metrics

import (
	"fmt"
	"math/rand"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

var errNoApp = fmt.Errorf("No NewRelic app defined")

// Middleware adds NewRelic middleware
func Middleware() (gin.HandlerFunc, error) {
	if app == nil {
		return emptyMW, errNoApp
	}

	if app.Config.InstrumentationRate == 0 {
		return emptyMW, nil
	}

	nrMiddleware := nrgin.Middleware(app)

	if app.Config.InstrumentationRate == 100 {
		return nrMiddleware, nil
	}

	rate := float64(app.Config.InstrumentationRate) / 100.0

	next := make(chan float64, 1000)
	go func(out chan<- float64) {
		for {
			out <- rand.Float64()
		}
	}(next)

	return func(c *gin.Context) {
		if n := <-next; n <= rate {
			nrMiddleware(c)
			return
		}
		emptyMW(c)
	}, nil
}

// HandlerFactory includes NewRelic transaction specific configuration endpoint naming
func HandlerFactory(handlerFactory krakendgin.HandlerFactory) krakendgin.HandlerFactory {
	if app == nil {
		return handlerFactory
	}
	return func(conf *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handler := handlerFactory(conf, p)
		return func(c *gin.Context) {
			if txn := nrgin.Transaction(c); txn != nil {
				txn.SetName(conf.Endpoint)
			}
			handler(c)
		}
	}
}

func emptyMW(c *gin.Context) {
	c.Next()
}
