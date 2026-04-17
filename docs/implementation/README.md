# Implementation Details

This section provides a deep dive into how gosaidno is implemented internally. Understanding these details can help you use the library more effectively and make informed decisions about when and how to use it.

## Table of Contents

- [Design Philosophy](./design-philosophy.md) - Why this approach was chosen and the trade-offs
- [Architecture Overview](./architecture.md) - How the components work together
- [Registry System](./registry.md) - Managing function mappings and thread safety
- [Advice Chain](./advice-chain.md) - Execution orchestration and priority system
- [Context Object](./context.md) - Communication between components
- [Wrapper Functions](./wrappers.md) - The bridge between your code and AOP
- [Execution Engine](./execution-engine.md) - The core orchestration logic
- [Performance Considerations](./performance.md) - Time/space complexity and overhead
- [Limitations and Trade-offs](./limitations.md) - What you need to know about constraints
- [Best Practices](./best-practices.md) - How to use implementation knowledge effectively
- [Extending gosaidno](./extending.md) - How to customize and extend the library