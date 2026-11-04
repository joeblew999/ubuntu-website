---
title: "Hoạt động hiện trường do AI điều khiển"
meta_title: "Hậu cần & Vận hành hiện trường ứng dụng AI | Ubuntu Software"
description: "Tối ưu hóa lộ trình bằng AI theo thời gian thực cho các hoạt động hiện trường. Tài xế xe tải, đội xe giao hàng và đội ngũ dịch vụ nhận được lộ trình do AI chỉ dẫn, có thể điều chỉnh tức thì qua giọng nói hoặc văn bản."
image: "/images/robotics.svg"
draft: false
---

## Khi ba nền tảng hội tụ

**Fleet + Publish + Foundation** kết hợp để tạo nên một điều mới: vận hành hiện trường do AI điều hướng, trong đó người vận hành trở thành các điểm biên thông minh trong một hệ thống được kết nối toàn diện.

Tài xế xe tải dùng iPad nhận được lộ trình thu gom do AI chỉ định. Hệ thống biết vị trí của từng phương tiện. Người vận hành có thể cập nhật lộ trình bằng giọng nói hoặc văn bản—và AI tái cấu trúc tất cả các lộ trình bị ảnh hưởng ngay lập tức.

---

## Vấn đề: Lộ trình tĩnh trong một thế giới động

Các hoạt động hiện trường truyền thống được lên kế hoạch từ tối hôm trước:

- **Lập kế hoạch lộ trình** — Điều phối viên tạo lộ trình thủ công, hoặc dùng tối ưu hóa theo lô qua đêm
- **Thay đổi phút chót** — Tài xế bị ốm? Đường bị chặn? Gọi điện, lập lại kế hoạch, hỗn loạn
- **Không có khả năng quan sát theo thời gian thực** — "Xe tải 47 đang ở đâu?" phải gọi bộ đàm mới biết
- **Quy trình dựa trên giấy tờ** — Lệnh công việc được in ra, trạng thái hoàn thành được nhập lại thủ công
- **Thông tin bị cô lập** — Tài xế biết một điều, điều phối viên không biết, khách hàng hoàn toàn không hay

Khi thực tế khác với kế hoạch, toàn bộ hệ thống chật vật để thích ứng.

---

## Giải pháp: Vận hành hiện trường do AI điều hướng

```
┌─────────────────────────────────────────────────────────────────────┐
│                    UBUNTU SOFTWARE UNIFIED STACK                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────┐  │
│  │   PUBLISH   │    │    FLEET    │    │      FOUNDATION         │  │
│  │             │    │             │    │                         │  │
│  │ Work orders │    │  Vehicles   │    │  Offline-first sync     │  │
│  │ Forms/PDFs  │◄──►│  Telemetry  │◄──►│  Automerge CRDT         │  │
│  │ Kiosks      │    │  Commands   │    │  NATS JetStream         │  │
│  │ Translation │    │  Digital    │    │  SQLite everywhere      │  │
│  │             │    │  Twins      │    │  WellKnown Gateway      │  │
│  └─────────────┘    └─────────────┘    └─────────────────────────┘  │
│         │                  │                       │                 │
│         └──────────────────┼───────────────────────┘                 │
│                            │                                         │
│                            ▼                                         │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                   FIELD OPERATOR (Human Edge)                │    │
│  │                                                               │    │
│  │   iPad/Phone with:                                           │    │
│  │   • Real-time route display (2D/3D)                          │    │
│  │   • AI-directed pickup notifications                         │    │
│  │   • SST/voice input for route changes                        │    │
│  │   • Offline capability (queues sync when connected)          │    │
│  │   • Position tracking (device GPS)                           │    │
│  │   • Work order completion/forms                              │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                            │                                         │
│                            ▼                                         │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                      VEHICLE (Truck)                         │    │
│  │                                                               │    │
│  │   • CAN bus telemetry (engine, fuel, diagnostics)            │    │
│  │   • GPS position (location tracking)                         │    │
│  │   • Gateway syncs to NATS JetStream                          │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Nhận định then chốt:** Người vận hành là các nút điện toán biên. Thiết bị của họ cung cấp định vị, nhận chỉ dẫn từ AI và cho phép cập nhật hệ thống ngay tức thì—giống như bất kỳ cảm biến nào khác trong mạng lưới.

---

## Ví dụ: Thu gom rác thải được tái định hình

Một ví dụ cụ thể về cách kiến trúc này hoạt động:

### Buổi sáng: AI lên kế hoạch cho cả ngày

1. Hệ thống biết tất cả các điểm thu gom, mức độ ưu tiên, mức đầy của thùng chứa (cảm biến IoT)
2. AI tính toán lộ trình tối ưu cho từng xe tải—giảm thiểu quãng đường, tuân thủ khung thời gian
3. Lộ trình được tự động đẩy đến máy tính bảng của tài xế trước khi bắt đầu ca
4. Mỗi tài xế xem lộ trình cá nhân hóa của mình dưới dạng 2D/3D với chỉ đường từng bước

### Trong quá trình vận hành: Thích ứng theo thời gian thực

**Tài xế 1 báo cáo:** "123 Main St bị chặn do thi công"

```
Driver speaks into tablet
    └─▶ SST converts to text
    └─▶ System marks location as blocked
    └─▶ AI refactors all affected routes
    └─▶ Updated routes push to affected drivers
    └─▶ Total time: seconds
