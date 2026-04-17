// Package main - variadic_example demonstrates slice-based variadic wrappers in gosaidno
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/seyallius/gosaidno/v2/aspect"
	"github.com/seyallius/gosaidno/v2/aspect/wrap"
)

// -------------------------------------------- Domain Models --------------------------------------------

// LogEntry represents a structured log entry
type LogEntry struct {
	Level     string
	Message   string
	Timestamp time.Time
	Fields    map[string]any
}

// -------------------------------------------- Setup with Variadic Wrappers --------------------------------------------

func setupAOP() {
	log.Println("=== Setting up AOP with Variadic Support ===")

	// Example 1: Logging with variable key-value pairs
	aspect.For("StructuredLog").
		WithBefore(func(c *aspect.Context) error {
			level := c.Args[0].(string)
			message := c.Args[1].(string)
			log.Printf("🟢 [LOG-START] %s: %s", level, message)
			return nil
		}).
		WithAfter(func(c *aspect.Context) error {
			level := c.Args[0].(string)
			message := c.Args[1].(string)
			log.Printf("🔵 [LOG-END] %s: %s - Duration: %v",
				level, message, time.Since(c.Metadata["start_time"].(time.Time)))
			return nil
		}).
		WithBeforeP(func(c *aspect.Context) error {
			// Record start time for timing
			c.Metadata["start_time"] = time.Now()
			return nil
		}, 90)

	// Example 2: Math operations with variable operands
	aspect.For("Calculate").
		WithAround(func(c *aspect.Context) error {
			// Caching: check if we've computed this before
			key := fmt.Sprintf("%v", c.Args)
			if cached, ok := cache[key]; ok {
				log.Printf("💾 [CACHE HIT] %v", key)
				c.SetResult(0, cached)
				c.Skipped = true
				return nil
			}
			log.Printf("🔍 [CACHE MISS] %v", key)
			return nil
		}).
		WithAfterReturning(func(c *aspect.Context) error {
			// Populate cache on success
			if !c.Skipped {
				key := fmt.Sprintf("%v", c.Args)
				cache[key] = c.Results[0]
			}
			return nil
		})

	// Example 3: String builder with variadic parts
	aspect.For("BuildString").
		WithBefore(func(c *aspect.Context) error {
			parts := c.Args[1].([]any)
			log.Printf("📝 [BUILD] Concatenating %d parts", len(parts))
			return nil
		})

	log.Println("=== Variadic AOP Setup Complete ===\n")
}

var cache = make(map[string]any)

// -------------------------------------------- Business Logic (Variadic Functions) --------------------------------------------

// StructuredLog logs a message with variable key-value fields
func structuredLogImpl(level, message string, fields []any) error {
	entry := &LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Fields:    make(map[string]any),
	}

	// Parse key-value pairs from fields slice
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			entry.Fields[key] = fields[i+1]
		}
	}

	log.Printf("📋 [LOGGED] %+v", entry)
	return nil
}

// Calculate performs math operations with variable operands
func calculateImpl(operation string, base float64, operands []any) (float64, error) {
	result := base

	for _, op := range operands {
		switch val := op.(type) {
		case int:
			switch operation {
			case "add":
				result += float64(val)
			case "mul":
				result *= float64(val)
			}
		case float64:
			switch operation {
			case "add":
				result += val
			case "mul":
				result *= val
			}
		}
	}

	return result, nil
}

// BuildString concatenates parts with a separator
func buildStringImpl(separator string, parts []any) string {
	strs := make([]string, 0, len(parts))
	for _, p := range parts {
		strs = append(strs, fmt.Sprintf("%v", p))
	}
	return strings.Join(strs, separator)
}

// -------------------------------------------- Wrapped Functions --------------------------------------------

