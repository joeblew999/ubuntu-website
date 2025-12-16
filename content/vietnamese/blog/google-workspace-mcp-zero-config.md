---
title: "Google Workspace thông qua MCP: Tích hợp Trí tuệ Nhân tạo (AI) không cần cấu hình"
meta_title: "Google Workspace thông qua MCP: Tích hợp Trí tuệ Nhân tạo (AI) không cần cấu hình | Ubuntu Software"
description: "Kết nối bất kỳ hệ thống AI nào—đám mây, cục bộ hoặc của riêng bạn—với Gmail, Lịch, Drive, Trang tính, Tài liệu và Trang trình bày chỉ với một lệnh. MCP có nghĩa là không bị ràng buộc bởi nhà cung cấp. Chỉ cần cung cấp địa chỉ email của bạn."
date: 2024-12-15T10:00:00Z
image: "/images/blog/google-mcp-integration.svg"
categories: ["Publish", "AI"]
author: "Gerard Webb"
tags: ["google", "mcp", "automation", "gmail", "calendar", "drive", "ai", "local-ai", "ollama", "vendor-lock-in", "spatial", "sensors"]
draft: true
---

Nếu việc kết nối trợ lý AI của bạn với toàn bộ Google Workspace chỉ cần một lệnh duy nhất thì sao?

Không cần "cấu hình các thông tin đăng nhập này, kích hoạt các API này, thiết lập các quyền này, rồi khởi động lại." Chỉ cần: cung cấp địa chỉ email của bạn, và mọi thứ sẽ hoạt động.

Đó chính là những gì chúng tôi đã xây dựng.

## Vấn đề với việc tích hợp trí tuệ nhân tạo (AI)

Mỗi câu chuyện về tích hợp trí tuệ nhân tạo (AI) đều tuân theo cùng một mô hình mệt mỏi:

1. Tạo thông tin đăng nhập API trong bảng điều khiển dành cho nhà phát triển
2. Kích hoạt một chục API theo cách thủ công
3. Cấu hình màn hình đồng ý OAuth
4. Cấu hình các URL chuyển hướng
5. Lưu trữ thông tin bí mật ở một nơi an toàn
6. Kết nối cấu hình kết nối
7. Khởi động lại các công cụ của bạn
8. Kiểm tra lỗi để xác định nguyên nhân tại sao nó không hoạt động
9. Bạn đã bỏ qua bước 3b
10. Bắt đầu lại

Khi bạn hoàn thành, bạn đã dành nhiều thời gian cho việc cài đặt hơn là thời gian bạn sẽ tiết kiệm được trong tháng tới. Và bạn vẫn chưa chắc chắn liệu nó có hoạt động đúng cách hay không.

## WellKnown + MCP: Một phương pháp tiếp cận khác biệt

[Model Context Protocol (MCP)]({{< relref "/platform/foundation" >}}) cung cấp một phương thức tiêu chuẩn cho các trợ lý AI truy cập các công cụ và dữ liệu bên ngoài. Dự án WellKnown mở rộng điều này: thiết lập tự động xử lý toàn bộ sự phức tạp phía sau hậu trường.

Dưới đây là quy trình cài đặt đầy đủ:

```bash
task google-mcp:setup
```

Đó là tất cả. Hệ thống:

1. **Cài đặt máy chủ MCP** (nếu chưa có)
2. **Mở xác thực có hướng dẫn** trong trình duyệt của bạn
3. **Lưu trữ thông tin đăng nhập một cách an toàn** trong hệ thống cục bộ của bạn*
4. **Cấu hình trợ lý AI của bạn** tự động
5. **Kiểm tra xem mọi thứ hoạt động bình thường*

Không kích hoạt API thủ công. Không chỉnh sửa tệp JSON. Không cần khởi động lại.

## Bất kỳ AI nào. Tích hợp giống nhau.

Đây là điều khiến MCP khác biệt cơ bản so với các giải pháp tích hợp độc quyền: **nó tương thích với bất kỳ hệ thống AI nào**.

| AI Type | Examples | Same MCP Integration |
|---------|----------|---------------------|
| **Cloud AI** | Claude, GPT-4, Gemini | ✓ |
| **Local AI** | Ollama, LM Studio, llama.cpp | ✓ |
| **Hybrid** | Cloud reasoning + local execution | ✓ |
| **Specialized** | Spatial reasoning, domain-specific models | ✓ |
| **Your own** | Custom models, fine-tuned deployments | ✓ |