```

Không cần gọi điện. Không cần điều phối viên can thiệp. Hệ thống tự xử lý.

**Trong khi đó:**
- Bảng điều khiển điều phối hiển thị tất cả xe tải theo thời gian thực
- Cổng thông tin khách hàng hiển thị thời gian đến dự kiến (ETA) cho lượt thu gom của họ
- Mọi lượt hoàn thành đều được ghi lại tự động
- Mọi thay đổi đều được kiểm toán với đầy đủ nguồn gốc

### Con người trong vòng lặp

Người vận hành có thể ghi đè AI bất cứ lúc nào:

| Phương thức nhập | Ví dụ | Phản hồi của hệ thống |
|--------------|---------|-----------------|
| **Giọng nói (SST)** | "Bỏ qua 456 Oak Ave, khách hàng không có nhà" | Đánh dấu bỏ qua, điều chỉnh lộ trình |
| **Chạm** | Chạm vào vị trí trên bản đồ, chọn "Bị chặn" | Xóa khỏi tất cả lộ trình |
| **Biểu mẫu** | Hoàn thành biểu mẫu dịch vụ kèm ghi chú | Ghi nhận hoàn thành, kích hoạt lập hóa đơn |
| **Ghi đè** | Kéo lộ trình để thay đổi thứ tự | Chấp nhận thay đổi, tối ưu lại các chặng phía sau |

AI tăng cường phán đoán của con người—chứ không thay thế nó.

---

## Con người trong vòng lặp: Người vận hành là trí tuệ biên

Bước đột phá không chỉ là AI—mà là **AI + phán đoán của con người cùng phối hợp**. Người vận hành hiện trường không phải là người thụ động nhận chỉ thị. Họ là các nút thông minh trong hệ thống.

### Thiết bị được khai thác toàn diện

iPad/điện thoại của mỗi tài xế trở thành vừa là cảm biến VỪA là bộ chấp hành:

```
┌─────────────────────────────────────────────────────────────────┐
│                    DRIVER DEVICE (iPad/Phone)                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  INPUTS (Sensing)              OUTPUTS (Acting)                  │
│  ─────────────────             ────────────────                  │
│  • GPS position                • Route display (2D/3D)           │
│  • Voice commands (SST)        • Turn-by-turn directions         │
│  • Touch interactions          • Pickup notifications            │
│  • Form submissions            • Alert sounds/haptics            │
│  • Camera (proof of service)   • Customer ETA updates            │
│  • Accelerometer (driving)     • Route change confirmations      │
│                                                                  │
│  LOCAL INTELLIGENCE                                              │
│  ──────────────────                                              │
│  • SQLite database (offline state)                               │
│  • Automerge CRDT (conflict-free sync)                           │
│  • Background location tracking                                  │
│  • Push notification handling                                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

Hệ thống biết mọi tài xế đang ở đâu, đang làm gì và có thể liên lạc với họ ngay lập tức. Nhưng điều quan trọng là—**tài xế có thể phản hồi lại**.

### Tái cấu trúc lộ trình bằng giọng nói

Chuyển giọng nói thành văn bản (SST) cho phép cập nhật hệ thống mà không cần dùng tay:

| Tài xế nói | Hệ thống thực hiện |
|-------------|-------------|
| "123 Main bị chặn do thi công" | Đánh dấu vị trí bị chặn, tái cấu trúc tất cả lộ trình bị ảnh hưởng |
| "Bỏ qua điểm thu gom tiếp theo, không có ai ở nhà" | Xóa khỏi lộ trình, thông báo cho khách hàng, điều chỉnh ETA |
| "Thùng chứa ở Oak Street bị tràn" | Gắn cờ ưu tiên, có thể điều chuyển sang xe tải lớn hơn |
| "Nghỉ giải lao 15 phút" | Điều chỉnh tất cả ETA phía sau, thông báo cho khách hàng bị ảnh hưởng |
| "Hoàn thành lộ trình" | Đóng lộ trình, cập nhật bảng điều khiển, kích hoạt lập hóa đơn |

