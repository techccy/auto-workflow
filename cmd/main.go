package main

import (
	"auto-workflow/pkg/engine"
	"auto-workflow/pkg/operations"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse command line flags
	workflowFile := flag.String("workflow", "", "Path to workflow JSON file")
	device := flag.String("device", "", "ADB device ID to use")
	listDevices := flag.Bool("list-devices", false, "List connected ADB devices")
	version := flag.Bool("version", false, "Show version information")

	flag.Parse()

	// Show version
	if *version {
		fmt.Println("Auto Workflow v1.0.0")
		fmt.Println("A flexible automation workflow system")
		os.Exit(0)
	}

	// List ADB devices
	if *listDevices {
		adbHandler := operations.NewADBHandler()
		devices, err := adbHandler.GetConnectedDevices()
		if err != nil {
			fmt.Printf("Error getting devices: %v\n", err)
			os.Exit(1)
		}

		if len(devices) == 0 {
			fmt.Println("No devices connected")
		} else {
			fmt.Println("Connected devices:")
			for _, device := range devices {
				fmt.Printf("  - %s\n", device)
			}
		}
		os.Exit(0)
	}

	// Check workflow file
	if *workflowFile == "" {
		fmt.Println("Error: workflow file is required")
		fmt.Println("Usage: auto-workflow -workflow <path>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parse workflow
	parser := engine.NewParser()
	workflow, referencedWorkflows, err := parser.ParseMainWorkflow(*workflowFile)
	if err != nil {
		fmt.Printf("Error parsing workflow: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded workflow: %s\n", workflow.Name)
	if len(referencedWorkflows) > 0 {
		fmt.Printf("Referenced workflows: %d\n", len(referencedWorkflows))
	}

	// Create executor
	executor := engine.NewExecutor()

	// Set ADB device if specified
	if *device != "" {
		executor.SetADBDevice(*device)
		fmt.Printf("Using ADB device: %s\n", *device)
	}

	// Execute workflow
	execCtx, err := executor.Execute(workflow)
	if err != nil {
		fmt.Printf("\nWorkflow execution failed: %v\n", err)
		os.Exit(1)
	}

	// Print execution summary
	fmt.Println("\nExecution Summary:")
	fmt.Printf("  Total steps: %d\n", len(workflow.Steps))
	fmt.Printf("  Completed steps: %d\n", execCtx.StepIndex+1)
	fmt.Printf("  Errors: %d\n", len(execCtx.Errors))

	// Print variables
	if len(execCtx.Variables) > 0 {
		fmt.Println("\nVariables:")
		for name, value := range execCtx.Variables {
			fmt.Printf("  %s: %v\n", name, value)
		}
	}

	fmt.Println("\nWorkflow completed successfully!")
}
