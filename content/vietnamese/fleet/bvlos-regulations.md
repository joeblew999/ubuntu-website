---
title: "Quy định BVLOS"
meta_title: "Quy định về máy bay không người lái bay ngoài tầm nhìn trực quan | Ubuntu Software"
description: "Quy định toàn cầu về máy bay không người lái BVLOS (Beyond Visual Line of Sight – bay ngoài tầm nhìn trực quan) theo từng quốc gia. Tìm hiểu các yêu cầu pháp lý đối với hoạt động vận hành đội bay thương mại."
image: "/images/robotics.svg"
draft: false
---

## Vì sao BVLOS quan trọng đối với vận hành đội bay

Các quy định về bay ngoài tầm nhìn trực quan (Beyond Visual Line of Sight – BVLOS) quyết định một cách căn bản những gì các đội drone thương mại có thể thực hiện. Hiểu rõ bối cảnh pháp lý là điều thiết yếu cho việc thiết kế nền tảng quản lý đội bay và chiến lược thị trường.

| Thách thức | Tác động đến nền tảng quản lý đội bay |
|-----------|-------------------------|
| **Giới hạn vận hành** | Các hạn chế VLOS (phạm vi 400-500m) làm giảm tính hữu dụng của đội bay. BVLOS cho phép kiểm tra hạ tầng, giao hàng và nông nghiệp ở quy mô lớn. |
| **Tích hợp tuân thủ** | Phần mềm quản lý đội bay phải thực thi các ranh giới pháp lý—các yêu cầu về Remote ID, phát hiện và tránh va chạm, cùng giấy phép vận hành định hướng cho các tính năng đo xa (telemetry) và an toàn. |
| **Thời điểm thị trường** | Các quốc gia đang ở những giai đoạn khác nhau trong việc cho phép BVLOS. Hiểu được nơi nào đang mở cửa về mặt pháp lý giúp ưu tiên các thị trường mục tiêu. |
| **Mức độ sẵn sàng của đối tác** | Các nhà vận hành đội bay cần giấy phép BVLOS cho các nhiệm vụ thực tế. Theo dõi ai đã có giấy phép giúp xác định khách hàng và đối tác tiềm năng. |

---

## Xu hướng toàn cầu

Trên tất cả các cơ quan hàng không lớn, bốn xu hướng đang nổi lên:

| Xu hướng | Mô tả |
|-------|-------------|
| **Áp dụng SORA** | FAA, EASA, CASA và Transport Canada sử dụng phương pháp Đánh giá Rủi ro Hoạt động Cụ thể (Specific Operations Risk Assessment) [[1]](#ref-1) |
| **Bắt buộc Remote ID** | Đang trở thành yêu cầu bắt buộc ở tất cả các khu vực pháp lý |
| **Miễn trừ → Quy tắc** | Chuyển từ các miễn trừ theo từng trường hợp sang các quy tắc tiêu chuẩn hóa dựa trên hiệu suất |
| **Tích hợp UTM** | Các quốc gia phát triển hệ thống Quản lý Giao thông UAS (UTM) cho BVLOS có khả năng mở rộng |

---

## Châu Mỹ

### Hoa Kỳ (FAA)

**Tình trạng:** Dựa trên miễn trừ, đang chuyển sang dựa trên quy tắc (Part 108) [[2]](#ref-2)

**Mốc thời gian** [[3]](#ref-3):

- 190 miễn trừ BVLOS đã được cấp tính đến tháng 10/2024 (tăng từ 6 vào năm 2020)
- Dự thảo quy tắc Part 108 dự kiến vào tháng 1/2026
- Đạo luật Tái ủy quyền FAA 2024 đã gia hạn chương trình BEYOND đến năm 2029 [[4]](#ref-4)

**Các nhà vận hành có giấy phép BVLOS:**

| Nhà vận hành | Trường hợp sử dụng | Ngày cấp phép | Tham chiếu |
|----------|----------|---------------|-----|
| DroneUp | Giao hàng y tế | Tháng 1/2024 | [[5]](#ref-5) |
| Amazon Prime Air | Giao hàng tiêu dùng | Tháng 5/2024 | [[6]](#ref-6) |
| Asylon | An ninh/giám sát | Tháng 4/2024 | [[7]](#ref-7) |
| Nhiều nhà vận hành (khu vực DFW) | Vùng trời đa nhà vận hành đầu tiên | Tháng 7/2024 | [[8]](#ref-8) |

---

### Canada (Transport Canada)

**Tình trạng:** Yêu cầu SFOC, đang chuyển sang BVLOS thường lệ [[9]](#ref-9)

**Khung pháp lý:**

- Tháng 4/2025: BVLOS thường lệ được phép mà không cần SFOC trong điều kiện rủi ro thấp (≤150kg, vùng trời không kiểm soát, dân cư thưa thớt)
- Yêu cầu Chứng chỉ Phi công Nâng cao (Advanced Pilot Certificate)
- Khung UTM quốc gia thông qua Nav Canada vào năm 2030

**Nhà vận hành:** Chưa xác định

---

## Châu Âu

### Liên minh Châu Âu (EASA)

**Tình trạng:** Hạng mục Specific với các lộ trình SORA/PDRA [[10]](#ref-10)

**Khung pháp lý** [[1]](#ref-1):

| Khía cạnh | Chi tiết |
|--------|---------|
| Kịch bản tiêu chuẩn | STS-02: BVLOS với quan sát viên vùng trời |
| Remote ID | Bắt buộc từ tháng 1/2024 [[11]](#ref-11) |
| Thời gian cấp phép | 4-9 tháng (y tế), 6-12 tháng (tiêu dùng) [[12]](#ref-12) |
| Lệ phí | €2.000-€15.000 |

**Nhà vận hành:** Nhiều nhà vận hành thông qua các sandbox quốc gia và chương trình cấp phép nhanh của các quốc gia thành viên.

---

### Vương quốc Anh (CAA)

**Tình trạng:** Giấy phép Vận hành (Operational Authorization – OA) thông qua UK SORA

**Khung pháp lý:**

- BVLOS yêu cầu Giấy phép Vận hành của CAA
- Quy trình UK SORA (từ tháng 4/2025)
- Thường yêu cầu công nghệ phát hiện và tránh va chạm (ADS-B, FLARM, radar mặt đất)

**Nhà vận hành:**

| Nhà vận hành | Thành tựu | Tham chiếu |
|----------|-------------|-----|
| sees.ai | BVLOS thường lệ đầu tiên tại Anh (2021), kiểm tra cho National Grid (2023) | [[13]](#ref-13) |

---

### Ireland (IAA)

**Tình trạng:** Lộ trình PDRA G-03

**Khung pháp lý:** Mở đường cho việc cấp phép EU được đơn giản hóa thông qua EASA.

**Nhà vận hành:**

| Nhà vận hành | Thành tựu | Tham chiếu |
|----------|-------------|-----|
| sees.ai | Giấy phép BVLOS đầu tiên tại Ireland (2025) | [[13]](#ref-13) |

---

## Châu Á - Thái Bình Dương

### Úc (CASA)

**Tình trạng:** Đánh giá dựa trên SORA theo Part 101 [[14]](#ref-14)

**Khung pháp lý** [[15]](#ref-15):

- Áp dụng SORA 2.5 với mô hình hóa rủi ro định lượng
- Có sẵn khung Tự đánh giá BVLOS Diện rộng (Broad Area BVLOS Self-Assessment)
- Quy trình được đơn giản hóa thu hút các công ty toàn cầu [[16]](#ref-16)

**Nhà vận hành:**

| Nhà vận hành | Thành tựu | Tham chiếu |
|----------|-------------|-----|
| Sphere Drones | Giấy phép BVLOS Diện rộng (tháng 5/2024) | [[17]](#ref-17) |

---

### Nhật Bản (MLIT/JCAB)

**Tình trạng:** Cho phép BVLOS Cấp độ 4 trên khu vực đông dân cư [[18]](#ref-18)

**Khung pháp lý** [[19]](#ref-19):

| Cấp độ | Mô tả | Giấy phép yêu cầu |
|-------|-------------|------------------|
| Cấp độ 1-2 | Hoạt động VLOS | Hạng Hai |
| Cấp độ 3 | BVLOS trên khu vực không có người ở | Hạng Hai |
| Cấp độ 4 | BVLOS trên khu vực đô thị/đông dân cư | Hạng Nhất |

- Giấy phép có hiệu lực 3 năm, yêu cầu trường đào tạo drone đã đăng ký + bài kiểm tra bay
- 422.879 UAV đã được đăng ký tính đến tháng 10/2024

**Nhà vận hành:**

- Chỉ có 1 loại UAV đạt chứng nhận Hạng Nhất (Cấp độ 4) tính đến tháng 10/2024
- 5 loại UAV có chứng nhận Hạng Hai

---

### Singapore (CAAS)

**Tình trạng:** BVLOS bị cấm nếu không có Giấy phép Nhà vận hành + Giấy phép Hoạt động Loại 1 [[20]](#ref-20)

**Khung pháp lý** [[21]](#ref-21):

- Yêu cầu cấp phép nghiêm ngặt cho mọi hoạt động BVLOS
- Drone >25kg hoặc BVLOS yêu cầu Giấy phép Nhà vận hành đầy đủ + giấy phép Loại 1
- B-RID (Broadcast Remote ID) bắt buộc từ tháng 12/2025 đối với drone >250g
- B-RID tạo nền tảng cho các hoạt động BVLOS/tự động hóa quy mô lớn trong tương lai
- Có sẵn các chương trình sandbox pháp lý để cấp phép nhanh

**Nhà vận hành:** Chưa xác định (đang thử nghiệm sandbox)

---

### Trung Quốc (CAAC)

**Tình trạng:** Yêu cầu giấy phép đặc biệt cho BVLOS [[22]](#ref-22)

**Khung pháp lý** [[23]](#ref-23):

- Quy định Tạm thời về Quản lý Chuyến bay Phương tiện Bay Không người lái (tháng 1/2024)
- BVLOS yêu cầu Chứng chỉ Hoạt động Nâng cao + Chứng chỉ Hoạt động Bay Đặc biệt (SFOC)
- Giới hạn độ cao 120m trên toàn quốc; cao hơn phải xin phép
- Quy trình nộp đơn kỹ thuật số được đơn giản hóa cho phép nhiều chuyến bay BVLOS hơn [[16]](#ref-16)

**Nhà vận hành:** Chưa xác định (các nhà vận hành trong nước theo khung CAAC)

---

### Ấn Độ (DGCA)

**Tình trạng:** Chưa có khung BVLOS chính thức; chỉ ở giai đoạn thử nghiệm [[24]](#ref-24)

**Khung pháp lý** [[25]](#ref-25):

- Hoạt động theo Quy tắc Drone 2021 với các miễn trừ theo dự án
- Tháng 4/2025: Ra mắt hạng mục giấy phép phi công BVLOS riêng cho khu vực đô thị/bán đô thị
- DGCA "đang ở giai đoạn cuối hoàn thiện các quy tắc BVLOS" [[26]](#ref-26)
- Các hành lang thử nghiệm đang hoạt động (ví dụ: "Medicine from the Sky" tại Telangana)

**Nhà vận hành:**

| Nhà vận hành | Thành tựu | Tham chiếu |
|----------|-------------|-----|
| Skye Air Mobility | Chuyến giao thuốc BVLOS dài nhất thế giới (104km, Baruipur-Medinipur) | [[27]](#ref-27) |
| Flipkart Health | Đối tác trong thử nghiệm giao hàng 104km | — |

---

## Yêu cầu tuân thủ

Trên tất cả các khu vực pháp lý, các quy định BVLOS hội tụ về những yêu cầu kỹ thuật chung. Phần này tóm lược những gì nền tảng quản lý đội bay của chúng tôi phải triển khai.

### 1. Remote ID (Bắt buộc ở mọi nơi)

Tất cả các cơ quan quản lý lớn đều yêu cầu drone phát sóng dữ liệu nhận dạng và vị trí.

| Khu vực pháp lý | Yêu cầu | Thời hạn | Tham chiếu |
|--------------|-------------|----------|-----|
| Hoa Kỳ (FAA) | Standard Remote ID hoặc mô-đun phát sóng | Tháng 9/2023 (đã có hiệu lực) | [[2]](#ref-2) |
| EU (EASA) | Remote ID cho hạng mục Specific + Class 1+ | Tháng 1/2024 (đã có hiệu lực) | [[11]](#ref-11) |
| Singapore | B-RID cho drone >250g | Tháng 12/2025 | [[21]](#ref-21) |
| Úc | Remote ID đang được xem xét | Chưa xác định | [[14]](#ref-14) |

**Yêu cầu đối với nền tảng:**
- Tích hợp phần cứng mô-đun phát sóng (tương thích OpenDroneID)
- Khả năng Network Remote ID qua NATS để theo dõi toàn đội bay
- Ghi nhật ký tất cả các lần truyền Remote ID vào JetStream để kiểm toán

---

### 2. Phát hiện và tránh va chạm (DAA)

Hoạt động BVLOS yêu cầu khả năng phát hiện và tránh các tàu bay và chướng ngại vật khác.

| Khu vực pháp lý | Yêu cầu | Ghi chú |
|--------------|-------------|-------|
| Hoa Kỳ (FAA) | Tiêu chuẩn "Well clear"; yêu cầu hệ thống DAA cho Part 108 | Dựa trên hiệu suất, không mang tính quy định cứng |
| Anh (CAA) | Thường yêu cầu công nghệ DAA (ADS-B, FLARM, radar mặt đất) | Theo từng trường hợp trong quy trình OA |
| Nhật Bản | Bắt buộc cho hoạt động đô thị Cấp độ 4 | [[18]](#ref-18) |

**Yêu cầu đối với nền tảng:**
- Tích hợp bộ thu ADS-B In (nhận biết giao thông)
- Hợp nhất cảm biến trên Jetson (tùy chọn camera, radar, lidar)
- Thao tác tránh va chạm tự động hoặc cảnh báo cho người vận hành
- Ghi nhật ký tất cả các sự kiện DAA vào JetStream

---

### 3. Hàng rào địa lý (Geofencing)

Các ranh giới được thực thi bằng phần cứng nhằm ngăn chặn bay vào khu vực hạn chế.

| Khu vực pháp lý | Yêu cầu | Ghi chú |
|--------------|-------------|-------|
| Hoa Kỳ (FAA) | Geofence là một phần của giấy phép vận hành | Được xác định trong miễn trừ/Part 108 |
| EU (EASA) | Yêu cầu nhận biết địa lý (geo-awareness) trong hạng mục Specific | Theo đánh giá rủi ro mặt đất của SORA |
| Tất cả | Giới hạn độ cao (thường là 120m/400ft AGL) | Nền tảng của mọi khung pháp lý |

**Yêu cầu đối với nền tảng:**
- Thực thi geofence trên bộ điều khiển bay (tham số PX4 `GF_*`)
- Cập nhật geofence động qua Gateway (các hạn chế tạm thời)
- Thực thi trần độ cao
- Ghi nhật ký vi phạm và cảnh báo cho đội bay

---

### 4. Liên kết Chỉ huy & Điều khiển (C2)

Liên lạc đáng tin cậy, có dự phòng giữa người vận hành và tàu bay.

| Khu vực pháp lý | Yêu cầu | Ghi chú |
|--------------|-------------|-------|
| Hoa Kỳ (FAA) | Tiêu chuẩn hiệu suất liên kết C2 trong Part 108 | Độ trễ, tính sẵn sàng, bảo mật |
| EU (EASA) | Độ tin cậy liên kết C2 theo SORA | Yêu cầu phụ thuộc vào mức rủi ro |
| Tất cả | Hành vi an toàn khi mất C2 | Quay về điểm xuất phát, hạ cánh hoặc bay lượn chờ |

**Yêu cầu đối với nền tảng:**
- C2 chính qua LTE/5G hoặc radio chuyên dụng
- Đường C2 dự phòng (ví dụ: vệ tinh cho tầm xa)
- Chế độ an toàn tự động khi mất liên kết (RTL được cấu hình trong PX4)
- Đo xa trạng thái liên kết C2 gửi về đội bay
- Ghi nhật ký chất lượng liên kết để kiểm toán

---

### 5. Hồ sơ cấp phép vận hành

Các cơ quan quản lý yêu cầu đánh giá rủi ro và quy trình vận hành được lập thành văn bản.

| Khung | Được sử dụng bởi | Tài liệu chính |
|-----------|---------|---------------|
| SORA | EASA, UK CAA, CASA, Transport Canada | ConOps, Rủi ro Mặt đất, Rủi ro Trên không, biện pháp giảm thiểu OSO |
| Miễn trừ Part 107 | FAA (hiện tại) | Hồ sơ an toàn, giới hạn vận hành |
| Part 108 | FAA (sắp tới) | Tuân thủ dựa trên hiệu suất |

**Yêu cầu đối với nền tảng:**
- Xuất nhật ký bay theo định dạng được cơ quan quản lý chấp nhận
- Tạo dữ liệu hỗ trợ ConOps (khu vực bay, thời gian, tần suất)
- Lưu trữ hồ sơ chứng chỉ của nhà vận hành/phi công
- Theo dõi tình trạng bảo dưỡng và kiểm tra tàu bay

---

### 6. Cấp phép phi công/nhà vận hành

| Khu vực pháp lý | Giấy phép BVLOS | Ghi chú |
|--------------|---------------|-------|
| Hoa Kỳ | Part 107 + miễn trừ (hiện tại), Part 108 (tương lai) | [[2]](#ref-2) |
| EU | Chứng chỉ A2/STS hoặc giấy phép Specific | [[1]](#ref-1) |
| Anh | GVC + Giấy phép Vận hành | |
| Nhật Bản | Hạng Nhất (Cấp độ 4), Hạng Hai (Cấp độ 3) | [[19]](#ref-19) |
| Canada | Chứng chỉ Phi công Nâng cao + SFOC | [[9]](#ref-9) |
| Ấn Độ | Giấy phép phi công BVLOS (từ tháng 4/2025) | [[25]](#ref-25) |

**Yêu cầu đối với nền tảng:**
- Xác thực phi công trước khi cấp phép bay
- Xác minh giấy phép/chứng chỉ và theo dõi ngày hết hạn
- Ghi nhật ký giờ bay theo từng phi công
- Phân quyền theo vai trò (phi công, quan sát viên, quản lý đội bay)

---

### 7. Nhật ký bay & Dấu vết kiểm toán

Tất cả các khu vực pháp lý đều yêu cầu hồ sơ bay toàn diện.

**Yêu cầu đối với nền tảng:**
- Ghi lại đầy đủ dữ liệu đo xa (vị trí, tư thế, pin, cảm biến)
- Ghi nhật ký sự kiện có dấu thời gian (kích hoạt, cất cánh, điểm tham chiếu, hạ cánh, ngắt kích hoạt)
- Lưu trữ bất biến (JetStream với chính sách lưu giữ)
- Khả năng xuất dữ liệu để nộp cho cơ quan quản lý
- Đánh dấu sự cố và hỗ trợ điều tra

---

### 8. Bảo hiểm & Trách nhiệm pháp lý

| Khu vực pháp lý | Yêu cầu |
|--------------|-------------|
| EU | Bắt buộc bảo hiểm trách nhiệm đối với bên thứ ba |
| Hoa Kỳ | Bảo hiểm thường được yêu cầu theo điều kiện miễn trừ |
| Tất cả | Thường yêu cầu bằng chứng bảo hiểm để được cấp phép |

**Yêu cầu đối với nền tảng:**
- Theo dõi hợp đồng bảo hiểm theo từng tàu bay
- Xác minh phạm vi bảo hiểm trước khi cấp phép bay
- Lập hồ sơ sự cố phục vụ yêu cầu bồi thường

---

## Tích hợp bên ngoài

Tuân thủ BVLOS yêu cầu trao đổi dữ liệu thời gian thực với các hệ thống bên ngoài. Đây là những tích hợp mà nền tảng của chúng tôi cần hỗ trợ.

### Vùng trời & Giao thông

| Hệ thống | Mục đích | Luồng dữ liệu | Ưu tiên |
|--------|---------|-----------|----------|
| **ADS-B In** | Nhận biết giao thông từ tàu bay có người lái | Nhận → Hệ thống DAA | Cao |
| **FLARM** | Nhận biết giao thông từ tàu lượn/hàng không chung | Nhận → Hệ thống DAA | Trung bình |
| **UTM/USS** | Quản lý giao thông UAS, cấp phép bay | Hai chiều | Cao |
| **LAANC** (Hoa Kỳ) | Cấp phép vùng trời tự động | Yêu cầu/Phản hồi | Cao (Hoa Kỳ) |
| **AIM/NOTAM** | Hạn chế bay tạm thời, thay đổi vùng trời | Nhận → Geofence | Cao |

**Nhà cung cấp dịch vụ UTM (USS):**
- Hoa Kỳ: Airspace Link, Aloft, DroneUp, FlyFreely (xem [danh sách LAANC của FAA](https://www.faa.gov/uas/getting_started/laanc))
- EU: Khác nhau theo từng quốc gia thành viên
- Úc: Wing, Airservices Australia

*Lưu ý: AirMap đã ngừng dịch vụ LAANC vào tháng 6/2023.*

---

### Mạng Remote ID

| Hệ thống | Mục đích | Giao thức | Ưu tiên |
|--------|---------|----------|----------|
| **OpenDroneID** | Tiêu chuẩn Broadcast Remote ID | Bluetooth 4/5, WiFi Beacon | Bắt buộc |
| **Network Remote ID** | Nhận dạng toàn đội bay qua internet | API của USS | Bắt buộc (Hoa Kỳ) |
| **FAA Remote ID** | Tuân thủ Remote ID tại Hoa Kỳ | Tích hợp USS | Bắt buộc (Hoa Kỳ) |

---

### Cơ sở dữ liệu pháp lý

| Hệ thống | Mục đích | Luồng dữ liệu | Ưu tiên |
|--------|---------|-----------|----------|
| **CSDL Giấy phép Phi công** | Xác minh chứng chỉ phi công | Truy vấn → Cấp phép | Cao |
| **Sổ đăng ký Tàu bay** | Xác minh tình trạng đăng ký | Truy vấn → Trước chuyến bay | Cao |
| **Xác minh Bảo hiểm** | Xác nhận hiệu lực bảo hiểm | Truy vấn → Trước chuyến bay | Trung bình |
| **Hồ sơ Bảo dưỡng** | Theo dõi khả năng bay | Hai chiều | Trung bình |

---

### Thời tiết & Môi trường

| Hệ thống | Mục đích | Luồng dữ liệu | Ưu tiên |
|--------|---------|-----------|----------|
| **Thời tiết hàng không (METAR/TAF)** | Gió, tầm nhìn, trần mây | Nhận → Lập kế hoạch bay | Cao |
| **Radar/Lượng mưa** | Các mối nguy thời tiết thời gian thực | Nhận → Cảnh báo trong chuyến bay | Trung bình |
| **Mặt trời/Địa từ** | Dự báo độ tin cậy GPS | Nhận → Lập kế hoạch bay | Thấp |

---

### Dịch vụ khẩn cấp

| Hệ thống | Mục đích | Luồng dữ liệu | Ưu tiên |
|--------|---------|-----------|----------|
| **ATC (nếu được yêu cầu)** | Phối hợp vùng trời có kiểm soát | Liên kết thoại/dữ liệu | Tùy tình huống |
| **Dịch vụ khẩn cấp (911/112)** | Báo cáo sự cố | Cảnh báo gửi đi | Bắt buộc |
| **Nộp NOTAM** | Công bố hoạt động bay | Gửi đi → AIM | Tùy tình huống |

---

### Kiến trúc tích hợp

```
┌─────────────────────────────────────────────────────────────────────┐
│                         FLEET PLATFORM                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │
│   │   Gateway   │    │  JetStream  │    │   Fleet     │             │
│   │  (Vehicle)  │    │  (Logging)  │    │  Manager    │             │
│   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘             │
│          │                  │                  │                     │
└──────────┼──────────────────┼──────────────────┼─────────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│                      INTEGRATION LAYER                               │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │   UTM    │  │  Remote  │  │ Weather  │  │ License  │             │
│  │   USS    │  │    ID    │  │   API    │  │   DB     │             │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘             │
│       │             │             │             │                    │
└───────┼─────────────┼─────────────┼─────────────┼────────────────────┘
        │             │             │             │
        ▼             ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │Airspace │   │   FAA   │   │  NOAA   │   │  CAA    │
   │  Link   │   │ DroneZone│  │ Aviation│   │ Registry│
   │ Aloft   │   │ NASA DIP│   │ Weather │   │         │
   └─────────┘   └─────────┘   └─────────┘   └─────────┘
```

---

### Lộ trình tích hợp

| Tích hợp | Tình trạng | Ghi chú |
|-------------|--------|-------|
| Phát sóng OpenDroneID | ✅ Đã hỗ trợ | Cần mô-đun phần cứng |
| Bộ thu ADS-B In | 🔧 Lộ trình | uAvionix pingRX hoặc tương tự |
| UTM/USS (Airspace Link) | 🔧 Lộ trình | AirHub API, giá theo cấp |
| UTM/USS (Aloft) | 🔧 Lộ trình | Tích hợp REST API |
| LAANC | 🔧 Lộ trình | Qua đối tác USS |
| Thời tiết hàng không | 🔧 Lộ trình | aviationweather.gov (mở, 100 yêu cầu/phút) |
| Nguồn cấp NOTAM | 🔧 Lộ trình | NASA DIP hoặc FAA SWIM |
| Xác minh giấy phép phi công | 🔧 Lộ trình | API theo từng khu vực pháp lý |
| Xác minh bảo hiểm | 🔧 Lộ trình | API của nhà cung cấp |

---

## Ma trận tuân thủ

| Yêu cầu | Thành phần nền tảng | Tình trạng |
|-------------|-------------------|--------|
| Phát sóng Remote ID | Mô-đun OpenDroneID + Gateway | ✅ Đã hỗ trợ |
| Ghi nhật ký Remote ID | Các luồng JetStream | ✅ Đã hỗ trợ |
| Phát hiện và tránh va chạm | Hợp nhất cảm biến Jetson | 🔧 Lộ trình |
| ADS-B In | Tích hợp bộ thu giao thông | 🔧 Lộ trình |
| Geofencing (tĩnh) | Tham số PX4 `GF_*` | ✅ Đã hỗ trợ |
| Geofencing (động) | Cập nhật chính sách Gateway | ✅ Đã hỗ trợ |
| Liên kết C2 chính | LTE/5G qua NATS | ✅ Đã hỗ trợ |
| Liên kết C2 dự phòng | Tích hợp vệ tinh | 🔧 Lộ trình |
| An toàn khi mất liên kết | Cấu hình failsafe PX4 | ✅ Đã hỗ trợ |
| Ghi nhật ký bay | Các luồng đo xa JetStream | ✅ Đã hỗ trợ |
| Xuất dữ liệu kiểm toán | Công cụ xuất nhật ký | ✅ Đã hỗ trợ |
| Xác thực phi công | Hệ thống cấp phép | ✅ Đã hỗ trợ |
| Theo dõi giấy phép | Quản lý đội bay | 🔧 Lộ trình |
| Theo dõi bảo hiểm | Quản lý đội bay | 🔧 Lộ trình |

Xem [Mô hình An toàn]({{< relref "/fleet/safety" >}}) để biết cách các thành phần này tích hợp với kiến trúc an toàn tổng thể.

---

## Tài liệu tham khảo

<span id="ref-1"></span>**[1]** EASA. "Drones (UAS) - Frequently Asked Questions." European Union Aviation Safety Agency. https://www.easa.europa.eu/en/the-agency/faqs/drones-uas

<span id="ref-2"></span>**[2]** Drone Pilot Ground School. "Part 108: What We Know So Far about the FAA's BVLOS Rule." https://www.dronepilotgroundschool.com/part-108/

<span id="ref-3"></span>**[3]** DroneLife. "FAA Progress on BVLOS Rules and Advanced Air Mobility Integration." Tháng 8/2024. https://dronelife.com/2024/08/01/faa-progress-on-bvlos-rules-and-advanced-air-mobility-integration/

<span id="ref-4"></span>**[4]** FAA. "BEYOND Program." Federal Aviation Administration. https://www.faa.gov/uas/programs_partnerships/beyond

<span id="ref-5"></span>**[5]** DroneUp. "DroneUp Awarded Landmark FAA Approval for Beyond Visual Line of Sight (BVLOS) in the U.S." Tháng 1/2024. https://www.droneup.com/news/droneup-awarded-bvlos

<span id="ref-6"></span>**[6]** DroneDJ. "Amazon gets FAA BVLOS approval for Prime Air drone delivery." Tháng 5/2024. https://dronedj.com/2024/05/30/amazon-gets-faa-bvlos-approval-for-prime-air-drone-delivery/

<span id="ref-7"></span>**[7]** DroneDJ. "Asylon bags new FAA approval for BVLOS drone operations." Tháng 4/2024. https://dronedj.com/2024/04/26/asylon-drone-faa-bvlos-dronesentry/

<span id="ref-8"></span>**[8]** DroneLife. "FAA Grants Historic BVLOS Approval for Multiple Operators in DFW Area." Tháng 7/2024. https://dronelife.com/2024/07/31/faa-grants-historic-bvlos-approval-for-multiple-operators-in-dfw-area/

<span id="ref-9"></span>**[9]** MLT Aikins. "International drone regulations: Key updates on BVLOS, UTM and use cases." https://www.mltaikins.com/insights/international-drone-regulations-key-updates-on-bvlos-utm-and-use-cases/

<span id="ref-10"></span>**[10]** FlytBase. "Understanding EU Drone Regulations for Docked BVLOS Operations." https://www.flytbase.com/blog/eu-drone-regulations-for-bvlos-operations

<span id="ref-11"></span>**[11]** Elsight. "EASA Remote ID Requirements: Compliance for Drone Operators." https://www.elsight.com/blog/primer-on-easa-remote-id-regulations/

<span id="ref-12"></span>**[12]** Sparkco. "Autonomous Drone Delivery Regulatory Approval." https://sparkco.ai/blog/autonomous-drone-delivery-regulatory-approval

<span id="ref-13"></span>**[13]** sees.ai. "Why Us - First UK routine BVLOS permission." https://www.sees.ai/why-us/

<span id="ref-14"></span>**[14]** CASA. "BVLOS drone operations in regional Australia - Consultation." Civil Aviation Safety Authority. https://consultation.casa.gov.au/stakeholder-engagement-group/consultation.2023-10-05.5578154857/

<span id="ref-15"></span>**[15]** Measure Australia. "New CASA Drone Regulations Incoming!" https://www.measureaustralia.com.au/news/new-casa-drone-regulations-incoming

<span id="ref-16"></span>**[16]** Droneii. "Commercial Drone Regulation Report 2025." Drone Industry Insights. https://droneii.com/product/commercial-drone-regulation-report

<span id="ref-17"></span>**[17]** DroneLife. "Sphere Drones Secures BVLOS Area Approval to Advance Commercial Drone Use in Australia." Tháng 5/2024. https://dronelife.com/2024/05/08/sphere-drones-secures-bvlos-area-approval-to-advance-commercial-drone-use-in-australia/

<span id="ref-18"></span>**[18]** Unmanned Airspace. "Japan's revised drone laws to permit BVLOS flights over people come into effect." https://www.unmannedairspace.info/emerging-regulations/japanese-revised-drone-laws-to-permit-bvlos-flights-over-people-come-into-effect/

<span id="ref-19"></span>**[19]** Fly Eye. "Japan Drone Laws (2025)." https://www.flyeye.io/japan-drone-laws/

<span id="ref-20"></span>**[20]** Drone Laws. "Singapore Drone Laws 2025." https://drone-laws.com/drone-laws-in-singapore/

<span id="ref-21"></span>**[21]** Heron AirBridge. "B-RID in Singapore: Full Guide to the New Drone Regulations." https://heron-airbridge.com/singapore-brid-regulation/

<span id="ref-22"></span>**[22]** Fly Eye. "Advanced Drone Operations in China." https://www.flyeye.io/chinese-drone-regulations-advanced-operations/

<span id="ref-23"></span>**[23]** Fly Eye. "Chinese Drone Laws (2025)." https://www.flyeye.io/chinese-drone-laws/

<span id="ref-24"></span>**[24]** Commercial UAV News. "India Makes a Giant Leap Toward BVLOS Operations." https://www.commercialuavnews.com/drone-delivery/india-makes-a-giant-leap-toward-bvlos-operations

<span id="ref-25"></span>**[25]** Fly and Tech. "Drone Rules in India 2025 – Latest DGCA Guidelines Update." https://flyandtech.com/drone-rules-india-2025/

<span id="ref-26"></span>**[26]** Sigma Chambers. "Drone Law Brief | Vol 1: September & October 2025." https://www.sigmachambers.in/post/drone-law-watch-vol-1-september-october-2025

<span id="ref-27"></span>**[27]** ITLN. "How will BVLOS drone delivery ops impact the Indian logistics industry?" https://www.itln.in/cargo-drones/how-will-bvlos-drone-delivery-ops-impact-the-indian-logistics-industry-1353269

---

## Quay lại Tổng quan

[← Kiến trúc Đội Phương tiện Tự hành]({{< relref "/fleet" >}})
