---
title: "Cảm biến"
meta_title: "Nền tảng Cảm biến & Nhận thức | Phần mềm Ubuntu"
description: "Tích hợp đa cảm biến cho trí tuệ không gian—LiDAR, camera và cảm biến công nghiệp được hợp nhất thông qua một tác tử biên (edge agent) duy nhất với kết nối 5G/eSIM."
image: "/images/spatial.svg"
draft: false
---

## Cảm biến không gian

Nhận thức thế giới thực cho bản sao số, robot và hệ thống tự hành. Kết nối LiDAR, camera và cảm biến công nghiệp với các mô hình không gian của bạn.

---

## Vấn đề

Cảm biến tạo ra dữ liệu. Để hiểu được dữ liệu đó, cần có:

- **Ngữ cảnh** — Cảm biến đặt ở đâu? Nó đang quan sát gì?
- **Hợp nhất** — Kết hợp nhiều luồng dữ liệu cảm biến thành một bức tranh thống nhất
- **Tích hợp** — Kết nối với các công cụ thiết kế, không chỉ với bảng điều khiển

Hầu hết các giải pháp cảm biến chỉ dừng lại ở việc thu thập dữ liệu. Chúng tôi kết nối cảm biến với các mô hình không gian.

---

## Chế độ triển khai

Cùng một tác tử biên. Cùng các cảm biến. Chỉ khác cấu hình.

| Chế độ | Nền tảng | Trường hợp sử dụng |
|------|----------|----------|
| **Trên không** | Drone doanh nghiệp DJI | Khảo sát, kiểm tra, lập bản đồ |
| **Mặt đất** | Chân máy, ba lô | Quét nội thất, xây dựng |
| **Robot** | Viam RDK, ROS2 | Điều hướng, gắp và đặt |
| **Cố định** | Lắp đặt cố định | Giao thông, an ninh, kho bãi |

### Điều hướng không cần GPS

Đo hành trình quán tính - thị giác (visual-inertial odometry) cho phép bay tự hành mà không cần GPS—điều thiết yếu cho môi trường trong nhà, hẻm núi đô thị và các khu vực không có GPS.

{{< video src="videos/visual-inertial-nav-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*Theo dõi vị trí thời gian thực chỉ bằng camera và IMU. Không cần cơ sở hạ tầng bên ngoài.*

{{< video src="videos/gps-free-vio-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*Bay tự hành vòng kín sử dụng đo hành trình quán tính - thị giác.*

Để biết thông số kỹ thuật phần cứng chi tiết, xem [Phần cứng đội drone →]({{< relref "/fleet/hardware" >}})

---

## Trừu tượng hóa cảm biến

**Không phụ thuộc phần cứng ngay từ khâu thiết kế.** Mã của bạn giao tiếp với API hợp nhất của chúng tôi, không phải với trình điều khiển của từng cảm biến.

### Các loại cảm biến được hỗ trợ

| Loại | Ví dụ |
|------|----------|
| **LiDAR** | [Livox Mid-360](https://www.livoxtech.com/mid-360) ([SDK](https://github.com/Livox-SDK/Livox-SDK2)), [Avia](https://www.livoxtech.com/avia) ([bản demo](https://terra-1-g.djicdn.com/65c028cd298f4669a7f0e40e50ba1131/Download/Avia/02_Intro%20Demo.mp4)) |
| **Camera RGB-D** | [Intel RealSense](https://www.intelrealsense.com/), [Luxonis OAK-D](https://shop.luxonis.com/products/oak-d) |
| **Định vị** | [u-blox RTK](https://www.u-blox.com/en/technologies/rtk-real-time-kinematic) GPS/GNSS, IMU |
| **Công nghiệp** | Cảm biến Modbus, bus CAN |

Thay cảm biến mà không cần sửa mã. Vận hành theo cấu hình, không phải theo mã.

---

## Kiến trúc tác tử biên

Tệp nhị phân Go chạy trên phần cứng của bạn—Raspberry Pi, Jetson, Linux công nghiệp hoặc ARM tùy chỉnh.

| Khả năng | Mô tả |
|------------|-------------|
| **Hệ thống plugin** | Thêm cảm biến qua cấu hình, không cần thay đổi mã |
| **Bộ đệm cục bộ** | Lưu trữ và chuyển tiếp khi ngoại tuyến |
| **Truyền phát thời gian thực** | NATS JetStream lên đám mây |
| **Gọn nhẹ** | Một tệp nhị phân duy nhất, không phụ thuộc runtime |

---

## Kết nối

### 5G/LTE với eSIM OTA

Không cần thay SIM. Không cần quét mã QR. Cấp phát được đẩy từ máy chủ.

- Modem được xuất xưởng kèm hồ sơ khởi tạo (bootstrap profile)
- Nền tảng của bạn kích hoạt việc tải xuống hồ sơ nhà mạng
- Chuyển đổi nhà mạng ngay giữa quá trình triển khai qua API

Hoạt động với drone đang bay, robot đang di chuyển và các hệ thống lắp đặt cố định ở vùng xa xôi.

---

## Tích hợp với Spatial

Cảm biến đưa dữ liệu trực tiếp vào các mô hình 3D của bạn:

| Luồng dữ liệu | Mục đích |
|-----------|---------|
| Đám mây điểm → Mô hình không gian | Số hóa hiện trạng |
| GPS/IMU → Định vị mô hình | Tham chiếu địa lý |
| Cảm biến môi trường → Bản sao số | Cập nhật trạng thái trực tiếp |
| I/O công nghiệp → Tự động hóa | Điều khiển vòng kín |

Không chỉ là bảng điều khiển. Cảm biến đặt trong ngữ cảnh.

---

## Xây dựng trên Foundation

Cảm biến tự động kế thừa mọi khả năng của [Foundation](/platform/foundation/):

| Khả năng | Ý nghĩa |
|------------|---------------|
| **Ưu tiên ngoại tuyến** | Thu thập dữ liệu không cần internet, đồng bộ khi có kết nối |
| **Triển khai mọi nơi** | Biên, di động, máy tính để bàn—cùng một tác tử |
| **Tự chủ** | Cảm biến của bạn, dữ liệu của bạn, máy chủ của bạn |
| **Đồng bộ thời gian thực** | Truyền phát đến nhiều đích cùng lúc |

[Tìm hiểu thêm về Foundation →](/platform/foundation/)

---

## Một phần của bức tranh lớn hơn

Cảm biến là lớp nhận thức của nền tảng Ubuntu Software.

Đối với các tổ chức cần thiết kế 3D và AI, nền tảng Spatial của chúng tôi cung cấp môi trường thiết kế—được tích hợp trực tiếp với dữ liệu cảm biến của bạn.

[Khám phá Spatial →](/platform/spatial/)

---

## Tham chiếu triển khai thực tế: Đội 1.000 drone

Xem cách các khả năng cảm biến này mở rộng lên quy mô sản xuất. Kiến trúc tham chiếu của chúng tôi minh họa việc triển khai trọn vẹn một đội bay sử dụng nền tảng phần cứng và phần mềm nêu trên—1.000 drone Holybro X500 chạy PX4, hai máy tính đồng hành và NATS JetStream cho bản sao số ở quy mô đội bay.

[Khám phá kiến trúc đội drone →]({{< relref "/fleet" >}})

---

## Cùng xây dựng với chúng tôi

Bạn đang triển khai cảm biến? Đang xây dựng hệ thống nhận thức? Hãy trao đổi với chúng tôi.

[Liên hệ →](/contact)
