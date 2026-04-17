# Wrapper Functions

Wrapper functions serve as the bridge between your code and the AOP system. They're the entry point that transforms regular function calls into AOP-enabled calls with all the associated advice execution.

## Wrapper Structure

```go
func Wrap1RE[A, R any](registry *Registry, funcKey FuncKey, fn func(A) (R, error)) func(A) (R, error) {
    return func(a A) (R, error) {
        var result R
        var err error

        c := executeWithAdvice(registry, funcKey, func(c *Context) {
            result, err = fn(a)           // Execute target function
            c.SetResult(0, result)      // Store result in context
            c.Error = err               // Store error in context
        }, a)

        // If Around advice set a result and skipped execution, use that result
        if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
            if res, ok := c.Results[0].(R); ok {
                result = res
            }
        }
        if c != nil && c.Error != nil {
            err = c.Error
        }

        return result, err
    }
}
```

## Available Wrapper Functions

gosaidno provides wrapper functions for various function signatures:

### No Arguments
- `Wrap0(registry *Registry, funcKey FuncKey, fn func()) func()` - No args, no returns
- `Wrap0R[R any](registry *Registry, funcKey FuncKey, fn func() R) func() R` - No args, one return
- `Wrap0RE[R any](registry *Registry, funcKey FuncKey, fn func() (R, error)) func() (R, error)` - No args, result + error
- `Wrap0E(registry *Registry, funcKey FuncKey, fn func() error) func() error` - No args, error only

### One Argument
- `Wrap1[A any](registry *Registry, funcKey FuncKey, fn func(A)) func(A)` - One arg, no returns
- `Wrap1R[A, R any](registry *Registry, funcKey FuncKey, fn func(A) R) func(A) R` - One arg, one return
- `Wrap1RE[A, R any](registry *Registry, funcKey FuncKey, fn func(A) (R, error)) func(A) (R, error)` - One arg, result + error
- `Wrap1E[A any](registry *Registry, funcKey FuncKey, fn func(A) error) func(A) error` - One arg, error only

### Two Arguments
- `Wrap2[A, B any](registry *Registry, funcKey FuncKey, fn func(A, B)) func(A, B)` - Two args, no returns
- `Wrap2R[A, B, R any](registry *Registry, funcKey FuncKey, fn func(A, B) R) func(A, B) R` - Two args, one return
- `Wrap2RE[A, B, R any](registry *Registry, funcKey FuncKey, fn func(A, B) (R, error)) func(A, B) (R, error)` - Two args, result + error
- `Wrap2E[A, B any](registry *Registry, funcKey FuncKey, fn func(A, B) error) func(A, B) error` - Two args, error only

### Three Arguments
- `Wrap3[A, B, C any](registry *Registry, funcKey FuncKey, fn func(A, B, C)) func(A, B, C)` - Three args, no returns
- `Wrap3R[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(A, B, C) R) func(A, B, C) R` - Three args, one return
- `Wrap3RE[A, B, C, R any](registry *Registry, funcKey FuncKey, fn func(A, B, C) (R, error)) func(A, B, C) (R, error)` - Three args, result + error
- `Wrap3E[A, B, C any](registry *Registry, funcKey FuncKey, fn func(A, B, C) error) func(A, B, C) error` - Three args, error only

## Integration with Fluent API

The wrapper functions work seamlessly with the fluent API. When using the fluent API, you retrieve the registry and function key from the builder:

```go
// Configure advice using fluent API
aspect.For("MyFunction").
    WithBefore(myBeforeAdvice).
    WithAfter(myAfterAdvice)

// Then wrap using the builder's registry and function key
builder := aspect.For("MyFunction")
wrappedFn := wrap.Wrap1RE[string, *User](
    builder.GetRegistry(),  // Get registry from builder
    builder.GetFuncKey(),   // Get function key from builder
    myOriginalFunction,
)
```

This pattern allows you to use the convenience of the fluent API for configuration while maintaining type safety in the wrapping process.

## Wrapper Execution Flow

Each wrapper follows the same pattern:

