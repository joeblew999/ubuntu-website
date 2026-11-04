---
title: "Landwirtschaft & Ackerbau"
meta_title: "IoT-Flotte für die Landwirtschaft | Ubuntu Software"
description: "Präzisionslandwirtschaft mit IoT-Flotten: Drohnen zur Überwachung, Traktoren und Erntemaschinen mit autonomer Lenkung sowie Feldsensoren für kontinuierliche Daten. Alles über eine einheitliche Architektur verwaltet."
image: "/images/robotics.svg"
draft: false
---

## Präzisionslandwirtschaft im großen Maßstab

Moderne Landwirtschaft verlangt Präzision – genau das auszubringen, was benötigt wird, genau dort, wo es benötigt wird, genau dann, wenn es benötigt wird. IoT-Flotten führen die Landwirtschaft weg von kalenderbasierten Abläufen hin zu datengetriebenen Entscheidungen.

---

## Flottentypen in der Landwirtschaft

| Flottentyp | Geräte | Rolle |
|------------|---------|------|
| **Luft** | Vermessungsdrohnen, Sprühdrohnen | Bestandsüberwachung, Schädlingserkennung, Präzisionsspritzung |
| **Bodenfahrzeuge** | Traktoren, Erntemaschinen, Bewässerungsroboter | Autonome Lenkung, teilflächenspezifische Ausbringung, Ernte |
| **Stationäres IoT** | Wetterstationen, Bodensensoren, Wasserzähler | Kontinuierliche Umweltüberwachung, Bewässerungssteuerung |

Alle Gerätetypen nutzen dieselbe Architektur: Edge-KI, NATS JetStream Messaging und ein einheitliches Flottenmanagement.

---

## Die Herausforderung

Die traditionelle Landwirtschaft arbeitet im Blindflug:

- **Flächendeckende Ausbringung** — Ganze Felder werden behandelt, obwohl nur einzelne Stellen Aufmerksamkeit brauchen
- **Verspätete Erkennung** — Probleme, die vom Boden aus sichtbar sind, sind bereits gravierend
- **Personalmangel** — Nicht genug Personal, um jeden Hektar zu kontrollieren
- **Wetterfenster** — Begrenzte Zeit zum Spritzen erfordert Tempo und Präzision
- **Dokumentation** — Gesetzliche Vorgaben verlangen lückenlose Aufzeichnungen

Satellitenbilder helfen, doch die Auflösung ist gering und der Zeitpunkt unvorhersehbar. Feldbegehungen lassen sich nicht skalieren. Landwirte brauchen Augen überall, und zwar jederzeit.

---

## Wie IoT-Flotten das lösen

### Luft: Kontinuierliche Überwachung

Automatisierte Drohnenflüge decken jeden Hektar nach Zeitplan ab:

- **Tägliche oder wöchentliche Befliegungen** — Probleme frühzeitig erkennen
- **Lückenlose Abdeckung** — Keine Lücken, keine ausgelassenen Bereiche
- **Wetterabhängige Planung** — Fliegen, wenn die Bedingungen es zulassen
- **Multispektrale Bildgebung** — Sehen, was das Auge nicht sieht

### Boden: Präzise Ausbringung

Traktoren und Feldspritzen mit teilflächenspezifischer Steuerung:

- **Punktbehandlung** — Unkraut spritzen, nicht die Kultur
- **Gezielte Schädlingsbekämpfung** — Nur dort ausbringen, wo ein Befall vorliegt
- **Düngeoptimierung** — Ausbringung am Bedarf der Pflanzen ausrichten
- **Autonome Lenkung** — RTK-GPS-Präzision für exakte Reihenführung

### Stationäres IoT: Umweltintelligenz

Sensoren im gesamten Betrieb:

- **Bodenfeuchte** — Bewässerungsentscheidungen in Echtzeit
- **Wetterstationen** — Lokale Mikroklimadaten
- **Durchflussmesser** — Verbrauchserfassung, Leckerkennung
- **Getreidesilo-Überwachung** — Temperatur- und Feuchtigkeitswarnungen

---

## KI an Bord (Edge-Verarbeitung)

Echtzeit-Inferenz auf NVIDIA Jetson ermöglicht sofortiges Handeln:

| Fähigkeit | Landwirtschaftliche Anwendung |
|------------|-------------------------|
| **Unkrauterkennung** | Invasive Arten in Pflanzenreihen in Echtzeit identifizieren |
| **Schädlingsidentifikation** | Fraßschadensmuster und Pilzinfektionen erkennen |
| **Bewertung der Pflanzengesundheit** | NDVI-Analyse, Stresserkennung aus Multispektraldaten |
| **Pflanzenzählung** | Bestandszählung, Auflaufkontrolle, Bestandsdichtekarten |
| **Spritzsteuerung** | Düsengenaue Präzision auf Basis der Erkennungsergebnisse |

