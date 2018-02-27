package metrics

import (
	"strings"
	"testing"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
)

func TestMiddleware_ok(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
			"license":  "1234567890123456789012345678901234567890",
			"debug":    true,
		},
	}
	registerNR(t, cfg)
	_, err := Middleware()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestMiddleware_koNoApp(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{},
	}
	registerNR(t, cfg)
	_, err := Middleware()
	if !strings.Contains(err.Error(), errNoApp.Error()) {
		t.Error("Should have given errNoApp error")
	}
}

func TestHandlerFactory_okAppNil(t *testing.T) {
	defaultHF := false
	app = nil
	HandlerFactory(defaultHandlerFunc(&defaultHF))
	if !defaultHF {
		t.Error("defaultHF should have been true")
	}
}

func defaultHandlerFunc(defaultHF *bool) krakendgin.HandlerFactory {
	*defaultHF = true
	return func(_ *config.EndpointConfig, _ proxy.Proxy) gin.HandlerFunc {
		return func(c *gin.Context) {
		}
	}
}
