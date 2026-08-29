---
page_title: "oneuptime_threat_intel_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  STIX/TAXII 2.1 threat-intelligence feeds. Indicators are polled on an interval and matched against incoming security events.
---

# oneuptime_threat_intel_feed (Resource)

STIX/TAXII 2.1 threat-intelligence feeds. Indicators are polled on an interval and matched against incoming security events.

## Example Usage

```terraform
resource "oneuptime_threat_intel_feed" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  api_root_url = "This is an example of longer text content that might be stored in this field."
  collection_id = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `api_root_url` (String) The TAXII 2.1 API root, e.g. https://taxii.example.com/api1/. Collections are addressed beneath it...
- `collection_id` (String) ID of the TAXII collection to poll for indicator objects...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) What this feed carries and why it is subscribed...
- `api_token` (String) Bearer token for token-authenticated collections. Encrypted at rest and never returned by the API. Leave empty for anonymous or basic-auth collections...
- `basic_auth_username` (String) Username for basic-auth collections. Leave empty for anonymous or token-authenticated collections...
- `basic_auth_password` (String) Password for basic-auth collections. Encrypted at rest and never returned by the API...
- `is_enabled` (Bool) Whether this feed is polled and matched...
- `poll_interval_in_minutes` (Number) How often the collection is polled for new indicators...
- `minimum_confidence` (Number) Skip indicators whose STIX confidence is below this (0-100). 0 ingests everything; indicators that carry no confidence always pass...
- `should_create_alert` (Bool) Whether indicator matches open OneUptime alerts...
- `should_write_detection_finding` (Bool) Whether matches also write a Detection Finding security event back into the events table...
- `should_create_incident` (Bool) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per feed...
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID..
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_polled_at` (String) A date time object..
- `cursor` (String) Poll cursor: the TAXII added_after timestamp already ingested, as an ISO string...
- `next_page_token` (String) Resume token for a poll that ended mid-pagination on a server that sends no X-TAXII-Date-Added-Last header. Cleared once the collection drains or the cursor advances...
- `last_poll_summary` (String) What the most recent successful poll did: objects fetched, indicators ingested, unsupported patterns skipped...
- `last_error` (String) The most recent poll error, if any. Cleared on the next successful poll...
- `last_evaluated_at` (String) A date time object..
- `last_match_at` (String) A date time object..
- `last_match_error` (String) The most recent matcher error, if any. Cleared on the next successful evaluation...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_threat_intel_feed.example <id>
```
