Session:
cmd/aios/runtime/session.go
hardware adapters
network listeners

Agent:
core/agent
signal interpretation
algorithm distillation

Optimizer:
core/optimization
platform-specific pruning

Router:
core/router
dispatch only

Modules:
modules/*
domain logic
storage interaction

Repositories:
storage/*
DB abstraction

13. Concurrency Model Summary
Modules run continuously under Supervisor.
Session receives external data concurrently.
Agent + Router are synchronous per message.
Storage may be asynchronous internally.
No cross-layer tight coupling.

14. Correct Mental Model
Think of it as:
Device Drivers → Signal Processor → Dispatcher → Domain Engine → Persistence
Each layer has exactly one responsibility.