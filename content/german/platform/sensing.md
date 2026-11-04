---
title: "Sensorik"
meta_title: "Plattform für Sensorik & Wahrnehmung | Ubuntu Software"
description: "Multisensor-Integration für räumliche Intelligenz – LiDAR, Kameras und Industriesensoren, vereint über einen einzigen Edge-Agenten mit 5G/eSIM-Konnektivität."
image: "/images/spatial.svg"
draft: false
---

## Spatial Sensing

Wahrnehmung der realen Welt für digitale Zwillinge, Robotik und autonome Systeme. Verbinden Sie LiDAR, Kameras und Industriesensoren mit Ihren räumlichen Modellen.

---

## Das Problem

Sensoren erzeugen Daten. Um diese Daten sinnvoll zu nutzen, braucht es:

- **Kontext** — Wo befindet sich der Sensor? Worauf ist er gerichtet?
- **Fusion** — Die Zusammenführung mehrerer Sensordatenströme zu einem stimmigen Gesamtbild
- **Integration** — Die Anbindung an Design-Tools, nicht nur an Dashboards

Die meisten Sensing-Lösungen enden bei der Datenerfassung. Wir verbinden Sensoren mit räumlichen Modellen.

---

## Einsatzmodi

Derselbe Edge-Agent. Dieselben Sensoren. Andere Konfiguration.

| Modus | Plattform | Anwendungsfall |
|------|----------|----------|
| **Luft** | DJI-Enterprise-Drohnen | Vermessung, Inspektion, Kartierung |
| **Boden** | Stativ, Rucksack | Innenraum-Scanning, Bauwesen |
| **Roboter** | Viam RDK, ROS2 | Navigation, Pick-and-Place |
| **Stationär** | Feste Montage | Verkehr, Sicherheit, Lager |

### GPS-freie Navigation

Visuell-inertiale Odometrie ermöglicht autonomen Flug ohne GPS – entscheidend für Innenräume, Häuserschluchten und Zonen ohne GPS-Empfang.

{{< video src="videos/visual-inertial-nav-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*Positionsverfolgung in Echtzeit allein mit Kameras und IMU. Keine externe Infrastruktur erforderlich.*

{{< video src="videos/gps-free-vio-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*Autonomer Flug im geschlossenen Regelkreis mithilfe visuell-inertialer Odometrie.*

Detaillierte Hardware-Spezifikationen finden Sie in unserer [Drohnenflotten-Hardware →]({{< relref "/fleet/hardware" >}})

---

## Sensor-Abstraktion

**Von Grund auf hardwareunabhängig.** Ihr Code kommuniziert mit unserer einheitlichen API, nicht mit einzelnen Sensortreibern.

### Unterstützte Sensortypen

| Typ | Beispiele |
|------|----------|
| **LiDAR** | [Livox Mid-360](https://www.livoxtech.com/mid-360) ([SDK](https://github.com/Livox-SDK/Livox-SDK2)), [Avia](https://www.livoxtech.com/avia) ([Demo](https://terra-1-g.djicdn.com/65c028cd298f4669a7f0e40e50ba1131/Download/Avia/02_Intro%20Demo.mp4)) |
| **RGB-D-Kameras** | [Intel RealSense](https://www.intelrealsense.com/), [Luxonis OAK-D](https://shop.luxonis.com/products/oak-d) |
| **Position** | [u-blox RTK](https://www.u-blox.com/en/technologies/rtk-real-time-kinematic) GPS/GNSS, IMU |
| **Industrie** | Modbus-Sensoren, CAN-Bus |

Sensoren austauschen, ohne den Code zu ändern. Konfigurationsgesteuert, nicht codegesteuert.

---

## Architektur des Edge-Agents

Go-Binary, das auf Ihrer Hardware läuft – Raspberry Pi, Jetson, Industrie-Linux oder individuelles ARM-System.

| Fähigkeit | Beschreibung |
|------------|-------------|
| **Plugin-System** | Sensoren per Konfiguration hinzufügen, ohne Codeänderungen |
| **Lokale Pufferung** | Store-and-Forward im Offline-Betrieb |
| **Echtzeit-Streaming** | NATS JetStream in die Cloud |
| **Leichtgewichtig** | Einzelnes Binary, keine Laufzeitabhängigkeiten |

---

## Konnektivität

### 5G/LTE mit eSIM OTA

Kein SIM-Kartenwechsel. Kein Scannen von QR-Codes. Bereitstellung per Server-Push.

- Das Modem wird mit einem Bootstrap-Profil ausgeliefert
- Ihre Plattform stößt den Download des Netzbetreiberprofils an
- Netzbetreiber während des laufenden Einsatzes per API wechseln

Funktioniert für Drohnen in der Luft, Roboter in Bewegung und stationäre Installationen an abgelegenen Standorten.

---

## Integration mit Spatial

Sensoren speisen ihre Daten direkt in Ihre 3D-Modelle ein:

| Datenfluss | Zweck |
|-----------|---------|
| Punktwolken → Räumliches Modell | Realitätserfassung |
| GPS/IMU → Modellpositionierung | Georeferenzierung |
| Umweltsensoren → Zwilling | Live-Zustandsaktualisierungen |
| Industrielle I/O → Automatisierung | Regelung im geschlossenen Kreis |

Nicht nur Dashboards. Sensoren im Kontext.

---

## Aufgebaut auf Foundation

Sensing übernimmt automatisch alle Fähigkeiten von [Foundation](/platform/foundation/):

| Fähigkeit | Was das bedeutet |
|------------|---------------|
| **Offline-first** | Erfassen ohne Internet, synchronisieren bei bestehender Verbindung |
| **Universeller Einsatz** | Edge, Mobilgerät, Desktop – derselbe Agent |
| **Selbstbestimmt** | Ihre Sensoren, Ihre Daten, Ihre Server |
| **Echtzeit-Synchronisierung** | Gleichzeitiges Streaming an mehrere Ziele |

[Mehr über Foundation erfahren →](/platform/foundation/)

---

## Teil von etwas Größerem

Sensing ist die Wahrnehmungsschicht der Ubuntu-Software-Plattform.

Für Organisationen, die 3D-Design und KI benötigen, bietet unsere Spatial-Plattform die Designumgebung – mit direkter Integration Ihrer Sensordaten.

[Spatial entdecken →](/platform/spatial/)

---

## Produktionsreferenz: Flotte mit 1.000 Drohnen

Sehen Sie, wie diese Sensing-Fähigkeiten auf Produktionsniveau skalieren. Unsere Referenzarchitektur zeigt einen vollständigen Flotteneinsatz mit dem oben beschriebenen Hardware- und Software-Stack – 1.000 Holybro-X500-Drohnen mit PX4, zwei Companion-Computern und NATS JetStream für digitale Zwillinge im Flottenmaßstab.

[Drohnenflotten-Architektur entdecken →]({{< relref "/fleet" >}})

---

## Bauen Sie mit uns

Sie setzen Sensoren ein? Sie entwickeln Wahrnehmungssysteme? Lassen Sie uns sprechen.

[Kontakt →](/contact)
