---
title: "Publish"
meta_title: "Content Management Platform | Ubuntu Software"
description: "Single-source CMS for multi-channel publishing. One DSL for text, graphics, forms, emails. Auto-translated outputs. Perfect branding by architecture."
image: "/images/publish.svg"
draft: false
---

## One Language. Every Output. Every Language.

A single-source content management system. Define everything in one DSL—**text (1D) that drives all 2D and 3D outputs**. Compose components that include components. Output to any channel. Auto-translate to any language. Capture data back.

All web-capable. All globally accessible. Perfect branding—not by discipline, by architecture.

---

## The Problem

Organizations publish content across multiple channels: websites, PDFs, emails, forms, signage.

Each channel maintained separately. Content drifts. Branding fractures. Translations lag. Form submissions get re-keyed manually.

**One change means updating five systems. Or it doesn't get updated at all.**

---

## One DSL

Everything is defined in the same language.

**Text content:**
```
# Welcome to Our Service

We help organizations {industry} achieve {outcome}.
```

**Graphics:**
```
@logo: svg {
  viewBox: "0 0 200 50"
  rect: { x: 0, y: 0, width: 50, height: 50, fill: "#2563eb" }
  text: { x: 60, y: 35, content: "Ubuntu", font: "bold 24px" }
}
```

**Form fields:**
```
Name: [_______________]{field: name, required: true}
Email: [_______________]{field: email, type: email}
Date: [__/__/____]{field: date, type: date}
```

**Emails:**
```
@welcome-email: email {
  to: {customer.email}
  subject: "Welcome, {customer.name}"
  body: include welcome-content
}
```

**One language. Text, graphics, forms, emails, documents.**

---

## Compose

Components include components.

```
@header: compose {
  include: logo
  include: navigation
  include: search-bar
}

@page-template: compose {
  include: header
  include: content
  include: footer
}

@welcome-letter: compose {
  include: page-template
  content: "Dear {customer.name}, welcome to..."
}

@application-form: compose {
  include: page-template
  content: include form-fields
}
```

**The cascade:**
- Change the logo → header updates
- Header updates → every page updates
- Every page updates → every PDF, email, form updates

**This is why branding stays perfect.** Not because someone remembers to update everything. Because the architecture makes it impossible not to.

Link to any existing content. Reuse across the system. One source of truth.

---

## Output

From one source, generate everything:

| Output | What Happens |
|--------|--------------|
| Web pages | Rendered to HTML, SEO-ready |
| PDFs | Print-ready, archival |
| Emails | Sent directly to recipients |
| Web forms | Interactive, data captures back |
| PDF forms | Fillable or printable, OCR captures back |
| SVG graphics | Vector assets, any size |
| Maps | Your location data visualized |
| Kiosks | Physical displays, real-time |

**Every output auto-translates.**

Write in English. German colleagues see German. Spanish partners see Spanish. Japanese customers see Japanese.

Not translated after the fact. Translated as part of rendering. Every output. Every language. Automatically.

---

## Translate

Translation isn't a feature. It's how the system works.

| Capability | What It Means |
|------------|---------------|
| Real-time in editor | Collaborate across languages simultaneously |
| Auto-translate outputs | Every format renders in every language |
| Offline AI | Works without internet |
| Contextual | Knows education vs. legal vs. medical terminology |
| Bi-directional | They edit in their language, you see yours |

**Write once. Publish in every language.**

---

## Capture

Forms capture data back. Because we generated the form, we know exactly where every field is.

| Channel | What Happens |
|---------|--------------|
| Web form submitted | Data flows to your database |
| PDF form filled digitally | Data extracts to your database |
| PDF form printed, filled by hand, scanned | OCR captures to your database |

No training the OCR. No field mapping. We generated the form—we know where everything is.

**Paper or digital. Same data. Same destination.**

---

## Human-in-Loop: Real-Time Field Operations

Publish isn't just for static documents. It's how field workers interact with your systems in real-time.

### The Fully-Tapped Device

Every field worker's device becomes both a sensor and a display:

| Input (from field) | Output (to field) |
|--------------------|-------------------|
| GPS position | Route maps, turn-by-turn |
| Voice commands (SST) | Task notifications |
| Touch interactions | Schedule updates |
| Form submissions | Confirmation alerts |
| Camera (photos/video) | Status displays |

The system knows where every worker is. Workers can talk back—via voice, touch, or forms.

### 2D Outputs

Publish renders all 2D content—web, documents, maps, displays:

| Output | Use Case | Example |
|--------|----------|---------|
| **Web pages** | Customer portals, status tracking | Service status, order tracking |
| **Route maps** | Navigation, delivery tracking | Waste collection, delivery stops |
| **Dashboards** | Status overview, KPIs | Dispatch console, team performance |
| **Documents** | Work orders, reports, forms | PDF manifests, compliance exports |
| **Kiosks** | Public displays, team boards | Depot status, restaurant menus |

