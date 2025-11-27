---
title: "One Source, Every Screen"
meta_title: "One Source, Every Screen | Ubuntu Software"
description: "From restaurant menus to government kiosks: how single-source publishing transforms physical displays alongside web and print."
date: 2024-11-27T05:00:00Z
image: "/images/blog/kiosk-system.svg"
categories: ["Publish", "Industry"]
author: "Gerard Webb"
tags: ["kiosk", "signage", "raspberry-pi", "single-source", "real-time"]
draft: false
---

A restaurant updates their menu. What happens next?

Someone edits the website. Someone else updates the PDF for takeout menus. A third person changes the point-of-sale system. Eventually, someone remembers the digital menu boards and manually updates those too.

Four systems. Four opportunities for errors. Four places where prices can drift.

**This is absurd. And it's everywhere.**

## The Display Fragmentation Problem

Walk into any organization with public-facing information:

**Restaurants**: Menu boards, table QR codes, website, delivery apps, printed menus
**Government offices**: Queue displays, wayfinding signs, forms counter, website
**Healthcare**: Waiting room displays, check-in kiosks, patient portals, printed materials
**Retail**: Price displays, promotional signage, inventory status, e-commerce

Each display typically runs its own system. Each requires separate updates. Each can show different information.

**Information drifts. Customers get confused. Staff waste time on manual sync.**

## The Single-Source Solution

What if every display—physical and digital—pulled from the same source?

```
Your Database (single source of truth)
         ↓
    Publish Engine
         ↓
┌────────┼────────┐
↓        ↓        ↓
Website  PDF    Kiosk Display
         ↓        ↓
      Print    Menu Board
```

Change the price once. Every display updates.

Add a new item once. It appears everywhere.

Mark something out of stock once. Customers see it immediately—on the website, the menu board, and the printed QR code menu.

## Real-Time, Offline-Capable Kiosks

Here's where it gets interesting.

Traditional digital signage requires constant connectivity. Network goes down? Displays go blank or show stale data.

Our approach uses Raspberry Pi devices running the same offline-first architecture as everything else:

- **Local data replica**: Each kiosk has a complete copy of relevant data
- **Automerge sync**: Updates flow automatically when connected
- **Graceful degradation**: Network outage? Display keeps working with local data
- **Queue for changes**: Updates made at the kiosk sync back when connectivity returns

**A menu board that works during internet outages. A check-in kiosk that queues submissions offline.**

## Use Cases That Make Sense

### Restaurant Systems

One source drives:
- Digital menu boards (real-time pricing, out-of-stock items greyed)
- QR code menus (same content, mobile-optimized)
- Website menu pages
- PDF menus for print
- Kitchen display integration

Change the price of a burger. Everywhere updates. The menu board by the counter, the PDF someone prints for catering, the website—all aligned.

### Government Service Centers

One source drives:
- Queue management displays ("Now serving B-47")
- Service directory kiosks
- Forms available for download (PDF) or completion (web/kiosk)
- Wayfinding signage
- Website service listings

Update office hours. Every display reflects the change. The wayfinding kiosk, the website, the printed handouts—all correct.

### Healthcare Waiting Rooms

One source drives:
- Check-in kiosks (form completion, queue joining)
- Waiting time displays
- Patient information screens
- Printed intake forms (same fields, aligned layout)

Patient fills form on kiosk → data flows to your EMR → no manual transcription.

### Retail Environments

One source drives:
- Price displays (electronic shelf labels)
- Promotional signage
- Inventory status screens
- Website product listings
- Printed catalogs

Mark item as "clearance" in inventory system. Price displays update. Website updates. Printed weekly ad shows clearance pricing.

## The Hardware Story

Why Raspberry Pi?

**Cost**: Under $100 per display node, including the Pi and basic enclosure
**Reliability**: No moving parts, runs Linux, runs for years
**Flexibility**: HDMI output drives any display size
**Offline capable**: Local compute and storage for true edge operation
**Updatable**: Remote management for software updates

A restaurant can deploy menu boards for a fraction of traditional digital signage costs. A government office can add kiosks without enterprise signage contracts.

## The Technical Architecture

Each kiosk runs:

| Component | Purpose |
|-----------|---------|
| **Local database** | SQLite replica of relevant data |
| **Sync engine** | Automerge CRDT for conflict-free updates |
| **Display renderer** | Web technologies (HTML/CSS/JS) for flexible layouts |
| **Input handlers** | Touch, keyboard, barcode scanner, card reader |
| **Offline queue** | Submissions stored locally until sync |

The same codebase that runs your web forms runs on the kiosk. Same validation. Same field layouts. Same data destination.

## Beyond Displays: Input Devices Too

Kiosks aren't just output devices. They're also input points.

**Form completion**: Customers complete forms on kiosk touchscreen
**Document scanning**: Kiosk scans paper documents, OCR extracts data
**Payment processing**: Integrated card readers
**ID verification**: Camera for document scanning
**Signature capture**: Touch signature pads

All captured data flows through the same system as web submissions.

A government form submitted via:
- Website web form → database
- PDF printed, filled, scanned at kiosk → database
- Kiosk touchscreen form → database

**Same data. Same destination. Different input method.**

## The Alignment Advantage

Here's what "aligned" really means:

A customer sees a form on your website. Same customer visits your office and sees the same form on a kiosk. Same fields, same layout, same questions.

They start filling on the kiosk, run out of time, scan a QR code to continue on their phone. Same form, same progress (synced via their session).

They complete at home, hit submit. Data flows to your database. No re-keying, no transcription errors, no misaligned field names.

**This is what single-source publishing enables when extended to physical displays.**

## The Business Case

**Before single-source displays:**
- Staff time updating multiple systems
- Errors from inconsistent information
- Customer confusion ("the sign says one price, the website says another")
- High cost of traditional digital signage systems
- Displays that fail when internet drops

**After single-source displays:**
- Update once, propagate everywhere
- Guaranteed consistency across all touchpoints
- Reduced customer complaints
- Commodity hardware at commodity prices
- Displays that work offline

**The math is straightforward: fewer errors, less staff time, better customer experience.**

## Getting Started

We're building this now. The same Publish engine that generates web pages, PDFs, and forms will generate kiosk displays.

Same Markdown source. Same DSL for fields. Same database connection.

The display is just another output format.

---

*Interested in early access to kiosk capabilities? [Contact us →](/contact)*

---

*Part of our Publish platform. [Learn more about single-source publishing →](/platform/publish)*
