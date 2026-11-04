---
title: "Xuất bản"
meta_title: "Nền tảng quản lý nội dung | Phần mềm Ubuntu"
description: "CMS nguồn đơn cho xuất bản đa kênh. Một DSL duy nhất cho văn bản, đồ họa, biểu mẫu, email. Đầu ra được dịch tự động. Thương hiệu hoàn hảo nhờ kiến trúc."
image: "/images/publish.svg"
draft: false
---

## Một ngôn ngữ. Mọi đầu ra. Mọi ngôn ngữ.

Một hệ thống quản lý nội dung từ một nguồn duy nhất. Định nghĩa mọi thứ trong một DSL—**văn bản (1D) điều khiển mọi đầu ra 2D và 3D**. Kết hợp các thành phần chứa các thành phần khác. Xuất ra mọi kênh. Tự động dịch sang mọi ngôn ngữ. Thu thập dữ liệu trở lại.

Tất cả đều chạy trên web. Tất cả đều truy cập được toàn cầu. Thương hiệu hoàn hảo—không nhờ kỷ luật, mà nhờ kiến trúc.

---

## Vấn đề

Các tổ chức xuất bản nội dung trên nhiều kênh: website, PDF, email, biểu mẫu, bảng hiệu.

Mỗi kênh được duy trì riêng biệt. Nội dung lệch nhau. Thương hiệu rạn nứt. Bản dịch chậm trễ. Dữ liệu biểu mẫu gửi về phải nhập lại thủ công.

**Một thay đổi nghĩa là phải cập nhật năm hệ thống. Hoặc chẳng được cập nhật gì cả.**

---

## Một DSL

Mọi thứ được định nghĩa trong cùng một ngôn ngữ.

**Nội dung văn bản:**
```
# Welcome to Our Service

We help organizations {industry} achieve {outcome}.
```

**Đồ họa:**
```
@logo: svg {
  viewBox: "0 0 200 50"
  rect: { x: 0, y: 0, width: 50, height: 50, fill: "#2563eb" }
  text: { x: 60, y: 35, content: "Ubuntu", font: "bold 24px" }
}
```

**Trường biểu mẫu:**
```
Name: [_______________]{field: name, required: true}
Email: [_______________]{field: email, type: email}
Date: [__/__/____]{field: date, type: date}
```

**Email:**
```
@welcome-email: email {
  to: {customer.email}
  subject: "Welcome, {customer.name}"
  body: include welcome-content
}
```

**Một ngôn ngữ. Văn bản, đồ họa, biểu mẫu, email, tài liệu.**

---

## Kết hợp

Thành phần chứa thành phần.

```
@header: compose {
  include: logo
  include: navigation
  include: search-bar
}

@page-template: compose {
  include: header
  include: content
  include: footer
}

@welcome-letter: compose {
  include: page-template
  content: "Dear {customer.name}, welcome to..."
}

@application-form: compose {
  include: page-template
  content: include form-fields
}
```

**Hiệu ứng lan truyền:**
- Thay đổi logo → header cập nhật
- Header cập nhật → mọi trang cập nhật
- Mọi trang cập nhật → mọi PDF, email, biểu mẫu cập nhật

**Đây là lý do thương hiệu luôn hoàn hảo.** Không phải vì ai đó nhớ cập nhật mọi thứ. Mà vì kiến trúc khiến việc không cập nhật trở nên bất khả thi.

Liên kết đến bất kỳ nội dung hiện có nào. Tái sử dụng trên toàn hệ thống. Một nguồn sự thật duy nhất.

---

## Đầu ra

Từ một nguồn, tạo ra mọi thứ:

| Đầu ra | Điều gì xảy ra |
|--------|--------------|
| Trang web | Kết xuất thành HTML, sẵn sàng cho SEO |
| PDF | Sẵn sàng in, lưu trữ lâu dài |
| Email | Gửi trực tiếp đến người nhận |
| Biểu mẫu web | Tương tác, dữ liệu được thu thập trở lại |
| Biểu mẫu PDF | Điền được hoặc in được, OCR thu thập trở lại |
| Đồ họa SVG | Tài nguyên vector, mọi kích thước |
| Bản đồ | Dữ liệu vị trí của bạn được trực quan hóa |
| Kiosk | Màn hình vật lý, thời gian thực |

