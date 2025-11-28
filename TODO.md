# todo




---

~~Make sure the claude.md has a note about makign sure the OG is right too via the frontmatter . assume site map naturally follows ?~~ DONE - See CLAUDE.md Bluesky section 



---

We need real names for testimonials ? 

We need the testimonals and the previsu proejcts to link well and i need to make the previous projects things good. The Climate Foundation system was a secure publisging system, data collection system and messaging and video conferencing system, so that all global projects had a unified and secure and self soverign way to collaborate and publish and for sciwnce divisiosn to collect sensor information.

a blog article about it too is good.

a testimonial about it too from the director.


--- 

In the Ourstory we need to say how we spe years builting Systsms and Apps for many years and that the company name existed before Ubuntu Linux happened, but we built for Linux and Open BSD for many many years. 

---

Now we need to add Images. But first work out best way. SVG since its easy for  CLaude to make off the text ? And then no scaling issues.

Need to consider colour. We use gray at the moment because its neutral.

Need to consider hugo plate and not working against it, and to ensure things keep working well when we update, in relation to images. 

Hugo plate elements page has some Nice Image tools that we can leverage perhaps too.

---

Some layouts are not used ? Can hugo tell us ? or do we even care ? or do we need a task file to help us. DO we really care ? its only EN we care about as its the single source of truth when checking what conrtent uses what layouts ? 


--

testomonials is  done ? get one from Dick of archethought.
reference his github too somwwhre in the site ? 

---

blogs need to be adjusted ? We can make some about the problems in the industry related to our product. Keep the ones about programming ?? 


---

Adjust static folder with the stuff i need to do to update Bluessky and GMail each time my logo and Strap line changes.



---

~~Contact page is stupid as we dont have a forms processor and i dont want to ad one. Better to have an email. maybe there is a decent obfuscation thing buuld in to Hugo and hugop plate to help ?~~ DONE - Using Web3Forms (free, unlimited, spam-protected) 


---

Pictures 

- we have default Grey ones right now that hugo plate defaults I dont want to make a big mess, but maybe we can easily add images ? 




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
