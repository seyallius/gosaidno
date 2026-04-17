// Package main - context_example demonstrates context propagation through AOP advice
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/seyallius/gosaidno/v2/aspect"
	"github.com/seyallius/gosaidno/v2/aspect/wrap"
)

// -------------------------------------------- Domain Models --------------------------------------------

type User struct {
	ID       string
	Username string
	Email    string
}

// -------------------------------------------- Setup with Context --------------------------------------------

func setupAOPWithContext() {
	log.Println("=== Setting up AOP with Context Support ===")

	// Example 1: Context-aware logging with request tracing
	aspect.For("GetUserWithContext").
		WithBefore(func(c *aspect.Context) error {
			// Access context values in advice
			requestID := c.Context().Value("request_id")
			log.Printf("🟢 [CONTEXT-LOG] Starting GetUserWithContext - Request ID: %v", requestID)
			return nil
		}).
		WithAfter(func(c *aspect.Context) error {
			requestID := c.Context().Value("request_id")
			log.Printf("🔵 [CONTEXT-LOG] Completed GetUserWithContext - Request ID: %v", requestID)
			return nil
		})

	// Example 2: Context deadline awareness
	aspect.For("SlowOperation").
		WithBefore(func(c *aspect.Context) error {
			log.Printf("🟠 [CONTEXT-DEADLINE] Checking deadline...")

			// Check if context has deadline
			if deadline, ok := c.Context().Deadline(); ok {
				log.Printf("⏰ [CONTEXT-DEADLINE] Deadline set: %v", deadline.Format("15:04:05.000"))

				// Check if already expired
				if time.Now().After(deadline) {
					log.Printf("❌ [CONTEXT-DEADLINE] Context already expired!")
					return context.DeadlineExceeded
				}
			} else {
				log.Printf("⏰ [CONTEXT-DEADLINE] No deadline set")
			}

			return nil
		})

	// Example 3: Context cancellation awareness
	aspect.For("CancelableOperation").
		WithBefore(func(c *aspect.Context) error {
			log.Printf("🟡 [CONTEXT-CANCEL] Checking for cancellation...")

			// Check if context is cancelled
			select {
			case <-c.Context().Done():
				log.Printf("🚫 [CONTEXT-CANCEL] Operation cancelled: %v", c.Context().Err())
				return c.Context().Err()
			default:
				log.Printf("✅ [CONTEXT-CANCEL] Context not cancelled, proceeding...")
			}

			return nil
		})

	log.Println("=== Context-Aware AOP Setup Complete ===")
	log.Println()
}

// -------------------------------------------- Business Logic (Context-Aware) --------------------------------------------

func getUserWithContextImpl(ctx context.Context, id string) (*User, error) {
	log.Printf("👨‍💼 [BUSINESS] getUserWithContextImpl executing with context and id: %s", id)

	// Simulate database query with context awareness
	done := make(chan struct{})
	go func() {
		defer close(done)
		// Simulate work
		time.Sleep(50 * time.Millisecond)
	}()

	select {
	case <-done:
		// Work completed
	case <-ctx.Done():
		// Context cancelled during work
		return nil, ctx.Err()
	}

	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	log.Printf("✅ [BUSINESS] getUserWithContextImpl completed successfully")
	return &User{
		ID:       id,
		Username: "john_doe",
		Email:    "john@example.com",
	}, nil
}

func slowOperationImpl(ctx context.Context, duration time.Duration) error {
	log.Printf("🐌 [BUSINESS] Slow operation starting for %v", duration)

	// Simulate slow operation with context awareness
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		log.Printf("✅ [BUSINESS] Slow operation completed after %v", duration)
		return nil
	case <-ctx.Done():
		log.Printf("❌ [BUSINESS] Slow operation cancelled: %v", ctx.Err())
		return ctx.Err()
	}
}

func cancelableOperationImpl(ctx context.Context) string {
	log.Printf("🔄 [BUSINESS] Cancelable operation starting...")

	// Simulate work that respects context cancellation
	select {
	case <-time.After(100 * time.Millisecond):
		log.Printf("✅ [BUSINESS] Cancelable operation completed")
		return "success"
	case <-ctx.Done():
		log.Printf("🚫 [BUSINESS] Cancelable operation interrupted: %v", ctx.Err())
		return "cancelled"
	}
}

// -------------------------------------------- Wrapped Functions (Context-Aware) --------------------------------------------

var (
	GetUserWithContext = func(ctx context.Context, id string) (*User, error) {
		builder := aspect.For("GetUserWithContext")
		return wrap.Wrap1RECtx[string, *User](builder.GetRegistry(), builder.GetFuncKey(), getUserWithContextImpl)(ctx, id)
	}

	SlowOperation = func(ctx context.Context, duration time.Duration) error {
		builder := aspect.For("SlowOperation")
		return wrap.Wrap1ECtx[time.Duration](builder.GetRegistry(), builder.GetFuncKey(), slowOperationImpl)(ctx, duration)
	}

	CancelableOperation = func(ctx context.Context) string {
		builder := aspect.For("CancelableOperation")
		return wrap.Wrap0RCtx[string](builder.GetRegistry(), builder.GetFuncKey(), cancelableOperationImpl)(ctx)
	}
)

// -------------------------------------------- Examples --------------------------------------------

func example1_ContextWithValue() {
	fmt.Println("\n========== Example 1: Context with Values ==========")

	// Create a context with values
	ctx := context.WithValue(context.Background(), "request_id", "req-12345")

	user, err := GetUserWithContext(ctx, "user_123")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("\n🎯 Result: Got user %s (%s) with request ID from context\n", user.Username, user.Email)
}

func example2_ContextWithDeadline() {
	fmt.Println("\n========== Example 2: Context with Deadline ==========")

	// Create a context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	log.Println("\n--- Slow operation that will exceed deadline ---")
	err := SlowOperation(ctx, 50*time.Millisecond) // Will take 50ms but deadline is 25ms
	if err != nil {
		fmt.Printf("\n❌ Operation failed due to deadline: %v\n", err)
	}

	log.Println("\n--- Slow operation that finishes before deadline ---")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	err2 := SlowOperation(ctx2, 50*time.Millisecond) // Will take 50ms with 100ms deadline
	if err2 != nil {
		fmt.Printf("\n❌ Unexpected error: %v\n", err2)
	} else {
		fmt.Printf("\n✅ Operation completed successfully\n")
	}
}

func example3_ContextCancellation() {
	fmt.Println("\n========== Example 3: Context Cancellation ==========")

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	log.Println("\n--- Starting operation, then cancelling ---")

	// Start the operation in a goroutine
	done := make(chan string, 1)
	go func() {
		result := CancelableOperation(ctx)
		done <- result
	}()

	// Wait a bit, then cancel
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for result
	result := <-done
	fmt.Printf("\n🔄 Operation result after cancellation: %s\n", result)

	log.Println("\n--- Operation with non-cancelled context ---")
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	result2 := CancelableOperation(ctx2)
	fmt.Printf("\n🔄 Operation result with active context: %s\n", result2)
}

// -------------------------------------------- Main --------------------------------------------

func main() {
	// Setup AOP with context support
	setupAOPWithContext()

	// Run examples
	example1_ContextWithValue()
	example2_ContextWithDeadline()
	example3_ContextCancellation()

	fmt.Println("\n========== Context Examples Complete ==========")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Context values propagate through all advice functions")
	fmt.Println("- Context deadlines are respected by advice and target functions")
	fmt.Println("- Context cancellation signals are detected and handled")
	fmt.Println("- AOP advice can access and react to context state")
}
