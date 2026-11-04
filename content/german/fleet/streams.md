---
title: "JetStream-Konfiguration"
meta_title: "JetStream-Streams für digitale Zwillinge von Drohnenflotten | Ubuntu Software"
description: "JetStream-Stream- und KV-Store-Konfiguration für digitale Zwillinge von Flotten: Zustandsaggregation, Ereignis-Audit-Trails, Befehlswarteschlangen und Shadow-State."
image: "/images/robotics.svg"
draft: false
---

## JetStream für digitale Zwillinge

JetStream stellt die Persistenzschicht für die digitalen Zwillinge der Flotte bereit. Jeder Stream-Typ erfüllt einen bestimmten Zweck mit maßgeschneiderten Einstellungen für Aufbewahrung und Replikation.

---

## Stream-Architektur

Fünf JetStream-Ressourcen unterstützen das Muster des digitalen Zwillings:

| Ressource | Typ | Zweck |
|----------|------|---------|
| **TWIN_STATE** | Stream | Telemetrie mit Rollup |
| **TWIN_EVENTS** | Stream | Dauerhafter Ereignis-Audit-Trail |
| **TWIN_CMD** | Stream | Befehlswarteschlange |
| **TWIN_SHADOW** | KV-Store | Gewünschter/gemeldeter Zustand |
| **TWIN_BLOBS** | Object Store | Große Artefakte |

---

## TWIN_STATE: Telemetrie-Stream

Speichert den kontinuierlichen Fahrzeugzustand mit **Subject-basiertem Rollup**:

```yaml
name: TWIN_STATE
subjects:
  - "fleet.prod.veh.*.state.>"
retention: limits
max_age: 1h
max_bytes: 10GB
storage: file
replicas: 3
discard: old
rollup_hdrs: true
```

### Konfiguration erklärt

| Einstellung | Wert | Begründung |
|---------|-------|-----------|
| **retention** | limits | Begrenzt durch Alter und Größe |
| **max_age** | 1h | Aktuellen Zustand für Replay vorhalten |
| **max_bytes** | 10GB | Speicher pro Hub begrenzen |
| **replicas** | 3 | Knotenausfälle überstehen |
| **rollup_hdrs** | true | Rollup pro Subject aktivieren |

### Subject-Rollup

Mit aktiviertem Rollup können Consumer **nur die neueste Nachricht pro Subject** anfordern:

```go
// Get latest position for all vehicles
sub, _ := js.Subscribe("fleet.prod.veh.*.state.position",
    nats.DeliverLastPerSubject())
```

Dadurch wird der Stream von einer Firehose zu einem abfragbaren Zustandsspeicher.

### Consumer-Muster

| Muster | Anwendungsfall |
|---------|----------|
| **DeliverAll** | Historisches Replay, Analysen |
| **DeliverLast** | Snapshot des aktuellen Zustands |
| **DeliverLastPerSubject** | Neuester Zustand pro Fahrzeug |
| **StartTime** | Replay ab einem bestimmten Zeitstempel |

---

## TWIN_EVENTS: Audit-Trail

Speichert diskrete Ereignisse mit **langer Aufbewahrung** für Compliance und Debugging:

```yaml
name: TWIN_EVENTS
subjects:
  - "fleet.prod.veh.*.evt.>"
retention: limits
max_age: 90d
max_bytes: 100GB
storage: file
replicas: 3
discard: old
```

### Konfiguration erklärt

| Einstellung | Wert | Begründung |
|---------|-------|-----------|
| **max_age** | 90d | Einhaltung regulatorischer Vorgaben, Vorfallsanalyse |
| **max_bytes** | 100GB | Ereignisse sind kleiner als Zustandsdaten |
| **storage** | file | Bleibt über Neustarts hinweg erhalten |

### Ereigniseigenschaften

Ereignisse unterscheiden sich von Zustandsdaten:

