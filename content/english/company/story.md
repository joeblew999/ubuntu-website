---
title: "Our Story"
meta_title: "The Ubuntu Software Story"
description: "From Flame on SGI to government sensor fusion to prefab factories in Sweden - why we're building what should exist for everyone."
draft: false
---

## It Started With Flame

Years ago, I worked with Flame on Silicon Graphics machines. If you don't know Flame, it's Autodesk's high-end compositing system - the kind of tool that Hollywood uses for visual effects, the kind of tool that costs serious money and runs on serious hardware.

What made Flame different wasn't just the power. It was how it thought about dimensions.

Flame didn't treat 2D as one thing, 3D as another thing, and time as yet another thing to bolt on later. It treated them all as first-class citizens in one unified environment. 2D, 3D, and 4D - where the fourth dimension is time, how things change, how they animate, how they evolve. You could work in 2D, move into 3D, scrub through time - and it all felt like the same tool, the same mental model, the same way of thinking.

That was my first real glimpse of what's possible when software doesn't artificially separate dimensions. When the tool matches how the real world actually works.

---

## The Tool Chaos

Over the years I worked with so many tools. AutoCAD. SketchUp. Blender. Photoshop. Illustrator. And dozens of others for specific industries and purposes.

Every single one of them does something well. And every single one of them is awful at working with the others.

You design something in one tool. You need to bring it into another tool. Export. Import. Lose information. Manually fix things. Re-enter data. Hope the conversion didn't break something subtle that you won't discover until much later.

The CAD tools don't talk to the desktop publishing tools. The 3D modeling tools don't talk to the vector graphics tools. The animation tools don't talk to the engineering tools. Everyone's working in their own silo, passing files back and forth, losing fidelity at every handoff.

And here's the thing I realized: a CAD system, a game engine, and a video editor are fundamentally the same thing. They're all dealing with geometry in space and time. They're all compositing - combining elements into a coherent whole. The only difference is when and where you're doing that compositing. Are you compositing for a building? A game level? A film sequence? A factory floor? A robot's path through space?

It's all the same problem. But we've built completely separate tools for each use case, and none of them talk to each other.

---

## Then Came the Government Work

Later I worked on government systems. Specifically, sensor fusion.

The problem was this: you have data coming in from multiple sources. Different sensors, different formats, different refresh rates, different levels of reliability. Radar. Cameras. GPS. Communications intercepts. All of it pouring in simultaneously.

The job was to take all of that chaos and render it into a coherent display. A 2D map view. A 3D spatial view. Something that a human operator could actually look at and understand what was happening.

If you've ever seen what the display looks like inside a Tesla - how it takes camera feeds and radar and ultrasonics and builds a real-time 3D model of everything around the car - that's the same idea. We were doing that kind of work, but in different contexts, for different purposes.

And here's the thing that hit me: this capability was locked away.

It was locked in expensive proprietary systems. It was locked behind classification levels. It was locked in organizations with massive budgets. The technology existed. The techniques existed. But they weren't available to anyone who wasn't already inside the walls.

I kept thinking: why can't architects have this? Why can't robotics engineers have this? Why can't factory planners or city designers or anyone else working in three dimensions have access to this kind of spatial reasoning?

The answer wasn't technical. It was just that nobody had built it for them yet.

---

## Then Came Building Systems in Germany

Later I worked on building facility management, IoT, and AI systems in Germany.

This is where the other piece clicked into place. Because buildings aren't just designs that get built and then you're done. Buildings are living systems. They have sensors. They have HVAC that needs to respond to conditions. They have occupancy patterns that change. They have maintenance schedules. They have energy consumption that needs to be optimized.

A building exists in the real world, and it keeps existing in the real world for decades. The design phase is just the beginning. After that comes operations - endless operations, maintenance, upgrades, renovations.

And all the same problems showed up again. The IoT sensors generating data that needed to be fused into a coherent picture. The AI systems trying to optimize operations but struggling because the data was scattered across incompatible systems. The facility managers working with drawings that didn't match what had actually been built because someone made changes during construction and never updated the documentation.

The digital twin concept makes sense - have a digital model that stays synchronized with the physical reality. But nobody had actually built the infrastructure to make that work properly. The tools didn't talk to each other. The data formats were incompatible. The real-time requirements were ignored.

This is what any real system faces. It's not just design and then done. It's design, build, operate, maintain, adapt, renovate, operate some more. The spatial model needs to stay alive and synchronized with reality for the entire lifetime of the thing being modeled.

---

## Prefab and the Factory Problem in Sweden

Then I worked in Sweden with large design-and-construct companies that wanted to vertically integrate. The idea was prefab modular construction - build components in factories, ship them to site, assemble them quickly. Faster, cheaper, better quality control than traditional construction.

The concept is brilliant. The execution is brutal.

Because now you've got three worlds that have to talk to each other perfectly: design, CAD, and factory.

The architects design something. The engineers turn it into detailed CAD models. The factory has to actually manufacture the components. And every single dimension, every tolerance, every connection point has to match up exactly. If the design says one thing and the factory builds another, you've got components that don't fit together on site. You've got trucks shipping the wrong things. You've got crews standing around waiting for parts that were built to the wrong spec.

The companies I worked with were trying to build mini factories - smaller, more flexible manufacturing facilities that could produce building components. The vision was beautiful. The reality was constant pain at the handoffs. Design changes that didn't propagate to the factory floor. Factory constraints that didn't feed back into design. CAD models that were technically correct but didn't account for how things actually get built.

