---
title: "Foundation"
meta_title: "Foundation Technology | Ubuntu Software"
description: "Offline-first architecture, universal deployment, and self-sovereign data. The technology foundation that powers both Publish and Spatial platforms."
image: "/images/foundation.svg"
draft: false
---

## Built for the Real World

Internet goes down. Teams span continents. Servers belong to you. We built for this reality.

Our foundation isn't about features—it's about how software should work. Offline-first. Self-sovereign. Universal deployment. These principles run through everything we build.

---

## Offline-First Architecture

**Work without internet. Sync when connected. Never lose data.**

Real work happens in places with bad WiFi—construction sites, factory floors, government buildings, hospital rooms, remote offices. Our foundation assumes disconnection, not connectivity.

### How It Works

| Component | What It Does |
|-----------|--------------|
| **Local-First Data** | Your data lives on your device first, not in a distant server |
| **Automerge CRDT** | Conflict-free merging when multiple people edit simultaneously |
| **Background Sync** | Automatic synchronization when connectivity returns |
| **Offline Queues** | Actions queue locally, execute when possible |

**No spinners. No "connection lost" errors. Just work.**

---

## Deploy Anywhere

**One codebase. Every platform. Native experience.**

Same application runs in the browser, on the desktop, and on mobile devices. Not three separate products—one codebase that deploys everywhere.

### Supported Platforms

| Platform | Delivery |
|----------|----------|
| **Web** | Any modern browser—Chrome, Firefox, Safari, Edge |
| **Desktop** | Native apps for Windows, macOS, Linux |
| **Mobile** | Native apps for iOS and Android |

Your team uses desktops in the office. Field workers use tablets. Customers use phones. Everyone accesses the same system with the same data.

**Write once. Deploy everywhere. Maintain one codebase.**

---

## Cloud Sync Options

**Your cloud. Our cloud. No cloud. Your choice.**

Sync doesn't require our servers. Connect to whatever infrastructure makes sense for your organization.

### Deployment Models

| Model | Best For |
|-------|----------|
| **Ubuntu Software Cloud** | Fastest setup, we handle operations |
| **Your Cloud** | AWS, Azure, GCP—your infrastructure, our software |
| **On-Premises** | Your data center, complete control |
| **Hybrid** | Some data in cloud, sensitive data on-prem |
| **Air-Gapped** | Fully disconnected networks, defense and secure environments |

Switching between models? Straightforward. Your data format doesn't change based on where it lives.

**No vendor lock-in. No forced cloud dependency. Real deployment flexibility.**

---

## Self-Sovereign

**Your data. Your servers. Your rules.**

Self-sovereign means you control your infrastructure. Not "your data hosted on our terms"—actually yours.

### What Self-Sovereign Means

- **Run anywhere** — Your data center, your cloud account, your laptop
- **No call-home** — Software works without contacting our servers
- **Export everything** — Standard formats, complete data portability
- **No usage tracking** — We don't see what you do with our software
- **Perpetual license options** — Keep running even if we disappear

**This isn't just privacy—it's operational independence.**

---

## Embeddable

**Integrate into existing systems. Don't rip and replace.**

Organizations have existing workflows, databases, identity systems. Our foundation integrates rather than replaces.

### Integration Patterns

| Pattern | Use Case |
|---------|----------|
| **API-First** | Everything accessible programmatically |
| **Database Connectors** | Read/write to your existing databases |
| **SSO Integration** | Your identity provider, not another login |
| **Webhook Events** | Push notifications to your systems |
| **White-Label** | Your branding, our engine |

**Extend your systems. Don't abandon them.**

---

## Wellknown Gateway

**Publish to Big Tech platforms. Own the relationship.**

The web's gateways—Google, Apple, YouTube—have massive reach. But publishing TO them doesn't mean being owned BY them.

### Reverse the Relationship

Traditional approach:
```
User → YouTube (owns everything) → Your content (captive)
```

Wellknown approach:
```
User → Your Gateway → Your System (primary)
                   ↳→ YouTube (mirror for discovery)
```

**You control the front door.** Big Tech becomes optional distribution, not a prison.

### How It Works

