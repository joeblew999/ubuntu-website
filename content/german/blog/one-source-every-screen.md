---
title: "Eine Quelle, jeder Bildschirm"
meta_title: "Eine Quelle, jeder Bildschirm | Ubuntu Software"
description: "Von Restaurant-Speisekarten bis zu Behördenterminals: Wie Single-Source-Publishing physische Displays neben Web und Print verändert."
date: 2024-11-27T05:00:00Z
image: "/images/blog/kiosk-system.svg"
categories: ["Publish", "Industry"]
author: "Gerard Webb"
tags: ["kiosk", "signage", "raspberry-pi", "single-source", "real-time"]
draft: false
---

Ein Restaurant aktualisiert seine Speisekarte. Was passiert als Nächstes?

Jemand bearbeitet die Website. Jemand anderes aktualisiert das PDF für die Speisekarten zum Mitnehmen. Eine dritte Person ändert das Kassensystem. Irgendwann denkt jemand an die digitalen Menütafeln und aktualisiert auch diese von Hand.

Vier Systeme. Vier Gelegenheiten für Fehler. Vier Stellen, an denen Preise auseinanderlaufen können.

**Das ist absurd. Und es ist überall so.**

## Das Problem der fragmentierten Anzeigen

Betreten Sie eine beliebige Organisation mit öffentlich sichtbaren Informationen:

**Restaurants**: Menütafeln, QR-Codes auf den Tischen, Website, Liefer-Apps, gedruckte Speisekarten
**Behörden**: Warteschlangenanzeigen, Wegweiser, Formularschalter, Website
**Gesundheitswesen**: Anzeigen im Wartezimmer, Check-in-Kioske, Patientenportale, gedruckte Materialien
**Einzelhandel**: Preisanzeigen, Werbebeschilderung, Lagerstatus, E-Commerce

Jede Anzeige läuft in der Regel auf einem eigenen System. Jede muss separat aktualisiert werden. Jede kann andere Informationen zeigen.

**Informationen laufen auseinander. Kunden sind verwirrt. Mitarbeiter verschwenden Zeit mit manuellem Abgleich.**

## Die Lösung: eine einzige Quelle

Was wäre, wenn jede Anzeige – physisch wie digital – aus derselben Quelle gespeist würde?

```
Your Database (single source of truth)
         ↓
    Publish Engine
         ↓
┌────────┼────────┐
↓        ↓        ↓
Website  PDF    Kiosk Display
         ↓        ↓
      Print    Menu Board
```

Ändern Sie den Preis einmal. Jede Anzeige wird aktualisiert.

Fügen Sie einen neuen Artikel einmal hinzu. Er erscheint überall.

Markieren Sie etwas einmal als ausverkauft. Kunden sehen es sofort – auf der Website, auf der Menütafel und auf der Speisekarte hinter dem gedruckten QR-Code.

## Echtzeitfähige, offlinefähige Kioske

Jetzt wird es interessant.

Herkömmliche Digital Signage erfordert eine ständige Verbindung. Fällt das Netzwerk aus? Die Anzeigen werden schwarz oder zeigen veraltete Daten.

Unser Ansatz nutzt Raspberry-Pi-Geräte, die dieselbe Offline-First-Architektur ausführen wie alles andere:

- **Lokale Datenreplik**: Jeder Kiosk hat eine vollständige Kopie der relevanten Daten
- **Automerge-Synchronisierung**: Aktualisierungen fließen automatisch, sobald eine Verbindung besteht
- **Graceful Degradation**: Netzwerkausfall? Die Anzeige arbeitet mit lokalen Daten weiter
- **Warteschlange für Änderungen**: Am Kiosk vorgenommene Aktualisierungen werden zurücksynchronisiert, sobald die Verbindung wiederhergestellt ist

**Eine Menütafel, die auch bei Internetausfällen funktioniert. Ein Check-in-Kiosk, der Eingaben offline zwischenspeichert.**

## Anwendungsfälle, die Sinn ergeben

### Restaurantsysteme

Eine Quelle speist:
- Digitale Menütafeln (Echtzeitpreise, ausverkaufte Artikel ausgegraut)
- QR-Code-Speisekarten (gleicher Inhalt, für Mobilgeräte optimiert)
- Speisekarten-Seiten auf der Website
- PDF-Speisekarten für den Druck
- Integration mit Küchendisplays

Ändern Sie den Preis eines Burgers. Überall wird aktualisiert. Die Menütafel an der Theke, das PDF, das jemand für das Catering ausdruckt, die Website – alles stimmt überein.

### Behördliche Servicezentren

Eine Quelle speist:
- Anzeigen für das Warteschlangenmanagement („Jetzt an der Reihe: B-47“)
- Kioske mit Serviceverzeichnis
- Formulare zum Herunterladen (PDF) oder Ausfüllen (Web/Kiosk)
- Wegweiser-Beschilderung
- Serviceübersichten auf der Website

Aktualisieren Sie die Öffnungszeiten. Jede Anzeige spiegelt die Änderung wider. Der Wegweiser-Kiosk, die Website, die gedruckten Handzettel – alles korrekt.

### Wartezimmer im Gesundheitswesen

Eine Quelle speist:
- Check-in-Kioske (Formulare ausfüllen, in die Warteschlange einreihen)
- Wartezeitanzeigen
- Bildschirme mit Patienteninformationen
- Gedruckte Aufnahmeformulare (gleiche Felder, einheitliches Layout)

Patient füllt das Formular am Kiosk aus → Daten fließen in Ihre elektronische Patientenakte → keine manuelle Abschrift.

