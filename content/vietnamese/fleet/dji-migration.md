---
title: "Lộ trình chuyển đổi từ DJI"
meta_title: "Chuyển đổi từ DJI sang drone mã nguồn mở | Ubuntu Software"
description: "Chuyển đổi từ DJI sang drone mã nguồn mở trước lệnh cấm của Mỹ. Hướng dẫn chuyển đổi toàn diện dành cho các nhà khai thác thương mại, cơ quan và đội bay doanh nghiệp."
image: "/images/robotics.svg"
draft: false
---

{{< notice "danger" >}}
**Lệnh cấm hiện đã có hiệu lực.** Ngày 22 tháng 12 năm 2025, FCC đã đưa DJI vào Danh sách Covered List. Các mẫu DJI mới và linh kiện RF không còn được FCC phê duyệt để bán tại Mỹ. Các drone hiện có vẫn được phép bay hợp pháp, nhưng đồng hồ chuyển đổi đang đếm ngược. [Đọc phân tích đầy đủ →]({{< relref "/blog/dji-fcc-ban-what-happened" >}})
{{< /notice >}}

## Tại sao điều này quan trọng

Nếu bạn đang bay drone DJI cho hoạt động thương mại (kiểm tra, lập bản đồ, khảo sát, an toàn công cộng), bạn cần một kế hoạch chuyển đổi. Không phải vì drone hiện có của bạn sẽ ngừng bay, mà vì:

- **Không có cập nhật firmware** sau lệnh cấm
- **Không có phụ tùng thay thế** từ các kênh chính hãng
- **Không có hỗ trợ SDK** cho các tích hợp tùy chỉnh
- **Không có dịch vụ đám mây** để quản lý đội bay
- **Không thể mua mới** để mở rộng đội bay

Những nhà vận hành chờ đợi sẽ phải chuyển đổi khẩn cấp với giá cao. Những ai bắt đầu ngay bây giờ có thể chuyển đổi một cách có phương pháp.

---

## Tình hình

### Điều gì đã xảy ra

Ngày 22 tháng 12 năm 2025, một ngày trước hạn chót, FCC đã công bố việc đưa DJI vào Danh sách Covered List. Chính quyền đã bỏ qua cuộc kiểm toán an ninh theo yêu cầu của Mục 1709 trong NDAA, thay vào đó ban hành một "quyết định an ninh quốc gia".

DJI đã gửi thư vào tháng 3, tháng 6 và tháng 12 năm 2025 đề nghị các cơ quan kiểm tra sản phẩm của họ. Không cơ quan nào phản hồi. Không cuộc điều tra nào được bắt đầu. Không bằng chứng nào được thu thập.

Quốc hội đã cố tình thiết kế Mục 1709 với những "cửa sập": không cơ quan nào được giao thực hiện kiểm toán, và lệnh cấm sẽ tự động có hiệu lực nếu không có cuộc kiểm toán nào diễn ra. Kết quả đã được định đoạt từ trước.

### Điều này có nghĩa là gì

| Tác động | Đội bay DJI hiện có của bạn | Mua DJI mới |
|--------|------------------------|-------------------|
| **Bay** | Vẫn hợp pháp theo FAA Part 107 | Bị chặn |
| **Firmware** | Không có cập nhật | N/A |
| **Phụ tùng** | Hạn chế/thị trường xám | N/A |
| **SDK/API** | Không được hỗ trợ | N/A |
| **Dịch vụ đám mây** | Không chắc chắn | N/A |
| **Mẫu mới** | N/A | Bị chặn |

### Ai bị ảnh hưởng

- **Cơ quan chính phủ**: Đã bị hạn chế mua sắm DJI
- **An toàn công cộng**: Hơn 80% lực lượng cứu hỏa, cảnh sát và tìm kiếm cứu nạn (SAR) tại Mỹ sử dụng DJI
- **Nhà vận hành thương mại**: Kiểm tra, lập bản đồ, khảo sát, nông nghiệp
- **Đội bay doanh nghiệp**: Cơ sở hạ tầng, tiện ích, xây dựng
- **Nhà thầu quốc phòng**: Yêu cầu Blue UAS

---

## Lộ trình chuyển đổi

### Bạn nhận được gì

Chuyển sang mã nguồn mở không chỉ để tránh lệnh cấm. Bạn còn nhận được:

| DJI | Mã nguồn mở (PX4/ArduPilot) |
|-----|----------------------------|
| Hệ sinh thái khép kín | Độc lập với nhà cung cấp |
| Dữ liệu độc quyền | Dữ liệu của bạn, máy chủ của bạn |
| Lộ trình tính năng do họ kiểm soát | Cộng đồng + lộ trình của bạn |
| Giới hạn SDK | Truy cập toàn bộ mã nguồn |
| Bị khóa phần cứng | Nhiều nhà sản xuất |
| Phụ thuộc đám mây | Tự lưu trữ hoặc được quản lý |

### Nền tảng tương thích

Các flight stack mã nguồn mở hoạt động với nhiều loại phần cứng:

| Nền tảng | Flight Stack | Điểm mạnh |
|----------|-------------|-----------|
| **Holybro X500/S500** | PX4/ArduPilot | Đã được kiểm chứng, sẵn có, tài liệu đầy đủ |
| **Bản tự lắp ráp** | PX4/ArduPilot | Linh hoạt tối đa |
| **Freefly** | PX4 | Quay phim/khảo sát chuyên nghiệp |
| **Autel** | Độc quyền | Lắp ráp tại Mỹ, nhưng cũng đang bị xem xét theo NDAA |
| **Skydio** | Độc quyền | Sản xuất tại Mỹ, tập trung vào khả năng tự hành |

Đối với vận hành đội bay, chúng tôi khuyến nghị các nền tảng dựa trên PX4 vì:
- Mã nguồn mở với sự hậu thuẫn thương mại (Dronecode Foundation)
- Giao thức MAVLink là một tiêu chuẩn mở
- Hệ sinh thái lớn các linh kiện tương thích
- Giấy phép BSD cho phép sửa đổi độc quyền

### Lựa chọn phần mềm điều khiển mặt đất

| Phần mềm | Giao thức | Điểm mạnh |
|----------|----------|-----------|
| **QGroundControl** | MAVLink | Miễn phí, đa nền tảng, PX4/ArduPilot |
| **Mission Planner** | MAVLink | Tập trung vào ArduPilot, Windows |
| **UgCS** | MAVLink | Thương mại, lập kế hoạch nâng cao |
| **Nền tảng Đội bay của chúng tôi** | NATS + MAVLink | Đa phương tiện, doanh nghiệp |

---

## Các cấp độ chuyển đổi

### Cấp 1: Đội bay song song (6-12 tháng)

Phù hợp nhất cho: Các cơ quan đang hoạt động không thể chấp nhận rủi ro gián đoạn.

1. **Mua máy bay thử nghiệm**: Bắt đầu với 1-3 nền tảng mã nguồn mở
2. **Đào tạo phi công**: Bộ điều khiển khác, phần mềm khác
3. **Xác thực quy trình**: Lập bản đồ, kiểm tra, các trường hợp sử dụng của bạn
4. **Xây dựng SOP**: Danh sách kiểm tra mới, quy trình bảo trì
5. **Chuyển đổi dần dần**: Chuyển hoạt động khi sự tự tin tăng lên

### Cấp 2: Chuyển đổi có kế hoạch (3-6 tháng)

Phù hợp nhất cho: Các nhà vận hành có thể lên lịch các khoảng thời gian chuyển đổi.

1. **Đánh giá tài sản**: Tài sản DJI nào, quy trình nào
2. **Lựa chọn nền tảng**: Đối chiếu khả năng với nhu cầu nhiệm vụ
3. **Mua sắm**: Đặt hàng phần cứng, tính đến thời gian giao hàng
4. **Đào tạo chuyên sâu**: Đào tạo phi công rút gọn
5. **Chuyển giao**: Loại bỏ tài sản DJI một cách có hệ thống

### Cấp 3: Chuyển đổi khẩn cấp (1-3 tháng)

Phù hợp nhất cho: Các nhà vận hành đã chờ đợi quá lâu.

1. **Mua sắm ngay lập tức**: Giá cao, lựa chọn hạn chế
2. **Đào tạo cấp tốc**: Chuyên sâu, có thể do nhà thầu cung cấp
3. **Năng lực giảm sút**: Chấp nhận những khoảng trống tạm thời trong quy trình
4. **Lặp lại cải tiến**: Cải thiện hoạt động theo thời gian

**Đừng để rơi vào Cấp 3.**

---

## Vai trò của chúng tôi

### Những gì chúng tôi cung cấp

**Cơ sở hạ tầng đội bay, không phải phần cứng.** Chúng tôi không cạnh tranh với các cửa hàng drone. Chúng tôi cung cấp lớp phần mềm giúp các đội bay thực sự vận hành được.

| Bạn nhận được | Chúng tôi xử lý |
|---------|-----------|
| Cơ sở hạ tầng NATS được quản lý | Vận hành cụm, cập nhật |
| Phần mềm Vehicle Gateway | Cầu nối MAVLink-sang-NATS |
| Quản lý trạng thái đội bay | Bản sao số (digital twin), dữ liệu đo xa |
| Điều phối đa phương tiện | Phân cấp subject, phân quyền |
| Tích hợp doanh nghiệp | API, xuất dữ liệu |

### Những gì bạn cung cấp

| Bạn xử lý | Chúng tôi hỗ trợ |
|------------|-----------|
| Mua sắm phần cứng | Giới thiệu cửa hàng đối tác |
| Đào tạo phi công | Tài liệu, thực hành tốt nhất |
| Vận hành tại chỗ | Giám sát đội bay từ xa |
| Tuân thủ quy định | Lưu giữ dữ liệu, nhật ký kiểm toán |

### Riêng cho lập bản đồ

