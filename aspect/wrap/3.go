package wrap

import (
	"context"

	"github.com/seyallius/gosaidno/v2/aspect"
)

// --- no return

// Wrap3 wraps a function with three arguments and no return values.
func Wrap3[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C)) func(A, B, C) {
	return func(a A, b B, c C) {
		executeWithAdvice(registry, funcKey, func(ct *aspect.Context) {
			fn(a, b, c)
		}, a, b, c)
	}
}

// Wrap3Ctx wraps a function with context, 3 args, no returns.
func Wrap3Ctx[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C)) func(context.Context, A, B, C) {
	return func(ctx context.Context, a A, b B, c C) {
		executeWithAdviceContext(registry, funcKey, ctx, func(ct *aspect.Context) {
			fn(ct.Context(), a, b, c)
		}, a, b, c)
	}
}

// Wrap3Slice wraps a function with three fixed arguments, variadic slice, and no return values.
func Wrap3Slice[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C, []any)) func(A, B, C, []any) {
	return func(a A, b B, c C, variadicArgs []any) {
		executeWithAdvice(registry, funcKey, func(ct *aspect.Context) {
			fn(a, b, c, variadicArgs)
		}, a, b, c, variadicArgs)
	}
}

// Wrap3SliceCtx wraps a function with context, 3 fixed args, variadic slice, no returns.
func Wrap3SliceCtx[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C, []any)) func(context.Context, A, B, C, []any) {
	return func(ctx context.Context, a A, b B, c C, variadicArgs []any) {
		executeWithAdviceContext(registry, funcKey, ctx, func(ct *aspect.Context) {
			fn(ct.Context(), a, b, c, variadicArgs)
		}, a, b, c, variadicArgs)
	}
}

// --- return value

// Wrap3R wraps a function with three arguments and one return value.
func Wrap3R[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C) R) func(A, B, C) R {
	return func(a A, b B, paramC C) R {
		var result R
		c := executeWithAdvice(registry, funcKey, func(ct *aspect.Context) {
			result = fn(a, b, paramC)
			ct.SetResult(0, result)
		}, a, b, paramC)
		return resolveResult(c, result)
	}
}

// Wrap3RCtx wraps a function with context, 3 args, one return.
func Wrap3RCtx[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C) R) func(context.Context, A, B, C) R {
	return func(ctx context.Context, a A, b B, paramC C) R {
		var result R
		c := executeWithAdviceContext(registry, funcKey, ctx, func(ct *aspect.Context) {
			result = fn(ct.Context(), a, b, paramC)
			ct.SetResult(0, result)
		}, a, b, paramC)
		return resolveResult(c, result)
	}
}

// Wrap3SliceR wraps a function with three fixed arguments, variadic slice, and one return value.
func Wrap3SliceR[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C, []any) R) func(A, B, C, []any) R {
	return func(a A, b B, paramC C, variadicArgs []any) R {
		var result R
		ct := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result = fn(a, b, paramC, variadicArgs)
			c.SetResult(0, result)
		}, a, b, paramC, variadicArgs)
		return resolveResult(ct, result)
	}
}

// Wrap3SliceRCtx wraps a function with context, 3 fixed args, variadic slice, one return.
func Wrap3SliceRCtx[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C, []any) R) func(context.Context, A, B, C, []any) R {
	return func(ctx context.Context, a A, b B, paramC C, variadicArgs []any) R {
		var result R
		ct := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result = fn(c.Context(), a, b, paramC, variadicArgs)
			c.SetResult(0, result)
		}, a, b, paramC, variadicArgs)
		return resolveResult(ct, result)
	}
}

// --- return error

// Wrap3E wraps a function with three arguments and returns error.
func Wrap3E[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C) error) func(A, B, C) error {
	return func(a A, b B, c C) error {
		var err error
		ctx := executeWithAdvice(registry, funcKey, func(ct *aspect.Context) {
			err = fn(a, b, c)
			ct.Error = err
		}, a, b, c)
		return resolveError(ctx, err)
	}
}

