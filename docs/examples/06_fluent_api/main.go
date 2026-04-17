// Package main - fluent_api_example demonstrates the new fluent API for aspect-oriented programming
package main

import (
	"errors"
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

type Order struct {
	ID     string
	UserID string
	Amount float64
}

// -------------------------------------------- Setup with Fluent API --------------------------------------------

func setupAOPWithFluentAPI() {
	log.Println("=== Setting up AOP with Fluent API ===")

	// Example 1: Simple logging with fluent API
	aspect.For("GetUserFluent").
		WithBefore(func(c *aspect.Context) error {
			log.Printf("🟢 [FLUENT-BEFORE] Starting GetUserFluent with args: %v", c.Args)
			return nil
		}).
		WithAfter(func(c *aspect.Context) error {
			log.Printf("🔵 [FLUENT-AFTER] Completed GetUserFluent - Error: %v", c.Error)
			return nil
		})

	// Example 2: Timing with fluent API
	aspect.For("CreateOrderFluent").
		WithBefore(func(c *aspect.Context) error {
			c.Metadata["start_time"] = time.Now()
			log.Printf("⏱️  [FLUENT-TIMING] Started timer for CreateOrderFluent")
			return nil
		}).
		WithAfter(func(c *aspect.Context) error {
			startTime, ok := c.Metadata["start_time"].(time.Time)
			if !ok {
				return nil
			}
			duration := time.Since(startTime)
			log.Printf("⏱️  [FLUENT-TIMING] CreateOrderFluent took %v", duration)
			return nil
		})

	// Example 3: Validation with fluent API
	aspect.For("CreateOrderFluent").
		WithBefore(func(c *aspect.Context) error {
			userID := c.Args[0].(string)
			amount := c.Args[1].(float64)

			if userID == "" {
				return errors.New("userID cannot be empty")
			}
			if amount <= 0 {
				return errors.New("amount must be positive")
			}
			log.Printf("✅ [FLUENT-VALIDATION] Order validation passed")
			return nil
		})

	// Example 4: Caching with Around advice using fluent API
	cache := make(map[string]*User)
	aspect.For("GetUserCached").
		WithAround(func(c *aspect.Context) error {
			userID := c.Args[0].(string)

			// Check cache first
			if cachedUser, exists := cache[userID]; exists {
				log.Printf("💾 [FLUENT-CACHE] Cache HIT for user %s", userID)
				c.SetResult(0, cachedUser)
				c.Skipped = true // Skip target execution
				return nil
			}

			log.Printf("🔍 [FLUENT-CACHE] Cache MISS for user %s", userID)
			return nil // Let target execute
		}).
		WithAfterReturning(func(c *aspect.Context) error {
			// Populate cache after successful execution
			userID := c.Args[0].(string)
			user := c.Results[0].(*User)
			cache[userID] = user
			log.Printf("💾 [FLUENT-CACHE] Cached user %s", userID)
			return nil
		})

	log.Println("=== Fluent API Setup Complete ===")
	log.Println()
}

// -------------------------------------------- Business Logic (Unwrapped) --------------------------------------------

func getUserImpl(id string) (*User, error) {
	log.Printf("👨‍💼 [BUSINESS] getUserImpl executing with id: %s", id)
	// Simulate database query
	time.Sleep(50 * time.Millisecond)

	if id == "" {
		return nil, errors.New("user ID is required")
	}

	log.Printf("✅ [BUSINESS] getUserImpl completed successfully")
	return &User{
		ID:       id,
		Username: "john_doe",
		Email:    "john@example.com",
	}, nil
}

func createOrderImpl(userID string, amount float64) (*Order, error) {
	log.Printf("🛒 [BUSINESS] createOrderImpl executing for user: %s, amount: %.2f", userID, amount)
	// Simulate order creation
	time.Sleep(100 * time.Millisecond)

	order := &Order{
		ID:     fmt.Sprintf("order_%d", time.Now().Unix()),
		UserID: userID,
		Amount: amount,
	}

	log.Printf("✅ [BUSINESS] createOrderImpl completed, order: %s", order.ID)
	return order, nil
}

// -------------------------------------------- Wrapped Functions using Fluent API --------------------------------------------

var (
	// Using the fluent API to wrap functions
	GetUserFluent = func(id string) (*User, error) {
		builder := aspect.For("GetUserFluent")
		return wrap.Wrap1RE[string, *User](builder.GetRegistry(), builder.GetFuncKey(), getUserImpl)(id)
	}

	CreateOrderFluent = func(userID string, amount float64) (*Order, error) {
		builder := aspect.For("CreateOrderFluent")
		return wrap.Wrap2RE[string, float64, *Order](builder.GetRegistry(), builder.GetFuncKey(), createOrderImpl)(userID, amount)
	}

	// Alternative approach: Create a helper function for cleaner usage
	GetUserCached = func(id string) (*User, error) {
		builder := aspect.For("GetUserCached")
		return wrap.Wrap1RE[string, *User](builder.GetRegistry(), builder.GetFuncKey(), getUserImpl)(id)
	}
)

// -------------------------------------------- Examples --------------------------------------------

func example1_BasicFluentUsage() {
	fmt.Println("\n========== Example 1: Basic Fluent API Usage ==========")

	log.Println("\n--- Calling GetUserFluent with valid ID ---")
	user, err := GetUserFluent("user_123")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("\n🎯 Result: Got user %s (%s)\n", user.Username, user.Email)
}

func example2_ValidationWithFluent() {
	fmt.Println("\n========== Example 2: Validation with Fluent API ==========")

	// This will fail validation
	log.Println("\n--- Attempting to create order with invalid data ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("\n❌ Order creation rejected by validation: %v\n", r)
			}
		}()
		_, _ = CreateOrderFluent("", -100)
	}()

	// This will succeed
	log.Println("\n--- Creating valid order ---")
	order, err := CreateOrderFluent("user_123", 99.99)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("\n✅ Order created: %s for $%.2f\n", order.ID, order.Amount)
}

