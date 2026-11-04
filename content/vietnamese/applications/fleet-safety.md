---
title: "An toàn Công cộng & Khẩn cấp"
meta_title: "Đội thiết bị IoT cho An toàn Công cộng | Ubuntu Software"
description: "An toàn công cộng với đội thiết bị IoT: drone giám sát trên không, phương tiện khẩn cấp với điều phối thời gian thực, và cảm biến tòa nhà để phát hiện mối đe dọa. Kiến trúc hợp nhất cho ứng phó phối hợp."
image: "/images/robotics.svg"
draft: false
---

## Khi từng phút đều quan trọng

Ứng phó khẩn cấp phụ thuộc vào thông tin. Đội thiết bị IoT mang lại nhận thức tình huống đúng lúc lực lượng ứng cứu cần nhất: triển khai nhanh hơn, bao phủ rộng hơn, vận hành an toàn hơn.

---

## Các loại đội thiết bị trong an toàn công cộng

| Loại đội thiết bị | Thiết bị | Vai trò |
|------------|---------|------|
| **Trên không** | Drone giám sát, drone tìm kiếm | Tìm kiếm cứu nạn, giám sát sự cố từ trên cao, lập bản đồ đám cháy |
| **Phương tiện mặt đất** | Xe khẩn cấp, xe tuần tra | Điều phối thời gian thực, viễn trắc (telematics), theo dõi tài sản |
| **IoT cố định** | Báo động tòa nhà, cảm biến giao thông, thiết bị phát hiện tiếng súng | Phát hiện mối đe dọa, nhận thức tình huống, thu thập bằng chứng |

Tất cả các loại thiết bị dùng chung một kiến trúc gồm AI tại biên, nhắn tin NATS JetStream và quản lý đội thiết bị hợp nhất.

---

## Thách thức

Các tình huống khẩn cấp thường gặp khoảng trống thông tin:

- **Tìm kiếm cứu nạn**: Khu vực cần bao phủ rộng lớn, nhân lực hạn chế
- **Đánh giá thảm họa**: Cơ sở hạ tầng bị hư hại, lối tiếp cận bị chặn
- **Sự cố đang diễn ra**: Cần quan sát hiện trường trước khi lực lượng đến nơi
- **Phối hợp**: Nhiều cơ quan, liên lạc rời rạc
- **Lập hồ sơ**: Bảo toàn bằng chứng, tái dựng sự cố

Đánh giá từ mặt đất thì chậm và đặt lực lượng ứng cứu vào nguy hiểm. Trực thăng thì đắt đỏ và số lượng có hạn. Những phút đầu tiên thường quyết định kết quả.

---

## Đội thiết bị IoT giải quyết vấn đề này như thế nào

### Trên không: Triển khai nhanh

Drone bố trí sẵn cất cánh ngay lập tức:

- **Điều phối tự động**: Drone đã lên đường trong khi đội ứng cứu còn đang huy động
- **Có mặt đầu tiên tại hiện trường**: Đánh giá bằng hình ảnh trước khi lực lượng đến nơi
- **Phát trực tiếp**: Sở chỉ huy thấy những gì drone thấy
- **Hoạt động ban đêm**: Ảnh nhiệt trong bóng tối hoàn toàn

### Mặt đất: Trí tuệ phương tiện

Phương tiện khẩn cấp phối hợp theo thời gian thực:

- **Điều phối tự động**: Chọn đơn vị gần nhất còn sẵn sàng dựa trên vị trí trực tiếp
- **Tối ưu hóa lộ trình**: Giao thông thời gian thực, đường nhanh nhất đến hiện trường
- **Viễn trắc**: Trạng thái xe, nhiên liệu, cảnh báo bảo dưỡng
- **Theo dõi nguồn lực**: Ai có thiết bị gì, ở đâu

### IoT cố định: Phát hiện mối đe dọa

Cảm biến mang lại nhận thức liên tục:

- **Phát hiện tiếng súng**: Cảm biến âm thanh xác định chính xác vị trí
- **Cảm biến giao thông**: Ùn tắc, tai nạn, ưu tiên tín hiệu
- **Báo động tòa nhà**: Hỏa hoạn, đột nhập, sự cố thang máy
- **Trạm thời tiết**: Điều kiện tại chỗ ảnh hưởng đến công tác ứng phó

### Bao phủ khu vực

Các mô hình tìm kiếm có hệ thống giúp bao phủ địa bàn nhanh chóng:

- **Tìm kiếm theo lưới**: Đường bay tự động, không bỏ sót
- **Phối hợp nhiều drone**: Bao phủ nhiều khu vực hơn cùng lúc
- **Giám sát liên tục**: Hàng giờ bao phủ cho mỗi lần thay pin
- **Hoạt động mọi thời tiết**: Gió, mưa, lạnh, trong giới hạn cho phép

---

## AI trên thiết bị (Xử lý tại biên)

Suy luận thời gian thực trên NVIDIA Jetson để phát hiện tức thì:

| Khả năng | Ứng dụng trong an toàn công cộng |
|------------|--------------------------|
| **Phát hiện người** | Định vị đối tượng ở vùng hoang dã, trong đống đổ nát sau thảm họa, giữa đám đông |
| **Nhận dạng phương tiện** | Theo dõi phương tiện cần chú ý, xác định xe bị bỏ lại |
| **Ảnh nhiệt** | Tìm dấu hiệu nhiệt xuyên qua khói, vào ban đêm, dưới tán lá |
| **Phát hiện đám cháy** | Xác định ngọn lửa đang cháy, điểm nóng, hướng lan của đám cháy |
| **Đánh giá kết cấu** | Phát hiện hư hại công trình, nguy cơ sập đổ, lối tiếp cận |