**Mọi đầu ra đều tự động dịch.**

Viết bằng tiếng Anh. Đồng nghiệp Đức thấy tiếng Đức. Đối tác Tây Ban Nha thấy tiếng Tây Ban Nha. Khách hàng Nhật thấy tiếng Nhật.

Không phải dịch sau khi hoàn tất. Mà dịch ngay trong quá trình kết xuất. Mọi đầu ra. Mọi ngôn ngữ. Tự động.

---

## Dịch

Dịch thuật không phải là một tính năng. Đó là cách hệ thống vận hành.

| Khả năng | Ý nghĩa |
|------------|---------------|
| Thời gian thực trong trình soạn thảo | Cộng tác đồng thời giữa nhiều ngôn ngữ |
| Tự động dịch đầu ra | Mọi định dạng được kết xuất bằng mọi ngôn ngữ |
| AI ngoại tuyến | Hoạt động không cần internet |
| Theo ngữ cảnh | Phân biệt thuật ngữ giáo dục, pháp lý, y tế |
| Hai chiều | Họ chỉnh sửa bằng ngôn ngữ của họ, bạn thấy ngôn ngữ của bạn |

**Viết một lần. Xuất bản bằng mọi ngôn ngữ.**

---

## Thu thập

Biểu mẫu thu thập dữ liệu trở lại. Vì chúng tôi tạo ra biểu mẫu, chúng tôi biết chính xác vị trí của từng trường.

| Kênh | Điều gì xảy ra |
|---------|--------------|
| Gửi biểu mẫu web | Dữ liệu chảy vào cơ sở dữ liệu của bạn |
| Điền biểu mẫu PDF kỹ thuật số | Dữ liệu được trích xuất vào cơ sở dữ liệu của bạn |
| In biểu mẫu PDF, điền tay, quét | OCR thu thập vào cơ sở dữ liệu của bạn |

Không cần huấn luyện OCR. Không cần ánh xạ trường. Chúng tôi tạo ra biểu mẫu—chúng tôi biết mọi thứ nằm ở đâu.

**Giấy hay kỹ thuật số. Cùng dữ liệu. Cùng đích đến.**

---

## Con người trong vòng lặp: Vận hành hiện trường thời gian thực

Publish không chỉ dành cho tài liệu tĩnh. Đó là cách nhân viên hiện trường tương tác với hệ thống của bạn theo thời gian thực.

### Thiết bị được khai thác toàn diện

Thiết bị của mỗi nhân viên hiện trường trở thành cả cảm biến lẫn màn hình hiển thị:

| Đầu vào (từ hiện trường) | Đầu ra (đến hiện trường) |
|--------------------|-------------------|
| Vị trí GPS | Bản đồ tuyến đường, chỉ đường từng bước |
| Lệnh giọng nói (SST) | Thông báo nhiệm vụ |
| Tương tác chạm | Cập nhật lịch trình |
| Gửi biểu mẫu | Cảnh báo xác nhận |
| Camera (ảnh/video) | Hiển thị trạng thái |

Hệ thống biết mọi nhân viên đang ở đâu. Nhân viên có thể phản hồi—qua giọng nói, chạm, hoặc biểu mẫu.

### Đầu ra 2D

Publish kết xuất toàn bộ nội dung 2D—web, tài liệu, bản đồ, màn hình hiển thị:

| Đầu ra | Trường hợp sử dụng | Ví dụ |
|--------|----------|---------|
| **Trang web** | Cổng khách hàng, theo dõi trạng thái | Trạng thái dịch vụ, theo dõi đơn hàng |
| **Bản đồ tuyến đường** | Điều hướng, theo dõi giao hàng | Thu gom rác thải, điểm giao hàng |
| **Bảng điều khiển** | Tổng quan trạng thái, KPI | Bảng điều phối, hiệu suất đội nhóm |
| **Tài liệu** | Lệnh công việc, báo cáo, biểu mẫu | Bản kê khai PDF, xuất dữ liệu tuân thủ |
| **Kiosk** | Màn hình công cộng, bảng đội nhóm | Trạng thái kho bãi, thực đơn nhà hàng |

