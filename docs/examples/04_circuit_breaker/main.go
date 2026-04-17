// Package main - circuit_breaker demonstrates Around advice for fault tolerance
package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/seyallius/gosaidno/v2/aspect"
	"github.com/seyallius/gosaidno/v2/aspect/wrap"
	"github.com/seyallius/gosaidno/v2/docs/examples/utils"
)

// -------------------------------------------- Circuit Breaker --------------------------------------------

type CircuitBreaker struct {
	mu              sync.RWMutex
	failures        int
	maxFailures     int
	state           string // CLOSED, OPEN, HALF_OPEN
	openedAt        time.Time
	resetTimeout    time.Duration
	successesNeeded int
	successes       int
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:     maxFailures,
		state:           "CLOSED",
		resetTimeout:    resetTimeout,
		successesNeeded: 2,
	}
}

var externalServiceCircuit = NewCircuitBreaker(3, 5*time.Second)

// -------------------------------------------- Setup --------------------------------------------

var registry = aspect.NewRegistry()

func setupAOP() {
	log.Println("=== Setting up Circuit Breaker AOP ===")

	registry.MustRegister("CallExternalService")

	// Around advice: circuit breaker
	registry.MustAddAdvice("CallExternalService", aspect.Advice{
		Type:     aspect.Around,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			utils.LogAround(c, 100, "CIRCUIT BREAKER")

			// Check circuit state
			externalServiceCircuit.mu.Lock()
			state := externalServiceCircuit.state
			failures := externalServiceCircuit.failures
			maxFailures := externalServiceCircuit.maxFailures

			// Check if we should transition from OPEN to HALF_OPEN
			if state == "OPEN" && time.Since(externalServiceCircuit.openedAt) > externalServiceCircuit.resetTimeout {
				log.Printf("   🔄 [CIRCUIT] Transitioning OPEN -> HALF_OPEN (timeout elapsed)")
				externalServiceCircuit.state = "HALF_OPEN"
				externalServiceCircuit.successes = 0
				state = "HALF_OPEN"
			}
			externalServiceCircuit.mu.Unlock()

			log.Printf("   🔌 [CIRCUIT] Current state: %s, Failures: %d/%d", state, failures, maxFailures)

			if state == "OPEN" {
				externalServiceCircuit.mu.RLock()
				remaining := externalServiceCircuit.resetTimeout - time.Since(externalServiceCircuit.openedAt)
				externalServiceCircuit.mu.RUnlock()

				log.Printf("   🚫 [CIRCUIT] Rejecting call - circuit OPEN")
				log.Printf("   ⏳ [CIRCUIT] Retry available in: %v", remaining.Round(time.Second))
				c.SetResult(0, "")
				c.Error = fmt.Errorf("circuit breaker OPEN, retry in %v", remaining.Round(time.Second))
				c.Skipped = true
				utils.LogAround(c, 100, "END (circuit open, execution blocked)")
				return nil
			}

			if state == "HALF_OPEN" {
				log.Printf("   ⚠️  [CIRCUIT] Circuit in HALF_OPEN state - testing service recovery")
			} else {
				log.Printf("   ✅ [CIRCUIT] Circuit CLOSED - allowing execution")
			}

			log.Printf("   ▶️  [CIRCUIT] Proceeding with function execution")
			utils.LogAround(c, 100, "END (allowing execution)")
			return nil // Allow execution
		},
	})

	// After advice: track failures
	registry.MustAddAdvice("CallExternalService", aspect.Advice{
		Type:     aspect.After,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			utils.LogAfter(c, 100, "CIRCUIT METRICS")

			if c.Skipped {
				log.Printf("   ⏭️  [CIRCUIT] Execution was skipped (circuit open)")
				return nil
			}

			if c.Error != nil {
				externalServiceCircuit.mu.Lock()
				externalServiceCircuit.failures++
				currentFailures := externalServiceCircuit.failures
				maxFailures := externalServiceCircuit.maxFailures
				log.Printf("   ❌ [CIRCUIT] Failure recorded (%d/%d)", currentFailures, maxFailures)

				if currentFailures >= maxFailures {
					externalServiceCircuit.state = "OPEN"
					externalServiceCircuit.openedAt = time.Now()
					log.Printf("   🚨 [CIRCUIT] Circuit OPENED - too many failures")
					log.Printf("   ⏰ [CIRCUIT] Reset timeout started: %v", externalServiceCircuit.resetTimeout)
				}
				externalServiceCircuit.mu.Unlock()
			} else {
				// Success
				externalServiceCircuit.mu.Lock()
				if externalServiceCircuit.state == "HALF_OPEN" {
					externalServiceCircuit.successes++
					currentSuccesses := externalServiceCircuit.successes
					neededSuccesses := externalServiceCircuit.successesNeeded
					log.Printf("   ✅ [CIRCUIT] Success in HALF_OPEN (%d/%d)", currentSuccesses, neededSuccesses)

					if currentSuccesses >= neededSuccesses {
						externalServiceCircuit.state = "CLOSED"
						externalServiceCircuit.failures = 0
						log.Printf("   🔄 [CIRCUIT] Circuit CLOSED - service recovered")
						log.Printf("   🎉 [CIRCUIT] Service is healthy again!")
					}
				} else {
					externalServiceCircuit.failures = 0
					log.Printf("   🔧 [CIRCUIT] Reset failure count - service is healthy")
				}
				externalServiceCircuit.mu.Unlock()
			}
			return nil
		},
	})

	// Before advice for logging
	registry.MustAddAdvice("CallExternalService", aspect.Advice{
		Type:     aspect.Before,
		Priority: 90,
		Handler: func(c *aspect.Context) error {
			utils.LogBefore(c, 90, "REQUEST LOG")
			endpoint := c.Args[0].(string)
			log.Printf("   🌐 [REQUEST] Calling external service: %s", endpoint)
			return nil
		},
	})

	log.Println("=== AOP Setup Complete ===\n")
}

