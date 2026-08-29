---
page_title: "oneuptime_threat_intel_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  STIX/TAXII 2.1 threat-intelligence feeds. Indicators are polled on an interval and matched against incoming security events.
---

# oneuptime_threat_intel_feed (Data Source)

STIX/TAXII 2.1 threat-intelligence feeds. Indicators are polled on an interval and matched against incoming security events. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_threat_intel_feed" "by_name" {
  name = "example-threat_intel_feed"
}

data "oneuptime_threat_intel_feed" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) What this feed carries and why it is subscribed... Computed.
- `api_root_url` (String) The TAXII 2.1 API root, e.g. https://taxii.example.com/api1/. Collections are addressed beneath it... Computed.
- `collection_id` (String) ID of the TAXII collection to poll for indicator objects... Computed.
- `basic_auth_username` (String) Username for basic-auth collections. Leave empty for anonymous or token-authenticated collections... Computed.
- `is_enabled` (Bool) Whether this feed is polled and matched... Computed.
- `poll_interval_in_minutes` (Number) How often the collection is polled for new indicators... Computed.
- `minimum_confidence` (Number) Skip indicators whose STIX confidence is below this (0-100). 0 ingests everything; indicators that carry no confidence always pass... Computed.
- `should_create_alert` (Bool) Whether indicator matches open OneUptime alerts... Computed.
- `should_write_detection_finding` (Bool) Whether matches also write a Detection Finding security event back into the events table... Computed.
- `should_create_incident` (Bool) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per feed... Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `last_polled_at` (String) A date time object.. Computed.
- `cursor` (String) Poll cursor: the TAXII added_after timestamp already ingested, as an ISO string... Computed.
- `next_page_token` (String) Resume token for a poll that ended mid-pagination on a server that sends no X-TAXII-Date-Added-Last header. Cleared once the collection drains or the cursor advances... Computed.
- `last_poll_summary` (String) What the most recent successful poll did: objects fetched, indicators ingested, unsupported patterns skipped... Computed.
- `last_error` (String) The most recent poll error, if any. Cleared on the next successful poll... Computed.
- `last_evaluated_at` (String) A date time object.. Computed.
- `last_match_at` (String) A date time object.. Computed.
- `last_match_error` (String) The most recent matcher error, if any. Cleared on the next successful evaluation... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
