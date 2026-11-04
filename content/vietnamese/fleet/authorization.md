---
title: "Ủy quyền & Cấp quyền"
meta_title: "Ủy quyền phi tập trung & Cấp quyền truy cập | Ubuntu Software"
description: "Mô hình bảo mật phi tập trung của NATS JetStream với cơ chế cấp quyền cho phép các bên thứ ba tham gia một cách an toàn vào các hệ thống thời gian thực mà không cần điều phối tập trung."
image: "/images/robotics.svg"
draft: false
---

## Bảo mật không có nút thắt cổ chai

Cơ chế ủy quyền tập trung truyền thống gây ra nhiều vấn đề khi mở rộng quy mô:

| Vấn đề | Tác động |
|---------|--------|
| **Điểm lỗi đơn lẻ** | Máy chủ xác thực ngừng hoạt động = toàn bộ hệ thống ngừng hoạt động |
| **Độ trễ trên mỗi yêu cầu** | Chuyến khứ hồi đến máy chủ xác thực làm tăng độ trễ |
| **Không thể hoạt động khi mất kết nối** | Không có mạng = không có ủy quyền |
| **Không mở rộng được** | Hàng nghìn thiết bị làm quá tải máy chủ trung tâm |

NATS tiếp cận theo cách khác: **ủy quyền phi tập trung** với thông tin xác thực tự chứa.

---

## Cách hoạt động của ủy quyền NATS

```
┌─────────────────────────────────────────────────────────────────┐
│                    CREDENTIAL ISSUANCE                          │
│   Account Authority issues signed JWTs with embedded perms      │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌─────────┐     ┌─────────┐     ┌─────────┐
        │ Vehicle │     │ Operator│     │ Partner │
        │  Creds  │     │  Creds  │     │  Creds  │
        └────┬────┘     └────┬────┘     └────┬────┘
             │               │               │
             ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS SERVER                                │
│   Validates JWT signature locally — no external call needed     │
│   Enforces permissions from embedded claims                     │
└─────────────────────────────────────────────────────────────────┘
```

**Các đặc tính chính:**

- **Xác thực cục bộ** — Máy chủ xác thực thông tin xác thực mà không cần gọi bất kỳ dịch vụ bên ngoài nào
- **Quyền được nhúng sẵn** — Những gì bạn có thể làm được mã hóa trong thông tin xác thực của bạn
- **Thu hồi tức thì** — Danh sách thu hồi cấp tài khoản được lan truyền khắp cụm
- **Zero trust** — Không có thông tin xác thực = không có quyền truy cập, chấm hết

Cơ chế này vẫn hoạt động khi mất kết nối. Không cần gọi về trung tâm để ủy quyền. Thông tin xác thực chính là bằng chứng.

---

## Phạm vi quyền

Thông tin xác thực chỉ định chính xác những gì một bên có thể làm:

| Phạm vi | Mô tả | Ví dụ |
|-------|-------------|---------|
| **Publish** | Gửi tin nhắn đến các subject | Phương tiện publish `fleet.prod.veh.V001.state.pos` |
| **Subscribe** | Nhận tin nhắn từ các subject | Người vận hành subscribe vào `fleet.prod.veh.*.state.*` |
| **Publish + Subscribe** | Cả hai chiều | Gateway xử lý lệnh/phản hồi |
| **Deny** | Chặn rõ ràng các mẫu | Chặn phương tiện subscribe vào các phương tiện khác |

**Ký tự đại diện (wildcard) của subject trong quyền:**

```
fleet.prod.veh.V001.*        → Only this vehicle's subjects
fleet.prod.veh.*.state.*     → All vehicles' state (read-only)
fleet.prod.veh.*.cmd.*       → All vehicles' commands (operator)
fleet.prod.veh.*.evt.*       → All events (auditor)
```

Đây là cách chúng tôi đảm bảo rằng phương tiện chỉ có thể publish dữ liệu của chính mình và không thể mạo danh phương tiện khác. Xem [Đặt tên Subject]({{< relref "/fleet/subjects" >}}) để biết toàn bộ cấu trúc phân cấp.

---

## Mô hình Grants

Grants cho phép **các quyền chi tiết, có thể ủy quyền lại**. Người vận hành có thể cấp cho bên thứ ba quyền truy cập vào dữ liệu cụ thể mà không cần trao cho họ toàn quyền truy cập hệ thống.

```
Account Authority (root)
    │
    ├── Fleet Operator Account
    │       │
    │       ├── Grant: Emergency Services (incident-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.pos
    │       │       └── Subscribe: fleet.prod.veh.*.stream.video
    │       │       └── Expires: 4 hours
    │       │
    │       ├── Grant: Research Partner (data-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.env.*
    │       │       └── Deny: fleet.prod.veh.*.state.pos
    │       │       └── Expires: 30 days
    │       │
    │       └── Grant: Auditor (read-only, full)
    │               └── Subscribe: fleet.prod.veh.*.state.*
    │               └── Subscribe: fleet.prod.veh.*.evt.*
    │               └── Deny: fleet.prod.veh.*.cmd.*
    │               └── Expires: 7 days
```

