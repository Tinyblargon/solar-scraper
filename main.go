package main

import (
	"solar-scraper/internal/config"
	"solar-scraper/internal/flags"
	"solar-scraper/internal/logger"
	"solar-scraper/internal/scheduler"
)

var version string // Set by build script

func main() {
	options := flags.Parse(version)
	log := logger.New(options.Log, options.Debug)
	config, err := config.Get(options.Config)
	if err != nil {
		log.Error.Fatal(err)
	}
	metricsWriter := config.InfluxDB.CreateWriter()
	if err = metricsWriter.Ping(); err != nil {
		log.Error.Fatal(err)
	}
	scheduler.Run(config.Time, config.Scraper, metricsWriter, log.Debug, log.Error)
}
