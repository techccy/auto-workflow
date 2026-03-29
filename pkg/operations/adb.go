package operations

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ADBHandler handles ADB operations
type ADBHandler struct {
	deviceID string
}

// NewADBHandler creates a new ADB handler
func NewADBHandler() *ADBHandler {
	return &ADBHandler{}
}

// SetDevice sets the target device ID
func (h *ADBHandler) SetDevice(deviceID string) {
	h.deviceID = deviceID
}

// GetType returns the operation type
func (h *ADBHandler) GetType() string {
	return "adb"
}

// Validate checks if the parameters are valid
func (h *ADBHandler) Validate(params map[string]interface{}) error {
	action, ok := params["action"].(string)
	if !ok {
		return fmt.Errorf("action parameter is required")
	}

	switch action {
	case "tap":
		if _, ok := params["x"]; !ok {
			return fmt.Errorf("x parameter is required for tap action")
		}
		if _, ok := params["y"]; !ok {
			return fmt.Errorf("y parameter is required for tap action")
		}
	case "swipe":
		if _, ok := params["x1"]; !ok {
			return fmt.Errorf("x1 parameter is required for swipe action")
		}
		if _, ok := params["y1"]; !ok {
			return fmt.Errorf("y1 parameter is required for swipe action")
		}
		if _, ok := params["x2"]; !ok {
			return fmt.Errorf("x2 parameter is required for swipe action")
		}
		if _, ok := params["y2"]; !ok {
			return fmt.Errorf("y2 parameter is required for swipe action")
		}
	case "input":
		if _, ok := params["text"]; !ok {
			return fmt.Errorf("text parameter is required for input action")
		}
	case "install":
		if _, ok := params["apk_path"]; !ok {
			return fmt.Errorf("apk_path parameter is required for install action")
		}
	case "uninstall":
		if _, ok := params["package"]; !ok {
			return fmt.Errorf("package parameter is required for uninstall action")
		}
	case "screencap":
		if _, ok := params["output_path"]; !ok {
			return fmt.Errorf("output_path parameter is required for screencap action")
		}
	case "shell":
		if _, ok := params["command"]; !ok {
			return fmt.Errorf("command parameter is required for shell action")
		}
	}

	return nil
}

// Execute performs the ADB operation
func (h *ADBHandler) Execute(ctx context.Context, params map[string]interface{}, execCtx interface{}) (map[string]interface{}, error) {
	action := params["action"].(string)

	// Build base command
	args := []string{}
	if h.deviceID != "" {
		args = append(args, "-s", h.deviceID)
	}

	var result map[string]interface{}
	var err error

	switch action {
	case "tap":
		result, err = h.executeTap(ctx, args, params)
	case "swipe":
		result, err = h.executeSwipe(ctx, args, params)
	case "input":
		result, err = h.executeInput(ctx, args, params)
	case "install":
		result, err = h.executeInstall(ctx, args, params)
	case "uninstall":
		result, err = h.executeUninstall(ctx, args, params)
	case "screencap":
		result, err = h.executeScreencap(ctx, args, params)
	case "shell":
		result, err = h.executeShell(ctx, args, params)
	case "click":
		result, err = h.executeClick(ctx, args, params)
	case "press":
		result, err = h.executePress(ctx, args, params)
	default:
		return nil, fmt.Errorf("unsupported ADB action: %s", action)
	}

	return result, err
}

func (h *ADBHandler) executeTap(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	x := getIntParam(params, "x")
	y := getIntParam(params, "y")

	args := append(baseArgs, "shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y))
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("tap failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"x":       x,
		"y":       y,
	}, nil
}

func (h *ADBHandler) executeSwipe(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	x1 := getIntParam(params, "x1")
	y1 := getIntParam(params, "y1")
	x2 := getIntParam(params, "x2")
	y2 := getIntParam(params, "y2")
	duration := getIntParamDefault(params, "duration", 300)

	args := append(baseArgs, "shell", "input", "swipe",
		strconv.Itoa(x1), strconv.Itoa(y1),
		strconv.Itoa(x2), strconv.Itoa(y2),
		strconv.Itoa(duration))
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("swipe failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success":  true,
		"from":     map[string]int{"x": x1, "y": y1},
		"to":       map[string]int{"x": x2, "y": y2},
		"duration": duration,
	}, nil
}

func (h *ADBHandler) executeInput(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	text := getStringParam(params, "text")

	args := append(baseArgs, "shell", "input", "text", text)
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("input failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"text":    text,
	}, nil
}

func (h *ADBHandler) executeInstall(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	apkPath := getStringParam(params, "apk_path")

	args := append(baseArgs, "install", apkPath)
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("install failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"apk":     apkPath,
		"output":  string(output),
	}, nil
}

func (h *ADBHandler) executeUninstall(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	pkg := getStringParam(params, "package")

	args := append(baseArgs, "uninstall", pkg)
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("uninstall failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"package": pkg,
		"output":  string(output),
	}, nil
}

func (h *ADBHandler) executeScreencap(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	outputPath := getStringParam(params, "output_path")

	args := append(baseArgs, "shell", "screencap", "-p", "/sdcard/screenshot.png")
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("screencap failed: %v, output: %s", err, string(output))
	}

	// Pull the screenshot
	pullArgs := append(baseArgs, "pull", "/sdcard/screenshot.png", outputPath)
	pullCmd := exec.CommandContext(ctx, "adb", pullArgs...)

	pullOutput, err := pullCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("pull screenshot failed: %v, output: %s", err, string(pullOutput))
	}

	return map[string]interface{}{
		"success":     true,
		"output_path": outputPath,
	}, nil
}

func (h *ADBHandler) executeShell(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	command := getStringParam(params, "command")

	args := append(baseArgs, "shell", command)
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("shell command failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"command": command,
		"output":  string(output),
	}, nil
}

func (h *ADBHandler) executeClick(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	x := getIntParam(params, "x")
	y := getIntParam(params, "y")

	args := append(baseArgs, "shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y))
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("click failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"x":       x,
		"y":       y,
	}, nil
}

func (h *ADBHandler) executePress(ctx context.Context, baseArgs []string, params map[string]interface{}) (map[string]interface{}, error) {
	key := getStringParam(params, "key")

	args := append(baseArgs, "shell", "input", "keyevent", key)
	cmd := exec.CommandContext(ctx, "adb", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("press failed: %v, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"success": true,
		"key":     key,
	}, nil
}

// Helper functions
func getIntParam(params map[string]interface{}, key string) int {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

func getIntParamDefault(params map[string]interface{}, key string, defaultValue int) int {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return defaultValue
}

func getStringParam(params map[string]interface{}, key string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WaitForDevice waits for device to be connected
func (h *ADBHandler) WaitForDevice(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []string{}
	if h.deviceID != "" {
		args = append(args, "-s", h.deviceID)
	}
	args = append(args, "wait-for-device")

	cmd := exec.CommandContext(ctx, "adb", args...)
	return cmd.Run()
}

// GetConnectedDevices returns list of connected devices
func (h *ADBHandler) GetConnectedDevices() ([]string, error) {
	cmd := exec.Command("adb", "devices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	devices := []string{}

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "daemon") {
			parts := strings.Fields(line)
			if len(parts) >= 2 && parts[1] == "device" {
				devices = append(devices, parts[0])
			}
		}
	}

	return devices, nil
}
