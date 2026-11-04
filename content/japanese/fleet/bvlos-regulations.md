---
title: "BVLOS（目視外飛行）規制"
meta_title: "目視外飛行ドローンの規制 | Ubuntu Software"
description: "世界各国のBVLOS（目視外飛行）ドローン規制。商用フリート運用における規制要件を理解する。"
image: "/images/robotics.svg"
draft: false
---

## フリート運用においてBVLOSが重要な理由

目視外飛行（BVLOS: Beyond Visual Line of Sight）に関する規制は、商用ドローンフリートが何を実現できるかを根本的に左右します。規制の全体像を理解することは、フリートプラットフォームの設計と市場戦略に不可欠です。

| 課題 | フリートプラットフォームへの影響 |
|-----------|-------------------------|
| **運用上の制限** | VLOSの制限（400〜500mの範囲）により、フリートの有用性は限られます。BVLOSにより、インフラ点検、配送、農業を大規模に展開できます。 |
| **コンプライアンスの統合** | フリートソフトウェアは規制上の境界を遵守させる必要があります。リモートID、検知・回避、運用許可の要件は、テレメトリーおよび安全機能の設計に反映されます。 |
| **市場参入のタイミング** | BVLOS解禁の段階は国によって異なります。どの国で規制が緩和されつつあるかを把握することで、ターゲット市場の優先順位付けに役立ちます。 |
| **パートナーの準備状況** | フリート運用者が実際のミッションを行うにはBVLOS承認が必要です。承認取得者を追跡することで、潜在的な顧客やパートナーを特定できます。 |

---

## グローバルな動向

主要な航空当局全体で、4つの動向が見られます：

