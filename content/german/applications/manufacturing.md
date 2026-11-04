---
title: "Fertigung"
meta_title: "KI-gestützte Fertigung | Ubuntu Software"
description: "Vom Design zur Fertigung. KI, die Geometrie, Toleranzen und Produktionsbeschränkungen versteht."
image: "/images/manufacturing.svg"
draft: false
---

## Vom Design zur Fertigung

Ein Design ist nur so gut wie seine Fertigbarkeit. Die eleganteste Geometrie ist wertlos, wenn sie sich nicht schneiden, umformen, drucken oder montieren lässt. In der Lücke zwischen dem, was konstruiert, und dem, was produziert wird, explodieren die Kosten, verschieben sich Termine und versagt die Qualität.

KI, die Fertigung versteht, zeichnet nicht nur Formen. Sie konstruiert für die Produktion.

---

## Die Lücke

Konstruktion und Fertigung sprechen unterschiedliche Sprachen:

**Die Konstruktion sagt:** „Hier ist die Geometrie.“

**Die Fertigung fragt:** „Wie stelle ich das her? Welche Toleranz? Welche Arbeitsgänge? In welcher Reihenfolge? Welche Werkzeuge? Was wird das kosten?“

Heute erfolgt die Übersetzung zwischen diesen Welten manuell:

- Ingenieure prüfen Konstruktionen auf Fertigbarkeit – langsam
- CAM-Programmierer setzen Geometrie in Werkzeugbahnen um – manuell
- Kalkulatoren berechnen Kosten – auf Basis von Erfahrung und Schätzungen
- Probleme werden in der Werkstatt entdeckt – zu spät

**Jede Übergabe ist eine Fehlerquelle. Jede Übersetzung verliert Informationen.**

---

## Was wir ermöglichen

### Fertigungsgerechte Konstruktion

KI, die Produktionsrestriktionen bereits bei der Konstruktion versteht.

Keine Prüfung am Ende. Kontinuierliches Feedback während der Konstruktion:

- **Prozessbewusstsein** — Funktioniert diese Geometrie für CNC? Blech? Guss? Additive Fertigung?
- **Toleranzanalyse** — Ist das einhaltbar? Zu welchen Kosten?
- **Feature-Erkennung** — Bohrungen, Taschen und Dome werden automatisch analysiert
- **Kostenauswirkungen** — Konstruktionsentscheidungen in Echtzeit den Produktionskosten zugeordnet

**Fertigbarkeitsprobleme in der Konstruktion erkennen, nicht in der Werkstatt.**

---

### Intelligentes CAM

Von der Geometrie zur Werkzeugbahn mit KI-Unterstützung.

- **Automatische Feature-Erkennung** — KI identifiziert Fertigungsmerkmale aus der Volumengeometrie
- **Arbeitsgangfolge** — Optimale Reihenfolge der Arbeitsgänge wird generiert
- **Werkzeugbahnoptimierung** — Zykluszeit minimieren, Standzeit maximieren
- **Mehrachsstrategien** — Komplexe Geometrie, intelligenter Ansatz

**STEP-Geometrie rein. Optimierter G-Code raus.**

---

### Arbeitsplanung

Die Brücke zwischen Entwicklung und Produktion.

- **Fertigungsroute** — Welche Maschinen, welche Reihenfolge
- **Rüstplanung** — Vorrichtungen, Spannmittel, Bezugspunkte
- **Zeitschätzung** — Genaue Zykluszeiten, bevor der erste Span fällt
- **Ressourcenzuweisung** — Aufträge an verfügbare Kapazitäten anpassen

**KI, die Ihre Fertigung versteht – keine generischen Vorlagen.**

---

### Additive Fertigung

3D-Druck im Produktionsmaßstab erfordert Produktionsdenken.

- **Optimierung der Bauausrichtung** — Stützstrukturen minimieren, Festigkeit maximieren
- **Nesting und Packing** — Bauräume effizient füllen
- **Stützstrukturgenerierung** — Intelligente Platzierung, einfaches Entfernen
- **Optimierung der Prozessparameter** — Materialspezifisch, geometriebewusst
- **Hybride Workflows** — Additiv und subtraktiv kombiniert

