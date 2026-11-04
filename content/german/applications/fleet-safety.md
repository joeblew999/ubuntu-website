---
title: "Öffentliche Sicherheit & Notfälle"
meta_title: "IoT-Flotte für öffentliche Sicherheit | Ubuntu Software"
description: "Öffentliche Sicherheit mit IoT-Flotten: Drohnen zur Luftüberwachung, Einsatzfahrzeuge mit Echtzeit-Disposition und Gebäudesensoren zur Bedrohungserkennung. Einheitliche Architektur für eine koordinierte Reaktion."
image: "/images/robotics.svg"
draft: false
---

## Wenn Minuten zählen

Notfalleinsätze hängen von Informationen ab. IoT-Flotten liefern Lagebewusstsein, wenn Einsatzkräfte es am dringendsten brauchen – schnellere Bereitstellung, bessere Abdeckung, sicherere Einsätze.

---

## Flottentypen in der öffentlichen Sicherheit

| Flottentyp | Geräte | Rolle |
|------------|---------|------|
| **Luft** | Überwachungsdrohnen, Suchdrohnen | Suche und Rettung, Einsatzüberwachung aus der Luft, Brandkartierung |
| **Bodenfahrzeuge** | Einsatzfahrzeuge, Streifenwagen | Echtzeit-Disposition, Telematik, Asset-Tracking |
| **Stationäres IoT** | Gebäudealarme, Verkehrssensoren, Schussdetektoren | Bedrohungserkennung, Lagebewusstsein, Beweissicherung |

Alle Gerätetypen teilen dieselbe Architektur: Edge-KI, NATS-JetStream-Messaging und einheitliches Flottenmanagement.

---

## Die Herausforderung

In Notfallsituationen fehlen oft entscheidende Informationen:

- **Suche und Rettung** — Riesige Gebiete, begrenztes Personal
- **Schadensbewertung** — Beschädigte Infrastruktur, blockierte Zugänge
- **Laufende Einsätze** — Lagebild nötig, bevor Personal eintrifft
- **Koordination** — Mehrere Behörden, fragmentierte Kommunikation
- **Dokumentation** — Beweissicherung, Rekonstruktion von Vorfällen

Die Lagebeurteilung am Boden ist langsam und gefährdet die Einsatzkräfte. Hubschrauber sind teuer und nur begrenzt verfügbar. Die ersten Minuten entscheiden oft über den Ausgang.

---

## Wie IoT-Flotten das lösen

### Luft: Schnelle Bereitstellung

Vorpositionierte Drohnen starten sofort:

- **Automatische Alarmierung** — Drohnen sind unterwegs, während sich die Teams formieren
- **Als Erste vor Ort** — Visuelle Lagebeurteilung, bevor Personal eintrifft
- **Live-Streaming** — Die Einsatzleitung sieht, was die Drohne sieht
- **Nachtfähigkeit** — Wärmebildgebung bei völliger Dunkelheit

### Boden: Fahrzeugintelligenz

Einsatzfahrzeuge mit Echtzeit-Koordination:

- **Automatische Disposition** — Nächstgelegene verfügbare Einheit anhand von Live-Positionen
- **Routenoptimierung** — Echtzeit-Verkehrsdaten, schnellster Weg zum Einsatzort
- **Telematik** — Fahrzeugstatus, Kraftstoff, Wartungswarnungen
- **Ressourcenverfolgung** — Wer welche Ausrüstung hat, und wo

### Stationäres IoT: Bedrohungserkennung

Sensoren für kontinuierliches Lagebewusstsein:

- **Schusserkennung** — Akustische Sensoren bestimmen den Ort präzise
- **Verkehrssensoren** — Staus, Unfälle, Ampelvorrang
- **Gebäudealarme** — Brand, Einbruch, Aufzugsnotfall
- **Wetterstationen** — Lokale Bedingungen, die den Einsatz beeinflussen

### Flächenabdeckung

Systematische Suchmuster decken Gelände schnell ab:

- **Rastersuchen** — Automatisierte Flugrouten, keine Lücken
- **Koordination mehrerer Drohnen** — Größere Flächen gleichzeitig abdecken
- **Dauerüberwachung** — Stundenlange Abdeckung pro Akkuwechsel
- **Allwetterbetrieb** — Wind, Regen, Kälte – innerhalb der Grenzen

---

## Onboard-KI (Edge-Verarbeitung)

Echtzeit-Inferenz auf NVIDIA Jetson für sofortige Erkennung:

| Fähigkeit | Anwendung in der öffentlichen Sicherheit |
|------------|--------------------------|
| **Personenerkennung** | Personen in der Wildnis, in Trümmern und in Menschenmengen lokalisieren |
| **Fahrzeugerkennung** | Relevante Fahrzeuge verfolgen, verlassene Autos identifizieren |
| **Wärmebildgebung** | Wärmesignaturen durch Rauch, bei Nacht und im Blattwerk finden |
| **Branderkennung** | Offene Flammen, Glutnester und Ausbreitungsrichtung des Feuers erkennen |
| **Strukturbewertung** | Gebäudeschäden, Einsturzgefahr und Zugangswege erkennen |

