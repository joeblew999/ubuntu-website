---
title: "Nền tảng"
meta_title: "Nền tảng Technology | Ubuntu Software"
description: "Kiến trúc ưu tiên ngoại tuyến, triển khai phổ quát và dữ liệu tự chủ. Nền tảng công nghệ hỗ trợ cả hai nền tảng Publish và Spatial."
image: "/images/foundation.svg"
draft: false
---

## Được thiết kế cho thế giới thực

Mạng Internet bị gián đoạn. Các đội ngũ làm việc trên nhiều châu lục. Máy chủ thuộc sở hữu của bạn. Chúng tôi đã xây dựng hệ thống này để đối phó với thực tế này.

Nền tảng của chúng tôi không tập trung vào các tính năng—mà là cách phần mềm nên hoạt động. Ưu tiên chế độ ngoại tuyến. Tự chủ. Triển khai phổ quát. Những nguyên tắc này được áp dụng trong mọi sản phẩm chúng tôi phát triển.

---

## Kiến trúc ưu tiên chế độ ngoại tuyến

**Làm việc mà không cần kết nối internet. Đồng bộ hóa khi kết nối. Không bao giờ mất dữ liệu.**

Công việc thực sự diễn ra ở những nơi có kết nối WiFi kém—công trường xây dựng, nhà máy, tòa nhà chính phủ, phòng bệnh viện, văn phòng từ xa. Nền tảng của chúng tôi dựa trên giả định về sự mất kết nối, chứ không phải kết nối.

### Cách thức hoạt động

| Component | What It Does |
|-----------|--------------|
| **Local-First Data** | Your data lives on your device first, not in a distant server |
| **Automerge CRDT** | Conflict-free merging when multiple people edit simultaneously |
| **Background Sync** | Automatic synchronization when connectivity returns |
| **Offline Queues** | Actions queue locally, execute when possible |

**Không có trình quay. Không có lỗi "mất kết nối". Chỉ cần làm việc.**

---

## Triển khai ở bất kỳ đâu

**Một mã nguồn duy nhất. Tất cả các nền tảng. Trải nghiệm bản địa.**

Cùng một ứng dụng chạy trên trình duyệt, trên máy tính để bàn và trên các thiết bị di động. Không phải ba sản phẩm riêng biệt—một mã nguồn duy nhất có thể triển khai trên mọi nền tảng.

### Các nền tảng được hỗ trợ

| Platform | Delivery |
|----------|----------|
| **Web** | Any modern browser—Chrome, Firefox, Safari, Edge |
| **Desktop** | Native apps for Windows, macOS, Linux |
| **Mobile** | Native apps for iOS and Android |

Đội ngũ của bạn sử dụng máy tính để bàn tại văn phòng. Nhân viên làm việc ngoài hiện trường sử dụng máy tính bảng. Khách hàng sử dụng điện thoại. Tất cả mọi người đều truy cập cùng một hệ thống với cùng một dữ liệu.

**Viết một lần. Triển khai mọi nơi. Duy trì một kho mã nguồn.**

---

## Tùy chọn đồng bộ hóa đám mây

**Đám mây của bạn. Đám mây của chúng tôi. Không có đám mây. Lựa chọn của bạn.**

Đồng bộ hóa không yêu cầu sử dụng máy chủ của chúng tôi. Kết nối với bất kỳ hạ tầng nào phù hợp với tổ chức của bạn.

### Các mô hình triển khai

| Model | Best For |
|-------|----------|
| **Ubuntu Software Cloud** | Fastest setup, we handle operations |
| **Your Cloud** | AWS, Azure, GCP—your infrastructure, our software |
| **On-Premises** | Your data center, complete control |
| **Hybrid** | Some data in cloud, sensitive data on-prem |
| **Air-Gapped** | Fully disconnected networks, defense and secure environments |

Chuyển đổi giữa các mô hình? Rất đơn giản. Định dạng dữ liệu của bạn không thay đổi tùy thuộc vào nơi nó được lưu trữ.

**Không bị ràng buộc bởi nhà cung cấp. Không bị ép buộc phải sử dụng đám mây. Tính linh hoạt thực sự trong triển khai.**

---

## Tự chủ

**Dữ liệu của bạn. Máy chủ của bạn. Quy tắc của bạn.**

Tự chủ có nghĩa là bạn kiểm soát cơ sở hạ tầng của mình. Không phải "dữ liệu của bạn được lưu trữ theo điều kiện của chúng tôi"—mà thực sự là của bạn.

