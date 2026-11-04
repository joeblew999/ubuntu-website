---
title: "車両ゲートウェイ"
meta_title: "自律走行フリート向け車両ゲートウェイ設計 | Ubuntu Software"
description: "Vehicle Gatewayのアーキテクチャ：プロトコル変換（MAVLink、CAN、ROS 2）、ステートのダウンサンプリング、イベント抽出、コマンド実行、およびシャドウの照合。"
image: "/images/robotics.svg"
draft: false
---

## 車両プロトコルとNATSをつなぐブリッジ

Vehicle Gatewayは、各車両のエッジコンピュータ上で動作するサービスです。車両ネイティブのプロトコル（MAVLink、CANバス、ROS 2）と、フリート連携に使用されるNATSメッセージングシステムとの間で変換を行います。

---

## なぜゲートウェイが必要なのか？

車両プロトコルとフリートメッセージングは、それぞれ異なる目的を持っています：

| 観点 | 車両プロトコル | NATS |
|--------|------------------|------|
| **スコープ** | 単一車両 | フリート全体 |
| **プロトコル** | バイナリ、コンパクト | JSON、人間が読める形式 |
| **レート** | 100Hz以上のテレメトリ | WAN向けにダウンサンプリング |
| **永続化** | なし | JetStreamストリーム |
| **パターン** | リクエスト・レスポンス | 永続化付きPub/Sub |

Gatewayはこれら2つの世界を橋渡しし、プロトコル変換、レート制限、状態管理を担います。

---

## サポートされるプロトコル

| プロトコル | 車両タイプ | 統合方式 |
|----------|--------------|-------------|
| **MAVLink** | ドローン（PX4、ArduPilot） | UDP/シリアル、メッセージ解析 |
| **CANバス** | 乗用車、トラック | SocketCAN、DBCデコード |
| **J1939** | 大型車両 | PGN/SPNマッピング |
| **ROS 2** | 自律プラットフォーム | トピック購読 |

各プロトコルにはプロトコル固有のアダプターが必要ですが、ゲートウェイのコアロジック（ダウンサンプリング、イベント、コマンド、シャドウ）は共通です。

---

## 主な責務

### 1. プロトコルの取り込み

車両ネイティブのプロトコルメッセージを受信して解析します。MAVLink（ドローン）の例：

```go
// Receive MAVLink frames from mavlink-router
conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 14550})

for {
    buf := make([]byte, 1024)
    n, _, _ := conn.ReadFromUDP(buf)

    frame, _ := mavlink.Parse(buf[:n])

    switch msg := frame.Message().(type) {
    case *mavlink.Heartbeat:
        handleHeartbeat(msg)
    case *mavlink.GlobalPositionInt:
        handlePosition(msg)
    case *mavlink.BatteryStatus:
        handleBattery(msg)
    // ... handle other message types
    }
}
```

**主なメッセージタイプ（MAVLinkの例）：**

| MAVLinkメッセージ | 内容 |
|-----------------|---------|
| `HEARTBEAT` | モード、アーム状態、システムステータス |
| `GLOBAL_POSITION_INT` | 緯度、経度、高度、速度 |
| `ATTITUDE` | ロール、ピッチ、ヨー |
| `BATTERY_STATUS` | 電圧、電流、残量 |
| `GPS_RAW_INT` | GPS測位状態、衛星数、HDOP |
| `SYS_STATUS` | CPU負荷、エラー、ヘルス |

**CANバスでの相当項目（地上車両）：**

| CAN信号 | 内容 |
|------------|---------|
| エンジン回転数 | 現在のエンジン回転速度 |
| 車速 | ホイールセンサーによる対地速度 |
| GPS位置 | テレマティクスモジュールからの緯度/経度 |
| 燃料残量 | タンクの充填率 |
| DTCコード | 故障診断コード |

### 2. 状態のダウンサンプリング

WAN送信のためにテレメトリのレートを削減します：

```go
// Position: downsample from 10Hz to 1Hz
positionSampler := NewDownsampler(100*time.Millisecond, func(msg *Position) {
    // Average positions over window
    return averagePosition(msg)
})

// Battery: only publish on change
batterySampler := NewChangeFilter(func(old, new *Battery) bool {
    return math.Abs(old.Remaining - new.Remaining) > 0.01
})
```

**ダウンサンプリング戦略：**

| 戦略 | ユースケース | 例 |
|----------|----------|---------|
| **時間ベース** | 定期的な状態 | 1Hzでの位置 |
| **変化ベース** | 離散値 | モード変更 |
| **しきい値** | 緩やかな変化 | Δ > 1%のときのバッテリー |
| **集約** | 高頻度データ | ウィンドウ内の最小/最大/平均 |

フルレートのテレメトリは、認識システムのために車両内に留まります。WANを通過するのはダウンサンプリングされたデータのみです。

### 3. イベント抽出

状態遷移からイベントを生成します：

