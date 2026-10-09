---
page_title: "oneuptime_huntress_connection Resource - oneuptime"
subcategory: "Other"
description: |-
  Huntress webhook endpoints. Every Huntress incident report opens an incident that pages on-call, and closing the report in Huntress resolves it.
---

# oneuptime_huntress_connection (Resource)

Huntress webhook endpoints. Every Huntress incident report opens an incident that pages on-call, and closing the report in Huntress resolves it.

## Example Usage

```terraform
resource "oneuptime_huntress_connection" "example" {
  name = "Example huntress connection"
}
```

## Schema

### Required

- `name` (String) A name for this connection, such as the Huntress account it receives from.

### Optional

- `critical_incident_severity_id` (String) ID of the incident severity critical reports open at. The ID of a `oneuptime_incident_severity`.
- `high_incident_severity_id` (String) ID of the incident severity high reports open at. The ID of a `oneuptime_incident_severity`.
- `is_signing_secret_set` (Boolean) Whether a signing secret is saved. Requests are only accepted once it is. Set from the signing secret itself; a value sent for it is ignored. Defaults to `false`.
- `labels` (Set of String) Labels added to every incident this connection opens, next to the label named after the report's organization. IDs of `oneuptime_label` resources.
- `low_incident_severity_id` (String) ID of the incident severity low reports open at. The ID of a `oneuptime_incident_severity`.
- `on_call_duty_policies` (Set of String) The on-call policies paged for an incident report at or above the Page On-Call For severity. IDs of `oneuptime_on_call_policy` resources.
- `page_on_call_for` (String) The lowest Huntress severity that pages the on-call policies: critical (critical reports only), high (high and critical reports) or low (every report). Every report opens an incident either way. Defaults to `high`.
- `resolve_incident_when_report_closes` (Boolean) Resolve the incident when the report is closed or dismissed in Huntress. When off, a note on the incident says so instead. Defaults to `true`.
- `signing_secret` (String) The endpoint's signing secret from Huntress (whsec_...), used to verify that a request really came from Huntress. Encrypted at rest and never returned by the API.
- `watched_organizations` (String) The Huntress organizations this connection opens incidents for, one per line, by name or by id. Empty watches every organization.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_error` (String) Why the last request was refused or could not be handled, if it was. Cleared by the next request that is handled.
- `last_error_at` (String) When the last error happened.
- `last_event_received_at` (String) When a request signed with the signing secret last arrived. Empty until Huntress sends the first one.
- `last_event_type` (String) The type of the last event received, such as incident_report.created.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing huntress connection by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_huntress_connection.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_huntress_connection.example <id>
```
