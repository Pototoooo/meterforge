# MeterForge API specs

This workspace contains two TypeSpec packages that generate OpenAPI specs:

| Package                        | Description                                                    | Output                                                    |
| ------------------------------ | -------------------------------------------------------------- | --------------------------------------------------------- |
| **Legacy** (`packages/legacy`) | MeterForge API (v1-v2) and the legacy cloud-compatible surface | `openapi.MeterForge.yaml`, `openapi.MeterForgeCloud.yaml` |
| **AIP** (`packages/aip`)       | MeterForge and Konnect metering & billing APIs (v3), AIP-style | `openapi.MeteringAndBilling.yaml` (MeterForge + Konnect)  |

From the repo root, run `make gen-api` (or `make -C api/spec generate`) to build both packages and copy/bundle artifacts into `api/`.

---

## Legacy API (`packages/legacy`)

Legacy specs follow MeterForge’s existing TypeSpec conventions. See [`packages/legacy/README.md`](packages/legacy/README.md) for patterns and guidelines.

---

## AIP (`packages/aip`)

The AIP package defines v3 metering and billing APIs in line with [Kong’s AIP (API Improvement Proposals)](https://kong-aip.netlify.app/list/).

- **MeterForge** (`meterforge.tsp`): MeterForge v3 API.
- **Konnect** (`konnect.tsp`): Konnect metering & billing API, same surface with Konnect-specific auth and servers.