**Additive Fertigung als echter Fertigungsprozess, nicht als Prototyping.**

---

### Qualitätsintegration

Die Prüfung kommt nicht erst zum Schluss.

- **GD&T-Interpretation** — Geometrische Bemaßung und Tolerierung nativ verstanden
- **Prüfplanung** — KMG-Programme aus der Konstruktionsabsicht generiert
- **Prozessbegleitende Prüfung** — Kritische Merkmale während der Produktion prüfen
- **Statistische Prozesslenkung** — Qualitätsdaten fließen zurück in die Konstruktion

**Von der Konstruktionsabsicht zum Prüfplan. Geschlossener Regelkreis.**

---

### Shopfloor-Anbindung

Konstruktionsentscheidungen treffen auf die Produktionsrealität.

**NATS JetStream** verbindet die Konstruktion mit der Fabrikhalle:

- **Maschinenüberwachung** — Status, Auslastung und Leistung in Echtzeit
- **Programmverteilung** — Das richtige Programm automatisch an die richtige Maschine
- **Produktionsrückmeldung** — Ist vs. Plan, kontinuierlich
- **Problemeskalation** — Probleme werden sofort sichtbar, nicht erst beim Versand

**Digitaler Faden von der Konstruktion bis zur Produktion.**

---

## Offene Standards

Ihre Fertigungsdaten gehören Ihnen.

**STEP (ISO 10303)** — Das Geometrie-Austauschformat, dem die Fertigung vertraut. AP203, AP214, AP242. Volle Genauigkeit, kein Konvertierungsverlust.

**Native Formate** — Import aus jedem CAD-System. Export in jedes CAM-System.

**Offene Integration** — Anbindung an Ihre ERP-, MES- und Qualitätssysteme über Standard-APIs.

Kein proprietärer Lock-in. Ihre Konstruktionen, Ihre Werkzeugbahnen, Ihre Daten.

---

## Anwendungsfälle

### CNC-Bearbeitung

Fräsmaschinen, Drehmaschinen, Mehrachs-Bearbeitungszentren. Wo Präzision auf Produktion trifft.

- Feature-basierte Bearbeitung aus Volumenmodellen
- Automatische Werkzeugbahngenerierung
- Zykluszeitoptimierung
- Flexible Postprozessoren

---

### Blechbearbeitung

Laser, Plasma, Stanzen, Biegen. Hohe Stückzahlen, knappe Margen.

- Abwicklungserzeugung mit Biegekompensation
- Nesting-Optimierung für maximale Materialausnutzung
- Biegefolgeplanung
- Integrierte Kostenkalkulation

---

### Metallbau & Schweißen

Tragwerke, Baugruppen, Schweißkonstruktionen.

- Nahtgestaltung und Zugänglichkeitsanalyse
- Optimierung der Schweißfolge
- Verzugsvorhersage
- Montagefolgeplanung

---

### Additive Fertigung

Pulverbett, FDM, Harz, Metall. Additive Serienfertigung.

- Baujobvorbereitung und Ausrichtung
- Stützstrukturoptimierung
- Mehrteile-Nesting
- Nachbearbeitungsplanung

---

### Montage

Aus Komponenten werden Produkte.

- Optimierung der Montagereihenfolge
- Überschneidungs- und Kollisionsprüfung
- Werkzeug- und Vorrichtungskonstruktion
- Erstellung von Arbeitsanweisungen

---

## Für Lohnfertiger

Jeder Auftrag anders. Jedes Angebot ein Glücksspiel.

- **Schnellere Angebotserstellung** — KI-gestützte Kalkulation aus 3D-Modellen
- **Kürzere Programmierzeit** — Intelligentes CAM macht aus Stunden Minuten
- **Weniger Überraschungen** — Fertigbarkeit wird vor der Auftragszusage geprüft
- **Bessere Margen** — Präzise kalkulieren, effizient produzieren

**Wettbewerb über Geschwindigkeit und Präzision, nicht nur über den Preis.**

---

## Für Serienfertiger

Stückzahlen erfordern Konsistenz.

