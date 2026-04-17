# Real-World Example

This example demonstrates a complete real-world usage of gosaidno with proper project structure and best practices.

## Key Concepts Demonstrated

### 1. Proper Project Structure
- Service layer with business logic
- Centralized AOP setup
- Wrapped functions organized by service
- Multiple approaches for organizing wrapped functions

### 2. Cross-Cutting Concerns
- **Logging**: Comprehensive request/response logging
- **Timing**: Performance measurement for all operations
- **Validation**: Input validation with detailed error messages
- **Caching**: Around advice for intelligent caching
- **Error Handling**: Panic recovery and error logging

### 3. Naming Conventions
- Function names follow format: `ServiceName.MethodName` (e.g., `"UserService.GetUser"`)
- Wrapped functions follow format: `[Service][Method]Wrapped` (e.g., `UserServiceGetUser`)
- Consistent naming makes debugging and monitoring easier

### 4. Organization Patterns
- **Centralized Setup**: All AOP configuration in one place
- **Service Structs**: Group related wrapped functions together
- **Global Variables**: Direct access to wrapped functions
- **Dependency Injection Ready**: Easy to inject wrapped services

## Real-World Usage Patterns

### Service Layer Pattern
```go
// Original service method
func (us *UserService) GetUser(username string) (*User, error) {
    // Business logic
}

// Wrapped function
UserServiceGetUser := func(username string) (*User, error) {
    builder := aspect.For("UserService.GetUser")
    return wrap.Wrap1RE[string, *User](
        builder.GetRegistry(), 
        builder.GetFuncKey(), 
        (&UserService{}).GetUser,
    )(username)
}
```

### Struct-Based Organization
```go
type WrappedUserService struct {
    GetUser    func(string) (*User, error)
    CreateUser func(*User) error
}

// Create organized service wrapper
wrappedServices := &WrappedUserService{
    GetUser: func(username string) (*User, error) {
        builder := aspect.For("UserService.GetUser")
        return wrap.Wrap1RE[string, *User](
            builder.GetRegistry(), 
            builder.GetFuncKey(), 
            (&UserService{}).GetUser,
        )(username)
    },
    // ... other methods
}
```

## Best Practices Shown

1. **One-Time Setup**: AOP configured once at application startup
2. **Separation of Concerns**: Business logic separate from cross-cutting concerns
3. **Consistent Naming**: Clear, descriptive function names
4. **Error Handling**: Proper error propagation and logging
5. **Performance Monitoring**: Built-in timing for all operations
6. **Input Validation**: Early validation to prevent invalid operations
7. **Caching Strategy**: Intelligent caching with cache-aside pattern

## Running the Example

```bash
go run docs/examples/07_real_world_example/main.go
```

This example shows how gosaidno can be integrated into a real application with multiple services, demonstrating the power and flexibility of the fluent API for managing complex cross-cutting concerns.