---
title: "Logistik & Lieferung"
meta_title: "IoT-Flotte für die Logistik | Ubuntu Software"
description: "Logistik mit IoT-Flotten: Lieferdrohnen, Lkw und fahrerlose Transportfahrzeuge (AGVs) mit Telematik, Lagersensoren für Bestandsführung und Rampenbetrieb. Einheitliche Architektur für flottenweite Optimierung."
image: "/images/robotics.svg"
draft: false
---

## Die Zukunft des Fulfillments

Die letzte Meile ist der teuerste Teil der Logistik. IoT-Flotten verändern die Wirtschaftlichkeit: schnellere Zustellung, niedrigere Kosten, weniger Verkehrsbelastung.

---

## Flottentypen in der Logistik

| Flottentyp | Geräte | Rolle |
|------------|---------|------|
| **Luft** | Lieferdrohnen, Inventurdrohnen | Zustellung auf der letzten Meile, Lagerinventuren |
| **Bodenfahrzeuge** | Lkw, Gabelstapler, FTS | Transport, Materialfluss, Lagerautomatisierung |
| **Stationäres IoT** | Rampensensoren, Fördersysteme, Inventar-Tags | Standortbetrieb, Durchsatzüberwachung, Asset-Tracking |

Alle Gerätetypen teilen dieselbe Architektur: Edge-KI, NATS-JetStream-Messaging und ein einheitliches Flottenmanagement.

---

## Die Herausforderung

Logistikbetriebe stehen unter wachsendem Druck:

- **Kosten der letzten Meile**: Bis zu 50 % der gesamten Versandkosten entfallen auf die finale Zustellung
- **Erwartungen an die Geschwindigkeit**: Forderungen nach Zustellung am selben Tag oder innerhalb der nächsten Stunde
- **Personalengpässe**: Fahrermangel, steigende Löhne
- **Verkehrsbelastung in Städten**: Staus, begrenzte Parkmöglichkeiten
- **Umweltdruck**: Vorgaben zur Emissionsreduzierung

Bodengebundene Zustellung skaliert schlecht. Jedes zusätzliche Paket bedeutet einen weiteren Stopp, eine weitere Fahrerminute, einen weiteren Kilometer Treibstoff.

---

## Wie IoT-Flotten das lösen

### Luft: Direkt von Punkt zu Punkt

Das Straßennetz komplett umgehen:

- **Routen in Luftlinie**: Kein Verkehr, keine Umwege
- **Gleichzeitige Zustellungen**: Mehrere Drohnen, mehrere Pakete, zur selben Zeit
- **Überall erreichbar**: Ländliche Gebiete, abgeschlossene Wohnanlagen, Dächer
- **Planbare Zeiten**: Keine Schwankungen durch den Verkehr

### Boden: Flottenintelligenz

Lkw, Gabelstapler und FTS mit Koordination in Echtzeit:

- **Routenoptimierung**: Dynamische Routen auf Basis von Verkehr, Zustellungen und Kapazität
- **Telematik**: Treibstoff, Wartung, Lenkzeiten, Beladungsstatus
- **FTS-Koordination**: Lagerroboter, die reibungslos zusammenarbeiten
- **Zustellnachweis**: Bestätigung mit GPS-Stempel und Fotos

### Stationäres IoT: Standortbetrieb

Sensoren entlang der gesamten Lieferkette:

- **Sensoren an Rampentoren**: Ankunft der Lkw, Abfahrtszeiten
- **Fördertechnik-Überwachung**: Durchsatz, Staus, Sortiergenauigkeit
- **Inventar-Tags**: Bestände in Echtzeit, Standortverfolgung
- **Umweltsensoren**: Temperatur und Luftfeuchtigkeit für verderbliche Waren

### Lagerintegration

Alle Flottentypen erweitern die Reichweite des Lagers:

- **Vom Verteilzentrum zum Kunden**: Leichte Pakete ganz ohne Lkw zustellen
- **Micro-Fulfillment**: Vorgelagerte Standorte bestücken, per Drohne nachfüllen
- **Warenbewegung**: Transfers zwischen Standorten ohne Lkw
- **Retourenabholung**: Pakete so einfach abholen, wie sie zugestellt werden

---

## KI an Bord (Edge-Verarbeitung)

Inferenz in Echtzeit auf NVIDIA Jetson für den autonomen Betrieb:

| Fähigkeit | Anwendung in der Logistik |
|------------|----------------------|
| **Hindernisvermeidung** | Bäume, Leitungen und Bauwerke in Echtzeit umfliegen |
| **Bewertung der Landezone** | Sichere Zustellpunkte erkennen, Hindernisse identifizieren |
| **Paketverifizierung** | Richtiges Paket bestätigen, Zustand bei der Zustellung prüfen |
| **Wetteranpassung** | Flugparameter an Wind und Niederschlag anpassen |
| **Ausfallsichere Navigation** | Sicher weiterfliegen, wenn das GPS-Signal nachlässt |

