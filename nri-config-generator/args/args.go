package args

import (
	"path/filepath"

	"github.com/newrelic/infra-integrations-sdk/v4/args"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	cfgFormat = "yaml"
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

func argVars(c *ArgumentList) (map[string]interface{}, error) {
	err := args.SetupArgs(c)
	if err != nil {
		return nil, err
	}
	configs := make(map[string]interface{})
	if c.ConfigPath != "" {
		cfg := viper.New()
		cfg.SetConfigType(cfgFormat)
		cfg.AddConfigPath(filepath.Dir(c.ConfigPath))
		cfg.SetConfigName(filepath.Base(c.ConfigPath))
		setViperDefaults(cfg)
		err = cfg.ReadInConfig()
		if err != nil {
			return nil, errors.Wrap(err, "could not read configuration")
		}
		for _, cfgName := range cfg.AllKeys() {
			configs[cfgName] = cfg.Get(cfgName)
		}
	}

	return configs, nil
}
