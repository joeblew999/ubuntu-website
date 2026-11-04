---
title: "ハードウェアスタック"
meta_title: "ドローンフリートのハードウェア仕様 | Ubuntu Software"
description: "1,000機のドローンフリート向け本番用ハードウェア仕様：Holybro X500機体、Pixhawk 6Xフライトコントローラー、Raspberry PiおよびJetsonコンパニオンコンピューター。"
image: "/images/robotics.svg"
draft: false
aliases:
  - "/fleet/hardware/"
---

{{< notice "info" >}}
**他の車両タイプをお探しですか？** このページではドローンのハードウェアについて説明します。地上車両（乗用車、トラック、AGV）については、[地上車両プラットフォーム]({{< relref "/fleet/platforms/ground" >}})をご覧ください。全体の概要については、[対応プラットフォーム]({{< relref "/fleet/platforms" >}})をご覧ください。
{{< /notice >}}

## 地元のショップで購入する

ドローンの組み立ては、慣れるまでに時間がかかります。**サポートしてくれるショップで購入しましょう**。ハードウェアはショップが取り揃え、フリートソフトウェアは私たちが提供します。

パートナーショップで購入すると、次のものが得られます：
- すぐに接続できるテスト済みのハードウェア
- ドローンに詳しいスタッフによる地域サポート
- 当社のマネージドフリートインフラへのアクセス

---

## 実店舗のあるショップ

組み立てについて直接サポートを受けたい場合は、店舗を訪れてください：

{{< featured-shops >}}

**ドローンショップの方ですか？** [パートナープログラムに参加]({{< relref "/partners/become-partner" >}})して、ハードウェア販売とあわせてフリートインフラを提供しましょう。

---

## 必要なもの

| キット | コンポーネント | 主な購入先 |
|-----|------------|-----------------|
| **完成ドローン** | X500フレーム、Pixhawk 6X、GPS、モーター、ESC、プロペラ | GetFPV, NewBeeDrone, Holybro |
| **コンパニオンコンピューター** | Raspberry Pi 4 + Jetson Orin Nano | Digi-Key, Mouser, SparkFun |
| **電源とバッテリー** | 4S LiPo 6000mAh、パワーモジュール | GetFPV, Tattu, Amazon |
| **無線機** | ExpressLRSレシーバー、SiKテレメトリー | GetFPV, HappyModel, Holybro |
| **ネットワーク** | 産業用イーサネットスイッチ、5Gモデム | PUSR, Quectel |

