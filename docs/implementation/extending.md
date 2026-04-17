# Extending gosaidno

gosaidno is designed to be extensible, allowing you to customize and extend its functionality to meet your specific needs. Understanding the extensibility points helps you adapt the library to your requirements.

## Extension Points

### 1. Custom Advice Types

While gosaidno provides 5 core advice types, you can conceptually extend the system by combining existing types:

```go
// Example: Conditional advice using existing types
func ConditionalAdvice(condition func(*aspect.Context) bool, advice aspect.Advice) aspect.Advice {
    return aspect.Advice{
        Type: advice.Type,
        Handler: func(c *aspect.Context) error {
            if condition(c) {
                return advice.Handler(c)
            }
            return nil
        },
        Priority: advice.Priority,
    }
}

// Usage
aspect.AddAdvice("MyFunc", ConditionalAdvice(
    func(c *aspect.Context) bool {
        return c.Args[0].(string) == "special"
    },
    aspect.Advice{
        Type: aspect.Before,
        Handler: specialLoggingHandler,
        Priority: 100,
    },
))
```

### 2. Custom Registry Implementations

The registry interface can be extended for specialized needs:

```go
type ExtendedRegistry struct {
    *aspect.Registry
    persistence PersistenceLayer  // Custom persistence layer
    metrics     MetricsCollector  // Custom metrics
}

func NewExtendedRegistry() *ExtendedRegistry {
    return &ExtendedRegistry{
        Registry: aspect.NewRegistry(),
        persistence: NewDatabasePersistence(),
        metrics: NewPrometheusCollector(),
    }
}

func (er *ExtendedRegistry) RegisterWithMetadata(name string, metadata map[string]interface{}) error {
    // Store metadata in persistence layer
    if err := er.persistence.StoreMetadata(name, metadata); err != nil {
        return err
    }
    
    // Register normally
    return er.Registry.Register(name)
}
```

### 3. Enhanced Context

While you can't directly modify the core Context, you can create utility functions:

```go
// Enhanced context utilities
type ContextHelper struct {
    c *aspect.Context
}

func NewContextHelper(c *aspect.Context) *ContextHelper {
    return &ContextHelper{c: c}
}

func (ch *ContextHelper) GetString(key string) (string, bool) {
    if val, exists := ch.c.Metadata[key]; exists {
        if str, ok := val.(string); ok {
            return str, true
        }
    }
    return "", false
}

func (ch *ContextHelper) SetTyped(key string, value interface{}) {
    ch.c.Metadata[key] = value
}

// Usage in advice
func myAdvice(c *aspect.Context) error {
    helper := NewContextHelper(c)
    userID, exists := helper.GetString("user_id")
    if !exists {
        return errors.New("user_id not found")
    }
    // Use userID safely
    return nil
}
```

### 4. Custom Wrapper Generators

For unsupported function signatures, create custom wrappers:

```go
// Custom wrapper for variadic functions
func WrapVariadic(name string, fn func(...interface{}) interface{}) func(...interface{}) interface{} {
    return func(args ...interface{}) interface{} {
        var result interface{}
        c := executeWithAdvice(name, func(c *aspect.Context) {
            result = fn(args...)
            c.SetResult(0, result)
        }, args...)
        
        // Handle result from Around advice if target was skipped
        if c.Skipped && len(c.Results) > 0 {
            result = c.Results[0]
        }
        
        return result
    }
}
```

### 5. Composite Advice Patterns

Combine multiple advice functions into reusable patterns:

```go
func StandardLoggingAndTiming() []aspect.Advice {
    return []aspect.Advice{
        {
            Type:     aspect.Before,
            Priority: 200,
            Handler:  timingStartHandler,
        },
        {
            Type:     aspect.After,
            Priority: 100,
            Handler:  timingEndAndLogHandler,
        },
    }
}

// Usage
for _, advice := range StandardLoggingAndTiming() {
    aspect.AddAdvice("MyFunc", advice)
}
```

## Advanced Extension Techniques

### 1. Decorator Pattern Integration

Integrate with Go's decorator pattern for more complex extensions:

```go
type Decorator func(aspect.AdviceFunc) aspect.AdviceFunc

// Retry decorator
func WithRetry(maxRetries int) Decorator {
    return func(next aspect.AdviceFunc) aspect.AdviceFunc {
        return func(c *aspect.Context) error {
            var lastErr error
            for i := 0; i <= maxRetries; i++ {
                if err := next(c); err != nil {
                    lastErr = err
                    if i < maxRetries {
                        time.Sleep(time.Duration(i+1) * time.Second)
                        continue
                    }
                } else {
                    return nil // Success
                }
            }
            return lastErr
        }
    }
}

// Apply decorator to advice
decoratedHandler := WithRetry(3)(originalHandler)
aspect.AddAdvice("MyFunc", aspect.Advice{
    Type:     aspect.Around,
    Handler:  decoratedHandler,
    Priority: 100,
})
```

### 2. Configuration-Driven Extensions

Build configuration systems that drive AOP behavior:

```go
type AOPRule struct {
    FunctionPattern string                 `json:"function_pattern"`
    Advices         []AdviceConfig        `json:"advices"`
    Conditions      map[string]interface{} `json:"conditions"`
}

type AdviceConfig struct {
    Type     string `json:"type"`
    Priority int    `json:"priority"`
    Config   map[string]interface{} `json:"config"`
}

func ApplyRules(rules []AOPRule) error {
    for _, rule := range rules {
        // Match functions based on pattern
        matchedFunctions := findFunctions(rule.FunctionPattern)
        
        for _, funcName := range matchedFunctions {
            for _, adviceCfg := range rule.Advices {
                advice, err := buildAdvice(adviceCfg)
                if err != nil {
                    return err
                }
                
                if err := aspect.AddAdvice(funcName, advice); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}
```

### 3. Plugin Architecture

Design plugins that can register their own AOP configurations:

```go
type AOPPlugin interface {
    Name() string
    Configure(*aspect.Registry) error
    Dependencies() []string
}

type LoggingPlugin struct{}

func (lp *LoggingPlugin) Name() string { return "logging" }

func (lp *LoggingPlugin) Configure(registry *aspect.Registry) error {
    // Register logging advice for all functions
    // This is conceptual - would need registry introspection
    return nil
}

func (lp *LoggingPlugin) Dependencies() []string { return []string{} }

// Plugin manager
type PluginManager struct {
    plugins map[string]AOPPlugin
}

func (pm *PluginManager) Register(plugin AOPPlugin) {
    pm.plugins[plugin.Name()] = plugin
}

func (pm *PluginManager) ApplyAll() error {
    // Topological sort based on dependencies
    sortedPlugins := topologicalSort(pm.plugins)
    
    for _, plugin := range sortedPlugins {
        if err := plugin.Configure(aspect.GlobalRegistry()); err != nil {
            return err
        }
    }
    return nil
}
```

## Considerations for Extensions

### Performance Impact
- Custom extensions may add overhead
- Measure performance after extending
- Consider caching for expensive operations

### Maintainability
- Keep extensions simple and focused
- Document custom extensions well
- Test extensions thoroughly

### Compatibility
- Extensions should work with future versions
- Follow the same patterns as core code
- Consider backward compatibility

## Community Extensions

Consider contributing useful extensions back to the community:
- Common advice patterns
- Integration with popular libraries
- Utility functions for common tasks

The extensibility of gosaidno allows you to adapt it to your specific requirements while maintaining the core benefits of the AOP approach.