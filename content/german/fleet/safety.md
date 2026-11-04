---
title: "Sicherheitsmodell"
meta_title: "Sicherheitsarchitektur für autonome Flotten | Ubuntu Software"
description: "Sicherheitsmodell für autonome Fahrzeugflotten: Befugnis zur manuellen Übersteuerung, Ausfallsicherungen im Fahrzeug, Netzwerkisolierung von Regelkreisen und schrittweise Funktionsreduzierung."
image: "/images/robotics.svg"
draft: false
---

## Sicherheit durch Design

Der Betrieb autonomer Fahrzeuge im Flottenmaßstab erfordert eine rigorose Sicherheitsarchitektur. Dieses Design stellt sicher, dass **Netzwerkausfälle, Softwarefehler und Systemstörungen niemals die Fahrzeugsicherheit gefährden**.

Die zentralen Sicherheitsprinzipien gelten für **alle Fahrzeugtypen**—Drohnen, PKW, LKW und AGVs. Die Details unterscheiden sich (RC vs. Lenkrad, RTL vs. Anhalten an Ort und Stelle), aber die Architektur ist einheitlich.

---

## Grundprinzip

**NATS ist niemals Teil des Echtzeit-Regelkreises.**

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SAFETY HIERARCHY                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   1. MANUAL OVERRIDE      ─────────────────────────▶  HIGHEST       │
│      Human control (RC, steering, e-stop), always available          │
│                                                                      │
│   2. VEHICLE FAILSAFES    ─────────────────────────▶  HIGH          │
│      Hardware-enforced, local to vehicle                             │
│                                                                      │
│   3. GATEWAY POLICY       ─────────────────────────▶  MEDIUM        │
│      Software enforcement, on-vehicle                                │
│                                                                      │
│   4. FLEET COMMANDS       ─────────────────────────▶  LOWEST        │
│      Network-dependent, advisory                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

Höhere Schichten können niedrigere Schichten jederzeit übersteuern. Netzwerkabhängige Befehle haben die niedrigste Autorität.

---

## Manuelle Übersteuerung hat oberste Autorität

Der menschliche Bediener verfügt über eine **bedingungslose Übersteuerung**—ob per RC-Sender, Lenkrad oder Not-Aus-Taster.

### Drohnen-Übersteuerung (RC)

| Aspekt | Umsetzung |
|--------|----------------|
| **Protokoll** | ExpressLRS direkt zum Empfänger |
| **Pfad** | TX → RX → Pixhawk (keine Software im Regelkreis) |
| **Modusschalter** | Hardware-Failsafe auf Manual/Stabilize |
| **Kill-Switch** | Sofortiges Disarmen der Motoren |

**RC ist niemals abhängig von:**

- der Verfügbarkeit des Jetson
- der NATS-Verbindung
- dem Vehicle Gateway
- irgendeiner Softwarekomponente

Selbst wenn alle Computer ausfallen, behält der Pilot die volle manuelle Kontrolle.

#### ExpressLRS-Failsafe

Bei Verlust des RC-Signals:

```
1. RX detects signal loss (100ms)
2. RX sends failsafe values to Pixhawk
3. Pixhawk enters RC_LOSS failsafe mode
4. Configured action executes (RTL, land, etc.)
```

Das Failsafe-Verhalten ist in PX4 konfiguriert und hängt von keiner Flottensoftware ab.

### Übersteuerung bei Bodenfahrzeugen

| Aspekt | PKW/LKW | AGVs |
|--------|-------------|------|
| **Primär** | Lenkrad + Pedale | Not-Aus-Taster |
| **Pfad** | Physische Verbindung zu den Aktoren | Hardware-Interrupt an die Motoren |
| **Backup** | Zündung aus / Leerlauf | Funk-Not-Aus |
| **Netzwerk** | Nie im Steuerungspfad | Nie im Steuerungspfad |

**Die Übersteuerung bei Bodenfahrzeugen ist niemals abhängig von:**

- der Verfügbarkeit des Bordcomputers
- der NATS-Verbindung
- irgendeiner autonomen Software