### Tự chủ là gì?

- **Chạy ở bất kỳ đâu** — Trung tâm dữ liệu của bạn, tài khoản đám mây của bạn, laptop của bạn
- **Không cần kết nối với máy chủ của chúng tôi* — Phần mềm hoạt động mà không cần kết nối với máy chủ của chúng tôi
- **Xuất tất cả dữ liệu** — Các định dạng tiêu chuẩn, khả năng di chuyển dữ liệu đầy đủ
- **Không theo dõi hoạt động sử dụng** — Chúng tôi không theo dõi cách bạn sử dụng phần mềm của chúng tôi
- **Tùy chọn giấy phép vĩnh viễn** — Tiếp tục hoạt động ngay cả khi chúng tôi không còn tồn tại

**Đây không chỉ là vấn đề riêng tư—đó là sự độc lập hoạt động.**

---

## Có thể nhúng

**Tích hợp vào các hệ thống hiện có. Không thay thế hoàn toàn.**

Các tổ chức đã có sẵn các quy trình làm việc, cơ sở dữ liệu và hệ thống xác thực. Nền tảng của chúng tôi tích hợp thay vì thay thế.

### Các mẫu tích hợp

| Pattern | Use Case |
|---------|----------|
| **API-First** | Everything accessible programmatically |
| **Database Connectors** | Read/write to your existing databases |
| **SSO Integration** | Your identity provider, not another login |
| **Webhook Events** | Push notifications to your systems |
| **White-Label** | Your branding, our engine |

**Mở rộng hệ thống của bạn. Đừng bỏ rơi chúng.**

---

## Cổng thông tin nổi tiếng

**Đăng tải lên các nền tảng công nghệ lớn. Kiểm soát mối quan hệ.**

Các cổng thông tin của mạng internet—Google, Apple, YouTube—có phạm vi tiếp cận khổng lồ. Tuy nhiên, việc đăng tải nội dung lên các nền tảng này không đồng nghĩa với việc bị sở hữu bởi họ.

### Đảo ngược mối quan hệ

Cách tiếp cận truyền thống:
```
User → YouTube (owns everything) → Your content (captive)
```

Phương pháp phổ biến:
```
User → Your Gateway → Your System (primary)
                   ↳→ YouTube (mirror for discovery)
```

**Bạn kiểm soát cửa chính.** Big Tech trở thành kênh phân phối tùy chọn, không phải là một nhà tù.

### Cách thức hoạt động

| Capability | What It Means |
|------------|---------------|
| **Your URIs everywhere** | Links point to YOUR gateway, not theirs |
| **Smart routing** | Send iOS users to Apple, Android to Google, web to your player |
| **Mirror publishing** | Auto-publish copies to YouTube, Google Maps, Apple Calendar |
| **Analytics you own** | See everything, track everyone, no data sharing |
| **Exit strategy built-in** | Remove any platform from routing without breaking links |

### Tương thích với mọi thứ

- **Video** — Tải lên máy chủ của bạn, đồng bộ hóa lên YouTube/Twitch để tăng phạm vi tiếp cận
- **Lịch** — Máy chủ CalDAV của bạn, đồng bộ hóa với Google/Apple để tiện lợi
- **Bản đồ** — Dữ liệu địa lý của bạn, tích hợp với Google/Apple Maps
- **Email** — Máy chủ email của bạn, tương thích với Gmail
- **Tệp tin** — Lưu trữ của bạn, chia sẻ chọn lọc lên Drive/Dropbox

**Đăng tải lên các nền tảng của họ. Đừng bao giờ bị khóa lại.**

---

## 25 năm kinh nghiệm trong lĩnh vực doanh nghiệp

**Ubuntu Software được thành lập vào năm 1999.** Chúng tôi đã xây dựng các hệ thống doanh nghiệp qua các giai đoạn dot-com, cách mạng di động, chuyển đổi đám mây và sự trỗi dậy của trí tuệ nhân tạo.

### Những bài học mà kinh nghiệm đã dạy cho chúng ta

- **Nhà cung cấp biến mất** — Xây dựng khả năng di chuyển dữ liệu từ ngày đầu tiên
- **Mạng lưới gặp sự cố** — Chế độ ưu tiên ngoại tuyến không phải là tùy chọn, mà là điều cần thiết
- **Yêu cầu thay đổi** — Tiêu chuẩn mở tồn tại lâu hơn các định dạng độc quyền
- **Quy mô mang lại bất ngờ** — Kiến trúc quan trọng hơn tối ưu hóa
- **Tích hợp là một thách thức** — Thiết kế cho nó, đừng chỉ thêm vào một cách gượng ép

