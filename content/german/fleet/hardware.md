---
title: "Hardware-Stack"
meta_title: "Drohnenflotten-Hardwarespezifikationen | Ubuntu Software"
description: "Produktions-Hardwarespezifikationen für eine Flotte von 1.000 Drohnen: Holybro-X500-Flugrahmen, Pixhawk-6X-Flugcontroller, Raspberry-Pi- und Jetson-Begleitcomputer."
image: "/images/robotics.svg"
draft: false
aliases:
  - "/fleet/hardware/"
---

{{< notice "info" >}}
**Suchen Sie andere Fahrzeugtypen?** Diese Seite behandelt Drohnen-Hardware. Für Bodenfahrzeuge (Pkw, Lkw, AGVs) siehe [Bodenfahrzeug-Plattform]({{< relref "/fleet/platforms/ground" >}}). Einen vollständigen Überblick finden Sie unter [Unterstützte Plattformen]({{< relref "/fleet/platforms" >}}).
{{< /notice >}}

## Bei einem lokalen Händler kaufen

Der Bau einer Drohne erfordert Einarbeitung. **Kaufen Sie bei einem Händler, der Ihnen helfen kann**—er führt die Hardware, wir liefern die Flottensoftware.

Wenn Sie bei einem Partnerhändler kaufen, erhalten Sie:
- Vorab getestete Hardware, die sofort verbunden werden kann
- Support vor Ort von Leuten, die sich mit Drohnen auskennen
- Zugang zu unserer verwalteten Flotteninfrastruktur

---

## Händler mit Ladengeschäft

Besuchen Sie einen Händler vor Ort und lassen Sie sich direkt beim Aufbau helfen:

{{< featured-shops >}}

**Sie betreiben einen Drohnenladen?** [Treten Sie unserem Partnerprogramm bei]({{< relref "/partners/become-partner" >}}), um zusammen mit Ihrer Hardware auch Flotteninfrastruktur anzubieten.

---

## Was Sie benötigen

| Kit | Komponenten | Typische Bezugsquellen |
|-----|------------|-----------------|
| **Komplette Drohne** | X500-Rahmen, Pixhawk 6X, GPS, Motoren, ESCs, Propeller | GetFPV, NewBeeDrone, Holybro |
| **Companion-Computer** | Raspberry Pi 4 + Jetson Orin Nano | Digi-Key, Mouser, SparkFun |
| **Stromversorgung & Akkus** | 4S LiPo 6000mAh, Strommodul | GetFPV, Tattu, Amazon |
| **Funk** | ExpressLRS-Empfänger, SiK-Telemetrie | GetFPV, HappyModel, Holybro |
| **Netzwerk** | Industrieller Ethernet-Switch, 5G-Modem | PUSR, Quectel |

