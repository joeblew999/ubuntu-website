---
title: "Sensorik"
meta_title: "Sensorik & Wahrnehmungsplattform | Ubuntu Software"
description: "Multi-Sensor-Integration für räumliche Intelligenz—LiDAR, Kameras und Industriesensoren vereint durch einen einzigen Edge-Agenten mit 5G/eSIM-Konnektivität."
image: "/images/spatial.svg"
draft: false
---

## Räumliche Sensorik

Echtwelt-Wahrnehmung für digitale Zwillinge, Robotik und autonome Systeme. Verbinden Sie LiDAR, Kameras und Industriesensoren mit Ihren räumlichen Modellen.

---

## Das Problem

Sensoren erzeugen Daten. Um diese Daten sinnvoll zu nutzen, braucht man:

- **Kontext** — Wo befindet sich der Sensor? Was erfasst er?
- **Fusion** — Kombination mehrerer Sensordatenströme zu einem kohärenten Bild
- **Integration** — Verbindung mit Design-Tools, nicht nur Dashboards

Die meisten Sensorlösungen enden bei der Datenerfassung. Wir verbinden Sensoren mit räumlichen Modellen.

---

## Einsatzmodi

Gleicher Edge-Agent. Gleiche Sensoren. Unterschiedliche Konfiguration.

| Modus | Plattform | Anwendungsfall |
|-------|-----------|----------------|
| **Luftgestützt** | DJI Enterprise-Drohnen | Vermessung, Inspektion, Kartierung |
| **Bodengestützt** | Stativ, Rucksack | Innenraumscanning, Bauwesen |
| **Roboter** | Viam RDK, ROS2 | Navigation, Pick-and-Place |
| **Fest montiert** | Permanente Montage | Verkehr, Sicherheit, Lager |

---

## Sensor-Abstraktion

**Hardware-agnostisch by Design.** Ihr Code spricht mit unserer einheitlichen API, nicht mit einzelnen Sensortreibern.

### Unterstützte Sensortypen

| Typ | Beispiele |
|-----|-----------|
| **LiDAR** | Livox Mid-360, Avia |
| **RGB-D Kameras** | Intel RealSense, Luxonis OAK-D |
| **Position** | GPS/GNSS (u-blox RTK), IMU |
| **Industrie** | Modbus-Sensoren, CAN-Bus |

Sensoren austauschen ohne Code-Änderungen. Konfigurationsgesteuert, nicht codegesteuert.

---

## Edge-Agent-Architektur

Go-Binary, das auf Ihrer Hardware läuft—Raspberry Pi, Jetson, industrielles Linux oder Custom ARM.

| Fähigkeit | Beschreibung |
|-----------|--------------|
| **Plugin-System** | Sensoren per Konfiguration hinzufügen, keine Code-Änderungen |
| **Lokales Puffern** | Store-and-Forward bei Offline-Betrieb |
| **Echtzeit-Streaming** | NATS JetStream zur Cloud |
| **Leichtgewichtig** | Einzelne Binary, keine Laufzeit-Abhängigkeiten |

---

## Konnektivität

### 5G/LTE mit eSIM OTA

Kein SIM-Tausch. Kein QR-Code-Scannen. Server-Push-Provisionierung.

- Modem wird mit Bootstrap-Profil geliefert
- Ihre Plattform löst den Carrier-Profil-Download aus
- Carrier-Wechsel während des Betriebs via API

Funktioniert für Drohnen in der Luft, Roboter in Bewegung, feste Installationen an abgelegenen Standorten.

---

## Integration mit Spatial

Sensoren speisen direkt in Ihre 3D-Modelle ein:

| Datenfluss | Zweck |
|------------|-------|
| Punktwolken → Spatial-Modell | Realitätserfassung |
| GPS/IMU → Modellpositionierung | Georeferenzierung |
| Umgebungssensoren → Zwilling | Live-Statusaktualisierungen |
| Industrie-I/O → Automatisierung | Closed-Loop-Steuerung |

Nicht nur Dashboards. Sensoren im Kontext.

---

## Aufgebaut auf Foundation

Sensorik erbt automatisch alle [Foundation](/platform/foundation/)-Fähigkeiten:

| Fähigkeit | Bedeutung |
|-----------|-----------|
| **Offline-first** | Erfassung ohne Internet, Sync bei Verbindung |
| **Universelles Deployment** | Edge, Mobil, Desktop—gleicher Agent |
| **Self-Sovereign** | Ihre Sensoren, Ihre Daten, Ihre Server |
| **Echtzeit-Sync** | Streaming zu mehreren Zielen gleichzeitig |

[Mehr über Foundation erfahren →](/platform/foundation/)

---

## Teil von etwas Größerem

Sensorik ist die Wahrnehmungsschicht der Ubuntu Software Plattform.

Für Organisationen, die 3D-Design und KI benötigen, bietet unsere Spatial-Plattform die Designumgebung—mit direkter Integration Ihrer Sensordaten.

[Spatial erkunden →](/platform/spatial/)

---

## Mit uns bauen

Sensoren einsetzen? Wahrnehmungssysteme aufbauen? Lassen Sie uns sprechen.

[Kontakt →](/contact)
