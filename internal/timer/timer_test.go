package timer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output timeObject
		err    bool
	}{
		// Valid
		{name: "Valid time Hour Minute Second",
			input:  "12:34:56",
			output: timeObject{hour: 12, minute: 34, second: 56},
		},
		{name: "Valid time Hour Minute",
			input:  "12:34",
			output: timeObject{hour: 12, minute: 34, second: 0},
		},
		{name: "Valid time Hour",
			input:  "12",
			output: timeObject{hour: 12, minute: 0, second: 0},
		},
		// Invalid
		{name: "Invalid time",
			input: "12:34:56:78",
			err:   true,
		},
		{name: "Invalid time Empty",
			input: "",
			err:   true,
		},
		{name: "Invalid time Hour > 23",
			input: "24",
			err:   true,
		},
		{name: "Invalid time Minute > 59",
			input: "00:60",
			err:   true,
		},
		{name: "Invalid time Second > 59",
			input: "00:00:60",
			err:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			tmpOut, tmpErr := parse(test.input)
			require.Equal(t, test.output, tmpOut, test.name)
			if test.err {
				require.Error(t, tmpErr, test.name)
			}
		})
	}
}
