---
title: "BVLOS-Vorschriften"
meta_title: "Drohnenvorschriften für Flüge außerhalb der Sichtweite | Ubuntu Software"
description: "Weltweite BVLOS-Drohnenvorschriften (Beyond Visual Line of Sight – außerhalb der Sichtweite) nach Ländern. Regulatorische Anforderungen für den kommerziellen Flottenbetrieb verstehen."
image: "/images/robotics.svg"
draft: false
---

## Warum BVLOS für den Flottenbetrieb wichtig ist

Vorschriften für Flüge außerhalb der Sichtweite (Beyond Visual Line of Sight, BVLOS) bestimmen grundlegend, was kommerzielle Drohnenflotten leisten können. Ein Verständnis der regulatorischen Landschaft ist für das Design der Flottenplattform und die Marktstrategie unerlässlich.

| Herausforderung | Auswirkung auf die Flottenplattform |
|-----------|-------------------------|
| **Betriebliche Grenzen** | VLOS-Beschränkungen (400-500 m Reichweite) schränken den Nutzen der Flotte ein. BVLOS ermöglicht Infrastrukturinspektion, Lieferungen und Landwirtschaft im großen Maßstab. |
| **Compliance-Integration** | Die Flottensoftware muss regulatorische Grenzen durchsetzen – Remote ID, Detect-and-Avoid und die Anforderungen an Betriebsgenehmigungen bestimmen die Telemetrie- und Sicherheitsfunktionen. |
| **Markt-Timing** | Die Länder befinden sich in unterschiedlichen Phasen der BVLOS-Freigabe. Zu wissen, wo sich die Vorschriften öffnen, hilft bei der Priorisierung der Zielmärkte. |
| **Partnerbereitschaft** | Flottenbetreiber benötigen BVLOS-Genehmigungen für echte Einsätze. Wer nachverfolgt, wer über Genehmigungen verfügt, findet potenzielle Kunden und Partner. |

---

## Globale Trends

Bei allen großen Luftfahrtbehörden zeichnen sich vier Trends ab:

