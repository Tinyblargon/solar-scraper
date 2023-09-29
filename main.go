package main

import (
	"solar-scraper/internal/config"
	"solar-scraper/internal/flags"
	"solar-scraper/internal/logger"
	"solar-scraper/internal/scheduler"
)

func main() {
	options := flags.Parse()
	log := logger.New(options.Log, options.Debug)
	config, err := config.Get(options.Config)
	if err != nil {
		log.Error.Fatal(err)
	}
	scheduler.Run(config.Time, config.Scraper, config.InfluxDB.CreateWriter(), log.Debug, log.Error)
}