Đối với **trực quan hóa không gian 3D**—công trường xây dựng, nhà kho, bản sao số (digital twin)—xem [Nền tảng Spatial](/platform/spatial/).

### Cập nhật bằng giọng nói

Chuyển giọng nói thành văn bản (SST) cho phép tương tác với hệ thống mà không cần dùng tay:

| Nhân viên nói | Hệ thống thực hiện |
|-------------|-------------|
| "123 Main bị chặn do thi công" | Đánh dấu vị trí, sắp xếp lại tuyến đường |
| "Bỏ qua nhiệm vụ tiếp theo, khách hàng không có mặt" | Xóa khỏi hàng đợi, điều chỉnh lịch trình |
| "Thiết bị cần sửa chữa" | Tạo phiếu bảo trì |
| "Hoàn thành nhiệm vụ" | Ghi nhận hoàn thành, kích hoạt quy trình tiếp theo |

Không cần gọi điện cho bộ phận điều phối. Không cần liên lạc qua bộ đàm. Nhân viên nói, AI lắng nghe, hệ thống thích ứng.

### Phân cấp thẩm quyền

Phán đoán của con người được ưu tiên hơn đề xuất của AI:

```
┌─────────────────────────────────────────────────────────────────┐
│                     DECISION AUTHORITY                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1. FIELD WORKER       ───────────────────────▶  HIGHEST       │
│      Eyes on ground, can override anything                       │
│                                                                  │
│   2. DISPATCHER         ───────────────────────▶  HIGH          │
│      Fleet-wide view, can reassign work                          │
│                                                                  │
│   3. AI OPTIMIZER       ───────────────────────▶  MEDIUM        │
│      Suggests routes/tasks, adapts to changes                    │
│                                                                  │
│   4. SCHEDULED PLAN     ───────────────────────▶  LOWEST        │
│      Starting point, always subject to reality                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**AI rất mạnh mẽ—nhưng người đang trực tiếp quan sát tại hiện trường mới có tiếng nói cuối cùng.**

### Hiện diện & Thông báo

Con người không phải lúc nào cũng sẵn sàng. Hệ thống biết ai đang sẵn sàng:

| Khả năng | Chức năng |
|------------|--------------|
| **Theo dõi hiện diện** | Ai trực tuyến, ai ngoại tuyến, ai đang bận |
| **Nhận biết vai trò** | Ai có thể phê duyệt, ai có thể ghi đè, ai đang trực |
| **Định tuyến thông minh** | Thông báo được gửi đến nhân sự đang sẵn sàng trước |
| **Leo thang** | Nếu người phụ trách chính không phản hồi, chuyển lên người dự phòng |
| **Xác nhận** | Theo dõi ai đã thấy gì, họ phản hồi khi nào |

**Kênh thông báo:**

- Thông báo đẩy (di động/máy tính)
- SMS cho cảnh báo khẩn cấp
- Email cho cập nhật không đồng bộ
- Cảnh báo trong ứng dụng kèm âm thanh/rung
- Màn hình kiosk để cả đội cùng theo dõi

Hệ thống thích ứng với mức độ sẵn sàng của con người. Vấn đề khẩn cấp tìm đến đúng người. Vấn đề không khẩn cấp được xếp hàng chờ cho đến khi có người sẵn sàng.

### Khả năng chống chịu khi ngoại tuyến

Công việc hiện trường diễn ra ở những nơi mạng chập chờn. Thiết bị vẫn tiếp tục hoạt động:

- **Mất kết nối** → Ứng dụng tiếp tục chạy với SQLite cục bộ
- **Thay đổi được xếp hàng cục bộ** → Đồng bộ khi kết nối lại
- **Không mất dữ liệu** → Automerge CRDT xử lý xung đột
- **Cập nhật hiện diện** → Đánh dấu ngoại tuyến, thông báo được xếp hàng

Vận hành hiện trường không thể dừng lại vì sự cố mạng. Kiến trúc này đảm bảo điều đó không xảy ra.

[Xem Vận hành hiện trường trong thực tế →](/applications/logistics/)

---

## Kết nối

Cơ sở dữ liệu của bạn. Lược đồ của bạn.

- Đọc dữ liệu để điền vào tài liệu: "Kính gửi {customer.name}..."
- Ghi dữ liệu thu thập được trở lại: dữ liệu biểu mẫu gửi về → các bảng của bạn
- Mọi cơ sở dữ liệu SQL: PostgreSQL, MySQL, SQLite
- Hạ tầng của bạn: tự lưu trữ hoặc đám mây

**Chúng tôi không thay thế backend của bạn. Chúng tôi kết nối với nó.**

---

## Các ngành

| Ngành | Vì sao chọn Publish |
|----------|-------------|
| [**Chính phủ**](/applications/government) | Phục vụ mọi công dân. Giấy + kỹ thuật số. Dễ tiếp cận. Tuân thủ quy định. |
| [**Y tế**](/applications/healthcare) | Ít thủ tục hành chính hơn. Ít sai sót hơn. Chăm sóc tốt hơn. |
| [**Tài chính**](/applications/financial) | Tuân thủ mà không phức tạp. Dấu vết kiểm toán đầy đủ. |
| [**Giáo dục**](/applications/education) | Quản trị hiện đại. Vẫn giữ lựa chọn truyền thống. |
| [**Bảo hiểm**](/applications/insurance) | Mọi kênh. Một hệ thống. Hiện trường và văn phòng hợp nhất. |

---

## Kiến trúc

| Tầng | Chức năng |
|-------|--------------|
| Bộ phân tích DSL | Hiểu văn bản, đồ họa, biểu mẫu, cấu trúc kết hợp |
| Bộ phân giải include | Kết hợp thành phần đệ quy |
| Bộ kết xuất | Đầu ra web, PDF, email, SVG, biểu mẫu |
| Công cụ dịch | AI thời gian thực, hoạt động được khi ngoại tuyến |
| Bộ ánh xạ trường | Biết chính xác vị trí của từng trường biểu mẫu |
| Công cụ OCR | Trích xuất chính xác nhờ bản đồ trường |
| Trình kết nối cơ sở dữ liệu | Đọc và ghi vào lược đồ của bạn |

**Không có định dạng độc quyền. Xuất mọi thứ. Làm chủ nội dung của bạn.**

---

## Chạy gốc trên mọi nền tảng

| Nền tảng | Trải nghiệm |
|----------|------------|
| **Linux** | Ứng dụng desktop gốc, kiosk nhúng, thiết bị biên |
| **Windows** | Ứng dụng Windows gốc |
| **macOS** | Ứng dụng Mac gốc |
| **iOS** | Ứng dụng di động gốc cho nhân viên hiện trường |
| **Android** | Ứng dụng di động gốc cho nhân viên hiện trường |
| **Web** | Trình duyệt hiện đại, đầy đủ chức năng |

Một codebase. Hiệu năng gốc trên mọi nền tảng. Người soạn thảo dùng máy tính để bàn. Nhân viên hiện trường dùng máy tính bảng. Khách hàng dùng điện thoại. Cùng một hệ thống. Cùng dữ liệu.

---

## Vượt ra ngoài tài liệu

Cùng cách tiếp cận một nguồn duy nhất áp dụng cho:

**Bản đồ** — Dữ liệu địa lý của bạn, phong cách của bạn, tích hợp với Google/Apple Maps để mở rộng phạm vi tiếp cận.

**Video** *(sắp ra mắt)* — Máy chủ video của bạn, phản chiếu lên YouTube/Vimeo để dễ được khám phá.

**Lịch** *(sắp ra mắt)* — Máy chủ CalDAV của bạn, đồng bộ với Google/Apple Calendar cho tiện lợi.

**Kiosk** *(sắp ra mắt)* — Màn hình Raspberry Pi, thực đơn nhà hàng, văn phòng chính phủ, bảng hiệu bán lẻ.

Mọi liên kết đều được định tuyến qua [cổng của bạn](#wellknown-gateway-data-sovereignty). Xuất bản LÊN các nền tảng Big Tech. Không bao giờ bị khóa chặt VÀO chúng.

---

## WellKnown Gateway: Chủ quyền dữ liệu

Nội dung bạn xuất bản sẽ đi đâu? Các nền tảng bên ngoài? Đó là lúc **WellKnown Gateway** phát huy tác dụng—đảm bảo bạn làm chủ dữ liệu của mình ngay cả khi xuất bản LÊN các hệ thống khác.

### Vấn đề của việc bị khóa chặt vào nền tảng

Cách tiếp cận truyền thống:
- Dữ liệu khách hàng nằm trong Salesforce
- Tài liệu nằm trong Google Drive
- Biểu mẫu nằm trong Typeform
- Bạn phụ thuộc vào API của từng nhà cung cấp
- Nếu bất kỳ nhà cung cấp nào thay đổi điều khoản, tăng giá hoặc ngừng hoạt động—bạn phải chật vật xoay xở

### Giải pháp WellKnown

```
┌─────────────────────────────────────────────────────────────────┐
│                    YOUR DATA (Source of Truth)                   │
│                                                                  │
│   SQLite on every device + NATS JetStream persistence            │
│   You own it. It's on your infrastructure.                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    WELLKNOWN GATEWAY                             │
│                                                                  │
│   Publishes TO external platforms                                │
│   Syncs FROM external platforms                                  │
│   You control what goes where                                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
   │   Google    │     │   Social    │     │  Industry   │
   │   Drive     │     │   Media     │     │  Platforms  │
   └─────────────┘     └─────────────┘     └─────────────┘
