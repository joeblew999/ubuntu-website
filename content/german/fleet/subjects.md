---
title: "Subject-Benennung"
meta_title: "NATS-Subject-Benennung für Drohnenflotten | Ubuntu Software"
description: "Hierarchische NATS-Subject-Namenskonventionen für Flotten-Telemetrie, Ereignisse, Befehle und die Zustandsverwaltung digitaler Zwillinge."
image: "/images/robotics.svg"
draft: false
---

## Subject-Hierarchie-Design

Die Benennung von Subjects ist die API Ihres Messaging-Systems. Eine gut durchdachte Hierarchie ermöglicht effizientes Routing, Filtern und Zugriffskontrolle und bleibt dabei für Entwickler intuitiv.

---

## Designprinzipien

| Prinzip | Begründung |
|-----------|-----------|
| **Hierarchisch** | Wildcard-Abonnements auf jeder Ebene ermöglichen |
| **Vorhersehbar** | Entwickler können Subjects ohne Nachschlagen konstruieren |
| **Umgebungsbewusst** | Produktion von Staging/Entwicklung trennen |
| **Filterbar** | Ein Fahrzeug, eine Gruppe oder die gesamte Flotte abonnieren |
| **Erweiterbar** | Neue Nachrichtentypen ohne Umstrukturierung hinzufügen |

---

## Zentrales Subject-Muster

Alle Flotten-Subjects folgen diesem Muster:

```
fleet.<env>.veh.<vid>.<category>.<type>
```

| Segment | Beschreibung | Beispiele |
|---------|-------------|----------|
| `fleet` | Root-Namespace | Immer "fleet" |
| `<env>` | Umgebung | `prod`, `staging`, `dev` |
| `veh` | Ressourcentyp | Fahrzeuge (zukünftig: `gcs` für die Bodenkontrolle) |
| `<vid>` | Fahrzeug-ID | `VID-001`, `drone-alpha-7` |
| `<category>` | Nachrichtenkategorie | `state`, `evt`, `cmd`, `cmdack` |
| `<type>` | Spezifischer Typ | `position`, `battery`, `takeoff` |

---

## Nachrichtenkategorien

### Zustand: `fleet.<env>.veh.<vid>.state.*`

Aktueller Fahrzeugzustand, wird kontinuierlich veröffentlicht:

| Subject | Beschreibung | Rate |
|---------|-------------|------|
| `state.position` | Lat/Lon/Alt, Geschwindigkeit | 10 Hz |
| `state.attitude` | Roll/Pitch/Yaw | 10 Hz |
| `state.battery` | Spannung, Strom, Restkapazität | 1 Hz |
| `state.mode` | Flugmodus (AUTO, MANUAL, RTL) | Bei Änderung |
| `state.health` | Zusammenfassung des Systemzustands | 1 Hz |

**Beispiel für eine Veröffentlichung:**

```go
subject := fmt.Sprintf("fleet.prod.veh.%s.state.position", vehicleID)
nc.Publish(subject, positionJSON)
```

### Ereignisse: `fleet.<env>.veh.<vid>.evt.*`

Einzelne Vorkommnisse, werden veröffentlicht, sobald sie eintreten:

| Subject | Beschreibung |
|---------|-------------|
| `evt.armed` | Fahrzeug scharfgeschaltet |
| `evt.disarmed` | Fahrzeug entschärft |
| `evt.takeoff` | Start eingeleitet |
| `evt.landed` | Landung abgeschlossen |
| `evt.waypoint` | Wegpunkt erreicht |
| `evt.failsafe` | Failsafe ausgelöst |
| `evt.geofence` | Geofence-Verletzung |
| `evt.battery.low` | Warnung bei niedrigem Akkustand |

**Ereignisse sind Zustandsübergänge**, keine kontinuierlichen Daten. Sie sind dauerhaft – sie werden für Audit-Trails in JetStream gespeichert.

### Befehle: `fleet.<env>.veh.<vid>.cmd.*`

An Fahrzeuge gesendete Anweisungen:

| Subject | Beschreibung |
|---------|-------------|
| `cmd.arm` | Motoren scharfschalten |
| `cmd.disarm` | Motoren entschärfen |
| `cmd.takeoff` | Start einleiten |
| `cmd.land` | Landung einleiten |
| `cmd.rtl` | Rückkehr zum Startpunkt |
| `cmd.goto` | Zu Position navigieren |
| `cmd.mission.upload` | Mission hochladen |
| `cmd.mission.start` | Mission starten |
| `cmd.param.set` | Parameter setzen |

**Befehle fließen abwärts** vom Flottenmanagement zu den Fahrzeugen.

### Befehlsbestätigungen: `fleet.<env>.veh.<vid>.cmdack.*`

Antworten auf Befehle:

| Subject | Beschreibung |
|---------|-------------|
| `cmdack.arm` | Ergebnis des Scharfschaltbefehls |
| `cmdack.takeoff` | Ergebnis des Startbefehls |
| `cmdack.goto` | Ergebnis des Goto-Befehls |

**Der Ack-Payload enthält:**

