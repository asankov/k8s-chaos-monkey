package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config is the config of the chaos monkey program.
type Config struct {
	// Namespace is the namespace in which the chaos monkey program will operate.
	Namespace string
	// PeriodInSeconds is the period in which the chaos monkey will run.
	PeriodInSeconds int
}

// NewConfigFromEnv returns a new config, the values of which are fetched from the environment (env variables).
//
// Neither variable is mandatory and they have sensible defaults.
func NewConfigFromEnv() (*Config, error) {
	ns := getenvOrDefault("K8S_CHAOS_NAMESPACE", "default")

	p := getenvOrDefault("K8S_CHAOS_PERIOD_SECONDS", "10")
	period, err := strconv.Atoi(p)
	if err != nil {
		return nil, fmt.Errorf("invalid value was provided for integer K8S_CHAOS_PERIOD_SECONDS: %w", err)
	}

	return &Config{
		Namespace:       ns,
		PeriodInSeconds: period,
	}, nil
}

func getenvOrDefault(env, def string) string {
	v := os.Getenv(env)
	if v == "" {
		return def
	}
	return v
}
