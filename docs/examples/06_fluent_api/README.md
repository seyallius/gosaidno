# Fluent API Example

This example demonstrates the new fluent/declarative API for configuring aspect-oriented programming in the gosaidno library.

## Key Features

- **Fluent Configuration**: Use method chaining to configure advice for functions
- **Type Safe**: No reflection used in the API design
- **Flexible**: Support for all advice types (Before, After, Around, AfterReturning, AfterThrowing)
- **Priority Support**: Configure priority levels for advice execution order

## API Usage

The new fluent API allows you to configure advice like this:

```go
// Configure advice using fluent API
aspect.For("GetUser").
    WithBefore(authCheck).
    WithAfter(logging).
    WithAround(caching)

// Then wrap your function using the registry
builder := aspect.For("GetUser")
wrappedFn := wrap.Wrap1RE[string,*User](builder.GetRegistry(), builder.GetFuncKey(), getUserImpl)
```

## Benefits

1. **Cleaner Configuration**: More readable and maintainable advice configuration
2. **IDE Friendly**: Better autocomplete and type checking support
3. **Consistent**: Follows Go idioms and patterns
4. **Performant**: No runtime reflection overhead

## Example Output

The example demonstrates:
- Basic fluent API usage with logging
- Validation with Before advice
- Timing measurements
- Caching with Around advice