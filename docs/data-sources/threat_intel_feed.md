---
page_title: "oneuptime_threat_intel_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  STIX/TAXII 2.1 threat-intelligence feeds. Indicators are polled on an interval and matched against incoming security events.
---

# oneuptime_threat_intel_feed (Data Source)

STIX/TAXII 2.1 threat-intelligence feeds. Indicators are polled on an interval and matched against incoming security events.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one threat intel feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_threat_intel_feed" "example" {
  name = "Example threat intel feed"
}

# Or by id:
data "oneuptime_threat_intel_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_severity_id` (String) ID of the alert severity for alerts opened by this feed. The ID of a `oneuptime_alert_severity`.
- `api_root_url` (String) The TAXII 2.1 API root, e.g. https://taxii.example.com/api1/. Collections are addressed beneath it.
- `basic_auth_username` (String) Username for basic-auth collections. Leave empty for anonymous or token-authenticated collections.
- `collection_id` (String) ID of the TAXII collection to poll for indicator objects.
- `created_by_user_id` (String) ID of the user who created this feed. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) Poll cursor: the TAXII added_after timestamp already ingested, as an ISO string.
- `description` (String) What this feed carries and why it is subscribed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_severity_id` (String) ID of the incident severity for incidents opened by this feed. The ID of a `oneuptime_incident_severity`.
- `is_enabled` (Boolean) Whether this feed is polled and matched.
- `last_error` (String) The most recent poll error, if any. Cleared on the next successful poll.
- `last_match_error` (String) The most recent matcher error, if any. Cleared on the next successful evaluation.
- `last_poll_summary` (String) What the most recent successful poll did: objects fetched, indicators ingested, unsupported patterns skipped.
- `minimum_confidence` (Number) Skip indicators whose STIX confidence is below this (0-100). 0 ingests everything; indicators that carry no confidence always pass.
- `name` (String) Friendly name for this feed, e.g. 'Corporate MISP'.
- `next_page_token` (String) Resume token for a poll that ended mid-pagination on a server that sends no X-TAXII-Date-Added-Last header. Cleared once the collection drains or the cursor advances.
- `poll_interval_in_minutes` (Number) How often the collection is polled for new indicators.
- `should_create_alert` (Boolean) Whether indicator matches open OneUptime alerts.
- `should_create_incident` (Boolean) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per feed.
- `should_write_detection_finding` (Boolean) Whether matches also write a Detection Finding security event back into the events table.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_evaluated_at` (String) When the matcher last evaluated security events against this feed's indicators. Null means it has never run.
- `last_match_at` (String) When this feed's indicators most recently matched security events. Null means they never have.
- `last_polled_at` (String) When this feed was last polled. Null means it has never run.
- `project_id` (String) ID of the project this feed belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