```go
type EventDetector struct {
    previousState *VehicleState
}

func (e *EventDetector) Process(current *VehicleState) []Event {
    var events []Event

    // Detect arm state change
    if current.Armed && !e.previousState.Armed {
        events = append(events, Event{
            Type:      "armed",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"reason": "manual"},
        })
    }

    // Detect mode change
    if current.Mode != e.previousState.Mode {
        events = append(events, Event{
            Type:      "mode_change",
            Timestamp: time.Now(),
            Data: map[string]interface{}{
                "from": e.previousState.Mode,
                "to":   current.Mode,
            },
        })
    }

    // Detect failsafe
    if current.FailsafeActive && !e.previousState.FailsafeActive {
        events = append(events, Event{
            Type:      "failsafe",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"type": current.FailsafeType},
        })
    }

    e.previousState = current
    return events
}
```

**検出されるイベント：**

| イベント | トリガー |
|-------|---------|
| `armed` | アーム状態 false → true |
| `disarmed` | アーム状態 true → false |
| `mode_change` | フライトモードの遷移 |
| `takeoff` | 飛行中状態 false → true |
| `landed` | 飛行中状態 true → false |
| `failsafe` | フェイルセーフの作動 |
| `geofence` | ジオフェンス違反の検出 |
| `battery.low` | バッテリーがしきい値を下回った |
| `battery.critical` | バッテリーが危険レベル |

### 4. ポリシーに基づくコマンド実行

NATSからコマンドを受信し、MAVLink経由で実行します：

```go
func (g *Gateway) handleCommand(cmd *Command) *CommandAck {
    // Validate command is allowed
    if !g.policy.Allows(cmd) {
        return &CommandAck{
            Status: "rejected",
            Error:  "command not allowed by policy",
        }
    }

    // Convert to MAVLink command
    mavCmd := cmd.ToMAVLink()

    // Send to flight controller
    if err := g.mavlink.Send(mavCmd); err != nil {
        return &CommandAck{
            Status: "failed",
            Error:  err.Error(),
        }
    }

    // Wait for MAVLink ACK
    mavAck, err := g.mavlink.WaitAck(mavCmd.CommandID, 5*time.Second)
    if err != nil {
        return &CommandAck{
            Status: "timeout",
            Error:  "no response from flight controller",
        }
    }

    return &CommandAck{
        Status: mavAck.Result.String(),
    }
}
```

**ポリシーの適用：**

| ポリシー | 説明 |
|--------|-------------|
| **ジオフェンス** | 境界外へのgotoコマンドを拒否 |
| **高度制限** | 高度コマンドの最大値を制限 |
| **モード制限** | 特定のモード遷移を禁止 |
| **レート制限** | コマンドの大量送信を防止 |
| **認証** | コマンド送信元を検証 |

### 5. シャドウの調整

目標状態を実際の車両状態と同期させます：

```go
func (g *Gateway) reconcileLoop() {
    ticker := time.NewTicker(1 * time.Second)

    for range ticker.C {
        // Read desired state from KV
        desired, _ := g.kv.Get(g.desiredKey())

        // Compare with actual state
        actual := g.getCurrentState()

        // Generate commands to reconcile
        commands := g.reconcile(desired, actual)

        for _, cmd := range commands {
            g.executeCommand(cmd)
        }

        // Update reported state
        reported := g.buildReportedState(actual)
        g.kv.Put(g.reportedKey(), reported)
    }
}

func (g *Gateway) reconcile(desired, actual *State) []*Command {
    var commands []*Command

    // Mode reconciliation
    if desired.Mode != actual.Mode {
        commands = append(commands, &Command{
            Type: "set_mode",
            Data: map[string]interface{}{"mode": desired.Mode},
        })
    }

    // Geofence reconciliation
    if desired.GeofenceEnabled != actual.GeofenceEnabled {
        commands = append(commands, &Command{
            Type: "set_geofence",
            Data: map[string]interface{}{"enabled": desired.GeofenceEnabled},
        })
    }

    return commands
}
```

**シャドウ状態により宣言的な管理が可能になります：**

- フリートオペレーターが目標状態を設定
- Gatewayが差分を検出
- Gatewayが収束させるためのコマンドを発行
- Gatewayが実際の状態を報告
- これを継続的に繰り返す

---

## アーキテクチャ

