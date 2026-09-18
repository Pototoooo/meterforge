# MeterForge Go SDK

## Install

```sh
go get github.com/Pototoooo/meterforge/api/client/go@v1.0.0-beta.53
```

## Usage

Initialize client.

```go
import (
  cloudevents "github.com/cloudevents/sdk-go/v2/event"
  mf "github.com/Pototoooo/meterforge/api/client/go"
)

func main() {
  // Initialize MeterForge client
  mf, err := meterforge.NewClientWithResponses("http://localhost:8888")
  if err != nil {
      panic(err.Error())
  }

  // Use MeterForge client
  // ...
}
```

### Ingest Event

Report usage to MeterForge.

```go
e := cloudevents.New()
e.SetID("00001")
e.SetSource("my-app")
e.SetType("tokens")
e.SetSubject("user-id")
e.SetTime(time.Now())
e.SetData("application/json", map[string]string{
  "tokens": "15",
  "model": "gpt-4",
})

_, err := client.IngestEventWithResponse(ctx, e)
```

### Query Meter

Retreive usage from MeterForge.

```go
slug := "token-usage"
subject := []string{"user-id"}
from, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
to, _ := time.Parse(time.RFC3339, "2023-01-02T00:00:00Z")
resp, _ := client.QueryMeterWithResponse(ctx, slug, &mf.QueryMeterParams{
    Subject: &subject,
    From:    &from,
    To:      &to,
})
// resp.JSON200.Data
```
