# DEEPLNK

## Code

Deeplink to work with GUI and also gioplugins WebView etc 

https://github.com/gioui/gio/pull/117

What I’d do for ubuntusoftware.net
Use custom scheme deep links for “inside the app” actions (works even without domain verification):
e.g. ubuntusoftware://open/doc/123
Use https deep links for marketing/email/web (https://links.ubuntusoftware.net/open/...) and bridge into the app via Universal/App Links.
If you tell me which Gio shell you’re shipping (pure Gio UI vs WebView shell) and your target platforms (iOS/Android/both), I’ll give you the exact deep-link wiring plan you should implement (including scheme names that won’t collide, and how to map incoming URLs to your router).

https://github.com/gioui/gio-cmd has stuff from him now

https://github.com/gioui/gio/commit/7bcb315ee174467d8e51c214ee434c41207cc718 has the deeplink stuff now !!

https://github.com/gioui/gio-cmd/pull/23 allows to see what apps are installed !

https://github.com/gioui-plugins/gio-plugins has the WebView and other stuff we need, but maybe is not compatatable yet ...

The current tested version is 1a17e9ea3725cf5bcb8bdd363e8c6b310669e936 (I think you can get it by go get gioui.org@1a17e9ea3725cf5bcb8bdd363e8c6b310669e936). That is after 0.8.0 tag (and before the current latest commit). Also, get the last version of Gio-Plugins (go get github.com/gioui-plugins/gio-plugins@main)

https://github.com/richzw/appstore useful ?

https://github.com/awa/go-iap useful ? go-iap verifies the purchase receipt via AppStore, GooglePlayStore, AmazonAppStore and Huawei HMS.

ALL THIS needs to be done in the goup-util repo first !!

## SETUP

Great — this is a *real* milestone. Let’s treat it as **production-grade App Store readiness**, not hand-wavy “hello world” deployment.

I’ll structure this so you can **actually ship** with `ubuntusoftware.net` as your canonical domain, with **deep linking working across iOS, Android, and web**.

---

## 1. Decide your app-store identity (do this first)

You want **one coherent identity across stores, web, and links**.

### Company / Publisher

* **Legal name**: Ubuntu Software (or whatever is on your incorporation docs)
* **Domain**: `ubuntusoftware.net` ✅
* **Support email**: `support@ubuntusoftware.net`
* **Developer contact**: `gerard.webb@ubuntusoftware.net`

You’ll reuse this **everywhere**:

* App Store Connect
* Google Play Console
* Apple Sign-in
* Universal Links / App Links
* OAuth callbacks

---

## 2. Apple App Store (iOS, iPadOS, macOS)

### Apple Developer Program

You must enroll as **Organization**, not Individual.

**Why Organization matters**

* Required for Universal Links
* Required for Apple Sign In
* Company name appears in App Store
* Proper entitlement handling

**What you need**

* D-U-N-S number (Apple verifies this)
* Legal entity docs
* `ubuntusoftware.net` domain control

---

### App ID & Bundle ID strategy (important)

Use **reverse-DNS tied to your domain**:

```
net.ubuntusoftware.appname
```

Examples:

```
net.ubuntusoftware.cad
net.ubuntusoftware.files
net.ubuntusoftware.control
```

This bundle ID **must never change** once shipped.

---

## 3. Google Play Store (Android)

### Google Play Console

* One-time fee
* Company profile using `ubuntusoftware.net`

### Android Application ID

Use the **same reverse domain**:

```
net.ubuntusoftware.appname
```

Apple + Google IDs **should match** (except platform specifics).

---

## 4. Deep linking architecture (this is the key part)

You should **not** use custom schemes alone (`myapp://`).
You want **real web-backed deep links**.

### Canonical deep link domain

Use a dedicated subdomain:

```
https://links.ubuntusoftware.net
```

Why:

* Clean separation
* Easier CDN / routing
* Safer long-term

---

## 5. iOS: Universal Links (production-grade)

### URL pattern

```
https://links.ubuntusoftware.net/open/...
```

Examples:

```
https://links.ubuntusoftware.net/open/doc/123
https://links.ubuntusoftware.net/open/project/abc
```

### Required file (hosted at root)

```
https://links.ubuntusoftware.net/.well-known/apple-app-site-association
```

**Example**

```json
{
  "applinks": {
    "apps": [],
    "details": [
      {
        "appID": "TEAMID.net.ubuntusoftware.appname",
        "paths": [ "/open/*" ]
      }
    ]
  }
}
```

Notes:

* No `.json` extension
* Served with `application/json`
* Must be HTTPS
* Apple caches aggressively (plan ahead)

---

## 6. Android: App Links

### Required file

```
https://links.ubuntusoftware.net/.well-known/assetlinks.json
```

**Example**

```json
[
  {
    "relation": ["delegate_permission/common.handle_all_urls"],
    "target": {
      "namespace": "android_app",
      "package_name": "net.ubuntusoftware.appname",
      "sha256_cert_fingerprints": [
        "AA:BB:CC:..."
      ]
    }
  }
]
```

This binds:

* Your **domain**
* Your **app**
* Your **signing key**

---

## 7. Web fallback behavior (don’t skip this)

Your link handler service should:

### Behavior matrix

| Situation         | Result                             |
| ----------------- | ---------------------------------- |
| App installed     | Open app directly                  |
| App not installed | Redirect to App Store / Play Store |
| Desktop browser   | Open web UI                        |
| Bot / crawler     | Serve metadata                     |

### Example flow

```
links.ubuntusoftware.net/open/doc/123
 ├─ iOS + app installed → app opens
 ├─ Android + app installed → app opens
 ├─ Mobile no app → store page
 └─ Desktop → https://app.ubuntusoftware.net/doc/123
```

---

## 8. How this fits your Go + WebView setup

Since you mentioned:

* **WebView**
* **Go backend**
* **Datastar / SSE**
* **NATS control plane**

### Best practice for you

* WebView loads:

  ```
  https://app.ubuntusoftware.net
  ```
* Native shell handles:

  * Universal/App Links
  * File providers
  * Auth handoff
  * Push notifications

Deep link arrives → native layer → inject URL into WebView.

This keeps **one UI**, many shells.

---

## 9. Apple Sign-in & OAuth alignment (important)

When you enable **Sign in with Apple**:

* Domain must be verified
* Callback URLs must be HTTPS

Use:

```
https://auth.ubuntusoftware.net/callback/apple
```

Same for:

* Google
* GitHub
* Your own auth

---

## 10. Store submission checklist (condensed)

### Required before submission

* Privacy Policy
  `https://ubuntusoftware.net/privacy`
* Terms
  `https://ubuntusoftware.net/terms`
* Support page
  `https://ubuntusoftware.net/support`
* App screenshots
* App icon (1024×1024)

### Apple-specific

* Export compliance (likely “No encryption” or “standard”)
* App Tracking Transparency (if applicable)
* Sign in with Apple (mandatory if you have other sign-ins)

---

## 11. Recommended domain layout

```
ubuntusoftware.net
├─ www.ubuntusoftware.net        (marketing)
├─ app.ubuntusoftware.net        (web app)
├─ links.ubuntusoftware.net      (deep links)
├─ auth.ubuntusoftware.net       (OAuth)
├─ api.ubuntusoftware.net        (Go APIs)
└─ status.ubuntusoftware.net
```

This scales cleanly for years.

---

## 12. What I can do next for you

If you want, I can:

1. **Generate**

   * `apple-app-site-association`
   * `assetlinks.json`
2. **Design**

   * Deep-link routing table
3. **Provide**

   * iOS Universal Link handler code
   * Android intent filter setup
4. **Create**

   * App Store & Play Store submission checklist tailored to your app
5. **Review**

   * Your current WebView + Go setup for store compliance

Just tell me **which app first** and **which platforms** (iOS / Android / both).
