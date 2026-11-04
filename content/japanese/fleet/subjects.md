---
title: "サブジェクトの命名"
meta_title: "ドローンフリートのためのNATSサブジェクト命名規則 | Ubuntu Software"
description: "フリートのテレメトリ、イベント、コマンド、およびデジタルツインの状態管理のための階層型NATSサブジェクト命名規則。"
image: "/images/robotics.svg"
draft: false
---

## サブジェクト階層の設計

サブジェクトの命名は、メッセージングシステムのAPIそのものです。適切に設計された階層は、開発者にとって直感的でありながら、効率的なルーティング、フィルタリング、アクセス制御を可能にします。

---

## 設計原則

| 原則 | 理由 |
|-----------|-----------|
| **階層的** | 任意のレベルでワイルドカードサブスクリプションを可能にする |
| **予測可能** | 開発者が調べることなくサブジェクトを組み立てられる |
| **環境を意識** | 本番環境とステージング/開発環境を分離する |
| **フィルタリング可能** | 1台の機体、グループ、またはフリート全体をサブスクライブできる |
| **拡張可能** | 構造を変更せずに新しいメッセージタイプを追加できる |

---

## コアとなるサブジェクトパターン

すべてのフリートサブジェクトは次のパターンに従います:

```
fleet.<env>.veh.<vid>.<category>.<type>
```

| セグメント | 説明 | 例 |
|---------|-------------|----------|
| `fleet` | ルート名前空間 | 常に「fleet」 |
| `<env>` | 環境 | `prod`、`staging`、`dev` |
| `veh` | リソースタイプ | 機体(将来: 地上管制用の`gcs`) |
| `<vid>` | 機体ID | `VID-001`、`drone-alpha-7` |
| `<category>` | メッセージカテゴリ | `state`、`evt`、`cmd`、`cmdack` |
| `<type>` | 具体的なタイプ | `position`、`battery`、`takeoff` |

---

## メッセージカテゴリ

### 状態: `fleet.<env>.veh.<vid>.state.*`

機体の現在の状態。継続的にパブリッシュされます:

| サブジェクト | 説明 | レート |
|---------|-------------|------|
| `state.position` | 緯度/経度/高度、速度 | 10 Hz |
| `state.attitude` | ロール/ピッチ/ヨー | 10 Hz |
| `state.battery` | 電圧、電流、残量 | 1 Hz |
| `state.mode` | フライトモード(AUTO、MANUAL、RTL) | 変更時 |
| `state.health` | システムヘルスの概要 | 1 Hz |

**パブリッシュの例:**

```go
subject := fmt.Sprintf("fleet.prod.veh.%s.state.position", vehicleID)
nc.Publish(subject, positionJSON)
```

### イベント: `fleet.<env>.veh.<vid>.evt.*`

個別の出来事。発生時にパブリッシュされます:

| サブジェクト | 説明 |
|---------|-------------|
| `evt.armed` | 機体がアームされた |
| `evt.disarmed` | 機体がディスアームされた |
| `evt.takeoff` | 離陸を開始した |
| `evt.landed` | 着陸が完了した |
| `evt.waypoint` | ウェイポイントに到達した |
| `evt.failsafe` | フェイルセーフが作動した |
| `evt.geofence` | ジオフェンス違反 |
| `evt.battery.low` | バッテリー残量低下の警告 |

**イベントは状態遷移であり**、継続的なデータではありません。イベントは永続化され、監査証跡としてJetStreamに保存されます。

### コマンド: `fleet.<env>.veh.<vid>.cmd.*`

機体に送信される指示:

| サブジェクト | 説明 |
|---------|-------------|
| `cmd.arm` | モーターをアームする |
| `cmd.disarm` | モーターをディスアームする |
| `cmd.takeoff` | 離陸を開始する |
| `cmd.land` | 着陸を開始する |
| `cmd.rtl` | 離陸地点に帰還する |
| `cmd.goto` | 指定位置へ移動する |
| `cmd.mission.upload` | ミッションをアップロードする |
| `cmd.mission.start` | ミッションを開始する |
| `cmd.param.set` | パラメータを設定する |

**コマンドはダウンストリームに流れます**。つまり、フリート管理から機体へと送られます。

### コマンド確認応答: `fleet.<env>.veh.<vid>.cmdack.*`

コマンドへの応答:

| サブジェクト | 説明 |
|---------|-------------|
| `cmdack.arm` | アームコマンドの結果 |
| `cmdack.takeoff` | 離陸コマンドの結果 |
| `cmdack.goto` | Gotoコマンドの結果 |

**Ackペイロードには以下が含まれます:**