[Vollständige Komponentenliste anzeigen ↓](#full-bill-of-materials)

---

## Hardware für den Flottenbetrieb

Der Flottenbetrieb erfordert Hardware, die Leistungsfähigkeit und Wartbarkeit vereint. Jede Komponente dieses Stacks wurde nach Zuverlässigkeit im Produktivbetrieb, Ersatzteilverfügbarkeit und Wartbarkeit im Feld ausgewählt.

---

## Flugrahmen

### Holybro X500 V2 ARF

Der **Holybro X500 V2** bietet eine ideale Plattform für den Flotteneinsatz:

| Spezifikation | Wert |
|---------------|-------|
| **SKU** | 30125 |
| **Radstand** | 500mm |
| **Rahmengewicht** | ~410g (610g mit Motoren/ESCs) |
| **Max. Startgewicht** | ~2kg |
| **Flugzeit** | 15-20 min (mit Standardnutzlast) |
| **Motorhalterung** | 16x16mm / 19x19mm |
| **Im Lieferumfang** | Motoren (2216 KV920), ESCs (20A BLHeli S), Propeller (1045) |

**Warum X500:**

- **Bewährtes Design** — Tausendfach weltweit im Einsatz, umfangreiches Community-Wissen
- **Ersatzteilverfügbarkeit** — Arme, Landegestell und Kleinteile sind leicht erhältlich
- **Wartungsfreundlich** — Modularer Aufbau, im Feld wartbar
- **Flexible Nutzlast** — Genug Kapazität für zwei Companion-Computer plus Sensoren

Das ARF-Kit (Almost Ready to Fly) enthält Rahmen, Motoren, ESCs und Propeller—das verringert Abweichungen bei der Montage innerhalb einer großen Flotte.

---

## Flugcontroller

### Pixhawk 6X

Der **Pixhawk 6X** mit **PX4 v1.14 LTS** übernimmt alle flugkritischen Funktionen:

| Spezifikation | Wert |
|---------------|-------|
| **Pixhawk 6X** | SKU 11073 |
| **Baseboard** | SKU 18117 |
| **Zubehör** | SKU 15011 (Kabel & Gehäuse) |
| **Prozessor** | STM32H753 (480MHz Cortex-M7) |
| **IMU** | Dreifach redundant (ICM-42688-P, ICM-45686, BMI088) |
| **Barometer** | Doppelt (MS5611, ICP-20100) |
| **Magnetometer** | IST8310 |
| **Schnittstellen** | 3x CAN, 6x UART, SPI, I2C, PWM |

### Flugcontroller-Zubehör

| Komponente | SKU | Zweck |
|-----------|-----|---------|
| **PM06 V2 Power Module** | 15019 | Strommessung + 5V/3A-BEC |
| **Holybro M9N GPS** | 12027 | GPS- und Kompassmodul |
| **GPS-Mast** | 12033 | Klappbarer Mast zur Verringerung von Störungen |
| **Sicherheitsschalter & Summer** | 12007 | Absicherung beim Scharfschalten + akustische Warnungen |
| **Vibrationsdämpfer** | 12010 | Schwingungsdämpfung |
| **Kabelsatz** | 12055 | JST-GH/CAN-Kabel |

**Warum Pixhawk 6X:**

- **Redundanz** — Dreifache IMU und doppeltes Barometer für Sensor-Voting
- **PX4-LTS-Unterstützung** — Langzeitstabilität, Sicherheitspatches
- **FMUv6X-Standard** — Austauschbar mit anderen konformen Boards
- **Vibrationsisolierung** — Integrierte IMU-Dämpfung

Der Flugcontroller übernimmt:

- Lageschätzung und Lageregelung
- Positionshaltung und Navigation
- Failsafe-Verhalten (RTL, Landung, Geofence)
- Verarbeitung der RC-Eingaben
- MAVLink-Telemetrie

---

## Companion-Computer

Flottenfahrzeuge haben **zwei Companion-Computer** mit klar getrennten Aufgaben:

### Raspberry Pi 4 / CM4 (Sensor-Companion)

Übernimmt die leichtgewichtige Sensorintegration und Datenerfassung:

| Spezifikation | Wert |
|---------------|-------|
| **Modell** | Raspberry Pi 4 Model B |
| **SKU** | RPI4-MODEL-B-8GB |
| **Prozessor** | Quad-Core Cortex-A72 @ 1.8GHz |
| **RAM** | 8GB |
| **Speicher** | Industrielle microSD 128GB (SDSDQED-128G-GNSIN) |
| **Stromversorgung** | 5V @ 3A typisch |

**Aufgaben:**

- Treiber für Umweltsensoren (Temperatur, Luftfeuchtigkeit, Luftqualität)
- Protokollierung von GPS/GNSS-Daten
- Kameraaufnahme (Workloads ohne KI)
- MAVLink-Routing zum Jetson
- Store-and-Forward, wenn der Jetson offline ist

### NVIDIA Jetson (KI-Companion)

Übernimmt rechenintensive Workloads:

| Modell | SKU | GPU-Kerne | KI-Leistung | Leistungsaufnahme |
|-------|-----|-----------|----------------|-------|
| **Jetson Orin Nano 8GB** ★ | 900-13767-0000-000 | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin Nano Super** | 945-13766-0000-000 | 1024 CUDA | 67 TOPS | 7-25W |
| **Jetson Orin NX** | 900-13767-0040-000 | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 900-13701-0000-000 | 2048 CUDA | 275 TOPS | 15-60W |

★ = Auswahl in der Referenz-Stückliste (BOM)

| Speicher | SKU | Kapazität |
|---------|-----|----------|
| **NVMe-SSD** | WD-SN530-256G | 256GB |

**Aufgaben:**

- Computer Vision (Objekterkennung, Tracking)
- Visuell-inertiale Odometrie
- Pfadplanung und Hindernisvermeidung
- Vehicle Gateway (NATS-Client)
- NATS-Leaf-Node

Im Flotteneinsatz kommt typischerweise der **Jetson Orin Nano** für kosteneffiziente Inferenz oder der **Orin NX** für anspruchsvolle Perception-Workloads zum Einsatz.

---

## Funksysteme

### RC-Steuerung: ExpressLRS

**ExpressLRS** stellt die Funkverbindung bereit, über die der Pilot eingreifen und übersteuern kann:

| Spezifikation | Wert |
|---------------|-------|
| **Empfänger** | ExpressLRS EP1 |
| **SKU** | HM-EP1-2400 |
| **Frequenz** | 2.4GHz |
| **Reichweite** | 10km+ (abhängig von Sendeleistung/Antenne) |
| **Latenz** | <5ms |
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
<summary><strong>Alle Verkäufer nach Komponente anzeigen (57 Artikel)</strong></summary>

{{< bom-sellers >}}

</details>

---

## Weiter

[Software-Stack →]({{< relref "/fleet/software" >}})
