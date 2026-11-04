---
title: "Cấu trúc liên kết NATS"
meta_title: "Cấu trúc liên kết NATS JetStream cho đội máy bay không người lái | Ubuntu Software"
description: "Cấu trúc liên kết NATS phân cấp cho đội bay 1.000 drone: leaf node trên từng phương tiện, cụm hub khu vực, cô lập WAN và phát lại bản sao số trên toàn đội bay."
image: "/images/robotics.svg"
draft: false
---

## Kiến trúc NATS cho quy mô đội bay

Cấu trúc liên kết nhắn tin quyết định cách dữ liệu đo xa (telemetry) luân chuyển, nơi trạng thái được lưu trữ và cách hệ thống hoạt động khi kết nối suy giảm. Kiến trúc này mở rộng từ hàng chục đến hàng nghìn phương tiện mà vẫn duy trì độ trễ thấp và độ tin cậy cao.

---

## Thách thức

Hệ thống nhắn tin cho đội bay phải giải quyết các yêu cầu mâu thuẫn nhau:

| Yêu cầu | Thách thức |
|-------------|-----------|
| **Độ trễ thấp** | Giao tiếp giữa các phương tiện cần tính bằng mili giây, không phải giây |
| **Khả năng chống chịu WAN** | Kết nối di động bị rớt, độ trễ tăng đột biến, băng thông dao động |
| **Lưu trữ trạng thái bền vững** | Dữ liệu bản sao số (digital twin) phải tồn tại qua các lần khởi động lại và kết nối lại |
| **Quy mô** | Hàng nghìn bên xuất bản và bên đăng ký đồng thời |
| **Cô lập** | Sự cố của một phương tiện không được lan sang toàn đội bay |

Các kiến trúc truyền thống buộc phải đánh đổi. Broker tập trung làm tăng độ trễ. Mạng lưới ngang hàng (peer-to-peer) không mở rộng được. Các cách tiếp cận chỉ ở biên (edge-only) làm mất khả năng quan sát toàn đội bay.

---

## Cấu trúc liên kết phân cấp

Giải pháp là một **hệ phân cấp ba tầng**:

```
┌─────────────────────────────────────────────────────────────────┐
│                         GLOBAL (optional)                        │
│                    Cross-region replication                      │
│                    Fleet-wide aggregation                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  REGIONAL HUB   │ │  REGIONAL HUB   │ │  REGIONAL HUB   │
│   Americas      │ │   Europe        │ │   Asia-Pacific  │
│  3-node cluster │ │  3-node cluster │ │  3-node cluster │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
    ┌────┼────┐         ┌────┼────┐         ┌────┼────┐
    │    │    │         │    │    │         │    │    │
    ▼    ▼    ▼         ▼    ▼    ▼         ▼    ▼    ▼
┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐
│LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF │
│VID-1││VID-2││VID-3││VID-4││VID-5││VID-6││VID-7││VID-8││VID-9│
└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘
```

---

## Tầng 1: Nút lá trên phương tiện

Mỗi drone chạy một **nút lá NATS (leaf node)** trên Jetson:

| Khía cạnh | Cấu hình |
|--------|---------------|
| **Tiến trình** | nats-server ở chế độ leaf |
| **Lưu trữ** | JetStream cục bộ trên SSD/eMMC |
| **Kết nối** | Một kết nối ngược dòng duy nhất tới hub khu vực |
| **Tài nguyên** | ~50MB RAM, CPU tối thiểu |

**Những gì nút lá cung cấp:**

### Pub/Sub cục bộ

Các dịch vụ trên phương tiện giao tiếp thông qua phiên bản NATS cục bộ:

```
Vehicle Gateway ──publish──▶ fleet.prod.veh.VID-001.state.position
                                       │
                                       ▼
                            ┌──────────────────┐
                            │   Local NATS     │
                            │   (Leaf Node)    │
                            └──────────────────┘
                                       │
                                       ▼
AI Perception ◀──subscribe── fleet.prod.veh.VID-001.state.position
```

Tin nhắn giữa các dịch vụ trên phương tiện không bao giờ rời khỏi phương tiện. Độ trễ dưới một mili giây.

