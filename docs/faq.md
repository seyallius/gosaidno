# Frequently Asked Questions (FAQ)

Find answers to common questions about gosaidno.

## General

### Q: What is gosaidno?

**A:** gosaidno is an Aspect-Oriented Programming (AOP) library for Go that allows you to modularize cross-cutting concerns like logging, authentication, caching, and error handling without cluttering your business logic. It provides a Go-idiomatic way to achieve separation of concerns without using reflection or code generation.

### Q: Why is it called gosaidno?

**A:** The name reflects the frustration of wanting to use annotations (like in Java) but Go saying "no". It's a playful reference to the lack of built-in AOP support in Go, while providing a solution anyway.

### Q: Is gosaidno production-ready?

**A:** Yes, gosaidno is designed for production use. It's thread-safe, has proper error handling, and includes comprehensive tests. However, as with any library, you should evaluate it in your specific context and requirements.

### Q: How does gosaidno compare to other AOP solutions?

**A:** Unlike other AOP libraries that rely on reflection or code generation, gosaidno uses simple function wrapping with generics. This makes it faster, more type-safe, and easier to debug than reflection-based approaches.

## Usage

### Q: Do I need to modify my existing functions to use gosaidno?

**A:** No, you don't need to modify your existing functions. You simply register them with gosaidno and wrap them using the appropriate wrapper function. Your original business logic remains unchanged.

### Q: Can I use gosaidno with methods on structs?

**A:** Yes, you can wrap methods by converting them to function values:

```go
type UserService struct{}

func (s *UserService) GetUser(id int) (User, error) {
    // implementation
}

service := &UserService{}
getUserFunc := func(id int) (User, error) {
    return service.GetUser(id)
}

wrappedGetUser := wrap.Wrap1RE[int, User]("UserService.GetUser", getUserFunc)
```

### Q: What happens if I forget to register a function?

**A:** If you call a wrapped function without registering it first, the function will still execute, but no advice will be applied. It will behave as if no AOP was configured.

### Q: Can I add advice to a function after it's been wrapped?

**A:** Yes, advice is looked up dynamically when the wrapped function is called, so you can add advice at any time before the function is invoked.

### Q: How do I handle errors in advice functions?

**A:** If an advice function returns an error, the execution flow depends on the advice type:
- **Before advice**: The error will prevent the target function from executing
- **Around advice**: The error will prevent the target function from executing
- **After advice**: Errors are typically logged but don't affect the target function
- **AfterReturning/AfterThrowing**: Errors are typically logged but don't affect the target function

### Q: Can advice functions modify function arguments or return values?

**A:** Yes, advice functions can modify the context, which allows them to modify arguments (through c.Args) and return values (through c.SetResult()). Around advice can also skip the target function entirely by setting c.Skipped = true.

### Q: How do I handle variadic functions?

**A:** Use the slice-based variadic wrappers (`Wrap*Slice*`) which accept a `[]any` as the final parameter. For ergonomic usage, create helper functions that accept variadic arguments and convert them to a slice:

```go
// Wrap with slice variant
wrapped := Wrap1SliceRE[int, int](registry, "Sum", sumFunc)

// Create helper for cleaner syntax
func Sum(base int, numbers ...int) (int, error) {
    args := make([]any, len(numbers))
    for i, n := range numbers {
        args[i] = n
    }
    return wrapped(base, args)
}
```

### Q: What's the performance difference between slice wrappers and fixed wrappers?

**A:** Slice wrappers have minimal overhead:
- **Memory**: ~24 bytes for empty slice allocation
- **CPU**: ~50-100ns additional processing time
- **Type assertions**: Required for extracting typed values from `[]any`

For performance-critical code with known argument counts, use fixed-arity wrappers. For flexible APIs or plugin systems, slice wrappers provide excellent ergonomics with acceptable overhead.

### Q: Can I mix fixed and variadic arguments?

**A:** Yes! That's exactly what the slice wrappers are designed for. The fixed arguments maintain full type safety, while the final `[]any` parameter provides flexibility:

```go
// Function signature: func(userID string, amount float64, metadata []any) error
wrapped := Wrap2SliceE[string, float64](registry, "ProcessPayment", processFunc)

// Call with fixed args + variable metadata
wrapped("user123", 99.99, []any{"currency", "USD", "source", "web"})
```

### Q: How do I access slice arguments in advice functions?

**A:** The slice is passed as a single argument in `c.Args`. Access it with type assertion:

```go
Handler: func(c *aspect.Context) error {
    // For Wrap1SliceRE: c.Args[0] is the fixed arg, c.Args[1] is the []any slice
    fixedArg := c.Args[0].(string)
    variadicArgs := c.Args[1].([]any)
    
    // Process variadic args
    for _, arg := range variadicArgs {
        // Handle each argument
    }
    return nil
}
```

### Q: Do slice wrappers work with the fluent API?

**A:** Yes! All slice wrapper variants are available through the fluent API:

```go
aspect.For("MyFunction").
    WithBefore(myAdvice).
    WithAfter(myAfterAdvice)

builder := aspect.For("MyFunction")
wrapped := builder.Wrap1SliceRE[string, int](myFunction)
```

## Technical

### Q: Does gosaidno use reflection?

**A:** No, gosaidno does not use reflection. It relies on Go generics and function wrapping to provide type-safe AOP capabilities.

### Q: What is the performance impact of using gosaidno?

