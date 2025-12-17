---
title: "Cuộc khủng hoảng hợp tác trong thiết kế 3D"
meta_title: "Cuộc khủng hoảng hợp tác trong thiết kế 3D | Ubuntu Software"
description: "Đội ngũ toàn cầu, chuyên môn phân tán, thời hạn thực tế—và các công cụ được thiết kế cho người dùng đơn lẻ trên các máy tính riêng lẻ. Điều gì đó cần phải thay đổi."
date: 2024-11-22T05:00:00Z
image: "/images/blog/collaboration-3d.svg"
categories: ["Industry", "Architecture"]
author: "Gerard Webb"
tags: ["collaboration", "distributed-teams", "crdt", "real-time"]
draft: false
---

Nhà máy nằm ở Việt Nam. Các kỹ sư kết cấu ở Đức. Các kiến trúc sư ở Úc. Khách hàng ở Singapore.

Mọi người cần làm việc trên cùng một mô hình. Mọi người cần xem cùng một trạng thái hiện tại. Mọi người cần các thay đổi được cập nhật ngay lập tức.

Và thế nhưng, vào năm 2024, quy trình làm việc chính vẫn là: gửi tệp qua email, chờ đợi, hy vọng không ai khác thay đổi gì, hợp nhất thủ công, lặp lại.

**Điều này bị hỏng.**

## Vấn đề tệp tin

CAD truyền thống được xây dựng dựa trên các tệp tin. Lưu. Đóng. Gửi email. Tải xuống. Mở. Chỉnh sửa. Lưu. Gửi email lại.

Mỗi lần chuyển giao đều tiềm ẩn rủi ro:
- "Phiên bản nào là phiên bản hiện tại?"
- "Bạn có thấy những thay đổi của tôi từ hôm qua không?"
- "Chúng ta đã làm việc trên cùng một phần—bây giờ thì sao?"
- "Tệp tin bị khóa — ai đang giữ nó?"

Hợp tác dựa trên tệp không thể mở rộng. Không thể vượt qua các múi giờ. Không thể vượt qua các tổ chức. Không thể đáp ứng tốc độ mà các dự án hiện đại yêu cầu.

## Cơn ác mộng của The Merge

Khi hai kỹ sư chỉnh sửa cùng một mô hình, ai đó phải giải quyết các khác biệt.

Trong mã nguồn (code), chúng ta có Git. So sánh (diff), hợp nhất (merge), giải quyết xung đột (resolve conflicts). Nó hoạt động—hầu hết.

Trong hình học 3D? Các công cụ gần như không tồn tại. So sánh thủ công. Kiểm tra bằng mắt. Hy vọng bạn phát hiện ra các sai sót. Cầu mong bạn không gây ra lỗi.

Mỗi lần hợp nhất đều tiềm ẩn rủi ro. Mỗi lần chuyển giao đều có thể dẫn đến thảm họa.

## Tại sao điều này lại quan trọng vào lúc này?

Áp lực đang gia tăng:

**Các dự án ngày càng phân tán hơn.** COVID-19 đã làm việc từ xa trở thành tiêu chuẩn. Chuỗi cung ứng đã trở nên toàn cầu. Nhân tài có mặt ở khắp mọi nơi — không chỉ trong văn phòng của bạn.

**Thời gian thực hiện bị rút ngắn.** Xây dựng mô-đun hứa hẹn tốc độ. Việc triển khai robot đòi hỏi sự lặp lại. "Chúng ta sẽ giải quyết vấn đề tại hiện trường" không khả thi khi nhà máy nằm cách xa đại dương.

**Các bên liên quan ngày càng gia tăng.** Kiến trúc sư, kỹ sư, nhà sản xuất, khách hàng, cơ quan quản lý—tất cả đều cần có tầm nhìn. Tất cả đều cần đóng góp ý kiến. Tất cả đều cần truy cập vào tình trạng hiện tại.

