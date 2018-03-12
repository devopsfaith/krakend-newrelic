package metrics

import (
	"os"

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

	if v, ok = tmp["debug"]; ok {
		if v, ok = v.(bool); ok {
			conf.IsDebugEnabled = v.(bool)
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
	if vs, ok = tmp["utilization"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["detect_aws"]; ok {
			if z, ok = z.(bool); ok {
				conf.Utilization.DetectAWS = z.(bool)
			}
		}

		if z, ok := w["detect_azure"]; ok {
			if z, ok = z.(bool); ok {
				conf.Utilization.DetectAzure = z.(bool)
			}
		}

		if z, ok := w["detect_pcf"]; ok {
			if z, ok = z.(bool); ok {
				conf.Utilization.DetectPCF = z.(bool)
			}
		}

		if z, ok := w["detect_gcp"]; ok {
			if z, ok = z.(bool); ok {
				conf.Utilization.DetectGCP = z.(bool)
			}
		}

		if z, ok := w["detect_docker"]; ok {
			if z, ok = z.(bool); ok {
				conf.Utilization.DetectDocker = z.(bool)
			}
		}
	}

	if vs, ok = tmp["enabled"]; ok {
		if z, ok := vs.(bool); ok {
			conf.Enabled = z
		}
	}

	if vs, ok = tmp["custom_insights_events"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["enabled"]; ok {
			if z, ok = z.(bool); ok {
				conf.CustomInsightsEvents.Enabled = z.(bool)
			}
		}
	}

	if vs, ok = tmp["transaction_events"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["enabled"]; ok {
			if z, ok = z.(bool); ok {
				conf.TransactionEvents.Enabled = z.(bool)
			}
		}
		if z, ok := w["attributes"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["enabled"]; ok {
				if x, ok = x.(bool); ok {
					conf.TransactionEvents.Attributes.Enabled = x.(bool)
				}
			}
		}
	}

	if vs, ok = tmp["high_security"]; ok {
		if z, ok := vs.(bool); ok {
			conf.HighSecurity = z
		}
	}

	if vs, ok = tmp["use_tls"]; ok {
		if z, ok := vs.(bool); ok {
			conf.UseTLS = z
		}
	}

	if vs, ok = tmp["error_collector"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["enabled"]; ok {
			if z, ok = z.(bool); ok {
				conf.ErrorCollector.Enabled = z.(bool)
			}
		}
		if z, ok := w["capture_events"]; ok {
			if z, ok = z.(bool); ok {
				conf.ErrorCollector.CaptureEvents = z.(bool)
			}
		}
		if z, ok := w["attributes"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["enabled"]; ok {
				if x, ok = x.(bool); ok {
					conf.ErrorCollector.Attributes.Enabled = x.(bool)
				}
			}
		}
		if z, ok := w["ignore_status_codes"]; ok {
			y, ok := z.([]interface{})
			if !ok {
				return
			}
			conf.ErrorCollector.IgnoreStatusCodes = make([]int, len(y))
			for i, v := range y {
				conf.ErrorCollector.IgnoreStatusCodes[i] = int(v.(float64))
			}
		}
	}

	if vs, ok = tmp["attributes"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["enabled"]; ok {
			if z, ok = z.(bool); ok {
				conf.Attributes.Enabled = z.(bool)
			}
		}
	}

	if vs, ok = tmp["runtime_sampler"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["enabled"]; ok {
			if z, ok = z.(bool); ok {
				conf.RuntimeSampler.Enabled = z.(bool)
			}
		}
	}

	if vs, ok = tmp["transaction_tracer"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["enabled"]; ok {
			if z, ok = z.(bool); ok {
				conf.TransactionTracer.Enabled = z.(bool)
			}
		}
		if z, ok := w["threshold"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["is_apdex_failing"]; ok {
				if x, ok = x.(bool); ok {
					conf.TransactionTracer.Threshold.IsApdexFailing = x.(bool)
				}
			}
			if x, ok := y["duration"]; ok {
				if x, ok = x.(time.Duration); ok {
					conf.TransactionTracer.Threshold.Duration = x.(time.Duration)
				}
			}
		}
		if z, ok := w["segment_threshold"]; ok {
			if z, ok = z.(time.Duration); ok {
				conf.TransactionTracer.SegmentThreshold = z.(time.Duration)
			}
		}
		if z, ok := w["stack_trace_threshold"]; ok {
			if z, ok = z.(time.Duration); ok {
				conf.TransactionTracer.StackTraceThreshold = z.(time.Duration)
			}
		}
		if z, ok := w["attributes"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["enabled"]; ok {
				if x, ok = x.(bool); ok {
					conf.TransactionTracer.Attributes.Enabled = x.(bool)
				}
			}
		}
	}

	if vs, ok = tmp["cross_application_tracer"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["enabled"]; ok {
			if z, ok = z.(bool); ok {
				conf.CrossApplicationTracer.Enabled = z.(bool)
			}
		}
	}

	if vs, ok := tmp["datastore_tracer"]; ok {
		w, ok := vs.(map[string]interface{})
		if !ok {
			return
		}
		if z, ok := w["instance_reporting"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["enabled"]; ok {
				if x, ok = x.(bool); ok {
					conf.DatastoreTracer.InstanceReporting.Enabled = x.(bool)
				}
			}
		}
		if z, ok := w["database_name_reporting"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["enabled"]; ok {
				if x, ok = x.(bool); ok {
					conf.DatastoreTracer.DatabaseNameReporting.Enabled = x.(bool)
				}
			}
		}
		if z, ok := w["query_parameters"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["enabled"]; ok {
				if x, ok = x.(bool); ok {
					conf.DatastoreTracer.QueryParameters.Enabled = x.(bool)
				}
			}
		}
		if z, ok := w["slow_query"]; ok {
			y, ok := z.(map[string]interface{})
			if !ok {
				return
			}
			if x, ok := y["enabled"]; ok {
				if x, ok = x.(bool); ok {
					conf.DatastoreTracer.SlowQuery.Enabled = x.(bool)
				}
			}
			if x, ok := y["threshold"]; ok {
				if x, ok = x.(time.Duration); ok {
					conf.DatastoreTracer.SlowQuery.Threshold = x.(time.Duration)
				}
			}
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