var (
	// Logging with variable fields
	StructuredLog = func(level, message string, fields ...any) error {
		builder := aspect.For("StructuredLog")
		return wrap.Wrap2SliceE[string, string](
			builder.GetRegistry(),
			builder.GetFuncKey(),
			structuredLogImpl,
		)(level, message, fields)
	}

	// Math with variable operands
	Calculate = func(operation string, base float64, operands ...any) (float64, error) {
		builder := aspect.For("Calculate")
		return wrap.Wrap2SliceRE[string, float64, float64](
			builder.GetRegistry(),
			builder.GetFuncKey(),
			calculateImpl,
		)(operation, base, operands)
	}

	// String building with variable parts
	BuildString = func(separator string, parts ...any) string {
		builder := aspect.For("BuildString")
		return wrap.Wrap1SliceR[string, string](
			builder.GetRegistry(),
			builder.GetFuncKey(),
			buildStringImpl,
		)(separator, parts)
	}
)

// -------------------------------------------- Examples --------------------------------------------

func example1_StructuredLogging() {
	fmt.Println("========== Example 1: Structured Logging with Variadic Fields ==========")

	// Log with different numbers of fields
	_ = StructuredLog("INFO", "User authenticated", "user_id", 123)
	_ = StructuredLog("WARN", "Rate limit approaching", "user_id", 456, "remaining", 10, "window", "1h")
	_ = StructuredLog("ERROR", "Database connection failed",
		"host", "db.example.com",
		"port", 5432,
		"error", "timeout",
		"retry_count", 3)
}

func example2_MathOperations() {
	fmt.Println("\n========== Example 2: Math Operations with Variable Operands ==========")

	// Addition with multiple operands
	result1, _ := Calculate("add", 10, 20, 30, 40)
	fmt.Printf("✅ Sum: 10 + 20 + 30 + 40 = %.0f\n", result1)

	// Multiplication with caching (second call should hit cache)
	result2, _ := Calculate("mul", 2, 3, 4, 5)
	fmt.Printf("✅ Product: 2 × 3 × 4 × 5 = %.0f\n", result2)

	// Same calculation again - should use cache
	result3, _ := Calculate("mul", 2, 3, 4, 5)
	fmt.Printf("✅ Cached: 2 × 3 × 4 × 5 = %.0f\n", result3)
}

func example3_StringBuilding() {
	fmt.Println("\n========== Example 3: String Building with Variadic Parts ==========")

	// Build path with variable segments
	path := BuildString("/", "api", "v1", "users", 123, "posts")
	fmt.Printf("✅ Path: %s\n", path)

	// Build CSV with variable values
	csv := BuildString(",", "John", "Doe", 30, "Engineer", true)
	fmt.Printf("✅ CSV: %s\n", csv)

	// Build URL query string
	query := BuildString("&", "page=1", "limit=20", "sort=name")
	fmt.Printf("✅ Query: %s\n", query)
}

func example4_ContextAwareVariadic() {
	fmt.Println("\n========== Example 4: Context-Aware Variadic Function ==========")

	// Context-aware logging function
	logWithContext := func(ctx context.Context, level, message string, fields []any) error {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Add request ID from context to fields
		if reqID := ctx.Value("request_id"); reqID != nil {
			fields = append(fields, "request_id", reqID)
		}

		return structuredLogImpl(level, message, fields)
	}

	// Wrap with context-aware slice wrapper
	builder := aspect.For("ContextLog")
	wrappedLog := wrap.Wrap2SliceECtx[string, string](
		builder.GetRegistry(),
		builder.GetFuncKey(),
		logWithContext,
	)

	// Create context with values
	ctx := context.WithValue(context.Background(), "request_id", "req-abc123")

	// Use the wrapped function
	_ = wrappedLog(ctx, "INFO", "Processing request", []any{"action", "create_user"})
}

// -------------------------------------------- Main --------------------------------------------

func main() {
	setupAOP()

	example1_StructuredLogging()
	example2_MathOperations()
	example3_StringBuilding()
	example4_ContextAwareVariadic()

	fmt.Println("\n========== Variadic Examples Complete ==========")
	fmt.Println("\n💡 Key Takeaways:")
	fmt.Println("  • Slice wrappers accept []any as final parameter for flexibility")
	fmt.Println("  • Create helper functions with ...T for ergonomic variadic syntax")
	fmt.Println("  • Type assertions needed when extracting values from []any")
	fmt.Println("  • All advice types work identically with slice wrappers")
	fmt.Println("  • Context propagation works seamlessly with slice variants")
}
