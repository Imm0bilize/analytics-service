package integration_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const (
	serviceUrl        = "127.0.0.1"
	servicePort       = "8080"
	longTimeLifeToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMSIsImV4cCI6MTY2NDIzNjM0MCwiaWF0IjoxNjU3MDM2MzQwfQ.D6VcRlT3fl-prnWZIu-HmdmRtEQEkQOANsMtAQtTNTw"
)

func (i *integrationTestSuite) TestRestApiPositive() {
	testTable := []struct {
		testName           string
		url                string
		method             string
		expectedStatusCode int
	}{
		{
			testName:           "get-accepted-tasks",
			url:                fmt.Sprintf("http://%s:%s/api/tasks/num-accepted", serviceUrl, servicePort),
			method:             "GET",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "get-rejected-tasks",
			url:                fmt.Sprintf("http://%s:%s/api/tasks/num-rejected", serviceUrl, servicePort),
			method:             "GET",
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "get-total-time",
			url:                fmt.Sprintf("http://%s:%s/api/tasks/total-time", serviceUrl, servicePort),
			method:             "GET",
			expectedStatusCode: http.StatusOK,
		},
	}
	c := http.Client{}
	cookieAccess := &http.Cookie{
		Name:   "access_token",
		Value:  longTimeLifeToken,
		MaxAge: 300,
	}
	cookieRefresh := &http.Cookie{
		Name:   "refresh_token",
		Value:  longTimeLifeToken,
		MaxAge: 300,
	}
	for _, testCase := range testTable {
		i.T().Run(testCase.testName, func(t *testing.T) {
			req, err := http.NewRequest(testCase.method, testCase.url, nil)
			req.AddCookie(cookieAccess)
			req.AddCookie(cookieRefresh)

			if err != nil {
				i.T().Fatal(err)
			}
			res, err := c.Do(req)
			if err != nil {
				i.T().Fatal(err)
			}
			require.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		})
	}
}

func (i *integrationTestSuite) TestAuth() {
	testTable := []struct {
		name               string
		accessToken        *http.Cookie
		refreshToken       *http.Cookie
		expectedStatusCode int
	}{
		{
			name:               "without all tokens",
			accessToken:        &http.Cookie{},
			refreshToken:       &http.Cookie{},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "only access token",
			accessToken: &http.Cookie{
				Name:   "access_token",
				Value:  longTimeLifeToken,
				MaxAge: 300,
			},
			refreshToken:       &http.Cookie{},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "incorrect tokens",
			accessToken: &http.Cookie{
				Name:   "access_token",
				Value:  "test",
				MaxAge: 300,
			},
			refreshToken: &http.Cookie{
				Name:   "refresh_token",
				Value:  "test",
				MaxAge: 300,
			},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name: "correct test",
			accessToken: &http.Cookie{
				Name:   "access_token",
				Value:  longTimeLifeToken,
				MaxAge: 300,
			},
			refreshToken: &http.Cookie{
				Name:   "refresh_token",
				Value:  longTimeLifeToken,
				MaxAge: 300,
			},
			expectedStatusCode: http.StatusOK,
		},
	}
	c := http.Client{}
	for _, testCase := range testTable {
		i.T().Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/api/tasks/total-time", serviceUrl, servicePort), nil)
			if err != nil {
				i.T().Fatal(err)
			}
			req.AddCookie(testCase.accessToken)
			req.AddCookie(testCase.refreshToken)

			res, err := c.Do(req)
			if err != nil {
				i.T().Fatal(err)
			}
			require.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		})
	}

}
