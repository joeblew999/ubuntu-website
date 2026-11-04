---
title: "Veröffentlichen"
meta_title: "Content-Management-Plattform | Ubuntu Software"
description: "Single-Source-CMS für Multi-Channel-Publishing. Eine DSL für Texte, Grafiken, Formulare und E-Mails. Automatisch übersetzte Ausgaben. Perfektes Branding durch die Architektur."
image: "/images/publish.svg"
draft: false
---

## Eine Sprache. Jede Ausgabe. Jede Sprache.

Ein Single-Source-Content-Management-System. Definieren Sie alles in einer DSL – **Text (1D), der alle 2D- und 3D-Ausgaben steuert**. Setzen Sie Komponenten zusammen, die wiederum Komponenten enthalten. Ausgabe auf jedem Kanal. Automatische Übersetzung in jede Sprache. Daten zurückerfassen.

Alles webfähig. Alles weltweit zugänglich. Perfektes Branding – nicht durch Disziplin, sondern durch Architektur.

---

## Das Problem

Organisationen veröffentlichen Inhalte über mehrere Kanäle: Websites, PDFs, E-Mails, Formulare, Beschilderung.

Jeder Kanal wird separat gepflegt. Inhalte driften auseinander. Das Branding zerfällt. Übersetzungen hinken hinterher. Formulareinsendungen werden manuell neu erfasst.

**Eine Änderung bedeutet, fünf Systeme zu aktualisieren. Oder sie wird gar nicht aktualisiert.**

---

## Eine DSL

Alles wird in derselben Sprache definiert.

**Textinhalte:**
```
# Welcome to Our Service

We help organizations {industry} achieve {outcome}.
```

**Grafiken:**
```
@logo: svg {
  viewBox: "0 0 200 50"
  rect: { x: 0, y: 0, width: 50, height: 50, fill: "#2563eb" }
  text: { x: 60, y: 35, content: "Ubuntu", font: "bold 24px" }
}
```

**Formularfelder:**
```
Name: [_______________]{field: name, required: true}
Email: [_______________]{field: email, type: email}
Date: [__/__/____]{field: date, type: date}
```

**E-Mails:**
```
@welcome-email: email {
  to: {customer.email}
  subject: "Welcome, {customer.name}"
  body: include welcome-content
}
```

**Eine Sprache. Texte, Grafiken, Formulare, E-Mails, Dokumente.**

---

## Zusammensetzen

Komponenten enthalten Komponenten.

```
@header: compose {
  include: logo
  include: navigation
  include: search-bar
}

@page-template: compose {
  include: header
  include: content
  include: footer
}

@welcome-letter: compose {
  include: page-template
  content: "Dear {customer.name}, welcome to..."
}

@application-form: compose {
  include: page-template
  content: include form-fields
}
```

**Die Kaskade:**
- Logo ändern → Header wird aktualisiert
- Header wird aktualisiert → jede Seite wird aktualisiert
- Jede Seite wird aktualisiert → jedes PDF, jede E-Mail, jedes Formular wird aktualisiert

**Deshalb bleibt das Branding perfekt.** Nicht, weil jemand daran denkt, alles zu aktualisieren. Sondern weil die Architektur es unmöglich macht, es nicht zu tun.

Verlinken Sie beliebige vorhandene Inhalte. Wiederverwendung im gesamten System. Eine einzige Quelle der Wahrheit.

---

## Ausgabe

Aus einer Quelle wird alles generiert:

| Ausgabe | Was passiert |
|--------|--------------|
| Webseiten | Als HTML gerendert, SEO-optimiert |
| PDFs | Druckfertig, archivierbar |
| E-Mails | Direkt an Empfänger gesendet |
| Webformulare | Interaktiv, Daten werden zurückerfasst |
| PDF-Formulare | Ausfüllbar oder druckbar, Rückerfassung per OCR |
| SVG-Grafiken | Vektor-Assets in beliebiger Größe |
| Karten | Ihre Standortdaten visualisiert |
| Kiosks | Physische Displays in Echtzeit |

**Jede Ausgabe wird automatisch übersetzt.**

Schreiben Sie auf Englisch. Deutsche Kollegen sehen Deutsch. Spanische Partner sehen Spanisch. Japanische Kunden sehen Japanisch.

Nicht nachträglich übersetzt. Übersetzt als Teil des Renderings. Jede Ausgabe. Jede Sprache. Automatisch.

---

## Übersetzen

