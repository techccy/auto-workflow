"""
System operations handler for the auto-workflow system.
"""
import os
import subprocess
import time
import platform
from typing import Dict, Any, List
from .handler import Handler
from ..models import ExecutionContext


class SystemHandler(Handler):
    """Handles system operations."""

    def get_type(self) -> str:
        """Return the operation type."""
        return "system"

    def validate(self, params: Dict[str, Any]) -> None:
        """Check if the parameters are valid."""
        action = params.get("action")
        if not action or not isinstance(action, str):
            raise ValueError("action parameter is required")

        if action in ["command", "shell"]:
            if "command" not in params:
                raise ValueError(f"command parameter is required for {action} action")
        elif action == "sleep":
            if "duration" not in params:
                raise ValueError("duration parameter is required for sleep action")
        elif action == "open":
            if "path" not in params:
                raise ValueError("path parameter is required for open action")
        elif action == "env":
            if "variable" not in params:
                raise ValueError("variable parameter is required for env action")

    def execute(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Perform the system operation."""
        action = params["action"]

        if action == "command":
            return self._execute_command(params)
        elif action == "shell":
            return self._execute_shell(params)
        elif action == "sleep":
            return self._execute_sleep(params)
        elif action == "open":
            return self._execute_open(params)
        elif action == "env":
            return self._execute_env(params)
        elif action == "kill":
            return self._execute_kill(params)
        else:
            raise ValueError(f"unsupported system action: {action}")

    def _execute_command(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute command operation."""
        command = self._get_string_param(params, "command")
        working_dir = self._get_string_param(params, "working_dir")

        # Parse command and arguments
        parts = command.split()
        if not parts:
            raise ValueError("empty command")

        cmd = parts[0]
        args = parts[1:]

        # Set environment variables if provided
        env = os.environ.copy()
        if "env" in params and isinstance(params["env"], dict):
            for k, v in params["env"].items():
                env[str(k)] = str(v)

        start_time = time.time()
        result = subprocess.run(
            [cmd] + args,
            capture_output=True,
            text=True,
            cwd=working_dir if working_dir else None,
            env=env
        )
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": result.returncode == 0,
            "command": command,
            "output": result.stdout,
            "error": result.stderr if result.returncode != 0 else None,
            "duration": duration,
        }

    def _execute_shell(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute shell command operation."""
        command = self._get_string_param(params, "command")
        working_dir = self._get_string_param(params, "working_dir")

        # Set environment variables if provided
        env = os.environ.copy()
        if "env" in params and isinstance(params["env"], dict):
            for k, v in params["env"].items():
                env[str(k)] = str(v)

        if working_dir:
            shell_command = f"cd {working_dir} && {command}"
        else:
            shell_command = command

        start_time = time.time()
        result = subprocess.run(
            shell_command,
            capture_output=True,
            text=True,
            shell=True,
            env=env
        )
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": result.returncode == 0,
            "command": command,
            "output": result.stdout,
            "error": result.stderr if result.returncode != 0 else None,
            "duration": duration,
        }

    def _execute_sleep(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute sleep operation."""
        duration = self._get_int_param(params, "duration")

        time.sleep(duration)

        return {
            "success": True,
            "duration": duration,
        }

    def _execute_open(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute open operation."""
        path = self._get_string_param(params, "path")

        system = platform.system()

        if system == "Darwin":  # macOS
            if path.endswith(".app"):
                subprocess.run(["open", "-a", path], check=True)
            else:
                subprocess.run(["open", path], check=True)
        elif system == "Windows":
            os.startfile(path)
        else:  # Linux and others
            if "://" in path:
                subprocess.run(["xdg-open", path], check=True)
            else:
                subprocess.run(["xdg-open", path], check=True)

        return {
            "success": True,
            "path": path,
        }

    def _execute_env(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute env operation."""
        variable = self._get_string_param(params, "variable")
        value = self._get_string_param(params, "value", "")

        if value == "":
            # Get environment variable
            env_value = os.environ.get(variable, "")
            return {
                "success": True,
                "variable": variable,
                "value": env_value,
            }
        else:
            # Set environment variable (only for current process)
            os.environ[variable] = value
            return {
                "success": True,
                "variable": variable,
                "value": value,
            }

    def _execute_kill(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute kill operation."""
        process_name = self._get_string_param(params, "process_name")
        pid = self._get_int_param(params, "pid")

        if pid > 0:
            # Kill by PID
            try:
                os.kill(pid, 9)  # SIGKILL
            except ProcessLookupError:
                raise RuntimeError(f"process with PID {pid} not found")
        elif process_name:
            # Kill by process name
            system = platform.system()
            if system == "Darwin" or system == "Linux":
                subprocess.run(["pkill", process_name], check=True)
            elif system == "Windows":
                subprocess.run(["taskkill", "/F", "/IM", process_name], check=True)
        else:
            raise ValueError("either process_name or pid must be specified")

        return {
            "success": True,
        }

    # Helper methods
    def _get_string_param(self, params: Dict[str, Any], key: str, default: str = "") -> str:
        """Get string parameter."""
        value = params.get(key)
        if isinstance(value, str):
            return value
        return str(value) if value is not None else default

    def _get_int_param(self, params: Dict[str, Any], key: str, default: int = 0) -> int:
        """Get integer parameter."""
        value = params.get(key)
        if value is None:
            return default
        if isinstance(value, int):
            return value
        if isinstance(value, float):
            return int(value)
        if isinstance(value, str):
            try:
                return int(value)
            except ValueError:
                return default
        return default
