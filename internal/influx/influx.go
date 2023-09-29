package influx

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb1 "github.com/influxdata/influxdb1-client/v2"

	"github.com/spf13/viper"
)

func debugMetrics(metrics SolarMetrics, reportTime time.Time, debug *log.Logger) {
	if metrics.NowNil {
		debug.Printf("Time: %s, YieldToday: %f, TotalYield: %f \n", reportTime.Format("20060102150405"), metrics.Today, metrics.Total)
	} else {
		debug.Printf("Time: %s, CurrentPower: %d, YieldToday: %f, TotalYield: %f \n", reportTime.Format("20060102150405"), metrics.Now, metrics.Today, metrics.Total)
	}
}

type MetricsWriter interface {
	Write(metrics SolarMetrics, reportTime time.Time, debug *log.Logger) error
}

type SolarMetrics struct {
	Now    uint
	NowNil bool
	Today  float64
	Total  float64
}

const (
	today       string = "YieldToday"
	total       string = "TotalYield"
	now         string = "CurrentPower"
	measurement string = "PowerYield"
	tagHost     string = "Host"
)

const (
	ErrorV1EmptyDatabase string = "empty database"
	ErrorV1EmptyUsername string = "empty username"
	ErrorV2EmptyOrg      string = "empty org"
	ErrorV2EmptyBucket   string = "empty bucket"
	ErrorEmptyUrl        string = "empty url"
	ErrorInvalidVersion  string = "invalid version"
)

const (
	v1 uint8 = 1
	v2 uint8 = 2
)

type Settings struct {
	Version            uint8      `mapstructure:"version"`
	InsecureSkipVerify bool       `mapstructure:"insecure_skip_verify"`
	Retry              uint       `mapstructure:"retry"`
	Tags               Tags       `mapstructure:"tags"`
	Url                string     `mapstructure:"url"`
	V1                 SettingsV1 `mapstructure:"v1"`
	V2                 SettingsV2 `mapstructure:"v2"`
}

func (s Settings) CreateWriter() MetricsWriter {
	switch s.Version {
	case v1:
		s.V1.url = s.Url
		s.V1.insecureSkipVerify = s.InsecureSkipVerify
		s.V1.tags = s.Tags
		s.V1.retry = s.Retry
		return &s.V1
	case v2:
		s.V2.url = s.Url
		s.V2.insecureSkipVerify = s.InsecureSkipVerify
		s.V2.tags = s.Tags
		s.V2.retry = s.Retry
		return &s.V2
	}
	return nil
}

func (s Settings) Defaults(setting string) {
	viper.SetDefault(setting+".retry", 2)
	viper.SetDefault(setting+".insecure_skip_verify", false)
	// Ignore the error, at worst the default will be empty
	hostname, _ := os.Hostname()
	viper.SetDefault(setting+".tags.host", hostname)

}

func (s Settings) Validate() error {
	if s.Version != 1 && s.Version != 2 {
		return errors.New(ErrorInvalidVersion)
	}
	if s.Url == "" {
		return errors.New(ErrorEmptyUrl)
	}
	switch s.Version {
	case v1:
		if err := s.V1.validate(); err != nil {
			return err
		}
	case v2:
		if err := s.V2.validate(); err != nil {
			return err
		}
	}
	return nil
}

type Tags struct {
	Host string `yml:"host"`
}

type SettingsV1 struct {
	Database           string `mapstructure:"database"`
	Password           string `mapstructure:"password"`
	Username           string `mapstructure:"username"`
	url                string
	insecureSkipVerify bool
	tags               Tags
	retry              uint
}

func (s SettingsV1) validate() error {
	if s.Database == "" {
		return errors.New(ErrorV1EmptyDatabase)
	}
	if s.Username == "" {
		return errors.New(ErrorV1EmptyUsername)
	}
	return nil
}

func (s SettingsV1) Write(metrics SolarMetrics, reportTime time.Time, debug *log.Logger) error {
	if debug != nil {
		debugMetrics(metrics, reportTime, debug)
	}

	client, err := influxdb1.NewHTTPClient(influxdb1.HTTPConfig{
		Addr:               s.url,
		Username:           s.Username,
		Password:           s.Password,
		InsecureSkipVerify: s.insecureSkipVerify,
	})
	if err != nil {
		return errors.New("Error creating InfluxDB Client: " + err.Error())
	}
	defer client.Close()

	bp, err := influxdb1.NewBatchPoints(influxdb1.BatchPointsConfig{
		Database:  s.Database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		today: metrics.Today,
		total: metrics.Total,
	}
	if !metrics.NowNil {
		fields[now] = metrics.Now
	}
	pt, err := influxdb1.NewPoint(measurement, map[string]string{tagHost: s.tags.Host}, fields, reportTime)
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := client.Write(bp); err != nil {
		return err
	}

	// Close client resources
	return client.Close()
}

type SettingsV2 struct {
	Organization       string `mapstructure:"org"`
	Bucket             string `mapstructure:"bucket"`
	AuthToken          string `mapstructure:"auth_token"`
	url                string
	insecureSkipVerify bool
	tags               Tags
	retry              uint
}

func (s SettingsV2) validate() error {
	if s.Organization == "" {
		return errors.New(ErrorV2EmptyOrg)
	}
	if s.Bucket == "" {
		return errors.New(ErrorV2EmptyBucket)
	}
	return nil
}

func (s SettingsV2) Write(metrics SolarMetrics, reportTime time.Time, debug *log.Logger) (err error) {
	if debug != nil {
		debugMetrics(metrics, reportTime, debug)
	}
	client := influxdb2.NewClient(s.url, s.AuthToken)
	defer client.Close()
	client.Options().SetTLSConfig(&tls.Config{InsecureSkipVerify: s.insecureSkipVerify})
	writeAPI := client.WriteAPIBlocking(s.Organization, s.Bucket)
	p := influxdb2.NewPointWithMeasurement(measurement).
		AddTag(tagHost, s.tags.Host).
		AddField(today, metrics.Today).
		AddField(total, metrics.Total).
		SetTime(reportTime)
	if !metrics.NowNil {
		p.AddField(now, metrics.Now)
	}
	for i := -1; i < int(s.retry); i++ {
		err = writeAPI.WritePoint(context.Background(), p)
		if err == nil {
			break
		}
	}
	return
}
