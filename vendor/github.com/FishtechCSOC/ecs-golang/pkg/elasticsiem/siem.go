package elasticsiem

import "time"

type ElasticSIEM struct {
	Signal Signal `json:"signal,omitempty" mapstructure:"signal"`
}

// nolint: tagliatelle
type Signal struct {
	Rule         Rule      `json:"rule,omitempty" mapstructure:"rule"`
	OriginalTime time.Time `json:"original_time,omitempty" mapstructure:"original_time"`
	Status       string    `json:"status,omitempty" mapstructure:"status"`
}

// nolint: tagliatelle
type Rule struct {
	ID             string    `json:"id,omitempty" mapstructure:"id"`
	RuleID         string    `json:"rule_id,omitempty" mapstructure:"rule_id"`
	Actions        []string  `json:"actions,omitempty" mapstructure:"actions"`
	FalsePositives []string  `json:"false_positives,omitempty" mapstructure:"false_positives"`
	MaxSignals     int64     `json:"max_signals,omitempty" mapstructure:"max_signals"`
	RiskScore      int64     `json:"risk_score,omitempty" mapstructure:"risk_score"`
	OutputIndex    string    `json:"output_index,omitempty" mapstructure:"output_index"`
	Description    string    `json:"description,omitempty" mapstructure:"description"`
	From           string    `json:"from,omitempty" mapstructure:"from"`
	Immutable      bool      `json:"immutable,omitempty" mapstructure:"immutable"`
	Index          []string  `json:"index,omitempty" mapstructure:"index"`
	Interval       string    `json:"interval,omitempty" mapstructure:"interval"`
	Language       string    `json:"language,omitempty" mapstructure:"language"`
	Name           string    `json:"name,omitempty" mapstructure:"name"`
	Query          string    `json:"query,omitempty" mapstructure:"query"`
	Reference      []string  `json:"references,omitempty" mapstructure:"references"`
	Severity       string    `json:"severity,omitempty" mapstructure:"severity"`
	Tags           []string  `json:"tags,omitempty" mapstructure:"tags"`
	Type           string    `json:"type,omitempty" mapstructure:"type"`
	To             string    `json:"to,omitempty" mapstructure:"to"`
	Enabled        bool      `json:"enabled,omitempty" mapstructure:"enabled"`
	CreatedBy      string    `json:"created_by,omitempty" mapstructure:"created_by"`
	UpdatedBy      string    `json:"updated_by,omitempty" mapstructure:"updated_by"`
	Threat         []string  `json:"threat,omitempty" mapstructure:"threat"`
	Version        int64     `json:"version,omitempty" mapstructure:"version"`
	CreatedAt      time.Time `json:"created_at,omitempty" mapstructure:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" mapstructure:"updated_at"`
}
