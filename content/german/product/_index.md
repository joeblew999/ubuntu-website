---
title: "Produkte"
meta_title: "KI + 3D Räumliche Intelligenz Plattform | Ubuntu Software"
description: "Wir geben KI die Fähigkeit, die physische 3D-Welt zu verstehen, darüber nachzudenken und mit ihr zu interagieren. Von Bauwesen über Fertigung bis zu digitalen Zwillingen."
draft: false
---

## KI denkt in 2D. Wir geben ihr räumliche Intelligenz.

Die heutige KI ist bemerkenswert leistungsfähig—sie kann lesen, schreiben, Bilder sehen und Gespräche führen. Aber sie denkt in flachen Dimensionen: Text ist 1D, Bilder sind 2D, Video ist 2D über Zeit.

**Die reale Welt ist 3D.**

Damit KI wirklich mit der physischen Realität interagieren kann—um Gebäude zu entwerfen, Roboter zu steuern, Anlagen zu verwalten oder Fertigung zu planen—braucht sie räumliche Intelligenz. Die Fähigkeit, Geometrie zu verstehen, über physische Einschränkungen nachzudenken und mit dreidimensionalem Raum zu arbeiten.

Das ist es, was wir bauen: eine Plattform, die KIs kognitive Fähigkeiten mit 3D-räumlichem Verständnis verbindet und Maschinen befähigt, mit der physischen Welt zu arbeiten.

---

## Anwendungen

{{< tabs >}}

{{< tab "Bauwesen" >}}

### Modulares & Vorgefertigtes Bauen

Die Bauindustrie verlagert sich hin zu fabrikgefertigten Modulhäusern. Vorfertigung bietet schnellere Bauzeiten, bessere Qualitätskontrolle und weniger Abfall—aber traditionelle CAD-Tools wurden nicht für diesen Workflow entwickelt.

**Die Herausforderung:**
- Verteilte Teams über Fabrikhallen, Ingenieurbüros und Baustellen
- Bedarf an Echtzeit-Zusammenarbeit, nicht dateibasierte Speichern-Zusammenführen-Commit-Zyklen
- Proprietäre Formate, die Herstellerbindung schaffen
- Altsysteme, die keine KI-Unterstützung integrieren können

**Unser Ansatz:**
- **Echtzeit-Zusammenarbeit** — Teams in Thailand, Deutschland und Australien arbeiten gleichzeitig am selben Modell mit Automerge CRDT-Technologie
- **100% Offene Standards** — STEP- und IFC-Formate stellen sicher, dass Ihre Entwürfe Ihnen gehören, unabhängig von Softwareänderungen zugänglich
- **KI-unterstütztes Design** — Befehle in natürlicher Sprache, Designoptimierungsvorschläge und automatische Fehlererkennung
- **Modul-zentrierter Workflow** — Erstklassige Unterstützung für Fertigbaukomponenten, Montagesequenzierung und Transportplanung

{{< /tab >}}

{{< tab "Digitale Zwillinge" >}}

### Digitale Zwillinge & Gebäudemanagement