### Lưu và chuyển tiếp (Store-and-Forward)

Khi kết nối ngược dòng bị rớt, nút lá sẽ:

1. Tiếp tục nhận các bản tin xuất bản từ các dịch vụ cục bộ
2. Lưu tin nhắn vào JetStream cục bộ
3. Xếp hàng các tin nhắn gửi đi dành cho hub
4. Kết nối lại và phát lại khi kết nối được khôi phục

Điều này có nghĩa là **phương tiện vẫn tiếp tục hoạt động khi mất kết nối**. Trạng thái sẽ được đồng bộ khi liên kết phục hồi.

### Cô lập WAN

Nút lá đóng vai trò như một vùng đệm giữa phương tiện và WAN:

- **Áp lực ngược (backpressure)** — Nếu hub không theo kịp, nút lá sẽ đệm dữ liệu cục bộ
- **Lọc** — Chỉ các subject đã đăng ký mới đi qua liên kết WAN
- **Nén** — Giao thức NATS nén hiệu quả
- **Kết nối lại** — Tự động kết nối lại với thời gian chờ tăng theo cấp số nhân (exponential backoff)

---

## Tầng 2: Cụm hub khu vực

Các hub khu vực tập hợp phương tiện theo khu vực địa lý:

| Khía cạnh | Cấu hình |
|--------|---------------|
| **Triển khai** | Cụm NATS 3 nút (tối thiểu) |
| **Vị trí** | Vùng đám mây hoặc trung tâm dữ liệu biên |
| **Lưu trữ** | JetStream trên SSD tốc độ cao |
| **Dung lượng** | Hàng trăm kết nối leaf trên mỗi cụm |

**Tại sao theo khu vực:**

- **Độ trễ** — Các nút lá kết nối tới các hub ở gần
- **Băng thông** — Theo mặc định, dữ liệu đo xa được giữ trong khu vực
- **Tuân thủ** — Đáp ứng yêu cầu lưu trú dữ liệu
- **Cô lập sự cố** — Sự cố ở một khu vực không ảnh hưởng đến các khu vực khác

**Trách nhiệm của hub:**

| Chức năng | Mô tả |
|----------|-------------|
| **Tổng hợp nút lá** | Chấp nhận kết nối từ các nút lá trên phương tiện |
| **Lưu trữ stream** | Lưu trữ bền vững dữ liệu bản sao số trong JetStream |
| **Hỗ trợ consumer** | Phục vụ bảng điều khiển đội bay, API, phân tích |
| **Định tuyến liên khu vực** | Chuyển tiếp lên tầng toàn cầu khi cần |

### Cấu hình cụm

Một cụm 3 nút cung cấp:

- **Tính sẵn sàng cao** — Chịu được sự cố của một nút đơn lẻ
- **Sao chép dữ liệu** — Các stream JetStream được sao chép trên các nút
- **Phân phối tải** — Các kết nối leaf được cân bằng giữa các nút

```yaml
# Example hub cluster configuration
cluster:
  name: hub-americas
  routes:
    - nats://hub-americas-1:6222
    - nats://hub-americas-2:6222
    - nats://hub-americas-3:6222

jetstream:
  store_dir: /data/jetstream
  max_memory_store: 4GB
  max_file_store: 100GB

leafnodes:
  port: 7422
  authorization:
    users:
      - user: vehicle
        password: $VEHICLE_PASSWORD
        allowed_connection_types: ["LEAFNODE"]
```

---

## Tầng 3: Toàn cầu (Tùy chọn)

Đối với các đội bay trải rộng trên nhiều khu vực, tầng toàn cầu cho phép:

| Chức năng | Mô tả |
|----------|-------------|
| **Phản chiếu liên khu vực** | Sao chép stream giữa các hub khu vực |
| **Tổng hợp toàn cầu** | Bảng điều khiển và phân tích trên toàn đội bay |
| **Định tuyến lệnh** | Gửi lệnh tới phương tiện ở bất kỳ khu vực nào |
| **Khôi phục sau thảm họa** | Chuyển đổi dự phòng giữa các khu vực |

**Các phương án triển khai:**

