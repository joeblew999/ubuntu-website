---
title: "Energie & Versorgung"
meta_title: "IoT-Flotte für Energie & Versorgung | Ubuntu Software"
description: "Energieinfrastruktur mit IoT-Flotten: Drohnen für die Leitungsinspektion, Servicefahrzeuge mit Telematik, intelligente Stromzähler und Sensoren in Umspannwerken. Einheitliche Architektur für netzweite Intelligenz."
image: "/images/robotics.svg"
draft: false
---

## Inspizieren, was am wichtigsten ist

Energieinfrastruktur erstreckt sich über Tausende von Kilometern. Herkömmliche Inspektionen sind langsam und gefährlich, und sie übersehen Probleme. IoT-Flotten verändern alles.

---

## Flottentypen in der Energiebranche

| Flottentyp | Geräte | Rolle |
|------------|---------|------|
| **Luft** | Inspektionsdrohnen, Vermessungsdrohnen | Leitungsinspektion, Thermalanalyse, Solarpark-Vermessungen |
| **Bodenfahrzeuge** | Servicefahrzeuge, Inspektionsfahrzeuge | Einsatzsteuerung von Teams, Ausrüstungstransport, Wartung vor Ort |
| **Stationäres IoT** | Intelligente Zähler, Umspannwerk-Monitore, Pipeline-Sensoren | Netzüberwachung in Echtzeit, Fehlererkennung, Verbrauchsanalysen |

Alle Gerätetypen nutzen dieselbe Architektur: Edge-KI, NATS-JetStream-Messaging und einheitliches Flottenmanagement.

---

## Die Herausforderung

Energieversorger und Versorgungsunternehmen stehen vor besonderen Herausforderungen bei der Inspektion:

- **Umfang**: Tausende Kilometer an Übertragungsleitungen, Pipelines und Verteilnetzen
- **Zugang**: Abgelegene Standorte, schwieriges Gelände, gefährliche Umgebungen
- **Risiko**: Sicherheit der Mitarbeiter bei Hochspannung, in großen Höhen und in engen Räumen
- **Häufigkeit**: Vorschriften verlangen regelmäßige Inspektionen alternder Infrastruktur
- **Präzision**: Kleine Defekte von heute werden zu katastrophalen Ausfällen von morgen

Manuelle Inspektionsteams können nicht genug Fläche abdecken. Hubschrauber-Befliegungen sind teuer und liefern nur begrenzte Daten. Satelliten fehlt die Auflösung für die Defekterkennung.

---

## Wie IoT-Flotten das lösen

### Luft: Systematische Abdeckung

Drohnen decken mehr Fläche bei gleichbleibender Qualität ab:

- **Automatisierte Flugrouten**: Folgen Übertragungsleitungen, Pipeline-Trassen und Begrenzungen
- **Einheitliche Inspektion**: Jedes Mal dieselben Winkel und dieselbe Auflösung
- **Paralleler Betrieb**: Mehrere Drohnen decken gleichzeitig verschiedene Abschnitte ab
- **Allwettertauglich**: Einsatz unter Bedingungen, bei denen Bodenteams nicht arbeiten können

### Boden: Intelligenz für die Serviceflotte

Fahrzeuge und Teams mit Koordination in Echtzeit:

- **Optimierte Einsatzplanung**: Techniker nach Qualifikation, Standort und Priorität einsetzen
- **Telematik**: Fahrzeugstatus, Kraftstoff, Einhaltung von Sicherheitsvorgaben
- **Integration von Arbeitsaufträgen**: Digitale Arbeitsaufträge, Nachverfolgung der Erledigung
- **Ersatzteilbestand**: Was sich auf jedem Fahrzeug befindet und was wo benötigt wird

### Stationäres IoT: Netztransparenz

Sensoren im gesamten Netz:

- **Intelligente Zähler**: Verbrauchsmuster, Ausfallerkennung, Lastprognosen
- **Umspannwerk-Monitore**: Transformatorzustand, Fehlererkennung, SCADA-Integration
- **Pipeline-Sensoren**: Druck, Durchfluss, Leckerkennung
- **Wetterstationen**: Lokale Bedingungen, die den Betrieb beeinflussen

### Detaillierte Analyse

Alle Flottentypen sehen, was Menschen übersehen:

- **Wärmebildgebung**: Hotspots erkennen, die auf elektrische Fehler hinweisen
- **Hochauflösende Bilder**: Korrosion, Schäden und überwuchernde Vegetation erkennen
- **LiDAR-Kartierung**: Abstände und Bodenbewegungen messen
- **Multispektral**: Zustand der Vegetation in der Nähe der Infrastruktur bewerten

---

## KI an Bord (Edge-Verarbeitung)

Inferenz in Echtzeit auf NVIDIA Jetson für sofortige Erkennung:

| Fähigkeit | Anwendung in der Energiebranche |
|------------|-------------------|
| **Erkennung thermischer Anomalien** | Überhitzte Verbindungen, Transformatoren und Isolatoren identifizieren |
| **Komponentenerkennung** | Masten, Türme, Isolatoren und Leiter automatisch klassifizieren |
| **Defekterkennung** | Korrosion, Risse, Vogelnester und beschädigte Komponenten erkennen |
| **Vegetationsanalyse** | Überwuchs messen, Bäume mit hohem Risiko identifizieren |
| **Erkennung von Sicherheitsgefahren** | Unbefugten Zutritt und Schäden an Anlagen erkennen |

