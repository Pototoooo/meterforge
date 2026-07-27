from os import environ
from typing import Optional
import datetime
import uuid

from meterforge import Client
from meterforge.models import Event
from corehttp.exceptions import HttpResponseError

ENDPOINT: str = environ.get("METERFORGE_ENDPOINT") or "http://localhost:8888"
token: Optional[str] = environ.get("METERFORGE_TOKEN")


client = Client(
    endpoint=ENDPOINT,
    token=token,
)


def main() -> None:
    try:
        # Create a CloudEvents event
        event = Event(
            id=str(uuid.uuid4()),
            source="my-app",
            specversion="1.0",
            type="prompt",
            subject="customer-1",
            time=datetime.datetime.now(datetime.timezone.utc),
            data={
                "tokens": 100,
                "model": "gpt-4o",
                "type": "input",
            },
        )

        # Ingest the event
        client.events.ingest_event(event)
        print("Event ingested successfully")
    except HttpResponseError as e:
        print(f"Error ingesting event: {e}")


main()
