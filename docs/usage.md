# Usage Guide

This guide provides comprehensive information about using gosaidsno effectively in your Go applications.

## Core Concepts

### Function Registration

Before you can apply advice to a function, you must register it with gosaidsno:

```go
import "github.com/seyallius/gosaidno/aspect"

// Register a function with a unique name
err := aspect.Register("UserService.GetUser")

// Or use MustRegister which panics on error
aspect.MustRegister("UserService.GetUser")
```

Registration creates an entry in the internal registry that associates the function name with an advice chain. If you try to register the same function twice, it will return an error.

### Advice Types

gosaidsno supports five distinct types of advice, each executing at different points in the function lifecycle:

#### Before Advice

Executes before the target function. Useful for logging, validation, authentication, etc.

```go
aspect.MustAddAdvice("UserService.GetUser", aspect.Advice{
    Type:     aspect.Before,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        log.Printf("About to call %s", c.FunctionName)
        return nil
    },
})
```

#### After Advice

Executes after the target function, regardless of whether it succeeded or panicked. Always runs. Useful for cleanup, logging completion, releasing resources, etc.

```go
aspect.MustAddAdvice("UserService.GetUser", aspect.Advice{
    Type:     aspect.After,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        log.Printf("Finished calling %s", c.FunctionName)
        return nil
    },
})
```

#### Around Advice

Wraps the target function execution. Can skip the target function entirely or modify arguments/results. Most powerful but also most complex advice type.

```go
aspect.MustAddAdvice("UserService.GetUser", aspect.Advice{
    Type:     aspect.Around,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        log.Printf("Around advice: About to call %s", c.FunctionName)

        // The target function executes here

        log.Printf("Around advice: Finished calling %s", c.FunctionName)
        return nil
    },
})
```

#### AfterReturning Advice

Executes only if the target function returns successfully (no panic). Useful for post-processing successful results, caching, etc.

```go
aspect.MustAddAdvice("UserService.GetUser", aspect.Advice{
    Type:     aspect.AfterReturning,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        if c.Error == nil {
            log.Printf("Function %s succeeded", c.FunctionName)
        }
        return nil
    },
})
```

#### AfterThrowing Advice

Executes only if the target function panics. Useful for error handling, cleanup on failure, etc.

```go
aspect.MustAddAdvice("UserService.GetUser", aspect.Advice{
    Type:     aspect.AfterThrowing,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        if c.PanicValue != nil {
            log.Printf("Function %s panicked: %v", c.FunctionName, c.PanicValue)
        }
        return nil
    },
})
```

### Priority System

Within each advice type, execution order is determined by priority. Higher priority values execute first:

```go
// This executes first (priority 200)
aspect.MustAddAdvice("MyFunc", aspect.Advice{
    Type:     aspect.Before,
    Priority: 200,
    Handler:  highPriorityHandler,
})

// This executes second (priority 100)
aspect.MustAddAdvice("MyFunc", aspect.Advice{
    Type:     aspect.Before,
    Priority: 100,
    Handler:  mediumPriorityHandler,
})

// This executes third (priority 50)
aspect.MustAddAdvice("MyFunc", aspect.Advice{
    Type:     aspect.Before,
    Priority: 50,
    Handler:  lowPriorityHandler,
})
```

### Context Object

The `Context` object is passed to every advice function and contains important information:

```go
type Context struct {
    FunctionName string         // Name of the registered function
    Args         []any          // Arguments passed to the function
    Results      []any          // Return values from the function
    Error        error          // Error returned by the function
    PanicValue   any            // Panic value if function panicked
    Metadata     map[string]any // Custom key-value storage for advice communication
    Skipped      bool           // Whether target function execution was skipped
}
```

## Function Wrapping

gosaidsno provides generic wrapper functions for different function signatures:

### No Arguments, No Return Values

```go
originalFunc := func() {
    // Your business logic
}

wrappedFunc := aspect.Wrap0("MyFunc", originalFunc)
```

### No Arguments, One Return Value

```go
originalFunc := func() string {
    return "result"
}

wrappedFunc := aspect.Wrap0R[string]("MyFunc", originalFunc)
```

### No Arguments, Result and Error

```go
originalFunc := func() (string, error) {
    return "result", nil
}

wrappedFunc := aspect.Wrap0RE[string]("MyFunc", originalFunc)
```

### One Argument, No Return Values

```go
originalFunc := func(userID int) {
    // Process user ID
}

wrappedFunc := aspect.Wrap1[int]("MyFunc", originalFunc)
```

### One Argument, One Return Value

```go
originalFunc := func(userID int) string {
    return fmt.Sprintf("user-%d", userID)
}

wrappedFunc := aspect.Wrap1R[int, string]("MyFunc", originalFunc)
```

### One Argument, Result and Error

```go
originalFunc := func(userID int) (string, error) {
    if userID <= 0 {
        return "", errors.New("invalid user ID")
    }
    return fmt.Sprintf("user-%d", userID), nil
}

wrappedFunc := aspect.Wrap1RE[int, string]("MyFunc", originalFunc)
```