Selbst wenn alle Computer ausfallen, behält der Fahrer/Bediener die physische Kontrolle, oder das Fahrzeug hält sicher an.

---

## Das Fahrzeugsteuerungssystem erzwingt Failsafes

Das Steuerungssystem des Fahrzeugs implementiert Failsafes lokal—kein Netzwerk erforderlich.

### Drohnen-Failsafes (PX4)

Der Flugcontroller implementiert mehrere Failsafes:

#### Failsafe bei RC-Verlust

| Parameter | Wert | Aktion |
|-----------|-------|--------|
| `COM_RC_LOSS_T` | 0.5s | Timeout vor dem Failsafe |
| `NAV_RCL_ACT` | 2 | RTL bei RC-Verlust |

#### Verlust der Datenverbindung

| Parameter | Wert | Aktion |
|-----------|-------|--------|
| `COM_DL_LOSS_T` | 10s | Timeout der GCS-Verbindung |
| `NAV_DLL_ACT` | 0 | Mission fortsetzen (Verbindung nicht kritisch) |

#### Akku-Failsafes

| Stufe | Parameter | Wert | Aktion |
|-------|-----------|-------|--------|
| Niedrig | `BAT_LOW_THR` | 0.25 | Warnung |
| Kritisch | `BAT_CRIT_THR` | 0.15 | RTL |
| Notfall | `BAT_EMERGEN_THR` | 0.07 | Sofort landen |

#### Geofence

| Parameter | Wert | Aktion |
|-----------|-------|--------|
| `GF_ACTION` | 3 | RTL bei Verletzung |
| `GF_MAX_HOR_DIST` | variiert | Horizontale Grenze |
| `GF_MAX_VER_DIST` | 120m | Höhengrenze |

#### Positionsverlust

| Parameter | Wert | Aktion |
|-----------|-------|--------|
| `COM_POS_FS_DELAY` | 1s | Positions-Timeout |
| `COM_POSCTL_NAVL` | 0 | Landen bei Positionsverlust |

**Alle Failsafes werden auf dem Pixhawk ausgeführt**—keine externe Abhängigkeit.

### Failsafes für Bodenfahrzeuge

Bodenfahrzeuge implementieren gleichwertige Sicherheitsmechanismen:

| Auslöser | Reaktion Drohne | Reaktion Bodenfahrzeug |
|---------|---------------|------------------------|
| **Kommunikationsverlust** | Return-to-Launch | An Ort und Stelle anhalten / rechts ranfahren |
| **Sensorausfall** | Sofort landen | Sicher anhalten, Parkbremse aktivieren |
| **Geofence-Verletzung** | RTL an der Grenze | An der Grenze anhalten |
| **Not-Aus ausgelöst** | Motorabschaltung | Sofortige Bremsung |
| **Niedriger Akku-/Kraftstoffstand** | RTL, dann landen | Rückkehr zum Depot / sicher anhalten |
| **Hindernis erkannt** | Ausweichen oder schweben | Anhalten und warten |

**Failsafes für Bodenfahrzeuge umfassen typischerweise:**

- **Lidar-/Radar-Hindernisse** — Notstopp, wenn der Pfad blockiert ist
- **Geschwindigkeitsbegrenzung** — Höchstgeschwindigkeit abhängig von der Umgebung
- **Zonenbeschränkungen** — Sperrzonen werden lokal durchgesetzt
- **Watchdog-Timer** — Anhalten, wenn die Autonomiesoftware hängt
- **Bremsredundanz** — Mehrere unabhängige Bremssysteme

**Alle Failsafes werden lokal ausgeführt**—keine Netzwerkabhängigkeit.

---

## NATS dient ausschließlich der Überwachung

Die Netzwerkschicht übernimmt:

| Funktion | Beschreibung |
|----------|-------------|
| **Telemetrieerfassung** | Fahrzeugzustand für Dashboards |
| **Ereignisprotokollierung** | Audit-Trail in JetStream |
| **Missionskoordination** | Planung auf Flottenebene |
| **Shadow State** | Synchronisierung von Soll- und Ist-Zustand |

Die Netzwerkschicht übernimmt **nicht:**

