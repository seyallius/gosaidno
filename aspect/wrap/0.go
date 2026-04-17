package wrap

import (
	"context"

	"github.com/seyallius/gosaidno/v2/aspect"
)

// --- no return

// Wrap0 wraps a function with no arguments and no return values.
func Wrap0(registry *aspect.Registry, funcKey aspect.FuncKey, fn func()) func() {
	return func() {
		executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			fn()
		})
	}
}

// Wrap0Ctx wraps a function with context, no arguments, no returns.
func Wrap0Ctx(registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			fn(c.Context())
		})
	}
}

// Wrap0Slice wraps a function with no fixed arguments and a variadic slice parameter.
// Use this when your function accepts only variable arguments: func([]any)
func Wrap0Slice(registry *aspect.Registry, funcKey aspect.FuncKey, fn func([]any)) func([]any) {
	return func(variadicArgs []any) {
		executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			fn(variadicArgs)
		}, variadicArgs)
	}
}

// Wrap0SliceCtx wraps a function with context, no fixed arguments, and a variadic slice.
func Wrap0SliceCtx(registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, []any)) func(context.Context, []any) {
	return func(ctx context.Context, variadicArgs []any) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			fn(c.Context(), variadicArgs)
		}, variadicArgs)
	}
}

// --- return value

// Wrap0R wraps a function with no arguments and one return value.
func Wrap0R[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func() R) func() R {
	return func() R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result = fn()
			c.SetResult(0, result)
		})
		return resolveResult(c, result)
	}
}

// Wrap0RCtx wraps a function with context, no arguments, one return.
func Wrap0RCtx[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context) R) func(context.Context) R {
	return func(ctx context.Context) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result = fn(c.Context())
			c.SetResult(0, result)
		})
		return resolveResult(c, result)
	}
}

// Wrap0SliceR wraps a function with no fixed arguments, variadic slice, and one return value.
func Wrap0SliceR[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func([]any) R) func([]any) R {
	return func(variadicArgs []any) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result = fn(variadicArgs)
			c.SetResult(0, result)
		}, variadicArgs)
		return resolveResult(c, result)
	}
}

// Wrap0SliceRCtx wraps a function with context, no fixed args, variadic slice, one return.
func Wrap0SliceRCtx[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, []any) R) func(context.Context, []any) R {
	return func(ctx context.Context, variadicArgs []any) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result = fn(c.Context(), variadicArgs)
			c.SetResult(0, result)
		}, variadicArgs)
		return resolveResult(c, result)
	}
}

// --- return error

// Wrap0E wraps a function with no arguments and returns error.
func Wrap0E(registry *aspect.Registry, funcKey aspect.FuncKey, fn func() error) func() error {
	return func() error {
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			err = fn()
			c.Error = err
		})
		return resolveError(c, err)
	}
}

// Wrap0ECtx wraps a function with context, no arguments, returns error.
func Wrap0ECtx(registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			err = fn(c.Context())
			c.Error = err
		})
		return resolveError(c, err)
	}
}

// Wrap0SliceE wraps a function with no fixed arguments, variadic slice, and error return.
func Wrap0SliceE(registry *aspect.Registry, funcKey aspect.FuncKey, fn func([]any) error) func([]any) error {
	return func(variadicArgs []any) error {
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			err = fn(variadicArgs)
			c.Error = err
		}, variadicArgs)
		return resolveError(c, err)
	}
}

// Wrap0SliceECtx wraps a function with context, no fixed args, variadic slice, error return.
func Wrap0SliceECtx(registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, []any) error) func(context.Context, []any) error {
	return func(ctx context.Context, variadicArgs []any) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			err = fn(c.Context(), variadicArgs)
			c.Error = err
		}, variadicArgs)
		return resolveError(c, err)
	}
}

// --- return value & error

// Wrap0RE wraps a function with no arguments and returns (result, error).
func Wrap0RE[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func() (R, error)) func() (R, error) {
	return func() (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result, err = fn()
			c.SetResult(0, result)
			c.Error = err
		})
		return resolveResultError(c, result, err)
	}
}

// Wrap0RECtx wraps a function with context, no arguments, returns (result, error).
func Wrap0RECtx[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context) (R, error)) func(context.Context) (R, error) {
	return func(ctx context.Context) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result, err = fn(c.Context())
			c.SetResult(0, result)
			c.Error = err
		})
		return resolveResultError(c, result, err)
	}
}

// Wrap0SliceRE wraps a function with no fixed arguments, variadic slice, and (result, error) return.
func Wrap0SliceRE[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func([]any) (R, error)) func([]any) (R, error) {
	return func(variadicArgs []any) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result, err = fn(variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, variadicArgs)
		return resolveResultError(c, result, err)
	}
}

// Wrap0SliceRECtx wraps a function with context, no fixed args, variadic slice, (result, error) return.
func Wrap0SliceRECtx[R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, []any) (R, error)) func(context.Context, []any) (R, error) {
	return func(ctx context.Context, variadicArgs []any) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result, err = fn(c.Context(), variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, variadicArgs)
		return resolveResultError(c, result, err)
	}
}
