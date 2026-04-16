# Quick Start Guide

This guide will help you get started with gosaidsno in just a few minutes. You'll learn how to set up AOP in your Go application and implement common cross-cutting concerns.

## Installation

First, add gosaidsno to your project:

```bash
go get github.com/seyallius/gosaidno
```

## Basic Setup

Here's a minimal example showing how to add logging to a function:

```go
package main

import (
    "fmt"
    "log"

    "github.com/seyallius/gosaidno/aspect"
)

func main() {
    // Step 1: Register your function
    err := aspect.Register("GreetUser")
    if err != nil {
        log.Fatal(err)
    }

    // Step 2: Add logging advice
    err = aspect.AddAdvice("GreetUser", aspect.Advice{
        Type:     aspect.Before,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            log.Printf("About to execute %s", c.FunctionName)
            return nil
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Add after advice
    err = aspect.AddAdvice("GreetUser", aspect.Advice{
        Type:     aspect.After,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            log.Printf("Finished executing %s", c.FunctionName)
            return nil
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Step 4: Wrap your function
    greetUser := aspect.Wrap1R[string]("GreetUser", func(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    })

    // Step 5: Use your enhanced function
    result := greetUser("World")
    fmt.Println(result)
}
```

## Understanding the Components

### 1. Registration
Every function you want to enhance with AOP must be registered first:

```go
err := aspect.Register("FunctionName")
```

Registration creates an entry in the internal registry that associates the function name with an advice chain.

### 2. Advice Types
gosaidsno supports five types of advice:

- **Before**: Executes before the target function
- **After**: Executes after the target function (always runs, even if the function panics)
- **Around**: Wraps the target function execution (can skip it)
- **AfterReturning**: Executes only if the function returns successfully (no panic/error)
- **AfterThrowing**: Executes only if the function panics

### 3. Priority System
Advice of the same type executes in priority order (higher numbers execute first):

```go
// This executes first (priority 200)
aspect.AddAdvice("MyFunc", aspect.Advice{
    Type:     aspect.Before,
    Priority: 200,
    Handler:  myHandler,
})

// This executes second (priority 100)
aspect.AddAdvice("MyFunc", aspect.Advice{
    Type:     aspect.Before,
    Priority: 100,
    Handler:  myOtherHandler,
})
```

### 4. Function Wrapping
Use the appropriate wrapper function based on your function signature:

- `Wrap0`: No arguments, no return values
- `Wrap0R`: No arguments, one return value
- `Wrap0RE`: No arguments, (result, error) return
- `Wrap1`: One argument, no return values
- `Wrap1R`: One argument, one return value
- `Wrap1RE`: One argument, (result, error) return
- `Wrap1E`: One argument, error return
- And more for functions with multiple arguments

## A More Complex Example

Let's create a function with error handling and timing:

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "math"
    "time"

    "github.com/seyallius/gosaidno/aspect"
)

func main() {
    // Register the function
    aspect.MustRegister("CalculateSquareRoot")

    // Add timing advice
    aspect.MustAddAdvice("CalculateSquareRoot", aspect.Advice{
        Type:     aspect.Around,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            startTime := time.Now()
            log.Printf("[%s] Starting execution", c.FunctionName)

            // Continue with function execution
            err := c.Next() // This would be the function call in a real implementation

            duration := time.Since(startTime)
            log.Printf("[%s] Completed in %v", c.FunctionName, duration)

            return err
        },
    })

    // Add error handling advice
    aspect.MustAddAdvice("CalculateSquareRoot", aspect.Advice{
        Type:     aspect.AfterThrowing,
        Priority: 100,
        Handler: func(c *aspect.Context) error {
            log.Printf("[%s] Function panicked with: %v", c.FunctionName, c.PanicValue)
            return nil
        },
    })

    // Wrap a function that calculates square root
    calculateSquareRoot := aspect.Wrap1RE[float64, float64]("CalculateSquareRoot",
        func(input float64) (float64, error) {
            if input < 0 {
                return 0, errors.New("cannot calculate square root of negative number")
            }
            return math.Sqrt(input), nil
        })

    // Use the enhanced function
    result, err := calculateSquareRoot(16)
    if err != nil {
        log.Printf("Error: %v", err)
    } else {
        fmt.Printf("Result: %f\n", result)
    }
}
```

## Best Practices

1. **Set up AOP once**: Initialize all your advice at application startup
2. **Use meaningful names**: Choose clear, descriptive names for your registered functions
3. **Handle errors gracefully**: Make sure your advice functions handle errors appropriately
4. **Be mindful of performance**: Each piece of advice adds overhead to function calls
5. **Use metadata for communication**: Use the context's metadata field to share data between advice

## Fluent API Alternative

gosaidsno also provides a fluent API that offers a more convenient way to configure advice:

```go
package main

import (
    "fmt"
    "log"

    "github.com/seyallius/gosaidno/aspect"
)

func main() {
    // Use the fluent API to configure advice
    aspect.For("GreetUser").
        WithBefore(func(c *aspect.Context) error {
            log.Printf("About to execute %s", c.FunctionName)
            return nil
        }).
        WithAfter(func(c *aspect.Context) error {
            log.Printf("Finished executing %s", c.FunctionName)
            return nil
        })

    // Wrap your function using the builder
    builder := aspect.For("GreetUser")
    greetUser := aspect.Wrap1R[string](
        builder.GetRegistry(),
        builder.GetFuncKey(),
        func(name string) string {
            return fmt.Sprintf("Hello, %s!", name)
        })

    // Use your enhanced function
    result := greetUser("World")
    fmt.Println(result)
}
```

The fluent API provides a more readable and concise way to configure multiple advice types for a function.

## Real-World Application Structure

For a complete real-world example showing proper project structure and multiple cross-cutting concerns, see the [real-world example](../examples/07_real_world_example/README.md). This example demonstrates:

- Service layer separation with business logic
- Centralized AOP setup using the fluent API
- Multiple organization approaches (globals vs structs)
- Complete cross-cutting concerns (logging, timing, validation, caching, error handling)

## Next Steps

Now that you've seen the basics, explore the [Usage Guide](./usage.md) for more advanced features and patterns, or check out the [Examples](../examples/README.md) for real-world use cases.