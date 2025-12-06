# Ubuntu Sensing — Technical Reference

## Overview

A spatial intelligence platform connecting LiDAR, cameras, and industrial sensors to real-time cloud analytics. Deployable on drones, tripods, robots, and fixed installations.

**Business Model:** Clients bring hardware, we provide edge agent + cloud platform + processing. Revenue from software subscriptions, not hardware sales.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│  EDGE (Client Hardware)                                                 │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Sensors                                                          │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │  │
│  │  │ Livox   │ │RealSense│ │  GPS    │ │  IMU    │ │ Modbus  │   │  │
│  │  │ LiDAR   │ │ RGB-D   │ │ GNSS    │ │         │ │ Sensors │   │  │
│  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘   │  │
│  │       └───────────┴───────────┴───────────┴───────────┘         │  │
│  └───────────────────────────────┬──────────────────────────────────┘  │
│                                  ▼                                      │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Edge Agent (Go binary)                                           │  │
│  │                                                                   │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │  │
│  │  │ Livox   │ │RealSense│ │  GPS    │ │ Modbus  │ │  CAN    │   │  │
│  │  │ Plugin  │ │ Plugin  │ │ Plugin  │ │ Plugin  │ │ Plugin  │   │  │
│  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘   │  │
│  │       └───────────┴───────────┴───────────┴───────────┘         │  │
│  │                               ▼                                  │  │
│  │                    ┌─────────────────┐                          │  │
│  │                    │  Unified Data   │                          │  │
│  │                    │  Model + Buffer │                          │  │
│  │                    └────────┬────────┘                          │  │
│  └─────────────────────────────┼────────────────────────────────────┘  │
│                                │                                        │
│  ┌─────────────────────────────┼────────────────────────────────────┐  │
│  │  Comms                      │                                     │  │
│  │  ┌──────────────────────────▼─────────────────────────────────┐  │  │
│  │  │  5G/LTE Modem (Quectel RM520N-GL) + eUICC eSIM             │  │  │
│  │  │  OTA provisioning — no SIM swapping                        │  │  │
│  │  └──────────────────────────┬─────────────────────────────────┘  │  │
│  └─────────────────────────────┼────────────────────────────────────┘  │
└────────────────────────────────┼────────────────────────────────────────┘
                                 │
                        5G/LTE + NATS
                                 │
┌────────────────────────────────▼────────────────────────────────────────┐
│  CLOUD                                                                  │
│                                                                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │ NATS        │  │ PocketBase  │  │ Processing  │  │ Datastar    │   │
│  │ JetStream   │──│ (DB, Auth)  │──│ Workers     │──│ Web UI      │   │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Deployment Modes

| Mode | Mount | Use Case |
|------|-------|----------|
| Aerial | DJI M300/M350 via SkyPort V2 | Surveying, inspection, mapping |
| Ground | Tripod or backpack | Interior scanning, construction |
| Robot | Viam, ROS2, custom | Navigation, pick-and-place, safety |
| Fixed | Permanent mount | Traffic, security, warehouse |

**Same edge agent binary, same sensors, different config.**

---

## DJI Integration

### PSDK (Payload SDK)

DJI's SDK for custom payloads on enterprise drones. C/C++ only — use a thin shim to bridge to Go edge agent.

| Component | Purpose | Cost |
|-----------|---------|------|
| SkyPort V2 | Payload adapter, twist-lock mount | ~$600 |
| X-Port | SkyPort + 3-axis gimbal | ~$800 |
| PSDK Expansion Board | Bench development, no drone needed | ~$100-150 |

### What PSDK Provides

- Power: 13-17V from drone
- Ethernet: 100Mbps to payload
- UART: Serial comms
- GPS/IMU: Drone position and attitude
- Time sync: PPS signal for georeferencing

### Architecture on Drone

