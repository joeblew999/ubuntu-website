---
title: "Drohnenplattform"
meta_title: "Drohnenflotten-Hardware & -Protokolle | Ubuntu Software"
description: "Hardwarespezifikationen für den Produktionseinsatz von Drohnenflotten: Holybro-X500-Flugrahmen, Pixhawk-6X-Flugcontroller, PX4, MAVLink-Protokoll und NVIDIA-Jetson-Companion-Computer."
image: "/images/robotics.svg"
draft: false
---

## Hardware für Drohnenbetrieb im Flottenmaßstab

Der Flottenbetrieb erfordert Hardware, die Leistungsfähigkeit und Wartbarkeit vereint. Jede Komponente dieses Stacks wurde nach drei Kriterien ausgewählt: Zuverlässigkeit im Produktivbetrieb, Verfügbarkeit von Ersatzteilen und Wartbarkeit im Feld.

---

## Protokoll: MAVLink

Die Drohnen unserer Flotte kommunizieren über **MAVLink** (Micro Air Vehicle Link):

| Aspekt | Details |
|--------|---------|
| **Version** | MAVLink 2.0 |
| **Transport** | Seriell (UART), UDP, TCP |
| **Nachrichtenrate** | 1-50Hz je nach Nachrichtentyp |
| **Verschlüsselung** | MAVLink-2.0-Signierung (optional) |

**MAVLink bietet:**
- Standardisierte Telemetrie (Position, Fluglage, Akku, Sensoren)
- Befehlsprotokoll (Scharfschalten, Start, Wegpunkte, RTL)
- Parameterverwaltung
- Dateiübertragung (Logs, Missionen)

Das [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) übersetzt MAVLink-Nachrichten in NATS-Subjects. So werden Telemetriedaten flottenweit zusammengeführt und Befehle an die gesamte Flotte verteilt.

---

## Flugrahmen

### Holybro X500 V2 ARF

Der **Holybro X500 V2** ist eine ideale Plattform für den Flotteneinsatz:

| Spezifikation | Wert |
|---------------|-------|
| **Achsabstand** | 500mm |
| **Rahmengewicht** | ~410g |
| **Max. Startgewicht** | ~2kg |
| **Flugzeit** | 15-20 Min. (mit Standard-Nutzlast) |
| **Motorhalterung** | 16x16mm / 19x19mm |

**Warum X500:**

- **Bewährtes Design** – Tausendfach weltweit im Einsatz, umfangreiches Community-Wissen
- **Ersatzteilverfügbarkeit** – Arme, Landegestell und Kleinteile leicht erhältlich
- **Wartungsfreundlich** – Modularer Aufbau, im Feld wartbar
- **Flexible Nutzlast** – Genug Kapazität für zwei Companion-Computer plus Sensoren

Das ARF-Kit (Almost Ready to Fly) enthält Rahmen, Motoren, ESCs und Propeller. Dadurch fallen Montageunterschiede in einer großen Flotte geringer aus.

---

## Flugsteuerung

### Pixhawk 6X

Der **Pixhawk 6X** mit **PX4 v1.14 LTS** übernimmt alle flugkritischen Funktionen:

| Spezifikation | Wert |
|---------------|-------|
| **Prozessor** | STM32H753 (480MHz Cortex-M7) |
| **IMU** | Dreifach redundant (ICM-42688-P, ICM-45686, BMI088) |
| **Barometer** | Doppelt (MS5611, ICP-20100) |
| **Magnetometer** | IST8310 |
| **Schnittstellen** | 3x CAN, 6x UART, SPI, I2C, PWM |

**Warum Pixhawk 6X:**

- **Redundanz** – Dreifach-IMU und Doppelbarometer für Sensor-Voting
- **PX4-LTS-Unterstützung** – Langfristige Stabilität, Sicherheitspatches
- **FMUv6X-Standard** – Austauschbar mit anderen konformen Boards
- **Vibrationsisolierung** – Integrierte IMU-Dämpfung

Der Flight Controller übernimmt:

- Lageschätzung und -regelung
- Positionshaltung und Navigation
- Failsafe-Verhalten (RTL, Landung, Geofence)
- Verarbeitung der RC-Eingaben
- MAVLink-Telemetrie

