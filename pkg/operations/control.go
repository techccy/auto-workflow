package operations

import (
	"auto-workflow/pkg/models"
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ControlHandler handles control flow operations
type ControlHandler struct{}

// NewControlHandler creates a new control handler
func NewControlHandler() *ControlHandler {
	return &ControlHandler{}
}

// GetType returns the operation type
func (h *ControlHandler) GetType() string {
	return "control"
}

// Validate checks if the parameters are valid
func (h *ControlHandler) Validate(params map[string]interface{}) error {
	action, ok := params["action"].(string)
	if !ok {
		return fmt.Errorf("action parameter is required")
	}

	switch action {
	case "set_variable":
		if _, ok := params["name"]; !ok {
			return fmt.Errorf("name parameter is required for set_variable action")
		}
		if _, ok := params["value"]; !ok {
			return fmt.Errorf("value parameter is required for set_variable action")
		}
	case "get_variable":
		if _, ok := params["name"]; !ok {
			return fmt.Errorf("name parameter is required for get_variable action")
		}
	case "increment":
		if _, ok := params["name"]; !ok {
			return fmt.Errorf("name parameter is required for increment action")
		}
	case "decrement":
		if _, ok := params["name"]; !ok {
			return fmt.Errorf("name parameter is required for decrement action")
		}
	case "append":
		if _, ok := params["name"]; !ok {
			return fmt.Errorf("name parameter is required for append action")
		}
		if _, ok := params["value"]; !ok {
			return fmt.Errorf("value parameter is required for append action")
		}
	case "log":
		if _, ok := params["message"]; !ok {
			return fmt.Errorf("message parameter is required for log action")
		}
	case "assert":
		if _, ok := params["condition"]; !ok {
			return fmt.Errorf("condition parameter is required for assert action")
		}
	}

	return nil
}

// Execute performs the control operation
func (h *ControlHandler) Execute(ctx context.Context, params map[string]interface{}, execCtx interface{}) (map[string]interface{}, error) {
	action := params["action"].(string)

	// Get execution context
	context, ok := execCtx.(*models.ExecutionContext)
	if !ok {
		return nil, fmt.Errorf("invalid execution context")
	}

	var result map[string]interface{}
	var err error

	switch action {
	case "set_variable":
		result, err = h.executeSetVariable(ctx, params, context)
	case "get_variable":
		result, err = h.executeGetVariable(ctx, params, context)
	case "increment":
		result, err = h.executeIncrement(ctx, params, context)
	case "decrement":
		result, err = h.executeDecrement(ctx, params, context)
	case "append":
		result, err = h.executeAppend(ctx, params, context)
	case "log":
		result, err = h.executeLog(ctx, params, context)
	case "assert":
		result, err = h.executeAssert(ctx, params, context)
	default:
		return nil, fmt.Errorf("unsupported control action: %s", action)
	}

	return result, err
}

func (h *ControlHandler) executeSetVariable(ctx context.Context, params map[string]interface{}, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	name := getStringParam(params, "name")
	value := params["value"]

	// Resolve variable references if needed
	resolvedValue, err := h.resolveValue(value, execCtx)
	if err != nil {
		return nil, fmt.Errorf("resolve value failed: %v", err)
	}

	if execCtx.Variables == nil {
		execCtx.Variables = make(map[string]interface{})
	}

	execCtx.Variables[name] = resolvedValue

	return map[string]interface{}{
		"success": true,
		"name":    name,
		"value":   resolvedValue,
	}, nil
}

func (h *ControlHandler) executeGetVariable(ctx context.Context, params map[string]interface{}, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	name := getStringParam(params, "name")

	value, exists := execCtx.Variables[name]
	if !exists {
		return nil, fmt.Errorf("variable '%s' not found", name)
	}

	return map[string]interface{}{
		"success": true,
		"name":    name,
		"value":   value,
	}, nil
}

func (h *ControlHandler) executeIncrement(ctx context.Context, params map[string]interface{}, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	name := getStringParam(params, "name")
	amount := getIntParamDefault(params, "amount", 1)

	value, exists := execCtx.Variables[name]
	if !exists {
		return nil, fmt.Errorf("variable '%s' not found", name)
	}

	var newValue int
	switch v := value.(type) {
	case float64:
		newValue = int(v) + amount
	case int:
		newValue = v + amount
	default:
		return nil, fmt.Errorf("variable '%s' is not a number", name)
	}

	execCtx.Variables[name] = newValue

	return map[string]interface{}{
		"success": true,
		"name":    name,
		"value":   newValue,
	}, nil
}

func (h *ControlHandler) executeDecrement(ctx context.Context, params map[string]interface{}, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	name := getStringParam(params, "name")
	amount := getIntParamDefault(params, "amount", 1)

	value, exists := execCtx.Variables[name]
	if !exists {
		return nil, fmt.Errorf("variable '%s' not found", name)
	}

	var newValue int
	switch v := value.(type) {
	case float64:
		newValue = int(v) - amount
	case int:
		newValue = v - amount
	default:
		return nil, fmt.Errorf("variable '%s' is not a number", name)
	}

	execCtx.Variables[name] = newValue

	return map[string]interface{}{
		"success": true,
		"name":    name,
		"value":   newValue,
	}, nil
}

func (h *ControlHandler) executeAppend(ctx context.Context, params map[string]interface{}, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	name := getStringParam(params, "name")
	value := params["value"]

	// Resolve variable references if needed
	resolvedValue, err := h.resolveValue(value, execCtx)
	if err != nil {
		return nil, fmt.Errorf("resolve value failed: %v", err)
	}

	existingValue, exists := execCtx.Variables[name]
	if !exists {
		return nil, fmt.Errorf("variable '%s' not found", name)
	}

	switch arr := existingValue.(type) {
	case []interface{}:
		execCtx.Variables[name] = append(arr, resolvedValue)
	default:
		return nil, fmt.Errorf("variable '%s' is not an array", name)
	}

	return map[string]interface{}{
		"success":  true,
		"name":     name,
		"appended": resolvedValue,
	}, nil
}

func (h *ControlHandler) executeLog(ctx context.Context, params map[string]interface{}, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	message := getStringParam(params, "message")
	level := getStringParamDefault(params, "level", "info")

	// Resolve variable references in message
	resolvedMessage, err := h.resolveVariables(message, execCtx)
	if err != nil {
		return nil, fmt.Errorf("resolve variables failed: %v", err)
	}

	fmt.Printf("[%s] %s\n", strings.ToUpper(level), resolvedMessage)

	return map[string]interface{}{
		"success": true,
		"level":   level,
		"message": resolvedMessage,
	}, nil
}

func (h *ControlHandler) executeAssert(ctx context.Context, params map[string]interface{}, execCtx *models.ExecutionContext) (map[string]interface{}, error) {
	condition := params["condition"]

	// Evaluate the condition
	result, err := h.evaluateCondition(condition, execCtx)
	if err != nil {
		return nil, fmt.Errorf("evaluate condition failed: %v", err)
	}

	if !result {
		message := getStringParamDefault(params, "message", "assertion failed")
		return nil, fmt.Errorf("assertion failed: %s", message)
	}

	return map[string]interface{}{
		"success":   true,
		"condition": condition,
	}, nil
}

// Helper functions

func (h *ControlHandler) resolveValue(value interface{}, execCtx *models.ExecutionContext) (interface{}, error) {
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
			varName := v[2 : len(v)-1]
			if val, exists := execCtx.Variables[varName]; exists {
				return val, nil
			}
			return nil, fmt.Errorf("variable '%s' not found", varName)
		}
		return v, nil
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			resolved, err := h.resolveValue(val, execCtx)
			if err != nil {
				return nil, err
			}
			result[k] = resolved
		}
		return result, nil
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			resolved, err := h.resolveValue(val, execCtx)
			if err != nil {
				return nil, err
			}
			result[i] = resolved
		}
		return result, nil
	default:
		return v, nil
	}
}