| Trend | Beschreibung |
|-------|-------------|
| **Einführung von SORA** | FAA, EASA, CASA und Transport Canada nutzen die Methodik des Specific Operations Risk Assessment [[1]](#ref-1) |
| **Remote-ID-Pflicht** | Wird in allen Rechtsordnungen verpflichtend |
| **Ausnahmegenehmigungen → Regeln** | Übergang von Ausnahmegenehmigungen im Einzelfall zu standardisierten, leistungsbasierten Regeln |
| **UTM-Integration** | Länder entwickeln UAS Traffic Management für skalierbaren BVLOS-Betrieb |

---

## Amerika

### Vereinigte Staaten (FAA)

**Status:** Auf Ausnahmegenehmigungen (Waivers) basierend, Übergang zu einem regelbasierten Ansatz (Part 108) [[2]](#ref-2)

**Zeitplan** [[3]](#ref-3):

- 190 BVLOS-Ausnahmegenehmigungen bis Okt. 2024 erteilt (gegenüber 6 im Jahr 2020)
- Regelentwurf für Part 108 für Jan. 2026 erwartet
- Der FAA Reauthorization Act 2024 hat das BEYOND-Programm bis 2029 verlängert [[4]](#ref-4)

**Betreiber mit BVLOS-Genehmigungen:**

| Betreiber | Anwendungsfall | Genehmigungsdatum | Ref. |
|----------|----------|---------------|-----|
| DroneUp | Medizinische Lieferungen | Jan. 2024 | [[5]](#ref-5) |
| Amazon Prime Air | Lieferungen an Endkunden | Mai 2024 | [[6]](#ref-6) |
| Asylon | Sicherheit/Überwachung | Apr. 2024 | [[7]](#ref-7) |
| Mehrere (Raum DFW) | Erster Luftraum für mehrere Betreiber | Juli 2024 | [[8]](#ref-8) |

---

### Kanada (Transport Canada)

**Status:** SFOC erforderlich, Übergang zu routinemäßigem BVLOS-Betrieb [[9]](#ref-9)

**Rahmenwerk:**

- Apr. 2025: Routinemäßiger BVLOS-Betrieb ohne SFOC unter Bedingungen mit geringem Risiko erlaubt (≤150 kg, unkontrollierter Luftraum, dünn besiedeltes Gebiet)
- Erfordert ein Advanced Pilot Certificate
- Nationales UTM-Rahmenwerk über Nav Canada bis 2030

**Betreiber:** Noch offen

---

## Europa

### Europäische Union (EASA)

**Status:** Kategorie „speziell“ mit SORA-/PDRA-Verfahren [[10]](#ref-10)

**Rahmenwerk** [[1]](#ref-1):

| Aspekt | Details |
|--------|---------|
| Standardszenario | STS-02: BVLOS mit Luftraumbeobachtern |
| Remote ID | Verpflichtend seit Jan. 2024 [[11]](#ref-11) |
| Genehmigungsdauer | 4-9 Monate (medizinisch), 6-12 Monate (Endkunden) [[12]](#ref-12) |
| Gebühren | 2.000-15.000 € |

**Betreiber:** Verschiedene Betreiber über nationale Sandboxes und beschleunigte Programme der Mitgliedstaaten.

---

### Vereinigtes Königreich (CAA)

**Status:** Betriebsgenehmigung (Operational Authorization, OA) über UK SORA

**Rahmenwerk:**

- BVLOS erfordert eine Betriebsgenehmigung (Operational Authorization) der CAA
- UK-SORA-Verfahren (seit Apr. 2025)
- Detect-and-Avoid-Technologie in der Regel erforderlich (ADS-B, FLARM, Bodenradar)

**Betreiber:**

| Betreiber | Erfolg | Ref. |
|----------|-------------|-----|
| sees.ai | Erster routinemäßiger BVLOS-Betrieb im Vereinigten Königreich (2021), Inspektionen für National Grid (2023) | [[13]](#ref-13) |

---

### Irland (IAA)

**Status:** Verfahren nach PDRA G-03

**Rahmenwerk:** Öffnet die Tür zu vereinfachten EU-Genehmigungen über die EASA.

**Betreiber:**

| Betreiber | Erfolg | Ref. |
|----------|-------------|-----|
| sees.ai | Erste BVLOS-Genehmigung in Irland (2025) | [[13]](#ref-13) |

---

## Asien-Pazifik

### Australien (CASA)

**Status:** SORA-basierte Bewertung gemäß Part 101 [[14]](#ref-14)

**Rahmenwerk** [[15]](#ref-15):

- Einführung von SORA 2.5 mit quantitativer Risikomodellierung
- Rahmenwerk zur Selbstbewertung für Broad Area BVLOS verfügbar
- Vereinfachtes Verfahren zieht globale Unternehmen an [[16]](#ref-16)

**Betreiber:**

| Betreiber | Erfolg | Ref. |
|----------|-------------|-----|
| Sphere Drones | Broad-Area-BVLOS-Genehmigung (Mai 2024) | [[17]](#ref-17) |

---

### Japan (MLIT/JCAB)

**Status:** Level-4-BVLOS über besiedelten Gebieten erlaubt [[18]](#ref-18)

**Rahmenwerk** [[19]](#ref-19):

| Level | Beschreibung | Erforderliche Lizenz |
|-------|-------------|------------------|
| Level 1-2 | VLOS-Betrieb | Zweite Klasse |
| Level 3 | BVLOS über unbewohnten Gebieten | Zweite Klasse |
| Level 4 | BVLOS über städtischen/besiedelten Gebieten | Erste Klasse |

- Lizenzen sind 3 Jahre gültig und erfordern eine registrierte Drohnenschule + Flugprüfung
- 422.879 registrierte UAVs (Stand Okt. 2024)

**Betreiber:**

- Bis Okt. 2024 hat nur 1 UAV-Typ eine Zertifizierung der Ersten Klasse (Level 4) erhalten
- 5 UAV-Typen verfügen über eine Zertifizierung der Zweiten Klasse

---

### Singapur (CAAS)

**Status:** BVLOS ohne Operator Permit + Class 1 Activity Permit verboten [[20]](#ref-20)

**Rahmenwerk** [[21]](#ref-21):

- Für jeden BVLOS-Betrieb ist eine strenge Genehmigung erforderlich
- Drohnen >25 kg oder BVLOS-Betrieb erfordern ein vollständiges Operator Permit + Class 1 Permit
- B-RID (Broadcast Remote ID) ab Dez. 2025 verpflichtend für Drohnen >250 g
- B-RID schafft die Grundlage für künftigen skalierten BVLOS- bzw. automatisierten Betrieb
- Regulatorische Sandbox-Programme für beschleunigte Genehmigungen verfügbar

**Betreiber:** Noch offen (Sandbox-Versuche laufen)

---

### China (CAAC)

**Status:** Sondergenehmigung für BVLOS erforderlich [[22]](#ref-22)

**Rahmenwerk** [[23]](#ref-23):

- Vorläufige Vorschriften zur Verwaltung von Flügen unbemannter Luftfahrzeuge (Jan. 2024)
- BVLOS erfordert ein Advanced Operations Certificate + Special Flight Operations Certificate (SFOC)
- Landesweite Höhenbegrenzung von 120 m; darüber ist eine Genehmigung erforderlich
- Vereinfachte digitale Antragsverfahren ermöglichen mehr BVLOS-Flüge [[16]](#ref-16)

**Betreiber:** Noch offen (inländische Betreiber im CAAC-Rahmenwerk)

---

### Indien (DGCA)

**Status:** Noch kein formales BVLOS-Rahmenwerk; nur Testphase [[24]](#ref-24)

**Rahmenwerk** [[25]](#ref-25):

- Betrieb gemäß den Drone Rules, 2021 mit projektbezogenen Ausnahmen
- Apr. 2025: Eigene Kategorie für BVLOS-Pilotenlizenzen für städtische/halbstädtische Gebiete eingeführt
- Die DGCA befindet sich „in einem fortgeschrittenen Stadium der Finalisierung der BVLOS-Regeln“ [[26]](#ref-26)
- Experimentelle Korridore aktiv (z. B. „Medicine from the Sky“ in Telangana)

**Betreiber:**

| Betreiber | Erfolg | Ref. |
|----------|-------------|-----|
| Skye Air Mobility | Weltweit längste BVLOS-Medikamentenlieferung (104 km, Baruipur–Medinipur) | [[27]](#ref-27) |
| Flipkart Health | Partner beim 104-km-Liefertest | — |

---

## Compliance-Anforderungen

In allen Rechtsordnungen laufen die BVLOS-Vorschriften auf gemeinsame technische Anforderungen hinaus. Dieser Abschnitt fasst zusammen, was unsere Flottenplattform implementieren muss.

### 1. Remote ID (überall verpflichtend)

Alle großen Regulierungsbehörden verlangen, dass Drohnen Identifikations- und Standortdaten aussenden.

| Rechtsordnung | Anforderung | Frist | Ref. |
|--------------|-------------|----------|-----|
| USA (FAA) | Standard Remote ID oder Broadcast-Modul | Sept. 2023 (in Kraft) | [[2]](#ref-2) |
| EU (EASA) | Remote ID für Kategorie „speziell“ + Klasse 1+ | Jan. 2024 (in Kraft) | [[11]](#ref-11) |
| Singapur | B-RID für Drohnen >250 g | Dez. 2025 | [[21]](#ref-21) |
| Australien | Remote ID in Prüfung | Noch offen | [[14]](#ref-14) |

**Plattformanforderungen:**
- Hardware-Integration des Broadcast-Moduls (OpenDroneID-kompatibel)
- Network-Remote-ID-Funktion über NATS für flottenweites Tracking
- Protokollierung aller Remote-ID-Übertragungen in JetStream für Audits

---

### 2. Detect-and-Avoid (DAA)

Der BVLOS-Betrieb erfordert die Fähigkeit, andere Luftfahrzeuge und Hindernisse zu erkennen und ihnen auszuweichen.

| Rechtsordnung | Anforderung | Hinweise |
|--------------|-------------|-------|
| USA (FAA) | „Well clear“-Standard; DAA-System für Part 108 erforderlich | Leistungsbasiert, nicht präskriptiv |
| UK (CAA) | DAA-Technologie in der Regel erforderlich (ADS-B, FLARM, Bodenradar) | Einzelfallprüfung im OA-Verfahren |
| Japan | Erforderlich für Level-4-Betrieb in Städten | [[18]](#ref-18) |

**Plattformanforderungen:**
- Integration eines ADS-B-In-Empfängers (Verkehrslagebewusstsein)
- Sensorfusion auf Jetson (Optionen: Kamera, Radar, Lidar)
- Autonome Ausweichmanöver oder Alarmierung des Bedieners
- Protokollierung aller DAA-Ereignisse in JetStream

---

### 3. Geofencing

Hardwareseitig durchgesetzte Grenzen, die Flüge in Sperrgebiete verhindern.

| Rechtsordnung | Anforderung | Hinweise |
|--------------|-------------|-------|
| USA (FAA) | Geofence ist Teil der Betriebsgenehmigung | Im Waiver/Part 108 definiert |
| EU (EASA) | Geo-Awareness in Kategorie „speziell“ erforderlich | Gemäß SORA-Bodenrisikobewertung |
| Alle | Höhenbegrenzungen (typischerweise 120 m/400 ft AGL) | Grundlegend für alle Rahmenwerke |

**Plattformanforderungen:**
- Geofence-Durchsetzung auf dem Flugcontroller (PX4 `GF_*`-Parameter)
- Dynamische Geofence-Aktualisierungen über das Gateway (temporäre Beschränkungen)
- Durchsetzung der Höhenobergrenze
- Protokollierung von Verletzungen und Alarmierung der Flotte

---

### 4. Command-and-Control-Verbindung (C2-Link)

Zuverlässige, redundante Kommunikation zwischen Bediener und Luftfahrzeug.

| Rechtsordnung | Anforderung | Hinweise |
|--------------|-------------|-------|
| USA (FAA) | Leistungsstandards für den C2-Link in Part 108 | Latenz, Verfügbarkeit, Sicherheit |
| EU (EASA) | Zuverlässigkeit des C2-Links gemäß SORA | Risikoabhängige Anforderungen |
| Alle | Failsafe-Verhalten bei C2-Verlust | Rückkehr zum Startpunkt, Landung oder Warteflug |

**Plattformanforderungen:**
- Primärer C2-Link über LTE/5G oder dedizierten Funk
- Backup-C2-Pfad (z. B. Satellit für große Reichweiten)
- Automatischer Failsafe bei Verbindungsverlust (RTL in PX4 konfiguriert)
- Telemetrie zum C2-Link-Status an die Flotte
- Protokollierung der Verbindungsqualität für Audits

---

### 5. Dokumentation der Betriebsgenehmigung

Regulierungsbehörden verlangen eine dokumentierte Risikobewertung und dokumentierte Betriebsverfahren.

| Rahmenwerk | Verwendet von | Wichtige Dokumente |
|-----------|---------|---------------|
| SORA | EASA, UK CAA, CASA, Transport Canada | ConOps, Bodenrisiko, Luftrisiko, OSO-Maßnahmen |
| Part 107 Waiver | FAA (aktuell) | Sicherheitsnachweis, Betriebsgrenzen |
| Part 108 | FAA (kommend) | Leistungsbasierte Compliance |

**Plattformanforderungen:**
- Export von Fluglogs in von den Behörden akzeptierten Formaten
- Erzeugung von Daten zur Unterstützung des ConOps (Fluggebiete, Zeiten, Häufigkeiten)
- Pflege der Zertifizierungsnachweise von Betreibern/Piloten
- Nachverfolgung des Wartungs- und Inspektionsstatus der Luftfahrzeuge

---

### 6. Lizenzierung von Piloten/Betreibern

| Rechtsordnung | BVLOS-Lizenz | Hinweise |
|--------------|---------------|-------|
| USA | Part 107 + Waiver (aktuell), Part 108 (künftig) | [[2]](#ref-2) |
| EU | A2/STS-Zertifikat oder Genehmigung in Kategorie „speziell“ | [[1]](#ref-1) |
| UK | GVC + Operational Authorization | |
| Japan | Erste Klasse (Level 4), Zweite Klasse (Level 3) | [[19]](#ref-19) |
| Kanada | Advanced Pilot Certificate + SFOC | [[9]](#ref-9) |
| Indien | BVLOS-Pilotenlizenz (ab Apr. 2025) | [[25]](#ref-25) |

**Plattformanforderungen:**
- Authentifizierung des Piloten vor der Fluggenehmigung
- Überprüfung von Lizenzen/Zertifizierungen und Nachverfolgung ihres Ablaufs
- Protokollierung der Flugstunden pro Pilot
- Rollenbasierter Zugriff (Pilot vs. Beobachter vs. Flottenmanager)

---

### 7. Flugprotokollierung & Audit-Trail

Alle Rechtsordnungen verlangen umfassende Flugaufzeichnungen.

**Plattformanforderungen:**
- Vollständige Aufzeichnung der Telemetrie (Position, Fluglage, Akku, Sensoren)
- Ereignisprotokollierung mit Zeitstempel (Scharfschalten, Start, Wegpunkte, Landung, Unscharfschalten)
- Unveränderliche Speicherung (JetStream mit Aufbewahrungsrichtlinie)
- Exportfunktion für behördliche Einreichungen
- Kennzeichnung von Vorfällen und Unterstützung bei Untersuchungen

---

### 8. Versicherung & Haftung

| Rechtsordnung | Anforderung |
|--------------|-------------|
| EU | Haftpflichtversicherung gegenüber Dritten verpflichtend |
| USA | Versicherung in der Regel durch die Auflagen der Ausnahmegenehmigung vorgeschrieben |
| Alle | Versicherungsnachweis für die Genehmigung häufig erforderlich |

**Plattformanforderungen:**
- Nachverfolgung der Versicherungspolice pro Luftfahrzeug
- Überprüfung des Versicherungsschutzes vor der Fluggenehmigung
- Dokumentation von Vorfällen für Schadensmeldungen

---

## Externe Integrationen

BVLOS-Compliance erfordert einen Datenaustausch in Echtzeit mit externen Systemen. Diese Integrationen muss unsere Plattform unterstützen.

### Luftraum & Verkehr

| System | Zweck | Datenfluss | Priorität |
|--------|---------|-----------|----------|
| **ADS-B In** | Verkehrslagebewusstsein gegenüber bemannten Luftfahrzeugen | Empfang → DAA-System | Hoch |
| **FLARM** | Verkehrslagebewusstsein gegenüber Segelflugzeugen/Allgemeiner Luftfahrt | Empfang → DAA-System | Mittel |
| **UTM/USS** | UAS-Verkehrsmanagement, Fluggenehmigung | Bidirektional | Hoch |
| **LAANC** (USA) | Automatisierte Luftraumfreigabe | Anfrage/Antwort | Hoch (USA) |
| **AIM/NOTAM** | Temporäre Flugbeschränkungen, Luftraumänderungen | Empfang → Geofence | Hoch |

**UTM Service Suppliers (USS):**
- USA: Airspace Link, Aloft, DroneUp, FlyFreely (siehe [FAA-LAANC-Liste](https://www.faa.gov/uas/getting_started/laanc))
- EU: Je nach Mitgliedstaat unterschiedlich
- Australien: Wing, Airservices Australia

*Hinweis: AirMap hat seinen LAANC-Dienst im Juni 2023 eingestellt.*

---

### Remote-ID-Netzwerke

| System | Zweck | Protokoll | Priorität |
|--------|---------|----------|----------|
| **OpenDroneID** | Standard für Broadcast Remote ID | Bluetooth 4/5, WiFi Beacon | Erforderlich |
| **Network Remote ID** | Flottenweite ID über das Internet | USS-APIs | Erforderlich (USA) |
| **FAA Remote ID** | Remote-ID-Compliance in den USA | USS-Integration | Erforderlich (USA) |

---

### Behördliche Datenbanken

| System | Zweck | Datenfluss | Priorität |
|--------|---------|-----------|----------|
| **Pilotenlizenz-DB** | Überprüfung der Pilotenzertifizierungen | Abfrage → Genehmigung | Hoch |
| **Luftfahrzeugregister** | Überprüfung des Registrierungsstatus | Abfrage → Vorflugkontrolle | Hoch |
| **Versicherungsprüfung** | Bestätigung der Gültigkeit des Versicherungsschutzes | Abfrage → Vorflugkontrolle | Mittel |
| **Wartungsunterlagen** | Nachverfolgung der Lufttüchtigkeit | Bidirektional | Mittel |

---

### Wetter & Umwelt

| System | Zweck | Datenfluss | Priorität |
|--------|---------|-----------|----------|
| **Flugwetter (METAR/TAF)** | Wind, Sicht, Wolkenuntergrenze | Empfang → Flugplanung | Hoch |
| **Radar/Niederschlag** | Wettergefahren in Echtzeit | Empfang → Warnungen während des Flugs | Mittel |
| **Solar/Geomagnetisch** | Prognosen zur GPS-Zuverlässigkeit | Empfang → Flugplanung | Niedrig |

---

### Notfalldienste

| System | Zweck | Datenfluss | Priorität |
|--------|---------|-----------|----------|
| **ATC (falls erforderlich)** | Koordination im kontrollierten Luftraum | Sprach-/Datenverbindung | Situationsabhängig |
| **Notrufdienste (911/112)** | Meldung von Vorfällen | Ausgehende Alarme | Erforderlich |
| **NOTAM-Einreichung** | Veröffentlichung von Flugbetrieb | Ausgehend → AIM | Situationsabhängig |

---

### Integrationsarchitektur

```
┌─────────────────────────────────────────────────────────────────────┐
│                         FLEET PLATFORM                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │
│   │   Gateway   │    │  JetStream  │    │   Fleet     │             │
│   │  (Vehicle)  │    │  (Logging)  │    │  Manager    │             │
│   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘             │
│          │                  │                  │                     │
└──────────┼──────────────────┼──────────────────┼─────────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│                      INTEGRATION LAYER                               │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │   UTM    │  │  Remote  │  │ Weather  │  │ License  │             │
│  │   USS    │  │    ID    │  │   API    │  │   DB     │             │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘             │
│       │             │             │             │                    │
└───────┼─────────────┼─────────────┼─────────────┼────────────────────┘
        │             │             │             │
        ▼             ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │Airspace │   │   FAA   │   │  NOAA   │   │  CAA    │
   │  Link   │   │ DroneZone│  │ Aviation│   │ Registry│
   │ Aloft   │   │ NASA DIP│   │ Weather │   │         │
   └─────────┘   └─────────┘   └─────────┘   └─────────┘
```

---

### Integrations-Roadmap

| Integration | Status | Hinweise |
|-------------|--------|-------|
| OpenDroneID-Broadcast | ✅ Unterstützt | Hardwaremodul erforderlich |
| ADS-B-In-Empfänger | 🔧 Roadmap | uAvionix pingRX o. Ä. |
| UTM/USS (Airspace Link) | 🔧 Roadmap | AirHub API, gestaffelte Preise |
| UTM/USS (Aloft) | 🔧 Roadmap | REST-API-Integration |
| LAANC | 🔧 Roadmap | Über USS-Partner |
| Flugwetter | 🔧 Roadmap | aviationweather.gov (offen, 100 Anfragen/Min.) |
| NOTAM-Feed | 🔧 Roadmap | NASA DIP oder FAA SWIM |
| Überprüfung von Pilotenlizenzen | 🔧 Roadmap | APIs je Rechtsordnung |
| Versicherungsprüfung | 🔧 Roadmap | APIs der Anbieter |

---

## Compliance-Matrix

| Anforderung | Plattformkomponente | Status |
|-------------|-------------------|--------|
| Remote-ID-Broadcast | OpenDroneID-Modul + Gateway | ✅ Unterstützt |
| Remote-ID-Protokollierung | JetStream-Streams | ✅ Unterstützt |
| Detect-and-Avoid | Sensorfusion auf Jetson | 🔧 Roadmap |
| ADS-B In | Integration des Verkehrsempfängers | 🔧 Roadmap |
| Geofencing (statisch) | PX4 `GF_*`-Parameter | ✅ Unterstützt |
| Geofencing (dynamisch) | Richtlinienaktualisierungen über das Gateway | ✅ Unterstützt |
| C2-Link primär | LTE/5G über NATS | ✅ Unterstützt |
| C2-Link Backup | Satellitenintegration | 🔧 Roadmap |
| Failsafe bei Verbindungsverlust | PX4-Failsafe-Konfiguration | ✅ Unterstützt |
| Flugprotokollierung | JetStream-Telemetrie-Streams | ✅ Unterstützt |
| Audit-Export | Tools für den Log-Export | ✅ Unterstützt |
| Pilotenauthentifizierung | Autorisierungssystem | ✅ Unterstützt |
| Lizenzverfolgung | Flottenmanagement | 🔧 Roadmap |
| Versicherungsverfolgung | Flottenmanagement | 🔧 Roadmap |

Siehe [Sicherheitsmodell]({{< relref "/fleet/safety" >}}), um zu erfahren, wie sich diese Komponenten in die gesamte Sicherheitsarchitektur einfügen.

---

## Quellen

<span id="ref-1"></span>**[1]** EASA. "Drones (UAS) - Frequently Asked Questions." European Union Aviation Safety Agency. https://www.easa.europa.eu/en/the-agency/faqs/drones-uas

<span id="ref-2"></span>**[2]** Drone Pilot Ground School. "Part 108: What We Know So Far about the FAA's BVLOS Rule." https://www.dronepilotgroundschool.com/part-108/

<span id="ref-3"></span>**[3]** DroneLife. "FAA Progress on BVLOS Rules and Advanced Air Mobility Integration." Aug. 2024. https://dronelife.com/2024/08/01/faa-progress-on-bvlos-rules-and-advanced-air-mobility-integration/

<span id="ref-4"></span>**[4]** FAA. "BEYOND Program." Federal Aviation Administration. https://www.faa.gov/uas/programs_partnerships/beyond

<span id="ref-5"></span>**[5]** DroneUp. "DroneUp Awarded Landmark FAA Approval for Beyond Visual Line of Sight (BVLOS) in the U.S." Jan. 2024. https://www.droneup.com/news/droneup-awarded-bvlos

<span id="ref-6"></span>**[6]** DroneDJ. "Amazon gets FAA BVLOS approval for Prime Air drone delivery." Mai 2024. https://dronedj.com/2024/05/30/amazon-gets-faa-bvlos-approval-for-prime-air-drone-delivery/

<span id="ref-7"></span>**[7]** DroneDJ. "Asylon bags new FAA approval for BVLOS drone operations." Apr. 2024. https://dronedj.com/2024/04/26/asylon-drone-faa-bvlos-dronesentry/

<span id="ref-8"></span>**[8]** DroneLife. "FAA Grants Historic BVLOS Approval for Multiple Operators in DFW Area." Juli 2024. https://dronelife.com/2024/07/31/faa-grants-historic-bvlos-approval-for-multiple-operators-in-dfw-area/

<span id="ref-9"></span>**[9]** MLT Aikins. "International drone regulations: Key updates on BVLOS, UTM and use cases." https://www.mltaikins.com/insights/international-drone-regulations-key-updates-on-bvlos-utm-and-use-cases/

<span id="ref-10"></span>**[10]** FlytBase. "Understanding EU Drone Regulations for Docked BVLOS Operations." https://www.flytbase.com/blog/eu-drone-regulations-for-bvlos-operations

<span id="ref-11"></span>**[11]** Elsight. "EASA Remote ID Requirements: Compliance for Drone Operators." https://www.elsight.com/blog/primer-on-easa-remote-id-regulations/

<span id="ref-12"></span>**[12]** Sparkco. "Autonomous Drone Delivery Regulatory Approval." https://sparkco.ai/blog/autonomous-drone-delivery-regulatory-approval

<span id="ref-13"></span>**[13]** sees.ai. "Why Us - First UK routine BVLOS permission." https://www.sees.ai/why-us/

<span id="ref-14"></span>**[14]** CASA. "BVLOS drone operations in regional Australia - Consultation." Civil Aviation Safety Authority. https://consultation.casa.gov.au/stakeholder-engagement-group/consultation.2023-10-05.5578154857/

<span id="ref-15"></span>**[15]** Measure Australia. "New CASA Drone Regulations Incoming!" https://www.measureaustralia.com.au/news/new-casa-drone-regulations-incoming

<span id="ref-16"></span>**[16]** Droneii. "Commercial Drone Regulation Report 2025." Drone Industry Insights. https://droneii.com/product/commercial-drone-regulation-report

<span id="ref-17"></span>**[17]** DroneLife. "Sphere Drones Secures BVLOS Area Approval to Advance Commercial Drone Use in Australia." Mai 2024. https://dronelife.com/2024/05/08/sphere-drones-secures-bvlos-area-approval-to-advance-commercial-drone-use-in-australia/

<span id="ref-18"></span>**[18]** Unmanned Airspace. "Japan's revised drone laws to permit BVLOS flights over people come into effect." https://www.unmannedairspace.info/emerging-regulations/japanese-revised-drone-laws-to-permit-bvlos-flights-over-people-come-into-effect/

<span id="ref-19"></span>**[19]** Fly Eye. "Japan Drone Laws (2025)." https://www.flyeye.io/japan-drone-laws/

<span id="ref-20"></span>**[20]** Drone Laws. "Singapore Drone Laws 2025." https://drone-laws.com/drone-laws-in-singapore/

<span id="ref-21"></span>**[21]** Heron AirBridge. "B-RID in Singapore: Full Guide to the New Drone Regulations." https://heron-airbridge.com/singapore-brid-regulation/

<span id="ref-22"></span>**[22]** Fly Eye. "Advanced Drone Operations in China." https://www.flyeye.io/chinese-drone-regulations-advanced-operations/

<span id="ref-23"></span>**[23]** Fly Eye. "Chinese Drone Laws (2025)." https://www.flyeye.io/chinese-drone-laws/

<span id="ref-24"></span>**[24]** Commercial UAV News. "India Makes a Giant Leap Toward BVLOS Operations." https://www.commercialuavnews.com/drone-delivery/india-makes-a-giant-leap-toward-bvlos-operations

<span id="ref-25"></span>**[25]** Fly and Tech. "Drone Rules in India 2025 – Latest DGCA Guidelines Update." https://flyandtech.com/drone-rules-india-2025/

<span id="ref-26"></span>**[26]** Sigma Chambers. "Drone Law Brief | Vol 1: September & October 2025." https://www.sigmachambers.in/post/drone-law-watch-vol-1-september-october-2025

<span id="ref-27"></span>**[27]** ITLN. "How will BVLOS drone delivery ops impact the Indian logistics industry?" https://www.itln.in/cargo-drones/how-will-bvlos-drone-delivery-ops-impact-the-indian-logistics-industry-1353269

---

## Zurück zur Übersicht

[← Architektur für autonome Fahrzeugflotten]({{< relref "/fleet" >}})
