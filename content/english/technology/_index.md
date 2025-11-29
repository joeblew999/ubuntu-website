---
title: "Technology"
meta_title: "Technology Stack | Ubuntu Software"
description: "Built on Go and NATS JetStream - our technology choices for performance, reliability, and simplicity."
image: "/images/spatial.svg"
draft: false
---

## Our Technology Stack

We build on technologies chosen for performance, reliability, and long-term maintainability.

---

## Go

**Why Go:**

- **Performance** — Compiled, statically typed, minimal runtime overhead
- **Simplicity** — One way to do things, readable by default
- **Concurrency** — Goroutines and channels built into the language
- **Deployment** — Single binary, no dependencies, cross-compilation
- **Ecosystem** — Strong standard library, excellent tooling

Go is our primary language across the stack—from backend services to CLI tools to robotics integration.

---

## NATS JetStream

**Why NATS JetStream:**

- **Persistence** — Durable message storage with replay
- **Exactly-once delivery** — Reliable message processing
- **Lightweight** — Single binary, minimal resource footprint
- **Scalable** — Clustering and horizontal scaling built-in
- **Real-time** — Sub-millisecond latency for pub/sub

NATS JetStream powers our event-driven architecture—connecting design tools, simulation, digital twins, and physical devices.

---

## Architecture Principles

| Principle | Implementation |
|-----------|----------------|
| **Offline-first** | Local-first data, sync when connected |
| **Event-driven** | NATS JetStream for all inter-service communication |
| **Open standards** | STEP, IFC, no proprietary formats |
| **Hardware-agnostic** | Abstraction layers for portability |
| **Self-sovereign** | Deploy anywhere—cloud, on-prem, air-gapped |

---

## Technology Areas

### Robotics

Our robotics stack is built on Viam RDK—an open-source robotics development kit that provides hardware abstraction, motion planning, and computer vision.

[Robotics Stack →](/technology/robotics/)

---

## Open Source Foundation

We build on and contribute to open source:

| Component | Technology |
|-----------|------------|
| Language | Go |
| Messaging | NATS JetStream |
| Collaboration | Automerge (CRDT) |
| 3D Formats | STEP, IFC |
| Robotics | Viam RDK |
| Vision | Intel RealSense, YOLOv8 |

---

## Learn More

- [Robotics Stack →](/technology/robotics/)
- [Spatial Platform →](/platform/spatial/)
- [Foundation →](/platform/foundation/)
- [Contact Us →](/contact/)
