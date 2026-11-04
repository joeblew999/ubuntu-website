---
title: "Nền tảng"
meta_title: "Công nghệ nền tảng | Ubuntu Software"
description: "Kiến trúc ưu tiên ngoại tuyến, triển khai phổ quát và dữ liệu tự chủ. Nền tảng công nghệ vận hành cả hai nền tảng Publish và Spatial."
image: "/images/foundation.svg"
draft: false
---

## Được xây dựng cho thế giới thực

Internet bị gián đoạn. Các nhóm làm việc trải rộng khắp các châu lục. Máy chủ thuộc về bạn. Chúng tôi xây dựng cho chính thực tế này.

Nền tảng cốt lõi của chúng tôi không xoay quanh tính năng—mà là về cách phần mềm nên vận hành. Ưu tiên ngoại tuyến. Tự chủ. Triển khai đa nền tảng. Những nguyên tắc này xuyên suốt mọi thứ chúng tôi xây dựng.

---

## Kiến trúc ưu tiên ngoại tuyến

**Làm việc không cần internet. Đồng bộ khi có kết nối. Không bao giờ mất dữ liệu.**

Công việc thực tế diễn ra ở những nơi WiFi kém—công trường xây dựng, xưởng sản xuất, tòa nhà chính phủ, phòng bệnh, văn phòng ở vùng xa. Nền tảng cốt lõi của chúng tôi mặc định tình trạng mất kết nối, chứ không phải luôn có kết nối.

### Cách hoạt động

| Thành phần | Chức năng |
|-----------|--------------|
| **Dữ liệu cục bộ là trên hết** | Dữ liệu của bạn nằm trên thiết bị của bạn trước tiên, không phải trên một máy chủ xa xôi |
| **Automerge CRDT** | Hợp nhất không xung đột khi nhiều người chỉnh sửa cùng lúc |
| **Đồng bộ nền** | Tự động đồng bộ khi kết nối được khôi phục |
| **Hàng đợi ngoại tuyến** | Các thao tác được xếp hàng cục bộ, thực thi khi có thể |

**Không biểu tượng chờ xoay vòng. Không lỗi "mất kết nối". Chỉ có làm việc.**

---

## Triển khai ở bất cứ đâu

**Một codebase. Mọi nền tảng. Trải nghiệm gốc.**

Cùng một ứng dụng chạy trên trình duyệt, máy tính để bàn và thiết bị di động. Không phải ba sản phẩm riêng biệt—một codebase duy nhất triển khai ở mọi nơi.

### Nền tảng được hỗ trợ

| Nền tảng | Hình thức phân phối |
|----------|----------|
| **Web** | Mọi trình duyệt hiện đại—Chrome, Firefox, Safari, Edge |
| **Máy tính để bàn** | Ứng dụng gốc cho Windows, macOS, Linux |
| **Di động** | Ứng dụng gốc cho iOS và Android |

Nhóm của bạn dùng máy tính để bàn tại văn phòng. Nhân viên hiện trường dùng máy tính bảng. Khách hàng dùng điện thoại. Mọi người truy cập cùng một hệ thống với cùng một dữ liệu.

**Viết một lần. Triển khai mọi nơi. Duy trì một codebase.**

---

## Tùy chọn đồng bộ đám mây

**Đám mây của bạn. Đám mây của chúng tôi. Không đám mây. Bạn chọn.**

Đồng bộ không yêu cầu máy chủ của chúng tôi. Kết nối với bất kỳ hạ tầng nào phù hợp với tổ chức của bạn.

### Mô hình triển khai

| Mô hình | Phù hợp nhất cho |
|-------|----------|
| **Ubuntu Software Cloud** | Thiết lập nhanh nhất, chúng tôi lo phần vận hành |
| **Đám mây của bạn** | AWS, Azure, GCP—hạ tầng của bạn, phần mềm của chúng tôi |
| **Tại chỗ (On-Premises)** | Trung tâm dữ liệu của bạn, toàn quyền kiểm soát |
| **Kết hợp (Hybrid)** | Một phần dữ liệu trên đám mây, dữ liệu nhạy cảm đặt tại chỗ |
| **Cách ly mạng (Air-Gapped)** | Mạng hoàn toàn tách biệt, môi trường quốc phòng và bảo mật |