```json
{
  "cmd_id": "cmd-123456",
  "status": "accepted|rejected|completed|failed",
  "error": "optional error message",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Wildcard-Abonnements

NATS-Wildcards ermöglichen flexible Abonnements:

| Wildcard | Bedeutung | Beispiel |
|----------|---------|---------|
| `*` | Einzelnes Token | `fleet.prod.veh.*.state.position` |
| `>` | Ein oder mehrere Tokens | `fleet.prod.veh.VID-001.>` |

### Gängige Abonnementmuster

| Muster | Anwendungsfall |
|---------|----------|
| `fleet.prod.veh.*.state.position` | Alle Fahrzeugpositionen |
| `fleet.prod.veh.VID-001.>` | Alles von einem Fahrzeug |
| `fleet.prod.veh.*.evt.>` | Alle Ereignisse von allen Fahrzeugen |
| `fleet.prod.veh.*.state.battery` | Alle Batteriezustände |
| `fleet.*.veh.VID-001.>` | Ein Fahrzeug über alle Umgebungen hinweg |

---

## KV-Stores: Digital-Twin-Shadow

JetStream-KV-Stores verwalten den **gewünschten** und den **gemeldeten** Zustand:

### Subject-Muster (KV)

```
fleet/<env>/veh/<vid>/desired
fleet/<env>/veh/<vid>/reported
```

Hinweis: KV verwendet das Trennzeichen `/` (Bucket-Pfad), während Pub/Sub das Trennzeichen `.` verwendet.

### Gewünschter Zustand

Was das Fahrzeug nach dem Willen des Flottenmanagementsystems tun **soll**:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "geofence_enabled": true,
  "max_altitude": 120,
  "home_position": {
    "lat": 37.7749,
    "lon": -122.4194
  }
}
```

### Gemeldeter Zustand

Was das Fahrzeug **tatsächlich** meldet:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "mission_progress": 0.45,
  "geofence_enabled": true,
  "armed": true,
  "in_flight": true,
  "battery_remaining": 0.72
}
```

### Shadow-Abgleich

Das Vehicle Gateway führt kontinuierlich folgende Schritte aus:

1. Liest den gewünschten Zustand aus dem KV
2. Vergleicht ihn mit dem tatsächlichen Fahrzeugzustand
3. Sendet Befehle, um die Realität an den gewünschten Zustand anzugleichen
4. Aktualisiert den gemeldeten Zustand

Dies ermöglicht **deklaratives Flottenmanagement** – den gewünschten Zustand festlegen und die Fahrzeuge darauf zusteuern lassen.

---

## Zugriffskontrolle

Die NATS-Autorisierung ordnet Benutzern die erlaubten Subjects zu:

```
# Vehicle can publish its own state/events
publish: fleet.prod.veh.VID-001.state.>
publish: fleet.prod.veh.VID-001.evt.>
publish: fleet.prod.veh.VID-001.cmdack.>

# Vehicle can subscribe to commands for itself
subscribe: fleet.prod.veh.VID-001.cmd.>

# Vehicle can read/write its own shadow
kv: fleet/prod/veh/VID-001/*
```

```
# Fleet operator can subscribe to all state/events
subscribe: fleet.prod.veh.*.state.>
subscribe: fleet.prod.veh.*.evt.>
subscribe: fleet.prod.veh.*.cmdack.>

# Fleet operator can publish commands
publish: fleet.prod.veh.*.cmd.>

# Fleet operator can manage shadows
kv: fleet/prod/veh/*/desired
```

Dies stellt sicher, dass **Fahrzeuge sich nicht als andere Fahrzeuge ausgeben** oder unautorisierte Befehle veröffentlichen können.

Eine ausführliche Darstellung des Autorisierungsmodells, einschließlich Berechtigungen für Dritte und Verwaltung von Anmeldedaten, finden Sie unter [Autorisierung & Berechtigungen]({{< relref "/fleet/authorization" >}}).

---

## Abwägungen beim Subject-Design

### Warum keine flachen Subjects?

Flache Subjects wie `vehicle-001-position` unterstützen Folgendes nicht:

- Wildcard-Abonnements
- Hierarchische Zugriffskontrolle
- Logische Gruppierung

### Warum keine tiefere Hierarchie?

Subjects wie `fleet.prod.region.us-west.site.warehouse-a.veh.VID-001.state.position`:

- Sind programmatisch schwieriger zu konstruieren
- Machen Wildcards komplex
- Führen dazu, dass sich Änderungen durch das gesamte System ziehen

Die vierstufige Hierarchie schafft ein Gleichgewicht zwischen Flexibilität und Einfachheit.

---

## Nachrichtenformat

Alle Nachrichten verwenden **JSON** mit einheitlicher Struktur:

```json
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00.123Z",
  "seq": 12345,
  "data": {
    // message-specific payload
  }
}
```

| Feld | Beschreibung |
|-------|-------------|
| `vid` | Fahrzeug-ID (redundant zum Subject, aber nützlich für Consumer) |
| `timestamp` | ISO 8601 mit Millisekundengenauigkeit |
| `seq` | Monotone Sequenznummer für die Reihenfolge |
| `data` | Nachrichtenspezifischer Payload |

---

## Zusammenfassung

| Kategorie | Subject-Muster | Richtung |
|----------|-----------------|-----------|
| **Zustand** | `fleet.<env>.veh.<vid>.state.<type>` | Fahrzeug → Hub |
| **Ereignisse** | `fleet.<env>.veh.<vid>.evt.<type>` | Fahrzeug → Hub |
| **Befehle** | `fleet.<env>.veh.<vid>.cmd.<type>` | Hub → Fahrzeug |
| **Acks** | `fleet.<env>.veh.<vid>.cmdack.<type>` | Fahrzeug → Hub |
| **Gewünscht** | `fleet/<env>/veh/<vid>/desired` (KV) | Hub → Fahrzeug |
| **Gemeldet** | `fleet/<env>/veh/<vid>/reported` (KV) | Fahrzeug → Hub |

---

## Weiter

[Stream-Konfiguration →]({{< relref "/fleet/streams" >}})
