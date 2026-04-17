// package wrap - wrap provides function wrapping utilities with AOP advice execution
package wrap

import (
	"context"
	"fmt"

	"github.com/seyallius/gosaidno/v2/aspect"
)

// -------------------------------------------- Internal Helpers --------------------------------------------

// resolveResult handles the logic for extracting a generic result from the context,
// checking for advice skips, and performing safe type assertions.
func resolveResult[R any](c *aspect.Context, original R) R {
	// If Around advice skipped execution and set a result, try to use it
	if c != nil && c.Skipped && len(c.Results) > 0 && c.Results[0] != nil {
		if res, ok := c.Results[0].(R); ok {
			return res
		}
	}
	return original
}

// resolveError handles the logic for extracting an error from the context,
// allowing advice chains to replace the original error.
func resolveError(c *aspect.Context, original error) error {
	if c != nil && c.Error != nil {
		return c.Error
	}
	return original
}

// resolveResultError combines result and error resolution for functions returning (R, error).
func resolveResultError[R any](c *aspect.Context, origRes R, origErr error) (R, error) {
	finalRes := resolveResult(c, origRes)
	finalErr := resolveError(c, origErr)
	return finalRes, finalErr
}

// executeWithAdvice executes a function with full advice chain support and returns the context.
func executeWithAdvice(registry *aspect.Registry, functionName aspect.FuncKey, targetFn func(*aspect.Context), args ...any) *aspect.Context {
	return executeWithAdviceContext(registry, functionName, context.Background(), targetFn, args...)
}

// executeWithAdviceContext executes a function with full advice chain support using a specific context.Context.
func executeWithAdviceContext(registry *aspect.Registry, functionName aspect.FuncKey, ctx context.Context, targetFn func(*aspect.Context), args ...any) *aspect.Context {
	// Get advice chain from registry
	chain, err := registry.GetAdviceChain(functionName)
	if err != nil {
		// No advice registered, just execute target function
		c := aspect.NewContextWithContext(ctx, functionName, args...)
		targetFn(c)
		return c
	}

	// Create execution context
	c := aspect.NewContextWithContext(ctx, functionName, args...)

	if err = executeWithChain(chain, targetFn, c); err != nil {
		c.Error = err
	}

	return c
}

// 1. Update your execution function to return errors instead of panicking
func executeWithChain(chain *aspect.AdviceChain, targetFn func(*aspect.Context), c *aspect.Context) (finalErr error) {
	// Always execute After advice (even on panic/error)
	defer func() {
		if afterErr := chain.ExecuteAfter(c); afterErr != nil {
			if finalErr != nil {
				finalErr = fmt.Errorf("%w, after advice error: %v", finalErr, afterErr)
			} else {
				finalErr = afterErr
			}
		}
	}()
	// Handle Panic Recovery and Throwing advice - convert panic to error
	defer func() {
		if r := recover(); r != nil {
			c.PanicValue = r

			// Execute AfterThrowing advice for panic
			if throwErr := chain.ExecuteAfterThrowing(c); throwErr != nil {
				// Combine errors
				finalErr = fmt.Errorf("panic: %v, afterThrowing error: %w", r, throwErr)
			} else {
				finalErr = fmt.Errorf("panic recovered: %v", r)
			}
		}
	}()

	// Execute Before advice
	if err := chain.ExecuteBefore(c); err != nil {
		return fmt.Errorf("before advice failed: %w", err)
	}

	// Execute Around advice
	if chain.HasAround() {
		if err := chain.ExecuteAround(c); err != nil {
			return fmt.Errorf("around advice failed: %w", err)
		}
		// If Around advice sets Skipped, we skip the target function
		if c.Skipped {
			// Execute AfterReturning if no error
			if c.Error == nil {
				if err := chain.ExecuteAfterReturning(c); err != nil {
					return fmt.Errorf("afterReturning advice failed: %w", err)
				}
			}
			return nil
		}
	}

	// Execute Target Function (may panic, which is caught by defer)
	targetFn(c)

	// Execute AfterReturning advice (only if no error and no panic occurred)
	if c.Error == nil && !c.HasPanic() {
		if err := chain.ExecuteAfterReturning(c); err != nil {
			return fmt.Errorf("afterReturning advice failed: %w", err)
		}
	}

	// Return any error from the target function
	return c.Error
}