func (h *ControlHandler) resolveVariables(text string, execCtx *models.ExecutionContext) (string, error) {
	result := text

	// Find all variable references
	start := 0
	for {
		startIdx := strings.Index(result[start:], "${")
		if startIdx == -1 {
			break
		}
		startIdx += start

		endIdx := strings.Index(result[startIdx:], "}")
		if endIdx == -1 {
			break
		}
		endIdx += startIdx

		varName := result[startIdx+2 : endIdx]
		if val, exists := execCtx.Variables[varName]; exists {
			varValue := fmt.Sprintf("%v", val)
			result = result[:startIdx] + varValue + result[endIdx+1:]
			start = startIdx + len(varValue)
		} else {
			start = endIdx + 1
		}
	}

	return result, nil
}

func (h *ControlHandler) evaluateCondition(condition interface{}, execCtx *models.ExecutionContext) (bool, error) {
	switch cond := condition.(type) {
	case bool:
		return cond, nil
	case string:
		// Try to resolve variable reference
		if strings.HasPrefix(cond, "${") && strings.HasSuffix(cond, "}") {
			varName := cond[2 : len(cond)-1]
			if val, exists := execCtx.Variables[varName]; exists {
				// For non-boolean values, return true if the value exists
				// The actual comparison should be done in complex conditions
				if b, ok := val.(bool); ok {
					return b, nil
				}
				// For numeric values, check if it's truthy (non-zero)
				if num, ok := toFloat64(val); ok {
					return num != 0, nil
				}
				// For strings, check if it's not empty
				if str, ok := val.(string); ok {
					return str != "", nil
				}
				// For other types, return true if it exists
				return true, nil
			}
			return false, fmt.Errorf("variable '%s' not found", varName)
		}
		// Try to parse as boolean
		switch strings.ToLower(cond) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
		return false, fmt.Errorf("invalid condition: %s", cond)
	case map[string]interface{}:
		// Complex condition with operator
		return h.evaluateComplexCondition(cond, execCtx)
	case float64, int, int32, int64:
		// Numeric values: true if non-zero
		return cond != 0, nil
	default:
		return false, fmt.Errorf("invalid condition type: %v", reflect.TypeOf(condition))
	}
}

