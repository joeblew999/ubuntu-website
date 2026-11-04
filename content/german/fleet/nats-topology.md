---
title: "NATS-Topologie"
meta_title: "NATS-JetStream-Topologie für Drohnenflotten | Ubuntu Software"
description: "Hierarchische NATS-Topologie für eine Flotte von 1.000 Drohnen: Leaf-Nodes pro Fahrzeug, regionale Hub-Cluster, WAN-Isolation und flottenweite Wiedergabe digitaler Zwillinge."
image: "/images/robotics.svg"
draft: false
---

## NATS-Architektur für skalierbare Flotten

Die Messaging-Topologie bestimmt, wie Telemetriedaten fließen, wo Zustände persistiert werden und wie sich das System bei nachlassender Konnektivität verhält. Diese Architektur skaliert von Dutzenden bis zu Tausenden von Fahrzeugen und gewährleistet dabei niedrige Latenz und hohe Zuverlässigkeit.

---

## Die Herausforderung

Flotten-Messaging muss konkurrierende Anforderungen erfüllen:

| Anforderung | Herausforderung |
|-------------|-----------|
| **Niedrige Latenz** | Fahrzeug-zu-Fahrzeug-Kommunikation benötigt Millisekunden, nicht Sekunden |
| **WAN-Resilienz** | Mobilfunkverbindungen brechen ab, Latenzspitzen treten auf, die Bandbreite schwankt |
| **Zustandspersistenz** | Digital-Twin-Daten müssen Neustarts und Wiederverbindungen überstehen |
| **Skalierung** | Tausende gleichzeitige Publisher und Subscriber |
| **Isolation** | Fahrzeugausfälle dürfen sich nicht auf die gesamte Flotte ausbreiten |

Traditionelle Architekturen erzwingen Kompromisse. Zentralisierte Broker erhöhen die Latenz. Peer-to-Peer-Meshes skalieren nicht. Reine Edge-Ansätze verlieren die Sicht auf die Flotte.

---

## Hierarchische Topologie

Die Lösung ist eine **dreistufige Hierarchie**:

```
┌─────────────────────────────────────────────────────────────────┐
│                         GLOBAL (optional)                        │
│                    Cross-region replication                      │
│                    Fleet-wide aggregation                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  REGIONAL HUB   │ │  REGIONAL HUB   │ │  REGIONAL HUB   │
│   Americas      │ │   Europe        │ │   Asia-Pacific  │
│  3-node cluster │ │  3-node cluster │ │  3-node cluster │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
    ┌────┼────┐         ┌────┼────┐         ┌────┼────┐
    │    │    │         │    │    │         │    │    │
    ▼    ▼    ▼         ▼    ▼    ▼         ▼    ▼    ▼
┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐
│LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF │
│VID-1││VID-2││VID-3││VID-4││VID-5││VID-6││VID-7││VID-8││VID-9│
└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘
```

---

## Stufe 1: Fahrzeug-Leaf-Nodes

Jede Drohne betreibt einen **NATS-Leaf-Node** auf dem Jetson:

| Aspekt | Konfiguration |
|--------|---------------|
| **Prozess** | nats-server im Leaf-Modus |
| **Speicher** | Lokales JetStream auf SSD/eMMC |
| **Verbindung** | Einzelne Upstream-Verbindung zum regionalen Hub |
| **Ressourcen** | ~50 MB RAM, minimale CPU-Last |

**Was Leaf-Nodes bieten:**

### Lokales Pub/Sub

Fahrzeuginterne Dienste kommunizieren über die lokale NATS-Instanz:

```
Vehicle Gateway ──publish──▶ fleet.prod.veh.VID-001.state.position
                                       │
                                       ▼
                            ┌──────────────────┐
                            │   Local NATS     │
                            │   (Leaf Node)    │
                            └──────────────────┘
                                       │
                                       ▼
AI Perception ◀──subscribe── fleet.prod.veh.VID-001.state.position
```

