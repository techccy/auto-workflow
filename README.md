# Auto Workflow

A flexible and powerful automation workflow system written in Go that allows you to define and execute complex automation tasks using JSON configuration files.

## Features

- **Multiple Operation Types**: Support for ADB, file system, system commands, web API, and control flow operations
- **Conditional Logic**: Execute steps based on conditions with support for various operators
- **Loops**: Implement count-based, while, and for-each loops for repetitive tasks
- **Retry Mechanism**: Automatic retry with configurable delay and exponential backoff
- **Parallel Execution**: Run multiple steps concurrently for better performance
- **Variable System**: Define and use variables throughout your workflow
- **Error Handling**: Comprehensive error handling with configurable error actions
- **JSON-based Configuration**: Easy-to-read and maintain workflow definitions

## Installation

### Prerequisites

- Go 1.21 or higher
- ADB (for Android device automation)

### Build

```bash
go build -o auto-workflow cmd/main.go
```

## Usage

### Basic Usage

```bash
# Run a workflow
./auto-workflow -workflow examples/simple_workflow.json

# Run with specific ADB device
./auto-workflow -workflow examples/adb_workflow.json -device <device_id>

# List connected ADB devices
./auto-workflow -list-devices

# Show version
./auto-workflow -version
```

### Command Line Options

- `-workflow <path>`: Path to the workflow JSON file (required)
- `-device <id>`: ADB device ID to use for ADB operations
- `-list-devices`: List all connected ADB devices
- `-version`: Show version information

## Workflow Structure

A workflow is defined in JSON format with the following structure:

```json
{
  "name": "Workflow Name",
  "description": "Workflow description",
  "version": "1.0.0",
  "variables": [
    {
      "name": "variable_name",
      "value": "variable_value",
      "type": "string"
    }
  ],
  "steps": [
    {
      "id": "step_id",
      "name": "Step Name",
      "type": "operation_type",
      "enabled": true,
      "parallel": false,
      "retry": {
        "max_retries": 3,
        "delay": 1,
        "backoff": true
      },
      "timeout": 30,
      "parameters": {
        "action": "specific_action",
        "param1": "value1"
      },
      "conditions": [
        {
          "variable": "var_name",
          "operator": "equals",
          "value": "expected_value"
        }
      ],
      "loop": {
        "type": "count",
        "count": 5,
        "variable": "iteration"
      }
    }
  ],
  "on_error": {
    "action": "stop",
    "notify": true
  }
}
```

For example, to finish the work below:
Step 1: Create a new folder named "test" in /Users/ccy/Downloads.
Step 2: Crawl https://cloud.techccy.dpdns.org and write the content to /Users/ccy/Downloads/test/in.html.
Step 3: Take a screenshot using ADB and save the photo to /Users/ccy/Downloads/test/in.html.

```
{
  "name": "CCY Auto Material Collection",     // Workflow name displayed during execution
  "description": "Create folder, sync web content and mobile screenshot", // Detailed functional description
  "version": "1.0.0",                        // Version number for script iteration management
  "variables": [                              // Global variable definitions to avoid hardcoding paths
    {
      "name": "target_dir",                   // Variable name: target directory
      "value": "/Users/ccy/Downloads/test",  // Variable value: the 'test' folder in your Downloads
      "type": "string"                        // Data type is string
    }
  ],
  "steps": [                                  // List of core steps executed in sequence
    {
      "id": "step1_mkdir",                    // Unique ID for logging and tracking
      "name": "Create Test Folder",           // Human-readable step name
      "type": "file",                         // Specifies the 'File System' module
      "parameters": {                         // Parameters passed to pkg/operations/file.go
        "action": "mkdir",                    // Action: Create directory
        "path": "${target_dir}"               // Uses the variable; auto-replaced at runtime
      }
    },
    {
      "id": "step2_web_get",                  // Step ID
      "name": "Fetch Web Content",            // Step name
      "type": "web",                          // Specifies the 'Network Request' module
      "parameters": {                         // Parameters passed to pkg/operations/web.go
        "action": "get",                      // Action: Execute HTTP GET request
        "url": "https://cloud.techccy.dpdns.org" // URL to fetch
      },
      "retry": {                              // Fault tolerance for network instability
        "max_retries": 2,                     // Automatically retry up to 2 times on failure
        "delay": 5                            // Wait 5 seconds between retries
      }
    },
    {
      "id": "step3_save_html",                // Step ID
      "name": "Write Web Content to File",    // Step name
      "type": "file",                         // Calls the file module again
      "parameters": {
        "action": "write",                    // Action: Write to file
        "path": "${target_dir}/in.html",      // Full path (Variable + Filename)
        "content": "Captured Content"         // Data to be written
      }
    },
    {
      "id": "step4_adb_screenshot",           // Step ID
      "name": "ADB Screenshot Sync",          // Step name
      "type": "adb",                          // Specifies the 'ADB Mobile' module
      "parameters": {                         // Parameters passed to pkg/operations/adb.go
        "action": "screencap",                // Action: Mobile screen capture
        "path": "${target_dir}/screenshot.png" // Path to save the image on the PC
      }
    }
  ],
  "on_error": {                               // Global error handling strategy
    "action": "stop"                          // Stop the workflow immediately if any step fails
  }
}
```
## Operation Types

