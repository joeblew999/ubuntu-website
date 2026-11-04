---
title: "Autorisierung & Berechtigungen"
meta_title: "Dezentrale Autorisierung & Zugriffsberechtigungen | Ubuntu Software"
description: "Das dezentrale Sicherheitsmodell von NATS JetStream mit Grants ermöglicht es Dritten, sich sicher und ohne zentrale Koordination an Echtzeitsystemen zu beteiligen."
image: "/images/robotics.svg"
draft: false
---

## Sicherheit ohne Engpässe

Herkömmliche zentralisierte Autorisierung verursacht bei großem Maßstab Probleme:

| Problem | Auswirkung |
|---------|--------|
| **Single Point of Failure** | Auth-Server ausgefallen = gesamtes System ausgefallen |
| **Latenz bei jeder Anfrage** | Der Roundtrip zum Auth-Server verursacht Verzögerung |
| **Kein Betrieb ohne Verbindung** | Kein Netzwerk = keine Autorisierung |
| **Skaliert nicht** | Tausende Geräte überlasten den zentralen Server |

NATS verfolgt einen anderen Ansatz: **dezentrale Autorisierung** mit eigenständigen Zugangsdaten.

---

## So funktioniert die NATS-Autorisierung

```
┌─────────────────────────────────────────────────────────────────┐
│                    CREDENTIAL ISSUANCE                          │
│   Account Authority issues signed JWTs with embedded perms      │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌─────────┐     ┌─────────┐     ┌─────────┐
        │ Vehicle │     │ Operator│     │ Partner │
        │  Creds  │     │  Creds  │     │  Creds  │
        └────┬────┘     └────┬────┘     └────┬────┘
             │               │               │
             ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS SERVER                                │
│   Validates JWT signature locally — no external call needed     │
│   Enforces permissions from embedded claims                     │
└─────────────────────────────────────────────────────────────────┘
```

**Wichtige Eigenschaften:**

- **Lokale Validierung** — Der Server validiert Zugangsdaten, ohne einen externen Dienst aufzurufen
- **Eingebettete Berechtigungen** — Was Sie tun dürfen, ist in Ihren Zugangsdaten kodiert
- **Sofortiger Widerruf** — Die Widerrufsliste auf Account-Ebene wird im gesamten Cluster verteilt
- **Zero Trust** — Keine Zugangsdaten = kein Zugriff, Punkt

Das funktioniert auch bei getrennter Verbindung. Kein Call-Home für die Autorisierung. Die Zugangsdaten sind der Nachweis.

---

## Berechtigungsbereiche

Zugangsdaten legen genau fest, was eine Partei tun darf:

| Bereich | Beschreibung | Beispiel |
|-------|-------------|---------|
| **Publish** | Nachrichten an Subjects senden | Fahrzeug veröffentlicht `fleet.prod.veh.V001.state.pos` |
| **Subscribe** | Nachrichten von Subjects empfangen | Operator abonniert `fleet.prod.veh.*.state.*` |
| **Publish + Subscribe** | Beide Richtungen | Gateway verarbeitet Befehl/Antwort |
| **Deny** | Muster explizit blockieren | Fahrzeug daran hindern, andere Fahrzeuge zu abonnieren |

**Subject-Wildcards in Berechtigungen:**

```
fleet.prod.veh.V001.*        → Only this vehicle's subjects
fleet.prod.veh.*.state.*     → All vehicles' state (read-only)
fleet.prod.veh.*.cmd.*       → All vehicles' commands (operator)
fleet.prod.veh.*.evt.*       → All events (auditor)
```

So stellen wir sicher, dass Fahrzeuge nur ihre eigenen Daten veröffentlichen und sich nicht als andere Fahrzeuge ausgeben können. Die vollständige Hierarchie finden Sie unter [Subject-Benennung]({{< relref "/fleet/subjects" >}}).

---

## Das Grants-Modell

Grants ermöglichen **feingranulare, delegierbare Berechtigungen**. Ein Operator kann einem Dritten Zugriff auf bestimmte Daten gewähren, ohne ihm vollen Systemzugriff zu geben.

```
Account Authority (root)
    │
    ├── Fleet Operator Account
    │       │
    │       ├── Grant: Emergency Services (incident-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.pos
    │       │       └── Subscribe: fleet.prod.veh.*.stream.video
    │       │       └── Expires: 4 hours
    │       │
    │       ├── Grant: Research Partner (data-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.env.*
    │       │       └── Deny: fleet.prod.veh.*.state.pos
    │       │       └── Expires: 30 days
    │       │
    │       └── Grant: Auditor (read-only, full)
    │               └── Subscribe: fleet.prod.veh.*.state.*
    │               └── Subscribe: fleet.prod.veh.*.evt.*
    │               └── Deny: fleet.prod.veh.*.cmd.*
    │               └── Expires: 7 days
```