Nachrichten zwischen fahrzeuginternen Diensten verlassen das Fahrzeug nie. Die Latenz liegt im Sub-Millisekundenbereich.

### Store-and-Forward

Wenn die Upstream-Verbindung abbricht, verhält sich der Leaf-Node wie folgt:

1. Er nimmt weiterhin Publishes von lokalen Diensten an
2. Er speichert Nachrichten im lokalen JetStream
3. Er stellt ausgehende Nachrichten für den Hub in eine Warteschlange
4. Er verbindet sich neu und spielt die Nachrichten erneut ab, sobald die Konnektivität zurückkehrt

Das bedeutet: **Fahrzeuge arbeiten auch bei Verbindungsverlust weiter**. Der Zustand wird synchronisiert, sobald die Verbindung wiederhergestellt ist.

### WAN-Isolation

Der Leaf-Node fungiert als Puffer zwischen Fahrzeug und WAN:

- **Backpressure** — Wenn der Hub nicht mithalten kann, puffert der Leaf-Node lokal
- **Filterung** — Nur abonnierte Subjects werden über die WAN-Verbindung übertragen
- **Komprimierung** — Das NATS-Protokoll komprimiert effizient
- **Wiederverbindung** — Automatische Wiederverbindung mit exponentiellem Backoff

---

## Stufe 2: Regionale Hub-Cluster

Regionale Hubs bündeln Fahrzeuge nach geografischem Gebiet:

| Aspekt | Konfiguration |
|--------|---------------|
| **Bereitstellung** | NATS-Cluster mit 3 Knoten (Minimum) |
| **Standort** | Cloud-Region oder Edge-Rechenzentrum |
| **Speicher** | JetStream auf schnellen SSDs |
| **Kapazität** | Hunderte von Leaf-Verbindungen pro Cluster |

**Warum regional:**

- **Latenz** — Leaf-Nodes verbinden sich mit nahegelegenen Hubs
- **Bandbreite** — Telemetriedaten bleiben standardmäßig in der Region
- **Compliance** — Anforderungen an die Datenresidenz
- **Fehlerisolation** — Regionale Ausfälle betreffen keine anderen Regionen

**Aufgaben des Hubs:**

| Funktion | Beschreibung |
|----------|-------------|
| **Leaf-Aggregation** | Nimmt Verbindungen von Fahrzeug-Leaf-Nodes an |
| **Stream-Speicherung** | Persistiert Digital-Twin-Daten in JetStream |
| **Consumer-Unterstützung** | Bedient Flotten-Dashboards, APIs und Analysen |
| **Regionsübergreifendes Routing** | Leitet bei Bedarf an die globale Stufe weiter |

### Cluster-Konfiguration

Ein Cluster mit 3 Knoten bietet:

- **Hochverfügbarkeit** — Übersteht den Ausfall eines einzelnen Knotens
- **Datenreplikation** — JetStream-Streams werden über die Knoten hinweg repliziert
- **Lastverteilung** — Leaf-Verbindungen werden gleichmäßig auf die Knoten verteilt

```yaml
# Example hub cluster configuration
cluster:
  name: hub-americas
  routes:
    - nats://hub-americas-1:6222
    - nats://hub-americas-2:6222
    - nats://hub-americas-3:6222

jetstream:
  store_dir: /data/jetstream
  max_memory_store: 4GB
  max_file_store: 100GB

leafnodes:
  port: 7422
  authorization:
    users:
      - user: vehicle
        password: $VEHICLE_PASSWORD
        allowed_connection_types: ["LEAFNODE"]
```

---

## Stufe 3: Global (optional)

Für Flotten, die sich über mehrere Regionen erstrecken, ermöglicht eine globale Stufe:

| Funktion | Beschreibung |
|----------|-------------|
| **Regionsübergreifendes Mirroring** | Repliziert Streams zwischen regionalen Hubs |
| **Globale Aggregation** | Flottenweite Dashboards und Analysen |
| **Befehls-Routing** | Sendet Befehle an Fahrzeuge in jeder beliebigen Region |
| **Disaster Recovery** | Failover zwischen Regionen |