// Wrap3ECtx wraps a function with context, 3 args, returns error.
func Wrap3ECtx[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C) error) func(context.Context, A, B, C) error {
	return func(ctx context.Context, a A, b B, c C) error {
		var err error
		ct := executeWithAdviceContext(registry, funcKey, ctx, func(ct *aspect.Context) {
			err = fn(ct.Context(), a, b, c)
			ct.Error = err
		}, a, b, c)
		return resolveError(ct, err)
	}
}

// Wrap3SliceE wraps a function with three fixed arguments, variadic slice, and error return.
func Wrap3SliceE[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C, []any) error) func(A, B, C, []any) error {
	return func(a A, b B, c C, variadicArgs []any) error {
		var err error
		ct := executeWithAdvice(registry, funcKey, func(aCtx *aspect.Context) {
			err = fn(a, b, c, variadicArgs)
			aCtx.Error = err
		}, a, b, c, variadicArgs)
		return resolveError(ct, err)
	}
}

// Wrap3SliceECtx wraps a function with context, 3 fixed args, variadic slice, error return.
func Wrap3SliceECtx[A, B, C any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C, []any) error) func(context.Context, A, B, C, []any) error {
	return func(ctx context.Context, a A, b B, c C, variadicArgs []any) error {
		var err error
		ct := executeWithAdviceContext(registry, funcKey, ctx, func(aCtx *aspect.Context) {
			err = fn(aCtx.Context(), a, b, c, variadicArgs)
			aCtx.Error = err
		}, a, b, c, variadicArgs)
		return resolveError(ct, err)
	}
}

// --- return value & error

// Wrap3RE wraps a function with three arguments and returns (result, error).
func Wrap3RE[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C) (R, error)) func(A, B, C) (R, error) {
	return func(a A, b B, paramC C) (R, error) {
		var result R
		var err error
		c := executeWithAdvice(registry, funcKey, func(ct *aspect.Context) {
			result, err = fn(a, b, paramC)
			ct.SetResult(0, result)
			ct.Error = err
		}, a, b, paramC)
		return resolveResultError(c, result, err)
	}
}

// Wrap3RECtx wraps a function with context, 3 args, returns (result, error).
func Wrap3RECtx[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C) (R, error)) func(context.Context, A, B, C) (R, error) {
	return func(ctx context.Context, a A, b B, paramC C) (R, error) {
		var result R
		var err error
		c := executeWithAdviceContext(registry, funcKey, ctx, func(ct *aspect.Context) {
			result, err = fn(ct.Context(), a, b, paramC)
			ct.SetResult(0, result)
			ct.Error = err
		}, a, b, paramC)
		return resolveResultError(c, result, err)
	}
}

// Wrap3SliceRE wraps a function with three fixed arguments, variadic slice, and (result, error) return.
func Wrap3SliceRE[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(A, B, C, []any) (R, error)) func(A, B, C, []any) (R, error) {
	return func(a A, b B, paramC C, variadicArgs []any) (R, error) {
		var result R
		var err error
		ct := executeWithAdvice(registry, funcKey, func(c *aspect.Context) {
			result, err = fn(a, b, paramC, variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, a, b, paramC, variadicArgs)
		return resolveResultError(ct, result, err)
	}
}

// Wrap3SliceRECtx wraps a function with context, 3 fixed args, variadic slice, (result, error) return.
func Wrap3SliceRECtx[A, B, C, R any](registry *aspect.Registry, funcKey aspect.FuncKey, fn func(context.Context, A, B, C, []any) (R, error)) func(context.Context, A, B, C, []any) (R, error) {
	return func(ctx context.Context, a A, b B, paramC C, variadicArgs []any) (R, error) {
		var result R
		var err error
		ct := executeWithAdviceContext(registry, funcKey, ctx, func(c *aspect.Context) {
			result, err = fn(c.Context(), a, b, paramC, variadicArgs)
			c.SetResult(0, result)
			c.Error = err
		}, a, b, paramC, variadicArgs)
		return resolveResultError(ct, result, err)
	}
}
