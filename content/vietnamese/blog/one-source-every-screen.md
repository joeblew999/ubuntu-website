---
title: "Một Nguồn, Mọi Màn Hình"
meta_title: "Một Nguồn, Mọi Màn Hình | Ubuntu Software"
description: "Từ thực đơn nhà hàng đến quầy thông tin của chính phủ: Cách xuất bản từ một nguồn duy nhất thay đổi cách hiển thị vật lý cùng với web và in ấn."
date: 2024-11-27T05:00:00Z
image: "/images/blog/kiosk-system.svg"
categories: ["Publish", "Industry"]
author: "Gerard Webb"
tags: ["kiosk", "signage", "raspberry-pi", "single-source", "real-time"]
draft: false
---

Một nhà hàng cập nhật thực đơn của mình. Điều gì sẽ xảy ra tiếp theo?

Ai đó chỉnh sửa trang web. Một người khác cập nhật tệp PDF cho menu mang về. Một người thứ ba thay đổi hệ thống điểm bán hàng. Cuối cùng, ai đó nhớ đến các bảng menu kỹ thuật số và cập nhật chúng thủ công.

Bốn hệ thống. Bốn cơ hội để xảy ra lỗi. Bốn nơi mà giá có thể biến động.

**Điều này thật荒谬. Và nó ở khắp mọi nơi.**

## Vấn đề phân mảnh hiển thị

Bước vào bất kỳ tổ chức nào có thông tin công khai:

**Restaurants**: Menu boards, table QR codes, website, delivery apps, printed menus
**Government offices**: Queue displays, wayfinding signs, forms counter, website
**Healthcare**: Waiting room displays, check-in kiosks, patient portals, printed materials
**Retail**: Price displays, promotional signage, inventory status, e-commerce

Mỗi màn hình thường chạy hệ thống riêng của mình. Mỗi màn hình yêu cầu cập nhật riêng biệt. Mỗi màn hình có thể hiển thị thông tin khác nhau.

**Thông tin bị lệch. Khách hàng bị nhầm lẫn. Nhân viên lãng phí thời gian vào việc đồng bộ hóa thủ công.**

## Giải pháp nguồn duy nhất

Nếu mọi màn hình—cả vật lý và kỹ thuật số—đều lấy dữ liệu từ cùng một nguồn thì sao?

```
Your Database (single source of truth)
         ↓
    Publish Engine
         ↓
┌────────┼────────┐
↓        ↓        ↓
Website  PDF    Kiosk Display
         ↓        ↓
      Print    Menu Board
```

Thay đổi giá một lần. Mọi màn hình hiển thị sẽ được cập nhật.

Thêm một mục mới một lần. Nó sẽ xuất hiện ở mọi nơi.

Ghi chú sản phẩm hết hàng một lần. Khách hàng sẽ thấy ngay lập tức—trên trang web, bảng menu và menu QR code in sẵn.

## Kiosk thời gian thực, có khả năng hoạt động offline

Đây là phần thú vị.

Biển hiệu kỹ thuật số truyền thống yêu cầu kết nối liên tục. Mạng bị ngắt kết nối? Màn hình sẽ tắt hoặc hiển thị dữ liệu cũ.

Phương pháp của chúng tôi sử dụng các thiết bị Raspberry Pi chạy trên cùng một kiến trúc ưu tiên chế độ ngoại tuyến như tất cả các hệ thống khác:

- **Bản sao dữ liệu cục bộ**: Mỗi kiosk đều có bản sao đầy đủ của dữ liệu liên quan
- **Tự động đồng bộ hóa*: Các bản cập nhật được đồng bộ hóa tự động khi kết nối
- **Sự suy giảm mềm mại**: Mạng bị gián đoạn? Màn hình vẫn hoạt động bình thường với dữ liệu cục bộ
- **Hàng đợi thay đổi**: Các cập nhật được thực hiện tại kiosk sẽ đồng bộ trở lại khi kết nối được khôi phục

**Một bảng menu hoạt động ngay cả khi mất kết nối internet. Một kiosk đăng ký cho phép xếp hàng các đơn đăng ký offline.**

## Các trường hợp sử dụng có ý nghĩa

### Hệ thống nhà hàng

Một nguồn cấp dữ liệu:
- Bảng menu kỹ thuật số (giá cả cập nhật theo thời gian thực, các mặt hàng hết hàng được hiển thị màu xám)
- Thực đơn mã QR (nội dung giống nhau, tối ưu hóa cho thiết bị di động)
- Trang menu của trang web
- Menu PDF để in
- Tích hợp hiển thị trong nhà bếp

Thay đổi giá của một chiếc burger. Tất cả các thông tin sẽ được cập nhật đồng bộ. Bảng menu tại quầy, tệp PDF mà ai đó in ra cho dịch vụ ẩm thực, trang web—tất cả đều được đồng bộ hóa.

### Trung tâm Dịch vụ Công

Một nguồn cấp dữ liệu:
- Hiển thị quản lý hàng đợi ("Đang phục vụ B-47")
- Các kiosk danh bạ dịch vụ
- Các biểu mẫu có thể tải xuống (PDF) hoặc điền trực tuyến (trang web/kiosk)
- Biển chỉ dẫn
- Danh sách dịch vụ website

Cập nhật giờ làm việc. Tất cả các màn hình hiển thị đều phản ánh sự thay đổi. Hệ thống chỉ dẫn, trang web, tài liệu in ấn—tất cả đều chính xác.

### Phòng chờ y tế

Một nguồn cấp dữ liệu:
- Quầy làm thủ tục tự động (điền biểu mẫu, xếp hàng)
- Thời gian chờ hiển thị
- Màn hình thông tin bệnh nhân
- Phiếu đăng ký in sẵn (các trường thông tin giống nhau, bố cục được căn chỉnh)

