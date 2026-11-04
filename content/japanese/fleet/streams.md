---
title: "JetStream の設定"
meta_title: "ドローンフリートのデジタルツインのためのJetStreamストリーム | Ubuntu Software"
description: "フリートのデジタルツインに向けたJetStreamストリームとKVストアの構成：状態のロールアップ、イベント監査証跡、コマンドキュー、シャドウステート。"
image: "/images/robotics.svg"
draft: false
---

## デジタルツインのためのJetStream

JetStreamは、フリートのデジタルツインの永続化レイヤーを提供します。各ストリームタイプはそれぞれ特定の目的を持ち、目的に合わせた保持設定とレプリケーション設定が適用されています。

---

## ストリームアーキテクチャ

デジタルツインパターンは、5つのJetStreamリソースで構成されます。

| リソース | タイプ | 目的 |
|----------|------|---------|
| **TWIN_STATE** | Stream | ロールアップ付きテレメトリ |
| **TWIN_EVENTS** | Stream | 永続的なイベント監査証跡 |
| **TWIN_CMD** | Stream | コマンドキュー |
| **TWIN_SHADOW** | KV Store | 望ましい状態/報告された状態 |
| **TWIN_BLOBS** | Object Store | 大容量アーティファクト |

---

## TWIN_STATE: テレメトリストリーム

車両の状態を継続的に保存し、**サブジェクトベースのロールアップ**を行います。

```yaml
name: TWIN_STATE
subjects:
  - "fleet.prod.veh.*.state.>"
retention: limits
max_age: 1h
max_bytes: 10GB
storage: file
replicas: 3
discard: old
rollup_hdrs: true
```

### 設定の解説

| 設定 | 値 | 理由 |
|---------|-------|-----------|
| **retention** | limits | 経過時間とサイズで上限を設定 |
| **max_age** | 1h | リプレイ用に直近の状態を保持 |
| **max_bytes** | 10GB | ハブごとのストレージ容量を制限 |
| **replicas** | 3 | ノード障害に耐える |
| **rollup_hdrs** | true | サブジェクトごとのロールアップを有効化 |

### サブジェクトロールアップ

ロールアップを有効にすると、コンシューマーは**サブジェクトごとの最新メッセージのみ**を要求できます。

```go
// Get latest position for all vehicles
sub, _ := js.Subscribe("fleet.prod.veh.*.state.position",
    nats.DeliverLastPerSubject())
```

これにより、ストリームは大量のデータが流れ続けるファイアホースから、クエリ可能な状態ストアに変わります。

### コンシューマーパターン

| パターン | ユースケース |
|---------|----------|
| **DeliverAll** | 履歴のリプレイ、分析 |
| **DeliverLast** | 現在の状態のスナップショット |
| **DeliverLastPerSubject** | 車両ごとの最新状態 |
| **StartTime** | 特定のタイムスタンプからのリプレイ |

---

## TWIN_EVENTS: 監査証跡

コンプライアンスとデバッグのために、個々のイベントを**長期間保持**します。

```yaml
name: TWIN_EVENTS
subjects:
  - "fleet.prod.veh.*.evt.>"
retention: limits
max_age: 90d
max_bytes: 100GB
storage: file
replicas: 3
discard: old
```

### 設定の解説

| 設定 | 値 | 理由 |
|---------|-------|-----------|
| **max_age** | 90d | 規制対応、インシデント分析 |
| **max_bytes** | 100GB | イベントは状態よりもデータ量が小さい |
| **storage** | file | 再起動後も保持される |

### イベントの特性

イベントと状態には次の違いがあります。

| 特性 | 状態 | イベント |
|----------|-------|--------|
| **頻度** | 連続(Hz) | 離散(変化時) |
| **保持期間** | 短期(数時間) | 長期(数か月) |
| **コンシューマー** | リアルタイムダッシュボード | 監査、分析、アラート |
| **リプレイ** | 現在の状態の再構築 | 調査、コンプライアンス |

### イベントの例

```json
// Armed event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "previous_state": "disarmed",
    "armed_by": "operator-123",
    "reason": "mission_start"
  }
}

// Failsafe event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T11:45:00Z",
  "data": {
    "failsafe_type": "low_battery",
    "action_taken": "rtl",
    "battery_remaining": 0.15
  }
}
```

---

## TWIN_CMD: コマンドキュー

コマンドを**ワークキュー保持ポリシー**で保存します。

```yaml
name: TWIN_CMD
subjects:
  - "fleet.prod.veh.*.cmd.>"
retention: workqueue
max_age: 5m
storage: file
replicas: 3
```

### 設定の解説

| 設定 | 値 | 理由 |
|---------|-------|-----------|
| **retention** | workqueue | ack後にメッセージを削除 |
| **max_age** | 5m | ackされないコマンドは期限切れになる |
| **replicas** | 3 | コマンドの配信を保証 |

### ワークキューのセマンティクス

1. コマンドがストリームにパブリッシュされる
2. 車両ゲートウェイが受信する(exactly-once配信)
3. ゲートウェイがコマンドを処理する
4. ゲートウェイがメッセージをackする
5. メッセージがストリームから削除される

**ackされていないコマンドは再配信される**ため、接続が切れている間もコマンドは失われません。

### コマンドフロー

```
Fleet Operator               TWIN_CMD Stream              Vehicle Gateway
     │                             │                             │
     │ Publish: cmd.takeoff        │                             │
     │ ─────────────────────────▶  │                             │
     │                             │ Deliver to consumer         │
     │                             │ ─────────────────────────▶  │
     │                             │                             │ Process
     │                             │                             │ command
     │                             │             Ack             │
     │                             │ ◀─────────────────────────  │
     │                             │ (message removed)           │
     │                             │                             │
```

