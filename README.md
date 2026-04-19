# gosaidno

> _I really wanted to use annotations…_
> _but Go said no._

<p align="center">
  <img src="./gosaidno1.png" alt="goxide logo" width="400" style="display:inline-block; margin-right:10px;"/>
  <img src="./gosaidno2.png" alt="goxide logo 2" width="400" style="display:inline-block;"/>
</p>

**AOP without annotations. Just function wrapping. The Go way.**

## What is gosaidno?

`gosaidno` is an Aspect-Oriented Programming (AOP) library for Go that allows you to modularize cross-cutting concerns
like logging, authentication, caching, and error handling without cluttering your business logic.

Instead of copy-pasting boilerplate code throughout your application, you can register functions and attach advice (
cross-cutting concerns) to them. The library provides a clean, Go-idiomatic way to achieve separation of concerns.

## Key Features

- **No magic**: No reflection, no code generation, no build tags
- **Simple API**: Register functions and attach advice with minimal setup
- **Flexible advice types**: Before, After, Around, AfterReturning, AfterThrowing
- **Priority-based execution**: Control the order of advice execution
- **Generic function wrappers**: Type-safe wrappers for functions with various signatures
- **Variadic Support**: Slice-based wrappers (`Wrap*Slice*`) for functions with dynamic argument counts, while maintaining type safety for fixed parameters
- **Metadata system**: Share data between different advice functions

## Quick Example

### Traditional API
```go
// Register your function
aspect.MustRegister("ProcessPayment")

// Add logging advice
aspect.MustAddAdvice("ProcessPayment", aspect.Advice{
    Type:     aspect.Before,
    Priority: 100,
    Handler: func (c *aspect.Context) error {
        log.Printf("Starting %s", c.FunctionName)
        return nil
    },
})

// Wrap your function
ProcessPayment := wrap.Wrap1E[int](
	"ProcessPayment",
	func (paymentID int) error {
        // Your business logic here
        return nil
    },
)

// Use as normal
ProcessPayment(123)
```

### Fluent API (Recommended Usage!)
```go
// Configure advice using fluent API
aspect.For("ProcessPayment").
    WithBefore(func(c *aspect.Context) error {
        log.Printf("Starting %s", c.FunctionName)
        return nil
    }).
    WithAfter(func(c *aspect.Context) error {
        log.Printf("Completed %s", c.FunctionName)
        return nil
    })

// Wrap your function using the builder
builder := aspect.For("ProcessPayment")
ProcessPayment := wrap.Wrap1E[int](
    builder.GetRegistry(),
    builder.GetFuncKey(),
    func (paymentID int) error {
        // Your business logic here
        return nil
    },
)

// Use as normal
ProcessPayment(123)
```

## Getting Started

Check out our [Quick Start Guide](./docs/quick-start.md) to begin using gosaidno in your project, or dive deeper into
the [Usage Documentation](./docs/usage.md) for comprehensive examples and best practices.
> Note: If the documentation page is not loaded properly due to GitHub pages, you can read them at [here](./docs)
