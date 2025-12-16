---
title: "Cảm biến"
meta_title: "Cảm biến & Perception Platform | Ubuntu Software"
description: "Tích hợp đa cảm biến cho trí tuệ không gian — LiDAR, camera và cảm biến công nghiệp được tích hợp thông qua một đại lý biên duy nhất với kết nối 5G/eSIM."
image: "/images/spatial.svg"
draft: false
---

## Nhận thức không gian

Nhận thức thực tế cho các bản sao kỹ thuật số, robot và hệ thống tự động. Kết nối LiDAR, camera và cảm biến công nghiệp với các mô hình không gian của bạn.

---

## Vấn đề

Các cảm biến tạo ra dữ liệu. Để hiểu được dữ liệu đó, cần phải:

- **Bối cảnh** — Cảm biến nằm ở đâu? Nó đang quan sát gì?
- **Fusion** — Kết hợp nhiều luồng dữ liệu cảm biến thành một hình ảnh thống nhất
- **Tích hợp** — Kết nối với các công cụ thiết kế, không chỉ các bảng điều khiển

Hầu hết các giải pháp cảm biến chỉ dừng lại ở việc thu thập dữ liệu. Chúng tôi kết nối cảm biến với các mô hình không gian.

---

## Chế độ triển khai

Cùng một thiết bị biên. Cùng các cảm biến. Cấu hình khác nhau.

| Mode | Platform | Use Case |
|------|----------|----------|
| **Aerial** | DJI enterprise drones | Surveying, inspection, mapping |
| **Ground** | Tripod, backpack | Interior scanning, construction |
| **Robot** | Viam RDK, ROS2 | Navigation, pick-and-place |
| **Fixed** | Permanent mount | Traffic, security, warehouse |

---

## Trừu tượng hóa cảm biến

**Thiết kế không phụ thuộc vào phần cứng.** Mã nguồn của bạn tương tác với API thống nhất của chúng tôi, không phải với các trình điều khiển cảm biến riêng lẻ.

### Các loại cảm biến được hỗ trợ

| Type | Examples |
|------|----------|
| **LiDAR** | Livox Mid-360, Avia |
| **RGB-D Cameras** | Intel RealSense, Luxonis OAK-D |
| **Position** | GPS/GNSS (u-blox RTK), IMU |
| **Industrial** | Modbus sensors, CAN bus |

Thay thế cảm biến mà không cần thay đổi mã nguồn. Dựa trên cấu hình, không dựa trên mã nguồn.

---

## Kiến trúc Edge Agent

Chạy chương trình nhị phân trên phần cứng của bạn — Raspberry Pi, Jetson, Linux công nghiệp hoặc ARM tùy chỉnh.

| Capability | Description |
|------------|-------------|
| **Plugin system** | Add sensors via config, not code changes |
| **Local buffering** | Store-and-forward when offline |
| **Real-time streaming** | NATS JetStream to cloud |
| **Lightweight** | Single binary, no runtime dependencies |

---

## Kết nối

### 5G/LTE với eSIM qua OTA

Không có việc đổi SIM. Không quét mã QR. Cài đặt tự động từ máy chủ.

- Modem được cung cấp kèm theo hồ sơ khởi động
- Nền tảng của bạn kích hoạt quá trình tải xuống hồ sơ nhà mạng
- Chuyển đổi nhà cung cấp dịch vụ trong quá trình triển khai thông qua API

Dành cho máy bay không người lái (drone) trong không trung, robot di động và các hệ thống cố định tại các vị trí xa xôi.

---

## Tích hợp với Spatial

Các cảm biến được kết nối trực tiếp với mô hình 3D của bạn:

| Data Flow | Purpose |
|-----------|---------|
| Point clouds → Spatial model | Reality capture |
| GPS/IMU → Model positioning | Georeferencing |
| Environmental sensors → Twin | Live state updates |
| Industrial I/O → Automation | Closed-loop control |

Không chỉ là bảng điều khiển. Cảm biến trong bối cảnh.

---

## Xây dựng trên nền tảng

Sensing inherits all [Foundation](/platform/foundation/) capabilities automatically:

| Capability | What It Means |
|------------|---------------|
| **Offline-first** | Capture without internet, sync when connected |
| **Universal deployment** | Edge, mobile, desktop—same agent |
| **Self-sovereign** | Your sensors, your data, your servers |
| **Real-time sync** | Stream to multiple destinations simultaneously |

[Tìm hiểu thêm về Quỹ →](/platform/foundation/)

---

## Một phần của điều gì đó lớn lao hơn

Sensing là lớp nhận thức của nền tảng phần mềm Ubuntu.

Đối với các tổ chức cần thiết kế 3D và trí tuệ nhân tạo (AI), nền tảng Spatial của chúng tôi cung cấp môi trường thiết kế với tích hợp trực tiếp vào dữ liệu cảm biến của bạn.

[Khám phá Không gian →](/platform/spatial/)

---

## Hãy cùng chúng tôi xây dựng

Đang triển khai cảm biến? Đang phát triển hệ thống nhận thức? Hãy cùng thảo luận.

[Liên hệ →](/contact)
