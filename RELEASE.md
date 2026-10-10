# iec104ctl Releases

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
