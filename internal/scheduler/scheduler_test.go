package scheduler

import (
	"errors"
	"solar-scraper/internal/influx"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Status_SubstituteCurrentStatus(t *testing.T) {
	tests := []struct {
		name               string
		input              status
		inputAfterRun      status
		maxSustainedErrors uint
		err                error
		output             bool
	}{
		{name: "Error, Last, ErrorCount < maxSustainedErrors",
			input: status{
				Last:       influx.SolarMetrics{Today: 20, Total: 200},
				Populated:  true,
				ErrorCount: 0,
			},
			inputAfterRun: status{
				Current:    influx.SolarMetrics{NowNil: true, Today: 20, Total: 200},
				Last:       influx.SolarMetrics{Today: 20, Total: 200},
				Populated:  true,
				ErrorCount: 1,
			},
			maxSustainedErrors: 2,
			err:                errors.New("test error"),
			output:             true,
		},
		{name: "Error, Last, ErrorCount = maxSustainedErrors",
			input: status{
				Last:       influx.SolarMetrics{Today: 20, Total: 200},
				Populated:  true,
				ErrorCount: 10,
			},
			inputAfterRun: status{
				Current:    influx.SolarMetrics{NowNil: true, Today: 20, Total: 200},
				Last:       influx.SolarMetrics{Today: 20, Total: 200},
				Populated:  true,
				ErrorCount: 11,
			},
			maxSustainedErrors: 10,
			err:                errors.New("test error"),
			output:             true,
		},
		{name: "Error, Last, ErrorCount > maxSustainedErrors",
			input: status{
				Last:       influx.SolarMetrics{Today: 20, Total: 200},
				Populated:  true,
				ErrorCount: 11,
			},
			inputAfterRun: status{
				Last:       influx.SolarMetrics{Today: 20, Total: 200},
				Populated:  true,
				ErrorCount: 12,
			},
			maxSustainedErrors: 10,
			err:                errors.New("test error"),
			output:             false,
		},
		{name: "Error, No Last, ErrorCount < maxSustainedErrors",
			input:              status{},
			inputAfterRun:      status{ErrorCount: 1},
			maxSustainedErrors: 2,
			err:                errors.New("test error"),
			output:             false,
		},
		{name: "Error, No Last, ErrorCount = maxSustainedErrors",
			input:              status{ErrorCount: 10},
			inputAfterRun:      status{ErrorCount: 11},
			maxSustainedErrors: 10,
			err:                errors.New("test error"),
			output:             false,
		},
		{name: "Error, No Last, ErrorCount > maxSustainedErrors",
			input:              status{ErrorCount: 11},
			inputAfterRun:      status{ErrorCount: 12},
			maxSustainedErrors: 10,
			err:                errors.New("test error"),
			output:             false,
		},
		{name: "No Error, Last",
			input: status{
				Current:    influx.SolarMetrics{Now: 20, Today: 120, Total: 1220},
				Last:       influx.SolarMetrics{Today: 1, Total: 1},
				Populated:  true,
				ErrorCount: 4,
			},
			inputAfterRun: status{
				Current:    influx.SolarMetrics{Now: 20, Today: 120, Total: 1220},
				Last:       influx.SolarMetrics{Today: 120, Total: 1220},
				Populated:  true,
				ErrorCount: 0,
			},
			output: true,
		},
		{name: "No Error, No Last",
			input: status{
				Current: influx.SolarMetrics{Now: 20, Today: 120, Total: 1220},
			},
			inputAfterRun: status{
				Current:   influx.SolarMetrics{Now: 20, Today: 120, Total: 1220},
				Last:      influx.SolarMetrics{Today: 120, Total: 1220},
				Populated: true,
			},
			output: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			output := test.input.SubstituteCurrentStatus(test.err, test.maxSustainedErrors)
			require.Equal(t, test.output, output, test.name)
			require.Equal(t, test.inputAfterRun, test.input, test.name)
		})
	}
}
