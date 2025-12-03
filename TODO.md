# todo

---

BLUESKY_APP_PASSWORD 

I was lazy and did not put it into .env !!!

---

logs folder is big and not gitignored ? is useful though to keep around ?

---

narun

https://github.com/akhenakh/narun

TO help buiuld out the stuff we are building into cmd to be genertic ?

---

Adjust static folder with the stuff i need to do to update Bluessky and GMail each time my logo and Strap line changes.

Is there a better way ? We have golang code and task and claude.md.

---

cmd/env is broken - DO NOT RUN

The env tool expects a specific .env structure and will corrupt your .env file.

Investigate https://github.com/helmfile/vals for proper secrets management.
vals supports multiple backends (Vault, AWS SSM, GCP Secrets, SOPS, etc.)

---

cmd/translate is broken - DO NOT RUN

Current workflow: Claude Code + Taskfile shell scripts (translate:status, translate:done, etc.)

The manual workflow is prototyping what cmd/translate should eventually do:
- Learning what prompts work best for translation quality
- Understanding file structure requirements
- Handling edge cases (missing files, orphaned translations)
- Defining the right workflow (status → translate → done checkpoint)

Once patterns are solid from manual use, codify them into cmd/translate.