**Vì sao xử lý tại biên quan trọng:**
- Cảnh báo lực lượng ứng cứu ngay khi phát hiện
- Hoạt động ở khu vực không có sóng di động
- Xử lý đồng thời ảnh nhiệt và ảnh quang học
- Tiếp tục hoạt động khi mạng bị nghẽn

---

## AI đám mây (Phân tích đội thiết bị)

Thông qua NATS JetStream, phối hợp ứng phó giữa các cơ quan:

| Khả năng | Ứng dụng trong an toàn công cộng |
|------------|--------------------------|
| **Phối hợp đa cơ quan** | Chia sẻ luồng dữ liệu giữa cứu hỏa, cảnh sát, cấp cứu (EMS), tìm kiếm cứu nạn (SAR) |
| **Điều phối nguồn lực** | Định tuyến tối ưu dựa trên điều kiện thời gian thực |
| **Phân tích mô hình** | Theo dõi diễn biến sự cố, dự đoán mức độ lan rộng |
| **Quản lý bằng chứng** | Chuỗi lưu giữ bằng chứng, xác minh dấu thời gian |
| **Đánh giá sau hành động** | Tái dựng sự cố, tài liệu huấn luyện |

**Vì sao đám mây quan trọng:**
- Góc nhìn hợp nhất cho trung tâm chỉ huy
- Phối hợp giữa các khu vực tài phán
- Dữ liệu lịch sử từ các sự cố tương tự
- Tích hợp với hệ thống CAD/điều phối

---

## Cấu hình phần cứng

Đối với các triển khai an toàn công cộng, chúng tôi khuyến nghị:

| Thành phần | Lựa chọn | Lý do |
|-----------|-----------|-----------|
| **Khung máy bay** | Holybro X500 hoặc biến thể gia cường | Đáng tin cậy trong điều kiện khắc nghiệt |
| **Bộ điều khiển bay** | Pixhawk 6X | Khả năng tự hành đã được kiểm chứng, có chế độ an toàn khi lỗi |
| **Máy tính đồng hành cảm biến** | Raspberry Pi CM4 | Đo xa (telemetry), cảm biến cơ bản |
| **Máy tính đồng hành AI** | Jetson Orin NX | Hiệu năng cao cho xử lý ảnh nhiệt + ảnh quang học |
| **Tải trọng** | Camera kép RGB + nhiệt | Hoạt động cả ngày lẫn đêm |
| **Kết nối** | LTE + vô tuyến mesh | Vẫn liên lạc được khi trạm phát sóng ngừng hoạt động |

[Thông số phần cứng đầy đủ →]({{< relref "/fleet/hardware" >}})

---

## Trường hợp sử dụng

### Tìm kiếm cứu nạn
Tìm kiếm người mất tích. Tìm kiếm cứu nạn ở vùng hoang dã. Hỗ trợ cứu hộ dưới nước.

### Ứng phó hỏa hoạn
Đánh giá cháy công trình. Lập bản đồ cháy rừng. Phát hiện điểm nóng.

### Thực thi pháp luật
Giám sát sự cố đang diễn ra từ trên cao. An ninh sự kiện. Quản lý sự cố giao thông.

### Ứng phó thảm họa
Đánh giá thiệt hại. Định vị người sống sót. Hướng dẫn bố trí điểm tập kết nguồn lực.

### Biên giới và hàng hải
Tuần tra ven biển. An ninh cảng. Giám sát biên giới.

---

## Mô hình triển khai

### Đội thiết bị bố trí sẵn
Drone được bố trí tại trạm cứu hỏa, đồn cảnh sát, căn cứ tìm kiếm cứu nạn. Cất cánh khi có lệnh điều phối.

### Chỉ huy di động
Hệ thống gắn trên phương tiện. Triển khai ở bất kỳ đâu. Thiết lập nhanh chóng.

### Bao phủ liên tục
Nhiệm vụ thời gian dài. Trạm thay pin hoặc trạm sạc. Hoạt động liên tục.

---

## Đi sâu vào kỹ thuật

Ứng dụng này được xây dựng trên kiến trúc đội thiết bị đạt chuẩn vận hành thực tế của chúng tôi:

| Thành phần | Vai trò trong an toàn công cộng |
|-----------|----------------------|
| [Cấu trúc liên kết NATS]({{< relref "/fleet/nats-topology" >}}) | Nhắn tin bền bỉ khi cơ sở hạ tầng bị hư hại |
| [Cổng phương tiện]({{< relref "/fleet/gateway" >}}) | Phát hiện thời gian thực kèm cảnh báo tức thì |
| [Luồng JetStream]({{< relref "/fleet/streams" >}}) | Lưu trữ video đạt chuẩn bằng chứng, có dấu thời gian |
| [Mô hình an toàn]({{< relref "/fleet/safety" >}}) | Hoạt động gần con người, công trình, tàu bay |

[Khám phá toàn bộ kiến trúc →]({{< relref "/fleet" >}})

---

## Tuân thủ

Hoạt động an toàn công cộng yêu cầu:

- **FAA Part 107**: hỗ trợ xin miễn trừ cho hoạt động ban đêm và BVLOS
- **Tuân thủ CJIS** khi xử lý dữ liệu thực thi pháp luật
- **Chuỗi lưu giữ bằng chứng** với xác minh bằng mật mã
- **Chia sẻ dữ liệu đa cơ quan** với kiểm soát truy cập

Kiến trúc của chúng tôi đáp ứng các yêu cầu này ngay từ khâu thiết kế.

---

## Bắt đầu

Khi từng giây đều quý giá, cơ sở hạ tầng đóng vai trò then chốt. Hãy triển khai đội drone phản ứng nhanh hơn và quan sát được nhiều hơn.

[Liên hệ với chúng tôi →](/contact)
