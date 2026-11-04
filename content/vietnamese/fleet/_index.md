---
title: "Kiến trúc đội xe tự hành"
meta_title: "Kiến trúc đội phương tiện cho Drone, Ô tô & Xe tải | Ubuntu Software"
description: "Kiến trúc tham chiếu đạt chuẩn sản xuất để quản lý đội phương tiện tự hành—drone, ô tô, xe tải, AGV—sử dụng NATS JetStream để tạo bản sao số (digital twin) ở quy mô toàn đội."
image: "/images/robotics.svg"
draft: false
---

## Quản lý đội phương tiện đạt chuẩn sản xuất ở quy mô lớn

Kiến trúc tham chiếu này minh họa cách triển khai và quản lý **đội phương tiện tự hành**—drone, ô tô, xe tải, AGV—bằng các thành phần mã nguồn mở và các mẫu thiết kế đã được kiểm chứng. Mọi quyết định thiết kế đều nhắm đến các yêu cầu vận hành thực tế: độ tin cậy, khả năng quan sát (observability) và suy giảm có kiểm soát (graceful degradation).

Kiến trúc lõi **không phụ thuộc vào loại phương tiện**. NATS JetStream, các mẫu digital twin và mô hình phân quyền hoạt động giống hệt nhau dù bạn đang quản lý drone hay xe tải giao hàng. Chỉ có phần cứng biên và giao thức là thay đổi.

---

## Thách thức

Vận hành ở quy mô đội phương tiện đòi hỏi giải quyết những vấn đề không tồn tại ở quy mô nhỏ:

- **Hàng nghìn luồng telemetry đồng thời** — Mỗi phương tiện báo cáo vị trí, vận tốc, tình trạng và dữ liệu cảm biến
- **Chỉ huy và điều khiển ở quy mô lớn** — Gửi lệnh đến từng phương tiện hoặc nhóm phương tiện cụ thể mà không làm tràn ngập mạng
- **Đồng bộ digital twin** — Duy trì biểu diễn trạng thái chính xác cho mọi phương tiện theo thời gian thực
- **Khả năng chống chịu khi ngoại tuyến** — Phương tiện tiếp tục hoạt động khi mất kết nối
- **Đảm bảo an toàn** — Bảo đảm sự cố mạng không bao giờ ảnh hưởng đến an toàn của phương tiện

Các cách tiếp cận truyền thống—kết nối giao thức trực tiếp, cơ sở dữ liệu tập trung, kiến trúc polling—đều sụp đổ trước những yêu cầu này.

---

## Cách tiếp cận của chúng tôi

Kiến trúc kết hợp ba công nghệ đã được kiểm chứng:

| Tầng | Công nghệ | Vai trò |
|-------|------------|------|
| **Điều khiển phương tiện** | Theo từng nền tảng (PX4, ECU, ROS 2) | Vận hành tự hành, failsafe, điều khiển thủ công |
| **Điện toán biên** | Jetson + PC công nghiệp | Xử lý cảm biến, suy luận AI, ra quyết định cục bộ |
| **Nhắn tin đội phương tiện** | NATS JetStream | Pub/sub, lưu trữ bền vững, trạng thái digital twin |

Mỗi tầng hoạt động độc lập. Sự cố mạng chỉ gây suy giảm có kiểm soát—phương tiện tiếp tục hoạt động, máy tính biên tiếp tục xử lý, và trạng thái được đồng bộ khi kết nối được khôi phục.

---

## Các thành phần kiến trúc

### Nền tảng được hỗ trợ

Kiến trúc hỗ trợ nhiều loại phương tiện với phần cứng và giao thức riêng cho từng nền tảng:

| Nền tảng | Giao thức | Hệ thống điều khiển | Trường hợp sử dụng |
|----------|----------|----------------|-----------|
| **Drone** | MAVLink | PX4/ArduPilot | Kiểm tra, lập bản đồ, giao hàng |
| **Phương tiện mặt đất** | CAN bus / J1939 | ECU / Autoware | Logistics, khai khoáng, nông nghiệp |
| **Robot ROS 2** | ROS 2 topics | Tùy chỉnh | Kho bãi, nghiên cứu, AGV |

