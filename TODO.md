# todo

---

pages use the old Via API with via.Context ? 

---

IOS Mobile trsut Issue.

- Leave for now as too hard and not important


---

DNS web step.

https://www.ubuntusoftware.net has some errors. I dont know if we can have a promote to the Domain as part of the http://localhost:3000/deploy page or how it weorks.

--- 

Versions

I noticed that wrangler gives us versions for each deploy. We can use this in our deploy page in a smart way ?

wrangler 4.22.0 (update available 4.47.0)
─────────────────────────────────────────────
┌──────────────────────────────────────┬─────────────┬────────┬─────────┬────────────────────────────────────┬────────────────┬──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Id                                   │ Environment │ Branch │ Source  │ Deployment                         │ Status         │ Build                                                                                                            │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ dc2dc9db-d13f-424c-aa5d-314b9622e9e1 │ Production  │ main   │ 9ad76a7 │ https://dc2dc9db.bbb-4ha.pages.dev │ 3 minutes ago  │ https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/pages/view/bbb/dc2dc9db-d13f-424c-aa5d-314b9622e9e1 │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ d7069f3a-3daa-4ea4-9a15-a81c1c3245a2 │ Production  │ main   │ 9ad76a7 │ https://d7069f3a.bbb-4ha.pages.dev │ 4 minutes ago  │ https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/pages/view/bbb/d7069f3a-3daa-4ea4-9a15-a81c1c3245a2 │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ a644bb65-8198-470e-9002-4a38b0e51d34 │ Production  │ main   │ 9ad76a7 │ https://a644bb65.bbb-4ha.pages.dev │ 4 minutes ago  │ https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/pages/view/bbb/a644bb65-8198-470e-9002-4a38b0e51d34 │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ 6413decb-a211-4896-8dc7-ab21ad34eea6 │ Production  │ main   │ 9ad76a7 │ https://6413decb.bbb-4ha.pages.dev │ 7 minutes ago  │ https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/pages/view/bbb/6413decb-a211-4896-8dc7-ab21ad34eea6 │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ d59e4628-b6f0-43dc-b446-6f89c5276133 │ Production  │ main   │ 9ad76a7 │ https://d59e4628.bbb-4ha.pages.dev │ 7 minutes ago  │ https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/pages/view/bbb/d59e4628-b6f0-43dc-b446-6f89c5276133 │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ a8fa6fee-c103-482b-b02a-229a452df8a0 │ Production  │ main   │ 9ad76a7 │ https://a8fa6fee.bbb-4ha.pages.dev │ 7 minutes ago  │ https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/pages/view/bbb/a8fa6fee-c103-482b-b02a-229a452df8a0 │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ 979b5dfd-588c-4a3c-9745-a6697d2d3827 │ Production  │ main   │ 9ad76a7 │ https://979b5dfd.bbb-4ha.pages.dev │ 7 minutes ago  │ https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/pages/view/bbb/979b5dfd-588c-4a3c-9745-a6697d2d3827 │
├──────────────────────────────────────┼─────────────┼────────┼─────────┼────────────────────────────────────┼────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤


---

We have deps like mkcert, hugo and wrangler.

We need a robust way to pin versions.
We need a robust way to update them if they dont match the versions. This is often tricky on devs laptops as wel as when depploye dto a server i have found.

---

Make it so that we can use claude locally using our existing claude code cli and login. We will use this for our translation system. This is code we have not worked on much in terms of the backend and the web gui. We really just want to make sure that we can call claude locally for now and that it works with use passing in a command to it.
YOu might have issues locating where claude code is installed to because we are usng vscode, and that, i think installed the binary somewhere that is not common.. 





