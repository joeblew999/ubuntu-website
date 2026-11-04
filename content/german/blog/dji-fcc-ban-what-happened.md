---
title: "DJI steht jetzt auf der FCC Covered List – was wirklich passiert ist"
meta_title: "DJI-FCC-Verbot: Was passiert ist und was es bedeutet | Ubuntu Software"
description: "Die Frist im Dezember 2025 ist verstrichen. DJI steht nun auf der FCC Covered List. So kam es dazu, das bedeutet es, und so geht es weiter."
date: 2025-12-24T06:00:00Z
image: "/images/blog/dji-fcc-ban.svg"
categories: ["Industry", "Drones"]
author: "Gerard Webb"
tags: ["drones", "dji", "fcc", "regulation", "open-source"]
draft: false
---

Am 22. Dezember 2025, einen Tag vor Ablauf der Frist, gab die FCC bekannt, dass „unbemannte Luftfahrzeugsysteme (UAS) und kritische UAS-Komponenten, die in einem fremden Land hergestellt werden“ in die FCC Covered List aufgenommen werden.

**Das Verbot ist jetzt in Kraft.**

## Wie es dazu kam

Section 1709 des NDAA für das Haushaltsjahr 2025 schrieb eine bundesbehördliche Sicherheitsprüfung von DJI-Drohnen vor. Die Frist dafür war der 23. Dezember 2025. Ohne eine solche Prüfung sollte DJI automatisch in die FCC Covered List aufgenommen werden.

So lief es ab:

**DJI beantragte die Prüfung selbst.** Im März, Juni und Dezember 2025 bat das Unternehmen Bundesbehörden schriftlich darum, seine Produkte zu untersuchen.

**Keine Behörde antwortete.** Es wurde keine Untersuchung eingeleitet, es wurden keine Beweise erhoben und keine Feststellungen getroffen.

**Die Regierung umging das Verfahren.** Am 21. Dezember 2025 erließ die Exekutive eine „Feststellung zur nationalen Sicherheit“, statt die gesetzlich vorgeschriebene Untersuchung abzuschließen.

**Die Frist verstrich.** Am 22. Dezember machte die FCC das Verbot offiziell.

## Die „Falltür“-Konstruktion

Der Kongress hat Section 1709 bewusst mit zwei Falltüren versehen:

1. **Keine zuständige Behörde.** Das Gesetz legte nicht fest, welche Behörde die Prüfung durchführen sollte.
2. **Automatisches Verbot.** Fand keine Prüfung statt, trat das Verbot automatisch in Kraft.

Das Ergebnis stand von vornherein fest. Das Verfahren sollte den Anschein eines rechtsstaatlichen Verfahrens erwecken und zugleich das Ergebnis sicherstellen.

## Was das konkret bedeutet

### Für bestehende DJI-Besitzer

Ihre Drohnen funktionieren weiterhin. Sie dürfen sie weiterhin nach FAA Part 107 fliegen. Im Betrieb ändert sich nichts, vorerst.

### Für Neuanschaffungen

Neue DJI-Modelle und betroffene HF-Komponenten können die FCC-Zulassungen nicht erhalten, die für den legalen Import, die Vermarktung oder den Verkauf in den USA erforderlich sind. Damit sind praktisch blockiert:

- Neue Drohnenmodelle für den US-Markt
- Neue Firmware, die eine FCC-Zertifizierung erfordert
- Neues Zubehör mit HF-Komponenten

### Der Weg über Ausnahmen

Laut FCC können Ausnahmen gewährt werden, wenn das Verteidigungsministerium oder das Heimatschutzministerium feststellt, dass ein bestimmtes UAS oder eine bestimmte Komponente kein Risiko darstellt.

Ob tatsächlich Ausnahmen gewährt werden, bleibt abzuwarten.

## DJIs Reaktion

DJI bezeichnete das Vorgehen der FCC als „enttäuschend“ und erklärte:

> „Die Bedenken hinsichtlich der Datensicherheit von DJI beruhten nicht auf Beweisen, sondern spiegeln Protektionismus wider, der den Grundsätzen eines offenen Marktes widerspricht.“