**Đặc tính của grant:**

| Đặc tính | Lợi ích |
|----------|---------|
| **Giới hạn phạm vi** | Chỉ các subject/mẫu cụ thể—không hơn |
| **Giới hạn thời gian** | Tự động hết hạn, không có quyền truy cập bị lãng quên |
| **Có thể thu hồi** | Thu hồi tức thì mà không cần can thiệp vào bên được cấp |
| **Có thể ủy quyền lại** | Người vận hành có thể cấp quyền con cho người khác |
| **Được kiểm toán** | Mọi lần cấp grant đều được ghi nhật ký |

---

## Các kịch bản tích hợp bên thứ ba

### Tích hợp dịch vụ khẩn cấp

```
Scenario: Wildfire response, fire department needs drone video

Grant issued:
  - Subject: fleet.prod.veh.*.stream.video
  - Subject: fleet.prod.veh.*.state.pos
  - Duration: 4 hours (incident duration)
  - Audit: All access logged with responder ID

Result: Fire department connects with their credential, receives
real-time video and position from all drones in area.
Grant expires automatically when incident closes.
```

**Tại sao cách này hiệu quả:**
- Sở cứu hỏa nhận được đúng những gì họ cần—video trực tiếp và vị trí
- Không có quyền truy cập lệnh—không thể can thiệp vào hoạt động của drone
- Giới hạn thời gian—quyền truy cập kết thúc khi sự cố kết thúc
- Có thể kiểm toán—mọi khung hình họ nhận được đều được ghi nhật ký

### Chia sẻ dữ liệu với đối tác

```
Scenario: University research on environmental monitoring

Grant issued:
  - Subject: fleet.prod.veh.*.state.env.* (temp, humidity, air quality)
  - Deny: fleet.prod.veh.*.state.pos (no location data)
  - Duration: 30 days (research period)
  - Audit: Data access logged for compliance

Result: Researcher receives environmental sensor data without
knowing vehicle locations—privacy preserved.
```

**Tại sao cách này hiệu quả:**
- Nhà nghiên cứu nhận được dữ liệu môi trường cho nghiên cứu của họ
- Dữ liệu vị trí bị từ chối rõ ràng—quyền riêng tư được bảo vệ
- Dấu vết kiểm toán rõ ràng phục vụ tuân thủ nghiên cứu
- Tự động dọn dẹp sau thời gian nghiên cứu

### Tuân thủ quy định

```
Scenario: Aviation authority audit of flight operations

Grant issued:
  - Subject: fleet.prod.veh.*.state.* (all telemetry)
  - Subject: fleet.prod.veh.*.evt.* (all events)
  - Deny: fleet.prod.veh.*.cmd.* (no command capability)
  - Duration: 7 days (audit period)
  - Audit: Every query logged for audit trail

Result: Auditor gets complete read access to historical data
but cannot issue any commands to vehicles.
```

**Tại sao cách này hiệu quả:**
- Minh bạch hoàn toàn với kiểm toán viên—không có gì bị che giấu
- Chỉ đọc—kiểm toán viên không thể ảnh hưởng đến hoạt động
- Dấu vết kiểm toán đầy đủ về quyền truy cập của kiểm toán viên
- Giới hạn thời gian—không có lỗ hổng truy cập vĩnh viễn

### Tích hợp khách hàng

```
Scenario: Enterprise customer wants API access to their fleet data

Grant issued:
  - Subject: fleet.prod.veh.CUST-*.state.* (only their vehicles)
  - Subject: fleet.prod.veh.CUST-*.evt.*
  - Deny: fleet.prod.veh.OTHER-*.* (other customers' data)
  - Duration: Persistent (until revoked)
  - Audit: API access logged

Result: Customer builds their own dashboards using real-time
data from their vehicles—and only their vehicles.
```

### Quyền truy cập của nhà thầu phụ

```
Scenario: Maintenance team needs to run diagnostics on specific vehicles

Grant issued:
  - Subject: fleet.prod.veh.V001.cmd.diag.* (diagnostic commands only)
  - Subject: fleet.prod.veh.V001.state.health.*
  - Deny: fleet.prod.veh.V001.cmd.flight.* (no flight commands)
  - Duration: 8 hours (maintenance window)
  - Audit: All commands logged with technician ID

Result: Technician can run diagnostics but cannot arm or fly the vehicle.
Access expires after maintenance window.
```

---

## Đặc tính bảo mật

