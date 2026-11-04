---
title: "DJI-Migrationspfad"
meta_title: "Migration von DJI zu Open-Source-Drohnen | Ubuntu Software"
description: "Umstieg von DJI auf Open-Source-Drohnen vor dem US-Verbot. Vollständiger Migrationsleitfaden für kommerzielle Betreiber, Behörden und Unternehmensflotten."
image: "/images/robotics.svg"
draft: false
---

{{< notice "danger" >}}
**Das Verbot ist jetzt in Kraft.** Am 22. Dezember 2025 hat die FCC DJI in die Covered List aufgenommen. Neue DJI-Modelle und HF-Komponenten können keine FCC-Zulassung für den Verkauf in den USA mehr erhalten. Bestehende Drohnen dürfen weiterhin legal geflogen werden, doch die Uhr für den Umstieg tickt. [Zur vollständigen Analyse →]({{< relref "/blog/dji-fcc-ban-what-happened" >}})
{{< /notice >}}

## Warum das wichtig ist

Wenn Sie DJI-Drohnen für kommerzielle Einsätze fliegen (Inspektion, Kartierung, Vermessung, öffentliche Sicherheit), brauchen Sie einen Umstiegsplan. Nicht weil Ihre bestehenden Drohnen den Betrieb einstellen, sondern aus folgenden Gründen:

- **Keine Firmware-Updates** nach dem Verbot
- **Keine Ersatzteile** über autorisierte Kanäle
- **Kein SDK-Support** für individuelle Integrationen
- **Keine Cloud-Dienste** für das Flottenmanagement
- **Keine Neuanschaffungen** zur Erweiterung Ihrer Flotte

Betreiber, die abwarten, müssen später einen Notfall-Umstieg zu Premiumpreisen bewältigen. Wer jetzt beginnt, kann den Umstieg planvoll vollziehen.

---

## Die Lage

### Was passiert ist

Am 22. Dezember 2025, einen Tag vor Ablauf der Frist, gab die FCC die Aufnahme von DJI in die Covered List bekannt. Die Regierung umging die in Section 1709 des NDAA vorgeschriebene Sicherheitsprüfung und erließ stattdessen eine „Feststellung zur nationalen Sicherheit“.

DJI schickte im März, Juni und Dezember 2025 Schreiben, in denen die Behörden gebeten wurden, die DJI-Produkte zu untersuchen. Keine Behörde antwortete. Es wurde keine Untersuchung eingeleitet. Es wurden keine Beweise gesammelt.

Der Kongress hat Section 1709 bewusst mit Falltüren versehen: Keiner Behörde wurde die Durchführung der Prüfung übertragen, und das Verbot sollte automatisch in Kraft treten, falls keine Prüfung stattfand. Das Ergebnis stand von vornherein fest.

### Was es bedeutet

| Auswirkung | Ihre bestehende DJI-Flotte | Neue DJI-Anschaffungen |
|--------|------------------------|-------------------|
| **Fliegen** | Weiterhin legal nach FAA Part 107 | Blockiert |
| **Firmware** | Keine Updates | Entfällt |
| **Ersatzteile** | Eingeschränkt/Graumarkt | Entfällt |
| **SDK/API** | Nicht unterstützt | Entfällt |
| **Cloud-Dienste** | Ungewiss | Entfällt |
| **Neue Modelle** | Entfällt | Blockiert |

### Wer betroffen ist

- **Regierungsbehörden** — Bereits von DJI-Beschaffungen ausgeschlossen
- **Öffentliche Sicherheit** — Über 80 % der US-Feuerwehren, Polizeibehörden und Such- und Rettungsdienste nutzen DJI
- **Kommerzielle Betreiber** — Inspektion, Kartierung, Vermessung, Landwirtschaft
- **Unternehmensflotten** — Infrastruktur, Versorgungsunternehmen, Bauwesen
- **Rüstungsunternehmen** — Blue-UAS-Anforderungen

---

## Der Migrationspfad

### Was Sie gewinnen

Beim Umstieg auf Open Source geht es nicht nur darum, dem Verbot auszuweichen. Sie gewinnen:

| DJI | Open Source (PX4/ArduPilot) |
|-----|----------------------------|
| Geschlossenes Ökosystem | Herstellerunabhängig |
| Proprietäre Daten | Ihre Daten, Ihre Server |
| Funktions-Roadmap unter fremder Kontrolle | Community- und Ihre eigene Roadmap |
| SDK-Einschränkungen | Vollständiger Quellcodezugriff |
| Hardware-Lock-in | Mehrere Hersteller |
| Cloud-Abhängigkeit | Selbst gehostet oder verwaltet |