**Trí tuệ nhân tạo (AI) đang đến.** Sắp tới, không chỉ con người mới tham gia hợp tác. Các tác nhân AI cũng sẽ tham gia vào quá trình thiết kế. Họ cũng cần truy cập thời gian thực.

Khoảng cách giữa những gì cần thiết và những gì có sẵn đang ngày càng gia tăng mỗi năm.

## Những Yếu Tố Cần Thiết Cho Hợp Tác Thực Sự

### Đồng bộ hóa thời gian thực

Không phải "đồng bộ hóa khi bạn lưu". Không phải "đẩy khi bạn hoàn tất". Đồng bộ hóa liên tục, thời gian thực, nơi các thay đổi xuất hiện ngay khi chúng được thực hiện.

Ai đó di chuyển một bức tường ở Sydney? Kỹ sư ở Munich thấy điều đó xảy ra. Không có độ trễ. Không có sự hợp nhất. Không có xung đột.

### Kiến trúc ưu tiên chế độ ngoại tuyến

Thời gian thực không có nghĩa là luôn kết nối. Khu vực sản xuất có tín hiệu wifi không ổn định. Công trường xây dựng không có tín hiệu. Máy bay không có internet.

Các công cụ hợp tác thực sự hoạt động offline và đồng bộ hóa tự động khi kết nối mạng được khôi phục. Khả năng đầy đủ, mọi lúc mọi nơi. Việc hợp nhất được hệ thống xử lý, không phải con người.

### Giải quyết xung đột hiệu quả

CRDT—Loại dữ liệu được sao chép không xung đột. Khoa học máy tính cuối cùng cũng đã bắt kịp nhu cầu.

Công nghệ Automerge và các công nghệ tương tự cho phép: nhiều người chỉnh sửa đồng thời, giải quyết xung đột tự động, không mất dữ liệu, không cần hợp nhất thủ công.

Điều này không phải là lý thuyết. Nó đã sẵn sàng cho sản xuất. Chỉ là nó chưa có trong hệ thống CAD của bạn mà thôi.

### Tiêu chuẩn mở

Hợp tác giữa các tổ chức có nghĩa là hợp tác giữa các phần mềm. Các kiến trúc sư của bạn sử dụng một công cụ. Các kỹ sư của bạn sử dụng một công cụ khác. Nhà sản xuất của bạn sử dụng một công cụ thứ ba.

Các định dạng độc quyền là kẻ thù của sự hợp tác. Các tiêu chuẩn mở (STEP, IFC) là công cụ hỗ trợ hợp tác.

Mô hình của bạn không nên bị giới hạn trong hệ sinh thái của một nhà cung cấp duy nhất. Sự hợp tác của bạn không nên phụ thuộc vào việc mọi người đều phải mua cùng một phần mềm.

## Cơ hội

Các nhà phát triển phần mềm đã giải quyết vấn đề này từ nhiều năm trước. Git, GitHub, hợp tác thời gian thực trong Google Docs—các mô hình này đã tồn tại.

Nhưng các công cụ thiết kế 3D không theo kịp. Hình học phức tạp hơn. Di sản sâu sắc hơn. Các động lực khuyến khích sự phụ thuộc vào một nền tảng duy nhất hơn là tính tương thích.

Điều này đang thay đổi. Công nghệ đã sẵn sàng:
- **Automerge CRDT** cho hợp tác thời gian thực, không xung đột
- **NATS JetStream** cho hệ thống nhắn tin phân tán quy mô lớn
- **Tiêu chuẩn mở** cho khả năng tương tác
- **Kiến trúc web bản địa* cho truy cập phổ quát

Điều cần thiết là có người để xây dựng nó. Được thiết kế chuyên biệt. Tích hợp trí tuệ nhân tạo (AI) từ đầu. Ưu tiên hợp tác.

---

*Working with distributed teams on 3D design projects? [See how we're solving this →](/vi/platform)*