| 動向 | 説明 |
|-------|-------------|
| **SORAの採用** | FAA、EASA、CASA、カナダ運輸省がSORA（Specific Operations Risk Assessment：特定運用リスク評価）手法を採用 [[1]](#ref-1) |
| **リモートIDの義務化** | すべての法域で義務化が進行中 |
| **ウェイバー → 規則** | 個別審査によるウェイバー（適用除外許可）から、標準化された性能ベースの規則へ移行中 |
| **UTMの統合** | 各国がスケーラブルなBVLOSに向けてUAS交通管理（UTM）を開発中 |

---

## 南北アメリカ

### 米国（FAA）

**ステータス：** ウェイバーベースから規則ベース（Part 108）へ移行中 [[2]](#ref-2)

**タイムライン** [[3]](#ref-3)：

- 2024年10月時点で190件のBVLOSウェイバーを発行（2020年の6件から増加）
- Part 108の規則案は2026年1月に公表予定
- 2024年FAA再授権法によりBEYONDプログラムが2029年まで延長 [[4]](#ref-4)

**BVLOS承認を取得した運用者：**

| 運用者 | ユースケース | 承認日 | 参照 |
|----------|----------|---------------|-----|
| DroneUp | 医療品配送 | 2024年1月 | [[5]](#ref-5) |
| Amazon Prime Air | 消費者向け配送 | 2024年5月 | [[6]](#ref-6) |
| Asylon | セキュリティ／監視 | 2024年4月 | [[7]](#ref-7) |
| 複数（DFW地域） | 初の複数運用者による共有空域 | 2024年7月 | [[8]](#ref-8) |

---

### カナダ（カナダ運輸省）

**ステータス：** SFOCが必要、定常的なBVLOSへ移行中 [[9]](#ref-9)

**枠組み：**

- 2025年4月：低リスク条件下（150kg以下、非管制空域、人口希薄地域）では、SFOCなしで定常的なBVLOSが許可
- 上級操縦士証明書（Advanced Pilot Certificate）が必要
- 2030年までにNav Canadaによる全国的なUTMの枠組みを整備

**運用者：** 未定

---

## ヨーロッパ

### 欧州連合（EASA）

**ステータス：** SORA／PDRAの経路を備えた特定（Specific）カテゴリー [[10]](#ref-10)

**枠組み** [[1]](#ref-1)：

| 項目 | 詳細 |
|--------|---------|
| 標準シナリオ | STS-02：空域監視員を配置したBVLOS |
| リモートID | 2024年1月より義務化 [[11]](#ref-11) |
| 承認期間 | 4〜9か月（医療）、6〜12か月（消費者向け） [[12]](#ref-12) |
| 手数料 | €2,000〜€15,000 |

**運用者：** 各国のサンドボックスや加盟国の迅速承認プログラムを通じて、さまざまな運用者が参加。

---

### 英国（CAA）

**ステータス：** UK SORAによる運用許可（OA）

**枠組み：**

- BVLOSにはCAAの運用許可が必要
- UK SORAプロセス（2025年4月以降）
- 通常、検知・回避技術（ADS-B、FLARM、地上レーダー）が必要

**運用者：**

| 運用者 | 実績 | 参照 |
|----------|-------------|-----|
| sees.ai | 英国初の定常的BVLOS（2021年）、National Gridの点検（2023年） | [[13]](#ref-13) |

---

### アイルランド（IAA）

**ステータス：** PDRA G-03の経路

**枠組み：** EASAを通じて、EU全域で簡素化された許可を得る道を開く。

**運用者：**

| 運用者 | 実績 | 参照 |
|----------|-------------|-----|
| sees.ai | アイルランド初のBVLOS許可（2025年） | [[13]](#ref-13) |

---

## アジア太平洋

### オーストラリア（CASA）

**ステータス：** Part 101に基づくSORAベースの評価 [[14]](#ref-14)

**枠組み** [[15]](#ref-15)：

- 定量的リスクモデリングを備えたSORA 2.5を採用
- 広域BVLOS自己評価（Broad Area BVLOS Self-Assessment）の枠組みを利用可能
- 簡素化されたプロセスがグローバル企業を惹きつけている [[16]](#ref-16)

**運用者：**

| 運用者 | 実績 | 参照 |
|----------|-------------|-----|
| Sphere Drones | 広域BVLOS承認（2024年5月） | [[17]](#ref-17) |

---

### 日本（国土交通省／航空局）

**ステータス：** 有人地帯上空でのレベル4 BVLOSが許可 [[18]](#ref-18)

**枠組み** [[19]](#ref-19)：

| レベル | 説明 | 必要なライセンス |
|-------|-------------|------------------|
| レベル1-2 | 目視内（VLOS）飛行 | 二等 |
| レベル3 | 無人地帯上空でのBVLOS | 二等 |
| レベル4 | 都市部／有人地帯上空でのBVLOS | 一等 |

- ライセンスの有効期間は3年で、登録講習機関での受講＋実地試験が必要
- 2024年10月時点で422,879機の無人航空機が登録済み

**運用者：**

- 2024年10月時点で、第一種（レベル4）型式認証を取得した機種は1機種のみ
- 第二種型式認証を取得した機種は5機種

---

### シンガポール（CAAS）

**ステータス：** 運用者許可（Operator Permit）＋クラス1活動許可（Class 1 Activity Permit）なしのBVLOSは禁止 [[20]](#ref-20)

**枠組み** [[21]](#ref-21)：

- いかなるBVLOSにも厳格な許可が必要
- 25kg超のドローンまたはBVLOSには、正式な運用者許可＋クラス1許可が必要
- 2025年12月から、250g超のドローンにB-RID（ブロードキャスト型リモートID）が義務化
- B-RIDは、将来の大規模BVLOS／自動化運用の基盤となる
- 迅速承認のための規制サンドボックスプログラムを利用可能

**運用者：** 未定（サンドボックスでの試験が進行中）

---

### 中国（CAAC）

**ステータス：** BVLOSには特別な許可が必要 [[22]](#ref-22)

**枠組み** [[23]](#ref-23)：

- 無人航空機飛行管理暫定条例（2024年1月）
- BVLOSには高度運用証明書（Advanced Operations Certificate）＋特別飛行運用証明書（SFOC）が必要
- 全国で高度120mの制限があり、それを超える場合は許可が必要
- デジタル申請プロセスの簡素化により、BVLOS飛行が増加 [[16]](#ref-16)

**運用者：** 未定（CAACの枠組み下にある国内運用者）

---

### インド（DGCA）

**ステータス：** 正式なBVLOSの枠組みはまだなく、試験段階のみ [[24]](#ref-24)

**枠組み** [[25]](#ref-25)：

- 2021年ドローン規則（Drone Rules, 2021）の下、プロジェクト単位の適用除外により運用
- 2025年4月：都市部／準都市部向けに独立したBVLOS操縦士ライセンス区分を導入
- DGCAは「BVLOS規則の最終決定に向けた最終段階にある」 [[26]](#ref-26)
- 実験的な飛行回廊が稼働中（例：テランガーナ州の「Medicine from the Sky」）

**運用者：**

| 運用者 | 実績 | 参照 |
|----------|-------------|-----|
| Skye Air Mobility | 世界最長のBVLOS医薬品配送（104km、Baruipur〜Medinipur） | [[27]](#ref-27) |
| Flipkart Health | 104km配送試験のパートナー | — |

---

## コンプライアンス要件

すべての法域において、BVLOS規制は共通の技術要件に収束しつつあります。本セクションでは、当社のフリートプラットフォームが実装すべき内容を整理します。

### 1. リモートID（全地域で義務化）

すべての主要規制当局は、ドローンに識別情報と位置データのブロードキャストを義務付けています。

| 法域 | 要件 | 期限 | 参照 |
|--------------|-------------|----------|-----|
| 米国（FAA） | 標準リモートIDまたはブロードキャストモジュール | 2023年9月（施行済み） | [[2]](#ref-2) |
| EU（EASA） | 特定カテゴリーおよびクラス1以上のリモートID | 2024年1月（施行済み） | [[11]](#ref-11) |
| シンガポール | 250g超のドローンにB-RID | 2025年12月 | [[21]](#ref-21) |
| オーストラリア | リモートIDを検討中 | 未定 | [[14]](#ref-14) |

**プラットフォーム要件：**
- ブロードキャストモジュールのハードウェア統合（OpenDroneID互換）
- フリート全体を追跡するための、NATS経由のネットワーク型リモートID機能
- 監査用に、すべてのリモートID送信をJetStreamに記録

---

### 2. 検知・回避（DAA）

BVLOS運用には、他の航空機や障害物を検知・回避する能力が求められます。

| 法域 | 要件 | 備考 |
|--------------|-------------|-------|
| 米国（FAA） | 「Well clear（十分な離隔）」基準。Part 108ではDAAシステムが必要 | 規定型ではなく性能ベース |
| 英国（CAA） | 通常、DAA技術（ADS-B、FLARM、地上レーダー）が必要 | OAプロセスで個別審査 |
| 日本 | レベル4の都市部運用で必要 | [[18]](#ref-18) |

**プラットフォーム要件：**
- ADS-B In受信機の統合（周辺交通の把握）
- Jetson上でのセンサーフュージョン（カメラ、レーダー、LiDARのオプション）
- 自律回避機動またはオペレーターへの警告
- すべてのDAAイベントをJetStreamに記録

---

### 3. ジオフェンシング

制限区域への飛行を防ぐ、ハードウェアで強制される境界。

| 法域 | 要件 | 備考 |
|--------------|-------------|-------|
| 米国（FAA） | ジオフェンスは運用承認の一部 | ウェイバー／Part 108で規定 |
| EU（EASA） | 特定カテゴリーではジオアウェアネスが必要 | SORAの地上リスク評価に準拠 |
| 全地域 | 高度制限（通常120m／400ft AGL） | すべての枠組みの基本 |

**プラットフォーム要件：**
- フライトコントローラーでのジオフェンス強制（PX4 `GF_*` パラメータ）
- Gateway経由の動的ジオフェンス更新（一時的な制限）
- 高度上限の強制
- 逸脱の記録とフリートへの警告

---

### 4. 指揮統制（C2）リンク

オペレーターと機体の間の、信頼性が高く冗長化された通信。

| 法域 | 要件 | 備考 |
|--------------|-------------|-------|
| 米国（FAA） | Part 108におけるC2リンクの性能基準 | 遅延、可用性、セキュリティ |
| EU（EASA） | SORAに基づくC2リンクの信頼性 | リスクに応じた要件 |
| 全地域 | C2喪失時のフェイルセーフ動作 | 帰還、着陸、またはロイター（旋回待機） |

**プラットフォーム要件：**
- LTE/5Gまたは専用無線によるプライマリC2
- バックアップC2経路（例：長距離向けの衛星通信）
- リンク喪失時の自動フェイルセーフ（PX4で設定したRTL）
- フリートへのC2リンク状態のテレメトリー送信
- 監査用のリンク品質記録

---

### 5. 運用許可の文書化

規制当局は、文書化されたリスク評価と運用手順を求めています。

| 枠組み | 採用機関 | 主要文書 |
|-----------|---------|---------------|
| SORA | EASA、英国CAA、CASA、カナダ運輸省 | ConOps、地上リスク、空中リスク、OSOによる緩和策 |
| Part 107 ウェイバー | FAA（現行） | セーフティケース、運用制限 |
| Part 108 | FAA（今後） | 性能ベースのコンプライアンス |

**プラットフォーム要件：**
- 規制当局が受け入れる形式での飛行ログのエクスポート
- ConOpsを裏付けるデータの生成（飛行エリア、時間、頻度）
- 運用者／操縦者の資格記録の管理
- 機体の整備・点検状況の追跡

---

### 6. 操縦者／運用者のライセンス

| 法域 | BVLOSライセンス | 備考 |
|--------------|---------------|-------|
| 米国 | Part 107＋ウェイバー（現行）、Part 108（将来） | [[2]](#ref-2) |
| EU | A2/STS証明書または特定カテゴリーの許可 | [[1]](#ref-1) |
| 英国 | GVC＋運用許可 | |
| 日本 | 一等（レベル4）、二等（レベル3） | [[19]](#ref-19) |
| カナダ | 上級操縦士証明書＋SFOC | [[9]](#ref-9) |
| インド | BVLOS操縦士ライセンス（2025年4月以降） | [[25]](#ref-25) |

**プラットフォーム要件：**
- 飛行許可前の操縦者認証
- ライセンス／資格の検証と有効期限の追跡
- 操縦者ごとの飛行時間の記録
- ロールベースのアクセス制御（操縦者／監視員／フリートマネージャー）

---

### 7. 飛行記録と監査証跡

すべての法域で、包括的な飛行記録が求められます。

**プラットフォーム要件：**
- テレメトリーの完全な記録（位置、姿勢、バッテリー、センサー）
- タイムスタンプ付きのイベント記録（アーム、離陸、ウェイポイント、着陸、ディスアーム）
- 改ざん不可能なストレージ（保持ポリシー付きJetStream）
- 規制当局への提出用エクスポート機能
- インシデントのフラグ付けと調査支援

---

### 8. 保険と賠償責任

| 法域 | 要件 |
|--------------|-------------|
| EU | 第三者賠償責任保険が義務 |
| 米国 | 通常、ウェイバーの条件として保険が必要 |
| 全地域 | 許可取得時に保険加入の証明を求められることが多い |

**プラットフォーム要件：**
- 機体ごとの保険契約の追跡
- 飛行許可前の補償範囲の検証
- 保険請求のためのインシデント文書化

---

## 外部連携

BVLOSコンプライアンスには、外部システムとのリアルタイムなデータ交換が必要です。以下は、当社プラットフォームがサポートすべき連携です。

### 空域と交通

| システム | 目的 | データフロー | 優先度 |
|--------|---------|-----------|----------|
| **ADS-B In** | 有人航空機の交通状況の把握 | 受信 → DAAシステム | 高 |
| **FLARM** | グライダー／一般航空機の交通状況の把握 | 受信 → DAAシステム | 中 |
| **UTM/USS** | UAS交通管理、飛行許可 | 双方向 | 高 |
| **LAANC**（米国） | 空域許可の自動化 | リクエスト／レスポンス | 高（米国） |
| **AIM/NOTAM** | 一時飛行制限、空域の変更 | 受信 → ジオフェンス | 高 |

**UTMサービスサプライヤー（USS）：**
- 米国：Airspace Link、Aloft、DroneUp、FlyFreely（[FAA LAANCリスト](https://www.faa.gov/uas/getting_started/laanc)を参照）
- EU：加盟国ごとにさまざま
- オーストラリア：Wing、Airservices Australia

*注：AirMapは2023年6月にLAANCサービスを終了しました。*

---

### リモートIDネットワーク

| システム | 目的 | プロトコル | 優先度 |
|--------|---------|----------|----------|
| **OpenDroneID** | ブロードキャスト型リモートIDの標準 | Bluetooth 4/5、WiFi Beacon | 必須 |
| **ネットワーク型リモートID** | インターネット経由でのフリート全体のID | USS API | 必須（米国） |
| **FAAリモートID** | 米国のリモートIDへの準拠 | USS連携 | 必須（米国） |

---

### 規制関連データベース

| システム | 目的 | データフロー | 優先度 |
|--------|---------|-----------|----------|
| **操縦者ライセンスDB** | 操縦者の資格の検証 | 照会 → 許可 | 高 |
| **機体登録簿** | 登録状況の検証 | 照会 → 飛行前点検 | 高 |
| **保険の検証** | 補償の有効性の確認 | 照会 → 飛行前点検 | 中 |
| **整備記録** | 耐空性の追跡 | 双方向 | 中 |

---

### 気象と環境

| システム | 目的 | データフロー | 優先度 |
|--------|---------|-----------|----------|
| **航空気象（METAR/TAF）** | 風、視程、雲底高度 | 受信 → 飛行計画 | 高 |
| **レーダー／降水** | リアルタイムの気象ハザード | 受信 → 飛行中の警告 | 中 |
| **太陽活動／地磁気** | GPS信頼性の予測 | 受信 → 飛行計画 | 低 |

---

### 緊急サービス

| システム | 目的 | データフロー | 優先度 |
|--------|---------|-----------|----------|
| **ATC（必要な場合）** | 管制空域での調整 | 音声／データリンク | 状況に応じて |
| **緊急サービス（911/112）** | インシデント報告 | 外部へのアラート送信 | 必須 |
| **NOTAMの提出** | 飛行運用の公示 | 送信 → AIM | 状況に応じて |

---

### 連携アーキテクチャ

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

### 連携ロードマップ

| 連携 | ステータス | 備考 |
|-------------|--------|-------|
| OpenDroneIDブロードキャスト | ✅ 対応済み | ハードウェアモジュールが必要 |
| ADS-B In受信機 | 🔧 ロードマップ | uAvionix pingRXまたは同等品 |
| UTM/USS（Airspace Link） | 🔧 ロードマップ | AirHub API、段階的な料金体系 |
| UTM/USS（Aloft） | 🔧 ロードマップ | REST API連携 |
| LAANC | 🔧 ロードマップ | USSパートナー経由 |
| 航空気象 | 🔧 ロードマップ | aviationweather.gov（オープン、100リクエスト/分） |
| NOTAMフィード | 🔧 ロードマップ | NASA DIPまたはFAA SWIM |
| 操縦者ライセンスの検証 | 🔧 ロードマップ | 法域ごとのAPI |
| 保険の検証 | 🔧 ロードマップ | プロバイダーのAPI |

---

## コンプライアンス・マトリクス

| 要件 | プラットフォームコンポーネント | ステータス |
|-------------|-------------------|--------|
| リモートIDのブロードキャスト | OpenDroneIDモジュール＋Gateway | ✅ 対応済み |
| リモートIDの記録 | JetStreamストリーム | ✅ 対応済み |
| 検知・回避 | Jetsonによるセンサーフュージョン | 🔧 ロードマップ |
| ADS-B In | 交通情報受信機の統合 | 🔧 ロードマップ |
| ジオフェンシング（静的） | PX4 `GF_*` パラメータ | ✅ 対応済み |
| ジオフェンシング（動的） | Gatewayのポリシー更新 | ✅ 対応済み |
| C2リンク（プライマリ） | NATS経由のLTE/5G | ✅ 対応済み |
| C2リンク（バックアップ） | 衛星通信の統合 | 🔧 ロードマップ |
| リンク喪失時のフェイルセーフ | PX4のフェイルセーフ設定 | ✅ 対応済み |
| 飛行記録 | JetStreamテレメトリーストリーム | ✅ 対応済み |
| 監査用エクスポート | ログエクスポートツール | ✅ 対応済み |
| 操縦者認証 | 認可システム | ✅ 対応済み |
| ライセンスの追跡 | フリート管理 | 🔧 ロードマップ |
| 保険の追跡 | フリート管理 | 🔧 ロードマップ |

これらが全体的な安全アーキテクチャとどのように統合されるかについては、[安全モデル]({{< relref "/fleet/safety" >}})を参照してください。

---

## 参考文献

<span id="ref-1"></span>**[1]** EASA. "Drones (UAS) - Frequently Asked Questions." European Union Aviation Safety Agency. https://www.easa.europa.eu/en/the-agency/faqs/drones-uas

<span id="ref-2"></span>**[2]** Drone Pilot Ground School. "Part 108: What We Know So Far about the FAA's BVLOS Rule." https://www.dronepilotgroundschool.com/part-108/

<span id="ref-3"></span>**[3]** DroneLife. "FAA Progress on BVLOS Rules and Advanced Air Mobility Integration." Aug 2024. https://dronelife.com/2024/08/01/faa-progress-on-bvlos-rules-and-advanced-air-mobility-integration/

<span id="ref-4"></span>**[4]** FAA. "BEYOND Program." Federal Aviation Administration. https://www.faa.gov/uas/programs_partnerships/beyond

<span id="ref-5"></span>**[5]** DroneUp. "DroneUp Awarded Landmark FAA Approval for Beyond Visual Line of Sight (BVLOS) in the U.S." Jan 2024. https://www.droneup.com/news/droneup-awarded-bvlos

<span id="ref-6"></span>**[6]** DroneDJ. "Amazon gets FAA BVLOS approval for Prime Air drone delivery." May 2024. https://dronedj.com/2024/05/30/amazon-gets-faa-bvlos-approval-for-prime-air-drone-delivery/

<span id="ref-7"></span>**[7]** DroneDJ. "Asylon bags new FAA approval for BVLOS drone operations." Apr 2024. https://dronedj.com/2024/04/26/asylon-drone-faa-bvlos-dronesentry/

<span id="ref-8"></span>**[8]** DroneLife. "FAA Grants Historic BVLOS Approval for Multiple Operators in DFW Area." Jul 2024. https://dronelife.com/2024/07/31/faa-grants-historic-bvlos-approval-for-multiple-operators-in-dfw-area/

<span id="ref-9"></span>**[9]** MLT Aikins. "International drone regulations: Key updates on BVLOS, UTM and use cases." https://www.mltaikins.com/insights/international-drone-regulations-key-updates-on-bvlos-utm-and-use-cases/

<span id="ref-10"></span>**[10]** FlytBase. "Understanding EU Drone Regulations for Docked BVLOS Operations." https://www.flytbase.com/blog/eu-drone-regulations-for-bvlos-operations

<span id="ref-11"></span>**[11]** Elsight. "EASA Remote ID Requirements: Compliance for Drone Operators." https://www.elsight.com/blog/primer-on-easa-remote-id-regulations/

<span id="ref-12"></span>**[12]** Sparkco. "Autonomous Drone Delivery Regulatory Approval." https://sparkco.ai/blog/autonomous-drone-delivery-regulatory-approval

<span id="ref-13"></span>**[13]** sees.ai. "Why Us - First UK routine BVLOS permission." https://www.sees.ai/why-us/

<span id="ref-14"></span>**[14]** CASA. "BVLOS drone operations in regional Australia - Consultation." Civil Aviation Safety Authority. https://consultation.casa.gov.au/stakeholder-engagement-group/consultation.2023-10-05.5578154857/

<span id="ref-15"></span>**[15]** Measure Australia. "New CASA Drone Regulations Incoming!" https://www.measureaustralia.com.au/news/new-casa-drone-regulations-incoming

<span id="ref-16"></span>**[16]** Droneii. "Commercial Drone Regulation Report 2025." Drone Industry Insights. https://droneii.com/product/commercial-drone-regulation-report

<span id="ref-17"></span>**[17]** DroneLife. "Sphere Drones Secures BVLOS Area Approval to Advance Commercial Drone Use in Australia." May 2024. https://dronelife.com/2024/05/08/sphere-drones-secures-bvlos-area-approval-to-advance-commercial-drone-use-in-australia/

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

## 概要に戻る

[← 自律型車両フリートアーキテクチャ]({{< relref "/fleet" >}})
