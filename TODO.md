# todo

The refactoring caused this at startup.... weird that we are calling Cloudflare at startup. this should not happen.

```sh
ubuntu-website % go run cmd/env/main.go web-gui
2025/11/18 14:08:07 
2025/11/18 14:08:07 Environment Setup GUI
2025/11/18 14:08:07 Opening in browser...
2025/11/18 14:08:07 
  http://localhost:3000

2025/11/18 14:08:07 Press Ctrl+C to stop
2025/11/18 14:08:07 
2025/11/18 14:08:07 Failed to fetch zones: failed to fetch zones: Get "https://api.cloudflare.com/client/v4/zones?per_page=50": dial tcp: lookup api.cloudflare.com: no such host
2025/11/18 14:08:07 Failed to fetch Pages projects: failed to fetch Pages projects: Get "https://api.cloudflare.com/client/v4/accounts/7384af54e33b8a54ff240371ea368440/pages/projects": dial tcp: lookup api.cloudflare.com: no such host
^Csignal: interrupt
```

---

When we get errors we dont know what page they occured on. We need via to have that .

```sh
ubuntu-website % go run cmd/env/main.go web-gui
2025/11/18 13:01:36 
2025/11/18 13:01:36 Environment Setup GUI
2025/11/18 13:01:36 Opening in browser...
2025/11/18 13:01:36 
  http://localhost:3000

2025/11/18 13:01:36 Press Ctrl+C to stop
2025/11/18 13:01:36 
2025/11/18 13:02:25 [error] msg="failed to handle session close: ctx '/_/dd297b1f' not found"
^Csignal: interrupt
apple@apples-MacBook-Pro ubuntu-website % go run cmd/env/main.go web-gui
2025/11/18 13:39:59 
2025/11/18 13:39:59 Environment Setup GUI
2025/11/18 13:39:59 Opening in browser...
2025/11/18 13:39:59 
  http://localhost:3000

2025/11/18 13:39:59 Press Ctrl+C to stop
2025/11/18 13:39:59 
2025/11/18 13:40:43 [error] via-ctx="/_/1610b467" msg="PatchElements failed: failed to send elements: context cancelled: context canceled"
^Csignal: interrupt
```

---

Make is so that we can use claude locally using our existing claude code cli, so that we can do the transaltion locally.

---

Make the local hugo use certs and then also publish the LAN URL too, so we can test from mobile. Mkcert can gen the certs for us. We do not need caddy complexity.



