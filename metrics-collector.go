package metrics

import (
	"os"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/newrelic/go-agent"
)

const Namespace = "github_com/letgoapp/krakend_newrelic"

var app newrelic.Application

// Config struct for NewRelic
type Config struct {
	License        string
	AppName        string
	IsDebugEnabled bool
}

// ConfigGetter gets config for NewRelic
func ConfigGetter(cfg config.ExtraConfig) interface{} {
	v, ok := cfg[Namespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	conf := Config{}
	var vs interface{}

	if vs, ok = tmp["license"]; !ok {
		return nil
	}
	conf.License, ok = vs.(string)
	if !ok {
		return nil
	}

	if vs, ok = tmp["app_name"]; !ok {
		return nil
	}
	conf.AppName, ok = vs.(string)
	if !ok {
		return nil
	}

	v, ok = tmp["debug"]
	if ok {
		if v, ok = v.(bool); ok {
			conf.IsDebugEnabled = true
		}
	}

	return conf
}

// Register registers the NewRelic app
func Register(cfg config.ExtraConfig, logger logging.Logger) {
	conf, ok := ConfigGetter(cfg).(Config)
	if !ok {
		logger.Debug("no config for the NR module")
		return
	}

	nrCfg := newrelic.NewConfig(conf.AppName, conf.License)
	if conf.IsDebugEnabled {
		nrCfg.Logger = newrelic.NewDebugLogger(os.Stdout)
	}
	nrApp, err := newrelic.NewApplication(nrCfg)
	if err != nil {
		logger.Debug("unable to start the NR module:", err.Error())
		return
	}

	app = nrApp
}