For **3D spatial visualization**—construction sites, warehouses, digital twins—see the [Spatial Platform](/platform/spatial/).

### Voice-Driven Updates

Speech-to-Text (SST) enables hands-free system interaction:

| Worker Says | System Does |
|-------------|-------------|
| "123 Main blocked by construction" | Marks location, refactors routes |
| "Skip next task, customer not available" | Removes from queue, adjusts schedule |
| "Equipment needs repair" | Creates maintenance ticket |
| "Task complete" | Records completion, triggers next workflow |

No phone calls to dispatch. No radio chatter. The worker speaks, the AI listens, the system adapts.

### Authority Hierarchy

Human judgment overrides AI recommendations:

```
┌─────────────────────────────────────────────────────────────────┐
│                     DECISION AUTHORITY                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1. FIELD WORKER       ───────────────────────▶  HIGHEST       │
│      Eyes on ground, can override anything                       │
│                                                                  │
│   2. DISPATCHER         ───────────────────────▶  HIGH          │
│      Fleet-wide view, can reassign work                          │
│                                                                  │
│   3. AI OPTIMIZER       ───────────────────────▶  MEDIUM        │
│      Suggests routes/tasks, adapts to changes                    │
│                                                                  │
│   4. SCHEDULED PLAN     ───────────────────────▶  LOWEST        │
│      Starting point, always subject to reality                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**The AI is powerful—but the human with eyes on the ground has final say.**

### Presence & Notifications

Humans aren't always available. The system knows who is:

| Capability | What It Does |
|------------|--------------|
| **Presence tracking** | Who's online, who's offline, who's busy |
| **Role awareness** | Who can approve, who can override, who's on-call |
| **Smart routing** | Notifications go to available personnel first |
| **Escalation** | If primary doesn't respond, escalate to backup |
| **Acknowledgment** | Track who saw what, when they responded |

**Notification channels:**

- Push notifications (mobile/desktop)
- SMS for critical alerts
- Email for async updates
- In-app alerts with sound/haptics
- Kiosk displays for team visibility

The system adapts to human availability. Urgent issues find the right person. Non-urgent issues queue until someone's ready.

### Offline Resilience

Field work happens where networks fail. The device keeps working:

- **Connectivity lost** → App continues with local SQLite
- **Changes queue locally** → Sync when reconnected
- **No data loss** → Automerge CRDT handles conflicts
- **Presence updates** → Marked offline, notifications queue

Field operations can't stop for network issues. This architecture ensures they don't.

[See Field Operations in action →](/applications/logistics/)

---

## Connect

Your database. Your schema.

- Read data to populate documents: "Dear {customer.name}..."
- Write captured data back: form submissions → your tables
- Any SQL database: PostgreSQL, MySQL, SQLite
- Your infrastructure: self-hosted or cloud

**We don't replace your backend. We connect to it.**

---

## Industries

| Industry | Why Publish |
|----------|-------------|
| [**Government**](/applications/government) | Serve every citizen. Paper + digital. Accessible. Compliant. |
| [**Healthcare**](/applications/healthcare) | Less administration. Fewer errors. Better care. |
| [**Financial**](/applications/financial) | Compliance without complexity. Complete audit trails. |
| [**Education**](/applications/education) | Modern administration. Traditional options. |
| [**Insurance**](/applications/insurance) | Every channel. One system. Field to office unified. |

---

## Architecture

| Layer | What It Does |
|-------|--------------|
| DSL parser | Understands text, graphics, forms, composition |
| Include resolver | Recursive component composition |
| Renderers | Web, PDF, email, SVG, form outputs |
| Translation engine | Real-time AI, offline capable |
| Field mapper | Knows exact position of every form field |
| OCR engine | Precise extraction using field map |
| Database connectors | Read and write to your schema |

**No proprietary formats. Export everything. Own your content.**

---

## Native on Every Platform

| Platform | Experience |
|----------|------------|
| **Linux** | Native desktop, embedded kiosks, edge devices |
| **Windows** | Native Windows application |
| **macOS** | Native Mac application |
| **iOS** | Native mobile app for field workers |
| **Android** | Native mobile app for field workers |
| **Web** | Modern browser, full functionality |

One codebase. Native performance on every platform. Authors use desktops. Field workers use tablets. Customers use phones. Same system. Same data.

---

## Beyond Documents

The same single-source approach works for:

**Maps** — Your geographic data, your styling, integrate with Google/Apple Maps for reach.

**Video** *(coming soon)* — Your video server, mirror to YouTube/Vimeo for discovery.

**Calendar** *(coming soon)* — Your CalDAV server, sync to Google/Apple Calendar for convenience.

**Kiosks** *(coming soon)* — Raspberry Pi displays, restaurant menus, government offices, retail signage.

All links route through [your gateway](#wellknown-gateway-data-sovereignty). Publish TO Big Tech platforms. Never be locked IN.

---

## WellKnown Gateway: Data Sovereignty

Where does your published content go? External platforms? That's where the **WellKnown Gateway** comes in—ensuring you own your data even when publishing TO other systems.

### The Problem with Platform Lock-In

Traditional approach:
- Customer data lives in Salesforce
- Documents live in Google Drive
- Forms live in Typeform
- You're dependent on each vendor's API
- If any vendor changes terms, raises prices, or shuts down—you scramble

### The WellKnown Solution

```
┌─────────────────────────────────────────────────────────────────┐
│                    YOUR DATA (Source of Truth)                   │
│                                                                  │
│   SQLite on every device + NATS JetStream persistence            │
│   You own it. It's on your infrastructure.                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    WELLKNOWN GATEWAY                             │
│                                                                  │
│   Publishes TO external platforms                                │
│   Syncs FROM external platforms                                  │
│   You control what goes where                                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
   │   Google    │     │   Social    │     │  Industry   │
   │   Drive     │     │   Media     │     │  Platforms  │
   └─────────────┘     └─────────────┘     └─────────────┘
