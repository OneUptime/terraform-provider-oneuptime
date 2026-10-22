---
page_title: "oneuptime_threat_intel_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  Threat-intelligence feeds: TAXII 2.1 collections, lists of indicators and Google SecOps. Indicators are polled and matched against security events.
---

# oneuptime_threat_intel_feed (Resource)

Threat-intelligence feeds: TAXII 2.1 collections, lists of indicators and Google SecOps. Indicators are polled and matched against security events.

## Example Usage

```terraform
resource "oneuptime_threat_intel_feed" "example" {
  name        = "Example threat intel feed"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this feed, e.g. 'Corporate MISP'.

### Optional

- `alert_severity_id` (String) ID of the alert severity for alerts opened by this feed. The ID of a `oneuptime_alert_severity`.
- `api_root_url` (String) The TAXII 2.1 API root, e.g. https://taxii.example.com/api1/. Collections are addressed beneath it.
- `api_token` (String, Sensitive) Bearer token for token-authenticated collections. Encrypted at rest and never returned by the API. Leave empty for anonymous or basic-auth collections.
- `basic_auth_password` (String, Sensitive) Password for basic-auth collections. Encrypted at rest and never returned by the API.
- `basic_auth_username` (String) Username for basic-auth collections. Leave empty for anonymous or token-authenticated collections.
- `collection_id` (String) ID of the TAXII collection to poll for indicator objects.
- `description` (String) What this feed carries and why it is subscribed.
- `feed_type` (String) Where the indicators come from: a TAXII 2.1 collection (taxii), a list at an https address (list) or Google SecOps (google-secops). Defaults to `taxii`.
- `incident_severity_id` (String) ID of the incident severity for incidents opened by this feed. The ID of a `oneuptime_incident_severity`.
- `indicator_confidence` (Number) For a list feed: how sure the list is, 1 to 100, given to each of its indicators. It sets the severity of a match.
- `is_enabled` (Boolean) Whether this feed is polled and matched. Defaults to `true`.
- `list_url` (String) For a list feed: the https address of the list, one IP address, domain, URL or file hash a line.
- `minimum_confidence` (Number) Skip indicators whose STIX confidence is below this (0-100). 0 ingests everything; indicators that carry no confidence always pass. Defaults to `0`.
- `poll_interval_in_minutes` (Number) How often the collection is polled for new indicators. Defaults to `60`.
- `security_event_connection_id` (String) For a Google SecOps feed: the ID of the connection whose key reads Google SecOps's indicator matches. The ID of a `oneuptime_security_event_connection`.
- `should_create_alert` (Boolean) Whether indicator matches open OneUptime alerts. Defaults to `true`.
- `should_create_case` (Boolean) Whether matches open security cases (or join the open case about the same entity), where the security team works them and playbooks respond. Defaults to `true`.
- `should_create_incident` (Boolean) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per feed. Defaults to `false`.
- `should_write_detection_finding` (Boolean) Whether matches also write a Detection Finding security event back into the events table. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this feed. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) Poll cursor: for a TAXII feed the added_after timestamp already ingested, for a list feed the hash of the list last written and when, for a Google SecOps feed the end of the window last read.
- `id` (String) Unique identifier for the resource.
- `last_error` (String) The most recent poll error, if any. Cleared on the next successful poll.
- `last_evaluated_at` (String) When the matcher last evaluated security events against this feed's indicators. Null means it has never run.
- `last_match_at` (String) When this feed's indicators most recently matched security events. Null means they never have.
- `last_match_error` (String) The most recent matcher error, if any. Cleared on the next successful evaluation.
- `last_poll_summary` (String) What the most recent successful poll did: objects fetched, indicators ingested, unsupported patterns skipped.
- `last_polled_at` (String) When this feed was last polled. Null means it has never run.
- `next_page_token` (String) Resume token for a poll that ended mid-pagination on a server that sends no X-TAXII-Date-Added-Last header. Cleared once the collection drains or the cursor advances.
- `project_id` (String) ID of the project this feed belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing threat intel feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_threat_intel_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_threat_intel_feed.example <id>
```