### Kompatible Plattformen

Open-Source-Flight-Stacks funktionieren mit einer großen Bandbreite an Hardware:

| Plattform | Flight Stack | Stärken |
|----------|-------------|-----------|
| **Holybro X500/S500** | PX4/ArduPilot | Bewährt, verfügbar, gut dokumentiert |
| **Eigenbauten** | PX4/ArduPilot | Maximale Flexibilität |
| **Freefly** | PX4 | Professionell für Film und Vermessung |
| **Autel** | Proprietär | In den USA montiert, steht aber ebenfalls vor einer NDAA-Prüfung |
| **Skydio** | Proprietär | In den USA hergestellt, Fokus auf Autonomie |

Für den Flottenbetrieb empfehlen wir PX4-basierte Plattformen aus folgenden Gründen:
- Open Source mit kommerziellem Rückhalt (Dronecode Foundation)
- Das MAVLink-Protokoll ist ein offener Standard
- Großes Ökosystem kompatibler Komponenten
- Die BSD-Lizenz erlaubt proprietäre Modifikationen

### Optionen für die Bodenkontrolle

| Software | Protokoll | Stärken |
|----------|----------|-----------|
| **QGroundControl** | MAVLink | Kostenlos, plattformübergreifend, PX4/ArduPilot |
| **Mission Planner** | MAVLink | Auf ArduPilot ausgerichtet, Windows |
| **UgCS** | MAVLink | Kommerziell, erweiterte Planung |
| **Unsere Flottenplattform** | NATS + MAVLink | Mehrere Fluggeräte, Enterprise |

---

## Umstiegsstufen

### Stufe 1: Parallele Flotte (6-12 Monate)

Ideal für: Behörden mit laufendem Betrieb, die sich keine Ausfallzeiten leisten können.

1. **Testfluggeräte beschaffen** — Beginnen Sie mit 1-3 Open-Source-Plattformen
2. **Piloten schulen** — Andere Fernsteuerung, andere Software
3. **Arbeitsabläufe validieren** — Kartierung, Inspektion, Ihre Anwendungsfälle
4. **SOPs erstellen** — Neue Checklisten, Wartungsverfahren
5. **Schrittweiser Umstieg** — Verlagern Sie den Betrieb, während das Vertrauen wächst

### Stufe 2: Geplante Migration (3-6 Monate)

Ideal für: Betreiber, die Umstiegsfenster einplanen können.

1. **Bestandsaufnahme** — Welche DJI-Geräte, welche Arbeitsabläufe
2. **Plattformauswahl** — Fähigkeiten auf die Einsatzanforderungen abstimmen
3. **Beschaffung** — Hardware bestellen, Lieferzeiten einplanen
4. **Intensivschulung** — Komprimierte Pilotenausbildung
5. **Umstellung** — DJI-Geräte systematisch ausmustern

### Stufe 3: Notfall-Migration (1-3 Monate)

Ideal für: Betreiber, die zu lange gewartet haben.

1. **Sofortige Beschaffung** — Premiumpreise, begrenzte Auswahl
2. **Crashkurs** — Intensiv, möglicherweise durch externe Dienstleister
3. **Eingeschränkte Leistungsfähigkeit** — Vorübergehende Lücken in den Arbeitsabläufen in Kauf nehmen
4. **Iterieren** — Den Betrieb mit der Zeit verbessern

**Landen Sie nicht in Stufe 3.**

---

## Unsere Rolle

### Was wir bieten

**Flotteninfrastruktur, keine Hardware.** Wir konkurrieren nicht mit Drohnenhändlern. Wir liefern die Softwareschicht, die Flotten im Alltag funktionsfähig macht.

| Sie erhalten | Wir übernehmen |
|---------|-----------|
| Verwaltete NATS-Infrastruktur | Cluster-Betrieb, Updates |
| Vehicle-Gateway-Software | MAVLink-zu-NATS-Bridge |
| Verwaltung des Flottenstatus | Digitale Zwillinge, Telemetrie |
| Koordination mehrerer Fluggeräte | Subject-Hierarchie, Autorisierung |
| Enterprise-Integration | APIs, Datenexport |

### Was Sie beitragen

| Sie übernehmen | Wir unterstützen |
|------------|-----------|
| Hardwarebeschaffung | Empfehlungen von Partnerhändlern |
| Pilotenschulung | Dokumentation, Best Practices |
| Betrieb vor Ort | Remote-Flottenüberwachung |
| Einhaltung gesetzlicher Vorschriften | Datenaufbewahrung, Audit-Trails |

### Speziell für die Kartierung

