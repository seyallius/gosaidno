// package wrap - wrap_slice_test validates slice variadic wrapper functionality
package wrap

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seyallius/gosaidno/v2/aspect"
)

// -------------------------------------------- Basic Slice Wrapper Tests --------------------------------------------

func TestWrap0Slice_BasicExecution(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSlice0")

	var capturedArgs []any
	var executed bool

	registry.MustAddAdvice("TestSlice0", aspect.Advice{
		Type:     aspect.Before,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			capturedArgs = make([]any, len(c.Args))
			copy(capturedArgs, c.Args)
			return nil
		},
	})

	targetFunc := func(args []any) {
		executed = true
	}

	wrapped := Wrap0Slice(registry, "TestSlice0", targetFunc)
	testArgs := []any{1, "hello", true}
	wrapped(testArgs)

	if !executed {
		t.Error("target function should have executed")
	}
	if len(capturedArgs) != 1 {
		t.Errorf("expected 1 arg in context (the slice), got %d", len(capturedArgs))
	}
	// The slice itself is passed as a single argument
	if slice, ok := capturedArgs[0].([]any); !ok || len(slice) != 3 {
		t.Error("expected slice argument to be preserved")
	}
}

func TestWrap1SliceRE_BasicExecution(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSlice1RE")

	var beforeCalled bool
	registry.MustAddAdvice("TestSlice1RE", aspect.Advice{
		Type:     aspect.Before,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			beforeCalled = true
			if len(c.Args) != 2 {
				t.Errorf("expected 2 args (fixed + slice), got %d", len(c.Args))
			}
			return nil
		},
	})

	targetFunc := func(base int, args []any) (int, error) {
		sum := base
		for _, a := range args {
			if val, ok := a.(int); ok {
				sum += val
			}
		}
		return sum, nil
	}

	wrapped := Wrap1SliceRE[int, int](registry, "TestSlice1RE", targetFunc)
	result, err := wrapped(10, []any{20, 30, 40})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != 100 {
		t.Errorf("expected 100, got %d", result)
	}
	if !beforeCalled {
		t.Error("Before advice should have executed")
	}
}

func TestWrap1SliceRE_ErrorPropagation(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSliceError")

	var capturedError error
	registry.MustAddAdvice("TestSliceError", aspect.Advice{
		Type:     aspect.After,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			capturedError = c.Error
			return nil
		},
	})

	targetFunc := func(prefix string, args []any) (string, error) {
		if len(args) == 0 {
			return "", errors.New("no arguments provided")
		}
		return prefix + args[0].(string), nil
	}

	wrapped := Wrap1SliceRE[string, string](registry, "TestSliceError", targetFunc)
	_, err := wrapped("test", []any{})

	if err == nil {
		t.Error("expected error from target function")
	}
	if capturedError == nil || capturedError.Error() != "no arguments provided" {
		t.Errorf("expected error to be captured in context")
	}
}

func TestWrap0SliceR_AroundAdviceSkip(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSliceSkip")

	registry.MustAddAdvice("TestSliceSkip", aspect.Advice{
		Type:     aspect.Around,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			// Skip target and return cached result
			c.SetResult(0, "cached")
			c.Skipped = true
			return nil
		},
	})

	var targetCalled bool
	targetFunc := func(args []any) string {
		targetCalled = true
		return "fresh"
	}

	wrapped := Wrap0SliceR[string](registry, "TestSliceSkip", targetFunc)
	result := wrapped([]any{1, 2, 3})

	if targetCalled {
		t.Error("target should not execute when Around advice skips")
	}
	if result != "cached" {
		t.Errorf("expected 'cached', got '%s'", result)
	}
}

// -------------------------------------------- Context-Aware Slice Wrapper Tests --------------------------------------------

func TestWrap1SliceRECtx_ContextPropagation(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSliceCtx")

	var capturedCtxValue string
	registry.MustAddAdvice("TestSliceCtx", aspect.Advice{
		Type:     aspect.Before,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			if val := c.Context().Value("request_id"); val != nil {
				capturedCtxValue = val.(string)
			}
			return nil
		},
	})

	targetFunc := func(ctx context.Context, userID string, args []any) (string, error) {
		// Verify context is propagated to target
		if ctx.Value("request_id") != "req-123" {
			return "", errors.New("context not propagated")
		}
		return "ok", nil
	}

	wrapped := Wrap1SliceRECtx[string, string](registry, "TestSliceCtx", targetFunc)
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	result, err := wrapped(ctx, "user1", []any{"arg1", "arg2"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected 'ok', got '%s'", result)
	}
	if capturedCtxValue != "req-123" {
		t.Errorf("expected context value in advice, got '%s'", capturedCtxValue)
	}
}

