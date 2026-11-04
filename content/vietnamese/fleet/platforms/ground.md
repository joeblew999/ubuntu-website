---
title: "Nền tảng phương tiện mặt đất"
meta_title: "Phần cứng & Giao thức cho Đội xe Mặt đất | Ubuntu Software"
description: "Kiến trúc đội xe cho ô tô, xe tải và AGV: tích hợp CAN bus, J1939, ROS 2, điện toán công nghiệp và các mô hình an toàn riêng cho từng loại phương tiện."
image: "/images/robotics.svg"
draft: false
---

## Phần cứng cho đội phương tiện mặt đất

Phương tiện mặt đất—ô tô, xe tải, AGV và các nền tảng tự hành—sử dụng cùng kiến trúc đội phương tiện như máy bay không người lái nhưng với phần cứng, giao thức và mô hình an toàn phù hợp với phương tiện.

---

## Tùy chọn giao thức

Phương tiện mặt đất giao tiếp bằng các giao thức tiêu chuẩn ngành:

### CAN Bus

| Khía cạnh | Chi tiết |
|--------|---------|
| **Tiêu chuẩn** | CAN 2.0B / CAN FD |
| **Tốc độ dữ liệu** | 500kbps (CAN) / 5Mbps (CAN FD) |
| **Định dạng thông điệp** | Tệp DBC định nghĩa cách mã hóa tín hiệu |
| **Lớp vật lý** | Cặp dây vi sai, bền bỉ trong môi trường nhiều nhiễu |

**CAN bus cung cấp:**
- Truy cập trực tiếp vào các hệ thống của xe (động cơ, hộp số, phanh)
- Dữ liệu telemetry thời gian thực (tốc độ, RPM, nhiên liệu, nhiệt độ)
- Mã lỗi chẩn đoán (OBD-II / J1939)
- Tích hợp thiết bị lắp thêm (aftermarket)

### J1939 (Phương tiện hạng nặng)

Dành cho xe tải, xe buýt và thiết bị hạng nặng:

| Khía cạnh | Chi tiết |
|--------|---------|
| **Tiêu chuẩn** | SAE J1939 |
| **Truyền tải** | CAN 2.0B ở 250kbps |
| **Định địa chỉ** | Định danh 29-bit với PGN/SPN |
| **Ứng dụng** | Động cơ, hộp số, phanh, rơ-moóc |

**J1939 cung cấp:**
- Các nhóm tham số (PGN) được chuẩn hóa giữa các nhà sản xuất
- Tính nhất quán của telemetry trên toàn đội phương tiện
- Tích hợp với các nền tảng telematics
- Tuân thủ các yêu cầu quy định

### ROS 2 (Nền tảng tự hành)

Dành cho các phương tiện tự hành được chế tạo chuyên dụng:

| Khía cạnh | Chi tiết |
|--------|---------|
| **Middleware** | ROS 2 Humble / Iron |
| **Truyền tải** | DDS (Cyclone DDS, Fast DDS) |
| **Loại thông điệp** | Tiêu chuẩn (sensor_msgs, geometry_msgs) + tùy chỉnh |
| **QoS** | Độ tin cậy, độ bền dữ liệu có thể cấu hình |

**ROS 2 cung cấp:**
- Tích hợp cảm biến gốc (LiDAR, camera, IMU)
- Các loại thông điệp tiêu chuẩn cho robot
- Bắc cầu topic trực tiếp sang NATS
- Tương thích với mô phỏng (Gazebo, CARLA)

---

## Ngăn xếp phần cứng

### Nền tảng tính toán

Phương tiện mặt đất thường sử dụng máy tính công nghiệp được gia cố:

| Nền tảng | Trường hợp sử dụng | Ghi chú |
|----------|----------|-------|
| **NVIDIA Jetson AGX Orin** | Phương tiện tự hành | Nhận thức hiệu năng cao |
| **NVIDIA Jetson Orin NX** | Phương tiện trong đội | Cân bằng hiệu năng/chi phí |
| **Máy tính công nghiệp** | Chỉ telematics | Chi phí thấp hơn, không cần GPU |
| **Neousys** | Ô tô, môi trường khắc nghiệt | Dải nhiệt độ rộng, điều khiển theo khóa điện |

**Nhiệm vụ của máy tính:**
- Vehicle Gateway (CAN/ROS 2 sang NATS)
- NATS leaf node
- Suy luận AI tại biên (tùy chọn)
- Lưu trữ và chuyển tiếp khi mất kết nối

### Giao diện CAN

| Thiết bị | Giao diện | Ghi chú |
|--------|-----------|-------|
| **PEAK PCAN-USB** | USB | Phổ biến, được hỗ trợ tốt |
| **Kvaser Leaf Light** | USB | Cấp công nghiệp |
| **Thiết bị SocketCAN** | USB/PCIe | Hỗ trợ gốc trên Linux |
| **CAN gốc của Jetson** | Tích hợp sẵn | Có trên AGX Orin |

