---
title: "DJI Is Now on the FCC Covered List — What Actually Happened"
meta_title: "DJI FCC Ban: What Happened and What It Means | Ubuntu Software"
description: "The December 2025 deadline passed. DJI is now on the FCC Covered List. Here's how it happened, what it means, and what comes next."
date: 2025-12-24T06:00:00Z
image: "/images/blog/dji-fcc-ban.svg"
categories: ["Industry", "Drones"]
author: "Gerard Webb"
tags: ["drones", "dji", "fcc", "regulation", "open-source"]
draft: false
---

On December 22, 2025—one day before the deadline—the FCC announced the addition of "uncrewed aircraft systems (UAS) and UAS critical components produced in a foreign country" to the FCC Covered List.

**The ban is now in effect.**

## How We Got Here

Section 1709 of the FY2025 NDAA required a federal security audit of DJI drones. The deadline was December 23, 2025. If no audit occurred, DJI would be automatically added to the FCC Covered List.

Here's what happened:

**DJI requested the audit.** They sent letters in March, June, and December 2025 asking federal agencies to examine their products.

**No agency responded.** No investigation was started. No evidence was gathered. No findings were made.

**The administration bypassed the process.** On December 21, 2025, the Executive Branch issued a "national security determination" instead of completing the mandated investigation.

**The deadline passed.** On December 22, the FCC made it official.

## The "Trap Door" Design

Congress deliberately built Section 1709 with two trap doors:

1. **No assigned agency.** The law didn't specify which agency should conduct the audit.
2. **Automatic ban.** If no audit occurred, the ban was automatic.

The outcome was predetermined. The process existed to create the appearance of due process while ensuring the result.

## What This Actually Means

### For Existing DJI Owners

Your drones still work. You can still fly them under FAA Part 107. Nothing changes operationally—for now.

### For New Purchases

New DJI models and covered RF components cannot obtain the FCC approvals required for legal import, marketing, or sale in the US. This effectively blocks:

- New drone models from entering the US market
- New firmware that requires FCC certification
- New accessories with RF components

### The Exception Path

The FCC notes that exceptions may be granted if the Department of Defense or Department of Homeland Security determines that a specific UAS or component does not pose a risk.

Whether any exceptions will actually be granted remains to be seen.

## DJI's Response

DJI called the FCC's action "disappointing" and stated:

> "Concerns about DJI's data security have not been grounded in evidence and instead reflect protectionism, contrary to the principles of an open market."

They confirmed existing products are not affected and that they remain committed to the US market.

## The Scale of Impact

This isn't a niche issue:

- **87%** of public safety drones in the United States are DJI
- **80%+** of state and local public safety agencies rely on DJI for wildfire response, search-and-rescue, and critical operations

The timing—right before Christmas, with wildfires and natural disasters ongoing—underscores the gap between policy and operational reality.

## What Happens Next

### Short Term (2025-2026)

- Existing DJI operators continue normal operations
- No new DJI models enter the US market
- Secondary market prices for existing DJI hardware increase
- Firmware updates become uncertain

### Medium Term (2026-2027)

- Supply chain for DJI parts and accessories tightens
- Fleet operators face decisions about replacement platforms
- Open source alternatives see increased adoption
- Other Chinese manufacturers face similar scrutiny

### Long Term

- US drone market fragments between domestic manufacturers and open source
- Public safety agencies transition to alternative platforms
- Enterprise operators invest in vendor-independent solutions

## The Open Source Alternative

We've been preparing for this. Our [DJI migration guide]({{< relref "/fleet/dji-migration" >}}) covers the transition path to open source platforms like PX4 and ArduPilot.

The case for open source is now stronger than ever:

| DJI (Post-Ban) | Open Source |
|----------------|-------------|
| No new models | Continuous innovation |
| Uncertain firmware | Community-maintained |
| Single vendor | Multiple manufacturers |
| Regulatory risk | Standards-based |
| Cloud dependency | Self-hosted option |

If you're operating a DJI fleet, the question isn't whether to transition—it's when and how.

## For Public Safety and Enterprise

The operators with the most at stake are public safety agencies and enterprise fleets. They need:

- Proven alternatives that match DJI capabilities
- Transition plans that don't disrupt operations
- Training and support for new platforms
- Fleet management infrastructure

We provide the infrastructure layer—the software that makes open source drone fleets actually work. Our [fleet architecture]({{< relref "/fleet" >}}) handles the coordination, telemetry, and management that enterprise operations require.

## The Bottom Line

The DJI ban is now reality, not speculation. The process was flawed, but the outcome is final.

Operators have two choices:

1. **Wait and see.** Continue with existing DJI hardware while supplies last. Hope for exceptions that may never come.

2. **Transition now.** Move to open source platforms while you can plan methodically instead of scrambling.

We recommend option two.

---

*Need help planning your transition? [Contact us →](/contact?subject=dji-migration)*

*Already have a DJI fleet? [See the migration guide →]({{< relref "/fleet/dji-migration" >}})*
