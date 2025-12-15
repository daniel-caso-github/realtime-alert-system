package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
)

func TestCreateAlert_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	resp := app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
		Title:    "Test Alert",
		Message:  "This is a test alert",
		Severity: "high",
		Source:   "integration-test",
	}, token)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var alertResp dto.AlertResponse
	err := json.Unmarshal(resp.Body.Bytes(), &alertResp)
	require.NoError(t, err)

	assert.NotEmpty(t, alertResp.ID)
	assert.Equal(t, "Test Alert", alertResp.Title)
	assert.Equal(t, "high", alertResp.Severity)
	assert.Equal(t, "active", alertResp.Status)
}

func TestCreateAlert_Unauthorized(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
		Title:    "Test Alert",
		Message:  "This is a test alert",
		Severity: "high",
	}, "")

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestCreateAlert_ValidationError(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	resp := app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
		Title:    "", // Empty title
		Message:  "This is a test alert",
		Severity: "high",
	}, token)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestListAlerts_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	// Create a few alerts
	for i := 0; i < 3; i++ {
		app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
			Title:    fmt.Sprintf("Test Alert %d", i),
			Message:  "Test message",
			Severity: "medium",
		}, token)
	}

	// List alerts
	resp := app.MakeRequest("GET", "/api/v1/alerts", nil, token)

	assert.Equal(t, http.StatusOK, resp.Code)

	var listResp dto.PaginatedAlertResponse
	err := json.Unmarshal(resp.Body.Bytes(), &listResp)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(listResp.Items), 3)
}

func TestGetAlert_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	// Create alert
	createResp := app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
		Title:    "Test Get Alert",
		Message:  "Test message",
		Severity: "low",
	}, token)

	var created dto.AlertResponse
	_ = json.Unmarshal(createResp.Body.Bytes(), &created)

	// Get alert
	resp := app.MakeRequest("GET", "/api/v1/alerts/"+created.ID, nil, token)

	assert.Equal(t, http.StatusOK, resp.Code)

	var alertResp dto.AlertResponse
	err := json.Unmarshal(resp.Body.Bytes(), &alertResp)
	require.NoError(t, err)

	assert.Equal(t, created.ID, alertResp.ID)
	assert.Equal(t, "Test Get Alert", alertResp.Title)
}

func TestGetAlert_NotFound(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	resp := app.MakeRequest("GET", "/api/v1/alerts/00000000-0000-0000-0000-000000000000", nil, token)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestAcknowledgeAlert_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	// Create alert
	createResp := app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
		Title:    "Test Acknowledge Alert",
		Message:  "Test message",
		Severity: "high",
	}, token)

	var created dto.AlertResponse
	_ = json.Unmarshal(createResp.Body.Bytes(), &created)

	// Acknowledge alert
	resp := app.MakeRequest("POST", "/api/v1/alerts/"+created.ID+"/acknowledge", nil, token)

	assert.Equal(t, http.StatusOK, resp.Code)

	var alertResp dto.AlertResponse
	err := json.Unmarshal(resp.Body.Bytes(), &alertResp)
	require.NoError(t, err)

	assert.Equal(t, "acknowledged", alertResp.Status)
	assert.NotNil(t, alertResp.AcknowledgedAt)
}

func TestResolveAlert_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	// Create alert
	createResp := app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
		Title:    "Test Resolve Alert",
		Message:  "Test message",
		Severity: "medium",
	}, token)

	var created dto.AlertResponse
	_ = json.Unmarshal(createResp.Body.Bytes(), &created)

	// Resolve alert
	resp := app.MakeRequest("POST", "/api/v1/alerts/"+created.ID+"/resolve", nil, token)

	assert.Equal(t, http.StatusOK, resp.Code)

	var alertResp dto.AlertResponse
	err := json.Unmarshal(resp.Body.Bytes(), &alertResp)
	require.NoError(t, err)

	assert.Equal(t, "resolved", alertResp.Status)
	assert.NotNil(t, alertResp.ResolvedAt)
}

func TestDeleteAlert_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	// Create alert
	createResp := app.MakeRequest("POST", "/api/v1/alerts", dto.CreateAlertRequest{
		Title:    "Test Delete Alert",
		Message:  "Test message",
		Severity: "low",
	}, token)

	var created dto.AlertResponse
	_ = json.Unmarshal(createResp.Body.Bytes(), &created)

	// Delete alert
	resp := app.MakeRequest("DELETE", "/api/v1/alerts/"+created.ID, nil, token)

	assert.Equal(t, http.StatusNoContent, resp.Code)

	// Verify deleted
	resp = app.MakeRequest("GET", "/api/v1/alerts/"+created.ID, nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetStatistics_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	resp := app.MakeRequest("GET", "/api/v1/alerts/statistics", nil, token)

	assert.Equal(t, http.StatusOK, resp.Code)

	var stats dto.AlertStatisticsResponse
	err := json.Unmarshal(resp.Body.Bytes(), &stats)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, stats.TotalAlerts, int64(0))
}
