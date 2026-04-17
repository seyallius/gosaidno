# Limitations and Trade-offs

Understanding the limitations and trade-offs of gosaidno helps you determine when it's the right tool for your use case and how to work within its constraints effectively.

## Architectural Limitations

### No Method Interception
**Limitation**: Cannot directly intercept methods on structs
**Workaround**: Convert methods to function values before wrapping
```go
type Service struct{}
func (s *Service) Method(x int) string { return "result" }

service := &Service{}
methodFunc := func(x int) string { return service.Method(x) }
wrapped := wrap.Wrap1R[int, string]("Service.Method", methodFunc)
```

### Static Registration Requirement
**Limitation**: Functions must be registered before use
**Impact**: Cannot dynamically add AOP to arbitrary functions at runtime
**Design Reason**: Enables performance optimizations and type safety

### Limited Return Type Support
**Limitation**: Only supports common function signatures (up to 3 args, standard return patterns)
**Impact**: Complex return types require custom wrappers
**Workaround**: Use struct return types to group multiple values

### No Conditional Advice
**Limitation**: Advice runs based on type, not runtime conditions
**Impact**: Cannot selectively apply advice based on function parameters
**Alternative**: Implement conditional logic inside advice functions

## Design Trade-offs

### Transparency vs Convenience
**Trade-off**: Explicit registration and wrapping vs annotation convenience
- **Pros**: Clear execution flow, easy debugging, no build-time dependencies
- **Cons**: More verbose setup, requires manual intervention

### Performance vs Flexibility
**Trade-off**: Compile-time generics vs runtime reflection
- **Pros**: Type safety, better performance, compile-time checking
- **Cons**: Limited to supported function signatures, more code generation

### Safety vs Power
**Trade-off**: Type safety vs dynamic behavior
- **Pros**: Compile-time error detection, IDE support, fewer runtime errors
- **Cons**: Less dynamic capability, more restrictive

### Simplicity vs Features
**Trade-off**: Minimal API vs extensive functionality
- **Pros**: Easier to understand, maintain, and reason about
- **Cons**: Fewer advanced features, less configurability

## Implementation-Specific Constraints

### Context Limitations
- **Untyped Args/Results**: Requires careful type assertions
- **Metadata Safety**: No compile-time checking for metadata keys/values
- **Memory Overhead**: Context created per function call

### Priority System Constraints
- **Manual Management**: No automatic priority assignment
- **Collision Risk**: Multiple teams may use overlapping ranges
- **Sorting Overhead**: O(n log n) complexity per advice type execution

### Registry Behavior
- **Monotonic Growth**: Registry only grows, never shrinks
- **Global State**: Default registry is global, may cause conflicts
- **Thread Safety**: Synchronization overhead for concurrent access

## When These Limitations Matter

### High-Performance Scenarios
- Very frequent function calls where microseconds matter
- Systems with strict memory constraints
- Real-time processing requirements

### Dynamic Requirements
- Applications that need to modify AOP behavior at runtime
- Plugin architectures with dynamic function interception
- Systems requiring complex conditional advice

### Complex Signatures
- Functions with many parameters (>3)
- Functions with complex return value combinations
- Variadic functions

## Working Within Constraints

### Naming Conventions
Establish clear naming for functions and metadata to avoid collisions:
```go
// Good: Descriptive names with namespace
aspect.Register("UserService.GetUserByID")
c.Metadata["user_service.authenticated_user"]
```

### Priority Management
Coordinate priority ranges across teams:
- 100-199: Authentication
- 200-299: Logging
- 300-399: Validation
- 400-499: Caching
- 500-599: Error handling

### Error Handling
Design robust error handling within advice:
```go
Handler: func(c *aspect.Context) error {
    // Validate context state
    if c.Args == nil || len(c.Args) == 0 {
        return errors.New("missing required arguments")
    }
    
    // Safe type assertion
    if userID, ok := c.Args[0].(int); !ok {
        return errors.New("invalid user ID type")
    }
    
    // Business logic
    return nil
}
```

## Alternatives for Different Needs

If gosaidno's limitations don't suit your needs, consider:

- **Decorator pattern**: Manual function composition for simple cases
- **Middleware**: For HTTP handlers or similar request-response patterns
- **Interface-based composition**: Using interfaces for dependency injection
- **Code generation**: Tools like `go generate` for compile-time AOP

Understanding these limitations and trade-offs helps you make informed decisions about when and how to use gosaidno effectively in your applications.