```
┌────────────────────────────────────────────────────────────────────┐
│                         Vehicle Gateway                             │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐     │
│  │   Protocol   │    │    State     │    │    NATS Client   │     │
│  │   Adapter    │───▶│   Machine    │───▶│    Publisher     │     │
│  └──────────────┘    └──────────────┘    └──────────────────┘     │
│         │                   │                      │               │
│         │            ┌──────▼──────┐               │               │
│         │            │   Event     │               │               │
│         │            │  Detector   │───────────────┤               │
│         │            └─────────────┘               │               │
│         │                                          │               │
│         │            ┌─────────────┐               │               │
│         │            │   Shadow    │◀──────────────┤               │
│         │            │ Reconciler  │               │               │
│         │            └──────┬──────┘               │               │
│         │                   │                      │               │
│         │            ┌──────▼──────┐    ┌─────────▼────────┐      │
│         │            │  Command    │    │   NATS Client    │      │
│         │◀───────────│  Executor   │◀───│   Subscriber     │      │
│         │            └─────────────┘    └──────────────────┘      │
│  ┌──────▼──────┐                                                   │
│  │   Protocol  │                                                   │
│  │   Sender    │                                                   │
│  └─────────────┘                                                   │
└────────────────────────────────────────────────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────┐                       ┌─────────────────┐
│  Vehicle Control│                       │   NATS Leaf     │
│ (Pixhawk/ECU/   │                       │   (localhost)   │
│  ROS 2 node)    │                       │                 │
└─────────────────┘                       └─────────────────┘
```

**プロトコルアダプター：**
- **MAVLink Adapter** — ドローン向け（PX4、ArduPilot）
- **CAN Adapter** — 乗用車/トラック向け（SocketCAN + DBC）
- **ROS 2 Adapter** — ROSベースのプラットフォーム向け

---

## 実装の詳細

### 依存関係

```go
import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "github.com/bluenviern/go-mavlink/v2"
)
```

**コアライブラリ：**

| ライブラリ | 用途 |
|---------|---------|
| **nats.go** | NATSクライアント、JetStream、KV |
| **go-mavlink** | MAVLinkプロトコルの実装 |
| **slog** | 構造化ロギング |

### 設定

```yaml
# gateway.yaml
vehicle_id: VID-001
environment: prod

mavlink:
  local_addr: ":14550"
  system_id: 1
  component_id: 1

nats:
  url: "nats://localhost:4222"
  credentials: "/etc/gateway/vehicle.creds"

sampling:
  position_hz: 1
  attitude_hz: 1
  battery_change_threshold: 0.01

policy:
  max_altitude: 120
  geofence_file: "/etc/gateway/geofence.json"
```

### デプロイ

```yaml
# systemd service
[Unit]
Description=Vehicle Gateway
After=network.target nats.service

[Service]
Type=simple
ExecStart=/usr/local/bin/vehicle-gateway --config /etc/gateway/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## エラー処理

### MAVLinkエラー

| エラー | 対応 |
|-------|----------|
| **接続断** | バックオフ付きで再接続 |
| **解析エラー** | ログに記録してフレームをスキップ |
| **タイムアウト** | 再試行またはコマンドを失敗扱い |

### NATSエラー

| エラー | 対応 |
|-------|----------|
| **接続断** | ローカルNATSは動作を継続し、ハブへ再接続 |
| **パブリッシュ失敗** | ローカルにバッファリングして再試行 |
| **ストリームエラー** | ログ記録、アラート、処理を継続 |

### コマンドエラー

| エラー | 対応 |
|-------|----------|
| **ポリシー違反** | 即座に拒否 |
| **MAVLinkによる拒否** | 失敗ACKを返す |
| **タイムアウト** | タイムアウトACKを返す |

---

## メトリクス

GatewayはPrometheusメトリクスを公開します：

```
# MAVLink
gateway_mavlink_messages_received_total{type="heartbeat"}
gateway_mavlink_messages_sent_total{type="command_long"}
gateway_mavlink_parse_errors_total

# NATS
gateway_nats_messages_published_total{subject="state.position"}
gateway_nats_messages_received_total{subject="cmd.takeoff"}
gateway_nats_publish_errors_total

# Commands
gateway_commands_received_total{type="takeoff"}
gateway_commands_executed_total{type="takeoff",result="success"}
gateway_command_latency_seconds{type="takeoff"}

# Shadow
gateway_shadow_reconciliations_total
gateway_shadow_commands_issued_total
```

---

## まとめ

| 責務 | 入力 | 出力 |
|----------------|-------|--------|
| **プロトコルの取り込み** | 車両メッセージ | 解析済みテレメトリ |
| **状態のダウンサンプリング** | 100Hzのテレメトリ | 1Hzの状態 |
| **イベント抽出** | 状態遷移 | 離散イベント |
| **コマンド実行** | NATSコマンド | 車両コマンド |
| **シャドウの調整** | 目標状態 | 収束コマンド |

Vehicle Gatewayは、基盤となる車両制御システムの安全性保証を維持しながら、フリート規模の運用を可能にする重要なコンポーネントです。

---

## 関連ドキュメント

- [サポートされるプラットフォーム]({{< relref "/fleet/platforms" >}}) — すべての車両タイプの概要
- [ドローンプラットフォーム]({{< relref "/fleet/platforms/drones" >}}) — MAVLink統合の詳細
- [地上車両]({{< relref "/fleet/platforms/ground" >}}) — CANバスおよびJ1939の統合

---

## 次へ

[安全モデル →]({{< relref "/fleet/safety" >}})
