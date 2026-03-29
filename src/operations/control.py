"""
Control flow operations handler for the auto-workflow system.
"""
import re
from typing import Dict, Any, Union
from .handler import Handler
from ..models import ExecutionContext


class ControlHandler(Handler):
    """Handles control flow operations."""

    def get_type(self) -> str:
        """Return the operation type."""
        return "control"

    def validate(self, params: Dict[str, Any]) -> None:
        """Check if the parameters are valid."""
        action = params.get("action")
        if not action or not isinstance(action, str):
            raise ValueError("action parameter is required")

        if action == "set_variable":
            if "name" not in params:
                raise ValueError("name parameter is required for set_variable action")
            if "value" not in params:
                raise ValueError("value parameter is required for set_variable action")
        elif action == "get_variable":
            if "name" not in params:
                raise ValueError("name parameter is required for get_variable action")
        elif action == "increment":
            if "name" not in params:
                raise ValueError("name parameter is required for increment action")
        elif action == "decrement":
            if "name" not in params:
                raise ValueError("name parameter is required for decrement action")
        elif action == "append":
            if "name" not in params:
                raise ValueError("name parameter is required for append action")
            if "value" not in params:
                raise ValueError("value parameter is required for append action")
        elif action == "log":
            if "message" not in params:
                raise ValueError("message parameter is required for log action")
        elif action == "assert":
            if "condition" not in params:
                raise ValueError("condition parameter is required for assert action")

    def execute(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Perform the control operation."""
        action = params["action"]

        if action == "set_variable":
            return self._execute_set_variable(params, exec_ctx)
        elif action == "get_variable":
            return self._execute_get_variable(params, exec_ctx)
        elif action == "increment":
            return self._execute_increment(params, exec_ctx)
        elif action == "decrement":
            return self._execute_decrement(params, exec_ctx)
        elif action == "append":
            return self._execute_append(params, exec_ctx)
        elif action == "log":
            return self._execute_log(params, exec_ctx)
        elif action == "assert":
            return self._execute_assert(params, exec_ctx)
        else:
            raise ValueError(f"unsupported control action: {action}")

    def _execute_set_variable(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute set_variable operation."""
        name = self._get_string_param(params, "name")
        value = params["value"]

        # Resolve variable references if needed
        resolved_value = self._resolve_value(value, exec_ctx)

        if exec_ctx.variables is None:
            exec_ctx.variables = {}

        exec_ctx.variables[name] = resolved_value

        return {
            "success": True,
            "name": name,
            "value": resolved_value,
        }

    def _execute_get_variable(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute get_variable operation."""
        name = self._get_string_param(params, "name")

        if name not in exec_ctx.variables:
            raise ValueError(f"variable '{name}' not found")

        value = exec_ctx.variables[name]

        return {
            "success": True,
            "name": name,
            "value": value,
        }

    def _execute_increment(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute increment operation."""
        name = self._get_string_param(params, "name")
        amount = self._get_int_param(params, "amount", 1)

        if name not in exec_ctx.variables:
            raise ValueError(f"variable '{name}' not found")

        value = exec_ctx.variables[name]

        if isinstance(value, (int, float)):
            new_value = value + amount
        elif isinstance(value, str):
            try:
                new_value = float(value) + amount
            except ValueError:
                raise ValueError(f"variable '{name}' is not a number")
        else:
            raise ValueError(f"variable '{name}' is not a number")

        exec_ctx.variables[name] = new_value

        return {
            "success": True,
            "name": name,
            "value": new_value,
        }

    def _execute_decrement(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute decrement operation."""
        name = self._get_string_param(params, "name")
        amount = self._get_int_param(params, "amount", 1)

        if name not in exec_ctx.variables:
            raise ValueError(f"variable '{name}' not found")

        value = exec_ctx.variables[name]

        if isinstance(value, (int, float)):
            new_value = value - amount
        elif isinstance(value, str):
            try:
                new_value = float(value) - amount
            except ValueError:
                raise ValueError(f"variable '{name}' is not a number")
        else:
            raise ValueError(f"variable '{name}' is not a number")

        exec_ctx.variables[name] = new_value

        return {
            "success": True,
            "name": name,
            "value": new_value,
        }

    def _execute_append(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute append operation."""
        name = self._get_string_param(params, "name")
        value = params["value"]

        # Resolve variable references if needed
        resolved_value = self._resolve_value(value, exec_ctx)

        if name not in exec_ctx.variables:
            raise ValueError(f"variable '{name}' not found")

        existing_value = exec_ctx.variables[name]

        if isinstance(existing_value, list):
            exec_ctx.variables[name] = existing_value + [resolved_value]
        else:
            raise ValueError(f"variable '{name}' is not an array")

        return {
            "success": True,
            "name": name,
            "appended": resolved_value,
        }

    def _execute_log(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute log operation."""
        message = self._get_string_param(params, "message")
        level = self._get_string_param(params, "level", "info")

        # Resolve variable references in message
        resolved_message = self._resolve_variables(message, exec_ctx)

        print(f"[{level.upper()}] {resolved_message}")

        return {
            "success": True,
            "level": level,
            "message": resolved_message,
        }

    def _execute_assert(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute assert operation."""
        condition = params["condition"]

        # Evaluate the condition
        result = self._evaluate_condition(condition, exec_ctx)

        if not result:
            message = self._get_string_param(params, "message", "assertion failed")
            raise AssertionError(f"assertion failed: {message}")

        return {
            "success": True,
            "condition": condition,
        }

    def _resolve_value(self, value: Any, exec_ctx: ExecutionContext) -> Any:
        """Resolve variable references in a value."""
        if isinstance(value, str):
            # Check for variable reference like ${variable_name}
            match = re.match(r'\$\{(.+?)\}', value)
            if match:
                var_name = match.group(1)
                if var_name in exec_ctx.variables:
                    return exec_ctx.variables[var_name]
                else:
                    raise ValueError(f"variable '{var_name}' not found")
            return value
        elif isinstance(value, dict):
            result = {}
            for k, v in value.items():
                result[k] = self._resolve_value(v, exec_ctx)
            return result
        elif isinstance(value, list):
            return [self._resolve_value(item, exec_ctx) for item in value]
        else:
            return value

    def _resolve_variables(self, text: str, exec_ctx: ExecutionContext) -> str:
        """Resolve variable references in text."""
        result = text

        # Find all variable references
        while True:
            match = re.search(r'\$\{(.+?)\}', result)
            if not match:
                break

            var_name = match.group(1)
            if var_name in exec_ctx.variables:
                var_value = str(exec_ctx.variables[var_name])
                result = result[:match.start()] + var_value + result[match.end():]
            else:
                # Keep the reference if variable not found
                result = result[:match.end()] + result[match.end():]

        return result

    def _evaluate_condition(self, condition: Any, exec_ctx: ExecutionContext) -> bool:
        """Evaluate a condition."""
        if isinstance(condition, bool):
            return condition
        elif isinstance(condition, str):
            # Try to resolve variable reference
            match = re.match(r'\$\{(.+?)\}', condition)
            if match:
                var_name = match.group(1)
                if var_name in exec_ctx.variables:
                    value = exec_ctx.variables[var_name]
                    # For non-boolean values, check if truthy
                    if isinstance(value, bool):
                        return value
                    elif isinstance(value, (int, float)):
                        return value != 0
                    elif isinstance(value, str):
                        return value != ""
                    else:
                        return value is not None
                else:
                    return False
            # Try to parse as boolean
            if condition.lower() == "true":
                return True
            elif condition.lower() == "false":
                return False
            else:
                raise ValueError(f"invalid condition: {condition}")
        elif isinstance(condition, dict):
            # Complex condition with operator
            return self._evaluate_complex_condition(condition, exec_ctx)
        elif isinstance(condition, (int, float)):
            # Numeric values: true if non-zero
            return condition != 0
        else:
            raise ValueError(f"invalid condition type: {type(condition)}")

    def _evaluate_complex_condition(self, condition: Dict[str, Any], exec_ctx: ExecutionContext) -> bool:
        """Evaluate a complex condition with operator."""
        operator = condition.get("operator")
        if not operator:
            raise ValueError("operator is required for complex condition")

        left = self._resolve_value(condition.get("left"), exec_ctx)
        right = self._resolve_value(condition.get("right"), exec_ctx)

        if operator == "equals":
            return self._compare_values(left, right) == 0
        elif operator == "not_equals":
            return self._compare_values(left, right) != 0
        elif operator == "greater":
            return self._compare_values(left, right) > 0
        elif operator == "less":
            return self._compare_values(left, right) < 0
        elif operator == "greater_equal":
            return self._compare_values(left, right) >= 0
        elif operator == "less_equal":
            return self._compare_values(left, right) <= 0
        elif operator == "contains":
            return self._contains_value(left, right)
        elif operator == "and":
            left_result = self._evaluate_condition(left, exec_ctx)
            if not left_result:
                return False
            return self._evaluate_condition(right, exec_ctx)
        elif operator == "or":
            left_result = self._evaluate_condition(left, exec_ctx)
            if left_result:
                return True
            return self._evaluate_condition(right, exec_ctx)
        elif operator == "not":
            return not self._evaluate_condition(left, exec_ctx)
        else:
            raise ValueError(f"unsupported operator: {operator}")

    def _compare_values(self, a: Any, b: Any) -> int:
        """Compare two values. Returns -1, 0, or 1."""
        # Try to convert to numbers
        a_num = self._to_float(a)
        b_num = self._to_float(b)

        if a_num is not None and b_num is not None:
            if a_num < b_num:
                return -1
            elif a_num > b_num:
                return 1
            else:
                return 0

        # Fall back to string comparison
        a_str = str(a)
        b_str = str(b)

        if a_str < b_str:
            return -1
        elif a_str > b_str:
            return 1
        else:
            return 0

    def _to_float(self, value: Any) -> Union[float, None]:
        """Convert value to float if possible."""
        if isinstance(value, (int, float)):
            return float(value)
        elif isinstance(value, str):
            try:
                return float(value)
            except ValueError:
                return None
        else:
            return None

    def _contains_value(self, container: Any, value: Any) -> bool:
        """Check if container contains value."""
        if isinstance(container, str):
            return str(value) in container
        elif isinstance(container, list):
            return value in container
        else:
            return False

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
