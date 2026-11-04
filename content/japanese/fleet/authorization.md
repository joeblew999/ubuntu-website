---
title: "認可とグラント"
meta_title: "分散型認可とアクセス許可 | Ubuntu Software"
description: "NATS JetStreamのグラントを用いた分散型セキュリティモデルにより、サードパーティは中央での調整を必要とせずに、リアルタイムシステムへ安全に参加できます。"
image: "/images/robotics.svg"
draft: false
---

## ボトルネックのないセキュリティ

従来の集中型認可は、大規模環境で問題を引き起こします。

| 問題 | 影響 |
|---------|--------|
| **単一障害点** | 認証サーバーのダウン = システム全体のダウン |
| **リクエストごとのレイテンシ** | 認証サーバーへのラウンドトリップが遅延を生む |
| **切断時に動作できない** | ネットワークなし = 認可なし |
| **スケールしない** | 数千台のデバイスが中央サーバーを圧倒する |

NATS は異なるアプローチを採用しています。自己完結型の認証情報による**分散型認可**です。

---

## NATS の認可の仕組み

```
┌─────────────────────────────────────────────────────────────────┐
│                    CREDENTIAL ISSUANCE                          │
│   Account Authority issues signed JWTs with embedded perms      │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌─────────┐     ┌─────────┐     ┌─────────┐
        │ Vehicle │     │ Operator│     │ Partner │
        │  Creds  │     │  Creds  │     │  Creds  │
        └────┬────┘     └────┬────┘     └────┬────┘
             │               │               │
             ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS SERVER                                │
│   Validates JWT signature locally — no external call needed     │
│   Enforces permissions from embedded claims                     │
└─────────────────────────────────────────────────────────────────┘
```

**主な特性:**

- **ローカル検証** — サーバーは外部サービスを呼び出すことなく認証情報を検証します
- **埋め込まれた権限** — 実行できる操作は認証情報にエンコードされています
- **即時失効** — アカウントレベルの失効リストがクラスター全体に伝播します
- **ゼロトラスト** — 認証情報がなければアクセスは一切不可です

これは切断中でも機能します。認可のためにホームへ問い合わせる必要はありません。認証情報そのものが証明となります。

---

## 権限スコープ

認証情報は、各当事者が実行できる操作を正確に指定します。

| スコープ | 説明 | 例 |
|-------|-------------|---------|
| **Publish** | サブジェクトへメッセージを送信 | 車両が `fleet.prod.veh.V001.state.pos` をパブリッシュ |
| **Subscribe** | サブジェクトからメッセージを受信 | オペレーターが `fleet.prod.veh.*.state.*` をサブスクライブ |
| **Publish + Subscribe** | 双方向 | ゲートウェイがコマンド/レスポンスを処理 |
| **Deny** | パターンを明示的にブロック | 車両が他の車両をサブスクライブすることをブロック |

**権限におけるサブジェクトのワイルドカード:**

```
fleet.prod.veh.V001.*        → Only this vehicle's subjects
fleet.prod.veh.*.state.*     → All vehicles' state (read-only)
fleet.prod.veh.*.cmd.*       → All vehicles' commands (operator)
fleet.prod.veh.*.evt.*       → All events (auditor)
```

これにより、車両が自身のデータのみをパブリッシュでき、他の車両になりすますことができないよう強制しています。完全な階層構造については [サブジェクト命名]({{< relref "/fleet/subjects" >}}) を参照してください。

---

## グラントモデル

グラントは**きめ細かく、委任可能な権限**を実現します。オペレーターは、システム全体へのアクセスを与えることなく、第三者に特定のデータへのアクセスを付与できます。

```
Account Authority (root)
    │
    ├── Fleet Operator Account
    │       │
    │       ├── Grant: Emergency Services (incident-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.pos
    │       │       └── Subscribe: fleet.prod.veh.*.stream.video
    │       │       └── Expires: 4 hours
    │       │
    │       ├── Grant: Research Partner (data-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.env.*
    │       │       └── Deny: fleet.prod.veh.*.state.pos
    │       │       └── Expires: 30 days
    │       │
    │       └── Grant: Auditor (read-only, full)
    │               └── Subscribe: fleet.prod.veh.*.state.*
    │               └── Subscribe: fleet.prod.veh.*.evt.*
    │               └── Deny: fleet.prod.veh.*.cmd.*
    │               └── Expires: 7 days
```

