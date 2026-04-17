package wrap

import (
	"testing"
	"time"

	"github.com/seyallius/gosaidno/v2/aspect"
)

// Benchmark_SliceWrapper_NoAdvice measures baseline overhead of slice wrappers
//
//	2447058	       482.7 ns/op	     424 B/op	       9 allocs/op
func Benchmark_SliceWrapper_NoAdvice(b *testing.B) {
	registry := aspect.NewRegistry()
	fn := func(args []any) int {
		sum := 0
		for _, a := range args {
			if val, ok := a.(int); ok {
				sum += val
			}
		}
		return sum
	}
	wrapped := Wrap0SliceR[int](registry, "bench", fn)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = wrapped([]any{1, 2, 3, 4, 5})
	}
}

// Benchmark_SliceWrapper_WithAdvice measures overhead with typical advice chain
//
//	992499	      1025 ns/op	     800 B/op	      14 allocs/op
func Benchmark_SliceWrapper_WithAdvice(b *testing.B) {
	registry := aspect.NewRegistry()
	registry.MustRegister("bench")

	// Add common advice pattern
	registry.MustAddAdvice("bench", aspect.Advice{
		Type: aspect.Before,
		Handler: func(c *aspect.Context) error {
			c.SetMetadataVal("start", time.Now())
			return nil
		},
	})
	registry.MustAddAdvice("bench", aspect.Advice{
		Type: aspect.After,
		Handler: func(c *aspect.Context) error {
			// Log timing
			return nil
		},
	})

	fn := func(args []any) int {
		sum := 0
		for _, a := range args {
			if val, ok := a.(int); ok {
				sum += val
			}
		}
		return sum
	}
	wrapped := Wrap0SliceR[int](registry, "bench", fn)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = wrapped([]any{1, 2, 3, 4, 5})
	}
}

// Benchmark_SliceWrapper_Vs_FixedArity compares slice vs fixed-arity performance
//
//	SliceWrapper         	 2657174	       503.4 ns/op
//	FixedArity           	 3070291	       368.7 ns/op
func Benchmark_SliceWrapper_Vs_FixedArity(b *testing.B) {
	registry := aspect.NewRegistry()

	// Slice version
	sliceFn := func(args []any) int {
		sum := 0
		for _, a := range args {
			if val, ok := a.(int); ok {
				sum += val
			}
		}
		return sum
	}
	wrappedSlice := Wrap0SliceR[int](registry, "slice", sliceFn)

	// Fixed arity version (3 args)
	fixedFn := func(a, b, c int) int {
		return a + b + c
	}
	wrappedFixed := Wrap3R[int, int, int, int](registry, "fixed", fixedFn)

	b.Run("SliceWrapper", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = wrappedSlice([]any{1, 2, 3})
		}
	})

	b.Run("FixedArity", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = wrappedFixed(1, 2, 3)
		}
	})
}
