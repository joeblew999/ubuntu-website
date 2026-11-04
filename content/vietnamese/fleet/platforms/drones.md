---
title: "Nền tảng Drone"
meta_title: "Phần cứng & Giao thức cho Đội Drone | Ubuntu Software"
description: "Thông số kỹ thuật phần cứng cấp sản xuất cho đội bay drone: khung máy bay Holybro X500, bộ điều khiển bay Pixhawk 6X, PX4, giao thức MAVLink và máy tính đồng hành NVIDIA Jetson."
image: "/images/robotics.svg"
draft: false
---

## Phần cứng cho vận hành drone quy mô đội bay

Vận hành đội bay đòi hỏi phần cứng cân bằng giữa năng lực và khả năng bảo trì. Mọi thành phần trong cấu hình này đều được lựa chọn dựa trên độ tin cậy trong môi trường sản xuất, khả năng cung ứng linh kiện và khả năng bảo dưỡng tại hiện trường.

---

## Giao thức: MAVLink

Các drone trong đội bay của chúng tôi giao tiếp bằng **MAVLink** (Micro Air Vehicle Link):

| Khía cạnh | Chi tiết |
|--------|---------|
| **Phiên bản** | MAVLink 2.0 |
| **Phương thức truyền** | Serial (UART), UDP, TCP |
| **Tần suất thông điệp** | 1-50Hz tùy theo loại thông điệp |
| **Mã hóa** | Ký số MAVLink 2.0 (tùy chọn) |

**MAVLink cung cấp:**
- Dữ liệu đo xa chuẩn hóa (vị trí, tư thế, pin, cảm biến)
- Giao thức lệnh (kích hoạt, cất cánh, điểm tham chiếu, RTL)
- Quản lý tham số
- Truyền tệp (nhật ký, nhiệm vụ)

[Cổng phương tiện]({{< relref "/fleet/gateway" >}}) chuyển đổi các thông điệp MAVLink thành các subject NATS, cho phép tổng hợp dữ liệu đo xa và phân phối lệnh trên toàn đội bay.

---

## Khung máy bay

### Holybro X500 V2 ARF

**Holybro X500 V2** là nền tảng lý tưởng để triển khai cho đội bay:

| Thông số kỹ thuật | Giá trị |
|---------------|-------|
| **Khoảng cách trục** | 500mm |
| **Trọng lượng khung** | ~410g |
| **Trọng lượng cất cánh tối đa** | ~2kg |
| **Thời gian bay** | 15-20 phút (với tải trọng tiêu chuẩn) |
| **Ngàm động cơ** | 16x16mm / 19x19mm |

**Tại sao chọn X500:**

- **Thiết kế đã được kiểm chứng** — Hàng nghìn chiếc đã được triển khai trên toàn thế giới, kho kiến thức cộng đồng phong phú
- **Linh kiện sẵn có** — Cánh tay, càng đáp và phụ kiện dễ dàng tìm mua
- **Dễ bảo trì** — Cấu trúc dạng mô-đun, có thể bảo dưỡng tại hiện trường
- **Linh hoạt về tải trọng** — Đủ sức chứa cho hai máy tính đồng hành cùng các cảm biến

Bộ kit ARF (Almost Ready to Fly — gần như sẵn sàng bay) bao gồm khung, động cơ, ESC và cánh quạt — giúp giảm sự khác biệt khi lắp ráp trên một đội bay lớn.

---

## Bộ điều khiển bay

### Pixhawk 6X

**Pixhawk 6X** chạy **PX4 v1.14 LTS** đảm nhận tất cả các chức năng quan trọng cho chuyến bay:

| Thông số kỹ thuật | Giá trị |
|---------------|-------|
| **Bộ xử lý** | STM32H753 (480MHz Cortex-M7) |
| **IMU** | Dự phòng ba lớp (ICM-42688-P, ICM-45686, BMI088) |
| **Khí áp kế** | Kép (MS5611, ICP-20100) |
| **Từ kế** | IST8310 |
| **Giao tiếp** | 3x CAN, 6x UART, SPI, I2C, PWM |

**Tại sao chọn Pixhawk 6X:**

- **Dự phòng** — Ba IMU, hai khí áp kế để biểu quyết cảm biến
- **Hỗ trợ PX4 LTS** — Ổn định dài hạn, các bản vá bảo mật
- **Tiêu chuẩn FMUv6X** — Có thể thay thế lẫn nhau với các bo mạch tương thích khác
- **Cách ly rung động** — Tích hợp giảm chấn cho IMU

Bộ điều khiển bay đảm nhận:

- Ước lượng và điều khiển tư thế
- Giữ vị trí và dẫn đường
- Các hành vi an toàn khi sự cố (RTL, hạ cánh, hàng rào địa lý)
- Xử lý tín hiệu đầu vào RC
- Dữ liệu đo xa MAVLink

---

## Máy tính đồng hành

Các phương tiện trong đội bay chạy **hai máy tính đồng hành** với nhiệm vụ riêng biệt:

### Raspberry Pi 4 / CM4 (Máy tính đồng hành cảm biến)

Đảm nhận việc tích hợp cảm biến nhẹ và thu thập dữ liệu:

| Thông số kỹ thuật | Giá trị |
|---------------|-------|
| **Bộ xử lý** | Cortex-A72 bốn nhân @ 1.8GHz |
| **RAM** | 4GB hoặc 8GB |
| **Lưu trữ** | microSD hoặc eMMC 32GB trở lên |
| **Nguồn điện** | 5V @ 3A thông thường |

**Nhiệm vụ:**

- Trình điều khiển cảm biến môi trường (nhiệt độ, độ ẩm, chất lượng không khí)
- Ghi nhật ký dữ liệu GPS/GNSS
- Chụp ảnh từ camera (các tác vụ không dùng AI)
- Định tuyến MAVLink tới Jetson
- Lưu trữ và chuyển tiếp khi Jetson ngoại tuyến

### NVIDIA Jetson (Máy tính đồng hành AI)

Đảm nhận các tác vụ tính toán nặng:

| Mẫu | Nhân GPU | Hiệu năng AI | Công suất |
|-------|-----------|----------------|-------|
| **Jetson Orin Nano** | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin NX** | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 2048 CUDA | 275 TOPS | 15-60W |
| **Jetson Xavier NX** | 384 CUDA | 21 TOPS | 10-20W |

**Nhiệm vụ:**

- Thị giác máy tính (phát hiện đối tượng, theo dõi)
- Đo hành trình quán tính - thị giác
- Lập kế hoạch đường bay và tránh chướng ngại vật
- Cổng phương tiện (NATS client)
- Nút lá NATS (leaf node)

Các đợt triển khai đội bay thường sử dụng **Jetson Orin Nano** để suy luận tiết kiệm chi phí hoặc **Orin NX** cho các tác vụ nhận thức đòi hỏi cao.

---

## Hệ thống vô tuyến

### Điều khiển RC: ExpressLRS

**ExpressLRS** cung cấp liên kết để phi công can thiệp điều khiển:

| Thông số kỹ thuật | Giá trị |
|---------------|-------|
| **Tần số** | 915MHz (Mỹ) / 868MHz (EU) |
| **Tầm hoạt động** | Trên 10km (tùy theo công suất/ăng-ten) |
| **Độ trễ** | <5ms |
| **Giao thức** | Mã nguồn mở |

**Tại sao chọn ExpressLRS:**

- **Mã nguồn mở** — Không phụ thuộc nhà cung cấp, phát triển dựa trên cộng đồng
- **Tầm hoạt động** — Liên kết tin cậy cho các hoạt động trong tầm nhìn trực quan
- **Độ trễ thấp** — Điều khiển thủ công phản hồi nhanh khi cần
- **Chi phí** — Bộ thu giá phải chăng cho quy mô đội bay

### Đo xa: Radio SiK

Radio **SiK 915 MHz** cung cấp dữ liệu đo xa cho trạm mặt đất:

| Thông số kỹ thuật | Giá trị |
|---------------|-------|
| **Tần số** | 915MHz (Mỹ) / 433MHz (EU) |
| **Tầm hoạt động** | 1-2km thông thường |
| **Tốc độ dữ liệu** | Lên đến 250kbps |
| **Giao tiếp** | Serial (UART) |

**Vai trò trong đội bay:**

- Dữ liệu đo xa dự phòng khi không có mạng di động
- Kết nối trạm điều khiển mặt đất
- Kiểm thử và gỡ lỗi tại chỗ

Đối với vận hành đội bay trong môi trường sản xuất, dữ liệu đo xa chủ yếu đi qua kết nối di động của Jetson tới NATS. Radio SiK đóng vai trò dự phòng và phục vụ chẩn đoán tại hiện trường.

---

## Hệ thống nguồn điện

### Pin: LiPo 4S

| Thông số kỹ thuật | Giá trị |
|---------------|-------|
| **Cấu hình** | 4S (14.8V danh định) |
| **Dung lượng** | 5000-6000mAh thông thường |
| **Tốc độ xả** | 30C trở lên |
| **Đầu nối** | XT60 |

**Lưu ý cho đội bay:**

- Chuẩn hóa một cấu hình pin duy nhất để thuận tiện hậu cần
- Sử dụng pin thông minh có dữ liệu đo xa khi có thể
- Triển khai luân chuyển pin và theo dõi tình trạng pin
- Lên kế hoạch tỷ lệ pin trên phương tiện là 3:1 để vận hành liên tục

