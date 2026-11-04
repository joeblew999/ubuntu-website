---
title: "Robotik & Automatisierung"
meta_title: "Räumliche Intelligenz für die Robotik | Ubuntu Software"
description: "Jeder Roboter muss den Raum verstehen. Darin navigieren. Objekte darin handhaben. Darin mit Menschen zusammenarbeiten. Wir bauen die Infrastruktur für die Explosion der Robotik."
image: "/images/robotics.svg"
draft: false
---

## Räumliche Intelligenz für die Maschinen, die unsere Welt prägen werden

Jeder Roboter muss den Raum verstehen. Sich darin bewegen. Objekte darin handhaben. Darin mit Menschen zusammenarbeiten.

Doch die heutigen Robotik-Toolchains wurden nicht für KI gebaut. Sie wurden nicht für Echtzeit-Zusammenarbeit gebaut. Sie wurden nicht für die Explosion gebaut, die bevorsteht.

---

## Der Moment

**Die Robotik tritt in ihre exponentielle Phase ein.**

- **Humanoide sind Realität.** Tesla Optimus. Figure. 1X. Sanctuary. Agility. Chinesische Hersteller im Wettlauf um den Markt. Milliardeninvestitionen. Millionen geplanter Einheiten.
- **Lagerautomatisierung skaliert.** Jedes Logistikunternehmen, jeder E-Commerce-Anbieter, jeder Hersteller automatisiert in hohem Tempo.
- **Chirurgische Roboter verbreiten sich rasant.** Präzisionsmedizin verlangt nach Präzisionsmaschinen.
- **Autonome Fahrzeuge brauchen Umgebungen.** Nicht nur Straßen – Laderampen, Lagerhallen, Baustellen, landwirtschaftliche Betriebe.

Sie alle haben eine gemeinsame Anforderung: **KI, die in drei Dimensionen versteht und konstruiert.**

---

## Das Problem

**Die heutigen Werkzeuge sind fragmentiert:**

- **CAD lebt in einer Welt.** Mechanische Konstruktion, statische Dateien, manuelle Übergabe.
- **Simulation lebt in einer anderen.** Getrennte Umgebungen, getrennte Formate, getrennte Teams.
- **KI lebt in einer dritten.** Nachträglich angeflanscht. Betrachtet Screenshots. Kein natives Verständnis.
- **Zusammenarbeit?** E-Mail-Anhänge. Albträume bei der Dateiversionierung. „Welches Modell ist aktuell?“

Die Robotikbranche baut die Zukunft mit Werkzeugen aus der Vergangenheit.

---

## Was wir ermöglichen

### KI-native räumliche Schlussfolgerung

**Roboter brauchen nicht nur Geometrie – sie brauchen KI, die sie versteht.**

Die Integration des Model Context Protocol gibt KI-Systemen direkten Zugriff auf:

- **Arbeitsraumgeometrie** — Keine Bilder von Arbeitsräumen. Die tatsächliche 3D-Struktur.
- **Kollisionsgrenzen** — Was kann der Arm erreichen? Was muss er vermeiden?
- **Kinematische Ketten** — Gelenkgrenzen, Freiheitsgrade, Erreichbarkeitsanalyse
- **Semantisches Verständnis** — Das ist ein Greifer. Das ist ein Förderband. Hier ist die Zone für Menschen.

KI, die schlussfolgern kann: *„Wenn wir diese Vorrichtung um 200 mm nach links verschieben, erreicht der Roboter die Rückwand ohne Umorientierung.“*

---

### Arbeitszellen-Design

**Der Roboter ist nur ein Teil des Systems.**

Konstruieren Sie komplette Arbeitszellen mit KI-Unterstützung:

- **Layout-Optimierung** — Roboter, Vorrichtungen und Förderbänder für maximale Effizienz positionieren
- **Reichweitenanalyse** — Vor dem Aufbau prüfen, dass jedes Ziel erreichbar ist
- **Taktzeitsimulation** — Bahnen und Abläufe optimieren
- **Planung von Sicherheitszonen** — Zonen für die Mensch-Roboter-Kollaboration, Lichtvorhänge, Verriegelungen
- **Umrüstungsdesign** — Flexible Zellen, die sich an Produktvarianten anpassen

---

### Bahnplanung & Simulation

**Vom Design zur Bewegung.**

- **Kollisionsfreie Bahngenerierung** — KI-gestützte Trajektorienplanung
- **Multi-Roboter-Koordination** — Synchronisierte Bewegung, keine Konflikte
- **Validierung per digitalem Zwilling** — In der Simulation testen, bevor physisch ausgerollt wird
- **Was-wäre-wenn-Szenarien** — Layout ändern, Auswirkungen sofort sehen

---

### Echtzeit-Zusammenarbeit

**Robotik ist Teamarbeit.**

Maschinenbauingenieure. Steuerungstechniker. Integrationsspezialisten. Endkunden. Oft weltweit verteilt.

Automerge-CRDT-Synchronisierung:

- Alle arbeiten gleichzeitig am selben Modell
- Keine Dateisperren. Keine Merge-Konflikte. Kein „Wer hat die neueste Version?“
- Änderungen sind in Echtzeit für das gesamte Team sichtbar
- Vollständige Versionshistorie für jede Entscheidung

---

### Offene Standards

