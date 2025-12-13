---
title: "The Economics of Everywhere Publishing"
meta_title: "The Economics of Everywhere Publishing | Ubuntu Software"
description: "A phased approach to single-source publishing: start with authoring everywhere, add connectivity, then integrate your legacy database. No rip-and-replace required."
date: 2024-12-13T05:00:00Z
image: "/images/blog/economics-publishing.svg"
categories: ["Publish", "Strategy"]
author: "Gerard Webb"
tags: ["publishing", "raspberry-pi", "esim", "database-integration", "strategy"]
draft: false
---

Here's a question that keeps organizations stuck: **How do we modernize without disrupting everything?**

The answer isn't a big-bang migration. It's a phased approach that delivers value at each step while building toward full integration.

I've been working on this problem for years. Let me walk you through how the economics actually work.

## Objective 1: Publishing System Everywhere

**This is the foundation. Get it into as many hands as possible.**

The first objective is deceptively simple: let people author content once and publish it everywhere. Web, PDF, kiosk—same source, multiple outputs.

Why start here?

**It's immediately useful.** No database integration required. No legacy system changes. Just a better way to create and publish content.

**The pain is universal.** Every organization has content scattered across Word docs, PDFs, websites, and various systems. None of them match. Updates require touching multiple places.

**Low barrier to entry.** Deploy a Raspberry Pi running the publishing system. Authors create content. Output flows to displays, web, print. Done.

```
Author creates content (Markdown)
         ↓
    Publish Engine
         ↓
┌────────┼────────┐
↓        ↓        ↓
Website  PDF    Kiosk Display
```

**The PDF round-tripping is particularly powerful.** A government office generates a form as PDF. Citizen fills it out. Scans it back. OCR extracts the data and round-trips it into the system.

This alone solves a massive problem for any government entity: paper forms that require manual data entry.

The kiosk use case extends this further. Same form, displayed on a touchscreen. Citizen fills it out directly. Data flows into the system without paper or scanning.

**Three input methods, same destination:**
- Web form → database
- Printed PDF, filled, scanned → database
- Kiosk touchscreen → database

This is objective one. Get the publishing system running. Let people author. Let content flow to every screen.

## Objective 2: Always-On Connectivity

**The Raspberry Pi needs a backhaul channel.**

Once you have publishing working, the next question is reliability. What happens when the internet goes down? What about remote locations with spotty connectivity?

The answer is eSIM.

Each Raspberry Pi gets an eSIM hardware attachment. This provides a cellular backhaul that's independent of local WiFi or ethernet. The device always has a connection path.

**Why this matters:**

- **Remote kiosks**: A government service point in a rural area doesn't have to depend on local ISP reliability
- **Fallback connectivity**: If the primary connection fails, eSIM takes over automatically
- **Secure channel**: Cellular connection can be configured for encrypted backhauling

The goal is programmatic provisioning—deploy a device, it connects automatically.

The technical architecture: Raspberry Pi + eSIM module = always-connected publishing node.

## Objective 3: Legacy Database Integration

**Now we connect to your real data.**

Once you have publishing working and connectivity solved, the obvious question emerges: "This is great, but our real data is in Oracle/SQL Server/PostgreSQL/whatever."

This is where the architecture gets interesting.

The publishing system uses SQLite as its origin database. This isn't a limitation—it's a design choice. SQLite is:
- Fast (it's local)
- Reliable (battle-tested)
- Portable (runs on the Pi)
- Zero configuration

But here's the key: **SQLite syncs to your legacy database.**

```
Publishing System (SQLite)
         ↓
    Sync Engine
         ↓
Legacy Database (Oracle, SQL Server, PostgreSQL, MySQL, etc.)
```

We have connectors for everything from Oracle to SQLite (yes, SQLite to SQLite for distributed setups). The sync is bidirectional where needed.

**What this means in practice:**

Your existing database doesn't change. Your existing applications keep running. The publishing system sits alongside, syncing data in and out.

- Product catalog in Oracle? Sync it to the publishing system for web/PDF/kiosk output
- Customer submissions from kiosks? Sync them back to your CRM
- Inventory updates in your ERP? Flow to digital signage automatically

**No rip and replace. No migration project. Just a sync layer.**

## Objective 4: Data Transformation

**Getting data back into legacy systems, structured correctly.**

The final piece is transforming data from the publishing system back into the format your legacy database expects.

This is partly working now. The challenge is that legacy systems have opinions about data structures. They expect specific field names, specific formats, specific relationships.

The transformation engine handles:
- Field mapping (our field names → your field names)
- Format conversion (dates, numbers, enumerations)
- Relationship resolution (foreign keys, lookups)
- Validation rules (your business logic)

When this is complete, the loop closes:

```
Legacy Database → Sync → Publishing System → Outputs
                                ↓
                          User Input
                                ↓
Publishing System → Transform → Sync → Legacy Database
```

**Data flows in both directions. Single source of truth. No manual data entry.**

## The Phased Economics

Here's why this approach works economically:

### Phase 1: Immediate Value
- Deploy publishing system
- Authors create content
- Multiple outputs generated automatically
- **ROI: Reduced content duplication, fewer errors, time savings**

### Phase 2: Infrastructure Investment
- Add eSIM connectivity
- Deploy to remote/unreliable locations
- Ensure always-on operation
- **ROI: Expand coverage, improve reliability**

### Phase 3: Integration Value
- Connect to legacy databases
- Eliminate manual data sync
- Single source of truth
- **ROI: Eliminate duplicate data entry, reduce errors, real-time accuracy**

### Phase 4: Transformation Value
- Close the loop
- User input flows to legacy systems
- Full round-trip automation
- **ROI: Complete elimination of manual data handling**

Each phase delivers value. Each phase funds the next. No massive upfront investment required.

## The Real Goal

Get the publishing system into many hands. Let people author and publish. Prove the value at the content layer.

Then add connectivity. Then add integration. Then add transformation.

**Start simple. Grow capability. No big bang.**

This is how you modernize organizations that can't afford to stop running while you rebuild everything.

---

*Ready to start with objective one? [Get in touch →]({{< relref "/contact" >}})*

---

*Part of our Publish platform. [Learn more about single-source publishing →]({{< relref "/platform/publish" >}})*