```json
{
  "cmd_id": "cmd-123456",
  "status": "accepted|rejected|completed|failed",
  "error": "optional error message",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## ワイルドカードサブスクリプション

NATSのワイルドカードにより、柔軟なサブスクリプションが可能になります:

| ワイルドカード | 意味 | 例 |
|----------|---------|---------|
| `*` | 単一トークン | `fleet.prod.veh.*.state.position` |
| `>` | 1つ以上のトークン | `fleet.prod.veh.VID-001.>` |

### 一般的なサブスクリプションパターン

| パターン | ユースケース |
|---------|----------|
| `fleet.prod.veh.*.state.position` | すべての機体の位置 |
| `fleet.prod.veh.VID-001.>` | 1台の機体からのすべてのメッセージ |
| `fleet.prod.veh.*.evt.>` | すべての機体からのすべてのイベント |
| `fleet.prod.veh.*.state.battery` | すべてのバッテリー状態 |
| `fleet.*.veh.VID-001.>` | すべての環境にわたる1台の機体 |

---

## KVストア: デジタルツインのシャドウ

JetStream KVストアは、**desired(期待)**状態と**reported(報告)**状態を保持します:

### サブジェクトパターン(KV)

```
fleet/<env>/veh/<vid>/desired
fleet/<env>/veh/<vid>/reported
```

注: KVは`/`区切り(バケットパス)を使用し、pub/subは`.`区切りを使用します。

### Desired状態

フリート管理システムが機体に**させたい**こと:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "geofence_enabled": true,
  "max_altitude": 120,
  "home_position": {
    "lat": 37.7749,
    "lon": -122.4194
  }
}
```

### Reported状態

機体が**実際に**報告する内容:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "mission_progress": 0.45,
  "geofence_enabled": true,
  "armed": true,
  "in_flight": true,
  "battery_remaining": 0.72
}
```

### シャドウの調整

Vehicle Gatewayは継続的に以下を行います:

1. KVからdesired状態を読み取る
2. 実際の機体の状態と比較する
3. 現実をdesired状態に合わせるためのコマンドを発行する
4. reported状態を更新する

これにより**宣言的なフリート管理**が可能になります。desired状態を設定すれば、機体がその状態に収束します。

---

## アクセス制御

NATSの認可は、ユーザーを許可されたサブジェクトにマッピングします:

```
# Vehicle can publish its own state/events
publish: fleet.prod.veh.VID-001.state.>
publish: fleet.prod.veh.VID-001.evt.>
publish: fleet.prod.veh.VID-001.cmdack.>

# Vehicle can subscribe to commands for itself
subscribe: fleet.prod.veh.VID-001.cmd.>

# Vehicle can read/write its own shadow
kv: fleet/prod/veh/VID-001/*
```

```
# Fleet operator can subscribe to all state/events
subscribe: fleet.prod.veh.*.state.>
subscribe: fleet.prod.veh.*.evt.>
subscribe: fleet.prod.veh.*.cmdack.>

# Fleet operator can publish commands
publish: fleet.prod.veh.*.cmd.>

# Fleet operator can manage shadows
kv: fleet/prod/veh/*/desired
```

これにより、**機体が他の機体になりすましたり**、許可されていないコマンドをパブリッシュしたりできないことが保証されます。

サードパーティへのグラントや認証情報の管理を含む認可モデルの詳細については、[認可とグラント]({{< relref "/fleet/authorization" >}})を参照してください。

---

## サブジェクト設計のトレードオフ

### フラットなサブジェクトを使わない理由

`vehicle-001-position`のようなフラットなサブジェクトは、以下をサポートしません:

- ワイルドカードサブスクリプション
- 階層的なアクセス制御
- 論理的なグループ化

### より深い階層を使わない理由

`fleet.prod.region.us-west.site.warehouse-a.veh.VID-001.state.position`のようなサブジェクトには、次の問題があります:

- プログラムで組み立てるのが難しい
- ワイルドカードが複雑になる
- 変更がシステム全体に波及する

4レベルの階層は、柔軟性とシンプルさのバランスを取っています。

---

## メッセージフォーマット

すべてのメッセージは、一貫した構造を持つ**JSON**を使用します:

```json
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00.123Z",
  "seq": 12345,
  "data": {
    // message-specific payload
  }
}
```

| フィールド | 説明 |
|-------|-------------|
| `vid` | 機体ID(サブジェクトと重複するが、コンシューマーにとって有用) |
| `timestamp` | ミリ秒精度のISO 8601 |
| `seq` | 順序付けのための単調増加シーケンス番号 |
| `data` | メッセージ固有のペイロード |

---

## まとめ

| カテゴリ | サブジェクトパターン | 方向 |
|----------|-----------------|-----------|
| **状態** | `fleet.<env>.veh.<vid>.state.<type>` | 機体 → ハブ |
| **イベント** | `fleet.<env>.veh.<vid>.evt.<type>` | 機体 → ハブ |
| **コマンド** | `fleet.<env>.veh.<vid>.cmd.<type>` | ハブ → 機体 |
| **Ack** | `fleet.<env>.veh.<vid>.cmdack.<type>` | 機体 → ハブ |
| **Desired** | `fleet/<env>/veh/<vid>/desired` (KV) | ハブ → 機体 |
| **Reported** | `fleet/<env>/veh/<vid>/reported` (KV) | 機体 → ハブ |

---

## 次へ

[ストリーム設定 →]({{< relref "/fleet/streams" >}})
