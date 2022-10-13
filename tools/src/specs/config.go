package specs

import (
	"github.com/spf13/viper"
	"strings"
)

// Config
type Config struct {
	IntegrationName              string       `mapstructure:"integrationName"`
	HumanReadableIntegrationName string       `mapstructure:"humanReadableIntegrationName"`
	Definitions                  []Definition `mapstructure:"definitions"` // List of conditions used find metrics that match.
}

// EntityRule contains rules to synthesis an entity
type Definition struct {
	EntityType string      `mapstructure:"type"`
	Conditions []Condition `mapstructure:"conditions"` // List of rules used to determining if a metric belongs to this entity.
}

// Condition decides whether a metric is suitable to be included into an entity based on the metric attributes and metric name.
type Condition struct {
	Attribute string `mapstructure:"attribute"`
	Prefix    string `mapstructure:"prefix"`
	Value     string `mapstructure:"value"`
}

// Match evaluates the condition an attribute by checking either its whole name against `Value` or if it starts with `Prefix`.
func (c Condition) Match(attribute string) bool {
	if c.Value != "" {
		return c.Value == attribute
	}
	// this returns true if c.Prefix is "" and is ok since the attribute exists
	return strings.HasPrefix(attribute, c.Prefix)
}

func LoadConfig() (Config, error) {
	c := Config{}
	cfg := viper.New()
	cfg.AddConfigPath("./")
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")

	err := cfg.ReadInConfig()
	if err != nil {
		return c, err
	}

	err = cfg.Unmarshal(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}
