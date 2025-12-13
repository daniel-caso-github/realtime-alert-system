package entity

import (
	"errors"
)

// AlertRule define las condiciones para disparar alertas automáticamente.
type AlertRule struct {
	ID              ID            `json:"id" db:"id"`
	Name            string        `json:"name" db:"name"`
	Description     string        `json:"description,omitempty" db:"description"`
	Condition       RuleCondition `json:"condition" db:"condition"`
	Severity        AlertSeverity `json:"severity" db:"severity"`
	IsEnabled       bool          `json:"is_enabled" db:"is_enabled"`
	CooldownMinutes int           `json:"cooldown_minutes" db:"cooldown_minutes"`
	CreatedBy       *ID           `json:"created_by,omitempty" db:"created_by"`
	Timestamps
}

// RuleCondition define la condición que dispara la regla.
// Se almacena como JSON en la base de datos.
type RuleCondition struct {
	Metric      string  `json:"metric"`
	Operator    string  `json:"operator"`
	Threshold   float64 `json:"threshold"`
	Consecutive int     `json:"consecutive"`
}

// Errores de validación de reglas.
var (
	ErrRuleNameRequired      = errors.New("rule name is required")
	ErrRuleNameTooLong       = errors.New("rule name must be less than 256 characters")
	ErrRuleInvalidSeverity   = errors.New("invalid rule severity")
	ErrRuleInvalidCooldown   = errors.New("cooldown must be between 0 and 1440 minutes")
	ErrRuleConditionRequired = errors.New("rule condition is required")
	ErrRuleInvalidOperator   = errors.New("invalid operator, must be one of: >, <, ==, >=, <=, !=")
	ErrRuleMetricRequired    = errors.New("condition metric is required")
)

// Operadores válidos para las condiciones.
var validOperators = map[string]bool{
	">":  true,
	"<":  true,
	"==": true,
	">=": true,
	"<=": true,
	"!=": true,
}

// NewAlertRule crea una nueva regla de alerta.
func NewAlertRule(name, description string, condition RuleCondition, severity AlertSeverity, createdBy *ID) (*AlertRule, error) {
	rule := &AlertRule{
		ID:              NewID(),
		Name:            name,
		Description:     description,
		Condition:       condition,
		Severity:        severity,
		IsEnabled:       true,
		CooldownMinutes: 5,
		CreatedBy:       createdBy,
		Timestamps:      NewTimestamps(),
	}

	if err := rule.Validate(); err != nil {
		return nil, err
	}

	return rule, nil
}

// Validate verifica que la regla sea válida.
func (r *AlertRule) Validate() error {
	if r.Name == "" {
		return ErrRuleNameRequired
	}

	if len(r.Name) > 255 {
		return ErrRuleNameTooLong
	}

	if !r.Severity.IsValid() {
		return ErrRuleInvalidSeverity
	}

	if r.CooldownMinutes < 0 || r.CooldownMinutes > 1440 {
		return ErrRuleInvalidCooldown
	}

	// Validar condición
	if err := r.Condition.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate verifica que la condición sea válida.
func (c *RuleCondition) Validate() error {
	if c.Metric == "" {
		return ErrRuleMetricRequired
	}

	if !validOperators[c.Operator] {
		return ErrRuleInvalidOperator
	}

	return nil
}

// Enable habilita la regla.
func (r *AlertRule) Enable() {
	r.IsEnabled = true
	r.Touch()
}

// Disable deshabilita la regla.
func (r *AlertRule) Disable() {
	r.IsEnabled = false
	r.Touch()
}

// SetCooldown establece el tiempo de cooldown.
func (r *AlertRule) SetCooldown(minutes int) error {
	if minutes < 0 || minutes > 1440 {
		return ErrRuleInvalidCooldown
	}
	r.CooldownMinutes = minutes
	r.Touch()
	return nil
}

// Evaluate evalúa si un valor cumple la condición de la regla.
// Retorna true si la condición se cumple (debería dispararse una alerta).
func (r *AlertRule) Evaluate(value float64) bool {
	if !r.IsEnabled {
		return false
	}

	switch r.Condition.Operator {
	case ">":
		return value > r.Condition.Threshold
	case "<":
		return value < r.Condition.Threshold
	case "==":
		return value == r.Condition.Threshold
	case ">=":
		return value >= r.Condition.Threshold
	case "<=":
		return value <= r.Condition.Threshold
	case "!=":
		return value != r.Condition.Threshold
	default:
		return false
	}
}
