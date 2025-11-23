# todo


Option 2: Fix Custom Domain Setup (Technical, unblocks production)
Debug and fix the Cloudflare domain API integration (currently failing)

---

CLAUDE: If your seeing this, you can just move anything done to the DONE area. KEEP IT KISS - the code is self documenting, so just move it and change anything you need to make it accurate but kiss.

## DONE

Preview URLs - FIXED by switching from HTTPS to HTTP-only development
- Removed mkcert self-signed certificate generation (~50 lines)
- Changed hugo.go StartHugoServer() to use HTTP instead of HTTPS
- Added --environment development to build command in BuildHugoSite()
- This ensures both server and build use config/development/config.toml with baseURL="/" and relativeURLs=true
- No more certificate trust issues on browsers or mobile devices
- Works instantly on all devices without manual CA installation
- ✅ Verified working from http://localhost:3000/deploy "Build site only"
🌐 Local Preview: http://localhost:1313
📱 LAN Preview (Mobile): http://192.168.1.49:1313

External Link Standardization - FIXED by creating reusable components
- Created RenderExternalLink() helper in components.go for standard "Visit: [Label] ↗" pattern
- Created RenderExternalLinkWithCustomPrefix() for custom prefix cases like "Add billing at:"
- Updated route_claude.go (4 links)
- Updated route_cloudflare_step1.go (1 link)
- Updated route_cloudflare_step2.go (1 link)
- Updated route_cloudflare_step3.go (1 link)
- Total: 7 external links standardized across all web GUI pages
- All external links now consistent with security best practices (target="_blank", rel="noopener noreferrer")
- ✅ Build verified successful


## NOT DONE 

---

DNS web step so that CUSTOM Domain works.

http://localhost:3000/cloudflare/step5 got me Failed to attach domain: failed to add domain ubuntusoftware.net (status: 400): { "result": null, "success": false, "errors": [ { "code": 8000006, "message": "Request body is incorrect. The request body may have an invalid JSON type or missing required keys. Refer to https://developers.cloudflare.com/api." } ], "messages": [] }


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
We need a robust way to update them if they dont match the versions. This is often tricky on devs laptops as well as when deployed to a server. We need a simple way for now using go install, and a local .bin perhaps ? so then it does not matter what is local on the os. BUT 

---

Make it so that we can use claude locally using our existing "claude code" cli and login. We will use this for our translation system. This is code we have not worked on much in terms of the backend and the web gui. We really just want to make sure that we can call claude locally for now and that it works with use passing in a command to it.
YOu might have issues locating where claude code is installed to because we are usng vscode, and that, i think installed the binary somewhere that is not common.. 

The translation system is currently designed to use the claude API. Maybe it can be adapted to call the claude code running locally. Just another way to skin the cat ..


---

Cleanup text used everywhere so it matches what we use at the env level.
I noticed the text in env web has a few inconsistencies. We want it to be really consistent with a source of truth.

---

cleanup cloudflare APi calls and Request and Response types ? Could be better ?
- Decided to leave as-is - code works fine and refactoring would add complexity without much benefit
- Current implementation is clear and functional 

## NOT VITAL

IOS Mobile trust Issue.

- Leave for now as too hard and not important