### Multiple Arguments

gosaidsno supports functions with up to 3 arguments:

```go
// Two arguments, result and error
wrappedFunc := aspect.Wrap2RE[string, int, User]("MyFunc",
    func(username string, age int) (User, error) {
        // Implementation
    })

// Three arguments, result and error
wrappedFunc := aspect.Wrap3RE[string, int, bool, User]("MyFunc",
    func(username string, age int, active bool) (User, error) {
        // Implementation
    })
```

## Advanced Patterns

### Using Metadata for Communication

Advice functions can communicate with each other using the context's metadata field:

```go
// Authentication advice stores user info in metadata
aspect.MustAddAdvice("UserService.GetUser", aspect.Advice{
    Type:     aspect.Before,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        token := c.Args[0].(string) // Assuming first arg is token
        user, err := authenticate(token)
        if err != nil {
            return err
        }
        c.Metadata["authenticatedUser"] = user
        return nil
    },
})

// Authorization advice reads from metadata
aspect.MustAddAdvice("UserService.GetUser", aspect.Advice{
    Type:     aspect.Before,
    Priority: 90, // Lower priority, runs after auth
    Handler: func(c *aspect.Context) error {
        user := c.Metadata["authenticatedUser"].(*User)
        if user.Role != "admin" {
            return errors.New("insufficient permissions")
        }
        return nil
    },
})
```

### Caching with Around Advice

Around advice can skip target function execution entirely:

```go
var cache = make(map[string]interface{})

aspect.MustAddAdvice("ExpensiveCalculation", aspect.Advice{
    Type:     aspect.Around,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        key := fmt.Sprintf("%v", c.Args[0]) // Simple key from first arg

        if cached, exists := cache[key]; exists {
            // Found in cache, skip target function
            c.SetResult(0, cached)
            c.Skipped = true
            return nil
        }

        // Not in cache, let target function execute
        // Result will be stored and available after function execution
        return nil
    },
})

aspect.MustAddAdvice("ExpensiveCalculation", aspect.Advice{
    Type:     aspect.AfterReturning,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        if !c.Skipped {
            // Cache the result only if function wasn't skipped
            key := fmt.Sprintf("%v", c.Args[0])
            cache[key] = c.Results[0]
        }
        return nil
    },
})
```

### Error Recovery and Retry Logic

Combine multiple advice types for robust error handling:

```go
aspect.MustAddAdvice("ExternalAPICall", aspect.Advice{
    Type:     aspect.Before,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        // Initialize retry counter
        c.Metadata["retryCount"] = 0
        return nil
    },
})

aspect.MustAddAdvice("ExternalAPICall", aspect.Advice{
    Type:     aspect.AfterReturning,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        if c.Error != nil {
            retryCount := c.Metadata["retryCount"].(int)
            if retryCount < 3 {
                c.Metadata["retryCount"] = retryCount + 1
                // In a real implementation, you'd trigger a retry here
            }
        }
        return nil
    },
})
```

## Best Practices

### 1. Centralized Setup

Set up all your AOP configuration in one place, typically during application initialization:

```go
// aop/setup.go
package aop

import "github.com/seyallius/gosaidno/aspect"

func Init() {
    setupLogging()
    setupAuthentication()
    setupCaching()
    setupErrorHandling()
}

func setupLogging() {
    // Register functions and add logging advice
    aspect.MustRegister("UserService.GetUser")
    aspect.MustAddAdvice("UserService.GetUser", loggingAdvice())
}

func setupAuthentication() {
    // Similar setup for auth
}
```

### 2. Error Handling in Advice

Always handle errors in your advice functions appropriately:

```go
aspect.MustAddAdvice("MyFunc", aspect.Advice{
    Type:     aspect.Before,
    Priority: 100,
    Handler: func(c *aspect.Context) error {
        // Always return an error if something goes wrong
        if someCondition {
            return errors.New("validation failed")
        }
        return nil
    },
})
```

### 3. Performance Considerations

Be mindful of the performance impact of advice:

- Minimize heavy computations in advice functions
- Use caching when appropriate
- Profile your application to understand the overhead

### 4. Testing

Test both your advice functions and the interaction between advice and target functions:

```go
func TestLoggingAdvice(t *testing.T) {
    var logOutput string

    // Set up test advice that captures log output
    aspect.MustAddAdvice("TestFunc", aspect.Advice{
        Type:     aspect.Before,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            logOutput = fmt.Sprintf("Called %s", c.FunctionName)
            return nil
        },
    })

    // Test the wrapped function
    wrappedFunc := aspect.Wrap0("TestFunc", func() {})
    wrappedFunc()

    if logOutput != "Called TestFunc" {
        t.Errorf("Expected logging, got %s", logOutput)
    }
}
```

## Common Use Cases

### Logging and Monitoring

