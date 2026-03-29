# Step Result Transitivity

## Overview

The auto-workflow system now supports transitivity, allowing you to use results from previous steps as variables in subsequent steps. This enables powerful data flow between workflow steps.

## Syntax

### Referencing Step Results

You can reference results from previous steps using the following syntax:

```
${step_id}
```

This references the entire result object from the step with ID `step_id`.

### Referencing Specific Fields

To reference a specific field from a step's result:

```
${step_id.field_name}
```

This accesses the `field_name` property from the result object of step `step_id`.

### Workflow Variables

Workflow-level variables continue to work as before:

```
${variable_name}
```

## Usage Examples

### Example 1: Web Scraping and Saving Content

```json
{
  "steps": [
    {
      "id": "fetch_page",
      "name": "Fetch Web Page",
      "type": "web",
      "parameters": {
        "action": "get",
        "url": "https://example.com"
      }
    },
    {
      "id": "save_content",
      "name": "Save HTML Content",
      "type": "file",
      "parameters": {
        "action": "write",
        "path": "./output/page.html",
        "content": "${fetch_page.body}"
      }
    }
  ]
}
```

### Example 2: Using Multiple Fields

```json
{
  "steps": [
    {
      "id": "fetch_data",
      "name": "Fetch API Data",
      "type": "web",
      "parameters": {
        "action": "get",
        "url": "https://api.example.com/data"
      }
    },
    {
      "id": "log_status",
      "name": "Log Response Status",
      "type": "file",
      "parameters": {
        "action": "write",
        "path": "./output/status.txt",
        "content": "Status Code: ${fetch_data.status_code}\nSuccess: ${fetch_data.success}\nDuration: ${fetch_data.duration}ms"
      }
    }
  ]
}
```

### Example 3: Conditional Execution Based on Step Results

```json
{
  "steps": [
    {
      "id": "check_api",
      "name": "Check API Status",
      "type": "web",
      "parameters": {
        "action": "get",
        "url": "https://api.example.com/health"
      }
    },
    {
      "id": "process_data",
      "name": "Process Data",
      "type": "file",
      "parameters": {
        "action": "write",
        "path": "./output/data.txt",
        "content": "API is healthy"
      },
      "conditions": [
        {
          "variable": "check_api.success",
          "operator": "equals",
          "value": "True"
        }
      ]
    }
  ]
}
```

### Example 4: Combining Workflow Variables and Step Results

```json
{
  "variables": [
    {
      "name": "output_dir",
      "value": "./output",
      "type": "string"
    }
  ],
  "steps": [
    {
      "id": "fetch_image",
      "name": "Fetch Image",
      "type": "web",
      "parameters": {
        "action": "get",
        "url": "https://example.com/image.png"
      }
    },
    {
      "id": "save_image",
      "name": "Save Image",
      "type": "file",
      "parameters": {
        "action": "write",
        "path": "${output_dir}/image.png",
        "content": "${fetch_image.body}"
      }
    }
  ]
}
```

## Available Fields by Step Type

### Web Steps

When you execute a web operation (get, post, put, patch, delete), the result contains:

- `success`: Boolean indicating if the request was successful (2xx status)
- `status_code`: HTTP status code (e.g., 200, 404, 500)
- `body`: Response body as text
- `headers`: Dictionary of response headers
- `duration`: Request duration in milliseconds

Example:
```json
{
  "id": "web_request",
  "name": "Make Request",
  "type": "web",
  "parameters": {
    "action": "get",
    "url": "https://api.example.com"
  }
}
```

Usage:
```
${web_request.success}        // true/false
${web_request.status_code}    // 200
${web_request.body}           // Response content
${web_request.headers}        // Headers dict
${web_request.duration}       // 1234
```

### File Steps

File operations return various fields depending on the action:

- `mkdir`: `{ "success": true, "path": "/path/to/dir" }`
- `write`: `{ "success": true, "path": "/path/to/file", "bytes_written": 1234 }`
- `read`: `{ "success": true, "content": "file content", "path": "/path/to/file" }`

### ADB Steps

ADB operations return:

- `screencap`: `{ "success": true, "output_path": "/path/to/screenshot.png" }`
- `install`: `{ "success": true, "package": "com.example.app" }`
- `uninstall`: `{ "success": true, "package": "com.example.app" }`

### System Steps

System operations return:

- `execute`: `{ "success": true, "exit_code": 0, "stdout": "output", "stderr": "" }`

## Best Practices

1. **Use Descriptive Step IDs**: Choose clear, descriptive IDs for steps to make references more readable.
   ```json
   { "id": "fetch_user_profile", ... }  // Good
   { "id": "step1", ... }               // Avoid
   ```

2. **Order Steps Correctly**: Ensure that steps are executed in the correct order. A step can only reference results from previous steps.

3. **Handle Missing Fields**: If you reference a field that doesn't exist in the step result, it will be replaced with an empty string.

4. **Use Conditions**: Leverage conditional execution to create more robust workflows that handle different scenarios.

5. **Combine with Loops**: Step result references work seamlessly with loop configurations.

## Limitations

- Steps can only reference results from **previous** steps in the workflow execution order.
- Circular references are not supported.
- If a step fails, its result will not be available to subsequent steps.

## Complete Example

See [`examples/transitivity_example.json`](../examples/transitivity_example.json) for a comprehensive example demonstrating all aspects of step result transitivity.

## Troubleshooting

### Variable Not Found

If you see a variable reference in your output (e.g., `${step_id.field}`), it means:
- The step hasn't been executed yet
- The step failed
- The field doesn't exist in the step result

### Step Execution Order

Ensure steps are ordered correctly. A step cannot reference results from steps that come after it in the workflow.

### Type Conversion

All values are converted to strings when used in parameters. If you need to preserve types for conditions, the system handles type conversion automatically.
