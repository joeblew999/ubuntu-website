# todo

---

can we have a ci task that checks that CI runs across a matrix ? I dont want to run everything across all OS. I just want to make sure the task files work on al OS within CI. this is to make sure out task file works for end users on different OS. some will be on Windows .

---

READY to split task files and then put them into .github repo ? 

so then all repos can pull versions tasks files at runtime as task can do thsi.

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
