---
title: "Ngăn xếp phần mềm"
meta_title: "Kiến trúc phần mềm đội bay drone | Ubuntu Software"
description: "Ngăn xếp phần mềm cho đội bay 1.000 drone: firmware điều khiển bay PX4, Ubuntu Server, JetPack, ROS 2 và Vehicle Gateway viết bằng Go."
image: "/images/robotics.svg"
draft: false
---

## Kiến trúc phần mềm

Mọi thành phần phần mềm đều chạy mã nguồn mở với trách nhiệm rõ ràng. Ngăn xếp phần mềm được thiết kế cho quản lý từ xa, cập nhật qua mạng (OTA) và hành vi nhất quán trên hơn 1.000 phương tiện.

---

## Bộ điều khiển bay: Pixhawk

### PX4 v1.14 LTS

Bộ điều khiển bay chạy **PX4 Autopilot v1.14 LTS**:

| Khía cạnh | Chi tiết |
|--------|---------|
| **Phiên bản** | v1.14.x LTS (Hỗ trợ dài hạn) |
| **Thời gian hỗ trợ** | Bản vá bảo mật và lỗi nghiêm trọng trong hơn 2 năm |
| **Cấu hình** | Tham số được chuẩn hóa trên toàn đội phương tiện |
| **Cập nhật** | Qua QGroundControl hoặc MAVLink |

**Tại sao chọn LTS:**

- **Ổn định** — Mã nguồn đã được kiểm chứng trong môi trường sản xuất, ít thay đổi
- **Bảo mật** — Các lỗ hổng nghiêm trọng được vá mà không thay đổi tính năng
- **Khả năng dự đoán** — Hành vi đã biết trên toàn bộ đội phương tiện
- **Hỗ trợ** — Có hỗ trợ thương mại từ hệ sinh thái Dronecode

**Các mô-đun chính:**

- **EKF2** — Bộ lọc Kalman mở rộng để ước lượng trạng thái
- **Position Controller** — Điều khiển vị trí và vận tốc cho máy bay đa cánh quạt
- **Navigator** — Thực thi nhiệm vụ, RTL, hàng rào địa lý (geofence)
- **Commander** — Máy trạng thái, xử lý an toàn khi sự cố (failsafe)
- **MAVLink** — Giao thức đo xa và điều khiển lệnh

---

## Máy tính đồng hành cảm biến: Raspberry Pi

### Ubuntu Server 22.04 LTS

Raspberry Pi chạy **Ubuntu Server 22.04**:

| Khía cạnh | Chi tiết |
|--------|---------|
| **Bản phân phối** | Ubuntu Server 22.04 LTS |
| **Nhân (Kernel)** | Linux mainline kèm các bản vá cho Pi |
| **Hệ thống khởi tạo** | systemd |
| **Thời gian hỗ trợ** | Đến tháng 4 năm 2027 |

**Tại sao chọn Ubuntu Server (không phải Raspberry Pi OS):**

- **Nhất quán** — Cùng một hệ điều hành trên Pi và Jetson
- **Đóng gói** — Kho apt tiêu chuẩn, có sẵn snap
- **Hỗ trợ LTS** — Cập nhật bảo mật trong 5 năm
- **Công cụ** — Quản trị quen thuộc với các đội vận hành

### Dịch vụ cốt lõi

| Dịch vụ | Mục đích |
|---------|---------|
| **mavlink-router** | Định tuyến thông điệp MAVLink giữa FCU và Jetson |
| **sensor-agent** | Thu thập dữ liệu cảm biến môi trường |
| **gpsd** | Daemon GPS/GNSS cung cấp dữ liệu vị trí |
| **chrony** | Đồng bộ thời gian NTP |

**mavlink-router** đóng vai trò then chốt—nó chuyển tiếp lưu lượng MAVLink từ cổng nối tiếp của Pixhawk đến Jetson qua UDP, cho phép Vehicle Gateway xử lý dữ liệu đo xa.

---

## Máy tính đồng hành AI: Jetson

### JetPack 6.x

Jetson chạy **NVIDIA JetPack 6.x**:

| Thành phần | Phiên bản |
|-----------|---------|
| **L4T (Linux for Tegra)** | Dựa trên Ubuntu 22.04 |
| **CUDA** | 12.x |
| **cuDNN** | 8.x |
| **TensorRT** | 8.x |
| **VPI (Vision Programming Interface)** | 3.x |

**Tại sao chọn JetPack 6:**

- **Hỗ trợ Orin mới nhất** — Tăng tốc phần cứng đầy đủ trên dòng Orin
- **Nền tảng Ubuntu 22.04** — Đồng bộ với Pi để vận hành nhất quán
- **Hỗ trợ container** — NVIDIA Container Runtime cho các khối lượng công việc được cô lập
- **Cập nhật OTA** — Công cụ của NVIDIA để cập nhật firmware trên toàn đội phương tiện

### CUDA + TensorRT

Suy luận AI chạy thông qua runtime **TensorRT**:

- **Tối ưu hóa mô hình** — Lượng tử hóa INT8/FP16 cho triển khai tại biên
- **Suy luận theo lô** — Xử lý nhiều khung hình một cách hiệu quả
- **Độ trễ thấp** — Truy cập trực tiếp bộ nhớ GPU

Các khối lượng công việc điển hình:

- Phát hiện đối tượng (YOLO, SSD)
- Phân đoạn ngữ nghĩa
- Luồng quang học (optical flow)
- Ước lượng độ sâu

### ROS 2 Humble

**ROS 2 Humble Hawksbill** cung cấp phần mềm trung gian (middleware) cho robot:

| Khía cạnh | Chi tiết |
|--------|---------|
| **Bản phân phối** | Humble Hawksbill |
| **Thời gian hỗ trợ** | Đến tháng 5 năm 2027 |
| **DDS** | CycloneDDS (mặc định) hoặc FastDDS |
| **Hệ thống build** | colcon |

**Tại sao chọn ROS 2 Humble:**

- **Bản phát hành LTS** — Hỗ trợ 5 năm, tương đồng với Ubuntu/JetPack
- **Khả năng thời gian thực** — Có sẵn các luồng thực thi tất định
- **Tích hợp cảm biến** — Nhiều gói driver phong phú
- **Mô phỏng** — Tích hợp Gazebo để kiểm thử

**Các gói chính:**

| Gói | Mục đích |
|---------|---------|
| **px4_msgs** | Định nghĩa thông điệp MAVLink cho ROS 2 |
| **px4_ros_com** | Cầu nối PX4-ROS 2 (uXRCE-DDS) |
| **image_transport** | Xử lý dữ liệu camera |
| **tf2** | Biến đổi hệ tọa độ |

---

## Vehicle Gateway

### Triển khai bằng Go

**Vehicle Gateway** là một dịch vụ Go tùy chỉnh chạy trên mỗi Jetson:

| Khía cạnh | Chi tiết |
|--------|---------|
| **Ngôn ngữ** | Go 1.21+ |
| **Phụ thuộc** | nats.go, thư viện mavlink |
| **Triển khai** | Dịch vụ systemd hoặc container |
| **Cấu hình** | Biến môi trường + tệp cấu hình |

**Tại sao chọn Go:**

- **Một tệp nhị phân duy nhất** — Không phụ thuộc runtime, triển khai đơn giản
- **Hiệu năng** — Độ trễ thấp, sử dụng bộ nhớ hiệu quả
- **Xử lý đồng thời** — Goroutine để xử lý song song MAVLink/NATS
- **Hệ sinh thái NATS** — Thư viện client nats.go được hỗ trợ hàng đầu

**Trách nhiệm:**

1. **Tiếp nhận MAVLink** — Nhận và phân tích thông điệp MAVLink từ mavlink-router
2. **Giảm tần suất lấy mẫu trạng thái** — Giảm tốc độ dữ liệu đo xa để truyền qua WAN
3. **Trích xuất sự kiện** — Tạo sự kiện từ các chuyển đổi trạng thái
4. **Thực thi lệnh** — Nhận lệnh NATS, xác thực, gửi đến FCU
5. **Đối chiếu bản sao (Shadow)** — Đồng bộ trạng thái mong muốn/đã báo cáo với backend của đội phương tiện

[Chi tiết Gateway →]({{< relref "/fleet/gateway" >}})

---

## Mặt phẳng điều khiển