Không gọi điện cho bộ phận điều phối. Không liên lạc bộ đàm rườm rà. Tài xế nói, AI lắng nghe, hệ thống thích ứng.

### Phân cấp thẩm quyền

Phán đoán của con người được ưu tiên hơn các khuyến nghị của AI:

```
┌─────────────────────────────────────────────────────────────────┐
│                     DECISION AUTHORITY                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1. DRIVER ON-SITE      ───────────────────────▶  HIGHEST      │
│      Sees reality, can override anything                         │
│                                                                  │
│   2. DISPATCHER          ───────────────────────▶  HIGH         │
│      Fleet-wide view, can reassign routes                        │
│                                                                  │
│   3. AI OPTIMIZER        ───────────────────────▶  MEDIUM       │
│      Suggests routes, adapts to changes                          │
│                                                                  │
│   4. SCHEDULED PLAN      ───────────────────────▶  LOWEST       │
│      Starting point, always subject to reality                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

AI rất mạnh mẽ—nhưng người trực tiếp quan sát tại hiện trường mới có tiếng nói cuối cùng.

---

## Xuất bản: Một nguồn duy nhất, mọi kênh

[Nền tảng Publish]({{< relref "/platform/publish" >}}) tạo ra tất cả đầu ra từ một nguồn sự thật duy nhất.

### Đầu ra logistics

| Định nghĩa nguồn | Đầu ra được tạo |
|-------------------|-------------------|
| **Bản ghi thu gom** | Màn hình máy tính bảng, bảng kê PDF, SMS cho khách hàng, màn hình kiosk |
| **Kế hoạch lộ trình** | Chỉ đường từng bước cho tài xế, bảng điều khiển điều phối, ETA cho khách hàng |
| **Biểu mẫu dịch vụ** | Biểu mẫu kỹ thuật số, PDF có thể in, giấy được quét OCR |
| **Tóm tắt hằng ngày** | Báo cáo qua email, widget bảng điều khiển, bản xuất tuân thủ |

**Một định nghĩa. Nhiều đầu ra.** Thay đổi nguồn, tất cả các kênh tự động cập nhật.

### Đội ngũ hiện trường đa ngôn ngữ

Nhân viên hiện trường nói các ngôn ngữ khác nhau? Publish xử lý việc dịch thuật:

- Máy tính bảng của Tài xế 1: Tiếng Anh
- Máy tính bảng của Tài xế 2: Tiếng Tây Ban Nha
- Cổng thông tin khách hàng: Theo lựa chọn của khách hàng
- Bảng kê PDF: Tiêu chuẩn của doanh nghiệp

Cùng dữ liệu, cách trình bày được bản địa hóa. Không cần quản lý nội dung trùng lặp.

### Biểu mẫu phản hồi ngược

Biểu mẫu hoàn thành dịch vụ được đưa ngược trở lại hệ thống:

```
Driver completes pickup form
    └─▶ Photo attached (proof of service)
    └─▶ Notes recorded ("Container damaged")
    └─▶ Time stamped automatically
    └─▶ GPS coordinates captured
    └─▶ Data syncs via Foundation
    └─▶ Billing system updated
    └─▶ Customer notified
    └─▶ Maintenance ticket created (if needed)
