---
title: "Cổng kết nối phương tiện"
meta_title: "Thiết kế cổng kết nối phương tiện cho đội xe tự hành | Ubuntu Software"
description: "Kiến trúc Vehicle Gateway: chuyển đổi giao thức (MAVLink, CAN, ROS 2), giảm tần số lấy mẫu trạng thái, trích xuất sự kiện, thực thi lệnh và đối chiếu shadow."
image: "/images/robotics.svg"
draft: false
---

## Cầu nối giữa các giao thức phương tiện và NATS

Vehicle Gateway là một dịch vụ chạy trên máy tính biên của mỗi phương tiện, có nhiệm vụ chuyển đổi giữa các giao thức gốc của phương tiện (MAVLink, CAN bus, ROS 2) và hệ thống nhắn tin NATS được dùng để điều phối đội phương tiện.

---

## Tại sao cần Gateway?

Giao thức phương tiện và hệ thống nhắn tin của đội phương tiện phục vụ những mục đích khác nhau:

| Khía cạnh | Giao thức phương tiện | NATS |
|--------|------------------|------|
| **Phạm vi** | Một phương tiện | Toàn đội phương tiện |
| **Giao thức** | Nhị phân, nhỏ gọn | JSON, con người đọc được |
| **Tần suất** | Telemetry 100+ Hz | Giảm mẫu cho WAN |
| **Lưu trữ bền vững** | Không có | Luồng JetStream |
| **Mô hình** | Yêu cầu-phản hồi | Pub/sub có lưu trữ bền vững |

Gateway kết nối hai thế giới này, đảm nhận việc chuyển đổi giao thức, giới hạn tần suất và quản lý trạng thái.

---

## Các giao thức được hỗ trợ

| Giao thức | Loại phương tiện | Tích hợp |
|----------|--------------|-------------|
| **MAVLink** | Drone (PX4, ArduPilot) | UDP/Serial, phân tích cú pháp thông điệp |
| **CAN bus** | Ô tô, xe tải | SocketCAN, giải mã DBC |
| **J1939** | Phương tiện hạng nặng | Ánh xạ PGN/SPN |
| **ROS 2** | Nền tảng tự hành | Đăng ký topic |

Mỗi giao thức cần một adapter riêng, nhưng logic cốt lõi của gateway (giảm mẫu, sự kiện, lệnh, shadow) được dùng chung.

---

## Trách nhiệm cốt lõi

### 1. Tiếp nhận giao thức

Nhận và phân tích cú pháp các thông điệp giao thức gốc của phương tiện. Ví dụ với MAVLink (drone):

```go
// Receive MAVLink frames from mavlink-router
conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 14550})

for {
    buf := make([]byte, 1024)
    n, _, _ := conn.ReadFromUDP(buf)

    frame, _ := mavlink.Parse(buf[:n])

    switch msg := frame.Message().(type) {
    case *mavlink.Heartbeat:
        handleHeartbeat(msg)
    case *mavlink.GlobalPositionInt:
        handlePosition(msg)
    case *mavlink.BatteryStatus:
        handleBattery(msg)
    // ... handle other message types
    }
}
```

**Các loại thông điệp chính (ví dụ MAVLink):**

| Thông điệp MAVLink | Nội dung |
|-----------------|---------|
| `HEARTBEAT` | Chế độ, trạng thái arm, trạng thái hệ thống |
| `GLOBAL_POSITION_INT` | Vĩ độ, kinh độ, độ cao, vận tốc |
| `ATTITUDE` | Roll, pitch, yaw |
| `BATTERY_STATUS` | Điện áp, dòng điện, dung lượng còn lại |
| `GPS_RAW_INT` | Định vị GPS, số vệ tinh, HDOP |
| `SYS_STATUS` | Tải CPU, lỗi, tình trạng hoạt động |

**Tương đương trên CAN bus (phương tiện mặt đất):**

| Tín hiệu CAN | Nội dung |
|------------|---------|
| Vòng tua động cơ | Tốc độ động cơ hiện tại |
| Tốc độ xe | Tốc độ so với mặt đất từ cảm biến bánh xe |
| Vị trí GPS | Vĩ độ/kinh độ từ module telematics |
| Mức nhiên liệu | Phần trăm nhiên liệu trong bình |
| Mã DTC | Mã lỗi chẩn đoán |

