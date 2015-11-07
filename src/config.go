package main

import (
	"time"

	"github.com/kardianos/osext"
	"github.com/spf13/viper"
)

var configLoaded bool

func init() {
	configLoaded = false
}

func loadConfig() error {
	if !configLoaded {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		if sourcePWD, err := osext.ExecutableFolder(); err == nil {
			viper.AddConfigPath(sourcePWD)
		}
		viper.SetDefault("interval", "1m")

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		configLoaded = true
	}
	return nil
}

func getChecks() []Check {
	var checks []Check
	if err := viper.UnmarshalKey("checks", &checks); err != nil {
		return nil
	}
	return checks
}

func getInterval() time.Duration {
	interval := viper.GetString("interval")
	if d, err := time.ParseDuration(interval); err == nil {
		return d
	}
	return 1 * time.Minute
}
