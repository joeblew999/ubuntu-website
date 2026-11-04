---
title: "Software-Stack"
meta_title: "Softwarearchitektur für Drohnenflotten | Ubuntu Software"
description: "Software-Stack für eine Flotte von 1.000 Drohnen: PX4-Flug-Firmware, Ubuntu Server, JetPack, ROS 2 und ein Go-basiertes Vehicle Gateway."
image: "/images/robotics.svg"
draft: false
---

## Softwarearchitektur

Jede Softwarekomponente läuft mit Open-Source-Code und hat klar definierte Verantwortlichkeiten. Der Stack ist auf Fernverwaltung, Over-the-Air-Updates und konsistentes Verhalten über mehr als 1.000 Fahrzeuge hinweg ausgelegt.

---

## Flugcontroller: Pixhawk

### PX4 v1.14 LTS

Auf dem Flugcontroller läuft **PX4 Autopilot v1.14 LTS**:

| Aspekt | Details |
|--------|---------|
| **Version** | v1.14.x LTS (Long Term Support) |
| **Supportzeitraum** | Sicherheits- und kritische Fehlerbehebungen für mehr als 2 Jahre |
| **Konfiguration** | Standardisierte Parameter in der gesamten Flotte |
| **Updates** | Über QGroundControl oder MAVLink |

**Warum LTS:**

- **Stabilität**: Produktionserprobte Codebasis mit minimalen Änderungen
- **Sicherheit**: Kritische Schwachstellen werden ohne Funktionsänderungen gepatcht
- **Vorhersagbarkeit**: Bekanntes Verhalten in der gesamten Flotte
- **Support**: Kommerzieller Support aus dem Dronecode-Ökosystem verfügbar

**Zentrale Module:**

- **EKF2**: Erweiterter Kalman-Filter zur Zustandsschätzung
- **Position Controller**: Positions- und Geschwindigkeitsregelung für Multicopter
- **Navigator**: Missionsausführung, RTL, Geofence
- **Commander**: Zustandsautomat, Failsafe-Behandlung
- **MAVLink**: Telemetrie- und Befehlsprotokoll

---

## Sensor-Companion: Raspberry Pi

### Ubuntu Server 22.04 LTS

Auf dem Raspberry Pi läuft **Ubuntu Server 22.04**:

| Aspekt | Details |
|--------|---------|
| **Distribution** | Ubuntu Server 22.04 LTS |
| **Kernel** | Mainline-Linux mit Pi-Patches |
| **Init-System** | systemd |
| **Supportzeitraum** | Bis April 2027 |

**Warum Ubuntu Server (und nicht Raspberry Pi OS):**

- **Konsistenz**: Gleiches Betriebssystem auf Pi und Jetson
- **Paketverwaltung**: Standard-apt-Repositories, Snaps verfügbar
- **LTS-Support**: Sicherheitsupdates für 5 Jahre
- **Tooling**: Vertraute Administration für Ops-Teams

### Kerndienste

| Dienst | Zweck |
|---------|---------|
| **mavlink-router** | MAVLink-Nachrichtenrouting zwischen FCU und Jetson |
| **sensor-agent** | Erfassung von Umweltsensordaten |
| **gpsd** | GPS/GNSS-Daemon für Positionsdaten |
| **chrony** | NTP-Zeitsynchronisation |

**mavlink-router** ist entscheidend: Er leitet den MAVLink-Verkehr vom seriellen Port des Pixhawk per UDP an den Jetson weiter. So kann das Vehicle Gateway die Telemetrie verarbeiten.

---

## KI-Companion: Jetson

### JetPack 6.x

Auf dem Jetson läuft **NVIDIA JetPack 6.x**:

| Komponente | Version |
|-----------|---------|
| **L4T (Linux for Tegra)** | Basiert auf Ubuntu 22.04 |
| **CUDA** | 12.x |
| **cuDNN** | 8.x |
| **TensorRT** | 8.x |
| **VPI (Vision Programming Interface)** | 3.x |

**Warum JetPack 6:**

- **Neueste Orin-Unterstützung**: Volle Hardwarebeschleunigung auf der Orin-Familie
- **Ubuntu-22.04-Basis**: Wie auf dem Pi, für einen einheitlichen Betrieb
- **Container-Unterstützung**: NVIDIA Container Runtime für isolierte Workloads
- **OTA-Updates**: NVIDIA-Tools für flottenweite Firmware-Updates

### CUDA + TensorRT

Die KI-Inferenz läuft über die **TensorRT**-Runtime:

- **Modelloptimierung**: INT8/FP16-Quantisierung für den Edge-Einsatz
- **Batch-Inferenz**: Effiziente Verarbeitung mehrerer Frames
- **Niedrige Latenz**: Direkter Zugriff auf den GPU-Speicher

Typische Workloads:

- Objekterkennung (YOLO, SSD)
- Semantische Segmentierung
- Optischer Fluss
- Tiefenschätzung

### ROS 2 Humble

**ROS 2 Humble Hawksbill** stellt die Robotik-Middleware bereit:

| Aspekt | Details |
|--------|---------|
| **Distribution** | Humble Hawksbill |
| **Supportzeitraum** | Bis Mai 2027 |
| **DDS** | CycloneDDS (Standard) oder FastDDS |
| **Build-System** | colcon |

**Warum ROS 2 Humble:**

- **LTS-Release**: 5 Jahre Support, passend zu Ubuntu/JetPack
- **Echtzeitfähig**: Deterministische Ausführungspfade verfügbar
- **Sensorintegration**: Umfangreiche Treiberpakete
- **Simulation**: Gazebo-Integration für Tests