### 2. Giảm mẫu trạng thái

Giảm tần suất telemetry để truyền qua WAN:

```go
// Position: downsample from 10Hz to 1Hz
positionSampler := NewDownsampler(100*time.Millisecond, func(msg *Position) {
    // Average positions over window
    return averagePosition(msg)
})

// Battery: only publish on change
batterySampler := NewChangeFilter(func(old, new *Battery) bool {
    return math.Abs(old.Remaining - new.Remaining) > 0.01
})
```

**Các chiến lược giảm mẫu:**

| Chiến lược | Trường hợp sử dụng | Ví dụ |
|----------|----------|---------|
| **Theo thời gian** | Trạng thái định kỳ | Vị trí ở 1Hz |
| **Theo thay đổi** | Giá trị rời rạc | Thay đổi chế độ |
| **Theo ngưỡng** | Thay đổi từ từ | Pin khi Δ > 1% |
| **Tổng hợp** | Dữ liệu tần suất cao | Min/max/trung bình trong một cửa sổ thời gian |

Telemetry tần suất đầy đủ được giữ lại trên phương tiện cho các hệ thống nhận thức. Chỉ dữ liệu đã giảm mẫu mới được truyền qua WAN.

### 3. Trích xuất sự kiện

Tạo sự kiện từ các chuyển đổi trạng thái:

```go
type EventDetector struct {
    previousState *VehicleState
}

func (e *EventDetector) Process(current *VehicleState) []Event {
    var events []Event

    // Detect arm state change
    if current.Armed && !e.previousState.Armed {
        events = append(events, Event{
            Type:      "armed",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"reason": "manual"},
        })
    }

    // Detect mode change
    if current.Mode != e.previousState.Mode {
        events = append(events, Event{
            Type:      "mode_change",
            Timestamp: time.Now(),
            Data: map[string]interface{}{
                "from": e.previousState.Mode,
                "to":   current.Mode,
            },
        })
    }

    // Detect failsafe
    if current.FailsafeActive && !e.previousState.FailsafeActive {
        events = append(events, Event{
            Type:      "failsafe",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"type": current.FailsafeType},
        })
    }

    e.previousState = current
    return events
}
```

**Các sự kiện được phát hiện:**

| Sự kiện | Điều kiện kích hoạt |
|-------|---------|
| `armed` | Trạng thái arm false → true |
| `disarmed` | Trạng thái arm true → false |
| `mode_change` | Chuyển đổi chế độ bay |
| `takeoff` | Trạng thái trên không false → true |
| `landed` | Trạng thái trên không true → false |
| `failsafe` | Failsafe được kích hoạt |
| `geofence` | Phát hiện vi phạm geofence |
| `battery.low` | Pin dưới ngưỡng |
| `battery.critical` | Pin ở mức nguy hiểm |

### 4. Thực thi lệnh theo chính sách

Nhận lệnh từ NATS và thực thi qua MAVLink:

```go
func (g *Gateway) handleCommand(cmd *Command) *CommandAck {
    // Validate command is allowed
    if !g.policy.Allows(cmd) {
        return &CommandAck{
            Status: "rejected",
            Error:  "command not allowed by policy",
        }
    }

    // Convert to MAVLink command
    mavCmd := cmd.ToMAVLink()

    // Send to flight controller
    if err := g.mavlink.Send(mavCmd); err != nil {
        return &CommandAck{
            Status: "failed",
            Error:  err.Error(),
        }
    }

    // Wait for MAVLink ACK
    mavAck, err := g.mavlink.WaitAck(mavCmd.CommandID, 5*time.Second)
    if err != nil {
        return &CommandAck{
            Status: "timeout",
            Error:  "no response from flight controller",
        }
    }

    return &CommandAck{
        Status: mavAck.Result.String(),
    }
}
```

**Thực thi chính sách:**

