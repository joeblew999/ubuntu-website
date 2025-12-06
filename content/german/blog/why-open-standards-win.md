---
title: "Warum offene Standards gewinnen"
meta_title: "Warum offene Standards gewinnen | Ubuntu Software"
description: "STEP, IFC und das Argument gegen proprietäre Abhängigkeit für 3D-Design und Engineering-Daten."
date: 2024-09-20T05:00:00Z
image: "/images/blog/open-standards.svg"
categories: ["Industry", "Standards"]
author: "Gerard Webb"
tags: ["open-standards", "step", "ifc", "interoperability", "cad"]
draft: false
---

Etwa alle zehn Jahre lernt eine Branche dieselbe Lektion: Proprietäre Abhängigkeit skaliert nicht.

Das Web hat es gelernt. Enterprise-Software hat es gelernt. Cloud Computing hat es gelernt.

Jetzt ist es Zeit für 3D-Design und Engineering, es zu lernen.

## Der aktuelle Stand der Dinge

Versuchen Sie dieses Experiment: Nehmen Sie ein 3D-Modell aus einem großen CAD-System und öffnen Sie es in einem anderen.

Was Sie finden werden:
- **Geometrieverlust**: Features werden nicht übersetzt. Constraints verschwinden.
- **Metadaten weg**: Alle Engineering-Informationen—Materialien, Toleranzen, Montagebeziehungen—verloren.
- **Manuelle Nacharbeit**: Jemand verbringt Stunden damit, das nachzubauen, was bereits existierte.

Das ist kein Fehler. Es ist ein Geschäftsmodell.

**Anbieter profitieren von Abhängigkeit. Nutzer leiden darunter.**

## Die Kosten proprietärer Formate

### Kollaborationsreibung

Wenn Ihr Roboterlieferant System A verwendet, Ihr Anlagenplaner System B und Ihr Simulationsteam System C, ist jede Übergabe eine Übersetzungsübung.

Informationen verschlechtern sich bei jeder Konvertierung. Ingenieure verbringen Zeit mit Dateiwrangling statt mit Engineering.

### Anbieterabhängigkeit

Ihre gesamte Designhistorie lebt in einem Format, das nur ein Anbieter kontrolliert. Sie legen die Preise fest. Sie bestimmen den Upgrade-Zeitplan. Sie entscheiden, wann Features eingestellt werden.

Ihr Engineering-IP wird als Geisel gehalten.

### Innovationsbarrieren

Möchten Sie KI-Tools auf Ihren Designdaten aufbauen? Viel Glück beim Zugriff über proprietäre APIs, die sich mit jeder Version ändern.

Möchten Sie mit der neuesten Physiksimulation integrieren? Hoffen Sie besser, dass Ihr CAD-Anbieter eine Partnerschaft hat.

Innovation stirbt an Formatgrenzen.

## Die offene Alternative

### STEP: Der Geometrie-Standard

ISO 10303 (STEP) existiert seit den 1990er Jahren. Es ist langweilig. Es funktioniert.

STEP erfasst:
- 3D-Geometrie mit voller Präzision
- Montagestrukturen und -beziehungen
- Produktfertigungsinformationen (PMI)
- Materialeigenschaften

Es ist nicht perfekt. Aber es ist universell.

### IFC: Gebäude, die sprechen

Industry Foundation Classes (IFC) macht für Gebäude, was STEP für Produkte macht.

Jede Wand, Tür, jeder Raum und jedes System—definiert in einem offenen Format, das jede Software lesen und schreiben kann.

BIM-Interoperabilität ist kein Traum. IFC macht es möglich.

### Der aufkommende Stack

Moderne offene Standards gehen über statische Geometrie hinaus:

- **glTF**: Leichtgewichtiges 3D für Visualisierung und AR/VR
- **USD**: Szenenbeschreibung für Simulation und Rendering
- **SDF**: Robot- und Umgebungsdefinition
- **URDF**: Robot-Beschreibungsformat

Ein Ökosystem bildet sich. Tools, die auf offenen Grundlagen aufgebaut sind, können teilnehmen.

## Was offene Standards ermöglichen

### Echter Wettbewerb

Wenn Ihre Daten nicht eingesperrt sind, können Sie Tools nach Fähigkeit wählen, nicht nach Gefangenschaft.

Anbieter konkurrieren um Features, nicht darum, wie schwierig sie die Migration machen.

### Ökosystem-Innovation

Offene Formate ermöglichen ein Ökosystem spezialisierter Tools:
- KI-Assistenten, die plattformübergreifend funktionieren
- Simulationsengines, die jede Geometrie akzeptieren
- Kollaborationstools, die nicht erfordern, dass jeder dieselbe Lizenz besitzt

### Zukunftssicherheit

Standardisierungsgremien bewegen sich langsam. Das ist ein Feature.

Eine STEP-Datei von 2000 öffnet sich heute noch. Wird sich Ihr proprietäres Format von 2020 im Jahr 2040 öffnen?

## Die hybride Realität

Seien wir praktisch: Reine Open-Standards-Workflows existieren noch nicht.

Die echte Strategie ist:
1. **Native Formate für aktive Arbeit**: Nutzen Sie das beste Tool für jeden Job
2. **Offene Formate für den Austausch**: Standardformate an jedem Übergabepunkt
3. **Offene Formate für die Archivierung**: Langfristige Speicherung in Formaten, die Sie kontrollieren

Das ist kein Idealismus. Es ist Risikomanagement.

## Der Branchenwandel

Der Schwung baut sich auf:

**Regierungsmandate**: Mehr Behörden fordern offene Formate für Beschaffung und Archivierung.

**Branchenkonsortien**: Organisationen wie buildingSMART treiben die IFC-Adoption voran.

**KI-Anforderungen**: Maschinelles Lernen braucht Trainingsdaten, die nicht weggesperrt sind.

**Cloud-Kollaboration**: Echtzeit-Kollaborationsplattformen wählen offene Grundlagen.

Die Anbieter, die offene Standards annehmen, werden gewinnen. Die, die dagegen kämpfen, werden umgangen.

## Den Übergang schaffen

Wenn Sie neu anfangen, bauen Sie auf offenen Grundlagen:
- Wählen Sie Tools mit starker Unterstützung offener Formate
- Verlangen Sie Standards-Compliance in Lieferantenverträgen
- Etablieren Sie Open-Format-Checkpoints in Ihren Workflows
- Archivieren Sie in Formaten, die Sie kontrollieren, nicht in Formaten, die Sie kontrollieren

Wenn Sie migrieren, beginnen Sie an den Grenzen:
- Neue Integrationen nutzen offene Formate
- Neue Projekte pilotieren offene Workflows
- Schrittweise Migration, während Tools und Workflows reifen

## Das größere Bild

Offene Standards handeln nicht von Technologie. Sie handeln von Macht.

Wer kontrolliert Ihre Engineering-Daten? Wer entscheidet, welche Tools Sie verwenden können? Wem gehört Ihre Designhistorie?

Proprietäre Formate antworten: dem Anbieter.

Offene Standards antworten: Ihnen.

Deshalb gewinnen offene Standards. Nicht weil sie technisch überlegen sind (obwohl sie es oft sind). Weil sie die Anreize richtig ausrichten.

Ihre Daten. Ihre Wahl. Ihre Zukunft.

---

*Wir haben unsere Plattform auf STEP, IFC und offenen APIs aufgebaut. [Sehen Sie, wie es funktioniert →](/platform)*
