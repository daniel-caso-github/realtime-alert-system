package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
)

func TestLogin_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	// Test with seeded admin user
	resp := app.MakeRequest("POST", "/api/v1/auth/login", dto.LoginRequest{
		Email:    "admin@alerting.local",
		Password: "Admin123!",
	}, "")

	assert.Equal(t, http.StatusOK, resp.Code)

	var loginResp dto.LoginResponse
	err := json.Unmarshal(resp.Body.Bytes(), &loginResp)
	require.NoError(t, err)

	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)
	assert.Equal(t, "admin@alerting.local", loginResp.User.Email)
	assert.Equal(t, "admin", loginResp.User.Role)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("POST", "/api/v1/auth/login", dto.LoginRequest{
		Email:    "admin@alerting.local",
		Password: "wrongpassword",
	}, "")

	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	var errResp dto.ErrorResponse
	err := json.Unmarshal(resp.Body.Bytes(), &errResp)
	require.NoError(t, err)

	assert.Equal(t, "Invalid email or password", errResp.Error)
}

func TestLogin_ValidationError(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("POST", "/api/v1/auth/login", dto.LoginRequest{
		Email:    "invalid-email",
		Password: "short",
	}, "")

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestRegister_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("POST", "/api/v1/auth/register", dto.RegisterRequest{
		Email:    "test.user@example.com",
		Password: "TestPassword123!",
		Name:     "Test User",
	}, "")

	assert.Equal(t, http.StatusCreated, resp.Code)

	var loginResp dto.LoginResponse
	err := json.Unmarshal(resp.Body.Bytes(), &loginResp)
	require.NoError(t, err)

	assert.NotEmpty(t, loginResp.AccessToken)
	assert.Equal(t, "test.user@example.com", loginResp.User.Email)
	assert.Equal(t, "viewer", loginResp.User.Role) // Default role
}

func TestRegister_DuplicateEmail(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	// First registration
	resp := app.MakeRequest("POST", "/api/v1/auth/register", dto.RegisterRequest{
		Email:    "test.duplicate@example.com",
		Password: "TestPassword123!",
		Name:     "Test User",
	}, "")
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Second registration with same email
	resp = app.MakeRequest("POST", "/api/v1/auth/register", dto.RegisterRequest{
		Email:    "test.duplicate@example.com",
		Password: "TestPassword123!",
		Name:     "Test User 2",
	}, "")

	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestMe_Success(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	token := app.Login(t, "admin@alerting.local", "Admin123!")

	resp := app.MakeRequest("GET", "/api/v1/auth/me", nil, token)

	assert.Equal(t, http.StatusOK, resp.Code)

	var userResp dto.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &userResp)
	require.NoError(t, err)

	assert.Equal(t, "admin@alerting.local", userResp.Email)
}

func TestMe_Unauthorized(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Cleanup(t)

	resp := app.MakeRequest("GET", "/api/v1/auth/me", nil, "")

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