**Keine Herstellerbindung. Volle Interoperabilität.**

STEP-Geometrie funktioniert mit jedem gängigen CAD- und Simulationssystem. Exportieren Sie in die Werkzeuge Ihres Roboter-OEMs. Importieren Sie die Konstruktionen Ihrer Kunden. Kein Konvertierungsverlust.

Ihr Arbeitszellen-Design ist nicht in einem proprietären Format gefangen. Es gehört Ihnen.

---

### Von der Simulation zur Realität

**Schließen Sie die Lücke.**

Verbinden Sie Ihren digitalen Zwilling mit der physischen Zelle:

- Echte Sensordaten in die Simulation streamen
- Modelle mit der tatsächlichen Leistung abgleichen
- Wartungsbedarf vor Ausfällen vorhersagen
- Kontinuierlicher Verbesserungskreislauf von der Produktion zurück ins Design

Für Luftrobotik und autonome Fahrzeuge siehe unsere [Drohnenflotten-Architektur]({{< relref "/fleet" >}}) – dieselbe NATS-JetStream-Infrastruktur treibt sowohl Boden- als auch Luftplattformen an.

---

## Anwendungsfälle

| Anwendung | Beschreibung |
|-------------|-------------|
| **Einsatz von Humanoiden** | Gestalten Sie die Umgebungen, in denen Humanoide arbeiten werden. KI-gestütztes Layout für die Mensch-Roboter-Kollaboration. Simulation vor dem Einsatz. |
| **Lagerautomatisierung** | Mobile Roboter, Ware-zur-Person-Systeme, automatisierte Lagerung. Konstruieren, simulieren und optimieren im großen Maßstab. |
| **Fertigungszellen** | Schweißen, Montage, Maschinenbeschickung, Palettierung. Komplettes Arbeitszellen-Design mit KI-Unterstützung. |
| **Chirurgische Robotik** | Präzisionsumgebungen für Präzisionsmaschinen. Räumliche Einschränkungen, sterile Zonen, Geräteintegration. |
| **Feldrobotik** | Landwirtschaft, Bauwesen, Bergbau. Konstruieren für unstrukturierte Umgebungen. |

---

## Für Roboterhersteller

**Sie bauen Roboter? Ihre Kunden müssen Umgebungen für sie gestalten.**

Geben Sie ihnen Werkzeuge, die:

- Ihre Robotermodelle nativ integrieren
- Mit präziser Kinematik simulieren
- In ihre bestehenden Systeme exportieren
- Echtzeit-Zusammenarbeit über Organisationsgrenzen hinweg ermöglichen

Werden Sie unser Partner, damit sich Ihre Roboter einfacher einsetzen lassen.

---

## Für Integratoren

**Sie bringen Roboter an Kundenstandorte?**

- Zellen schneller konstruieren mit KI-Unterstützung
- Validieren, bevor Sie bauen
- In Echtzeit mit Kunden zusammenarbeiten
- Saubere Ergebnisse in offenen Standards übergeben

Integrationszeit verkürzen. Margen steigern.

---

## Für Endanwender

**Sie holen die Automatisierung ins eigene Haus?**

- Visualisieren, wie die Automatisierung in Ihrer Anlage aussehen wird
- Sich direkt am Designprozess beteiligen
- Zielkonflikte verstehen, bevor Sie sich festlegen
- Ihre Konstruktionen in offenen Formaten besitzen

Ihre Anlage. Ihre Automatisierung. Ihre Daten.

---

## Den Kreis schließen mit Publish

Robotik besteht nicht nur aus 3D-Modellen – sondern auch aus Dokumentation.

- **Arbeitsanweisungen** — Montageanleitungen, die aus dem Modell generiert werden
- **Sicherheitsdokumentation** — Verfahren, Checklisten, Compliance-Formulare
- **Wartungshandbücher** — Gerätespezifikationen, Teilenummern, Wartungsintervalle
- **Schulungsmaterialien** — Bedienerhandbücher mit Visualisierungen aus der Simulation

Alles aus einer einzigen Quelle. Alles abgestimmt auf Ihr 3D-Arbeitszellen-Design.

[Publish entdecken →]({{< relref "/platform/publish" >}})

---

## Datensouveränität

Ihre Robotikdaten bleiben Ihre:

| Traditioneller Ansatz | Mit WellKnown |
|---------------------|----------------|
| Arbeitszellen-Designs in der Cloud des Anbieters | Designs auf Ihrer Infrastruktur |
| Roboter-OEM besitzt die Integrationsdaten | Sie besitzen Ihr Automatisierungs-IP |
| Plattformwechsel = Neuaufbau | Wechsel = Gateway neu verbinden |
| IP über Anbieter verstreut | Sie kontrollieren Ihren Wettbewerbsvorteil |

Export nach STEP. Veröffentlichen AUF Plattformen. Niemals darin GEFANGEN.

[Mehr über WellKnown erfahren →]({{< relref "/platform/publish#wellknown-gateway-data-sovereignty" >}})

---

## Jetzt starten

Die Robotik-Explosion braucht Infrastruktur. Wir bauen sie.

Ob Sie Humanoide konstruieren, Cobots einsetzen oder Lager automatisieren – lassen Sie uns sprechen.

[Kontakt aufnehmen →](/contact)
