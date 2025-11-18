# todo


Make the local hugo use certs and then also publish the LAN URL too, so we can test from mobile. Mkcert can gen the certs for us. We do not need caddy complexity. We can make it gen the certs everytime to make the code simpler ? 

- Leave for now as too hard and not important


---

The Build and deploy with wrangler told me its at: https://f971d136.bbb-4ha.pages.dev

I dont know if you can easily display https://f971d136.bbb-4ha.pages.dev as a hyperlink like we do for local builds ?


https://www.ubuntusoftware.net has some errors. I dont know if we can have a promote to the Domain as part of the http://localhost:3000/deploy page or how it weorks.



---

We have deps like mkcert, hugo and wrangler.

We need a robust way to pin versions.
We need a robust way to update them if they dont match the versions. This is often tricky on devs laptops as wel as when depploye dto a server i have found.

---

Make it so that we can use claude locally using our existing claude code cli and login. We will use this for our translation system. This is code we have not worked on much in terms of the backend and the web gui. We really just want to make sure that we can call claude locally for now and that it works with use passing in a command to it.
YOu might have issues locating where claude code is installed to because we are usng vscode, and that, i think installed the binary somewhere that is not common.. 