- **Optimierte Prozesse** — Jede Zykluszeit ausgereizt
- **Eingebaute Qualität** — Nicht nachträglich erprüft
- **Kontinuierliche Verbesserung** — Produktionsdaten treiben die Weiterentwicklung der Konstruktion
- **Flexible Automatisierung** — Anpassung an Produktänderungen, ohne von vorn zu beginnen

**Skalieren ohne Qualitätseinbußen.**

---

## Für OEMs

Konstruktion und Fertigung unter einem Dach – oder über eine Lieferkette hinweg.

- **Simultaneous Engineering** — Konstruktion und Fertigung arbeiten gleichzeitig
- **Zusammenarbeit mit Lieferanten** — Konstruktionen teilen, ohne die Kontrolle zu verlieren
- **Kostengerechte Konstruktion** — Fertigungskosten in Echtzeit sichtbar
- **Digitaler Faden** — Rückverfolgbarkeit von der Anforderung bis zum ausgelieferten Teil

**Konstruieren Sie, was Sie bauen können. Bauen Sie, was Sie konstruiert haben.**

---

## Industrie 4.0

Die Fabrik der Zukunft braucht KI, die den Raum versteht.

Keine Dashboards und Schlagworte. Echte Fähigkeiten:

- **Digitaler Zwilling der Produktion** — Ihre Fabrikhalle auf der Plattform
- **Vorausschauende Wartung** — Anlagenzustand aus Sensordaten
- **Adaptive Planung** — In Echtzeit auf Änderungen reagieren
- **Kontinuierliche Optimierung** — KI findet Verbesserungen, die Menschen übersehen

**Industrie 4.0 auf Basis räumlicher Intelligenz, nicht nur Vernetzung.**

---

## Die Architektur

Infrastruktur in Fertigungsqualität.

| Ebene | Technologie | Zweck |
|-------|------------|---------|
| Geometrie | STEP nativ | Präzise CAD-Daten, volle Genauigkeit |
| Zusammenarbeit | Automerge CRDT | Entwicklung und Fertigung synchron |
| Messaging | NATS JetStream | Maschinenanbindung, Echtzeit-Ereignisse |
| KI | Model Context Protocol | Fertigungsbewusste Intelligenz |
| Integration | Offene APIs | ERP-, MES- und Qualitätssysteme |

**Gebaut für die Anforderungen der Produktion.**

---

## Den Kreis schließen mit Publish

Die Fertigung läuft ebenso mit Papierkram wie mit Maschinen.

- **Arbeitsanweisungen** — Schritt-für-Schritt-Montageanleitungen aus CAD generiert
- **Qualitätschecklisten** — Prüfformulare digital oder auf Papier erfasst, zurück in Ihr QMS
- **Stücklisten und Zuschnittlisten** — Teiledokumentation abgestimmt auf das 3D-Modell
- **Compliance-Nachweise** — Zertifizierungen, Prüfergebnisse, Rückverfolgbarkeitsdokumentation

Alles aus einer Quelle. Alles mit Ihren Fertigungsdaten verbunden.

[Publish entdecken →]({{< relref "/platform/publish" >}})

---

## Datensouveränität

Ihre Fertigungsdaten bleiben Ihre:

| Traditioneller Ansatz | Mit WellKnown |
|---------------------|----------------|
| Konstruktionen in der Cloud des Anbieters | Konstruktionen auf Ihrer Infrastruktur |
| Anbieter besitzt die Kundenbeziehung | Sie besitzen die Kundenbeziehung |
| CAD-Wechsel = Neuaufbau | Wechsel = Gateway neu verbinden |
| IP gegenüber Anbietern offengelegt | Sie kontrollieren den IP-Schutz |

Export nach STEP. Veröffentlichen AUF Plattformen. Niemals darin EINGESPERRT.

[Mehr über WellKnown erfahren →]({{< relref "/platform/publish#wellknown-gateway-data-sovereignty" >}})

---

## Jetzt starten

In der Fertigung trifft Design auf Realität. Dort wird Geometrie zum Produkt. Dort zählt Präzision.

KI, die sowohl Konstruktion als auch Produktion versteht. Genau das bauen wir.

[Kontaktieren Sie uns →](/contact/)
