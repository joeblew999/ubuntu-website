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
