package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
)

func TestHealthCheck(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("GET", "/health", nil, "")

	assert.Equal(t, http.StatusOK, resp.Code)

	var health dto.HealthResponse
	err := json.Unmarshal(resp.Body.Bytes(), &health)
	require.NoError(t, err)

	assert.Equal(t, "healthy", health.Status)
	assert.NotEmpty(t, health.Services)
}

func TestReadyCheck(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("GET", "/ready", nil, "")

	assert.Equal(t, http.StatusOK, resp.Code)

	var ready dto.ReadyResponse
	err := json.Unmarshal(resp.Body.Bytes(), &ready)
	require.NoError(t, err)

	assert.Equal(t, "ready", ready.Status)
}

func TestLiveCheck(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("GET", "/live", nil, "")

	assert.Equal(t, http.StatusOK, resp.Code)

	var live dto.LiveResponse
	err := json.Unmarshal(resp.Body.Bytes(), &live)
	require.NoError(t, err)

	assert.Equal(t, "alive", live.Status)
}