**グラントの特性:**

| 特性 | メリット |
|----------|---------|
| **スコープ指定** | 特定のサブジェクト/パターンのみ — それ以上は不可 |
| **期間限定** | 自動失効により、放置されたアクセスが残らない |
| **失効可能** | 付与先に手を加えることなく即時失効 |
| **委任可能** | オペレーターが他者にサブ権限を付与できる |
| **監査対応** | すべてのグラント発行が記録される |

---

## 第三者連携のシナリオ

### 緊急サービスとの連携

```
Scenario: Wildfire response, fire department needs drone video

Grant issued:
  - Subject: fleet.prod.veh.*.stream.video
  - Subject: fleet.prod.veh.*.state.pos
  - Duration: 4 hours (incident duration)
  - Audit: All access logged with responder ID

Result: Fire department connects with their credential, receives
real-time video and position from all drones in area.
Grant expires automatically when incident closes.
```

**これが機能する理由:**
- 消防署は必要なもの — ライブ映像と位置情報 — を過不足なく取得できます
- コマンドへのアクセスはなく、ドローンの運用を妨げることはできません
- 期間限定 — インシデントの終了とともにアクセスも終了します
- 監査可能 — 受信したすべてのフレームが記録されます

### パートナーとのデータ共有

```
Scenario: University research on environmental monitoring

Grant issued:
  - Subject: fleet.prod.veh.*.state.env.* (temp, humidity, air quality)
  - Deny: fleet.prod.veh.*.state.pos (no location data)
  - Duration: 30 days (research period)
  - Audit: Data access logged for compliance

Result: Researcher receives environmental sensor data without
knowing vehicle locations—privacy preserved.
```

**これが機能する理由:**
- 研究者は研究に必要な環境データを取得できます
- 位置データは明示的に拒否され、プライバシーが保護されます
- 研究コンプライアンスのための明確な監査証跡があります
- 研究期間終了後に自動でクリーンアップされます

### 規制コンプライアンス

```
Scenario: Aviation authority audit of flight operations

Grant issued:
  - Subject: fleet.prod.veh.*.state.* (all telemetry)
  - Subject: fleet.prod.veh.*.evt.* (all events)
  - Deny: fleet.prod.veh.*.cmd.* (no command capability)
  - Duration: 7 days (audit period)
  - Audit: Every query logged for audit trail

Result: Auditor gets complete read access to historical data
but cannot issue any commands to vehicles.
```

**これが機能する理由:**
- 監査人に対して完全な透明性 — 隠し事はありません
- 読み取り専用 — 監査人は運用に影響を与えられません
- 監査人のアクセスに関する完全な監査証跡があります
- 期間限定 — 恒久的なアクセスの抜け穴は残りません

### 顧客との連携

```
Scenario: Enterprise customer wants API access to their fleet data

Grant issued:
  - Subject: fleet.prod.veh.CUST-*.state.* (only their vehicles)
  - Subject: fleet.prod.veh.CUST-*.evt.*
  - Deny: fleet.prod.veh.OTHER-*.* (other customers' data)
  - Duration: Persistent (until revoked)
  - Audit: API access logged

Result: Customer builds their own dashboards using real-time
data from their vehicles—and only their vehicles.
```

### 下請け業者のアクセス

```
Scenario: Maintenance team needs to run diagnostics on specific vehicles

Grant issued:
  - Subject: fleet.prod.veh.V001.cmd.diag.* (diagnostic commands only)
  - Subject: fleet.prod.veh.V001.state.health.*
  - Deny: fleet.prod.veh.V001.cmd.flight.* (no flight commands)
  - Duration: 8 hours (maintenance window)
  - Audit: All commands logged with technician ID

Result: Technician can run diagnostics but cannot arm or fly the vehicle.
Access expires after maintenance window.
```

