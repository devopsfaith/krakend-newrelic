package metrics

import (
	"github.com/devopsfaith/krakend/config"
	newrelic "github.com/newrelic/go-agent"
)

var app *newrelic.Application

type Config struct {
	License string
	AppName string
}

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

	if vs, ok = tmp["appName"]; !ok {
		return nil
	}
	conf.AppName, ok = vs.(string)
	if !ok {
		return nil
	}

	return conf
}

func Register(cfg config.ExtraConfig) {
	conf, ok := ConfigGetter(cfg).(Config)
	if !ok {
		return
	}

	nrCfg := newrelic.NewConfig(conf.AppName, conf.License)
	nrApp, err := newrelic.NewApplication(nrCfg)
	if err != nil {
		return
	}

	app = &nrApp
}
