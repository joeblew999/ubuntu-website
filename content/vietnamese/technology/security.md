---
title: "An ninh"
meta_title: "An ninh & Compliance | Ubuntu Software"
description: "Kiến trúc tự chủ giúp đơn giản hóa việc tuân thủ SOC 2, FedRAMP, HIPAA và ISO 27001. Triển khai cách ly mạng, không kết nối về máy chủ, kiểm soát hoàn toàn dữ liệu."
draft: false
---

## Bảo mật từ thiết kế

Chúng tôi không áp dụng các biện pháp bảo mật một cách gượng ép để vượt qua các cuộc kiểm toán. Chúng tôi đã xây dựng một kiến trúc loại bỏ hoàn toàn bề mặt tấn công.

**Kết quả:** Tuân thủ trở nên đơn giản hơn vì có ít hơn các yếu tố cần kiểm toán.

---

## Kiến trúc Tự chủ

| Principle | What It Means |
|-----------|---------------|
| **No Call-Home** | Zero telemetry. No usage tracking. No phone-home to external servers. |
| **Air-Gapped Ready** | Full functionality without internet connectivity. |
| **Data Never Leaves** | Your data stays on your infrastructure. Always. |
| **Single Audit Scope** | No third-party cloud providers to investigate. |
| **No Supply Chain Risk** | Minimal dependencies. Single-binary deployment. |

---

## Tuân thủ trở nên đơn giản hơn

Kiến trúc tự chủ giảm thiểu độ phức tạp của quá trình kiểm toán trên mọi khung công nghệ chính:

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

## Công nghệ nền tảng

An ninh của chúng tôi không chỉ là chính sách—nó được tích hợp vào các lựa chọn công nghệ:

| Component | Security Benefit |
|-----------|------------------|
| **OpenBSD** | Security-focused OS for critical deployments |
| **Go** | Memory-safe, compiled binaries (no runtime vulnerabilities) |
| **NATS JetStream** | mTLS encryption, zero-trust messaging |
| **Single Binary** | No dependency chain. No supply chain attacks. |

---

## Các tùy chọn triển khai

Chọn mức độ bảo mật phù hợp với yêu cầu của bạn:

| Deployment | Use Case | Security Posture |
|------------|----------|------------------|
| **Cloud** | Development, demos | Standard controls |
| **On-Premise** | Enterprise, regulated industries | Enhanced isolation |
| **Air-Gapped** | Government, critical infrastructure | Maximum security |

Tất cả các bản triển khai đều sử dụng cùng một mã nguồn. Không có sự thỏa hiệp về tính năng để đổi lấy mức độ bảo mật cao hơn.

---

## Tuân thủ của nhà cung cấp

NATS JetStream là phụ thuộc bên ngoài duy nhất của chúng tôi. Synadia (công ty đứng sau NATS) duy trì chứng nhận SOC 2 Loại II.

| Vendor | Certification | Compliance Platform |
|--------|---------------|---------------------|
| **Synadia (NATS)** | SOC 2 Type II | Vanta |

- [Thông báo về Chứng nhận Synadia SOC 2 Loại II →](https://www.synadia.com/blog/synadia-soc2-type2-compliant)
- [Xem Báo cáo SOC 2 Loại II của Synadia (PDF) →](/vi/doc/synadia-soc2-type2-report.pdf)

**Một nhà cung cấp. Một cuộc kiểm toán.** Đó là lợi thế của việc giảm thiểu sự phụ thuộc.

---

## Sẵn sàng đơn giản hóa việc tuân thủ?

Kiến trúc của chúng tôi đảm nhận phần việc khó khăn. Các cuộc kiểm toán của bạn sẽ trở nên dễ dàng hơn.

[Bắt đầu →](/vi/get-started/)

---

## Báo cáo lỗ hổng bảo mật

Phát hiện vấn đề bảo mật? Chúng tôi coi trọng vấn đề bảo mật và đánh giá cao việc tiết lộ có trách nhiệm.

[Chính sách công bố lỗ hổng bảo mật →](/vi/security/)

---

## Có câu hỏi nào không?

Cần thảo luận về các yêu cầu tuân thủ cụ thể?

[Liên hệ với chúng tôi →](/vi/contact/)
