---
title: "Đăng"
meta_title: "Content Management Platform | Ubuntu Software"
description: "Single-source CMS for multi-channel publishing. Một DSL for text, graphics, forms, emails. Auto-translated outputs. Perfect branding by architecture."
image: "/images/publish.svg"
draft: false
---

## One Language. Every Đầu ra. Every Language.

A single-source content management system. Define everything in one DSL. Soạn thảo components that include components. Output to any channel. Auto-translate to any language. Chụp data back.

Thương hiệu hoàn hảo. Không phải bằng kỷ luật. Mà bằng kiến trúc.

---

## Vấn đề

Các tổ chức đăng tải nội dung trên nhiều kênh khác nhau: trang web, tệp PDF, email, biểu mẫu, biển báo.

Mỗi kênh được quản lý riêng biệt. Nội dung không đồng nhất. Thương hiệu bị phân mảnh. Dịch thuật chậm trễ. Các biểu mẫu gửi đi phải nhập lại thủ công.

**Một thay đổi có nghĩa là phải cập nhật năm hệ thống. Hoặc là không được cập nhật gì cả.**

---

## One DSL

Tất cả đều được định nghĩa bằng cùng một ngôn ngữ.

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

**Các trường biểu mẫu:**
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

## Compose

Các thành phần bao gồm các thành phần.

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

**Dòng chảy:**
- Thay đổi logo → Cập nhật tiêu đề
- Cập nhật tiêu đề → Cập nhật trên mọi trang
- Mỗi trang được cập nhật → mỗi tệp PDF, email, biểu mẫu được cập nhật

**Đó là lý do tại sao thương hiệu luôn hoàn hảo.** Không phải vì ai đó nhớ cập nhật mọi thứ. Mà vì kiến trúc khiến việc không cập nhật trở nên bất khả thi.

Liên kết đến bất kỳ nội dung hiện có nào. Tái sử dụng trên toàn hệ thống. Một nguồn thông tin chính xác duy nhất.

---

## Output

Từ một nguồn duy nhất, tạo ra mọi thứ:

| Output | What Happens |
|--------|--------------|
| Web pages | Rendered to HTML, SEO-ready |
| PDFs | Print-ready, archival |
| Emails | Sent directly to recipients |
| Web forms | Interactive, data captures back |
| PDF forms | Fillable or printable, OCR captures back |
| SVG graphics | Vector assets, any size |
| Maps | Your location data visualized |
| Kiosks | Physical displays, real-time |

**Mọi kết quả đầu ra đều được dịch tự động.**

Viết bằng tiếng Anh. Đồng nghiệp người Đức sẽ thấy tiếng Đức. Đối tác người Tây Ban Nha sẽ thấy tiếng Tây Ban Nha. Khách hàng người Nhật sẽ thấy tiếng Nhật.

Not translated after the fact. Dịchd as part of rendering. Every output. Every language. Automatically.

---

## Translate

Dịch thuật không phải là một tính năng. Đó là cách hệ thống hoạt động.

| Capability | What It Means |
|------------|---------------|
| Real-time in editor | Collaborate across languages simultaneously |
| Auto-translate outputs | Every format renders in every language |
| Offline AI | Works without internet |
| Contextual | Knows education vs. legal vs. medical terminology |
| Bi-directional | They edit in their language, you see yours |

**Viết một lần. Xuất bản bằng mọi ngôn ngữ.**

---

## Capture

Các biểu mẫu thu thập dữ liệu. Vì chúng tôi đã tạo ra biểu mẫu, chúng tôi biết chính xác vị trí của từng trường.

| Channel | What Happens |
|---------|--------------|
| Web form submitted | Data flows to your database |
| PDF form filled digitally | Data extracts to your database |
| PDF form printed, filled by hand, scanned | OCR captures to your database |

Không cần đào tạo OCR. Không cần ánh xạ trường. Chúng tôi đã tạo biểu mẫu—chúng tôi biết mọi thứ nằm ở đâu.

**Giấy hay kỹ thuật số. Dữ liệu giống nhau. Điểm đến giống nhau.**

---

## Kết nối

Cơ sở dữ liệu của bạn. Cấu trúc cơ sở dữ liệu của bạn.

- Đọc dữ liệu để điền vào tài liệu: "Kính gửi {customer.name}..."
- Ghi lại dữ liệu đã thu thập: các biểu mẫu gửi → các bảng của bạn
- Bất kỳ cơ sở dữ liệu SQL nào: PostgreSQL, MySQL, SQLite
- Hạ tầng của bạn: tự quản lý hoặc đám mây

