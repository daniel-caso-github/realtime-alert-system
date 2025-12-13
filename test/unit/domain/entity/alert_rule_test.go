package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

func TestNewAlertRule_Success(t *testing.T) {
	// Arrange
	condition := entity.RuleCondition{
		Metric:    "cpu_usage",
		Operator:  ">",
		Threshold: 90,
	}

	// Act
	rule, err := entity.NewAlertRule("High CPU", "Triggers when CPU > 90%", condition, entity.AlertSeverityHigh, nil)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "High CPU", rule.Name)
	assert.Equal(t, entity.AlertSeverityHigh, rule.Severity)
	assert.True(t, rule.IsEnabled)
	assert.Equal(t, 5, rule.CooldownMinutes)
}

func TestNewAlertRule_ValidationErrors(t *testing.T) {
	validCondition := entity.RuleCondition{Metric: "cpu", Operator: ">", Threshold: 90}

	testCases := []struct {
		name        string
		ruleName    string
		condition   entity.RuleCondition
		severity    entity.AlertSeverity
		expectedErr error
	}{
		{
			name:        "empty name",
			ruleName:    "",
			condition:   validCondition,
			severity:    entity.AlertSeverityHigh,
			expectedErr: entity.ErrRuleNameRequired,
		},
		{
			name:        "invalid severity",
			ruleName:    "Test Rule",
			condition:   validCondition,
			severity:    entity.AlertSeverity("invalid"),
			expectedErr: entity.ErrRuleInvalidSeverity,
		},
		{
			name:        "empty metric",
			ruleName:    "Test Rule",
			condition:   entity.RuleCondition{Metric: "", Operator: ">", Threshold: 90},
			severity:    entity.AlertSeverityHigh,
			expectedErr: entity.ErrRuleMetricRequired,
		},
		{
			name:        "invalid operator",
			ruleName:    "Test Rule",
			condition:   entity.RuleCondition{Metric: "cpu", Operator: "invalid", Threshold: 90},
			severity:    entity.AlertSeverityHigh,
			expectedErr: entity.ErrRuleInvalidOperator,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule, err := entity.NewAlertRule(tc.ruleName, "desc", tc.condition, tc.severity, nil)

			assert.Nil(t, rule)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestAlertRule_Evaluate(t *testing.T) {
	testCases := []struct {
		name      string
		operator  string
		threshold float64
		value     float64
		expected  bool
	}{
		{"greater than - true", ">", 90, 95, true},
		{"greater than - false", ">", 90, 85, false},
		{"less than - true", "<", 10, 5, true},
		{"less than - false", "<", 10, 15, false},
		{"equal - true", "==", 50, 50, true},
		{"equal - false", "==", 50, 51, false},
		{"greater or equal - true (greater)", ">=", 90, 95, true},
		{"greater or equal - true (equal)", ">=", 90, 90, true},
		{"greater or equal - false", ">=", 90, 85, false},
		{"less or equal - true", "<=", 10, 5, true},
		{"less or equal - false", "<=", 10, 15, false},
		{"not equal - true", "!=", 50, 51, true},
		{"not equal - false", "!=", 50, 50, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			condition := entity.RuleCondition{
				Metric:    "test",
				Operator:  tc.operator,
				Threshold: tc.threshold,
			}
			rule, _ := entity.NewAlertRule("Test", "desc", condition, entity.AlertSeverityMedium, nil)

			result := rule.Evaluate(tc.value)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAlertRule_Evaluate_Disabled(t *testing.T) {
	condition := entity.RuleCondition{Metric: "cpu", Operator: ">", Threshold: 90}
	rule, _ := entity.NewAlertRule("Test", "desc", condition, entity.AlertSeverityMedium, nil)

	assert.True(t, rule.Evaluate(95))

	rule.Disable()
	assert.False(t, rule.Evaluate(95))
}

func TestAlertRule_SetCooldown(t *testing.T) {
	condition := entity.RuleCondition{Metric: "cpu", Operator: ">", Threshold: 90}
	rule, _ := entity.NewAlertRule("Test", "desc", condition, entity.AlertSeverityMedium, nil)

	// Valid cooldown
	err := rule.SetCooldown(30)
	assert.NoError(t, err)
	assert.Equal(t, 30, rule.CooldownMinutes)

	// Invalid cooldown - negative
	err = rule.SetCooldown(-1)
	assert.ErrorIs(t, err, entity.ErrRuleInvalidCooldown)

	// Invalid cooldown - too large
	err = rule.SetCooldown(1500)
	assert.ErrorIs(t, err, entity.ErrRuleInvalidCooldown)
}

func TestAlertRule_EnableDisable(t *testing.T) {
	condition := entity.RuleCondition{Metric: "cpu", Operator: ">", Threshold: 90}
	rule, _ := entity.NewAlertRule("Test", "desc", condition, entity.AlertSeverityMedium, nil)

	assert.True(t, rule.IsEnabled)

	rule.Disable()
	assert.False(t, rule.IsEnabled)

	rule.Enable()
	assert.True(t, rule.IsEnabled)
}
