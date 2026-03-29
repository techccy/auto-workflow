#!/usr/bin/env python3
"""
Main CLI entry point for auto-workflow.
"""
import argparse
import sys
from src.engine import Parser, Executor
from src.operations import ADBHandler


def main():
    """Main entry point for the CLI."""
    parser = argparse.ArgumentParser(
        description="Auto Workflow - A flexible automation workflow system"
    )

    parser.add_argument(
        "-w", "--workflow",
        type=str,
        help="Path to workflow JSON file"
    )

    parser.add_argument(
        "-d", "--device",
        type=str,
        help="ADB device ID to use"
    )

    parser.add_argument(
        "-l", "--list-devices",
        action="store_true",
        help="List connected ADB devices"
    )

    parser.add_argument(
        "-v", "--version",
        action="store_true",
        help="Show version information"
    )

    args = parser.parse_args()

    # Show version
    if args.version:
        print("Auto Workflow v1.0.0")
        print("A flexible automation workflow system")
        sys.exit(0)

    # List ADB devices
    if args.list_devices:
        adb_handler = ADBHandler()
        try:
            devices = adb_handler.get_connected_devices()
            if not devices:
                print("No devices connected")
            else:
                print("Connected devices:")
                for device in devices:
                    print(f"  - {device}")
        except Exception as e:
            print(f"Error getting devices: {e}")
            sys.exit(1)
        sys.exit(0)

    # Check workflow file
    if not args.workflow:
        print("Error: workflow file is required")
        print("Usage: python main.py -workflow <path>")
        parser.print_help()
        sys.exit(1)

    # Parse workflow
    workflow_parser = Parser()
    try:
        workflow, referenced_workflows = workflow_parser.parse_main_workflow(args.workflow)
    except Exception as e:
        print(f"Error parsing workflow: {e}")
        sys.exit(1)

    print(f"Loaded workflow: {workflow.name}")
    if referenced_workflows:
        print(f"Referenced workflows: {len(referenced_workflows)}")

    # Create executor
    executor = Executor()

    # Set ADB device if specified
    if args.device:
        executor.set_adb_device(args.device)
        print(f"Using ADB device: {args.device}")

    # Execute workflow
    try:
        exec_ctx = executor.execute(workflow)
    except Exception as e:
        print(f"\nWorkflow execution failed: {e}")
        sys.exit(1)

    # Print execution summary
    print("\nExecution Summary:")
    print(f"  Total steps: {len(workflow.steps)}")
    print(f"  Completed steps: {exec_ctx.step_index + 1}")
    print(f"  Errors: {len(exec_ctx.errors)}")

    # Print variables
    if exec_ctx.variables:
        print("\nVariables:")
        for name, value in exec_ctx.variables.items():
            print(f"  {name}: {value}")

    print("\nWorkflow completed successfully!")


if __name__ == "__main__":
    main()