| Đặc tính | Triển khai |
|----------|----------------|
| **Zero Trust** | Không có sự tin tưởng ngầm định—mọi kết nối đều được xác thực |
| **Đặc quyền tối thiểu** | Thông tin xác thực chỉ chứa các quyền tối thiểu cần thiết |
| **Phòng thủ nhiều lớp** | Nhiều lớp: TLS + JWT + ACL theo subject |
| **Hoạt động ngoại tuyến** | Ủy quyền hoạt động mà không cần kết nối mạng |
| **Có thể kiểm toán** | Mọi truy cập đều được ghi nhật ký kèm danh tính của thông tin xác thực |
| **Có thể thu hồi** | Thu hồi tức thì qua danh sách thu hồi của tài khoản |

---

## Lợi ích vận hành

### Dành cho người vận hành đội phương tiện

- **Thêm/xóa bên thứ ba mà không cần thay đổi hệ thống** — Cấp hoặc thu hồi thông tin xác thực, không cần chỉnh sửa tệp cấu hình
- **Cấp quyền truy cập tạm thời mà không cần chia sẻ thông tin xác thực** — Mỗi bên có thông tin xác thực riêng
- **Kiểm toán ai đã truy cập gì và khi nào** — Nhật ký truy cập đầy đủ
- **Thu hồi quyền truy cập ngay lập tức khi cần** — Có hiệu lực từ lần kết nối tiếp theo

### Dành cho bên thứ ba

- **Tự kết nối** — Sử dụng thông tin xác thực được cung cấp, kết nối khi sẵn sàng
- **Phạm vi rõ ràng** — Biết chính xác những gì bạn có thể và không thể truy cập
- **Không cần hiểu kiến trúc nội bộ** — Chỉ cần kết nối và subscribe
- **Tự động hết hạn** — Bạn không cần dọn dẹp gì

### Dành cho tuân thủ

- **Dấu vết kiểm toán đầy đủ** — Mọi truy cập đều truy ngược được đến thông tin xác thực
- **Ranh giới quyền có thể chứng minh** — Chứng minh ai đã có quyền truy cập vào cái gì
- **Quyền truy cập có thời hạn cho kiểm toán viên** — Không có lỗ hổng truy cập vĩnh viễn
- **Không dùng chung thông tin xác thực** — Trách nhiệm giải trình của từng cá nhân

---

## Ghi chú triển khai

### Định dạng thông tin xác thực

NATS sử dụng **NKeys** và **JWT**:

- **Chữ ký Ed25519** — An toàn về mật mã, xác minh nhanh
- **Claims tự chứa** — Quyền được nhúng trong token
- **Không cần xác thực bên ngoài** — Máy chủ xác minh chữ ký cục bộ

### Thu hồi

Khi bạn cần thu hồi quyền truy cập:

1. Thêm thông tin xác thực vào danh sách thu hồi của tài khoản
2. Danh sách được lan truyền đến tất cả máy chủ trong cụm
3. Có hiệu lực từ lần thử kết nối tiếp theo
4. Không cần can thiệp vào hệ thống của bên được cấp

### Xoay vòng thông tin xác thực

Đối với các tích hợp dài hạn:

1. Cấp thông tin xác thực mới trước khi thông tin cũ hết hạn
2. Khoảng thời gian chồng lấn cho phép chuyển đổi liền mạch
3. Thu hồi thông tin xác thực cũ sau khi hoàn tất chuyển đổi
4. Không gián đoạn dịch vụ

---

## Tại sao điều này quan trọng với hệ thống thực tế

Ủy quyền truyền thống đòi hỏi:
- Máy chủ xác thực trung tâm (điểm lỗi đơn lẻ)
- Lệnh gọi mạng trên mỗi yêu cầu (độ trễ)
- Tích hợp phức tạp cho bên thứ ba (hàng tháng công sức)
- Dọn dẹp thủ công quyền truy cập đã hết hạn (thông tin xác thực bị lãng quên)

Ủy quyền phi tập trung của NATS mang lại:
- **Khả năng phục hồi** — Hoạt động ngoại tuyến, không có nút thắt trung tâm
- **Tốc độ** — Xác thực cục bộ, không cần khứ hồi
- **Tính mở rộng** — Thêm bên thứ ba trong vài phút, không phải vài tháng
- **An toàn** — Tự động hết hạn, thu hồi tức thì

Đây là cách xây dựng hệ thống vừa **bảo mật cao** vừa **rất thân thiện với người dùng**.

---

## Tài liệu liên quan

- [Đặt tên Subject]({{< relref "/fleet/subjects" >}}) — Cách cấu trúc phân cấp subject cho phép ACL chi tiết
- [Cấu trúc liên kết NATS]({{< relref "/fleet/nats-topology" >}}) — Cách tổ chức hạ tầng nhắn tin
- [Mô hình an toàn]({{< relref "/fleet/safety" >}}) — Cách ủy quyền lệnh khớp với hệ thống phân cấp an toàn
- [Giám sát]({{< relref "/company/surveillance" >}}) — Cách grants cho phép bên thứ ba truy cập để giám sát

---

## Tiếp theo

[Streams & Sự kiện →]({{< relref "/fleet/streams" >}})