**Chúng tôi không thay thế hệ thống backend của bạn. Chúng tôi kết nối với nó.**

---

## Các ngành công nghiệp

| Industry | Why Publish |
|----------|-------------|
| [**Government**](/applications/government) | Serve every citizen. Paper + digital. Accessible. Compliant. |
| [**Healthcare**](/applications/healthcare) | Less administration. Fewer errors. Better care. |
| [**Financial**](/applications/financial) | Compliance without complexity. Complete audit trails. |
| [**Education**](/applications/education) | Modern administration. Traditional options. |
| [**Insurance**](/applications/insurance) | Every channel. One system. Field to office unified. |

---

## Kiến trúc

| Layer | What It Does |
|-------|--------------|
| DSL parser | Understands text, graphics, forms, composition |
| Include resolver | Recursive component composition |
| Renderers | Web, PDF, email, SVG, form outputs |
| Translation engine | Real-time AI, offline capable |
| Field mapper | Knows exact position of every form field |
| OCR engine | Precise extraction using field map |
| Database connectors | Read and write to your schema |

**Không có định dạng độc quyền. Xuất tất cả. Sở hữu nội dung của bạn.**

---

## Hỗ trợ trên mọi nền tảng

| Platform | Experience |
|----------|------------|
| **Linux** | Native desktop, embedded kiosks, edge devices |
| **Windows** | Native Windows application |
| **macOS** | Native Mac application |
| **iOS** | Native mobile app for field workers |
| **Android** | Native mobile app for field workers |
| **Web** | Modern browser, full functionality |

Một mã nguồn duy nhất. Hiệu suất bản địa trên mọi nền tảng. Người dùng văn phòng sử dụng máy tính để bàn. Nhân viên hiện trường sử dụng máy tính bảng. Khách hàng sử dụng điện thoại di động. Cùng một hệ thống. Cùng một dữ liệu.

---

## Ngoài các tài liệu

Cùng một phương pháp nguồn duy nhất cũng áp dụng cho:

**Bản đồ** — Dữ liệu địa lý của bạn, phong cách thiết kế của bạn, tích hợp với Google/Apple Maps để mở rộng phạm vi tiếp cận.

**Video** *(sắp ra mắt)* — Máy chủ video của bạn, gương cho YouTube/Vimeo để phát hiện.

**Lịch** *(sắp ra mắt)* — Máy chủ CalDAV của bạn, đồng bộ với Google/Apple Calendar cho tiện lợi.

**Kiosk** *(sắp ra mắt)* — Màn hình Raspberry Pi, thực đơn nhà hàng, văn phòng chính phủ, biển hiệu bán lẻ.

All links route through [your gateway](/platform/foundation/#wellknown-gateway). Publish TO Big Tech platforms. Never be locked IN.

---

## Xây dựng trên nền tảng

Publish inherits all [Foundation](/platform/foundation/) capabilities automatically:

| Capability | What It Means |
|------------|---------------|
| **Offline-first** | Work without internet, sync when connected |
| **Universal deployment** | Web, desktop, mobile—one codebase |
| **Self-sovereign** | Your servers, your data, your rules |
| **Real-time sync** | Multiple editors, automatic conflict resolution |
| **Wellknown Gateway** | Publish TO Big Tech, never locked IN |

[Tìm hiểu thêm về Quỹ →](/platform/foundation/)

---

## Bắt đầu

Nếu bạn đang quản lý nội dung ở nhiều nơi, theo dõi sự thay đổi thương hiệu, chờ đợi bản dịch hoặc xử lý thủ công các bản nộp—có một cách tốt hơn.

[Liên hệ với chúng tôi →](/contact/)

---

## Một phần của điều gì đó lớn lao hơn

Publish là lớp 2D của nền tảng phần mềm Ubuntu.

Đối với các tổ chức hoạt động trong lĩnh vực 3D—robot, sản xuất, xây dựng, mô hình số—Publish cung cấp lớp tài liệu. Bản vẽ kỹ thuật, danh sách vật liệu (BOM), hướng dẫn công việc, biểu mẫu tuân thủ—tất cả đều được tạo ra từ một nguồn duy nhất và đồng bộ với mô hình 3D.

[Khám phá Nền tảng Không gian →](/platform/spatial/)
