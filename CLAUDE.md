# Claude Assistant Notes

## CRITICAL REMINDERS

USE TASKFILE - it makes conventions for development...

### Branding Assets

Source: `cmd/genlogo/main.go` → `task generate:assets`

After regenerating, manually update: Bluesky, Gmail signature

### Translation Workflow

Tasks:
- `task translate:status` - what English files changed since last translation
- `task translate:missing` - which languages are missing content files
- `task translate:done` - mark translations complete (updates checkpoint)

Workflow:
1. `task translate:status` → translate changed files to all 5 languages
2. `task translate:done` → update checkpoint

Languages: de (German), sv (Swedish), zh (Chinese), ja (Japanese), th (Thai)

### Path Convention
**ALWAYS use `joeblew999` (with three 9s), NEVER `joeblew99` (with two 9s)**

Correct paths:
- `/Users/apple/workspace/go/src/github.com/joeblew999/ubuntu-website`
- `github.com/joeblew999/ubuntu-website`

### Via Framework

When working on Via coding aspects:

- **Source code**: `/Users/apple/workspace/go/src/github.com/joeblew999/wellknown/.src/via`

- use a closure variable pattern instead of trying to access the signal's value.
- use PicoCSS

### Hugo Source code

When working on Hugo coding aspects use the source at:

/Users/apple/workspace/go/src/github.com/joeblew999/wellknown/.src/hugo

### Hugo Plate Source code

When working on Hugo Plate, which is a Hugo theme we use coding aspects use the source at:

/Users/apple/workspace/go/src/github.com/joeblew999/wellknown/.src/hugoplate

Only use Hugo template things properly. I do not want to steer away from the standard Hugo Plate way of doing things !! 

## Information about the company 

This proejct is a web site for my company. I have various info here about me and the company.

/Users/apple/Library/Mobile Documents/com~apple~CloudDocs/Thailand /dtv-employment-proof