Chúng tôi đã thấy những gì hiệu quả và những gì không. Nền tảng này phản ánh những bài học được rút ra từ hàng thập kỷ kinh nghiệm trong môi trường sản xuất.

---

## Quỹ tài trợ là nguồn động lực cho mọi hoạt động

Cả hai nền tảng Publish và Spatial đều tự động kế thừa các tính năng này:

| Capability | Publish | Spatial |
|------------|---------|---------|
| Offline editing | Edit documents without internet | Design 3D models without internet |
| Real-time sync | Multiple editors, one document | Multiple designers, one model |
| Universal deploy | Forms on any device | 3D viewer on any device |
| Self-hosted | Your document server | Your design server |
| Cloud options | Managed or self-managed | Managed or self-managed |

**Chọn nền tảng của bạn. Nhận nền tảng miễn phí.**

---

## Chi tiết kỹ thuật

Đối với các đội đang đánh giá kiến trúc của chúng tôi:

| Layer | Technology |
|-------|------------|
| **Storage** | SQLite everywhere—local devices and servers |
| **Sync Engine** | CRDT-based replication via NATS JetStream |
| **Messaging** | NATS JetStream |
| **UI Framework** | Cross-platform native rendering |
| **API** | HTTP REST + SSE (Server-Sent Events) |
| **AI Integration** | Model Context Protocol (MCP) |
| **Auth** | OIDC-compatible, bring your own IdP |

### SQLite phân tán

Mỗi nút — máy tính xách tay của bạn, điện thoại của bạn, máy chủ của bạn — đều chạy SQLite. Các thay đổi được đồng bộ hóa qua NATS JetStream sử dụng ngữ nghĩa CRDT.

- **Bất kỳ máy chủ nào cũng có thể gặp sự cố** — Các máy chủ khác vẫn tiếp tục hoạt động
- **Bất kỳ máy chủ nào cũng có thể ngừng hoạt động** — Đồng bộ hóa khi kết nối lại
- **Không có điểm yếu duy nhất* — Kiến trúc phân tán thực sự
- **Cùng một cơ sở dữ liệu ở mọi nơi** — Từ thiết bị cục bộ đến cụm toàn cầu

**Công nghệ tiêu chuẩn. Không có sự ràng buộc độc quyền.**

---

## Mở rộng không giới hạn

### Không có điểm yếu duy nhất (SPOF)

Mỗi thành phần đều có tính dự phòng. Không có máy chủ, dịch vụ hoặc trung tâm dữ liệu nào có thể làm hệ thống ngừng hoạt động. Nếu bất kỳ nút nào gặp sự cố, hệ thống vẫn tiếp tục hoạt động.

### Không có điểm yếu duy nhất (SPOP)

Hệ thống tính toán mở rộng theo chiều ngang. Tăng khả năng xử lý bằng cách thêm các nút (nodes), không cần mua các máy chủ có cấu hình cao hơn. Các tác vụ được phân phối tự động trên các tài nguyên có sẵn.

### Hàng trăm trung tâm dữ liệu

Kiến trúc được thiết kế cho phân phối toàn cầu:

| Capability | What It Means |
|------------|---------------|
| **Deploy anywhere** | Cloud, on-prem, edge, air-gapped |
| **Deploy close to users** | Low latency, local compliance |
| **Replicate for redundancy** | Survive regional outages |
| **Partition tolerance** | Operate independently when networks split |

### Các tùy chọn triển khai

| Method | Use Case |
|--------|----------|
| **Binaries** | Single-file deployment, minimal dependencies |
| **Docker** | Containerized, reproducible environments |
| **Kubernetes** | Orchestrated, auto-scaling clusters |

**Từ một chiếc laptop duy nhất đến hàng trăm trung tâm dữ liệu. Cùng một kiến trúc. Cùng một mã nguồn.**

---

## Bắt đầu

Nền tảng đã được tích hợp sẵn. Khi bạn sử dụng Publish hoặc Spatial, bạn sẽ tự động nhận được các tùy chọn triển khai ưu tiên ngoại tuyến, triển khai đa nền tảng và tự chủ dữ liệu.

[Explore Publish →](/vi/platform/publish/) | [Explore Spatial →](/vi/platform/spatial/) | [Linux & Cross-Platform →](/vi/technology/linux/) | [Contact Us →](/vi/contact/)
