---
title: "Giám sát"
meta_title: "Giám sát thời gian thực & Nhận thức tình huống | Ubuntu Software"
description: "Biến luồng dữ liệu khổng lồ từ drone, vệ tinh và cảm biến IoT thành những cảnh báo có thể hành động ngay. Tự chủ, bảo mật, với khả năng truy xuất nguồn gốc đầy đủ."
image: "/images/robotics.svg"
draft: false
---

## Chìm ngập trong dữ liệu

Các hoạt động hiện đại tạo ra nhiều dữ liệu hơn khả năng xử lý của con người. Drone truyền phát video. Vệ tinh thu thập hình ảnh. Cảm biến báo cáo liên tục. Tín hiệu vẫn ở đó—bị chôn vùi trong nhiễu.

Các phương pháp truyền thống thất bại khi mở rộng quy mô:
- **Giám sát thủ công** không thể mở rộng quá vài nguồn dữ liệu
- **Ngưỡng đơn giản** gây ra tình trạng mệt mỏi vì cảnh báo
- **Hệ thống rời rạc** bỏ sót các mẫu hình liên nguồn
- **Giải pháp chỉ dùng đám mây** không thể xử lý các hoạt động nhạy cảm

Bạn cần một hệ thống theo dõi mọi thứ, hiểu ngữ cảnh và chỉ cảnh báo bạn khi thực sự cần thiết.

---

## Từ dòng dữ liệu ồ ạt đến thông tin chi tiết

Kiến trúc giám sát của chúng tôi xử lý dữ liệu từ nhiều nguồn theo thời gian thực, áp dụng AI ở cả biên (edge) và đám mây, sau đó lọc qua các quy tắc do người vận hành định nghĩa để đưa ra các sự kiện có thể hành động.

```
┌─────────────────────────────────────────────────────────────────┐
│                         DATA SOURCES                            │
├─────────────────┬─────────────────┬─────────────────────────────┤
│     Drones      │   Satellites    │        IoT Sensors          │
│  (telemetry,    │  (imagery,      │  (environmental,            │
│   video, AI)    │   positioning)  │   industrial)               │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    FIREHOSE INGESTION                           │
│            NATS JetStream — persistent, ordered streams         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AI PROCESSING                              │
│   Edge (Jetson): Real-time inference, event detection           │
│   Cloud: Pattern recognition, correlation, anomaly detection    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    OPERATOR RULES ENGINE                        │
│   "Alert me when..." — thresholds, patterns, conditions         │
│   Customizable per user, role, or organization                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    ALERT + CONTEXT                              │
│   • What happened (event description)                           │
│   • Where (2D map, 3D spatial view)                             │
│   • Why it matters (AI reasoning)                               │
│   • Provenance (full data lineage, evidence chain)              │
└─────────────────────────────────────────────────────────────────┘
```

---

## Nguồn dữ liệu

Chúng tôi coi mọi nguồn dữ liệu đều quan trọng ngang nhau—thu nhận thống nhất, xử lý nhất quán, phân tích tương quan.

| Nguồn | Loại dữ liệu | Ví dụ |
|--------|------------|----------|
| **Drone** | Dữ liệu đo xa, video, số đọc cảm biến, kết quả suy luận AI | Vị trí đội bay, ảnh nhiệt, phát hiện đối tượng |
| **Vệ tinh** | Hình ảnh, định vị, dữ liệu thời tiết | Phát hiện thay đổi, phân tích vùng phủ sóng, điều kiện khí quyển |
| **Cảm biến IoT** | Môi trường, công nghiệp, cơ sở hạ tầng | Nhiệt độ, độ rung, mực nước, chất lượng không khí |

Cùng một kiến trúc xử lý dữ liệu đo xa của drone cũng xử lý hình ảnh vệ tinh và cảm biến mặt đất. Một hệ thống. Bức tranh toàn cảnh.

---

## Bạn xác định điều gì quan trọng

Người vận hành cấu hình những gì họ quan tâm. Hệ thống theo dõi mọi thứ và chỉ cảnh báo khi xảy ra các điều kiện bạn chỉ định.

| Người vận hành nói | Hệ thống thực hiện |
|---------------|-------------|
| "Cảnh báo tôi nếu bất kỳ drone nào mất GPS" | Giám sát toàn bộ dữ liệu đo xa của đội bay, thông báo khi mất GPS |
| "Đánh dấu bất thường nhiệt trên đường dây điện" | Chạy phân tích nhiệt, cảnh báo các điểm nóng |
| "Theo dõi phương tiện đi vào khu vực hạn chế" | Giám sát hàng rào địa lý với cảnh báo tức thì |
| "Thông báo cho tôi khi mức độ căng thẳng của cây trồng vượt quá 20%" | Tổng hợp dữ liệu đa phổ, kích hoạt khi vượt ngưỡng |
| "Tương quan thay đổi từ vệ tinh với cảm biến mặt đất" | Đối chiếu chéo nhiều nguồn, phát hiện các mẫu hình |

Quy tắc có thể là các ngưỡng đơn giản hoặc tương quan đa nguồn phức tạp. Bạn kiểm soát độ nhạy.

---

## Trực quan hóa với nguồn gốc dữ liệu

*Tính năng trong lộ trình phát triển*

Khi một cảnh báo được kích hoạt, người vận hành thấy được bức tranh toàn cảnh:

- **Chế độ xem bản đồ 2D** — Ngữ cảnh vị trí kèm lịch sử hành trình
- **Chế độ xem không gian 3D** — Toàn bộ ngữ cảnh môi trường, địa hình, công trình
- **Nguồn gốc dữ liệu** — Chuỗi hoàn chỉnh từ dữ liệu cảm biến thô đến cảnh báo
- **Gói bằng chứng** — Dữ liệu gốc, kết quả đã xử lý, điểm tin cậy của AI

**Mọi quyết định đều có thể kiểm toán. Mọi cảnh báo đều có thể truy vết.**

Người vận hành có thể đi sâu từ cảnh báo đến dữ liệu thô, hiểu chính xác lý do hệ thống đánh dấu một sự kiện. Không có hộp đen.

---

## Bảo mật & Chủ quyền

Các hệ thống này xử lý những hoạt động nhạy cảm. Bảo mật không phải là tùy chọn—đó là nền tảng.

| Nguyên tắc | Triển khai |
|-----------|----------------|
| **Tự chủ quyền** | Bạn sở hữu và kiểm soát toàn bộ dữ liệu—chúng tôi không bao giờ truy cập khi chưa có sự cho phép rõ ràng |
| **Mã hóa đầu cuối** | Dữ liệu được mã hóa khi truyền và khi lưu trữ, với khóa do bạn kiểm soát |
| **Zero Trust** | Mọi yêu cầu đều được xác thực, mọi hành động đều được ủy quyền, không có sự tin cậy ngầm định |
| **Sẵn sàng cách ly mạng (Air-Gap)** | Đầy đủ chức năng mà không cần kết nối internet cho các triển khai nhạy cảm |
| **Nhật ký kiểm toán đầy đủ** | Mọi truy cập, mọi truy vấn, mọi cảnh báo đều được ghi lại kèm dấu thời gian |
| **Lưu trú dữ liệu** | Chọn nơi dữ liệu của bạn được lưu trữ—cơ sở hạ tầng của bạn, khu vực pháp lý của bạn |

**Dữ liệu của bạn. Cơ sở hạ tầng của bạn. Quyền kiểm soát của bạn.**

Triển khai tự lưu trữ chạy hoàn toàn trên phần cứng của bạn. Triển khai được quản lý cách ly dữ liệu của bạn bằng các khóa mã hóa chuyên dụng. Dù theo cách nào, bạn vẫn giữ quyền chủ quyền.

---

## Trường hợp sử dụng

| Lĩnh vực | Ứng dụng |
|--------|-------------|
| **An ninh biên giới** | Giám sát vành đai đa cảm biến, phát hiện xâm nhập, điều phối phản ứng |
| **Cơ sở hạ tầng trọng yếu** | Giám sát đường ống và đường dây điện, bảo trì dự đoán, phát hiện bất thường |
| **Môi trường** | Theo dõi động vật hoang dã, giám sát ô nhiễm, phát hiện cháy rừng, sức khỏe hệ sinh thái |
| **Nông nghiệp** | Giám sát sức khỏe cây trồng, tối ưu hóa tưới tiêu, cảnh báo sớm sâu bệnh |
| **Hàng hải** | Theo dõi tàu thuyền, phát hiện đánh bắt cá bất hợp pháp, an ninh cảng, tìm kiếm cứu nạn |
| **Vận hành đô thị** | Phân tích lưu lượng giao thông, giám sát đám đông, điều phối ứng phó khẩn cấp |

Mỗi trường hợp sử dụng đều tận dụng cùng một kiến trúc cốt lõi—cảm biến khác nhau, quy tắc khác nhau, cùng một nền tảng đáng tin cậy.

---

## Nền tảng kỹ thuật

Năng lực giám sát này được xây dựng trên kiến trúc đội bay đạt chuẩn sản xuất của chúng tôi:

- [Kiến trúc đội bay]({{< relref "/fleet" >}}) — Cách dữ liệu luân chuyển từ biên đến đám mây
- [Cấu trúc liên kết NATS]({{< relref "/fleet/nats-topology" >}}) — Xương sống nhắn tin xử lý dòng dữ liệu ồ ạt
- [Luồng & Sự kiện]({{< relref "/fleet/streams" >}}) — Cách chúng tôi biến các luồng liên tục thành trạng thái có thể truy vấn
- [Cổng phương tiện]({{< relref "/fleet/gateway" >}}) — Xử lý tại biên và phát hiện sự kiện
- [Ủy quyền & Cấp quyền]({{< relref "/fleet/authorization" >}}) — Cách bên thứ ba có được quyền truy cập an toàn, có phạm vi giới hạn

Cùng một cơ sở hạ tầng quản lý đội drone cũng vận hành các hoạt động giám sát. Đã được kiểm chứng thực tế ở quy mô lớn.

---

## Bắt đầu

Sẵn sàng biến dòng dữ liệu ồ ạt của bạn thành thông tin tình báo có thể hành động?

[Liên hệ với chúng tôi →](/contact) để thảo luận về các yêu cầu giám sát của bạn.

Chúng tôi sẽ giúp bạn hiểu rõ:
- Nên tích hợp những nguồn dữ liệu nào
- Những quy tắc vận hành nào phù hợp với lĩnh vực của bạn
- Các tùy chọn triển khai tự lưu trữ so với được quản lý
- Các cân nhắc về bảo mật và tuân thủ