func example3_CachingWithFluent() {
	fmt.Println("\n========== Example 3: Caching with Fluent API ==========")

	log.Println("\n--- First call to GetUserCached (cache miss) ---")
	user1, err := GetUserCached("cached_user_123")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("🎯 First call result: Got user %s\n", user1.Username)

	log.Println("\n--- Second call to GetUserCached (cache hit) ---")
	user2, err := GetUserCached("cached_user_123")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("🎯 Second call result: Got user %s\n", user2.Username)
}

// Helper function to create a cleaner API for users
func wrapWithFluent[Fn any](funcName aspect.FuncKey, fn Fn) Fn {
	builder := aspect.For(funcName)
	// This is a simplified version - in practice, we'd need specific functions for each signature
	// For demonstration purposes, we'll use the registry directly
	registry := builder.GetRegistry()

	// Convert the generic function to specific types based on the actual function
	// This is a simplified example - a real implementation would need more sophisticated type handling
	switch v := any(fn).(type) {
	case func(string) (*User, error):
		wrapped := wrap.Wrap1RE[string, *User](registry, builder.GetFuncKey(), v)
		return any(wrapped).(Fn)
	case func(string, float64) (*Order, error):
		wrapped := wrap.Wrap2RE[string, float64, *Order](registry, builder.GetFuncKey(), v)
		return any(wrapped).(Fn)
	default:
		// For this example, we'll just return the original function
		// In a real implementation, we'd need to handle all possible function signatures
		return fn
	}
}

// -------------------------------------------- Main --------------------------------------------

func main() {
	// Setup AOP with the new fluent API
	setupAOPWithFluentAPI()

	// Run examples
	example1_BasicFluentUsage()
	example2_ValidationWithFluent()
	example3_CachingWithFluent()

	fmt.Println("\n========== All Fluent API Examples Complete ==========")

	// Demonstrate the proposed API from the requirements
	fmt.Println("\n========== Demonstrating the Proposed API Syntax ==========")

	// This is the API that was requested in the requirements:
	fmt.Println("// The fluent API allows this kind of usage:")
	fmt.Println("aspect.For(\"GetUser\").")
	fmt.Println("    WithBefore(authCheck).")
	fmt.Println("    WithAfter(logging).")
	fmt.Println("    WithAround(caching).")
	fmt.Println("    // Then wrap with: wrap.Wrap1RE[string,*User](builder.GetRegistry(), builder.GetFuncKey(), getUserImpl)")
}
