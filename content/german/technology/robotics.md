---
title: "Robotik"
meta_title: "Robotik-Stack | Ubuntu-Software"
description: "Unsere Robotik-Architektur auf Basis von Viam RDK: vom Design bis zum Deployment, mit Hardware-Abstraktion, Computer Vision und industrieller Integration."
image: "/images/spatial.svg"
draft: false
---

## Robotik-Stack

Vom Design bis zum Deployment. Spatial liefert die Ebene für 3D-Design und Simulation. Viam RDK liefert die Laufzeitumgebung, mit der reale Roboter gesteuert und betrieben werden.

---

## Embedded-Linux-Laufzeitumgebung

Robotik läuft auf Linux. Wir haben jahrzehntelange Erfahrung mit Linux und Embedded-Systemen. Deshalb entwickeln wir für die reale Welt – für Fabrikhallen, Außenumgebungen und Edge-Geräte mit begrenzten Ressourcen.

### Viam-RDK-Laufzeitumgebung

**Warum wir uns für Viam RDK entschieden haben:**

- **Open Source** — Passt zu unserer Philosophie offener Standards
- **Go-nativ** — Passt zu unserem Technologie-Stack
- **Modular** — Komponenten lassen sich austauschen und erweitern
- **Cloud optional** — Funktioniert offline-first, wie unsere Produkte
- **Hardware-agnostisch** — Nicht an bestimmte Roboterhersteller gebunden
- **Läuft auf Embedded Linux** — Raspberry Pi, NVIDIA Jetson, Industriesteuerungen

Viam RDK ist ein Open-Source-Entwicklungskit für Robotik. Es bietet sprachunabhängige SDKs (Go, Python, TypeScript) und eine modulare Komponentenarchitektur.

[Viam-RDK-Dokumentation →](https://docs.viam.com/)

### Embedded-Plattformen

| Plattform | Anwendungsfall |
|----------|----------|
| **Raspberry Pi** | Prototyping, Edge Computing, Kiosksysteme |
| **NVIDIA Jetson** | GPU-beschleunigte Bildverarbeitung und ML |
| **Industrie-Linux** | Fabrikautomatisierung, raue Umgebungen |
| **Kundenspezifische ARM-Boards** | Anwendungsspezifische Deployments |

Dieselbe Go-Codebasis. Dasselbe Viam RDK. Deployment von der Entwicklungs-Workstation bis zum Embedded-Edge-Gerät.

---

## Funktionen

| Funktion | Viam-Service |
|------------|--------------|
| Bewegungsplanung mit Kollisionsvermeidung | motion service |
| Objekterkennung und Segmentierung | vision service + ML-Modelle |
| Punktwolken aus Tiefensensoren | Camera-Komponente |
| Hand-Auge-Kalibrierung | Frame system |
| SLAM mit IMU | slam service |

---

## Bildverarbeitung & Wahrnehmung

### Kamera-Abstraktion

Viam bietet eine einheitliche `rdk:component:camera`-API, die mit vielen verschiedenen Kameratypen funktioniert – mit Webcams, IP-Kameras, LiDAR und Tiefenkameras wie der Intel RealSense D435i.

**So funktioniert es:**

- **Standardisierte API** — Ihr Code kommuniziert mit der Viam-Kameraschnittstelle, nicht direkt mit der Hardware
- **Integrierte Treiber** — Das `viam-camera-realsense`-Modul übernimmt die RealSense-Integration
- **Hardware-agnostisch** — Kameras austauschen, ohne die Anwendungslogik zu ändern
- **Konfigurationsgesteuert** — Auflösung, Sensoren und Streams werden per JSON konfiguriert, nicht im Code

Verfügbare Methoden: `GetImage()`, `GetImages()`, `GetPointCloud()` sowie intrinsische Kameraparameter.

### Intel RealSense D435i

RGB-D-Tiefenerfassung zur Erzeugung von Punktwolken und für die räumliche Wahrnehmung.

[librealsense auf GitHub →](https://github.com/realsenseai/librealsense)

### ML-Pipeline

- Objekterkennung und Segmentierung mit YOLOv8
- Integration mit dem Viam vision service
- Echtzeit-Inferenz auf Edge-Geräten

---

## Hardware-Abstraktion

**Der zentrale Nutzen: Konfigurationsänderung statt Codeänderung.**

- **Gleiche Software, größerer Arm** — xArm-Konfiguration gegen UR5e oder KUKA tauschen und neu deployen
- **xArm als Teach-Pendant** — Bewegungen am xArm vorführen, der größere Arm ahmt sie nach
- **Koordination mehrerer Roboter** — Der xArm übernimmt Kleinteile, der KUKA das Heben schwerer Lasten
- **Entwicklung mit digitalem Zwilling** — Mit dem xArm entwickeln, mit KUKA-Kinematik simulieren

---

## Unterstützte Arme

| Modell | Traglast | Reichweite | Anwendungsfall |
|-------|---------|-------|----------|
| xArm 6 | 5kg | 700mm | Entwicklung |
| UR5e | 5kg | 850mm | Produktion |
| KUKA KR6 | 6kg | 900mm | Kleinteile |
| KUKA KR10 | 10kg | 900-1100mm | Mittlere Montage |
| KUKA KR16 | 16kg | 1600mm | Schweißen, Palettieren |
| KUKA KR30 | 30kg | 2000mm+ | Handhabung schwerer Lasten |

[viam-kuka-Modul →](https://github.com/viam-soleng/viam-kuka)

---

## Architektur

```
SOFTWARE STACK
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Vision   │  │ Motion   │  │   ML     │  │ Business │
│ Pipeline │  │ Planning │  │ Models   │  │  Logic   │
└────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
     └─────────────┴─────────────┴─────────────┘
                         │
                  Viam Arm API
              (rdk:component:arm)
                         │
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐
    │  xArm 6  │   │  UR5e    │   │  KUKA    │
    │  (dev)   │   │  (prod)  │   │  (heavy) │
    │ 5kg/700mm│   │ 5kg/850mm│   │ 30kg+    │
    └──────────┘   └──────────┘   └──────────┘
```

---

## Industrielle Integration

**I/O-Integration:**

Modbus → SPS zur Anbindung an die Industrieautomatisierung. Förderbänder, Sensoren, Aktoren und Sicherheitssysteme lassen sich aus derselben Laufzeitumgebung steuern.

---

## Wie alles zusammenpasst

| Ebene | Technologie | Zweck |
|-------|------------|---------|
| Design | Spatial | 3D-Design und Simulation von Arbeitszellen |
| Laufzeit | Viam RDK | Steuerung und Betrieb von Robotern |
| Bildverarbeitung | RealSense + YOLOv8 | Wahrnehmung und Objekterkennung |
| Arme | xArm / UR / KUKA | Physische Manipulation |
| I/O | Modbus / SPS | Industrielle Integration |

---

## Mehr erfahren

- [Technologieüberblick →](/technology/)
- [Spatial-Plattform →](/platform/spatial/)
- [Kontaktieren Sie uns →](/contact/)