Übersetzung ist kein Feature. Sie ist die Art, wie das System funktioniert.

| Fähigkeit | Was das bedeutet |
|------------|---------------|
| Echtzeit im Editor | Gleichzeitige Zusammenarbeit über Sprachgrenzen hinweg |
| Automatisch übersetzte Ausgaben | Jedes Format wird in jeder Sprache gerendert |
| Offline-KI | Funktioniert ohne Internet |
| Kontextbezogen | Unterscheidet zwischen Bildungs-, Rechts- und medizinischer Terminologie |
| Bidirektional | Andere bearbeiten in ihrer Sprache, Sie sehen Ihre |

**Einmal schreiben. In jeder Sprache veröffentlichen.**

---

## Erfassen

Formulare erfassen Daten zurück. Da wir das Formular generiert haben, wissen wir genau, wo sich jedes Feld befindet.

| Kanal | Was passiert |
|---------|--------------|
| Webformular abgeschickt | Daten fließen in Ihre Datenbank |
| PDF-Formular digital ausgefüllt | Daten werden in Ihre Datenbank extrahiert |
| PDF-Formular gedruckt, von Hand ausgefüllt, gescannt | OCR erfasst die Daten in Ihre Datenbank |

Kein Training der OCR. Kein Feld-Mapping. Wir haben das Formular generiert – wir wissen, wo alles ist.

**Papier oder digital. Dieselben Daten. Dasselbe Ziel.**

---

## Human-in-Loop: Außendiensteinsätze in Echtzeit

Publish ist nicht nur für statische Dokumente gedacht. Es ist die Art, wie Außendienstmitarbeiter in Echtzeit mit Ihren Systemen interagieren.

### Das voll ausgeschöpfte Gerät

Das Gerät jedes Außendienstmitarbeiters wird zugleich zu Sensor und Anzeige:

| Eingabe (aus dem Außendienst) | Ausgabe (an den Außendienst) |
|--------------------|-------------------|
| GPS-Position | Routenkarten, Turn-by-Turn-Navigation |
| Sprachbefehle (SST) | Aufgabenbenachrichtigungen |
| Touch-Interaktionen | Aktualisierungen des Zeitplans |
| Formulareinsendungen | Bestätigungsmeldungen |
| Kamera (Fotos/Video) | Statusanzeigen |

Das System weiß, wo sich jeder Mitarbeiter befindet. Mitarbeiter können antworten – per Sprache, Touch oder Formular.

### 2D-Ausgaben

Publish rendert alle 2D-Inhalte – Web, Dokumente, Karten, Displays:

| Ausgabe | Anwendungsfall | Beispiel |
|--------|----------|---------|
| **Webseiten** | Kundenportale, Statusverfolgung | Servicestatus, Bestellverfolgung |
| **Routenkarten** | Navigation, Lieferverfolgung | Müllabfuhr, Lieferstopps |
| **Dashboards** | Statusübersicht, KPIs | Dispositionskonsole, Teamleistung |
| **Dokumente** | Arbeitsaufträge, Berichte, Formulare | PDF-Manifeste, Compliance-Exporte |
| **Kiosks** | Öffentliche Displays, Team-Boards | Depotstatus, Speisekarten |

Für **räumliche 3D-Visualisierung** – Baustellen, Lagerhallen, digitale Zwillinge – siehe die [Spatial-Plattform](/platform/spatial/).

### Sprachgesteuerte Aktualisierungen

Speech-to-Text (SST) ermöglicht eine freihändige Interaktion mit dem System:

| Mitarbeiter sagt | System reagiert |
|-------------|-------------|
| „123 Main wegen Baustelle gesperrt“ | Markiert den Standort, plant Routen neu |
| „Nächste Aufgabe überspringen, Kunde nicht erreichbar“ | Entfernt sie aus der Warteschlange, passt den Zeitplan an |
| „Gerät muss repariert werden“ | Erstellt ein Wartungsticket |
| „Aufgabe erledigt“ | Erfasst den Abschluss, startet den nächsten Workflow |

Keine Anrufe bei der Disposition. Kein Funkverkehr. Der Mitarbeiter spricht, die KI hört zu, das System passt sich an.

### Entscheidungshierarchie

Menschliches Urteilsvermögen hat Vorrang vor KI-Empfehlungen:

```
┌─────────────────────────────────────────────────────────────────┐
│                     DECISION AUTHORITY                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1. FIELD WORKER       ───────────────────────▶  HIGHEST       │
│      Eyes on ground, can override anything                       │
│                                                                  │
│   2. DISPATCHER         ───────────────────────▶  HIGH          │
│      Fleet-wide view, can reassign work                          │
│                                                                  │
│   3. AI OPTIMIZER       ───────────────────────▶  MEDIUM        │
│      Suggests routes/tasks, adapts to changes                    │
│                                                                  │
│   4. SCHEDULED PLAN     ───────────────────────▶  LOWEST        │
│      Starting point, always subject to reality                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Die KI ist leistungsstark – aber der Mensch vor Ort hat das letzte Wort.**

### Präsenz & Benachrichtigungen

Menschen sind nicht immer verfügbar. Das System weiß, wer es ist:

| Fähigkeit | Was sie bewirkt |
|------------|--------------|
| **Präsenzerkennung** | Wer ist online, wer offline, wer beschäftigt |
| **Rollenbewusstsein** | Wer genehmigen darf, wer übersteuern darf, wer Bereitschaft hat |
| **Intelligentes Routing** | Benachrichtigungen gehen zuerst an verfügbares Personal |
| **Eskalation** | Reagiert der Hauptverantwortliche nicht, wird an die Vertretung eskaliert |
| **Bestätigung** | Nachverfolgen, wer was gesehen und wann reagiert hat |

**Benachrichtigungskanäle:**

- Push-Benachrichtigungen (Mobil/Desktop)
- SMS für kritische Warnmeldungen
- E-Mail für asynchrone Updates
- In-App-Hinweise mit Ton/Haptik
- Kiosk-Displays für Sichtbarkeit im Team

Das System passt sich der Verfügbarkeit der Menschen an. Dringende Anliegen finden die richtige Person. Nicht dringende Anliegen warten in der Warteschlange, bis jemand bereit ist.

### Offline-Resilienz

Außendienstarbeit findet dort statt, wo Netze versagen. Das Gerät arbeitet weiter:

- **Verbindung verloren** → App arbeitet mit lokalem SQLite weiter
- **Änderungen werden lokal gepuffert** → Synchronisierung bei erneuter Verbindung
- **Kein Datenverlust** → Automerge CRDT löst Konflikte
- **Präsenz-Updates** → Als offline markiert, Benachrichtigungen werden gepuffert

Außendiensteinsätze dürfen wegen Netzwerkproblemen nicht stillstehen. Diese Architektur stellt sicher, dass sie es nicht tun.

[Außendiensteinsätze in Aktion sehen →](/applications/logistics/)

---

## Verbinden

Ihre Datenbank. Ihr Schema.

- Daten lesen, um Dokumente zu befüllen: „Sehr geehrte(r) {customer.name}...“
- Erfasste Daten zurückschreiben: Formulareinsendungen → Ihre Tabellen
- Jede SQL-Datenbank: PostgreSQL, MySQL, SQLite
- Ihre Infrastruktur: selbst gehostet oder in der Cloud

**Wir ersetzen Ihr Backend nicht. Wir verbinden uns damit.**

---

## Branchen

| Branche | Warum Publish |
|----------|-------------|
| [**Öffentliche Verwaltung**](/applications/government) | Jeden Bürger erreichen. Papier + digital. Barrierefrei. Regelkonform. |
| [**Gesundheitswesen**](/applications/healthcare) | Weniger Verwaltung. Weniger Fehler. Bessere Versorgung. |
| [**Finanzwesen**](/applications/financial) | Compliance ohne Komplexität. Lückenlose Audit-Trails. |
| [**Bildung**](/applications/education) | Moderne Verwaltung. Traditionelle Optionen. |
| [**Versicherungen**](/applications/insurance) | Jeder Kanal. Ein System. Außendienst und Büro vereint. |

---

## Architektur

| Schicht | Was sie tut |
|-------|--------------|
| DSL-Parser | Versteht Text, Grafiken, Formulare, Komposition |
| Include-Resolver | Rekursive Komposition von Komponenten |
| Renderer | Web-, PDF-, E-Mail-, SVG- und Formularausgaben |
| Übersetzungs-Engine | Echtzeit-KI, offlinefähig |
| Feld-Mapper | Kennt die exakte Position jedes Formularfelds |
| OCR-Engine | Präzise Extraktion anhand der Feldzuordnung |
| Datenbank-Konnektoren | Lesen aus und Schreiben in Ihr Schema |

**Keine proprietären Formate. Alles exportieren. Ihre Inhalte gehören Ihnen.**

---

## Nativ auf jeder Plattform

| Plattform | Erlebnis |
|----------|------------|
| **Linux** | Nativer Desktop, eingebettete Kiosks, Edge-Geräte |
| **Windows** | Native Windows-Anwendung |
| **macOS** | Native Mac-Anwendung |
| **iOS** | Native Mobil-App für Außendienstmitarbeiter |
| **Android** | Native Mobil-App für Außendienstmitarbeiter |
| **Web** | Moderner Browser, voller Funktionsumfang |

Eine Codebasis. Native Performance auf jeder Plattform. Autoren nutzen Desktops. Außendienstmitarbeiter nutzen Tablets. Kunden nutzen Smartphones. Dasselbe System. Dieselben Daten.

---

## Mehr als Dokumente

Derselbe Single-Source-Ansatz funktioniert für:

**Karten** — Ihre Geodaten, Ihr Styling, Integration mit Google/Apple Maps für mehr Reichweite.

**Video** *(demnächst)* — Ihr Videoserver, Spiegelung auf YouTube/Vimeo für bessere Auffindbarkeit.

**Kalender** *(demnächst)* — Ihr CalDAV-Server, Synchronisierung mit Google/Apple Kalender für mehr Komfort.

**Kiosks** *(demnächst)* — Raspberry-Pi-Displays, Speisekarten, Behörden, Beschilderung im Einzelhandel.

Alle Links laufen über [Ihr Gateway](#wellknown-gateway-data-sovereignty). Veröffentlichen Sie AUF Big-Tech-Plattformen. Seien Sie nie an sie GEBUNDEN.

---

## WellKnown Gateway: Datensouveränität

Wohin gehen Ihre veröffentlichten Inhalte? Auf externe Plattformen? Genau hier kommt das **WellKnown Gateway** ins Spiel – es stellt sicher, dass Ihre Daten Ihnen gehören, selbst wenn Sie AUF anderen Systemen veröffentlichen.

### Das Problem mit dem Plattform-Lock-in

Herkömmlicher Ansatz:
- Kundendaten liegen in Salesforce
- Dokumente liegen in Google Drive
- Formulare liegen in Typeform
- Sie sind von der API jedes einzelnen Anbieters abhängig
- Ändert ein Anbieter seine Bedingungen, erhöht er die Preise oder stellt er den Betrieb ein – müssen Sie hektisch reagieren

### Die WellKnown-Lösung

```
┌─────────────────────────────────────────────────────────────────┐
│                    YOUR DATA (Source of Truth)                   │
│                                                                  │
│   SQLite on every device + NATS JetStream persistence            │
│   You own it. It's on your infrastructure.                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    WELLKNOWN GATEWAY                             │
│                                                                  │
│   Publishes TO external platforms                                │
│   Syncs FROM external platforms                                  │
│   You control what goes where                                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
   │   Google    │     │   Social    │     │  Industry   │
   │   Drive     │     │   Media     │     │  Platforms  │
   └─────────────┘     └─────────────┘     └─────────────┘