| Chính sách | Mô tả |
|--------|-------------|
| **Geofence** | Từ chối lệnh goto nằm ngoài ranh giới |
| **Giới hạn độ cao** | Giới hạn độ cao tối đa trong các lệnh |
| **Hạn chế chế độ** | Không cho phép một số chuyển đổi chế độ nhất định |
| **Giới hạn tần suất** | Ngăn chặn tình trạng gửi lệnh ồ ạt |
| **Xác thực** | Xác minh nguồn gửi lệnh |

### 5. Đồng bộ shadow

Đồng bộ trạng thái mong muốn với trạng thái thực tế của phương tiện:

```go
func (g *Gateway) reconcileLoop() {
    ticker := time.NewTicker(1 * time.Second)

    for range ticker.C {
        // Read desired state from KV
        desired, _ := g.kv.Get(g.desiredKey())

        // Compare with actual state
        actual := g.getCurrentState()

        // Generate commands to reconcile
        commands := g.reconcile(desired, actual)

        for _, cmd := range commands {
            g.executeCommand(cmd)
        }

        // Update reported state
        reported := g.buildReportedState(actual)
        g.kv.Put(g.reportedKey(), reported)
    }
}

func (g *Gateway) reconcile(desired, actual *State) []*Command {
    var commands []*Command

    // Mode reconciliation
    if desired.Mode != actual.Mode {
        commands = append(commands, &Command{
            Type: "set_mode",
            Data: map[string]interface{}{"mode": desired.Mode},
        })
    }

    // Geofence reconciliation
    if desired.GeofenceEnabled != actual.GeofenceEnabled {
        commands = append(commands, &Command{
            Type: "set_geofence",
            Data: map[string]interface{}{"enabled": desired.GeofenceEnabled},
        })
    }

    return commands
}
```

**Trạng thái shadow cho phép quản lý theo kiểu khai báo:**

- Người vận hành đội phương tiện đặt trạng thái mong muốn
- Gateway phát hiện các khác biệt
- Gateway phát lệnh để hội tụ về trạng thái mong muốn
- Gateway báo cáo trạng thái thực tế
- Lặp lại liên tục

---

## Kiến trúc

```
┌────────────────────────────────────────────────────────────────────┐
│                         Vehicle Gateway                             │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐     │
│  │   Protocol   │    │    State     │    │    NATS Client   │     │
│  │   Adapter    │───▶│   Machine    │───▶│    Publisher     │     │
│  └──────────────┘    └──────────────┘    └──────────────────┘     │
│         │                   │                      │               │
│         │            ┌──────▼──────┐               │               │
│         │            │   Event     │               │               │
│         │            │  Detector   │───────────────┤               │
│         │            └─────────────┘               │               │
│         │                                          │               │
│         │            ┌─────────────┐               │               │
│         │            │   Shadow    │◀──────────────┤               │
│         │            │ Reconciler  │               │               │
│         │            └──────┬──────┘               │               │
│         │                   │                      │               │
│         │            ┌──────▼──────┐    ┌─────────▼────────┐      │
│         │            │  Command    │    │   NATS Client    │      │
│         │◀───────────│  Executor   │◀───│   Subscriber     │      │
│         │            └─────────────┘    └──────────────────┘      │
│  ┌──────▼──────┐                                                   │
│  │   Protocol  │                                                   │
│  │   Sender    │                                                   │
│  └─────────────┘                                                   │
└────────────────────────────────────────────────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────┐                       ┌─────────────────┐
│  Vehicle Control│                       │   NATS Leaf     │
│ (Pixhawk/ECU/   │                       │   (localhost)   │
│  ROS 2 node)    │                       │                 │
└─────────────────┘                       └─────────────────┘
```

**Các adapter giao thức:**
- **MAVLink Adapter** — Dành cho drone (PX4, ArduPilot)
- **CAN Adapter** — Dành cho ô tô/xe tải (SocketCAN + DBC)
- **ROS 2 Adapter** — Dành cho các nền tảng dựa trên ROS

---

## Chi tiết hiện thực

### Phụ thuộc

```go
import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "github.com/bluenviern/go-mavlink/v2"
)
```

**Thư viện cốt lõi:**

