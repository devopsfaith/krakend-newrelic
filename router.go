package metrics

import (
	"fmt"

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

	return nrgin.Middleware(app), nil
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
