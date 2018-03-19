package metrics

import (
	"encoding/json"

	"fmt"
	"os"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/newrelic/go-agent"
)

// Namespace for krakend_newrelic
const Namespace = "github_com/letgoapp/krakend_newrelic"

var (
	app            newrelic.Application
	isDebugEnabled bool
)

// Config struct for NewRelic
type Config struct {
	License        string
	AppName        string
	IsDebugEnabled bool
}

// ConfigGetter gets config for NewRelic
func ConfigGetter(cfg config.ExtraConfig) (newrelic.Config, error) {
	v, ok := cfg[Namespace]
	if !ok {
		return newrelic.Config{}, fmt.Errorf("unknown Namespace %s", Namespace)
	}

	tmp, ok := v.(map[string]interface{})
	if !ok {
		return newrelic.Config{}, fmt.Errorf("Cannot map config to map string interface")
	}

	// check whether compulsory fields are present
	if _, ok := tmp["license"]; !ok {
		return newrelic.Config{}, fmt.Errorf("Config should have the field license defined")
	}

	if _, ok = tmp["appName"]; !ok {
		return newrelic.Config{}, fmt.Errorf("Config should have the field appName defined")
	}

	// check whether debug enabled
	var val interface{}
	if val, ok = tmp["debugEnabled"]; !ok {
		isDebugEnabled = false
	} else {
		valB, ok := val.(bool)
		if !ok || !valB {
			isDebugEnabled = false
		} else {
			isDebugEnabled = true
		}
	}

	marshaledConf, err := json.Marshal(tmp)
	if err != nil {
		return newrelic.Config{}, err
	}

	var unmarshaledConf newrelic.Config
	err = json.Unmarshal(marshaledConf, &unmarshaledConf)

	return unmarshaledConf, err
}

// Register registers the NewRelic app
func Register(cfg config.ExtraConfig, logger logging.Logger) {
	conf, err := ConfigGetter(cfg)
	if err != nil {
		logger.Debug("no config for the NR module")
		return
	}

	nrCfg := newrelic.NewConfig(conf.AppName, conf.License)

	if isDebugEnabled {
		nrCfg.Logger = newrelic.NewDebugLogger(os.Stdout)
	}

	nrApp, err := newrelic.NewApplication(nrCfg)
	if err != nil {
		logger.Debug("unable to start the NR module:", err.Error())
		return
	}

	app = nrApp
}
