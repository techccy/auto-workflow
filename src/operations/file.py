"""
File system operations handler for the auto-workflow system.
"""
import os
import shutil
from typing import Dict, Any, List
from .handler import Handler
from ..models import ExecutionContext


class FileHandler(Handler):
    """Handles file system operations."""

    def get_type(self) -> str:
        """Return the operation type."""
        return "file"

    def validate(self, params: Dict[str, Any]) -> None:
        """Check if the parameters are valid."""
        action = params.get("action")
        if not action or not isinstance(action, str):
            raise ValueError("action parameter is required")

        if action in ["read", "delete", "exists"]:
            if "path" not in params:
                raise ValueError(f"path parameter is required for {action} action")
        elif action in ["write", "append"]:
            if "path" not in params:
                raise ValueError(f"path parameter is required for {action} action")
            if "content" not in params:
                raise ValueError(f"content parameter is required for {action} action")
        elif action == "copy":
            if "source" not in params:
                raise ValueError("source parameter is required for copy action")
            if "destination" not in params:
                raise ValueError("destination parameter is required for copy action")
        elif action in ["move", "rename"]:
            if "source" not in params:
                raise ValueError(f"source parameter is required for {action} action")
            if "destination" not in params:
                raise ValueError(f"destination parameter is required for {action} action")
        elif action == "mkdir":
            if "path" not in params:
                raise ValueError("path parameter is required for mkdir action")
        elif action == "list":
            if "path" not in params:
                raise ValueError("path parameter is required for list action")

    def execute(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Perform the file operation."""
        action = params["action"]

        if action == "read":
            return self._execute_read(params)
        elif action == "write":
            return self._execute_write(params)
        elif action == "append":
            return self._execute_append(params)
        elif action == "delete":
            return self._execute_delete(params)
        elif action == "copy":
            return self._execute_copy(params)
        elif action == "move":
            return self._execute_move(params)
        elif action == "rename":
            return self._execute_rename(params)
        elif action == "exists":
            return self._execute_exists(params)
        elif action == "mkdir":
            return self._execute_mkdir(params)
        elif action == "list":
            return self._execute_list(params)
        else:
            raise ValueError(f"unsupported file action: {action}")

    def _execute_read(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute read operation."""
        path = self._get_string_param(params, "path")

        with open(path, 'r', encoding='utf-8') as f:
            content = f.read()

        return {
            "success": True,
            "path": path,
            "content": content,
            "size": len(content),
        }

    def _execute_write(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute write operation."""
        path = self._get_string_param(params, "path")
        content = self._get_string_param(params, "content")

        # Create directory if it doesn't exist
        dir_path = os.path.dirname(path)
        if dir_path:
            os.makedirs(dir_path, exist_ok=True)

        with open(path, 'w', encoding='utf-8') as f:
            f.write(content)

        return {
            "success": True,
            "path": path,
            "size": len(content),
        }

    def _execute_append(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute append operation."""
        path = self._get_string_param(params, "path")
        content = self._get_string_param(params, "content")

        with open(path, 'a', encoding='utf-8') as f:
            f.write(content)

        return {
            "success": True,
            "path": path,
            "appended": len(content),
        }

    def _execute_delete(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute delete operation."""
        path = self._get_string_param(params, "path")

        if os.path.isdir(path):
            shutil.rmtree(path)
        else:
            os.remove(path)

        return {
            "success": True,
            "path": path,
        }

    def _execute_copy(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute copy operation."""
        source = self._get_string_param(params, "source")
        destination = self._get_string_param(params, "destination")

        # Create destination directory if it doesn't exist
        dir_path = os.path.dirname(destination)
        if dir_path:
            os.makedirs(dir_path, exist_ok=True)

        # Check if source is a directory
        if os.path.isdir(source):
            shutil.copytree(source, destination)
        else:
            shutil.copy2(source, destination)

        return {
            "success": True,
            "source": source,
            "destination": destination,
        }

    def _execute_move(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute move operation."""
        source = self._get_string_param(params, "source")
        destination = self._get_string_param(params, "destination")

        # Create destination directory if it doesn't exist
        dir_path = os.path.dirname(destination)
        if dir_path:
            os.makedirs(dir_path, exist_ok=True)

        shutil.move(source, destination)

        return {
            "success": True,
            "source": source,
            "destination": destination,
        }

    def _execute_rename(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute rename operation."""
        source = self._get_string_param(params, "source")
        destination = self._get_string_param(params, "destination")

        os.rename(source, destination)

        return {
            "success": True,
            "source": source,
            "destination": destination,
        }

    def _execute_exists(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute exists operation."""
        path = self._get_string_param(params, "path")

        exists = os.path.exists(path)

        return {
            "success": True,
            "path": path,
            "exists": exists,
        }

    def _execute_mkdir(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute mkdir operation."""
        path = self._get_string_param(params, "path")
        recursive = self._get_bool_param(params, "recursive", True)

        if recursive:
            os.makedirs(path, exist_ok=True)
        else:
            os.makedirs(path, exist_ok=False)

        return {
            "success": True,
            "path": path,
        }

    def _execute_list(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute list operation."""
        path = self._get_string_param(params, "path")

        entries = os.listdir(path)
        files = []

        for entry in entries:
            entry_path = os.path.join(path, entry)
            is_dir = os.path.isdir(entry_path)
            size = os.path.getsize(entry_path) if not is_dir else 0

            files.append({
                "name": entry,
                "is_dir": is_dir,
                "size": size,
            })

        return {
            "success": True,
            "path": path,
            "files": files,
            "count": len(files),
        }

    # Helper methods
    def _get_string_param(self, params: Dict[str, Any], key: str) -> str:
        """Get string parameter."""
        value = params.get(key)
        if isinstance(value, str):
            return value
        return str(value) if value is not None else ""

    def _get_bool_param(self, params: Dict[str, Any], key: str, default: bool = False) -> bool:
        """Get boolean parameter."""
        value = params.get(key)
        if isinstance(value, bool):
            return value
        return default
