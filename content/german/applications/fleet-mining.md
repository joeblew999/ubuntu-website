---
title: "Bergbau & Rohstoffe"
meta_title: "IoT-Flotte für Bergbau & Rohstoffe | Ubuntu Software"
description: "Bergbaubetriebe mit IoT-Flotten: Drohnen für Vermessungen, Muldenkipper und Radlader mit Telematik, Grubensensoren für Sicherheit und Bewetterung. Einheitliche Architektur für standortweite Intelligenz."
image: "/images/robotics.svg"
draft: false
---

## Messen, was sich bewegt

Bergbaubetriebe bewegen täglich enorme Mengen. Genaue Messungen sind für Bestandsführung, Planung und Compliance unerlässlich. IoT-Flotten liefern Präzision im großen Maßstab.

---

## Flottentypen im Bergbau

| Flottentyp | Geräte | Rolle |
|------------|---------|------|
| **Luftgestützt** | Vermessungsdrohnen, Inspektionsdrohnen | Haldenvermessung, Grubenvermessung, Sprenganalyse |
| **Bodenfahrzeuge** | Muldenkipper, Radlader, Bohrgeräte | Geräteortung, Telematik, Produktionsüberwachung |
| **Stationäres IoT** | Grubensensoren, Bewetterungsmonitore, Gasdetektoren | Sicherheitsüberwachung, Einhaltung von Umweltauflagen |

Alle Gerätetypen nutzen dieselbe Architektur: Edge-KI, NATS-JetStream-Messaging und ein einheitliches Flottenmanagement.

---

## Die Herausforderung

Betriebe in Bergbau und Rohstoffgewinnung haben bei der Messung mit mehreren Problemen zu kämpfen:

- **Volumengenauigkeit**: Herkömmliche Vermessungen haben 10–15 % Fehler. Das sind Millionen an falsch erfasstem Bestand
- **Vermessungshäufigkeit**: Bodenvermessungen sind langsam. Die Bedingungen ändern sich schneller, als man messen kann
- **Sicherer Zugang**: Instabile Halden, befahrene Förderstraßen und Sprengzonen schränken den Zugang am Boden ein
- **Größenordnung**: Tagebaue erstrecken sich über Kilometer. In Untertagebetrieben ist die Sicht begrenzt
- **Umwelt-Compliance**: Aufsichtsbehörden verlangen eine genaue Erfassung der in Anspruch genommenen Flächen

Vierteljährliche Bodenvermessungen können mit dem täglichen Betrieb nicht Schritt halten. Satellitenbildern fehlt die Auflösung für genaue Volumenberechnungen.

---

## Wie IoT-Flotten das lösen

### Luftgestützt: Hochfrequente Messung

Häufiger und genauer vermessen:

- **Tägliche Überflüge**: Veränderungen an Halden in Echtzeit verfolgen
- **Zentimetergenauigkeit**: RTK-GPS und Photogrammetrie liefern einen Volumenfehler unter 2 %
- **Vollständige Abdeckung**: Das gesamte Gelände vermessen, nicht nur die zugänglichen Bereiche
- **Einheitliche Methodik**: Jedes Mal derselbe Ablauf, vergleichbare Ergebnisse

### Boden: Geräteintelligenz

Muldenkipper, Radlader und Bohrgeräte mit Echtzeit-Tracking:

- **Flottentelematik**: Position, Geschwindigkeit, Kraftstoff, Betriebsstunden, Nutzlast
- **Produktionserfassung**: Bewegte Tonnen, abgeschlossene Zyklen, Effizienzkennzahlen
- **Wartungswarnungen**: Motorfehlercodes, Serviceintervalle, Reifendruck
- **Kollisionsvermeidung**: Näherungserkennung, Personenschutz

### Stationäres IoT: Sicherheit & Umwelt

Sensoren im gesamten Betrieb:

- **Gasdetektoren**: Methan-, CO- und Sauerstoffwerte in engen Räumen
- **Bewetterungsmonitore**: Luftstrom in Untertagebetrieben
- **Bodenbewegung**: Böschungsstabilität, Setzungserkennung
- **Staubmonitore**: Einhaltung von Umweltauflagen, Arbeitssicherheit

### Betriebliche Intelligenz

Mehr als nur Messung:

- **Sprengplanung**: Vermessungen vor und nach der Sprengung für die Fragmentierungsanalyse
- **Optimierung von Förderstraßen**: Steigungsanalyse und Priorisierung der Instandhaltung
- **Wassermanagement**: Pegelstände von Becken und Entwässerungsmuster überwachen
- **Rekultivierungsnachweis**: Die fortschreitende Rekultivierung dokumentieren

---

## Onboard-KI (Edge-Verarbeitung)

Echtzeit-Inferenz auf NVIDIA Jetson für sofortige Erkenntnisse:

| Fähigkeit | Anwendung im Bergbau |
|------------|-------------------|
| **Haldenerkennung** | Materialhalden automatisch erkennen und klassifizieren |
| **Volumenschätzung** | Grobe Volumenschätzung in Echtzeit während des Flugs |
| **Fahrzeugerkennung** | Muldenkipper, Radlader und Personal verfolgen |
| **Sicherheitszonenüberwachung** | Unbefugten Zutritt zu Sperrbereichen erkennen |
| **Änderungserkennung** | Wesentliche Änderungen seit der letzten Vermessung erkennen |

