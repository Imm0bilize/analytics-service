package httpServer

import (
	"analytic-service/internal/config"
	"github.com/stretchr/testify/require"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	absPath, err := filepath.Abs("../../configs/dev-config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.New(absPath)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		caseName     string
		port         string
		readTimeout  time.Duration
		writeTimeout time.Duration
	}{
		{
			caseName: "incorrect port",
			port:     "!",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.caseName, func(t *testing.T) {
			cfg.Http.Port = testCase.port

			server := New(cfg, http.NewServeMux())
			server.Run()
			err := <-server.notify
			require.NotNil(t, err)
		})
	}

}
