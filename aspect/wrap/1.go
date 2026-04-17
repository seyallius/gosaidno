package wrap

import (
	"context"

	"github.com/seyallius/gosaidno/v2/aspect"
)

// --- no return

// Wrap1 wraps a function with one argument and no return values.
func Wrap1[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A)) func(A) {
	return func(a A) {
		executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			fn(a)
		}, a)
	}
}

// Wrap1Ctx wraps a function with context, 1 arg, no returns.
func Wrap1Ctx[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A)) func(context.Context, A) {
	return func(ctx context.Context, a A) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			fn(c.Context(), a)
		}, a)
	}
}

// Wrap1Slice wraps a function with one fixed argument, variadic slice, and no return values.
func Wrap1Slice[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, []any)) func(A, []any) {
	return func(a A, variadicArgs []any) {
		executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			fn(a, variadicArgs)
		}, a, variadicArgs)
	}
}

// Wrap1SliceCtx wraps a function with context, 1 fixed arg, variadic slice, no returns.
func Wrap1SliceCtx[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, []any)) func(context.Context, A, []any) {
	return func(ctx context.Context, a A, variadicArgs []any) {
		executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			fn(c.Context(), a, variadicArgs)
		}, a, variadicArgs)
	}
}

// --- return value

// Wrap1R wraps a function with one argument and one return value.
func Wrap1R[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A) R) func(A) R {
	return func(a A) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result = fn(a)
			c.SetResult(0, result)
		}, a)
		return resolveResult(c, result)
	}
}

// Wrap1RCtx wraps a function with context, 1 arg, one return.
func Wrap1RCtx[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A) R) func(context.Context, A) R {
	return func(ctx context.Context, a A) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result = fn(c.Context(), a)
			c.SetResult(0, result)
		}, a)
		return resolveResult(c, result)
	}
}

// Wrap1SliceR wraps a function with one fixed argument, variadic slice, and one return value.
func Wrap1SliceR[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, []any) R) func(A, []any) R {
	return func(a A, variadicArgs []any) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result = fn(a, variadicArgs)
			c.SetResult(0, result)
		}, a, variadicArgs)
		return resolveResult(c, result)
	}
}

// Wrap1SliceRCtx wraps a function with context, 1 fixed arg, variadic slice, one return.
func Wrap1SliceRCtx[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, []any) R) func(context.Context, A, []any) R {
	return func(ctx context.Context, a A, variadicArgs []any) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result = fn(c.Context(), a, variadicArgs)
			c.SetResult(0, result)
		}, a, variadicArgs)
		return resolveResult(c, result)
	}
}

// --- return error

// Wrap1E wraps a function with one argument and returns error.
func Wrap1E[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A) error) func(A) error {
	return func(a A) error {
		var err error
		executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			err = fn(a)
			c.Error = err
		}, a)
		return err // executeWithAdvice returns *aspect.Context, but simplistic wrapper can just return captured err if no complex mutation needed
	}
}

// Wrap1ECtx wraps a function with context, 1 arg, returns error.
func Wrap1ECtx[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A) error) func(context.Context, A) error {
	return func(ctx context.Context, a A) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			err = fn(c.Context(), a)
			c.Error = err
		}, a)
		return resolveError(c, err)
	}
}

// Wrap1SliceE wraps a function with one fixed argument, variadic slice, and error return.
func Wrap1SliceE[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, []any) error) func(A, []any) error {
	return func(a A, variadicArgs []any) error {
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			err = fn(a, variadicArgs)
			c.Error = err
		}, a, variadicArgs)
		return resolveError(c, err)
	}
}

// Wrap1SliceECtx wraps a function with context, 1 fixed arg, variadic slice, error return.
func Wrap1SliceECtx[A any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, []any) error) func(context.Context, A, []any) error {
	return func(ctx context.Context, a A, variadicArgs []any) error {
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			err = fn(c.Context(), a, variadicArgs)
			c.Error = err
		}, a, variadicArgs)
		return resolveError(c, err)
	}
}

// --- return value & error

// Wrap1RE wraps a function with one argument and returns (result, error).
func Wrap1RE[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A) (R, error)) func(A) (R, error) {
	return func(a A) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result, err = fn(a)
			c.SetResult(0, result)
			c.Error = err
		}, a)
		return resolveResultError(c, result, err)
	}
}

// Wrap1RECtx wraps a function with context, 1 arg, returns (result, error).
func Wrap1RECtx[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A) (R, error)) func(context.Context, A) (R, error) {
	return func(ctx context.Context, a A) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result, err = fn(c.Context(), a)
			c.SetResult(0, result)
			c.Error = err
		}, a)
		return resolveResultError(c, result, err)
	}
}

// Wrap1SliceRE wraps a function with one fixed argument, variadic slice, and (result, error) return.
func Wrap1SliceRE[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, []any) (R, error)) func(A, []any) (R, error) {
	return func(a A, variadicArgs []any) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result, err = fn(a, variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, a, variadicArgs)
		return resolveResultError(c, result, err)
	}
}

// Wrap1SliceRECtx wraps a function with context, 1 fixed arg, variadic slice, (result, error) return.
func Wrap1SliceRECtx[A, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, []any) (R, error)) func(context.Context, A, []any) (R, error) {
	return func(ctx context.Context, a A, variadicArgs []any) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result, err = fn(c.Context(), a, variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, a, variadicArgs)
		return resolveResultError(c, result, err)
	}
}
