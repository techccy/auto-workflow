package operations

import (
	"context"
)

// Handler defines the interface for operation handlers
type Handler interface {
	// Execute performs the operation with given parameters
	Execute(ctx context.Context, params map[string]interface{}, execCtx interface{}) (map[string]interface{}, error)

	// Validate checks if the parameters are valid for this operation
	Validate(params map[string]interface{}) error

	// GetType returns the operation type this handler supports
	GetType() string
}

// HandlerRegistry manages operation handlers
type HandlerRegistry struct {
	handlers map[string]Handler
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]Handler),
	}
}

// Register registers a new handler for an operation type
func (r *HandlerRegistry) Register(handler Handler) {
	r.handlers[handler.GetType()] = handler
}

// Get retrieves a handler by operation type
func (r *HandlerRegistry) Get(opType string) (Handler, bool) {
	handler, exists := r.handlers[opType]
	return handler, exists
}

// List returns all registered operation types
func (r *HandlerRegistry) List() []string {
	types := make([]string, 0, len(r.handlers))
	for t := range r.handlers {
		types = append(types, t)
	}
	return types
}
