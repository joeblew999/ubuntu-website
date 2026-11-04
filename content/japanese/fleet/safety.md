---
title: "安全モデル"
meta_title: "自律型フリート安全アーキテクチャ | Ubuntu Software"
description: "自律走行車両フリートの安全モデル：手動オーバーライド権限、車両フェイルセーフ、制御ループからのネットワーク分離、そしてグレースフル・デグラデーション。"
image: "/images/robotics.svg"
draft: false
---

## 設計による安全性

フリート規模の自律型車両の運用には、厳格な安全アーキテクチャが求められます。この設計では、**ネットワーク障害、ソフトウェアのバグ、システムの誤動作によって車両の安全性が損なわれることは決してありません**。

中核となる安全原則は、ドローン、乗用車、トラック、AGVといった**すべての車両タイプ**に適用されます。具体的な内容（RCかステアリングホイールか、RTLかその場停止か）は異なりますが、アーキテクチャは共通です。

---

## 基本原則

**NATSがリアルタイム制御ループに入ることは決してありません。**

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SAFETY HIERARCHY                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   1. MANUAL OVERRIDE      ─────────────────────────▶  HIGHEST       │
│      Human control (RC, steering, e-stop), always available          │
│                                                                      │
│   2. VEHICLE FAILSAFES    ─────────────────────────▶  HIGH          │
│      Hardware-enforced, local to vehicle                             │
│                                                                      │
│   3. GATEWAY POLICY       ─────────────────────────▶  MEDIUM        │
│      Software enforcement, on-vehicle                                │
│                                                                      │
│   4. FLEET COMMANDS       ─────────────────────────▶  LOWEST        │
│      Network-dependent, advisory                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

上位レイヤーは常に下位レイヤーをオーバーライドできます。ネットワークに依存するコマンドの権限は最も低くなります。

---

## 手動オーバーライドが最優先の権限

人間のオペレーターは、RC送信機、ステアリングホイール、緊急停止ボタンのいずれを使う場合でも、**無条件のオーバーライド権限**を持ちます。

### ドローンのオーバーライド（RC）

| 項目 | 実装 |
|--------|----------------|
| **プロトコル** | ExpressLRSで受信機へ直接接続 |
| **経路** | TX → RX → Pixhawk（ループ内にソフトウェアなし） |
| **モードスイッチ** | マニュアル/スタビライズへのハードウェアフェイルセーフ |
| **キルスイッチ** | モーターを即時ディスアーム |

**RCは以下のいずれにも依存しません：**

- Jetsonがオンラインであること
- NATS接続
- Vehicle Gateway
- あらゆるソフトウェアコンポーネント

すべてのコンピューターが故障しても、パイロットは完全な手動制御を維持します。

#### ExpressLRSのフェイルセーフ

RC信号が失われた場合：

```
1. RX detects signal loss (100ms)
2. RX sends failsafe values to Pixhawk
3. Pixhawk enters RC_LOSS failsafe mode
4. Configured action executes (RTL, land, etc.)
```

フェイルセーフの動作はPX4で設定され、フリートソフトウェアには一切依存しません。

### 地上車両のオーバーライド

| 項目 | 乗用車/トラック | AGV |
|--------|-------------|------|
| **プライマリ** | ステアリングホイール＋ペダル | 緊急停止ボタン |
| **経路** | アクチュエーターへの物理的接続 | モーターへのハードウェア割り込み |
| **バックアップ** | キーオフ / ギアをニュートラル | ワイヤレス緊急停止 |
| **ネットワーク** | 制御経路に一切含まれない | 制御経路に一切含まれない |

**地上車両のオーバーライドは以下のいずれにも依存しません：**

- オンボードコンピューターがオンラインであること
- NATS接続
- あらゆる自律ソフトウェア

すべてのコンピューターが故障しても、ドライバー/オペレーターは物理的な制御を維持するか、車両が安全に停止します。

---

## 車両制御システムがフェイルセーフを強制

車両の制御システムはフェイルセーフをローカルで実装しており、ネットワークは不要です。

### ドローンのフェイルセーフ（PX4）

フライトコントローラーは複数のフェイルセーフを実装しています：

#### RC喪失フェイルセーフ

| パラメーター | 値 | アクション |
|-----------|-------|--------|
| `COM_RC_LOSS_T` | 0.5s | フェイルセーフ発動までのタイムアウト |
| `NAV_RCL_ACT` | 2 | RC喪失時にRTL |

#### データリンク喪失

| パラメーター | 値 | アクション |
|-----------|-------|--------|
| `COM_DL_LOSS_T` | 10s | GCS接続タイムアウト |
| `NAV_DLL_ACT` | 0 | ミッションを継続（リンクはクリティカルではない） |

#### バッテリーフェイルセーフ

| レベル | パラメーター | 値 | アクション |
|-------|-----------|-------|--------|
| 低 | `BAT_LOW_THR` | 0.25 | 警告 |
| 危険 | `BAT_CRIT_THR` | 0.15 | RTL |
| 緊急 | `BAT_EMERGEN_THR` | 0.07 | 即時着陸 |

