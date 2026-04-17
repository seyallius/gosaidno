package wrap

import (
	"context"

	"github.com/seyallius/gosaidno/v2/aspect"
)

// --- no return

// Wrap2 wraps a function with two arguments and no return values.
func Wrap2[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B)) func(A, B) {
	return func(a A, b B) {
		executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			fn(a, b)
		}, a, b)
	}
}

// Wrap2Ctx wraps a function with context, 2 args, no returns.
func Wrap2Ctx[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B)) func(context.Context, A, B) {
	return func(ctx context.Context, a A, b B) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			fn(c.Context(), a, b)
		}, a, b)
	}
}

// Wrap2Slice wraps a function with two fixed arguments, variadic slice, and no return values.
func Wrap2Slice[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, []any)) func(A, B, []any) {
	return func(a A, b B, variadicArgs []any) {
		executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			fn(a, b, variadicArgs)
		}, a, b, variadicArgs)
	}
}

// Wrap2SliceCtx wraps a function with context, 2 fixed args, variadic slice, no returns.
func Wrap2SliceCtx[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, []any)) func(context.Context, A, B, []any) {
	return func(ctx context.Context, a A, b B, variadicArgs []any) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			fn(c.Context(), a, b, variadicArgs)
		}, a, b, variadicArgs)
	}
}

// --- return value

// Wrap2R wraps a function with two arguments and one return value.
func Wrap2R[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B) R) func(A, B) R {
	return func(a A, b B) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result = fn(a, b)
			c.SetResult(0, result)
		}, a, b)
		return resolveResult(c, result)
	}
}

// Wrap2RCtx wraps a function with context, 2 args, one return.
func Wrap2RCtx[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B) R) func(context.Context, A, B) R {
	return func(ctx context.Context, a A, b B) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result = fn(c.Context(), a, b)
			c.SetResult(0, result)
		}, a, b)
		return resolveResult(c, result)
	}
}

// Wrap2SliceR wraps a function with two fixed arguments, variadic slice, and one return value.
func Wrap2SliceR[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, []any) R) func(A, B, []any) R {
	return func(a A, b B, variadicArgs []any) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result = fn(a, b, variadicArgs)
			c.SetResult(0, result)
		}, a, b, variadicArgs)
		return resolveResult(c, result)
	}
}

// Wrap2SliceRCtx wraps a function with context, 2 fixed args, variadic slice, one return.
func Wrap2SliceRCtx[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, []any) R) func(context.Context, A, B, []any) R {
	return func(ctx context.Context, a A, b B, variadicArgs []any) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result = fn(c.Context(), a, b, variadicArgs)
			c.SetResult(0, result)
		}, a, b, variadicArgs)
		return resolveResult(c, result)
	}
}

// --- return error

// Wrap2E wraps a function with two arguments and returns error.
func Wrap2E[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B) error) func(A, B) error {
	return func(a A, b B) error {
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			err = fn(a, b)
			c.Error = err
		}, a, b)
		return resolveError(c, err)
	}
}

// Wrap2ECtx wraps a function with context, 2 args, returns error.
func Wrap2ECtx[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B) error) func(context.Context, A, B) error {
	return func(ctx context.Context, a A, b B) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			err = fn(c.Context(), a, b)
			c.Error = err
		}, a, b)
		return resolveError(c, err)
	}
}

// Wrap2SliceE wraps a function with two fixed arguments, variadic slice, and error return.
func Wrap2SliceE[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, []any) error) func(A, B, []any) error {
	return func(a A, b B, variadicArgs []any) error {
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			err = fn(a, b, variadicArgs)
			c.Error = err
		}, a, b, variadicArgs)
		return resolveError(c, err)
	}
}

// Wrap2SliceECtx wraps a function with context, 2 fixed args, variadic slice, error return.
func Wrap2SliceECtx[A, B any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, []any) error) func(context.Context, A, B, []any) error {
	return func(ctx context.Context, a A, b B, variadicArgs []any) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			err = fn(c.Context(), a, b, variadicArgs)
			c.Error = err
		}, a, b, variadicArgs)
		return resolveError(c, err)
	}
}

// --- return value & error

// Wrap2RE wraps a function with two arguments and returns (result, error).
func Wrap2RE[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B) (R, error)) func(A, B) (R, error) {
	return func(a A, b B) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result, err = fn(a, b)
			c.SetResult(0, result)
			c.Error = err
		}, a, b)
		return resolveResultError(c, result, err)
	}
}

// Wrap2RECtx wraps a function with context, 2 args, returns (result, error).
func Wrap2RECtx[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B) (R, error)) func(context.Context, A, B) (R, error) {
	return func(ctx context.Context, a A, b B) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result, err = fn(c.Context(), a, b)
			c.SetResult(0, result)
			c.Error = err
		}, a, b)
		return resolveResultError(c, result, err)
	}
}

// Wrap2SliceRE wraps a function with two fixed arguments, variadic slice, and (result, error) return.
func Wrap2SliceRE[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, []any) (R, error)) func(A, B, []any) (R, error) {
	return func(a A, b B, variadicArgs []any) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result, err = fn(a, b, variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, a, b, variadicArgs)
		return resolveResultError(c, result, err)
	}
}

// Wrap2SliceRECtx wraps a function with context, 2 fixed args, variadic slice, (result, error) return.
func Wrap2SliceRECtx[A, B, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, []any) (R, error)) func(context.Context, A, B, []any) (R, error) {
	return func(ctx context.Context, a A, b B, variadicArgs []any) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result, err = fn(c.Context(), a, b, variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, a, b, variadicArgs)
		return resolveResultError(c, result, err)
	}
}
