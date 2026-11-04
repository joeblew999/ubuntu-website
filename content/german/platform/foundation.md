---
title: "Grundlage"
meta_title: "Grundlagentechnologie | Ubuntu Software"
description: "Offline-First-Architektur, universelle Bereitstellung und selbstsouveräne Daten. Das technologische Fundament, das sowohl die Publish- als auch die Spatial-Plattform antreibt."
image: "/images/foundation.svg"
draft: false
---

## Für die reale Welt gebaut

Das Internet fällt aus. Teams arbeiten über Kontinente hinweg. Die Server gehören Ihnen. Wir haben für genau diese Realität gebaut.

Bei unserem Fundament geht es nicht um Features—es geht darum, wie Software funktionieren sollte. Offline-first. Selbstsouverän. Universelles Deployment. Diese Prinzipien ziehen sich durch alles, was wir entwickeln.

---

## Offline-First-Architektur

**Ohne Internet arbeiten. Bei Verbindung synchronisieren. Niemals Daten verlieren.**

Echte Arbeit findet an Orten mit schlechtem WLAN statt—auf Baustellen, in Fabrikhallen, in Behördengebäuden, in Krankenzimmern, in abgelegenen Büros. Unser Fundament geht von fehlender Verbindung aus, nicht von ständiger Konnektivität.

### So funktioniert es

| Komponente | Was sie leistet |
|-----------|--------------|
| **Local-First-Daten** | Ihre Daten liegen zuerst auf Ihrem Gerät, nicht auf einem entfernten Server |
| **Automerge CRDT** | Konfliktfreies Zusammenführen, wenn mehrere Personen gleichzeitig bearbeiten |
| **Hintergrund-Synchronisierung** | Automatische Synchronisierung, sobald die Verbindung wiederhergestellt ist |
| **Offline-Warteschlangen** | Aktionen werden lokal eingereiht und ausgeführt, sobald es möglich ist |

**Keine Ladekreisel. Keine „Verbindung verloren“-Fehler. Einfach arbeiten.**

---

## Überall bereitstellen

**Eine Codebasis. Jede Plattform. Native Erfahrung.**

Dieselbe Anwendung läuft im Browser, auf dem Desktop und auf Mobilgeräten. Nicht drei separate Produkte—eine Codebasis, die überall bereitgestellt wird.

### Unterstützte Plattformen

| Plattform | Bereitstellung |
|----------|----------|
| **Web** | Jeder moderne Browser—Chrome, Firefox, Safari, Edge |
| **Desktop** | Native Apps für Windows, macOS, Linux |
| **Mobil** | Native Apps für iOS und Android |

Ihr Team nutzt Desktop-Rechner im Büro. Außendienstmitarbeiter nutzen Tablets. Kunden nutzen Smartphones. Alle greifen auf dasselbe System mit denselben Daten zu.

**Einmal schreiben. Überall bereitstellen. Eine Codebasis pflegen.**

---

## Cloud-Sync-Optionen

**Ihre Cloud. Unsere Cloud. Keine Cloud. Sie entscheiden.**

Die Synchronisierung setzt unsere Server nicht voraus. Verbinden Sie sich mit der Infrastruktur, die für Ihre Organisation sinnvoll ist.

### Bereitstellungsmodelle

| Modell | Ideal für |
|-------|----------|
| **Ubuntu Software Cloud** | Schnellste Einrichtung, wir übernehmen den Betrieb |
| **Ihre Cloud** | AWS, Azure, GCP—Ihre Infrastruktur, unsere Software |
| **On-Premises** | Ihr Rechenzentrum, volle Kontrolle |
| **Hybrid** | Einige Daten in der Cloud, sensible Daten vor Ort |
| **Air-Gapped** | Vollständig getrennte Netzwerke, Verteidigung und Hochsicherheitsumgebungen |

Wechsel zwischen den Modellen? Unkompliziert. Ihr Datenformat ändert sich nicht, je nachdem, wo die Daten liegen.

**Kein Vendor-Lock-in. Keine erzwungene Cloud-Abhängigkeit. Echte Flexibilität bei der Bereitstellung.**

---

## Selbstsouverän

**Ihre Daten. Ihre Server. Ihre Regeln.**

Selbstsouverän bedeutet, dass Sie Ihre Infrastruktur kontrollieren. Nicht „Ihre Daten, gehostet zu unseren Bedingungen“—sondern wirklich Ihre.

### Was Selbstsouveränität bedeutet