```
┌─────────────────────────────────────────┐
│  DJI M350 RTK                           │
│                                         │
│         Gimbal Port (underneath)        │
└────────────────┬────────────────────────┘
                 │
          ┌──────▼──────┐
          │  SkyPort V2 │  ← Twist-lock mount
          └──────┬──────┘
                 │ Power + Ethernet + PPS
          ┌──────▼──────────────────────────┐
          │  Your Payload Module            │
          │                                 │
          │  ┌─────────┐  ┌──────────────┐ │
          │  │ Jetson  │  │ Livox        │ │
          │  │ Orin    │  │ Mid-360      │ │
          │  │ + LTE   │  │              │ │
          │  └─────────┘  └──────────────┘ │
          └─────────────────────────────────┘
```

### PSDK Shim Architecture

```
┌─────────────────────────────────────────┐
│  Jetson                                 │
│                                         │
│  ┌──────────────┐    ┌───────────────┐ │
│  │ PSDK C Shim  │◄──►│ Go Edge Agent │ │
│  │ (thin daemon)│    │ (your code)   │ │
│  └──────┬───────┘    └───────────────┘ │
│         │ Unix socket / local NATS      │
│         ▼                               │
│  SkyPort Interface                      │
└─────────────────────────────────────────┘
```

No Go bindings for PSDK exist — run minimal C daemon, communicate via IPC.

---

## Sensors

### LiDAR

| Sensor | Range | FOV | Weight | Cost | Use Case |
|--------|-------|-----|--------|------|----------|
| Livox Mid-360 | 40m | 360° | 265g | ~$1,000 | Indoor, robotics, drone |
| Livox Avia | 450m | 70° | 498g | ~$1,500 | Aerial surveying |
| Livox HAP | 150m | 120° | 550g | ~$1,200 | Automotive |

**Livox SDK2** — UDP protocol, can be simulated.

### Cameras

| Sensor | Type | Interface | Cost | Use Case |
|--------|------|-----------|------|----------|
| Intel RealSense D455 | RGB-D | USB | ~$300 | Short range, indoor |
| Luxonis OAK-D | RGB-D + AI | USB | ~$300 | Edge inference |

### IMU

| Sensor | Grade | Interface | Cost |
|--------|-------|-----------|------|
| BNO055 | Consumer | I2C | ~$30 |
| Xsens MTi-630 | Industrial | SPI/UART | ~$2,000 |
| SBG Ellipse-D | Survey | Serial | ~$5,000 |

### GPS/GNSS

| Sensor | Accuracy | RTK | Cost |
|--------|----------|-----|------|
| u-blox ZED-F9P | 1cm | Yes | ~$200-400 |
| u-blox M8N | 2m | No | ~$50 |

---

## Industrial Protocols

| Protocol | Transport | Use Case | Simulation |
|----------|-----------|----------|------------|
| Modbus RTU | RS-485 serial | PLCs, meters, inverters | USB-RS485 adapter (~$15) + diagslave |
| Modbus TCP | Ethernet | Modern PLCs | pymodbus on localhost |
| CAN bus | CAN | Vehicles, AGVs, robots | CANable (~$50) + vcan |
| OPC-UA | Ethernet | Industrial automation | open62541 simulator |

### Go Libraries

```go
// Modbus
import "github.com/goburrow/modbus"

// CAN
import "go.einride.tech/can"
```

---

## Connectivity

### 5G Modem with eSIM OTA

**Problem:** Drone in air — can't scan QR codes for eSIM provisioning.

**Solution:** M2M eSIM with server-push OTA provisioning.

| Type | Provisioning | User Interaction |
|------|--------------|------------------|
| Consumer (SGP.22) | QR code scan | Required |
| M2M (SGP.02/32) | Server push OTA | None |

### How M2M eSIM Works

1. Modem ships with **bootstrap profile** (pre-installed by provider)
2. Device powers on → connects to bootstrap network
3. Your platform calls provider API → triggers profile download
4. SM-SR pushes operational profile OTA
5. Modem switches carrier — no human interaction

### Recommended Modem

**Quectel RM520N-GL**
- 5G Sub-6, fallback to 4G LTE
- M.2 form factor
- eUICC support
- AT commands for eSIM management

