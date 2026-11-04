---
title: "Bau & Infrastruktur"
meta_title: "IoT-Flotte für das Baugewerbe | Ubuntu Software"
description: "Baustellenüberwachung mit IoT-Flotten: Drohnen für Vermessungen, Bagger und Lader mit Telematik, Baustellensensoren für Sicherheit und die Einhaltung von Umweltauflagen. Alles verwaltet mit einer einheitlichen Architektur."
image: "/images/robotics.svg"
draft: false
---

## Augen auf jeder Baustelle, jeden Tag

Bauprojekte scheitern, wenn die Transparenz fehlt. IoT-Flotten sorgen für kontinuierlichen Überblick – sie verfolgen den Fortschritt, dokumentieren Zustände und erkennen Probleme, bevor sie teuer werden.

---

## Flottentypen im Bauwesen

| Flottentyp | Geräte | Rolle |
|------------|---------|------|
| **Luftgestützt** | Vermessungsdrohnen, Inspektionsdrohnen | Fortschrittsverfolgung, Baustellenkartierung, Bauwerksinspektion |
| **Bodenfahrzeuge** | Bagger, Lader, Materialtransporter | Geräteortung, Telematik, Auslastungsüberwachung |
| **Stationäres IoT** | Baustellensensoren, Zutrittskontrolle, Geräteüberwachung | Einhaltung von Sicherheitsvorschriften, Umweltüberwachung, Objektschutz |

Alle Gerätetypen nutzen dieselbe Architektur: Edge-KI, NATS JetStream-Messaging und einheitliches Flottenmanagement.

---

## Die Herausforderung

Bau und Infrastruktur kämpfen mit Transparenzlücken:

- **Fortschrittsverfolgung** — Manuelle Berichte kommen verspätet und sind subjektiv
- **Einhaltung von Sicherheitsvorschriften** — Inspektoren können nicht überall sein
- **Dokumentation** — Fotos erfassen Kontext und Zeitstempel nicht zuverlässig
- **Änderungserkennung** — Unerlaubte Arbeiten oder Abweichungen von den Plänen erkennen
- **Anlageninspektion** — Brücken, Türme und Pipelines, die sich über riesige Entfernungen erstrecken

Bei herkömmlichen Methoden klettern Menschen auf Bauwerke, fahren Strecken ab oder verlassen sich auf planmäßige Inspektionen. Zwischen den Besuchen bleiben Probleme verborgen.

---

## Wie IoT-Flotten das lösen

### Luftgestützt: Kontinuierliche Baustellenüberwachung

Automatisierte tägliche Drohnenflüge erfassen:

- **Orthomosaik-Karten** — Zentimetergenaue Luftbilder der Baustelle
- **3D-Punktwolken** — Volumendaten zur Verfolgung von Erdarbeiten
- **Wärmebildgebung** — Erkennung von Feuchtigkeit, Dämmungsproblemen und Wärmeentwicklung an Geräten
- **360°-Dokumentation** — Vollständige visuelle Aufzeichnung mit Zeitstempeln

### Boden: Intelligente Geräteüberwachung

Telematik von jeder Maschine auf der Baustelle:

- **Auslastungsverfolgung** — Betriebsstunden, Leerlaufzeit, Kraftstoffverbrauch
- **Standorterfassung** — Echtzeitposition aller Geräte
- **Wartungsmeldungen** — Motorfehlercodes, Serviceintervalle, Fehlererkennung
- **Geofencing** — Einhaltung von Einsatzgrenzen, Diebstahlschutz

### Stationäres IoT: Überblick über die Baustelle

Sensoren im gesamten Projekt:

- **Luftqualitätsmessgeräte** — Staub, Emissionen, Nachweis der Grenzwerteinhaltung
- **Lärmsensoren** — Einhaltung gesetzlicher Vorgaben, gute Nachbarschaftsbeziehungen
- **Zutrittskontrolle** — Wer ist wann und wie lange auf der Baustelle
- **Bauwerksüberwachung** — Erkennung von Setzungen, Erschütterungen und Bewegungen

### Infrastrukturinspektion

Systematische Inspektion ohne Gerüste oder Fahrspursperrungen:

- **Brückeninspektion** — Dokumentation von Fahrbahn, Überbau und Unterbau
- **Turminspektion** — Telekommunikation, Stromübertragung, Windkraftanlagen
- **Pipeline-Befliegung** — Trassenüberwachung, Leckerkennung
- **Gebäudehülle** — Fassadeninspektion, Bewertung des Dachzustands

---

## Onboard-KI (Edge-Verarbeitung)

Echtzeit-Inferenz auf NVIDIA Jetson für sofortige Erkenntnisse:

| Fähigkeit | Anwendung im Bauwesen |
|------------|-------------------------|
| **Fortschrittserkennung** | Abgeschlossene vs. geplante Bauelemente erkennen |
| **Sicherheitsverstöße** | Fehlende PSA, unzureichende Absperrungen und Absturzgefahren erkennen |
| **Rissdetektion** | Oberflächenschäden an Beton, Asphalt und Stahl |
| **Korrosionserkennung** | Rostmuster an Brücken, Tanks und Stahlkonstruktionen |
| **Thermische Anomalien** | Hotspots, die auf elektrische Probleme oder eindringende Feuchtigkeit hinweisen |

**Warum Edge wichtig ist:**
- Sicherheitsprobleme in Echtzeit melden, damit sofort reagiert werden kann
- Große Orthomosaik-Datensätze lokal verarbeiten
- An abgelegenen Standorten mit eingeschränkter Konnektivität arbeiten
- Datenkosten bei häufigen Flügen senken

---

## Cloud-KI (Flottenanalyse)

Über NATS JetStream werden Erkenntnisse projektübergreifend zusammengeführt:

| Fähigkeit | Anwendung im Bauwesen |
|------------|-------------------------|
| **BIM-Abgleich** | As-built-Punktwolken mit Planungsmodellen überlagern |
| **Fortschrittsanalyse** | Fertigstellungsgrad, Terminabweichung, Trendanalyse |
| **Mängelverfolgung** | Mängel erfassen, Behebung verfolgen, Reparaturen überprüfen |
| **Flotteneinsatzplanung** | Drohnen effizient über mehrere Baustellen hinweg einsetzen |
| **Compliance-Berichte** | Inspektionsberichte und Sicherheitsdokumentation erstellen |

**Warum die Cloud wichtig ist:**
- Aktuellen Zustand mit der Planung vergleichen
- Fortschritt über ein ganzes Projektportfolio hinweg verfolgen
- Historische Analysen für Nachtragsforderungen und Streitfälle
- Integration in Projektmanagementsysteme

---

## Hardware-Konfiguration

Für Einsätze im Bauwesen empfehlen wir:

| Komponente | Auswahl | Begründung |
|-----------|-----------|-----------|
| **Flugwerk** | Holybro X500 | Stabile Plattform für Kartierungs-Nutzlasten |
| **Flugcontroller** | Pixhawk 6X | RTK-GPS für Vermessungsgenauigkeit |
| **Sensor-Companion** | Raspberry Pi CM4 | GPS/IMU-Datenaufzeichnung |
| **KI-Companion** | Jetson Orin NX | Mehr Rechenleistung für die 3D-Verarbeitung |
| **Nutzlast** | RGB- + Wärmebildkamera | Visuelle und thermische Dokumentation |
| **Konnektivität** | 4G/LTE + eSIM | Abdeckung städtischer und abgelegener Standorte |

[Vollständige Hardware-Spezifikationen →]({{< relref "/fleet/hardware" >}})

---

## Anwendungsfälle

### Aktive Baustellen
Tägliche Fortschrittsdokumentation. Verfolgung von Erdbewegungsvolumen. Sicherheitsüberwachung.

### Brückeninspektion
NBIS-konforme Dokumentation. Erfassung von Fahrbahn, Überbau und Unterbau.

### Linieninfrastruktur
Pipelines, Stromleitungen, Autobahnen. Befliegung und Inspektion im großen Maßstab.

### Gebäudeinspektion
Bewertung des Fassadenzustands. Dachvermessungen. Thermische Analyse der Gebäudehülle.

### Bergbau
Haldenvermessung. Abbaufortschritt in der Grube. Überwachung der Transportwege.

---

## Integration

Verbinden Sie die Daten Ihrer Drohnenflotte mit Ihren bestehenden Systemen:

- **BIM-Plattformen** — Autodesk Construction Cloud, Bentley, Trimble
- **Projektmanagement** — Procore, PlanGrid, Oracle Aconex
- **GIS-Systeme** — Esri ArcGIS, QGIS
- **Asset-Management** — IBM Maximo, SAP PM

---

## Technischer Einblick

Diese Anwendung basiert auf unserer produktionsreifen Flottenarchitektur:

| Komponente | Rolle im Bauwesen |
|-----------|---------------------|
| [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) | Synchronisiert Daten aus Baustellencontainern und verarbeitet sie in regionalen Hubs |
| [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) | Übernimmt die Photogrammetrie-Vorverarbeitung direkt auf der Drohne |
| [JetStream-Streams]({{< relref "/fleet/streams" >}}) | Speichert historische Bilddaten als Referenz für die Änderungserkennung |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Geofencing um aktive Arbeitsbereiche, Kräne und Personal |

[Gesamte Architektur erkunden →]({{< relref "/fleet" >}})

---

## Jetzt starten

Jedes Projekt profitiert von besserer Transparenz. Ob Sie eine Baustelle oder Hunderte verwalten – Drohnenflotten liefern die Dokumentation, die Sie brauchen.

[Kontakt aufnehmen →](/contact)
