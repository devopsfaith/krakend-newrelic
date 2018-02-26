package metrics

import (
	"testing"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/config"
	"time"
	"github.com/devopsfaith/krakend/encoding"
)

func TestBackendFactory_okAppNil(t *testing.T) {
	app = nil
	cfg := &config.Backend{
			URLPattern: "/",
			Host: []string{
				"localhost:8080",
			},
			Timeout: time.Second,
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				"github.com/devopsfaith/krakend-martian": map[string]interface{}{
					"header.ToBody": struct{}{},
				},
			},
		}

	bf := BackendFactory("segm", func(cf *config.Backend) proxy.Proxy {
		return  proxy.NewRoundRobinLoadBalancedMiddleware(cf)()
	})

	println(bf(cfg))
}
