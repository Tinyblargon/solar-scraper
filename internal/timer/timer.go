package timer

import (
	"errors"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	ErrorTimeFrameTooSmall string = "time frame is smaller than the polling interval"
	ErrorInvalidTimeFormat string = "invalid time format, expected HH:MM:SS"
	ErrorPollingInterval   string = "polling interval must be greater than 0"
)

type Settings struct {
	end                      timeObject `mapstructure:"-"`
	End                      string     `mapstructure:"end"`
	start                    timeObject `mapstructure:"-"`
	Start                    string     `mapstructure:"start"`
	PollingIntervalInSeconds uint       `mapstructure:"polling_interval"`
}

func (s Settings) Defaults(setting string) {
	viper.SetDefault(setting+".end", "23:59:59")
	viper.SetDefault(setting+".start", "00:00:00")
	viper.SetDefault(setting+".polling_interval", uint(60))
}

func (s Settings) GetEndTime() timeObject {
	return s.end
}

func (s Settings) GetStartTime() timeObject {
	return s.start
}

func (s *Settings) Validate() (err error) {
	if s.PollingIntervalInSeconds == 0 {
		return errors.New(ErrorPollingInterval)
	}
	s.end, err = parse(s.End)
	if err != nil {
		return
	}
	s.start, err = parse(s.Start)
	if err != nil {
		return
	}
	return s.validateTimeFrameSize()
}

// Check if difference between start and end time is bigger than the polling interval
func (s Settings) validateTimeFrameSize() error {
	end := s.end.Hour()*3600 + s.end.Minute()*60 + s.end.Second()
	start := s.start.Hour()*3600 + s.start.Minute()*60 + s.start.Second()
	if end-start > int(s.PollingIntervalInSeconds) {
		return nil
	}
	return errors.New(ErrorTimeFrameTooSmall)
}

type timeObject struct {
	hour   int
	minute int
	second int
}

func (t timeObject) Hour() int {
	return t.hour
}

func (t timeObject) Minute() int {
	return t.minute
}

func (t timeObject) Second() int {
	return t.second
}

func parse(rawTime string) (timeObject timeObject, err error) {
	timeArray := strings.Split(rawTime, ":")
	if len(timeArray) > 3 {
		return timeObject, errors.New(ErrorInvalidTimeFormat)
	}
	var tmpTime time.Time
	if len(timeArray) > 0 {
		tmpTime, err = time.Parse("15", timeArray[0])
		if err != nil {
			return
		}
		timeObject.hour = tmpTime.Hour()
	}
	if len(timeArray) > 1 {
		tmpTime, err = time.Parse("04", timeArray[1])
		if err != nil {
			return
		}
		timeObject.minute = tmpTime.Minute()
	}
	if len(timeArray) > 2 {
		tmpTime, err = time.Parse("05", timeArray[2])
		if err != nil {
			return
		}
		timeObject.second = tmpTime.Second()
		return
	}
	return
}