Bệnh nhân điền vào biểu mẫu trên kiosk → dữ liệu được truyền trực tiếp vào hệ thống quản lý hồ sơ y tế điện tử (EMR) của bạn → không cần nhập liệu thủ công.

### Môi trường bán lẻ

Một nguồn cấp dữ liệu:
- Hiển thị giá (nhãn giá điện tử)
- Biển quảng cáo
- Màn hình trạng thái kho hàng
- Danh sách sản phẩm trên trang web
- Catalog in ấn

Đánh dấu mặt hàng là "giảm giá" trong hệ thống kho. Giá hiển thị được cập nhật. Trang web được cập nhật. Quảng cáo in hàng tuần hiển thị giá giảm giá.

## Câu chuyện về phần cứng

Tại sao lại là Raspberry Pi?

**Cost**: Under $100 per display node, including the Pi and basic enclosure
**Reliability**: No moving parts, runs Linux, runs for years
**Flexibility**: HDMI output drives any display size
**Offline capable**: Local compute and storage for true edge operation
**Updatable**: Remote management for software updates

Một nhà hàng có thể triển khai bảng menu với chi phí chỉ bằng một phần nhỏ so với chi phí của bảng hiệu kỹ thuật số truyền thống. Một cơ quan chính phủ có thể lắp đặt kiosk mà không cần hợp đồng bảng hiệu doanh nghiệp.

## Kiến trúc kỹ thuật

Mỗi kiosk hoạt động:

| Component | Purpose |
|-----------|---------|
| **Local database** | SQLite replica of relevant data |
| **Sync engine** | Automerge CRDT for conflict-free updates |
| **Display renderer** | Web technologies (HTML/CSS/JS) for flexible layouts |
| **Input handlers** | Touch, keyboard, barcode scanner, card reader |
| **Offline queue** | Submissions stored locally until sync |

Cùng một mã nguồn chạy các biểu mẫu web của bạn cũng chạy trên kiosk. Cùng một quy trình xác thực. Cùng một bố cục trường. Cùng một đích dữ liệu.

## Không chỉ màn hình: Thiết bị nhập liệu cũng vậy

Kiosk không chỉ là thiết bị đầu ra. Chúng cũng là điểm nhập liệu.

**Form completion**: Customers complete forms on kiosk touchscreen
**Document scanning**: Kiosk scans paper documents, OCR extracts data
**Payment processing**: Integrated card readers
**ID verification**: Camera for document scanning
**Signature capture**: Touch signature pads

Tất cả dữ liệu được thu thập đều được xử lý qua cùng một hệ thống như các biểu mẫu trực tuyến.

Mẫu đơn của chính phủ được nộp qua:
- Biểu mẫu trang web → cơ sở dữ liệu
- Tệp PDF được in, điền thông tin, quét tại kiosk → cơ sở dữ liệu
- Biểu mẫu màn hình cảm ứng kiosk → cơ sở dữ liệu

**Cùng dữ liệu. Cùng đích đến. Phương pháp nhập liệu khác nhau.**

## Lợi thế về sự đồng bộ

Đây là ý nghĩa thực sự của từ "aligned":

Khách hàng nhìn thấy một biểu mẫu trên trang web của bạn. Cùng một khách hàng đến văn phòng của bạn và nhìn thấy cùng một biểu mẫu trên một kiosk. Các trường thông tin, bố cục và câu hỏi đều giống nhau.

Họ bắt đầu điền thông tin trên kiosk, hết thời gian, quét mã QR để tiếp tục trên điện thoại. Cùng một biểu mẫu, cùng tiến độ (được đồng bộ hóa qua phiên làm việc của họ).

Họ hoàn thành tại nhà, nhấn nút "Gửi". Dữ liệu được truyền trực tiếp vào cơ sở dữ liệu của bạn. Không cần nhập lại, không có lỗi nhập liệu, không có tên trường không khớp.

**Đây là những gì mà xuất bản từ một nguồn duy nhất mang lại khi được mở rộng sang các màn hình vật lý.**

## Trường hợp kinh doanh

**Trước khi có màn hình hiển thị từ một nguồn duy nhất:**
- Thời gian của nhân viên dành cho việc cập nhật nhiều hệ thống
- Lỗi do thông tin không nhất quán
- Sự nhầm lẫn của khách hàng ("biển hiệu ghi một giá, trang web lại ghi giá khác")
- Chi phí cao của các hệ thống biển báo kỹ thuật số truyền thống
- Các màn hình bị lỗi khi kết nối internet bị gián đoạn

**Sau khi hiển thị từ một nguồn duy nhất:**
- Cập nhật một lần, áp dụng cho mọi nơi
- Đảm bảo tính nhất quán trên tất cả các điểm tiếp xúc
- Giảm số lượng khiếu nại của khách hàng
- Phần cứng thông dụng với giá cả thông dụng
- Màn hình hoạt động offline

**Toán học rất đơn giản: ít lỗi hơn, tiết kiệm thời gian của nhân viên, trải nghiệm khách hàng tốt hơn.**

## Bắt đầu

Chúng tôi đang phát triển điều này ngay bây giờ. Cùng một công cụ xuất bản (Publish engine) tạo ra các trang web, tệp PDF và biểu mẫu cũng sẽ tạo ra các màn hình kiosk.

Cùng một nguồn Markdown. Cùng một DSL cho các trường. Cùng một kết nối cơ sở dữ liệu.

Màn hình chỉ là một định dạng đầu ra khác.

---

*Interested in early access to kiosk capabilities? [Contact us →](/contact)*

---

*Part of our Publish platform. [Learn more about single-source publishing →](/platform/publish)*
