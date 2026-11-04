---
title: "Unterstützte Plattformen"
meta_title: "Unterstützte Fahrzeugplattformen | Ubuntu Software"
description: "Unsere Flottenarchitektur unterstützt Drohnen, Autos, Lkw und autonome Bodenfahrzeuge. Dieselbe Kerninfrastruktur, fahrzeugspezifische Integration."
image: "/images/robotics.svg"
draft: false
---

## Eine Architektur, viele Fahrzeuge

Die Flottenarchitektur von Ubuntu Software ist **von Grund auf fahrzeugunabhängig**. Dieselbe NATS-JetStream-Infrastruktur, dieselben Digital-Twin-Muster und dasselbe Autorisierungsmodell, mit denen Drohnenflotten verwaltet werden, verwalten auch Autos, Lkw und autonome Bodenfahrzeuge.

---

## Was universell ist

Diese Komponenten funktionieren bei allen Fahrzeugtypen identisch:

| Komponente | Beschreibung |
|-----------|-------------|
| **NATS JetStream** | Messaging-Backbone, funktioniert für jede Art von Telemetrie |
| **Subject-Hierarchie** | `fleet.{env}.veh.{id}.*` — die Fahrzeug-ID ist nur ein Bezeichner |
| **Digitale Zwillinge** | State-Streams, Event-Streams, Shadow-Stores |
| **Autorisierung** | Dezentrale Sicherheit, Berechtigungen, Zugriff durch Dritte |
| **Streams & Events** | Telemetrie-Rollup, Audit-Trails, abfragbarer Zustand |

**Für die Architektur spielt es keine Rolle, was die Telemetrie erzeugt.** Ein Positionsupdate einer Drohne und ein Positionsupdate eines Lkw fließen durch dieselbe Infrastruktur.

---

## Was fahrzeugspezifisch ist

Jede Plattform hat eigene Anforderungen:

| Aspekt | Variiert je nach Plattform |
|--------|-------------------|
| **Hardware** | Flugzellen vs. Fahrgestelle, Flugsteuerungen vs. Steuergeräte (ECUs) |
| **Protokolle** | MAVLink, CAN-Bus, ROS 2, J1939 |
| **Sicherheitsmodell** | RC-Override vs. Not-Aus, RTL vs. Anhalten an Ort und Stelle |
| **Sensoren** | Flughöhe vs. Bodenfreiheit, Fluggeschwindigkeit vs. Raddrehzahl |
| **Vorschriften** | Luftfahrtbehörden vs. Straßenverkehrsbehörden |

---

## Unterstützte Plattformen

### [Drohnen]({{< relref "/fleet/platforms/drones" >}})

Unbemannte Luftfahrzeuge mit PX4 oder ArduPilot:

- **Hardware**: Holybro X500, Pixhawk 6X, NVIDIA Jetson
- **Protokoll**: MAVLink über Seriell/UDP
- **Sicherheit**: RC-Override, Return-to-Launch, Geofencing
- **Anwendungsfälle**: Inspektion, Kartierung, Überwachung, Zustellung

[Details zur Drohnenplattform →]({{< relref "/fleet/platforms/drones" >}})

---

### [Bodenfahrzeuge]({{< relref "/fleet/platforms/ground" >}})

Autos, Lkw, fahrerlose Transportfahrzeuge (AGVs) und autonome Bodenplattformen:

- **Hardware**: Fahrzeug-Steuergeräte (ECUs), Industrie-PCs, Jetson
- **Protokoll**: CAN-Bus, J1939, ROS 2
- **Sicherheit**: Not-Aus, Anhalten an Ort und Stelle, Geofencing
- **Anwendungsfälle**: Logistik, Bergbau, Landwirtschaft, Last-Mile-Zustellung

[Details zu Bodenfahrzeugen →]({{< relref "/fleet/platforms/ground" >}})

---

## Gateway-Muster

Das [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) verbindet fahrzeugeigene Protokolle mit NATS. Für jede Plattform gibt es eine eigene Gateway-Implementierung:

```
┌─────────────────────────────────────────────────────────────────┐
│                     VEHICLE GATEWAY                              │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   MAVLink GW    │    CAN GW       │      ROS 2 GW               │
│   (Drones)      │  (Cars/Trucks)  │   (ROS Robots)              │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS JETSTREAM                              │
│            Same subjects, same streams, same patterns            │
└─────────────────────────────────────────────────────────────────┘
```

Das Gateway:
- Übersetzt protokollspezifische Nachrichten in ein gemeinsames Telemetrieformat
- Veröffentlicht in der standardisierten Subject-Hierarchie
- Setzt Sicherheitsrichtlinien durch
- Übernimmt Store-and-Forward bei Verbindungsverlust

---

## Neue Plattformen hinzufügen

Die Architektur unterstützt jedes Fahrzeug, das Folgendes kann:

1. **Telemetrie erzeugen** — Position, Geschwindigkeit, Zustandsdaten
2. **Befehle annehmen** — Start, Stopp, Wegpunktnavigation
3. **Ein Gateway ausführen** — das native Protokoll nach NATS übersetzen

Eigene Plattformen werden integriert, indem ein Gateway für ihr Protokoll implementiert wird. Die übrige Infrastruktur (NATS, Streams, Autorisierung) funktioniert unverändert.

---

## Mehr erfahren

| Plattform | Dokumentation |
|----------|---------------|
| [Drohnen]({{< relref "/fleet/platforms/drones" >}}) | PX4-, Pixhawk- und MAVLink-Integration |
| [Bodenfahrzeuge]({{< relref "/fleet/platforms/ground" >}}) | CAN-Bus-, J1939- und Automotive-Integration |
| [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) | Architektur für Protokoll-Bridging |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Plattformspezifische Sicherheitsanforderungen |

---

## Jetzt starten

Sie bauen eine gemischte Flotte auf? Sie setzen Bodenfahrzeuge neben Drohnen ein? Wir unterstützen Sie bei Architektur und Implementierung.

[Kontakt aufnehmen →](/contact)