**Grant-Eigenschaften:**

| Eigenschaft | Vorteil |
|----------|---------|
| **Begrenzter Umfang** | Nur bestimmte Subjects/Muster—nicht mehr |
| **Zeitlich begrenzt** | Automatischer Ablauf, kein vergessener Zugriff |
| **Widerrufbar** | Sofortiger Widerruf, ohne den Empfänger einzubeziehen |
| **Delegierbar** | Operatoren können Unterberechtigungen an andere vergeben |
| **Auditiert** | Jede Grant-Ausstellung wird protokolliert |

---

## Szenarien für die Integration von Drittanbietern

### Integration von Rettungsdiensten

```
Scenario: Wildfire response, fire department needs drone video

Grant issued:
  - Subject: fleet.prod.veh.*.stream.video
  - Subject: fleet.prod.veh.*.state.pos
  - Duration: 4 hours (incident duration)
  - Audit: All access logged with responder ID

Result: Fire department connects with their credential, receives
real-time video and position from all drones in area.
Grant expires automatically when incident closes.
```

**Warum das funktioniert:**
- Die Feuerwehr erhält genau das, was sie braucht—Live-Video und Positionen
- Kein Zugriff auf Befehle—kann den Drohnenbetrieb nicht stören
- Zeitlich begrenzt—der Zugriff endet, wenn der Einsatz endet
- Auditierbar—jedes empfangene Frame wird protokolliert

### Datenaustausch mit Partnern

```
Scenario: University research on environmental monitoring

Grant issued:
  - Subject: fleet.prod.veh.*.state.env.* (temp, humidity, air quality)
  - Deny: fleet.prod.veh.*.state.pos (no location data)
  - Duration: 30 days (research period)
  - Audit: Data access logged for compliance

Result: Researcher receives environmental sensor data without
knowing vehicle locations—privacy preserved.
```

**Warum das funktioniert:**
- Forschende erhalten Umweltdaten für ihre Studie
- Standortdaten explizit verweigert—Privatsphäre geschützt
- Klarer Audit-Trail für die Einhaltung von Forschungsvorgaben
- Automatische Bereinigung nach Ende des Forschungszeitraums

### Einhaltung regulatorischer Vorgaben

```
Scenario: Aviation authority audit of flight operations

Grant issued:
  - Subject: fleet.prod.veh.*.state.* (all telemetry)
  - Subject: fleet.prod.veh.*.evt.* (all events)
  - Deny: fleet.prod.veh.*.cmd.* (no command capability)
  - Duration: 7 days (audit period)
  - Audit: Every query logged for audit trail

Result: Auditor gets complete read access to historical data
but cannot issue any commands to vehicles.
```

**Warum das funktioniert:**
- Volle Transparenz für den Auditor—nichts wird verborgen
- Nur Lesezugriff—der Auditor kann den Betrieb nicht beeinflussen
- Vollständiger Audit-Trail der Zugriffe des Auditors
- Zeitlich begrenzt—keine dauerhafte Zugriffslücke

### Kundenintegration

```
Scenario: Enterprise customer wants API access to their fleet data

Grant issued:
  - Subject: fleet.prod.veh.CUST-*.state.* (only their vehicles)
  - Subject: fleet.prod.veh.CUST-*.evt.*
  - Deny: fleet.prod.veh.OTHER-*.* (other customers' data)
  - Duration: Persistent (until revoked)
  - Audit: API access logged

Result: Customer builds their own dashboards using real-time
data from their vehicles—and only their vehicles.
```

### Zugriff für Subunternehmer

```
Scenario: Maintenance team needs to run diagnostics on specific vehicles

Grant issued:
  - Subject: fleet.prod.veh.V001.cmd.diag.* (diagnostic commands only)
  - Subject: fleet.prod.veh.V001.state.health.*
  - Deny: fleet.prod.veh.V001.cmd.flight.* (no flight commands)
  - Duration: 8 hours (maintenance window)
  - Audit: All commands logged with technician ID

Result: Technician can run diagnostics but cannot arm or fly the vehicle.
Access expires after maintenance window.
```

---

## Sicherheitseigenschaften