| Capability | What It Means |
|------------|---------------|
| **Your URIs everywhere** | Links point to YOUR gateway, not theirs |
| **Smart routing** | Send iOS users to Apple, Android to Google, web to your player |
| **Mirror publishing** | Auto-publish copies to YouTube, Google Maps, Apple Calendar |
| **Analytics you own** | See everything, track everyone, no data sharing |
| **Exit strategy built-in** | Remove any platform from routing without breaking links |

### Works With Everything

- **Video** — Host on your server, mirror to YouTube/Twitch for reach
- **Calendar** — Your CalDAV server, sync to Google/Apple for convenience
- **Maps** — Your geographic data, integrate with Google/Apple Maps
- **Email** — Your mail server, compatible with Gmail
- **Files** — Your storage, selective sharing to Drive/Dropbox

**Publish TO their platforms. Never be locked IN.**

---

## 25 Years of Enterprise Experience

**Ubuntu Software started in 1999.** We've built enterprise systems through dot-com, mobile revolution, cloud transition, and AI emergence.

### What Experience Taught Us

- **Vendors disappear** — Build for data portability from day one
- **Networks fail** — Offline-first isn't optional, it's essential
- **Requirements change** — Open standards outlast proprietary formats
- **Scale surprises** — Architecture matters more than optimization
- **Integration is hard** — Design for it, don't bolt it on

We've seen what works and what doesn't. This foundation reflects decades of lessons learned in production environments.

---

## The Foundation Powers Everything

Both Publish and Spatial platforms inherit these capabilities automatically:

| Capability | Publish | Spatial |
|------------|---------|---------|
| Offline editing | Edit documents without internet | Design 3D models without internet |
| Real-time sync | Multiple editors, one document | Multiple designers, one model |
| Universal deploy | Forms on any device | 3D viewer on any device |
| Self-hosted | Your document server | Your design server |
| Cloud options | Managed or self-managed | Managed or self-managed |

**Choose your platform. Get the foundation free.**

---

## Technical Details

For teams evaluating our architecture:

| Layer | Technology |
|-------|------------|
| **Storage** | SQLite everywhere—local devices and servers |
| **Sync Engine** | CRDT-based replication via NATS JetStream |
| **Messaging** | NATS JetStream |
| **UI Framework** | Cross-platform native rendering |
| **API** | HTTP REST + SSE (Server-Sent Events) |
| **AI Integration** | Model Context Protocol (MCP) |
| **Auth** | OIDC-compatible, bring your own IdP |

### Distributed SQLite

Every node—your laptop, your phone, your servers—runs SQLite. Changes replicate via NATS JetStream using CRDT semantics.

- **Any server can fail** — Others keep serving
- **Any server can go offline** — Sync when reconnected
- **No single point of failure** — True distributed architecture
- **Same database everywhere** — Local device to global cluster

**Standard technologies. No proprietary lock-in.**

---

## Scale Without Limits

### No Single Point of Failure (SPOF)

Every component is redundant. No single server, service, or data center can take down the system. Fail any node—the system keeps running.

### No Single Point of Performance (SPOP)

Computation scales horizontally. Add capacity by adding nodes, not by buying bigger servers. Workloads distribute automatically across available resources.

### Hundreds of Data Centers

Architecture designed for global distribution:

| Capability | What It Means |
|------------|---------------|
| **Deploy anywhere** | Cloud, on-prem, edge, air-gapped |
| **Deploy close to users** | Low latency, local compliance |
| **Replicate for redundancy** | Survive regional outages |
| **Partition tolerance** | Operate independently when networks split |

### Deployment Options

| Method | Use Case |
|--------|----------|
| **Binaries** | Single-file deployment, minimal dependencies |
| **Docker** | Containerized, reproducible environments |
| **Kubernetes** | Orchestrated, auto-scaling clusters |

**From a single laptop to hundreds of data centers. Same architecture. Same codebase.**

---

## Get Started

The foundation is built in. When you use Publish or Spatial, you get offline-first, universal deployment, and self-sovereign options automatically.

[Explore Publish →](/platform/publish/) | [Explore Spatial →](/platform/spatial/) | [Linux & Cross-Platform →](/technology/linux/) | [Contact Us →](/contact/)
