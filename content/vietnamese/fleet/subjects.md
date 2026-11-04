---
title: "Đặt tên chủ đề"
meta_title: "Đặt tên Subject NATS cho đội máy bay không người lái | Ubuntu Software"
description: "Quy ước đặt tên subject NATS theo cấu trúc phân cấp cho dữ liệu telemetry, sự kiện, lệnh của đội phương tiện và quản lý trạng thái bản sao số (digital twin)."
image: "/images/robotics.svg"
draft: false
---

## Thiết kế phân cấp Subject

Cách đặt tên subject chính là API của hệ thống nhắn tin của bạn. Một hệ thống phân cấp được thiết kế tốt cho phép định tuyến, lọc và kiểm soát truy cập hiệu quả, đồng thời vẫn trực quan đối với các nhà phát triển.

---

## Nguyên tắc thiết kế

| Nguyên tắc | Lý do |
|-----------|-----------|
| **Phân cấp** | Cho phép đăng ký wildcard ở bất kỳ cấp nào |
| **Dễ dự đoán** | Nhà phát triển có thể tự xây dựng subject mà không cần tra cứu |
| **Nhận biết môi trường** | Tách biệt production khỏi staging/development |
| **Có thể lọc** | Đăng ký theo một phương tiện, một nhóm hoặc toàn bộ đội phương tiện |
| **Có thể mở rộng** | Thêm loại thông điệp mới mà không cần tái cấu trúc |

---

## Mẫu Subject cốt lõi

Tất cả subject của đội phương tiện đều tuân theo mẫu này:

```
fleet.<env>.veh.<vid>.<category>.<type>
```

| Phân đoạn | Mô tả | Ví dụ |
|---------|-------------|----------|
| `fleet` | Không gian tên gốc | Luôn là "fleet" |
| `<env>` | Môi trường | `prod`, `staging`, `dev` |
| `veh` | Loại tài nguyên | Phương tiện (tương lai: `gcs` cho trạm điều khiển mặt đất) |
| `<vid>` | ID phương tiện | `VID-001`, `drone-alpha-7` |
| `<category>` | Danh mục thông điệp | `state`, `evt`, `cmd`, `cmdack` |
| `<type>` | Loại cụ thể | `position`, `battery`, `takeoff` |

---

## Danh mục thông điệp

### Trạng thái: `fleet.<env>.veh.<vid>.state.*`

Trạng thái hiện tại của phương tiện, được publish liên tục:

| Subject | Mô tả | Tần suất |
|---------|-------------|------|
| `state.position` | Vĩ độ/kinh độ/độ cao, vận tốc | 10 Hz |
| `state.attitude` | Góc roll/pitch/yaw | 10 Hz |
| `state.battery` | Điện áp, dòng điện, dung lượng còn lại | 1 Hz |
| `state.mode` | Chế độ bay (AUTO, MANUAL, RTL) | Khi thay đổi |
| `state.health` | Tóm tắt tình trạng hệ thống | 1 Hz |

**Ví dụ publish:**

```go
subject := fmt.Sprintf("fleet.prod.veh.%s.state.position", vehicleID)
nc.Publish(subject, positionJSON)
```

### Sự kiện: `fleet.<env>.veh.<vid>.evt.*`

Các sự kiện rời rạc, được publish khi chúng xảy ra:

| Subject | Mô tả |
|---------|-------------|
| `evt.armed` | Phương tiện đã được kích hoạt (arm) |
| `evt.disarmed` | Phương tiện đã được vô hiệu hóa (disarm) |
| `evt.takeoff` | Bắt đầu cất cánh |
| `evt.landed` | Hạ cánh hoàn tất |
| `evt.waypoint` | Đã đến điểm tham chiếu (waypoint) |
| `evt.failsafe` | Cơ chế an toàn (failsafe) được kích hoạt |
| `evt.geofence` | Vi phạm hàng rào địa lý (geofence) |
| `evt.battery.low` | Cảnh báo pin yếu |

**Sự kiện là các chuyển đổi trạng thái**, không phải dữ liệu liên tục. Chúng có tính bền vững—được lưu trong JetStream để phục vụ nhật ký kiểm toán.

### Lệnh: `fleet.<env>.veh.<vid>.cmd.*`

Các chỉ thị được gửi đến phương tiện:

| Subject | Mô tả |
|---------|-------------|
| `cmd.arm` | Kích hoạt (arm) động cơ |
| `cmd.disarm` | Vô hiệu hóa (disarm) động cơ |
| `cmd.takeoff` | Bắt đầu cất cánh |
| `cmd.land` | Bắt đầu hạ cánh |
| `cmd.rtl` | Quay về điểm xuất phát |
| `cmd.goto` | Di chuyển đến vị trí |
| `cmd.mission.upload` | Tải lên nhiệm vụ |
| `cmd.mission.start` | Bắt đầu nhiệm vụ |
| `cmd.param.set` | Thiết lập tham số |

**Lệnh đi theo chiều xuống (downstream)** từ hệ thống quản lý đội phương tiện đến các phương tiện.

### Xác nhận lệnh: `fleet.<env>.veh.<vid>.cmdack.*`

Phản hồi cho các lệnh:

| Subject | Mô tả |
|---------|-------------|
| `cmdack.arm` | Kết quả lệnh arm |
| `cmdack.takeoff` | Kết quả lệnh cất cánh |
| `cmdack.goto` | Kết quả lệnh goto |

**Payload của ack bao gồm:**

