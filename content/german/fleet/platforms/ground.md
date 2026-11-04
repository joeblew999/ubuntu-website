---
title: "Bodenfahrzeugplattform"
meta_title: "Hardware & Protokolle für Bodenfahrzeugflotten | Ubuntu Software"
description: "Flottenarchitektur für Pkw, Lkw und fahrerlose Transportfahrzeuge (AGVs): CAN-Bus-Integration, J1939, ROS 2, industrielle Rechenplattformen und fahrzeugspezifische Sicherheitsmodelle."
image: "/images/robotics.svg"
draft: false
---

## Hardware für Bodenfahrzeugflotten

Bodenfahrzeuge – Pkw, Lkw, AGVs und autonome Plattformen – nutzen dieselbe Flottenarchitektur wie Drohnen, jedoch mit fahrzeuggerechter Hardware, passenden Protokollen und eigenen Sicherheitsmodellen.

---

## Protokolloptionen

Bodenfahrzeuge kommunizieren über branchenübliche Standardprotokolle:

### CAN-Bus

| Aspekt | Details |
|--------|---------|
| **Standard** | CAN 2.0B / CAN FD |
| **Datenrate** | 500kbps (CAN) / 5Mbps (CAN FD) |
| **Nachrichtenformat** | DBC-Dateien definieren die Signalkodierung |
| **Physikalisch** | Differenzielles Leitungspaar, robust in störungsreichen Umgebungen |

**Der CAN-Bus bietet:**
- Direkten Zugriff auf Fahrzeugsysteme (Motor, Getriebe, Bremsen)
- Echtzeit-Telemetrie (Geschwindigkeit, Drehzahl, Kraftstoff, Temperaturen)
- Diagnosefehlercodes (OBD-II / J1939)
- Integration von Nachrüstgeräten

### J1939 (Nutzfahrzeuge)

Für Lkw, Busse und schweres Gerät:

| Aspekt | Details |
|--------|---------|
| **Standard** | SAE J1939 |
| **Transport** | CAN 2.0B mit 250kbps |
| **Adressierung** | 29-Bit-Identifier mit PGN/SPN |
| **Anwendungen** | Motor, Getriebe, Bremsen, Anhänger |

**J1939 bietet:**
- Herstellerübergreifend standardisierte Parametergruppen (PGNs)
- Flottenweit konsistente Telemetrie
- Integration mit Telematikplattformen
- Einhaltung gesetzlicher Vorschriften

### ROS 2 (Autonome Plattformen)

Für speziell entwickelte autonome Fahrzeuge:

| Aspekt | Details |
|--------|---------|
| **Middleware** | ROS 2 Humble / Iron |
| **Transport** | DDS (Cyclone DDS, Fast DDS) |
| **Nachrichtentypen** | Standard (sensor_msgs, geometry_msgs) + benutzerdefiniert |
| **QoS** | Konfigurierbare Zuverlässigkeit und Dauerhaftigkeit |

**ROS 2 bietet:**
- Native Sensorintegration (LiDAR, Kameras, IMUs)
- Standard-Nachrichtentypen für die Robotik
- Direktes Topic-Bridging zu NATS
- Simulationskompatibilität (Gazebo, CARLA)

---

## Hardware-Stack

### Rechenplattform

Bodenfahrzeuge nutzen typischerweise robuste Industrierechner:

| Plattform | Anwendungsfall | Hinweise |
|----------|----------|-------|
| **NVIDIA Jetson AGX Orin** | Autonome Fahrzeuge | Hochleistungsfähige Umgebungswahrnehmung |
| **NVIDIA Jetson Orin NX** | Flottenfahrzeuge | Ausgewogenes Preis-Leistungs-Verhältnis |
| **Industrie-PC** | Nur Telematik | Geringere Kosten, keine GPU erforderlich |
| **Neousys** | Robuster Automobileinsatz | Großer Temperaturbereich, Zündungssteuerung |

**Aufgaben der Recheneinheit:**
- Vehicle Gateway (CAN/ROS 2 zu NATS)
- NATS-Leaf-Node
- Edge-KI-Inferenz (optional)
- Store-and-Forward bei Verbindungsverlust

### CAN-Schnittstelle

| Gerät | Schnittstelle | Hinweise |
|--------|-----------|-------|
| **PEAK PCAN-USB** | USB | Beliebt, gut unterstützt |
| **Kvaser Leaf Light** | USB | Industrietauglich |
| **SocketCAN-Geräte** | USB/PCIe | Native Linux-Unterstützung |
| **Jetson-natives CAN** | Integriert | Verfügbar auf AGX Orin |