Viele DJI-Betreiber benötigen Kartierungslösungen. Unsere Architektur unterstützt:

- **Pipelines für Sensorrohdaten** — Bilder zum Verarbeitungs-Backend
- **Echtzeit-Telemetrie** — Position, Höhe, Kameraereignisse
- **Flottenkoordination** — Vermessungsmuster mit mehreren Fluggeräten
- **Datenübergabe** — Integration mit ODM, Pix4D, DroneDeploy

Wir sind keine vollständige Kartierungslösung. Wir sind die Flotteninfrastruktur, die Ihre Drohnen mit Ihrer Verarbeitungspipeline verbindet.

---

## Hardware-Partner

Diese Händler helfen Ihnen bei der Beschaffung von Open-Source-Hardware:

{{< partners-list status="active" type="manufacturer,reseller" >}}

Suchen Sie einen Händler vor Ort? [Finden Sie einen Partner in Ihrer Nähe]({{< relref "/fleet/hardware" >}}).

---

## Erste Schritte

### Beratungsgespräch

Ein kostenloses 30-minütiges Gespräch, um Ihren aktuellen Betrieb zu verstehen:

- Flottengröße und -zusammensetzung
- Hauptanwendungsfälle (Kartierung, Inspektion usw.)
- Zeitplan und Dringlichkeit
- Budgetvorgaben

[Beratung vereinbaren →](/contact?subject=dji-migration)

### Pilotprogramm

Für Behörden und Unternehmen mit großen DJI-Flotten:

- Proof of Concept mit Ihren Arbeitsabläufen
- Parallelbetrieb neben der bestehenden DJI-Flotte
- Leistungsvergleich und Validierung
- Go/No-Go-Entscheidung auf Basis realer Daten

[Pilotprogramm anfragen →](/contact?subject=dji-pilot)

---

## Ressourcen

- [Hardware-Stack →]({{< relref "/fleet/hardware" >}}) — Referenz-Stückliste für PX4-Drohnen
- [Software-Stack →]({{< relref "/fleet/software" >}}) — Betriebssysteme und Flight Stacks
- [Flottenarchitektur →]({{< relref "/fleet" >}}) — NATS-basiertes Flottenmanagement
- [Partnerprogramm →]({{< relref "/partners" >}}) — Für Drohnenhändler, die DJI-Umsteiger betreuen

---

## FAQ

**Werden meine bestehenden DJI-Drohnen aufhören zu funktionieren?**

Nein. Ihre aktuellen Drohnen dürfen weiterhin nach FAA Part 107 fliegen. Sie erhalten jedoch keine Firmware-Updates mehr, Ersatzteile könnten knapp werden und Cloud-Dienste könnten eingestellt werden.

**Ist Autel eine sichere Alternative?**

Autel wird in der Prüfanforderung von NDAA Section 1709 ebenfalls genannt. Für Autel gilt derselbe Zeitplan wie für DJI.

**Ist Skydio eine sichere Alternative?**

Skydio stellt in den USA her und unterliegt nicht denselben Beschränkungen. Allerdings ist auch das Ökosystem von Skydio proprietär.

**Warum Open Source statt einer anderen proprietären Plattform?**

Open Source gibt Ihnen die Kontrolle. Sie tauschen nicht eine Herstellerabhängigkeit gegen eine andere. Sie können die Software nach Belieben anpassen, die Daten gehören Ihnen und die Hardware ist austauschbar.

**Wie lange dauert der Umstieg?**

Das hängt von Ihrem Betrieb ab. Planen Sie 3-6 Monate für einen methodischen Umstieg mit minimalen Unterbrechungen ein. Notfall-Umstiege sind innerhalb weniger Wochen möglich, jedoch mit höheren Kosten und Risiken.

**Was ist mit unseren bestehenden Arbeitsabläufen für die Kartierung?**

Die meisten Verarbeitungsprogramme (Pix4D, DroneDeploy, ODM) funktionieren mit jeder Drohne, die georeferenzierte Bilder liefert. Der Umstieg betrifft die Aufnahmeplattform, nicht das Verarbeitungs-Backend.

---

## Die Zeit ist abgelaufen

Das Verbot ist jetzt in Kraft. Wenn Ihre Flotte noch ausschließlich aus DJI-Drohnen besteht, arbeiten Sie auf geliehener Zeit. Sie haben keine Möglichkeit zur Flottenerweiterung, der Firmware-Support ist ungewiss und die Lieferketten werden immer knapper.

Die Händler, die sich mit Open Source auskennen, werden bald alle Hände voll zu tun haben.

[Kontakt aufnehmen →](/contact?subject=dji-migration)