```

**Key principle:** You publish TO platforms. They don't own your data—you do.

### Customer Relationship Sovereignty

Your customers interact with your brand, not a vendor's:

| Without WellKnown | With WellKnown |
|-------------------|----------------|
| Customer logs into vendor portal | Customer logs into YOUR portal |
| Vendor owns the relationship | You own the relationship |
| Vendor has your customer list | You control customer data |
| Switching vendors = rebuilding | Switching vendors = reconnect gateway |

The customer thinks they're talking to you. They are. The plumbing is invisible.

### Built on Internet Standards

WellKnown uses the same standards that power email, calendars, and federated social networks:

| Standard | Purpose | Example |
|----------|---------|---------|
| **RFC 8615** | `/.well-known/` URI discovery | `yourdomain.com/.well-known/...` |
| **WebFinger** | Account/resource discovery | Find users across domains |
| **OAuth Discovery** | Authorization endpoints | Secure third-party access |
| **CalDAV/CardDAV** | Calendar/contacts sync | Works with any compliant app |
| **ActivityPub** | Federation protocol | Interoperability with any system |

**Your domain is your identity.** Just like email (`you@yourdomain.com`), your data lives at URIs you control. Standards-based federation means you can connect to any compliant system—not just ours.

This is how the open web works. We're just applying it to publishing.

### Data Portability

Everything syncs to standard formats you control:

- **Documents** → Your SQLite database
- **Forms** → Your database tables
- **Media** → Your storage
- **Audit logs** → Your compliance records

Export anytime. Migrate anytime. You're never trapped.

---

## Built on Foundation

Publish inherits all [Foundation](/platform/foundation/) capabilities automatically:

| Capability | What It Means |
|------------|---------------|
| **Offline-first** | Work without internet, sync when connected |
| **Universal deployment** | Web, desktop, mobile—one codebase |
| **Self-sovereign** | Your servers, your data, your rules |
| **Real-time sync** | Multiple editors, automatic conflict resolution |

[Learn more about Foundation →](/platform/foundation/)

---

## Get Started

If you're maintaining content in multiple places, watching branding drift, waiting for translations, or manually processing submissions—there's a better way.

[Contact Us →](/contact/)

---

## Part of Something Bigger

Publish is the **2D layer** of a unified dimensional stack—all web-capable, globally accessible.

### The 1D → 2D → 3D Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                    SINGLE SOURCE OF TRUTH                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1D: TEXT (Markdown/DSL)                                        │
│       └── Content definition, structured data, semantic markup   │
│                          │                                       │
│                          ▼                                       │
│   2D: PUBLISH                                                    │
│       └── Web pages, documents, forms, maps, dashboards, kiosks  │
│                          │                                       │
│                          ▼                                       │
│   3D: SPATIAL                                                    │
│       └── CAD, digital twins, robotics, simulation               │
│                                                                  │
│   ALL LAYERS: Web-capable, globally accessible                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Layer | Platform | What It Handles |
|-------|----------|-----------------|
| **1D** | DSL/Markdown | Text content, structured definitions, the source |
| **2D** | Publish | Web, documents, forms, maps, kiosks, emails |
| **3D** | [Spatial](/platform/spatial/) | CAD, digital twins, robotics, simulation |
| **Shared** | [Foundation](/platform/foundation/) | Offline sync, CRDT, NATS, deployment |

**Everything is web-capable.** The 1D text drives 2D outputs. The 2D outputs can embed 3D views. The 3D models can generate 2D documentation. All accessible via browser, all synchronized, all from the same source of truth.

### How They Work Together

For organizations working in 3D—robotics, manufacturing, construction—Publish provides the documentation layer:

- **Technical drawings** generated from 3D models
- **BOMs (Bills of Materials)** extracted from assemblies
- **Work instructions** with step-by-step visuals
- **Compliance forms** linked to CAD revisions
- **Field inspection** with 2D checklists + 3D context

Same source of truth. 2D and 3D outputs stay synchronized. All web-accessible.

[Explore the Spatial Platform →](/platform/spatial/)
