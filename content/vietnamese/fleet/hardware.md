---
title: "Ngăn xếp phần cứng"
meta_title: "Thông số kỹ thuật phần cứng đội drone | Ubuntu Software"
description: "Thông số kỹ thuật phần cứng sản xuất cho đội bay 1.000 drone: khung Holybro X500, bộ điều khiển bay Pixhawk 6X, máy tính đồng hành Raspberry Pi và Jetson."
image: "/images/robotics.svg"
draft: false
aliases:
  - "/fleet/hardware/"
---

{{< notice "info" >}}
**Bạn đang tìm các loại phương tiện khác?** Trang này đề cập đến phần cứng drone. Đối với phương tiện mặt đất (ô tô, xe tải, AGV), xem [Nền tảng Phương tiện Mặt đất]({{< relref "/fleet/platforms/ground" >}}). Để có cái nhìn tổng quan đầy đủ, xem [Các Nền tảng Được Hỗ trợ]({{< relref "/fleet/platforms" >}}).
{{< /notice >}}

## Mua từ Cửa hàng Địa phương

Việc chế tạo drone cần thời gian để làm quen. **Hãy mua từ một cửa hàng có thể hỗ trợ bạn**—họ cung cấp phần cứng, chúng tôi cung cấp phần mềm quản lý đội bay.

Khi mua hàng từ cửa hàng đối tác, bạn nhận được:
- Phần cứng đã được kiểm tra trước, sẵn sàng kết nối
- Hỗ trợ tại chỗ từ những người am hiểu về drone
- Quyền truy cập vào hạ tầng đội bay được quản lý của chúng tôi

---

## Cửa hàng có Địa điểm Thực tế

Ghé thăm trực tiếp một cửa hàng để được hỗ trợ tận tay khi lắp ráp:

{{< featured-shops >}}

**Bạn là cửa hàng drone?** [Tham gia chương trình đối tác của chúng tôi]({{< relref "/partners/become-partner" >}}) để cung cấp hạ tầng đội bay kèm theo phần cứng bạn bán.

---

## Những gì Bạn Cần

| Bộ kit | Linh kiện | Nguồn cung cấp phổ biến |
|-----|------------|-----------------|
| **Drone Hoàn chỉnh** | Khung X500, Pixhawk 6X, GPS, động cơ, ESC, cánh quạt | GetFPV, NewBeeDrone, Holybro |
| **Máy tính Đồng hành** | Raspberry Pi 4 + Jetson Orin Nano | Digi-Key, Mouser, SparkFun |
| **Nguồn & Pin** | LiPo 4S 6000mAh, mô-đun nguồn | GetFPV, Tattu, Amazon |
| **Radio** | Bộ thu ExpressLRS, telemetry SiK | GetFPV, HappyModel, Holybro |
| **Mạng** | Switch ethernet công nghiệp, modem 5G | PUSR, Quectel |

