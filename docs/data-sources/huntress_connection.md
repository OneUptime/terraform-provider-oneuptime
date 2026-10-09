---
page_title: "oneuptime_huntress_connection Data Source - oneuptime"
subcategory: "Other"
description: |-
  Huntress webhook endpoints. Every Huntress incident report opens an incident that pages on-call, and closing the report in Huntress resolves it.
---

# oneuptime_huntress_connection (Data Source)

Huntress webhook endpoints. Every Huntress incident report opens an incident that pages on-call, and closing the report in Huntress resolves it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one huntress connection may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_huntress_connection" "example" {
  name = "Example huntress connection"
}

# Or by id:
data "oneuptime_huntress_connection" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `critical_incident_severity_id` (String) ID of the incident severity critical reports open at. The ID of a `oneuptime_incident_severity`.
- `high_incident_severity_id` (String) ID of the incident severity high reports open at. The ID of a `oneuptime_incident_severity`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_signing_secret_set` (Boolean) Whether a signing secret is saved. Requests are only accepted once it is. Set from the signing secret itself; a value sent for it is ignored.
- `last_error` (String) Why the last request was refused or could not be handled, if it was. Cleared by the next request that is handled.
- `last_event_type` (String) The type of the last event received, such as incident_report.created.
- `low_incident_severity_id` (String) ID of the incident severity low reports open at. The ID of a `oneuptime_incident_severity`.
- `name` (String) A name for this connection, such as the Huntress account it receives from.
- `page_on_call_for` (String) The lowest Huntress severity that pages the on-call policies: critical (critical reports only), high (high and critical reports) or low (every report). Every report opens an incident either way.
- `resolve_incident_when_report_closes` (Boolean) Resolve the incident when the report is closed or dismissed in Huntress. When off, a note on the incident says so instead.
- `watched_organizations` (String) The Huntress organizations this connection opens incidents for, one per line, by name or by id. Empty watches every organization.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Labels added to every incident this connection opens, next to the label named after the report's organization. IDs of `oneuptime_label` resources.
- `last_error_at` (String) When the last error happened.
- `last_event_received_at` (String) When a request signed with the signing secret last arrived. Empty until Huntress sends the first one.
- `on_call_duty_policies` (Set of String) The on-call policies paged for an incident report at or above the Page On-Call For severity. IDs of `oneuptime_on_call_policy` resources.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
