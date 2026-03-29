# Quick Start Guide

This guide will help you get started with Auto Workflow quickly.

## Prerequisites

1. **Go 1.21 or higher**
   ```bash
   go version
   ```

2. **ADB (for Android device automation)**
   ```bash
   adb version
   ```

## Installation

### 1. Clone or Download the Project

```bash
cd /path/to/your/workspace
# If you have the project already, navigate to it
cd auto-workflow
```

### 2. Build the Application

```bash
make build
```

Or manually:

```bash
go build -o build/auto-workflow cmd/main.go
```

### 3. Verify Installation

```bash
./build/auto-workflow -version
```

You should see:
```
Auto Workflow v1.0.0
A flexible automation workflow system
```

## Your First Workflow

Let's create a simple workflow to get you started.

### Step 1: Create a Simple Workflow

Create a file named `my_first_workflow.json`:

```json
{
  "name": "My First Workflow",
  "description": "A simple workflow to get started",
  "version": "1.0.0",
  "variables": [
    {
      "name": "message",
      "value": "Hello from Auto Workflow!",
      "type": "string"
    }
  ],
  "steps": [
    {
      "id": "step1",
      "name": "Log message",
      "type": "control",
      "enabled": true,
      "parallel": false,
      "parameters": {
        "action": "log",
        "message": "${message}",
        "level": "info"
      }
    },
    {
      "id": "step2",
      "name": "Create output directory",
      "type": "file",
      "enabled": true,
      "parallel": false,
      "parameters": {
        "action": "mkdir",
        "path": "./my_output",
        "recursive": true
      }
    },
    {
      "id": "step3",
      "name": "Write to file",
      "type": "file",
      "enabled": true,
      "parallel": false,
      "parameters": {
        "action": "write",
        "path": "./my_output/hello.txt",
        "content": "${message}"
      }
    },
    {
      "id": "step4",
      "name": "Read file",
      "type": "file",
      "enabled": true,
      "parallel": false,
      "parameters": {
        "action": "read",
        "path": "./my_output/hello.txt"
      }
    }
  ],
  "on_error": {
    "action": "stop",
    "notify": true
  }
}
```

### Step 2: Run Your Workflow

```bash
./build/auto-workflow -workflow my_first_workflow.json
```

### Step 3: Check the Output

You should see output similar to:

```
Starting workflow: My First Workflow (version: 1.0.0)
Description: A simple workflow to get started
==========================================
Step 1 (step1): Success
Step 2 (step2): Success
Step 3 (step3): Success
Step 4 (step4): Success
==========================================
Workflow completed in Xms
Workflow completed successfully!

Execution Summary:
  Total steps: 4
  Completed steps: 4
  Errors: 0

Variables:
  message: Hello from Auto Workflow!
```

Check the created file:

```bash
cat my_output/hello.txt
```

## Try the Examples

The project includes several example workflows:

### 1. Simple File Operations

```bash
./build/auto-workflow -workflow examples/simple_workflow.json
```

### 2. ADB Device Automation

First, connect your Android device and list devices:

```bash
./build/auto-workflow -list-devices
```

Then run the ADB workflow:

```bash
./build/auto-workflow -workflow examples/adb_workflow.json -device <your_device_id>
```

### 3. Web API Automation

```bash
./build/auto-workflow -workflow examples/web_workflow.json
```

### 4. Advanced Control Flow

```bash
./build/auto-workflow -workflow examples/advanced_workflow.json
```

## Common Workflow Patterns

### Pattern 1: Sequential Steps

Execute steps one after another:

```json
{
  "steps": [
    {
      "id": "step1",
      "type": "file",
      "parameters": {
        "action": "write",
        "path": "./file1.txt",
        "content": "First"
      }
    },
    {
      "id": "step2",
      "type": "file",
      "parameters": {
        "action": "write",
        "path": "./file2.txt",
        "content": "Second"
      }
    }
  ]
}
```

### Pattern 2: Using Variables

Define variables and reference them:

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
      "type": "file",
      "parameters": {
        "action": "mkdir",
        "path": "${output_dir}"
      }
    }
  ]
}
```

### Pattern 3: Conditional Execution

Execute steps based on conditions:

```json
{
  "steps": [
    {
      "id": "conditional_step",
      "type": "control",
      "conditions": [
        {
          "variable": "counter",
          "operator": "greater",
          "value": 5
        }
      ],
      "parameters": {
        "action": "log",
        "message": "Counter is greater than 5"
      }
    }
  ]
}
```

### Pattern 4: Retry with Backoff

Configure retry for unreliable operations:

```json
{
  "steps": [
    {
      "id": "retry_step",
      "type": "web",
      "retry": {
        "max_retries": 3,
        "delay": 2,
        "backoff": true
      },
      "parameters": {
        "action": "get",
        "url": "https://api.example.com/data"
      }
    }
  ]
}
```

### Pattern 5: Loop

Repeat steps multiple times:

```json
{
  "steps": [
    {
      "id": "loop_step",
      "type": "control",
      "loop": {
        "type": "count",
        "count": 5,
        "variable": "i"
      },
      "parameters": {
        "action": "log",
        "message": "Iteration ${i}"
      }
    }
  ]
}
```

## Tips and Best Practices

1. **Start Simple**: Begin with basic operations and gradually add complexity
2. **Use Variables**: Define variables for reusable values
3. **Test Incrementally**: Test each step before adding more
4. **Handle Errors**: Configure error handling appropriately
5. **Use Logging**: Add log steps to track execution progress
6. **Set Timeouts**: Configure timeouts for network operations
7. **Organize Workflows**: Keep related steps together
8. **Comment Your Workflows**: Use the `description` field to document

## Troubleshooting

### Build Errors

If you encounter build errors:

```bash
# Clean and rebuild
make clean
make build

# Or manually
rm -rf build/
go build -o build/auto-workflow cmd/main.go
```

### ADB Device Not Found

```bash
# Check ADB server
adb start-server

# List devices
adb devices

# If needed, restart ADB
adb kill-server
adb start-server
```

### Permission Errors

```bash
# Make the binary executable
chmod +x build/auto-workflow

# Or run with go run directly
go run cmd/main.go -workflow examples/simple_workflow.json
```

### JSON Syntax Errors

Use a JSON validator to check your workflow file:

```bash
# Using Python
python -m json.tool my_workflow.json

# Or use online JSON validators
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Explore the [examples/](examples/) directory for more workflows
- Experiment with different operation types
- Create your own custom workflows

## Getting Help

If you encounter issues:

1. Check the [README.md](README.md) for detailed documentation
2. Review the example workflows
3. Ensure all prerequisites are installed
4. Check the workflow JSON syntax

Happy automating! 🚀