This is where I really saw the cost of disconnected systems. Not just inconvenience. Actual money. Actual delays. Actual buildings that couldn't be assembled because the pieces didn't fit.

The vertical integration they wanted - design through manufacturing through assembly - required a level of data synchronization that simply didn't exist in the tools available. Everyone was working in their own system, exporting files, hoping the translations were accurate, discovering problems too late.

---

## The Document Problem

But there was another thread running through my career, and it took me a while to see how it connected.

In government work - and honestly in every large organization I ever worked with - there was always the document problem.

You have content that needs to go out through multiple channels. The same information needs to appear on a website, in a PDF, in a printed form, maybe in an app. And it needs to look exactly right in every single one of those outputs. Branding has to be perfect. Legal language has to be precise. One wrong word and you've got a compliance problem.

Then add languages.

Now you're publishing to five countries, ten countries, more. Different languages, different writing systems, different cultural expectations. The same content, translated, but still needing to be perfectly consistent with the original.

Then add time zones.

You've got teams in different countries working on the same documents. Someone in London makes a change. Someone in Singapore needs to see it. Someone in New York needs to approve it. And they're all working at different times, often passing changes back and forth without ever being online simultaneously.

When this goes wrong - and it goes wrong all the time - it's not just annoying. It's expensive and dangerous.

Think about any decent engineering project. You've got architects working on designs. Engineers doing structural calculations. Suppliers quoting on materials. Stakeholders reviewing progress. Planning departments requiring compliance documents. Finance tracking costs. Legal reviewing contracts.

All of these people need to work from the same information. And they're all in different companies, different countries, different time zones. Speaking different languages. Using different tools.

Someone in London updates a specification. Someone in Singapore needs to see it but they're asleep. Someone in New York approved the old version yesterday. The supplier in Germany already started manufacturing based on last week's drawings. The planning submission uses a document from three versions ago.

In engineering documents, a translation error can mean something gets built wrong. In government documents, an inconsistency can mean legal exposure. In medical or financial documents, mistakes can genuinely hurt people.

The pain is real. I've watched teams struggle with it for years. Spreadsheets tracking versions. Email chains trying to coordinate changes. Review processes that add weeks to timelines. Meetings just to figure out who has the latest file. And still things slip through. Still mistakes get made. Still projects get delayed and budgets get blown because information didn't flow properly.

---

## The Same Problem

Here's what I eventually realized: Spatial and Publish are the same problem.

They don't look the same on the surface. One is about 3D geometry and the other is about documents. But underneath, they share the same fundamental challenge.

**Complex outputs, multiple channels, distributed teams, zero room for error.**

In Spatial, you're taking data from multiple sources - CAD models, sensor feeds, simulations, whatever - and rendering it into a coherent 3D view. Teams across different time zones need to collaborate on it in real time. If things get out of sync, a robot arm doesn't go where it should. A building component doesn't fit. A simulation gives you wrong answers.

In Publish, you're taking content from a single source and pushing it out to multiple outputs - web, PDF, print, forms. Teams across different time zones need to collaborate on it in real time. If things get out of sync, a legal document has the wrong clause. A translated manual gives dangerous instructions. A form collects data that doesn't match the database.

Same pattern. Same pain. Same need for a single source of truth that stays synchronized across every output and every collaborator.

---

## Why Now

The technology finally caught up. CRDTs make real-time collaboration work without merge conflicts. AI can actually reason about spatial data, not just text. Open standards like STEP and IFC won the format wars. All the pieces existed separately - what didn't exist was something that put them together.

The name Ubuntu comes from the philosophy: "I am because we are." Software doesn't exist in isolation. It exists in the connections between people, between systems, between design and reality. Everything I've built has been about making those connections work.

---

## What We're Actually Building

The system is vector-based at its core. Massively scalable vectors that can represent anything from a logo to a building to a city. It's CAD and desktop publishing unified - because they're the same problem at different scales.

**Spatial** is for anyone working in three dimensions - robotics, simulation, digital twins, manufacturing, construction. It's a CAD system. It's a game engine. It's a video editor. It's all of them, because they're all just compositing geometry in space and time. The only difference is what you're compositing for.

**Publish** is for anyone who needs documents to work - government, healthcare, legal, financial, education. Write once, output everywhere. Translate without losing your mind. Keep distributed teams synchronized without the version control nightmare. Vector graphics that scale from a business card to a billboard without losing a single detail.

Both are built on the same foundation. Both handle the same fundamental problem: keeping complex information synchronized across teams, tools, channels, and languages.

Both are built offline-first because real work happens in real places, not just in offices with perfect internet. Both use open standards because lock-in is a trap. Both have AI built in from the start because that's the only way it actually helps instead of getting in the way.

---

## Where We Are

This is early. The foundation is built. There's a long way to go.

If you've felt this pain - spatial tools that don't collaborate, documents that won't stay synchronized, systems that should talk to each other but don't - you know what we're building and why it matters.

We're looking for people who see the same problems. Early customers who want to shape the direction. Partners who've hit the same walls. Engineers who want to build something that actually solves this.

[Get in Touch →](/contact/)

---

## More About Us

[Experience](/company/experience/) — 25 years of mission-critical systems for global enterprises.

[Founder](/company/founder/) — Meet the person behind Ubuntu Software.

[Advisors](/company/advisors/) — Industry leaders who collaborate with us.
