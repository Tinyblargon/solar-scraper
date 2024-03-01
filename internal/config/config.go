package config

import (
	"fmt"
	"path"
	"solar-scraper/internal/influx"
	"solar-scraper/internal/scraper"
	"solar-scraper/internal/timer"

	"github.com/spf13/viper"
)

// Settings is the configuration for the application
type Settings struct {
	Time     timer.Settings   `mapstructure:"time"`
	Scraper  scraper.Settings `mapstructure:"scraper"`
	InfluxDB influx.Settings  `mapstructure:"influxdb"`
}

func (s *Settings) validate() error {
	if err := s.Time.Validate(); err != nil {
		return err
	}
	if err := s.Scraper.Validate(); err != nil {
		return err
	}
	return s.InfluxDB.Validate()
}

func (s Settings) defaults() {
	s.Time.Defaults("time")
	s.Scraper.Defaults("scraper")
	s.InfluxDB.Defaults("influxdb")
}

func Get(configPath string) (Settings, error) {

	// Set the path to look for the configurations file
	if len(configPath) == 0 {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigName(path.Base(configPath))
		viper.AddConfigPath(path.Dir(configPath))
	}

	viper.SetConfigType("yml")

	var configuration Settings
	if err := viper.ReadInConfig(); err != nil {
		return configuration, fmt.Errorf("error reading config file, %s", err)
	}

	configuration.defaults()
	if err := viper.Unmarshal(&configuration); err != nil {
		return configuration, fmt.Errorf("unable to decode into struct, %v", err)
	}
	if err := configuration.validate(); err != nil {
		return configuration, fmt.Errorf("unable to validate config, %v", err)
	}
	return configuration, nil
}
