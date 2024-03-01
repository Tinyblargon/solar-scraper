package scraper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"math"
	"net/http"
	"solar-scraper/internal/influx"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

const (
	now           string = `var webdata_now_p = "`
	today         string = `var webdata_today_e = "`
	total         string = `var webdata_total_e = "`
	ErrorEmptyUrl string = "empty url"
)

func EncodeCredentials(username, password string) credentials {
	return credentials(base64.StdEncoding.EncodeToString([]byte(username + ":" + password)))
}

func extractStatusValues(body []byte) (stats influx.SolarMetrics, err error) {
	value, err := getValue(body, now)
	if err != nil {
		return
	}
	stats.Now = uint(math.Round(value))
	stats.Today, err = getValue(body, today)
	if err != nil {
		return
	}
	stats.Total, err = getValue(body, total)
	return
}

func errorSearchKeyNotFound(search string) error {
	return errors.New("string (" + search + ") not found")
}

func GetStatus(url string, encoded credentials, retry uint) (stats influx.SolarMetrics, reportingTime time.Time, err error) {
	for i := -1; i < int(retry); i++ {
		stats, reportingTime, err = retryStatus(url, encoded)
		if err == nil {
			break
		}
	}
	return
}

func retryStatus(url string, encoded credentials) (stats influx.SolarMetrics, reportingTime time.Time, err error) {
	var resp *http.Response
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Basic "+string(encoded))
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	// must be done directly after getting response to give the most accurate reporting time
	reportingTime = time.Now()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	stats, err = extractStatusValues(body)
	return
}

func getValue(body []byte, search string) (float64, error) {
	split := bytes.SplitAfter(body, []byte(search))
	if len(split) < 2 {
		return 0, errorSearchKeyNotFound(search)
	}
	numberAsByte := []byte{}
	for i, e := range split[1] {
		if e == '"' {
			numberAsByte = split[1][:i]
			break
		}
	}
	return strconv.ParseFloat(string(numberAsByte), 64)
}

type credentials string

type Settings struct {
	MaxSustainedErrors uint   `mapstructure:"sustained_errors"` // The amount of consecutive errors that are applicable for a status substitution. If this value is exceeded nothing wil be written to the database, until valid data is received.
	Password           string `mapstructure:"password"`
	Retry              uint   `mapstructure:"retry"`
	URL                string `mapstructure:"url"`
	Username           string `mapstructure:"username"`
}

func (s Settings) Defaults(setting string) {
	viper.SetDefault(setting+".sustained_errors", uint(5))
	viper.SetDefault(setting+".retry", uint(2))
}

func (s Settings) Validate() error {
	if s.URL == "" {
		return errors.New(ErrorEmptyUrl)
	}
	return nil
}