### Phân phối nguồn điện

- **PDB/PMS** — Bo mạch quản lý nguồn cho các máy tính đồng hành
- **BEC** — Nguồn ổn áp 5V cho Pi, servo
- **Nguồn Jetson** — Lấy trực tiếp từ pin qua bộ ổn áp (12V thông thường)

---

## Kết nối

### Di động: 4G/5G + eSIM

Mỗi Jetson kết nối qua modem di động:

| Thông số kỹ thuật | Giá trị |
|---------------|-------|
| **Modem** | Quectel RM520N-GL hoặc tương đương |
| **Băng tần** | 4G LTE + 5G NR |
| **SIM** | eSIM với cấp phát OTA |
| **Giao tiếp** | USB 3.0 hoặc M.2 |

**Lợi ích của eSIM:**

- Không cần thay SIM vật lý trên một đội bay lớn
- Cấp phát nhà mạng từ xa
- Chuyển đổi nhà mạng để tối ưu vùng phủ sóng
- Quản lý thuê bao tập trung

---

## Mô hình an toàn

Các tính năng an toàn dành riêng cho drone được thực thi ở nhiều cấp độ:

| Cấp độ | Cơ chế | Hành vi |
|-------|-----------|----------|
| **Can thiệp RC** | ExpressLRS | Phi công luôn có quyền điều khiển thủ công |
| **Bộ điều khiển bay** | Cơ chế an toàn PX4 | Thực thi RTL, hạ cánh, hàng rào địa lý |
| **Máy tính đồng hành** | Cổng phương tiện | Xác thực lệnh, kiểm tra an toàn |
| **Mạng** | NATS không bao giờ nằm trong vòng điều khiển | Chỉ giám sát, không ảnh hưởng đến chuyến bay |

Xem [Mô hình an toàn]({{< relref "/fleet/safety" >}}) để biết kiến trúc an toàn đầy đủ.

---

## Tóm tắt

| Thành phần | Lựa chọn | Lý do |
|-----------|-----------|-----------|
| **Khung** | Holybro X500 V2 ARF | Đã kiểm chứng, sẵn có, dễ bảo trì |
| **FCU** | Pixhawk 6X | Dự phòng, hỗ trợ PX4 LTS |
| **Máy tính đồng hành cảm biến** | Raspberry Pi 4/CM4 | Tiết kiệm chi phí, nhiều GPIO |
| **Máy tính đồng hành AI** | Jetson Orin/Xavier | Suy luận trên GPU, hệ sinh thái CUDA |
| **RC** | ExpressLRS | Mã nguồn mở, tầm xa |
| **Đo xa** | SiK 915MHz | Liên kết dự phòng, gỡ lỗi |
| **Nguồn điện** | LiPo 4S 5000mAh | Chuẩn hóa đội bay |
| **Di động** | 4G/5G + eSIM | Cấp phát OTA, vùng phủ sóng |
| **Giao thức** | MAVLink 2.0 | Tiêu chuẩn ngành, tích hợp sẵn trong PX4 |

---

## Nơi mua

Chúng tôi không bán phần cứng — các đối tác của chúng tôi bán. **Các cửa hàng drone được chứng nhận** có thể cung cấp bộ kit hoàn chỉnh kèm theo hạ tầng đội bay.

**Dành cho người mua drone:**
Mua từ một [đối tác được chứng nhận]({{< relref "/partners" >}}) và nhận được:
- Phần cứng được cấu hình sẵn, sẵn sàng kết nối
- Thông tin xác thực cho hạ tầng được quản lý của chúng tôi
- Hỗ trợ tại chỗ từ cửa hàng của bạn

**Dành cho cửa hàng drone:**
[Tham gia chương trình đối tác của chúng tôi]({{< relref "/partners/become-partner" >}}) để cung cấp hạ tầng đội bay cùng với việc bán phần cứng của bạn.

---

## Tài liệu liên quan

- [Nền tảng được hỗ trợ]({{< relref "/fleet/platforms" >}}) — Tổng quan về tất cả các loại phương tiện
- [Phương tiện mặt đất]({{< relref "/fleet/platforms/ground" >}}) — Ô tô, xe tải và AGV
- [Cổng phương tiện]({{< relref "/fleet/gateway" >}}) — Cầu nối MAVLink sang NATS
- [Mô hình an toàn]({{< relref "/fleet/safety" >}}) — Kiến trúc an toàn khi sự cố
- [Ngăn xếp phần mềm]({{< relref "/fleet/software" >}}) — Hệ điều hành và phần mềm trung gian
