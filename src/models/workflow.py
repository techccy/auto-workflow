"""
Workflow data models for the auto-workflow system.
"""
from dataclasses import dataclass, field
from typing import Dict, List, Any, Optional
from datetime import datetime


@dataclass
class Variable:
    """Represents a workflow-level variable."""
    name: str
    value: Any
    type: str = "string"  # string, number, boolean, array, object


@dataclass
class RetryPolicy:
    """Defines retry behavior for a step."""
    max_retries: int
    delay: int  # in seconds
    backoff: bool = False  # exponential backoff


@dataclass
class Condition:
    """Represents a conditional check."""
    variable: str
    operator: str  # equals, not_equals, greater, less, contains, exists
    value: Any


@dataclass
class LoopConfig:
    """Defines loop behavior for a step."""
    type: str  # count, while, for_each
    count: Optional[int] = None
    variable: Optional[str] = None
    values: Optional[List[Any]] = None


@dataclass
class ErrorHandler:
    """Defines error handling behavior."""
    action: str  # continue, stop, retry
    steps: Optional[List[str]] = None
    notify: bool = False


@dataclass
class Step:
    """Represents a single operation in the workflow."""
    id: str
    name: str
    type: str  # adb, file, system, web, control
    parameters: Dict[str, Any]
    enabled: bool = True
    parallel: bool = False
    retry: Optional[RetryPolicy] = None
    timeout: Optional[int] = None  # in seconds
    conditions: Optional[List[Condition]] = None
    loop: Optional[LoopConfig] = None
    on_success: Optional[List[str]] = None
    on_failure: Optional[List[str]] = None


@dataclass
class ExecutionContext:
    """Holds runtime context for workflow execution."""
    variables: Dict[str, Any] = field(default_factory=dict)
    results: Dict[str, Any] = field(default_factory=dict)
    step_index: int = 0
    start_time: datetime = field(default_factory=datetime.now)
    errors: List[Exception] = field(default_factory=list)


@dataclass
class Workflow:
    """Represents a complete automation workflow."""
    name: str
    description: str
    version: str
    steps: List[Step]
    variables: Optional[List[Variable]] = None
    on_error: Optional[ErrorHandler] = None