```

**Grundprinzip:** Sie veröffentlichen AUF Plattformen. Nicht die Plattformen besitzen Ihre Daten – sondern Sie.

### Souveränität über Kundenbeziehungen

Ihre Kunden interagieren mit Ihrer Marke, nicht mit der eines Anbieters:

| Ohne WellKnown | Mit WellKnown |
|-------------------|----------------|
| Kunde meldet sich im Anbieterportal an | Kunde meldet sich in IHREM Portal an |
| Anbieter besitzt die Kundenbeziehung | Sie besitzen die Kundenbeziehung |
| Anbieter hat Ihre Kundenliste | Sie kontrollieren die Kundendaten |
| Anbieterwechsel = Neuaufbau | Anbieterwechsel = Gateway neu verbinden |

Der Kunde glaubt, mit Ihnen zu sprechen. Das tut er auch. Die Technik dahinter ist unsichtbar.

### Basierend auf Internetstandards

WellKnown nutzt dieselben Standards, auf denen E-Mail, Kalender und föderierte soziale Netzwerke basieren:

| Standard | Zweck | Beispiel |
|----------|---------|---------|
| **RFC 8615** | `/.well-known/` URI-Discovery | `yourdomain.com/.well-known/...` |
| **WebFinger** | Konto-/Ressourcen-Discovery | Nutzer domänenübergreifend finden |
| **OAuth Discovery** | Autorisierungs-Endpunkte | Sicherer Zugriff durch Dritte |
| **CalDAV/CardDAV** | Kalender-/Kontaktsynchronisierung | Funktioniert mit jeder konformen App |
| **ActivityPub** | Föderationsprotokoll | Interoperabilität mit jedem System |

**Ihre Domain ist Ihre Identität.** Genau wie bei E-Mail (`you@yourdomain.com`) liegen Ihre Daten unter URIs, die Sie kontrollieren. Standardbasierte Föderation bedeutet, dass Sie sich mit jedem konformen System verbinden können – nicht nur mit unserem.

So funktioniert das offene Web. Wir wenden es lediglich auf das Publishing an.

### Datenportabilität

Alles wird in Standardformate synchronisiert, die Sie kontrollieren:

- **Dokumente** → Ihre SQLite-Datenbank
- **Formulare** → Ihre Datenbanktabellen
- **Medien** → Ihr Speicher
- **Audit-Logs** → Ihre Compliance-Aufzeichnungen

Jederzeit exportieren. Jederzeit migrieren. Sie sind nie gefangen.

---

## Basierend auf Foundation

Publish erbt automatisch alle Fähigkeiten von [Foundation](/platform/foundation/):

| Fähigkeit | Was das bedeutet |
|------------|---------------|
| **Offline-first** | Ohne Internet arbeiten, bei Verbindung synchronisieren |
| **Universelles Deployment** | Web, Desktop, Mobil – eine Codebasis |
| **Selbstbestimmt** | Ihre Server, Ihre Daten, Ihre Regeln |
| **Echtzeit-Synchronisierung** | Mehrere Bearbeiter, automatische Konfliktlösung |

[Mehr über Foundation erfahren →](/platform/foundation/)

---

## Jetzt starten

Wenn Sie Inhalte an mehreren Orten pflegen, zusehen, wie Ihr Branding auseinanderdriftet, auf Übersetzungen warten oder Einsendungen manuell verarbeiten – es gibt einen besseren Weg.

[Kontakt aufnehmen →](/contact/)

---

## Teil von etwas Größerem

Publish ist die **2D-Schicht** eines einheitlichen dimensionalen Stacks – vollständig webfähig, weltweit zugänglich.

### Der 1D → 2D → 3D-Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                    SINGLE SOURCE OF TRUTH                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1D: TEXT (Markdown/DSL)                                        │
│       └── Content definition, structured data, semantic markup   │
│                          │                                       │
│                          ▼                                       │
│   2D: PUBLISH                                                    │
│       └── Web pages, documents, forms, maps, dashboards, kiosks  │
│                          │                                       │
│                          ▼                                       │
│   3D: SPATIAL                                                    │
│       └── CAD, digital twins, robotics, simulation               │
│                                                                  │
│   ALL LAYERS: Web-capable, globally accessible                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Schicht | Plattform | Wofür sie zuständig ist |
|-------|----------|-----------------|
| **1D** | DSL/Markdown | Textinhalte, strukturierte Definitionen, die Quelle |
| **2D** | Publish | Web, Dokumente, Formulare, Karten, Kiosks, E-Mails |
| **3D** | [Spatial](/platform/spatial/) | CAD, digitale Zwillinge, Robotik, Simulation |
| **Gemeinsam** | [Foundation](/platform/foundation/) | Offline-Synchronisierung, CRDT, NATS, Deployment |

**Alles ist webfähig.** Der 1D-Text steuert die 2D-Ausgaben. Die 2D-Ausgaben können 3D-Ansichten einbetten. Die 3D-Modelle können 2D-Dokumentation generieren. Alles über den Browser zugänglich, alles synchronisiert, alles aus derselben Quelle der Wahrheit.

### Wie sie zusammenarbeiten

Für Organisationen, die in 3D arbeiten – Robotik, Fertigung, Bauwesen –, liefert Publish die Dokumentationsschicht:

- **Technische Zeichnungen**, generiert aus 3D-Modellen
- **Stücklisten (BOMs)**, extrahiert aus Baugruppen
- **Arbeitsanweisungen** mit Schritt-für-Schritt-Visualisierungen
- **Compliance-Formulare**, verknüpft mit CAD-Revisionen
- **Vor-Ort-Inspektionen** mit 2D-Checklisten + 3D-Kontext

Dieselbe Quelle der Wahrheit. 2D- und 3D-Ausgaben bleiben synchronisiert. Alles über das Web zugänglich.

[Spatial-Plattform entdecken →](/platform/spatial/)