#### ジオフェンス

| パラメーター | 値 | アクション |
|-----------|-------|--------|
| `GF_ACTION` | 3 | 逸脱時にRTL |
| `GF_MAX_HOR_DIST` | 可変 | 水平方向の制限 |
| `GF_MAX_VER_DIST` | 120m | 高度制限 |

#### 位置喪失

| パラメーター | 値 | アクション |
|-----------|-------|--------|
| `COM_POS_FS_DELAY` | 1s | 位置タイムアウト |
| `COM_POSCTL_NAVL` | 0 | 位置喪失時に着陸 |

**すべてのフェイルセーフはPixhawk上で実行されます**。外部への依存はありません。

### 地上車両のフェイルセーフ

地上車両も同等の安全メカニズムを実装しています：

| トリガー | ドローンの対応 | 地上車両の対応 |
|---------|---------------|------------------------|
| **通信喪失** | Return-to-Launch | その場で停止 / 路肩に停車 |
| **センサー故障** | 即時着陸 | 安全に停止し、パーキングを作動 |
| **ジオフェンス逸脱** | 境界でRTL | 境界で停止 |
| **緊急停止の作動** | モーター遮断 | 即時ブレーキ |
| **バッテリー/燃料低下** | RTL後に着陸 | デポへ帰還 / 安全に停止 |
| **障害物検知** | 回避またはホバリング | 停止して待機 |

**地上車両のフェイルセーフには通常以下が含まれます：**

- **LiDAR/レーダーによる障害物検知** — 経路が塞がれた場合は緊急停止
- **速度制限** — 環境に応じた最高速度
- **ゾーン制限** — 進入禁止エリアをローカルで強制
- **ウォッチドッグタイマー** — 自律ソフトウェアがハングした場合に停止
- **ブレーキの冗長化** — 複数の独立したブレーキシステム

**すべてのフェイルセーフはローカルで実行されます**。ネットワークへの依存はありません。

---

## NATSはモニタリング専用

ネットワークレイヤーが担当するもの：

| 機能 | 説明 |
|----------|-------------|
| **テレメトリー収集** | ダッシュボード向けの車両状態 |
| **イベントロギング** | JetStreamによる監査証跡 |
| **ミッション調整** | フリートレベルの計画 |
| **シャドウステート** | desired/reportedの同期 |

ネットワークレイヤーが**担当しないもの：**

| 機能 | 理由 |
|----------|---------|
| **姿勢制御** | レイテンシーが許容できず、安全上クリティカル |
| **モーターコマンド** | FCUでローカルに処理する必要がある |
| **センサーフュージョン** | リアルタイム要件がある |
| **フェイルセーフの判断** | ネットワークなしで動作する必要がある |

### コマンド検証

NATS経由でコマンドが到着した場合：

```go
func (g *Gateway) validateCommand(cmd *Command) error {
    // Check geofence
    if cmd.Type == "goto" {
        if !g.geofence.Contains(cmd.Target) {
            return errors.New("target outside geofence")
        }
    }

    // Check altitude
    if cmd.Altitude > g.policy.MaxAltitude {
        return errors.New("altitude exceeds limit")
    }

    // Check vehicle state
    if cmd.RequiresArmed && !g.state.Armed {
        return errors.New("vehicle not armed")
    }

    // Check rate limiting
    if g.rateLimiter.Exceeded(cmd.Type) {
        return errors.New("command rate exceeded")
    }

    return nil
}
```

コマンドは**実行前に検証されます**。仮に不正なコマンドが実行されたとしても、最終的な保護はPX4のフェイルセーフが担います。

---

## グレースフルデグラデーション

コンポーネントが故障した場合、システムは安全に縮退します：

### Jetsonの故障

```
Jetson offline
    └─▶ No perception, no AI inference
    └─▶ Pi continues MAVLink routing
    └─▶ RC control fully operational
    └─▶ PX4 failsafes active
    └─▶ Vehicle lands safely or continues mission
```

### NATS接続の喪失

```
WAN disconnected
    └─▶ Leaf node buffers locally
    └─▶ Gateway continues operation
    └─▶ Vehicle follows current mission
    └─▶ State syncs when connection returns
```

### Vehicle Gatewayのクラッシュ

```
Gateway process dies
    └─▶ systemd restarts within 5 seconds
    └─▶ PX4 continues autonomous flight
    └─▶ RC override always available
    └─▶ No fleet commands until restart
```

### ハブクラスターの障害

```
Regional hub down
    └─▶ Leaf nodes queue messages
    └─▶ Vehicles continue operating
    └─▶ RC and PX4 failsafes unaffected
    └─▶ Fleet visibility lost temporarily
```

### 連鎖障害の防止