[全コンポーネントリストを見る ↓](#full-bill-of-materials)

---

## フリート規模の運用向けハードウェア

フリート運用には、性能と保守性を両立したハードウェアが欠かせません。このスタックのコンポーネントはすべて、本番環境での信頼性、部品の入手しやすさ、現場での整備のしやすさを基準に選定しています。

---

## 機体

### Holybro X500 V2 ARF

**Holybro X500 V2** は、フリート展開に最適なプラットフォームです：

| 仕様 | 値 |
|---------------|-------|
| **SKU** | 30125 |
| **ホイールベース** | 500mm |
| **フレーム重量** | 約410g（モーター/ESC込みで610g） |
| **最大離陸重量** | 約2kg |
| **飛行時間** | 15〜20分（標準ペイロード時） |
| **モーターマウント** | 16x16mm / 19x19mm |
| **付属品** | モーター（2216 KV920）、ESC（20A BLHeli S）、プロペラ（1045） |

**X500を選ぶ理由：**

- **実績のある設計** — 世界中で数千台が運用されており、コミュニティの知見が豊富
- **部品の入手しやすさ** — アーム、ランディングギア、各種ハードウェアを容易に入手可能
- **メンテナンスしやすい** — モジュール構造で、現場での整備が可能
- **ペイロードの柔軟性** — デュアルコンパニオンコンピューターとセンサーを搭載できる十分な積載能力

ARF（Almost Ready to Fly）キットにはフレーム、モーター、ESC、プロペラが含まれているため、大規模フリート全体で組み立て品質のばらつきを抑えられます。

---

## フライトコントローラー

### Pixhawk 6X

**PX4 v1.14 LTS** を搭載した **Pixhawk 6X** が、飛行に不可欠なすべての機能を担います：

| 仕様 | 値 |
|---------------|-------|
| **Pixhawk 6X** | SKU 11073 |
| **ベースボード** | SKU 18117 |
| **アクセサリー** | SKU 15011（ケーブル＆ケース） |
| **プロセッサー** | STM32H753（480MHz Cortex-M7） |
| **IMU** | 三重冗長（ICM-42688-P, ICM-45686, BMI088） |
| **気圧計** | デュアル（MS5611, ICP-20100） |
| **磁力計** | IST8310 |
| **インターフェース** | 3x CAN, 6x UART, SPI, I2C, PWM |

### フライトコントローラーのアクセサリー

| コンポーネント | SKU | 用途 |
|-----------|-----|---------|
| **PM06 V2 パワーモジュール** | 15019 | 電力検出 + 5V/3A BEC |
| **Holybro M9N GPS** | 12027 | GPS + コンパスモジュール |
| **GPSマスト** | 12033 | 干渉を低減する折りたたみ式マスト |
| **セーフティスイッチ＆ブザー** | 12007 | アーム時の安全確保 + 音声アラート |
| **防振マウント** | 12010 | 振動の減衰 |
| **ケーブルキット** | 12055 | JST-GH/CANケーブル |

**Pixhawk 6Xを選ぶ理由：**

- **冗長性** — センサー投票に対応する三重IMUとデュアル気圧計
- **PX4 LTSサポート** — 長期的な安定性とセキュリティパッチ
- **FMUv6X規格** — 他の準拠ボードと交換可能
- **防振構造** — IMUダンピングを内蔵

フライトコントローラーは以下を担当します：

- 姿勢推定と制御
- 位置保持とナビゲーション
- フェイルセーフ動作（RTL、着陸、ジオフェンス）
- RC入力処理
- MAVLinkテレメトリー

---

## コンパニオンコンピューター

フリート車両には、それぞれ異なる役割を持つ **デュアルコンパニオンコンピューター** が搭載されています：

### Raspberry Pi 4 / CM4（センサーコンパニオン）

軽量なセンサー統合とデータ収集を担当します：

| 仕様 | 値 |
|---------------|-------|
| **モデル** | Raspberry Pi 4 Model B |
| **SKU** | RPI4-MODEL-B-8GB |
| **プロセッサー** | クアッドコア Cortex-A72 @ 1.8GHz |
| **RAM** | 8GB |
| **ストレージ** | 産業用microSD 128GB（SDSDQED-128G-GNSIN） |
| **電源** | 5V @ 3A（標準） |

**役割：**

- 環境センサードライバー（温度、湿度、空気質）
- GPS/GNSSデータロギング
- カメラキャプチャ（AI以外のワークロード）
- JetsonへのMAVLinkルーティング
- Jetsonがオフラインのときのストア＆フォワード

### NVIDIA Jetson（AIコンパニオン）

計算負荷の高いワークロードを担当します：

| モデル | SKU | GPUコア | AI性能 | 消費電力 |
|-------|-----|-----------|----------------|-------|
| **Jetson Orin Nano 8GB** ★ | 900-13767-0000-000 | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin Nano Super** | 945-13766-0000-000 | 1024 CUDA | 67 TOPS | 7-25W |
| **Jetson Orin NX** | 900-13767-0040-000 | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 900-13701-0000-000 | 2048 CUDA | 275 TOPS | 15-60W |

★ = リファレンスBOMの選定品

| ストレージ | SKU | 容量 |
|---------|-----|----------|
| **NVMe SSD** | WD-SN530-256G | 256GB |

**役割：**

- コンピュータービジョン（物体検出、追跡）
- ビジュアル慣性オドメトリ
- 経路計画と障害物回避
- Vehicle Gateway（NATSクライアント）
- NATSリーフノード

フリート展開では通常、コスト効率の高い推論には **Jetson Orin Nano**、負荷の高い認識ワークロードには **Orin NX** を使用します。

---

## 無線システム

### RC制御：ExpressLRS

**ExpressLRS** は、パイロットによるオーバーライド用のリンクを提供します：

| 仕様 | 値 |
|---------------|-------|
| **レシーバー** | ExpressLRS EP1 |
| **SKU** | HM-EP1-2400 |
| **周波数** | 2.4GHz |
| **通信距離** | 10km以上（出力/アンテナによる） |
| **遅延** | <5ms |
| **Protocol** | Open-source (CRSF) |

**Why ExpressLRS:**

- **Open-source** — No vendor lock-in, community-driven development
- **Range** — Reliable link for visual line-of-sight operations
- **Low latency** — Responsive manual control when needed
- **Cost** — Affordable receivers for fleet scale

### Telemetry: SiK Radio

**SiK 915 MHz** radios provide ground station telemetry:

| Specification | Value |
|---------------|-------|
| **SKU** | 17012 (100mW) / 17013 (500mW) |
| **Frequency** | 915MHz (US) / 433MHz (EU) |
| **Range** | 300m-2km (depending on power) |
| **Data Rate** | Up to 250kbps |
| **Interface** | Serial (UART) |
| **Protocol** | MAVLink framing, FHSS |

**Role in Fleet:**

- Backup telemetry when cellular unavailable
- Ground control station connectivity
- Local testing and debugging

For production fleet operations, telemetry primarily flows through the Jetson's cellular connection to NATS. SiK radios serve as backup and for field diagnostics.

---

## Power System

### Battery: 4S LiPo

| Specification | Value |
|---------------|-------|
| **SKU** | TAA60004S35 |
| **Configuration** | 4S (14.8V nominal) |
| **Capacity** | 6000mAh |
| **Discharge Rate** | 35C |
| **Connector** | XT60 |

**Fleet Considerations:**

- Standardize on single battery configuration for logistics
- Use smart batteries with telemetry when available
- Implement battery rotation and health tracking
- Plan for 3:1 battery-to-vehicle ratio for continuous operations

### Power Distribution

| Component | SKU | Purpose |
|-----------|-----|---------|
| **PM06 V2 Power Module** | 15019 | Power sensing + 5V/3A BEC |

- **BEC** — 5V regulated supply for Pi, servos
- **Jetson power** — Direct from battery through regulator (12V typical)

### Networking

| Component | SKU | Purpose |
|-----------|-----|---------|
| **Industrial Ethernet Switch** | USR-ES105 | 3-port GbE between Pi, Jetson, payload |

Enables high-bandwidth communication between companion computers and connected payloads (cameras, sensors).

---

## Connectivity

### Cellular: 4G/5G + eSIM

Each Jetson connects via cellular modem:

| Specification | Value |
|---------------|-------|
| **Modem** | Quectel RM520N-GL |
| **SKU** | RM520NGLAA-M20-SGASA |
| **Bands** | 4G LTE + 5G NR Sub-6GHz |
| **SIM** | eSIM with OTA provisioning |
| **Interface** | M.2 (USB 3.0 / PCIe) |
| **Form Factor** | M.2 NGFF |

**eSIM Benefits:**

- No physical SIM swapping across large fleet
- Remote carrier provisioning
- Carrier switching for coverage optimization
- Centralized subscription management

See [Sensing platform](/platform/sensing/) for details on 5G/LTE connectivity architecture.

---

## Full Bill of Materials

| Category | Component | SKU |
|----------|-----------|-----|
| **Airframe** | Holybro X500 V2 ARF Kit | SKU30125 |
| **Flight Control** | Pixhawk 6X | SKU11073 |
| **Flight Control** | Pixhawk 6X Baseboard | SKU18117 |
| **Flight Control** | Pixhawk 6X Accessories | SKU15011 |
| **Power** | PM06 V2 Power Module | SKU15019 |
| **Navigation** | Holybro M9N GPS | SKU12027 |
| **Navigation** | GPS Mast | SKU12033 |
| **Safety** | Safety Switch & Buzzer | SKU12007 |
| **Companion** | Raspberry Pi 4 8GB | RPI4-MODEL-B-8GB |
| **Storage** | Industrial microSD 128GB | SDSDQED-128G-GNSIN |
| **Companion** | Jetson Orin Nano 8GB | 900-13767-0000-000 |
| **Storage** | NVMe SSD 256GB | WD-SN530-256G |
| **RC** | ExpressLRS EP1 | HM-EP1-2400 |
| **Telemetry** | SiK Telemetry Radio V3 | SKU17012 |
| **Networking** | Industrial Ethernet Switch | USR-ES105 |
| **Power** | LiPo Battery 4S 6000mAh | TAA60004S35 |
| **Mounting** | Pixhawk Vibration Mount | SKU12010 |
| **Cabling** | Holybro Cable Kit | SKU12055 |
| **Connectivity** | Quectel RM520N-GL 5G Modem | RM520NGLAA-M20-SGASA |

---

## Component Sellers

<details>
<summary><strong>コンポーネント別にすべての販売店を表示（57点）</strong></summary>

{{< bom-sellers >}}

</details>

---

## 次へ

[ソフトウェアスタック →]({{< relref "/fleet/software" >}})