**Warum Edge wichtig ist:**
- Sofortige Warnungen bei kritischen Defekten
- Weniger Datenübertragung über abgelegene Netze
- Thermische Daten in Echtzeit verarbeiten
- Inspektion auch bei schlechter Verbindung fortsetzen

---

## Cloud-KI (Flottenanalysen)

Über NATS JetStream die Wartung im gesamten Netz optimieren:

| Fähigkeit | Anwendung in der Energiebranche |
|------------|-------------------|
| **Vorausschauende Wartung** | Komponentenausfälle anhand thermischer Trends vorhersagen |
| **Bewertung des Anlagenzustands** | Infrastruktur nach Zustand und Risiko einstufen |
| **Änderungserkennung** | Inspektionen über die Zeit vergleichen, Verschleiß erkennen |
| **Compliance-Berichte** | Regulatorische Dokumentation automatisch erstellen |
| **Ausfallkorrelation** | Inspektionsbefunde mit Versorgungsunterbrechungen verknüpfen |

**Warum die Cloud wichtig ist:**
- Befunde über die gesamte Infrastruktur zusammenführen
- Modelle mit historischen Defektdaten trainieren
- Integration mit GIS- und Asset-Management-Systemen
- Anforderungen regulatorischer Audits unterstützen

---

## Hardwarekonfiguration

Für Einsätze bei Energie- und Versorgungsunternehmen empfehlen wir:

| Komponente | Auswahl | Begründung |
|-----------|-----------|-----------|
| **Flugwerk** | Industrieller Multikopter oder Starrflügler | Stabilität für Nahinspektionen, Ausdauer für lange Strecken |
| **Flugsteuerung** | Pixhawk 6X | Zuverlässige Autonomie in EMV-belasteten Umgebungen |
| **Sensor-Begleitrechner** | Raspberry Pi CM4 | Koordination mehrerer Sensoren |
| **KI-Begleitrechner** | Jetson Orin Nano | Thermische und visuelle KI-Verarbeitung |
| **Wärmebildkamera** | Radiometrische Wärmebildkamera | Temperaturmessung, nicht nur Bildgebung |
| **Kamera** | 42 MP+ mit Zoom | Details zur Defekterkennung |
| **Konnektivität** | 4G/LTE + Satelliten-Backup | Abdeckung in abgelegenen Gebieten |

[Vollständige Hardware-Spezifikationen →]({{< relref "/fleet/hardware" >}})

---

## Anwendungsfälle

### Inspektion von Übertragungsleitungen
Hochspannungsleitungen durch abgelegenes Gelände. Schäden an Isolatoren, Verschleiß an Leitern und Korrosion an Masten erkennen.

### Überwachung von Umspannwerken
Thermalanalyse von Transformatoren, Schaltanlagen und Verbindungen. Probleme erkennen, bevor es zu Ausfällen kommt.

### Analyse von Solarparks
Thermische Inspektion Modul für Modul. Hotspots, Verschmutzung und Degradation bei Tausenden von Modulen erkennen.

### Inspektion von Windkraftanlagen
Rotorblattinspektion ohne Industriekletterer. Risse, Erosion und Blitzschäden erkennen.

### Pipeline-Überwachung
Überwachung der Trasse, Leckerkennung, Eingriffe durch Dritte. Täglich Hunderte von Kilometern abdecken.

### Verteilnetz
Mastinspektion, Vegetationsmanagement und Bewertung von Hausanschlüssen in städtischen und vorstädtischen Gebieten.

---

## Einhaltung regulatorischer Vorgaben

Die Inspektion von Energieinfrastruktur erfordert die Einhaltung von:

- **NERC-Standards**: Anforderungen an die Zuverlässigkeit von Übertragungsnetzen
- **FAA Part 107**: Drohnenbetrieb, Ausnahmegenehmigungen für BVLOS
- **Pipeline-Sicherheit**: PHMSA-Inspektionsanforderungen
- **Umwelt**: Einschränkungen in Schutzgebieten, Rücksicht auf Wildtiere

Unsere Architektur unterstützt die Compliance mit:
- Vollständiger Flugprotokollierung mit Zeitstempeln und GPS
- Automatischer Umgehung von Flugverbotszonen
- Prüfpfaden für jede Inspektion
- Exportierbaren Berichten für die Einreichung bei Behörden

---

## Technischer Tiefgang

Diese Anwendung basiert auf unserer produktionsreifen Flottenarchitektur:

| Komponente | Rolle in der Energiebranche |
|-----------|---------------|
| [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) | Koordiniert Inspektionen über abgelegene Regionen hinweg |
| [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) | Thermische und visuelle KI-Verarbeitung |
| [JetStream-Streams]({{< relref "/fleet/streams" >}}) | Anlagenhistorie, Defektverfolgung |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Ausfallsicherer Betrieb in der Nähe von Hochspannung |

[Gesamte Architektur erkunden →]({{< relref "/fleet" >}})

---

## Jetzt starten

Jede Stunde Inspektionsflug ersetzt Tage manueller Arbeit. Lassen Sie uns über Ihren Bedarf an Infrastrukturinspektionen sprechen.

[Kontakt aufnehmen →](/contact)
