---
title: "Cấu hình JetStream"
meta_title: "Luồng JetStream cho bản sao số đội máy bay không người lái | Ubuntu Software"
description: "Cấu hình stream JetStream và kho KV cho bản sao số (digital twin) của đội phương tiện: tổng hợp trạng thái, nhật ký kiểm toán sự kiện, hàng đợi lệnh và trạng thái bóng (shadow state)."
image: "/images/robotics.svg"
draft: false
---

## JetStream cho Bản sao số

JetStream cung cấp lớp lưu trữ bền vững cho bản sao số (digital twin) của đội phương tiện. Mỗi loại stream phục vụ một mục đích cụ thể, với các thiết lập lưu giữ (retention) và sao chép (replication) phù hợp.

---

## Kiến trúc Stream

Năm tài nguyên JetStream hỗ trợ mô hình bản sao số:

| Tài nguyên | Loại | Mục đích |
|----------|------|---------|
| **TWIN_STATE** | Stream | Dữ liệu telemetry với rollup |
| **TWIN_EVENTS** | Stream | Nhật ký kiểm toán sự kiện bền vững |
| **TWIN_CMD** | Stream | Hàng đợi lệnh |
| **TWIN_SHADOW** | KV Store | Trạng thái mong muốn/được báo cáo |
| **TWIN_BLOBS** | Object Store | Tệp dữ liệu lớn |

---

## TWIN_STATE: Stream Telemetry

Lưu trữ trạng thái liên tục của phương tiện với **rollup theo subject**:

```yaml
name: TWIN_STATE
subjects:
  - "fleet.prod.veh.*.state.>"
retention: limits
max_age: 1h
max_bytes: 10GB
storage: file
replicas: 3
discard: old
rollup_hdrs: true
```

### Giải thích cấu hình

| Thiết lập | Giá trị | Lý do |
|---------|-------|-----------|
| **retention** | limits | Giới hạn theo thời gian và dung lượng |
| **max_age** | 1h | Giữ trạng thái gần đây để phát lại |
| **max_bytes** | 10GB | Giới hạn dung lượng lưu trữ cho mỗi hub |
| **replicas** | 3 | Chịu được sự cố node |
| **rollup_hdrs** | true | Bật rollup theo từng subject |

### Rollup theo Subject

Khi bật rollup, consumer có thể yêu cầu **chỉ thông điệp mới nhất của mỗi subject**:

```go
// Get latest position for all vehicles
sub, _ := js.Subscribe("fleet.prod.veh.*.state.position",
    nats.DeliverLastPerSubject())
```

Điều này biến stream từ một luồng dữ liệu ồ ạt thành một kho trạng thái có thể truy vấn.

### Các mẫu Consumer

| Mẫu | Trường hợp sử dụng |
|---------|----------|
| **DeliverAll** | Phát lại lịch sử, phân tích |
| **DeliverLast** | Ảnh chụp trạng thái hiện tại |
| **DeliverLastPerSubject** | Trạng thái mới nhất của từng phương tiện |
| **StartTime** | Phát lại từ một mốc thời gian cụ thể |

---

## TWIN_EVENTS: Nhật ký kiểm toán

Lưu trữ các sự kiện rời rạc với **thời gian lưu giữ dài** phục vụ tuân thủ và gỡ lỗi:

```yaml
name: TWIN_EVENTS
subjects:
  - "fleet.prod.veh.*.evt.>"
retention: limits
max_age: 90d
max_bytes: 100GB
storage: file
replicas: 3
discard: old
```

### Giải thích cấu hình

| Thiết lập | Giá trị | Lý do |
|---------|-------|-----------|
| **max_age** | 90d | Tuân thủ quy định, phân tích sự cố |
| **max_bytes** | 100GB | Sự kiện nhỏ hơn dữ liệu trạng thái |
| **storage** | file | Bền vững qua các lần khởi động lại |

### Thuộc tính sự kiện

Sự kiện khác với trạng thái:

| Thuộc tính | Trạng thái | Sự kiện |
|----------|-------|--------|
| **Tần suất** | Liên tục (Hz) | Rời rạc (khi có thay đổi) |
| **Lưu giữ** | Ngắn (vài giờ) | Dài (vài tháng) |
| **Consumer** | Dashboard thời gian thực | Kiểm toán, phân tích, cảnh báo |
| **Phát lại** | Tái tạo trạng thái hiện tại | Điều tra, tuân thủ |

### Ví dụ về sự kiện

```json
// Armed event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "previous_state": "disarmed",
    "armed_by": "operator-123",
    "reason": "mission_start"
  }
}

// Failsafe event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T11:45:00Z",
  "data": {
    "failsafe_type": "low_battery",
    "action_taken": "rtl",
    "battery_remaining": 0.15
  }
}
```

---

## TWIN_CMD: Hàng đợi lệnh

Lưu trữ các lệnh với **retention kiểu workqueue**:

```yaml
name: TWIN_CMD
subjects:
  - "fleet.prod.veh.*.cmd.>"
retention: workqueue
max_age: 5m
storage: file
replicas: 3
```

### Giải thích cấu hình

| Thiết lập | Giá trị | Lý do |
|---------|-------|-----------|
| **retention** | workqueue | Thông điệp bị xóa sau khi ack |
| **max_age** | 5m | Lệnh chưa được ack sẽ hết hạn |
| **replicas** | 3 | Đảm bảo lệnh được chuyển giao |

### Ngữ nghĩa Workqueue

1. Lệnh được publish vào stream
2. Vehicle Gateway nhận lệnh (chuyển giao đúng một lần)
3. Gateway xử lý lệnh
4. Gateway ack thông điệp
5. Thông điệp bị xóa khỏi stream

**Lệnh chưa được ack sẽ được gửi lại** để đảm bảo lệnh không bị mất khi mất kết nối.

### Luồng lệnh

```
Fleet Operator               TWIN_CMD Stream              Vehicle Gateway
     │                             │                             │
     │ Publish: cmd.takeoff        │                             │
     │ ─────────────────────────▶  │                             │
     │                             │ Deliver to consumer         │
     │                             │ ─────────────────────────▶  │
     │                             │                             │ Process
     │                             │                             │ command
     │                             │             Ack             │
     │                             │ ◀─────────────────────────  │
     │                             │ (message removed)           │
     │                             │                             │
```

---

## TWIN_SHADOW: KV Store

JetStream KV store duy trì trạng thái shadow của bản sao số:

```yaml
bucket: TWIN_SHADOW
max_value_size: 1MB
history: 5
ttl: 0
replicas: 3
storage: file
```

### Giải thích cấu hình

| Thiết lập | Giá trị | Lý do |
|---------|-------|-----------|
| **history** | 5 | Giữ 5 phiên bản trước đó |
| **ttl** | 0 | Không hết hạn (vĩnh viễn) |
| **replicas** | 3 | Tính sẵn sàng cao |

### Cấu trúc khóa

```
fleet/prod/veh/VID-001/desired    → Desired state JSON
fleet/prod/veh/VID-001/reported   → Reported state JSON
fleet/prod/veh/VID-001/config     → Vehicle configuration
fleet/prod/veh/VID-001/meta       → Metadata (serial, type, etc.)
```

### Thao tác Shadow

```go
// Write desired state
kv.Put("fleet/prod/veh/VID-001/desired", desiredJSON)

// Read reported state
entry, _ := kv.Get("fleet/prod/veh/VID-001/reported")

// Watch for changes
watcher, _ := kv.Watch("fleet/prod/veh/VID-001/desired")
for entry := range watcher.Updates() {
    // Desired state changed, reconcile
}
```

### Lịch sử để gỡ lỗi

Khi bật history, bạn có thể kiểm tra các trạng thái trước đó:

```go
// Get all historical values
history, _ := kv.History("fleet/prod/veh/VID-001/desired")
for _, entry := range history {
    fmt.Printf("Rev %d at %s: %s\n",
        entry.Revision(), entry.Created(), entry.Value())
}
```

---

## TWIN_BLOBS: Object Store

Các tệp dữ liệu lớn không vừa trong thông điệp:

```yaml
bucket: TWIN_BLOBS
max_chunk_size: 128KB
storage: file
replicas: 3
```

