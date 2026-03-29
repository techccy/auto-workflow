package operations

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileHandler handles file system operations
type FileHandler struct{}

// NewFileHandler creates a new file handler
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// GetType returns the operation type
func (h *FileHandler) GetType() string {
	return "file"
}

// Validate checks if the parameters are valid
func (h *FileHandler) Validate(params map[string]interface{}) error {
	action, ok := params["action"].(string)
	if !ok {
		return fmt.Errorf("action parameter is required")
	}

	switch action {
	case "read", "delete", "exists":
		if _, ok := params["path"]; !ok {
			return fmt.Errorf("path parameter is required for %s action", action)
		}
	case "write", "append":
		if _, ok := params["path"]; !ok {
			return fmt.Errorf("path parameter is required for %s action", action)
		}
		if _, ok := params["content"]; !ok {
			return fmt.Errorf("content parameter is required for %s action", action)
		}
	case "copy":
		if _, ok := params["source"]; !ok {
			return fmt.Errorf("source parameter is required for copy action")
		}
		if _, ok := params["destination"]; !ok {
			return fmt.Errorf("destination parameter is required for copy action")
		}
	case "move", "rename":
		if _, ok := params["source"]; !ok {
			return fmt.Errorf("source parameter is required for %s action", action)
		}
		if _, ok := params["destination"]; !ok {
			return fmt.Errorf("destination parameter is required for %s action", action)
		}
	case "mkdir":
		if _, ok := params["path"]; !ok {
			return fmt.Errorf("path parameter is required for mkdir action")
		}
	case "list":
		if _, ok := params["path"]; !ok {
			return fmt.Errorf("path parameter is required for list action")
		}
	}

	return nil
}

// Execute performs the file operation
func (h *FileHandler) Execute(ctx context.Context, params map[string]interface{}, execCtx interface{}) (map[string]interface{}, error) {
	action := params["action"].(string)

	var result map[string]interface{}
	var err error

	switch action {
	case "read":
		result, err = h.executeRead(ctx, params)
	case "write":
		result, err = h.executeWrite(ctx, params)
	case "append":
		result, err = h.executeAppend(ctx, params)
	case "delete":
		result, err = h.executeDelete(ctx, params)
	case "copy":
		result, err = h.executeCopy(ctx, params)
	case "move":
		result, err = h.executeMove(ctx, params)
	case "rename":
		result, err = h.executeRename(ctx, params)
	case "exists":
		result, err = h.executeExists(ctx, params)
	case "mkdir":
		result, err = h.executeMkdir(ctx, params)
	case "list":
		result, err = h.executeList(ctx, params)
	default:
		return nil, fmt.Errorf("unsupported file action: %s", action)
	}

	return result, err
}

func (h *FileHandler) executeRead(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"path":    path,
		"content": string(content),
		"size":    len(content),
	}, nil
}

func (h *FileHandler) executeWrite(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")
	content := getStringParam(params, "content")

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create directory failed: %v", err)
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return nil, fmt.Errorf("write failed: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"path":    path,
		"size":    len(content),
	}, nil
}

func (h *FileHandler) executeAppend(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")
	content := getStringParam(params, "content")

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return nil, fmt.Errorf("append failed: %v", err)
	}

	return map[string]interface{}{
		"success":  true,
		"path":     path,
		"appended": len(content),
	}, nil
}

func (h *FileHandler) executeDelete(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")

	err := os.RemoveAll(path)
	if err != nil {
		return nil, fmt.Errorf("delete failed: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"path":    path,
	}, nil
}

func (h *FileHandler) executeCopy(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	source := getStringParam(params, "source")
	destination := getStringParam(params, "destination")

	// Create destination directory if it doesn't exist
	dir := filepath.Dir(destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create destination directory failed: %v", err)
	}

	// Check if source is a directory
	info, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("stat source failed: %v", err)
	}

	if info.IsDir() {
		err = copyDir(source, destination)
	} else {
		err = copyFile(source, destination)
	}

	if err != nil {
		return nil, fmt.Errorf("copy failed: %v", err)
	}

	return map[string]interface{}{
		"success":     true,
		"source":      source,
		"destination": destination,
	}, nil
}

func (h *FileHandler) executeMove(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	source := getStringParam(params, "source")
	destination := getStringParam(params, "destination")

	// Create destination directory if it doesn't exist
	dir := filepath.Dir(destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create destination directory failed: %v", err)
	}

	err := os.Rename(source, destination)
	if err != nil {
		return nil, fmt.Errorf("move failed: %v", err)
	}

	return map[string]interface{}{
		"success":     true,
		"source":      source,
		"destination": destination,
	}, nil
}

func (h *FileHandler) executeRename(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	source := getStringParam(params, "source")
	destination := getStringParam(params, "destination")

	err := os.Rename(source, destination)
	if err != nil {
		return nil, fmt.Errorf("rename failed: %v", err)
	}

	return map[string]interface{}{
		"success":     true,
		"source":      source,
		"destination": destination,
	}, nil
}

func (h *FileHandler) executeExists(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")

	_, err := os.Stat(path)
	exists := err == nil

	return map[string]interface{}{
		"success": true,
		"path":    path,
		"exists":  exists,
	}, nil
}

func (h *FileHandler) executeMkdir(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")
	recursive := getBoolParamDefault(params, "recursive", true)

	var err error
	if recursive {
		err = os.MkdirAll(path, 0755)
	} else {
		err = os.Mkdir(path, 0755)
	}

	if err != nil {
		return nil, fmt.Errorf("mkdir failed: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"path":    path,
	}, nil
}

func (h *FileHandler) executeList(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path := getStringParam(params, "path")

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("list failed: %v", err)
	}

	files := []map[string]interface{}{}
	for _, entry := range entries {
		info, _ := entry.Info()
		files = append(files, map[string]interface{}{
			"name":   entry.Name(),
			"is_dir": entry.IsDir(),
			"size":   info.Size(),
		})
	}

	return map[string]interface{}{
		"success": true,
		"path":    path,
		"files":   files,
		"count":   len(files),
	}, nil
}

// Helper functions
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

func copyDir(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, sourceInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
		} else {
			err = copyFile(srcPath, dstPath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func getBoolParamDefault(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}
