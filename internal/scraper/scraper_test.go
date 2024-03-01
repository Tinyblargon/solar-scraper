package scraper

import (
	"errors"
	"os"
	"solar-scraper/internal/influx"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EncodeCredentials(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		output   credentials
	}{
		{name: "Valid",
			username: "admin",
			password: "enter123!",
			output:   "YWRtaW46ZW50ZXIxMjMh",
		},
		{name: "Valid no username",
			password: "enter123!",
			output:   "OmVudGVyMTIzIQ==",
		},
		{name: "Valid no password",
			username: "admin",
			output:   "YWRtaW46",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			require.Equal(t, test.output, EncodeCredentials(test.username, test.password), test.name)
		})
	}
}

func Test_extractStatusValues(t *testing.T) {
	type testOutput struct {
		stats influx.SolarMetrics
		err   error
	}
	tests := []struct {
		name   string
		input  func() []byte
		output testOutput
	}{
		{name: "Valid minimal",
			input: func() []byte {
				return []byte(`var webdata_now_p = "123"; var webdata_today_e = "123.45"; var webdata_total_e = "1234.56";`)
			},
			output: testOutput{stats: influx.SolarMetrics{Now: 123, Today: 123.45, Total: 1234.56}},
		},
		{name: "Valid real values",
			input: func() []byte {
				data, _ := os.ReadFile("../../test/data/sample.html")
				return data
			},
			output: testOutput{stats: influx.SolarMetrics{Now: 150, Today: 3.10, Total: 4756.2}},
		},
		{name: "Error no now",
			input: func() []byte {
				return []byte(`var webdata_today_e = "123.45"; var webdata_total_e = "1234.56";`)
			},
			output: testOutput{err: errorSearchKeyNotFound(now)},
		},
		{name: "Error no today",
			input: func() []byte {
				return []byte(`var webdata_now_p = "123"; var webdata_total_e = "1234.56";`)
			},
			output: testOutput{stats: influx.SolarMetrics{Now: 123}, err: errorSearchKeyNotFound(today)},
		},
		{name: "Error no total",
			input: func() []byte {
				return []byte(`var webdata_now_p = "123"; var webdata_today_e = "123.45";`)
			},
			output: testOutput{stats: influx.SolarMetrics{Now: 123, Today: 123.45}, err: errorSearchKeyNotFound(total)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			stats, err := extractStatusValues(test.input())
			require.Equal(t, test.output.stats, stats, test.name)
			require.Equal(t, test.output.err, err, test.name)
		})
	}
}

func Test_Validate(t *testing.T) {
	tests := []struct {
		name   string
		input  Settings
		output error
	}{
		{name: "Valid",
			input: Settings{URL: "http://localhost:8086"},
		},
		{name: "ErrorEmptyUrl",
			input:  Settings{},
			output: errors.New(ErrorEmptyUrl),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			require.Equal(t, test.output, test.input.Validate(), test.name)
		})
	}
}
