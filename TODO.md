# todo


---

The home page in English has a "Learn More" button / link, but on the other language pages when you click the same button / link, it goes to the English page for some reason. It should not. 

Yet the links on the blog home page do go to the same language blog pages. 

Work out whats wrong.

---

make sure Privacy page is decent. We need to make sure we are not breaking any laws in the EU, for example, and telling our users that we respect their rights to privacy and to not be tracked.   

I think that there is some law about take downs and they if we do not do it we are in trouble. But we are an Australain Comapny, so maybe that means we are OK there ? 

---

Option 2: Fix Custom Domain Setup (Technical, unblocks production)
Debug and fix the Cloudflare domain API integration (currently failing)

---


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
