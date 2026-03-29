"""
Web operations handler for the auto-workflow system.
"""
import json
import time
from typing import Dict, Any
import requests
from .handler import Handler
from ..models import ExecutionContext


class WebHandler(Handler):
    """Handles web automation operations."""

    def __init__(self):
        self.client = requests.Session()
        self.client.timeout = 30

    def set_timeout(self, timeout: int) -> None:
        """Set the HTTP client timeout."""
        self.client.timeout = timeout

    def get_type(self) -> str:
        """Return the operation type."""
        return "web"

    def validate(self, params: Dict[str, Any]) -> None:
        """Check if the parameters are valid."""
        action = params.get("action")
        if not action or not isinstance(action, str):
            raise ValueError("action parameter is required")

        if action in ["get", "post", "put", "patch", "delete"]:
            if "url" not in params:
                raise ValueError(f"url parameter is required for {action} action")
        elif action == "request":
            if "url" not in params:
                raise ValueError("url parameter is required for request action")
            if "method" not in params:
                raise ValueError("method parameter is required for request action")

    def execute(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Perform the web operation."""
        action = params["action"]

        if action == "get":
            return self._execute_get(params)
        elif action == "post":
            return self._execute_post(params)
        elif action == "put":
            return self._execute_put(params)
        elif action == "patch":
            return self._execute_patch(params)
        elif action == "delete":
            return self._execute_delete(params)
        elif action == "request":
            return self._execute_request(params)
        else:
            raise ValueError(f"unsupported web action: {action}")

    def _execute_get(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute GET request."""
        url = self._get_string_param(params, "url")
        headers = self._get_headers(params)

        start_time = time.time()
        response = self.client.get(url, headers=headers)
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": 200 <= response.status_code < 300,
            "status_code": response.status_code,
            "body": response.text,
            "headers": dict(response.headers),
            "duration": duration,
        }

    def _execute_post(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute POST request."""
        url = self._get_string_param(params, "url")
        headers = self._get_headers(params)

        # Prepare body
        body = None
        if "body" in params:
            body = params["body"]
            if isinstance(body, dict):
                if headers.get("Content-Type") is None:
                    headers["Content-Type"] = "application/json"
                body = json.dumps(body)
            elif isinstance(body, str):
                pass
            else:
                body = str(body)

        start_time = time.time()
        response = self.client.post(url, headers=headers, data=body)
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": 200 <= response.status_code < 300,
            "status_code": response.status_code,
            "body": response.text,
            "headers": dict(response.headers),
            "duration": duration,
        }

    def _execute_put(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute PUT request."""
        url = self._get_string_param(params, "url")
        headers = self._get_headers(params)

        # Prepare body
        body = None
        if "body" in params:
            body = params["body"]
            if isinstance(body, dict):
                if headers.get("Content-Type") is None:
                    headers["Content-Type"] = "application/json"
                body = json.dumps(body)
            elif isinstance(body, str):
                pass
            else:
                body = str(body)

        start_time = time.time()
        response = self.client.put(url, headers=headers, data=body)
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": 200 <= response.status_code < 300,
            "status_code": response.status_code,
            "body": response.text,
            "headers": dict(response.headers),
            "duration": duration,
        }

    def _execute_patch(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute PATCH request."""
        url = self._get_string_param(params, "url")
        headers = self._get_headers(params)

        # Prepare body
        body = None
        if "body" in params:
            body = params["body"]
            if isinstance(body, dict):
                if headers.get("Content-Type") is None:
                    headers["Content-Type"] = "application/json"
                body = json.dumps(body)
            elif isinstance(body, str):
                pass
            else:
                body = str(body)

        start_time = time.time()
        response = self.client.patch(url, headers=headers, data=body)
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": 200 <= response.status_code < 300,
            "status_code": response.status_code,
            "body": response.text,
            "headers": dict(response.headers),
            "duration": duration,
        }

    def _execute_delete(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute DELETE request."""
        url = self._get_string_param(params, "url")
        headers = self._get_headers(params)

        start_time = time.time()
        response = self.client.delete(url, headers=headers)
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": 200 <= response.status_code < 300,
            "status_code": response.status_code,
            "body": response.text,
            "headers": dict(response.headers),
            "duration": duration,
        }

    def _execute_request(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute custom HTTP request."""
        url = self._get_string_param(params, "url")
        method = self._get_string_param(params, "method").upper()
        headers = self._get_headers(params)

        # Prepare body
        body = None
        if "body" in params:
            body = params["body"]
            if isinstance(body, dict):
                if headers.get("Content-Type") is None:
                    headers["Content-Type"] = "application/json"
                body = json.dumps(body)
            elif isinstance(body, str):
                pass
            else:
                body = str(body)

        start_time = time.time()
        response = self.client.request(method, url, headers=headers, data=body)
        duration = int((time.time() - start_time) * 1000)

        return {
            "success": 200 <= response.status_code < 300,
            "status_code": response.status_code,
            "body": response.text,
            "headers": dict(response.headers),
            "duration": duration,
        }

    def _get_headers(self, params: Dict[str, Any]) -> Dict[str, str]:
        """Get headers from parameters."""
        headers = {}
        if "headers" in params and isinstance(params["headers"], dict):
            for k, v in params["headers"].items():
                if isinstance(v, str):
                    headers[k] = v
                else:
                    headers[k] = str(v)
        return headers

    # Helper methods
    def _get_string_param(self, params: Dict[str, Any], key: str, default: str = "") -> str:
        """Get string parameter."""
        value = params.get(key)
        if isinstance(value, str):
            return value
        return str(value) if value is not None else default