### Trường hợp sử dụng

| Tệp dữ liệu | Kích thước | Ví dụ |
|----------|------|---------|
| **Tệp nhiệm vụ** | 10-100KB | Định nghĩa điểm tham chiếu (waypoint) |
| **Tệp log** | 1-10MB | Log chuyến bay, chẩn đoán |
| **Firmware** | 10-50MB | Ảnh firmware PX4 |
| **Bản đồ** | 1-100MB | Ô bản đồ ngoại tuyến |
| **Mô hình ML** | 10-500MB | Tệp engine TensorRT |

### Thao tác Object

```go
// Upload mission file
obj, _ := obs.PutFile("missions/mission-456.json", file)

// Download to vehicle
reader, _ := obs.Get("missions/mission-456.json")

// List available firmware
objects := obs.List(nats.ObjectSearchPrefix("firmware/"))
```

### Truyền theo khối

Object Store xử lý tệp lớn bằng cách:

1. Chia thành các khối (mặc định 128KB)
2. Lưu các khối dưới dạng thông điệp stream
3. Theo dõi metadata (kích thước, hash, các khối)
4. Ghép lại khi truy xuất

Điều này cho phép **truyền tải có thể tiếp tục** qua các kết nối không ổn định.

---

## Định cỡ Stream

### Ước tính cho mỗi phương tiện

| Stream | Kích thước thông điệp | Tần suất | Khối lượng hằng ngày |
|--------|--------------|------|--------------|
| **STATE** | 200 byte | 10 Hz | ~170 MB |
| **EVENTS** | 500 byte | 10/ngày | ~5 KB |
| **CMD** | 300 byte | 100/ngày | ~30 KB |

### Ước tính cho đội phương tiện (1.000 phương tiện)

| Tài nguyên | 1 giờ | 1 ngày | 90 ngày |
|----------|--------|-------|--------|
| **TWIN_STATE** | 7 GB | - | - |
| **TWIN_EVENTS** | - | 5 MB | 450 MB |
| **TWIN_SHADOW** | - | 1 GB | 1 GB |
| **TWIN_BLOBS** | - | - | 50 GB |

Đây là các ước tính cơ sở. Mức sử dụng thực tế phụ thuộc vào tần suất và kích thước thông điệp.

---

## Các cân nhắc vận hành

### Tình trạng Stream

Theo dõi các chỉ số sau:

| Chỉ số | Ngưỡng cảnh báo |
|--------|-----------------|
| **Độ trễ consumer** | > 10.000 thông điệp |
| **Dung lượng stream** | > 80% mức tối đa |
| **Đồng bộ bản sao** | Bất kỳ bản sao nào bị chậm |
| **Ack đang chờ** | > 1.000 mỗi consumer |

### Chiến lược sao lưu

| Tài nguyên | Phương pháp sao lưu | Tần suất |
|----------|---------------|-----------|
| **Stream** | Snapshot JetStream | Hằng ngày |
| **KV Store** | Xuất khóa | Hằng ngày |
| **Object Store** | Sao lưu cấp tệp | Hằng tuần |

### Di chuyển

Stream có thể được di chuyển giữa các cluster:

1. Tạo stream trên cluster đích với cùng cấu hình
2. Tạo mirror lấy nguồn từ stream gốc
3. Chờ đồng bộ
4. Chuyển producer/consumer sang cluster đích
5. Xóa stream gốc

---

## Tóm tắt

| Tài nguyên | Loại | Lưu giữ | Mục đích |
|----------|------|-----------|---------|
| **TWIN_STATE** | Stream | 1 giờ | Dữ liệu telemetry với rollup |
| **TWIN_EVENTS** | Stream | 90 ngày | Nhật ký kiểm toán |
| **TWIN_CMD** | Stream | Workqueue | Chuyển giao lệnh |
| **TWIN_SHADOW** | KV | Vĩnh viễn | Trạng thái mong muốn/được báo cáo |
| **TWIN_BLOBS** | Object | Vĩnh viễn | Tệp dữ liệu lớn |

---

## Tiếp theo

[Vehicle Gateway →]({{< relref "/fleet/gateway" >}})
