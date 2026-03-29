package engine

import (
	"auto-workflow/pkg/models"
	"encoding/json"
	"fmt"
	"os"
)

// Parser handles workflow JSON parsing
type Parser struct{}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses a workflow from a JSON file
func (p *Parser) ParseFile(filePath string) (*models.Workflow, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file failed: %v", err)
	}

	return p.ParseJSON(data)
}

// ParseJSON parses a workflow from JSON data
func (p *Parser) ParseJSON(data []byte) (*models.Workflow, error) {
	var workflow models.Workflow
	if err := json.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("unmarshal JSON failed: %v", err)
	}

	// Validate workflow
	if err := p.Validate(&workflow); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return &workflow, nil
}

// Validate validates a workflow
func (p *Parser) Validate(workflow *models.Workflow) error {
	if workflow.Name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if len(workflow.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	// Validate each step
	for i, step := range workflow.Steps {
		if step.ID == "" {
			return fmt.Errorf("step %d: ID is required", i)
		}
		if step.Type == "" {
			return fmt.Errorf("step %d (%s): type is required", i, step.ID)
		}
		if len(step.Parameters) == 0 {
			return fmt.Errorf("step %d (%s): parameters are required", i, step.ID)
		}
	}

	return nil
}

// ParseMainWorkflow parses the main workflow file that may reference other workflows
func (p *Parser) ParseMainWorkflow(filePath string) (*models.Workflow, map[string]*models.Workflow, error) {
	mainWorkflow, err := p.ParseFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	// Parse referenced workflows
	referencedWorkflows := make(map[string]*models.Workflow)

	// Check if there are referenced workflows in parameters
	for _, step := range mainWorkflow.Steps {
		if workflowPath, ok := step.Parameters["workflow_path"].(string); ok {
			referencedWorkflow, err := p.ParseFile(workflowPath)
			if err != nil {
				return nil, nil, fmt.Errorf("parse referenced workflow '%s' failed: %v", workflowPath, err)
			}
			referencedWorkflows[referencedWorkflow.Name] = referencedWorkflow
		}
	}

	return mainWorkflow, referencedWorkflows, nil
}