| Funktion | Warum nicht |
|----------|---------|
| **Lageregelung** | Latenz inakzeptabel, sicherheitskritisch |
| **Motorbefehle** | Müssen lokal auf der FCU erfolgen |
| **Sensorfusion** | Echtzeitanforderungen |
| **Failsafe-Entscheidungen** | Müssen ohne Netzwerk funktionieren |

### Befehlsvalidierung

Wenn Befehle über NATS eintreffen:

```go
func (g *Gateway) validateCommand(cmd *Command) error {
    // Check geofence
    if cmd.Type == "goto" {
        if !g.geofence.Contains(cmd.Target) {
            return errors.New("target outside geofence")
        }
    }

    // Check altitude
    if cmd.Altitude > g.policy.MaxAltitude {
        return errors.New("altitude exceeds limit")
    }

    // Check vehicle state
    if cmd.RequiresArmed && !g.state.Armed {
        return errors.New("vehicle not armed")
    }

    // Check rate limiting
    if g.rateLimiter.Exceeded(cmd.Type) {
        return errors.New("command rate exceeded")
    }

    return nil
}
```

Befehle werden **vor der Ausführung validiert**—doch selbst wenn ein fehlerhafter Befehl ausgeführt wird, bieten die PX4-Failsafes den letzten Schutz.

---

## Kontrollierte Degradation

Das System degradiert sicher, wenn Komponenten ausfallen:

### Ausfall des Jetson

```
Jetson offline
    └─▶ No perception, no AI inference
    └─▶ Pi continues MAVLink routing
    └─▶ RC control fully operational
    └─▶ PX4 failsafes active
    └─▶ Vehicle lands safely or continues mission
```

### NATS-Verbindung unterbrochen

```
WAN disconnected
    └─▶ Leaf node buffers locally
    └─▶ Gateway continues operation
    └─▶ Vehicle follows current mission
    └─▶ State syncs when connection returns
```

### Absturz des Vehicle Gateway

```
Gateway process dies
    └─▶ systemd restarts within 5 seconds
    └─▶ PX4 continues autonomous flight
    └─▶ RC override always available
    └─▶ No fleet commands until restart
```

### Ausfall des Hub-Clusters

```
Regional hub down
    └─▶ Leaf nodes queue messages
    └─▶ Vehicles continue operating
    └─▶ RC and PX4 failsafes unaffected
    └─▶ Fleet visibility lost temporarily
```

### Vermeidung kaskadierender Ausfälle

| Ausfall | Auswirkungsbereich | Gegenmaßnahme |
|---------|--------------|------------|
| Einzelnes Fahrzeug | Nur dieses Fahrzeug | Failsafe zum Landen |
| Hub-Knoten | Keine Auswirkung (repliziert) | Cluster heilt sich selbst |
| Hub-Cluster | Regionale Flottenkoordination | Fahrzeuge arbeiten autonom weiter |
| Globale Ebene | Regionsübergreifende Sichtbarkeit | Regionaler Betrieb läuft weiter |

---

## Betriebliche Sicherheitsverfahren

### Vorflug-Checkliste

Vor dem Scharfschalten eines Flottenfahrzeugs:

1. **RC-Verbindung geprüft** — Pilot bestätigt die Steuerungsautorität
2. **Failsafe-Test** — RTL-/Landeverhalten überprüfen
3. **Geofence geladen** — Grenzen in PX4 bestätigt
4. **Akku-Schwellenwert** — Failsafe-Stufen angemessen
5. **Wettercheck** — Wind innerhalb der Grenzwerte
6. **Luftraum frei** — Keine Konflikte

### Notfallverfahren

| Notfall | Reaktion |
|-----------|----------|
| **Sichtverlust** | RTL per RC aktivieren |
| **Flyaway** | RC-Kill-Switch oder Geofence löst RTL aus |
| **Konflikt mehrerer Fahrzeuge** | Flottenweiter Landebefehl |
| **Vollständiger Kommunikationsverlust** | PX4-Failsafe RC_LOSS |
| **Brand/Absturz** | RC-Disarm, Rettungsdienste |

### Abschaltkette

Mehrere Methoden, um ein Fahrzeug zu stoppen:

```
1. RC kill switch      → Immediate motor stop
2. RC mode switch      → Land mode
3. GCS command         → RTL or land
4. Fleet command       → RTL or land
5. Geofence breach     → Auto-RTL
6. Battery critical    → Auto-land
```

Mindestens drei Methoden stehen jederzeit zur Verfügung (RC, PX4-Failsafe, lokaler Timeout).

---

## Überlegungen zur Zertifizierung

Diese Architektur unterstützt die Einhaltung regulatorischer Vorgaben:

| Anforderung | Umsetzung |
|-------------|----------------|
| **Verantwortlicher Pilot** | RC-Übersteuerung jederzeit |
| **Failsafe-Verhalten** | Zertifizierte PX4-Failsafes |
| **Geofencing** | Hardwareseitig durchgesetzte Grenzen |
| **Protokollierung** | JetStream-Audit-Trail |
| **Remote ID** | Aussendung über separates Modul |

**Hinweis:** Regulatorische Anforderungen variieren je nach Rechtsordnung. Diese Architektur bildet die technische Grundlage—die konkrete Einhaltung erfordert zusätzliche Maßnahmen gemäß den lokalen Vorschriften.

---

## Testen der Sicherheitssysteme

### Unit-Tests

```go
func TestGeofenceRejection(t *testing.T) {
    gw := NewGateway(config)
    gw.SetGeofence(testGeofence)

    cmd := &Command{
        Type:   "goto",
        Target: Position{Lat: 0, Lon: 0}, // Outside geofence
    }

    err := gw.validateCommand(cmd)
    assert.ErrorContains(t, err, "outside geofence")
}
```

### Integrationstests

```
1. Simulate RC loss → Verify RTL activates
2. Simulate battery drain → Verify land at threshold
3. Simulate geofence breach → Verify return behavior
4. Simulate Jetson crash → Verify continued flight
5. Simulate NATS disconnect → Verify local operation
```

### Flugtests

Vor dem Flotteneinsatz:

1. **Validierung der Failsafes an einem einzelnen Fahrzeug**
2. **Überprüfung der Reaktion auf RC-Verlust**
3. **Tests zur Geofence-Verletzung**
4. **Validierung der Akku-Failsafes**
5. **Tests der Koordination mehrerer Fahrzeuge**

---

## Zusammenfassung

| Schicht | Autorität | Drohnen | Bodenfahrzeuge |
|-------|-----------|--------|-----------------|
| **Manuell** | Höchste | RC-Sender | Lenkung / Not-Aus |
| **Fahrzeug** | Hoch | PX4-Failsafes | ECU / Sicherheitssteuerung |
| **Gateway** | Mittel | Durchsetzung von Richtlinien | Durchsetzung von Richtlinien |
| **Flotte** | Niedrigste | Nur Koordination | Nur Koordination |

Das Sicherheitsmodell gewährleistet:

- **Das Netzwerk ist niemals sicherheitskritisch**
- **Der Bediener hat stets Übersteuerungsautorität**
- **Hardware-Failsafes werden lokal ausgeführt**
- **Ausfälle führen zu kontrollierter Degradation**
- **Kein einzelner Fehlerpunkt hat flottenweite Auswirkungen**

---

## Verwandte Dokumentation

- [Unterstützte Plattformen]({{< relref "/fleet/platforms" >}}) — Übersicht aller Fahrzeugtypen
- [Drohnenplattform]({{< relref "/fleet/platforms/drones" >}}) — Drohnenspezifische Sicherheitsdetails
- [Bodenfahrzeuge]({{< relref "/fleet/platforms/ground" >}}) — Sicherheitsdetails für Bodenfahrzeuge

---

## Weiterführende Literatur

- [PX4-Sicherheitskonfiguration](https://docs.px4.io/main/en/config/safety.html)
- [MAVLink Common Message Set](https://mavlink.io/en/messages/common.html)
- [NATS-Sicherheitsmodell](https://docs.nats.io/running-a-nats-service/configuration/securing_nats)

---

## Zurück zur Übersicht

[← Architektur der autonomen Fahrzeugflotte]({{< relref "/fleet" >}})
