"""
Operations module for auto-workflow system.
"""
from .handler import Handler, HandlerRegistry
from .adb import ADBHandler
from .file import FileHandler
from .system import SystemHandler
from .web import WebHandler
from .control import ControlHandler

__all__ = [
    'Handler',
    'HandlerRegistry',
    'ADBHandler',
    'FileHandler',
    'SystemHandler',
    'WebHandler',
    'ControlHandler',
]