Überbrücken Sie die Lücke zwischen Design und der realen Welt. Unser Ansatz basiert auf Erfahrungen mit KI-gestützten Gebäudemanagementsystemen, einschließlich der Zusammenarbeit mit [Bilfinger](https://www.bilfinger.com/), einem der führenden deutschen Ingenieur- und Dienstleistungsunternehmen.

**Fähigkeiten:**
- **Live-Sensor-Integration** — Verbinden Sie Ihre 3D-Modelle mit Echtzeit-IoT-Daten von Gebäuden und Infrastruktur
- **Vorausschauende Wartung** — KI-gestützte Analyse zur Identifizierung von Geräteproblemen, bevor sie zu Ausfällen werden
- **Was-Wäre-Wenn-Simulation** — Testen Sie Szenarien in Ihrem digitalen Zwilling, bevor Sie Änderungen in der Realität umsetzen
- **Gebäudemanagement** — Verwalten Sie Gebäudesysteme, verfolgen Sie Anlagen und optimieren Sie Abläufe über eine räumliche Schnittstelle

**Architektur:**
Aufgebaut auf NATS JetStream für ereignisgesteuerten Echtzeit-Datenfluss zwischen physischen Sensoren und digitalen Modellen. Dieselbe räumliche Intelligenz, die beim Entwurf eines Gebäudes hilft, kann es während seines gesamten Lebenszyklus überwachen und verwalten.

{{< /tab >}}

{{< tab "Fertigung" >}}

### Fertigung & Automatisierung

Wenn KI den 3D-Raum versteht, kann sie an der physischen Produktion teilnehmen—vom Design über die Fertigung bis zur Qualitätskontrolle.

**Anwendungen:**
- **Robotik-Integration** — Räumlichen Kontext für automatisierte Systeme bereitstellen, ermöglicht intelligentere Bahnplanung und Montageoperationen
- **Produktionsplanung** — KI-Unterstützung für Layoutoptimierung, Prozesssequenzierung und Ressourcenzuweisung
- **Qualitätskontrolle** — Ist-Geometrie mit Designabsicht vergleichen, automatische Erkennung von Abweichungen
- **Montagedokumentation** — Arbeitsanweisungen generieren, die dreidimensionale Beziehungen verstehen

**Über das Bauwesen hinaus:**
Während unser anfänglicher Fokus auf modularem Bauen liegt, gilt die räumliche Intelligenz der Plattform für jeden Fertigungsbereich—Automobil, Luft- und Raumfahrt, Konsumgüter, Industrieausrüstung. Überall, wo KI über physische Geometrie nachdenken muss.

{{< /tab >}}

{{< tab "Technologie" >}}

### Plattform-Architektur

Moderner Technologie-Stack für räumliche KI-Anwendungen.

**Offene Standards:**
- **STEP-Format** — ISO 10303 Standard für 3D-Geometrieaustausch
- **IFC-Unterstützung** — ISO 16739 Building Information Modeling Interoperabilität
- Keine proprietäre Bindung. Ihre Daten bleiben zugänglich.

**Echtzeit-Zusammenarbeit:**
- **Automerge** — CRDT-basierte Synchronisation für konfliktfreie Zusammenarbeit
- **Offline-First** — Volle Funktionalität ohne Internet, automatische Synchronisation bei Wiederverbindung
- **Versionskontrolle** — Vollständige Revisionshistorie mit Verzweigung und Zusammenführung

**KI-Integration:**
- **MCP-Protokoll** — Model Context Protocol für nahtlose KI-Tool-Integration
- **Natürliche Sprache** — Designabsicht in einfacher Sprache beschreiben
- **Konversationsschnittstelle** — Mit Ihren räumlichen Modellen chatten

**Ereignisgesteuerte Architektur:**
- **NATS JetStream** — Echtzeit-Messaging für digitale Zwillings-Fähigkeiten
- **Sensor-Konnektivität** — Brücke zwischen physischem IoT und digitalen Modellen

**Plattformübergreifend:**
- Web-native GUI mit Desktop- und Mobile-Unterstützung
- Windows, macOS, Linux, iOS, Android
- Self-Hosted oder Cloud-Deployment-Optionen

{{< /tab >}}

{{< /tabs >}}

---

## Aktueller Status

Wir entwickeln aktiv die Plattform und konzentrieren uns auf modulares Bauen als unseren ersten Bereich, während wir die zugrundeliegenden räumlichen Intelligenzfähigkeiten aufbauen, die branchenübergreifend anwendbar sind.

**Wir suchen:**
- Modulhaushersteller, die bessere Designwerkzeuge suchen
- Gebäudemanager, die an KI-gestützten digitalen Zwillingen interessiert sind
- Fertigungsunternehmen, die räumliche KI-Anwendungen erkunden

---

## Interessiert?

[Kontaktieren Sie uns](/de/contact), um mehr über unseren Fortschritt zu erfahren und zu besprechen, wie räumliche KI auf Ihren Bereich anwendbar sein könnte.
