package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-crowdfunding/auth"
	"backend-crowdfunding/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type configResponse struct {
	Data struct {
		PaymentsEnabled bool   `json:"payments_enabled"`
		SupportEmail    string `json:"support_email"`
	} `json:"data"`
}

// getConfig hits the public /config route on a router built without real
// handlers — the route only depends on the config, so the DB-backed handlers
// can stay nil. It also pins the JSON shape the web app's zod schema expects.
func getConfig(t *testing.T, cfg *config.Config) (*httptest.ResponseRecorder, configResponse) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := newRouter(cfg, nil, nil, nil, auth.NewService("test-secret"), nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/config", nil)
	router.ServeHTTP(w, req)

	var body configResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	return w, body
}

func TestConfigRouteReportsDemoMode(t *testing.T) {
	w, body := getConfig(t, &config.Config{
		SupportEmail: "hi@wikukarno.dev",
		Midtrans:     config.MidtransConfig{Enabled: false, ServerKey: "k", ClientKey: "c"},
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, body.Data.PaymentsEnabled)
	assert.Equal(t, "hi@wikukarno.dev", body.Data.SupportEmail)
}

func TestConfigRouteReportsLiveWhenEnabledWithKeys(t *testing.T) {
	_, body := getConfig(t, &config.Config{
		SupportEmail: "hi@wikukarno.dev",
		Midtrans:     config.MidtransConfig{Enabled: true, ServerKey: "k", ClientKey: "c"},
	})

	assert.True(t, body.Data.PaymentsEnabled)
}

func TestConfigRouteStaysDemoWhenEnabledButKeysMissing(t *testing.T) {
	// Flag on but no keys: Available() must keep us in demo mode so we never
	// advertise a checkout that would fail.
	_, body := getConfig(t, &config.Config{
		SupportEmail: "hi@wikukarno.dev",
		Midtrans:     config.MidtransConfig{Enabled: true, ServerKey: "", ClientKey: ""},
	})

	assert.False(t, body.Data.PaymentsEnabled)
}
