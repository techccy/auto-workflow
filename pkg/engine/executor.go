package engine

import (
	"auto-workflow/pkg/models"
	"auto-workflow/pkg/operations"
	"context"
	"fmt"
	"sync"
	"time"
)

// Executor handles workflow execution
type Executor struct {
	registry   *operations.HandlerRegistry
	adbHandler *operations.ADBHandler
}

// NewExecutor creates a new executor
func NewExecutor() *Executor {
	registry := operations.NewHandlerRegistry()

	// Register all handlers
	registry.Register(operations.NewADBHandler())
	registry.Register(operations.NewFileHandler())
	registry.Register(operations.NewSystemHandler())
	registry.Register(operations.NewWebHandler())
	registry.Register(operations.NewControlHandler())

	executor := &Executor{
		registry: registry,
	}

	// Get reference to ADB handler for device management
	if handler, ok := registry.Get("adb"); ok {
		executor.adbHandler = handler.(*operations.ADBHandler)
	}

	return executor
}

// SetADBDevice sets the target ADB device
func (e *Executor) SetADBDevice(deviceID string) {
	if e.adbHandler != nil {
		e.adbHandler.SetDevice(deviceID)
	}
}

// Execute executes a workflow
func (e *Executor) Execute(workflow *models.Workflow) (*models.ExecutionContext, error) {
	ctx := context.Background()
	return e.ExecuteWithContext(ctx, workflow)
}

// ExecuteWithContext executes a workflow with context
func (e *Executor) ExecuteWithContext(ctx context.Context, workflow *models.Workflow) (*models.ExecutionContext, error) {
	// Initialize execution context
	execCtx := &models.ExecutionContext{
		Variables: make(map[string]interface{}),
		Results:   make(map[string]interface{}),
		StartTime: time.Now(),
		Errors:    make([]error, 0),
	}

	// Initialize workflow variables
	for _, variable := range workflow.Variables {
		execCtx.Variables[variable.Name] = variable.Value
	}

	fmt.Printf("Starting workflow: %s (version: %s)\n", workflow.Name, workflow.Version)
	fmt.Printf("Description: %s\n", workflow.Description)
	fmt.Println("==========================================")

	// Execute steps
	for i, step := range workflow.Steps {
		execCtx.StepIndex = i

		// Check if step is enabled
		if !step.Enabled {
			fmt.Printf("Step %d (%s): Skipped (disabled)\n", i+1, step.ID)
			continue
		}

		// Check conditions
		if !e.evaluateConditions(step.Conditions, execCtx) {
			fmt.Printf("Step %d (%s): Skipped (conditions not met)\n", i+1, step.ID)
			continue
		}

		// Execute step with retry logic
		result, err := e.executeStepWithRetry(ctx, step, execCtx)

		// Store result
		execCtx.Results[step.ID] = result

		if err != nil {
			execCtx.Errors = append(execCtx.Errors, err)
			fmt.Printf("Step %d (%s): Failed - %v\n", i+1, step.ID, err)

			// Handle error
			if workflow.OnError != nil {
				if workflow.OnError.Action == "stop" {
					return execCtx, fmt.Errorf("workflow stopped due to error: %v", err)
				} else if workflow.OnError.Action == "continue" {
					continue
				}
			} else {
				return execCtx, fmt.Errorf("workflow failed at step %d: %v", i+1, err)
			}
		} else {
			fmt.Printf("Step %d (%s): Success\n", i+1, step.ID)
		}
	}

	duration := time.Since(execCtx.StartTime)
	fmt.Println("==========================================")
	fmt.Printf("Workflow completed in %v\n", duration)
	if len(execCtx.Errors) > 0 {
		fmt.Printf("Errors: %d\n", len(execCtx.Errors))
	}

	return execCtx, nil
}

// executeStepWithRetry executes a step with retry logic
func (e *Executor) executeStepWithRetry(ctx context.Context, step models.Step, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	var lastError error

	// Handle loop if configured
	if step.Loop != nil {
		return e.executeStepWithLoop(ctx, step, execCtx)
	}

	// Execute with retry
	maxRetries := 1
	if step.Retry != nil {
		maxRetries = step.Retry.MaxRetries + 1
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err := e.executeStep(ctx, step, execCtx)
		if err == nil {
			return result, nil
		}

		lastError = err

		if attempt < maxRetries {
			delay := time.Second
			if step.Retry != nil {
				delay = time.Duration(step.Retry.Delay) * time.Second
				if step.Retry.Backoff {
					delay = time.Duration(step.Retry.Delay*attempt) * time.Second
				}
			}

			fmt.Printf("  Retry %d/%d in %v...\n", attempt, maxRetries-1, delay)
			time.Sleep(delay)
		}
	}

	return nil, lastError
}