[Nền tảng được hỗ trợ →]({{< relref "/fleet/platforms" >}})

### Bộ phần cứng (Drone)

Các lựa chọn phần cứng tiêu chuẩn, cân bằng giữa năng lực, tính sẵn có và khả năng bảo trì:

- **Khung máy bay**: Holybro X500 V2 ARF — nền tảng đã được kiểm chứng, linh kiện thay thế rất sẵn có
- **Bộ điều khiển bay**: Pixhawk 6X chạy PX4 v1.14 LTS
- **Máy tính đồng hành cảm biến**: Raspberry Pi 4 / CM4 để tích hợp cảm biến gọn nhẹ
- **Máy tính đồng hành AI**: NVIDIA Jetson (Orin/Xavier/Nano) cho thị giác máy tính và suy luận

[Chi tiết phần cứng drone →]({{< relref "/fleet/platforms/drones" >}}) | [Chi tiết phương tiện mặt đất →]({{< relref "/fleet/platforms/ground" >}})

### Bộ phần mềm

Phần mềm mã nguồn mở ở mọi tầng:

- **PX4 v1.14 LTS** — Điều khiển bay với độ ổn định đã được chứng minh
- **Ubuntu Server 22.04** — Môi trường Linux nhất quán trên tất cả máy tính đồng hành
- **ROS 2 Humble** — Middleware robot để tích hợp cảm biến
- **Go** — Ngôn ngữ triển khai Vehicle Gateway

[Chi tiết phần mềm →]({{< relref "/fleet/software" >}})

### Kiến trúc NATS

Topology nhắn tin phân cấp được thiết kế cho triển khai qua WAN:

- **Leaf node** trên mỗi phương tiện — pub/sub cục bộ, lưu và chuyển tiếp (store-and-forward)
- **Cụm hub khu vực** — tập hợp phương tiện theo khu vực địa lý
- **Mirroring toàn cầu** — sao chép liên khu vực khi cần

[Topology NATS →]({{< relref "/fleet/nats-topology" >}})

### Thiết kế Digital Twin

Các stream JetStream và KV store duy trì trạng thái của đội phương tiện:

- **Phân cấp subject** để định tuyến và lọc
- **State stream** để tổng hợp telemetry
- **Event stream** cho nhật ký kiểm toán
- **Shadow store** để đối chiếu trạng thái mong muốn/được báo cáo

[Đặt tên Subject →]({{< relref "/fleet/subjects" >}}) | [Cấu hình Stream →]({{< relref "/fleet/streams" >}})

### Vehicle Gateway

Dịch vụ chạy trên máy tính biên của mỗi phương tiện, làm cầu nối giữa các giao thức gốc và NATS:

- Chuyển đổi giao thức (MAVLink, CAN bus, ROS 2)
- Giảm tần suất lấy mẫu và tổng hợp trạng thái
- Trích xuất sự kiện từ telemetry
- Thực thi lệnh kèm áp dụng chính sách
- Đối chiếu trạng thái shadow

[Thiết kế Gateway →]({{< relref "/fleet/gateway" >}})

### Mô hình an toàn

Kết nối mạng không bao giờ được tin cậy cho an toàn của phương tiện:

- **Điều khiển thủ công là quyền ưu tiên cao nhất** — Người vận hành luôn nắm quyền kiểm soát (RC cho drone, vô lăng/dừng khẩn cấp cho phương tiện mặt đất)
- **Hệ thống điều khiển phương tiện thực thi failsafe** — RTL/dừng tại chỗ, geofence
- **NATS không bao giờ nằm trong vòng điều khiển** — Chỉ dùng để giám sát và điều phối
- **Suy giảm có kiểm soát** — Mất AI hoặc mất mạng sẽ kích hoạt chế độ an toàn

[Chi tiết an toàn →]({{< relref "/fleet/safety" >}})

