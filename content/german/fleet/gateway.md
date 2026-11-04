---
title: "Fahrzeug-Gateway"
meta_title: "Fahrzeug-Gateway-Design für autonome Flotten | Ubuntu Software"
description: "Vehicle-Gateway-Architektur: Protokollübersetzung (MAVLink, CAN, ROS 2), Zustands-Downsampling, Ereignisextraktion, Befehlsausführung und Shadow-Abgleich."
image: "/images/robotics.svg"
draft: false
---

## Die Brücke zwischen Fahrzeugprotokollen und NATS

Das Vehicle Gateway ist ein Dienst, der auf dem Edge-Computer jedes Fahrzeugs läuft und zwischen fahrzeugeigenen Protokollen (MAVLink, CAN-Bus, ROS 2) und dem für die Flottenkoordination verwendeten NATS-Messaging-System übersetzt.

---

## Warum ein Gateway?

Fahrzeugprotokolle und Flotten-Messaging dienen unterschiedlichen Zwecken:

| Aspekt | Fahrzeugprotokoll | NATS |
|--------|------------------|------|
| **Umfang** | Einzelnes Fahrzeug | Flottenweit |
| **Protokoll** | Binär, kompakt | JSON, menschenlesbar |
| **Rate** | 100+ Hz Telemetrie | Für WAN reduziert |
| **Persistenz** | Keine | JetStream-Streams |
| **Muster** | Request-Response | Pub/Sub mit Persistenz |

Das Gateway verbindet diese Welten und übernimmt Protokollübersetzung, Ratenbegrenzung und Zustandsverwaltung.

---

## Unterstützte Protokolle

| Protokoll | Fahrzeugtyp | Integration |
|----------|--------------|-------------|
| **MAVLink** | Drohnen (PX4, ArduPilot) | UDP/Seriell, Nachrichten-Parsing |
| **CAN-Bus** | Autos, Lkw | SocketCAN, DBC-Dekodierung |
| **J1939** | Schwerfahrzeuge | PGN/SPN-Mapping |
| **ROS 2** | Autonome Plattformen | Topic-Abonnement |

Jedes Protokoll erfordert einen protokollspezifischen Adapter, die Kernlogik des Gateways (Downsampling, Ereignisse, Befehle, Shadow) wird jedoch gemeinsam genutzt.

---

## Kernaufgaben

### 1. Protokoll-Ingest

Fahrzeugeigene Protokollnachrichten empfangen und parsen. Beispiel für MAVLink (Drohnen):

```go
// Receive MAVLink frames from mavlink-router
conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 14550})

for {
    buf := make([]byte, 1024)
    n, _, _ := conn.ReadFromUDP(buf)

    frame, _ := mavlink.Parse(buf[:n])

    switch msg := frame.Message().(type) {
    case *mavlink.Heartbeat:
        handleHeartbeat(msg)
    case *mavlink.GlobalPositionInt:
        handlePosition(msg)
    case *mavlink.BatteryStatus:
        handleBattery(msg)
    // ... handle other message types
    }
}
```

**Wichtige Nachrichtentypen (MAVLink-Beispiel):**

| MAVLink-Nachricht | Inhalt |
|-----------------|---------|
| `HEARTBEAT` | Modus, Scharfschaltstatus, Systemstatus |
| `GLOBAL_POSITION_INT` | Breite, Länge, Höhe, Geschwindigkeit |
| `ATTITUDE` | Rollen, Nicken, Gieren |
| `BATTERY_STATUS` | Spannung, Strom, Restkapazität |
| `GPS_RAW_INT` | GPS-Fix, Satelliten, HDOP |
| `SYS_STATUS` | CPU-Last, Fehler, Zustand |

**CAN-Bus-Äquivalent (Bodenfahrzeuge):**

| CAN-Signal | Inhalt |
|------------|---------|
| Motordrehzahl | Aktuelle Drehzahl des Motors |
| Fahrzeuggeschwindigkeit | Geschwindigkeit über Grund aus Radsensoren |
| GPS-Position | Breite/Länge vom Telematikmodul |
| Kraftstoffstand | Tankfüllstand in Prozent |
| DTC-Codes | Diagnose-Fehlercodes |

