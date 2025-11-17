# Testing Guide

This document describes the test suite for the Witness screen capture tool.

## Overview

The test suite provides comprehensive coverage of all core functionality with a focus on:
- Unit testing of individual components
- Integration testing with mock implementations
- Platform-agnostic testing for macOS-specific code
- Fixture-based testing to prevent regressions

## Test Coverage

Current test coverage by package:

- **pkg/capture**: 92.6% - Core capture interfaces and mock implementations
- **pkg/encoder**: 94.3% - GIF encoding with multiple quality levels
- **pkg/selector**: 53.4% - Region selection, parsing, and configuration management

## Running Tests

### Run All Tests
```bash
make test
```

### Run Tests with Verbose Output
```bash
make test-verbose
```

### Run Tests with Coverage Report
```bash
make test-coverage
```

### Run Tests for a Specific Package
```bash
go test ./pkg/capture/... -v
go test ./pkg/encoder/... -v
go test ./pkg/selector/... -v
```

## Test Structure

### Package: `pkg/capture`

**Files:**
- `capture_test.go` - Tests for Region, Config, and Frame structs
- `mock_capturer.go` - Mock implementation of the Capturer interface
- `mock_capturer_test.go` - Tests for the mock capturer

**Key Features Tested:**
- Region validation and configuration
- Frame capture and timestamp handling
- Mock capturer with configurable behavior
- Frame generation with custom colors and patterns
- Error simulation for testing error handling paths

### Package: `pkg/encoder`

**Files:**
- `gif_test.go` - Comprehensive GIF encoder tests

**Key Features Tested:**
- GIF encoder initialization with various FPS and quality settings
- Frame addition and validation
- Multi-frame GIF encoding
- Quality level impact on palette selection
- Frame count tracking
- File size estimation
- Error handling (nil frames, invalid paths, no frames)
- Different frame sizes and formats

**Test Helpers:**
- `createTestFrame()` - Creates solid color test frames
- `createGradientFrame()` - Creates gradient pattern frames for color testing

### Package: `pkg/selector`

**Files:**
- `selector_test.go` - Tests for region parsing and formatting
- `config_test.go` - Tests for region configuration management
- `selector_darwin_test.go` - Platform-specific selector tests with mocks
- `system_command.go` - System command wrapper interface for testing

**Key Features Tested:**
- Region string parsing (`x,y,w,h` format)
- Region string formatting
- Round-trip parsing and formatting
- Config file persistence (JSON)
- Multi-region management
- Default region selection
- Region CRUD operations (save, load, delete, list)
- macOS selector with mocked system commands
- System command execution mocking

**Test Helpers:**
- `setupTestConfig()` - Creates temporary config directories
- `MockSystemCommand` - Mocks system commands like `screencapture` and `defaults`

## Mocking Strategy

### macOS System Commands

The selector package uses a wrapper interface (`SystemCommand`) to abstract system command execution. This allows testing macOS-specific functionality on any platform:

```go
type SystemCommand interface {
    Run(name string, args ...string) ([]byte, error)
    RunInteractive(name string, args ...string) error
}
```

**Real Implementation:**
- `RealSystemCommand` - Uses `os/exec` for actual command execution

**Mock Implementation:**
- `MockSystemCommand` - Configurable mock with output/error injection
- Tracks all commands executed for verification
- Allows setting custom outputs and errors per command

### Screen Capture

The capture package provides a mock capturer (`MockCapturer`) for testing without actual screen capture:

**Features:**
- Configurable frame generation
- Customizable frame dimensions and colors
- FPS simulation
- Error injection
- Frame counting and limits
- Custom frame generation functions

## Writing New Tests

### Testing Guidelines

1. **Use Table-Driven Tests**: For testing multiple scenarios
   ```go
   tests := []struct {
       name    string
       input   string
       want    *Region
       wantErr bool
   }{
       // test cases...
   }
   ```

2. **Use Test Helpers**: Create helper functions for common setup
   ```go
   func setupTestConfig(t *testing.T) (string, func()) {
       // Setup and return cleanup function
   }
   ```

3. **Test Error Cases**: Always test error handling paths
   ```go
   if err == nil {
       t.Error("expected error for invalid input")
   }
   ```

4. **Clean Up Resources**: Use `defer` for cleanup
   ```go
   tmpDir, cleanup := setupTestConfig(t)
   defer cleanup()
   ```

5. **Use Descriptive Test Names**: Make failures easy to understand
   ```go
   t.Run("invalid format - missing value", func(t *testing.T) {
       // test code
   })
   ```

### Adding Tests for New Features

When adding new functionality:

1. Write tests first (TDD approach recommended)
2. Create mock implementations for external dependencies
3. Test both success and failure cases
4. Add integration tests if the feature spans multiple packages
5. Update this documentation with new test coverage

## Continuous Integration

The test suite is designed to run in CI environments:

- No macOS-specific requirements (all platform code is mocked)
- Deterministic test execution
- Fast execution (< 2 seconds for full suite)
- No external dependencies required

## Fixtures and Test Data

### Configuration Files

Tests create temporary config directories to avoid interfering with user data:
- Uses `os.TempDir()` for isolation
- Cleans up after test completion
- Sets `HOME` environment variable for config path testing

### Image Data

GIF encoder tests generate test frames programmatically:
- Solid color frames for basic testing
- Gradient frames for color palette testing
- Configurable dimensions for size testing

## Troubleshooting

### Tests Failing on macOS

If macOS-specific tests fail:
1. Verify the mock command outputs match actual macOS format
2. Check that system commands return expected output format
3. Update mocks if macOS behavior has changed

### Coverage Gaps

To identify untested code:
```bash
go test ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Race Conditions

To detect race conditions:
```bash
go test ./... -race
```

## Future Improvements

- [ ] Add integration tests for full capture-to-GIF pipeline
- [ ] Add benchmark tests for encoder performance
- [ ] Add tests for video encoding (when implemented)
- [ ] Increase selector coverage with more edge cases
- [ ] Add performance regression tests
- [ ] Add fuzzing tests for region parsing

## Contributing

When contributing tests:

1. Maintain or improve coverage percentages
2. Follow existing test patterns and naming conventions
3. Add documentation for complex test scenarios
4. Ensure tests are platform-agnostic when possible
5. Include both positive and negative test cases
