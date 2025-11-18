package selector

import (
	"bytes"
	"os/exec"
)

// SystemCommand is an interface for executing system commands
// This allows us to mock system commands for testing
type SystemCommand interface {
	// Run executes a command and returns the output
	Run(name string, args ...string) ([]byte, error)

	// RunInteractive executes a command without capturing output
	RunInteractive(name string, args ...string) error
}

// RealSystemCommand implements SystemCommand using actual os/exec
type RealSystemCommand struct{}

// NewRealSystemCommand creates a new real system command executor
func NewRealSystemCommand() SystemCommand {
	return &RealSystemCommand{}
}

// Run executes a command and returns the output
func (r *RealSystemCommand) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}

// RunInteractive executes a command without capturing output
func (r *RealSystemCommand) RunInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// MockSystemCommand is a mock implementation for testing
type MockSystemCommand struct {
	// Output maps command names to their mock output
	Output map[string][]byte

	// Errors maps command names to errors they should return
	Errors map[string]error

	// CallLog records all commands that were executed
	CallLog []CommandCall
}

// CommandCall records a command execution
type CommandCall struct {
	Name string
	Args []string
}

// NewMockSystemCommand creates a new mock system command executor
func NewMockSystemCommand() *MockSystemCommand {
	return &MockSystemCommand{
		Output:  make(map[string][]byte),
		Errors:  make(map[string]error),
		CallLog: make([]CommandCall, 0),
	}
}

// Run executes a mock command and returns configured output
func (m *MockSystemCommand) Run(name string, args ...string) ([]byte, error) {
	m.CallLog = append(m.CallLog, CommandCall{Name: name, Args: args})

	if err, ok := m.Errors[name]; ok {
		return nil, err
	}

	if output, ok := m.Output[name]; ok {
		return output, nil
	}

	return []byte{}, nil
}

// RunInteractive executes a mock command without output
func (m *MockSystemCommand) RunInteractive(name string, args ...string) error {
	m.CallLog = append(m.CallLog, CommandCall{Name: name, Args: args})

	if err, ok := m.Errors[name]; ok {
		return err
	}

	return nil
}

// SetOutput configures the output for a command
func (m *MockSystemCommand) SetOutput(name string, output []byte) {
	m.Output[name] = output
}

// SetError configures an error for a command
func (m *MockSystemCommand) SetError(name string, err error) {
	m.Errors[name] = err
}

// GetCallCount returns the number of times a command was called
func (m *MockSystemCommand) GetCallCount(name string) int {
	count := 0
	for _, call := range m.CallLog {
		if call.Name == name {
			count++
		}
	}
	return count
}

// WasCalled checks if a command was called with specific arguments
func (m *MockSystemCommand) WasCalled(name string, args ...string) bool {
	for _, call := range m.CallLog {
		if call.Name != name {
			continue
		}
		if len(call.Args) != len(args) {
			continue
		}
		match := true
		for i, arg := range args {
			if call.Args[i] != arg {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// Reset clears the call log and configured outputs/errors
func (m *MockSystemCommand) Reset() {
	m.Output = make(map[string][]byte)
	m.Errors = make(map[string]error)
	m.CallLog = make([]CommandCall, 0)
}
