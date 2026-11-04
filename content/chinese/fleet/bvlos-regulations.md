---
title: "BVLOS 法规"
meta_title: "超视距无人机法规 | Ubuntu Software"
description: "全球各国BVLOS（超视距）无人机法规。了解商业机队运营的监管要求。"
image: "/images/robotics.svg"
draft: false
---

## 为什么 BVLOS 对机队运营至关重要

超视距（BVLOS）法规从根本上决定了商用无人机机队能够完成哪些任务。了解监管格局对于机队平台设计和市场战略至关重要。

| 挑战 | 对机队平台的影响 |
|-----------|-------------------------|
| **运营限制** | VLOS 限制（400-500 米范围）制约了机队的实用性。BVLOS 使基础设施巡检、配送和农业作业得以规模化开展。 |
| **合规集成** | 机队软件必须执行监管边界——远程识别（Remote ID）、探测与避让以及运营授权要求决定了遥测和安全功能的设计。 |
| **市场时机** | 各国在 BVLOS 开放方面处于不同阶段。了解哪些地区的法规正在放开，有助于确定目标市场的优先级。 |
| **合作伙伴就绪度** | 机队运营商需要获得 BVLOS 批准才能执行实际任务。跟踪哪些运营商已获批准，有助于识别潜在客户和合作伙伴。 |

---

## 全球趋势

在所有主要航空管理机构中，正在显现四大趋势：

