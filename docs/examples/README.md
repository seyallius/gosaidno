# Examples Directory

## Running Examples

```bash
# Basic usage - logging, timing, validation
go run examples/01_basic_usage/main.go

# Caching with Around advice - Broken... needs fix
go run examples/02_caching_pattern/main.go

# Authentication and authorization
go run examples/03_authentication/main.go

# Circuit breaker for fault tolerance
go run examples/04_circuit_breaker/main.go

# Retry with exponential backoff
go run examples/05_retry_pattern/main.go

# Fluent API for declarative advice configuration
go run examples/06_fluent_api/main.go

# Real-world example with proper project structure
go run examples/07_real_world_example/main.go

# Context example demonstrating context propagation through AOP
go run examples/08_context_example/main.go
```

## Examples Overview

### 01_basic_usage

**Real-world use cases:**

- Request/response logging
- Method execution timing
- Input validation (Before advice)
- Success-only actions (AfterReturning)

**Key patterns:**

- Setup AOP once at startup
- Wrap functions during initialization
- Use metadata to pass data between advice

### 02_caching_pattern - Broken... needs fix

**Real-world use cases:**

- Database query caching
- API response caching
- TTL-based cache expiration
- Cache hit/miss metrics

**Key patterns:**

- Around advice checks cache
- Skip execution on cache hit
- AfterReturning populates cache

### 03_authentication

**Real-world use cases:**

- Token validation
- Permission checks
- Session management

**Key patterns:**

- Before advice for auth checks
- Early termination on auth failure
- Secure parameter access

### 04_circuit_breaker

**Real-world use cases:**

- Fault tolerance for external services
- Prevent cascading failures
- Graceful degradation

**Key patterns:**

- Around advice controls execution flow
- State management for circuit breaker
- Failure counting and timeouts

### 05_retry_pattern

**Real-world use cases:**

- Network resilience
- Transient error handling
- Automatic recovery

**Key patterns:**

- Wrapper functions for retry logic
- Exponential backoff
- Conditional retry based on error types

### 06_fluent_api

**Real-world use cases:**

- Declarative advice configuration
- Cleaner API for aspect setup
- Improved readability and maintainability

**Key patterns:**

- Fluent builder pattern for advice configuration
- Method chaining for multiple advice types
- Type-safe configuration without reflection

### 07_real_world_example

**Real-world use cases:**

- Complete application structure with services
- Multiple cross-cutting concerns (logging, timing, validation, caching)
- Proper organization of wrapped functions
- Realistic error handling and recovery

**Key patterns:**

- Centralized AOP setup
- Service layer separation
- Multiple organization approaches (globals vs structs)
- Realistic validation and error handling

### 08_context_example

**Real-world use cases:**

- Context value propagation through advice functions
- Context deadline awareness in AOP workflows
- Context cancellation detection and handling
- Integration with Go's context package for cancellation and timeouts

**Key patterns:**

- Context propagation through all advice types
- Deadline checking in advice functions
- Cancellation signal detection
- Proper error handling when context expires

## Project Setup Pattern

All examples follow this pattern:

```go
// 1. Setup AOP once at startup
func setupAOP() {
    aspect.MustRegister("FunctionName")
    aspect.MustAddAdvice("FunctionName", /* ... */)
}

// 2. Implement business logic
func businessLogicImpl(args...) result {
    // Pure business logic
}

// 3. Wrap functions (once, during initialization)
var BusinessLogic = wrap.Wrap*("FunctionName", businessLogicImpl)

// 4. Use normally throughout application
func main() {
    setupAOP()  // Once

    // Use wrapped functions
    result := BusinessLogic(args...)
}
```

## Advice Execution Order

For a function with all advice types:

1. **Before** (high priority → low)
2. **Around** (can skip step 3)
3. Target function
4. **AfterReturning** (only if success)
5. **AfterThrowing** (only if panic)
6. **After** (always runs)

## Common Metadata Keys

Examples use these metadata conventions:

- `startTime` - For timing (type: `time.Time`)
- `userID` - For auth/audit (type: `string`)
- `role` - For authorization (type: `string`)
- `attempt` - For retry logic (type: `int`)
- `maxRetries` - Retry configuration (type: `int`)

## Performance Notes

Each example includes:

- Timing measurements
- Success/failure counts
- Performance impact of advice
- Cache hit rates (where applicable)