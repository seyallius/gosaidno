# Performance Considerations

Understanding the performance implications of using gosaidno helps you make informed decisions about when and how to use it effectively in your applications.

## Time Complexity Analysis

### Function Registration
- **Operation**: `aspect.Register("FunctionName")`
- **Time Complexity**: O(1) average case
- **Details**: Hash map lookup and insertion
- **Frequency**: Per application startup, not per function call

### Advice Addition
- **Operation**: `aspect.AddAdvice("FunctionName", advice)`
- **Time Complexity**: O(1) - simple slice append
- **Frequency**: Per application startup, not per function call

### Function Call with AOP
- **Operation**: Calling a wrapped function
- **Time Complexity**: O(P log P + A) where:
  - P = number of advice with the same priority level
  - A = total number of advice for the function
- **Details**: Includes priority sorting and advice execution

### Context Operations
- **Creation**: O(1) for basic fields, O(M) for M metadata entries
- **Result setting**: O(1) amortized with slice growth
- **Metadata access**: O(1) average case for map operations

## Memory Overhead

### Per Registered Function
- **Registry entry**: ~100 bytes for the advice chain structure
- **Storage**: Hash map entry overhead
- **Growth**: Monotonic - registry only grows, never shrinks

### Per Piece of Advice
- **Advice structure**: ~50 bytes per advice
- **Storage**: Slice storage overhead
- **Components**: Type, handler function pointer, priority

### Per Function Call
- **Context object**: ~200 bytes for basic fields
- **Arguments**: Size varies with function signature
- **Results**: Size varies with return types
- **Metadata**: Additional overhead if used

### Total Memory Impact
- **Registry**: O(F + A) where F = registered functions, A = total advice
- **Per call**: O(1) additional memory (context) regardless of advice count
- **Peak usage**: During function execution when context exists

## Performance Bottlenecks

### Priority Sorting
The most significant performance consideration is the priority sorting that happens during execution:

```go
sort.Slice(sortedAdviceList, func(i, j int) bool {
    return sortedAdviceList[i].Priority > sortedAdviceList[j].Priority
})
```

This occurs O(log P) times per advice type where P is the number of advice of that type.

### Registry Contention
In high-concurrency scenarios, registry access could become a bottleneck due to the RWMutex, though read operations predominate.

### Context Creation
Creating a new context object for each function call adds allocation overhead, though this is typically minor compared to the function's actual work.

## When Performance Matters

### High-Frequency Functions
For functions called thousands of times per second:
- Consider limiting the number of advice
- Profile to measure actual overhead
- Consider bypassing AOP for performance-critical paths

### Memory-Constrained Environments
Monitor registry growth:
- All registered functions remain in memory
- Large numbers of registered functions increase memory usage
- Consider using local registries for temporary scenarios

### Latency-Sensitive Applications
Profile advice execution time:
- Each piece of advice adds measurable overhead
- Priority sorting adds logarithmic factor
- Consider the cumulative effect of multiple advice

## Optimization Strategies

### Minimize Advice Count
- Only add essential advice to frequently called functions
- Consolidate similar advice into fewer, more comprehensive functions
- Use conditional logic within advice rather than multiple advice

### Efficient Advice Implementation
- Keep advice functions lightweight
- Avoid expensive operations in advice
- Cache expensive computations when possible

### Strategic Registration
- Register all functions during startup
- Avoid dynamic registration in hot paths
- Consider lazy registration for rarely used functions

## Benchmarking Considerations

When benchmarking code that uses gosaidno:

```go
func BenchmarkWithAOP(b *testing.B) {
    // Setup AOP once
    aspect.MustRegister("BenchmarkFunc")
    aspect.MustAddAdvice("BenchmarkFunc", loggingAdvice())
    
    wrappedFunc := wrap.Wrap0("BenchmarkFunc", func() { /* work */ })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wrappedFunc()
    }
}
```

Compare performance with and without AOP to quantify the overhead in your specific use case.

## Typical Overhead

In most applications, the overhead is negligible:
- **Simple functions**: 1-5 microseconds per call with moderate advice
- **Complex functions**: Overhead becomes insignificant compared to function work
- **Memory**: ~300 bytes per call temporarily, plus registry overhead

The performance impact is usually acceptable given the benefits of cleaner, more maintainable code.