func (h *ControlHandler) evaluateComplexCondition(condition map[string]interface{}, execCtx *models.ExecutionContext) (bool, error) {
	operator, ok := condition["operator"].(string)
	if !ok {
		return false, fmt.Errorf("operator is required for complex condition")
	}

	left, err := h.resolveValue(condition["left"], execCtx)
	if err != nil {
		return false, err
	}

	right, err := h.resolveValue(condition["right"], execCtx)
	if err != nil {
		return false, err
	}

	switch operator {
	case "equals":
		// For numeric values, compare after type conversion
		if leftNum, leftOk := toFloat64(left); leftOk {
			if rightNum, rightOk := toFloat64(right); rightOk {
				return leftNum == rightNum, nil
			}
		}
		return reflect.DeepEqual(left, right), nil
	case "not_equals":
		// For numeric values, compare after type conversion
		if leftNum, leftOk := toFloat64(left); leftOk {
			if rightNum, rightOk := toFloat64(right); rightOk {
				return leftNum != rightNum, nil
			}
		}
		return !reflect.DeepEqual(left, right), nil
	case "greater":
		return compareNumbers(left, right) > 0, nil
	case "less":
		return compareNumbers(left, right) < 0, nil
	case "greater_equal":
		return compareNumbers(left, right) >= 0, nil
	case "less_equal":
		return compareNumbers(left, right) <= 0, nil
	case "contains":
		return containsValue(left, right), nil
	case "and":
		leftBool, err := h.evaluateCondition(left, execCtx)
		if err != nil || !leftBool {
			return false, err
		}
		return h.evaluateCondition(right, execCtx)
	case "or":
		leftBool, err := h.evaluateCondition(left, execCtx)
		if err == nil && leftBool {
			return true, nil
		}
		return h.evaluateCondition(right, execCtx)
	case "not":
		result, err := h.evaluateCondition(left, execCtx)
		return !result, err
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func compareNumbers(a, b interface{}) int {
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
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func containsValue(container, value interface{}) bool {
	switch c := container.(type) {
	case string:
		return strings.Contains(c, fmt.Sprintf("%v", value))
	case []interface{}:
		for _, item := range c {
			if reflect.DeepEqual(item, value) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Helper functions for parameter extraction
func getStringParamDefault(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}