---

## Companion-Computer

In jedem Flottenfahrzeug arbeiten **zwei Companion-Computer** mit getrennten Aufgaben:

### Raspberry Pi 4 / CM4 (Sensor-Companion)

Übernimmt die einfache Sensorintegration und die Datenerfassung:

| Spezifikation | Wert |
|---------------|-------|
| **Prozessor** | Quad-Core Cortex-A72 @ 1.8GHz |
| **RAM** | 4GB oder 8GB |
| **Speicher** | 32GB+ microSD oder eMMC |
| **Stromversorgung** | 5V @ 3A typisch |

**Aufgaben:**

- Treiber für Umweltsensoren (Temperatur, Luftfeuchtigkeit, Luftqualität)
- Aufzeichnung von GPS/GNSS-Daten
- Kameraaufnahmen (ohne KI-Verarbeitung)
- MAVLink-Routing zum Jetson
- Store-and-Forward, wenn der Jetson offline ist

### NVIDIA Jetson (KI-Companion)

Übernimmt rechenintensive Workloads:

| Modell | GPU-Kerne | KI-Leistung | Leistungsaufnahme |
|-------|-----------|----------------|-------|
| **Jetson Orin Nano** | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin NX** | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 2048 CUDA | 275 TOPS | 15-60W |
| **Jetson Xavier NX** | 384 CUDA | 21 TOPS | 10-20W |

**Aufgaben:**

- Computer Vision (Objekterkennung, Tracking)
- Visuell-inertiale Odometrie
- Pfadplanung und Hindernisvermeidung
- Vehicle Gateway (NATS-Client)
- NATS-Leaf-Node

In Flotten kommt meist der **Jetson Orin Nano** für kostengünstige Inferenz zum Einsatz oder der **Orin NX** für anspruchsvolle Wahrnehmungsaufgaben.

---

## Funksysteme

### RC-Steuerung: ExpressLRS

**ExpressLRS** stellt die Funkverbindung bereit, über die der Pilot jederzeit eingreifen kann:

| Spezifikation | Wert |
|---------------|-------|
| **Frequenz** | 915MHz (US) / 868MHz (EU) |
| **Reichweite** | 10km+ (abhängig von Sendeleistung/Antenne) |
| **Latenz** | <5ms |
| **Protokoll** | Open Source |

**Warum ExpressLRS:**

- **Open Source** – Kein Vendor-Lock-in, Entwicklung durch die Community
- **Reichweite** – Zuverlässige Verbindung für Einsätze in Sichtweite
- **Geringe Latenz** – Präzise manuelle Steuerung, wenn nötig
- **Kosten** – Günstige Empfänger, auch für große Flotten

### Telemetrie: SiK-Funk

**SiK-Funkmodule mit 915 MHz** übertragen die Telemetrie zur Bodenstation:

| Spezifikation | Wert |
|---------------|-------|
| **Frequenz** | 915MHz (US) / 433MHz (EU) |
| **Reichweite** | typisch 1-2km |
| **Datenrate** | Bis zu 250kbps |
| **Schnittstelle** | Seriell (UART) |

**Rolle in der Flotte:**

- Ersatz-Telemetrie, wenn kein Mobilfunk verfügbar ist
- Verbindung zur Bodenkontrollstation
- Lokale Tests und Debugging

Im produktiven Flottenbetrieb laufen die Telemetriedaten hauptsächlich über die Mobilfunkverbindung des Jetson zu NATS. Die SiK-Funkmodule dienen als Rückfallebene und zur Diagnose im Feld.

---

## Energiesystem

### Akku: 4S LiPo

| Spezifikation | Wert |
|---------------|-------|
| **Konfiguration** | 4S (14.8V nominal) |
| **Kapazität** | 5000-6000mAh typisch |
| **Entladerate** | 30C+ |
| **Stecker** | XT60 |

**Hinweise für den Flottenbetrieb:**

- Eine einheitliche Akkukonfiguration für die gesamte Flotte festlegen, um die Logistik zu vereinfachen
- Smart-Akkus mit Telemetrie verwenden, sofern verfügbar
- Akkus im Wechsel einsetzen und ihren Zustand überwachen
- Für den Dauerbetrieb drei Akkus pro Fahrzeug einplanen

