package metrics

import (
	"testing"

	"github.com/devopsfaith/krakend-gologging"
	"github.com/devopsfaith/krakend/config"
)

func TestConfigGetter_ok(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName": "test",
			"license": "123456",
			"rate":    75,
		},
	}

	res, err := ConfigGetter(cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if res.InstrumentationRate != 75 {
		t.Errorf("unexpected rate. have: %d, want: 75", res.InstrumentationRate)
	}
}

func TestConfigGetter_okIgnoreDebugNotBool(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName":      "test",
			"debugEnabled": "true",
			"license":      "123456",
		},
	}

	_, err := ConfigGetter(cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestConfigGetter_okDebugIsBool(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName":      "test",
			"license":      "123456",
			"debugEnabled": true,
		},
	}

	_, err := ConfigGetter(cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestConfigGetter_okWithIgnoreStatusCodes(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName": "test",
			"license": "123456",
			"errorCollector": map[string]interface{}{
				"ignoreStatusCodes": []int{
					400,
					401,
					402,
				},
			},
		},
	}

	conf, err := ConfigGetter(cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	for _, v := range conf.ErrorCollector.IgnoreStatusCodes {
		if v != 400 && v != 401 && v != 402 {
			t.Errorf("unexpected value in conf.ErrorCollector.IgnoreStatusCodes: %d", v)
			return
		}
	}
}

func TestConfigGetter_koWrongNamespace(t *testing.T) {
	cfg := config.ExtraConfig{
		"WrongNamespace": map[string]interface{}{
			"appName": "test",
			"license": "123456",
		},
	}

	_, err := ConfigGetter(cfg)
	if err == nil {
		t.Error("it should have errored")
	}
}

func TestConfigGetter_koWrongConfigType(t *testing.T) {

	cfg := config.ExtraConfig{
		Namespace: map[int]interface{}{
			123: "test",
		},
	}

	_, err := ConfigGetter(cfg)
	if err == nil {
		t.Error("it should have errored")
	}
}

func TestConfigGetter_koNoLicenseKey(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName": "test",
		},
	}

	_, err := ConfigGetter(cfg)
	if err == nil {
		t.Error("it should have errored")
	}
}

func TestConfigGetter_koLicenseNotString(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName": "test",
			"license": 123456,
		},
	}

	_, err := ConfigGetter(cfg)
	if err == nil {
		t.Error("it should have errored")
	}
}

func TestConfigGetter_koNoAppNameKey(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"license": "123456",
		},
	}

	_, err := ConfigGetter(cfg)
	if err == nil {
		t.Error("it should have errored")
	}
}

func TestConfigGetter_koAppNameNotString(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName": 11,
			"license": "123456",
		},
	}

	_, err := ConfigGetter(cfg)
	if err == nil {
		t.Error("it should have errored")
	}
}

func TestRegister_ok(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName":      "test",
			"license":      "1234567890123456789012345678901234567890",
			"debugEnabled": true,
		},
	}
	registerNR(t, cfg)
	if app == nil {
		t.Error("it should have errored")
	}
}
func TestRegister_koWrongConfig(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{},
	}
	registerNR(t, cfg)
	if app != nil {
		t.Errorf("app should be nil, instead it has the value %v", app)
	}
}

func TestRegister_koUnableToStartNR(t *testing.T) {
	cfg := config.ExtraConfig{
		Namespace: map[string]interface{}{
			"appName":      "test",
			"license":      "12345",
			"debugEnabled": true,
		},
	}
	registerNR(t, cfg)
	if app != nil {
		t.Errorf("app should be nil, instead it has the value %v", app)
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
		t.Errorf("unexpected error: %s", err.Error())
		return
	}
	Register(cfg, logger)
}