### NATS + JetStream

Việc điều phối đội phương tiện chạy trên **NATS** với khả năng lưu trữ bền vững **JetStream**:

| Thành phần | Vai trò |
|-----------|------|
| **NATS Server** | Nhắn tin pub/sub cốt lõi |
| **JetStream** | Lưu trữ bền vững luồng, kho KV |
| **Leaf Nodes** | Các phiên bản NATS cục bộ trên phương tiện |
| **Hub Clusters** | Điểm tổng hợp theo khu vực |

**Tại sao chọn NATS:**

- **Hiệu năng** — Hàng triệu thông điệp mỗi giây trên mỗi nút
- **Đơn giản** — Một tệp nhị phân duy nhất, cấu hình tối thiểu
- **Lưu trữ bền vững** — JetStream cho các luồng bền vững và KV
- **Cấu trúc liên kết** — Nút lá, cụm, siêu cụm cho mọi quy mô

**Cấu trúc liên kết của đội phương tiện:**

```
┌─────────────────────────────────────────────────────┐
│                   Global Hub (optional)             │
│                   NATS Supercluster                 │
└─────────────────┬───────────────────┬───────────────┘
                  │                   │
      ┌───────────▼───────┐ ┌─────────▼───────────┐
      │   Regional Hub A  │ │   Regional Hub B    │
      │   3-node cluster  │ │   3-node cluster    │
      └─────────┬─────────┘ └─────────┬───────────┘
                │                     │
    ┌───────────┼───────────┐         │
    │           │           │         │
┌───▼───┐ ┌─────▼───┐ ┌─────▼───┐ ┌───▼───┐
│ Leaf  │ │  Leaf   │ │  Leaf   │ │ Leaf  │
│ VID-1 │ │ VID-2   │ │ VID-3   │ │VID-N  │
└───────┘ └─────────┘ └─────────┘ └───────┘
```

Mỗi phương tiện chạy một **nút lá (leaf node)** kết nối với hub khu vực của nó. Nút lá xử lý pub/sub cục bộ và lưu trữ-chuyển tiếp (store-and-forward) khi mất kết nối.

[Chi tiết cấu trúc liên kết NATS →]({{< relref "/fleet/nats-topology" >}})

---

## Chiến lược cập nhật phần mềm

### OTA trên toàn đội phương tiện

| Thành phần | Cơ chế cập nhật |
|-----------|------------------|
| **PX4** | Tải firmware lên qua MAVLink bằng QGC hoặc công cụ tùy chỉnh |
| **Pi/Ubuntu** | apt + unattended-upgrades, hoặc triển khai image tùy chỉnh |
| **Jetson/JetPack** | Công cụ OTA của NVIDIA, cập nhật container |
| **Vehicle Gateway** | Thay thế tệp nhị phân qua systemd, hoặc cập nhật container |
| **NATS Leaf** | Cập nhật tệp nhị phân với khởi động lại an toàn |

**Nguyên tắc:**

- **Triển khai theo giai đoạn** — Cập nhật một nhóm nhỏ, xác minh, rồi mở rộng
- **Khả năng khôi phục** — Phiên bản trước luôn sẵn có
- **Kiểm tra tình trạng** — Tự động xác minh sau khi cập nhật
- **Thời gian ngừng hoạt động tối thiểu** — Cập nhật trong khung giờ bảo trì

---

## Tóm tắt

| Lớp | Phần mềm | Phiên bản |
|-------|----------|---------|
| **Bộ điều khiển bay** | PX4 | v1.14 LTS |
| **Máy tính đồng hành cảm biến** | Ubuntu Server | 22.04 LTS |
| **Máy tính đồng hành cảm biến** | mavlink-router | Mới nhất |
| **Máy tính đồng hành AI** | JetPack | 6.x |
| **Máy tính đồng hành AI** | CUDA/TensorRT | 12.x/8.x |
| **Máy tính đồng hành AI** | ROS 2 | Humble |
| **Máy tính đồng hành AI** | Vehicle Gateway | Go 1.21+ |
| **Mặt phẳng điều khiển** | NATS + JetStream | Mới nhất |

---

## Tiếp theo

[Cấu trúc liên kết NATS →]({{< relref "/fleet/nats-topology" >}})
