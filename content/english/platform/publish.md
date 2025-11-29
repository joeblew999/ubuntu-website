---
title: "Publish"
meta_title: "Content Management Platform | Ubuntu Software"
description: "Single-source CMS for multi-channel publishing. One DSL for text, graphics, forms, emails. Auto-translated outputs. Perfect branding by architecture."
image: "/images/publish.svg"
draft: false
---

## One Language. Every Output. Every Language.

A single-source content management system. Define everything in one DSL. Compose components that include components. Output to any channel. Auto-translate to any language. Capture data back.

Perfect branding. Not by discipline. By architecture.

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

All links route through [your gateway](/platform/foundation/#wellknown-gateway). Publish TO Big Tech platforms. Never be locked IN.

---

## Built on Foundation

Publish inherits all [Foundation](/platform/foundation/) capabilities automatically:

| Capability | What It Means |
|------------|---------------|
| **Offline-first** | Work without internet, sync when connected |
| **Universal deployment** | Web, desktop, mobile—one codebase |
| **Self-sovereign** | Your servers, your data, your rules |
| **Real-time sync** | Multiple editors, automatic conflict resolution |
| **Wellknown Gateway** | Publish TO Big Tech, never locked IN |

[Learn more about Foundation →](/platform/foundation/)

---

## Get Started

If you're maintaining content in multiple places, watching branding drift, waiting for translations, or manually processing submissions—there's a better way.

[Contact Us →](/contact/)

---

## Part of Something Bigger

Publish is the 2D layer of the Ubuntu Software platform.

For organizations working in 3D—robotics, manufacturing, construction, digital twins—Publish provides the documentation layer. Technical drawings, BOMs, work instructions, compliance forms—all from single source, all aligned with the 3D model.

[Explore the Spatial Platform →](/platform/spatial/)