### ADB Operations

Automate Android device operations using ADB.

**Actions:**
- `tap`: Tap on screen coordinates
- `swipe`: Swipe gesture between points
- `input`: Input text
- `install`: Install APK
- `uninstall`: Uninstall app
- `screencap`: Take screenshot
- `shell`: Execute shell command
- `click`: Click on coordinates (alias for tap)
- `press`: Press key event

**Example:**
```json
{
  "type": "adb",
  "parameters": {
    "action": "tap",
    "x": 500,
    "y": 1000
  }
}
```

### File Operations

Perform file system operations.

**Actions:**
- `read`: Read file contents
- `write`: Write content to file
- `append`: Append content to file
- `delete`: Delete file or directory
- `copy`: Copy file or directory
- `move`: Move file or directory
- `rename`: Rename file or directory
- `exists`: Check if file exists
- `mkdir`: Create directory
- `list`: List directory contents

**Example:**
```json
{
  "type": "file",
  "parameters": {
    "action": "write",
    "path": "./output/test.txt",
    "content": "Hello World!"
  }
}
```

### System Operations

Execute system commands and operations.

**Actions:**
- `command`: Execute system command
- `shell`: Execute shell command
- `sleep`: Wait for specified duration
- `open`: Open file, directory, or URL
- `env`: Get or set environment variable
- `kill`: Kill process by name or PID

**Example:**
```json
{
  "type": "system",
  "parameters": {
    "action": "sleep",
    "duration": 5
  }
}
```

### Web Operations

Make HTTP requests to web APIs.

**Actions:**
- `get`: HTTP GET request
- `post`: HTTP POST request
- `put`: HTTP PUT request
- `patch`: HTTP PATCH request
- `delete`: HTTP DELETE request
- `request`: Custom HTTP request

**Example:**
```json
{
  "type": "web",
  "parameters": {
    "action": "get",
    "url": "https://api.example.com/data",
    "headers": {
      "Authorization": "Bearer token"
    }
  }
}
```

### Control Operations

Manage workflow control flow and variables.

**Actions:**
- `set_variable`: Set a variable value
- `get_variable`: Get a variable value
- `increment`: Increment a numeric variable
- `decrement`: Decrement a numeric variable
- `append`: Append value to array variable
- `log`: Log a message
- `assert`: Assert a condition

**Example:**
```json
{
  "type": "control",
  "parameters": {
    "action": "set_variable",
    "name": "counter",
    "value": 0
  }
}
```

## Advanced Features

### Variables

Variables can be defined at the workflow level and referenced using `${variable_name}` syntax:

```json
{
  "variables": [
    {
      "name": "api_url",
      "value": "https://api.example.com",
      "type": "string"
    }
  ],
  "steps": [
    {
      "type": "web",
      "parameters": {
        "action": "get",
        "url": "${api_url}/data"
      }
    }
  ]
}
```

### Conditions

Execute steps conditionally:

```json
{
  "conditions": [
    {
      "variable": "counter",
      "operator": "greater",
      "value": 5
    }
  ]
}
```

**Operators:**
- `equals`: Check equality
- `not_equals`: Check inequality
- `greater`: Greater than
- `less`: Less than
- `contains`: Check if value contains another value
- `exists`: Check if variable exists

### Loops

Repeat steps multiple times:

```json
{
  "loop": {
    "type": "count",
    "count": 10,
    "variable": "iteration"
  }
}
```

**Loop Types:**
- `count`: Fixed number of iterations
- `while`: Continue while condition is true
- `for_each`: Iterate over array values

### Retry Mechanism

Configure retry behavior for unreliable operations:

```json
{
  "retry": {
    "max_retries": 3,
    "delay": 2,
    "backoff": true
  }
}
```

### Parallel Execution

Run multiple steps concurrently:

```json
{
  "parallel": true
}
```

## Examples

Check the `examples/` directory for various workflow examples:

- `simple_workflow.json`: Basic file operations
- `adb_workflow.json`: Android device automation
- `web_workflow.json`: Web API automation
- `advanced_workflow.json`: Advanced control flow

## Project Structure

```
auto-workflow/
├── cmd/
│   └── main.go              # Main CLI application
├── pkg/
│   ├── models/
│   │   └── workflow.go      # Workflow data models
│   ├── operations/
│   │   ├── handler.go       # Handler interface and registry
│   │   ├── adb.go           # ADB operations
│   │   ├── file.go          # File operations
│   │   ├── system.go        # System operations
│   │   ├── web.go           # Web operations
│   │   └── control.go       # Control flow operations
│   └── engine/
│       ├── parser.go        # Workflow parser
│       └── executor.go      # Workflow executor
├── examples/                # Example workflows
├── go.mod                   # Go module definition
└── README.md                # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please open an issue on the project repository.
