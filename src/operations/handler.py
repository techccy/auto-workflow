"""
Handler interface and registry for operation handlers.
"""
from abc import ABC, abstractmethod
from typing import Dict, Any, List
from ..models import ExecutionContext


class Handler(ABC):
    """Abstract base class for operation handlers."""

    @abstractmethod
    def execute(self, params: Dict[str, Any], exec_ctx: ExecutionContext) -> Dict[str, Any]:
        """Perform the operation with given parameters."""
        pass

    @abstractmethod
    def validate(self, params: Dict[str, Any]) -> None:
        """Check if the parameters are valid for this operation."""
        pass

    @abstractmethod
    def get_type(self) -> str:
        """Return the operation type this handler supports."""
        pass


class HandlerRegistry:
    """Manages operation handlers."""

    def __init__(self):
        self._handlers: Dict[str, Handler] = {}

    def register(self, handler: Handler) -> None:
        """Register a new handler for an operation type."""
        self._handlers[handler.get_type()] = handler

    def get(self, op_type: str) -> Handler:
        """Retrieve a handler by operation type."""
        if op_type not in self._handlers:
            raise ValueError(f"No handler found for type: {op_type}")
        return self._handlers[op_type]

    def list(self) -> List[str]:
        """Return all registered operation types."""
        return list(self._handlers.keys())