```

Biểu mẫu không chỉ để thu thập dữ liệu—chúng là yếu tố kích hoạt tự động hóa ở các bước tiếp theo.

### Chủ quyền dữ liệu

Dữ liệu của bạn đi đâu? [WellKnown Gateway]({{< relref "/platform/publish#wellknown-gateway-data-sovereignty" >}}) đảm bảo bạn sở hữu dữ liệu của mình ngay cả khi xuất bản LÊN các nền tảng bên ngoài như Salesforce, hệ thống theo dõi đội xe hoặc nhà cung cấp dịch vụ lập hóa đơn.

**Nguyên tắc then chốt:** Bạn xuất bản LÊN các nền tảng. Họ không sở hữu dữ liệu của bạn—bạn mới là người sở hữu.

[Toàn bộ Nền tảng Publish →]({{< relref "/platform/publish" >}})

---

## Tích hợp Fleet: Phương tiện là nguồn dữ liệu

**Kiến trúc Fleet** cung cấp lớp phương tiện—xe tải trở thành các bản sao số giàu dữ liệu.

### Những gì xe tải biết

Mọi phương tiện đều truyền dữ liệu đo xa qua cổng CAN bus:

| Điểm dữ liệu | Nguồn | Công dụng |
|------------|--------|-----|
| **Vị trí GPS** | Mô-đun viễn thông (telematics) | Theo dõi thời gian thực, tính toán ETA |
| **Trạng thái động cơ** | ECU | Giám sát tình trạng, cảnh báo bảo trì |
| **Mức nhiên liệu** | Cảm biến nhiên liệu | Ước tính phạm vi hoạt động, lên lịch tiếp nhiên liệu |
| **Tốc độ** | Cảm biến bánh xe | Giám sát an toàn, tính thời gian lộ trình |
| **Chẩn đoán** | OBD-II / J1939 | Bảo trì dự đoán |
| **Trạng thái cửa** | Bộ điều khiển thân xe | Xác minh dịch vụ |
| **Trọng lượng** | Cảm biến tải trọng | Quản lý tải trọng |

Xe tải không chỉ là phương tiện vận chuyển—nó là một nền tảng cảm biến di động.

### Đồng bộ bản sao số

Mỗi xe tải đều có một bản sao số trong hệ thống:

```
Physical Truck                    Digital Twin
─────────────                    ────────────
Position: 34.05°, -118.24°   →   state.position
Fuel: 67%                    →   state.fuel
Engine: Running              →   state.engine
Speed: 25 mph                →   state.speed
Door: Closed                 →   state.door
Weight: 3,400 lbs            →   state.weight
```

Điều phối viên nhìn thấy thực tế. AI tối ưu hóa dựa trên thực tế. Khách hàng nhận được ETA chính xác.

[Kiến trúc Fleet →]({{< relref "/fleet" >}}) | [Phương tiện mặt đất →]({{< relref "/fleet/platforms/ground" >}}) | [Cổng phương tiện →]({{< relref "/fleet/gateway" >}})

---

## Foundation: Ưu tiên ngoại tuyến cho mọi thứ

**Nền tảng Foundation** đảm bảo hệ thống vẫn hoạt động ngay cả khi mất kết nối.

### Tại sao ngoại tuyến lại quan trọng

Hoạt động hiện trường diễn ra ở những nơi mạng không ổn định:

- Bãi đỗ xe ngầm
- Vùng nông thôn có sóng chập chờn
- Đường hầm và đường chui
- Khu vực đông đúc khi có sự kiện
- Biên giới quốc tế (khoảng trống chuyển vùng)

Nếu ứng dụng của bạn bị treo khi ngoại tuyến, hoạt động của bạn cũng bị đình trệ. Không thể chấp nhận được.

### Automerge CRDT: Đồng bộ không xung đột

Các thay đổi thực hiện khi ngoại tuyến được hợp nhất gọn gàng khi có kết nối trở lại:

```
Driver A (online):    Marks pickup complete at 9:15am
Driver B (offline):   Marks same pickup skipped at 9:14am (tunnel)

Driver B reconnects:
    └─▶ Automerge detects conflict
    └─▶ Applies "last write wins" or custom merge rule
    └─▶ Both changes preserved in history
    └─▶ Alert triggers for dispatcher review
    └─▶ No data loss, no app crash
```

Cơ sở dữ liệu truyền thống sẽ làm hỏng hoặc mất dữ liệu. CRDT xử lý điều này một cách êm thấm.

### SQLite ở mọi nơi

Mỗi thiết bị đều có một cơ sở dữ liệu cục bộ hoàn chỉnh:

- **Máy tính bảng của tài xế**: Lộ trình, khách hàng, biểu mẫu, lịch sử
- **Bảng điều khiển điều phối**: Trạng thái đội xe, cảnh báo, phân tích
- **Cổng thông tin khách hàng**: Hồ sơ dịch vụ của họ, các lượt thu gom sắp tới

Không có mạng? Không vấn đề gì. Ứng dụng có mọi thứ cần thiết ngay trên thiết bị.

[Nền tảng Foundation →]({{< relref "/platform/foundation" >}})

---

## Nhận biết toàn diện: Trực quan hóa 2D/3D

Người vận hành và điều phối viên cùng nhìn thấy một thực tế:

### Góc nhìn tài xế (Máy tính bảng)
- Lộ trình hiện tại với chỉ đường từng bước
- Chi tiết điểm thu gom tiếp theo và ghi chú của khách hàng
- Tất cả các xe tải khác trong khu vực (nhận biết tình huống)
- Lớp phủ thông tin giao thông và tình trạng đường

### Góc nhìn điều phối (Bảng điều khiển)
- Tất cả phương tiện được định vị trên bản đồ
- Tiến độ lộ trình hiện tại của từng tài xế
- Cảnh báo và ngoại lệ được làm nổi bật
- Lịch sử hành trình và nhật ký kiểm toán

### Góc nhìn khách hàng (Cổng thông tin)
- ETA cho dịch vụ của họ
- Thông báo khi tài xế sắp đến
- Xác nhận hoàn thành dịch vụ

Tất cả các góc nhìn đều bắt nguồn từ cùng một luồng dữ liệu thời gian thực. Tất cả đều được đồng bộ tự động.

---

## Nền tảng kỹ thuật

### Khả năng phục hồi ngoại tuyến

Hệ thống vẫn tiếp tục hoạt động khi mất kết nối:

```
Connectivity lost
    └─▶ App continues with local SQLite
    └─▶ Changes queue in Automerge CRDT
    └─▶ Position caches locally

