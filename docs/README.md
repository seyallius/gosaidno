# _gosaidno_

Welcome to the comprehensive documentation for _gosaidno_, an Aspect-Oriented Programming (AOP) library for Go.

## What is AOP?

Aspect-Oriented Programming (AOP) is a programming paradigm that aims to increase modularity by allowing the separation of cross-cutting concerns. It does so by adding additional behavior to existing code without modifying the code itself.

In simpler terms, AOP helps you handle concerns that cut across multiple parts of your application (like logging, security, caching) separately from your core business logic.

## Why _gosaidno_?

Go doesn't have built-in support for annotations or aspects like Java or other languages. However, many developers still need to handle cross-cutting concerns cleanly. _gosaidno_ provides a Go-idiomatic solution that:

- Requires no code generation or build tools
- Uses simple function wrapping instead of complex reflection
- Maintains type safety through generics
- Provides flexible advice execution order
- Integrates seamlessly with existing Go code

## Core Concepts

- **Function Registration**: Register functions that you want to enhance with cross-cutting concerns
- **Advice**: Code that implements cross-cutting concerns (logging, caching, etc.)
- **Advice Types**: Different types of advice that execute at different points in the function lifecycle
- **Wrapping**: The mechanism that connects your original function with its advice
- **Context**: The shared state that allows advice functions to communicate with each other
- **Fluent API**: Declarative, type-safe configuration using method chaining

## Fluent API

gosaidno now includes a fluent/declarative API that makes it easy to configure advice:

```go
// Configure advice using fluent API
aspect.For("GetUser").
    WithBefore(authCheck).
    WithAfter(logging).
    WithAround(caching)

// Then wrap your function
builder := aspect.For("GetUser")
getUserBusinessLogicFn := wrap.Wrap1RE[string,*User](builder.GetRegistry(), builder.GetFuncKey(), GetUserBusinessLogicFn)
user, err := getUserBusinessLogicFn("userId_1") // `GetUserBusinessLogicFn` will run with AOP support!
```

Ready to get started? Check out the [Quick Start Guide](./quick-start.md)!