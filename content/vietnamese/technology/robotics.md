---
title: "Robotics"
meta_title: "Hệ thống robot | Ubuntu Software"
description: "Kiến trúc robotics của chúng tôi được xây dựng trên nền tảng Viam RDK - từ thiết kế đến triển khai với khả năng trừu tượng hóa phần cứng, thị giác máy tính và tích hợp công nghiệp."
image: "/images/spatial.svg"
draft: false
---

## Robotics Stack

Từ thiết kế đến triển khai. Spatial cung cấp lớp thiết kế và mô phỏng 3D. Viam RDK cung cấp môi trường chạy để điều khiển và vận hành robot thực tế.

---

## Hệ điều hành Linux nhúng

Robotics hoạt động trên hệ điều hành Linux. Với hàng thập kỷ kinh nghiệm trong lĩnh vực Linux và hệ thống nhúng, chúng tôi thiết kế các giải pháp phù hợp với môi trường thực tế - từ nhà máy sản xuất, môi trường ngoài trời cho đến các thiết bị biên có tài nguyên hạn chế.

### Viam RDK Thời gian chạy

**Tại sao chúng tôi chọn Viam RDK:**

- **Nguồn mở** — Phù hợp với triết lý tiêu chuẩn mở của chúng tôi
- **Go-native** — Phù hợp với bộ công nghệ của chúng tôi
- **Modular** — Các thành phần có thể được thay thế và mở rộng
- **Tùy chọn đám mây* — Hoạt động ưu tiên chế độ ngoại tuyến, giống như các sản phẩm của chúng tôi
- **Không phụ thuộc vào phần cứng** — Không bị giới hạn bởi các nhà sản xuất robot cụ thể
- **Hoạt động trên hệ điều hành Linux nhúng** — Raspberry Pi, NVIDIA Jetson, bộ điều khiển công nghiệp

Viam RDK là bộ công cụ phát triển robot mã nguồn mở với các SDK không phụ thuộc vào ngôn ngữ (Go, Python, TypeScript) và kiến trúc thành phần mô-đun.

[Tài liệu Viam RDK →](https://docs.viam.com/)

### Nền tảng nhúng

| Platform | Use Case |
|----------|----------|
| **Raspberry Pi** | Prototyping, edge compute, kiosks |
| **NVIDIA Jetson** | GPU-accelerated vision and ML |
| **Industrial Linux** | Factory automation, harsh environments |
| **Custom ARM boards** | Application-specific deployments |

Cùng một mã nguồn Go. Cùng một Viam RDK. Triển khai từ máy trạm phát triển sang thiết bị biên nhúng.

---

## Khả năng

| Capability | Viam Service |
|------------|--------------|
| Motion planning with collision avoidance | motion service |
| Object detection and segmentation | vision service + ML models |
| Point cloud from depth sensors | Camera component |
| Hand-eye calibration | Frame system |
| SLAM with IMU | slam service |

---

## Tầm nhìn & Nhận thức

### Trừu tượng hóa camera

Viam cung cấp một giao diện lập trình ứng dụng (API) thống nhất `rdk:component:camera` hoạt động trên nhiều loại camera khác nhau—webcam, camera IP, LiDAR và camera độ sâu như Intel RealSense D435i.

**Cách thức hoạt động:**

- **Giao diện API chuẩn hóa** — Mã nguồn của bạn tương tác với giao diện camera Viam, không phải với phần cứng
- **Driver tích hợp** — Module `viam-camera-realsense` quản lý tích hợp RealSense
- **Không phụ thuộc vào phần cứng** — Thay đổi camera mà không cần thay đổi logic ứng dụng
- **Dựa trên cấu hình* — Độ phân giải, cảm biến và luồng được cấu hình thông qua JSON, không phải mã nguồn

Các phương thức có sẵn: `GetImage()`, `GetImages()`, `GetPointCloud()`, cùng với các thông số nội tại của camera.

### Intel RealSense D435i

Cảm biến độ sâu RGB-D cho việc tạo đám mây điểm và nhận thức không gian.

[librealsense trên GitHub →](https://github.com/realsenseai/librealsense)

### Dòng công việc Machine Learning

- Phát hiện và phân đoạn đối tượng YOLOv8
- Tích hợp với dịch vụ Viam Vision
- Xử lý suy luận thời gian thực trên các thiết bị biên

---

## Trừu tượng hóa phần cứng

**Giá trị cốt lõi: Thay đổi cấu hình, không phải thay đổi mã nguồn.**

- **Cùng phần mềm, cánh tay lớn hơn** — Thay thế cấu hình xArm cho UR5e hoặc KUKA, triển khai lại
- **xArm như thiết bị điều khiển giảng dạy** — Thực hiện các động tác trên xArm, cánh tay lớn hơn mô phỏng
- **Tổ chức phối hợp nhiều robot* — xArm xử lý các bộ phận nhỏ, KUKA xử lý các tác vụ nâng hạ nặng
- **Phát triển mô hình số** — Phát triển trên nền tảng xArm, mô phỏng dựa trên cơ học KUKA

---

## Vũ khí được hỗ trợ

| Model | Payload | Reach | Use Case |
|-------|---------|-------|----------|
| xArm 6 | 5kg | 700mm | Development |
| UR5e | 5kg | 850mm | Production |
| KUKA KR6 | 6kg | 900mm | Small parts |
| KUKA KR10 | 10kg | 900-1100mm | Medium assembly |
| KUKA KR16 | 16kg | 1600mm | Welding, palletizing |
| KUKA KR30 | 30kg | 2000mm+ | Heavy handling |

[module viam-kuka →](https://github.com/viam-soleng/viam-kuka)

---

## Kiến trúc

```
SOFTWARE STACK
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Vision   │  │ Motion   │  │   ML     │  │ Business │
│ Pipeline │  │ Planning │  │ Models   │  │  Logic   │
└────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
     └─────────────┴─────────────┴─────────────┘
                         │
                  Viam Arm API
              (rdk:component:arm)
                         │
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐
    │  xArm 6  │   │  UR5e    │   │  KUKA    │
    │  (dev)   │   │  (prod)  │   │  (heavy) │
    │ 5kg/700mm│   │ 5kg/850mm│   │ 30kg+    │
    └──────────┘   └──────────┘   └──────────┘
```

---

## Tích hợp công nghiệp

**Tích hợp I/O:**

Modbus → PLC cho tích hợp tự động hóa công nghiệp. Điều khiển băng tải, cảm biến, bộ truyền động và hệ thống an toàn từ cùng một môi trường chạy.

---

## Cách các thành phần kết hợp với nhau

| Layer | Technology | Purpose |
|-------|------------|---------|
| Design | Spatial | 3D work cell design and simulation |
| Runtime | Viam RDK | Robot control and operation |
| Vision | RealSense + YOLOv8 | Perception and object detection |
| Arms | xArm / UR / KUKA | Physical manipulation |
| I/O | Modbus / PLC | Industrial integration |

---

## Tìm hiểu thêm

- [Tổng quan về công nghệ →](/vi/technology/)
- [Nền tảng không gian →](/vi/platform/spatial/)
- [Liên hệ với chúng tôi →](/vi/contact/)
