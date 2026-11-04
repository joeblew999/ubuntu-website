---
title: "Technologie"
meta_title: "Technologie-Stack | Ubuntu Software"
description: "Basierend auf Go und NATS JetStream - unsere Technologieentscheidungen für Performance, Zuverlässigkeit und Einfachheit."
image: "/images/spatial.svg"
draft: false
---

## Unser Technologie-Stack

Wir setzen auf Technologien, die wir nach Leistung, Zuverlässigkeit und langfristiger Wartbarkeit auswählen.

---

## Go

**Warum Go:**

- **Leistung** — Kompiliert, statisch typisiert, minimaler Laufzeit-Overhead
- **Einfachheit** — Ein Weg, Dinge zu tun, standardmäßig gut lesbar
- **Nebenläufigkeit** — Goroutines und Channels sind fest in die Sprache integriert
- **Deployment** — Einzelne Binärdatei, keine Abhängigkeiten, Cross-Kompilierung
- **Ökosystem** — Starke Standardbibliothek, hervorragende Werkzeuge

Go ist unsere Hauptsprache im gesamten Stack, von Backend-Diensten über CLI-Tools bis hin zur Robotik-Integration.

---

## NATS JetStream

**Warum NATS JetStream:**

- **Persistenz** — Dauerhafte Nachrichtenspeicherung mit Replay
- **Exactly-once-Zustellung** — Zuverlässige Nachrichtenverarbeitung
- **Leichtgewichtig** — Einzelne Binärdatei, minimaler Ressourcenbedarf
- **Skalierbar** — Clustering und horizontale Skalierung integriert
- **Echtzeit** — Latenz unter einer Millisekunde für Pub/Sub

NATS JetStream treibt unsere ereignisgesteuerte Architektur an und verbindet Design-Tools, Simulation, digitale Zwillinge und physische Geräte.

---

## Architekturprinzipien

| Prinzip | Umsetzung |
|-----------|----------------|
| **Offline-first** | Lokale Datenhaltung, Synchronisierung bei bestehender Verbindung |
| **Ereignisgesteuert** | NATS JetStream für die gesamte Kommunikation zwischen Diensten |
| **Offene Standards** | STEP, IFC, keine proprietären Formate |
| **Hardware-unabhängig** | Abstraktionsschichten für Portabilität |
| **Selbstbestimmt** | Überall bereitstellbar: Cloud, On-Premises, Air-Gapped |

---

## Technologiebereiche

### Robotik

Unser Robotik-Stack basiert auf Viam RDK, einem Open-Source-Robotik-Entwicklungskit, das Hardware-Abstraktion, Bewegungsplanung und Computer Vision bereitstellt.

[Robotik-Stack →](/technology/robotics/)

### Linux & plattformübergreifend

Jahrzehntelange Linux-Erfahrung. Plattformübergreifende Anwendungen für Windows, Mac, Linux, iOS und Android. Ein eigenes GUI-Framework nach Jahren mit Qt, Flutter und Electron.

[Linux & plattformübergreifend →](/technology/linux/)

### Sicherheit & Compliance

Eine selbstbestimmte Architektur, die die Compliance mit SOC 2, FedRAMP, HIPAA und ISO 27001 vereinfacht. Air-Gapped-Bereitstellung, kein Call-Home, vollständige Datenkontrolle.

[Sicherheit & Compliance →](/technology/security/)

---

## Bereitstellungsplattformen

Wir stellen überall bereit:

| Plattform | Anwendungsfall |
|----------|----------|
| **Linux** | Server, Desktop, Embedded |
| **OpenBSD** | Sicherheitskritische Systeme |
| **Windows** | Unternehmens-Desktop |
| **macOS** | Design und Entwicklung |
| **iOS/Android** | Mobile Anwendungen |
| **Docker/Kubernetes** | Container-Orchestrierung |

---

## Open-Source-Fundament

Wir bauen auf Open Source auf und tragen dazu bei:

| Komponente | Technologie |
|-----------|------------|
| Sprache | Go |
| Messaging | NATS JetStream |
| Zusammenarbeit | Automerge (CRDT) |
| 3D-Formate | STEP, IFC |
| Robotik | Viam RDK |
| Vision | Intel RealSense, YOLOv8 |

---

## Mehr erfahren

- [Robotik-Stack →](/technology/robotics/)
- [Linux & plattformübergreifend →](/technology/linux/)
- [Sicherheit & Compliance →](/technology/security/)
- [Spatial-Plattform →](/platform/spatial/)
- [Foundation →](/platform/foundation/)
- [Kontakt →](/contact/)
