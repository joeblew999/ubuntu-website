---
title: "Tại sao Tiêu chuẩn Mở lại chiến thắng"
meta_title: "Tại sao Tiêu chuẩn Mở lại chiến thắng | Ubuntu Software"
description: "STEP, IFC và vụ kiện chống lại việc khóa độc quyền đối với dữ liệu thiết kế và kỹ thuật 3D."
date: 2024-09-20T05:00:00Z
image: "/images/blog/open-standards.svg"
categories: ["Industry", "Standards"]
author: "Gerard Webb"
tags: ["open-standards", "step", "ifc", "interoperability", "cad"]
draft: false
---

Mỗi thập kỷ một lần, một ngành công nghiệp lại học được bài học tương tự: việc khóa chặt công nghệ độc quyền không thể mở rộng quy mô.

Internet đã học được điều đó. Phần mềm doanh nghiệp đã học được điều đó. Điện toán đám mây đã học được điều đó.

Bây giờ là lúc để thiết kế và kỹ thuật 3D học hỏi.

## Tình hình hiện tại

Hãy thử thí nghiệm này: Lấy một mô hình 3D từ một hệ thống CAD chính và mở nó trong một hệ thống CAD khác.

Bạn sẽ tìm thấy:
- **Mất hình học**: Các đặc trưng không được chuyển đổi. Các ràng buộc biến mất.
- **Thông tin metadata bị mất**: Tất cả thông tin kỹ thuật—vật liệu, dung sai, mối quan hệ lắp ráp—đã bị mất.
- **Sửa chữa thủ công**: Ai đó dành hàng giờ để tái tạo lại những gì đã tồn tại.

Đây không phải là lỗi. Đây là mô hình kinh doanh.

**Nhà cung cấp hưởng lợi từ việc khóa khách hàng. Người dùng phải chịu thiệt hại từ điều đó.**

## Chi phí của các định dạng độc quyền

### Ma sát trong hợp tác

Khi nhà cung cấp robot của bạn sử dụng Hệ thống A, nhà thiết kế cơ sở của bạn sử dụng Hệ thống B và nhóm mô phỏng của bạn sử dụng Hệ thống C, mỗi lần chuyển giao đều là một quá trình dịch thuật.

Thông tin bị suy giảm sau mỗi lần chuyển đổi. Các kỹ sư phải dành thời gian cho việc xử lý tệp thay vì tập trung vào công việc kỹ thuật.

### Sự phụ thuộc vào nhà cung cấp

Toàn bộ lịch sử thiết kế của bạn được lưu trữ trong một định dạng mà chỉ một nhà cung cấp duy nhất kiểm soát. Họ quyết định giá cả. Họ quyết định thời gian nâng cấp. Họ quyết định khi nào các tính năng sẽ bị loại bỏ.

Quyền sở hữu trí tuệ (IP) về kỹ thuật của bạn đang bị giữ làm con tin.

### Rào cản đổi mới

Bạn muốn phát triển các công cụ AI dựa trên dữ liệu thiết kế của mình? Chúc may mắn khi truy cập dữ liệu đó thông qua các API độc quyền thay đổi theo từng phiên bản.

Muốn tích hợp với mô phỏng vật lý mới nhất? Tốt nhất là hy vọng nhà cung cấp CAD của bạn có hợp tác.

Sự sáng tạo chết tại ranh giới định dạng.

## Giải pháp mở

### STEP: Tiêu chuẩn hình học

ISO 10303 (STEP) đã tồn tại từ những năm 1990. Nó nhàm chán. Nó hoạt động.

STEP ghi lại:
- hình học 3D với độ chính xác tuyệt đối
- Cấu trúc và mối quan hệ của các bộ phận lắp ráp
- Thông tin sản xuất sản phẩm (PMI)
- Tính chất vật liệu

Nó không hoàn hảo. Nhưng nó phổ quát.

### IFC: Các tòa nhà biết nói

Các Lớp Cơ sở Ngành (IFC) đóng vai trò tương tự như STEP đối với sản phẩm, nhưng áp dụng cho các công trình xây dựng.

Mọi bức tường, cửa, không gian và hệ thống—được định nghĩa theo định dạng mở mà bất kỳ phần mềm nào cũng có thể đọc và ghi.

Khả năng tương tác BIM không phải là điều không thể. IFC làm cho điều đó trở nên khả thi.

### Công nghệ mới nổi

Các tiêu chuẩn mở hiện đại vượt ra ngoài hình học tĩnh:

- **glTF**: Định dạng 3D nhẹ cho hiển thị và thực tế tăng cường/thực tế ảo (AR/VR)
- **USD**: Mô tả cảnh cho mô phỏng và hiển thị
- **SDF**: Định nghĩa robot và môi trường
- **URDF**: Định dạng mô tả robot

Một hệ sinh thái đang hình thành. Các công cụ được xây dựng trên nền tảng mở có thể tham gia.

## Những gì Tiêu chuẩn Mở cho phép

### Cạnh tranh thực sự

Khi dữ liệu của bạn không bị khóa, bạn có thể lựa chọn công cụ dựa trên khả năng, chứ không phải sự phụ thuộc.

Các nhà cung cấp cạnh tranh dựa trên tính năng, chứ không phải dựa trên mức độ phức tạp mà họ gây ra cho quá trình di chuyển.

### Sáng tạo hệ sinh thái

Các định dạng mở cho phép phát triển một hệ sinh thái các công cụ chuyên dụng:
- Trợ lý AI hoạt động trên nhiều nền tảng
- Các bộ mô phỏng có khả năng xử lý bất kỳ hình học nào
- Các công cụ hợp tác không yêu cầu mọi người phải sở hữu cùng một giấy phép

### Chuẩn bị cho tương lai

Các tổ chức tiêu chuẩn hoạt động chậm chạp. Đó là một đặc điểm.

Tệp STEP từ năm 2000 vẫn có thể mở được cho đến ngày nay. Liệu định dạng độc quyền của bạn từ năm 2020 có thể mở được vào năm 2040?

## Thực tế lai

Hãy thực tế: các quy trình làm việc dựa hoàn toàn trên tiêu chuẩn mở hiện vẫn chưa tồn tại.

Chiến lược thực sự là:
1. **Định dạng gốc cho công việc thực tế**: Sử dụng công cụ phù hợp nhất cho từng công việc
2. **Định dạng mở cho việc trao đổi**: Định dạng tiêu chuẩn tại mọi điểm chuyển giao
3. **Định dạng mở cho lưu trữ lâu dài: Lưu trữ lâu dài trong các định dạng mà bạn kiểm soát*

Đây không phải là chủ nghĩa lý tưởng. Đây là quản lý rủi ro.

## Sự chuyển đổi trong ngành

Sự hứng khởi đang ngày càng tăng cao:

**Quy định của chính phủ**: Ngày càng có nhiều cơ quan yêu cầu sử dụng định dạng mở cho việc mua sắm và lưu trữ.

**Liên minh ngành**: Các tổ chức như buildingSMART đang thúc đẩy việc áp dụng IFC.

**Yêu cầu về AI**: Học máy cần dữ liệu đào tạo không bị khóa lại.

**Hợp tác đám mây**: Các nền tảng hợp tác thời gian thực lựa chọn nền tảng mở.

Các nhà cung cấp chấp nhận các tiêu chuẩn mở sẽ thành công. Những nhà cung cấp chống đối sẽ bị loại bỏ.

## Chuyển đổi

Nếu bạn đang bắt đầu từ đầu, hãy xây dựng trên nền tảng mở:
- Chọn các công cụ có hỗ trợ định dạng mở mạnh mẽ
- Yêu cầu tuân thủ tiêu chuẩn trong hợp đồng với nhà cung cấp
- Thiết lập các điểm kiểm tra định dạng mở trong quy trình làm việc của bạn
- Lưu trữ dữ liệu theo định dạng bạn kiểm soát, không phải định dạng kiểm soát bạn

Nếu bạn đang di chuyển, hãy bắt đầu từ các ranh giới:
- Các tích hợp mới sử dụng các định dạng mở
- Các dự án mới thử nghiệm quy trình làm việc mở
- Di chuyển dần dần khi các công cụ và quy trình làm việc ngày càng hoàn thiện

## Bức tranh toàn cảnh

Tiêu chuẩn mở không liên quan đến công nghệ. Chúng liên quan đến quyền lực.

Ai kiểm soát dữ liệu kỹ thuật của bạn? Ai quyết định bạn có thể sử dụng công cụ nào? Ai sở hữu lịch sử thiết kế của bạn?

Định dạng độc quyền: nhà cung cấp.

Câu trả lời cho tiêu chuẩn mở: Bạn làm.

Đó là lý do tại sao các tiêu chuẩn mở chiến thắng. Không phải vì chúng kỹ thuật vượt trội (mặc dù thường là như vậy). Mà vì chúng điều chỉnh động lực một cách chính xác.

Dữ liệu của bạn. Lựa chọn của bạn. Tương lai của bạn.

---

*We built our platform on STEP, IFC, and open APIs. [See how it works →](/vi/platform)*