| Eigenschaft | Umsetzung |
|----------|----------------|
| **Zero Trust** | Kein implizites Vertrauen—jede Verbindung wird authentifiziert |
| **Least Privilege** | Zugangsdaten enthalten nur die minimal notwendigen Berechtigungen |
| **Defense in Depth** | Mehrere Schichten: TLS + JWT + Subject-ACLs |
| **Offline-fähig** | Autorisierung funktioniert ohne Netzwerkverbindung |
| **Auditierbar** | Jeder Zugriff wird mit der Identität der Zugangsdaten protokolliert |
| **Widerrufbar** | Sofortiger Widerruf über die Account-Widerrufsliste |

---

## Betriebliche Vorteile

### Für Flottenbetreiber

- **Dritte ohne Systemänderungen hinzufügen/entfernen** — Zugangsdaten ausstellen oder widerrufen, keine Konfigurationsdateien bearbeiten
- **Temporären Zugriff gewähren, ohne Zugangsdaten zu teilen** — Jede Partei erhält eigene Zugangsdaten
- **Nachvollziehen, wer wann auf was zugegriffen hat** — Vollständiges Zugriffsprotokoll
- **Zugriff bei Bedarf sofort widerrufen** — Wirksam bei der nächsten Verbindung

### Für Dritte

- **Self-Service-Verbindung** — Bereitgestellte Zugangsdaten nutzen und verbinden, sobald Sie bereit sind
- **Klarer Umfang** — Genau wissen, worauf Sie zugreifen können und worauf nicht
- **Keine Kenntnis der internen Architektur nötig** — Einfach verbinden und abonnieren
- **Automatischer Ablauf** — Keine Bereinigung auf Ihrer Seite erforderlich

### Für Compliance

- **Vollständiger Audit-Trail** — Jeder Zugriff ist auf Zugangsdaten zurückführbar
- **Nachweisbare Berechtigungsgrenzen** — Belegen, wer Zugriff auf was hatte
- **Zeitlich begrenzter Zugriff für Auditoren** — Keine dauerhaften Zugriffslücken
- **Keine geteilten Zugangsdaten** — Individuelle Verantwortlichkeit

---

## Hinweise zur Implementierung

### Format der Zugangsdaten

NATS verwendet **NKeys** und **JWTs**:

- **Ed25519-Signaturen** — Kryptografisch sicher, schnelle Verifizierung
- **Eigenständige Claims** — Berechtigungen sind im Token eingebettet
- **Keine externe Validierung** — Der Server verifiziert die Signatur lokal

### Widerruf

Wenn Sie Zugriff widerrufen müssen:

1. Zugangsdaten zur Account-Widerrufsliste hinzufügen
2. Die Liste wird an alle Server im Cluster verteilt
3. Wirksam beim nächsten Verbindungsversuch
4. Kein Eingriff in die Systeme des Empfängers nötig

### Rotation von Zugangsdaten

Für langlebige Integrationen:

1. Neue Zugangsdaten ausstellen, bevor die alten ablaufen
2. Ein Überlappungszeitraum ermöglicht einen nahtlosen Übergang
3. Alte Zugangsdaten nach abgeschlossenem Übergang widerrufen
4. Keine Dienstunterbrechung

---

## Warum das für reale Systeme wichtig ist

Herkömmliche Autorisierung erfordert:
- Einen zentralen Auth-Server (Single Point of Failure)
- Einen Netzwerkaufruf bei jeder Anfrage (Latenz)
- Eine komplexe Integration für Dritte (monatelange Arbeit)
- Eine manuelle Bereinigung abgelaufener Zugriffe (vergessene Zugangsdaten)

Die dezentrale NATS-Autorisierung bietet:
- **Resilienz** — Funktioniert offline, kein zentraler Engpass
- **Geschwindigkeit** — Lokale Validierung, kein Roundtrip
- **Erweiterbarkeit** — Dritte in Minuten statt Monaten hinzufügen
- **Sicherheit** — Automatischer Ablauf, sofortiger Widerruf

So baut man Systeme, die sowohl **hochsicher** als auch **äußerst benutzerfreundlich** sind.

---

## Verwandte Dokumentation

- [Subject-Benennung]({{< relref "/fleet/subjects" >}}) — Wie die Subject-Hierarchie feingranulare ACLs ermöglicht
- [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) — Wie die Messaging-Infrastruktur aufgebaut ist
- [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) — Wie sich die Befehlsautorisierung in die Sicherheitshierarchie einfügt
- [Überwachung]({{< relref "/company/surveillance" >}}) — Wie Grants Dritten Zugriff für das Monitoring ermöglichen

---

## Weiter

[Streams & Events →]({{< relref "/fleet/streams" >}})
