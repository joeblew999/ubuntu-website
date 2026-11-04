---
title: "Mô hình an toàn"
meta_title: "Kiến trúc An toàn Đội phương tiện Tự hành | Ubuntu Software"
description: "Mô hình an toàn cho đội xe tự hành: quyền can thiệp thủ công, cơ chế an toàn dự phòng trên phương tiện, cách ly mạng khỏi các vòng điều khiển, và suy giảm chức năng có kiểm soát."
image: "/images/robotics.svg"
draft: false
---

## An toàn ngay từ thiết kế

Vận hành đội phương tiện tự hành ở quy mô lớn đòi hỏi kiến trúc an toàn nghiêm ngặt. Thiết kế này đảm bảo rằng **sự cố mạng, lỗi phần mềm và trục trặc hệ thống không bao giờ ảnh hưởng đến an toàn của phương tiện**.

Các nguyên tắc an toàn cốt lõi áp dụng cho **mọi loại phương tiện**—drone, ô tô, xe tải và AGV. Chi tiết cụ thể có thể khác nhau (RC so với vô lăng, RTL so với dừng tại chỗ), nhưng kiến trúc vẫn nhất quán.

---

## Nguyên tắc cốt lõi

**NATS không bao giờ nằm trong vòng điều khiển thời gian thực.**

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SAFETY HIERARCHY                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   1. MANUAL OVERRIDE      ─────────────────────────▶  HIGHEST       │
│      Human control (RC, steering, e-stop), always available          │
│                                                                      │
│   2. VEHICLE FAILSAFES    ─────────────────────────▶  HIGH          │
│      Hardware-enforced, local to vehicle                             │
│                                                                      │
│   3. GATEWAY POLICY       ─────────────────────────▶  MEDIUM        │
│      Software enforcement, on-vehicle                                │
│                                                                      │
│   4. FLEET COMMANDS       ─────────────────────────▶  LOWEST        │
│      Network-dependent, advisory                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

Các lớp cao hơn luôn có thể ghi đè các lớp thấp hơn. Các lệnh phụ thuộc vào mạng có quyền hạn thấp nhất.

---

## Điều khiển thủ công là quyền hạn cao nhất

Người vận hành có **quyền ghi đè vô điều kiện**—dù thông qua bộ phát RC, vô lăng hay nút dừng khẩn cấp.

### Ghi đè trên drone (RC)

| Khía cạnh | Triển khai |
|--------|----------------|
| **Giao thức** | ExpressLRS trực tiếp đến bộ thu |
| **Đường truyền** | TX → RX → Pixhawk (không có phần mềm trong vòng lặp) |
| **Chuyển chế độ** | Failsafe phần cứng về manual/stabilize |
| **Công tắc ngắt** | Ngắt động cơ (disarm) ngay lập tức |

**RC không bao giờ phụ thuộc vào:**

- Jetson đang trực tuyến
- Kết nối NATS
- Vehicle Gateway
- Bất kỳ thành phần phần mềm nào

Nếu mọi máy tính đều hỏng, phi công vẫn giữ toàn quyền điều khiển thủ công.

#### Failsafe của ExpressLRS

Khi mất tín hiệu RC:

```
1. RX detects signal loss (100ms)
2. RX sends failsafe values to Pixhawk
3. Pixhawk enters RC_LOSS failsafe mode
4. Configured action executes (RTL, land, etc.)
```

Hành vi failsafe được cấu hình trong PX4, không phụ thuộc vào bất kỳ phần mềm quản lý đội nào.

### Ghi đè trên phương tiện mặt đất

| Khía cạnh | Ô tô/Xe tải | AGV |
|--------|-------------|------|
| **Chính** | Vô lăng + bàn đạp | Nút dừng khẩn cấp |
| **Đường truyền** | Kết nối vật lý đến bộ chấp hành | Ngắt phần cứng đến động cơ |
| **Dự phòng** | Tắt khóa điện / về số N | Dừng khẩn cấp không dây |
| **Mạng** | Không bao giờ nằm trong đường điều khiển | Không bao giờ nằm trong đường điều khiển |

**Ghi đè trên phương tiện mặt đất không bao giờ phụ thuộc vào:**

- Máy tính trên xe đang trực tuyến
- Kết nối NATS
- Bất kỳ phần mềm tự hành nào