### Kết nối

| Phương thức | Trường hợp sử dụng |
|--------|----------|
| **Mạng di động 4G/5G** | Kết nối chính cho đội phương tiện |
| **WiFi** | Hoạt động tại bãi xe/cơ sở |
| **Ethernet** | Kết nối có dây tại trạm đỗ |
| **Vệ tinh** | Hoạt động ở vùng xa/nông thôn |

---

## Các loại phương tiện

### Xe chở khách (Ô tô)

| Thành phần | Lựa chọn điển hình |
|-----------|-------------------|
| **Giao thức** | CAN bus (cổng OBD-II hoặc trực tiếp) |
| **Tính toán** | Jetson Orin NX hoặc máy tính công nghiệp |
| **Nguồn điện** | Hệ thống điện xe 12V |
| **Kết nối** | Mạng di động 4G/5G |
| **An toàn** | Người lái ghi đè, dừng khẩn cấp |

**Trường hợp sử dụng:** Đội xe gọi xe/chia sẻ chuyến đi, xe cho thuê, xe giao hàng

### Xe thương mại (Xe tải)

| Thành phần | Lựa chọn điển hình |
|-----------|-------------------|
| **Giao thức** | J1939 (dựa trên CAN) |
| **Tính toán** | Máy tính công nghiệp gia cố hoặc Jetson |
| **Nguồn điện** | Hệ thống điện xe 12V/24V |
| **Kết nối** | 4G/5G + vệ tinh dự phòng |
| **An toàn** | Người lái ghi đè, tuân thủ ELD |

**Trường hợp sử dụng:** Vận tải đường dài, đội xe giao hàng, xe dịch vụ

### Phương tiện mặt đất tự hành (AGV)

| Thành phần | Lựa chọn điển hình |
|-----------|-------------------|
| **Giao thức** | ROS 2 gốc hoặc giao thức độc quyền |
| **Tính toán** | Jetson AGX Orin |
| **Nguồn điện** | Hệ thống pin 24V/48V |
| **Kết nối** | WiFi (trong cơ sở) + mạng di động (ngoài trời) |
| **An toàn** | Dừng khẩn cấp, cảm biến đạt chuẩn an toàn |

**Trường hợp sử dụng:** Logistics kho bãi, khai thác mỏ, nông nghiệp, vận hành cảng

---

## Mô hình an toàn

An toàn của phương tiện mặt đất khác với máy bay không người lái nhưng tuân theo cùng các nguyên tắc:

### Phân cấp thẩm quyền

```
Manual Override (highest authority)
    │
    ├── Steering wheel / brake pedal (driver present)
    ├── E-stop button (AGVs)
    └── Remote operator console
    │
    ▼
Autonomous System (software control)
    │
    ├── Perception → Planning → Control
    └── Geofencing, speed limits, operational domain
    │
    ▼
Fleet Management (lowest authority)
    │
    └── Mission assignment, routing, monitoring
```

### Hành vi an toàn khi gặp sự cố (Failsafe)

| Tác nhân kích hoạt | Tương đương trên drone | Phản ứng của phương tiện mặt đất |
|---------|------------------|------------------------|
| Mất liên lạc | Quay về điểm cất cánh (Return-to-Launch) | Dừng tại chỗ / tấp vào lề |
| Lỗi cảm biến | Hạ cánh ngay lập tức | Dừng an toàn, cảnh báo người vận hành |
| Vượt hàng rào địa lý | Quay về ranh giới | Dừng tại ranh giới |
| Kích hoạt dừng khẩn cấp | Ngắt động cơ | Phanh ngay lập tức |
| Pin yếu | RTL / hạ cánh | Quay về bãi xe / dừng an toàn |

### An toàn mạng

**Cùng nguyên tắc như drone:** NATS không bao giờ nằm trong vòng điều khiển.

| NATS làm gì | NATS không làm gì |
|----------------|---------------------|
| Thu thập telemetry | Lái thời gian thực |
| Phân công nhiệm vụ | Lệnh phanh |
| Điều phối đội phương tiện | Điều khiển ga |
| Giám sát & cảnh báo | Quyết định quan trọng về an toàn |

Các hệ thống trên xe xử lý mọi chức năng quan trọng về an toàn. Mất mạng nghĩa là mất khả năng giám sát, không phải mất quyền điều khiển phương tiện.

---

## Triển khai Gateway

Vehicle Gateway bắc cầu các giao thức gốc của phương tiện sang NATS:

```
┌─────────────────────────────────────────────────────────────────┐
│                    GROUND VEHICLE GATEWAY                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   CAN Bus ─────► DBC Parser ─────► Telemetry Subjects           │
│                                                                  │
│   J1939 ───────► PGN Decoder ────► Telemetry Subjects           │
│                                                                  │
│   ROS 2 ───────► Topic Bridge ───► Telemetry Subjects           │
│                                                                  │
│   Commands ◄───────────────────── Command Subjects              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS JETSTREAM                              │
│   fleet.prod.veh.{vehicle_id}.state.*                           │
│   fleet.prod.veh.{vehicle_id}.evt.*                             │
│   fleet.prod.veh.{vehicle_id}.cmd.*                             │
└─────────────────────────────────────────────────────────────────┘
```

**Nhiệm vụ của gateway:**
- Chuyển đổi giao thức (CAN/J1939/ROS 2 → định dạng telemetry chung)
- Giảm tần số lấy mẫu trạng thái (giảm dữ liệu CAN tần số cao xuống mức dễ quản lý)
- Trích xuất sự kiện (phát hiện thay đổi trạng thái, vượt ngưỡng)
- Xác thực lệnh (kiểm tra an toàn trước khi thực thi)
- Lưu trữ và chuyển tiếp (đệm dữ liệu khi mất kết nối)

---

## Ánh xạ Subject

Phương tiện mặt đất sử dụng cùng phân cấp subject như drone:

| Subject | Nội dung với phương tiện mặt đất |
|---------|----------------------|
| `fleet.prod.veh.{id}.state.pos` | Vị trí GPS, hướng di chuyển |
| `fleet.prod.veh.{id}.state.vel` | Tốc độ, gia tốc |
| `fleet.prod.veh.{id}.state.health` | Trạng thái động cơ, chẩn đoán |
| `fleet.prod.veh.{id}.state.fuel` | Mức nhiên liệu / SOC pin |
| `fleet.prod.veh.{id}.evt.trip` | Sự kiện bắt đầu/kết thúc chuyến đi |
| `fleet.prod.veh.{id}.evt.geofence` | Vào/ra khu vực |
| `fleet.prod.veh.{id}.cmd.mission` | Phân công tuyến đường |

Cấu trúc subject không phụ thuộc vào loại phương tiện. Một đội phương tiện kết hợp drone và xe tải sử dụng các mẫu giống hệt nhau.

---

## Mẫu tích hợp

### Tuân thủ ELD (Xe tải tại Hoa Kỳ)

Yêu cầu về Thiết bị Ghi nhật ký Điện tử (Electronic Logging Device):
- Theo dõi Giờ làm việc (Hours of Service)
- Nhận dạng tài xế
- Ghi nhật ký vị trí
- Truyền dữ liệu cho cơ quan chức năng

Gateway thu thập dữ liệu động cơ J1939 và tích hợp với các hệ thống tuân thủ ELD.

### Tích hợp Telematics

Các nền tảng telematics hiện có có thể kết nối qua:
- NATS grants (đăng ký theo dõi các phương tiện cụ thể)
- REST API (truy vấn dữ liệu lịch sử)
- Sự kiện webhook (cảnh báo thời gian thực)

### Hệ thống quản lý đội phương tiện

Tích hợp với TMS/FMS doanh nghiệp:
- Phân công phương tiện
- Tối ưu hóa tuyến đường
- Quản lý nhiên liệu
- Lập lịch bảo trì

---

## Tóm tắt

| Khía cạnh | Drone | Phương tiện mặt đất |
|--------|--------|-----------------|
| **Giao thức** | MAVLink | CAN bus / J1939 / ROS 2 |
| **Tính toán** | Jetson Orin Nano | Jetson / Máy tính công nghiệp |
| **Nguồn điện** | 4S LiPo (14.8V) | Hệ thống điện xe 12V/24V/48V |
| **Ghi đè an toàn** | Bộ phát RC | Vô lăng / Dừng khẩn cấp |
| **Failsafe** | Quay về điểm cất cánh (Return-to-Launch) | Dừng tại chỗ |
| **Quy định** | Hàng không (FAA/EASA) | Vận tải đường bộ (DOT) |
| **Kiến trúc đội phương tiện** | **Giống hệt** | **Giống hệt** |

---

## Tài liệu liên quan

- [Các nền tảng được hỗ trợ]({{< relref "/fleet/platforms" >}}) — Tổng quan về tất cả các loại phương tiện
- [Nền tảng Drone]({{< relref "/fleet/platforms/drones" >}}) — Phần cứng UAV và MAVLink
- [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) — Kiến trúc bắc cầu giao thức
- [Mô hình an toàn]({{< relref "/fleet/safety" >}}) — Nguyên tắc an toàn không phụ thuộc loại phương tiện
- [Đặt tên Subject]({{< relref "/fleet/subjects" >}}) — Phân cấp subject telemetry

---

## Bắt đầu

Đang triển khai đội phương tiện mặt đất? Chúng tôi có thể hỗ trợ về kiến trúc và tích hợp.

[Liên hệ với chúng tôi →](/contact)