### Einzelhandelsumgebungen

Eine Quelle speist:
- Preisanzeigen (elektronische Regaletiketten)
- Werbebeschilderung
- Bildschirme mit Lagerstatus
- Produktübersichten auf der Website
- Gedruckte Kataloge

Markieren Sie einen Artikel im Warenwirtschaftssystem als „Abverkauf“. Die Preisanzeigen werden aktualisiert. Die Website wird aktualisiert. Der gedruckte Wochenprospekt zeigt die Abverkaufspreise.

## Die Hardware

Warum Raspberry Pi?

**Kosten**: Unter 100 $ pro Anzeigeknoten, einschließlich Pi und einfachem Gehäuse
**Zuverlässigkeit**: Keine beweglichen Teile, läuft unter Linux, läuft jahrelang
**Flexibilität**: Der HDMI-Ausgang steuert Anzeigen jeder Größe an
**Offlinefähig**: Lokale Rechenleistung und Speicher für echten Edge-Betrieb
**Aktualisierbar**: Fernverwaltung für Software-Updates

Ein Restaurant kann Menütafeln zu einem Bruchteil der Kosten herkömmlicher Digital Signage einsetzen. Eine Behörde kann Kioske ohne teure Signage-Verträge für Unternehmen hinzufügen.

## Die technische Architektur

Auf jedem Kiosk laufen:

| Komponente | Zweck |
|-----------|---------|
| **Lokale Datenbank** | SQLite-Replik der relevanten Daten |
| **Sync-Engine** | Automerge-CRDT für konfliktfreie Aktualisierungen |
| **Display-Renderer** | Webtechnologien (HTML/CSS/JS) für flexible Layouts |
| **Eingabe-Handler** | Touch, Tastatur, Barcodescanner, Kartenleser |
| **Offline-Warteschlange** | Eingaben werden bis zur Synchronisierung lokal gespeichert |

Dieselbe Codebasis, die Ihre Webformulare ausführt, läuft auch auf dem Kiosk. Gleiche Validierung. Gleiche Feldlayouts. Gleiches Datenziel.

## Mehr als Anzeigen: auch Eingabegeräte

Kioske sind nicht nur Ausgabegeräte. Sie sind auch Eingabepunkte.

**Formulare ausfüllen**: Kunden füllen Formulare auf dem Touchscreen des Kiosks aus
**Dokumente scannen**: Der Kiosk scannt Papierdokumente, OCR extrahiert die Daten
**Zahlungsabwicklung**: Integrierte Kartenleser
**Identitätsprüfung**: Kamera zum Scannen von Ausweisdokumenten
**Unterschriftserfassung**: Touch-Unterschriftenpads

Alle erfassten Daten fließen durch dasselbe System wie Web-Eingaben.

Ein Behördenformular, eingereicht über:
- Webformular auf der Website → Datenbank
- PDF ausgedruckt, ausgefüllt, am Kiosk gescannt → Datenbank
- Touchscreen-Formular am Kiosk → Datenbank

**Gleiche Daten. Gleiches Ziel. Unterschiedliche Eingabemethode.**

## Der Vorteil der Einheitlichkeit

Das bedeutet „einheitlich“ wirklich:

Ein Kunde sieht ein Formular auf Ihrer Website. Derselbe Kunde besucht Ihr Büro und sieht dasselbe Formular an einem Kiosk. Gleiche Felder, gleiches Layout, gleiche Fragen.

Er beginnt am Kiosk mit dem Ausfüllen, hat keine Zeit mehr und scannt einen QR-Code, um auf seinem Smartphone weiterzumachen. Gleiches Formular, gleicher Fortschritt (synchronisiert über seine Sitzung).

Er schließt es zu Hause ab und klickt auf „Absenden“. Die Daten fließen in Ihre Datenbank. Kein erneutes Abtippen, keine Übertragungsfehler, keine abweichenden Feldnamen.

**Das ermöglicht Single-Source-Publishing, wenn es auf physische Anzeigen ausgeweitet wird.**

## Der Business Case

**Vor Single-Source-Anzeigen:**
- Arbeitszeit für die Aktualisierung mehrerer Systeme
- Fehler durch inkonsistente Informationen
- Verwirrte Kunden („Auf dem Schild steht ein Preis, auf der Website ein anderer“)
- Hohe Kosten herkömmlicher Digital-Signage-Systeme
- Anzeigen, die ausfallen, wenn das Internet weg ist

**Nach Single-Source-Anzeigen:**
- Einmal aktualisieren, überall verbreiten
- Garantierte Konsistenz an allen Kontaktpunkten
- Weniger Kundenbeschwerden
- Standard-Hardware zu Standardpreisen
- Anzeigen, die offline funktionieren

**Die Rechnung ist einfach: weniger Fehler, weniger Arbeitszeit, bessere Kundenerfahrung.**

## Erste Schritte

Wir bauen das gerade. Dieselbe Publish-Engine, die Webseiten, PDFs und Formulare erzeugt, wird auch Kiosk-Anzeigen erzeugen.

Gleiche Markdown-Quelle. Gleiche DSL für Felder. Gleiche Datenbankverbindung.

Die Anzeige ist einfach ein weiteres Ausgabeformat.

---

*Interesse an frühem Zugang zu den Kiosk-Funktionen? [Kontaktieren Sie uns →](/contact)*

---

*Teil unserer Publish-Plattform. [Mehr über Single-Source-Publishing erfahren →](/platform/publish)*