### Stromverteilung

- **PDB/PMS** – Power-Management-Board für die Companion-Computer
- **BEC** – Geregelte 5V-Versorgung für Pi und Servos
- **Jetson-Stromversorgung** – Direkt vom Akku über einen Spannungsregler (typisch 12V)

---

## Konnektivität

### Mobilfunk: 4G/5G + eSIM

Jeder Jetson verbindet sich über ein Mobilfunkmodem:

| Spezifikation | Wert |
|---------------|-------|
| **Modem** | Quectel RM520N-GL oder vergleichbar |
| **Frequenzbänder** | 4G LTE + 5G NR |
| **SIM** | eSIM mit OTA-Provisionierung |
| **Schnittstelle** | USB 3.0 oder M.2 |

**Vorteile der eSIM:**

- Kein Austausch physischer SIM-Karten in einer großen Flotte
- Provisionierung beim Netzbetreiber aus der Ferne
- Wechsel des Netzbetreibers für bessere Netzabdeckung
- Zentrale Verwaltung der Mobilfunkverträge

---

## Sicherheitsmodell

Drohnenspezifische Sicherheitsfunktionen greifen auf mehreren Ebenen:

| Ebene | Mechanismus | Verhalten |
|-------|-----------|----------|
| **RC-Übersteuerung** | ExpressLRS | Der Pilot behält stets die manuelle Kontrolle |
| **Flight Controller** | PX4-Failsafes | Erzwingt RTL, Landung und Geofence |
| **Companion** | Vehicle Gateway | Befehlsvalidierung, Sicherheitsprüfungen |
| **Netzwerk** | NATS ist nie Teil des Regelkreises | Nur Überwachung, nicht flugkritisch |

Die vollständige Sicherheitsarchitektur finden Sie unter [Sicherheitsmodell]({{< relref "/fleet/safety" >}}).

---

## Zusammenfassung

| Komponente | Auswahl | Begründung |
|-----------|-----------|-----------|
| **Rahmen** | Holybro X500 V2 ARF | Bewährt, verfügbar, wartbar |
| **FCU** | Pixhawk 6X | Redundanz, PX4-LTS-Unterstützung |
| **Sensor-Companion** | Raspberry Pi 4/CM4 | Kostengünstig, viele GPIOs |
| **KI-Companion** | Jetson Orin/Xavier | GPU-Inferenz, CUDA-Ökosystem |
| **RC** | ExpressLRS | Open Source, hohe Reichweite |
| **Telemetrie** | SiK 915MHz | Ersatzverbindung, Debugging |
| **Stromversorgung** | 4S LiPo 5000mAh | Einheitlicher Flottenstandard |
| **Mobilfunk** | 4G/5G + eSIM | OTA-Provisionierung, Netzabdeckung |
| **Protokoll** | MAVLink 2.0 | Industriestandard, nativ in PX4 |

---

## Bezugsquellen

Wir verkaufen keine Hardware – das übernehmen unsere Partner. **Zertifizierte Drohnenshops** liefern komplette Kits inklusive Flotteninfrastruktur.

**Für Drohnenkäufer:**
Kaufen Sie bei einem [zertifizierten Partner]({{< relref "/partners" >}}) und erhalten Sie:
- Vorkonfigurierte Hardware, die sofort verbunden werden kann
- Zugangsdaten für unsere verwaltete Infrastruktur
- Support vor Ort durch Ihren Shop

**Für Drohnenshops:**
[Werden Sie Teil unseres Partnerprogramms]({{< relref "/partners/become-partner" >}}) und bieten Sie Flotteninfrastruktur zusammen mit Ihrer Hardware an.

---

## Weiterführende Dokumentation

- [Unterstützte Plattformen]({{< relref "/fleet/platforms" >}}) – Übersicht aller Fahrzeugtypen
- [Bodenfahrzeuge]({{< relref "/fleet/platforms/ground" >}}) – Autos, Lkw und AGVs
- [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) – Brücke von MAVLink zu NATS
- [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) – Failsafe-Architektur
- [Software-Stack]({{< relref "/fleet/software" >}}) – Betriebssysteme und Middleware