Nếu mọi máy tính đều hỏng, người lái/người vận hành vẫn giữ quyền điều khiển vật lý hoặc phương tiện dừng lại an toàn.

---

## Hệ thống điều khiển phương tiện thực thi các failsafe

Hệ thống điều khiển của phương tiện triển khai các failsafe ngay tại chỗ—không cần mạng.

### Failsafe của drone (PX4)

Bộ điều khiển bay triển khai nhiều failsafe:

#### Failsafe khi mất RC

| Tham số | Giá trị | Hành động |
|-----------|-------|--------|
| `COM_RC_LOSS_T` | 0.5s | Thời gian chờ trước khi kích hoạt failsafe |
| `NAV_RCL_ACT` | 2 | RTL khi mất RC |

#### Mất liên kết dữ liệu

| Tham số | Giá trị | Hành động |
|-----------|-------|--------|
| `COM_DL_LOSS_T` | 10s | Thời gian chờ kết nối GCS |
| `NAV_DLL_ACT` | 0 | Tiếp tục nhiệm vụ (liên kết không quan trọng) |

#### Failsafe pin

| Mức | Tham số | Giá trị | Hành động |
|-------|-----------|-------|--------|
| Thấp | `BAT_LOW_THR` | 0.25 | Cảnh báo |
| Nghiêm trọng | `BAT_CRIT_THR` | 0.15 | RTL |
| Khẩn cấp | `BAT_EMERGEN_THR` | 0.07 | Hạ cánh ngay lập tức |

#### Hàng rào địa lý (Geofence)

| Tham số | Giá trị | Hành động |
|-----------|-------|--------|
| `GF_ACTION` | 3 | RTL khi vi phạm |
| `GF_MAX_HOR_DIST` | thay đổi | Giới hạn ngang |
| `GF_MAX_VER_DIST` | 120m | Giới hạn độ cao |

#### Mất định vị

| Tham số | Giá trị | Hành động |
|-----------|-------|--------|
| `COM_POS_FS_DELAY` | 1s | Thời gian chờ định vị |
| `COM_POSCTL_NAVL` | 0 | Hạ cánh khi mất định vị |

**Mọi failsafe đều thực thi trên Pixhawk**—không phụ thuộc bên ngoài.

### Failsafe của phương tiện mặt đất

Phương tiện mặt đất triển khai các cơ chế an toàn tương đương:

| Tác nhân kích hoạt | Phản ứng của drone | Phản ứng của phương tiện mặt đất |
|---------|---------------|------------------------|
| **Mất liên lạc** | Quay về điểm cất cánh (RTL) | Dừng tại chỗ / tấp vào lề |
| **Hỏng cảm biến** | Hạ cánh ngay lập tức | Dừng an toàn, kích hoạt phanh đỗ |
| **Vi phạm geofence** | RTL tại ranh giới | Dừng tại ranh giới |
| **Kích hoạt dừng khẩn cấp** | Ngắt động cơ | Phanh ngay lập tức |
| **Pin/nhiên liệu thấp** | RTL rồi hạ cánh | Quay về kho / dừng an toàn |
| **Phát hiện chướng ngại vật** | Tránh hoặc bay lơ lửng | Dừng và chờ |

**Failsafe của phương tiện mặt đất thường bao gồm:**

- **Chướng ngại vật qua lidar/radar** — Dừng khẩn cấp nếu đường đi bị chặn
- **Giới hạn tốc độ** — Tốc độ tối đa dựa trên môi trường
- **Hạn chế khu vực** — Các vùng cấm được thực thi tại chỗ
- **Bộ định thời watchdog** — Dừng nếu phần mềm tự hành bị treo
- **Dự phòng phanh** — Nhiều hệ thống phanh độc lập

**Mọi failsafe đều thực thi tại chỗ**—không phụ thuộc vào mạng.

---

## NATS chỉ dùng để giám sát

Lớp mạng đảm nhận:

| Chức năng | Mô tả |
|----------|-------------|
| **Thu thập telemetry** | Trạng thái phương tiện cho dashboard |
| **Ghi log sự kiện** | Nhật ký kiểm toán trong JetStream |
| **Điều phối nhiệm vụ** | Lập kế hoạch cấp đội |
| **Shadow state** | Đồng bộ trạng thái mong muốn/được báo cáo |