[Xem danh sách linh kiện đầy đủ ↓](#full-bill-of-materials)

---

## Phần cứng cho Vận hành Quy mô Đội bay

Vận hành đội bay đòi hỏi phần cứng cân bằng giữa năng lực và khả năng bảo trì. Mỗi linh kiện trong bộ này đều được chọn vì độ tin cậy trong sản xuất, nguồn linh kiện sẵn có và khả năng bảo dưỡng tại hiện trường.

---

## Khung máy bay

### Holybro X500 V2 ARF

**Holybro X500 V2** là nền tảng lý tưởng để triển khai đội bay:

| Thông số | Giá trị |
|---------------|-------|
| **SKU** | 30125 |
| **Chiều dài cơ sở** | 500mm |
| **Trọng lượng khung** | ~410g (610g kèm động cơ/ESC) |
| **Trọng lượng cất cánh tối đa** | ~2kg |
| **Thời gian bay** | 15-20 phút (với tải trọng tiêu chuẩn) |
| **Ngàm động cơ** | 16x16mm / 19x19mm |
| **Bao gồm** | Động cơ (2216 KV920), ESC (20A BLHeli S), Cánh quạt (1045) |

**Tại sao chọn X500:**

- **Thiết kế đã được kiểm chứng** — Hàng nghìn chiếc đã được triển khai trên toàn thế giới, kiến thức cộng đồng phong phú
- **Linh kiện sẵn có** — Tay đòn, càng hạ cánh và phụ kiện dễ dàng tìm mua
- **Dễ bảo trì** — Cấu trúc mô-đun, có thể sửa chữa tại hiện trường
- **Tải trọng linh hoạt** — Đủ sức chứa hai máy tính đồng hành cùng các cảm biến

Bộ kit ARF (Almost Ready to Fly – Gần sẵn sàng bay) bao gồm khung, động cơ, ESC và cánh quạt—giảm sự khác biệt khi lắp ráp trên toàn bộ một đội bay lớn.

---

## Bộ điều khiển bay

### Pixhawk 6X

**Pixhawk 6X** chạy **PX4 v1.14 LTS** đảm nhận mọi chức năng quan trọng cho chuyến bay:

| Thông số | Giá trị |
|---------------|-------|
| **Pixhawk 6X** | SKU 11073 |
| **Bo mạch nền** | SKU 18117 |
| **Phụ kiện** | SKU 15011 (cáp & vỏ) |
| **Bộ xử lý** | STM32H753 (480MHz Cortex-M7) |
| **IMU** | Dự phòng ba lớp (ICM-42688-P, ICM-45686, BMI088) |
| **Khí áp kế** | Kép (MS5611, ICP-20100) |
| **Từ kế** | IST8310 |
| **Giao tiếp** | 3x CAN, 6x UART, SPI, I2C, PWM |

### Phụ kiện Bộ điều khiển bay

| Linh kiện | SKU | Mục đích |
|-----------|-----|---------|
| **Mô-đun nguồn PM06 V2** | 15019 | Cảm biến nguồn + BEC 5V/3A |
| **GPS Holybro M9N** | 12027 | Mô-đun GPS + la bàn |
| **Cột GPS** | 12033 | Cột gập giúp giảm nhiễu |
| **Công tắc an toàn & Còi** | 12007 | An toàn khi kích hoạt + cảnh báo âm thanh |
| **Đế chống rung** | 12010 | Giảm chấn chống rung |
| **Bộ cáp** | 12055 | Cáp JST-GH/CAN |

**Tại sao chọn Pixhawk 6X:**

- **Dự phòng** — Ba IMU, hai khí áp kế để bỏ phiếu cảm biến
- **Hỗ trợ PX4 LTS** — Ổn định lâu dài, bản vá bảo mật
- **Chuẩn FMUv6X** — Có thể thay thế lẫn nhau với các bo mạch tương thích khác
- **Cách ly rung động** — Tích hợp giảm chấn IMU

Bộ điều khiển bay đảm nhận:

- Ước lượng và điều khiển tư thế
- Giữ vị trí và điều hướng
- Hành vi an toàn dự phòng (RTL, hạ cánh, hàng rào địa lý)
- Xử lý tín hiệu đầu vào RC
- Telemetry MAVLink

---

## Máy tính Đồng hành

Các phương tiện trong đội bay chạy **hai máy tính đồng hành** với trách nhiệm riêng biệt:

### Raspberry Pi 4 / CM4 (Máy tính Đồng hành Cảm biến)

Đảm nhận tích hợp cảm biến nhẹ và thu thập dữ liệu:

| Thông số | Giá trị |
|---------------|-------|
| **Mẫu** | Raspberry Pi 4 Model B |
| **SKU** | RPI4-MODEL-B-8GB |
| **Bộ xử lý** | Cortex-A72 lõi tứ @ 1.8GHz |
| **RAM** | 8GB |
| **Lưu trữ** | microSD công nghiệp 128GB (SDSDQED-128G-GNSIN) |
| **Nguồn** | 5V @ 3A điển hình |

**Trách nhiệm:**

- Driver cảm biến môi trường (nhiệt độ, độ ẩm, chất lượng không khí)
- Ghi dữ liệu GPS/GNSS
- Chụp ảnh từ camera (tác vụ không dùng AI)
- Định tuyến MAVLink đến Jetson
- Lưu trữ và chuyển tiếp khi Jetson ngoại tuyến

### NVIDIA Jetson (Máy tính Đồng hành AI)

Đảm nhận các tác vụ đòi hỏi tính toán cao:

| Mẫu | SKU | Lõi GPU | Hiệu năng AI | Công suất |
|-------|-----|-----------|----------------|-------|
| **Jetson Orin Nano 8GB** ★ | 900-13767-0000-000 | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin Nano Super** | 945-13766-0000-000 | 1024 CUDA | 67 TOPS | 7-25W |
| **Jetson Orin NX** | 900-13767-0040-000 | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 900-13701-0000-000 | 2048 CUDA | 275 TOPS | 15-60W |

★ = Lựa chọn trong BOM tham chiếu

| Lưu trữ | SKU | Dung lượng |
|---------|-----|----------|
| **SSD NVMe** | WD-SN530-256G | 256GB |

**Trách nhiệm:**

- Thị giác máy tính (phát hiện, theo dõi đối tượng)
- Đo hành trình quán tính-thị giác
- Lập kế hoạch đường bay và tránh chướng ngại vật
- Vehicle Gateway (NATS client)
- NATS leaf node

Các triển khai đội bay thường sử dụng **Jetson Orin Nano** để suy luận tiết kiệm chi phí hoặc **Orin NX** cho các tác vụ nhận thức đòi hỏi cao.

---

## Hệ thống Radio

### Điều khiển RC: ExpressLRS

**ExpressLRS** cung cấp liên kết điều khiển ghi đè cho phi công:

| Thông số | Giá trị |
|---------------|-------|
| **Bộ thu** | ExpressLRS EP1 |
| **SKU** | HM-EP1-2400 |
| **Tần số** | 2.4GHz |
| **Tầm hoạt động** | 10km+ (tùy thuộc vào công suất/ăng-ten) |
| **Độ trễ** | <5ms |
| **Protocol** | Open-source (CRSF) |

**Why ExpressLRS:**

- **Open-source** — No vendor lock-in, community-driven development
- **Range** — Reliable link for visual line-of-sight operations
- **Low latency** — Responsive manual control when needed
- **Cost** — Affordable receivers for fleet scale

### Telemetry: SiK Radio

**SiK 915 MHz** radios provide ground station telemetry:

| Specification | Value |
|---------------|-------|
| **SKU** | 17012 (100mW) / 17013 (500mW) |
| **Frequency** | 915MHz (US) / 433MHz (EU) |
| **Range** | 300m-2km (depending on power) |
| **Data Rate** | Up to 250kbps |
| **Interface** | Serial (UART) |
| **Protocol** | MAVLink framing, FHSS |

**Role in Fleet:**

- Backup telemetry when cellular unavailable
- Ground control station connectivity
- Local testing and debugging

For production fleet operations, telemetry primarily flows through the Jetson's cellular connection to NATS. SiK radios serve as backup and for field diagnostics.

---

## Power System

### Battery: 4S LiPo

| Specification | Value |
|---------------|-------|
| **SKU** | TAA60004S35 |
| **Configuration** | 4S (14.8V nominal) |
| **Capacity** | 6000mAh |
| **Discharge Rate** | 35C |
| **Connector** | XT60 |

**Fleet Considerations:**

- Standardize on single battery configuration for logistics
- Use smart batteries with telemetry when available
- Implement battery rotation and health tracking
- Plan for 3:1 battery-to-vehicle ratio for continuous operations

### Power Distribution

| Component | SKU | Purpose |
|-----------|-----|---------|
| **PM06 V2 Power Module** | 15019 | Power sensing + 5V/3A BEC |

- **BEC** — 5V regulated supply for Pi, servos
- **Jetson power** — Direct from battery through regulator (12V typical)

### Networking

| Component | SKU | Purpose |
|-----------|-----|---------|
| **Industrial Ethernet Switch** | USR-ES105 | 3-port GbE between Pi, Jetson, payload |

Enables high-bandwidth communication between companion computers and connected payloads (cameras, sensors).

---

## Connectivity

### Cellular: 4G/5G + eSIM

Each Jetson connects via cellular modem:

| Specification | Value |
|---------------|-------|
| **Modem** | Quectel RM520N-GL |
| **SKU** | RM520NGLAA-M20-SGASA |
| **Bands** | 4G LTE + 5G NR Sub-6GHz |
| **SIM** | eSIM with OTA provisioning |
| **Interface** | M.2 (USB 3.0 / PCIe) |
| **Form Factor** | M.2 NGFF |

**eSIM Benefits:**

- No physical SIM swapping across large fleet
- Remote carrier provisioning
- Carrier switching for coverage optimization
- Centralized subscription management

See [Sensing platform](/platform/sensing/) for details on 5G/LTE connectivity architecture.

---

## Full Bill of Materials

| Category | Component | SKU |
|----------|-----------|-----|
| **Airframe** | Holybro X500 V2 ARF Kit | SKU30125 |
| **Flight Control** | Pixhawk 6X | SKU11073 |
| **Flight Control** | Pixhawk 6X Baseboard | SKU18117 |
| **Flight Control** | Pixhawk 6X Accessories | SKU15011 |
| **Power** | PM06 V2 Power Module | SKU15019 |
| **Navigation** | Holybro M9N GPS | SKU12027 |
| **Navigation** | GPS Mast | SKU12033 |
| **Safety** | Safety Switch & Buzzer | SKU12007 |
| **Companion** | Raspberry Pi 4 8GB | RPI4-MODEL-B-8GB |
| **Storage** | Industrial microSD 128GB | SDSDQED-128G-GNSIN |
| **Companion** | Jetson Orin Nano 8GB | 900-13767-0000-000 |
| **Storage** | NVMe SSD 256GB | WD-SN530-256G |
| **RC** | ExpressLRS EP1 | HM-EP1-2400 |
| **Telemetry** | SiK Telemetry Radio V3 | SKU17012 |
| **Networking** | Industrial Ethernet Switch | USR-ES105 |
| **Power** | LiPo Battery 4S 6000mAh | TAA60004S35 |
| **Mounting** | Pixhawk Vibration Mount | SKU12010 |
| **Cabling** | Holybro Cable Kit | SKU12055 |
| **Connectivity** | Quectel RM520N-GL 5G Modem | RM520NGLAA-M20-SGASA |

---

## Component Sellers

<details>
<summary><strong>Xem tất cả người bán theo linh kiện (57 mục)</strong></summary>

{{< bom-sellers >}}

</details>

---

## Tiếp theo

[Ngăn xếp Phần mềm →]({{< relref "/fleet/software" >}})
