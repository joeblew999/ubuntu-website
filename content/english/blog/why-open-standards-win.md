---
title: "Why Open Standards Win"
meta_title: "Why Open Standards Win | Ubuntu Software"
description: "STEP, IFC, and the case against proprietary lock-in for 3D design and engineering data."
date: 2024-09-20T05:00:00Z
image: "/images/blog/open-standards.svg"
categories: ["Industry", "Standards"]
author: "Gerard Webb"
tags: ["open-standards", "step", "ifc", "interoperability", "cad"]
draft: false
---

Every decade or so, an industry learns the same lesson: proprietary lock-in doesn't scale.

The web learned it. Enterprise software learned it. Cloud computing learned it.

Now it's time for 3D design and engineering to learn it.

## The Current State of Affairs

Try this experiment: Take a 3D model from one major CAD system and open it in another.

What you'll find:
- **Geometry loss**: Features don't translate. Constraints disappear.
- **Metadata gone**: All the engineering information—materials, tolerances, assembly relationships—lost.
- **Manual rework**: Someone spends hours recreating what already existed.

This isn't a bug. It's a business model.

**Vendors profit from lock-in. Users suffer from it.**

## The Cost of Proprietary Formats

### Collaboration Friction

When your robot vendor uses System A, your facility designer uses System B, and your simulation team uses System C, every handoff is a translation exercise.

Information degrades with each conversion. Engineers spend time on file wrangling instead of engineering.

### Vendor Dependency

Your entire design history lives in a format only one vendor controls. They set the prices. They set the upgrade timeline. They decide when features get deprecated.

Your engineering IP is held hostage.

### Innovation Barriers

Want to build AI tools on top of your design data? Good luck accessing it through proprietary APIs that change with every version.

Want to integrate with the latest physics simulation? Better hope your CAD vendor has a partnership.

Innovation dies at format boundaries.

## The Open Alternative

### STEP: The Geometry Standard

ISO 10303 (STEP) has existed since the 1990s. It's boring. It works.

STEP captures:
- 3D geometry with full precision
- Assembly structures and relationships
- Product manufacturing information (PMI)
- Material properties

It's not perfect. But it's universal.

### IFC: Buildings that Talk

Industry Foundation Classes (IFC) does for buildings what STEP does for products.

Every wall, door, space, and system—defined in an open format that any software can read and write.

BIM interoperability isn't a dream. IFC makes it possible.

### The Emerging Stack

Modern open standards go beyond static geometry:

- **glTF**: Lightweight 3D for visualization and AR/VR
- **USD**: Scene description for simulation and rendering
- **SDF**: Robot and environment definition
- **URDF**: Robot description format

An ecosystem is forming. Tools built on open foundations can participate.

## What Open Standards Enable

### Real Competition

When your data isn't locked in, you can choose tools based on capability, not captivity.

Vendors compete on features, not on how difficult they make migration.

### Ecosystem Innovation

Open formats enable an ecosystem of specialized tools:
- AI assistants that work across platforms
- Simulation engines that take any geometry
- Collaboration tools that don't require everyone to own the same license

### Future-Proofing

Standards bodies move slowly. That's a feature.

A STEP file from 2000 still opens today. Will your proprietary format from 2020 open in 2040?

## The Hybrid Reality

Let's be practical: pure open-standards workflows don't exist yet.

The real strategy is:
1. **Native formats for active work**: Use the best tool for each job
2. **Open formats for exchange**: Standard formats at every handoff point
3. **Open formats for archival**: Long-term storage in formats you control

This isn't idealism. It's risk management.

## The Industry Shift

The momentum is building:

**Government mandates**: More agencies requiring open formats for procurement and archival.

**Industry consortia**: Organizations like buildingSMART pushing IFC adoption.

**AI requirements**: Machine learning needs training data that isn't locked away.

**Cloud collaboration**: Real-time collaboration platforms choosing open foundations.

The vendors that embrace open standards will win. The ones that fight it will be routed around.

## Making the Transition

If you're starting fresh, build on open foundations:
- Choose tools with strong open format support
- Require standards compliance in vendor contracts
- Establish open-format checkpoints in your workflows
- Archive in formats you control, not formats that control you

If you're migrating, start at the boundaries:
- New integrations use open formats
- New projects pilot open workflows
- Gradual migration as tools and workflows mature

## The Bigger Picture

Open standards aren't about technology. They're about power.

Who controls your engineering data? Who decides what tools you can use? Who owns your design history?

Proprietary formats answer: the vendor.

Open standards answer: you do.

That's why open standards win. Not because they're technically superior (though they often are). Because they align incentives correctly.

Your data. Your choice. Your future.

---

*We built our platform on STEP, IFC, and open APIs. [See how it works →](/platform)*
