package models

import "time"

// Workflow represents a complete automation workflow
type Workflow struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Version     string        `json:"version"`
	Variables   []Variable    `json:"variables,omitempty"`
	Steps       []Step        `json:"steps"`
	OnError     *ErrorHandler `json:"on_error,omitempty"`
}

// Variable defines a workflow-level variable
type Variable struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"` // string, number, boolean, array, object
}

// Step represents a single operation in the workflow
type Step struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"` // adb, file, system, web, control
	Enabled    bool                   `json:"enabled"`
	Parallel   bool                   `json:"parallel"`
	Retry      *RetryPolicy           `json:"retry,omitempty"`
	Timeout    int                    `json:"timeout,omitempty"` // in seconds
	Parameters map[string]interface{} `json:"parameters"`
	Conditions []Condition            `json:"conditions,omitempty"`
	Loop       *LoopConfig            `json:"loop,omitempty"`
	OnSuccess  []string               `json:"on_success,omitempty"` // step IDs to run on success
	OnFailure  []string               `json:"on_failure,omitempty"` // step IDs to run on failure
}

// RetryPolicy defines retry behavior for a step
type RetryPolicy struct {
	MaxRetries int  `json:"max_retries"`
	Delay      int  `json:"delay"`   // in seconds
	Backoff    bool `json:"backoff"` // exponential backoff
}

// Condition represents a conditional check
type Condition struct {
	Variable string      `json:"variable"`
	Operator string      `json:"operator"` // equals, not_equals, greater, less, contains, exists
	Value    interface{} `json:"value"`
}

// LoopConfig defines loop behavior for a step
type LoopConfig struct {
	Type     string        `json:"type"` // count, while, for_each
	Count    int           `json:"count,omitempty"`
	Variable string        `json:"variable,omitempty"` // for variable-based loops
	Values   []interface{} `json:"values,omitempty"`   // for for_each loops
}

// ErrorHandler defines error handling behavior
type ErrorHandler struct {
	Action string   `json:"action"`          // continue, stop, retry
	Steps  []string `json:"steps,omitempty"` // step IDs to run on error
	Notify bool     `json:"notify"`
}

// ExecutionContext holds runtime context for workflow execution
type ExecutionContext struct {
	Variables map[string]interface{}
	Results   map[string]interface{}
	StepIndex int
	StartTime time.Time
	Errors    []error
}