Máy chủ MCP không quan tâm đến việc AI nào đang gọi nó. Cloud Claude yêu cầu lịch của bạn? Hoạt động bình thường. Local Llama chạy trên laptop của bạn? Cùng một tích hợp. Mô hình riêng của công ty bạn chạy trong trung tâm dữ liệu của bạn? Cấu hình giống hệt.

Đây là điều ngược lại với tình trạng bị khóa vào nhà cung cấp. Bạn chỉ cần tích hợp Google Workspace một lần, và nó sẽ hoạt động với bất kỳ mô hình AI nào bạn chọn - hôm nay, ngày mai hoặc thậm chí sau năm năm nữa khi có những mô hình hoàn toàn mới xuất hiện.

**Chuyển đổi nhà cung cấp AI mà không cần thay đổi các tích hợp hiện có.** Lớp MCP vẫn giữ nguyên trong khi bạn thử nghiệm các mô hình khác nhau, nâng cấp khả năng hoặc chuyển đổi giữa xử lý trên đám mây và xử lý tại chỗ dựa trên yêu cầu về quyền riêng tư.

### Trí tuệ nhân tạo chuyên biệt: Suy luận không gian

Nền tảng không gian của chúng tôi [Spatial Platform]({{< relref "/platform/spatial" >}}) bao gồm trí tuệ nhân tạo (AI) cục bộ với các khả năng mà các mô hình đám mây thông thường không có:

- **Tư duy không gian** - Hiểu các mối quan hệ 3D, các ràng buộc hình học, và trình tự lắp ráp
- **Tích hợp cảm biến* - Xử lý dữ liệu thời gian thực từ các thiết bị IoT, camera và thiết bị công nghiệp
- **Kiến thức chuyên môn** - Định dạng CAD, dung sai sản xuất, tiêu chuẩn kỹ thuật xây dựng

Khi trí tuệ nhân tạo (AI) chuyên biệt địa phương này kết nối với Google Workspace của bạn thông qua MCP, bạn sẽ có được những kết hợp mạnh mẽ:

> "Trích xuất dữ liệu cảm biến từ Drive, phân tích các mẫu nhiệt, và soạn thảo cảnh báo bảo trì gửi đến đội vận hành."

> "Tạo bản trình bày từ bản vẽ CAD, ghi chú các thành phần vượt quá giới hạn dung sai dựa trên bảng kiểm tra."

Cùng một tích hợp MCP hoạt động với Claude hoặc GPT cũng hoạt động với các mô hình cục bộ có khả năng nhận thức không gian — ngoại trừ bây giờ AI của bạn hiểu được hình học, không chỉ văn bản.

### Bạn sẽ nhận được gì

Sau khi kết nối, trợ lý AI của bạn có quyền truy cập gốc vào:

| Service | Capabilities |
|---------|-------------|
| **Gmail** | Read, search, send, draft emails |
| **Calendar** | View events, create meetings, check availability |
| **Drive** | List files, search, download, upload |
| **Sheets** | Read data, update cells, query ranges |
| **Docs** | Read content, create documents, edit text |
| **Slides** | Create presentations, add slides, export PDF |

Trong suốt cuộc trò chuyện tự nhiên:

> "Hãy cho tôi xem các email từ tuần trước liên quan đến đề xuất dự án."

> "Tôi có những cuộc họp nào vào ngày mai?"

> "Tạo một tài liệu Google Doc tóm tắt tệp PDF đính kèm"

> "Thêm một hàng vào bảng tính ngân sách với các số liệu sau"

## Hỗ trợ nhiều tài khoản

Đây là phần thú vị cho các doanh nghiệp: **bạn có thể kết nối nhiều tài khoản Google**.

Tài khoản Gmail cá nhân. Tài khoản Google Workspace cho công việc. Tài khoản khách hàng. Mỗi tài khoản được xác thực độc lập và có thể truy cập bởi trợ lý AI của bạn với ngữ cảnh phù hợp.

```
task google-mcp:auth
# Authenticate first account

task google-mcp:auth
# Authenticate second account
# The system tracks both
```

Trợ lý AI của bạn có thể hoạt động trên nhiều tài khoản:

> "Kiểm tra lịch làm việc và lịch cá nhân của tôi để xem có xung đột nào vào thứ Ba tuần sau không."

> "Chuyển tiếp email của khách hàng đó đến tài khoản công việc của tôi"

> "Hiển thị danh sách các tệp gần đây từ cả Drive cá nhân và Drive doanh nghiệp của tôi"

## Kết nối Xuất bản

Điều này tích hợp trực tiếp với [hệ thống xuất bản nguồn duy nhất của WellKnown]({{< relref "/blog/one-source-every-screen" >}}). Quy trình xuất bản nội dung của bạn giờ đây có Google Workspace là đích đến chính thức:

**Nguồn Markdown → Chuyển đổi WellKnown → Google Docs**

Viết bằng Markdown, xuất bản lên Google Docs với định dạng được giữ nguyên. Không sao chép và dán. Không chỉnh sửa định dạng thủ công. Nguồn dữ liệu tạo ra trang web của bạn cũng có thể tự động cập nhật vào Google Drive.

**Dữ liệu bảng tính → Quy trình xuất bản**

Nhập dữ liệu trực tiếp từ Sheets vào nội dung của bạn. Danh mục sản phẩm, bảng giá, danh sách đội ngũ—luôn cập nhật, luôn đồng bộ.

**Lịch → Lập lịch tự động**

Các bài đăng trên blog được lên lịch cho các ngày trong tương lai có thể kích hoạt các sự kiện lịch, lời nhắc nhở và thậm chí là các email nháp gửi đến các bên liên quan.

## Tự chủ, không phụ thuộc

Đây là điểm khác biệt quan trọng so với các tích hợp đám mây thông thường: **thông tin đăng nhập và dữ liệu của bạn vẫn được lưu trữ trên máy tính của bạn**.

Máy chủ MCP chạy trên máy tính cục bộ. Các token xác thực được lưu trữ trong thư mục cá nhân của bạn. Không có dịch vụ đám mây trung gian nào truy cập dữ liệu Google của bạn. Không có bên thứ ba nào có quyền truy cập vào email hoặc tệp tin của bạn.

Bạn đang sử dụng các dịch vụ của Google, nhưng thông qua hạ tầng do bạn kiểm soát. Nếu bạn quyết định ngắt kết nối, chỉ cần xóa các thông tin đăng nhập cục bộ và bạn đã hoàn tất. Không cần liên hệ với nhà cung cấp. Không cần yêu cầu xuất dữ liệu.

Điều này phù hợp với các nguyên tắc chúng ta đã thảo luận trong [hệ thống email tự chủ]({{< relref "/blog/self-sovereign-email-ai-automation" >}}): sử dụng các dịch vụ bạn cần, nhưng vẫn duy trì quyền kiểm soát đối với lớp tích hợp.

## Các cấp độ tự động hóa AI được áp dụng

Bạn còn nhớ các [mức độ tự động hóa]({{< relref "/blog/self-sovereign-email-ai-automation" >}})? Chúng được áp dụng trực tiếp tại đây:

**Cấp độ 1 - Nhận dạng mẫu**
> "Bạn có 3 email chưa đọc từ đội của mình. Một trong số đó có ghi 'cấp bách'. Bạn có muốn tôi tóm tắt nội dung của chúng không?"

**Cấp độ 2 - Tạo bản nháp**
> dựa trên email của Sarah về hạn chót nộp đề xuất, tôi đã soạn thảo một phản hồi xác nhận việc giao hàng vào thứ Tư. Bạn có thể xem qua không?

**Cấp độ 3 - Hành động tự chủ**
> "Bảng tính báo cáo hàng tuần đã được cập nhật. Tôi đã xuất bản tóm tắt sang tài liệu cho các bên liên quan và lên lịch gửi email phân phối vào lúc 9 giờ sáng thứ Hai."

**Cấp độ 4 - Trí tuệ chủ động**
> "Tôi thấy anh có cuộc họp với khách hàng vào ngày mai nhưng không có tài liệu chuẩn bị. Tôi có nên tạo một tài liệu từ ba chuỗi email gần đây và lịch trình dự án trong Drive không?"

Mỗi cấp độ sẽ trở nên mạnh mẽ hơn khi AI có thể truy cập toàn bộ Google Workspace của bạn, không chỉ các dịch vụ riêng lẻ.