**Warum Edge wichtig ist:**
- Daten vor Ort verarbeiten und weniger hochladen müssen
- Sofortige Warnungen bei Sicherheitsverstößen
- Weiterarbeiten auch bei eingeschränkter Konnektivität
- Schnelle Rückmeldung, ob die Vermessung vollständig ist

---

## Cloud-KI (Flottenanalyse)

Mit NATS JetStream den gesamten Betrieb optimieren:

| Fähigkeit | Anwendung im Bergbau |
|------------|-------------------|
| **Volumetrische Trendanalyse** | Materialbewegungen über Tage, Wochen und Monate verfolgen |
| **Bestandsabgleich** | Drohnenmessungen mit Produktionsaufzeichnungen vergleichen |
| **Auf-/Abtragsanalyse** | Erdbewegungen für eine genaue Planung berechnen |
| **Fortschrittsüberwachung** | Den Abbau mit dem Abbauplan abgleichen |
| **Compliance-Berichte** | Berichte für Umweltbehörden erstellen |

**Warum die Cloud wichtig ist:**
- Messungen mehrerer Standorte zusammenführen
- Historische Vergleiche für Trendanalysen
- Integration in Bergbauplanungssoftware
- Unterstützung bei Audit- und Compliance-Anforderungen

---

## Hardwarekonfiguration

Für Einsätze im Bergbau empfehlen wir:

| Komponente | Auswahl | Begründung |
|-----------|-----------|-----------|
| **Flugzeugzelle** | Industrieller Starrflügler oder Schwerlast-Multirotor | Lange Flugzeit für große Gelände, windtolerant |
| **Flugsteuerung** | Pixhawk 6X | Zuverlässig in staubigen Umgebungen mit starken Vibrationen |
| **Sensor-Companion** | Raspberry Pi CM4 | Koordination von RTK-GPS und Kamera |
| **KI-Companion** | Jetson Orin Nano | Verarbeitung vor Ort, Fahrzeugerkennung |
| **Vermessungskamera** | 45 MP+ mit mechanischem Verschluss | Scharfe Bilder für die Photogrammetrie |
| **RTK-GPS** | Zweifrequenzempfänger | Positionsgenauigkeit im Zentimeterbereich |
| **LiDAR** | Optional, Solid-State | Direkte 3D-Messung durch Vegetation hindurch |

[Vollständige Hardwarespezifikationen →]({{< relref "/fleet/hardware" >}})

---

## Anwendungsfälle

### Haldenvermessung
Tägliche Volumenvermessungen für das Bestandsmanagement. ROM-, Produkt- und Abraumhalden mit Materialklassifizierung.

### Grubenvermessung
Regelmäßige topografische Vermessungen für die Abbauplanung. Den Abbaufortschritt mit der Planung abgleichen.

### Inspektion von Förderstraßen
Steigungsanalyse, Bewertung des Fahrbahnzustands und Instandhaltungsplanung.

### Sprenganalyse
Vermessungen vor der Sprengung für die Planung. Vermessungen nach der Sprengung für Fragmentierung und Haufwerksvolumen.

### Überwachung von Absetzanlagen (Tailings)
Damminspektion, Vermessung des Spülstrands und Überwachung des Wasserstands für Sicherheit und Compliance.

### Unterstützung der Exploration
Schnelle Geländekartierung neuer Gebiete. Erkennung geologischer Merkmale.

---

## Integration

Bergbaubetriebe müssen an bestehende Systeme angebunden werden:

- **Bergbauplanungssoftware**: Export nach Deswik, Surpac, MineSight
- **Flottenmanagement**: Abgleich mit Systemen zur Lkw-Disposition
- **ERP-Integration**: Bestandsdaten an SAP und Oracle übermitteln
- **Umweltsysteme**: Compliance-Dokumentation und Berichterstattung

Unsere Architektur bietet:
- Exporte in Standardformaten (LAS, DXF, GeoTIFF)
- API-Integration für automatisierte Workflows
- Geplante Vermessungen mit automatischer Verarbeitung
- Aufbewahrung historischer Daten für Audits

---

## Technischer Einblick

Diese Anwendung baut auf unserer produktionsreifen Flottenarchitektur auf:

| Komponente | Rolle im Bergbau |
|-----------|---------------|
| [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) | Koordiniert Vermessungen über mehrere Gruben und Standorte |
| [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) | RTK-Positionierung, Verarbeitung der Vermessungsdaten |
| [JetStream-Streams]({{< relref "/fleet/streams" >}}) | Vermessungshistorie, Volumentrends |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Ausfallsicherer Betrieb in aktiven Abbaubereichen |

[Vollständige Architektur erkunden →]({{< relref "/fleet" >}})

---

## Jetzt starten

Jeder Prozentpunkt an Volumengenauigkeit ist bares Geld wert. Lassen Sie uns über Ihre Anforderungen an die Vermessung sprechen.

[Kontakt aufnehmen →](/contact)