**Implementierungsoptionen:**

1. **NATS Supercluster** — Föderiert regionale Cluster über Gateway-Verbindungen
2. **JetStream Mirroring** — Repliziert bestimmte Streams an einen zentralen Standort
3. **Anwendungsebene** — Individuelle Synchronisation zwischen regionalen APIs

Die meisten Flotten operieren innerhalb einer einzigen Region und benötigen die globale Stufe anfangs nicht.

---

## Vorteile dieser Topologie

### Lokales Pub/Sub mit niedriger Latenz

Fahrzeuginterne Kommunikation erfolgt im **Sub-Millisekundenbereich**:

```
MAVLink message received
    └─▶ Vehicle Gateway publishes to local NATS
        └─▶ AI Perception subscribes from local NATS
            └─▶ Total latency: <1ms
```

Kein WAN-Roundtrip für fahrzeuginterne Kommunikation.

### WAN-Isolation

Der Fahrzeugbetrieb **läuft auch bei Verbindungsverlust weiter**:

- Perzeptionssysteme laufen weiter
- Zustände werden lokal gesammelt
- Befehle werden zur Ausführung in eine Warteschlange gestellt
- Die Wiederverbindung synchronisiert automatisch

### Flottenweites Digital-Twin-Replay

JetStream ermöglicht **zeitbezogene Abfragen** über die gesamte Flotte:

```
# Replay vehicle state from 10 minutes ago
nats stream get TWIN_STATE --start-time="2024-01-15T10:00:00Z"

# Subscribe to real-time state updates
nats sub "fleet.prod.veh.*.state.position"
```

Der Zustandsverlauf jedes Fahrzeugs ist über den regionalen Hub abfragbar.

### Skalierbare Architektur

Die Topologie skaliert horizontal:

| Skalierungspunkt | Ansatz |
|-------------|----------|
| **Mehr Fahrzeuge** | Leaf-Verbindungen zum bestehenden Hub hinzufügen |
| **Höherer Durchsatz** | Knoten zum Hub-Cluster hinzufügen |
| **Mehr Regionen** | Zusätzliche regionale Hubs bereitstellen |
| **Globale Reichweite** | Hubs über einen Supercluster verbinden |

---

## Verbindungsablauf

Wenn ein Fahrzeug eingeschaltet wird:

```
1. Jetson boots, NATS leaf node starts
2. Leaf node reads hub address from config
3. TLS connection established to regional hub
4. Authentication via credentials
5. Leaf node advertises local subjects
6. Hub routes relevant subscriptions to leaf
7. Bidirectional message flow begins
```

Wenn die Verbindung abbricht:

```
1. Leaf node detects connection loss
2. Local pub/sub continues uninterrupted
3. Outbound messages queue locally
4. Reconnection attempts with backoff
5. On reconnect, queued messages replay
6. Stream state synchronizes
```

---

## Sicherheit

Jede Verbindung ist authentifiziert und verschlüsselt:

| Verbindung | Sicherheit |
|------------|----------|
| **Leaf → Hub** | TLS 1.3, Authentifizierung über Credentials |
| **Hub-Cluster** | TLS, authentifizierte Cluster-Routen |
| **Hub → Global** | TLS, Gateway-Authentifizierung |

Credentials werden pro Fahrzeug bereitgestellt und ermöglichen:

- Widerruf des Zugriffs einzelner Fahrzeuge
- Audit-Trails pro Fahrzeug
- Rate Limiting pro Verbindung

Das Autorisierungsmodell unterstützt **dezentrale Sicherheit mit Grants** — Dritte können eingeschränkten, zeitlich begrenzten Zugriff ohne zentrale Koordination erhalten. Siehe [Autorisierung & Grants]({{< relref "/fleet/authorization" >}}) für Details zur Verwaltung von Credentials und zur Integration von Drittanbietern.

---

## Zusammenfassung

