"""
Workflow executor for the auto-workflow system.
"""
import time
import threading
from typing import Dict, Any, List
from ..models import Workflow, ExecutionContext, Step
from ..operations import HandlerRegistry, ADBHandler, FileHandler, SystemHandler, WebHandler, ControlHandler


class Executor:
    """Handles workflow execution."""

    def __init__(self):
        self.registry = HandlerRegistry()

        # Register all handlers
        self.registry.register(ADBHandler())
        self.registry.register(FileHandler())
        self.registry.register(SystemHandler())
        self.registry.register(WebHandler())
        self.registry.register(ControlHandler())

        # Get reference to ADB handler for device management
        self.adb_handler = None
        try:
            self.adb_handler = self.registry.get("adb")
        except ValueError:
            pass

    def set_adb_device(self, device_id: str) -> None:
        """Set the target ADB device."""
        if self.adb_handler:
            self.adb_handler.set_device(device_id)

    def execute(self, workflow: Workflow) -> ExecutionContext:
        """Execute a workflow."""
        # Initialize execution context
        exec_ctx = ExecutionContext()
        exec_ctx.variables = {}
        exec_ctx.results = {}
        exec_ctx.errors = []
        exec_ctx.start_time = time.time()

        # Initialize workflow variables
        if workflow.variables:
            for variable in workflow.variables:
                exec_ctx.variables[variable.name] = variable.value

        print(f"Starting workflow: {workflow.name} (version: {workflow.version})")
        print(f"Description: {workflow.description}")
        print("=" * 50)

        # Execute steps
        for i, step in enumerate(workflow.steps):
            exec_ctx.step_index = i

            # Check if step is enabled
            if not step.enabled:
                print(f"Step {i+1} ({step.id}): Skipped (disabled)")
                continue

            # Check conditions
            if not self._evaluate_conditions(step.conditions, exec_ctx):
                print(f"Step {i+1} ({step.id}): Skipped (conditions not met)")
                continue

            # Execute step with retry logic
            try:
                result = self._execute_step_with_retry(step, exec_ctx)
                exec_ctx.results[step.id] = result
                print(f"Step {i+1} ({step.id}): Success")
            except Exception as e:
                exec_ctx.errors.append(e)
                print(f"Step {i+1} ({step.id}): Failed - {e}")

                # Handle error
                if workflow.on_error:
                    if workflow.on_error.action == "stop":
                        raise RuntimeError(f"workflow stopped due to error: {e}") from e
                    elif workflow.on_error.action == "continue":
                        continue
                else:
                    raise RuntimeError(f"workflow failed at step {i+1}: {e}") from e

        duration = time.time() - exec_ctx.start_time
        print("=" * 50)
        print(f"Workflow completed in {duration:.2f} seconds")
        if exec_ctx.errors:
            print(f"Errors: {len(exec_ctx.errors)}")

        return exec_ctx

    def _execute_step_with_retry(self, step: Step, exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute a step with retry logic."""
        # Handle loop if configured
        if step.loop:
            return self._execute_step_with_loop(step, exec_ctx)

        # Execute with retry
        max_retries = 1
        if step.retry:
            max_retries = step.retry.max_retries + 1

        last_error = None

        for attempt in range(1, max_retries + 1):
            try:
                return self._execute_step(step, exec_ctx)
            except Exception as e:
                last_error = e

                if attempt < max_retries:
                    delay = 1  # Default delay
                    if step.retry:
                        delay = step.retry.delay
                        if step.retry.backoff:
                            delay = step.retry.delay * attempt

                    print(f"  Retry {attempt}/{max_retries-1} in {delay} seconds...")
                    time.sleep(delay)

        raise last_error

    def _execute_step_with_loop(self, step: Step, exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute a step with loop logic."""
        iterations = 1
        loop_values = []

        if step.loop.type == "count":
            iterations = step.loop.count if step.loop.count else 1
        elif step.loop.type == "for_each":
            loop_values = step.loop.values if step.loop.values else []
            iterations = len(loop_values)
        elif step.loop.type == "while":
            iterations = -1  # Infinite loop until condition is false

        results = []

        for i in range(iterations if iterations != -1 else 1000):  # Safety limit
            # Set loop variable if specified
            if step.loop.variable:
                if step.loop.type == "for_each" and i < len(loop_values):
                    exec_ctx.variables[step.loop.variable] = loop_values[i]
                else:
                    exec_ctx.variables[step.loop.variable] = i

            # Check while condition
            if step.loop.type == "while":
                condition = step.loop.variable
                if not condition:
                    break
                # Evaluate condition
                if condition in exec_ctx.variables:
                    if not exec_ctx.variables[condition]:
                        break
                else:
                    break

            result = self._execute_step(step, exec_ctx)
            results.append(result)

            if step.loop.type == "while" and i >= 999:
                break  # Safety limit reached

        return {
            "iterations": len(results),
            "results": results,
        }

    def _execute_step(self, step: Step, exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Execute a single step."""
        # Get handler for step type
        handler = self.registry.get(step.type)

        # Resolve variables in parameters
        resolved_params = self._resolve_variables_in_params(step.parameters, exec_ctx)

        # Validate parameters
        handler.validate(resolved_params)

        # Execute step
        result = handler.execute(resolved_params, exec_ctx)

        # Handle success/failure steps (simplified implementation)
        if step.on_success:
            print(f"  Executing success steps: {step.on_success}")
        elif step.on_failure:
            print(f"  Executing failure steps: {step.on_failure}")

        return result

    def _resolve_variables_in_params(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Resolve variable references in parameters."""
        import re

        def resolve_value(value: Any) -> Any:
            if isinstance(value, str):
                # Check for variable reference like ${variable_name}
                result = value
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
            elif isinstance(value, dict):
                return {k: resolve_value(v) for k, v in value.items()}
            elif isinstance(value, list):
                return [resolve_value(item) for item in value]
            else:
                return value

        return {k: resolve_value(v) for k, v in params.items()}

    def _evaluate_conditions(self, conditions: List[Any], exec_ctx: ExecutionContext) -> bool:
        """Evaluate step conditions."""
        if not conditions:
            return True

        for condition in conditions:
            if not self._evaluate_condition(condition, exec_ctx):
                return False

        return True

    def _evaluate_condition(self, condition: Any, exec_ctx: ExecutionContext) -> bool:
        """Evaluate a single condition."""
        if not condition:
            return True

        variable = condition.get("variable")
        operator = condition.get("operator")
        value = condition.get("value")

        if variable not in exec_ctx.variables:
            return False

        var_value = exec_ctx.variables[variable]

        if operator == "equals":
            return str(var_value) == str(value)
        elif operator == "not_equals":
            return str(var_value) != str(value)
        elif operator == "greater":
            return self._compare_values(var_value, value) > 0
        elif operator == "less":
            return self._compare_values(var_value, value) < 0
        elif operator == "contains":
            return self._contains_value(var_value, value)
        elif operator == "exists":
            return True
        else:
            return False

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

    def _to_float(self, value: Any) -> float:
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

    def execute_parallel(self, steps: List[Step], exec_ctx: ExecutionContext) -> None:
        """Execute steps in parallel."""
        threads = []
        errors = []

        def execute_step_thread(step: Step):
            try:
                self._execute_step_with_retry(step, exec_ctx)
            except Exception as e:
                errors.append(e)

        for step in steps:
            if step.parallel:
                thread = threading.Thread(target=execute_step_thread, args=(step,))
                threads.append(thread)
                thread.start()

        # Wait for all threads to complete
        for thread in threads:
            thread.join()

        # Check for errors
        if errors:
            raise RuntimeError(f"parallel execution failed: {errors[0]}")