**Warum Edge wichtig ist:**
- Navigationsentscheidungen im Bruchteil einer Sekunde
- Betrieb in Häuserschluchten mit schlechtem GPS-Empfang
- Unerwartete Hindernisse autonom bewältigen
- Sicherheit auch ohne Verbindung zur Bodenstation

---

## Cloud-KI (Flottenanalyse)

Über NATS JetStream das gesamte Netzwerk optimieren:

| Fähigkeit | Anwendung in der Logistik |
|------------|----------------------|
| **Routenoptimierung** | Gesamtflugzeit über alle Zustellungen minimieren |
| **Bedarfsprognose** — Zustellvolumen vorhersagen, Bestände vorab positionieren |
| **Flottenplanung** | Auslastung, Akkumanagement und Wartung ausbalancieren |
| **Dynamische Umverteilung** | Kapazitäten nach Echtzeitbedarf verlagern |
| **Leistungsanalyse** | Zustellzeiten, Erfolgsquoten, Effizienzkennzahlen |

**Warum die Cloud wichtig ist:**
- Tausende gleichzeitige Zustellungen optimieren
- Aus jedem Flug lernen, um künftige Routen zu verbessern
- Last über das gesamte Distributionsnetz verteilen
- Integration mit Auftragsmanagementsystemen

---

## Hardware-Konfiguration

Für Logistikeinsätze empfehlen wir:

| Komponente | Auswahl | Begründung |
|-----------|-----------|-----------|
| **Flugzeugzelle** | Speziell gebaute Lieferdrohne | Nutzlastschacht, Abwurfmechanismus |
| **Flugsteuerung** | Pixhawk 6X | Zuverlässige Autonomie |
| **Sensor-Companion** | Raspberry Pi CM4 | Paketsensoren, Status-Telemetrie |
| **KI-Companion** | Jetson Orin Nano | KI für Navigation und Landung |
| **Nutzlast** | Paketkapazität 1–5 kg | Deckt den Großteil des E-Commerce ab |
| **Konnektivität** | 4G/LTE + eSIM | Abdeckung in Städten, Redundanz über mehrere Netzbetreiber |

[Vollständige Hardware-Spezifikationen →]({{< relref "/fleet/hardware" >}})

---

## Anwendungsfälle

### Zustellung auf der letzten Meile
E-Commerce-Pakete. Essenslieferungen. Apotheken- und medizinischer Bedarf.

### Zustellung auf dem Campus
Firmengelände. Universitätsgelände. Krankenhauskomplexe.

### Lagerbetrieb
Bestandstransfers. Inventurzählungen. Unterstützung bei der Entnahme aus Hochregalen.

### Notfallversorgung
Medizinischer Bedarf für abgelegene Orte. Ersatzteile für liegengebliebene Fahrzeuge.

### Zustellung im ländlichen Raum
Gebiete, die von der klassischen Logistik unterversorgt sind. Inseln, Berge, abgelegene Gemeinden.

---

## Regulatorischer Weg

Drohnenzustellung erfordert die Einhaltung gesetzlicher Vorschriften:

- **FAA Part 135**: Zertifizierung als Luftfahrtunternehmen für die Paketzustellung
- **Remote ID**: Ausgestrahlte Identifikation bei allen Einsätzen
- **BVLOS-Betrieb**: Ausnahmegenehmigungen für Flüge außerhalb der Sichtweite
- **Luftraumintegration**: Teilnahme an UTM, Vereinbarungen über Flugkorridore

Unsere Architektur unterstützt die Compliance durch:
- Vollständige Flugprotokollierung und Speicherung der Telemetrie
- Remote-ID-Ausstrahlung von jedem Fahrzeug
- Positionsmeldungen in Echtzeit an UTM-Systeme
- Prüfpfade für behördliche Kontrollen

---

## Technischer Tiefgang

Diese Anwendung basiert auf unserer produktionsreifen Flottenarchitektur:

| Komponente | Rolle in der Logistik |
|-----------|------------------|
| [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) | Koordiniert Tausende gleichzeitige Zustellungen |
| [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) | Navigations-KI, Bewertung der Landezone |
| [JetStream-Streams]({{< relref "/fleet/streams" >}}) | Zustellbestätigung, Nachweis der Übergabekette |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Ausfallsicherer Betrieb in besiedelten Gebieten |

[Gesamte Architektur entdecken →]({{< relref "/fleet" >}})

---

## Jetzt starten

Die Wirtschaftlichkeit der Drohnenzustellung überzeugt. Die Technologie ist bereit. Lassen Sie uns über Ihre logistischen Herausforderungen sprechen.

[Kontakt aufnehmen →](/contact)
