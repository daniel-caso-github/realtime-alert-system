package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

func TestNewAlert_Success(t *testing.T) {
	// Act
	alert, err := entity.NewAlert("High CPU", "CPU usage at 95%", entity.AlertSeverityCritical, "server-01")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, alert)
	assert.NotEqual(t, entity.ID{}, alert.ID)
	assert.Equal(t, "High CPU", alert.Title)
	assert.Equal(t, "CPU usage at 95%", alert.Message)
	assert.Equal(t, entity.AlertSeverityCritical, alert.Severity)
	assert.Equal(t, entity.AlertStatusActive, alert.Status)
	assert.Equal(t, "server-01", alert.Source)
	assert.NotNil(t, alert.Metadata)
}

func TestNewAlert_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name        string
		title       string
		message     string
		severity    entity.AlertSeverity
		expectedErr error
	}{
		{
			name:        "empty title",
			title:       "",
			message:     "message",
			severity:    entity.AlertSeverityMedium,
			expectedErr: entity.ErrAlertTitleRequired,
		},
		{
			name:        "empty message",
			title:       "title",
			message:     "",
			severity:    entity.AlertSeverityMedium,
			expectedErr: entity.ErrAlertMessageRequired,
		},
		{
			name:        "invalid severity",
			title:       "title",
			message:     "message",
			severity:    entity.AlertSeverity("invalid"),
			expectedErr: entity.ErrAlertInvalidSeverity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			alert, err := entity.NewAlert(tc.title, tc.message, tc.severity, "source")

			assert.Nil(t, alert)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestAlertSeverity_Priority(t *testing.T) {
	assert.Equal(t, 1, entity.AlertSeverityCritical.Priority())
	assert.Equal(t, 2, entity.AlertSeverityHigh.Priority())
	assert.Equal(t, 3, entity.AlertSeverityMedium.Priority())
	assert.Equal(t, 4, entity.AlertSeverityLow.Priority())
	assert.Equal(t, 5, entity.AlertSeverityInfo.Priority())

	assert.Less(t, entity.AlertSeverityCritical.Priority(), entity.AlertSeverityHigh.Priority())
}

func TestAlert_Acknowledge(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")
	userID := entity.NewID()

	// Act
	err := alert.Acknowledge(userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, entity.AlertStatusAcknowledged, alert.Status)
	assert.NotNil(t, alert.AcknowledgedBy)
	assert.Equal(t, userID, *alert.AcknowledgedBy)
	assert.NotNil(t, alert.AcknowledgedAt)
}

func TestAlert_Acknowledge_AlreadyAcknowledged(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")
	_ = alert.Acknowledge(entity.NewID())

	// Act
	err := alert.Acknowledge(entity.NewID())

	// Assert
	assert.ErrorIs(t, err, entity.ErrAlertAlreadyAcknowledged)
}

func TestAlert_Acknowledge_AlreadyResolved(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")
	_ = alert.Resolve(entity.NewID())

	// Act
	err := alert.Acknowledge(entity.NewID())

	// Assert
	assert.ErrorIs(t, err, entity.ErrAlertAlreadyResolved)
}

func TestAlert_Resolve(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")
	userID := entity.NewID()

	// Act
	err := alert.Resolve(userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, entity.AlertStatusResolved, alert.Status)
	assert.NotNil(t, alert.ResolvedBy)
	assert.Equal(t, userID, *alert.ResolvedBy)
	assert.NotNil(t, alert.ResolvedAt)
}

func TestAlert_Resolve_FromAcknowledged(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")
	userID := entity.NewID()
	_ = alert.Acknowledge(userID)

	// Act
	err := alert.Resolve(userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, entity.AlertStatusResolved, alert.Status)
}

func TestAlert_Resolve_AlreadyResolved(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")
	_ = alert.Resolve(entity.NewID())

	// Act
	err := alert.Resolve(entity.NewID())

	// Assert
	assert.ErrorIs(t, err, entity.ErrAlertAlreadyResolved)
}

func TestAlert_Expire(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")

	// Act
	alert.Expire()

	// Assert
	assert.Equal(t, entity.AlertStatusExpired, alert.Status)
}

func TestAlert_IsExpired(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")

	// No expiration set
	assert.False(t, alert.IsExpired())

	// Set future expiration
	future := time.Now().Add(1 * time.Hour)
	alert.SetExpiration(future)
	assert.False(t, alert.IsExpired())

	// Set past expiration
	past := time.Now().Add(-1 * time.Hour)
	alert.SetExpiration(past)
	assert.True(t, alert.IsExpired())
}

func TestAlert_AddMetadata(t *testing.T) {
	// Arrange
	alert, _ := entity.NewAlert("Test", "Message", entity.AlertSeverityMedium, "source")

	// Act
	alert.AddMetadata("cpu_percent", 95.5)
	alert.AddMetadata("hostname", "server-01")

	// Assert
	assert.Equal(t, 95.5, alert.Metadata["cpu_percent"])
	assert.Equal(t, "server-01", alert.Metadata["hostname"])
}

func TestAlert_NeedsImmediateAttention(t *testing.T) {
	testCases := []struct {
		name     string
		severity entity.AlertSeverity
		status   entity.AlertStatus
		expected bool
	}{
		{"critical active", entity.AlertSeverityCritical, entity.AlertStatusActive, true},
		{"high active", entity.AlertSeverityHigh, entity.AlertStatusActive, true},
		{"medium active", entity.AlertSeverityMedium, entity.AlertStatusActive, false},
		{"critical acknowledged", entity.AlertSeverityCritical, entity.AlertStatusAcknowledged, false},
		{"critical resolved", entity.AlertSeverityCritical, entity.AlertStatusResolved, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			alert, _ := entity.NewAlert("Test", "Message", tc.severity, "source")

			switch tc.status {
			case entity.AlertStatusAcknowledged:
				_ = alert.Acknowledge(entity.NewID())
			case entity.AlertStatusResolved:
				_ = alert.Resolve(entity.NewID())
			}

			assert.Equal(t, tc.expected, alert.NeedsImmediateAttention())
		})
	}
}