## Kiến trúc kỹ thuật

Đối với những ai muốn hiểu rõ những gì đang diễn ra bên trong:

```
┌─────────────────────────────────────────────────────────────┐
│                    Your Local Machine                        │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │   Claude     │───▶│  MCP Server  │───▶│   Google     │  │
│  │   Code       │◀───│  (local)     │◀───│   APIs       │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│         │                   │                              │
│         │                   ▼                              │
│         │           ┌──────────────┐                       │
│         │           │   Tokens     │                       │
│         │           │   (~/.google-│                       │
│         │           │   mcp-...)   │                       │
│         │           └──────────────┘                       │
│         │                                                  │
│         ▼                                                  │
│  ┌──────────────────────────────────────────────────────┐ │
│  │              .vscode/mcp.json                         │ │
│  │   {                                                   │ │
│  │     "servers": {                                      │ │
│  │       "google": { "command": "google-mcp-server" }   │ │
│  │     }                                                 │ │
│  │   }                                                   │ │
│  └──────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

MCP server đóng vai trò là cầu nối giữa trợ lý AI của bạn và các API của Google. Nó xử lý:

- **Cập nhật token OAuth** (tự động, không cần can thiệp)
- **Gộp yêu cầu** (sử dụng API hiệu quả)
- **Định dạng phản hồi** (được cấu trúc để AI có thể xử lý)
- **Định tuyến nhiều tài khoản** (hướng các yêu cầu đến tài khoản đúng)

Tất cả đều chạy trên máy cục bộ. Tất cả đều nằm dưới sự kiểm soát của bạn.

## Bắt đầu

### Bắt đầu nhanh (Đa số người dùng)

```bash
# From your project directory
task google-mcp:setup
```

Làm theo các hướng dẫn. Khởi động lại trợ lý AI của bạn. Xong.

### Bước từng bước (Nếu bạn muốn kiểm soát)

```bash
# 1. Install the MCP server
task google-mcp:install

# 2. Authenticate with Google
task google-mcp:auth

# 3. Add to your AI assistant config
task google-mcp:claude:add

# 4. Verify everything works
task google-mcp:status
```

### Kiểm tra các thiết bị đang kết nối

```bash
task google-mcp:accounts:list
# Shows all authenticated Google accounts

task google-mcp:status
# Full status: binary, accounts, configuration
```

### Xóa quyền truy cập

```bash
task google-mcp:reset CONFIRM=y
# Removes all local credentials and configuration
# Also opens Google Console to revoke OAuth access
```

## Ngoài Google

Cùng một mô hình áp dụng cho các dịch vụ khác. MCP cung cấp giao thức; WellKnown cung cấp cài đặt tự động. Chúng tôi đang phát triển tích hợp cho:

- **GitHub** - Vấn đề, Yêu cầu kéo, Tìm kiếm mã nguồn
- **Notion** - Trang, cơ sở dữ liệu, không gian làm việc
- **Slack** - Kênh, tin nhắn, chuỗi tin nhắn
- **Linear** - Vấn đề, dự án, lộ trình

Mỗi thiết bị đều tuân theo nguyên tắc tương tự: một lệnh để kết nối, truy cập đầy đủ vào AI, và điều khiển cục bộ.

## Những thay đổi

Khi trợ lý AI của bạn có thể truy cập mượt mà vào Google Workspace, quy trình làm việc sẽ thay đổi một cách cơ bản:

**Trước đây**: Bạn phải chuyển đổi giữa các công cụ. Kiểm tra email. Mở lịch. Tìm kiếm trên Drive. Sao chép dữ liệu vào AI. Nhận phản hồi. Dán lại.

**Sau đó:** Bạn mô tả những gì bạn cần. Trí tuệ nhân tạo (AI) sẽ điều hướng giữa các dịch vụ, thu thập thông tin bối cảnh, thực hiện các hành động và báo cáo kết quả.

Chi phí nhận thức khi chuyển đổi công cụ biến mất. Bạn suy nghĩ về kết quả, không phải giao diện.

## Quý vị quan tâm?

We're currently onboarding early access partners to this integration. If you're using [Claude Code](https://docs.anthropic.com/en/docs/claude-code/overview) or another MCP-compatible assistant and want zero-config Google Workspace access, [get in touch]({{< relref "/contact" >}})—we'll help you get connected.