- **Überall ausführen** — Ihr Rechenzentrum, Ihr Cloud-Konto, Ihr Laptop
- **Kein Call-Home** — Die Software funktioniert, ohne unsere Server zu kontaktieren
- **Alles exportieren** — Standardformate, vollständige Datenportabilität
- **Kein Nutzungs-Tracking** — Wir sehen nicht, was Sie mit unserer Software tun
- **Optionen für unbefristete Lizenzen** — Läuft weiter, selbst wenn es uns nicht mehr gibt

**Das ist nicht nur Datenschutz—es ist betriebliche Unabhängigkeit.**

---

## Einbettbar

**In bestehende Systeme integrieren. Nicht herausreißen und ersetzen.**

Organisationen haben bestehende Workflows, Datenbanken und Identitätssysteme. Unser Fundament integriert, statt zu ersetzen.

### Integrationsmuster

| Muster | Anwendungsfall |
|---------|----------|
| **API-First** | Alles ist programmatisch zugänglich |
| **Datenbank-Konnektoren** | Lesen/Schreiben in Ihre bestehenden Datenbanken |
| **SSO-Integration** | Ihr Identitätsanbieter, kein weiteres Login |
| **Webhook-Events** | Push-Benachrichtigungen an Ihre Systeme |
| **White-Label** | Ihr Branding, unsere Engine |

**Erweitern Sie Ihre Systeme. Geben Sie sie nicht auf.**

---

## Wellknown Gateway

**Auf Big-Tech-Plattformen veröffentlichen. Die Beziehung selbst besitzen.**

Die Gateways des Webs—Google, Apple, YouTube—haben eine enorme Reichweite. Doch AUF ihren Plattformen zu veröffentlichen bedeutet nicht, IHNEN zu gehören.

### Die Beziehung umkehren

Traditioneller Ansatz:
```
User → YouTube (owns everything) → Your content (captive)
```

Wellknown-Ansatz:
```
User → Your Gateway → Your System (primary)
                   ↳→ YouTube (mirror for discovery)
```

**Sie kontrollieren die Eingangstür.** Big Tech wird zu einem optionalen Vertriebskanal, nicht zu einem Gefängnis.

### So funktioniert es

| Fähigkeit | Was das bedeutet |
|------------|---------------|
| **Ihre URIs überall** | Links verweisen auf IHR Gateway, nicht auf deren |
| **Intelligentes Routing** | iOS-Nutzer zu Apple, Android-Nutzer zu Google, Web-Nutzer zu Ihrem Player |
| **Spiegel-Veröffentlichung** | Automatische Veröffentlichung von Kopien auf YouTube, Google Maps, Apple Kalender |
| **Analysen in Ihrem Besitz** | Alles sehen, alle erfassen, keine Datenweitergabe |
| **Integrierte Exit-Strategie** | Jede Plattform aus dem Routing entfernen, ohne Links zu brechen |

### Funktioniert mit allem

- **Video** — Hosting auf Ihrem Server, Spiegelung zu YouTube/Twitch für mehr Reichweite
- **Kalender** — Ihr CalDAV-Server, Synchronisierung mit Google/Apple für mehr Komfort
- **Karten** — Ihre Geodaten, integriert mit Google/Apple Maps
- **E-Mail** — Ihr Mailserver, kompatibel mit Gmail
- **Dateien** — Ihr Speicher, selektives Teilen über Drive/Dropbox

**Veröffentlichen Sie AUF ihren Plattformen. Lassen Sie sich niemals EINSPERREN.**

---

## 25 Jahre Enterprise-Erfahrung

**Ubuntu Software wurde 1999 gegründet.** Wir haben Unternehmenssysteme durch die Dotcom-Ära, die mobile Revolution, den Wandel zur Cloud und den Aufstieg der KI hindurch entwickelt.

### Was uns die Erfahrung gelehrt hat

- **Anbieter verschwinden** — Von Tag eins an auf Datenportabilität auslegen
- **Netzwerke fallen aus** — Offline-first ist nicht optional, sondern unverzichtbar
- **Anforderungen ändern sich** — Offene Standards überdauern proprietäre Formate
- **Skalierung überrascht** — Architektur ist wichtiger als Optimierung
- **Integration ist schwierig** — Von Anfang an dafür entwerfen, nicht nachträglich anflanschen

Wir haben gesehen, was funktioniert und was nicht. Dieses Fundament spiegelt jahrzehntelange Erfahrungen aus Produktionsumgebungen wider.

---

## Das Fundament treibt alles an

Sowohl die Publish- als auch die Spatial-Plattform erben diese Fähigkeiten automatisch:

| Fähigkeit | Publish | Spatial |
|------------|---------|---------|
| Offline-Bearbeitung | Dokumente ohne Internet bearbeiten | 3D-Modelle ohne Internet entwerfen |
| Echtzeit-Synchronisierung | Mehrere Bearbeiter, ein Dokument | Mehrere Designer, ein Modell |
| Universelles Deployment | Formulare auf jedem Gerät | 3D-Viewer auf jedem Gerät |
| Selbst gehostet | Ihr Dokumentenserver | Ihr Design-Server |
| Cloud-Optionen | Verwaltet oder selbst verwaltet | Verwaltet oder selbst verwaltet |

**Wählen Sie Ihre Plattform. Das Fundament gibt es gratis dazu.**

---

## Technische Details

Für Teams, die unsere Architektur evaluieren:

| Schicht | Technologie |
|-------|------------|
| **Speicher** | SQLite überall—auf lokalen Geräten und Servern |
| **Sync-Engine** | CRDT-basierte Replikation über NATS JetStream |
| **Messaging** | NATS JetStream |
| **UI-Framework** | Plattformübergreifendes natives Rendering |
| **API** | HTTP REST + SSE (Server-Sent Events) |
| **KI-Integration** | Model Context Protocol (MCP) |
| **Auth** | OIDC-kompatibel, bringen Sie Ihren eigenen IdP mit |

### NATS JetStream

Unser Messaging-Backbone ist ein Incubating-Projekt der **CNCF (Cloud Native Computing Foundation)**—derselben Foundation, unter deren Dach auch Kubernetes steht.

| Eigenschaft | Detail |
|----------|--------|
| **Governance** | CNCF-Incubating-Projekt |
| **Lizenz** | Apache 2.0 (100 % Open Source) |
| **Zertifizierung** | Synadia verfügt über SOC 2 Type II |
| **Community** | Tausende Produktivbereitstellungen weltweit |

Open Source. Für den Unternehmenseinsatz gehärtet. Auditiert.

[Details zu Sicherheit & Compliance →]({{< relref "/technology/security" >}})

### Verteiltes SQLite

Jeder Knoten—Ihr Laptop, Ihr Smartphone, Ihre Server—führt SQLite aus. Änderungen werden über NATS JetStream mit CRDT-Semantik repliziert.

- **Jeder Server kann ausfallen** — Die anderen bedienen weiter
- **Jeder Server kann offline gehen** — Synchronisierung bei erneuter Verbindung
- **Kein Single Point of Failure** — Echte verteilte Architektur
- **Überall dieselbe Datenbank** — Vom lokalen Gerät bis zum globalen Cluster

**Standardtechnologien. Kein proprietärer Lock-in.**

---

## Skalierung ohne Grenzen

### Kein Single Point of Failure (SPOF)

Jede Komponente ist redundant. Kein einzelner Server, Dienst oder Rechenzentrum kann das System lahmlegen. Fällt ein beliebiger Knoten aus—das System läuft weiter.

### Kein Single Point of Performance (SPOP)

Die Rechenleistung skaliert horizontal. Sie erweitern die Kapazität durch zusätzliche Knoten, nicht durch den Kauf größerer Server. Workloads verteilen sich automatisch auf die verfügbaren Ressourcen.

### Hunderte Rechenzentren

Eine Architektur, die für globale Verteilung konzipiert ist:

| Fähigkeit | Was das bedeutet |
|------------|---------------|
| **Überall bereitstellen** | Cloud, On-Premises, Edge, Air-Gapped |
| **Nah an den Nutzern bereitstellen** | Geringe Latenz, lokale Compliance |
| **Für Redundanz replizieren** | Regionale Ausfälle überstehen |
| **Partitionstoleranz** | Unabhängiger Betrieb bei Netzwerktrennung |

### Bereitstellungsoptionen

| Methode | Anwendungsfall |
|--------|----------|
| **Binaries** | Deployment als einzelne Datei, minimale Abhängigkeiten |
| **Docker** | Containerisierte, reproduzierbare Umgebungen |
| **Kubernetes** | Orchestrierte, automatisch skalierende Cluster |

**Vom einzelnen Laptop bis zu Hunderten Rechenzentren. Dieselbe Architektur. Dieselbe Codebasis.**

---

## Jetzt starten

Das Fundament ist bereits integriert. Wenn Sie Publish oder Spatial nutzen, erhalten Sie automatisch Offline-first, universelles Deployment und selbstsouveräne Optionen.

[Publish entdecken →](/platform/publish/) | [Spatial entdecken →](/platform/spatial/) | [Linux & plattformübergreifend →](/technology/linux/) | [Kontakt →](/contact/)