```
AT+QESIM="ota"      — Download profile OTA
AT+QESIM="list"     — List installed profiles
AT+QESIM="enable"   — Enable profile
AT+QESIM="disable"  — Disable profile
AT+QESIM="delete"   — Delete profile
```

### M2M eSIM Providers

| Provider | Bootstrap + OTA | API | Notes |
|----------|-----------------|-----|-------|
| 1NCE | Yes | Yes | €10/10 years, 500MB |
| Hologram | Yes | Yes | Hyper eUICC |
| Eseye | Yes | Yes | AnyNet bootstrap |
| Twilio IoT | Yes | Yes | Super SIM |
| EMnify | Yes | Yes | European focus |
| Soracom | Yes | Yes | APAC focus |

### Carrier Switching via API

```go
// Switch drone to local carrier mid-flight
func SwitchCarrier(deviceID, carrier string) error {
    return esimProvider.PushProfile(deviceID, carrier)
    // Provider's SM-SR pushes OTA
    // Modem receives, installs, switches
}
```

---

## Simulation & Development

### Philosophy

Simulate everything on desk. Code against simulators. Deploy same code with real hardware.

### Simulation Hardware

| Real Hardware | Simulation | Cost |
|---------------|------------|------|
| DJI Drone | PSDK Expansion Board | ~$150 |
| Livox LiDAR | Livox SDK simulator (free) | $0 |
| RealSense | .bag file replay (free) | $0 |
| Modbus devices | USB-RS485 + diagslave | ~$15 |
| CAN bus | CANable + vcan | ~$50 |
| GPS/GNSS | u-blox eval kit or gpsfake | ~$50 |
| IMU | BNO055 breakout | ~$30 |
| 5G connectivity | USB LTE dongle or eSunFi hotspot | ~$100-150 |

### Minimum Dev Kit

| Item | Cost |
|------|------|
| PSDK Expansion Board | $150 |
| USB to RS-485 adapter | $15 |
| CANable | $50 |
| BNO055 IMU breakout | $30 |
| u-blox GPS eval | $50 |
| eSunFi 4G eSIM hotspot | $150 |
| **Total** | **~$450** |

### Network Simulation

```bash
# Simulate 200ms latency + 5% packet loss
sudo tc qdisc add dev eth0 root netem delay 200ms loss 5%

# Simulate 1Mbps bandwidth limit
sudo tc qdisc add dev eth0 root tbf rate 1mbit burst 32kbit latency 400ms

# Remove simulation
sudo tc qdisc del dev eth0 root
```

### Dev Setup

```
┌─────────────────────────────────────────────────────────────┐
│  Your Desk                                                  │
│                                                             │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────────┐  │
│  │ PSDK    │ │ Livox   │ │ Modbus  │ │ GPS/IMU         │  │
│  │ Exp Bd  │ │ Sim     │ │ Sim     │ │ Eval Boards     │  │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────────┬────────┘  │
│       └───────────┴───────────┴────────────────┘           │
│                           │                                 │
│                    ┌──────▼──────┐                         │
│                    │ Go Edge     │                         │
│                    │ Agent       │                         │
│                    └──────┬──────┘                         │
│                           │                                 │
│              ┌────────────┼────────────┐                   │
│              ▼            ▼            ▼                   │
│        ┌─────────┐ ┌───────────┐ ┌──────────┐            │
│        │ Local   │ │ PocketBase│ │ Datastar │            │
│        │ NATS    │ │ (Docker)  │ │ UI       │            │
│        └─────────┘ └───────────┘ └──────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### Public Datasets for Testing

| Dataset | Contents | Use |
|---------|----------|-----|
| KITTI | LiDAR + camera + GPS | Autonomous driving |
| NuScenes | Multi-sensor | 3D object detection |
| Newer College | Handheld LiDAR | Indoor SLAM |
| TUM RGB-D | RGB-D | Indoor mapping |

---

## Edge Agent Plugin Interface

```go
type SensorPlugin interface {
    Init(config PluginConfig) error
    Start() error
    Stop() error
    Subscribe(handler func(DataPacket)) error
    Status() PluginStatus
}

