# todo

you are going to need a cross platform golang tool for some of the things you doing maybe ? all task files might need to depend on it. For example i was using go-which  ( a golang installable thng ) so that the which thing was cross platform. I call these base tools because local and ci task files need them, and again they need to be idempotently checked and installed.    also the golang version of jq too.   so just like oyu have runtime , you have tools.  Also all this CI is going to cost me on my CI plan ?  Maybe i am over baking the cake too early too ?   Task is good but you need to file in the holes ? 

Its normal to recactor back to golang.

task is good and i plan to iuse it in production on servers and desktops with real users editing it.   this is why its worth getting this right .


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
