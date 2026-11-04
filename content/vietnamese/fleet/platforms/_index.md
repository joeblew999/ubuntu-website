---
title: "Các nền tảng được hỗ trợ"
meta_title: "Các nền tảng phương tiện được hỗ trợ | Ubuntu Software"
description: "Kiến trúc đội phương tiện của chúng tôi hỗ trợ drone, ô tô, xe tải và phương tiện mặt đất tự hành. Cùng một hạ tầng cốt lõi, tích hợp riêng cho từng loại phương tiện."
image: "/images/robotics.svg"
draft: false
---

## Một kiến trúc, nhiều phương tiện

Kiến trúc đội phương tiện của Ubuntu Software **được thiết kế để không phụ thuộc vào loại phương tiện**. Cùng một hạ tầng NATS JetStream, các mẫu bản sao số (digital twin) và mô hình phân quyền dùng để quản lý đội drone cũng quản lý ô tô, xe tải và phương tiện mặt đất tự hành.

---

## Những gì dùng chung

Các thành phần này hoạt động giống hệt nhau trên mọi loại phương tiện:

| Thành phần | Mô tả |
|-----------|-------------|
| **NATS JetStream** | Xương sống truyền thông điệp, hoạt động với mọi dữ liệu đo xa |
| **Phân cấp Subject** | `fleet.{env}.veh.{id}.*` — ID phương tiện chỉ là một định danh |
| **Bản sao số** | Luồng trạng thái, luồng sự kiện, kho shadow |
| **Phân quyền** | Bảo mật phi tập trung, cấp quyền, truy cập của bên thứ ba |
| **Luồng & Sự kiện** | Tổng hợp dữ liệu đo xa, nhật ký kiểm toán, trạng thái có thể truy vấn |

**Kiến trúc không quan tâm thứ gì đang tạo ra dữ liệu đo xa.** Một bản cập nhật vị trí từ drone và một bản cập nhật vị trí từ xe tải đều đi qua cùng một hạ tầng.

---

## Những gì riêng cho từng loại phương tiện

Mỗi nền tảng có những yêu cầu riêng:

| Khía cạnh | Khác nhau theo nền tảng |
|--------|-------------------|
| **Phần cứng** | Khung máy bay so với khung gầm, bộ điều khiển bay so với ECU |
| **Giao thức** | MAVLink, CAN bus, ROS 2, J1939 |
| **Mô hình an toàn** | Điều khiển RC can thiệp so với dừng khẩn cấp (e-stop), RTL so với dừng tại chỗ |
| **Cảm biến** | Độ cao so với khoảng sáng gầm, tốc độ không khí so với tốc độ bánh xe |
| **Quy định** | Cơ quan hàng không so với cơ quan vận tải đường bộ |

---

## Các nền tảng được hỗ trợ

### [Drone]({{< relref "/fleet/platforms/drones" >}})

Phương tiện bay không người lái chạy PX4 hoặc ArduPilot:

- **Phần cứng**: Holybro X500, Pixhawk 6X, NVIDIA Jetson
- **Giao thức**: MAVLink qua serial/UDP
- **An toàn**: Điều khiển RC can thiệp, Return-to-Launch, hàng rào địa lý (geofencing)
- **Trường hợp sử dụng**: Kiểm tra, lập bản đồ, giám sát, giao hàng

[Chi tiết nền tảng Drone →]({{< relref "/fleet/platforms/drones" >}})

---

### [Phương tiện mặt đất]({{< relref "/fleet/platforms/ground" >}})

Ô tô, xe tải, AGV và các nền tảng mặt đất tự hành:

- **Phần cứng**: ECU phương tiện, máy tính công nghiệp, Jetson
- **Giao thức**: CAN bus, J1939, ROS 2
- **An toàn**: Dừng khẩn cấp (e-stop), dừng tại chỗ, hàng rào địa lý
- **Trường hợp sử dụng**: Logistics, khai khoáng, nông nghiệp, giao hàng chặng cuối

[Chi tiết phương tiện mặt đất →]({{< relref "/fleet/platforms/ground" >}})

---

## Mẫu Gateway

[Gateway phương tiện]({{< relref "/fleet/gateway" >}}) làm cầu nối giữa các giao thức gốc của phương tiện và NATS. Mỗi nền tảng có một bản triển khai gateway riêng:

```
┌─────────────────────────────────────────────────────────────────┐
│                     VEHICLE GATEWAY                              │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   MAVLink GW    │    CAN GW       │      ROS 2 GW               │
│   (Drones)      │  (Cars/Trucks)  │   (ROS Robots)              │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS JETSTREAM                              │
│            Same subjects, same streams, same patterns            │
└─────────────────────────────────────────────────────────────────┘
```

Gateway:
- Chuyển đổi các thông điệp đặc thù của từng giao thức sang định dạng dữ liệu đo xa chung
- Xuất bản lên phân cấp subject tiêu chuẩn
- Thực thi các chính sách an toàn
- Xử lý lưu trữ và chuyển tiếp (store-and-forward) khi mất kết nối

---

## Bổ sung nền tảng mới

Kiến trúc hỗ trợ bất kỳ phương tiện nào có thể:

1. **Tạo dữ liệu đo xa** — Dữ liệu vị trí, vận tốc, tình trạng
2. **Nhận lệnh** — Khởi động, dừng, điều hướng theo điểm tham chiếu (waypoint)
3. **Chạy một gateway** — Chuyển đổi giao thức gốc sang NATS

Các nền tảng tùy chỉnh được tích hợp bằng cách triển khai một gateway cho giao thức của chúng. Phần còn lại của hạ tầng (NATS, luồng, phân quyền) vẫn hoạt động mà không cần thay đổi.

---

## Tìm hiểu thêm

| Nền tảng | Tài liệu |
|----------|---------------|
| [Drone]({{< relref "/fleet/platforms/drones" >}}) | Tích hợp PX4, Pixhawk, MAVLink |
| [Phương tiện mặt đất]({{< relref "/fleet/platforms/ground" >}}) | Tích hợp CAN bus, J1939, ô tô |
| [Gateway phương tiện]({{< relref "/fleet/gateway" >}}) | Kiến trúc cầu nối giao thức |
| [Mô hình an toàn]({{< relref "/fleet/safety" >}}) | Yêu cầu an toàn riêng cho từng nền tảng |

---

## Bắt đầu

Đang xây dựng một đội phương tiện hỗn hợp? Triển khai phương tiện mặt đất cùng với drone? Chúng tôi có thể hỗ trợ về kiến trúc và triển khai.

[Liên hệ với chúng tôi →](/contact)