1. **Setup**: Initialize variables to capture results/errors
2. **Context Creation**: Create execution context with arguments
3. **Advice Execution**: Execute all applicable advice through the execution engine
4. **Target Execution**: Execute the original function (unless skipped by Around advice)
5. **Result Handling**: Extract results from context if execution was skipped
6. **Return**: Return results to caller

## Performance Considerations

- **Function Call Overhead**: Each wrapper adds one function call layer
- **Context Creation**: New context created on each call (memory allocation)
- **Advice Execution**: Each piece of advice adds execution time
- **Type Assertion**: When Around advice skips execution, type assertions occur

## Error Propagation

Wrapper functions properly propagate errors from both the target function and advice functions:
- Errors from target functions are preserved in the context
- Errors from Before/Around advice cause early termination
- Errors from After advice are logged but don't affect return values
- Panic recovery is handled by the execution engine

## Why This Two-Step Process?

The wrapper performs several critical functions:

### 1. Closure Creation
The wrapper creates a closure that captures the original function, preserving it for later execution while adding AOP capabilities.

### 2. Type Safety Maintenance
Through generics, the wrapper maintains type safety for the function signature while adding AOP functionality.

### 3. Context Management
The wrapper handles the creation and management of the execution context, ensuring proper data flow between the original function and its advice.

### 4. Result Handling
The wrapper manages the return of results from the execution engine back to the caller.

## Generic Implementation Benefits

### Type Safety
- Compile-time checking of function signatures
- No runtime type errors for correct usage
- Full IDE support with autocompletion

### Performance
- No reflection overhead
- Direct function calls without type inspection
- Optimized by the compiler for specific types

### Maintainability
- Clear function signatures
- Easy to understand type relationships
- Good error messages from the compiler

## Limitations of Generic Wrappers

### Limited Arity Support
Currently, only functions with up to 3 arguments are supported:
- `Wrap0`, `Wrap0R`, `Wrap0RE` - No arguments
- `Wrap1`, `Wrap1R`, `Wrap1RE`, `Wrap1E` - One argument
- `Wrap2`, `Wrap2R`, `Wrap2RE`, `Wrap2E` - Two arguments
- `Wrap3RE` - Three arguments

For functions with more arguments, you can:
- Create custom wrappers
- Refactor to use a single struct parameter
- Use variadic functions with manual handling

### Complex Signature Challenges
Functions with complex return types or multiple return values beyond (result, error) patterns require custom handling.

### Code Duplication
Each argument count requires its own set of wrapper functions, leading to some code duplication. This is a trade-off for type safety and performance.

## Wrapper Generation Pattern

Different wrapper types handle different function signatures:

### No Return Values
```go
func Wrap1[A any](name string, fn func(A)) func(A) {
    return func(a A) {
        executeWithAdvice(name, func(c *Context) {
            fn(a)
        }, a)
    }
}
```

### Single Return Value
```go
func Wrap1R[A, R any](name string, fn func(A) R) func(A) R {
    return func(a A) R {
        var result R
        c := executeWithAdvice(name, func(c *Context) {
            result = fn(a)
            c.SetResult(0, result)
        }, a)
        
        // Handle result from Around advice if target was skipped
        if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
            if res, ok := c.Results[0].(R); ok {
                result = res
            }
        }
        
        return result
    }
}
```

### Result and Error
```go
func Wrap1RE[A, R any](name string, fn func(A) (R, error)) func(A) (R, error) {
    return func(a A) (R, error) {
        var result R
        var err error
        executeWithAdvice(name, func(c *Context) {
            result, err = fn(a)
            c.SetResult(0, result)
            c.Error = err
        }, a)
        return result, err
    }
}
```

## Performance Considerations

- **Function Call Overhead**: One additional function call per wrapped function
- **Closure Creation**: Minimal overhead during wrapper creation
- **Type Assertion**: When Around advice provides results, type assertions occur
- **Memory Allocation**: Context creation per call

## Memory Management

Each wrapper creates a closure that captures the original function, but this is typically a one-time cost during application initialization. The actual function call involves creating a context and executing the advice chain.

The wrapper functions are the user-facing API of the AOP system, providing a clean, type-safe interface that hides the complexity of advice execution while maintaining performance and usability.