| Stufe | Zweck | Bereitstellung |
|------|---------|------------|
| **Leaf (Fahrzeug)** | Lokales Pub/Sub, Store-and-Forward | Auf dem Jetson |
| **Hub (regional)** | Aggregation, Persistenz, Consumer | Cloud-/Edge-Rechenzentrum |
| **Global (optional)** | Regionsübergreifende, flottenweite Ansichten | Zentrale Cloud |

Diese Topologie stellt sicher, dass Fahrzeuge unabhängig arbeiten, während flottenweite Transparenz und Kontrolle erhalten bleiben.

---

## Bereitstellungsoptionen

### Open-Source-Grundlage

NATS JetStream ist **zu 100 % Open Source** unter der Apache-2.0-Lizenz:

- **Quellcode**: [github.com/nats-io/nats-server](https://github.com/nats-io/nats-server)
- **Keine kommerzielle Lizenz erforderlich** — frei im Produktivbetrieb nutzbar
- **Vollständige Funktionsparität** — die Open-Source-Version enthält alle Funktionen
- **Aktive Entwicklung** — unterstützt von Synadia, von Tausenden genutzt

Sie betreiben exakt denselben Code wie wir. Kein Lock-in.

---

### Option 1: Selbst gehostet

Stellen Sie den gesamten Stack auf Ihrer eigenen Infrastruktur bereit:

| Komponente | Ihre Verantwortung |
|-----------|---------------------|
| Hub-Cluster | Bereitstellen, konfigurieren, warten |
| Leaf-Nodes | Auf jedem Fahrzeug bereitstellen |
| Monitoring | Observability einrichten |
| Updates | Upgrades verwalten |

**Wir bieten:**
- Referenzkonfigurationen (diese Dokumentation)
- Architekturberatung
- Unterstützung bei der Implementierung

Ideal für: Organisationen mit Ops-Teams, Anforderungen an die Datenresidenz oder bestehender Infrastruktur.

---

### Option 2: Anbindung an unsere Infrastruktur

Nutzen Sie die verwaltete NATS-Infrastruktur von Ubuntu Software:

| Komponente | Verantwortung |
|-----------|----------------|
| Hub-Cluster | **Verwalten wir** |
| Leaf-Nodes | Stellen Sie auf den Fahrzeugen bereit |
| Monitoring | **Wir stellen Dashboards bereit** |
| Skalierung | **Übernehmen wir** |

**So funktioniert es:**
1. Wir stellen Credentials für Ihre Flotte bereit
2. Ihre Leaf-Nodes verbinden sich mit unseren regionalen Hubs
3. Ihre Daten bleiben isoliert (dedizierte Subjects)
4. Sie erhalten Flotten-Dashboards und API-Zugriff

**Vorteile:**
- Keine Infrastruktur, die aufgebaut oder gewartet werden muss
- Vorab optimiert für die Muster von Drohnenflotten
- Globale Abdeckung über mehrere Regionen
- Konzentrieren Sie sich auf Ihre Drohnen, nicht auf Ihr Messaging

**Preise:**
- **Free Tier** — Bis zu 10 Fahrzeuge, ideal für Entwicklung und kleine Bereitstellungen
- **Scale Tier** — Preise pro Fahrzeug, die mit Ihrem Wachstum mitwachsen

[Kontakt aufnehmen →](/contact), um loszulegen.

---

### Vergleich

| Aspekt | Selbst gehostet | Verwaltet |
|--------|-------------|---------|
| **Einrichtungszeit** | Tage bis Wochen | Stunden |
| **Betriebsaufwand** | Sie warten | Wir warten |
| **Kostenmodell** | Infrastruktur + Zeit | Pro Fahrzeug |
| **Anpassbarkeit** | Volle Kontrolle | Standardkonfiguration |
| **Datenstandort** | Ihre Wahl | Unsere Regionen |
| **Ideal für** | Große Ops-Teams | Schnelle Bereitstellung |

---

## Weiter

[Subject-Benennung →]({{< relref "/fleet/subjects" >}})