---

## Tại sao chọn kiến trúc này

### Đã được kiểm chứng ở quy mô lớn

Mọi thành phần đều đã được triển khai trong môi trường sản xuất:

- **NATS** vận hành hạ tầng nhắn tin toàn cầu của Synadia
- **PX4** bay trên hàng nghìn drone thương mại trên toàn thế giới
- **Jetson** chạy suy luận trong các phương tiện tự hành và robot

### Tiêu chuẩn mở

Không bị khóa chặt vào nhà cung cấp:

- MAVLink là giao thức mở với nhiều bản triển khai
- NATS là mã nguồn mở và có sẵn hỗ trợ thương mại
- PX4 chạy trên nhiều nền tảng phần cứng điều khiển bay

### Hạ tầng của bạn hoặc của chúng tôi

NATS JetStream là **100% mã nguồn mở** (Apache 2.0). Tự vận hành hoặc kết nối với hạ tầng được quản lý của chúng tôi—cùng giao thức, cùng mã nguồn, bạn tự lựa chọn.

- **Tự lưu trữ** — Triển khai trên hạ tầng của bạn với các cấu hình tham chiếu của chúng tôi
- **Được quản lý** — Kết nối leaf node của bạn với các hub khu vực của chúng tôi

Miễn phí cho đội phương tiện nhỏ. Mở rộng cùng bạn.

[Tùy chọn triển khai →]({{< relref "/fleet/nats-topology#deployment-options" >}})

### Thực tế vận hành

Được thiết kế cho thế giới thực:

- Phương tiện có thể được bảo dưỡng mà không cần công cụ chuyên dụng
- Cập nhật phần mềm được triển khai qua mạng (OTA)
- Dữ liệu telemetry được đưa vào các bộ công cụ observability tiêu chuẩn
- Trạng thái đội phương tiện được xuất sang các hệ thống doanh nghiệp hiện có

---

## Tìm hiểu thêm

Khám phá chi tiết từng thành phần:

| Mục | Mô tả |
|---------|-------------|
| [Nền tảng được hỗ trợ]({{< relref "/fleet/platforms" >}}) | Drone, phương tiện mặt đất, AGV |
| [Nền tảng Drone]({{< relref "/fleet/platforms/drones" >}}) | Khung máy bay, bộ điều khiển bay, MAVLink |
| [Phương tiện mặt đất]({{< relref "/fleet/platforms/ground" >}}) | Ô tô, xe tải, CAN bus, J1939 |
| [Phần mềm]({{< relref "/fleet/software" >}}) | Hệ điều hành, middleware, ứng dụng |
| [Topology NATS]({{< relref "/fleet/nats-topology" >}}) | Leaf node, hub, kết nối WAN |
| [Đặt tên Subject]({{< relref "/fleet/subjects" >}}) | Cấu trúc subject phân cấp |
| [Cấu hình Stream]({{< relref "/fleet/streams" >}}) | Thiết lập JetStream cho digital twin |
| [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) | Cầu nối từ giao thức sang NATS |
| [Phân quyền & Cấp quyền]({{< relref "/fleet/authorization" >}}) | Bảo mật phi tập trung, truy cập của bên thứ ba |
| [Mô hình an toàn]({{< relref "/fleet/safety" >}}) | Failsafe và suy giảm có kiểm soát |
| [Quy định BVLOS]({{< relref "/fleet/bvlos-regulations" >}}) | Quy định toàn cầu về bay ngoài tầm nhìn trực tiếp (Beyond Visual Line of Sight) |
| [Trình xem không phận]({{< relref "/fleet/airspace-demo" >}}) | Trực quan hóa tương tác không phận Hoa Kỳ |

---

## Bắt đầu

Bạn đang xây dựng đội phương tiện tự hành? Triển khai drone, xe tải hoặc AGV ở quy mô lớn? Chúng tôi có thể hỗ trợ về kiến trúc, triển khai và vận hành.

[Liên hệ với chúng tôi →](/contact)