**Zentrale Pakete:**

| Paket | Zweck |
|---------|---------|
| **px4_msgs** | MAVLink-Nachrichtendefinitionen für ROS 2 |
| **px4_ros_com** | PX4-ROS-2-Bridge (uXRCE-DDS) |
| **image_transport** | Verarbeitung von Kameradaten |
| **tf2** | Transformationen zwischen Koordinatensystemen |

---

## Vehicle Gateway

### Go-Implementierung

Das **Vehicle Gateway** ist ein selbst entwickelter Go-Dienst, der auf jedem Jetson läuft:

| Aspekt | Details |
|--------|---------|
| **Sprache** | Go 1.21+ |
| **Abhängigkeiten** | nats.go, mavlink-Bibliothek |
| **Deployment** | systemd-Dienst oder Container |
| **Konfiguration** | Umgebungsvariablen + Konfigurationsdatei |

**Warum Go:**

- **Einzelnes Binary**: Keine Laufzeitabhängigkeiten, einfaches Deployment
- **Performance**: Niedrige Latenz, effiziente Speichernutzung
- **Nebenläufigkeit**: Goroutines für die parallele Verarbeitung von MAVLink und NATS
- **NATS-Ökosystem**: Erstklassige nats.go-Clientbibliothek

**Aufgaben:**

1. **MAVLink-Ingest**: MAVLink-Nachrichten vom mavlink-router empfangen und parsen
2. **Zustands-Downsampling**: Telemetrierate für die WAN-Übertragung reduzieren
3. **Event-Extraktion**: Events aus Zustandsübergängen erzeugen
4. **Befehlsausführung**: NATS-Befehle empfangen, validieren und an die FCU senden
5. **Shadow-Abgleich**: Soll- und Ist-Zustand mit dem Flotten-Backend synchronisieren

[Gateway-Details →]({{< relref "/fleet/gateway" >}})

---

## Control Plane

### NATS + JetStream

Die Flottenkoordination läuft über **NATS** mit **JetStream**-Persistenz:

| Komponente | Rolle |
|-----------|------|
| **NATS Server** | Zentrales Pub/Sub-Messaging |
| **JetStream** | Stream-Persistenz, KV-Stores |
| **Leaf Nodes** | Fahrzeuglokale NATS-Instanzen |
| **Hub Clusters** | Regionale Aggregationspunkte |

**Warum NATS:**

- **Performance**: Millionen von Nachrichten pro Sekunde und Node
- **Einfachheit**: Einzelnes Binary, minimale Konfiguration
- **Persistenz**: JetStream für dauerhafte Streams und KV
- **Topologie**: Leaf Nodes, Cluster und Supercluster für jede Größenordnung

**Flottentopologie:**

```
┌─────────────────────────────────────────────────────┐
│                   Global Hub (optional)             │
│                   NATS Supercluster                 │
└─────────────────┬───────────────────┬───────────────┘
                  │                   │
      ┌───────────▼───────┐ ┌─────────▼───────────┐
      │   Regional Hub A  │ │   Regional Hub B    │
      │   3-node cluster  │ │   3-node cluster    │
      └─────────┬─────────┘ └─────────┬───────────┘
                │                     │
    ┌───────────┼───────────┐         │
    │           │           │         │
┌───▼───┐ ┌─────▼───┐ ┌─────▼───┐ ┌───▼───┐
│ Leaf  │ │  Leaf   │ │  Leaf   │ │ Leaf  │
│ VID-1 │ │ VID-2   │ │ VID-3   │ │VID-N  │
└───────┘ └─────────┘ └─────────┘ └───────┘
```

Auf jedem Fahrzeug läuft ein **Leaf Node**, der sich mit seinem regionalen Hub verbindet. Der Leaf Node übernimmt lokales Pub/Sub und bei unterbrochener Verbindung Store-and-Forward.

[NATS-Topologie-Details →]({{< relref "/fleet/nats-topology" >}})

---

## Software-Update-Strategie

### Flottenweites OTA

| Komponente | Update-Mechanismus |
|-----------|------------------|
| **PX4** | MAVLink-Firmware-Upload über QGC oder eigenes Tooling |
| **Pi/Ubuntu** | apt + unattended-upgrades oder Deployment eigener Images |
| **Jetson/JetPack** | NVIDIA-OTA-Tools, Container-Updates |
| **Vehicle Gateway** | Austausch des Binarys über systemd oder Container-Update |
| **NATS Leaf** | Binary-Update mit Graceful Restart |

**Grundsätze:**

- **Gestaffelter Rollout**: Teilmenge aktualisieren, prüfen, ausweiten
- **Rollback-Fähigkeit**: Vorherige Version stets verfügbar
- **Health Checks**: Automatisierte Prüfung nach dem Update
- **Minimale Ausfallzeit**: Updates während der Wartungsfenster

---

## Zusammenfassung

| Ebene | Software | Version |
|-------|----------|---------|
| **Flugcontroller** | PX4 | v1.14 LTS |
| **Sensor-Companion** | Ubuntu Server | 22.04 LTS |
| **Sensor-Companion** | mavlink-router | Aktuell |
| **KI-Companion** | JetPack | 6.x |
| **KI-Companion** | CUDA/TensorRT | 12.x/8.x |
| **KI-Companion** | ROS 2 | Humble |
| **KI-Companion** | Vehicle Gateway | Go 1.21+ |
| **Control Plane** | NATS + JetStream | Aktuell |

---

## Weiter

[NATS-Topologie →]({{< relref "/fleet/nats-topology" >}})
