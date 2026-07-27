# Quickstart

## Prerequisites

- Docker (with Compose)
- curl
- jq

Clone the repository:

```sh
git clone git@github.com:Pototoooo/meterforge.git
cd meterforge/quickstart
```

## 1. Launch MeterForge

Launch MeterForge and its dependencies via:

```sh
docker compose up -d
```

Open the local Metering & Billing console at:

<http://localhost:3000>

The console follows MeterForge's core workflow:

1. Create or inspect a meter.
2. Create a feature and publish a plan.
3. Create a customer and start a subscription.
4. Send usage events and query the aggregated usage.
5. Inspect the resulting entitlement balance and invoice.

## 2. Ingest usage event(s)

Before running the examples, open **Billing → New Customer** in the console and
set its Usage Attribution subject to `customer-1`. MeterForge keeps unmatched
subjects in the event stream as invalid events and does not aggregate them.

Ingest usage events in [CloudEvents](https://cloudevents.io/) format:

```sh
curl -X POST http://localhost:48888/api/v1/events \
-H 'Content-Type: application/cloudevents+json' \
--data-raw '
{
  "specversion" : "1.0",
  "type": "request",
  "id": "00001",
  "time": "2026-07-07T00:00:00.001Z",
  "source": "service-0",
  "subject": "customer-1",
  "data": {
    "method": "GET",
    "route": "/hello",
    "duration_ms": 10
  }
}
'
```

Note how ID is different:

```sh
curl -X POST http://localhost:48888/api/v1/events \
-H 'Content-Type: application/cloudevents+json' \
--data-raw '
{
  "specversion" : "1.0",
  "type": "request",
  "id": "00002",
  "time": "2026-07-07T00:00:00.001Z",
  "source": "service-0",
  "subject": "customer-1",
  "data": {
    "method": "GET",
    "route": "/hello",
    "duration_ms": 20
  }
}
'
```

Note how ID and time are different:

```sh
curl -X POST http://localhost:48888/api/v1/events \
-H 'Content-Type: application/cloudevents+json' \
--data-raw '
{
  "specversion" : "1.0",
  "type": "request",
  "id": "00003",
  "time": "2026-07-08T00:00:00.001Z",
  "source": "service-0",
  "subject": "customer-1",
  "data": {
    "method": "GET",
    "route": "/hello",
    "duration_ms": 30
  }
}
'
```

## 3. Query Usage

Query the usage hourly:

```sh
curl 'http://localhost:48888/api/v1/meters/api_requests_total/query?windowSize=HOUR&groupBy=method&groupBy=route' | jq
```

```json
{
  "windowSize": "HOUR",
  "data": [
    {
      "value": 2,
      "windowStart": "2026-07-07T00:00:00Z",
      "windowEnd": "2026-07-07T01:00:00Z",
      "subject": null,
      "groupBy": {
        "method": "GET",
        "route": "/hello"
      }
    },
    {
      "value": 1,
      "windowStart": "2026-07-08T00:00:00Z",
      "windowEnd": "2026-07-08T01:00:00Z",
      "subject": null,
      "groupBy": {
        "method": "GET",
        "route": "/hello"
      }
    }
  ]
}
```

Query the total usage for `customer-1`:

```sh
curl 'http://localhost:48888/api/v1/meters/api_requests_total/query?subject=customer-1' | jq
```

```json
{
  "data": [
    {
      "value": 3,
      "windowStart": "2026-07-07T00:00:00Z",
      "windowEnd": "2026-07-08T00:01:00Z",
      "subject": "customer-1",
      "groupBy": {}
    }
  ]
}
```

## 4. Configure additional meter(s) _(optional)_

In this example we will meter LLM token usage, groupped by AI model and prompt type.
You can think about it how OpenAI [charges](https://openai.com/pricing) by tokens for ChatGPT.

Configure how MeterForge should process your usage events in this new `tokens_total` meter.

```yaml
# ...

meters:
  # Sample meter to count LLM Token Usage
  - slug: tokens_total
    description: AI Token Usage
    eventType: prompt               # Filter events by type
    aggregation: SUM
    valueProperty: $.tokens         # JSONPath to parse usage value
    groupBy:
      model: $.model                # AI model used: gpt4-turbo, etc.
      type: $.type                  # Prompt type: input, output, system

```

## Cleanup

Once you are done, stop any running instances:

```sh
docker compose down -v
```
