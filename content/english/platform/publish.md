---
title: "Publish"
meta_title: "Single-Source Publishing | Ubuntu Software"
description: "Single-source publishing. Markdown to web, PDF, and forms. OCR capture back to your database. Full circle."
image: "/images/publish.svg"
draft: false
---

## Write Once. Output Everywhere. Capture Back.

Markdown in. Web pages, PDFs, and forms out—all perfectly aligned. Capture data from digital or paper submissions. Back to your database. Full circle.

No more maintaining multiple versions. No more re-keying paper forms. No more content drift.

---

## The Problem

Every organization publishes information in multiple formats.

The website says one thing. The PDF says another. The form has different fields. Paper submissions get typed in manually. Updates mean changing everything in three places.

**Content drifts. Errors multiply. Staff waste hours on manual data entry.**

---

## The Solution

### Single Source of Truth

Everything starts as Markdown and a simple DSL.
```
# Application Form

Name: [_______________]{field: name, required: true}
Email: [_______________]{field: email, type: email}
Date: [__/__/____]{field: date, type: date}
```

Human-readable. Version-controlled. The one source that drives everything.

---

### Aligned Outputs

From that single source, generate:

| Output | Use |
|--------|-----|
| Web page | Online viewing |
| PDF document | Print, email, archive |
| Web form | Digital submissions |
| PDF form | Fillable digital or printable |

**All outputs are 100% aligned.** Same content. Same field positions. Same structure.

Change the source. Everything updates.

---

### Full-Circle Capture

Here's the breakthrough:

**Web form** → Data flows to your database

**PDF form (filled digitally)** → Data extracts to your database

**PDF form (printed, handwritten, scanned)** → OCR captures to your database

Because we generated the form, OCR knows exactly where every field is. No training. No guessing. Precise extraction.

**Paper or digital. Same data. Same destination.**

---

### Any Database

We don't replace your backend. We connect to it.

- Your existing database
- Your schema
- Read data to populate documents
- Write captured data back

**Your infrastructure. Our publishing and capture layer.**

---

## How It Works

### Step 1: Write

Create your content in Markdown with our DSL for fields and structure.

Plain text files. Version control friendly. Human readable.

### Step 2: Publish

Generate all outputs from the source:

- Static web pages
- PDF documents
- Interactive web forms
- Fillable PDF forms

Deploy to your website. Print PDFs. Distribute however you need.

### Step 3: Capture

Receive submissions through any channel:

- Web form submissions
- Digitally filled PDFs
- Scanned paper forms

All data validates, transforms, and flows to your database.

### Step 4: Close the Loop

Data in your database can populate new documents. Generate personalized outputs. Drive workflows.

**Source → Output → Capture → Database → Source**

The circle is complete.

---

## Use Cases

### Government & Public Sector

Citizens need options. Some file online. Some mail paper forms. Both need to work.

- Accessible web and PDF from same source
- Paper submissions processed automatically
- Compliance with accessibility standards
- Archival formats supported

**Serve every citizen. Process every submission.**

---

### Healthcare

Patient intake. Consent forms. Medical history. Insurance paperwork.

- Waiting room paper forms, patient portal digital forms—same source
- OCR eliminates transcription errors
- Audit trail for compliance
- Integration with existing EMR/EHR systems

**Less administration. Fewer errors. Better care.**

---

### Financial Services

Applications. Claims. Disclosures. Account forms.

- Regulatory-compliant PDF output
- Digital and paper channels unified
- Data capture without manual re-entry
- Version control for audit requirements

**Compliance without complexity.**

---

### Education

Enrollment. Registration. Permission slips. Transcripts.

- Parents choose paper or digital
- Administrative staff freed from data entry
- Records stay synchronized
- Bulk document generation from student data

**Modern administration. Traditional options.**

---

### Insurance

Policy applications. Claims forms. Assessments. Renewals.

- Field agents work on paper
- Customers submit online
- Both captured identically
- Back-office integration seamless

**Every channel. One system.**

---

### Enterprise

Internal forms. HR paperwork. Procedures. Checklists.

- Intranet and print from same source
- Field operations work offline on paper
- Data returns to enterprise systems
- Documentation always current

**Operations that actually work.**

---

## What You Get

### For Content Teams

- Write in Markdown—no design tools required
- One source to maintain, not three
- Updates propagate automatically
- Version control built in

### For Operations

- Paper and digital submissions unified
- No more manual data entry
- Error rates plummet
- Processing time drops

### For IT

- Connects to existing databases
- Standard formats (PDF, HTML)
- No proprietary lock-in
- Self-hosted or cloud

### For Compliance

- Single source of truth
- Complete audit trail
- Version history preserved
- Archival formats supported

---

## Architecture

Clean. Simple. Yours.

| Layer | What It Does |
|-------|--------------|
| Source files | Markdown + DSL, version controlled |
| Parser | Extracts content, structure, field definitions |
| Renderers | Generates web, PDF, form outputs |
| Field map | Knows exact position of every field |
| OCR engine | Precise extraction using field map |
| Connectors | Read/write to your database |

**No proprietary formats. Export everything. Own your content.**

---

## Integration

Works with what you have:

- **Databases:** PostgreSQL, MySQL, SQLite, or any SQL database
- **Storage:** Local filesystem, S3, or any object storage
- **Auth:** Your existing identity provider
- **Deployment:** Your servers, your cloud, your choice

**We adapt to your infrastructure. Not the other way around.**

---

## The ROI

**Before:**
- Staff hours re-keying paper submissions
- Errors from manual data entry
- Time updating multiple document versions
- Content drift between web and PDF

**After:**
- Paper submissions processed automatically
- Near-zero data entry errors
- Single source updates everything
- Perfect alignment across all outputs

**The math is simple. The savings are real.**

---

## Get Started

If you're maintaining forms in multiple places, manually processing paper submissions, or watching content drift across formats—there's a better way.

[Contact Us →](/contact/)

---

## Part of Something Bigger

Publish is the 2D foundation of the Ubuntu Software platform.

For organizations working in 3D—robotics, manufacturing, construction, digital twins—Publish provides the documentation layer. Technical drawings, BOMs, work instructions, compliance forms—all from single source, all aligned with the 3D model.

[Explore the Spatial Platform →](/platform/spatial/)
