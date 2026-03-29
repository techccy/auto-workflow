# Migration Guide: Go to Python

This guide helps you migrate from the Go version of Auto Workflow to the Python version.

## Overview

The Python version maintains the same JSON workflow format and functionality as the Go version, making migration straightforward. All existing workflow files should work without modification.

## Key Differences

### Installation

**Go Version:**
```bash
go build -o auto-workflow cmd/main.go
```

**Python Version:**
```bash
pip install -r requirements.txt
```

### Execution

**Go Version:**
```bash
./auto-workflow -workflow examples/simple_workflow.json
```

**Python Version:**
```bash
python main.py -w examples/simple_workflow.json
```

## Compatibility

### Workflow Format

✅ **Fully Compatible** - All existing JSON workflow files work without changes

### Features

| Feature | Go Version | Python Version | Status |
|----------|------------|----------------|---------|
| ADB Operations | ✅ | ✅ | Compatible |
| File Operations | ✅ | ✅ | Compatible |
| System Operations | ✅ | ✅ | Compatible |
| Web Operations | ✅ | ✅ | Compatible |
| Control Operations | ✅ | ✅ | Compatible |
| Variables | ✅ | ✅ | Compatible |
| Conditions | ✅ | ✅ | Compatible |
| Loops | ✅ | ✅ | Compatible |
| Retry Mechanism | ✅ | ✅ | Compatible |
| Parallel Execution | ✅ | ✅ | Compatible |

### Operation Actions

All operation actions from the Go version are supported in the Python version:

#### ADB Operations
- `tap`, `swipe`, `input`, `install`, `uninstall`, `screencap`, `shell`, `click`, `press`

#### File Operations
- `read`, `write`, `append`, `delete`, `copy`, `move`, `rename`, `exists`, `mkdir`, `list`

#### System Operations
- `command`, `shell`, `sleep`, `open`, `env`, `kill`

#### Web Operations
- `get`, `post`, `put`, `patch`, `delete`, `request`

#### Control Operations
- `set_variable`, `get_variable`, `increment`, `decrement`, `append`, `log`, `assert`

## Migration Steps

### 1. Install Python Dependencies

```bash
pip install -r requirements.txt
```

### 2. Update Command Line Usage

Replace Go command invocations with Python equivalents:

```bash
# Old (Go)
./auto-workflow -workflow examples/simple_workflow.json

# New (Python)
python main.py -w examples/simple_workflow.json
```

### 3. Optional: Install as Package

For easier access, you can install the Python version as a package:

```bash
pip install -e .
```

Then use the command directly:
```bash
auto-workflow -w examples/simple_workflow.json
```

## Testing Your Workflows

After migration, test your workflows to ensure they work correctly:

```bash
# Test a simple workflow
python main.py -w examples/simple_workflow.json

# Test with ADB device
python main.py -w examples/adb_workflow.json -d <device_id>

# List ADB devices
python main.py -l
```

## Performance Considerations

The Python version has similar performance characteristics to the Go version for most workflows. However, consider:

- **I/O-bound operations**: Python performs similarly to Go
- **CPU-bound operations**: Go may be slightly faster
- **Parallel execution**: Both versions support parallel steps

## Troubleshooting

### Common Issues

**Issue:** "Module not found" error
```bash
# Solution: Ensure you're in the project directory
cd /path/to/auto-workflow
python main.py -w examples/simple_workflow.json
```

**Issue:** Variable substitution not working
```bash
# Solution: Ensure variable names are correctly defined in the workflow
# Check that variables are referenced as ${variable_name}
```

**Issue:** ADB commands failing
```bash
# Solution: Ensure ADB is installed and in your PATH
adb version
```

## Advanced Usage

### Custom Handlers

Both Go and Python versions support custom handlers. The Python version uses a simpler interface:

```python
from src.operations import Handler
from src.models import ExecutionContext

class CustomHandler(Handler):
    def get_type(self) -> str:
        return "custom"

    def validate(self, params: Dict[str, Any]) -> None:
        # Validate parameters
        pass

    def execute(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        # Execute operation
        return {"success": True}
```

### Programmatic Usage

You can also use the Python version programmatically:

```python
from src.engine import Parser, Executor

# Parse workflow
parser = Parser()
workflow = parser.parse_file("workflow.json")

# Execute workflow
executor = Executor()
exec_ctx = executor.execute(workflow)

# Access results
print(f"Variables: {exec_ctx.variables}")
print(f"Results: {exec_ctx.results}")
```

## Support

If you encounter issues during migration:

1. Check the [README_PYTHON.md](README_PYTHON.md) for Python-specific documentation
2. Review the [examples/](examples/) directory for sample workflows
3. Open an issue on the project repository

## Conclusion

The Python version of Auto Workflow provides full compatibility with the Go version while offering the benefits of Python's ecosystem and ease of use. Most workflows will work without any modifications, and the migration process is straightforward.