| Thư viện | Mục đích |
|---------|---------|
| **nats.go** | Client NATS, JetStream, KV |
| **go-mavlink** | Hiện thực giao thức MAVLink |
| **slog** | Ghi log có cấu trúc |

### Cấu hình

```yaml
# gateway.yaml
vehicle_id: VID-001
environment: prod

mavlink:
  local_addr: ":14550"
  system_id: 1
  component_id: 1

nats:
  url: "nats://localhost:4222"
  credentials: "/etc/gateway/vehicle.creds"

sampling:
  position_hz: 1
  attitude_hz: 1
  battery_change_threshold: 0.01

policy:
  max_altitude: 120
  geofence_file: "/etc/gateway/geofence.json"
```

### Triển khai

```yaml
# systemd service
[Unit]
Description=Vehicle Gateway
After=network.target nats.service

[Service]
Type=simple
ExecStart=/usr/local/bin/vehicle-gateway --config /etc/gateway/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## Xử lý lỗi

### Lỗi MAVLink

| Lỗi | Phản hồi |
|-------|----------|
| **Mất kết nối** | Kết nối lại với backoff |
| **Lỗi phân tích cú pháp** | Ghi log và bỏ qua frame |
| **Hết thời gian chờ** | Thử lại hoặc báo lệnh thất bại |

### Lỗi NATS

| Lỗi | Phản hồi |
|-------|----------|
| **Mất kết nối** | NATS cục bộ tiếp tục hoạt động, kết nối lại với hub |
| **Publish thất bại** | Lưu đệm cục bộ, thử lại |
| **Lỗi stream** | Ghi log, cảnh báo, tiếp tục |

### Lỗi lệnh

| Lỗi | Phản hồi |
|-------|----------|
| **Vi phạm chính sách** | Từ chối ngay lập tức |
| **MAVLink từ chối** | Trả về ACK thất bại |
| **Hết thời gian chờ** | Trả về ACK hết thời gian chờ |

---

## Chỉ số

Gateway cung cấp các chỉ số Prometheus:

```
# MAVLink
gateway_mavlink_messages_received_total{type="heartbeat"}
gateway_mavlink_messages_sent_total{type="command_long"}
gateway_mavlink_parse_errors_total

# NATS
gateway_nats_messages_published_total{subject="state.position"}
gateway_nats_messages_received_total{subject="cmd.takeoff"}
gateway_nats_publish_errors_total

# Commands
gateway_commands_received_total{type="takeoff"}
gateway_commands_executed_total{type="takeoff",result="success"}
gateway_command_latency_seconds{type="takeoff"}

# Shadow
gateway_shadow_reconciliations_total
gateway_shadow_commands_issued_total
```

---

## Tóm tắt

| Trách nhiệm | Đầu vào | Đầu ra |
|----------------|-------|--------|
| **Tiếp nhận giao thức** | Thông điệp phương tiện | Telemetry đã phân tích |
| **Giảm mẫu trạng thái** | Telemetry 100Hz | Trạng thái 1Hz |
| **Trích xuất sự kiện** | Chuyển đổi trạng thái | Sự kiện rời rạc |
| **Thực thi lệnh** | Lệnh NATS | Lệnh phương tiện |
| **Đồng bộ shadow** | Trạng thái mong muốn | Lệnh hội tụ |

Vehicle Gateway là thành phần then chốt giúp vận hành ở quy mô đội phương tiện trở nên khả thi, đồng thời vẫn giữ nguyên các đảm bảo an toàn của hệ thống điều khiển phương tiện bên dưới.

---

## Tài liệu liên quan

- [Các nền tảng được hỗ trợ]({{< relref "/fleet/platforms" >}}) — Tổng quan về tất cả các loại phương tiện
- [Nền tảng drone]({{< relref "/fleet/platforms/drones" >}}) — Chi tiết tích hợp MAVLink
- [Phương tiện mặt đất]({{< relref "/fleet/platforms/ground" >}}) — Tích hợp CAN bus và J1939

---

## Tiếp theo

[Mô hình an toàn →]({{< relref "/fleet/safety" >}})