1. **NATS Supercluster** — Liên kết các cụm khu vực thông qua kết nối gateway
2. **JetStream Mirroring** — Sao chép các stream cụ thể về một vị trí trung tâm
3. **Cấp ứng dụng** — Đồng bộ tùy chỉnh giữa các API khu vực

Hầu hết các đội bay hoạt động trong một khu vực duy nhất và ban đầu không cần tầng toàn cầu.

---

## Lợi ích của cấu trúc liên kết này

### Pub/Sub cục bộ độ trễ thấp

Giao tiếp trên phương tiện có độ trễ **dưới một mili giây**:

```
MAVLink message received
    └─▶ Vehicle Gateway publishes to local NATS
        └─▶ AI Perception subscribes from local NATS
            └─▶ Total latency: <1ms
```

Không có vòng truyền qua WAN đối với giao tiếp trên phương tiện.

### Cô lập WAN

Hoạt động của phương tiện **vẫn tiếp tục khi mất kết nối**:

- Các hệ thống nhận thức tiếp tục chạy
- Trạng thái được tích lũy cục bộ
- Lệnh được xếp hàng chờ thực thi
- Việc kết nối lại tự động đồng bộ

### Phát lại bản sao số trên toàn đội bay

JetStream cho phép **truy vấn theo thời gian** trên toàn đội bay:

```
# Replay vehicle state from 10 minutes ago
nats stream get TWIN_STATE --start-time="2024-01-15T10:00:00Z"

# Subscribe to real-time state updates
nats sub "fleet.prod.veh.*.state.position"
```

Lịch sử trạng thái của mọi phương tiện đều có thể truy vấn từ hub khu vực.

### Kiến trúc có khả năng mở rộng

Cấu trúc liên kết mở rộng theo chiều ngang:

| Điểm mở rộng | Cách tiếp cận |
|-------------|----------|
| **Thêm phương tiện** | Thêm kết nối leaf vào hub hiện có |
| **Thông lượng cao hơn** | Thêm nút vào cụm hub |
| **Thêm khu vực** | Triển khai thêm các hub khu vực |
| **Phạm vi toàn cầu** | Kết nối các hub qua supercluster |

---

## Luồng kết nối

Khi một phương tiện khởi động:

```
1. Jetson boots, NATS leaf node starts
2. Leaf node reads hub address from config
3. TLS connection established to regional hub
4. Authentication via credentials
5. Leaf node advertises local subjects
6. Hub routes relevant subscriptions to leaf
7. Bidirectional message flow begins
```

Khi kết nối bị rớt:

```
1. Leaf node detects connection loss
2. Local pub/sub continues uninterrupted
3. Outbound messages queue locally
4. Reconnection attempts with backoff
5. On reconnect, queued messages replay
6. Stream state synchronizes
```

---

## Bảo mật

Mọi kết nối đều được xác thực và mã hóa:

| Kết nối | Bảo mật |
|------------|----------|
| **Leaf → Hub** | TLS 1.3, xác thực bằng thông tin xác thực |
| **Cụm hub** | TLS, các tuyến cụm được xác thực |
| **Hub → Toàn cầu** | TLS, xác thực gateway |

Thông tin xác thực được cấp riêng cho từng phương tiện, cho phép:

- Thu hồi quyền của từng phương tiện riêng lẻ
- Nhật ký kiểm toán cho từng phương tiện
- Giới hạn tốc độ cho từng kết nối

Mô hình phân quyền hỗ trợ **bảo mật phi tập trung với cơ chế cấp quyền (grants)**—bên thứ ba có thể nhận quyền truy cập giới hạn theo phạm vi và thời gian mà không cần điều phối trung tâm. Xem [Phân quyền & Cấp quyền]({{< relref "/fleet/authorization" >}}) để biết chi tiết về quản lý thông tin xác thực và tích hợp bên thứ ba.

---

## Tóm tắt

| Tầng | Mục đích | Triển khai |
|------|---------|------------|
| **Leaf (Phương tiện)** | Pub/sub cục bộ, lưu và chuyển tiếp | Trên Jetson |
| **Hub (Khu vực)** | Tổng hợp, lưu trữ bền vững, consumer | Trung tâm dữ liệu đám mây/biên |
| **Toàn cầu (Tùy chọn)** | Liên khu vực, góc nhìn toàn đội bay | Đám mây trung tâm |

