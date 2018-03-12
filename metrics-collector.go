package metrics

import (
	"os"

	"net/http"
	"time"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/newrelic/go-agent"
)

// Namespace for krakend_newrelic
const Namespace = "github_com/letgoapp/krakend_newrelic"

var app newrelic.Application

// Config struct for NewRelic
type Config struct {
	License              string
	AppName              string
	IsDebugEnabled       bool
	Logger               newrelic.Logger
	Enabled              bool
	Labels               map[string]string
	HighSecurity         bool
	CustomInsightsEvents struct {
		Enabled bool
	}

	TransactionEvents struct {
		Enabled    bool
		Attributes newrelic.AttributeDestinationConfig
	}

	ErrorCollector struct {
		Enabled           bool
		CaptureEvents     bool
		IgnoreStatusCodes []int
		Attributes        newrelic.AttributeDestinationConfig
	}

	TransactionTracer struct {
		Enabled   bool
		Threshold struct {
			IsApdexFailing bool
			Duration       time.Duration
		}
		SegmentThreshold    time.Duration
		StackTraceThreshold time.Duration
		Attributes          newrelic.AttributeDestinationConfig
	}

	HostDisplayName string
	UseTLS          bool
	Transport       http.RoundTripper
	Utilization     struct {
		DetectAWS         bool
		DetectAzure       bool
		DetectPCF         bool
		DetectGCP         bool
		DetectDocker      bool
		LogicalProcessors int
		TotalRAMMIB       int
		BillingHostname   string
	}

	CrossApplicationTracer struct {
		Enabled bool
	}

	DatastoreTracer struct {
		InstanceReporting struct {
			Enabled bool
		}
		DatabaseNameReporting struct {
			Enabled bool
		}
		QueryParameters struct {
			Enabled bool
		}
		SlowQuery struct {
			Enabled   bool
			Threshold time.Duration
		}
	}

	Attributes     newrelic.AttributeDestinationConfig
	RuntimeSampler struct {
		Enabled bool
	}
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

func extraParamsSetter(extraCfg config.ExtraConfig, conf *newrelic.Config) {
	v, ok := extraCfg[Namespace]
	if !ok {
		return
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return
	}
	var vs interface{}
	vs, ok = tmp["utilization"]
	if ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		z, ok := w["detect_aws"]
		if z, ok = z.(bool); ok {
			conf.Utilization.DetectAWS = z.(bool)
		}

		z, ok = w["detect_azure"]
		if z, ok = z.(bool); ok {
			conf.Utilization.DetectAzure = z.(bool)
		}

		z, ok = w["detect_pcf"]
		if z, ok = z.(bool); ok {
			conf.Utilization.DetectPCF = z.(bool)
		}

		z, ok = w["detect_gcp"]
		if z, ok = z.(bool); ok {
			conf.Utilization.DetectGCP = z.(bool)
		}

		z, ok = w["detect_docker"]
		if z, ok = z.(bool); ok {
			conf.Utilization.DetectDocker = z.(bool)
		}
	}
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

	extraParamsSetter(cfg, &nrCfg)

	nrApp, err := newrelic.NewApplication(nrCfg)
	if err != nil {
		logger.Debug("unable to start the NR module:", err.Error())
		return
	}

	app = nrApp
}
