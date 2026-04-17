# Architecture Overview

gosaidno's architecture consists of interconnected components that work together to provide AOP functionality:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Target        │───▶│   AdviceChain    │───▶│   Registry      │
│   Function      │    │   (per function) │    │   (global)      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        │
                        ┌─────────────┐                │
                        │   Context   │                │
                        │   (state)   │                │
                        └─────────────┘                │
                                                        │
┌─────────────────┐                                     │
│  Fluent API     │─────────────────────────────────────┘
│  (optional)     │
└─────────────────┘
```

## Component Roles

### 1. Registry
Acts as a lookup table mapping function names to their advice chains.
- Maintains global state of all registered functions
- Provides thread-safe access to advice chains
- Uses string keys for human-readable function identification
- Includes a default global registry for convenience
- Supports multiple registries for different contexts

### 2. AdviceChain
Manages the collection and execution order of advice for a single function.
- Stores different types of advice separately
- Handles priority-based execution within each advice type
- Orchestrates the execution flow
- Ensures proper execution order based on advice type semantics

### 3. Context
Serves as the communication channel between all components.
- Carries function arguments and results
- Holds error and panic information
- Provides metadata space for advice communication
- Signals execution state (skipped, etc.)
- Thread-safe for metadata access

### 4. Wrapper Functions
The entry point that orchestrates the entire AOP process.
- Create closures around target functions
- Handle type conversion and generics
- Call the execution engine with proper context
- Support various function signatures (0-3 arguments, with/without return values)

### 5. Fluent API (Optional Layer)
Provides a declarative, fluent interface for configuring advice.
- Offers method chaining for cleaner configuration
- Integrates with the Registry system
- Maintains type safety without runtime reflection
- Simplifies common configuration patterns

## How Components Interact

1. **Initialization Phase**:
   - Developer registers functions with the Registry (directly or via Fluent API)
   - Advice is added to specific AdviceChains
   - Target functions are wrapped with appropriate wrapper functions
   - Fluent API provides convenience methods for common patterns

2. **Execution Phase**:
   - Wrapped function is called
   - Wrapper creates Context with arguments
   - Registry is queried for the function's AdviceChain
   - AdviceChain executes advice in proper order
   - Target function executes (unless skipped by Around advice)
   - Results are returned through the wrapper

## Design Benefits

This architecture provides several advantages:

- **Separation of Concerns**: Each component has a clear responsibility
- **Scalability**: Registry can handle many functions efficiently
- **Flexibility**: Different advice types can be executed at different times
- **Communication**: Context enables data sharing between advice
- **Type Safety**: Generics maintain type information throughout
- **Usability**: Fluent API provides cleaner configuration options
- **Modularity**: Components can be extended or replaced independently

## Potential Bottlenecks

Understanding the architecture helps identify potential performance considerations:

- **Registry Access**: Global lookup during each function call
- **Sorting Overhead**: Priority sorting happens on each execution
- **Context Creation**: New context per function call
- **Memory Usage**: Each registered function maintains an advice chain
- **Fluent API Overhead**: Additional method calls during configuration (one-time cost)

This modular design allows each component to be optimized independently while maintaining the overall system integrity.