```

**Nguyên tắc cốt lõi:** Bạn xuất bản LÊN các nền tảng. Họ không sở hữu dữ liệu của bạn—bạn mới là người sở hữu.

### Chủ quyền quan hệ khách hàng

Khách hàng của bạn tương tác với thương hiệu của bạn, không phải của nhà cung cấp:

| Không có WellKnown | Có WellKnown |
|-------------------|----------------|
| Khách hàng đăng nhập vào cổng của nhà cung cấp | Khách hàng đăng nhập vào cổng CỦA BẠN |
| Nhà cung cấp sở hữu mối quan hệ | Bạn sở hữu mối quan hệ |
| Nhà cung cấp nắm danh sách khách hàng của bạn | Bạn kiểm soát dữ liệu khách hàng |
| Đổi nhà cung cấp = xây dựng lại | Đổi nhà cung cấp = kết nối lại cổng |

Khách hàng nghĩ họ đang nói chuyện với bạn. Và đúng là vậy. Hệ thống phía sau hoàn toàn vô hình.

### Xây dựng trên các tiêu chuẩn Internet

WellKnown sử dụng chính các tiêu chuẩn vận hành email, lịch và mạng xã hội liên bang:

| Tiêu chuẩn | Mục đích | Ví dụ |
|----------|---------|---------|
| **RFC 8615** | Khám phá URI `/.well-known/` | `yourdomain.com/.well-known/...` |
| **WebFinger** | Khám phá tài khoản/tài nguyên | Tìm người dùng trên nhiều tên miền |
| **OAuth Discovery** | Điểm cuối ủy quyền | Truy cập an toàn cho bên thứ ba |
| **CalDAV/CardDAV** | Đồng bộ lịch/danh bạ | Hoạt động với mọi ứng dụng tuân thủ |
| **ActivityPub** | Giao thức liên bang | Khả năng tương tác với mọi hệ thống |

**Tên miền của bạn là danh tính của bạn.** Giống như email (`you@yourdomain.com`), dữ liệu của bạn nằm tại các URI do bạn kiểm soát. Liên kết liên bang dựa trên tiêu chuẩn nghĩa là bạn có thể kết nối với bất kỳ hệ thống tuân thủ nào—không chỉ của chúng tôi.

Đây là cách web mở vận hành. Chúng tôi chỉ áp dụng nó vào việc xuất bản.

### Khả năng di chuyển dữ liệu

Mọi thứ đồng bộ sang các định dạng tiêu chuẩn do bạn kiểm soát:

- **Tài liệu** → Cơ sở dữ liệu SQLite của bạn
- **Biểu mẫu** → Các bảng cơ sở dữ liệu của bạn
- **Phương tiện** → Kho lưu trữ của bạn
- **Nhật ký kiểm toán** → Hồ sơ tuân thủ của bạn

Xuất bất cứ lúc nào. Di chuyển bất cứ lúc nào. Bạn không bao giờ bị mắc kẹt.

---

## Xây dựng trên Foundation

Publish tự động kế thừa mọi khả năng của [Foundation](/platform/foundation/):

| Khả năng | Ý nghĩa |
|------------|---------------|
| **Ưu tiên ngoại tuyến** | Làm việc không cần internet, đồng bộ khi có kết nối |
| **Triển khai toàn diện** | Web, desktop, di động—một codebase |
| **Tự chủ** | Máy chủ của bạn, dữ liệu của bạn, quy tắc của bạn |
| **Đồng bộ thời gian thực** | Nhiều người cùng chỉnh sửa, tự động giải quyết xung đột |

[Tìm hiểu thêm về Foundation →](/platform/foundation/)

---

## Bắt đầu

Nếu bạn đang duy trì nội dung ở nhiều nơi, chứng kiến thương hiệu dần lệch lạc, chờ đợi bản dịch, hoặc xử lý dữ liệu gửi về một cách thủ công—có một cách tốt hơn.

[Liên hệ với chúng tôi →](/contact/)

---

## Một phần của điều lớn lao hơn

Publish là **tầng 2D** của một ngăn xếp đa chiều hợp nhất—tất cả đều chạy trên web, truy cập được toàn cầu.

### Ngăn xếp 1D → 2D → 3D

```
┌─────────────────────────────────────────────────────────────────┐
│                    SINGLE SOURCE OF TRUTH                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1D: TEXT (Markdown/DSL)                                        │
│       └── Content definition, structured data, semantic markup   │
│                          │                                       │
│                          ▼                                       │
│   2D: PUBLISH                                                    │
│       └── Web pages, documents, forms, maps, dashboards, kiosks  │
│                          │                                       │
│                          ▼                                       │
│   3D: SPATIAL                                                    │
│       └── CAD, digital twins, robotics, simulation               │
│                                                                  │
│   ALL LAYERS: Web-capable, globally accessible                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Tầng | Nền tảng | Phạm vi xử lý |
|-------|----------|-----------------|
| **1D** | DSL/Markdown | Nội dung văn bản, định nghĩa có cấu trúc, nguồn gốc |
| **2D** | Publish | Web, tài liệu, biểu mẫu, bản đồ, kiosk, email |
| **3D** | [Spatial](/platform/spatial/) | CAD, bản sao số, robot, mô phỏng |
| **Dùng chung** | [Foundation](/platform/foundation/) | Đồng bộ ngoại tuyến, CRDT, NATS, triển khai |

**Mọi thứ đều chạy trên web.** Văn bản 1D điều khiển các đầu ra 2D. Các đầu ra 2D có thể nhúng chế độ xem 3D. Các mô hình 3D có thể tạo tài liệu 2D. Tất cả truy cập qua trình duyệt, tất cả được đồng bộ, tất cả từ cùng một nguồn sự thật duy nhất.

### Cách chúng phối hợp với nhau

Đối với các tổ chức làm việc trong không gian 3D—robot, sản xuất, xây dựng—Publish cung cấp tầng tài liệu:

- **Bản vẽ kỹ thuật** được tạo từ mô hình 3D
- **BOM (Định mức vật tư)** được trích xuất từ các cụm lắp ráp
- **Hướng dẫn công việc** với hình ảnh minh họa từng bước
- **Biểu mẫu tuân thủ** liên kết với các bản sửa đổi CAD
- **Kiểm tra hiện trường** với danh sách kiểm tra 2D + ngữ cảnh 3D

Cùng một nguồn sự thật. Đầu ra 2D và 3D luôn đồng bộ. Tất cả đều truy cập được qua web.

[Khám phá Nền tảng Spatial →](/platform/spatial/)
