package modulex

import "context"

// Application defines the interface for the main application
type Application interface {
	// Bannder
	Bannder() string
	// Name
	Name() string
	// Init initializes the application
	Init(context.Context) error
	// Start starts the application
	Start(context.Context) error
	// Stop stops the application
	Stop(context.Context) error
	// AddModule adds a module to the application
	AddModule(Module) error
	// GetModule returns a module by name
	GetModule(name string) (Module, bool)

	// AddHook adds a hook to a specific lifecycle phase
	Hookable
}