| Eigenschaft | Zustand | Ereignisse |
|----------|-------|--------|
| **Häufigkeit** | Kontinuierlich (Hz) | Diskret (bei Änderung) |
| **Aufbewahrung** | Kurz (Stunden) | Lang (Monate) |
| **Consumer** | Echtzeit-Dashboards | Audit, Analysen, Alarmierung |
| **Replay** | Rekonstruktion des aktuellen Zustands | Untersuchung, Compliance |

### Beispielereignisse

```json
// Armed event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "previous_state": "disarmed",
    "armed_by": "operator-123",
    "reason": "mission_start"
  }
}

// Failsafe event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T11:45:00Z",
  "data": {
    "failsafe_type": "low_battery",
    "action_taken": "rtl",
    "battery_remaining": 0.15
  }
}
```

---

## TWIN_CMD: Befehlswarteschlange

Speichert Befehle mit **Workqueue-Aufbewahrung**:

```yaml
name: TWIN_CMD
subjects:
  - "fleet.prod.veh.*.cmd.>"
retention: workqueue
max_age: 5m
storage: file
replicas: 3
```

### Konfiguration erklärt

| Einstellung | Wert | Begründung |
|---------|-------|-----------|
| **retention** | workqueue | Nachrichten werden nach dem Ack gelöscht |
| **max_age** | 5m | Nicht bestätigte Befehle laufen ab |
| **replicas** | 3 | Garantierte Befehlszustellung |

### Workqueue-Semantik

1. Befehl wird im Stream veröffentlicht
2. Vehicle Gateway empfängt ihn (Exactly-once-Zustellung)
3. Gateway verarbeitet den Befehl
4. Gateway bestätigt die Nachricht (Ack)
5. Nachricht wird aus dem Stream entfernt

**Nicht bestätigte Befehle werden erneut zugestellt**, damit bei Verbindungsabbrüchen keine Befehle verloren gehen.

### Befehlsablauf

```
Fleet Operator               TWIN_CMD Stream              Vehicle Gateway
     │                             │                             │
     │ Publish: cmd.takeoff        │                             │
     │ ─────────────────────────▶  │                             │
     │                             │ Deliver to consumer         │
     │                             │ ─────────────────────────▶  │
     │                             │                             │ Process
     │                             │                             │ command
     │                             │             Ack             │
     │                             │ ◀─────────────────────────  │
     │                             │ (message removed)           │
     │                             │                             │
```

---

## TWIN_SHADOW: KV-Store

JetStream-KV-Stores verwalten den Shadow-Zustand des digitalen Zwillings:

```yaml
bucket: TWIN_SHADOW
max_value_size: 1MB
history: 5
ttl: 0
replicas: 3
storage: file
```

### Konfiguration erklärt

| Einstellung | Wert | Begründung |
|---------|-------|-----------|
| **history** | 5 | 5 vorherige Versionen aufbewahren |
| **ttl** | 0 | Kein Ablauf (dauerhaft) |
| **replicas** | 3 | Hochverfügbarkeit |

### Schlüsselstruktur

```
fleet/prod/veh/VID-001/desired    → Desired state JSON
fleet/prod/veh/VID-001/reported   → Reported state JSON
fleet/prod/veh/VID-001/config     → Vehicle configuration
fleet/prod/veh/VID-001/meta       → Metadata (serial, type, etc.)
```

### Shadow-Operationen

```go
// Write desired state
kv.Put("fleet/prod/veh/VID-001/desired", desiredJSON)

// Read reported state
entry, _ := kv.Get("fleet/prod/veh/VID-001/reported")

// Watch for changes
watcher, _ := kv.Watch("fleet/prod/veh/VID-001/desired")
for entry := range watcher.Updates() {
    // Desired state changed, reconcile
}
```

### Historie für das Debugging

Mit aktivierter Historie können Sie frühere Zustände untersuchen:

```go
// Get all historical values
history, _ := kv.History("fleet/prod/veh/VID-001/desired")
for _, entry := range history {
    fmt.Printf("Rev %d at %s: %s\n",
        entry.Revision(), entry.Created(), entry.Value())
}
```

---

## TWIN_BLOBS: Object Store

Große Artefakte, die nicht in Nachrichten passen:

```yaml
bucket: TWIN_BLOBS
max_chunk_size: 128KB
storage: file
replicas: 3
```

### Anwendungsfälle

| Artefakt | Größe | Beispiel |
|----------|------|---------|
| **Missionsdateien** | 10-100KB | Wegpunktdefinitionen |
| **Logdateien** | 1-10MB | Fluglogs, Diagnosedaten |
| **Firmware** | 10-50MB | PX4-Firmware-Images |
| **Karten** | 1-100MB | Offline-Kartenkacheln |
| **ML-Modelle** | 10-500MB | TensorRT-Engine-Dateien |

### Objektoperationen

```go
// Upload mission file
obj, _ := obs.PutFile("missions/mission-456.json", file)

// Download to vehicle
reader, _ := obs.Get("missions/mission-456.json")

// List available firmware
objects := obs.List(nats.ObjectSearchPrefix("firmware/"))
```

### Übertragung in Chunks

Object Store verarbeitet große Dateien durch:

1. Aufteilen in Chunks (standardmäßig 128KB)
2. Speichern der Chunks als Stream-Nachrichten
3. Erfassen von Metadaten (Größe, Hash, Chunks)
4. Zusammensetzen beim Abruf

Dies ermöglicht **fortsetzbare Übertragungen** über unzuverlässige Verbindungen.

---

## Stream-Dimensionierung

### Schätzungen pro Fahrzeug

| Stream | Nachrichtengröße | Rate | Tägliches Volumen |
|--------|--------------|------|--------------|
| **STATE** | 200 Bytes | 10 Hz | ~170 MB |
| **EVENTS** | 500 Bytes | 10/Tag | ~5 KB |
| **CMD** | 300 Bytes | 100/Tag | ~30 KB |

### Schätzungen für die Flotte (1.000 Fahrzeuge)

| Ressource | 1 Stunde | 1 Tag | 90 Tage |
|--------|--------|-------|--------|
| **TWIN_STATE** | 7 GB | - | - |
| **TWIN_EVENTS** | - | 5 MB | 450 MB |
| **TWIN_SHADOW** | - | 1 GB | 1 GB |
| **TWIN_BLOBS** | - | - | 50 GB |

Dies sind Basisschätzungen. Die tatsächliche Nutzung hängt von Nachrichtenraten und -größen ab.

---

## Betriebliche Aspekte

### Stream-Zustand

Überwachen Sie diese Metriken:

| Metrik | Alarmschwelle |
|--------|-----------------|
| **Consumer-Verzögerung** | > 10.000 Nachrichten |
| **Stream-Bytes** | > 80 % des Maximums |
| **Replikat-Synchronisierung** | Beliebiges Replikat im Rückstand |
| **Ausstehende Acks** | > 1.000 pro Consumer |

### Backup-Strategie

| Ressource | Backup-Methode | Häufigkeit |
|----------|---------------|-----------|
| **Streams** | JetStream-Snapshot | Täglich |
| **KV-Store** | Schlüsselexport | Täglich |
| **Object Store** | Backup auf Dateiebene | Wöchentlich |

### Migration

Streams können zwischen Clustern migriert werden:

1. Stream auf dem Ziel mit derselben Konfiguration erstellen
2. Mirror erstellen, der aus dem Original bezieht
3. Auf die Synchronisierung warten
4. Producer/Consumer auf das Ziel umstellen
5. Ursprünglichen Stream entfernen

---

## Zusammenfassung

| Ressource | Typ | Aufbewahrung | Zweck |
|----------|------|-----------|---------|
| **TWIN_STATE** | Stream | 1 Stunde | Telemetrie mit Rollup |
| **TWIN_EVENTS** | Stream | 90 Tage | Audit-Trail |
| **TWIN_CMD** | Stream | Workqueue | Befehlszustellung |
| **TWIN_SHADOW** | KV | Dauerhaft | Gewünschter/gemeldeter Zustand |
| **TWIN_BLOBS** | Object | Dauerhaft | Große Artefakte |

---

## Weiter

[Vehicle Gateway →]({{< relref "/fleet/gateway" >}})
