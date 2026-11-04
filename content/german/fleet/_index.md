---
title: "Architektur einer autonomen Fahrzeugflotte"
meta_title: "Flottenarchitektur für Drohnen, Autos & Lkw | Ubuntu Software"
description: "Produktionsreife Referenzarchitektur für die Verwaltung autonomer Fahrzeugflotten—Drohnen, Autos, Lkw, AGVs—mit NATS JetStream für digitale Zwillinge im Flottenmaßstab."
image: "/images/robotics.svg"
draft: false
---

## Produktionsreifes Flottenmanagement im großen Maßstab

Diese Referenzarchitektur zeigt, wie sich **Flotten autonomer Fahrzeuge** – Drohnen, Autos, Lkw, AGVs – mit Open-Source-Komponenten und bewährten Mustern bereitstellen und verwalten lassen. Jede Designentscheidung zielt auf reale Produktionsanforderungen ab: Zuverlässigkeit, Beobachtbarkeit und kontrollierte Degradation.

Die Kernarchitektur ist **fahrzeugunabhängig**. NATS JetStream, Digital-Twin-Muster und das Autorisierungsmodell funktionieren identisch, egal ob Sie Drohnen oder Lieferwagen verwalten. Nur die Edge-Hardware und die Protokolle ändern sich.

---

## Die Herausforderung

Der Betrieb im Flottenmaßstab erfordert die Lösung von Problemen, die im kleinen Maßstab nicht existieren:

- **Tausende gleichzeitige Telemetriestreams** — Jedes Fahrzeug meldet Position, Geschwindigkeit, Zustand und Sensordaten
- **Befehl und Steuerung im großen Maßstab** — Anweisungen an bestimmte Fahrzeuge oder Gruppen senden, ohne das Netzwerk zu überfluten
- **Digital-Twin-Synchronisation** — Eine genaue Zustandsdarstellung für jedes Fahrzeug in Echtzeit aufrechterhalten
- **Offline-Resilienz** — Fahrzeuge, die weiterarbeiten, wenn die Verbindung abbricht
- **Sicherheitsgarantien** — Sicherstellen, dass Netzwerkausfälle niemals die Fahrzeugsicherheit gefährden

Traditionelle Ansätze – direkte Protokollverbindungen, zentralisierte Datenbanken, Polling-Architekturen – brechen unter diesen Anforderungen zusammen.

---

## Unser Ansatz

Die Architektur kombiniert drei bewährte Technologien:

| Ebene | Technologie | Rolle |
|-------|------------|------|
| **Fahrzeugsteuerung** | Plattformspezifisch (PX4, ECU, ROS 2) | Autonomer Betrieb, Failsafes, manuelle Übersteuerung |
| **Edge Computing** | Jetson + Industrie-PC | Sensorverarbeitung, KI-Inferenz, lokale Entscheidungen |
| **Flotten-Messaging** | NATS JetStream | Pub/Sub, Persistenz, Digital-Twin-Zustand |

Jede Ebene arbeitet unabhängig. Netzwerkausfälle führen zu einer kontrollierten Degradation – Fahrzeuge arbeiten weiter, Edge-Computer verarbeiten weiter, und der Zustand wird synchronisiert, sobald die Verbindung wiederhergestellt ist.

---

## Architekturkomponenten

### Unterstützte Plattformen

Die Architektur unterstützt mehrere Fahrzeugtypen mit plattformspezifischer Hardware und Protokollen:

| Plattform | Protokoll | Steuerungssystem | Anwendungsfälle |
|----------|----------|----------------|-----------|
| **Drohnen** | MAVLink | PX4/ArduPilot | Inspektion, Kartierung, Lieferung |
| **Bodenfahrzeuge** | CAN-Bus / J1939 | ECU / Autoware | Logistik, Bergbau, Landwirtschaft |
| **ROS-2-Roboter** | ROS-2-Topics | Individuell | Lager, Forschung, AGVs |

[Unterstützte Plattformen →]({{< relref "/fleet/platforms" >}})

### Hardware-Stack (Drohnen)

Standard-Hardware, die Leistungsfähigkeit, Verfügbarkeit und Wartbarkeit in Einklang bringt:

- **Rahmen**: Holybro X500 V2 ARF — bewährte Plattform, hervorragende Ersatzteilverfügbarkeit
- **Flugcontroller**: Pixhawk 6X mit PX4 v1.14 LTS
- **Sensor-Begleitrechner**: Raspberry Pi 4 / CM4 für schlanke Sensorintegration
- **KI-Begleitrechner**: NVIDIA Jetson (Orin/Xavier/Nano) für Computer Vision und Inferenz

[Drohnen-Hardware-Details →]({{< relref "/fleet/platforms/drones" >}}) | [Details zu Bodenfahrzeugen →]({{< relref "/fleet/platforms/ground" >}})

### Software-Stack

Open-Source-Software auf jeder Ebene:

- **PX4 v1.14 LTS** — Flugsteuerung mit bewährter Stabilität
- **Ubuntu Server 22.04** — Einheitliche Linux-Umgebung auf allen Begleitrechnern
- **ROS 2 Humble** — Roboter-Middleware für die Sensorintegration
- **Go** — Implementierung des Vehicle Gateway

[Software-Details →]({{< relref "/fleet/software" >}})

### NATS-Architektur

Hierarchische Messaging-Topologie, ausgelegt für den WAN-Einsatz:

- **Leaf Nodes** auf jedem Fahrzeug — lokales Pub/Sub, Store-and-Forward
- **Regionale Hub-Cluster** — fassen Fahrzeuge nach geografischer Region zusammen
- **Globales Mirroring** — regionsübergreifende Replikation bei Bedarf