// executeStepWithLoop executes a step with loop logic
func (e *Executor) executeStepWithLoop(ctx context.Context, step models.Step, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	iterations := 1
	var loopValues []interface{}

	switch step.Loop.Type {
	case "count":
		iterations = step.Loop.Count
	case "for_each":
		loopValues = step.Loop.Values
		iterations = len(loopValues)
	case "while":
		// While loop - check condition each iteration
		iterations = -1 // Infinite loop until condition is false
	}

	results := []map[string]interface{}{}

	for i := 0; iterations == -1 || i < iterations; i++ {
		// Set loop variable if specified
		if step.Loop.Variable != "" {
			if step.Loop.Type == "for_each" && i < len(loopValues) {
				execCtx.Variables[step.Loop.Variable] = loopValues[i]
			} else {
				execCtx.Variables[step.Loop.Variable] = i
			}
		}

		// Check while condition
		if step.Loop.Type == "while" {
			condition := step.Loop.Variable
			if condition == "" {
				break
			}
			// Evaluate condition (simplified - in real implementation would be more robust)
			if val, ok := execCtx.Variables[condition].(bool); ok && !val {
				break
			}
		}

		result, err := e.executeStep(ctx, step, execCtx)
		if err != nil {
			return nil, fmt.Errorf("loop iteration %d failed: %v", i+1, err)
		}

		results = append(results, result)
	}

	return map[string]interface{}{
		"iterations": len(results),
		"results":    results,
	}, nil
}

// executeStep executes a single step
func (e *Executor) executeStep(ctx context.Context, step models.Step, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	// Get handler for step type
	handler, ok := e.registry.Get(step.Type)
	if !ok {
		return nil, fmt.Errorf("no handler found for type: %s", step.Type)
	}

	// Validate parameters
	if err := handler.Validate(step.Parameters); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %v", err)
	}

	// Set timeout if specified
	if step.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(step.Timeout)*time.Second)
		defer cancel()
	}

	// Execute step
	result, err := handler.Execute(ctx, step.Parameters, execCtx)

	// Handle success/failure steps
	if err == nil && len(step.OnSuccess) > 0 {
		// Execute success steps (simplified - would need proper implementation)
		fmt.Printf("  Executing success steps: %v\n", step.OnSuccess)
	} else if err != nil && len(step.OnFailure) > 0 {
		// Execute failure steps (simplified - would need proper implementation)
		fmt.Printf("  Executing failure steps: %v\n", step.OnFailure)
	}

	return result, err
}

// evaluateConditions evaluates step conditions
func (e *Executor) evaluateConditions(conditions []models.Condition, execCtx *models.ExecutionContext) bool {
	if len(conditions) == 0 {
		return true
	}

	for _, condition := range conditions {
		if !e.evaluateCondition(condition, execCtx) {
			return false
		}
	}

	return true
}

// evaluateCondition evaluates a single condition
func (e *Executor) evaluateCondition(condition models.Condition, execCtx *models.ExecutionContext) bool {
	value, exists := execCtx.Variables[condition.Variable]
	if !exists {
		return false
	}

	switch condition.Operator {
	case "equals":
		return fmt.Sprintf("%v", value) == fmt.Sprintf("%v", condition.Value)
	case "not_equals":
		return fmt.Sprintf("%v", value) != fmt.Sprintf("%v", condition.Value)
	case "greater":
		return compareValues(value, condition.Value) > 0
	case "less":
		return compareValues(value, condition.Value) < 0
	case "contains":
		return containsValue(value, condition.Value)
	case "exists":
		return true
	default:
		return false
	}
}

// ExecuteParallel executes steps in parallel
func (e *Executor) ExecuteParallel(ctx context.Context, steps []models.Step, execCtx *models.ExecutionContext) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(steps))

	for _, step := range steps {
		if !step.Parallel {
			continue
		}

		wg.Add(1)
		go func(s models.Step) {
			defer wg.Done()

			_, err := e.executeStepWithRetry(ctx, s, execCtx)
			if err != nil {
				errChan <- fmt.Errorf("step %s failed: %v", s.ID, err)
			}
		}(step)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		return err
	}

	return nil
}

// Helper functions
func compareValues(a, b interface{}) int {
	aVal, aOk := toFloat64(a)
	bVal, bOk := toFloat64(b)

	if !aOk || !bOk {
		return 0
	}

	if aVal < bVal {
		return -1
	} else if aVal > bVal {
		return 1
	}
	return 0
}

func containsValue(container, value interface{}) bool {
	switch c := container.(type) {
	case string:
		return fmt.Sprintf("%v", c) == fmt.Sprintf("%v", value)
	case []interface{}:
		for _, item := range c {
			if fmt.Sprintf("%v", item) == fmt.Sprintf("%v", value) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}
