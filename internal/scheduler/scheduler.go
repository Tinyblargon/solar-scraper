package scheduler

import (
	"context"
	"log"
	"solar-scraper/internal/influx"
	"solar-scraper/internal/scraper"
	"solar-scraper/internal/timer"
	"time"

	"github.com/procyon-projects/chrono"
)

func Run(timeS timer.Settings, scraperS scraper.Settings, metricsWriter influx.MetricsWriter, debugLog, errorLog *log.Logger) {
	end := timeS.GetEndTime()
	start := timeS.GetStartTime()

	credentials := scraper.EncodeCredentials(scraperS.Username, scraperS.Password)
	pollingInterval := time.Duration(timeS.PollingIntervalInSeconds) * time.Second
	for {
		currentTime := time.Now()
		startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), start.Hour(), start.Minute(), start.Second(), 0, currentTime.Location())
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), end.Hour(), end.Minute(), end.Second(), 0, currentTime.Location())
		nextStartTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, start.Hour(), start.Minute(), start.Second(), 0, currentTime.Location())
		if currentTime.Before(startTime) {
			// if current time is before start time wait till start time
			time.Sleep(startTime.Sub(currentTime))
		} else if currentTime.After(endTime) {
			// if current time is after end time wait till next start time
			time.Sleep(nextStartTime.Sub(currentTime))
			continue
		}

		runStatus := status{}
		var err error
		var reportingTime time.Time
		task, _ := chrono.NewDefaultTaskScheduler().ScheduleAtFixedRate(func(ctx context.Context) {
			runStatus.Current, reportingTime, err = scraper.GetStatus(scraperS.URL, credentials, scraperS.Retry)
			if err != nil {
				errorLog.Println(err)
			}
			if runStatus.SubstituteCurrentStatus(err, scraperS.MaxSustainedErrors) {
				if err = metricsWriter.Write(runStatus.Current, reportingTime, debugLog); err != nil {
					errorLog.Println(err)
				}
			}
		}, pollingInterval)
		// if current time is after start time and before end time
		time.Sleep(time.Until(endTime))
		task.Cancel()
	}
}

type status struct {
	Current    influx.SolarMetrics
	Last       influx.SolarMetrics
	Populated  bool
	ErrorCount uint
}

func (stat *status) SubstituteCurrentStatus(err error, maxSustainedErrors uint) bool {
	if err != nil {
		if stat.ErrorCount <= maxSustainedErrors {
			if stat.Populated {
				stat.Current = influx.SolarMetrics{
					NowNil: true,
					Today:  stat.Last.Today,
					Total:  stat.Last.Total,
				}
				stat.ErrorCount++
				return true
			}
		}
		stat.ErrorCount++
		return false
	}
	stat.Last = influx.SolarMetrics{
		Today: stat.Current.Today,
		Total: stat.Current.Total,
	}
	stat.Populated = true
	stat.ErrorCount = 0
	return true
}