### 2. Zustands-Downsampling

Telemetrierate für die WAN-Übertragung reduzieren:

```go
// Position: downsample from 10Hz to 1Hz
positionSampler := NewDownsampler(100*time.Millisecond, func(msg *Position) {
    // Average positions over window
    return averagePosition(msg)
})

// Battery: only publish on change
batterySampler := NewChangeFilter(func(old, new *Battery) bool {
    return math.Abs(old.Remaining - new.Remaining) > 0.01
})
```

**Downsampling-Strategien:**

| Strategie | Anwendungsfall | Beispiel |
|----------|----------|---------|
| **Zeitbasiert** | Periodischer Zustand | Position mit 1Hz |
| **Änderungsbasiert** | Diskrete Werte | Modusänderungen |
| **Schwellenwert** | Graduelle Änderungen | Akku bei Δ > 1% |
| **Aggregation** | Hochfrequente Daten | Min/Max/Mittelwert über ein Zeitfenster |

Telemetrie mit voller Rate verbleibt für Wahrnehmungssysteme im Fahrzeug. Nur reduzierte Daten werden über das WAN übertragen.

### 3. Ereignisextraktion

Ereignisse aus Zustandsübergängen generieren:

```go
type EventDetector struct {
    previousState *VehicleState
}

func (e *EventDetector) Process(current *VehicleState) []Event {
    var events []Event

    // Detect arm state change
    if current.Armed && !e.previousState.Armed {
        events = append(events, Event{
            Type:      "armed",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"reason": "manual"},
        })
    }

    // Detect mode change
    if current.Mode != e.previousState.Mode {
        events = append(events, Event{
            Type:      "mode_change",
            Timestamp: time.Now(),
            Data: map[string]interface{}{
                "from": e.previousState.Mode,
                "to":   current.Mode,
            },
        })
    }

    // Detect failsafe
    if current.FailsafeActive && !e.previousState.FailsafeActive {
        events = append(events, Event{
            Type:      "failsafe",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"type": current.FailsafeType},
        })
    }

    e.previousState = current
    return events
}
```

**Erkannte Ereignisse:**

| Ereignis | Auslöser |
|-------|---------|
| `armed` | Scharfschaltstatus false → true |
| `disarmed` | Scharfschaltstatus true → false |
| `mode_change` | Flugmoduswechsel |
| `takeoff` | In-Air-Status false → true |
| `landed` | In-Air-Status true → false |
| `failsafe` | Failsafe aktiviert |
| `geofence` | Geofence-Verletzung erkannt |
| `battery.low` | Akku unter Schwellenwert |
| `battery.critical` | Kritischer Akkustand |

### 4. Befehlsausführung mit Richtlinien

Befehle von NATS empfangen und über MAVLink ausführen:

```go
func (g *Gateway) handleCommand(cmd *Command) *CommandAck {
    // Validate command is allowed
    if !g.policy.Allows(cmd) {
        return &CommandAck{
            Status: "rejected",
            Error:  "command not allowed by policy",
        }
    }

    // Convert to MAVLink command
    mavCmd := cmd.ToMAVLink()

    // Send to flight controller
    if err := g.mavlink.Send(mavCmd); err != nil {
        return &CommandAck{
            Status: "failed",
            Error:  err.Error(),
        }
    }

    // Wait for MAVLink ACK
    mavAck, err := g.mavlink.WaitAck(mavCmd.CommandID, 5*time.Second)
    if err != nil {
        return &CommandAck{
            Status: "timeout",
            Error:  "no response from flight controller",
        }
    }

    return &CommandAck{
        Status: mavAck.Result.String(),
    }
}
```

**Durchsetzung von Richtlinien:**

