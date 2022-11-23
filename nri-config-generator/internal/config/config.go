package config

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	cfgFormat                  = "yaml"
	exportersConfigPathSetting = "config_path"
	defaultExportersConfigPath = "/tmp/"
)

// setViperDefaults loads the default configuration into the given Viper registry.
func setViperDefaults(viper *viper.Viper) {
	viper.SetDefault("debug", false)
	viper.SetDefault("verbose", false)
	viper.SetDefault("audit", false)
}

// ArgumentList Available Arguments
type ArgumentList struct {
	ConfigPath   string `default:"" help:"Path to the config file"`
	ShowVersion  bool   `default:"false" help:"Print build information and exit"`
	ShortRunning bool   `default:"false" help:"By default execution is long running, but this can be override"`
}

func getConfig(c *ArgumentList) (map[string]interface{}, string, error) {
	configs := make(map[string]interface{})
	configPathFound := false
	exporterConfigPath := ""
	if c.ConfigPath != "" {
		cfg := viper.New()
		cfg.SetConfigType(cfgFormat)
		cfg.AddConfigPath(filepath.Dir(c.ConfigPath))
		cfg.SetConfigName(filepath.Base(c.ConfigPath))
		setViperDefaults(cfg)
		err := cfg.ReadInConfig()
		if err != nil {
			return nil, exporterConfigPath, errors.Wrap(err, "could not read configuration")
		}
		for _, cfgName := range cfg.AllKeys() {
			configs[cfgName] = cfg.Get(cfgName)
			if cfgName == exportersConfigPathSetting {
				configPathFound = true
				exporterConfigPath = configs[cfgName].(string)
			}
		}
		if !configPathFound {
			configs[exportersConfigPathSetting] = defaultExportersConfigPath
			exporterConfigPath = defaultExportersConfigPath
		}
	}
	return configs, exporterConfigPath, nil
}

// ProcessingRule is subset of the rules supported by nri-prometheus.
type ProcessingRule struct {
	IgnoreMetrics []IgnoreRule `mapstructure:"ignore_metrics" json:"ignore_metrics,omitempty"`
}

// IgnoreRule skips for processing metrics that match any of the Prefixes.
// Metrics that match any of the Except are never skipped.
// If Prefixes is empty and Except is not, then all metrics that do not
// match Except will be skipped.
type IgnoreRule struct {
	Prefixes []string `mapstructure:"prefixes" json:"prefixes,omitempty"`
	Except   []string `mapstructure:"except" json:"except,omitempty"`
}
