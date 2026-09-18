# MeterForge Collector

MeterForge Collector is a configurable, production-ready data pipeline for **usage metering**. It helps you **collect, transform, buffer, and reliably deliver** usage events into MeterForge, especially in distributed and network-unreliable environments.

Learn more in the docs:

* Overview: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors)
* Quickstart: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/quickstart](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/quickstart)
* How it works: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/how-it-works](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/how-it-works)

---

## Capabilities

* **Multiple ingestion sources**: HTTP/event ingestion and a growing set of presets and integrations

  * Kubernetes: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/kubernetes](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/kubernetes)
  * Prometheus: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/prometheus](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/prometheus)
  * OpenTelemetry: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/otel](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/otel)
  * ClickHouse: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/clickhouse](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/clickhouse)
  * Postgres: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/postgres](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/postgres)
  * S3: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/s3](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/s3)
  * NVIDIA Run:ai: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/nvidia-run-ai](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/nvidia-run-ai)

* **Reliable delivery with buffering**: disk-backed buffering for network resilience and backpressure handling
  [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/buffer](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/buffer)

* **High availability patterns**: deploy and operate collectors in HA setups
  [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/high-availability](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/high-availability)

* **Built-in observability**: metrics and operational visibility for pipelines
  [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/observability](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/observability)

---

## Architecture

```text
+-------------------+   +-------------------+   +-------------------+
|  App / SDK Events  |   |  Infra Signals     |   |  Data Stores       |
|  (HTTP, webhooks)  |   |  (K8s/Prom/OTEL)   |   |  (PG/CH/S3/...)    |
+---------+---------+   +---------+---------+   +---------+---------+
          \                   |                     /
           \                  |                    /
            \                 |                   /
             v                v                  v
        +------------------------------------------------+
        |                MeterForge Collector              |
        |------------------------------------------------|
        |  Ingest -> Validate/Transform -> Batch/Retry    |
        |                 |                               |
        |                 v                               |
        |        Disk Buffer (durable replay)             |
        +-----------------+------------------------------+
                          |
                          v
              +------------------------------+
              |          MeterForge           |
              |   metering • usage • billing |
              +------------------------------+

```

See detailed architecture and concepts: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/how-it-works](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/how-it-works)

---

## Getting Started

### 1. Choose a deployment model

* Kubernetes (Helm): [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/kubernetes](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/kubernetes)
* Other environments and presets: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/quickstart](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/quickstart)

### 2. Install the Collector

Follow the official quickstart for step-by-step installation and configuration examples:
[https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/quickstart](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/quickstart)

### 3. Send usage events

Configure your SDKs or usage producers to send events to the Collector endpoint.

### 4. Run in production

* Buffering & delivery guarantees: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/buffer](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/buffer)
* High availability: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/high-availability](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/high-availability)
* Observability & metrics: [https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/observability](https://github.com/Pototoooo/meterforge/tree/main/docs/collectors/observability)

---

## Repository

Collector source code lives in the main MeterForge repository:
[https://github.com/Pototoooo/meterforge/tree/main/collector](https://github.com/Pototoooo/meterforge/tree/main/collector)
