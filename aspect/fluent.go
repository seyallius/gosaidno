// Package aspect. fluent provides a fluent/declarative API for registering advice
package aspect

import "context"

// -------------------------------------------- Types --------------------------------------------

// FluentBuilder provides a fluent API for configuring advice for a function.
type FluentBuilder struct {
	registry *Registry
	funcKey  FuncKey
}

// -------------------------------------------- Public Functions --------------------------------------------

// For creates a new fluent builder for the given function name.
// It ensures the function is registered in the registry and returns a builder
// for adding advice in a fluent manner.
// Use ForWithRegistry if you don't want to use the default Registry instance (recommended).
func For(funcName FuncKey) *FluentBuilder {
	registry := DefaultRegistry()
	return &FluentBuilder{
		registry: registry,
		funcKey:  funcName,
	}
}

// ForWithRegistry creates a new fluent builder for the given function name using a specific registry.
// It ensures the function is registered in the provided registry and returns a builder
// for adding advice in a fluent manner.
func ForWithRegistry(registry *Registry, funcName FuncKey) *FluentBuilder {
	return &FluentBuilder{
		registry: registry,
		funcKey:  funcName,
	}
}

// WithBefore adds a Before advice to the function.
func (fb *FluentBuilder) WithBefore(handler AdviceFunc) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:    Before,
		Handler: handler,
	})
	return fb
}

// WithBeforeP adds a Before advice with a specific priority to the function.
func (fb *FluentBuilder) WithBeforeP(handler AdviceFunc, priority int) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:     Before,
		Handler:  handler,
		Priority: priority,
	})
	return fb
}

// WithAfter adds an After advice to the function.
func (fb *FluentBuilder) WithAfter(handler AdviceFunc) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:    After,
		Handler: handler,
	})
	return fb
}

// WithAfterP adds an After advice with a specific priority to the function.
func (fb *FluentBuilder) WithAfterP(handler AdviceFunc, priority int) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:     After,
		Handler:  handler,
		Priority: priority,
	})
	return fb
}

// WithAround adds an Around advice to the function.
func (fb *FluentBuilder) WithAround(handler AdviceFunc) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:    Around,
		Handler: handler,
	})
	return fb
}

// WithAroundP adds an Around advice with a specific priority to the function.
func (fb *FluentBuilder) WithAroundP(handler AdviceFunc, priority int) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:     Around,
		Handler:  handler,
		Priority: priority,
	})
	return fb
}

// WithAfterReturning adds an AfterReturning advice to the function.
func (fb *FluentBuilder) WithAfterReturning(handler AdviceFunc) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:    AfterReturning,
		Handler: handler,
	})
	return fb
}

// WithAfterReturningP adds an AfterReturning advice with a specific priority to the function.
func (fb *FluentBuilder) WithAfterReturningP(handler AdviceFunc, priority int) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:     AfterReturning,
		Handler:  handler,
		Priority: priority,
	})
	return fb
}

// WithAfterThrowing adds an AfterThrowing advice to the function.
func (fb *FluentBuilder) WithAfterThrowing(handler AdviceFunc) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:    AfterThrowing,
		Handler: handler,
	})
	return fb
}

// WithAfterThrowingP adds an AfterThrowing advice with a specific priority to the function.
func (fb *FluentBuilder) WithAfterThrowingP(handler AdviceFunc, priority int) *FluentBuilder {
	fb.registry.RegisterOrGet(fb.funcKey)
	fb.registry.MustAddAdvice(fb.funcKey, Advice{
		Type:     AfterThrowing,
		Handler:  handler,
		Priority: priority,
	})
	return fb
}

// GetRegistry returns the registry used by this fluent builder.
// This allows users to call the appropriate Wrap methods on the registry.
func (fb *FluentBuilder) GetRegistry() *Registry {
	return fb.registry
}

// GetFuncKey returns the function key used by this fluent builder.
func (fb *FluentBuilder) GetFuncKey() FuncKey {
	return fb.funcKey
}

// The Wrap method is intentionally omitted to avoid reflection usage.
// Users should use the specific typed Wrap methods like Wrap0, Wrap1, etc.
// based on their function signatures for type safety.

// Context-aware wrapper methods

// Wrap0Ctx wraps a function with no arguments and no return values with context.
func (fb *FluentBuilder) Wrap0Ctx(fn func(context.Context)) func(context.Context) {
	fb.registry.MustRegister(fb.funcKey)
	return Wrap0Ctx(fb.registry, fb.funcKey, fn)
}

// Wrap0ECtx wraps a function with no arguments and returns error with context.
func (fb *FluentBuilder) Wrap0ECtx(fn func(context.Context) error) func(context.Context) error {
	fb.registry.MustRegister(fb.funcKey)
	return Wrap0ECtx(fb.registry, fb.funcKey, fn)
}

// The generic context-aware methods are removed since Go doesn't support generic methods.
// Users should call the registry's context-aware Wrap functions directly after configuring advice.
// Example:
// aspect.For("MyFunction").
//     WithBefore(myBeforeAdvice).
//     WithAfter(myAfterAdvice)
//
// builder := aspect.For("MyFunction")
// wrappedFn := wrap.Wrap1Ctx(builder.GetRegistry(), builder.GetFuncKey(), myFunction)
