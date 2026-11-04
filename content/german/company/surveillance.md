---
title: "Überwachung"
meta_title: "Echtzeitüberwachung & Lagebewusstsein | Ubuntu Software"
description: "Verwandeln Sie die Datenflut von Drohnen, Satelliten und IoT-Sensoren in handlungsrelevante Warnmeldungen. Selbstbestimmt, sicher, mit lückenlosem Herkunftsnachweis."
image: "/images/robotics.svg"
draft: false
---

## In Daten ertrinken

Moderne Einsätze erzeugen mehr Daten, als Menschen verarbeiten können. Drohnen streamen Video. Satelliten erfassen Bilder. Sensoren melden kontinuierlich. Das Signal ist da – begraben im Rauschen.

Traditionelle Ansätze scheitern bei großem Umfang:
- **Manuelle Überwachung** skaliert nicht über einige wenige Feeds hinaus
- **Einfache Schwellenwerte** führen zu Alarmmüdigkeit
- **Isolierte Systeme** übersehen quellenübergreifende Muster
- **Reine Cloud-Lösungen** können sensible Einsätze nicht abdecken

Sie brauchen ein System, das alles beobachtet, Zusammenhänge versteht und Sie nur dann alarmiert, wenn es darauf ankommt.

---

## Von der Datenflut zur Erkenntnis

Unsere Überwachungsarchitektur verarbeitet Daten aus mehreren Quellen in Echtzeit, setzt KI sowohl am Edge als auch in der Cloud ein und filtert anschließend anhand betreiberdefinierter Regeln, um handlungsrelevante Ereignisse sichtbar zu machen.

```
┌─────────────────────────────────────────────────────────────────┐
│                         DATA SOURCES                            │
├─────────────────┬─────────────────┬─────────────────────────────┤
│     Drones      │   Satellites    │        IoT Sensors          │
│  (telemetry,    │  (imagery,      │  (environmental,            │
│   video, AI)    │   positioning)  │   industrial)               │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    FIREHOSE INGESTION                           │
│            NATS JetStream — persistent, ordered streams         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AI PROCESSING                              │
│   Edge (Jetson): Real-time inference, event detection           │
│   Cloud: Pattern recognition, correlation, anomaly detection    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    OPERATOR RULES ENGINE                        │
│   "Alert me when..." — thresholds, patterns, conditions         │
│   Customizable per user, role, or organization                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    ALERT + CONTEXT                              │
│   • What happened (event description)                           │
│   • Where (2D map, 3D spatial view)                             │
│   • Why it matters (AI reasoning)                               │
│   • Provenance (full data lineage, evidence chain)              │
└─────────────────────────────────────────────────────────────────┘
```

---

## Datenquellen

Wir behandeln alle Datenquellen gleichrangig – einheitliche Erfassung, konsistente Verarbeitung, korrelierte Analyse.

| Quelle | Datentypen | Beispiele |
|--------|------------|----------|
| **Drohnen** | Telemetrie, Video, Sensormesswerte, KI-Inferenzergebnisse | Flottenposition, Wärmebilder, Objekterkennung |
| **Satelliten** | Bilddaten, Positionsbestimmung, Wetterdaten | Änderungserkennung, Abdeckungsanalyse, atmosphärische Bedingungen |
| **IoT-Sensoren** | Umwelt, Industrie, Infrastruktur | Temperatur, Vibration, Wasserstände, Luftqualität |

Dieselbe Architektur, die Drohnentelemetrie verarbeitet, verarbeitet auch Satellitenbilder und Bodensensoren. Ein System. Das vollständige Bild.

---

## Sie legen fest, was zählt

Betreiber konfigurieren, was für sie relevant ist. Das System beobachtet alles und alarmiert nur bei den von Ihnen festgelegten Bedingungen.

| Der Betreiber sagt | Das System tut |
|---------------|-------------|
| „Alarmiere mich, wenn eine Drohne das GPS-Signal verliert“ | Überwacht die gesamte Flottentelemetrie, benachrichtigt bei GPS-Verlust |
| „Markiere thermische Anomalien an Stromleitungen“ | Führt Wärmeanalysen durch, alarmiert bei Hotspots |
| „Verfolge Fahrzeuge, die in Sperrzonen einfahren“ | Geofence-Überwachung mit sofortigen Alarmen |
| „Benachrichtige mich bei Pflanzenstress über 20 %“ | Aggregiert multispektrale Daten, löst bei Überschreiten des Schwellenwerts aus |
| „Korreliere Satellitenveränderungen mit Bodensensoren“ | Gleicht mehrere Quellen ab, erkennt Muster |

Regeln können einfache Schwellenwerte oder komplexe Korrelationen über mehrere Quellen sein. Sie steuern die Empfindlichkeit.

---

## Visualisierung mit Herkunftsnachweis