| 趋势 | 说明 |
|-------|-------------|
| **采用 SORA** | FAA、EASA、CASA 和加拿大交通部均采用特定运行风险评估（SORA）方法 [[1]](#ref-1) |
| **远程识别强制要求** | 正在成为所有司法管辖区的强制要求 |
| **豁免 → 规则** | 从逐案审批豁免过渡到标准化的基于性能的规则 |
| **UTM 集成** | 各国正在开发无人机交通管理（UTM）系统，以实现可扩展的 BVLOS |

---

## 美洲

### 美国（FAA）

**现状：** 基于豁免，正在转向基于规则（Part 108） [[2]](#ref-2)

**时间线** [[3]](#ref-3)：

- 截至 2024 年 10 月已颁发 190 项 BVLOS 豁免（2020 年仅为 6 项）
- Part 108 拟议规则预计于 2026 年 1 月发布
- 《2024 年 FAA 重新授权法案》将 BEYOND 计划延长至 2029 年 [[4]](#ref-4)

**获得 BVLOS 批准的运营商：**

| 运营商 | 用例 | 批准日期 | 参考 |
|----------|----------|---------------|-----|
| DroneUp | 医疗配送 | 2024 年 1 月 | [[5]](#ref-5) |
| Amazon Prime Air | 消费者配送 | 2024 年 5 月 | [[6]](#ref-6) |
| Asylon | 安防/监控 | 2024 年 4 月 | [[7]](#ref-7) |
| 多家（DFW 地区） | 首个多运营商共享空域 | 2024 年 7 月 | [[8]](#ref-8) |

---

### 加拿大（Transport Canada）

**现状：** 需要 SFOC，正在转向常规化 BVLOS [[9]](#ref-9)

**框架：**

- 2025 年 4 月：在低风险条件下（≤150 公斤、非管制空域、人口稀少地区）无需 SFOC 即可开展常规 BVLOS 运行
- 需要高级飞行员证书
- 到 2030 年通过 Nav Canada 建立国家 UTM 框架

**运营商：** 待定

---

## 欧洲

### 欧盟（EASA）

**现状：** 特定类别，采用 SORA/PDRA 路径 [[10]](#ref-10)

**框架** [[1]](#ref-1)：

| 方面 | 详情 |
|--------|---------|
| 标准场景 | STS-02：配备空域观察员的 BVLOS |
| 远程识别 | 自 2024 年 1 月起强制执行 [[11]](#ref-11) |
| 审批周期 | 4-9 个月（医疗），6-12 个月（消费者） [[12]](#ref-12) |
| 费用 | €2,000-€15,000 |

**运营商：** 多家运营商通过各国监管沙盒和成员国快速通道计划获批。

---

### 英国（CAA）

**现状：** 通过英国 SORA 获得运营授权（OA）

**框架：**

- BVLOS 需要 CAA 运营授权
- 英国 SORA 流程（自 2025 年 4 月起）
- 通常需要探测与避让技术（ADS-B、FLARM、地面雷达）

**运营商：**

| 运营商 | 成就 | 参考 |
|----------|-------------|-----|
| sees.ai | 英国首个常规化 BVLOS（2021 年），National Grid 巡检（2023 年） | [[13]](#ref-13) |

---

### 爱尔兰（IAA）

**现状：** PDRA G-03 路径

**框架：** 为通过 EASA 获得简化的欧盟授权打开了大门。

**运营商：**

| 运营商 | 成就 | 参考 |
|----------|-------------|-----|
| sees.ai | 爱尔兰首个 BVLOS 许可（2025 年） | [[13]](#ref-13) |

---

## 亚太地区

### 澳大利亚（CASA）

**现状：** 依据 Part 101 进行基于 SORA 的评估 [[14]](#ref-14)

**框架** [[15]](#ref-15)：

- 采用 SORA 2.5 及定量风险建模
- 提供广域 BVLOS 自我评估框架
- 简化的流程正在吸引全球企业 [[16]](#ref-16)

**运营商：**

| 运营商 | 成就 | 参考 |
|----------|-------------|-----|
| Sphere Drones | 广域 BVLOS 批准（2024 年 5 月） | [[17]](#ref-17) |

---

### 日本（MLIT/JCAB）

**现状：** 允许在人口密集区上空进行 4 级 BVLOS [[18]](#ref-18)

**框架** [[19]](#ref-19)：

| 等级 | 说明 | 所需执照 |
|-------|-------------|------------------|
| 1-2 级 | VLOS 运行 | 二级 |
| 3 级 | 无人区上空 BVLOS | 二级 |
| 4 级 | 城市/人口密集区上空 BVLOS | 一级 |

- 执照有效期 3 年，需通过注册无人机学校培训 + 飞行考试
- 截至 2024 年 10 月已登记 422,879 架无人机

**运营商：**

- 截至 2024 年 10 月，仅有 1 款无人机型号获得一级（4 级）认证
- 5 款无人机型号获得二级认证

---

### 新加坡（CAAS）

**现状：** 未取得运营商许可 + 1 类活动许可，禁止 BVLOS [[20]](#ref-20)

**框架** [[21]](#ref-21)：

- 任何 BVLOS 均需严格授权
- 超过 25 公斤或 BVLOS 的无人机需要完整的运营商许可 + 1 类许可
- 自 2025 年 12 月起，超过 250 克的无人机强制使用 B-RID（广播式远程识别）
- B-RID 为未来规模化 BVLOS/自动化运行奠定基础
- 提供监管沙盒计划以实现快速审批

**运营商：** 待定（沙盒试验进行中）

---

### 中国（CAAC）

**现状：** BVLOS 需要特别授权 [[22]](#ref-22)

**框架** [[23]](#ref-23)：

- 《无人驾驶航空器飞行管理暂行条例》（2024 年 1 月）
- BVLOS 需要高级运行证书 + 特殊飞行运行证书（SFOC）
- 全国高度限制为 120 米；超出需获得授权
- 简化的数字化申请流程使更多 BVLOS 飞行成为可能 [[16]](#ref-16)

**运营商：** 待定（CAAC 框架下的国内运营商）

---

### 印度（DGCA）

**现状：** 尚无正式 BVLOS 框架；仅处于试验阶段 [[24]](#ref-24)

**框架** [[25]](#ref-25)：

- 依据《2021 年无人机规则》运行，并采用基于项目的豁免
- 2025 年 4 月：针对城市/半城市地区引入单独的 BVLOS 飞行员执照类别
- DGCA“正处于 BVLOS 规则最终确定的后期阶段” [[26]](#ref-26)
- 实验走廊已投入使用（例如特伦甘纳邦的“Medicine from the Sky”项目）

**运营商：**

| 运营商 | 成就 | 参考 |
|----------|-------------|-----|
| Skye Air Mobility | 全球最长 BVLOS 药品配送（104 公里，Baruipur-Medinipur） | [[27]](#ref-27) |
| Flipkart Health | 104 公里配送试验合作伙伴 | — |

---

## 合规要求

在所有司法管辖区，BVLOS 法规都趋向于一套共同的技术要求。本节提炼出我们的机队平台必须实现的内容。

### 1. 远程识别（全球强制）

所有主要监管机构都要求无人机广播身份和位置数据。

| 司法管辖区 | 要求 | 截止日期 | 参考 |
|--------------|-------------|----------|-----|
| 美国（FAA） | 标准远程识别或广播模块 | 2023 年 9 月（已生效） | [[2]](#ref-2) |
| 欧盟（EASA） | 特定类别 + 1 类及以上需远程识别 | 2024 年 1 月（已生效） | [[11]](#ref-11) |
| 新加坡 | 超过 250 克的无人机需 B-RID | 2025 年 12 月 | [[21]](#ref-21) |
| 澳大利亚 | 远程识别正在考虑中 | 待定 | [[14]](#ref-14) |

**平台要求：**
- 广播模块硬件集成（兼容 OpenDroneID）
- 通过 NATS 实现网络远程识别能力，用于全机队跟踪
- 将所有远程识别传输记录到 JetStream 以供审计

---

### 2. 探测与避让（DAA）

BVLOS 运行要求具备探测并避让其他航空器和障碍物的能力。

| 司法管辖区 | 要求 | 备注 |
|--------------|-------------|-------|
| 美国（FAA） | “保持充分间隔”标准；Part 108 要求配备 DAA 系统 | 基于性能，而非规定性要求 |
| 英国（CAA） | 通常需要 DAA 技术（ADS-B、FLARM、地面雷达） | 在 OA 流程中逐案审定 |
| 日本 | 4 级城市运行必需 | [[18]](#ref-18) |

**平台要求：**
- ADS-B In 接收机集成（交通态势感知）
- 在 Jetson 上进行传感器融合（可选摄像头、雷达、激光雷达）
- 自主避让机动或向操作员发出警报
- 将所有 DAA 事件记录到 JetStream

---

### 3. 地理围栏

通过硬件强制执行的边界，防止飞入限制区域。

| 司法管辖区 | 要求 | 备注 |
|--------------|-------------|-------|
| 美国（FAA） | 地理围栏是运营批准的一部分 | 在豁免/Part 108 中定义 |
| 欧盟（EASA） | 特定类别要求具备地理感知能力 | 依据 SORA 地面风险评估 |
| 全部 | 高度限制（通常为离地 120 米/400 英尺） | 所有框架的基本要求 |

**平台要求：**
- 在飞控上执行地理围栏（PX4 `GF_*` 参数）
- 通过 Gateway 动态更新地理围栏（临时限制）
- 执行高度上限
- 越界记录并向机队发出警报

---

### 4. 指挥与控制（C2）链路

操作员与航空器之间可靠、冗余的通信。

| 司法管辖区 | 要求 | 备注 |
|--------------|-------------|-------|
| 美国（FAA） | Part 108 中的 C2 链路性能标准 | 延迟、可用性、安全性 |
| 欧盟（EASA） | 依据 SORA 的 C2 链路可靠性 | 要求取决于风险等级 |
| 全部 | C2 链路丢失时的失效保护行为 | 自动返航、降落或悬停盘旋 |

**平台要求：**
- 通过 LTE/5G 或专用无线电作为主 C2 链路
- 备用 C2 路径（例如用于远程作业的卫星链路）
- 链路丢失时自动失效保护（在 PX4 中配置 RTL）
- 向机队发送 C2 链路状态遥测
- 记录链路质量以供审计

---

### 5. 运营授权文件

监管机构要求提供书面的风险评估和运行程序。

| 框架 | 采用方 | 关键文件 |
|-----------|---------|---------------|
| SORA | EASA、英国 CAA、CASA、加拿大交通部 | ConOps、地面风险、空中风险、OSO 缓解措施 |
| Part 107 豁免 | FAA（现行） | 安全论证、运行限制 |
| Part 108 | FAA（即将推出） | 基于性能的合规 |

**平台要求：**
- 以监管机构认可的格式导出飞行日志
- 生成支持 ConOps 的数据（飞行区域、时间、频次）
- 维护运营商/飞行员认证记录
- 跟踪航空器维护和检查状态

---

### 6. 飞行员/运营商执照

| 司法管辖区 | BVLOS 执照 | 备注 |
|--------------|---------------|-------|
| 美国 | Part 107 + 豁免（现行），Part 108（未来） | [[2]](#ref-2) |
| 欧盟 | A2/STS 证书或特定类别授权 | [[1]](#ref-1) |
| 英国 | GVC + 运营授权 | |
| 日本 | 一级（4 级），二级（3 级） | [[19]](#ref-19) |
| 加拿大 | 高级飞行员证书 + SFOC | [[9]](#ref-9) |
| 印度 | BVLOS 飞行员执照（2025 年 4 月起） | [[25]](#ref-25) |

**平台要求：**
- 飞行授权前进行飞行员身份验证
- 执照/认证核验及到期跟踪
- 按飞行员记录飞行小时数
- 基于角色的访问控制（飞行员、观察员、机队管理员）

---

### 7. 飞行日志与审计追踪

所有司法管辖区都要求完整的飞行记录。

**平台要求：**
- 完整的遥测记录（位置、姿态、电池、传感器）
- 带时间戳的事件记录（解锁、起飞、航点、降落、上锁）
- 不可篡改的存储（带保留策略的 JetStream）
- 支持导出以便向监管机构提交
- 事件标记与调查支持

---

### 8. 保险与责任

| 司法管辖区 | 要求 |
|--------------|-------------|
| 欧盟 | 强制要求第三方责任保险 |
| 美国 | 豁免条件通常要求投保 |
| 全部 | 授权通常需要提供保险证明 |

**平台要求：**
- 按航空器跟踪保险单
- 飞行授权前核验保险覆盖范围
- 为理赔提供事件记录

---

## 外部集成

BVLOS 合规需要与外部系统进行实时数据交换。以下是我们的平台需要支持的集成。

### 空域与交通

| 系统 | 用途 | 数据流 | 优先级 |
|--------|---------|-----------|----------|
| **ADS-B In** | 感知有人驾驶航空器交通 | 接收 → DAA 系统 | 高 |
| **FLARM** | 感知滑翔机/通用航空交通 | 接收 → DAA 系统 | 中 |
| **UTM/USS** | 无人机交通管理、飞行授权 | 双向 | 高 |
| **LAANC**（美国） | 自动化空域授权 | 请求/响应 | 高（美国） |
| **AIM/NOTAM** | 临时飞行限制、空域变更 | 接收 → 地理围栏 | 高 |

**UTM 服务提供商（USS）：**
- 美国：Airspace Link、Aloft、DroneUp、FlyFreely（参见 [FAA LAANC 列表](https://www.faa.gov/uas/getting_started/laanc)）
- 欧盟：因成员国而异
- 澳大利亚：Wing、Airservices Australia

*注：AirMap 已于 2023 年 6 月关闭 LAANC 服务。*

---

### 远程识别网络

| 系统 | 用途 | 协议 | 优先级 |
|--------|---------|----------|----------|
| **OpenDroneID** | 广播式远程识别标准 | Bluetooth 4/5、WiFi Beacon | 必需 |
| **网络远程识别** | 通过互联网实现全机队识别 | USS API | 必需（美国） |
| **FAA 远程识别** | 美国远程识别合规 | USS 集成 | 必需（美国） |

---

### 监管数据库

| 系统 | 用途 | 数据流 | 优先级 |
|--------|---------|-----------|----------|
| **飞行员执照数据库** | 核验飞行员认证 | 查询 → 授权 | 高 |
| **航空器登记** | 核验登记状态 | 查询 → 飞行前检查 | 高 |
| **保险核验** | 确认保险有效性 | 查询 → 飞行前检查 | 中 |
| **维护记录** | 跟踪适航性 | 双向 | 中 |

---

### 气象与环境

| 系统 | 用途 | 数据流 | 优先级 |
|--------|---------|-----------|----------|
| **航空气象（METAR/TAF）** | 风、能见度、云底高 | 接收 → 飞行规划 | 高 |
| **雷达/降水** | 实时天气危险 | 接收 → 飞行中警报 | 中 |
| **太阳/地磁** | GPS 可靠性预报 | 接收 → 飞行规划 | 低 |

---

### 应急服务

| 系统 | 用途 | 数据流 | 优先级 |
|--------|---------|-----------|----------|
| **ATC（如需要）** | 管制空域协调 | 语音/数据链路 | 视情况而定 |
| **应急服务（911/112）** | 事件报告 | 外发警报 | 必需 |
| **NOTAM 申报** | 发布飞行运行信息 | 外发 → AIM | 视情况而定 |

---

### 集成架构

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

### 集成路线图

| 集成 | 状态 | 备注 |
|-------------|--------|-------|
| OpenDroneID 广播 | ✅ 已支持 | 需要硬件模块 |
| ADS-B In 接收机 | 🔧 路线图 | uAvionix pingRX 或同类产品 |
| UTM/USS（Airspace Link） | 🔧 路线图 | AirHub API，分级定价 |
| UTM/USS（Aloft） | 🔧 路线图 | REST API 集成 |
| LAANC | 🔧 路线图 | 通过 USS 合作伙伴 |
| 航空气象 | 🔧 路线图 | aviationweather.gov（开放，100 次请求/分钟） |
| NOTAM 数据源 | 🔧 路线图 | NASA DIP 或 FAA SWIM |
| 飞行员执照核验 | 🔧 路线图 | 各司法管辖区 API |
| 保险核验 | 🔧 路线图 | 保险提供商 API |

---

## 合规矩阵

| 要求 | 平台组件 | 状态 |
|-------------|-------------------|--------|
| 远程识别广播 | OpenDroneID 模块 + Gateway | ✅ 已支持 |
| 远程识别记录 | JetStream 流 | ✅ 已支持 |
| 探测与避让 | Jetson 传感器融合 | 🔧 路线图 |
| ADS-B In | 交通接收机集成 | 🔧 路线图 |
| 地理围栏（静态） | PX4 `GF_*` 参数 | ✅ 已支持 |
| 地理围栏（动态） | Gateway 策略更新 | ✅ 已支持 |
| C2 主链路 | 通过 NATS 的 LTE/5G | ✅ 已支持 |
| C2 备用链路 | 卫星集成 | 🔧 路线图 |
| 链路丢失失效保护 | PX4 失效保护配置 | ✅ 已支持 |
| 飞行日志 | JetStream 遥测流 | ✅ 已支持 |
| 审计导出 | 日志导出工具 | ✅ 已支持 |
| 飞行员身份验证 | 授权系统 | ✅ 已支持 |
| 执照跟踪 | 机队管理 | 🔧 路线图 |
| 保险跟踪 | 机队管理 | 🔧 路线图 |

请参阅 [安全模型]({{< relref "/fleet/safety" >}})，了解这些内容如何与整体安全架构集成。

---

## 参考文献

<span id="ref-1"></span>**[1]** EASA. "Drones (UAS) - Frequently Asked Questions." European Union Aviation Safety Agency. https://www.easa.europa.eu/en/the-agency/faqs/drones-uas

<span id="ref-2"></span>**[2]** Drone Pilot Ground School. "Part 108: What We Know So Far about the FAA's BVLOS Rule." https://www.dronepilotgroundschool.com/part-108/

<span id="ref-3"></span>**[3]** DroneLife. "FAA Progress on BVLOS Rules and Advanced Air Mobility Integration." 2024 年 8 月. https://dronelife.com/2024/08/01/faa-progress-on-bvlos-rules-and-advanced-air-mobility-integration/

<span id="ref-4"></span>**[4]** FAA. "BEYOND Program." Federal Aviation Administration. https://www.faa.gov/uas/programs_partnerships/beyond

<span id="ref-5"></span>**[5]** DroneUp. "DroneUp Awarded Landmark FAA Approval for Beyond Visual Line of Sight (BVLOS) in the U.S." 2024 年 1 月. https://www.droneup.com/news/droneup-awarded-bvlos

<span id="ref-6"></span>**[6]** DroneDJ. "Amazon gets FAA BVLOS approval for Prime Air drone delivery." 2024 年 5 月. https://dronedj.com/2024/05/30/amazon-gets-faa-bvlos-approval-for-prime-air-drone-delivery/

<span id="ref-7"></span>**[7]** DroneDJ. "Asylon bags new FAA approval for BVLOS drone operations." 2024 年 4 月. https://dronedj.com/2024/04/26/asylon-drone-faa-bvlos-dronesentry/

<span id="ref-8"></span>**[8]** DroneLife. "FAA Grants Historic BVLOS Approval for Multiple Operators in DFW Area." 2024 年 7 月. https://dronelife.com/2024/07/31/faa-grants-historic-bvlos-approval-for-multiple-operators-in-dfw-area/

<span id="ref-9"></span>**[9]** MLT Aikins. "International drone regulations: Key updates on BVLOS, UTM and use cases." https://www.mltaikins.com/insights/international-drone-regulations-key-updates-on-bvlos-utm-and-use-cases/

<span id="ref-10"></span>**[10]** FlytBase. "Understanding EU Drone Regulations for Docked BVLOS Operations." https://www.flytbase.com/blog/eu-drone-regulations-for-bvlos-operations

<span id="ref-11"></span>**[11]** Elsight. "EASA Remote ID Requirements: Compliance for Drone Operators." https://www.elsight.com/blog/primer-on-easa-remote-id-regulations/

<span id="ref-12"></span>**[12]** Sparkco. "Autonomous Drone Delivery Regulatory Approval." https://sparkco.ai/blog/autonomous-drone-delivery-regulatory-approval

<span id="ref-13"></span>**[13]** sees.ai. "Why Us - First UK routine BVLOS permission." https://www.sees.ai/why-us/

<span id="ref-14"></span>**[14]** CASA. "BVLOS drone operations in regional Australia - Consultation." Civil Aviation Safety Authority. https://consultation.casa.gov.au/stakeholder-engagement-group/consultation.2023-10-05.5578154857/

<span id="ref-15"></span>**[15]** Measure Australia. "New CASA Drone Regulations Incoming!" https://www.measureaustralia.com.au/news/new-casa-drone-regulations-incoming

<span id="ref-16"></span>**[16]** Droneii. "Commercial Drone Regulation Report 2025." Drone Industry Insights. https://droneii.com/product/commercial-drone-regulation-report

<span id="ref-17"></span>**[17]** DroneLife. "Sphere Drones Secures BVLOS Area Approval to Advance Commercial Drone Use in Australia." 2024 年 5 月. https://dronelife.com/2024/05/08/sphere-drones-secures-bvlos-area-approval-to-advance-commercial-drone-use-in-australia/

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

## 返回概览

[← 自主载具机队架构]({{< relref "/fleet" >}})