```go
func loggingAdvice() aspect.Advice {
    return aspect.Advice{
        Type:     aspect.Around,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            start := time.Now()
            log.Printf("Starting %s with args: %v", c.FunctionName, c.Args)

            // Function executes here

            duration := time.Since(start)
            log.Printf("Completed %s in %v, result: %v, error: %v",
                      c.FunctionName, duration, c.Results, c.Error)
            return nil
        },
    }
}
```

### Authentication and Authorization

```go
func authAdvice() aspect.Advice {
    return aspect.Advice{
        Type:     aspect.Before,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            token := extractToken(c.Args)
            user, err := validateToken(token)
            if err != nil {
                return errors.New("unauthorized")
            }
            c.Metadata["user"] = user
            return nil
        },
    }
}
```

### Rate Limiting

```go
func rateLimitingAdvice() aspect.Advice {
    return aspect.Advice{
        Type:     aspect.Before,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            user := c.Metadata["user"].(*User)
            if !rateLimiter.Allow(user.ID) {
                return errors.New("rate limit exceeded")
            }
            return nil
        },
    }
}
```

## Fluent API

gosaidsno now includes a fluent/declarative API that provides a more convenient and readable way to configure advice:

### Basic Usage

Instead of manually registering functions and adding advice separately, you can use the fluent API:

```go
// Old way
aspect.MustRegister("GetUser")
aspect.MustAddAdvice("GetUser", aspect.Advice{
    Type:     aspect.Before,
    Handler:  authCheck,
})
aspect.MustAddAdvice("GetUser", aspect.Advice{
    Type:     aspect.After,
    Handler:  logging,
})

// New fluent way
aspect.For("GetUser").
    WithBefore(authCheck).
    WithAfter(logging)
```

### Fluent API Methods

The fluent API provides methods for all advice types:

- `WithBefore(handler)` - Add Before advice
- `WithAfter(handler)` - Add After advice  
- `WithAround(handler)` - Add Around advice
- `WithAfterReturning(handler)` - Add AfterReturning advice
- `WithAfterThrowing(handler)` - Add AfterThrowing advice

Each method also has a priority variant:

- `WithBeforeP(handler, priority)` - Add Before advice with priority
- `WithAfterP(handler, priority)` - Add After advice with priority
- `WithAroundP(handler, priority)` - Add Around advice with priority
- `WithAfterReturningP(handler, priority)` - Add AfterReturning advice with priority
- `WithAfterThrowingP(handler, priority)` - Add AfterThrowing advice with priority

### Using Custom Registries

You can also use the fluent API with custom registries:

```go
registry := aspect.NewRegistry()
aspect.ForWithRegistry(registry, "GetUser").
    WithBefore(authCheck).
    WithAfter(logging)
```

### Combining with Function Wrapping

After configuring advice with the fluent API, wrap your functions using the registry:

```go
// Configure advice
aspect.For("GetUser").
    WithBefore(authCheck).
    WithAfter(logging).
    WithAround(caching)

// Then wrap your function using the builder
builder := aspect.For("GetUser")
wrappedFn := aspect.Wrap1RE[string,*User](builder.GetRegistry(), builder.GetFuncKey(), getUserImpl)
```

### Complete Example

Here's a complete example using the fluent API:

```go
package main

import "github.com/seyallius/gosaidno/aspect"

func main() {
    // Configure advice using fluent API
    aspect.For("GetUser").
        WithBefore(func(c *aspect.Context) error {
            // Authentication check
            return nil
        }).
        WithAfter(func(c *aspect.Context) error {
            // Logging
            return nil
        }).
        WithAround(func(c *aspect.Context) error {
            // Caching logic
            return nil
        })

    // Wrap your function
    builder := aspect.For("GetUser")
    wrappedGetUser := aspect.Wrap1RE[string,*User](
        builder.GetRegistry(), 
        builder.GetFuncKey(), 
        getUserImpl,
    )

    // Use the wrapped function normally
    user, err := wrappedGetUser("user123")
}

func getUserImpl(id string) (*User, error) {
    // Your business logic here
    return &User{ID: id}, nil
}

type User struct {
    ID string
}
```

## Best Practices with Fluent API

### 1. Group Related Configuration

Use the fluent API to group related advice configuration:

```go
// Configure all security-related advice together
aspect.For("SensitiveOperation").
    WithBefore(authentication).
    WithBefore(authorization).
    WithAfter(auditLog)
```

### 2. Use Priority Variants for Ordering

When you need specific execution order, use the priority variants:

```go
aspect.For("CriticalOperation").
    WithBeforeP(highPriorityValidation, 100).
    WithBeforeP(normalValidation, 50).
    WithBeforeP(lowPriorityValidation, 10)
```

### 3. Combine with Centralized Setup

Integrate the fluent API into your centralized setup:

```go
func setupSecurity() {
    aspect.For("UserService.GetUser").
        WithBefore(authentication).
        WithBefore(authorization)
    
    aspect.For("PaymentService.Process").
        WithBefore(fraudDetection).
        WithAfter(paymentAudit)
}
```

This guide covers the essential aspects of using gosaidsno. For more specific examples, check out the [Examples](../examples/README.md) directory.