type DataPacket struct {
    Timestamp   time.Time
    SensorID    string
    Type        DataType  // PointCloud, Image, Position, Telemetry
    Payload     []byte
    Metadata    map[string]any
}
```

### Example Plugin Config

```yaml
sensors:
  - type: livox-mid360
    mode: hardware  # or "simulator"
    ip: 192.168.1.100

  - type: realsense
    mode: hardware
    serial: 943222071234

  - type: modbus-tcp
    mode: hardware
    host: 10.0.0.50:502
    registers:
      - name: inverter_power
        address: 40001

  - type: gps
    mode: hardware
    port: /dev/ttyUSB0
    baud: 115200

transport:
  nats:
    url: nats://cloud.ubuntusoftware.com
    token: ${UBUNTU_TOKEN}
  buffer:
    path: /var/lib/ubuntu-sensing/buffer
    max_size: 1GB

processing:
  downsample: 0.05  # 5cm voxel grid
  ground_plane: true
```

---

## Payload Module BOM

For DJI M300/M350 drone mount:

| Component | Model | Cost |
|-----------|-------|------|
| LiDAR | Livox Mid-360 | $1,000 |
| Compute | Jetson Orin Nano | $500 |
| DJI Interface | SkyPort V2 | $600 |
| 5G Modem | Quectel RM520N-GL | $150 |
| LTE Modem (if no 5G) | Quectel EM06 | $80 |
| eSIM | M2M provider | $10-50 |
| Housing | 3D printed / machined | $300 |
| Misc | Cables, power reg, antenna | $100 |
| **Total** | | **~$2,700** |

**Weight:** ~800g-1kg (within M350's 2.7kg payload limit)

---

## Client Deployment Flow

```
1. Discovery
   └── Understand their hardware, workflow, use case

2. Pilot
   ├── Ship edge agent binary
   ├── Client installs on their compute
   ├── Configure sensors via YAML
   └── Connect to our cloud (or their on-prem)

3. Integration
   ├── Custom processing pipelines
   ├── API integration with their systems
   └── eSIM provisioning for their region

4. Scale
   ├── Fleet management dashboard
   ├── Multi-device coordination
   └── Usage-based billing
```

---

## Go-to-Market: Developer Hook

**Entry point:** Free, open-source edge agent with simulator support.

```bash
# Developer tries it in 5 minutes
./ubuntu-sensing --config=demo.yaml --simulate

# Sees live point cloud in browser
# Gets hooked
# Wants real sensors → needs cloud → subscription
```

**Land and expand:**
1. Developer downloads, tries simulator
2. Connects real sensor, streams to free tier
3. Needs more storage/devices → paid tier
4. Enterprise needs custom integration → enterprise deal

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Edge Agent | Go |
| Sensor SDKs | Livox SDK2, librealsense, cgo shims |
| Transport | NATS JetStream |
| Database | PocketBase |
| Web UI | Datastar (HTMX), templ, Tailwind |
| Processing | PDAL, PCL, custom Go |
| Deployment | Fly.io, Cloudflare |
| Point Cloud Viz | Potree, Three.js |

---

## Key Differentiators

| Us | Them (ROCK, Emesent, etc.) |
|----|----------------------------|
| Bring your own hardware | Buy our hardware |
| Open data formats | Proprietary cloud lock-in |
| API-first | GUI-only workflows |
| Real-time streaming | Post-flight processing |
| $299/mo software | $20k+ hardware + subscription |
| M2M eSIM OTA | Manual SIM management |

---

## Links & Resources

- Livox SDK2: github.com/Livox-SDK/Livox-SDK2
- Intel RealSense: github.com/IntelRealSense/librealsense
- DJI PSDK: github.com/dji-sdk/Payload-SDK
- Viam RDK: github.com/viamrobotics/rdk
- PDAL: pdal.io
- Potree: github.com/potree/potree
