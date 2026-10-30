# iec104ctl Releases

## v0.3.0

**Date:** unreleased
**Previous release:** v0.2.0

### Summary

**Breaking for scripts:** the exit codes follow the OT Fabric organisation table (ADR-0006 of otfabric/engineering). Nothing else changes in the commands.

### Exit codes

| Case | Old | New | JSON error |
|---|---|---|---|
| `context_canceled` (interrupted before a result) | 10 | 130 | kind `interrupted`, was `internal` |

Every other exit code is unchanged. See [docs/exitcodes.md](docs/exitcodes.md).

### Environment

`IEC104CTL_FORMAT` is now the canonical name of the format variable; `IEC104CTL_OUTPUT_FORMAT` stays accepted, and `IEC104CTL_FORMAT` wins when both are set.

---

## v0.2.0

**Date:** 2026-10-26  
**Previous release:** v0.1.0

### Summary

The application-layer security of **IEC 62351-5**, the A-Profile of IEC TS 60870-5-7, on every client command and on the simulated station. After STARTDT both ends agree on keys in a Station Association authenticated with certificates, and from then on every ASDU travels as Secure Data, authenticated with HMAC-SHA-256 or encrypted with AES-256-GCM. Built on [go-iec104](https://github.com/otfabric/go-iec104) v0.0.4 (was v0.0.3), whose package `sec` implements it.

### Highlights

- **`--secure`:** On every client command, through the common connection flags. **`--security-cert`** and **`--security-key`** are the station's own certificate and key (P-256, or X25519 signed by an authority); **`--security-trust`** holds the stations to trust: a CA certificate is an authority whose station certificates are accepted, any other certificate is a pinned self-signed station certificate. **`--association-id`** (the AIM, default 1), **`--mal`** (MAC algorithm, default HMAC-SHA-256) and **`--dpa`** (data protection, e.g. **`aes-256-gcm`**; default `--mal`), by code or by name. The command starts once the session is secure. Combines with TLS.
- **`server start --secure`:** The same certificate, key and trust flags; **`--association-id`** is the AIS. Every controlling station (AIM) gets its own association. **`--roles viewer,operator,…`** (the roles of IEC TS 60870-5-7) applies role-based access control: a request no role permits is refused with cause 17 (`not-authorized`). New events **`secured`** and **`security-refused`**, and **`refused`** for a request the roles do not permit.
- **Error code `security_failed`:** Kind `transport`, exit code **4**: a refused certificate, an algorithm the station does not accept, a MAC that did not verify, or no answer to the key management within `--timeout`.
- **Validation first:** A security flag without `--secure`, `--secure` without its three files, a file that cannot be read or used, a non-CA certificate in `--security-trust` that is not self-signed, or an unknown algorithm or role exits **3** before anything is sent or the server listens.
- **`--debug`:** Traces the security events (`STAS_PROC_SUCC`, `NODE_NOT_AUTR`, …) and the algorithms of the session along with the frames.
- **Environment and completion:** Every new flag has an **`IEC104CTL_*`** variable (`IEC104CTL_SECURE`, `IEC104CTL_SECURITY_CERT`, `IEC104CTL_SECURITY_KEY`, `IEC104CTL_SECURITY_TRUST`, `IEC104CTL_ASSOCIATION_ID`, `IEC104CTL_MAL`, `IEC104CTL_DPA`, `IEC104CTL_ROLES`); `--mal`, `--dpa` and `--roles` complete to their values, the key material to PEM files.
- **Docs:** A security section in the README with a worked openssl example for pinned self-signed certificates and for a plant authority; **`docs/exitcodes.md`**, **`docs/json-output.md`** and **`SECURITY.md`** updated.

### Limitations

- The association is kept in memory: every run of `iec104ctl` associates anew, and the simulated station forgets its associations when it stops.
- The simulated station gives every controlling station the same trust and roles; per-AIM trust and roles are not configurable.
- The timers and key usage limits of IEC 62351-5 are the library's defaults and have no flags.
- secp256k1, Curve448, brainpoolP256r1 and BLAKE2s are not supported (go-iec104 has no implementation of them).
- The conformance of go-iec104's implementation to the standards, which are not public, is not established independently; it is tested for interoperability with cB-Sigma/iec60870.

### Quality

- **Tests:** Client commands against the simulated station with authority and pinned trust, HMAC-SHA-256, HMAC-SHA3-256 and AES-256-GCM, other association IDs; a viewer refused a command with cause 17 while an operator may; untrusted certificates on either side (`security_failed`); a client without security against a secure station (timeout) and security against a station without it; `server start --secure` end to end with its events; every validation error exits 3. Whole-module coverage stays above the floor of 85%.

---

## v0.1.0

**Date:** 2026-10-10  
**Previous release:** none

### Summary

First release. iec104ctl is a command-line tool for **IEC 60870-5-104** in both roles: a **controlling station** that interrogates, monitors and commands a station and transfers files, and a **simulated controlled station** served from a model file. It also finds stations on a network. Built on [go-iec104](https://github.com/otfabric/go-iec104) v0.0.3.

### Highlights

- **`client interrogate`:** General and group interrogation (`C_IC_NA_1`). Every value with quality, cause of transmission and time tag; **`--type`** and **`--ioa`** filter what is shown.
- **`client counters`:** Counter interrogation (`C_CI_NA_1`), general or per counter group, with **`--freeze read|freeze|freeze-reset|reset`**.
- **`client read`:** Read command (`C_RD_NA_1`) for a list or range of addresses. A refused address is reported in its row and the others are still read (exit code **7**).
- **`client monitor`:** Stays connected and prints every value as it arrives, with the time it arrived and its time tag. **`--interrogate`**, **`--ca`**, **`--type`**, **`--ioa`**, **`--count`**, **`--duration`**, **`--all`**; reconnects unless **`--no-reconnect`**.
- **`client command`:** Single, double, regulating step, set point (normalized, scaled, short floating point) and bitstring commands. **`--mode execute|sbo|select|cancel`**, **`--qualifier`**, **`--time-tag`** (the CP56Time2a variant), **`--wait-termination`**. Everything is validated before anything is sent, and a command is sent exactly once.
- **`client clock`** and **`client reset`:** Clock synchronization (`C_CS_NA_1`) with **`--time`** and **`--timezone`**; reset process command (`C_RP_NA_1`) with a required **`--qualifier process|events`**.
- **`client ping`:** Link test frames (TESTFR) or test commands (`C_TS_TA_1`) at an interval, with round-trip times and a summary.
- **`client scan`:** Finds the common addresses behind one connection: **`--method broadcast`** (one interrogation to address 65535) or **`sweep`**.
- **`client file list|get|put`:** File directory, download and upload, with the checksums of every section and of the file verified and the SHA-256 of the content in the result.
- **`client types`:** The type identifications iec104ctl knows, with direction and description.
- **`discover`:** Sweeps subnets for IEC 104 stations with a single link test frame per address; data transfer is not started. Safety cap of 4096 hosts (**`--force-large-scan`**).
- **`server start`:** A simulated controlled station from a YAML or JSON model, or the built-in demo. Interrogations, read, counters with freeze and reset, commands direct or select-before-operate, clock, test command, reset, files and uploads; values that change by themselves (**`toggle`**, **`ramp`**, **`sine`**, **`random`**, **`count`**), cyclic measured values, queued events, redundant connections; TLS, **`--allow`**, **`--max-sessions`**, **`--duration`**. Every event is logged to stdout.
- **`server generate-config`:** Writes the commented demo model, as YAML or JSON.
- **Transports:** TCP and TLS (IEC 62351-3) with **`--url tcp://…|tls://…`** or **`--ip`** / **`--port`**; **`--tls-ca`**, **`--tls-cert`**, **`--tls-key`**, **`--tls-server-name`**, **`--tls-insecure`**.
- **Protocol parameters:** **`--k`**, **`--w`**, **`--t0`** to **`--t3`**, **`--originator`**, and the ASDU field sizes **`--cot-size`**, **`--ca-size`**, **`--ioa-size`** for stations behind IEC 60870-5-101 gateways.
- **Output:** **`--format text|json|table`** on every command; streams write JSON Lines. A single JSON error object on failure with `--format json`, and stable exit codes (0, 2, 3, 4, 5, 6, 7, 10). The JSON shapes and exit codes are documented in **`docs/json-output.md`** and **`docs/exitcodes.md`**.
- **Environment:** Every flag has an **`IEC104CTL_*`** variable; a flag wins over the environment.
- **`--debug`:** Traces every frame sent and received, with the ASDU decoded and in hex, and shows the causes of an error.
- **Completion:** bash, zsh, fish and PowerShell; flags with a fixed set of values complete to those values.

### Limitations

- The application-layer authentication of IEC 62351-5 is not implemented; its ASDU types are recognized by name.
- No serial link layer: IEC 60870-5-101 itself is not spoken, only its field sizes over IEC 104.
- The simulated station serves one common address per instance.
- `client command` has no parameter commands (`P_ME_…`, `P_AC_NA_1`); parameters a station sends are decoded and shown.

### Quality

- **Tests:** Whole-module statement coverage of **94%**, with a floor of 85% enforced by **`make coverage`** (`COVERAGE_MIN`). Every command is tested through the real command tree against an in-process simulated station (**`internal/stationtest`**), over TCP and TLS, in all three output formats, with its exit codes.
- **Interoperability:** The client commands were run against the servers of lib60870-C, j60870 and wendy512/iec104, and the simulated station against their clients (interrogation, counters, read, clock, test command, commands direct / select-before-operate / cancel, file download).
- **Hardware-free:** The test suite needs no station and no network beyond the loopback interface.

---