**Warum Edge entscheidend ist:**
- Spritzentscheidungen in Millisekunden statt Sekunden
- Funktioniert auch in Gebieten ohne Mobilfunkabdeckung
- Verarbeitet hochauflösende Bilder lokal
- Senkt die Kosten für die Datenübertragung

---

## Cloud-KI (Flottenanalyse)

Über NATS JetStream gewinnen Sie Erkenntnisse aus Ihrem gesamten Betrieb:

| Fähigkeit | Landwirtschaftliche Anwendung |
|------------|-------------------------|
| **Ertragsprognose** | Historische Muster + aktuelle Bedingungen = Ernteschätzungen |
| **Behandlungsoptimierung** | Welche Maßnahmen unter welchen Bedingungen am besten wirken |
| **Flottenkoordination** | Mehrere Drohnen für eine effiziente Abdeckung planen |
| **Compliance-Berichte** | Automatisierte Spritzprotokolle, Ausbringungsnachweise |
| **Trendanalyse** | Saisonübergreifende Muster, Karten der Bodenvariabilität |

**Warum die Cloud entscheidend ist:**
- Erkenntnisse über Tausende Hektar hinweg
- Historischer Vergleich mit früheren Saisons
- Integration mit Farm-Management-Systemen
- Automatisierung der behördlichen Berichterstattung

---

## Hardware-Konfiguration

Für landwirtschaftliche Einsätze empfehlen wir:

| Komponente | Auswahl | Begründung |
|-----------|-----------|-----------|
| **Flugrahmen** | Holybro X500 oder landwirtschaftliche Sprühdrohne | Nutzlastkapazität für Sensoren + Tank |
| **Flugsteuerung** | Pixhawk 6X | Zuverlässige Autonomie, RTK-GPS-Unterstützung |
| **Sensor-Begleitrechner** | Raspberry Pi CM4 | Umweltsensoren, grundlegende Telemetrie |
| **KI-Begleitrechner** | Jetson Orin Nano | Echtzeit-Inferenz für die Erkennung |
| **Nutzlast** | Multispektralkamera + Sprühsystem | Erkennungs- und Behandlungsfähigkeit |
| **Konnektivität** | 4G/LTE + eSIM | Abdeckung in ländlichen Gebieten |

[Vollständige Hardware-Spezifikationen →]({{< relref "/fleet/hardware" >}})

---

## Anwendungsfälle

### Überwachung von Reihenkulturen
Mais, Sojabohnen, Weizen, Baumwolle. Großbetriebe, in denen Effizienz am meisten zählt.

### Weinberg- und Obstplantagenmanagement
Präzisionsspritzung in hochwertigen Kulturen. Behandlung pro Rebstock oder pro Baum.

### Sonderkulturen
Gemüse, Beeren, Baumschulware. Kulturen, bei denen die Arbeitskosten dominieren.

### Viehhaltung
Weideüberwachung, Zaunkontrolle, Herdenverfolgung.

---

## Integration

Verbinden Sie die Daten Ihrer Drohnenflotte mit Ihren bestehenden Systemen:

- **Farm-Management-Software** — John Deere Operations Center, Climate FieldView
- **Precision-Farming-Plattformen** — Applikationskarten für die teilflächenspezifische Ausbringung
- **Werkzeuge für Agronomen** — Feldbegehungsberichte, Empfehlungssysteme
- **Compliance-Systeme** — Spritzprotokolle, Zertifizierungsunterlagen

---

## Technischer Tiefgang

Diese Anwendung basiert auf unserer produktionsreifen Flottenarchitektur:

| Komponente | Rolle in der Landwirtschaft |
|-----------|---------------------|
| [NATS-Topologie]({{< relref "/fleet/nats-topology" >}}) | Verarbeitet Telemetrie von abgelegenen Feldern mit lückenhafter Konnektivität |
| [Fahrzeug-Gateway]({{< relref "/fleet/gateway" >}}) | Verarbeitet Multispektraldaten, löst Spritzentscheidungen aus |
| [JetStream-Streams]({{< relref "/fleet/streams" >}}) | Speichert historische Bilddaten und Behandlungsnachweise |
| [Sicherheitsmodell]({{< relref "/fleet/safety" >}}) | Gewährleistet einen sicheren Betrieb in der Nähe von Arbeitskräften, Tieren und Gebäuden |

[Vollständige Architektur erkunden →]({{< relref "/fleet" >}})

---

## Jetzt starten

Präzisionslandwirtschaft erfordert eine präzise Infrastruktur. Ob Sie Hunderte oder Hunderttausende Hektar bewirtschaften – Drohnenflotten wachsen mit Ihnen.

[Kontakt aufnehmen →](/contact)
