package operations

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// SystemHandler handles system operations
type SystemHandler struct{}

// NewSystemHandler creates a new system handler
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// GetType returns the operation type
func (h *SystemHandler) GetType() string {
	return "system"
}

// Validate checks if the parameters are valid
func (h *SystemHandler) Validate(params map[string]interface{}) error {
	action, ok := params["action"].(string)
	if !ok {
		return fmt.Errorf("action parameter is required")
	}

	switch action {
	case "command", "shell":
		if _, ok := params["command"]; !ok {
			return fmt.Errorf("command parameter is required for %s action", action)
		}
	case "sleep":
		if _, ok := params["duration"]; !ok {
			return fmt.Errorf("duration parameter is required for sleep action")
		}
	case "open":
		if _, ok := params["path"]; !ok {
			return fmt.Errorf("path parameter is required for open action")
		}
	case "env":
		if _, ok := params["variable"]; !ok {
			return fmt.Errorf("variable parameter is required for env action")
		}
	}

	return nil
}

// Execute performs the system operation
func (h *SystemHandler) Execute(ctx context.Context, params map[string]interface{}, execCtx interface{}) (map[string]interface{}, error) {
	action := params["action"].(string)

	var result map[string]interface{}
	var err error

	switch action {
	case "command":
		result, err = h.executeCommand(ctx, params)
	case "shell":
		result, err = h.executeShell(ctx, params)
	case "sleep":
		result, err = h.executeSleep(ctx, params)
	case "open":
		result, err = h.executeOpen(ctx, params)
	case "env":
		result, err = h.executeEnv(ctx, params)
	case "kill":
		result, err = h.executeKill(ctx, params)
	default:
		return nil, fmt.Errorf("unsupported system action: %s", action)
	}

	return result, err
}

func (h *SystemHandler) executeCommand(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	command := getStringParam(params, "command")
	workingDir := getStringParam(params, "working_dir")

	// Parse command and arguments
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	if workingDir != "" {
		cmd.Dir = workingDir
	}

	// Set environment variables if provided
	if env, ok := params["env"].(map[string]interface{}); ok {
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%v", k, v))
		}
	}

	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	if err != nil {
		return map[string]interface{}{
			"success":  false,
			"command":  command,
			"output":   string(output),
			"error":    err.Error(),
			"duration": duration.Milliseconds(),
		}, nil
	}

	return map[string]interface{}{
		"success":  true,
		"command":  command,
		"output":   string(output),
		"duration": duration.Milliseconds(),
	}, nil
}

func (h *SystemHandler) executeShell(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	command := getStringParam(params, "command")
	workingDir := getStringParam(params, "working_dir")

	var cmd *exec.Cmd
	if workingDir != "" {
		cmd = exec.CommandContext(ctx, "sh", "-c", "cd "+workingDir+" && "+command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	// Set environment variables if provided
	if env, ok := params["env"].(map[string]interface{}); ok {
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%v", k, v))
		}
	}

	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	if err != nil {
		return map[string]interface{}{
			"success":  false,
			"command":  command,
			"output":   string(output),
			"error":    err.Error(),
			"duration": duration.Milliseconds(),
		}, nil
	}

	return map[string]interface{}{
		"success":  true,
		"command":  command,
		"output":   string(output),
		"duration": duration.Milliseconds(),
	}, nil
}

func (h *SystemHandler) executeSleep(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	duration := getIntParam(params, "duration")

	select {
	case <-time.After(time.Duration(duration) * time.Second):
		return map[string]interface{}{
			"success":  true,
			"duration": duration,
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (h *SystemHandler) executeOpen(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")

	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".app"):
		// macOS application
		cmd = exec.Command("open", "-a", path)
	case strings.Contains(path, "://"):
		// URL
		cmd = exec.Command("open", path)
	default:
		// File or directory
		cmd = exec.Command("open", path)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("open failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"path":    path,
	}, nil
}

func (h *SystemHandler) executeEnv(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	variable := getStringParam(params, "variable")
	value := getStringParam(params, "value")

	if value == "" {
		// Get environment variable
		envValue := os.Getenv(variable)
		return map[string]interface{}{
			"success":  true,
			"variable": variable,
			"value":    envValue,
		}, nil
	} else {
		// Set environment variable (only for current process)
		os.Setenv(variable, value)
		return map[string]interface{}{
			"success":  true,
			"variable": variable,
			"value":    value,
		}, nil
	}
}

func (h *SystemHandler) executeKill(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	processName := getStringParam(params, "process_name")
	pid := getIntParam(params, "pid")

	var err error
	if pid > 0 {
		// Kill by PID
		process, errFind := os.FindProcess(pid)
		if errFind != nil {
			return nil, fmt.Errorf("find process failed: %v", errFind)
		}
		err = process.Kill()
	} else if processName != "" {
		// Kill by process name (macOS/Linux)
		cmd := exec.Command("pkill", processName)
		output, errKill := cmd.CombinedOutput()
		if errKill != nil {
			return nil, fmt.Errorf("kill process failed: %v, output: %s", errKill, string(output))
		}
	} else {
		return nil, fmt.Errorf("either process_name or pid must be specified")
	}

	if err != nil {
		return nil, fmt.Errorf("kill failed: %v", err)
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}