Chuyển đổi giữa các mô hình? Đơn giản. Định dạng dữ liệu của bạn không thay đổi dù nó được lưu ở đâu.

**Không bị khóa chặt vào nhà cung cấp. Không bắt buộc phụ thuộc đám mây. Linh hoạt triển khai thực sự.**

---

## Tự chủ

**Dữ liệu của bạn. Máy chủ của bạn. Quy tắc của bạn.**

Tự chủ nghĩa là bạn kiểm soát hạ tầng của mình. Không phải "dữ liệu của bạn được lưu trữ theo điều khoản của chúng tôi"—mà thực sự là của bạn.

### Tự chủ nghĩa là gì

- **Chạy ở bất cứ đâu** — Trung tâm dữ liệu của bạn, tài khoản đám mây của bạn, máy tính xách tay của bạn
- **Không "gọi về nhà"** — Phần mềm hoạt động mà không cần liên lạc với máy chủ của chúng tôi
- **Xuất mọi thứ** — Định dạng tiêu chuẩn, khả năng di chuyển dữ liệu hoàn toàn
- **Không theo dõi sử dụng** — Chúng tôi không thấy bạn làm gì với phần mềm của chúng tôi
- **Tùy chọn giấy phép vĩnh viễn** — Tiếp tục vận hành ngay cả khi chúng tôi biến mất

**Đây không chỉ là quyền riêng tư—mà là sự độc lập trong vận hành.**

---

## Có thể nhúng

**Tích hợp vào hệ thống hiện có. Không cần gỡ bỏ và thay thế.**

Các tổ chức đều có sẵn quy trình làm việc, cơ sở dữ liệu, hệ thống định danh. Nền tảng cốt lõi của chúng tôi tích hợp thay vì thay thế.

### Mô hình tích hợp

| Mô hình | Trường hợp sử dụng |
|---------|----------|
| **Ưu tiên API** | Mọi thứ đều truy cập được bằng lập trình |
| **Trình kết nối cơ sở dữ liệu** | Đọc/ghi vào cơ sở dữ liệu hiện có của bạn |
| **Tích hợp SSO** | Nhà cung cấp định danh của bạn, không thêm một lần đăng nhập nữa |
| **Sự kiện Webhook** | Đẩy thông báo đến hệ thống của bạn |
| **Nhãn trắng (White-Label)** | Thương hiệu của bạn, động cơ của chúng tôi |

**Mở rộng hệ thống của bạn. Đừng từ bỏ chúng.**

---

## Wellknown Gateway

**Xuất bản lên các nền tảng Big Tech. Làm chủ mối quan hệ.**

Các cổng của web—Google, Apple, YouTube—có phạm vi tiếp cận khổng lồ. Nhưng xuất bản LÊN chúng không có nghĩa là bị chúng SỞ HỮU.

### Đảo ngược mối quan hệ

Cách tiếp cận truyền thống:
```
User → YouTube (owns everything) → Your content (captive)
```

Cách tiếp cận Wellknown:
```
User → Your Gateway → Your System (primary)
                   ↳→ YouTube (mirror for discovery)
```

**Bạn kiểm soát cửa chính.** Big Tech trở thành kênh phân phối tùy chọn, không phải nhà tù.

### Cách hoạt động

| Khả năng | Ý nghĩa |
|------------|---------------|
| **URI của bạn ở khắp nơi** | Liên kết trỏ đến cổng CỦA BẠN, không phải của họ |
| **Định tuyến thông minh** | Đưa người dùng iOS đến Apple, Android đến Google, web đến trình phát của bạn |
| **Xuất bản bản sao** | Tự động xuất bản bản sao lên YouTube, Google Maps, Apple Calendar |
| **Phân tích do bạn sở hữu** | Thấy mọi thứ, theo dõi tất cả mọi người, không chia sẻ dữ liệu |
| **Tích hợp sẵn chiến lược rút lui** | Gỡ bất kỳ nền tảng nào khỏi định tuyến mà không làm hỏng liên kết |