// -------------------------------------------- Business Logic --------------------------------------------

var simulateFailure = true
var callCount = 0

func callExternalServiceImpl(endpoint string) (string, error) {
	callCount++
	log.Printf("   🌐 [BUSINESS] callExternalServiceImpl executing - call #%d", callCount)
	log.Printf("   💻 [BUSINESS] Making API call to: %s", endpoint)

	// Simulate flaky service
	time.Sleep(100 * time.Millisecond)

	if simulateFailure && callCount <= 3 {
		log.Printf("   ❌ [BUSINESS] Service unavailable (simulated failure)")
		return "", errors.New("service unavailable (500)")
	}

	response := fmt.Sprintf("Response from %s", endpoint)
	log.Printf("   ✅ [BUSINESS] API call successful: %s", response)
	return response, nil
}

// -------------------------------------------- Wrapped Functions --------------------------------------------

var CallExternalService = wrap.Wrap1RE(registry, "CallExternalService", callExternalServiceImpl)

// -------------------------------------------- Examples --------------------------------------------

func example1_CircuitOpenDueToFailures() {
	fmt.Println("\n========== Example 1: Circuit Opens Due to Failures ==========\n")

	simulateFailure = true
	callCount = 0

	for i := 1; i <= 5; i++ {
		fmt.Printf("\n--- Call %d ---\n", i)
		result, err := CallExternalService("/api/data")
		if err != nil {
			fmt.Printf("❌ Call failed: %v\n", err)
		} else {
			fmt.Printf("✅ Call succeeded: %s\n", result)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func example2_CircuitRecovery() {
	fmt.Println("\n========== Example 2: Circuit Recovery ==========\n")

	fmt.Println("--- Waiting for circuit to enter HALF_OPEN (5s) ---")
	time.Sleep(6 * time.Second)

	// Service is now healthy
	simulateFailure = false
	callCount = 0

	for i := 1; i <= 4; i++ {
		fmt.Printf("\n--- Recovery Call %d ---\n", i)
		result, err := CallExternalService("/api/data")
		if err != nil {
			fmt.Printf("❌ Call failed: %v\n", err)
		} else {
			fmt.Printf("✅ Call succeeded: %s\n", result)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func example3_RealWorldScenario() {
	fmt.Println("\n========== Example 3: Real World Scenario ==========\n")

	// Reset circuit
	externalServiceCircuit.mu.Lock()
	externalServiceCircuit.state = "CLOSED"
	externalServiceCircuit.failures = 0
	externalServiceCircuit.mu.Unlock()

	simulateFailure = false
	callCount = 0

	// Normal operation
	fmt.Println("--- Normal operation ---")
	for i := 1; i <= 2; i++ {
		result, _ := CallExternalService("/api/users")
		fmt.Printf("✅ %s\n", result)
		time.Sleep(100 * time.Millisecond)
	}

	// Service becomes unhealthy
	fmt.Println("\n--- Service degradation ---")
	simulateFailure = true
	for i := 1; i <= 4; i++ {
		_, err := CallExternalService("/api/users")
		if err != nil {
			fmt.Printf("❌ %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Service recovers
	fmt.Println("\n--- Waiting for recovery window (5s) ---")
	time.Sleep(6 * time.Second)

	simulateFailure = false
	callCount = 0
	fmt.Println("\n--- Service recovered ---")
	for i := 1; i <= 3; i++ {
		result, err := CallExternalService("/api/users")
		if err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Printf("✅ %s\n", result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// -------------------------------------------- Main --------------------------------------------

func main() {
	setupAOP()

	example1_CircuitOpenDueToFailures()
	example2_CircuitRecovery()
	example3_RealWorldScenario()

	fmt.Println("\n========== Circuit Breaker Examples Complete ==========")
}
