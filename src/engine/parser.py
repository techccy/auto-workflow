"""
Workflow JSON parser for the auto-workflow system.
"""
import json
from typing import Dict, Any, Tuple
from ..models import (
    Workflow,
    Step,
    Variable,
    RetryPolicy,
    Condition,
    LoopConfig,
    ErrorHandler,
)


class Parser:
    """Handles workflow JSON parsing."""

    def parse_file(self, file_path: str) -> Workflow:
        """Parse a workflow from a JSON file."""
        with open(file_path, 'r', encoding='utf-8') as f:
            data = json.load(f)
        return self.parse_json(data)

    def parse_json(self, data: Dict[str, Any]) -> Workflow:
        """Parse a workflow from JSON data."""
        workflow = self._parse_workflow(data)
        self.validate(workflow)
        return workflow

    def parse_main_workflow(self, file_path: str) -> Tuple[Workflow, Dict[str, Workflow]]:
        """Parse the main workflow file that may reference other workflows."""
        main_workflow = self.parse_file(file_path)
        referenced_workflows = {}

        # Parse referenced workflows
        for step in main_workflow.steps:
            if "workflow_path" in step.parameters:
                workflow_path = step.parameters["workflow_path"]
                if isinstance(workflow_path, str):
                    referenced_workflow = self.parse_file(workflow_path)
                    referenced_workflows[referenced_workflow.name] = referenced_workflow

        return main_workflow, referenced_workflows

    def _parse_workflow(self, data: Dict[str, Any]) -> Workflow:
        """Parse workflow data from JSON."""
        variables = None
        if "variables" in data:
            variables = [self._parse_variable(v) for v in data["variables"]]

        steps = [self._parse_step(s) for s in data["steps"]]

        on_error = None
        if "on_error" in data:
            on_error = self._parse_error_handler(data["on_error"])

        return Workflow(
            name=data.get("name", ""),
            description=data.get("description", ""),
            version=data.get("version", "1.0.0"),
            variables=variables,
            steps=steps,
            on_error=on_error,
        )

    def _parse_variable(self, data: Dict[str, Any]) -> Variable:
        """Parse variable data from JSON."""
        return Variable(
            name=data.get("name", ""),
            value=data.get("value", None),
            type=data.get("type", "string"),
        )

    def _parse_step(self, data: Dict[str, Any]) -> Step:
        """Parse step data from JSON."""
        retry = None
        if "retry" in data:
            retry = self._parse_retry_policy(data["retry"])

        conditions = None
        if "conditions" in data:
            conditions = [self._parse_condition(c) for c in data["conditions"]]

        loop = None
        if "loop" in data:
            loop = self._parse_loop_config(data["loop"])

        return Step(
            id=data.get("id", ""),
            name=data.get("name", ""),
            type=data.get("type", ""),
            parameters=data.get("parameters", {}),
            enabled=data.get("enabled", True),
            parallel=data.get("parallel", False),
            retry=retry,
            timeout=data.get("timeout"),
            conditions=conditions,
            loop=loop,
            on_success=data.get("on_success"),
            on_failure=data.get("on_failure"),
        )

    def _parse_retry_policy(self, data: Dict[str, Any]) -> RetryPolicy:
        """Parse retry policy data from JSON."""
        return RetryPolicy(
            max_retries=data.get("max_retries", 0),
            delay=data.get("delay", 1),
            backoff=data.get("backoff", False),
        )

    def _parse_condition(self, data: Dict[str, Any]) -> Condition:
        """Parse condition data from JSON."""
        return Condition(
            variable=data.get("variable", ""),
            operator=data.get("operator", ""),
            value=data.get("value", None),
        )

    def _parse_loop_config(self, data: Dict[str, Any]) -> LoopConfig:
        """Parse loop config data from JSON."""
        return LoopConfig(
            type=data.get("type", "count"),
            count=data.get("count"),
            variable=data.get("variable"),
            values=data.get("values"),
        )

    def _parse_error_handler(self, data: Dict[str, Any]) -> ErrorHandler:
        """Parse error handler data from JSON."""
        return ErrorHandler(
            action=data.get("action", "stop"),
            steps=data.get("steps"),
            notify=data.get("notify", False),
        )

    def validate(self, workflow: Workflow) -> None:
        """Validate a workflow."""
        if not workflow.name:
            raise ValueError("workflow name is required")

        if not workflow.steps:
            raise ValueError("workflow must have at least one step")

        # Validate each step
        for i, step in enumerate(workflow.steps):
            if not step.id:
                raise ValueError(f"step {i}: ID is required")
            if not step.type:
                raise ValueError(f"step {i} ({step.id}): type is required")
            if not step.parameters:
                raise ValueError(f"step {i} ({step.id}): parameters are required")