### Hoạt động với mọi thứ

- **Video** — Lưu trữ trên máy chủ của bạn, sao chép sang YouTube/Twitch để mở rộng tiếp cận
- **Lịch** — Máy chủ CalDAV của bạn, đồng bộ với Google/Apple cho tiện lợi
- **Bản đồ** — Dữ liệu địa lý của bạn, tích hợp với Google/Apple Maps
- **Email** — Máy chủ thư của bạn, tương thích với Gmail
- **Tệp** — Kho lưu trữ của bạn, chia sẻ có chọn lọc lên Drive/Dropbox

**Xuất bản LÊN nền tảng của họ. Không bao giờ bị khóa chặt TRONG đó.**

---

## 25 năm kinh nghiệm doanh nghiệp

**Ubuntu Software bắt đầu từ năm 1999.** Chúng tôi đã xây dựng các hệ thống doanh nghiệp qua thời kỳ dot-com, cuộc cách mạng di động, quá trình chuyển đổi lên đám mây và sự trỗi dậy của AI.

### Kinh nghiệm đã dạy chúng tôi điều gì

- **Nhà cung cấp biến mất** — Xây dựng cho khả năng di chuyển dữ liệu ngay từ ngày đầu
- **Mạng bị lỗi** — Ưu tiên ngoại tuyến không phải tùy chọn, mà là thiết yếu
- **Yêu cầu thay đổi** — Tiêu chuẩn mở tồn tại lâu hơn định dạng độc quyền
- **Quy mô gây bất ngờ** — Kiến trúc quan trọng hơn tối ưu hóa
- **Tích hợp rất khó** — Thiết kế cho nó ngay từ đầu, đừng chắp vá về sau

Chúng tôi đã thấy điều gì hiệu quả và điều gì không. Nền tảng cốt lõi này phản ánh hàng thập kỷ bài học rút ra từ môi trường sản xuất thực tế.

---

## Nền tảng cốt lõi vận hành mọi thứ

Cả hai nền tảng Publish và Spatial đều tự động kế thừa các khả năng này:

| Khả năng | Publish | Spatial |
|------------|---------|---------|
| Chỉnh sửa ngoại tuyến | Chỉnh sửa tài liệu không cần internet | Thiết kế mô hình 3D không cần internet |
| Đồng bộ thời gian thực | Nhiều người chỉnh sửa, một tài liệu | Nhiều nhà thiết kế, một mô hình |
| Triển khai đa nền tảng | Biểu mẫu trên mọi thiết bị | Trình xem 3D trên mọi thiết bị |
| Tự lưu trữ | Máy chủ tài liệu của bạn | Máy chủ thiết kế của bạn |
| Tùy chọn đám mây | Được quản lý hoặc tự quản lý | Được quản lý hoặc tự quản lý |

**Chọn nền tảng của bạn. Nhận nền tảng cốt lõi miễn phí.**

---

## Chi tiết kỹ thuật

Dành cho các nhóm đang đánh giá kiến trúc của chúng tôi:

| Tầng | Công nghệ |
|-------|------------|
| **Lưu trữ** | SQLite ở mọi nơi—thiết bị cục bộ và máy chủ |
| **Công cụ đồng bộ** | Sao chép dựa trên CRDT qua NATS JetStream |
| **Nhắn tin** | NATS JetStream |
| **Khung giao diện (UI Framework)** | Kết xuất gốc đa nền tảng |
| **API** | HTTP REST + SSE (Server-Sent Events) |
| **Tích hợp AI** | Model Context Protocol (MCP) |
| **Xác thực** | Tương thích OIDC, dùng IdP của riêng bạn |