*Geplante Funktion (Roadmap)*

Wenn ein Alarm ausgelöst wird, sehen Betreiber das vollständige Bild:

- **2D-Kartenansicht** — Standortkontext mit historischem Bewegungsverlauf
- **3D-Raumansicht** — Vollständiger Umgebungskontext, Gelände, Bauwerke
- **Datenherkunft** — Lückenlose Kette von den Sensor-Rohdaten bis zum Alarm
- **Beweispaket** — Originaldaten, verarbeitete Ergebnisse, KI-Konfidenzwerte

**Jede Entscheidung ist nachprüfbar. Jeder Alarm ist rückverfolgbar.**

Betreiber können vom Alarm bis zu den Rohdaten vordringen und genau nachvollziehen, warum das System ein Ereignis markiert hat. Keine Blackboxes.

---

## Sicherheit & Souveränität

Diese Systeme unterstützen sensible Einsätze. Sicherheit ist nicht optional – sie ist grundlegend.

| Prinzip | Umsetzung |
|-----------|----------------|
| **Selbstsouverän** | Sie besitzen und kontrollieren alle Daten – wir greifen niemals ohne ausdrückliche Genehmigung darauf zu |
| **Ende-zu-Ende-Verschlüsselung** | Daten werden bei der Übertragung und im Ruhezustand verschlüsselt, mit Schlüsseln unter Ihrer Kontrolle |
| **Zero Trust** | Jede Anfrage authentifiziert, jede Aktion autorisiert, kein implizites Vertrauen |
| **Air-Gap-fähig** | Voller Funktionsumfang ohne Internetverbindung für sensible Bereitstellungen |
| **Lückenloser Audit-Trail** | Jeder Zugriff, jede Abfrage, jeder Alarm wird mit Zeitstempel protokolliert |
| **Datenresidenz** | Wählen Sie, wo Ihre Daten liegen – Ihre Infrastruktur, Ihre Rechtsordnung |

**Ihre Daten. Ihre Infrastruktur. Ihre Kontrolle.**

Selbst gehostete Bereitstellungen laufen vollständig auf Ihrer Hardware. Verwaltete Bereitstellungen isolieren Ihre Daten mit dedizierten Verschlüsselungsschlüsseln. So oder so behalten Sie die Souveränität.

---

## Anwendungsfälle

| Bereich | Anwendung |
|--------|-------------|
| **Grenzsicherung** | Perimeterüberwachung mit mehreren Sensoren, Eindringlingserkennung, Einsatzkoordination |
| **Kritische Infrastruktur** | Überwachung von Pipelines und Stromleitungen, vorausschauende Wartung, Anomalieerkennung |
| **Umwelt** | Wildtierverfolgung, Schadstoffüberwachung, Waldbranderkennung, Zustand von Ökosystemen |
| **Landwirtschaft** | Überwachung der Pflanzengesundheit, Bewässerungsoptimierung, Frühwarnung vor Schädlingen/Krankheiten |
| **Maritim** | Schiffsverfolgung, Erkennung illegaler Fischerei, Hafensicherheit, Suche und Rettung |
| **Urbane Einsätze** | Verkehrsflussanalyse, Überwachung von Menschenmengen, Koordination von Notfalleinsätzen |

Jeder Anwendungsfall nutzt dieselbe Kernarchitektur – andere Sensoren, andere Regeln, dieselbe zuverlässige Plattform.

---

## Technisches Fundament

Diese Überwachungsfunktion basiert auf unserer produktionsreifen Flottenarchitektur:

- [Flottenarchitektur]({{< relref "/fleet" >}}) — Wie Daten vom Edge in die Cloud fließen
- [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) — Das Messaging-Rückgrat, das die Datenflut bewältigt
- [Streams & Ereignisse]({{< relref "/fleet/streams" >}}) — Wie wir kontinuierliche Datenströme in abfragbaren Zustand umwandeln
- [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) — Edge-Verarbeitung und Ereigniserkennung
- [Autorisierung & Berechtigungen]({{< relref "/fleet/authorization" >}}) — Wie Dritte eingeschränkten, sicheren Zugriff erhalten

Dieselbe Infrastruktur, die Drohnenflotten verwaltet, treibt auch Überwachungseinsätze an. Im großen Maßstab praxiserprobt.

---

## Erste Schritte

Bereit, Ihre Datenflut in verwertbare Erkenntnisse zu verwandeln?

[Kontaktieren Sie uns →](/contact), um Ihre Überwachungsanforderungen zu besprechen.

Wir helfen Ihnen zu klären:
- Welche Datenquellen integriert werden sollten
- Welche Betreiberregeln für Ihren Bereich sinnvoll sind
- Selbst gehostete vs. verwaltete Bereitstellungsoptionen
- Sicherheits- und Compliance-Aspekte
