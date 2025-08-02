package cmd

import (
	"log"
)

// ServerCommand represents the server command
type ServerCommand struct {
	Port int
	Host string
}

// NewServerCommand creates a new server command instance
func NewServerCommand() *ServerCommand {
	return &ServerCommand{
		Port: 8080,
		Host: "localhost",
	}
}

// Run executes the server command
func (s *ServerCommand) Run() error {
	log.Printf("Starting server on %s:%d", s.Host, s.Port)
	// TODO: Implement server startup logic
	return nil
}