---

## TWIN_SHADOW: KVストア

JetStreamのKVストアは、デジタルツインのシャドウ状態を保持します。

```yaml
bucket: TWIN_SHADOW
max_value_size: 1MB
history: 5
ttl: 0
replicas: 3
storage: file
```

### 設定の解説

| 設定 | 値 | 理由 |
|---------|-------|-----------|
| **history** | 5 | 過去5バージョンを保持 |
| **ttl** | 0 | 有効期限なし(永続) |
| **replicas** | 3 | 高可用性 |

### キー構造

```
fleet/prod/veh/VID-001/desired    → Desired state JSON
fleet/prod/veh/VID-001/reported   → Reported state JSON
fleet/prod/veh/VID-001/config     → Vehicle configuration
fleet/prod/veh/VID-001/meta       → Metadata (serial, type, etc.)
```

### シャドウ操作

```go
// Write desired state
kv.Put("fleet/prod/veh/VID-001/desired", desiredJSON)

// Read reported state
entry, _ := kv.Get("fleet/prod/veh/VID-001/reported")

// Watch for changes
watcher, _ := kv.Watch("fleet/prod/veh/VID-001/desired")
for entry := range watcher.Updates() {
    // Desired state changed, reconcile
}
```

### デバッグのための履歴

履歴を有効にすると、過去の状態を確認できます。

```go
// Get all historical values
history, _ := kv.History("fleet/prod/veh/VID-001/desired")
for _, entry := range history {
    fmt.Printf("Rev %d at %s: %s\n",
        entry.Revision(), entry.Created(), entry.Value())
}
```

---

## TWIN_BLOBS: オブジェクトストア

メッセージに収まらない大容量アーティファクトを保存します。

```yaml
bucket: TWIN_BLOBS
max_chunk_size: 128KB
storage: file
replicas: 3
```

### ユースケース

| アーティファクト | サイズ | 例 |
|----------|------|---------|
| **ミッションファイル** | 10-100KB | ウェイポイント定義 |
| **ログファイル** | 1-10MB | フライトログ、診断情報 |
| **ファームウェア** | 10-50MB | PX4ファームウェアイメージ |
| **地図** | 1-100MB | オフライン地図タイル |
| **MLモデル** | 10-500MB | TensorRTエンジンファイル |

### オブジェクト操作

```go
// Upload mission file
obj, _ := obs.PutFile("missions/mission-456.json", file)

// Download to vehicle
reader, _ := obs.Get("missions/mission-456.json")

// List available firmware
objects := obs.List(nats.ObjectSearchPrefix("firmware/"))
```

### チャンク転送

オブジェクトストアは、大容量ファイルを次の手順で扱います。

1. チャンクに分割する(デフォルトは128KB)
2. チャンクをストリームメッセージとして保存する
3. メタデータ(サイズ、ハッシュ、チャンク)を記録する
4. 取得時に再構成する

これにより、不安定な接続でも**再開可能な転送**が可能になります。

---

## ストリームのサイジング

### 車両あたりの見積もり

| ストリーム | メッセージサイズ | レート | 1日あたりの量 |
|--------|--------------|------|--------------|
| **STATE** | 200バイト | 10 Hz | 約170 MB |
| **EVENTS** | 500バイト | 10件/日 | 約5 KB |
| **CMD** | 300バイト | 100件/日 | 約30 KB |

### フリート全体の見積もり(1,000台)

| リソース | 1時間 | 1日 | 90日 |
|----------|--------|-------|--------|
| **TWIN_STATE** | 7 GB | - | - |
| **TWIN_EVENTS** | - | 5 MB | 450 MB |
| **TWIN_SHADOW** | - | 1 GB | 1 GB |
| **TWIN_BLOBS** | - | - | 50 GB |

これらは基準となる見積もりです。実際の使用量は、メッセージのレートとサイズによって変わります。

---

## 運用上の考慮事項

### ストリームの健全性

次のメトリクスを監視します。

| メトリクス | アラートしきい値 |
|--------|-----------------|
| **コンシューマーラグ** | 10,000メッセージ超 |
| **ストリームのバイト数** | 上限の80%超 |
| **レプリカの同期** | 遅れているレプリカがある |
| **ack保留** | コンシューマーあたり1,000件超 |

### バックアップ戦略

| リソース | バックアップ方法 | 頻度 |
|----------|---------------|-----------|
| **ストリーム** | JetStreamスナップショット | 毎日 |
| **KVストア** | キーのエクスポート | 毎日 |
| **オブジェクトストア** | ファイルレベルのバックアップ | 毎週 |

### 移行

ストリームはクラスター間で移行できます。

1. 移行先に同じ設定でストリームを作成する
2. 元のストリームをソースとするミラーを作成する
3. 同期が完了するまで待つ
4. プロデューサー/コンシューマーを移行先に切り替える
5. 元のストリームを削除する

---

## まとめ

| リソース | タイプ | 保持期間 | 目的 |
|----------|------|-----------|---------|
| **TWIN_STATE** | Stream | 1時間 | ロールアップ付きテレメトリ |
| **TWIN_EVENTS** | Stream | 90日 | 監査証跡 |
| **TWIN_CMD** | Stream | Workqueue | コマンド配信 |
| **TWIN_SHADOW** | KV | 永続 | 望ましい状態/報告された状態 |
| **TWIN_BLOBS** | Object | 永続 | 大容量アーティファクト |

---

## 次へ

[車両ゲートウェイ →]({{< relref "/fleet/gateway" >}})