Connectivity returns
    └─▶ Automerge merges all changes
    └─▶ No conflicts, no data loss
    └─▶ System catches up automatically
```

Hoạt động hiện trường không thể dừng lại vì sự cố mạng. Kiến trúc này đảm bảo điều đó không xảy ra.

### Đồng bộ thời gian thực

Khi có kết nối, các thay đổi được lan truyền ngay lập tức:

```
Driver completes pickup
    └─▶ Form data saves to local SQLite
    └─▶ Automerge syncs to NATS
    └─▶ JetStream persists event
    └─▶ Dashboard updates
    └─▶ Customer notified
    └─▶ Total latency: <1 second
```

### Kiến trúc thông báo

Mỗi thiết bị đăng ký các subject NATS liên quan:

| Mẫu subject | Ai nhận |
|-----------------|--------------|
| `fleet.prod.veh.{id}.route.*` | Thiết bị của tài xế đó |
| `fleet.prod.veh.*.state.pos` | Bảng điều khiển điều phối |
| `fleet.prod.customer.{id}.*` | Cổng thông tin khách hàng |

Thay đổi được xuất bản một lần, tự động gửi đến tất cả các bên quan tâm.

---

## Không chỉ thu gom rác thải

Cùng một mô hình áp dụng cho mọi hoạt động hiện trường:

| Ngành | Ứng dụng |
|----------|-------------|
| **Quản lý chất thải** | Thu gom tối ưu lộ trình, cảm biến mức đầy |
| **Giao hàng & Logistics** | Tối ưu chặng cuối, bằng chứng giao hàng |
| **Dịch vụ hiện trường** | Điều phối kỹ thuật viên, tồn kho linh kiện |
| **Nông nghiệp** | Điều phối thiết bị nông nghiệp, định tuyến thu hoạch |
| **Xây dựng** | Giao thiết bị/vật liệu, logistics công trường |
| **Tiện ích công cộng** | Lộ trình đọc đồng hồ, cuộc gọi dịch vụ |
| **Y tế** | Thăm khám chăm sóc tại nhà, thu thập mẫu bệnh phẩm |

Cùng kiến trúc. Cùng mô hình. Khác dữ liệu lĩnh vực.

---

## Bảo mật & Chủ quyền dữ liệu

Dữ liệu vận hành của bạn vẫn thuộc về bạn:

| Nguyên tắc | Triển khai |
|-----------|----------------|
| **Tự chủ** | Dữ liệu được lưu trữ trên hạ tầng của bạn hoặc của chúng tôi—tùy bạn lựa chọn |
| **WellKnown Gateway** | Xuất bản LÊN các nền tảng mà không bị phụ thuộc nhà cung cấp |
| **Mã hóa đầu cuối** | Dữ liệu được mã hóa khi truyền và khi lưu trữ |
| **Zero Trust** | Mọi thiết bị đều được xác thực, mọi yêu cầu đều được ủy quyền |
| **Nhật ký kiểm toán** | Mọi thay đổi đều được ghi lại với đầy đủ nguồn gốc |
| **Ưu tiên ngoại tuyến** | Dữ liệu nằm trên thiết bị, không chỉ trên đám mây |

[WellKnown Gateway →]({{< relref "/platform/publish#wellknown-gateway-data-sovereignty" >}}) | [Ủy quyền & Cấp quyền →]({{< relref "/fleet/authorization" >}})

---

## Bắt đầu

Vận hành hiện trường do AI điều hướng thay đổi cách công việc được thực hiện. Người vận hành của bạn trở thành các điểm biên thông minh trong một hệ thống thích ứng theo thời gian thực.

Sẵn sàng khám phá cách Fleet + Publish + Foundation có thể chuyển đổi hoạt động hiện trường của bạn?

[Liên hệ với chúng tôi →](/contact)