Lớp mạng **không đảm nhận:**

| Chức năng | Lý do |
|----------|---------|
| **Điều khiển tư thế** | Độ trễ không chấp nhận được, quan trọng về an toàn |
| **Lệnh động cơ** | Phải nằm cục bộ tại FCU |
| **Hợp nhất cảm biến** | Yêu cầu thời gian thực |
| **Quyết định failsafe** | Phải hoạt động khi không có mạng |

### Xác thực lệnh

Khi lệnh đến qua NATS:

```go
func (g *Gateway) validateCommand(cmd *Command) error {
    // Check geofence
    if cmd.Type == "goto" {
        if !g.geofence.Contains(cmd.Target) {
            return errors.New("target outside geofence")
        }
    }

    // Check altitude
    if cmd.Altitude > g.policy.MaxAltitude {
        return errors.New("altitude exceeds limit")
    }

    // Check vehicle state
    if cmd.RequiresArmed && !g.state.Armed {
        return errors.New("vehicle not armed")
    }

    // Check rate limiting
    if g.rateLimiter.Exceeded(cmd.Type) {
        return errors.New("command rate exceeded")
    }

    return nil
}
```

Lệnh được **xác thực trước khi thực thi**—nhưng ngay cả khi một lệnh sai được thực thi, các failsafe của PX4 vẫn cung cấp lớp bảo vệ cuối cùng.

---

## Suy giảm có kiểm soát

Hệ thống suy giảm một cách an toàn khi các thành phần gặp sự cố:

### Jetson gặp sự cố

```
Jetson offline
    └─▶ No perception, no AI inference
    └─▶ Pi continues MAVLink routing
    └─▶ RC control fully operational
    └─▶ PX4 failsafes active
    └─▶ Vehicle lands safely or continues mission
```

### Mất kết nối NATS

```
WAN disconnected
    └─▶ Leaf node buffers locally
    └─▶ Gateway continues operation
    └─▶ Vehicle follows current mission
    └─▶ State syncs when connection returns
```

### Vehicle Gateway bị sập

```
Gateway process dies
    └─▶ systemd restarts within 5 seconds
    └─▶ PX4 continues autonomous flight
    └─▶ RC override always available
    └─▶ No fleet commands until restart
```

### Cụm Hub gặp sự cố

```
Regional hub down
    └─▶ Leaf nodes queue messages
    └─▶ Vehicles continue operating
    └─▶ RC and PX4 failsafes unaffected
    └─▶ Fleet visibility lost temporarily
```

### Ngăn chặn sự cố dây chuyền

| Sự cố | Phạm vi ảnh hưởng | Biện pháp giảm thiểu |
|---------|--------------|------------|
| Một phương tiện | Chỉ phương tiện đó | Failsafe để hạ cánh |
| Nút Hub | Không ảnh hưởng (được sao chép) | Cụm tự phục hồi |
| Cụm Hub | Điều phối đội khu vực | Phương tiện tiếp tục tự hành |
| Tầng toàn cầu | Khả năng quan sát liên khu vực | Hoạt động khu vực vẫn tiếp tục |

---

## Quy trình an toàn vận hành

### Danh sách kiểm tra trước chuyến bay

Trước khi arm bất kỳ phương tiện nào trong đội:

1. **Xác minh liên kết RC** — Phi công xác nhận quyền điều khiển
2. **Kiểm tra failsafe** — Xác minh hành vi RTL/hạ cánh
3. **Đã tải geofence** — Ranh giới đã được xác nhận trong PX4
4. **Ngưỡng pin** — Các mức failsafe phù hợp
5. **Kiểm tra thời tiết** — Gió trong giới hạn
6. **Không phận thông thoáng** — Không có xung đột

### Quy trình khẩn cấp

| Tình huống khẩn cấp | Phản ứng |
|-----------|----------|
| **Mất tầm nhìn** | Kích hoạt RTL qua RC |
| **Mất kiểm soát (flyaway)** | Công tắc ngắt RC hoặc geofence kích hoạt RTL |
| **Xung đột nhiều phương tiện** | Lệnh hạ cánh toàn đội |
| **Mất liên lạc hoàn toàn** | Failsafe RC_LOSS của PX4 |
| **Cháy/rơi** | Disarm qua RC, gọi dịch vụ khẩn cấp |

