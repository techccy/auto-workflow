# Auto Workflow (Python Version)

A flexible and powerful automation workflow system written in Python that allows you to define and execute complex automation tasks using JSON configuration files.

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

- Python 3.8 or higher
- ADB (for Android device automation)

### Setup

```bash
# Clone the repository
git clone <repository-url>
cd auto-workflow

# Install dependencies
pip install -r requirements.txt
```

## Usage

### Basic Usage

```bash
# Run a workflow
python main.py -w examples/simple_workflow.json

# Run with specific ADB device
python main.py -w examples/adb_workflow.json -device <device_id>

# List connected ADB devices
python main.py -list-devices

# Show version
python main.py -version
```

### Command Line Options

- `-w, --workflow <path>`: Path to the workflow JSON file (required)
- `-d, --device <id>`: ADB device ID to use for ADB operations
- `-l, --list-devices`: List all connected ADB devices
- `-v, --version`: Show version information

## Workflow Structure

A workflow is defined in JSON format with the following structure:
Step 1: Create a new folder named "test——output"
Step 2: Crawl https://cloud.techccy.dpdns.org and write the content to ./test_output/in.html.

```json
{
  "name": "CCY",
  "version": "1.0.0",
  "variables": [
    {
      "name": "target_dir",
      "value": "./test_output",
      "type": "string"
    }
  ],
  "steps": [
    {
      "id": "step1_mkdir",
      "name": "创建文件夹",
      "type": "file",
      "enabled": true,
      "parameters": {
        "action": "mkdir",
        "path": "${target_dir}" 
      }
    },
    {
      "id": "step2_fetch",
      "name": "抓取网页",
      "type": "web",
      "enabled": true,
      "parameters": {
        "action": "get",
        "url": "https://cloud.techccy.dpdns.org"
      }
    },
    {
      "id": "step3_save",
      "name": "保存网页",
      "type": "file",
      "enabled": true,
      "parameters": {
        "action": "write",
        "path": "${target_dir}/in.html",
        "content": "${step2_fetch.body}"
      }
    },
    {
      "id": "step4_screenshot",
      "name": "手机截屏",
      "type": "adb",
      "enabled": false,
      "parameters": {
        "action": "screencap",
        "output_path": "${target_dir}/screenshot.png"
      }
    },
    {
      "id": "step5_log",
      "name": "记录状态码",
      "type": "file",
      "enabled": true,
      "parameters": {
        "action": "write",
        "path": "${target_dir}/status.txt",
        "content": "HTTP Status Code: ${step2_fetch.status_code}\nSuccess: ${step2_fetch.success}\nDuration: ${step2_fetch.duration}ms"
      }
    }
  ]
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
├── src/
│   ├── models/
│   │   ├── __init__.py
│   │   └── workflow.py          # Workflow data models
│   ├── operations/
│   │   ├── __init__.py
│   │   ├── handler.py           # Handler interface and registry
│   │   ├── adb.py               # ADB operations
│   │   ├── file.py              # File operations
│   │   ├── system.py            # System operations
│   │   ├── web.py               # Web operations
│   │   └── control.py           # Control flow operations
│   └── engine/
│       ├── __init__.py
│       ├── parser.py            # Workflow parser
│       └── executor.py          # Workflow executor
├── examples/                    # Example workflows
├── main.py                      # Main CLI application
├── requirements.txt             # Python dependencies
└── README_PYTHON.md            # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please open an issue on the project repository.
