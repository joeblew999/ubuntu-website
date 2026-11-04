# FLEET stack

## DRONES

### MAVLink in Go (foundation layer)
gomavlib
Repo: https://github.com/bluenviron/gomavlib
Why it matters:
Clean, production-grade MAVLink implementation
PX4 & ArduPilot compatible
Works over:
Serial
UDP
TCP
Used for:
Parameter push/pull
Health monitoring
Mission upload
Command & control
This is the core control API for Go-based provisioning.


### PX4 parameter automation (critical for zero-touch)
px4-params (generated from PX4 metadata)
Repo: https://github.com/PX4/px4_msgs (indirect)
Used via: gomavlib + generated bindings
Pattern used in real systems:
Versioned .params files
Applied programmatically
No QGroundControl required
This is how you eliminate snowflake drones.

### ingest the RTSP stream from a camera mounted on a Holybro PX4

REPO: https://github.com/bluenviron/mediamtx

mediamtx sits between your camera hardware and the rest of your infrastructure, managing streams efficiently.

1. Real-Time Video Streaming
Mapping drones often need live video streams:
Telemetry overlay
Operator monitoring
Payload verification
mediamtx can ingest the RTSP stream from a camera mounted on a Holybro PX4 build
It can then:
Distribute to ground stations
Send to UI dashboards
Store it for post-flight review
This is exactly the kind of infrastructure DJI has built into its ecosystem, but you’ll implement it yourself.

2. Payload Camera Data Capture
Whether you’re:
Taking pictures for SfM (Structure from Motion)
Streaming FPV back to base
Using thermal / multispectral feeds
mediamtx gives you:
✔ Centralized media pipeline
✔ Recording + playback
✔ Bandwidth management
✔ No reliance on 3rd-party servers


