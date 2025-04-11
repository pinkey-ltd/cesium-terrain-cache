package cache

import (
	"os"
)

// Persistence represents the mechanism for persisting the cache commands to a file.
// It uses an append-only file (AOF) to persist commands that modify the cache state.
type Persistence struct {
	file *os.File // The file used to store the commands
}

// NewAOF initializes and returns a new Persistence instance for appending commands to a file.
// It opens the specified file in append mode, creating it if it doesn't exist.
func NewAOF(filename string) (*Persistence, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Persistence{file: file}, nil
}

// Append writes a command to the AOF file.
// The command is written as a string followed by a newline character.
func (p *Persistence) Append(cmd string) error {
	_, err := p.file.WriteString(cmd + "\n")
	return err
}

// Close closes the AOF file when the server is shutting down.
func (p *Persistence) Close() {
	p.file.Close()
}