```json
{
  "cmd_id": "cmd-123456",
  "status": "accepted|rejected|completed|failed",
  "error": "optional error message",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Đăng ký wildcard

Wildcard của NATS cho phép đăng ký linh hoạt:

| Wildcard | Ý nghĩa | Ví dụ |
|----------|---------|---------|
| `*` | Một token | `fleet.prod.veh.*.state.position` |
| `>` | Một hoặc nhiều token | `fleet.prod.veh.VID-001.>` |

### Các mẫu đăng ký phổ biến

| Mẫu | Trường hợp sử dụng |
|---------|----------|
| `fleet.prod.veh.*.state.position` | Vị trí của tất cả phương tiện |
| `fleet.prod.veh.VID-001.>` | Mọi thứ từ một phương tiện |
| `fleet.prod.veh.*.evt.>` | Tất cả sự kiện từ tất cả phương tiện |
| `fleet.prod.veh.*.state.battery` | Tất cả trạng thái pin |
| `fleet.*.veh.VID-001.>` | Một phương tiện trên tất cả môi trường |

---

## KV Store: Shadow cho Digital Twin

JetStream KV store duy trì trạng thái **mong muốn** (desired) và **được báo cáo** (reported):

### Mẫu Subject (KV)

```
fleet/<env>/veh/<vid>/desired
fleet/<env>/veh/<vid>/reported
```

Lưu ý: KV sử dụng dấu phân cách `/` (đường dẫn bucket), trong khi pub/sub sử dụng dấu phân cách `.`.

### Trạng thái mong muốn

Những gì hệ thống quản lý đội phương tiện **muốn** phương tiện thực hiện:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "geofence_enabled": true,
  "max_altitude": 120,
  "home_position": {
    "lat": 37.7749,
    "lon": -122.4194
  }
}
```

### Trạng thái được báo cáo

Những gì phương tiện **thực sự** báo cáo:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "mission_progress": 0.45,
  "geofence_enabled": true,
  "armed": true,
  "in_flight": true,
  "battery_remaining": 0.72
}
```

### Đối chiếu Shadow

Vehicle Gateway liên tục:

1. Đọc trạng thái mong muốn từ KV
2. So sánh với trạng thái thực tế của phương tiện
3. Phát lệnh để đưa thực tế về đúng trạng thái mong muốn
4. Cập nhật trạng thái được báo cáo

Điều này cho phép **quản lý đội phương tiện theo kiểu khai báo**—thiết lập trạng thái mong muốn, để các phương tiện tự hội tụ về trạng thái đó.

---

## Kiểm soát truy cập

Cơ chế phân quyền của NATS ánh xạ người dùng tới các subject được phép:

```
# Vehicle can publish its own state/events
publish: fleet.prod.veh.VID-001.state.>
publish: fleet.prod.veh.VID-001.evt.>
publish: fleet.prod.veh.VID-001.cmdack.>

# Vehicle can subscribe to commands for itself
subscribe: fleet.prod.veh.VID-001.cmd.>

# Vehicle can read/write its own shadow
kv: fleet/prod/veh/VID-001/*
```

```
# Fleet operator can subscribe to all state/events
subscribe: fleet.prod.veh.*.state.>
subscribe: fleet.prod.veh.*.evt.>
subscribe: fleet.prod.veh.*.cmdack.>

# Fleet operator can publish commands
publish: fleet.prod.veh.*.cmd.>

# Fleet operator can manage shadows
kv: fleet/prod/veh/*/desired
```

Điều này đảm bảo **phương tiện không thể mạo danh phương tiện khác** hoặc publish các lệnh trái phép.

Để tìm hiểu chi tiết về mô hình phân quyền, bao gồm cấp quyền cho bên thứ ba và quản lý thông tin xác thực, xem [Phân quyền & Cấp quyền]({{< relref "/fleet/authorization" >}}).

---

## Đánh đổi trong thiết kế Subject

### Tại sao không dùng subject phẳng?

Các subject phẳng như `vehicle-001-position` không hỗ trợ:

- Đăng ký wildcard
- Kiểm soát truy cập theo phân cấp
- Nhóm theo logic

### Tại sao không dùng phân cấp sâu hơn?

Các subject như `fleet.prod.region.us-west.site.warehouse-a.veh.VID-001.state.position`:

- Khó xây dựng bằng lập trình hơn
- Wildcard trở nên phức tạp
- Thay đổi lan truyền khắp hệ thống

Phân cấp bốn cấp cân bằng giữa tính linh hoạt và sự đơn giản.

---

## Định dạng thông điệp

Tất cả thông điệp đều sử dụng **JSON** với cấu trúc nhất quán:

```json
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00.123Z",
  "seq": 12345,
  "data": {
    // message-specific payload
  }
}
```

| Trường | Mô tả |
|-------|-------------|
| `vid` | ID phương tiện (trùng lặp với subject, nhưng hữu ích cho consumer) |
| `timestamp` | ISO 8601 với độ chính xác đến mili giây |
| `seq` | Số thứ tự tăng đơn điệu để sắp xếp thứ tự |
| `data` | Payload riêng của từng loại thông điệp |

---

## Tóm tắt

| Danh mục | Mẫu Subject | Hướng |
|----------|-----------------|-----------|
| **Trạng thái** | `fleet.<env>.veh.<vid>.state.<type>` | Phương tiện → Hub |
| **Sự kiện** | `fleet.<env>.veh.<vid>.evt.<type>` | Phương tiện → Hub |
| **Lệnh** | `fleet.<env>.veh.<vid>.cmd.<type>` | Hub → Phương tiện |
| **Ack** | `fleet.<env>.veh.<vid>.cmdack.<type>` | Phương tiện → Hub |
| **Mong muốn** | `fleet/<env>/veh/<vid>/desired` (KV) | Hub → Phương tiện |
| **Được báo cáo** | `fleet/<env>/veh/<vid>/reported` (KV) | Phương tiện → Hub |

---

## Tiếp theo

[Cấu hình Stream →]({{< relref "/fleet/streams" >}})
