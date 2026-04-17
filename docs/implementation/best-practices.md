# Best Practices

Based on the implementation details, here are best practices for using gosaidno effectively and avoiding common pitfalls.

## When to Use gosaidno

### Ideal Use Cases
- **Cross-cutting concerns**: Logging, monitoring, authentication, authorization
- **Consistent patterns**: Error handling, timing, caching, rate limiting
- **Team collaboration**: Shared AOP rules across codebase
- **Performance-critical code**: When reflection overhead matters
- **Clean architecture**: Separating business logic from infrastructure concerns

### When to Avoid gosaidno
- **Simple applications**: Overhead for basic use cases without complex cross-cutting concerns
- **Dynamic requirements**: When advice needs to change frequently at runtime
- **Method-heavy designs**: When most logic is in methods (requires conversion)
- **Very high-frequency calls**: When every microsecond of overhead matters
- **Small codebases**: Where the setup overhead isn't justified

## Implementation-Specific Best Practices

### 1. Initialize Early
Set up all AOP configuration during application startup:
```go
// aop/setup.go
func Init() {
    setupLogging()
    setupAuthentication()
    setupCaching()
    setupErrorHandling()
}

func setupLogging() {
    aspect.MustRegister("UserService.GetUser")
    aspect.MustAddAdvice("UserService.GetUser", loggingAdvice())
}
```

### 2. Use Meaningful Names
Function names appear in logs and debugging, so use descriptive names:
```go
// Good: Descriptive and hierarchical
aspect.Register("UserService.GetUserByID")
aspect.Register("PaymentService.ProcessPayment")

// Avoid: Generic or unclear names
aspect.Register("func1")
aspect.Register("api")
```

### 3. Mind the Priority System
Coordinate priority ranges across teams and modules:
```go
// Establish consistent priority ranges
const (
    PriorityAuth       = 1000  // Authentication runs first
    PriorityAuthz      = 900   // Authorization after auth
    PriorityRateLimit  = 800   // Rate limiting after auth
    PriorityLogging    = 100   // Logging runs last
)
```

### 4. Handle Errors Carefully
Advice errors can halt execution, so handle them appropriately:
```go
aspect.AddAdvice("MyFunc", aspect.Advice{
    Type: aspect.Before,
    Handler: func(c *aspect.Context) error {
        // Always validate inputs
        if c.Args == nil || len(c.Args) < 1 {
            return errors.New("missing required arguments")
        }
        
        // Safe type assertions
        userID, ok := c.Args[0].(int)
        if !ok {
            return errors.New("invalid user ID type")
        }
        
        // Business logic with error handling
        user, err := getUser(userID)
        if err != nil {
            return fmt.Errorf("failed to get user: %w", err)
        }
        
        c.Metadata["user"] = user
        return nil
    },
})
```

### 5. Monitor Memory Usage
The registry grows monotonically, so monitor it in long-running applications:
```go
// For applications with many dynamic registrations
// Consider using local registries for temporary scenarios
localRegistry := aspect.NewRegistry()
// Use local registry for specific operations
// Registry will be garbage collected when out of scope
```

### 6. Use Metadata Safely
Since metadata is untyped, establish conventions and validate:
```go
// Define constants for metadata keys
const (
    MetadataUser     = "user"
    MetadataStart    = "start_time"
    MetadataAttempts = "retry_attempts"
)

// Safe metadata access
func getUserFromContext(c *aspect.Context) (*User, error) {
    userData, exists := c.Metadata[MetadataUser]
    if !exists {
        return nil, errors.New("user not found in context")
    }
    
    user, ok := userData.(*User)
    if !ok {
        return nil, errors.New("invalid user type in context")
    }
    
    return user, nil
}
```

### 7. Test AOP Separately
Test advice functions independently and integration tests for the whole system:
```go
func TestLoggingAdvice(t *testing.T) {
    executed := false
    var capturedLog string
    
    // Mock logger
    logCapture := func(format string, args ...interface{}) {
        capturedLog = fmt.Sprintf(format, args...)
        executed = true
    }
    
    // Create advice with mock
    advice := loggingAdviceWithLogger(logCapture)
    
    // Test the advice directly
    c := aspect.NewContext("TestFunc", "arg1")
    err := advice.Handler(c)
    
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    
    if !executed {
        t.Error("Expected logging to execute")
    }
}
```

### 8. Profile Performance Impact
Measure the actual overhead in your specific use case:
```go
func BenchmarkAOPOverhead(b *testing.B) {
    // Setup AOP once
    aspect.MustRegister("BenchmarkFunc")
    aspect.MustAddAdvice("BenchmarkFunc", loggingAdvice())
    aspect.MustAddAdvice("BenchmarkFunc", timingAdvice())
    
    wrappedFunc := wrap.Wrap0("BenchmarkFunc", func() {
        // Simulate actual work
        time.Sleep(time.Microsecond)
    })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wrappedFunc()
    }
}
```

### 9. Document Your AOP Configuration
Maintain documentation about which advice applies to which functions:
```go
// AOP Configuration Documentation
//
// UserService.GetUser:
// - LoggingAdvice (Priority 100) - Logs entry/exit
// - AuthAdvice (Priority 200) - Validates authentication
// - CacheAdvice (Priority 150) - Implements caching
//
// PaymentService.ProcessPayment:
// - ValidationAdvice (Priority 300) - Validates payment data
// - LoggingAdvice (Priority 100) - Logs payment processing
```

### 10. Plan for Evolution
Structure your AOP setup to accommodate changes:
```go
// Use configuration-driven AOP setup
type AOPConfig struct {
    FunctionName string
    Advices      []aspect.Advice
}

func ApplyConfig(configs []AOPConfig) error {
    for _, config := range configs {
        if err := aspect.Register(config.FunctionName); err != nil {
            return err
        }
        
        for _, advice := range config.Advices {
            if err := aspect.AddAdvice(config.FunctionName, advice); err != nil {
                return err
            }
        }
    }
    return nil
}
```

Following these best practices will help you leverage gosaidno effectively while avoiding common pitfalls and performance issues.