Cấu trúc liên kết này đảm bảo các phương tiện hoạt động độc lập trong khi vẫn duy trì khả năng quan sát và kiểm soát trên toàn đội bay.

---

## Các phương án triển khai

### Nền tảng mã nguồn mở

NATS JetStream là **100% mã nguồn mở** theo giấy phép Apache 2.0:

- **Mã nguồn**: [github.com/nats-io/nats-server](https://github.com/nats-io/nats-server)
- **Không cần giấy phép thương mại** — sử dụng tự do trong môi trường production
- **Đầy đủ tính năng** — phiên bản mã nguồn mở có tất cả các tính năng
- **Phát triển tích cực** — được Synadia hỗ trợ, được hàng nghìn tổ chức sử dụng

Bạn chạy chính xác cùng mã nguồn mà chúng tôi chạy. Không bị ràng buộc nhà cung cấp.

---

### Phương án 1: Tự lưu trữ

Triển khai toàn bộ hệ thống trên hạ tầng của bạn:

| Thành phần | Trách nhiệm của bạn |
|-----------|---------------------|
| Cụm hub | Cấp phát, cấu hình, bảo trì |
| Nút lá | Triển khai trên từng phương tiện |
| Giám sát | Thiết lập khả năng quan sát |
| Cập nhật | Quản lý nâng cấp |

**Chúng tôi cung cấp:**
- Cấu hình tham chiếu (tài liệu này)
- Tư vấn kiến trúc
- Hỗ trợ triển khai

Phù hợp nhất cho: Các tổ chức có đội vận hành, có yêu cầu lưu trú dữ liệu hoặc đã có hạ tầng sẵn.

---

### Phương án 2: Kết nối với hạ tầng của chúng tôi

Tận dụng hạ tầng NATS được quản lý của Ubuntu Software:

| Thành phần | Trách nhiệm |
|-----------|----------------|
| Cụm hub | **Chúng tôi quản lý** |
| Nút lá | Bạn triển khai trên phương tiện |
| Giám sát | **Chúng tôi cung cấp bảng điều khiển** |
| Mở rộng | **Chúng tôi xử lý** |

**Cách hoạt động:**
1. Chúng tôi cấp thông tin xác thực cho đội bay của bạn
2. Các nút lá của bạn kết nối tới các hub khu vực của chúng tôi
3. Dữ liệu của bạn được cô lập (subject chuyên biệt)
4. Bạn nhận được bảng điều khiển đội bay và quyền truy cập API

**Lợi ích:**
- Không cần xây dựng hay bảo trì hạ tầng
- Được tinh chỉnh sẵn cho các mô hình hoạt động của đội drone
- Phủ sóng toàn cầu trên nhiều khu vực
- Tập trung vào drone của bạn, không phải hệ thống nhắn tin

**Giá:**
- **Gói miễn phí** — Tối đa 10 phương tiện, lý tưởng cho phát triển và các triển khai nhỏ
- **Gói mở rộng** — Tính giá theo từng phương tiện khi bạn phát triển

[Liên hệ với chúng tôi →](/contact) để bắt đầu.

---

### So sánh

| Khía cạnh | Tự lưu trữ | Được quản lý |
|--------|-------------|---------|
| **Thời gian thiết lập** | Vài ngày đến vài tuần | Vài giờ |
| **Gánh nặng vận hành** | Bạn bảo trì | Chúng tôi bảo trì |
| **Mô hình chi phí** | Hạ tầng + thời gian | Theo từng phương tiện |
| **Tùy biến** | Toàn quyền kiểm soát | Cấu hình tiêu chuẩn |
| **Vị trí dữ liệu** | Bạn lựa chọn | Các khu vực của chúng tôi |
| **Phù hợp nhất cho** | Đội vận hành lớn | Triển khai nhanh |

---

## Tiếp theo

[Đặt tên Subject →]({{< relref "/fleet/subjects" >}})