---

## セキュリティ特性

| 特性 | 実装 |
|----------|----------------|
| **ゼロトラスト** | 暗黙の信頼なし — すべての接続を認証 |
| **最小権限** | 認証情報には必要最小限の権限のみを含む |
| **多層防御** | 複数のレイヤー: TLS + JWT + サブジェクト ACL |
| **オフライン対応** | ネットワーク接続なしで認可が機能 |
| **監査対応** | すべてのアクセスを認証情報の ID とともに記録 |
| **失効可能** | アカウント失効リストによる即時失効 |

---

## 運用上のメリット

### フリートオペレーター向け

- **システムを変更せずに第三者を追加/削除** — 認証情報を発行または失効させるだけで、設定ファイルの編集は不要です
- **認証情報を共有せずに一時的なアクセスを付与** — 各当事者が独自の認証情報を取得します
- **誰がいつ何にアクセスしたかを監査** — 完全なアクセスログがあります
- **必要に応じて即座にアクセスを失効** — 次回の接続時に反映されます

### 第三者向け

- **セルフサービスでの接続** — 提供された認証情報を使い、準備ができたら接続できます
- **明確なスコープ** — アクセスできるもの、できないものが正確にわかります
- **内部アーキテクチャの知識は不要** — 接続してサブスクライブするだけです
- **自動失効** — 自分の側でクリーンアップする必要はありません

### コンプライアンス向け

- **完全な監査証跡** — すべてのアクセスが認証情報まで追跡できます
- **証明可能な権限境界** — 誰が何にアクセスできたかを実証できます
- **監査人向けの期間限定アクセス** — 恒久的なアクセスの抜け穴は残りません
- **共有認証情報なし** — 個人ごとの説明責任が確保されます

---

## 実装に関する注意事項

### 認証情報の形式

NATS は **NKeys** と **JWT** を使用します。

- **Ed25519 署名** — 暗号学的に安全で、高速に検証できます
- **自己完結型のクレーム** — 権限がトークンに埋め込まれています
- **外部検証不要** — サーバーがローカルで署名を検証します

### 失効

アクセスを失効させる必要がある場合:

1. 認証情報をアカウント失効リストに追加します
2. リストがクラスター内のすべてのサーバーに伝播します
3. 次回の接続試行時に反映されます
4. 付与先のシステムに手を加える必要はありません

### 認証情報のローテーション

長期間にわたる連携の場合:

1. 古い認証情報の有効期限が切れる前に新しい認証情報を発行します
2. 重複期間を設けることで、シームレスな移行が可能になります
3. 移行完了後に古い認証情報を失効させます
4. サービスの中断はありません

---

## 実システムにおいてこれが重要な理由

従来の認可には以下が必要です。
- 中央認証サーバー(単一障害点)
- リクエストごとのネットワーク呼び出し(レイテンシ)
- 第三者向けの複雑な連携(数か月の作業)
- 期限切れアクセスの手動クリーンアップ(放置された認証情報)

NATS の分散型認可が提供するもの:
- **レジリエンス** — オフラインで動作し、中央のボトルネックがありません
- **スピード** — ローカル検証により、ラウンドトリップがありません
- **拡張性** — 第三者の追加は数か月ではなく数分で完了します
- **安全性** — 自動失効と即時失効

これこそが、**高いセキュリティ**と**高い使いやすさ**を両立したシステムを構築する方法です。

---

## 関連ドキュメント

- [サブジェクト命名]({{< relref "/fleet/subjects" >}}) — サブジェクト階層がきめ細かな ACL をどのように実現するか
- [NATS トポロジー]({{< relref "/fleet/nats-topology" >}}) — メッセージング基盤がどのように構成されているか
- [安全モデル]({{< relref "/fleet/safety" >}}) — コマンド認可が安全階層にどのように組み込まれるか
- [監視]({{< relref "/company/surveillance" >}}) — グラントが監視のための第三者アクセスをどのように実現するか

---

## 次へ

[ストリームとイベント →]({{< relref "/fleet/streams" >}})
