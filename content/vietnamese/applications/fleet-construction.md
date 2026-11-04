---
title: "Xây dựng & Hạ tầng"
meta_title: "Đội thiết bị IoT cho Xây dựng | Ubuntu Software"
description: "Giám sát công trường xây dựng với các đội thiết bị IoT: drone phục vụ khảo sát, máy đào và máy xúc lật tích hợp telematics, cảm biến công trường đảm bảo an toàn và tuân thủ quy định về môi trường. Tất cả được quản lý bằng một kiến trúc thống nhất."
image: "/images/robotics.svg"
draft: false
---

## Quan sát mọi công trường, mỗi ngày

Các dự án xây dựng thất bại khi khả năng quan sát thất bại. Các đội thiết bị IoT mang lại nhận thức liên tục—theo dõi tiến độ, ghi nhận hiện trạng và phát hiện vấn đề trước khi chúng trở nên tốn kém.

---

## Các loại đội thiết bị trong xây dựng

| Loại đội thiết bị | Thiết bị | Vai trò |
|------------|---------|------|
| **Trên không** | Drone khảo sát, drone kiểm tra | Theo dõi tiến độ, lập bản đồ công trường, kiểm tra kết cấu |
| **Phương tiện mặt đất** | Máy đào, máy xúc lật, xe tải chở vật liệu | Theo dõi thiết bị, viễn trắc (telematics), giám sát mức độ sử dụng |
| **IoT cố định** | Cảm biến công trường, kiểm soát ra vào, thiết bị giám sát máy móc | Tuân thủ an toàn, giám sát môi trường, an ninh |

Tất cả các loại thiết bị đều dùng chung một kiến trúc: AI tại biên, nhắn tin qua NATS JetStream và quản lý đội thiết bị thống nhất.

---

## Thách thức

Ngành xây dựng và hạ tầng đối mặt với những lỗ hổng về khả năng quan sát:

- **Theo dõi tiến độ** — Báo cáo thủ công bị chậm trễ và mang tính chủ quan
- **Tuân thủ an toàn** — Thanh tra viên không thể có mặt ở mọi nơi
- **Lưu trữ hồ sơ** — Ảnh chụp không ghi lại bối cảnh hoặc dấu thời gian một cách đáng tin cậy
- **Phát hiện thay đổi** — Nhận ra công việc trái phép hoặc sai lệch so với bản vẽ
- **Kiểm tra tài sản** — Cầu, tháp, đường ống trải dài trên những khoảng cách rất lớn

Các phương pháp truyền thống đòi hỏi con người leo lên công trình, lái xe dọc tuyến hoặc dựa vào các đợt kiểm tra theo lịch. Vấn đề ẩn nấp giữa các lần kiểm tra.

---

## Cách các đội thiết bị IoT giải quyết vấn đề

### Trên không: Giám sát công trường liên tục

Các chuyến bay drone tự động hằng ngày ghi lại:

- **Bản đồ trực ảnh (orthomosaic)** — Hình ảnh công trường chính xác đến từng centimet
- **Đám mây điểm 3D** — Dữ liệu thể tích để theo dõi công tác đất
- **Ảnh nhiệt** — Phát hiện độ ẩm, vấn đề cách nhiệt, nhiệt tỏa ra từ thiết bị
- **Hồ sơ 360°** — Bản ghi hình ảnh đầy đủ kèm dấu thời gian

### Mặt đất: Trí tuệ thiết bị

Dữ liệu viễn trắc từ mọi máy móc trên công trường:

- **Theo dõi mức độ sử dụng** — Số giờ hoạt động, thời gian chạy không tải, mức tiêu thụ nhiên liệu
- **Nhận biết vị trí** — Vị trí theo thời gian thực của toàn bộ thiết bị
- **Cảnh báo bảo trì** — Mã lỗi động cơ, chu kỳ bảo dưỡng, phát hiện sự cố
- **Hàng rào địa lý (geofencing)** — Tuân thủ ranh giới hoạt động của thiết bị, chống trộm

### IoT cố định: Nhận thức công trường

Cảm biến khắp dự án:

- **Thiết bị giám sát chất lượng không khí** — Bụi, khí thải, theo dõi tuân thủ
- **Cảm biến tiếng ồn** — Tuân thủ quy định, quan hệ với cư dân lân cận
- **Kiểm soát ra vào** — Ai có mặt trên công trường, khi nào, trong bao lâu
- **Thiết bị giám sát kết cấu** — Phát hiện lún, rung động, dịch chuyển

### Kiểm tra hạ tầng

Kiểm tra có hệ thống mà không cần giàn giáo hay phong tỏa làn đường:

- **Kiểm tra cầu** — Lập hồ sơ mặt cầu, kết cấu phần trên, kết cấu phần dưới
- **Kiểm tra tháp** — Viễn thông, truyền tải điện, tua-bin gió
- **Tuần tra đường ống** — Giám sát hành lang tuyến, phát hiện rò rỉ
- **Vỏ bao che công trình** — Kiểm tra mặt dựng, đánh giá tình trạng mái

---

## AI tích hợp trên thiết bị (Xử lý tại biên)

Suy luận thời gian thực trên NVIDIA Jetson để có thông tin chuyên sâu tức thì:

| Khả năng | Ứng dụng trong xây dựng |
|------------|-------------------------|
| **Phát hiện tiến độ** | Nhận diện các hạng mục đã hoàn thành so với kế hoạch |
| **Vi phạm an toàn** | Xác định thiếu PPE, rào chắn không đúng quy cách, nguy cơ té ngã |
| **Phát hiện vết nứt** | Khuyết tật bề mặt trên bê tông, nhựa đường, thép |
| **Nhận diện ăn mòn** | Dấu hiệu gỉ sét trên cầu, bồn chứa, thép kết cấu |
| **Bất thường nhiệt** | Điểm nóng cho thấy sự cố điện, thấm ẩm |

**Vì sao xử lý tại biên quan trọng:**
- Gắn cờ các vấn đề an toàn theo thời gian thực để phản ứng ngay lập tức
- Xử lý các tập dữ liệu trực ảnh lớn ngay tại chỗ
- Hoạt động tại các công trường xa xôi có kết nối hạn chế
- Giảm chi phí dữ liệu cho các chuyến bay thường xuyên

---

## AI đám mây (Phân tích đội thiết bị)

Thông qua NATS JetStream, tổng hợp thông tin chuyên sâu trên nhiều dự án:

| Khả năng | Ứng dụng trong xây dựng |
|------------|-------------------------|
| **So sánh BIM** | Chồng đám mây điểm hiện trạng lên mô hình thiết kế |
| **Phân tích tiến độ** | Phần trăm hoàn thành, sai lệch tiến độ, phân tích xu hướng |
| **Theo dõi khuyết tật** | Ghi nhận sự cố, theo dõi khắc phục, xác minh sửa chữa |
| **Lập lịch đội thiết bị** | Điều phối drone qua nhiều công trường một cách hiệu quả |
| **Báo cáo tuân thủ** | Tạo báo cáo kiểm tra, hồ sơ an toàn |

**Vì sao đám mây quan trọng:**
- So sánh hiện trạng với ý đồ thiết kế
- Theo dõi tiến độ trên toàn bộ danh mục dự án
- Phân tích dữ liệu lịch sử phục vụ khiếu nại và tranh chấp
- Tích hợp với hệ thống quản lý dự án

---

## Cấu hình phần cứng

Đối với các triển khai trong xây dựng, chúng tôi khuyến nghị:

| Thành phần | Lựa chọn | Lý do |
|-----------|-----------|-----------|
| **Khung máy bay** | Holybro X500 | Nền tảng ổn định cho tải trọng lập bản đồ |
| **Bộ điều khiển bay** | Pixhawk 6X | GPS RTK cho độ chính xác cấp khảo sát |
| **Máy tính đồng hành cảm biến** | Raspberry Pi CM4 | Ghi dữ liệu GPS/IMU |
| **Máy tính đồng hành AI** | Jetson Orin NX | Năng lực tính toán cao hơn cho xử lý 3D |
| **Tải trọng** | Camera RGB + camera nhiệt | Lập hồ sơ hình ảnh và nhiệt |
| **Kết nối** | 4G/LTE + eSIM | Phủ sóng công trường đô thị và vùng xa |

[Thông số phần cứng đầy đủ →]({{< relref "/fleet/hardware" >}})

---

## Trường hợp sử dụng

### Công trường đang thi công
Lập hồ sơ tiến độ hằng ngày. Theo dõi khối lượng công tác đất. Giám sát an toàn.

### Kiểm tra cầu
Hồ sơ tuân thủ NBIS. Bao quát mặt cầu, kết cấu phần trên, kết cấu phần dưới.

### Hạ tầng tuyến tính
Đường ống, đường dây điện, đường cao tốc. Tuần tra và kiểm tra ở quy mô lớn.

### Kiểm tra tòa nhà
Đánh giá tình trạng mặt dựng. Khảo sát mái. Phân tích nhiệt vỏ bao che.

### Hoạt động khai thác mỏ
Đo lường bãi chứa vật liệu. Tiến độ khai thác moong. Giám sát đường vận chuyển.

---

## Tích hợp

Kết nối dữ liệu đội drone với các hệ thống hiện có của bạn:

- **Nền tảng BIM** — Autodesk Construction Cloud, Bentley, Trimble
- **Quản lý dự án** — Procore, PlanGrid, Oracle Aconex
- **Hệ thống GIS** — Esri ArcGIS, QGIS
- **Quản lý tài sản** — IBM Maximo, SAP PM

---

## Tìm hiểu kỹ thuật chuyên sâu

Ứng dụng này được xây dựng trên kiến trúc đội thiết bị cấp sản xuất của chúng tôi:

| Thành phần | Vai trò trong xây dựng |
|-----------|---------------------|
| [Cấu trúc liên kết NATS]({{< relref "/fleet/nats-topology" >}}) | Đồng bộ dữ liệu từ văn phòng container tại công trường, xử lý tại các trung tâm khu vực |
| [Cổng phương tiện]({{< relref "/fleet/gateway" >}}) | Tiền xử lý dữ liệu đo ảnh (photogrammetry) ngay trên drone |
| [Luồng JetStream]({{< relref "/fleet/streams" >}}) | Lưu trữ hình ảnh lịch sử, dữ liệu nền để phát hiện thay đổi |
| [Mô hình an toàn]({{< relref "/fleet/safety" >}}) | Hàng rào địa lý quanh khu vực đang thi công, cần cẩu, nhân sự |

[Khám phá toàn bộ kiến trúc →]({{< relref "/fleet" >}})

---

## Bắt đầu

Mọi dự án đều được hưởng lợi từ khả năng quan sát tốt hơn. Dù bạn quản lý một công trường hay hàng trăm công trường, các đội drone đều mang đến hồ sơ bạn cần.

[Liên hệ với chúng tôi →](/contact)
