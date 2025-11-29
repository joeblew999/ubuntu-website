---
title: "Security"
meta_title: "Security & Compliance | Ubuntu Software"
description: "Self-sovereign architecture that simplifies SOC 2, FedRAMP, HIPAA, and ISO 27001 compliance. Air-gapped deployment, no call-home, complete data control."
draft: false
---

## Security by Design

We don't bolt on security controls to pass audits. We built an architecture that eliminates the attack surface.

**The result:** Compliance becomes simpler because there's less to audit.

---

## Self-Sovereign Architecture

| Principle | What It Means |
|-----------|---------------|
| **No Call-Home** | Zero telemetry. No usage tracking. No phone-home to external servers. |
| **Air-Gapped Ready** | Full functionality without internet connectivity. |
| **Data Never Leaves** | Your data stays on your infrastructure. Always. |
| **Single Audit Scope** | No third-party cloud providers to investigate. |
| **No Supply Chain Risk** | Minimal dependencies. Single-binary deployment. |

---

## Compliance Made Simpler

Self-sovereign architecture reduces audit complexity across every major framework:

| Framework | Who Needs It | Why Self-Sovereign Helps |
|-----------|--------------|-------------------------|
| **SOC 2 Type II** | Enterprise B2B | Reduced attack surface. Simpler scope. Faster audits. |
| **FedRAMP** | U.S. Government | Agencies prefer isolated deployments. Complete data sovereignty. |
| **HIPAA** | Healthcare | No external BAAs needed. ePHI never leaves your systems. |
| **PCI DSS** | Financial Services | Encryption keys under your control. No cloud intermediaries. |
| **ISO 27001** | International | Fewer external dependencies to document and monitor. |
| **IEC 62443** | Manufacturing | Meets air-gapped security level requirements for industrial control. |
| **GDPR/CCPA** | Data Privacy | No cross-border transfers. Supports data residency requirements. |

---

## Technology Stack

Our security isn't just policy—it's built into the technology choices:

| Component | Security Benefit |
|-----------|------------------|
| **OpenBSD** | Security-focused OS for critical deployments |
| **Go** | Memory-safe, compiled binaries (no runtime vulnerabilities) |
| **NATS JetStream** | mTLS encryption, zero-trust messaging |
| **Single Binary** | No dependency chain. No supply chain attacks. |

---

## Deployment Options

Choose the security level that matches your requirements:

| Deployment | Use Case | Security Posture |
|------------|----------|------------------|
| **Cloud** | Development, demos | Standard controls |
| **On-Premise** | Enterprise, regulated industries | Enhanced isolation |
| **Air-Gapped** | Government, critical infrastructure | Maximum security |

All deployments use the same codebase. No feature compromises for higher security.

---

## Vendor Compliance

NATS JetStream is our only external dependency. Synadia (the company behind NATS) maintains SOC 2 Type I certification.

| Vendor | Certification | Compliance Platform |
|--------|---------------|---------------------|
| **Synadia (NATS)** | SOC 2 Type I | Vanta |

- [View Synadia SOC 2 Report (PDF) →](/doc/synadia-soc2-type1-report.pdf)
- [Synadia SOC 2 Type II Announcement →](https://www.synadia.com/blog/synadia-soc2-type2-compliant)

**One vendor. One audit.** That's the advantage of minimal dependencies.

---

## Ready to Simplify Compliance?

Our architecture does the hard work. Your audits get easier.

[Get Started →](/get-started/)

---

## Questions?

Need to discuss specific compliance requirements?

[Contact Us →](/contact/)