### Chuỗi dừng (Kill Chain)

Nhiều phương pháp để dừng một phương tiện:

```
1. RC kill switch      → Immediate motor stop
2. RC mode switch      → Land mode
3. GCS command         → RTL or land
4. Fleet command       → RTL or land
5. Geofence breach     → Auto-RTL
6. Battery critical    → Auto-land
```

Luôn có sẵn ít nhất ba phương pháp (RC, failsafe PX4, timeout cục bộ).

---

## Các cân nhắc về chứng nhận

Kiến trúc này hỗ trợ tuân thủ quy định:

| Yêu cầu | Triển khai |
|-------------|----------------|
| **Phi công chỉ huy** | Ghi đè RC mọi lúc |
| **Hành vi failsafe** | Failsafe được chứng nhận của PX4 |
| **Geofencing** | Ranh giới được thực thi bằng phần cứng |
| **Ghi log** | Nhật ký kiểm toán JetStream |
| **Remote ID** | Phát sóng qua mô-đun riêng |

**Lưu ý:** Yêu cầu pháp lý khác nhau tùy theo khu vực tài phán. Kiến trúc này cung cấp nền tảng kỹ thuật—việc tuân thủ cụ thể đòi hỏi các biện pháp bổ sung theo quy định địa phương.

---

## Kiểm thử hệ thống an toàn

### Kiểm thử đơn vị

```go
func TestGeofenceRejection(t *testing.T) {
    gw := NewGateway(config)
    gw.SetGeofence(testGeofence)

    cmd := &Command{
        Type:   "goto",
        Target: Position{Lat: 0, Lon: 0}, // Outside geofence
    }

    err := gw.validateCommand(cmd)
    assert.ErrorContains(t, err, "outside geofence")
}
```

### Kiểm thử tích hợp

```
1. Simulate RC loss → Verify RTL activates
2. Simulate battery drain → Verify land at threshold
3. Simulate geofence breach → Verify return behavior
4. Simulate Jetson crash → Verify continued flight
5. Simulate NATS disconnect → Verify local operation
```

### Kiểm thử bay

Trước khi triển khai đội:

1. **Xác thực failsafe trên từng phương tiện**
2. **Xác minh phản ứng khi mất RC**
3. **Kiểm thử vi phạm geofence**
4. **Xác thực failsafe pin**
5. **Kiểm thử điều phối nhiều phương tiện**

---

## Tóm tắt

| Lớp | Quyền hạn | Drone | Phương tiện mặt đất |
|-------|-----------|--------|-----------------|
| **Thủ công** | Cao nhất | Bộ phát RC | Vô lăng / Dừng khẩn cấp |
| **Phương tiện** | Cao | Failsafe PX4 | ECU / bộ điều khiển an toàn |
| **Gateway** | Trung bình | Thực thi chính sách | Thực thi chính sách |
| **Đội** | Thấp nhất | Chỉ điều phối | Chỉ điều phối |

Mô hình an toàn đảm bảo:

- **Mạng không bao giờ là thành phần quan trọng về an toàn**
- **Người vận hành luôn có quyền ghi đè**
- **Failsafe phần cứng thực thi tại chỗ**
- **Sự cố được suy giảm có kiểm soát**
- **Không có điểm đơn lẻ nào gây ảnh hưởng toàn đội**

---

## Tài liệu liên quan

- [Các nền tảng được hỗ trợ]({{< relref "/fleet/platforms" >}}) — Tổng quan về mọi loại phương tiện
- [Nền tảng Drone]({{< relref "/fleet/platforms/drones" >}}) — Chi tiết an toàn dành riêng cho drone
- [Phương tiện mặt đất]({{< relref "/fleet/platforms/ground" >}}) — Chi tiết an toàn cho phương tiện mặt đất

---

## Đọc thêm

- [Cấu hình an toàn PX4](https://docs.px4.io/main/en/config/safety.html)
- [Bộ thông điệp chung MAVLink](https://mavlink.io/en/messages/common.html)
- [Mô hình bảo mật NATS](https://docs.nats.io/running-a-nats-service/configuration/securing_nats)

---

## Quay lại tổng quan

[← Kiến trúc đội phương tiện tự hành]({{< relref "/fleet" >}})