func TestWrap0SliceECtx_ContextCancellation(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSliceCancel")

	registry.MustAddAdvice("TestSliceCancel", aspect.Advice{
		Type:     aspect.Before,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			// Check for cancellation
			select {
			case <-c.Context().Done():
				return c.Context().Err()
			default:
				return nil
			}
		},
	})

	targetFunc := func(ctx context.Context, args []any) error {
		// Simulate work
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	wrapped := Wrap0SliceECtx(registry, "TestSliceCancel", targetFunc)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := wrapped(ctx, []any{1, 2, 3})
	if err == nil {
		t.Error("expected context cancellation error")
	}
}

// -------------------------------------------- Panic Recovery Tests --------------------------------------------

func TestWrap1SliceE_PanicRecovery(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSlicePanic")

	var panicCaught bool
	var panicValue any
	registry.MustAddAdvice("TestSlicePanic", aspect.Advice{
		Type:     aspect.AfterThrowing,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			panicCaught = true
			panicValue = c.PanicValue
			return nil
		},
	})
	t.Logf("use me or remove me: %v", panicValue)

	targetFunc := func(prefix string, args []any) error {
		if len(args) > 0 {
			panic("intentional panic")
		}
		return nil
	}

	wrapped := Wrap1SliceE[string](registry, "TestSlicePanic", targetFunc)

	// Catch the re-panic (our execution engine converts panic to error now)
	defer func() {
		// With the updated execution engine, panics are converted to errors
		// so we shouldn't see a re-panic here
	}()

	err := wrapped("test", []any{1})
	// The error should contain the panic information
	if err == nil {
		t.Error("expected error from panic recovery")
	}
	if !panicCaught {
		t.Error("AfterThrowing advice should have executed")
	}
}

// -------------------------------------------- Metadata Passing Tests --------------------------------------------

func TestWrap2SliceRE_MetadataSharing(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("TestSliceMetadata")

	// First advice sets metadata
	registry.MustAddAdvice("TestSliceMetadata", aspect.Advice{
		Type:     aspect.Before,
		Priority: 100,
		Handler: func(c *aspect.Context) error {
			c.SetMetadataVal("processed_count", len(c.Args[1].([]any)))
			return nil
		},
	})

	// Second advice reads metadata
	var capturedCount int
	registry.MustAddAdvice("TestSliceMetadata", aspect.Advice{
		Type:     aspect.After,
		Priority: 90,
		Handler: func(c *aspect.Context) error {
			if val, ok := c.GetMetadataVal("processed_count"); ok {
				capturedCount = val.(int)
			}
			return nil
		},
	})

	targetFunc := func(userID string, tags []any, args []any) (bool, error) {
		return true, nil
	}

	wrapped := Wrap2SliceRE[string, []any, bool](registry, "TestSliceMetadata", targetFunc)
	_, _ = wrapped("user1", []any{"tag1", "tag2"}, []any{"arg1"})

	if capturedCount != 2 {
		t.Errorf("expected metadata count 2, got %d", capturedCount)
	}
}

// -------------------------------------------- Helper Function Tests --------------------------------------------

func TestSliceWrapper_HelperPattern(t *testing.T) {
	registry := aspect.NewRegistry()
	registry.MustRegister("SumHelper")

	// Original function with slice
	sumFunc := func(base int, numbers []any) (int, error) {
		sum := base
		for _, n := range numbers {
			if val, ok := n.(int); ok {
				sum += val
			}
		}
		return sum, nil
	}

	// Wrap it
	wrapped := Wrap1SliceRE[int, int](registry, "SumHelper", sumFunc)

	// Create a helper for cleaner usage (variadic syntax)
	Sum := func(base int, numbers ...int) (int, error) {
		args := make([]any, len(numbers))
		for i, n := range numbers {
			args[i] = n
		}
		return wrapped(base, args)
	}

	result, err := Sum(10, 20, 30, 40)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != 100 {
		t.Errorf("expected 100, got %d", result)
	}
}