| Richtlinie | Beschreibung |
|--------|-------------|
| **Geofence** | Goto-Befehle außerhalb der Begrenzung ablehnen |
| **Höhenbegrenzung** | Maximale Höhe in Befehlen begrenzen |
| **Modusbeschränkungen** | Bestimmte Moduswechsel verbieten |
| **Ratenbegrenzung** | Befehlsflut verhindern |
| **Authentifizierung** | Befehlsquelle verifizieren |

### 5. Shadow-Abgleich

Soll-Zustand mit dem tatsächlichen Fahrzeugzustand synchronisieren:

```go
func (g *Gateway) reconcileLoop() {
    ticker := time.NewTicker(1 * time.Second)

    for range ticker.C {
        // Read desired state from KV
        desired, _ := g.kv.Get(g.desiredKey())

        // Compare with actual state
        actual := g.getCurrentState()

        // Generate commands to reconcile
        commands := g.reconcile(desired, actual)

        for _, cmd := range commands {
            g.executeCommand(cmd)
        }

        // Update reported state
        reported := g.buildReportedState(actual)
        g.kv.Put(g.reportedKey(), reported)
    }
}

func (g *Gateway) reconcile(desired, actual *State) []*Command {
    var commands []*Command

    // Mode reconciliation
    if desired.Mode != actual.Mode {
        commands = append(commands, &Command{
            Type: "set_mode",
            Data: map[string]interface{}{"mode": desired.Mode},
        })
    }

    // Geofence reconciliation
    if desired.GeofenceEnabled != actual.GeofenceEnabled {
        commands = append(commands, &Command{
            Type: "set_geofence",
            Data: map[string]interface{}{"enabled": desired.GeofenceEnabled},
        })
    }

    return commands
}
```

**Der Shadow-Zustand ermöglicht deklaratives Management:**

- Flottenbetreiber legt Soll-Zustand fest
- Gateway erkennt Abweichungen
- Gateway sendet Befehle zur Konvergenz
- Gateway meldet Ist-Zustand
- Kontinuierlich wiederholen

---

## Architektur

```
┌────────────────────────────────────────────────────────────────────┐
│                         Vehicle Gateway                             │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐     │
│  │   Protocol   │    │    State     │    │    NATS Client   │     │
│  │   Adapter    │───▶│   Machine    │───▶│    Publisher     │     │
│  └──────────────┘    └──────────────┘    └──────────────────┘     │
│         │                   │                      │               │
│         │            ┌──────▼──────┐               │               │
│         │            │   Event     │               │               │
│         │            │  Detector   │───────────────┤               │
│         │            └─────────────┘               │               │
│         │                                          │               │
│         │            ┌─────────────┐               │               │
│         │            │   Shadow    │◀──────────────┤               │
│         │            │ Reconciler  │               │               │
│         │            └──────┬──────┘               │               │
│         │                   │                      │               │
│         │            ┌──────▼──────┐    ┌─────────▼────────┐      │
│         │            │  Command    │    │   NATS Client    │      │
│         │◀───────────│  Executor   │◀───│   Subscriber     │      │
│         │            └─────────────┘    └──────────────────┘      │
│  ┌──────▼──────┐                                                   │
│  │   Protocol  │                                                   │
│  │   Sender    │                                                   │
│  └─────────────┘                                                   │
└────────────────────────────────────────────────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────┐                       ┌─────────────────┐
│  Vehicle Control│                       │   NATS Leaf     │
│ (Pixhawk/ECU/   │                       │   (localhost)   │
│  ROS 2 node)    │                       │                 │
└─────────────────┘                       └─────────────────┘
```

**Protokolladapter:**
- **MAVLink-Adapter** — Für Drohnen (PX4, ArduPilot)
- **CAN-Adapter** — Für Autos/Lkw (SocketCAN + DBC)
- **ROS-2-Adapter** — Für ROS-basierte Plattformen

---

## Implementierungsdetails

### Abhängigkeiten

```go
import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "github.com/bluenviern/go-mavlink/v2"
)
```

**Kernbibliotheken:**

