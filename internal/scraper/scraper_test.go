package scraper

import (
	"errors"
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

func Test_Validate(t *testing.T) {
	tests := []struct {
		name   string
		input  Settings
		output error
	}{
		// TODO add valid case
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