Nhiều nhà vận hành DJI cần giải pháp lập bản đồ. Kiến trúc của chúng tôi hỗ trợ:

- **Đường ống dữ liệu cảm biến thô**: Hình ảnh đến hệ thống xử lý phía sau
- **Dữ liệu đo xa thời gian thực**: Vị trí, độ cao, sự kiện camera
- **Điều phối đội bay**: Các mẫu khảo sát đa phương tiện
- **Bàn giao dữ liệu**: Tích hợp với ODM, Pix4D, DroneDeploy

Chúng tôi không phải là giải pháp lập bản đồ hoàn chỉnh. Chúng tôi là cơ sở hạ tầng đội bay kết nối drone của bạn với đường ống xử lý của bạn.

---

## Đối tác phần cứng

Các cửa hàng này có thể giúp bạn tìm nguồn phần cứng mã nguồn mở:

{{< partners-list status="active" type="manufacturer,reseller" >}}

Đang tìm cửa hàng địa phương? [Tìm đối tác gần bạn]({{< relref "/fleet/hardware" >}}).

---

## Bắt đầu

### Cuộc gọi đánh giá

Cuộc gọi miễn phí 30 phút để hiểu rõ hoạt động hiện tại của bạn:

- Quy mô và cơ cấu đội bay
- Các trường hợp sử dụng chính (lập bản đồ, kiểm tra, v.v.)
- Thời gian và mức độ cấp bách
- Giới hạn ngân sách

[Đặt lịch đánh giá →](/contact?subject=dji-migration)

### Chương trình thí điểm

Dành cho các cơ quan và doanh nghiệp có đội bay DJI lớn:

- Thử nghiệm khái niệm (proof-of-concept) với quy trình của bạn
- Vận hành song song cùng đội bay DJI hiện có
- So sánh và xác thực hiệu suất
- Quyết định tiếp tục hay dừng dựa trên dữ liệu thực tế

[Yêu cầu Chương trình thí điểm →](/contact?subject=dji-pilot)

---

## Tài nguyên

- [Hardware Stack →]({{< relref "/fleet/hardware" >}}): BOM tham chiếu cho drone PX4
- [Software Stack →]({{< relref "/fleet/software" >}}): Hệ điều hành và flight stack
- [Kiến trúc Đội bay →]({{< relref "/fleet" >}}): Quản lý đội bay dựa trên NATS
- [Chương trình Đối tác →]({{< relref "/partners" >}}): Dành cho các cửa hàng drone phục vụ những người đang rời bỏ DJI

---

## Câu hỏi thường gặp

**Drone DJI hiện có của tôi có ngừng hoạt động không?**

Không. Drone hiện tại của bạn sẽ tiếp tục bay theo FAA Part 107. Nhưng bạn sẽ không nhận được cập nhật firmware, phụ tùng có thể trở nên khan hiếm và dịch vụ đám mây có thể bị ngừng.

**Autel có phải là lựa chọn thay thế an toàn không?**

Autel cũng có tên trong yêu cầu kiểm toán của Mục 1709 NDAA. Họ phải đối mặt với cùng mốc thời gian như DJI.

**Skydio có phải là lựa chọn thay thế an toàn không?**

Skydio được sản xuất tại Mỹ và không chịu các hạn chế tương tự. Tuy nhiên, hệ sinh thái của họ cũng là độc quyền.

**Tại sao chọn mã nguồn mở thay vì một nền tảng độc quyền khác?**

Mã nguồn mở trao cho bạn quyền kiểm soát. Bạn không đánh đổi sự phụ thuộc vào nhà cung cấp này lấy sự phụ thuộc vào nhà cung cấp khác. Phần mềm là của bạn để sửa đổi, dữ liệu là của bạn để lưu giữ, và phần cứng có thể thay thế lẫn nhau.

**Quá trình chuyển đổi mất bao lâu?**

Tùy thuộc vào hoạt động của bạn. Hãy dự trù 3-6 tháng cho một quá trình chuyển đổi có phương pháp với gián đoạn tối thiểu. Chuyển đổi khẩn cấp có thể diễn ra trong vài tuần nhưng với chi phí và rủi ro cao hơn.

**Còn các quy trình lập bản đồ hiện có của chúng tôi thì sao?**

Hầu hết phần mềm xử lý (Pix4D, DroneDeploy, ODM) hoạt động với bất kỳ drone nào tạo ra hình ảnh được gắn thẻ địa lý. Việc chuyển đổi nằm ở nền tảng thu thập, không phải ở hệ thống xử lý phía sau.

---

## Thời gian đã hết

Lệnh cấm hiện đã có hiệu lực. Nếu bạn vẫn đang bay toàn bộ bằng DJI, bạn đang vận hành trong thời gian vay mượn, không có lộ trình mở rộng đội bay, hỗ trợ firmware không chắc chắn và chuỗi cung ứng ngày càng thắt chặt.

Các cửa hàng am hiểu mã nguồn mở sắp trở nên rất bận rộn.

[Liên hệ với chúng tôi →](/contact?subject=dji-migration)
