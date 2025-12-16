---
title: "Công nghệ"
meta_title: "Công nghệ Stack | Ubuntu Software"
description: "Built on Đi and NATS JetStream - our technology choices for performance, reliability, and simplicity."
image: "/images/spatial.svg"
draft: false
---

## Công nghệ của chúng tôi

Chúng tôi xây dựng trên các công nghệ được lựa chọn dựa trên hiệu suất, độ tin cậy và khả năng bảo trì lâu dài.

---

## Go

**Tại sao nên đi:**

- **Hiệu suất** — Được biên dịch, có kiểu tĩnh, chi phí thời gian chạy tối thiểu
- **Đơn giản** — Một cách để thực hiện, dễ đọc theo mặc định
- **Đồng thời** — Goroutines và kênh được tích hợp sẵn trong ngôn ngữ
- **Triển khai** — Tệp nhị phân duy nhất, không có phụ thuộc, biên dịch chéo
- **Hệ sinh thái** — Thư viện tiêu chuẩn mạnh mẽ, công cụ hỗ trợ xuất sắc

Go là ngôn ngữ lập trình chính của chúng tôi trên toàn bộ hệ thống—từ các dịch vụ phía server đến các công cụ dòng lệnh (CLI) cho đến tích hợp robotics.

---

## NATS JetStream

**Tại sao chọn NATS JetStream:**

- **Kiên trì** — Lưu trữ tin nhắn bền vững với khả năng phát lại
- **Giao hàng chính xác một lần* — Xử lý tin nhắn đáng tin cậy
- **Nhẹ nhàng** — Tệp nhị phân duy nhất, tiêu thụ tài nguyên tối thiểu
- **Mở rộng quy mô** — Tích hợp sẵn khả năng tạo cụm và mở rộng theo chiều ngang
- **Thời gian thực** — Độ trễ dưới một mili giây cho mô hình pub/sub

NATS JetStream cung cấp nền tảng cho kiến trúc hướng sự kiện của chúng tôi—kết nối các công cụ thiết kế, mô phỏng, bản sao kỹ thuật số và các thiết bị vật lý.

---

## Nguyên tắc kiến trúc

| Principle | Implementation |
|-----------|----------------|
| **Offline-first** | Local-first data, sync when connected |
| **Event-driven** | NATS JetStream for all inter-service communication |
| **Open standards** | STEP, IFC, no proprietary formats |
| **Hardware-agnostic** | Abstraction layers for portability |
| **Self-sovereign** | Deploy anywhere—cloud, on-prem, air-gapped |

---

## Các lĩnh vực công nghệ

### Robotics

Hệ thống robotics của chúng tôi được xây dựng trên Viam RDK—một bộ công cụ phát triển robotics mã nguồn mở cung cấp trừu tượng hóa phần cứng, lập kế hoạch chuyển động và thị giác máy tính.

[Hệ thống robotics →](/technology/robotics/)

### Linux & Đa nền tảng

Hàng chục năm kinh nghiệm với Linux. Ứng dụng đa nền tảng cho Windows, Mac, Linux, iOS và Android. Khung giao diện người dùng (GUI) riêng sau nhiều năm làm việc với Qt, Flutter và Electron.

[Linux & Đa nền tảng →](/technology/linux/)

### An ninh & Tuân thủ

Kiến trúc tự chủ giúp đơn giản hóa việc tuân thủ SOC 2, FedRAMP, HIPAA và ISO 27001. Triển khai cách ly mạng, không kết nối về máy chủ, kiểm soát hoàn toàn dữ liệu.

[An ninh & Tuân thủ →](/technology/security/)

---

## Nền tảng triển khai

Chúng tôi triển khai ở mọi nơi:

| Platform | Use Case |
|----------|----------|
| **Linux** | Server, desktop, embedded |
| **OpenBSD** | Security-critical systems |
| **Windows** | Enterprise desktop |
| **macOS** | Design and development |
| **iOS/Android** | Mobile applications |
| **Docker/Kubernetes** | Container orchestration |

---

## Quỹ Nguồn Mở

Chúng tôi phát triển và đóng góp cho mã nguồn mở:

| Component | Technology |
|-----------|------------|
| Language | Go |
| Messaging | NATS JetStream |
| Collaboration | Automerge (CRDT) |
| 3D Formats | STEP, IFC |
| Robotics | Viam RDK |
| Vision | Intel RealSense, YOLOv8 |

---

## Tìm hiểu thêm

- [Hệ thống robotics →](/technology/robotics/)
- [Linux & Đa nền tảng →](/technology/linux/)
- [An ninh & Tuân thủ →](/technology/security/)
- [Nền tảng không gian →](/platform/spatial/)
- [Nền tảng →](/platform/foundation/)
- [Liên hệ với chúng tôi →](/contact/)