### Konnektivität

| Methode | Anwendungsfall |
|--------|----------|
| **4G/5G-Mobilfunk** | Primäre Flottenkonnektivität |
| **WLAN** | Betrieb im Depot bzw. auf dem Betriebsgelände |
| **Ethernet** | Kabelgebundene Verbindung am Dock |
| **Satellit** | Betrieb in abgelegenen bzw. ländlichen Gebieten |

---

## Fahrzeugtypen

### Personenkraftwagen (Pkw)

| Komponente | Typische Auswahl |
|-----------|-------------------|
| **Protokoll** | CAN-Bus (OBD-II-Port oder direkt) |
| **Recheneinheit** | Jetson Orin NX oder Industrie-PC |
| **Stromversorgung** | 12-V-Bordnetz |
| **Konnektivität** | 4G/5G-Mobilfunk |
| **Sicherheit** | Fahrerübersteuerung, Not-Aus |

**Anwendungsfälle:** Ride-Sharing-Flotten, Mietwagen, Lieferfahrzeuge

### Nutzfahrzeuge (Lkw)

| Komponente | Typische Auswahl |
|-----------|-------------------|
| **Protokoll** | J1939 (CAN-basiert) |
| **Recheneinheit** | Robuster Industrie-PC oder Jetson |
| **Stromversorgung** | 12-V-/24-V-Bordnetz |
| **Konnektivität** | 4G/5G + Satelliten-Backup |
| **Sicherheit** | Fahrerübersteuerung, ELD-Konformität |

**Anwendungsfälle:** Fernverkehrslogistik, Lieferflotten, Servicefahrzeuge

### Autonome Bodenfahrzeuge (AGVs)

| Komponente | Typische Auswahl |
|-----------|-------------------|
| **Protokoll** | ROS 2 nativ oder proprietär |
| **Recheneinheit** | Jetson AGX Orin |
| **Stromversorgung** | 24-V-/48-V-Batteriesystem |
| **Konnektivität** | WLAN (Gelände) + Mobilfunk (im Freien) |
| **Sicherheit** | Not-Aus, sicherheitszertifizierte Sensoren |

**Anwendungsfälle:** Lagerlogistik, Bergbau, Landwirtschaft, Hafenbetrieb

---

## Sicherheitsmodell

Die Sicherheit von Bodenfahrzeugen unterscheidet sich von der von Drohnen, folgt aber denselben Prinzipien:

### Autoritätshierarchie

```
Manual Override (highest authority)
    │
    ├── Steering wheel / brake pedal (driver present)
    ├── E-stop button (AGVs)
    └── Remote operator console
    │
    ▼
Autonomous System (software control)
    │
    ├── Perception → Planning → Control
    └── Geofencing, speed limits, operational domain
    │
    ▼
Fleet Management (lowest authority)
    │
    └── Mission assignment, routing, monitoring
```

### Failsafe-Verhalten

| Auslöser | Drohnen-Äquivalent | Reaktion des Bodenfahrzeugs |
|---------|------------------|------------------------|
| Kommunikationsverlust | Return-to-Launch | An Ort und Stelle anhalten / an den Rand fahren |
| Sensorausfall | Sofort landen | Sicher anhalten, Bediener alarmieren |
| Geofence-Verletzung | Rückkehr zur Grenze | An der Grenze anhalten |
| Not-Aus aktiviert | Motorabschaltung | Sofortiges Bremsen |
| Niedriger Akkustand | RTL / landen | Rückkehr zum Depot / sicher anhalten |

### Netzwerksicherheit

**Dasselbe Prinzip wie bei Drohnen:** NATS ist niemals Teil des Regelkreises.

| Was NATS tut | Was NATS nicht tut |
|----------------|---------------------|
| Telemetrieerfassung | Lenkung in Echtzeit |
| Missionszuweisung | Bremsbefehle |
| Flottenkoordination | Gassteuerung |
| Überwachung & Alarmierung | Sicherheitskritische Entscheidungen |

Die Bordsysteme des Fahrzeugs übernehmen alle sicherheitskritischen Funktionen. Ein Netzwerkverlust bedeutet den Verlust der Überwachung, nicht den Verlust der Fahrzeugkontrolle.

---

## Gateway-Implementierung

Das Vehicle Gateway verbindet fahrzeugeigene Protokolle mit NATS:

```
┌─────────────────────────────────────────────────────────────────┐
│                    GROUND VEHICLE GATEWAY                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   CAN Bus ─────► DBC Parser ─────► Telemetry Subjects           │
│                                                                  │
│   J1939 ───────► PGN Decoder ────► Telemetry Subjects           │
│                                                                  │
│   ROS 2 ───────► Topic Bridge ───► Telemetry Subjects           │
│                                                                  │
│   Commands ◄───────────────────── Command Subjects              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS JETSTREAM                              │
│   fleet.prod.veh.{vehicle_id}.state.*                           │
│   fleet.prod.veh.{vehicle_id}.evt.*                             │
│   fleet.prod.veh.{vehicle_id}.cmd.*                             │
└─────────────────────────────────────────────────────────────────┘
```

**Aufgaben des Gateways:**
- Protokollübersetzung (CAN/J1939/ROS 2 → gemeinsames Telemetrieformat)
- Zustands-Downsampling (Reduzierung hochfrequenter CAN-Daten auf eine handhabbare Rate)
- Ereignisextraktion (Erkennung von Zustandsänderungen und Schwellenwerten)
- Befehlsvalidierung (Sicherheitsprüfungen vor der Ausführung)
- Store-and-Forward (Pufferung bei Verbindungsverlust)

---

## Subject-Mapping

Bodenfahrzeuge verwenden dieselbe Subject-Hierarchie wie Drohnen:

| Subject | Inhalt beim Bodenfahrzeug |
|---------|----------------------|
| `fleet.prod.veh.{id}.state.pos` | GPS-Position, Fahrtrichtung |
| `fleet.prod.veh.{id}.state.vel` | Geschwindigkeit, Beschleunigung |
| `fleet.prod.veh.{id}.state.health` | Motorstatus, Diagnose |
| `fleet.prod.veh.{id}.state.fuel` | Kraftstoffstand / Batterie-SOC |
| `fleet.prod.veh.{id}.evt.trip` | Ereignisse zu Fahrtbeginn/-ende |
| `fleet.prod.veh.{id}.evt.geofence` | Zoneneintritt/-austritt |
| `fleet.prod.veh.{id}.cmd.mission` | Routenzuweisung |

Die Subject-Struktur ist fahrzeugunabhängig. Eine gemischte Flotte aus Drohnen und Lkw verwendet identische Muster.

---

## Integrationsmuster

### ELD-Konformität (US-Lkw)

Anforderungen an Electronic Logging Devices:
- Erfassung der Lenk- und Ruhezeiten (Hours of Service)
- Fahreridentifikation
- Standortprotokollierung
- Datenübermittlung an Behörden

Das Gateway erfasst J1939-Motordaten und lässt sich in ELD-Compliance-Systeme integrieren.

### Telematik-Integration

Bestehende Telematikplattformen können sich über folgende Wege verbinden:
- NATS-Grants (Abonnement bestimmter Fahrzeuge)
- REST-API (Abfrage historischer Daten)
- Webhook-Ereignisse (Echtzeitwarnungen)

### Flottenmanagementsysteme

Integration mit Enterprise-TMS/FMS:
- Fahrzeugzuweisungen
- Routenoptimierung
- Kraftstoffmanagement
- Wartungsplanung

---

## Zusammenfassung

| Aspekt | Drohnen | Bodenfahrzeuge |
|--------|--------|-----------------|
| **Protokoll** | MAVLink | CAN-Bus / J1939 / ROS 2 |
| **Recheneinheit** | Jetson Orin Nano | Jetson / Industrie-PC |
| **Stromversorgung** | 4S LiPo (14,8V) | 12V/24V/48V Fahrzeugbordnetz |
| **Sicherheits-Override** | RC-Sender | Lenkrad / Not-Aus |
| **Failsafe** | Return-to-Launch | An Ort und Stelle anhalten |
| **Vorschriften** | Luftfahrt (FAA/EASA) | Straßenverkehr (DOT) |
| **Flottenarchitektur** | **Identisch** | **Identisch** |

---

## Verwandte Dokumentation

- [Unterstützte Plattformen]({{< relref "/fleet/platforms" >}}) — Übersicht aller Fahrzeugtypen
- [Drohnenplattform]({{< relref "/fleet/platforms/drones" >}}) — UAV-Hardware und MAVLink
- [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) — Architektur der Protokollüberbrückung
- [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) — Fahrzeugunabhängige Sicherheitsprinzipien
- [Subject-Benennung]({{< relref "/fleet/subjects" >}}) — Hierarchie der Telemetrie-Subjects

---

## Erste Schritte

Sie setzen eine Bodenfahrzeugflotte ein? Wir unterstützen Sie bei Architektur und Integration.

[Kontakt aufnehmen →](/contact)