**Warum Edge wichtig ist:**
- Einsatzkräfte sofort bei Erkennung alarmieren
- Arbeiten in Gebieten ohne Mobilfunkabdeckung
- Wärme- und Bilddaten gleichzeitig verarbeiten
- Weiterbetrieb bei Netzüberlastung

---

## Cloud-KI (Flottenanalyse)

Über NATS JetStream den Einsatz behördenübergreifend koordinieren:

| Fähigkeit | Anwendung in der öffentlichen Sicherheit |
|------------|--------------------------|
| **Behördenübergreifende Koordination** | Feeds zwischen Feuerwehr, Polizei, Rettungsdienst und SAR teilen |
| **Ressourcendisposition** | Optimale Routenführung auf Basis von Echtzeitbedingungen |
| **Musteranalyse** | Verlauf von Vorfällen verfolgen, Ausbreitung vorhersagen |
| **Beweismittelverwaltung** | Beweismittelkette, Zeitstempelverifizierung |
| **Einsatznachbesprechung** | Rekonstruktion von Vorfällen, Schulungsmaterial |

**Warum die Cloud wichtig ist:**
- Einheitliche Sicht für Leitstellen
- Koordination über Zuständigkeitsgrenzen hinweg
- Historische Daten zu ähnlichen Vorfällen
- Integration mit CAD-/Leitstellensystemen

---

## Hardware-Konfiguration

Für Einsätze in der öffentlichen Sicherheit empfehlen wir:

| Komponente | Auswahl | Begründung |
|-----------|-----------|-----------|
| **Flugwerk** | Holybro X500 oder robuste Variante | Zuverlässigkeit unter rauen Bedingungen |
| **Flugsteuerung** | Pixhawk 6X | Bewährte Autonomie, Fail-Safe-Modi |
| **Sensor-Companion** | Raspberry Pi CM4 | Telemetrie, einfache Sensoren |
| **KI-Companion** | Jetson Orin NX | Hohe Leistung für Wärme- und Bilddaten |
| **Nutzlast** | Duale RGB- + Wärmebildkamera | Tag- und Nachtfähigkeit |
| **Konnektivität** | LTE + Mesh-Funk | Ausfallsicherheit, wenn Funkmasten ausfallen |

[Vollständige Hardware-Spezifikationen →]({{< relref "/fleet/hardware" >}})

---

## Anwendungsfälle

### Suche und Rettung
Vermisstensuche. SAR in der Wildnis. Unterstützung bei der Wasserrettung.

### Brandbekämpfung
Lagebeurteilung bei Gebäudebränden. Kartierung von Waldbränden. Erkennung von Glutnestern.

### Strafverfolgung
Überwachung laufender Einsätze aus der Luft. Veranstaltungssicherheit. Management von Verkehrsunfällen.

### Katastrophenhilfe
Schadensbewertung. Ortung von Überlebenden. Unterstützung bei der Bereitstellung von Ressourcen.

### Grenz- und Seeüberwachung
Küstenpatrouille. Hafensicherheit. Grenzüberwachung.

---

## Einsatzmodelle

### Vorpositionierte Flotten
Drohnen stationiert an Feuerwachen, Polizeirevieren und SAR-Stützpunkten. Start bei Alarmierung.

### Mobile Einsatzleitung
Fahrzeuggestützte Systeme. Überall einsetzbar. Schneller Aufbau.

### Dauerhafte Abdeckung
Langzeitmissionen. Akkuwechsel- oder Ladestationen. Kontinuierlicher Betrieb.

---

## Technischer Tiefgang

Diese Anwendung basiert auf unserer produktionsreifen Flottenarchitektur:

| Komponente | Rolle in der öffentlichen Sicherheit |
|-----------|----------------------|
| [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) | Ausfallsicheres Messaging bei beschädigter Infrastruktur |
| [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) | Echtzeit-Erkennung mit sofortiger Alarmierung |
| [JetStream-Streams]({{< relref "/fleet/streams" >}}) | Beweissichere Videospeicherung mit Zeitstempeln |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Betrieb in der Nähe von Menschen, Gebäuden und Luftfahrzeugen |

[Vollständige Architektur erkunden →]({{< relref "/fleet" >}})

---

## Compliance

Einsätze in der öffentlichen Sicherheit erfordern:

- **FAA Part 107** – Unterstützung bei Ausnahmegenehmigungen für Nacht- und BVLOS-Einsätze
- **CJIS-Compliance** für den Umgang mit Daten der Strafverfolgung
- **Lückenlose Beweismittelkette** mit kryptografischer Verifizierung
- **Behördenübergreifender Datenaustausch** mit Zugriffskontrollen

Unsere Architektur unterstützt diese Anforderungen von Grund auf.

---

## Jetzt starten

Wenn Sekunden zählen, kommt es auf die Infrastruktur an. Setzen Sie Drohnenflotten ein, die schneller reagieren und mehr sehen.

[Kontaktieren Sie uns →](/contact)
