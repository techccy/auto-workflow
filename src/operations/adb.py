"""
ADB operations handler for the auto-workflow system.
"""
import subprocess
import time
from typing import Dict, Any, List
from .handler import Handler
from ..models import ExecutionContext


class ADBHandler(Handler):
    """Handles ADB operations."""

    def __init__(self):
        self.device_id: str = ""

    def set_device(self, device_id: str) -> None:
        """Set the target device ID."""
        self.device_id = device_id

    def get_type(self) -> str:
        """Return the operation type."""
        return "adb"

    def validate(self, params: Dict[str, Any]) -> None:
        """Check if the parameters are valid."""
        action = params.get("action")
        if not action or not isinstance(action, str):
            raise ValueError("action parameter is required")

        if action == "tap":
            if "x" not in params:
                raise ValueError("x parameter is required for tap action")
            if "y" not in params:
                raise ValueError("y parameter is required for tap action")
        elif action == "swipe":
            if "x1" not in params:
                raise ValueError("x1 parameter is required for swipe action")
            if "y1" not in params:
                raise ValueError("y1 parameter is required for swipe action")
            if "x2" not in params:
                raise ValueError("x2 parameter is required for swipe action")
            if "y2" not in params:
                raise ValueError("y2 parameter is required for swipe action")
        elif action == "input":
            if "text" not in params:
                raise ValueError("text parameter is required for input action")
        elif action == "install":
            if "apk_path" not in params:
                raise ValueError("apk_path parameter is required for install action")
        elif action == "uninstall":
            if "package" not in params:
                raise ValueError("package parameter is required for uninstall action")
        elif action == "screencap":
            if "output_path" not in params:
                raise ValueError("output_path parameter is required for screencap action")
        elif action == "shell":
            if "command" not in params:
                raise ValueError("command parameter is required for shell action")

    def execute(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Perform the ADB operation."""
        action = params["action"]

        # Build base command
        args = []
        if self.device_id:
            args.extend(["-s", self.device_id])

        if action == "tap":
            return self._execute_tap(args, params)
        elif action == "swipe":
            return self._execute_swipe(args, params)
        elif action == "input":
            return self._execute_input(args, params)
        elif action == "install":
            return self._execute_install(args, params)
        elif action == "uninstall":
            return self._execute_uninstall(args, params)
        elif action == "screencap":
            return self._execute_screencap(args, params)
        elif action == "shell":
            return self._execute_shell(args, params)
        elif action == "click":
            return self._execute_click(args, params)
        elif action == "press":
            return self._execute_press(args, params)
        else:
            raise ValueError(f"unsupported ADB action: {action}")

    def _execute_tap(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute tap operation."""
        x = self._get_int_param(params, "x")
        y = self._get_int_param(params, "y")

        args = base_args + ["shell", "input", "tap", str(x), str(y)]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"tap failed: {result.stderr}")

        return {
            "success": True,
            "x": x,
            "y": y,
        }

    def _execute_swipe(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute swipe operation."""
        x1 = self._get_int_param(params, "x1")
        y1 = self._get_int_param(params, "y1")
        x2 = self._get_int_param(params, "x2")
        y2 = self._get_int_param(params, "y2")
        duration = self._get_int_param(params, "duration", 300)

        args = base_args + [
            "shell", "input", "swipe",
            str(x1), str(y1), str(x2), str(y2), str(duration)
        ]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"swipe failed: {result.stderr}")

        return {
            "success": True,
            "from": {"x": x1, "y": y1},
            "to": {"x": x2, "y": y2},
            "duration": duration,
        }

    def _execute_input(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute input operation."""
        text = self._get_string_param(params, "text")

        args = base_args + ["shell", "input", "text", text]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"input failed: {result.stderr}")

        return {
            "success": True,
            "text": text,
        }

    def _execute_install(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute install operation."""
        apk_path = self._get_string_param(params, "apk_path")

        args = base_args + ["install", apk_path]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"install failed: {result.stderr}")

        return {
            "success": True,
            "apk": apk_path,
            "output": result.stdout,
        }

    def _execute_uninstall(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute uninstall operation."""
        package = self._get_string_param(params, "package")

        args = base_args + ["uninstall", package]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"uninstall failed: {result.stderr}")

        return {
            "success": True,
            "package": package,
            "output": result.stdout,
        }

    def _execute_screencap(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute screencap operation."""
        output_path = self._get_string_param(params, "output_path")

        # Take screenshot on device
        args = base_args + ["shell", "screencap", "-p", "/sdcard/screenshot.png"]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"screencap failed: {result.stderr}")

        # Pull screenshot to local
        pull_args = base_args + ["pull", "/sdcard/screenshot.png", output_path]
        pull_result = subprocess.run(["adb"] + pull_args, capture_output=True, text=True)

        if pull_result.returncode != 0:
            raise RuntimeError(f"pull screenshot failed: {pull_result.stderr}")

        return {
            "success": True,
            "output_path": output_path,
        }

    def _execute_shell(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute shell command operation."""
        command = self._get_string_param(params, "command")

        args = base_args + ["shell", command]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"shell command failed: {result.stderr}")

        return {
            "success": True,
            "command": command,
            "output": result.stdout,
        }

    def _execute_click(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute click operation (alias for tap)."""
        x = self._get_int_param(params, "x")
        y = self._get_int_param(params, "y")

        args = base_args + ["shell", "input", "tap", str(x), str(y)]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"click failed: {result.stderr}")

        return {
            "success": True,
            "x": x,
            "y": y,
        }

    def _execute_press(self, base_args: List[str], params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute key press operation."""
        key = self._get_string_param(params, "key")

        args = base_args + ["shell", "input", "keyevent", key]
        result = subprocess.run(["adb"] + args, capture_output=True, text=True)

        if result.returncode != 0:
            raise RuntimeError(f"press failed: {result.stderr}")

        return {
            "success": True,
            "key": key,
        }

    def wait_for_device(self, timeout: int = 30) -> None:
        """Wait for device to be connected."""
        args = []
        if self.device_id:
            args.extend(["-s", self.device_id])
        args.append("wait-for-device")

        result = subprocess.run(
            ["adb"] + args,
            capture_output=True,
            text=True,
            timeout=timeout
        )

        if result.returncode != 0:
            raise RuntimeError(f"wait for device failed: {result.stderr}")

    def get_connected_devices(self) -> List[str]:
        """Return list of connected devices."""
        result = subprocess.run(
            ["adb", "devices"],
            capture_output=True,
            text=True
        )

        if result.returncode != 0:
            raise RuntimeError(f"failed to get devices: {result.stderr}")

        devices = []
        lines = result.stdout.split("\n")

        for line in lines[1:]:
            line = line.strip()
            if line and "daemon" not in line:
                parts = line.split()
                if len(parts) >= 2 and parts[1] == "device":
                    devices.append(parts[0])

        return devices

    # Helper methods
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

    def _get_string_param(self, params: Dict[str, Any], key: str) -> str:
        """Get string parameter."""
        value = params.get(key)
        if isinstance(value, str):
            return value
        return str(value) if value is not None else ""