DJI bestätigte, dass bestehende Produkte nicht betroffen sind und das Unternehmen weiterhin am US-Markt festhält.

## Das Ausmaß der Auswirkungen

Dies ist kein Nischenthema:

- **87 %** der Drohnen im Bereich öffentliche Sicherheit in den USA stammen von DJI
- **Über 80 %** der Behörden für öffentliche Sicherheit auf Bundesstaats- und Kommunalebene setzen bei der Waldbrandbekämpfung, bei Such- und Rettungseinsätzen und bei kritischen Einsätzen auf DJI

Der Zeitpunkt kurz vor Weihnachten, mitten in andauernden Waldbränden und Naturkatastrophen, zeigt deutlich, wie weit Politik und betriebliche Realität auseinanderliegen.

## Wie es weitergeht

### Kurzfristig (2025–2026)

- Bestehende DJI-Betreiber setzen den normalen Betrieb fort
- Keine neuen DJI-Modelle kommen auf den US-Markt
- Die Gebrauchtpreise für vorhandene DJI-Hardware steigen
- Firmware-Updates werden unsicher

### Mittelfristig (2026–2027)

- Die Lieferkette für DJI-Ersatzteile und -Zubehör wird knapper
- Flottenbetreiber müssen über Ersatzplattformen entscheiden
- Open-Source-Alternativen werden stärker genutzt
- Andere chinesische Hersteller geraten ähnlich unter Druck

### Langfristig

- Der US-Drohnenmarkt spaltet sich zwischen inländischen Herstellern und Open Source auf
- Behörden für öffentliche Sicherheit wechseln auf alternative Plattformen
- Unternehmen investieren in herstellerunabhängige Lösungen

## Die Open-Source-Alternative

Wir haben uns darauf vorbereitet. Unser [Leitfaden zur DJI-Migration]({{< relref "/fleet/dji-migration" >}}) beschreibt den Weg zu Open-Source-Plattformen wie PX4 und ArduPilot.

Die Argumente für Open Source sind heute stärker denn je:

| DJI (nach dem Verbot) | Open Source |
|----------------|-------------|
| Keine neuen Modelle | Kontinuierliche Innovation |
| Ungewisse Firmware | Von der Community gepflegt |
| Ein einziger Anbieter | Mehrere Hersteller |
| Regulatorisches Risiko | Auf Standards basierend |
| Cloud-Abhängigkeit | Self-Hosting möglich |

Wenn Sie eine DJI-Flotte betreiben, lautet die Frage nicht, ob Sie umsteigen, sondern wann und wie.

## Für öffentliche Sicherheit und Unternehmen

Am meisten steht für Behörden der öffentlichen Sicherheit und Unternehmensflotten auf dem Spiel. Sie brauchen:

- Bewährte Alternativen, die mit den Fähigkeiten von DJI mithalten
- Umstiegspläne, die den Betrieb nicht stören
- Schulung und Support für neue Plattformen
- Infrastruktur für das Flottenmanagement

Wir liefern die Infrastrukturebene, also die Software, mit der Open-Source-Drohnenflotten tatsächlich funktionieren. Unsere [Flottenarchitektur]({{< relref "/fleet" >}}) übernimmt die Koordination, Telemetrie und Verwaltung, die der Unternehmensbetrieb erfordert.

## Das Fazit

Das DJI-Verbot ist jetzt Realität und keine Spekulation mehr. Das Verfahren war fehlerhaft, aber das Ergebnis ist endgültig.

Betreiber haben zwei Möglichkeiten:

1. **Abwarten.** Mit vorhandener DJI-Hardware weitermachen, solange der Vorrat reicht, und auf Ausnahmen hoffen, die vielleicht nie kommen.

2. **Jetzt umsteigen.** Auf Open-Source-Plattformen wechseln, solange Sie noch planvoll vorgehen können, statt unter Zeitdruck improvisieren zu müssen.

Wir empfehlen die zweite Option.

---

*Brauchen Sie Hilfe bei der Planung Ihres Umstiegs? [Kontaktieren Sie uns →](/contact?subject=dji-migration)*

*Sie haben bereits eine DJI-Flotte? [Zum Migrationsleitfaden →]({{< relref "/fleet/dji-migration" >}})*
