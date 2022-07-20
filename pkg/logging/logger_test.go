package logging

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPositive(t *testing.T) {
	testTable := []struct {
		level      string
		timeFormat string
	}{
		{
			level:      "debug",
			timeFormat: time.RFC3339,
		},
		{
			level:      "debug",
			timeFormat: time.ANSIC,
		},
		{
			level:      "fatal",
			timeFormat: time.Kitchen,
		},
	}
	for _, testCase := range testTable {
		_, err := New(testCase.level, testCase.timeFormat)
		require.Nil(t, err)
	}
}

func TestNewNegative(t *testing.T) {
	testTable := []struct {
		name        string
		timeFormate string
		level       string
	}{
		{
			name:        "incorrect log level",
			timeFormate: time.RFC3339,
			level:       "negative case",
		},
		{
			name:        "number in log level",
			timeFormate: time.RFC3339,
			level:       "INF0",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := New(testCase.level, testCase.timeFormate)
			require.NotNil(t, err)
		})
	}
}
