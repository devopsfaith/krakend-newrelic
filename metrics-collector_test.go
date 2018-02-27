package metrics

import (
	"testing"

	"github.com/devopsfaith/krakend-gologging"
	"github.com/devopsfaith/krakend/config"
)

func TestConfigGetter_ok(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
			"license":  "123456",
		},
	}

	conf := ConfigGetter(cfg)
	if conf == nil {
		t.Error("conf shouldn't be nil")
	}
}

func TestConfigGetter_okIgnoreDebugNotBool(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
			"license":  "123456",
			"debug":    "true",
		},
	}

	conf := ConfigGetter(cfg)
	if conf == nil {
		t.Error("conf shouldn't be nil")
	}
}

func TestConfigGetter_okDebugIsBool(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
			"license":  "123456",
			"debug":    true,
		},
	}

	conf := ConfigGetter(cfg)
	if conf == nil {
		t.Error("conf shouldn't be nil")
	}
}

func TestConfigGetter_koWrongNamespace(t *testing.T) {
	cfg := config.ExtraConfig{
		"WrongNamespace": map[string]interface{}{
			"app_name": "test",
			"license":  "123456",
		},
	}

	conf := ConfigGetter(cfg)
	if conf != nil {
		t.Errorf("conf should be nil, %v", conf)
	}
}

func TestConfigGetter_koWrongConfigType(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[int]interface{}{
			123: "test",
		},
	}

	conf := ConfigGetter(cfg)
	if conf != nil {
		t.Errorf("conf should be nil, %v", conf)
	}
}

func TestConfigGetter_koNoLicenseKey(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
		},
	}

	conf := ConfigGetter(cfg)
	if conf != nil {
		t.Errorf("conf should be nil, %v", conf)
	}
}

func TestConfigGetter_koLicenseNotString(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
			"license":  123456,
		},
	}

	conf := ConfigGetter(cfg)
	if conf != nil {
		t.Errorf("conf should be nil, %v", conf)
	}
}

func TestConfigGetter_koNoAppNameKey(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"license": "123456",
		},
	}

	conf := ConfigGetter(cfg)
	if conf != nil {
		t.Errorf("conf should be nil, %v", conf)
	}
}

func TestConfigGetter_koAppNameNotString(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": 11,
			"license":  "123456",
		},
	}

	conf := ConfigGetter(cfg)
	if conf != nil {
		t.Errorf("conf should be nil, %v", conf)
	}
}

func TestRegister_ok(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
			"license":  "1234567890123456789012345678901234567890",
			"debug":    true,
		},
	}
	registerNR(t, cfg)
	if app == nil {
		t.Error("app shouldn't be nil")
	}
}
func TestRegister_koWrongConfig(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{},
	}
	registerNR(t, cfg)
	if app != nil {
		t.Errorf("app should be nil, %v", app)
	}
}

func TestRegister_koUnableToStartNR(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"app_name": "test",
			"license":  "12345",
			"debug":    true,
		},
	}
	registerNR(t, cfg)
	if app != nil {
		t.Errorf("app should be nil, %v", app)
	}
}

func registerNR(t *testing.T, cfg config.ExtraConfig) {
	app = nil
	logger, err := gologging.NewLogger(config.ExtraConfig{
		gologging.Namespace: map[string]interface{}{
			"level":  "DEBUG",
			"stdout": true,
		},
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	Register(cfg, logger)
}