| 障害 | 影響範囲 | 緩和策 |
|---------|--------------|------------|
| 単一車両 | その車両のみ | フェイルセーフで着陸 |
| ハブノード | 影響なし（レプリケーション済み） | クラスターが自己修復 |
| ハブクラスター | 地域フリートの調整 | 車両は自律的に動作を継続 |
| グローバル層 | リージョン間の可視性 | 地域の運用は継続 |

---

## 運用安全手順

### 飛行前チェックリスト

フリート車両をアームする前に：

1. **RCリンクの確認** — パイロットが制御権限を確認
2. **フェイルセーフテスト** — RTL/着陸動作を検証
3. **ジオフェンスの読み込み** — PX4で境界を確認
4. **バッテリーしきい値** — フェイルセーフレベルが適切であること
5. **気象確認** — 風速が制限内であること
6. **空域のクリア** — 競合がないこと

### 緊急時の手順

| 緊急事態 | 対応 |
|-----------|----------|
| **目視の喪失** | RC経由でRTLを作動 |
| **フライアウェイ** | RCキルスイッチ、またはジオフェンスによるRTL発動 |
| **複数車両の競合** | フリート全体に着陸コマンド |
| **通信の完全喪失** | PX4のRC_LOSSフェイルセーフ |
| **火災/墜落** | RCでディスアーム、緊急サービスへ連絡 |

### キルチェーン

車両を停止させる複数の方法：

```
1. RC kill switch      → Immediate motor stop
2. RC mode switch      → Land mode
3. GCS command         → RTL or land
4. Fleet command       → RTL or land
5. Geofence breach     → Auto-RTL
6. Battery critical    → Auto-land
```

少なくとも3つの方法（RC、PX4フェイルセーフ、ローカルタイムアウト）が常に利用可能です。

---

## 認証に関する考慮事項

このアーキテクチャは規制への準拠を支援します：

| 要件 | 実装 |
|-------------|----------------|
| **機長（Pilot in command）** | 常時RCによるオーバーライド |
| **フェイルセーフ動作** | PX4の認証済みフェイルセーフ |
| **ジオフェンシング** | ハードウェアで強制される境界 |
| **ロギング** | JetStreamによる監査証跡 |
| **リモートID** | 独立したモジュールでブロードキャスト |

**注：** 規制要件は管轄区域によって異なります。このアーキテクチャは技術的な基盤を提供するものであり、具体的な準拠には現地の規制に応じた追加措置が必要です。

---

## 安全システムのテスト

### ユニットテスト

```go
func TestGeofenceRejection(t *testing.T) {
    gw := NewGateway(config)
    gw.SetGeofence(testGeofence)

    cmd := &Command{
        Type:   "goto",
        Target: Position{Lat: 0, Lon: 0}, // Outside geofence
    }

    err := gw.validateCommand(cmd)
    assert.ErrorContains(t, err, "outside geofence")
}
```

### 統合テスト

```
1. Simulate RC loss → Verify RTL activates
2. Simulate battery drain → Verify land at threshold
3. Simulate geofence breach → Verify return behavior
4. Simulate Jetson crash → Verify continued flight
5. Simulate NATS disconnect → Verify local operation
```

### 飛行テスト

フリート展開前に：

1. **単一車両のフェイルセーフ検証**
2. **RC喪失時の応答確認**
3. **ジオフェンス逸脱テスト**
4. **バッテリーフェイルセーフの検証**
5. **複数車両の連携テスト**

---

## まとめ

| レイヤー | 権限 | ドローン | 地上車両 |
|-------|-----------|--------|-----------------|
| **手動** | 最高 | RC送信機 | ステアリング / 緊急停止 |
| **車両** | 高 | PX4フェイルセーフ | ECU / 安全コントローラー |
| **ゲートウェイ** | 中 | ポリシーの強制 | ポリシーの強制 |
| **フリート** | 最低 | 調整のみ | 調整のみ |

この安全モデルは以下を保証します：

- **ネットワークが安全上クリティカルになることはない**
- **オペレーターは常にオーバーライド権限を持つ**
- **ハードウェアフェイルセーフはローカルで実行される**
- **障害時は段階的に縮退する**
- **単一障害点がフリート全体に影響を及ぼすことはない**

---

## 関連ドキュメント

- [対応プラットフォーム]({{< relref "/fleet/platforms" >}}) — すべての車両タイプの概要
- [ドローンプラットフォーム]({{< relref "/fleet/platforms/drones" >}}) — ドローン固有の安全性の詳細
- [地上車両]({{< relref "/fleet/platforms/ground" >}}) — 地上車両の安全性の詳細

---

## 参考資料

- [PX4の安全設定](https://docs.px4.io/main/en/config/safety.html)
- [MAVLink共通メッセージセット](https://mavlink.io/en/messages/common.html)
- [NATSセキュリティモデル](https://docs.nats.io/running-a-nats-service/configuration/securing_nats)

---

## 概要に戻る

[← 自律型車両フリートアーキテクチャ]({{< relref "/fleet" >}})