**A:** The performance impact depends on how much advice you have attached to each function. Each piece of advice adds a small overhead (typically microseconds). For most applications, this overhead is negligible compared to the benefits of cleaner code organization.

### Q: Is gosaidno thread-safe?

**A:** Yes, gosaidno is designed to be thread-safe. The registry uses appropriate synchronization mechanisms, and context objects are not shared between goroutines.

### Q: How does the priority system work?

**A:** Within each advice type, advice is executed in descending order of priority (higher numbers execute first). For example, if you have three Before advice functions with priorities 100, 50, and 200, they will execute in the order: 200, 100, 50.

### Q: Can I remove advice after adding it?

**A:** Currently, gosaidno doesn't provide a direct way to remove advice. However, you can clear the entire registry using `aspect.Clear()` to reset all registrations and advice.

### Q: What happens if a target function panics?

**A:** If a target function panics, the AfterThrowing advice will execute, followed by the After advice. The panic will then be re-thrown, maintaining the original panic behavior.

### Q: How does the metadata system work?

**A:** The context's Metadata field is a map[string]any that allows advice functions to communicate with each other. Data stored by one advice function can be accessed by others in the same execution chain.

## Best Practices

### Q: How should I organize my AOP setup in a large application?

**A:** It's recommended to centralize your AOP setup in a single location, typically during application initialization. Create a dedicated package or function for setting up all your advice, and use consistent naming conventions for registered functions.

### Q: Should I use MustRegister and MustAddAdvice?

**A:** Use `MustRegister` and `MustAddAdvice` when you're certain the operations should succeed (e.g., during application startup with hardcoded function names). Use the regular versions when you need to handle errors gracefully.

### Q: How do I test functions that use gosaidno?

**A:** You can test your business logic independently of the advice, and test your advice functions separately. For integration tests, you can set up the AOP configuration in your test setup and verify that the expected advice is executed.

## Limitations

### Q: Are there any limitations on function signatures?

**A:** gosaidno provides wrapper functions for functions with up to 3 arguments. If you need to wrap functions with more arguments, you can create custom wrappers or refactor your functions to accept a single struct parameter.

### Q: Can I use gosaidno with third-party packages?

**A:** Yes, you can wrap functions from third-party packages as long as you can reference them. Simply register and wrap the functions as you would with your own code.

### Q: Does gosaidno support async/await patterns?

**A:** Since Go doesn't have async/await, gosaidno works with regular Go functions. You can wrap functions that return channels or work with goroutines, but the advice will execute in the calling goroutine.

## Troubleshooting

### Q: My advice isn't executing, what could be wrong?

**A:** Common issues include:
- Forgetting to register the function before adding advice
- Using the wrong function name when adding advice
- Calling the original function instead of the wrapped version
- Adding advice after the function has already been called

### Q: I'm getting priority conflicts, how do I resolve them?

**A:** Use a consistent priority scheme across your application. For example, use priorities in ranges like 100-199 for authentication, 200-299 for logging, etc. Document your priority scheme to avoid conflicts.

### Q: How do I debug issues with advice execution?

**A:** Add logging to your advice functions to trace execution flow. You can also inspect the context object to see what data is available at each step.

## Fluent API

### Q: What is the fluent API?

**A:** The fluent API provides a declarative, method-chaining interface for configuring advice. Instead of manually registering functions and adding advice separately, you can use method chaining:

```go
aspect.For("GetUser").
    WithBefore(authCheck).
    WithAfter(logging).
    WithAround(caching)
```

### Q: How does the fluent API work with function wrapping?

**A:** After configuring advice with the fluent API, you still need to wrap your functions using the registry from the builder:

```go
aspect.For("GetUser").
    WithBefore(authCheck).
    WithAfter(logging)

builder := aspect.For("GetUser")
wrappedFn := wrap.Wrap1RE[string,*User](builder.GetRegistry(), builder.GetFuncKey(), getUserImpl)
```

### Q: What are the benefits of the fluent API?

**A:** The fluent API offers several benefits:
- **Readability**: Method chaining makes configuration more readable
- **Convenience**: Groups related advice configuration together
- **Type Safety**: Still maintains full type safety without reflection
- **Consistency**: Provides a consistent way to configure multiple advice types

### Q: Can I use both the traditional API and the fluent API?

**A:** Yes, both APIs work with the same underlying registry system. You can use either approach or mix them as needed in your application.

### Q: Does the fluent API use reflection?

**A:** No, the fluent API does not use reflection. It uses the same underlying registry and advice chain system as the traditional API, maintaining the same performance characteristics and type safety.

### Q: How do I use the fluent API with custom registries?

**A:** You can use the fluent API with custom registries using `ForWithRegistry`:

```go
registry := aspect.NewRegistry()
aspect.ForWithRegistry(registry, "MyFunction").
    WithBefore(myAdvice)
```

## Development

### Q: How can I contribute to gosaidno?

**A:** Contributions are welcome! Check the CONTRIBUTING.md file for guidelines on submitting issues, pull requests, and code style. The project follows standard Go practices and welcomes improvements to documentation, examples, and functionality.

### Q: Where can I find more examples?

**A:** Look in the examples/ directory for various use cases including logging, caching, authentication, circuit breakers, and retry patterns.

## License

### Q: What license is gosaidno released under?

**A:** gosaidno is released under the MIT License. See the LICENSE file for complete licensing information.