[NATS-Topologie →]({{< relref "/fleet/nats-topology" >}})

### Digital-Twin-Design

JetStream-Streams und KV-Stores verwalten den Flottenzustand:

- **Subject-Hierarchie** für Routing und Filterung
- **State-Streams** für die Telemetrie-Aggregation
- **Event-Streams** für Audit-Trails
- **Shadow-Stores** für den Abgleich von Soll- und Ist-Zustand

[Subject-Benennung →]({{< relref "/fleet/subjects" >}}) | [Stream-Konfiguration →]({{< relref "/fleet/streams" >}})

### Vehicle Gateway

Dienst, der auf dem Edge-Computer jedes Fahrzeugs läuft und native Protokolle mit NATS verbindet:

- Protokollübersetzung (MAVLink, CAN-Bus, ROS 2)
- Downsampling und Aggregation des Zustands
- Extraktion von Ereignissen aus der Telemetrie
- Befehlsausführung mit Durchsetzung von Richtlinien
- Abgleich des Shadow-Zustands

[Gateway-Design →]({{< relref "/fleet/gateway" >}})

### Sicherheitsmodell

Für die Fahrzeugsicherheit wird der Netzwerkverbindung niemals vertraut:

- **Manuelle Übersteuerung hat oberste Autorität** — Der Bediener hat immer die Kontrolle (RC-Fernsteuerung bei Drohnen, Lenkung/Not-Aus bei Bodenfahrzeugen)
- **Das Fahrzeugsteuerungssystem setzt Failsafes durch** — RTL/Anhalten an Ort und Stelle, Geofence
- **NATS ist niemals Teil des Regelkreises** — Nur Überwachung und Koordination
- **Kontrollierte Degradation** — Der Ausfall von KI oder Netzwerk löst sichere Modi aus

[Sicherheitsdetails →]({{< relref "/fleet/safety" >}})

---

## Warum diese Architektur

### Im großen Maßstab bewährt

Jede Komponente wird bereits in Produktionsumgebungen eingesetzt:

- **NATS** betreibt die globale Messaging-Infrastruktur von Synadia
- **PX4** fliegt weltweit auf Tausenden kommerzieller Drohnen
- **Jetson** führt Inferenz in autonomen Fahrzeugen und Robotern aus

### Offene Standards

Kein Vendor-Lock-in:

- MAVLink ist ein offenes Protokoll mit mehreren Implementierungen
- NATS ist Open Source, kommerzieller Support ist verfügbar
- PX4 läuft auf mehreren Flugcontroller-Hardwareplattformen

### Ihre Infrastruktur oder unsere

NATS JetStream ist **zu 100 % Open Source** (Apache 2.0). Betreiben Sie es selbst oder verbinden Sie sich mit unserer verwalteten Infrastruktur – gleiches Protokoll, gleicher Code, Ihre Wahl.

- **Selbst gehostet** — Bereitstellung auf Ihrer Infrastruktur mit unseren Referenzkonfigurationen
- **Managed** — Verbinden Sie Ihre Leaf Nodes mit unseren regionalen Hubs

Kostenlos für kleine Flotten. Wächst mit Ihnen mit.

[Bereitstellungsoptionen →]({{< relref "/fleet/nats-topology#deployment-options" >}})

### Betriebliche Realität

Für die reale Welt konzipiert:

- Fahrzeuge können ohne Spezialwerkzeug gewartet werden
- Software-Updates werden over-the-air ausgerollt
- Telemetriedaten fließen in Standard-Observability-Stacks
- Der Flottenzustand lässt sich in bestehende Unternehmenssysteme exportieren

---

## Mehr erfahren

Erkunden Sie jede Komponente im Detail:

| Abschnitt | Beschreibung |
|---------|-------------|
| [Unterstützte Plattformen]({{< relref "/fleet/platforms" >}}) | Drohnen, Bodenfahrzeuge, AGVs |
| [Drohnenplattform]({{< relref "/fleet/platforms/drones" >}}) | Rahmen, Flugcontroller, MAVLink |
| [Bodenfahrzeuge]({{< relref "/fleet/platforms/ground" >}}) | Autos, Lkw, CAN-Bus, J1939 |
| [Software]({{< relref "/fleet/software" >}}) | Betriebssysteme, Middleware, Anwendungen |
| [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) | Leaf Nodes, Hubs, WAN-Konnektivität |
| [Subject-Benennung]({{< relref "/fleet/subjects" >}}) | Hierarchische Subject-Struktur |
| [Stream-Konfiguration]({{< relref "/fleet/streams" >}}) | JetStream-Einrichtung für digitale Zwillinge |
| [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) | Brücke von Protokollen zu NATS |
| [Autorisierung & Grants]({{< relref "/fleet/authorization" >}}) | Dezentrale Sicherheit, Zugriff für Dritte |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Failsafes und kontrollierte Degradation |
| [BVLOS-Vorschriften]({{< relref "/fleet/bvlos-regulations" >}}) | Weltweite Vorschriften für Flüge außerhalb der Sichtweite (Beyond Visual Line of Sight) |
| [Luftraum-Viewer]({{< relref "/fleet/airspace-demo" >}}) | Interaktive Visualisierung des US-Luftraums |

---

## Jetzt starten

Sie bauen eine autonome Fahrzeugflotte auf? Sie setzen Drohnen, Lkw oder AGVs im großen Maßstab ein? Wir unterstützen Sie bei Architektur, Implementierung und Betrieb.

[Kontakt →](/contact)