| Bibliothek | Zweck |
|---------|---------|
| **nats.go** | NATS-Client, JetStream, KV |
| **go-mavlink** | MAVLink-Protokollimplementierung |
| **slog** | Strukturiertes Logging |

### Konfiguration

```yaml
# gateway.yaml
vehicle_id: VID-001
environment: prod

mavlink:
  local_addr: ":14550"
  system_id: 1
  component_id: 1

nats:
  url: "nats://localhost:4222"
  credentials: "/etc/gateway/vehicle.creds"

sampling:
  position_hz: 1
  attitude_hz: 1
  battery_change_threshold: 0.01

policy:
  max_altitude: 120
  geofence_file: "/etc/gateway/geofence.json"
```

### Bereitstellung

```yaml
# systemd service
[Unit]
Description=Vehicle Gateway
After=network.target nats.service

[Service]
Type=simple
ExecStart=/usr/local/bin/vehicle-gateway --config /etc/gateway/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## Fehlerbehandlung

### MAVLink-Fehler

| Fehler | Reaktion |
|-------|----------|
| **Verbindung verloren** | Erneut verbinden mit Backoff |
| **Parse-Fehler** | Protokollieren und Frame überspringen |
| **Timeout** | Erneut versuchen oder Befehl fehlschlagen lassen |

### NATS-Fehler

| Fehler | Reaktion |
|-------|----------|
| **Verbindung verloren** | Lokales NATS läuft weiter, erneute Verbindung zum Hub |
| **Veröffentlichung fehlgeschlagen** | Lokal puffern, erneut versuchen |
| **Stream-Fehler** | Protokollieren, alarmieren, fortfahren |

### Befehlsfehler

| Fehler | Reaktion |
|-------|----------|
| **Richtlinienverstoß** | Sofort ablehnen |
| **MAVLink-Ablehnung** | Fehler-ACK zurückgeben |
| **Timeout** | Timeout-ACK zurückgeben |

---

## Metriken

Das Gateway stellt Prometheus-Metriken bereit:

```
# MAVLink
gateway_mavlink_messages_received_total{type="heartbeat"}
gateway_mavlink_messages_sent_total{type="command_long"}
gateway_mavlink_parse_errors_total

# NATS
gateway_nats_messages_published_total{subject="state.position"}
gateway_nats_messages_received_total{subject="cmd.takeoff"}
gateway_nats_publish_errors_total

# Commands
gateway_commands_received_total{type="takeoff"}
gateway_commands_executed_total{type="takeoff",result="success"}
gateway_command_latency_seconds{type="takeoff"}

# Shadow
gateway_shadow_reconciliations_total
gateway_shadow_commands_issued_total
```

---

## Zusammenfassung

| Aufgabe | Eingabe | Ausgabe |
|----------------|-------|--------|
| **Protokoll-Ingest** | Fahrzeugnachrichten | Geparste Telemetrie |
| **Zustands-Downsampling** | 100Hz-Telemetrie | 1Hz-Zustand |
| **Ereignisextraktion** | Zustandsübergänge | Diskrete Ereignisse |
| **Befehlsausführung** | NATS-Befehle | Fahrzeugbefehle |
| **Shadow-Abgleich** | Soll-Zustand | Konvergenzbefehle |

Das Vehicle Gateway ist die entscheidende Komponente, die Operationen im Flottenmaßstab ermöglicht und gleichzeitig die Sicherheitsgarantien des zugrunde liegenden Fahrzeugsteuerungssystems bewahrt.

---

## Verwandte Dokumentation

- [Unterstützte Plattformen]({{< relref "/fleet/platforms" >}}) — Übersicht über alle Fahrzeugtypen
- [Drohnenplattform]({{< relref "/fleet/platforms/drones" >}}) — Details zur MAVLink-Integration
- [Bodenfahrzeuge]({{< relref "/fleet/platforms/ground" >}}) — CAN-Bus- und J1939-Integration

---

## Weiter

[Sicherheitsmodell →]({{< relref "/fleet/safety" >}})
