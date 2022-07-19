package postgre

import (
	"analytic-service/pkg/logging"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	testTable := []struct {
		caseName  string
		user      string
		password  string
		host      string
		port      string
		nAttempts int
		expErr    error
	}{
		{
			caseName: "incorrect config params",
			user:     "+",
			password: "%",
			host:     "+",
			port:     "-",
			expErr:   ErrParseConfigFile,
		},
		{
			caseName:  "incorrect nAttempts (negative)",
			user:      "admin",
			password:  "secret",
			host:      "127.0.0.1",
			port:      "5432",
			nAttempts: -2,
			expErr:    ErrNAttempts,
		},

		{
			caseName:  "no connection to db",
			user:      "admin",
			password:  "secret",
			host:      "127.0.0.1",
			port:      "5432",
			nAttempts: 1,
			expErr:    ErrConnectionToDb,
		},
	}

	logger := logging.NewLoggerStub()

	for _, testCase := range testTable {
		t.Run(testCase.caseName, func(t *testing.T) {
			_, err := New(logger, testCase.user, testCase.password, testCase.host, testCase.port, testCase.nAttempts)
			require.Equal(t, testCase.expErr, err)
		})
	}
}