### NATS JetStream

Xương sống nhắn tin của chúng tôi là một dự án đang trong giai đoạn ươm tạo của **CNCF (Cloud Native Computing Foundation)**—cũng chính là tổ chức quản lý Kubernetes.

| Thuộc tính | Chi tiết |
|----------|--------|
| **Quản trị** | Dự án ươm tạo của CNCF |
| **Giấy phép** | Apache 2.0 (100% mã nguồn mở) |
| **Chứng nhận** | Synadia duy trì chứng nhận SOC 2 Type II |
| **Cộng đồng** | Hàng nghìn hệ thống triển khai thực tế trên toàn cầu |

Mã nguồn mở. Được gia cố cho doanh nghiệp. Đã được kiểm toán.

[Chi tiết bảo mật & tuân thủ →]({{< relref "/technology/security" >}})

### SQLite phân tán

Mọi nút—máy tính xách tay, điện thoại, máy chủ của bạn—đều chạy SQLite. Các thay đổi được sao chép qua NATS JetStream theo ngữ nghĩa CRDT.

- **Bất kỳ máy chủ nào cũng có thể gặp sự cố** — Các máy chủ khác vẫn tiếp tục phục vụ
- **Bất kỳ máy chủ nào cũng có thể ngoại tuyến** — Đồng bộ khi kết nối lại
- **Không có điểm lỗi đơn lẻ** — Kiến trúc phân tán thực sự
- **Cùng một cơ sở dữ liệu ở mọi nơi** — Từ thiết bị cục bộ đến cụm toàn cầu

**Công nghệ tiêu chuẩn. Không bị khóa chặt bởi công nghệ độc quyền.**

---

## Mở rộng không giới hạn

### Không có điểm lỗi đơn lẻ (SPOF)

Mọi thành phần đều có dự phòng. Không một máy chủ, dịch vụ hay trung tâm dữ liệu đơn lẻ nào có thể làm sập hệ thống. Bất kỳ nút nào gặp sự cố—hệ thống vẫn tiếp tục chạy.

### Không có điểm nghẽn hiệu năng đơn lẻ (SPOP)

Năng lực tính toán mở rộng theo chiều ngang. Tăng công suất bằng cách thêm nút, không phải mua máy chủ lớn hơn. Khối lượng công việc tự động phân bổ trên các tài nguyên sẵn có.

### Hàng trăm trung tâm dữ liệu

Kiến trúc được thiết kế cho phân phối toàn cầu:

| Khả năng | Ý nghĩa |
|------------|---------------|
| **Triển khai ở bất cứ đâu** | Đám mây, tại chỗ, biên (edge), cách ly mạng |
| **Triển khai gần người dùng** | Độ trễ thấp, tuân thủ quy định địa phương |
| **Sao chép để dự phòng** | Vượt qua sự cố mất dịch vụ theo khu vực |
| **Chịu được phân mảnh mạng** | Vận hành độc lập khi mạng bị chia tách |

### Tùy chọn triển khai

| Phương thức | Trường hợp sử dụng |
|--------|----------|
| **Tệp nhị phân** | Triển khai bằng một tệp duy nhất, phụ thuộc tối thiểu |
| **Docker** | Môi trường đóng gói trong container, có thể tái tạo |
| **Kubernetes** | Cụm được điều phối, tự động mở rộng |

**Từ một máy tính xách tay đến hàng trăm trung tâm dữ liệu. Cùng một kiến trúc. Cùng một codebase.**

---

## Bắt đầu

Nền tảng cốt lõi được tích hợp sẵn. Khi sử dụng Publish hoặc Spatial, bạn tự động có được khả năng ưu tiên ngoại tuyến, triển khai đa nền tảng và các tùy chọn tự chủ.

[Khám phá Publish →](/platform/publish/) | [Khám phá Spatial →](/platform/spatial/) | [Linux & Đa nền tảng →](/technology/linux/) | [Liên hệ với chúng tôi →](/contact/)
